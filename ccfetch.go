package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Secrets struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Options struct {
	Host string
	BuildTypeId string
	PropertyNames []string
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
	options, err := getOptions()
	if err != nil {
		os.Exit(1)
	}

	secrets, err := readSecrets()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	statistics, err := fetchStatistics(options.Host, options.BuildTypeId, secrets)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	codeCoverage, err := statistics.extractCodeCoverage(options.PropertyNames)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	codeCoverageJson, err := json.MarshalIndent(codeCoverage, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(string(codeCoverageJson))
}

func getOptions() (Options, error) {
	var options Options
	flag.StringVar(&options.Host, "host", "http://127.0.0.1/", "The address of your TeamCity server")
	flag.StringVar(&options.BuildTypeId, "id", "", "The ID of the build type to fetch properties for")
	propertyNamesList := flag.String("props", "", "A colon-seperated list of the names of the properties to fetch")

	flag.Parse()

	if options.BuildTypeId == "" || *propertyNamesList == "" {
		flag.Usage()
		return Options{}, errors.New("Invalid command-line options")
	}

	options.PropertyNames = strings.Split(*propertyNamesList, ":")

	return options, nil
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

func fetchStatistics(host, buildTypeId string, secrets Secrets) (Statistics, error) {
	endpoint := "app/rest/builds/status:SUCCESS,branch:master,buildType:(id:" + buildTypeId + "),count:1/statistics"
	req, err := http.NewRequest("GET", host + endpoint, nil)
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

func (s Statistics) extractCodeCoverage(propertyNames []string) (map[string]float64, error) {
	codeCoverage := make(map[string]float64)

	for _, propertyName := range(propertyNames) {
		property, err := s.findProperty(propertyName)
		if err != nil {
			return codeCoverage, err
		}
		value, err := strconv.ParseFloat(property.Value, 64)
		if err != nil {
			return codeCoverage, err
		}

		codeCoverage[propertyName] = value
	}

	return codeCoverage, nil
}

func (s Statistics) findProperty(propertyName string) (Property, error) {
	for _, property := range s.Properties {
		if property.Name == propertyName {
			return property, nil
		}
	}
	return Property{}, errors.New("Could not find property " + propertyName)
}
