package utils

type Args struct {
	Command       string   `arg:"positional, required"`
	InPaths       []string `arg:"-i, required" help:"Input path(s)."`
	OutPath       string   `arg:"-o" help:"Output path. Path will be made if it doesn't already exist."`
	Threads       int      `arg:"-t" default:"10" help:"Max threads (1-50). Be careful; memory intensive."`
	NoCompression bool     `arg:"-n" help:"Don't compress any files when packing. Might be a bit more stable."`
}
