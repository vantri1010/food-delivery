package subscriber

import (
	"context"
	"food-delivery/common"
	"food-delivery/component/appctx"
	"food-delivery/component/asyncjob"
	"food-delivery/pubsub"
	"log"
)

type consumerJob struct {
	Title string
	Hdl   func(ctx context.Context, message *pubsub.Message) error
}

type consumerEngine struct {
	appCtx appctx.AppContext
}

func NewEngine(appContext appctx.AppContext) *consumerEngine {
	return &consumerEngine{appCtx: appContext}
}

func (engine *consumerEngine) Start() error {
	engine.startSubTopic(
		common.TopicUserLikeRestaurant,
		true,
		IncreaseLikeCountAfterUserLikeRestaurant(engine.appCtx),
		PushNotiUserLikeRestaurant(engine.appCtx),
		EmitRealtimeAfterUserLikeRestaurant(engine.appCtx),
	)

	engine.startSubTopic(
		common.TopicUserDisLikeRestaurant,
		true,
		DecreaseLikeCountAfterUserDisLikeRestaurant(engine.appCtx),
	)

	return nil
}

type GroupJob interface {
	Run(ctx context.Context) error
}

func (engine *consumerEngine) startSubTopic(topic pubsub.Topic, isConcurrent bool, hdls ...consumerJob) error {
	// forever: listen message and execute group
	// hdls => []job
	// new group ([]job)

	c, _ := engine.appCtx.GetPubSub().Subscribe(context.Background(), topic)

	for _, item := range hdls {
		log.Println("Setup subscriber for:", item.Title)
	}

	getHld := func(job *consumerJob, message *pubsub.Message) asyncjob.JobHandler {
		return func(ctx context.Context) error {
			log.Println("running job for ", job.Title, ". Value: ", message.Data())
			return job.Hdl(ctx, message)
		}
	}

	go func() {
		for {
			msg := <-c

			jobHdlArr := make([]asyncjob.Job, len(hdls))

			for i := range hdls {

				// capture msg & hlds[i]
				//jobHdl := func(ctx context.Context) error {
				//	log.Println("running job for ", hdls[i].Title, ". Value: ", msg.Data())
				//	return hdls[i].Hdl(ctx, msg)
				//}

				jobHdl := getHld(&hdls[i], msg)
				jobHdlArr[i] = asyncjob.NewJob(jobHdl)
			}

			group := asyncjob.NewGroup(isConcurrent, jobHdlArr...)

			if err := group.Run(context.Background()); err != nil {
				log.Println(err)
			}
		}
	}()

	return nil
}
