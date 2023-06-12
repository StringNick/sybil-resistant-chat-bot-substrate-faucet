package entity

type Message struct {
	// ChannelID - channel id of message
	ChannelID string
	// FromID - from id of user that send message
	FromID string
	// Body - message body
	Body []byte
}
