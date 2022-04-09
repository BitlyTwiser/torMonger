package links

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"tor/src/database"
	"tor/src/logging"
	"tor/src/tor"
)

func Extract(url, port string) ([]string, error) {
	db := database.DB{}
	resp, err := tor.ConnectToProxy(url, port)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		logging.LogError(fmt.Errorf("getting %s: %s", url, resp.Status))
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}

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
