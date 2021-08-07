package main

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/garyburd/go-oauth/oauth"
	"github.com/joho/godotenv"
)

var conn net.Conn

func dial(netw, addr string) (net.Conn, error) {
	if conn != nil {
		conn.Close()
		conn = nil
	}
	netc, err := net.DialTimeout(netw, addr, 5*time.Second)
	if err != nil {
		return nil, err
	}
	conn = netc
	return netc, nil
}

var reader io.ReadCloser

func closeConn() {
	if conn != nil {
		conn.Close()
	}
	if reader != nil {
		reader.Close()
	}
}

var (
	authClient *oauth.Client
	creds      *oauth.Credentials
)

// Fetch Creds for Twitter API.
func setupTwitterAuth() {
	err := godotenv.Load("twittervotes.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ConsumerKey := os.Getenv("SP_TWITTER_KEY")
	ConsumerSecret := os.Getenv("SP_TWITTER_SECRET")
	AccessToken := os.Getenv("SP_TWITTER_ACCESSTOKEN")
	AccessSecret := os.Getenv("SP_TWITTER_ACCESSSECRET")

	creds = &oauth.Credentials{
		Token:  AccessToken,
		Secret: AccessSecret,
	}
	authClient = &oauth.Client{
		Credentials: oauth.Credentials{
			Token:  ConsumerKey,
			Secret: ConsumerSecret,
		},
	}
}

var (
	authSetupOnce sync.Once
	httpClient    *http.Client
)

func makeRequest(req *http.Request, params url.Values) (*http.Response, error) {
	authSetupOnce.Do(func() {
		setupTwitterAuth()
		httpClient = &http.Client{
			Transport: &http.Transport{
				Dial: dial,
			},
		}
	})
	formEnc := params.Encode()
	req.Header.Set("Content-Type", "application/x-www-form- urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(formEnc)))
	req.Header.Set("Authorization", authClient.AuthorizationHeader(creds,
		"POST", req.URL, params))

	return httpClient.Do(req)
}

type tweet struct {
	Text string
}

// Reads option data into a channel.
func readFromTwitter(votes chan<- string) {
	options, err := loadOptions()
	if err != nil {
		log.Println("failed to load options: ", err)
		return
	}

	u, err := url.Parse("https://stream.twitter.com/1.1/statuses/filter.json")
	if err != nil {
		log.Println("creating filter request failed: ", err)
		return
	}
	query := make(url.Values)
	query.Set("track", strings.Join(options, ","))
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(query.Encode()))
	if err != nil {
		log.Println("creating filter request failed: ", err)
		return
	}
	resp, err := makeRequest(req, query)
	if err != nil {
		log.Println("making request failed: ", err)
		return
	}
	reader := resp.Body
	decoder := json.NewDecoder(reader)
	for {
		var t tweet
		if err := decoder.Decode(&t); err != nil {
			break
		}
		for _, option := range options {
			if strings.Contains(
				strings.ToLower(t.Text),
				strings.ToLower(option),
			) {
				log.Println("vote: ", option)
				votes <- option
			}
		}
	}
}
