package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/charmbracelet/bubbles/table"
)

const (
	baseApiURL = "https://api.stackexchange.com/2.3"
)

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    120 * time.Second,
		DisableCompression: true,
	},
}

type ResponseItem struct {
	Tags    []string `json:"tags"`
	Answers []struct {
		Comments []struct {
			Score     int `json:"score"`
			PostID    int `json:"post_id"`
			CommentID int `json:"comment_id"`
		} `json:"comments,omitempty"`
		CommentCount int    `json:"comment_count"`
		IsAccepted   bool   `json:"is_accepted"`
		Score        int    `json:"score"`
		LastEditDate int    `json:"last_edit_date,omitempty"`
		AnswerID     int    `json:"answer_id"`
		QuestionID   int    `json:"question_id"`
		BodyMarkdown string `json:"body_markdown"`
	} `json:"answers"`
	ViewCount        int    `json:"view_count"`
	AcceptedAnswerID int    `json:"accepted_answer_id,omitempty"`
	AnswerCount      int    `json:"answer_count"`
	Score            int    `json:"score"`
	LastEditDate     int    `json:"last_edit_date,omitempty"`
	QuestionID       int    `json:"question_id"`
	BodyMarkdown     string `json:"body_markdown"`
	Link             string `json:"link"`
	Title            string `json:"title"`
}

type SEResponse struct {
	Items          []ResponseItem `json:"items"`
	HasMore        bool           `json:"has_more"`
	QuotaMax       int            `json:"quota_max"`
	QuotaRemaining int            `json:"quota_remaining"`
}

func (resp SEResponse) ToRows() []table.Row {
	rows := []table.Row{}

	for _, item := range resp.Items {
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", item.QuestionID),
			item.Title,
			fmt.Sprintf("%d", item.Score),
			fmt.Sprintf("%d", item.ViewCount),
		})
	}

	return rows
}

type RequestOptions struct {
	IDs    string
	Sort   string
	Order  string
	Site   string
	Filter string
}

func (opts RequestOptions) GetURL() string {
	return fmt.Sprintf("%s/%s/%s?site=%s&sort=%s&order=%s&filter=%s&access_token=%s&key=%s", baseApiURL, "questions", opts.IDs, opts.Site, opts.Sort, opts.Order, opts.Filter, GetToken(), authKey)
}

func MakeRequest(opts RequestOptions) SEResponse {
	url := opts.GetURL()
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "sotui")
	req.Header.Set("Connection", "keep-alive")

	resp, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respBytes, _ := ioutil.ReadAll(resp.Body)
	gzipReader, _ := gzip.NewReader(bytes.NewReader(respBytes))
	decompressedData, _ := ioutil.ReadAll(gzipReader)
	response := SEResponse{}

	err = json.Unmarshal([]byte(string(decompressedData)), &response)

	if err != nil {
		panic(err)
	}

	return response
}
