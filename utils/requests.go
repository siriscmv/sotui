package utils

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	baseApiURL = "https://api.stackexchange.com/2.3"
)

var client = &http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    120 * time.Second,
		DisableCompression: true,
	},
}

func MakeRequest(endpoint string, parmas url.Values) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s?%s&site=stackoverflow&access_token=%s&key=%s", baseApiURL, endpoint, parmas.Encode(), GetToken(), key)
	fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "sotui")

	req.Header.Set("Connection", "keep-alive")
	return client.Do(req)
}
