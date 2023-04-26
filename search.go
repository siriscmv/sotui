package main

import (
	"regexp"

	googlesearch "github.com/rocketlaunchr/google-search"
)

func Search(query string, site string, sort string, order string, filter string) SEResponse {
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
		filter = "!m()D0hHD1-.c61_vXxpH8BorZ9taft2)4vH6)J2QabmX)URKjC*VS(z2"
	}

	return MakeRequest(RequestOptions{
		IDs:    ids,
		Sort:   sort,
		Order:  order,
		Site:   site,
		Filter: filter,
	})
}
