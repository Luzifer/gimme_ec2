package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func findRecentUbuntuAMI() (string, error) {
	data := struct {
		Data [][]string `json:"aaData"`
	}{}
	resp, err := http.Get(ubuntuReleaseTableURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// JSON is invalid, has ',' after last array element
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	raw = []byte(strings.Replace(string(raw), "\n", "", -1))
	raw = []byte(strings.Replace(string(raw), "],]}", "]]}", -1))

	if err := json.Unmarshal(raw, &data); err != nil {
		return "", err
	}

	for _, d := range data.Data {
		if d[0] != defaultRegion || d[1] != defaultUbuntuVersion || d[3] != "amd64" || d[4] != "hvm:ebs" {
			continue
		}

		// Table is not intended for our use, we need to parse a string...
		// d[6] = "<a href=\"https://console.aws.amazon.com/ec2/home?region=ap-northeast-1#launchAmi=ami-2a0fa42b\">ami-2a0fa42b</a>"

		rex := regexp.MustCompile(".*\">([^<]+)</a>")
		img := rex.FindStringSubmatch(d[6])

		if len(img) != 2 {
			return "", errors.New("Unable to parse image table")
		}

		return img[1], nil
	}

	return "", errors.New("No suitable image found")
}
