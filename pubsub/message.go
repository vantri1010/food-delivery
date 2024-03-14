package pubsub

import (
	"fmt"
	"time"
)

type Message struct {
	id        string
	topic     Topic
	data      interface{}
	createdAt time.Time
}

func NewMessage(data interface{}) *Message {
	now := time.Now().UTC()

	return &Message{
		id:        fmt.Sprintf("%d", now.UnixNano()),
		data:      data,
		createdAt: time.Now().UTC(),
	}
}

func (evt *Message) String() string {
	return fmt.Sprintf("Message %s", evt.topic)
}

func (evt *Message) Topic() Topic {
	return evt.topic
}

func (evt *Message) SetTopic(channel Topic) {
	evt.topic = channel
}

func (evt *Message) Data() interface{} {
	return evt.data
}
