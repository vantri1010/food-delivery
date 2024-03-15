package main

import (
	"context"
	"errors"
	"food-delivery/component/asyncjob"
	"log"
	"time"
)

func main() {
	job1 := asyncjob.NewJob(func(ctx context.Context) error {
		time.Sleep(time.Second)
		log.Println("I am job 1")

		//return nil
		return errors.New("something went wrong at job 1")
	})

	//if err := job1.Execute(context.Background()); err != nil {
	//	log.Println("Job  is : ", job1.State(), err)
	//
	//	for {
	//		if err := job1.Retry(context.Background()); err != nil {
	//			log.Println("Retry with:", job1.RetryIndex(), "with state: ", job1.State(), err)
	//		}
	//
	//		if job1.State() == asyncjob.StateRetryFailed {
	//			break
	//		}
	//
	//		if job1.State() == asyncjob.StateCompleted {
	//			break
	//		}
	//	}
	//}

	job2 := asyncjob.NewJob(func(ctx context.Context) error {
		time.Sleep(time.Second * 2)
		log.Println("I am job 2")

		return nil
	})

	job3 := asyncjob.NewJob(func(ctx context.Context) error {
		time.Sleep(time.Second * 3)
		log.Println("I am job 3")

		return nil
	})

	group := asyncjob.NewGroup(true, job1, job2, job3)

	if err := group.Run(context.Background()); err != nil {
		log.Println(err)
	}
}
