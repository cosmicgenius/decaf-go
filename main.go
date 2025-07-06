package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/pflag"

	"github.com/cosmicgenius/decaf-go/config"
	"github.com/cosmicgenius/decaf-go/output"
)

func main() {
	var (
		stage      = pflag.String("target", "", "Target stage (scan, parse, inter, assemble, or empty to run all)")
		outputPath = pflag.String("output", "", "Output file. If not set, output to stdout.")
		opt        = pflag.String("opt", "", "Optimizations")
		debug      = pflag.Bool("debug", false, "Debug mode")
	)

	// Parse flags
	pflag.Parse()

	// Check for exactly one non-flag argument
	if pflag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s <filename> [flags]\n", os.Args[0])
		os.Exit(1)
	}

	// Get the filename argument
	filename := pflag.Arg(0)

	config := config.Load(filename, stage, outputPath, opt, debug)

	if config.OutputPath != "" {
		outputFile, err := os.Create(config.OutputPath)
		if err != nil {
			log.Fatalf("Failed to open output file: %s", err)
		}
		defer outputFile.Close()

		output.SetWriter(outputFile)
	}

	output.Writef("Config: %+v\n", config)

	// TODO
}
