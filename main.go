package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/psanchezg/inbox-extract-data/config"
	"github.com/psanchezg/inbox-extract-data/extractors"
	"github.com/psanchezg/inbox-extract-data/modules/bolt"
	"github.com/psanchezg/inbox-extract-data/outputs"
	"github.com/spf13/viper"
)

var (
	// destinationDir is the path to the directory where the attachments will be saved.
	afterDate = os.Getenv("AFTER_DATE")
)

// func writeFile(msg *gmail.Message) {
// 	time, err := inboxer.ReceivedTime(msg.InternalDate)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	f, err := os.Create(fmt.Sprintf("./dump/%s-%s.txt", time.Format("2006-02-01"), msg.Id))
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	decoded, err := base64.URLEncoding.DecodeString(msg.Payload.Body.Data)
// 	if err != nil {
// 		fmt.Println(err)
// 		return

// 	}
// 	if _, err := f.WriteString(string(decoded)); err != nil {
// 		fmt.Println(err)
// 		f.Close()
// 		return
// 	}
// 	err = f.Close()
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// }

func processMails(configurations config.Configurations) {
	// TODO: Iterate over all processes
	process := configurations.Processes[0]
	query := process.Query
	msgs, err := extractors.ExtractMails(query)
	if err != nil {
		fmt.Println(err)
	}

	if process.Module == "bolt" {
		fmt.Println("========================================================")
		fmt.Printf("PROCESSING... %v\n", process.Name)
		fmt.Println("========================================================")
		planes, err := bolt.ProcessRawData(msgs)
		if err != nil {
			fmt.Println(err)
			return
		}
		var serialized []map[string]interface{}
		inrec, _ := json.Marshal(planes)
		json.Unmarshal(inrec, &serialized)
		// Human export
		var ret []string
		if ret, err = bolt.ExportDataAsStrings[map[string]interface{}](serialized); err != nil {
			fmt.Println(err)
		}
		for _, output := range process.Outputs {
			if output.Type == "stdout" {
				outputs.ConsoleOutput(ret)
			} else if output.Type == "file" {
				outputs.FileOutput(ret, output.Path)
			}
		}
		// Machine export
		// TODO: CSV, JSON, XML, Sheets, etc.
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
