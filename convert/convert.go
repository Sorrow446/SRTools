package convert

import (
	"errors"
	"main/convert/scribe"
	"main/utils"
	"strings"
)

func checkSuffixes(inPath, outPath, suffOne, suffTwo string) bool {
	return strings.HasSuffix(inPath, suffOne) && strings.HasSuffix(outPath, suffTwo)
}

func Run(args *utils.Args) error {
	var err error
	inPath := args.InPaths[0]
	outPath := args.OutPath
	switch {
	case checkSuffixes(inPath, outPath, ".json", ".scribe_pad"):
		err = scribe.To(args)
	case checkSuffixes(inPath, outPath, ".scribe_pad", ".json"):
		err = scribe.From(args)
	default:
		err = errors.New("Invalid input and output file extension combination.")
	}
	return err
}
