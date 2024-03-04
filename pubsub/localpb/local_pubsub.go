package localpb

import (
	"context"
	"food-delivery/pubsub"
	"log"
	"sync"
)

// A pb run locally (in-mem)
// It has a queue (buffer channel) at it's core and many group of subscribers.
// Because we want to send a message with a specific topic for many subscribers in a group can handle.

type localPubSub struct {
	messageQueue chan *pubsub.Message
	mapChannel   map[pubsub.Topic][]chan *pubsub.Message // one topic has many subscribers
	locker       *sync.RWMutex                           // use when concurrent
}

func NewPubsub() *localPubSub {
	pb := &localPubSub{
		messageQueue: make(chan *pubsub.Message, 10000),
		mapChannel:   make(map[pubsub.Topic][]chan *pubsub.Message),
		locker:       new(sync.RWMutex),
	}

	pb.run()

	return pb
}

// Publish publishes a message to a specific topic.
// It is a fire-and-forget operation, meaning that it does not wait for the message to be delivered to all subscribers.
// Instead, it immediately dispatches the message to the message queue and returns.
// The message is delivered to the subscribers in a separate goroutine.
func (ps *localPubSub) Publish(ctx context.Context, channel pubsub.Topic, data *pubsub.Message) error {
	data.SetChannel(channel)

	go func() {
		ps.messageQueue <- data
		log.Println("New event published: ", data.String(), "with data: ", data.Data())
	}()
	return nil
}

// Subscribe subscribes to a topic and returns a channel that receives messages on that topic.
// The returned function can be called to unsubscribe from the topic.
//
// The context is used to allow cancellation of the subscription. If the context is cancelled,
// the returned function will unsubscribe and close the channel.
//
// The function is thread-safe and can be called from multiple goroutines.
func (ps *localPubSub) Subscribe(ctx context.Context, topic pubsub.Topic) (ch <-chan *pubsub.Message, close func()) {
	c := make(chan *pubsub.Message)

	ps.locker.Lock()

	if val, ok := ps.mapChannel[topic]; ok {
		val = append(ps.mapChannel[topic], c)
		ps.mapChannel[topic] = val
	} else {
		ps.mapChannel[topic] = []chan *pubsub.Message{c}
	}

	ps.locker.Unlock()

	return c, func() {
		log.Println("Unsubscribe")

		if chans, ok := ps.mapChannel[topic]; ok {
			for i := range chans {
				if chans[i] == c {
					// remove element at index i (of c) in slice  chans
					chans = append(chans[:i], chans[i+1:]...)

					ps.locker.Lock()
					ps.mapChannel[topic] = chans
					ps.locker.Unlock()
					break
				}
			}
		}
	}

}

func (ps *localPubSub) run() error {
	log.Println("Pubsub started")

	go func() {
		for {
			mess := <-ps.messageQueue
			log.Println("Message dequeue", mess)

			if subs, ok := ps.mapChannel[mess.Channel()]; ok {
				for i := range subs {
					go func(c chan *pubsub.Message) {
						c <- mess
					}(subs[i])
				}
			}
			//else {
			//	ps.messageQueue <- mess
			//}
		}
	}()

	return nil
}
