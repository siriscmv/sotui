package main

import (
	"fmt"

	"github.com/Siris01/sotui/utils"
)

func main() {
	fmt.Println("Visit this URL to authenticate: " + utils.GetAuthURL())
	utils.Oauth2()
	fmt.Println("Token: " + utils.GetToken())
}
