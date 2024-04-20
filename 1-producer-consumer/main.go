//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer scenario, where a producer reads in
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

func producer(stream Stream, channel chan *Tweet) (tweets []*Tweet) {
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			channel <- nil
			return tweets
		}
		channel <- tweet
		tweets = append(tweets, tweet)
	}
}

func consumer(channel chan *Tweet, wg *sync.WaitGroup) {
	for {
		tweet := <-channel
		if tweet == nil {
			wg.Done()
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
	stream := GetMockStream()

	channel := make(chan *Tweet)
	wg := sync.WaitGroup{}
	wg.Add(1)
	// Producer
	go producer(stream, channel)

	// Consumer
	go consumer(channel, &wg)
	wg.Wait()
	fmt.Printf("Process took %s\n", time.Since(start))

}
