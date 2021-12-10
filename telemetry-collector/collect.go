package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

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

func loadDotEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// https://mholt.github.io/json-to-go/
type SecretScanResults struct {
	Count    int         `json:"count"`
	Next     interface{} `json:"next"`
	Previous interface{} `json:"previous"`
	Results  []struct {
		ID                  int    `json:"id"`
		Monitored           bool   `json:"monitored"`
		Visibility          string `json:"visibility"`
		DisplayName         string `json:"display_name"`
		IntegrationName     string `json:"integrationName"`
		Type                string `json:"type"`
		URL                 string `json:"url"`
		Business            bool   `json:"business"`
		AlreadyScanned      bool   `json:"already_scanned"`
		AlreadyFinishedScan bool   `json:"already_finished_scan"`
		LastScan            struct {
			Date            time.Time `json:"date"`
			Status          string    `json:"status"`
			CommitsScanned  int       `json:"commits_scanned"`
			Duration        string    `json:"duration"`
			BranchesScanned int       `json:"branches_scanned"`
			SecretsCount    int       `json:"secrets_count"`
			TaskID          string    `json:"task_id"`
			SourceType      string    `json:"source_type"`
		} `json:"last_scan"`
		OpenIssuesCount               int    `json:"open_issues_count"`
		ClosedIssuesCount             int    `json:"closed_issues_count"`
		Health                        string `json:"health"`
		SpecificSource                int    `json:"specific_source"`
		OpenIssuesWithPresenceCount   int    `json:"open_issues_with_presence_count"`
		ClosedIssuesWithPresenceCount int    `json:"closed_issues_with_presence_count"`
	} `json:"results"`
}

func main() {
	loadDotEnv()

	client := &http.Client{}
	url := formWebEndPointForRepo(os.Getenv("TEST_REPO_NAME"))
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

	var secretScanResults SecretScanResults
	json.Unmarshal(body, &secretScanResults)
	log.Printf("API Response as struct %+v\n", secretScanResults)
}
