package main

import "fmt"
import "io"
import "log"
import "net/http"
import "os"
import "strings"
import "sync"

type Predator struct {
	TwitterClient  *TwitterClient
	PostgresClient *PostgresClient
	Conf           *Configuration
	Wg             *sync.WaitGroup
}

// Downloads and deduplicates an image
func (p *Predator) HandleImage(url string, source string, sourceId string) {
	defer p.Wg.Done()
	fmt.Println("Downloading image from " + url)
	resp, err := http.Get(url)
	if err != nil {
		// Don't exit program here since we expect a baseline of HTTP errors
		log.Print(err)

		// No response body to close
		return
	}
	defer resp.Body.Close()

	// Get the filename
	spl := strings.Split(url, "/")
	fname := spl[len(spl)-1]

	// Create the file and copy the response body into it
	file, err := os.Create(p.Conf.UnclassifiedWorkDir + fname)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	p.PostgresClient.InsertImage(fname, url, source, sourceId)
}

// Hits twitter api and downloads images
func (p *Predator) ProcessTwitterTimeline(handle string) {
	defer p.Wg.Done()
	res := p.TwitterClient.GetTweets(handle)
	for _, tweet := range res {
		// If tweet has already been handled, skip
		if p.PostgresClient.ImageExists(tweet.Id_str) {
			continue
		}

		// Otherwise, grab the media URLs and process them
		medias := tweet.Entities.Media
		for _, media := range medias {
			url := media.Media_url
			p.Wg.Add(1)
			go p.HandleImage(url, "twitter", tweet.Id_str)
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
	p.TwitterClient = NewTwitterClient(p.Conf.TwitterConsumerKey,
		p.Conf.TwitterConsumerSecret)
	p.PostgresClient = NewPostgresClient(p.Conf.PGHost, p.Conf.PGPort,
		p.Conf.PGUser, p.Conf.PGPassword, p.Conf.PGDbname)
	var wg sync.WaitGroup
	p.Wg = &wg
	return p
}
