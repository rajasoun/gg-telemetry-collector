package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func formWebEndPointForRepo(repoName string) string {
	domain := "gitscan-wwwin.cisco.com"
	api_end_point := "/api/v1/accounts/2/sources/?"
	params := url.Values{}
	params.Add("monitored", "true")
	params.Add("page", "1")
	params.Add("page_size", "10")
	params.Add("ordering", "-open_issues_count")
	params.Add("search", repoName)
	url := "https://" + domain + api_end_point + params.Encode()
	return url
}

func getContentFromFile(fileName string) string {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	// Convert byte[] to string
	cookieContent := string(content)
	return cookieContent
}

func main() {
	client := &http.Client{}

	url := formWebEndPointForRepo("LearningAtCisco/Cairo")
	method := "GET"

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Get content from cookie.txt
	cookieContent := getContentFromFile("cookie.txt")
	// Remove cookie header
	cookie := strings.ReplaceAll(cookieContent, "cookie: ", "")

	req.Header.Add("cookie", cookie)

	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
		return
	}
	json := string(body)
	log.Println(json)
}
