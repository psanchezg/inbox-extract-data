package utils

import (
	"encoding/base64"
	"fmt"
	"os"

	"gitlab.com/hartsfield/inboxer"
	"google.golang.org/api/gmail/v1"
)

func WriteFile(msg *gmail.Message) {
	time, err := inboxer.ReceivedTime(msg.InternalDate)
	if err != nil {
		fmt.Println(err)
	}
	f, err := os.Create(fmt.Sprintf("./dump/%s-%s.txt", time.Format("2006-02-01"), msg.Id))
	if err != nil {
		fmt.Println(err)
		return
	}
	decoded, err := base64.URLEncoding.DecodeString(msg.Payload.Body.Data)
	if err != nil {
		fmt.Println(err)
		return

	}
	if _, err := f.WriteString(string(decoded)); err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
