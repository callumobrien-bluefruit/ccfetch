package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Secrets struct {
	Host string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	host, username, password := readSecrets()
	statisticsJson := getStatisticsJson(host, username, password)
	fmt.Println(statisticsJson)
}

func readSecrets() (host string, username string, password string) {
	const secretsPath string = "secrets.json"

	secretsJson, err := ioutil.ReadFile(secretsPath)
	if err != nil {
		log.Fatal("Failed to read secrets.json. Error: ", err)
	}

	var secrets Secrets
	err = json.Unmarshal(secretsJson, &secrets)
	if err != nil {
		log.Fatal("Failed to parse secrets.json. Error: ", err)
	}

	return secrets.Host, secrets.Username, secrets.Password
}

func getStatisticsJson(host string, username string, password string) []byte {
	const endpoint string = "app/rest/builds/status:SUCCESS,branch:master,buildType:(id:WatsonMarlowPims_Absw),count:1/statistics"

	req, err := http.NewRequest("GET", host + endpoint, nil)
	if err != nil {
		log.Fatal("Failed to create GET request. Error: ", err)
	}
	req.SetBasicAuth(username, password)
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Failed to get statistics from TeamCity. Error: ", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Failed to read response body. Error: ", err)
	}

	return body
}
