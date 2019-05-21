package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type Secrets struct {
	Host string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	host, username, password := readSecrets()
	fmt.Println(host, username, password)
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
