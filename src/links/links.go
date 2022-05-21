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

type tormongerDataValues struct {
	tormongerDataId       string
	tormongerDataSubDirId string
	foundValues           bool
}

var db = database.DatabaseInit()

//([^http:\/\/||https:\/\/||.onion])([a-zA-Z1-9]+)
// Extract extracts the html from the onion site, parses html and stores link and data in the database.
func Extract(url, port string, overrideHtml bool) ([]string, error) {
	resp, err := tor.ConnectToProxy(url, port)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		logging.LogError(fmt.Errorf("getting %s: %s", url, resp.Status))
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}

	recordId := stripLinkCheckForDuplicates(url)
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}

	var links []string
	visitNode := func(n *html.Node) {
		//Look at parsing all attributes here, building the html, saving off into the data table with the link.
		// Only parse if we have an object id.
		if recordId.foundValues || overrideHtml {
			htmlAttributeParser(n, recordId.tormongerDataId)
		}

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

// Will only update data for newfound values unless overrideHtml is
func stripLinkCheckForDuplicates(link string) tormongerDataValues {
	tormongerData := tormongerDataValues{}
	var hasSubDir bool = false
	var tormongerSubDirId string

	//Parse out original onion link
	regex, err := regexp.Compile("([^http:\\/\\/||https:\\/\\/||.onion])([a-zA-Z1-9]+)")
	if err != nil {
		logging.LogError(fmt.Errorf("error parsing regex: %s", err.Error()))
		tormongerData.foundValues = false
		return tormongerData
	}

	if !regex.MatchString(link) {
		logging.LogError(fmt.Errorf("error matching onion url to regular expression: %s", link))
		tormongerData.foundValues = false
		return tormongerData
	} else {
		match := regex.FindString(link)
		hasReference, tormongerDataId := linkReferenceInDatabase(base64EncodeString(match))
		hasSubdirectories, subDirsMatch := linkHasSubdirectories(link)

		if hasSubdirectories {
			hasSubDir, tormongerSubDirId = subDirExists(link, subDirsMatch)
		}

		if !hasReference {
			//Log value in database then capture html
			tormongerDataId = createTormongerDataRecord(link, match)
			tormongerData.foundValues = true
		}
		if !hasSubDir {
			//Strip subdomain and check if it already exists as well
			tormongerSubDirId = createSubDirectoryRecord(link, subDirsMatch, tormongerDataId)
			tormongerData.foundValues = true
		}
		if !tormongerData.foundValues {
			// Add all values to struct in case overrideHTML was thrown.
			logging.Log(fmt.Sprintf("All data already exists for %s in database. No new data will be added.", link))
			tormongerData.tormongerDataSubDirId = tormongerSubDirId
			tormongerData.foundValues = false
			return tormongerData
		}

		tormongerData.tormongerDataId = tormongerDataId
		tormongerData.tormongerDataSubDirId = tormongerSubDirId

		return tormongerData
	}
}

func linkHasSubdirectories(link string) (bool, string) {
	regex, err := regexp.Compile(".*?.onion")
	if err != nil {
		logging.LogError(fmt.Errorf("error parsing regex: %s", err.Error()))
		return false, ""
	}

	if !regex.MatchString(link) {
		logging.LogError(fmt.Errorf("error matching onion url to regular expression: %s", link))
		return false, ""
	}

	found := regex.FindString(link)
	fmt.Println(found)
	subdirectoriesMatch := regex.ReplaceAllString(link, "")
	if len(subdirectoriesMatch) > 0 {
		return true, subdirectoriesMatch
	}

	return false, ""
}

func subDirExists(link, subdirectoriesMatch string) (bool, string) {
	//Parse Subdirectories
	//Find the URL up until the .onion, then cut away leaving only subdirs.

	values, err := db.FindSubDirectoryMatch(base64EncodeString(subdirectoriesMatch), database.SubdirctoryReference{})
	if err != nil {
		logging.LogError(fmt.Errorf("error obtaining value from subdirectory match: %s", err.Error()))
	}

	if len(values.TormongerDataId) > 0 {
		return true, values.TormongerDataId
	}

	return false, ""
}

func createSubDirectoryRecord(link, subdirectoriesMatch, tormonger_id string) string {
	return db.CreateSubDirectoryRecord(link, subdirectoriesMatch, tormonger_id)
}

func createTormongerDataRecord(link, match string) string {
	return db.CreateTormongDataRecord(base64EncodeString(match), link)
}

func base64EncodeString(stringToEncode string) string {
	return base64.StdEncoding.EncodeToString([]byte(stringToEncode))
}

func linkReferenceInDatabase(link string) (bool, string) {
	values, err := db.FindLinkReference(link, database.LinkReference{})
	if err != nil {
		logging.LogError(fmt.Errorf("errror retreiving values from database: %s", err))
	}

	if len(values.Id) > 0 {
		return true, values.Id
	}

	return false, ""
}

// Parses, then re-assembles the html node values in an attempt to re-build a snapshot of the html from the onion site.
//Uses the ID's of the created database elements and inserts html data for these records in the html data table.
// stores data then into html_data table.
func htmlAttributeParser(n *html.Node, recordId string) {
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
