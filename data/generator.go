package data

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/NeowayLabs/clinit-cfn-tool/utils"
)

var CONCURRENT_ROUTINES int = 3

type RecordConfig struct {
	Name  string
	Type  string
	Chars string
	Min   int
	Max   int
}

type DataConfig struct {
	Records    []RecordConfig
	Length     int32
	OutputFile string
}

func CSVLineCreate(record []RecordConfig) string {
	return "bleh; blih; bloh; bluh\n"
}

func feedWorkers(recordChan chan []RecordConfig, dataConfig DataConfig) {
	var i int32
	for i = 0; i < dataConfig.Length; i++ {
		recordChan <- dataConfig.Records
	}

	close(recordChan)

	fmt.Println("All records scheduled to be created. Waiting workers...")
}

func pushRecords(workerIdx int, recordChan chan []RecordConfig, outputChan chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for record := range recordChan {
		outLine := CSVLineCreate(record)
		outputChan <- outLine
		time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
	}

	fmt.Printf("Worker %d finished executing.\n", workerIdx)
}

func outputData(workerIdx int, outputChan chan string, config DataConfig, fileOut *os.File, wg *sync.WaitGroup) {
	defer wg.Done()

	for outStr := range outputChan {
		if _, err := fileOut.Write([]byte(outStr)); err != nil {
			panic(err)
		}

		time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
	}

	fmt.Printf("Output worker %d finished execution.\n", workerIdx)
}

func GenerateCsv(config DataConfig) error {
	var wgRecords, wgOutput sync.WaitGroup
	recordChan := make(chan []RecordConfig)
	outputChan := make(chan string)

	ncpu := runtime.NumCPU()
	runtime.GOMAXPROCS(ncpu)

	go feedWorkers(recordChan, config)

	files := make([]*os.File, ncpu)
	for i := 0; i < ncpu; i++ {
		wgRecords.Add(1)
		go pushRecords(i, recordChan, outputChan, &wgRecords)

		file, err := os.Create(config.OutputFile + "_" + strconv.Itoa(i) + ".csv")
		files[i] = file
		utils.Check(err)
		wgOutput.Add(1)
		go outputData(i, outputChan, config, files[i], &wgOutput)
	}

	wgRecords.Wait()
	close(outputChan)
	wgOutput.Wait()

	// close fo on exit and check for its returned error
	defer func() {
		for _, fileOut := range files {
			if err := fileOut.Close(); err != nil {
				panic(err)
			}
		}
	}()

	return nil
}

func Generator(configFile string, outputFile string, format string, length int32) error {
	var dataConfig DataConfig

	if length == 0 {
		return errors.New("Number of records need be greater than zero.")
	}

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
	recordConfig := make([]RecordConfig, len(fields))

	for idx, field := range fields {
		for name, config := range field.(map[interface{}]interface{}) {
			cfg := config.(map[interface{}]interface{})
			rConfig := RecordConfig{
				Name:  name.(string),
				Type:  cfg["type"].(string),
				Chars: cfg["chars"].(string),
				Min:   cfg["min"].(int),
				Max:   cfg["max"].(int),
			}

			recordConfig[idx] = rConfig
		}
	}

	dataConfig.Records = recordConfig
	dataConfig.Length = length
	dataConfig.OutputFile = outputFile

	switch format {
	case "csv":
		err = GenerateCsv(dataConfig)
	default:
		fmt.Println("Unknown format: " + format)
	}

	return err
}
