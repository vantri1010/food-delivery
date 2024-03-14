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
	buffChanMsg chan *pubsub.Message
	mapChannel  map[pubsub.Topic][]chan *pubsub.Message // one topic has many subscribers
	locker      *sync.RWMutex                           // use when subscribe (non-concurrent)
}

func NewPubsub() *localPubSub {
	pb := &localPubSub{
		buffChanMsg: make(chan *pubsub.Message, 10000),
		mapChannel:  make(map[pubsub.Topic][]chan *pubsub.Message),
		locker:      new(sync.RWMutex),
	}

	pb.run()

	return pb
}

// Publish publishes a message to a specific topic.
// It is a fire-and-forget operation, meaning that it does not wait for the message to be delivered to all subscribers.
// Instead, it immediately dispatches the message to the message queue and returns.
// The message is delivered to the subscribers in a separate goroutine.
func (ps *localPubSub) Publish(ctx context.Context, topic pubsub.Topic, data *pubsub.Message) error {
	data.SetTopic(topic)

	go func() {
		log.Println("New event published: ", data.String(), "with data: ", data.Data())
		ps.buffChanMsg <- data
	}()
	return nil
}

// Subscribe function subscribes to a topic by create a channel of message slot in a list channels identical by topic we call mapChannel
// Then returns a channel that receives messages on that topic.
// The returned channel and the channel appended to the map are identical
// The returned function can be called to unsubscribe from the topic.
func (ps *localPubSub) Subscribe(ctx context.Context, topic pubsub.Topic) (ch <-chan *pubsub.Message, close func()) {
	newsubcrb := make(chan *pubsub.Message)

	ps.locker.Lock()

	if chans, ok := ps.mapChannel[topic]; ok {
		chans = append(ps.mapChannel[topic], newsubcrb)
		ps.mapChannel[topic] = chans
	} else {
		ps.mapChannel[topic] = []chan *pubsub.Message{newsubcrb}
	}

	ps.locker.Unlock()

	return newsubcrb, func() {
		log.Println("Unsubscribe")

		if chans, ok := ps.mapChannel[topic]; ok {
			for i := range chans {
				if chans[i] == newsubcrb {
					// remove element at index i in slice chans
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

// The run function in the localPubSub struct is responsible for starting the pubsub service.
// It does this by creating a goroutine that listens for messages on the buffChanMsg channel
// and dispatches them to the appropriate subscribers that registered by Subscribe function.
func (ps *localPubSub) run() error {
	log.Println("Pubsub started")

	go func() {
		for {
			mess := <-ps.buffChanMsg // wait (blocking) until buffChanMsg has message
			log.Println("Message dequeue", mess.String(), "with data: ", mess.Data())

			if subs, ok := ps.mapChannel[mess.Topic()]; ok {
				// use goroutine to send messages to each subscriber concurrently
				// to avoid blocking if the buffer channel somehow broken
				for i := range subs {
					go func(subi chan *pubsub.Message) {
						subi <- mess
					}(subs[i])
				}
			}
			//else {
			//	ps.buffChanMsg <- mess
			//}
		}
	}()

	return nil
}
