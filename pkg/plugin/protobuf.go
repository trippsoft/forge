// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package plugin

import (
	"encoding/binary"
	"io"

	"google.golang.org/protobuf/proto"
)

func Write(w io.Writer, m proto.Message) error {
	data, err := proto.Marshal(m)
	if err != nil {
		return err
	}

	length := uint32(len(data))
	err = binary.Write(w, binary.BigEndian, length)
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}

func Read(r io.Reader, m proto.Message) error {
	var length uint32
	err := binary.Read(r, binary.BigEndian, &length)
	if err != nil {
		return err
	}

	data := make([]byte, length)
	_, err = io.ReadFull(r, data)
	if err != nil {
		return err
	}

	return proto.Unmarshal(data, m)
}
