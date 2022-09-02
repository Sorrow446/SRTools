package pack

type File struct {
	Name           string
	Size           int64
	NameOffset     int64
	DataOffset     int64
	FullPath       string
	CompressedPath string
	CompressedSize int64
	ShouldCompress bool
	Flag           int16
	Alignment      int16
}

type Dir struct {
	NameOffset int64
	Name       string
	Files      []*File
}

type Dirs struct {
	FileTotal int
	Dirs      []*Dir
}
