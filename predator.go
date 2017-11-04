package main

import "fmt"

type Predator struct {
	TwitterClient *TwitterClient
	Conf          *Configuration
}

// Entry point for a single run across all image sources
func (p *Predator) Run() {
	res := p.TwitterClient.GetTweets("corgsbot")
	for _, el := range res {
		fmt.Println(el.Text)
		fmt.Println("\n")
	}
}

func NewPredator(conf *Configuration) *Predator {
	p := new(Predator)
	p.Conf = conf
	p.TwitterClient = NewTwitterClient(p.Conf.TwitterConsumerKey, p.Conf.TwitterConsumerSecret)
	return p
}
