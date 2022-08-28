package main

type Args struct {
	Command string   `arg:"positional, required"`
	InPaths []string `arg:"-i, required" help:"Input path of packfile."`
	OutPath string   `arg:"-o" help:"Output path. Path will be made if it doesn't already exist."`
	Threads int      `arg:"-t" default:"10" help:"Max threads (1-100)."`
}

type Header struct {
	DirEntryCount uint32
	DirCount      uint32
	NamesOffset   uint32
	BaseOffset    uint32
}

type FileEntry struct {
	NameOffset   uint64
	DirOffset    uint64
	DataOffset   uint64
	CompSize     uint64
	UncompSize   uint64
	IsCompressed bool
	IsCondensed  bool
	Name         string
	Directory    string
}
