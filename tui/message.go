package tui

import "github.com/flipez/rcterm/ws"

import (
	"fmt"
	"time"
)

func CreateMessageItem(message ws.Message) MessageListItem {
	messageTime := time.Unix(0, int64(message.Date.Timestamp)*int64(time.Millisecond))
	title := fmt.Sprintf("%s @ %s", message.Sender.Name, messageTime.Format("3:04 PM, 02 Jan 2006"))

	return NewMessageListItem(title, styleMessageBody(message.Message))
}

func styleMessageTitle()

func styleMessageBody(body string) string {
	// TODO: Implement Markdown here
	//r, _ := glamour.NewTermRenderer(glamour.WithAutoStyle(), glamour.WithPreservedNewLines(), glamour.WithWordWrap(40), glamour.WithEmoji())

	//markdown, err := r.Render(body)
	//if err != nil {
	//	return body
	//}

	//return markdown

	return body
}
