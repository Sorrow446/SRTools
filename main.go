package main

import (
	"fmt"
	"main/pack"
	"main/unpack"
	"main/utils"
	"strings"
	"time"

	"github.com/alexflint/go-arg"
)

func parseArgs() (*utils.Args, error) {
	var args utils.Args
	arg.MustParse(&args)
	args.Command = strings.ToLower(args.Command)
	return &args, nil
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
		err = unpack.Run(args)
	case "pack":
		err = pack.Run(args)
	default:
		panic("Unknown command: " + command)
	}
	if err != nil {
		panic(err)
	}
	fmt.Println("Finished in " + time.Since(now).String() + ".")
}
