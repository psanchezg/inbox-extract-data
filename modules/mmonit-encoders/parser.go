package mmonitencoders

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"

	"github.com/psanchezg/inbox-extract-data/utils"
	"google.golang.org/api/gmail/v1"
)

/**
 * Parses text with the given regular expression and returns the
 * group values defined in the expression.
 *
 */
func getParams(regEx, text string) (paramsMap []map[string]string) {

	var compRegEx = regexp.MustCompile(regEx)
	matches := compRegEx.FindAllStringSubmatch(text, -1)
	// fmt.Printf("%v", text)
	// fmt.Println(matches)

	paramsMap = make([]map[string]string, 0)
	index := 0
	for _, match := range matches {
		paramsMap = append(paramsMap, make(map[string]string))
		for i, name := range compRegEx.SubexpNames() {
			if name != "" && i >= 0 && i <= len(match) {
				paramsMap[index][name] = match[i]
			}
		}
		aux := strings.Split(paramsMap[index]["Channel"], "_")
		if len(aux) > 2 {
			if index > 0 {
				paramsMap = append(paramsMap[:index], paramsMap[index+1:]...)
			}
		} else {
			index++
		}
	}
	return paramsMap
}

func createMmonitAlert(params map[string]string, raw string, body string) MmonitAlert {
	alert := MmonitAlert{
		Snippet: raw,
	}
	// Convert date
	l := "02 Jan 15:04:05 -0700"
	tt, err := time.Parse(l, params["Date"])
	if err == nil {
		alert.Date = tt
	} else {
		fmt.Println("error", params["Date"], err, raw)
	}
	alert.Action = params["Action"]
	alert.Channel = params["Channel"]
	alert.Encoder = params["Encoder"]
	alert.Snippet = raw
	alert.Body = body
	return alert
}

func getSentDate(msg *gmail.Message) time.Time {
	var dateHeader string
	for _, header := range msg.Payload.Headers {
		if header.Name == "Date" {
			dateHeader = header.Value
			break
		}
	}
	t, _ := time.Parse(time.RFC3339, dateHeader)
	return t
}

// ProcessRawData process raw data
func ProcessRawData(msgs []*gmail.Message, currentYear int) (map[string][]MmonitUsage, error) {
	var firstMessage time.Time
	var lastMessage time.Time
	alerts := []MmonitAlert{}
	usages := map[string][]MmonitUsage{}
	openedUsages := map[string]MmonitUsage{}

	rxbody := `Date:\s+(?P<Date>\d{1,2} (Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec) \d{2}:\d{2}:\d{2} \+0000)\nHost:\s+(?P<Encoder>[a-z0-9\.]+)\nService:\s+(?P<Channel>[a-z0-9_]+)\sAction:\s+Alert\nDescription:\s+(?P<Action>(stop|start){1}) action done`

	// Range over the messages
	for _, msg := range msgs {
		// utils.WriteFile(msg)
		body, errbody := utils.GetMsgBody(msg)
		if errbody == nil {
			fullparams := getParams(rxbody, body)
			// fmt.Println(fullparams)
			for _, params := range fullparams {
				if params["Date"] != "" {
					alert := createMmonitAlert(params, msg.Snippet, body)
					sendDate := getSentDate(msg)
					localLocation, err := time.LoadLocation("Local")
					if err != nil {
						fmt.Println("Error al cargar la zona horaria local:", err)
					}
					// Convertir la fecha UTC a la zona horaria local
					sendDate = sendDate.In(localLocation)
					newDate := time.Date(currentYear, alert.Date.Month(), alert.Date.Day(), alert.Date.Hour(), alert.Date.Minute(), alert.Date.Second(), alert.Date.Nanosecond(), sendDate.Location())
					alert.Date = newDate
					alerts = append(alerts, alert)
					if firstMessage.IsZero() || firstMessage.After(newDate) {
						firstMessage = newDate
					}
					lastMessage = newDate
					// if errbody == nil {
					// }
					// fmt.Println(params["Date"], params["Encoder"], params["Channel"], params["Action"])
					if !alert.Date.IsZero() {
						// fmt.Println(alert.Date.Format("02-01-2006 15:04:05 MST"), alert.Encoder, alert.Channel, alert.Action)
						key := fmt.Sprintf("%s-%s", alert.Encoder, alert.Channel)
						if alert.Action == "start" {
							if currentUsage, exists := openedUsages[key]; exists {
								fmt.Println("*****************ERROR START****************", key, currentUsage)
							} else {
								openedUsages[key] = MmonitUsage{
									Start:    alert.Date,
									Stop:     alert.Date,
									Encoder:  alert.Encoder,
									Channel:  alert.Channel,
									Duration: 0,
									Minutes:  0,
								}
							}
						} else if alert.Action == "stop" {
							if currentUsage, exists := openedUsages[key]; exists {
								aux := strings.Split(alert.Channel, "_")
								simplekey := aux[0] + "_" + aux[1]
								currentUsage.Stop = alert.Date
								currentUsage.Channel = simplekey
								currentUsage.Duration = float64(currentUsage.Stop.Sub(currentUsage.Start).Seconds())
								currentUsage.Minutes = float64(math.Round(currentUsage.Duration*100/60.0) / 100)
								openedUsages[key] = currentUsage
								usages[simplekey] = append(usages[simplekey], currentUsage)
								delete(openedUsages, key)
							} else {
								fmt.Println("*****************ERROR STOP****************", key, "start not found", alert.Date.Format("02/01/2006 15:04:05"))
							}
						}
					}
				}
			}
		}
	}
	if !firstMessage.IsZero() {
		fmt.Printf("Primera alerta detectada: %v\n", firstMessage.Format("02-01-2006 15:04:05 MST"))
	}
	if !lastMessage.IsZero() {
		fmt.Printf("Última alerta detectada: %v\n", lastMessage.Format("02-01-2006 15:04:05 MST"))
	}
	fmt.Printf("Total número de alertas procesadas: %v\n", len(alerts))
	return usages, nil
}
