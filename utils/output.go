package utils

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"gitlab.com/hartsfield/inboxer"
	"google.golang.org/api/gmail/v1"
)

func WriteFile(msg *gmail.Message) error {
	time, err := inboxer.ReceivedTime(msg.InternalDate)
	if err != nil {
		fmt.Println(err)
	}
	f, err := os.Create(fmt.Sprintf("./dump/%s-%s.txt", time.Format("2006-02-01"), msg.Id))
	if err != nil {
		fmt.Println(err)
		return f.Close()
	}
	decoded, err := GetMsgBody(msg)
	if err != nil {
		fmt.Println(err)
		return f.Close()
	}
	if _, err := f.WriteString(decoded + "\n"); err != nil {
		fmt.Println(err)
		return f.Close()
	}
	if _, err := f.WriteString("=============================================\n"); err != nil {
		fmt.Println(err)
		return f.Close()
	}
	if _, err := f.WriteString(msg.Snippet + "\n"); err != nil {
		fmt.Println(err)
		return f.Close()
	}
	return f.Close()
}

func ParseDateWithFormat(date string) (time.Time, error) {
	if strings.HasSuffix(date, "Z") {
		date = date[:len(date)-1] + "+00:00"
	}
	layout := "2006-01-02T15:04:05-07:00"
	return time.Parse(layout, date)
}

func Round(value any) float64 {
	return math.Round(value.(float64)*100) / 100
}
