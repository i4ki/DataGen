package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/NeowayLabs/clinit-cfn-tool/utils"
	"github.com/jteeuwen/go-pkg-optarg"
	"github.com/tiago4orion/DataGen/data"
)

func main() {
	var dataConfigPath, outputFile, format string
	var helpOpt, missingOpts bool
	var nrecords int32

	optarg.Add("h", "help", "Displays this help", false)
	optarg.Add("c", "data-config", "Data configuration file", "")
	optarg.Add("o", "output-file", "Output file", "")
	optarg.Add("f", "format", "Output format", "")
	optarg.Add("n", "number-records", "Number of records", "")

	for opt := range optarg.Parse() {
		switch opt.ShortName {
		case "c":
			dataConfigPath = opt.String()
		case "o":
			outputFile = opt.String()
		case "f":
			format = opt.String()
		case "n":
			ntmp, err := strconv.Atoi(opt.String())
			if err != nil {
				fmt.Println("Invalid -n value: ", opt.String())
			} else {
				nrecords = int32(ntmp)
			}
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

	if nrecords == 0 {
		fmt.Println("-n is required...")
		missingOpts = true
	}

	if missingOpts {
		optarg.Usage()
		os.Exit(1)
	}

	err := data.Generator(dataConfigPath, outputFile, format, nrecords)
	utils.Check(err)
}
