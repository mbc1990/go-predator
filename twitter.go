package main

import "encoding/base64"
import "encoding/json"
import "net/http"
import "io/ioutil"
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
	Created_at string
	Text       string
	Id_str     string
	Media_url  string
	Entities   struct {
		Media []struct {
			Media_url string
		}
	}
}

type AuthResponse struct {
	Token_type   string
	Access_token string
}

// Returns a slice of tweets by the requested user
func (tc *TwitterClient) GetTweets(user string) []Tweet {
	url := tc.timelineUrl + "?screen_name=" + user
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+tc.bearerToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8.")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	// Unmarshall response into a slice of Tweet structs
	tweets := make([]Tweet, 0)
	json.Unmarshal(body, &tweets)
	return tweets
}

// Gets a bearer token for app-only authentication and sets it on the receiver struct
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
