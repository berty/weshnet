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
	"encoding/binary"
	"io"

	"google.golang.org/protobuf/proto"
)

const uint32BinaryLen = 4

func NewUint32DelimitedWriter(w io.Writer, byteOrder binary.ByteOrder) WriteCloser {
	return &uint32Writer{w, byteOrder, nil, make([]byte, uint32BinaryLen)}
}

func NewSizeUint32DelimitedWriter(w io.Writer, byteOrder binary.ByteOrder, size int) WriteCloser {
	return &uint32Writer{w, byteOrder, make([]byte, size), make([]byte, uint32BinaryLen)}
}

type uint32Writer struct {
	w         io.Writer
	byteOrder binary.ByteOrder
	buffer    []byte
	lenBuf    []byte
}

func (writer *uint32Writer) writeFallback(msg proto.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	length := uint32(len(data))
	writer.byteOrder.PutUint32(writer.lenBuf, length)
	if _, err = writer.w.Write(writer.lenBuf); err != nil {
		return err
	}
	_, err = writer.w.Write(data)
	return err
}

func (writer *uint32Writer) WriteMsg(msg proto.Message) error {
	m, ok := msg.(marshaler)
	if !ok {
		return writer.writeFallback(msg)
	}

	n, ok := getSize(m)
	if !ok {
		return writer.writeFallback(msg)
	}

	size := n + uint32BinaryLen
	if size > len(writer.buffer) {
		writer.buffer = make([]byte, size)
	}

	writer.byteOrder.PutUint32(writer.buffer, uint32(n))
	if _, err := m.MarshalTo(writer.buffer[uint32BinaryLen:]); err != nil {
		return err
	}

	_, err := writer.w.Write(writer.buffer[:size])
	return err
}

func (writer *uint32Writer) Close() error {
	if closer, ok := writer.w.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

type uint32Reader struct {
	r         io.Reader
	byteOrder binary.ByteOrder
	lenBuf    []byte
	buf       []byte
	maxSize   int
}

func NewUint32DelimitedReader(r io.Reader, byteOrder binary.ByteOrder, maxSize int) ReadCloser {
	return &uint32Reader{r, byteOrder, make([]byte, 4), nil, maxSize}
}

func (reader *uint32Reader) ReadMsg(msg proto.Message) error {
	if _, err := io.ReadFull(reader.r, reader.lenBuf); err != nil {
		return err
	}
	length32 := reader.byteOrder.Uint32(reader.lenBuf)
	length := int(length32)
	if length < 0 || length > reader.maxSize {
		return io.ErrShortBuffer
	}
	if length > len(reader.buf) {
		reader.buf = make([]byte, length)
	}
	_, err := io.ReadFull(reader.r, reader.buf[:length])
	if err != nil {
		return err
	}
	return proto.Unmarshal(reader.buf[:length], msg)
}

func (reader *uint32Reader) Close() error {
	if closer, ok := reader.r.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
