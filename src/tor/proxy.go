package tor

import (
	"fmt"
	"net/http"

	"golang.org/x/net/proxy"
)

//Tunnels traffic through the socks5 proxy.
func ConnectToProxy(url, port string) (*http.Response, error) {
	//Look into changing the hard coded port, allow user to enter that as a flag and 9150 is the default.
	dialSocksProxy, err := proxy.SOCKS5("tcp", fmt.Sprintf("127.0.0.1:%v", port), nil, proxy.Direct)
	if err != nil {
		fmt.Println("Error connecting to proxy:", err)
	}
	tr := &http.Transport{Dial: dialSocksProxy.Dial}

	// Create client
	myClient := &http.Client{
		Transport: tr,
	}
	//Chage IP/Addr to be fed in via links.
	resp, err := myClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Error calling onion service. Error: %v", err)
	} else {
		return resp, nil
	}
}
