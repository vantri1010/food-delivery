package main

import (
	"context"
	"food-delivery/pubsub"
	"food-delivery/pubsub/localpb"
	"log"
	"time"
)

func main() {
	var localPb pubsub.Pubsub = localpb.NewPubsub()
	topic := pubsub.Topic("OrderCreated")

	sub1, close1 := localPb.Subscribe(context.Background(), topic)
	sub2, close2 := localPb.Subscribe(context.Background(), topic)

	localPb.Publish(context.Background(), topic, pubsub.NewMessage(1))
	localPb.Publish(context.Background(), topic, pubsub.NewMessage(2))

	go func() {
		for {
			log.Println("Sub 1:", (<-sub1).Data())
			time.Sleep(time.Millisecond * 400)
		}
	}()

	go func() {
		for {
			log.Println("Sub 2:", (<-sub2).Data())
			time.Sleep(time.Millisecond * 400)
		}
	}()

	time.Sleep(time.Second * 7)
	close1()
	close2()
}
