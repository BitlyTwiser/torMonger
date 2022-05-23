package links

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"tor/src/database"
	"tor/src/logging"
	"tor/src/tor"
	"tor/src/types"

	"golang.org/x/net/html"
)

var db = database.DatabaseInit()

//([^http:\/\/||https:\/\/||.onion])([a-zA-Z1-9]+)
// Extract extracts the html from the onion site, parses html and stores link and data in the database.
func Extract(url, port string, overrideHtml bool) ([]string, error) {
	tormongerData := types.TormongerDataValues{}
	htmlReferenceData := database.HtmlDataReference{}
	resp, err := tor.ConnectToProxy(url, port)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		logging.LogError(fmt.Errorf("getting %s: %s", url, resp.Status))
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}

	parseLinkAttributesFindOrCreate(url, &tormongerData, &htmlReferenceData)
	if tormongerData.FoundValues || overrideHtml {
		db.CreateOrUpdateHtmlData(returnRawHtmlData(resp), tormongerData, htmlReferenceData)
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}

	var links []string
	visitNode := func(n *html.Node) {
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

// This is due to my own laziness, I could reconstruct the HTML from the html.Node itself and parse the values of each
// node attributes into a string, but alas, I did not want to.
func returnRawHtmlData(response *http.Response) string {
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logging.LogError(fmt.Errorf("getting error from pasrsing data%s:", err.Error))
	}

	responseString := string(responseData)

	return responseString
}

// Will only update data for newfound values unless overrideHtml is
func parseLinkAttributesFindOrCreate(link string, tormongerData *types.TormongerDataValues, htmlReferenceData *database.HtmlDataReference) {
	var hasSubDirInDatabase bool = false
	var tormongerSubDirId string

	//Parse out original onion link
	regex, err := regexp.Compile("([^http:\\/\\/||https:\\/\\/||.onion])([a-zA-Z1-9]+)")
	if err != nil {
		logging.LogError(fmt.Errorf("error parsing regex: %s", err.Error()))
		tormongerData.FoundValues = false
		return
	}

	if !regex.MatchString(link) {
		logging.LogError(fmt.Errorf("error matching onion url to regular expression: %s", link))
		tormongerData.FoundValues = false
		return
	} else {
		match := regex.FindString(link)
		hasReference, tormongerDataId := linkReferenceInDatabase(base64EncodeString(match))
		hasSubdirectories, subDirsMatch := linkHasSubdirectories(link)

		if hasSubdirectories {
			hasSubDirInDatabase, tormongerSubDirId = subDirExistsInDatabase(tormongerDataId, subDirsMatch)
		}

		if !hasReference {
			//Log value in database then capture html
			tormongerDataId = createTormongerDataRecord(link, match)
			tormongerData.FoundValues = true
		}
		if hasSubdirectories && !hasSubDirInDatabase {
			//Strip subdomain and check if it already exists as well
			tormongerSubDirId = createSubDirectoryRecord(link, subDirsMatch, tormongerDataId)
			tormongerData.FoundValues = true
		}
		if !tormongerData.FoundValues && !linkHasHtmlRecords(tormongerDataId, htmlReferenceData) {
			// Add all values to struct in case overrideHTML was thrown.
			logging.Log(fmt.Sprintf("All data already exists for %s in database. No new data will be added.", link))
			tormongerData.TormongerDataId = tormongerDataId
			tormongerData.TormongerDataSubDirId = tormongerSubDirId
			tormongerData.FoundValues = false
			return
		}

		tormongerData.TormongerDataId = tormongerDataId
		tormongerData.TormongerDataSubDirId = tormongerSubDirId
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

	subdirectoriesMatch := regex.ReplaceAllString(link, "")
	if len(subdirectoriesMatch) > 0 {
		return true, subdirectoriesMatch
	}

	return false, ""
}

func subDirExistsInDatabase(tormongerDataId, subdirectoriesMatch string) (bool, string) {
	//Parse Subdirectories
	//Find the URL up until the .onion, then cut away leaving only subdirs.

	values, err := db.FindSubDirectoryMatch(tormongerDataId, subdirectoriesMatch, database.SubdirctoryReference{})
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

func linkHasHtmlRecords(tormongerDataId string, htmlData *database.HtmlDataReference) bool {
	values, err := db.FindHtmlRecordForLink(tormongerDataId, database.HtmlDataReference{})
	if err != nil {
		logging.LogError(fmt.Errorf("error obtaining value from subdirectory match: %s", err.Error()))
	}

	if len(values.TormongerDataId) > 0 {
		htmlData.FoundValues = true
		htmlData.Id = values.Id
		htmlData.TormongerDataId = values.TormongerDataId
		htmlData.TormongerDataSubDirectoriesId = values.TormongerDataSubDirectoriesId
		return true
	}

	htmlData.FoundValues = false
	return false
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
func htmlAttributeParser(n *html.Node, tormongerData string) {
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
