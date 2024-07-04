package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/psanchezg/inbox-extract-data/config"
	"github.com/psanchezg/inbox-extract-data/extractors"
	"github.com/psanchezg/inbox-extract-data/modules/bolt"
	mmonitencoders "github.com/psanchezg/inbox-extract-data/modules/mmonit-encoders"
	"github.com/psanchezg/inbox-extract-data/outputs"
	"github.com/psanchezg/inbox-extract-data/utils"
	"github.com/spf13/viper"
)

var (
	// destinationDir is the path to the directory where the attachments will be saved.
	afterDate = os.Getenv("AFTER_DATE")
)

func processMails(configurations config.Configurations) {
	// TODO: Iterate over all processes
	process := configurations.Processes[0]
	query := process.Query
	currentYear := time.Now().Year()
	// Extract date
	rx := `after:(?P<After>\d{4}\/\d{2}\/\d{2})`
	params := utils.GetParams(rx, query)
	// Convert date
	l := "2006/01/02"
	tt, err := time.Parse(l, params["After"])
	if err == nil {
		currentYear = tt.Year()
	}
	msgs, err := extractors.ExtractMails(query)
	if err != nil {
		fmt.Println(err)
	}

	if process.Module == "bolt" {
		fmt.Println("========================================================")
		fmt.Printf("PROCESSING... %v\n", process.Name)
		fmt.Println("========================================================")
		planes, err := bolt.ProcessRawData(msgs, currentYear)
		if err != nil {
			fmt.Println(err)
			return
		}
		var serialized []map[string]interface{}
		inrec, _ := json.Marshal(planes)
		json.Unmarshal(inrec, &serialized)
		// Human export
		var lines []string
		var values [][]interface{}
		if lines, values, err = bolt.ExportData[map[string]interface{}](serialized); err != nil {
			fmt.Println(err)
		}
		for _, output := range process.Outputs {
			if output.Type == "stdout" {
				outputs.ConsoleOutput(lines)
			} else if output.Type == "file" {
				outputs.FileOutput(lines, output.Path)
			} else if output.Type == "sheet" {
				outputs.SheetsOutput(values, output.Path)
			}
		}
		// Machine export
		// TODO: CSV, JSON, XML, Sheets, etc.
	} else if process.Module == "mmonit-encoders" {
		fmt.Println("========================================================")
		fmt.Printf("PROCESSING... %v\n", process.Name)
		fmt.Println("========================================================")
		data, err := mmonitencoders.ProcessRawData(msgs, currentYear)
		if err != nil {
			fmt.Println(err)
			return
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		var serialized map[string][]map[string]interface{}
		inrec, _ := json.Marshal(data)
		json.Unmarshal(inrec, &serialized)
		// Human export

		var ret []string
		stats := map[string]mmonitencoders.MmonitUsageStats{}
		for key := range data {
			if ret2, err := mmonitencoders.ExportDataAsStrings[map[string]interface{}](serialized[key]); err != nil {
				fmt.Println(err)
			} else {
				aux := mmonitencoders.GetAggregateStats[map[string]interface{}](serialized[key])
				ret = append(ret, "========================================================\n")
				ret = append(ret, fmt.Sprintf("Total uso del canal %s: %v minutos\n", key, aux.Minutes))
				stats[key] = aux
				ret = append(ret, ret2...)
			}
		}
		for _, output := range process.Outputs {
			if output.Type == "stdout" {
				outputs.ConsoleOutput(ret)
			} else if output.Type == "file" {
				outputs.FileOutput(ret, output.Path)
			}
		}
	}

	fmt.Println("========================================================")
}

func main() {
	// Set the file name of the configurations file
	viper.SetConfigName("config")

	// Set the path to look for the configurations file
	viper.AddConfigPath(".")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType("yml")
	var configuration config.Configurations

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	// Set undefined variables
	// viper.SetDefault("database.dbname", "test_db")

	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	processMails(configuration)
}
