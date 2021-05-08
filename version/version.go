package version

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	version "github.com/mcuadros/go-version"
)

var Version string

type Release struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

type Asset struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func getLatestRelease() (*Release, error) {
	body, err := downloadFile("https://api.github.com/repos/liamg/aminal/releases/latest")
	if err != nil {
		return nil, err
	}

	release := Release{}
	if err := json.Unmarshal(body, &release); err != nil {
		return nil, err
	}

	return &release, nil
}

func GetNewerRelease() (*Release, error) {
	release, err := getLatestRelease()
	if err != nil {
		return nil, err
	}

	if version.Compare(Version, release.TagName, "<") {
		return release, nil
	}

	return nil, nil
}

func downloadFile(url string) ([]byte, error) {
	spaceClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := spaceClient.Do(req)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(res.Body)
}
