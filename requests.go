package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
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

func MakeRequest(opts RequestOptions) {
	url := opts.GetURL()
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "sotui")
	req.Header.Set("Connection", "keep-alive")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respBytes, _ := ioutil.ReadAll(resp.Body)
	gzipReader, _ := gzip.NewReader(bytes.NewReader(respBytes))
	decompressedData, _ := ioutil.ReadAll(gzipReader)
	json := string(decompressedData)

	fmt.Println(json)
}

//TODO: Create type based on response json and send it back, also handle errors like backoff etc
