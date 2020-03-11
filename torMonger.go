package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"tor/links"
)

type urls []string

var depth int
var urlFlag urls
var port string
var apiEndpoint string

func init() {
	//Init the command line arguments.
	flag.Var(&urlFlag, "url", "Base URL's to initiate the crawler.")
	flag.IntVar(&depth, "depth", 20, "Recursion depth. How deep do you want to go?")
	flag.StringVar(&port, "port", "9150", "The socks5 port to send the requests to.")
	flag.StringVar(&apiEndpoint, "api", nil, "The API endpoint to POST JSON data to. Format: <ip>:<port>")
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
	if apiEndpoint != nil {
		sendtoApi(url)
	}
	fmt.Println(url)
	list, err := links.Extract(url, port)
	if err != nil {
		log.Print(err)
	}
	return list
}

func sendToApi(url string) {
	client := http.Client{}
	req, err := http.NewRequest("POST", apiEndpoint, nil)
	if err != nil {
		log.Printf("Error staging request to API endpoint: %v. Error: %v", apiEndpoint, err)
	}
	data := url.Values{}
	data.Add("link", url)
	data.Add("data", time.Now().Format("Mon Jan 2 15:04:05 MST 2006"))
	req.PostForm = data
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending data to API. Error: %v", err)
	}
	log.Println("Data sent to API.")
}

func main() {
	//necessary to call within main for parsing of flags.
	flag.Parse()

	worklist := make(chan []string)  // lists of URLs, may have duplicates
	unseenLinks := make(chan string) // de-duplicated URLs

	// Add command-line arguments to worklist.
	go func() {
		var returnedLinks []string
		if len(urlFlag) > 0 {
			for _, link := range urlFlag {
				if strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://") {
					returnedLinks = append(returnedLinks, link)
				} else {
					log.Println("No protocol Scheme, defaulting to  http.")
					returnedLinks = append(returnedLinks, fmt.Sprintf("http://%v", link))
				}
			}
			worklist <- returnedLinks
		} else {
			fmt.Println("It appears that you did not provide a URL, Please provide a starting URL.")
			os.Exit(0)
		}
	}()

	for i := 0; i < depth; i++ {
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
