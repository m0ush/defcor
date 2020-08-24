package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
)

const baseURL = "https://cloud.iexapis.com"

func endpoint(strs ...string) string {
	base, err := url.Parse(baseURL)
	if err != nil {
		panic(err)
	}
	relative, err := url.Parse(path.Join(strs...))
	if err != nil {
		panic(err)
	}
	params := relative.Query()
	params.Set("token", os.Getenv("IEXCLOUD_SECRET"))
	relative.RawQuery = params.Encode()
	urlString := base.ResolveReference(relative).String()
	return urlString

}

func makeRequest() ([]byte, error) {
	resp, err := http.Get(endpoint("stable", "ref-data", "symbols"))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Securities returns all active securities and their accompanied data
func Securities() ([]*Stock, error) {
	data, err := makeRequest()
	if err != nil {
		return nil, err
	}
	var secs []*Stock
	if err := json.Unmarshal(data, &secs); err != nil {
		return nil, err
	}
	return secs, nil
}
