package main

import "encoding/base64"
import "encoding/json"
import "net/http"
import "io/ioutil"
import "fmt"
import "bytes"

// Make twitter api calls
type TwitterClient struct {
	consumerKey    string
	consumerSecret string
	bearerToken    string
	timelineUrl    string
	appAuthUrl     string
}

// Single tweet
type Tweet struct {
	user      string
	content   string
	timestamp int
	hashtags  []string
}

type AuthResponse struct {
	Token_type   string
	Access_token string
}

// Returns a slice of tweets by the requested user
func (tc *TwitterClient) GetTweets(user string) []Tweet {
	// url := tc.timelineUrl + "?screen_name=" + user
	res := make([]Tweet, 1)
	return res
}

// Gets a bearer token for app-only authentication
func (tc *TwitterClient) AppOnlyAuth() {
	toEncode := []byte(tc.consumerKey + ":" + tc.consumerSecret)
	toSend := base64.StdEncoding.EncodeToString(toEncode)
	client := &http.Client{}

	bodyToSend := bytes.NewBuffer([]byte("grant_type=client_credentials"))
	req, _ := http.NewRequest("POST", tc.appAuthUrl, bodyToSend)
	req.Header.Add("Authorization", "Basic "+toSend)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8.")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var auth AuthResponse = AuthResponse{}
	err = json.Unmarshal(body, &auth)
	if err != nil {
		panic(err)
	}
	fmt.Println("access token: " + auth.Access_token)
	tc.bearerToken = auth.Access_token
}

func NewTwitterClient(consumerKey string, consumerSecret string) *TwitterClient {
	client := new(TwitterClient)

	// Set tokens
	client.consumerKey = consumerKey
	client.consumerSecret = consumerSecret

	// Set endpoint urls
	client.timelineUrl = "https://api.twitter.com/1.1/statuses/user_timeline.json"
	client.appAuthUrl = "https://api.twitter.com/oauth2/token"

	// Authenticate
	client.AppOnlyAuth()

	return client
}
