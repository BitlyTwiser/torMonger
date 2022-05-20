package links

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"regexp"
	"tor/src/database"
	"tor/src/logging"
	"tor/src/tor"

	"golang.org/x/net/html"
)

var db = database.DatabaseInit()

//([^http:\/\/||https:\/\/||.onion])([a-zA-Z1-9]+)
// Extract extracts the html from the onion site, parses html and stores link and data in the database.
func Extract(url, port string) ([]string, error) {
	resp, err := tor.ConnectToProxy(url, port)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		logging.LogError(fmt.Errorf("getting %s: %s", url, resp.Status))
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}

	stripLinkCheckForDuplicates(url)
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}

	var links []string
	visitNode := func(n *html.Node) {
		//Look at parsing all attributes here, building the html, saving off into the data table with the link.
		htmlAttributeParser(n)
		db.Insert()
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue // ignore bad URLs
				}
				links = append(links, link.String())
			}
		}
	}
	forEachNode(doc, visitNode, nil)
	return links, nil
}

//
func stripLinkCheckForDuplicates(link string) {
	regex, err := regexp.Compile("([^http:\\/\\/||https:\\/\\/||.onion])([a-zA-Z1-9]+)")
	if err != nil {
		logging.LogError(fmt.Errorf("error parsing regex: %s", err.Error()))
	}

	if !regex.MatchString(link) {
		logging.LogError(fmt.Errorf("error matching onion url to regular expression: %s", link))
	} else {
		match := regex.FindString(link)
		encoded := base64.StdEncoding.EncodeToString([]byte(match))
		if !linkReferenceInDatabase(encoded) {
			fmt.Println("Sup")
			//Log value in database
		}
		fmt.Println(match)
		fmt.Println(encoded)
	}
}

func linkReferenceInDatabase(link string) bool {
	values, err := db.FindLinkReference(link, database.LinkReference{})
	if err != nil {
		logging.LogError(fmt.Errorf("errror retreiving values from database: %s", err))
	}
	fmt.Println(values)
	return true
}

// Parses, then re-assembles the html node values in an attempt to re-build a snapshot of the html from the onion site.
func htmlAttributeParser(n *html.Node) {
	for i, val := range n.Attr {
		fmt.Sprintf("%d, val: %s", i, val)
	}
}

func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}
