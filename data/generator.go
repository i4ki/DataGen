package data

import (
	"errors"
	"fmt"
	"math"
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

	fmt.Println("Feeding data to process!")
	for i = 0; i < dataConfig.Length; i++ {
		fmt.Println("feeding data: ", i)
		recordChan <- dataConfig.Records
		time.Sleep(time.Millisecond)
	}

	close(recordChan)

	fmt.Println("All records scheduled to be created. Waiting workers...")
}

func pushRecords(workerIdx int, nrecords int32, config *DataConfig, outputChan chan string, wg *sync.WaitGroup) chan float64 {
	totalChan := make(chan float64, 100)

	go func() {
		wg.Add(1)
		for i := int32(0); i < nrecords; i++ {
			outLine := CSVLineCreate(config.Records)
			outputChan <- outLine
			totalChan <- 100.0 * float64(i+1) / float64(nrecords)
			time.Sleep(time.Millisecond)
		}

		defer wg.Done()
	}()

	return totalChan
}

func outputData(workerIdx int, outputChan chan string, config *DataConfig, fileOut *os.File, wg *sync.WaitGroup) {
	defer wg.Done()

	for outStr := range outputChan {
		if _, err := fileOut.Write([]byte(outStr)); err != nil {
			panic(err)
		}

		time.Sleep(time.Millisecond)
	}
}

func GenerateCsv(config *DataConfig) error {
	var wgRecords, wgOutput sync.WaitGroup
	outputChan := make(chan string)

	ncpu := runtime.NumCPU()
	runtime.GOMAXPROCS(ncpu)

	recordsPerCore := config.Length / int32(ncpu)

	files := make([]*os.File, ncpu)
	totalPWorkers := make([]chan float64, ncpu)
	for i := 0; i < ncpu; i++ {
		if i == (ncpu - 1) {
			recordsPerCore += int32(math.Remainder(float64(config.Length), float64(ncpu)))
		}

		fmt.Println("Scheduling create of ", recordsPerCore)
		totalPWorkers[i] = pushRecords(i, recordsPerCore, config, outputChan, &wgRecords)

		file, err := os.Create(config.OutputFile + "_" + strconv.Itoa(i) + ".csv")
		files[i] = file
		utils.Check(err)
		wgOutput.Add(1)
		go outputData(i, outputChan, config, files[i], &wgOutput)
	}

	var wt1, wt2, wt3, wt4 float64

	for !(wt1 == 100 && wt2 == 100 && wt3 == 100 && wt4 == 100) {
		fmt.Printf("Workers status: %f%%, %f%%, %f%%, %f%%\r", wt1, wt2, wt3, wt4)
		select {
		case wt1 = <-totalPWorkers[0]:
		case wt2 = <-totalPWorkers[1]:
		case wt3 = <-totalPWorkers[2]:
		case wt4 = <-totalPWorkers[3]:
		}
	}

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
		err = GenerateCsv(&dataConfig)
	default:
		fmt.Println("Unknown format: " + format)
	}

	return err
}
