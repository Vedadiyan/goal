package codecs

import (
	"bytes"
	"fmt"
	"io"

	"github.com/klauspost/compress/zstd"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type ProtoConn struct{}

type CompressedProtoConn struct{}

func NewProtoConn(conn *nats.Conn) (*nats.EncodedConn, error) {
	const encType = "ProtoConn"
	nats.RegisterEncoder(encType, ProtoConn{})
	return nats.NewEncodedConn(conn, encType)
}

func NewCompressedProtoConn(conn *nats.Conn) (*nats.EncodedConn, error) {
	const encType = "CompressedProtoConn"
	nats.RegisterEncoder(encType, CompressedProtoConn{})
	return nats.NewEncodedConn(conn, encType)
}

func (ProtoConn) Encode(subject string, v interface{}) ([]byte, error) {
	value, ok := v.(protoreflect.ProtoMessage)
	if !ok {
		return nil, fmt.Errorf("the type '%T' is not a registered protobuffer type", v)
	}
	return proto.Marshal(value)
}

func (ProtoConn) Decode(subject string, data []byte, vPtr interface{}) error {
	value, ok := vPtr.(protoreflect.ProtoMessage)
	if !ok {
		return fmt.Errorf("the type '%T' is not a registered protobuffer type", vPtr)
	}
	data, err := decompress(data)
	if err != nil {
		return err
	}
	return proto.Unmarshal(data, value)
}

func (CompressedProtoConn) Encode(subject string, v interface{}) ([]byte, error) {
	value, ok := v.(protoreflect.ProtoMessage)
	if !ok {
		return nil, fmt.Errorf("the type '%T' is not a registered protobuffer type", v)
	}
	output, err := proto.Marshal(value)
	if err != nil {
		return nil, err
	}
	return compress(output)
}

func (CompressedProtoConn) Decode(subject string, data []byte, vPtr interface{}) error {
	value, ok := vPtr.(protoreflect.ProtoMessage)
	if !ok {
		return fmt.Errorf("the type '%T' is not a registered protobuffer type", vPtr)
	}
	data, err := decompress(data)
	if err != nil {
		return err
	}
	err = proto.Unmarshal(data, value)
	if err != nil {
		return err
	}
	return nil
}

func compress(in []byte) ([]byte, error) {
	var out bytes.Buffer
	enc, err := zstd.NewWriter(&out)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(enc, bytes.NewBuffer(in))
	if err != nil {
		_ = enc.Close()
		return nil, err
	}
	err = enc.Close()
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func decompress(in []byte) ([]byte, error) {
	var out bytes.Buffer
	enc, err := zstd.NewReader(bytes.NewBuffer(in))
	if err != nil {
		return nil, err
	}
	defer enc.Close()
	_, err = io.Copy(&out, enc)
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
