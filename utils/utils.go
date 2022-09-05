package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"io"
	"os"
)

func GetCurrentPos(f *os.File) (int64, error) {
	return f.Seek(0, io.SeekCurrent)
}

func ReadUint32(f *os.File) (int32, error) {
	buf := make([]byte, 4)
	_, err := f.Read(buf)
	if err != nil {
		return 0, err
	}
	value := binary.LittleEndian.Uint32(buf)
	return int32(value), nil
}

func ReadUint64(f *os.File) (int64, error) {
	buf := make([]byte, 8)
	_, err := f.Read(buf)
	if err != nil {
		return 0, err
	}
	value := binary.LittleEndian.Uint64(buf)
	return int64(value), nil
}

func WriteUint32(f *os.File, value int32) error {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(value))
	_, err := f.Write(buf)
	return err
}

func ReadBytes(f *os.File, bytesLen int64) ([]byte, error) {
	buf := make([]byte, bytesLen)
	_, err := f.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func WriteNull(f *os.File, rep int) error {
	_, err := f.Write(bytes.Repeat([]byte{0x0}, rep))
	return err
}

func B64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func B64Decode(str string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(str)
}
