package socket

import (
	"encoding/binary"
	"fmt"
	"io"
)

var magicBytes = [3]byte{88, 77, 68}

var emptySignature = make([]byte, 128)

const PROTOCOL_VERSION byte = 101

// 日志包
type LogPackage struct {
	Magic          [3]byte // 3 byte 魔数
	Version        byte    // 1 byte 协议版本
	Type           byte    // 1 byte 包类型
	HostNameLength int32   //  主机名称长度
	BodySize       int32   // 4 byte 数据部分长度,仅指 body 长度
	TimeStamp      int64   // 8 byte 时间戳
	HostName       []byte  //  主机部分长度
	Signature      []byte  // 128 byte 签名
	Body           []byte  // 数据部分长度
}

// Pack 编码
func (l *LogPackage) Pack(writer io.Writer) error {
	var err error
	err = binary.Write(writer, binary.BigEndian, &l.Magic)
	err = binary.Write(writer, binary.BigEndian, &l.Version)
	err = binary.Write(writer, binary.BigEndian, &l.Type)
	err = binary.Write(writer, binary.BigEndian, &l.HostNameLength)
	err = binary.Write(writer, binary.BigEndian, &l.BodySize)
	err = binary.Write(writer, binary.BigEndian, &l.TimeStamp)
	err = binary.Write(writer, binary.BigEndian, &l.HostName)
	err = binary.Write(writer, binary.BigEndian, &l.Signature)
	err = binary.Write(writer, binary.BigEndian, &l.Body)
	return err
}

func (l *LogPackage) Unpack(reader io.Reader) error {
	var err error
	err = binary.Read(reader, binary.BigEndian, &l.Magic)
	err = binary.Read(reader, binary.BigEndian, &l.Version)
	err = binary.Read(reader, binary.BigEndian, &l.Type)
	err = binary.Read(reader, binary.BigEndian, &l.HostNameLength)
	err = binary.Read(reader, binary.BigEndian, &l.BodySize)
	err = binary.Read(reader, binary.BigEndian, &l.TimeStamp)

	l.HostName = make([]byte, l.HostNameLength)
	err = binary.Read(reader, binary.BigEndian, &l.HostName)

	l.Signature = emptySignature
	err = binary.Read(reader, binary.BigEndian, &l.Signature)

	l.Body = make([]byte, l.BodySize)
	err = binary.Read(reader, binary.BigEndian, &l.Body)
	return err
}

func (l *LogPackage) String() string {
	return fmt.Sprintf("magic:%s version:%b type:%d bodySize:%d timestamp:%d signature:%s body:%s", l.Magic, l.Version, l.Type, l.BodySize, l.TimeStamp, l.Signature, l.Body)
}
