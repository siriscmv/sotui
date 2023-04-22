package commands

import (
	"fmt"
	"net/url"

	"github.com/Siris01/sotui/utils"
)

func Search(query string) {
	res, err := utils.MakeRequest("search", url.Values{
		"order": {"desc"},
		"sort": {"votes"},
		"intitle": {query},
	}) //TODO: Search google first then get ids and search stackoverflow, stackoverflow search api barely works

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	fmt.Println(res)
}