package extractors

import (
	"context"
	"sort"

	"gitlab.com/hartsfield/gmailAPI"
	"gitlab.com/hartsfield/inboxer"
	"google.golang.org/api/gmail/v1"
)

func ExtractMails(query string) ([]*gmail.Message, error) {
	// Connect to the gmail API service.
	ctx := context.Background()
	srv := gmailAPI.ConnectToService(ctx, gmail.GmailReadonlyScope)
	msgs, err := inboxer.Query(srv, query)
	if err != nil {
		return nil, err
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
