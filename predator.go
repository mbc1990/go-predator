package main

import "fmt"
import "io"
import "log"
import "net/http"
import "os"
import "strings"
import "sync"

type Predator struct {
	TwitterClient *TwitterClient
	Conf          *Configuration
	Wg            *sync.WaitGroup
}

// Downloads and deduplicates an image
func (p *Predator) HandleImage(url string) {
	defer p.Wg.Done()
	fmt.Println("Downloading image from " + url)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Get the filename
	spl := strings.Split(url, "/")
	fname := spl[len(spl)-1]

	// Create the file and copy the response body into it
	file, err := os.Create(p.Conf.UnclassifiedWorkDir + fname)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
}

// Hits twitter api and downloads images
func (p *Predator) ProcessTwitterTimeline(handle string) {
	defer p.Wg.Done()
	res := p.TwitterClient.GetTweets(handle)
	for _, tweet := range res {
		medias := tweet.Entities.Media
		for _, media := range medias {
			url := media.Media_url
			// TODO: If URL in already queried, skip
			p.Wg.Add(1)
			go p.HandleImage(url)
		}
	}
}

// Entry point for a single run across all image sources
func (p *Predator) Run() {
	for _, handle := range p.Conf.TwitterSources {
		p.Wg.Add(1)
		go p.ProcessTwitterTimeline(handle)
	}
	p.Wg.Wait()
}

func NewPredator(conf *Configuration) *Predator {
	p := new(Predator)
	p.Conf = conf
	p.TwitterClient = NewTwitterClient(p.Conf.TwitterConsumerKey, p.Conf.TwitterConsumerSecret)
	var wg sync.WaitGroup
	p.Wg = &wg
	return p
}
