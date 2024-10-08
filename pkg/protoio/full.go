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
	"io"

	"google.golang.org/protobuf/proto"
)

func NewFullWriter(w io.Writer) WriteCloser {
	return &fullWriter{w, nil}
}

type fullWriter struct {
	w      io.Writer
	buffer []byte
}

func (writer *fullWriter) WriteMsg(msg proto.Message) (err error) {
	var data []byte
	if m, ok := msg.(marshaler); ok {
		n, ok := getSize(m)
		if !ok {
			_, err = proto.Marshal(msg)
			if err != nil {
				return err
			}
		}
		if n >= len(writer.buffer) {
			writer.buffer = make([]byte, n)
		}
		_, err = m.MarshalTo(writer.buffer)
		if err != nil {
			return err
		}
		data = writer.buffer[:n]
	} else {
		data, err = proto.Marshal(msg)
		if err != nil {
			return err
		}
	}
	_, err = writer.w.Write(data)
	return err
}

func (writer *fullWriter) Close() error {
	if closer, ok := writer.w.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

type fullReader struct {
	r   io.Reader
	buf []byte
}

func NewFullReader(r io.Reader, maxSize int) ReadCloser {
	return &fullReader{r, make([]byte, maxSize)}
}

func (reader *fullReader) ReadMsg(msg proto.Message) error {
	length, err := reader.r.Read(reader.buf)
	if err != nil {
		return err
	}
	return proto.Unmarshal(reader.buf[:length], msg)
}

func (reader *fullReader) Close() error {
	if closer, ok := reader.r.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
