package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/alexflint/go-arg"
)

const (
	dirEntriesOffset = 0x78
	nullTerminator   = 0x0
	defaultOutPath   = "SRTools_extracted"
)

var magic = [4]byte{'\xCE', '\x0A', '\x89', '\x51'}

func contains(arr []string, v string) bool {
	for _, value := range arr {
		if strings.EqualFold(value, v) {
			return true
		}
	}
	return false
}

func filterInPaths(paths []string) ([]string, error) {
	var filtered []string
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	pathSep := fmt.Sprintf("%c", os.PathSeparator)
	for _, path := range paths {
		if !filepath.IsAbs(path) {
			path = filepath.Join(cwd, path)
		}
		path = strings.TrimSuffix(path, pathSep)
		if !contains(filtered, path) {
			filtered = append(filtered, path)
		} else {
			fmt.Println("Duplicate path filtered:", path)
		}
	}
	return filtered, nil
}

func parseArgs() (*Args, error) {
	var args Args
	arg.MustParse(&args)
	if !(args.Threads >= 1 && args.Threads <= 100) {
		return nil, errors.New("Max threads must be between 1 and 100.")
	}
	if args.OutPath == "" {
		args.OutPath = defaultOutPath
	}
	args.OutPath = filepath.Join(args.OutPath, "sr5")
	args.Command = strings.ToLower(args.Command)
	filteredPaths, err := filterInPaths(args.InPaths)
	if err != nil {
		return nil, err
	}
	args.InPaths = filteredPaths
	return &args, nil
}

func makeDirs(path string) error {
	err := os.MkdirAll(path, 0755)
	return err
}

func readUint32(f *os.File) (uint32, error) {
	buf := make([]byte, 4)
	_, err := f.Read(buf)
	if err != nil {
		return 0, err
	}
	value := binary.LittleEndian.Uint32(buf)
	return value, nil
}

func readUint64(f *os.File) (uint64, error) {
	buf := make([]byte, 8)
	_, err := f.Read(buf)
	if err != nil {
		return 0, err
	}
	value := binary.LittleEndian.Uint64(buf)
	return value, nil
}

func parseHeader(f *os.File) (*Header, error) {
	buf := make([]byte, 4)
	_, err := f.Read(buf)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(buf, magic[:]) {
		return nil, errors.New("File is not a packfile.")
	}
	version, err := readUint32(f)
	if err != nil {
		return nil, err
	}
	if version != 17 {
		return nil, errors.New("Unsupported packfile version.")
	}
	_, err = f.Seek(8, io.SeekCurrent)
	if err != nil {
		return nil, err
	}
	dirEntryCount, err := readUint32(f)
	if err != nil {
		return nil, err
	}
	dirCount, err := readUint32(f)
	if err != nil {
		return nil, err
	}
	namesOffset, err := readUint32(f)
	if err != nil {
		return nil, err
	}
	_, err = f.Seek(36, io.SeekCurrent)
	if err != nil {
		return nil, err
	}
	baseOffset, err := readUint32(f)
	if err != nil {
		return nil, err
	}
	header := &Header{
		DirEntryCount: dirEntryCount,
		DirCount:      dirCount,
		NamesOffset:   dirEntriesOffset + namesOffset,
		BaseOffset:    baseOffset,
	}
	return header, nil
}

func parseEntries(f *os.File, header *Header) ([]*FileEntry, error) {
	var entries []*FileEntry
	_, err := f.Seek(dirEntriesOffset, io.SeekStart)
	if err != nil {
		return nil, err
	}
	dirEntryCount := int(header.DirEntryCount)
	for i := 0; i < dirEntryCount; i++ {
		nameOffset, err := readUint64(f)
		if err != nil {
			return nil, err
		}
		dirOffset, err := readUint64(f)
		if err != nil {
			return nil, err
		}
		dataOffset, err := readUint64(f)
		if err != nil {
			return nil, err
		}
		uncompSize, err := readUint64(f)
		if err != nil {
			return nil, err
		}
		compSize, err := readUint64(f)
		if err != nil {
			return nil, err
		}
		// comp, err := readFlag(f)
		// if err != nil {
		// 	return nil, err
		// }
		// cond, err := readFlag(f)
		// if err != nil {
		// 	return nil, err
		// }
		isComp := compSize != math.MaxUint64
		if !isComp {
			compSize = uncompSize
		}
		_, err = f.Seek(8, io.SeekCurrent)
		if err != nil {
			return nil, err
		}
		entry := &FileEntry{
			NameOffset:   nameOffset,
			DirOffset:    dirOffset,
			DataOffset:   dataOffset,
			UncompSize:   uncompSize,
			CompSize:     compSize,
			IsCompressed: isComp,
			// IsCondensed:  cond == 1,
		}
		entries = append(entries, entry)
	}
	// _, err = f.Seek(int64(header.NamesOffset), io.SeekStart)
	// //_, err = f.Seek(int64(header.DirCount)*8, io.SeekCurrent)
	// if err != nil {
	// 	return nil, err
	// }
	return entries, nil
}

func getCurrentPos(f *os.File) (int64, error) {
	return f.Seek(0, io.SeekCurrent)
}

func readString(f *os.File, offset int64) (string, error) {
	var value string
	buf := make([]byte, 1)
	_, err := f.Seek(offset, io.SeekStart)
	if err != nil {
		return "", err
	}
	for {
		_, err := f.Read(buf)
		if err != nil {
			return "", err
		}
		if buf[0] == nullTerminator {
			break
		}
		value += string(buf[:])
	}
	return value, nil
}

func parseNamesAndDirs(f *os.File, entries []*FileEntry, namesOffset uint32) error {
	for _, entry := range entries {
		namesOffset := int64(namesOffset)
		offset := namesOffset + int64(entry.NameOffset)
		name, err := readString(f, offset)
		if err != nil {
			return err
		}
		entry.Name = name
		offset = namesOffset + int64(entry.DirOffset)
		dir, err := readString(f, offset)
		if err != nil {
			return err
		}
		entry.Directory = dir
	}
	return nil
}

// block dependency not supported for some.
// func decompress(data []byte, uncompSize uint64) ([]byte, error) {
// 	r := bytes.NewReader(data)
// 	buf := make([]byte, uncompSize)
// 	zr := lz4.NewReader(r)
// 	_, err := zr.Read(buf)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return buf, err
// }

func decompress(outPath string) error {
	decOutPath := outPath + "_dec"
	var (
		errBuffer bytes.Buffer
		args      = []string{"-d", outPath, decOutPath, "--rm", "-f"}
	)
	cmd := exec.Command("lz4", args...)
	cmd.Stderr = &errBuffer
	err := cmd.Run()
	if err != nil {
		errString := err.Error() + "\n" + errBuffer.String()
		return errors.New(errString)
	}
	err = os.Rename(decOutPath, outPath)
	return err
}

func writeFile(buf []byte, outPath string, isComp bool, uncompSize uint64) error {
	outFile, err := os.OpenFile(outPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	_, err = outFile.Write(buf)
	outFile.Close()
	if err != nil {
		return err
	}
	if isComp {
		err = decompress(outPath)
	}
	return err
}

func writeFiles(f *os.File, entries []*FileEntry, _outPath string, baseOffset, namesOffset uint32, threads int) error {
	var wg sync.WaitGroup
	ch := make(chan struct{}, threads)
	for _, entry := range entries {
		ch <- struct{}{}
		//i++
		//fmt.Printf("%d of %d:\n", i, entryTotal)
		outPath := filepath.Join(_outPath, entry.Directory)
		err := makeDirs(outPath)
		if err != nil {
			panic(err)
		}
		name := entry.Name
		isComp := entry.IsCompressed
		fullOutPath := filepath.Join(outPath, name)
		uncompSize := entry.UncompSize
		dataOffset := int64(baseOffset) + int64(entry.DataOffset)
		fmt.Println(filepath.Join(entry.Directory, name))
		fmt.Println("Start offset:", fmt.Sprintf("0x%X", dataOffset))
		fmt.Println("End offset:", fmt.Sprintf("0x%X", dataOffset+int64(uncompSize)))
		fmt.Println("Compressed size:", entry.CompSize, "bytes")
		fmt.Println("Uncompressed size:", uncompSize, "bytes")
		fmt.Println("Compressed:", isComp)
		fmt.Println("")
		wg.Add(1)
		go func(entry *FileEntry) {
			defer wg.Done()
			buf := make([]byte, entry.CompSize)
			_, err = f.ReadAt(buf, dataOffset)
			if err != nil {
				panic(err)
			}
			err = writeFile(buf, fullOutPath, isComp, uncompSize)
			if err != nil {
				panic(err)
			}
			<-ch
		}(entry)
	}
	wg.Wait()
	return nil
}

func unpack(args *Args) error {
	outPath := args.OutPath
	err := makeDirs(outPath)
	if err != nil {
		return err
	}
	for _, path := range args.InPaths {
		f, err := os.OpenFile(path, os.O_RDONLY, 0755)
		if err != nil {
			return err
		}
		defer f.Close()
		fmt.Println("Parsing header...")
		header, err := parseHeader(f)
		if err != nil {
			return err
		}
		fmt.Println("Parsing entries...")
		entries, err := parseEntries(f, header)
		if err != nil {
			return err
		}
		fmt.Println("Parsing name and directory strings...\n")
		err = parseNamesAndDirs(f, entries, header.NamesOffset)
		if err != nil {
			return err
		}
		err = writeFiles(
			f, entries, outPath, header.BaseOffset, header.NamesOffset, args.Threads)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	args, err := parseArgs()
	if err != nil {
		panic(err)
	}
	command := args.Command
	now := time.Now()
	switch command {
	case "unpack", "extract":
		err := unpack(args)
		if err != nil {
			fmt.Println("Unpacking failed.")
			panic(err)
		}
	default:
		panic("Unknown command: " + command)
	}
	fmt.Println("Finished in " + time.Since(now).String() + ".")
}
