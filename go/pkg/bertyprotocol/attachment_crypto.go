package bertyprotocol

import (
	"bytes"
	"crypto/cipher"
	"fmt"
	"io"
	"math/big"

	libp2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
	"go.uber.org/zap"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/sha3"

	"berty.tech/berty/v2/go/internal/cryptoutil"
	"berty.tech/berty/v2/go/pkg/bertytypes"
	"berty.tech/berty/v2/go/pkg/errcode"
)

// ⚠⚠⚠ FIXME: Needs thorough security review ⚠⚠⚠

const (
	attachmentCipherblockSize = 64 * 1024
	attachmentNonceIV         = 0
	attachmentKeyV0Prefix     = "/libp2psk+xchacha20poly1305_64_0/" // TODO: replace when multikey rolls out
)

// - DATA ENCRYPTION

type attachmentCipher struct {
	nonce    *big.Int
	nonceBuf [chacha20poly1305.NonceSizeX]byte
	aead     cipher.AEAD
}

func attachmentNewCipher(sk libp2pcrypto.PrivKey) (*attachmentCipher, error) {
	key, err := sk.Raw()
	if err != nil {
		return nil, errcode.ErrInvalidInput.Wrap(err)
	}
	if len(key) < chacha20poly1305.KeySize {
		return nil, errcode.ErrInvalidInput.Wrap(fmt.Errorf("bad key size, got %d, expected to be >= %d", len(key), chacha20poly1305.KeySize))
	}

	aead, err := chacha20poly1305.NewX(key[:chacha20poly1305.KeySize])
	if err != nil {
		return nil, errcode.ErrCryptoCipherInit.Wrap(err)
	}

	ac := attachmentCipher{
		aead:  aead,
		nonce: big.NewInt(attachmentNonceIV),
	}

	bigIntFillBytes(ac.nonce, ac.nonceBuf[:])

	return &ac, nil
}

var bigOne = big.NewInt(1)

func attachmentSealer(plaintext io.Reader, l *zap.Logger) (libp2pcrypto.PrivKey, *io.PipeReader, error) {
	sk, _, err := libp2pcrypto.GenerateKeyPair(libp2pcrypto.Ed25519, 0)
	if err != nil {
		return nil, nil, errcode.ErrCryptoKeyGeneration.Wrap(err)
	}

	ac, err := attachmentNewCipher(sk)
	if err != nil {
		return nil, nil, errcode.ErrCryptoCipherInit.Wrap(err)
	}

	in, out := io.Pipe()

	buf := make([]byte, attachmentCipherblockSize-16)
	go func() {
		err := func() error {
			for {
				n, readErr := io.ReadFull(plaintext, buf)
				if readErr == io.EOF {
					return nil
				}
				if readErr != nil && readErr != io.ErrUnexpectedEOF {
					return errcode.ErrStreamRead.Wrap(readErr)
				}

				ciphertext := ac.aead.Seal([]byte(nil), ac.nonceBuf[:], buf[:n], []byte(nil))

				if _, err := out.Write(ciphertext); err != nil {
					return errcode.ErrStreamWrite.Wrap(err)
				}

				if readErr == io.ErrUnexpectedEOF {
					return nil // last block can be smaller
				}

				ac.nonce.Add(ac.nonce, bigOne)
				bigIntFillBytes(ac.nonce, ac.nonceBuf[:])
			}
		}()
		closePipeOut(out, err, "attachmentSealer: failed to properly close ciphertext out", l)
	}()

	return sk, in, nil
}

func attachmentOpener(ciphertext io.Reader, sk libp2pcrypto.PrivKey, l *zap.Logger) (*io.PipeReader, error) {
	ac, err := attachmentNewCipher(sk)
	if err != nil {
		return nil, errcode.ErrCryptoCipherInit.Wrap(err)
	}

	in, out := io.Pipe()

	buf := make([]byte, attachmentCipherblockSize)
	go func() {
		err := func() error {
			for {
				n, readErr := io.ReadFull(ciphertext, buf)
				if readErr == io.EOF {
					return nil
				}
				if readErr != nil && readErr != io.ErrUnexpectedEOF {
					return errcode.ErrStreamRead.Wrap(readErr)
				}

				plaintext, err := ac.aead.Open([]byte(nil), ac.nonceBuf[:], buf[:n], []byte(nil))
				if err != nil {
					return errcode.ErrCryptoDecrypt.Wrap(err)
				}

				if _, err := out.Write(plaintext); err != nil {
					return errcode.ErrStreamWrite.Wrap(err)
				}

				if readErr == io.ErrUnexpectedEOF {
					return nil // last block can be smaller
				}

				ac.nonce.Add(ac.nonce, bigOne)
				bigIntFillBytes(ac.nonce, ac.nonceBuf[:])
			}
		}()
		closePipeOut(out, err, "attachmentOpener: failed to properly close plaintext out", l)
	}()

	return in, nil
}

// - KEY SERIALIZATION

func attachmentKeyMarshal(sk libp2pcrypto.PrivKey) ([]byte, error) {
	skBytes, err := libp2pcrypto.MarshalPrivateKey(sk)
	if err != nil {
		return nil, errcode.ErrSerialization.Wrap(err)
	}

	return append([]byte(attachmentKeyV0Prefix), skBytes...), nil
}

func attachmentKeyUnmarshal(s []byte) (libp2pcrypto.PrivKey, error) {
	if len(s) <= len(attachmentKeyV0Prefix) || !bytes.Equal(s[:len(attachmentKeyV0Prefix)], []byte(attachmentKeyV0Prefix)) {
		return nil, errcode.ErrInvalidInput.Wrap(fmt.Errorf("invalid secret prefix %s", s))
	}
	skBytes := s[len(attachmentKeyV0Prefix):]

	sk, err := libp2pcrypto.UnmarshalPrivateKey(skBytes)
	if err != nil {
		return nil, errcode.ErrDeserialization.Wrap(err)
	}

	return sk, nil
}

// - CID ENCRYPTION

func attachmentCIDEncryptionKey(source *[cryptoutil.KeySize]byte) (*[cryptoutil.KeySize]byte, error) {
	hkdf := hkdf.New(sha3.New256, source[:], nil, []byte("cid encryption v0"))

	var key [cryptoutil.KeySize]byte
	if _, err := io.ReadFull(hkdf, key[:]); err != nil {
		return nil, errcode.ErrStreamRead.Wrap(err)
	}

	return &key, nil
}

func attachmentCIDEncrypt(sk *[cryptoutil.KeySize]byte, cid []byte) ([]byte, error) {
	nonce, err := cryptoutil.GenerateNonce()
	if err != nil {
		return nil, errcode.ErrCryptoNonceGeneration.Wrap(err)
	}

	return append(nonce[:], secretbox.Seal(nil, cid, nonce, sk)...), nil
}

func attachmentCIDDecrypt(sk *[cryptoutil.KeySize]byte, eCID []byte) ([]byte, error) {
	if len(eCID) <= cryptoutil.NonceSize {
		return nil, errcode.ErrInvalidInput.Wrap(fmt.Errorf("encrypted cid too small, got %v, expected to be > %v", len(eCID), cryptoutil.NonceSize))
	}

	var nonce [cryptoutil.NonceSize]byte
	_ = copy(nonce[:], eCID[:cryptoutil.NonceSize])

	cid, ok := secretbox.Open(nil, eCID[cryptoutil.NonceSize:], &nonce, sk)
	if !ok {
		return nil, errcode.ErrCryptoDecrypt
	}

	return cid, nil
}

func attachmentCIDSliceEncrypt(g *bertytypes.Group, cids [][]byte) ([][]byte, error) {
	sk, err := attachmentCIDEncryptionKey(g.GetSharedSecret())
	if err != nil {
		return nil, errcode.ErrCryptoKeyDerivation.Wrap(err)
	}
	return mapBufArray(cids, func(cid []byte) ([]byte, error) { return attachmentCIDEncrypt(sk, cid) })
}

func attachmentCIDSliceDecrypt(g *bertytypes.Group, eCIDs [][]byte) ([][]byte, error) {
	sk, err := attachmentCIDEncryptionKey(g.GetSharedSecret())
	if err != nil {
		return nil, errcode.ErrCryptoKeyDerivation.Wrap(err)
	}
	return mapBufArray(eCIDs, func(eCID []byte) ([]byte, error) { return attachmentCIDDecrypt(sk, eCID) })
}
