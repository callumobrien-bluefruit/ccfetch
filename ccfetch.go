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

type Statistics struct {
	Count int `json:"count"`
	Properties []Property `json:"property"`
}

type Property struct {
	Name string `json:"name"`
	Value string `json:"value"`
}

func main() {
	secrets, err := readSecrets()
	if err != nil {
		log.Fatal(err)
	}

	statistics, err := fetchStatistics(secrets)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(statistics)
}

func readSecrets() (Secrets, error) {
	const secretsPath string = "secrets.json"

	secretsJson, err := ioutil.ReadFile(secretsPath)
	if err != nil {
		return Secrets{}, err
	}

	var secrets Secrets
	err = json.Unmarshal(secretsJson, &secrets)
	if err != nil {
		return Secrets{}, err
	}

	return secrets, nil
}

func fetchStatistics(secrets Secrets) (Statistics, error) {
	const endpoint string = "app/rest/builds/status:SUCCESS,branch:master,buildType:(id:WatsonMarlowPims_Absw),count:1/statistics"

	req, err := http.NewRequest("GET", secrets.Host + endpoint, nil)
	if err != nil {
		return Statistics{}, err
	}
	req.SetBasicAuth(secrets.Username, secrets.Password)
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Statistics{}, err
	}

	statisticsJson, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Statistics{}, err
	}

	var statistics Statistics
	err = json.Unmarshal(statisticsJson, &statistics)
	if err != nil {
		return Statistics{}, err
	}

	return statistics, nil
}
