package http

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var CustomClient *http.Client

// InitHttpClient change default http client parameters.
func InitHttpClient() {
	defaultRoundTripper := http.DefaultTransport
	defaultTransportPointer, ok := defaultRoundTripper.(*http.Transport)
	if !ok {
		panic(fmt.Sprintf("defaultRoundTripper not an *http.Transport"))
	}
	defaultTransport := *defaultTransportPointer
	defaultTransport.MaxIdleConns = 100
	defaultTransport.MaxIdleConnsPerHost = 100
	CustomClient = &http.Client{Transport: &defaultTransport}

}

func Get(url string) ([]byte, error) {
	resp, err := CustomClient.Get(url)
	if err != nil {
		log.Println("error in loading config from config url: ", err)
		return nil, nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}
