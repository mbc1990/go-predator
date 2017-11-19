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
	FacebookClient *FacebookClient
	PostgresClient *PostgresClient
	Conf           *Configuration
	Wg             *sync.WaitGroup
	ExistingImages map[string]bool
}

// Downloads an image and writes its metadata to postgres
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
	fmt.Println("Processing timeline")
	defer p.Wg.Done()
	res := p.TwitterClient.GetTweets(handle)
	for _, tweet := range res {
		// If we've processed this tweet already, continue
		if _, ok := p.ExistingImages[tweet.Id_str]; ok {
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

func (p *Predator) ProcessFacebookPage(groupId int) {
	fmt.Println("Processing facebook page")
	defer p.Wg.Done()
	res := p.FacebookClient.GetFeed(groupId)
	for _, item := range res.Data {
		p.Wg.Add(1)
		go func(id string) {
			defer p.Wg.Done()
			imageInfo := p.FacebookClient.GetImageUrlsFromPostId(id)

			for _, info := range imageInfo {
				p.Wg.Add(1)
				go p.HandleImage(info.Url, "facebook", info.Id)
			}
		}(item.Id)
	}
}

// Entry point for a single run across all image sources
func (p *Predator) Run() {
	// TODO: Uncomment
	/*
		for _, handle := range p.Conf.TwitterSources {
			p.Wg.Add(1)
			go p.ProcessTwitterTimeline(handle)
		}
	*/
	for _, groupId := range p.Conf.FacebookSources {
		p.Wg.Add(1)
		go p.ProcessFacebookPage(groupId)
	}
	p.Wg.Wait()
}

func NewPredator(conf *Configuration) *Predator {
	p := new(Predator)
	p.Conf = conf

	// Twitter
	p.TwitterClient = NewTwitterClient(p.Conf.TwitterConsumerKey,
		p.Conf.TwitterConsumerSecret)

	// Facebook
	p.FacebookClient = NewFacebookClient(p.Conf.FacebookAccessToken)

	// Postgres
	p.PostgresClient = NewPostgresClient(p.Conf.PGHost, p.Conf.PGPort,
		p.Conf.PGUser, p.Conf.PGPassword, p.Conf.PGDbname)

	p.ExistingImages = p.PostgresClient.GetExistingImages()
	var wg sync.WaitGroup
	p.Wg = &wg
	return p
}
