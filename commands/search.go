package commands

import (
	"fmt"
	"io"
	"regexp"

	"github.com/Siris01/sotui/utils"
	googlesearch "github.com/rocketlaunchr/google-search"
)

func Search(query string, site string, sort string, order string, filter string) {
	searchResults, err := googlesearch.Search(nil, query+" site:stackoverflow.com") //TODO: Fix this so that it works for all sites. Note: site is NOT the full domain
	if err != nil {
		panic(err)
	}

	ids := ""
	re := regexp.MustCompile("/questions/([0-9]+)/")

	for _, result := range searchResults {
		questionId := re.FindStringSubmatch(result.URL)[1]

		if ids == "" {
			ids = questionId
		} else {
			ids = ids + ";" + questionId
		}
	}

	if site == "" {
		site = "stackoverflow"
	}
	if sort == "" {
		sort = "votes"
	}
	if order == "" {
		order = "desc"
	}
	if filter == "" {
		filter = "!szz.51ErE5dRYIAadZEuxVMHA5r6Nj7"
	}

	res, err := utils.MakeRequest(utils.RequestOptions{
		IDs:    ids,
		Sort:   sort,
		Order:  order,
		Site:   site,
		Filter: filter,
	})

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}
