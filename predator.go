package main

import "fmt"

type Predator struct {
	TwitterClient *TwitterClient
	Conf          *Configuration
}

// Downloads an image
func (p *Predator) GetImage(url string) {

}

// Entry point for a single run across all image sources
func (p *Predator) Run() {
	for _, handle := range p.Conf.TwitterSources {
		res := p.TwitterClient.GetTweets(handle)
		for _, tweet := range res {
			medias := tweet.Entities.Media
			for _, media := range medias {
				url := media.Media_url
				// TODO: Get image
				fmt.Println(url)
			}
		}
	}
}

func NewPredator(conf *Configuration) *Predator {
	p := new(Predator)
	p.Conf = conf
	p.TwitterClient = NewTwitterClient(p.Conf.TwitterConsumerKey, p.Conf.TwitterConsumerSecret)
	return p
}
