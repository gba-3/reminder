package main

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gba-3/reminder/models"
	"github.com/gba-3/reminder/notify"
	"github.com/pkg/errors"
)

const (
	LOCALE = "Asia/Tokyo"
)

func init() {
	lcl, err := time.LoadLocation(LOCALE)
	if err != nil {
		log.Fatal(errors.WithStack(err))
	}
	time.Local = lcl
}

func sendReminder(ctx context.Context, task models.Task, chReminder chan<- models.Task) {
	pd, err := task.PublicDate()
	if err != nil {
		log.Println(errors.WithStack(err))
		return
	}

	for {
		now := time.Now()
		if pd.Before(now) {
			chReminder <- task
			break
		}
	}
}

func remind(ctx context.Context, wg *sync.WaitGroup, notify notify.Notify, chReminder <-chan models.Task) {
	for {
		select {
		case task := <-chReminder:
			notify.SendMessage(task.Name)
		case <-ctx.Done():
			wg.Done()
		default:
		}
	}
}

func main() {
	ctx := context.Background()

	url := os.Getenv("SLACK_WEBHOOK_URL")
	if url == "" {
		log.Fatalln("unexpected token: SLACK_WEBHOOK_URL is empty.")
	}

	sw, err := notify.NewSlackWebhook(url)
	if err != nil {
		log.Fatalln(errors.WithStack(err))
	}

	tasks := []models.Task{
		{
			Name:   "task1",
			Date:   "2022-05-14 14:48:24",
			Status: false,
		},
		{
			Name:   "task2",
			Date:   "2022-05-14 14:09:24",
			Status: false,
		},
	}

	wg := &sync.WaitGroup{}
	chReminder := make(chan models.Task)
	for _, task := range tasks {
		go sendReminder(ctx, task, chReminder)
	}

	wg.Add(1)
	go remind(ctx, wg, sw, chReminder)
	wg.Wait()
}
