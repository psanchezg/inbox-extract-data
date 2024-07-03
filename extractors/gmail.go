package extractors

import (
	"context"
	"log"
	"sort"

	"gitlab.com/hartsfield/gmailAPI"
	"gitlab.com/hartsfield/inboxer"
	"google.golang.org/api/gmail/v1"
)

// getByID gets emails individually by ID. This is necessary because this is
// how the gmail API is set up apparently [0][1] (but why?).
// [0] https://developers.google.com/gmail/api/v1/reference/users/messages/get
// [1] https://stackoverflow.com/questions/36365172/message-payload-is-always-null-for-all-messages-how-do-i-get-this-data
func getByID(srv *gmail.Service, msgs *gmail.ListMessagesResponse) ([]*gmail.Message, error) {
	var msgSlice []*gmail.Message
	for _, v := range msgs.Messages {
		msg, err := srv.Users.Messages.Get("me", v.Id).Do()
		if err != nil {
			return msgSlice, err
		}
		msgSlice = append(msgSlice, msg)
	}
	return msgSlice, nil
}

func ExtractMails(query string) ([]*gmail.Message, error) {
	// Connect to the gmail API service.
	ctx := context.Background()
	srv := gmailAPI.ConnectToService(ctx, gmail.GmailReadonlyScope)
	// Recuperar todos los mensajes de la consulta
	var msgs []*gmail.Message
	nextPageToken := ""
	for {
		inbox, err := srv.Users.Messages.List("me").Q(query).PageToken(nextPageToken).Do()
		if err != nil {
			log.Fatalf("Error al obtener mensajes: %v", err)
			return nil, err
		}
		messagesList, err := getByID(srv, inbox)
		if err != nil {
			return msgs, err
		}

		msgs = append(msgs, messagesList...)
		nextPageToken = inbox.NextPageToken
		if nextPageToken == "" {
			break
		}
	}
	sort.Slice(msgs, func(a, b int) bool {
		timea, err := inboxer.ReceivedTime(msgs[a].InternalDate)
		if err != nil {
			return true
		}
		timeb, err := inboxer.ReceivedTime(msgs[b].InternalDate)
		if err != nil {
			return false
		}
		return timea.Before(timeb)
	})
	return msgs, nil
}
