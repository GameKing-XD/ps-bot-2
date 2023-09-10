package main

import (
	"github.com/mitaka8/dca"
	"os"

	"io"
)

func main() {
	tmpFile := os.Args[1]
	outfile := os.Args[2]
	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 96
	options.Application = "lowdelay"

	encodingSession, err := dca.EncodeFile(tmpFile, options)
	if err != nil {
		panic(err)
	}
	defer encodingSession.Cleanup()

	output, err := os.Create(outfile)
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(output, encodingSession)

}
