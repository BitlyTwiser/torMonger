package tor

import (
	"fmt"
	"golang.org/x/net/proxy"
	"net/http"
	"strings"
	"tor/src/logging"
)

func peelOnion(url string) (string, error) {
	//cheap but effective, we should build a regex to validate that http/https is used and that the .onion value is present.
	if strings.Contains(url, ".onion") {
		return url, nil
	} else {
		logging.LogError(fmt.Errorf("the provided link does not appear to be an onion link, ignoring link: %v", url))

		return "", fmt.Errorf("The provided link does not appear to be an onion link, ignoring link: %v.", url)
	}
}

//Tunnels traffic through the socks5 proxy.
func ConnectToProxy(url, port string) (*http.Response, error) {
	//Look into changing the hard coded port, allow user to enter that as a flag and 9150 is the default.
	dialSocksProxy, err := proxy.SOCKS5("tcp", fmt.Sprintf("127.0.0.1:%v", port), nil, proxy.Direct)
	if err != nil {
		logging.LogError(fmt.Errorf("error connecting to proxy: %s", err.Error()))
	}
	tr := &http.Transport{Dial: dialSocksProxy.Dial}

	// Create client using transport layer created above.
	proxyClient := &http.Client{
		Transport: tr,
	}
	u, err := peelOnion(url)
	if err != nil {
		return nil, err
	} else {
		resp, err := proxyClient.Get(u)
		if err != nil {
			logging.LogError(fmt.Errorf("error calling onion service. Error: %s", err.Error()))

			return nil, fmt.Errorf("Error calling onion service. Error: %v", err)
		} else {
			return resp, nil
		}
	}
}
