package unpack

type Header struct {
	DirEntryCount int32
	DirCount      int32
	NamesOffset   int32
	BaseOffset    int32
}

type FileEntry struct {
	NameOffset   int64
	DirOffset    int64
	DataOffset   int64
	CompSize     int64
	UncompSize   int64
	IsCompressed bool
	Name         string
	Directory    string
}
