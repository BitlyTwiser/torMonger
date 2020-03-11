package tor

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/proxy"
)

func peelOnion(url string) error {
	//cheap but effective, we should build a regex to validate that http/https is used and that the .onion value is present.
	if !strings.Contains(url, ".onion") {
		return errors.New("Provided link is not an onion site. Moving on..")
	}
	return nil
}

//Tunnels traffic through the socks5 proxy.
func ConnectToProxy(url, port string) (*http.Response, error) {
	//Look into changing the hard coded port, allow user to enter that as a flag and 9150 is the default.
	dialSocksProxy, err := proxy.SOCKS5("tcp", fmt.Sprintf("127.0.0.1:%v", port), nil, proxy.Direct)
	if err != nil {
		fmt.Println("Error connecting to proxy:", err)
	}
	tr := &http.Transport{Dial: dialSocksProxy.Dial}

	// Create client using transport layer created above.
	proxyClient := &http.Client{
		Transport: tr,
	}
	err = peelOnion(url)
	if err != nil {
		return nil, err
	} else {
		resp, err := proxyClient.Get(url)
		if err != nil {
			return nil, fmt.Errorf("Error calling onion service. Error: %v", err)
		} else {
			return resp, nil
		}
	}
}
