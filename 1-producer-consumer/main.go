//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer szenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"sync"
	"time"
)

func producer(stream Stream, wg *sync.WaitGroup, ch chan *Tweet) (tweets []*Tweet) {
	defer wg.Done()
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			close(ch)
			return tweets
		}
		ch <- tweet

		tweets = append(tweets, tweet)
	}
}

func consumer(wg *sync.WaitGroup, ch chan *Tweet) {
	defer wg.Done()
	for{
		tweet, err := <- ch
		if(!err){
			return
		}

		if tweet.IsTalkingAboutGo() {
			fmt.Println(tweet.Username, "\ttweets about golang")
		} else {
			fmt.Println(tweet.Username, "\tdoes not tweet about golang")
		}

	}
}

func main() {
	start := time.Now()
	var waitGroup sync.WaitGroup
	channel := make(chan *Tweet)

	waitGroup.Add(2)
	stream := GetMockStream()

	// Producer
	go producer(stream, &waitGroup, channel)

	// Consumer
	go consumer(&waitGroup, channel)
	waitGroup.Wait()

	fmt.Printf("Process took %s\n", time.Since(start))
}
