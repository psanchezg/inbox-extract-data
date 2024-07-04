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

func processData(configurations config.Configurations) {
	// TODO: Iterate over all processes
	for _, process := range configurations.Processes {
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

		// Human export
		var lines []string
		var values [][]interface{}

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
			if lines, values, err = bolt.ExportData[map[string]interface{}](serialized); err != nil {
				fmt.Println(err)
			}
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
			if lines, values, err = mmonitencoders.ExportData[map[string]interface{}](serialized); err != nil {
				fmt.Println(err)
			}
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

		fmt.Println("========================================================")
	}
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

	processData(configuration)
}
