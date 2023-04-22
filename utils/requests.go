package utils

import (
	"fmt"
	"net/http"
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

type RequestOptions struct {
	IDs    string
	Sort   string
	Order  string
	Site   string
	Filter string
}

func (opts RequestOptions) GetURL() string {
	return fmt.Sprintf("%s/%s/%s/%s?site=%s&sort=%s&order=%s&filter=%s&access_token=%s&key=%s", baseApiURL, "questions", opts.IDs, "answers", opts.Site, opts.Sort, opts.Order, opts.Filter, GetToken(), key)
}

func MakeRequest(opts RequestOptions) (*http.Response, error) {
	url := opts.GetURL()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept-Charset", "utf-8")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Host", "api.stackexchange.com")
	req.Header.Set("User-Agent", "sotui")
	req.Header.Set("Connection", "keep-alive")
	return client.Do(req)
}
