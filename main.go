package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

const SESSION_COUNT = 3

var times []time.Time

func setTimer() error {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return err
	}
	now := time.Now().In(jst)
	rand.Seed(now.UnixNano())
	i := rand.Intn(40)

	target := now.Add(time.Duration(i) * time.Second)
	times = append(times, target)

	return nil
}

func main() {
	for i := 0; i < SESSION_COUNT; i++ {
		if err := setTimer(); err != nil {
			log.Println(err.Error())
		}
	}

	wg := sync.WaitGroup{}
	for i, t := range times {
		si := i + 1
		target := t

		wg.Add(1)
		go func(tt time.Time, ti int) {
			fmt.Println(ti)
			for {
				jst, err := time.LoadLocation("Asia/Tokyo")
				if err != nil {
					log.Println(err.Error())
					wg.Done()
					break
				}
				now := time.Now().In(jst)

				if tt.Before(now) {
					fmt.Println("session:", ti, " time:", tt.Format("2006-01-02 15-04-05"))
					wg.Done()
					break
				}
			}
		}(target, si)

		wg.Wait()
	}
}
