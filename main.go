package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/jinzhu/now"
	"github.com/kurrik/oauth1a"
	"github.com/kurrik/twittergo"
)

const (
	// BANNER ...
	BANNER = `
go-horoscope - Get your horoscope via the @poetastrologers twitter account
`
)

var (
	twitterConsumerKey    string
	twitterConsumerSecret string
	twitterAccount        string
	numTweets             string
	tweet                 []byte
	astroSign             string
)

func init() {
	signs := []string{"Aries", "Taurus", "Gemini", "Cancer", "Leo", "Virgo", "Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces"}

	// declare flags for consumer key & secret (hint: register for a developer account at Twitter to receive keys)
	flag.StringVar(&twitterConsumerKey, "consumer-key", os.Getenv("CONSUMER_KEY"), "Twitter consumer key")
	flag.StringVar(&twitterConsumerSecret, "consumer-secret", os.Getenv("CONSUMER_KEY"), "Twitter consumer secret")

	// print usage information when -h or -help flag is invoked
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, fmt.Sprintf(BANNER))
		fmt.Println()
		fmt.Println("Usage:  go-horoscope [options] <sign>")
		fmt.Println()
		flag.PrintDefaults()
	}

	// parse flags
	flag.Parse()

	// check consumer key
	if twitterConsumerKey == "" {
		if twitterConsumerKey = os.Getenv("CONSUMER_KEY"); twitterConsumerKey == "" {
			flag.PrintDefaults()
			os.Exit(1)
		}
	}

	// check consumer secret
	if twitterConsumerSecret == "" {
		if twitterConsumerSecret = os.Getenv("CONSUMER_SECRET"); twitterConsumerSecret == "" {
			flag.PrintDefaults()
			os.Exit(1)
		}
	}

	// get argumemnt
	initialArgument := flag.Args()[0]

	// validate user input
	for _, value := range signs {
		if strings.EqualFold(initialArgument, value) == true {
			astroSign = initialArgument
			break
		}
	}
	if astroSign == "" {
		fmt.Println("Enter a valid astrological sign.")
		os.Exit(1)
	}

}

func main() {
	// we'll hardcode the twitter account to pull from (for now) and the number of tweets to query
	twitterAccount = "poetastrologers"
	numTweets = "100"

	// we use jinzhu's Now library to calcuate the beginning of the week (i.e. when @poetastrologers post their weekly horoscopes)
	current := now.BeginningOfWeek()
	month := current.Month()
	day := current.Day()

	// join strings from slice
	stringAr := []string{"Week of ", strconv.Itoa(int(month)), "/", strconv.Itoa(day), " in ", astroSign}
	dateStr := strings.Join(stringAr, "")

	// create config
	config := &oauth1a.ClientConfig{
		ConsumerKey:    twitterConsumerKey,
		ConsumerSecret: twitterConsumerSecret,
	}

	// create client
	client := twittergo.NewClient(config, nil)
	if err := client.FetchAppToken(); err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't fetch app token: %v\n", err)
		os.Exit(2)
	}

	// set url
	value := url.Values{}
	value.Set("count", numTweets)
	value.Set("screen_name", twitterAccount)
	request, err := http.NewRequest("GET", "/1.1/statuses/user_timeline.json?"+value.Encode(), nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't parse request: %v\n", err)
		os.Exit(2)
	}

	// send request
	response, err := client.SendRequest(request)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't send request: %v\n", err)
		os.Exit(2)
	}

	// parse response
	results := &twittergo.Timeline{}
	if err := response.Parse(results); err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't parse response: %v\n", err)
		os.Exit(2)
	}

	// encode json from response
	for _, value := range *results {
		if tweet, err = json.Marshal(*results); err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't encode tweets: %v\n", err)
			os.Exit(2)
		}

		// compare string from tweet.Text() to date, print resulting tweet
		if (strings.Contains(value.Text(), dateStr)) == true {
			fmt.Println(value.Text())
		}
	}
}
