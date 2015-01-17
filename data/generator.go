package data

import (
	"errors"
	"fmt"

	"github.com/NeowayLabs/clinit-cfn-tool/utils"
)

type DataConfig struct {
	Name  string
	Chars string
	Min   int
	Max   int
}

func Generator(configFile string, outputFile string, format string) error {
	cfgContent := utils.ReadFile(configFile)
	cfgYaml, err := utils.DecodeYaml([]byte(cfgContent))

	utils.Check(err)

	if format == "" {
		if cfgYaml["format"].(string) != "" {
			format = cfgYaml["format"].(string)
			fmt.Printf("Output format: %s\n", format)
		} else {
			return errors.New("No output format chosen...")
		}
	}

	fields := cfgYaml["fields"].([]interface{})
	dataConfigs := make([]DataConfig, len(fields))

	for idx, field := range fields {
		for name, config := range field.(map[interface{}]interface{}) {
			cfg := config.(map[interface{}]interface{})
			dConfig := DataConfig{
				Name:  name.(string),
				Chars: cfg["chars"].(string),
				Min:   cfg["min"].(int),
				Max:   cfg["max"].(int),
			}

			fmt.Println(dConfig)

			dataConfigs[idx] = dConfig
		}
	}

	switch format {
	case "csv":
		err = GenerateCsv(dataConfigs)
	default:
		fmt.Println("Unknown format: " + format)
	}

	return err
}
