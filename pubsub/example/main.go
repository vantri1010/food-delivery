package example

import (
	"context"
	"food-delivery/pubsub"
	"food-delivery/pubsub/localpb"
	"log"
	"time"
)

func main() {
	var localPb pubsub.Pubsub = localpb.NewPubsub()
	chn := pubsub.Topic("OrderCreated")

	sub1, close1 := localPb.Subscribe(context.Background(), chn)
	sub2, close2 := localPb.Subscribe(context.Background(), chn)

	localPb.Publish(context.Background(), chn, pubsub.NewMessage(1))
	localPb.Publish(context.Background(), chn, pubsub.NewMessage(2))

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
