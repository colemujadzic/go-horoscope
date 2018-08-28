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
go-horoscope - Get your horoscope via the @poetastrologers (Astro Poets) twitter account
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
	flag.StringVar(&twitterConsumerKey, "consumer-key", "", "Twitter consumer key")
	flag.StringVar(&twitterConsumerSecret, "consumer-secret", "", "Twitter consumer secret")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, fmt.Sprintf(BANNER))
		flag.PrintDefaults()
	}

	flag.Parse()

	if twitterConsumerKey == "" {
		if twitterConsumerKey = os.Getenv("CONSUMER_KEY"); twitterConsumerKey == "" {
			flag.PrintDefaults()
			os.Exit(1)
		}
	}

	if twitterConsumerSecret == "" {
		if twitterConsumerSecret = os.Getenv("CONSUMER_SECRET"); twitterConsumerSecret == "" {
			flag.PrintDefaults()
			os.Exit(1)
		}
	}

	initialArgument := flag.Args()[0]
	// secondArgument := flag.Args()[1]

	astroSign = initialArgument
}

func main() {
	// we'll hardcode the twitter account to pull from (for now)
	twitterAccount = "poetastrologers"
	numTweets = "100"

	// ttime := time.Now().UTC()
	// month := ttime.Month()
	// day := ttime.Day() - 1
	// fmt.Println(int(month), int(day))
	// fmt.Println("Week of ", int(month), "/", day)
	// fmt.Printf("%s%d%s%d", "Week of ", int(month), "/", day)
	// stringAr := []string{"Week of ", strconv.Itoa(int(month)), "/", strconv.Itoa(day), " in ", astroSign}
	// dateStr := strings.Join(stringAr, "")
	// fmt.Println(dateStr)

	current := now.BeginningOfWeek()
	month := current.Month()
	day := current.Day()
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

	// send request
	value := url.Values{}
	value.Set("count", numTweets)
	value.Set("screen_name", twitterAccount)
	request, err := http.NewRequest("GET", "/1.1/statuses/user_timeline.json?"+value.Encode(), nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't parse request: %v\n", err)
		os.Exit(2)
	}

	response, err := client.SendRequest(request)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't send request: %v\n", err)
		os.Exit(2)
	}

	// get response
	results := &twittergo.Timeline{}
	if err := response.Parse(results); err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't parse response: %v\n", err)
		os.Exit(2)
	}

	for _, value := range *results {
		if tweet, err = json.Marshal(*results); err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't encode tweets: %v\n", err)
			os.Exit(2)
		}
		if (strings.Contains(value.Text(), dateStr)) == true {
			fmt.Println(value.Text())
		}
	}
}
