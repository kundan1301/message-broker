package http

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var CustomClient *http.Client

// InitHTTPClient change default http client parameters.
func InitHTTPClient() {
	defaultRoundTripper := http.DefaultTransport
	defaultTransportPointer, ok := defaultRoundTripper.(*http.Transport)
	if !ok {
		panic(fmt.Sprintf("defaultRoundTripper not an *http.Transport"))
	}
	defaultTransport := *defaultTransportPointer
	defaultTransport.MaxIdleConns = 500
	defaultTransport.MaxIdleConnsPerHost = 50
	CustomClient = &http.Client{Transport: &defaultTransport}

}

func Get(url string) ([]byte, error) {
	resp, err := CustomClient.Get(url)
	if err != nil {
		log.Println("error", err)
		return nil, nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

func ClosePrevConnection(clientID string) {

}

func CheckAuth(authUrl, userName, password, clientID string) bool {
	authUrl = authUrl + "?username=" + userName + "&password=" + password + "&clientid=" + clientID
	parsedURL, err := url.ParseRequestURI(authUrl)
	if err != nil {
		log.Println("error in parsing url", err)
		return false
	}
	resp, err := CustomClient.Get(parsedURL.String())
	if err != nil {
		log.Println("error in auth http call", err)
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		return true
	}
	return false
}
