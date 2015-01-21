package data

import (
	"errors"
	"fmt"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/NeowayLabs/clinit-cfn-tool/utils"
	utilsg "github.com/tiago4orion/DataGen/utils"
)

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

type WorkerStatus struct {
	Id    int
	Total float64
}

func isWorkersComplete(workStats []float64) bool {
	var complete bool = true

	for i := 0; i < len(workStats); i++ {
		complete = complete && (workStats[i] == 100)
	}

	return complete
}

func CSVLineCreate(record []RecordConfig) string {
	fields := make([]string, len(record))
	var err error

	for idx, recordField := range record {
		switch recordField.Type {
		case "string":
			fields[idx], err = utilsg.GeneratorString(recordField.Chars,
				recordField.Min, recordField.Max)
			if err != nil {
				panic(err)
			}
		case "integer":
			tmpInt, err := utilsg.GeneratorInteger(recordField.Min, recordField.Max)
			if err != nil {
				panic(err)
			}

			fields[idx] = strconv.Itoa(tmpInt)
		}
	}

	return strings.Join(fields, ",") + "\n"
}

func pushRecords(workerIdx int, nrecords int32, config *DataConfig, outputChan chan string, workStatChan chan WorkerStatus, wg *sync.WaitGroup) {
	wg.Add(1)
	for i := int32(0); i < nrecords; i++ {
		outLine := CSVLineCreate(config.Records)
		outputChan <- outLine
		workStatChan <- WorkerStatus{
			Id:    workerIdx,
			Total: 100.0 * float64(i+1) / float64(nrecords),
		}

		time.Sleep(time.Millisecond)
	}

	defer wg.Done()
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

func workersResumeTotal(workStats []float64) float64 {
	total := float64(0)

	for i := 0; i < len(workStats); i++ {
		total += workStats[i]
	}

	return total / float64(len(workStats))
}

func GenerateCsv(config *DataConfig, concurrent int) error {
	var wgRecords, wgOutput sync.WaitGroup
	outputChan := make(chan string)
	workStatChan := make(chan WorkerStatus)
	ncpu := concurrent

	if concurrent == 0 {
		ncpu = runtime.NumCPU()
	}

	runtime.GOMAXPROCS(ncpu)

	recordsPerCore := config.Length / int32(ncpu)

	files := make([]*os.File, ncpu)
	for i := 0; i < ncpu; i++ {
		if i == (ncpu - 1) {
			recordsPerCore += int32(math.Remainder(float64(config.Length), float64(ncpu)))
		}

		fmt.Println("Scheduling create of ", recordsPerCore)
		go pushRecords(i, recordsPerCore, config, outputChan, workStatChan, &wgRecords)

		file, err := os.Create(config.OutputFile + "_" + strconv.Itoa(i) + ".csv")
		files[i] = file
		utils.Check(err)
		wgOutput.Add(1)
		go outputData(i, outputChan, config, files[i], &wgOutput)
	}

	workStats := make([]float64, ncpu)

	for !isWorkersComplete(workStats) {
		fmt.Printf("Workers status: %.2f%%                      \r", workersResumeTotal(workStats))
		select {
		case status := <-workStatChan:
			workStats[status.Id] = status.Total
		}
	}

	fmt.Printf("Workers status: %.2f%%                      \n", workersResumeTotal(workStats))

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

func Generator(configFile string, outputFile string, format string, length int32, concurrent int) error {
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
			chars, ok := cfg["chars"].(string)

			if !ok {
				chars = ""
			}

			rConfig := RecordConfig{
				Name:  name.(string),
				Type:  cfg["type"].(string),
				Chars: chars,
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
		err = GenerateCsv(&dataConfig, concurrent)
	default:
		fmt.Println("Unknown format: " + format)
	}

	return err
}
