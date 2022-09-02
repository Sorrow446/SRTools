package pack

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"main/utils"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	null           = []byte{'\x00'}
	pathSep        = fmt.Sprintf("%c", os.PathSeparator)
	defaultOutPath = "SRTools_packed.vpp_pc"
)

var (
	extensions = []string{
		".bk2", ".bik", ".str2_pc", ".strh_pc", ".cvbm_pc", ".gvbm_pc",
		".vpp_pc", ".vpkg", ".ttf", ".xml", ".ridv_pc", ".lua",
		".wem_pad", ".hkcmp_64m", ".refl_xbox3_gdk", ".refl_pc",
		".sr2_h", ".bin", ".txt", ".vint_proj",
	}
	extensionsTwo = []string{".fxo_dx11_pc", ".fxo_dx12_pc", ".fxo_vk_pc"}
	//extensionsTwo   = []string{".sr2_h"}
	extensionsThree = []string{".bk2"}
)

func processArgs(args *utils.Args) (*utils.Args, error) {
	if args.OutPath == "" {
		args.OutPath = defaultOutPath
	}
	if !(strings.HasSuffix(args.OutPath, ".vpp_pc") || strings.HasSuffix(args.OutPath, ".str2_pc")) {
		return nil, errors.New("Invalid output file file extension.")
	}
	return args, nil
}

func contains(arr []string, v string) bool {
	for _, value := range arr {
		if value == v {
			return true
		}
	}
	return false
}

func add(dirs *Dirs, path string, file *File) {
	for i, dir := range dirs.Dirs {
		if dir.Name == path {
			dirs.Dirs[i].Files = append(dirs.Dirs[i].Files, file)
			return
		}
	}
	dir := &Dir{
		Name:  path,
		Files: []*File{file},
	}
	dirs.Dirs = append(dirs.Dirs, dir)
}

func hasExtension(fname string, extensions []string) bool {
	for _, ext := range extensions {
		if strings.HasSuffix(fname, ext) {
			return true
		}
	}
	return false
}

func getAlign(fname string) int16 {
	var align int16
	switch {
	case hasExtension(fname, extensionsTwo):
		align = 16
	case hasExtension(fname, extensionsThree):
		align = 2048
	default:
		align = 1
	}
	return align
}

func populateDirs(packFolder, tempPath string, compressAll, noCompression bool) (*Dirs, error) {
	var fileTotal int
	dirs := &Dirs{
		Dirs: []*Dir{},
	}
	if !noCompression {
		fmt.Println("Compression is enabled, this may take a while for large packfiles.")
		fmt.Println("Compressing files...")
	}
	err := filepath.Walk(packFolder, func(path string, f os.FileInfo, err error) error {
		if path == packFolder {
			return nil
		}
		if !f.IsDir() {
			var (
				flag           int16
				align          int16
				shouldCompress bool
			)
			folder := "sr5"
			idx := strings.Index(path, pathSep+"data"+pathSep)
			if idx == -1 {
				return nil
			}
			path = path[idx+1:]
			path = filepath.Dir(path)
			if strings.HasPrefix(path, "data"+pathSep+"engine") {
				path = `..\ctg\` + path
				folder = "ctg"
			}
			fname := f.Name()
			if noCompression {
				shouldCompress = false
			} else if compressAll {
				shouldCompress = true
			} else {
				// ".bnk_pad"
				shouldCompress = !hasExtension(fname, extensions)
			}
			if shouldCompress {
				flag = 1
				if noCompression {
					align = 1
				} else {
					align = getAlign(fname)
				}
			} else {
				align = getAlign(fname)
			}

			file := &File{
				Name:           fname,
				Size:           f.Size(),
				FullPath:       filepath.Join(packFolder, folder, path, fname),
				ShouldCompress: shouldCompress,
				Flag:           flag,
				Alignment:      align,
			}
			if shouldCompress {
				compPath, compSize, err := compress(file.FullPath, tempPath)
				if err != nil {
					return err
				}
				file.CompressedPath = compPath
				file.CompressedSize = compSize
			}
			add(dirs, path, file)
			fileTotal++
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	dirs.FileTotal = fileTotal
	return dirs, nil
}

func writeUint16(f *os.File, value int16) error {
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, uint16(value))
	_, err := f.Write(buf)
	return err
}

func writeUint32(f *os.File, value int32) error {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(value))
	_, err := f.Write(buf)
	return err
}

func writeUint64(f *os.File, value int64) error {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(value))
	_, err := f.Write(buf)
	return err
}

func getCurrentPos(f *os.File) (int64, error) {
	return f.Seek(0, io.SeekCurrent)
}

func getPackFolder(path string) string {
	lastIdx := strings.LastIndex(path, pathSep)
	if lastIdx == -1 {
		return path
	} else {
		return path[lastIdx+len(pathSep):]
	}
}

func compress(path, tempPath string) (string, int64, error) {
	outPath := filepath.Join(tempPath, path)
	err := os.MkdirAll(filepath.Dir(outPath), 0755)
	if err != nil {
		return "", -1, err
	}
	var (
		errBuffer bytes.Buffer
		args      = []string{"-9", "-BS5", path, outPath}
	)
	cmd := exec.Command("lz4", args...)
	cmd.Stderr = &errBuffer
	err = cmd.Run()
	if err != nil {
		errString := err.Error() + "\n" + errBuffer.String()
		return "", -1, errors.New(errString)
	}
	f, err := os.Stat(outPath)
	if err != nil {
		return "", -1, err
	}
	return outPath, f.Size(), nil
}

func getTempPath() (string, error) {
	return os.MkdirTemp(os.TempDir(), "")
}

// Clean up.
func Run(args *utils.Args) error {
	args, err := processArgs(args)
	if err != nil {
		return err
	}
	outPath := args.OutPath
	compressAll := strings.HasSuffix(outPath, ".str2_pc")
	tempPath, err := getTempPath()
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempPath)
	packFolder := getPackFolder(args.InPaths[0])
	fmt.Println("Populating paths...")
	dirs, err := populateDirs(packFolder, tempPath, compressAll, args.NoCompression)
	if err != nil {
		return err
	}
	var (
		nameTable []byte
		dataSize  int64
	)
	fmt.Println("Building name and directory string table...")
	for idx, dir := range dirs.Dirs {
		dirs.Dirs[idx].NameOffset = int64(len(nameTable))
		nameTable = append(nameTable, []byte(dir.Name)...)
		nameTable = append(nameTable, null...)
		for fileIdx, file := range dir.Files {
			dirs.Dirs[idx].Files[fileIdx].NameOffset = int64(len(nameTable))
			nameTable = append(nameTable, []byte(file.Name)...)
			nameTable = append(nameTable, null...)
			dirs.Dirs[idx].Files[fileIdx].DataOffset = dataSize
			if file.ShouldCompress {
				dataSize += file.CompressedSize
			} else {
				dataSize += file.Size
			}
		}
	}

	f, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	fmt.Println("Writing header...")
	// magic
	_, err = f.Write([]byte{'\xCE', '\x0A', '\x89', '\x51'})
	if err != nil {
		return err
	}
	// version
	err = writeUint32(f, 17)
	if err != nil {
		return err
	}
	// crc
	_, err = f.Write(bytes.Repeat(null, 4))
	if err != nil {
		return err
	}
	// flags
	err = writeUint32(f, 20481)
	if err != nil {
		return err
	}
	// file count
	err = writeUint32(f, int32(dirs.FileTotal))
	if err != nil {
		return err
	}
	// dir count
	err = writeUint32(f, int32(len(dirs.Dirs)))
	if err != nil {
		return err
	}
	namesOffset, err := getCurrentPos(f)
	if err != nil {
		return err
	}
	// names offset
	err = writeUint32(f, 0)
	if err != nil {
		return err
	}
	// names size
	err = writeUint32(f, 0)
	if err != nil {
		return err
	}
	packSize, err := getCurrentPos(f)
	if err != nil {
		return err
	}
	// pack size
	err = writeUint64(f, 0)
	if err != nil {
		return err
	}
	// data size
	err = writeUint64(f, dataSize)
	if err != nil {
		return err
	}
	// compressed data size
	err = writeUint64(f, dataSize)
	if err != nil {
		return err
	}
	// epoch timestamp
	err = writeUint64(f, 0)
	if err != nil {
		return err
	}
	dataOffsetBase, err := getCurrentPos(f)
	if err != nil {
		return err
	}
	// data offset base
	err = writeUint64(f, 0)
	if err != nil {
		return err
	}
	// reserved
	_, err = f.Write(bytes.Repeat(null, 48))
	if err != nil {
		return err
	}
	fmt.Println("Writing file entries...")
	for _, dir := range dirs.Dirs {
		for _, file := range dir.Files {
			err = writeUint64(f, file.NameOffset)
			if err != nil {
				return err
			}
			err = writeUint64(f, dir.NameOffset)
			if err != nil {
				return err
			}
			err = writeUint64(f, file.DataOffset)
			if err != nil {
				return err
			}
			err = writeUint64(f, file.Size)
			if err != nil {
				return err
			}
			if file.ShouldCompress {
				err = writeUint64(f, file.CompressedSize)
			} else {
				_, err = f.Write(bytes.Repeat([]byte{'\xFF'}, 8))
			}
			if err != nil {
				return err
			}
			err = writeUint16(f, file.Flag)
			if err != nil {
				return err
			}
			err = writeUint16(f, file.Alignment)
			if err != nil {
				return err
			}
			_, err = f.Write(bytes.Repeat(null, 4))
			if err != nil {
				return err
			}
		}
	}
	fmt.Println("Writing directory name offsets...")
	for _, dir := range dirs.Dirs {
		err = writeUint64(f, dir.NameOffset)
		if err != nil {
			return err
		}
	}
	curPos, err := getCurrentPos(f)
	if err != nil {
		return err
	}
	_, err = f.Seek(namesOffset, io.SeekStart)
	if err != nil {
		return err
	}
	err = writeUint32(f, int32(curPos-120))
	if err != nil {
		return err
	}
	err = writeUint32(f, int32(len(nameTable)))
	if err != nil {
		return err
	}
	_, err = f.Seek(curPos, io.SeekStart)
	if err != nil {
		return err
	}
	_, err = f.Write(nameTable)
	if err != nil {
		return err
	}
	curPos, err = getCurrentPos(f)
	if err != nil {
		return err
	}
	_, err = f.Seek(dataOffsetBase, io.SeekStart)
	if err != nil {
		return err
	}
	err = writeUint32(f, int32(curPos))
	if err != nil {
		return err
	}
	_, err = f.Seek(curPos, io.SeekStart)
	if err != nil {
		return err
	}
	fmt.Println("Writing files...")
	i := 1
	for _, dir := range dirs.Dirs {
		for _, file := range dir.Files {
			fmt.Printf("\r%d of %d.", i, dirs.FileTotal)
			var path string
			if file.ShouldCompress {
				path = file.CompressedPath
			} else {
				path = file.FullPath
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			// Crashes game.
			// align := int64(file.Alignment)
			// if align > 0 && align != 16 {
			// 	fmt.Println(file.Name)
			// 	curPos, err := getCurrentPos(f)
			// 	if err != nil {
			// 		return err
			// 	}
			// 	pos := (curPos + align - 1) / align * align
			// 	_, err = f.Seek(pos, io.SeekStart)
			// 	if err != nil {
			// 		return err
			// 	}
			// }
			_, err = f.Write(data)
			if err != nil {
				return err
			}
			if file.ShouldCompress {
				err = os.Remove(path)
				if err != nil {
					fmt.Println("Failed to delete compressed file:", path)
				}
			}
			i++
		}
	}
	fmt.Println("")
	curPos, err = getCurrentPos(f)
	if err != nil {
		return err
	}
	_, err = f.Seek(packSize, io.SeekStart)
	if err != nil {
		return err
	}
	err = writeUint32(f, int32(curPos))
	return err
}
