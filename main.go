package main

import (
    "fmt"
	"log"
    "os"

	"github.com/spf13/pflag"

	"github.com/cosmicgenius/decaf-go/config"
)

func main() {
    var (
		stage = pflag.String("target", "", "Target stage (scan, parse, inter, assemble, or empty to run all)")
		output = pflag.String("output", "", "Output file")
		opt = pflag.String("opt", "", "Optimizations")
		debug = pflag.Bool("debug", false, "Debug mode")
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

    fmt.Printf("Processing file: %s\n", filename)

	config := config.Load(filename, stage, output, opt, debug)

	log.Printf("Config: %+v", config)

	// TODO
}
