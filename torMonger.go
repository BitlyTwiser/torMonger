package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
	"tor/src/links"
	"tor/src/logging"
)

type urls []string

var threads int
var urlFlag urls
var port string

func init() {
	//Init the command line arguments.
	flag.Var(&urlFlag, "url", "Base URL to initiate the crawler.")
	flag.IntVar(&threads, "threads", 1, "how many threads to spawn. Set at 1 initially, but can run as many as your hardware allows")
	flag.StringVar(&port, "port", "9050", "The socks5 port to send the requests to. When one runs tor from CLI, the initial port is 9050")
}

//Part of the flag.value interface.
func (s *urls) String() string {
	return fmt.Sprint(*s)
}

//Need to have satisfy the flag value interface when using Var
func (i *urls) Set(url string) error {
	for _, u := range strings.Split(url, ",") {
		*i = append(*i, u)
	}
	return nil
}

//Call the imported links library and crawl the network.
func crawl(url string) []string {
	fmt.Println(url)
	list, err := links.Extract(url, port)
	if err != nil {
		logging.LogError(fmt.Errorf("error in crawl function: %s", err))
	}
	return list
}

// ([^http:\/\/||https:\/\/||.onion])([a-zA-Z1-9]+)
func stripLinkCheckForDuplicates(link string) {

}

func loadEnvFile() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	//necessary to call within main for parsing of flags.
	flag.Parse()

	//Load Env Vars
	loadEnvFile()

	worklist := make(chan []string)  // lists of URLs, may have duplicates
	unseenLinks := make(chan string) // de-duplicated URLs

	// Add command-line arguments to worklist.
	go func() {
		//Remember to remove me as I am just for tsting so we can use delve easily as we are lazy.
		urlFlag = append(urlFlag, "http://zqktlwiuavvvqqt4ybvgvi7tyo4hjl5xgfuvpdf6otjiycgwqbym2qad.onion")
		var returnedLinks []string
		if len(urlFlag) > 0 {
			for _, link := range urlFlag {
				if strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://") {
					returnedLinks = append(returnedLinks, link)
				} else {
					logging.Log("\"No protocol Scheme, defaulting to  http.\"")
					returnedLinks = append(returnedLinks, fmt.Sprintf("http://%v", link))
				}
			}
			worklist <- returnedLinks
		} else {
			logging.Log("It appears that you did not provide a URL, Please provide a starting URL.")
			os.Exit(0)
		}
	}()

	for i := 0; i < threads; i++ {
		go func() {
			for link := range unseenLinks {
				foundLinks := crawl(link)
				go func() { worklist <- foundLinks }()
			}
		}()
	}

	// The main goroutine de-duplicates worklist items
	// and sends the unseen ones to the crawlers.
	seen := make(map[string]bool)
	for list := range worklist {
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				unseenLinks <- link
			}
		}
	}
}
