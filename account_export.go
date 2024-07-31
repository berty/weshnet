package weshnet

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/ipfs/go-cid"
	cbornode "github.com/ipfs/go-ipld-cbor"
	coreiface "github.com/ipfs/kubo/core/coreiface"
	mh "github.com/multiformats/go-multihash"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	orbitdb "berty.tech/go-orbit-db"
	"berty.tech/weshnet/pkg/errcode"
	"berty.tech/weshnet/pkg/protocoltypes"
)

const (
	exportAccountKeyFilename      = "account.key"
	exportAccountProofKeyFilename = "account_proof.key"
	exportOrbitDBEntriesPrefix    = "entries/"
	exportOrbitDBHeadsPrefix      = "heads/"
)

func (s *service) export(ctx context.Context, output io.Writer) error {
	tw := tar.NewWriter(output)
	defer tw.Close()

	if err := s.exportAccountKeys(tw); err != nil {
		return errcode.ErrCode_ErrInternal.Wrap(err)
	}

	s.lock.RLock()
	groups := make([]*GroupContext, len(s.openedGroups))
	i := 0
	for _, gc := range s.openedGroups {
		groups[i] = gc
		i++
	}
	s.lock.RUnlock()

	for _, gc := range groups {
		if err := s.exportGroupContext(ctx, gc, tw); err != nil {
			return errcode.ErrCode_ErrInternal.Wrap(err)
		}
	}

	return nil
}

func (s *service) exportGroupContext(ctx context.Context, gc *GroupContext, tw *tar.Writer) error {
	if err := s.exportOrbitDBStore(ctx, gc.metadataStore, tw); err != nil {
		return errcode.ErrCode_ErrInternal.Wrap(err)
	}

	if err := s.exportOrbitDBStore(ctx, gc.messageStore, tw); err != nil {
		return errcode.ErrCode_ErrInternal.Wrap(err)
	}

	metaRawHeads := gc.metadataStore.OpLog().RawHeads()
	cidsMeta := make([]cid.Cid, metaRawHeads.Len())
	for i, raw := range metaRawHeads.Slice() {
		cidsMeta[i] = raw.GetHash()
	}

	messagesRawHeads := gc.messageStore.OpLog().RawHeads()
	cidsMessages := make([]cid.Cid, messagesRawHeads.Len())
	for i, raw := range messagesRawHeads.Slice() {
		cidsMessages[i] = raw.GetHash()
	}

	if err := s.exportOrbitDBGroupHeads(gc, cidsMeta, cidsMessages, tw); err != nil {
		return errcode.ErrCode_ErrInternal.Wrap(err)
	}

	return nil
}

func (s *service) exportOrbitDBStore(ctx context.Context, store orbitdb.Store, tw *tar.Writer) error {
	allCIDs := store.OpLog().GetEntries().Keys()

	if len(allCIDs) == 0 {
		return nil
	}

	for _, idStr := range allCIDs {
		if err := s.exportOrbitDBEntry(ctx, tw, idStr); err != nil {
			if clErr := tw.Close(); clErr != nil {
				err = multierr.Append(err, clErr)
			}

			return errcode.ErrCode_ErrInternal.Wrap(err)
		}
	}

	return nil
}

func (s *service) exportAccountKeys(tw *tar.Writer) error {
	accountPrivateKeyBytes, accountProofPrivateKeyBytes, err := s.secretStore.ExportAccountKeysForBackup()
	if err != nil {
		return errcode.ErrCode_ErrInternal.Wrap(err)
	}

	err = exportPrivateKey(tw, accountPrivateKeyBytes, exportAccountKeyFilename)
	if err != nil {
		return errcode.ErrCode_ErrStreamWrite.Wrap(err)
	}

	err = exportPrivateKey(tw, accountProofPrivateKeyBytes, exportAccountProofKeyFilename)
	if err != nil {
		return errcode.ErrCode_ErrStreamWrite.Wrap(err)
	}

	return nil
}

func (s *service) exportOrbitDBGroupHeads(gc *GroupContext, headsMetadata []cid.Cid, headsMessages []cid.Cid, tw *tar.Writer) error {
	cidsMeta := make([][]byte, len(headsMetadata))
	for i, id := range headsMetadata {
		cidsMeta[i] = id.Bytes()
	}

	cidsMessages := make([][]byte, len(headsMessages))
	for i, id := range headsMessages {
		cidsMessages[i] = id.Bytes()
	}

	spk, err := gc.group.GetSigningPubKey()
	if err != nil {
		return errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	spkBytes, err := spk.Raw()
	if err != nil {
		return errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	linkKeyArr, err := gc.group.GetLinkKeyArray()
	if err != nil {
		return errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	headsExport := &protocoltypes.GroupHeadsExport{
		PublicKey:         gc.group.PublicKey,
		SignPub:           spkBytes,
		MetadataHeadsCids: cidsMeta,
		MessagesHeadsCids: cidsMessages,
		LinkKey:           linkKeyArr[:],
	}

	entryName := base64.RawURLEncoding.EncodeToString(gc.group.PublicKey)

	data, err := proto.Marshal(headsExport)
	if err != nil {
		return errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	if err := tw.WriteHeader(&tar.Header{
		Typeflag: tar.TypeReg,
		Name:     fmt.Sprintf("%s%s", exportOrbitDBHeadsPrefix, entryName),
		Mode:     0o600,
		Size:     int64(len(data)),
	}); err != nil {
		return errcode.ErrCode_ErrStreamWrite.Wrap(err)
	}

	size, err := tw.Write(data)
	if err != nil {
		return errcode.ErrCode_ErrStreamWrite.Wrap(err)
	}

	if size != len(data) {
		return errcode.ErrCode_ErrStreamWrite.Wrap(fmt.Errorf("wrote %d bytes instead of %d", size, len(data)))
	}

	return nil
}

func exportPrivateKey(tw *tar.Writer, marshalledPrivateKey []byte, filename string) error {
	if err := tw.WriteHeader(&tar.Header{
		Typeflag: tar.TypeReg,
		Name:     filename,
		Mode:     0o600,
		Size:     int64(len(marshalledPrivateKey)),
	}); err != nil {
		return errcode.ErrCode_ErrStreamWrite.Wrap(err)
	}

	size, err := tw.Write(marshalledPrivateKey)
	if err != nil {
		return errcode.ErrCode_ErrStreamWrite.Wrap(err)
	}

	if size != len(marshalledPrivateKey) {
		return errcode.ErrCode_ErrStreamWrite.Wrap(fmt.Errorf("wrote %d bytes instead of %d", size, len(marshalledPrivateKey)))
	}

	return nil
}

func (s *service) exportOrbitDBEntry(ctx context.Context, tw *tar.Writer, idStr string) error {
	id, err := cid.Parse(idStr)
	if err != nil {
		return errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	dagNode, err := s.ipfsCoreAPI.Dag().Get(ctx, id)
	if err != nil {
		return errcode.ErrCode_ErrInternal.Wrap(err)
	}

	dagNodeBytes := dagNode.RawData()

	if err := tw.WriteHeader(&tar.Header{
		Typeflag: tar.TypeReg,
		Name:     fmt.Sprintf("%s%s", exportOrbitDBEntriesPrefix, idStr),
		Mode:     0o600,
		Size:     int64(len(dagNodeBytes)),
	}); err != nil {
		return errcode.ErrCode_ErrStreamWrite.Wrap(err)
	}

	size, err := tw.Write(dagNodeBytes)
	if err != nil {
		return errcode.ErrCode_ErrStreamWrite.Wrap(err)
	}

	if size != len(dagNodeBytes) {
		return errcode.ErrCode_ErrStreamWrite.Wrap(fmt.Errorf("wrote %d bytes instead of %d", size, len(dagNodeBytes)))
	}

	return nil
}

func readExportSecretKeyFile(expectedSize int64, reader *tar.Reader) ([]byte, error) {
	if expectedSize == 0 {
		return nil, errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("invalid expected key size"))
	}

	keyContents := new(bytes.Buffer)
	size, err := io.Copy(keyContents, reader)
	if err != nil {
		return nil, errcode.ErrCode_ErrInternal.Wrap(fmt.Errorf("unable to read %d bytes: %w", expectedSize, err))
	}

	if size != expectedSize {
		return nil, errcode.ErrCode_ErrInternal.Wrap(fmt.Errorf("unexpected file size"))
	}

	return keyContents.Bytes(), nil
}

func readExportOrbitDBGroupHeads(expectedSize int64, reader *tar.Reader) (*protocoltypes.GroupHeadsExport, []cid.Cid, []cid.Cid, error) {
	if expectedSize == 0 {
		return nil, nil, nil, errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("invalid expected node size"))
	}

	nodeContents := new(bytes.Buffer)
	size, err := io.Copy(nodeContents, reader)
	if err != nil {
		return nil, nil, nil, errcode.ErrCode_ErrInternal.Wrap(fmt.Errorf("unable to read %d bytes: %w", expectedSize, err))
	}

	if size != expectedSize {
		return nil, nil, nil, errcode.ErrCode_ErrInternal.Wrap(fmt.Errorf("unexpected file size"))
	}

	groupHeads := &protocoltypes.GroupHeadsExport{}
	if err := proto.Unmarshal(nodeContents.Bytes(), groupHeads); err != nil {
		return nil, nil, nil, errcode.ErrCode_ErrDeserialization.Wrap(err)
	}

	messagesCIDs := make([]cid.Cid, len(groupHeads.MessagesHeadsCids))
	for i, cidBytes := range groupHeads.MessagesHeadsCids {
		messagesCIDs[i], err = cid.Parse(cidBytes)
		if err != nil {
			return nil, nil, nil, errcode.ErrCode_ErrDeserialization.Wrap(err)
		}
	}

	metaCIDs := make([]cid.Cid, len(groupHeads.MetadataHeadsCids))
	for i, cidBytes := range groupHeads.MetadataHeadsCids {
		metaCIDs[i], err = cid.Parse(cidBytes)
		if err != nil {
			return nil, nil, nil, errcode.ErrCode_ErrDeserialization.Wrap(err)
		}
	}

	return groupHeads, metaCIDs, messagesCIDs, nil
}

func readExportCBORNode(expectedSize int64, cidStr string, reader *tar.Reader) (*cbornode.Node, error) {
	if expectedSize == 0 {
		return nil, errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("invalid expected node size"))
	}

	nodeContents := new(bytes.Buffer)
	expectedCID, err := cid.Parse(cidStr)
	if err != nil {
		return nil, errcode.ErrCode_ErrDeserialization.Wrap(fmt.Errorf("unable to parse CID in filename"))
	}

	size, err := io.Copy(nodeContents, reader)
	if err != nil {
		return nil, errcode.ErrCode_ErrInternal.Wrap(fmt.Errorf("unable to read %d bytes: %w", expectedSize, err))
	}

	if size != expectedSize {
		return nil, errcode.ErrCode_ErrInternal.Wrap(fmt.Errorf("unexpected file size"))
	}

	node, err := cbornode.Decode(nodeContents.Bytes(), mh.SHA2_256, -1)
	if err != nil {
		return nil, errcode.ErrCode_ErrDeserialization.Wrap(err)
	}

	if !node.Cid().Equals(expectedCID) {
		return nil, errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("entry CID doesn't match file CID"))
	}

	return node, nil
}

type RestoreAccountHandler struct {
	Handler     func(header *tar.Header, reader *tar.Reader) (bool, error)
	PostProcess func() error
}

type restoreAccountState struct {
	keys map[string][]byte
}

func (state *restoreAccountState) readKey(keyName string) RestoreAccountHandler {
	return RestoreAccountHandler{
		Handler: func(header *tar.Header, reader *tar.Reader) (bool, error) {
			if header.Name != keyName {
				return false, nil
			}

			if state.keys[keyName] != nil {
				return false, errcode.ErrCode_ErrInternal.Wrap(fmt.Errorf("multiple keys found in archive"))
			}

			var err error

			state.keys[keyName], err = readExportSecretKeyFile(header.Size, reader)
			if err != nil {
				return true, errcode.ErrCode_ErrInternal.Wrap(err)
			}

			return true, nil
		},
	}
}

func (state *restoreAccountState) restoreKeys(odb *WeshOrbitDB) RestoreAccountHandler {
	return RestoreAccountHandler{
		PostProcess: func() error {
			if err := odb.secretStore.ImportAccountKeys(state.keys[exportAccountKeyFilename], state.keys[exportAccountProofKeyFilename]); err != nil {
				return errcode.ErrCode_ErrInternal.Wrap(err)
			}

			return nil
		},
	}
}

func restoreOrbitDBEntry(ctx context.Context, coreAPI coreiface.CoreAPI) RestoreAccountHandler {
	return RestoreAccountHandler{
		Handler: func(header *tar.Header, reader *tar.Reader) (bool, error) {
			if !strings.HasPrefix(header.Name, exportOrbitDBEntriesPrefix) {
				return false, nil
			}

			cidStr := strings.TrimPrefix(header.Name, exportOrbitDBEntriesPrefix)

			node, err := readExportCBORNode(header.Size, cidStr, reader)
			if err != nil {
				return true, errcode.ErrCode_ErrInternal.Wrap(err)
			}

			if err := coreAPI.Dag().Add(ctx, node); err != nil {
				return true, errcode.ErrCode_ErrInternal.Wrap(err)
			}

			return true, nil
		},
	}
}

func restoreOrbitDBHeads(ctx context.Context, odb *WeshOrbitDB) RestoreAccountHandler {
	return RestoreAccountHandler{
		Handler: func(header *tar.Header, reader *tar.Reader) (bool, error) {
			if !strings.HasPrefix(header.Name, exportOrbitDBHeadsPrefix) {
				return false, nil
			}

			heads, metaCIDs, messageCIDs, err := readExportOrbitDBGroupHeads(header.Size, reader)
			if err != nil {
				return true, errcode.ErrCode_ErrInternal.Wrap(err)
			}

			if err := odb.setHeadsForGroup(ctx, &protocoltypes.Group{
				PublicKey: heads.PublicKey,
				SignPub:   heads.SignPub,
				LinkKey:   heads.LinkKey,
			}, metaCIDs, messageCIDs); err != nil {
				return true, errcode.ErrCode_ErrOrbitDBAppend.Wrap(fmt.Errorf("error while restoring db head: %w", err))
			}

			return true, nil
		},
	}
}

func RestoreAccountExport(ctx context.Context, reader io.Reader, coreAPI coreiface.CoreAPI, odb *WeshOrbitDB, logger *zap.Logger, handlers ...RestoreAccountHandler) error {
	tr := tar.NewReader(reader)
	state := restoreAccountState{
		keys: map[string][]byte{},
	}

	handlers = append(
		[]RestoreAccountHandler{
			state.readKey(exportAccountKeyFilename),
			state.readKey(exportAccountProofKeyFilename),
			state.restoreKeys(odb),
			restoreOrbitDBEntry(ctx, coreAPI),
			restoreOrbitDBHeads(ctx, odb),
		},
		handlers...,
	)

	for {
		header, err := tr.Next()

		if err == io.EOF {
			break
		} else if err != nil {
			return errcode.ErrCode_ErrInternal.Wrap(err)
		}

		if header.Typeflag != tar.TypeReg {
			logger.Warn("invalid entry type", zap.String("filename", header.Name), zap.Any("filename", header.Typeflag))
			continue
		}

		notHandled := true

		for _, h := range handlers {
			if h.Handler == nil {
				continue
			}

			handled, err := h.Handler(header, tr)
			if err != nil {
				return errcode.ErrCode_ErrInternal.Wrap(err)
			}

			if handled {
				notHandled = false
				break
			}
		}

		if notHandled {
			logger.Warn("unknown export entry", zap.String("filename", header.Name))
		}
	}

	for _, h := range handlers {
		if h.PostProcess == nil {
			continue
		}

		if err := h.PostProcess(); err != nil {
			return errcode.ErrCode_ErrInternal.Wrap(err)
		}
	}

	return nil
}
