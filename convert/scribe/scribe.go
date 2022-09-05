package scribe

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"main/utils"
	"os"
	"strconv"
)

const startOffset = 0x140

var magic = [4]byte{0x54, 0x38, 0x09, 0x00}

func checkMagic(f *os.File) (bool, error) {
	buf := make([]byte, 4)
	_, err := f.Read(buf)
	if err != nil {
		return false, err
	}
	return bytes.Equal(buf, magic[:]), nil
}

func readString(f *os.File, length int32, isType bool) (string, error) {
	buf := make([]byte, length)
	_, err := f.Read(buf)
	if err != nil {
		return "", err
	}
	if isType {
		return string(buf[:len(buf)-1]), nil
	}
	return string(buf), nil
}

func setTextPos(f *os.File) error {
	buf := make([]byte, 1)
	for {
		_, err := f.Read(buf)
		if err != nil {
			return err
		}
		if buf[0] == 0x20 {
			_, err := f.Seek(-6, io.SeekCurrent)
			return err
		}
	}
}

func parseEntryHeader(f *os.File, scribe *Scribe) error {
	unk, err := utils.ReadBytes(f, 8)
	if err != nil {
		return err
	}
	entriesSize, err := utils.ReadUint32(f)
	if err != nil {
		return err
	}
	unkTwo, err := utils.ReadUint32(f)
	if err != nil {
		return err
	}
	scribe.EntriesHeader = &EntriesHeader{
		Unk:         unk,
		EntriesSize: entriesSize,
		UnkTwo:      unkTwo,
	}
	return nil
}

func parseEntries(f *os.File, scribe *Scribe, endPos int64) error {
	for {
		startPos, err := utils.GetCurrentPos(f)
		if err != nil {
			return err
		}
		entrySize, err := utils.ReadUint32(f)
		if err != nil {
			return err
		}
		unk, err := utils.ReadUint32(f)
		if err != nil {
			return err
		}
		unkTwo, err := utils.ReadBytes(f, 4)
		if err != nil {
			return err
		}
		_, err = f.Seek(12, io.SeekCurrent)
		if err != nil {
			return err
		}
		TypeStringLen, err := utils.ReadUint32(f)
		if err != nil {
			return err
		}
		typeString, err := readString(f, TypeStringLen, true)
		if err != nil {
			return err
		}
		fmt.Println(typeString)
		err = setTextPos(f)
		if err != nil {
			return err
		}
		textLen, err := utils.ReadUint32(f)
		if err != nil {
			return err
		}
		_, err = f.Seek(4, io.SeekCurrent)
		if err != nil {
			return err
		}
		text, err := readString(f, textLen, false)
		if err != nil {
			return err
		}
		scribe.Entries = append(scribe.Entries, &Entry{
			TypeString: typeString,
			Text:       text,
			Unk:        unk,
			UnkTwo:     unkTwo,
		})
		fmt.Println(text)
		nextEntryOffset := startPos + int64(entrySize)
		if nextEntryOffset-startOffset-8 == int64(scribe.EntriesHeader.EntriesSize) {
			_, err = f.Seek(4, io.SeekCurrent)
			if err != nil {
				return err
			}
			break
		}
		_, err = f.Seek(nextEntryOffset, io.SeekStart)
		if err != nil {
			return err
		}
	}
	curPos, err := utils.GetCurrentPos(f)
	if err != nil {
		return err
	}
	endData, err := utils.ReadBytes(f, endPos-curPos)
	if err != nil {
		return err
	}
	scribe.EndData = endData
	return nil
}

func test(scribe *Scribe) {
	var entriesSize int32 = 8
	for idx, _ := range scribe.Entries {
		typeString := "abc_" + strconv.Itoa(idx)
		TypeStringLen := int32(len(typeString)) + 1
		text := "def_" + strconv.Itoa(idx)
		textLen := int32(len(text))
		scribe.Entries[idx].TypeString = typeString
		scribe.Entries[idx].TypeStringLen = TypeStringLen
		scribe.Entries[idx].Text = text
		scribe.Entries[idx].TextLen = textLen
		//entrySize := 28 + 4 + TypeStringLen + 2 + 4 + textLen + 2
		entrySize := TypeStringLen + textLen + 40
		scribe.Entries[idx].Size = entrySize
		entriesSize += entrySize
	}
	scribe.EntriesHeader.EntriesSize = entriesSize
}

func writeScribe(scribe *Scribe, outPath string) error {
	f, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(magic[:])
	if err != nil {
		return err
	}
	_, err = f.Write(bytes.Repeat([]byte{0x0}, 316))
	if err != nil {
		return err
	}
	_, err = f.Write(scribe.EntriesHeader.Unk)
	if err != nil {
		return err
	}
	err = utils.WriteUint32(f, scribe.EntriesHeader.EntriesSize)
	if err != nil {
		return err
	}
	err = utils.WriteUint32(f, scribe.EntriesHeader.UnkTwo)
	if err != nil {
		return err
	}
	for _, entry := range scribe.Entries {
		err = utils.WriteUint32(f, entry.Size)
		if err != nil {
			return err
		}
		err = utils.WriteUint32(f, entry.Unk)
		if err != nil {
			return err
		}
		_, err = f.Write(entry.UnkTwo)
		if err != nil {
			return err
		}
		err = utils.WriteNull(f, 12)
		if err != nil {
			return err
		}
		err = utils.WriteUint32(f, entry.TypeStringLen)
		if err != nil {
			return err
		}
		_, err = f.WriteString(entry.TypeString)
		if err != nil {
			return err
		}
		err = utils.WriteNull(f, 3)
		if err != nil {
			return err
		}
		err = utils.WriteUint32(f, entry.TextLen)
		if err != nil {
			return err
		}
		err = utils.WriteUint32(f, 8192)
		if err != nil {
			return err
		}
		_, err = f.WriteString(entry.Text)
		if err != nil {
			return err
		}
		err = utils.WriteNull(f, 2)
		if err != nil {
			return err
		}
	}
	_, err = f.Write(scribe.EndData)
	return err
}

func parseScribeJson(path string) (*Scribe, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var obj Scribe
	err = json.Unmarshal(data, &obj)
	if err != nil {
		return nil, err
	}
	decUnk, err := utils.B64Decode(obj.EntriesHeader.UnkB64)
	if err != nil {
		return nil, err
	}
	decEndData, err := utils.B64Decode(obj.EndDataB64)
	if err != nil {
		return nil, err
	}
	obj.EntriesHeader.Unk = decUnk
	obj.EndData = decEndData
	var entriesSize int32 = 8
	for idx, entry := range obj.Entries {
		decUnkTwo, err := utils.B64Decode(entry.UnkTwoB64)
		if err != nil {
			return nil, err
		}
		obj.Entries[idx].UnkTwo = decUnkTwo
		typeStringLen := int32(len(entry.TypeString)) + 1
		textLen := int32(len(entry.Text))
		entrySize := typeStringLen + textLen + 40
		obj.Entries[idx].TypeStringLen = typeStringLen
		obj.Entries[idx].TextLen = textLen
		obj.Entries[idx].Size = entrySize
		entriesSize += entrySize
	}
	obj.EntriesHeader.EntriesSize = entriesSize
	return &obj, nil
}

func writeJson(scribe *Scribe, outPath string) error {
	outScribe := scribe
	outScribe.EntriesHeader.EntriesSize = 0
	outScribe.EntriesHeader.UnkB64 = utils.B64Encode(outScribe.EntriesHeader.Unk)
	outScribe.EntriesHeader.Unk = nil
	outScribe.EndDataB64 = utils.B64Encode(outScribe.EndData)
	outScribe.EndData = nil
	for idx, entry := range outScribe.Entries {
		outScribe.Entries[idx].UnkTwoB64 = utils.B64Encode(entry.UnkTwo)
		outScribe.Entries[idx].UnkTwo = nil
	}
	m, err := json.MarshalIndent(outScribe, "", "\t")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(outPath, m, 0755)
	return err
}

// JSON to scribe.
func To(args *utils.Args) error {
	scribe, err := parseScribeJson(args.InPaths[0])
	if err != nil {
		return err
	}
	err = writeScribe(scribe, args.OutPath)
	return err
}

// Scribe to JSON.
func From(args *utils.Args) error {
	scribe := &Scribe{
		EntriesHeader: &EntriesHeader{},
		Entries:       []*Entry{},
	}
	f, err := os.OpenFile(args.InPaths[0], os.O_RDONLY, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	ok, err := checkMagic(f)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("File is not a scribe file.")
	}
	stat, err := f.Stat()
	if err != nil {
		return err
	}
	_, err = f.Seek(startOffset, io.SeekStart)
	if err != nil {
		return err
	}
	err = parseEntryHeader(f, scribe)
	if err != nil {
		return err
	}
	err = parseEntries(f, scribe, stat.Size())
	if err != nil {
		return err
	}
	err = writeJson(scribe, args.OutPath)
	return err
}
