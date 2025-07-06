package config

import (
	"log"
	"strings"
)

type TargetStage int

const (
	TargetStageAll TargetStage = iota
	TargetStageScan
	TargetStageParse
	TargetStageInter
	TargetStageAssembly
)

var stageNameToTargetStage = map[string]TargetStage{
	"":         TargetStageAll,
	"scan":     TargetStageScan,
	"parse":    TargetStageParse,
	"inter":    TargetStageInter,
	"assembly": TargetStageAssembly,
}

type Config struct {
	InputPath     string
	Stage         TargetStage
	OutputPath    string
	Optimizations []string
	Debug         bool
}

func Load(
	inputPath string,
	stagePtr *string,
	outputPathPtr *string,
	optPtr *string,
	debugPtr *bool,
) *Config {
	config := &Config{
		InputPath: inputPath,
	}

	if stagePtr != nil {
		stage, ok := stageNameToTargetStage[*stagePtr]
		if !ok {
			log.Fatalf("invalid target stage %s", *stagePtr)
		}
		config.Stage = stage
	}

	if outputPathPtr != nil {
		config.OutputPath = *outputPathPtr
	}

	if optPtr != nil {
		config.Optimizations = strings.Split(*optPtr, ",")
	}

	if debugPtr != nil {
		config.Debug = *debugPtr
	}

	return config
}
