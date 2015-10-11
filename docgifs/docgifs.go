package docgifs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/kurrik/oauth1a"
	"github.com/kurrik/twittergo"
)

var (
	twitterConsumerKey       = os.Getenv("TWITTER_CONSUMER_KEY")
	twitterConsumerSecret    = os.Getenv("TWITTER_CONSUMER_SECRET")
	twitterAccessToken       = os.Getenv("TWITTER_ACCESS_TOKEN")
	twitterAccessTokenSecret = os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")

	giphyURL   = "http://media4.giphy.com/media/14aUO0Mf7dWDXW/giphy.gif"
	searchText = "oh no"
	mu         sync.Mutex
)

type TemplatePage struct {
	GiphyURL   string
	SearchText string
}

// Can only hit the api 180 times / 15 mins (So about every 5 seconds)
func PeriodicallyRefresh() {
	if err := refresh(); err != nil {
		log.Println("ERROR:", "docgifs:", err.Error())
	}
	for range time.Tick(15 * time.Minute / 180) {
		if err := refresh(); err != nil {
			log.Println("ERROR:", "docgifs:", err.Error())
		}
	}
}

func refresh() error {
	config := &oauth1a.ClientConfig{
		ConsumerKey:    twitterConsumerKey,
		ConsumerSecret: twitterConsumerSecret,
	}
	user := oauth1a.NewAuthorizedConfig(twitterAccessToken, twitterAccessTokenSecret)
	client := twittergo.NewClient(config, user)
	query := url.Values{
		"screen_name": []string{"docdocdocbrown"},
		"count":       []string{"1"},
	}
	url := fmt.Sprintf("/1.1/statuses/user_timeline.json?%v", query.Encode())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	resp, err := client.SendRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var statuses []map[string]interface{}
	err = json.Unmarshal(body, &statuses)
	if err != nil {
		return err
	}

	if len(statuses) < 1 {
		return errors.New("couldn't find any doc gifs")
	}

	text, ok := statuses[0]["text"].(string)
	if !ok {
		return errors.New("couldn't find any text in the tweet")
	}

	lastWS := strings.LastIndex(text, " ")
	if lastWS < 0 {
		return errors.New("couldn't find any search text in the tweet")
	}

	fields := strings.Fields(text)
	shortURL := fields[len(fields)-1]

	req, err = http.NewRequest("GET", shortURL, nil)
	if err != nil {
		return err
	}

	// Use jquery-like syntax to parse the giphy page for the direct gif url.
	doc, err := goquery.NewDocument(shortURL)
	if err != nil {
		return err
	}

	mu.Lock()
	defer mu.Unlock()
	giphyURL, ok = doc.Find("meta[property='og:image']").Attr("content")
	if !ok {
		return errors.New("couldn't locate giphy url")
	}
	searchText = strings.Join(fields[:len(fields)-1], " ")

	return nil
}

func CurrentPage() TemplatePage {
	return TemplatePage{
		GiphyURL:   giphyURL,
		SearchText: searchText,
	}
}
