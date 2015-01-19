package main

import (
	"fmt"
	"os"

	"github.com/NeowayLabs/clinit-cfn-tool/utils"
	"github.com/jteeuwen/go-pkg-optarg"
	"github.com/tiago4orion/datagen/data"
)

func main() {
	var dataConfigPath, outputFile, format string
	var helpOpt, missingOpts bool

	optarg.Add("h", "help", "Displays this help", false)
	optarg.Add("c", "data-config", "Data configuration file", "")
	optarg.Add("o", "output-file", "Output file", "")
	optarg.Add("f", "format", "Output format", "")

	for opt := range optarg.Parse() {
		switch opt.ShortName {
		case "c":
			dataConfigPath = opt.String()
		case "o":
			outputFile = opt.String()
		case "f":
			format = opt.String()
		case "h":
			helpOpt = opt.Bool()

		default:
			fmt.Println("Invalid flag: ", opt)
			optarg.Usage()
			os.Exit(1)
		}
	}

	if helpOpt {
		optarg.Usage()
		os.Exit(0)
	}

	if dataConfigPath == "" {
		fmt.Println("-c is required...")
		missingOpts = true
	}

	if missingOpts {
		optarg.Usage()
		os.Exit(1)
	}

	if os.Getenv("DEBUG_OPTS") != "" {
		fmt.Println("Data config: ", dataConfigPath)
		fmt.Println("Output file: ", outputFile)
		fmt.Println("Format: ", format)
	}

	err := data.Generator(dataConfigPath, outputFile, format, 100)
	utils.Check(err)
}
