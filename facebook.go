package main

import "fmt"
import "encoding/json"
import "io/ioutil"
import "net/http"

type FacebookClient struct {
	accessToken   string
	feedUrl       string
	attachmentUrl string
}

type AttachmentResponse struct {
	Data []struct {
		Description string
		Type        string
		Url         string
		Media       struct {
			Image struct {
				Src string
			}
		}
		Target struct {
			Id string
		}
	}
}

type FeedResponse struct {
	Data []struct {
		Message string
		Id      string
	}
}

type ImageInfo struct {
	Url string
	Id  string
}

// Hits the /attachment url and gets the image url from the post
func (fc *FacebookClient) GetImageUrlsFromPostId(postId string) []ImageInfo {
	// TODO: Post ID is always the last post id when this is called
	fmt.Println("Getting images for post id: " + postId)
	url := fmt.Sprintf(fc.attachmentUrl, postId) + "?access_token=" + fc.accessToken
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	// Unmarshall response
	att := new(AttachmentResponse)
	json.Unmarshal(body, &att)
	ret := make([]ImageInfo, 0)
	for _, data := range att.Data {
		info := ImageInfo{}
		info.Url = data.Media.Image.Src
		info.Id = data.Target.Id
		ret = append(ret, info)
	}
	return ret
}

func (fc *FacebookClient) GetFeed(groupId int) *FeedResponse {
	url := fmt.Sprintf(fc.feedUrl, groupId) + "?access_token=" + fc.accessToken
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	// Unmarshall response into a slice of FeedItems
	feed := new(FeedResponse)
	json.Unmarshal(body, &feed)
	return feed
}

func NewFacebookClient(accessToken string) *FacebookClient {
	client := new(FacebookClient)

	// Set tokens
	client.accessToken = accessToken

	// Set endpoint urls
	client.feedUrl = "https://graph.facebook.com/v2.11/%d/feed"
	client.attachmentUrl = "https://graph.facebook.com/v2.11/%s/attachments"

	return client
}
