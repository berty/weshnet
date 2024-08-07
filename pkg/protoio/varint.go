// Protocol Buffers for Go with Gadgets
//
// Copyright (c) 2013, The GoGo Authors. All rights reserved.
// http://github.com/gogo/protobuf
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package protoio

import (
	"bufio"
	"encoding/binary"
	"io"

	"google.golang.org/protobuf/proto"
)

func NewDelimitedWriter(w io.Writer) WriteCloser {
	return &varintWriter{w, make([]byte, binary.MaxVarintLen64), nil}
}

type varintWriter struct {
	w      io.Writer
	lenBuf []byte
	buffer []byte
}

func (writer *varintWriter) WriteMsg(msg proto.Message) (err error) {
	var data []byte
	if m, ok := msg.(marshaler); ok {
		n, ok := getSize(m)
		if ok {
			if n+binary.MaxVarintLen64 >= len(writer.buffer) {
				writer.buffer = make([]byte, n+binary.MaxVarintLen64)
			}
			lenOff := binary.PutUvarint(writer.buffer, uint64(n))
			_, err = m.MarshalTo(writer.buffer[lenOff:])
			if err != nil {
				return err
			}
			_, err = writer.w.Write(writer.buffer[:lenOff+n])
			return err
		}
	}

	// fallback
	data, err = proto.Marshal(msg)
	if err != nil {
		return err
	}
	length := uint64(len(data))
	n := binary.PutUvarint(writer.lenBuf, length)
	_, err = writer.w.Write(writer.lenBuf[:n])
	if err != nil {
		return err
	}
	_, err = writer.w.Write(data)
	return err
}

func (writer *varintWriter) Close() error {
	if closer, ok := writer.w.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

func NewDelimitedReader(r io.Reader, maxSize int) ReadCloser {
	var closer io.Closer
	if c, ok := r.(io.Closer); ok {
		closer = c
	}
	return &varintReader{bufio.NewReader(r), nil, maxSize, closer}
}

type varintReader struct {
	r       *bufio.Reader
	buf     []byte
	maxSize int
	closer  io.Closer
}

func (reader *varintReader) ReadMsg(msg proto.Message) error {
	length64, err := binary.ReadUvarint(reader.r)
	if err != nil {
		return err
	}
	length := int(length64)
	if length < 0 || length > reader.maxSize {
		return io.ErrShortBuffer
	}
	if len(reader.buf) < length {
		reader.buf = make([]byte, length)
	}
	buf := reader.buf[:length]
	if _, err := io.ReadFull(reader.r, buf); err != nil {
		return err
	}
	return proto.Unmarshal(buf, msg)
}

func (reader *varintReader) Close() error {
	if reader.closer != nil {
		return reader.closer.Close()
	}
	return nil
}
