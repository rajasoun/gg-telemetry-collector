package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/tidwall/gjson"
)

func loadDotEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func createHttpClient() (*http.Client, *http.Request, bool) {
	client := &http.Client{}
	url := formWebEndPointForRepo(os.Getenv("TEST_REPO_NAME"))
	method := "GET"

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Println(err)
		return nil, nil, true
	}
	cookieContent := getContentFromFile("cookie.txt")
	cookie := strings.ReplaceAll(cookieContent, "cookie: ", "")

	req.Header.Add("cookie", cookie)
	return client, req, false
}

func executeHttpRequest(client *http.Client, req *http.Request) ([]byte, bool) {
	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
		return nil, true
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
		return nil, true
	}
	return body, false
}

func getSecretsCount(body []byte, repoName string) string {
	json := string(body)
	query := "results.#(url%\"*/" + repoName + "\").open_issues_count"
	println(query)
	value := gjson.Parse(json).Get(query)
	return value.String()
}

func formWebEndPointForRepo(repoName string) string {
	domain := os.Getenv("GITGUARDIAN_URL")
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
	loadDotEnv()

	client, req, onErr := createHttpClient()
	if onErr {
		return
	}

	body, onErr := executeHttpRequest(client, req)
	if onErr {
		return
	}

	secretsCount := getSecretsCount(body, os.Getenv("TEST_REPO_NAME"))
	println(secretsCount)
}
