package artifact

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Maven struct {
	host string
}

const (
	defaultHost string = "https://search.maven.org"
	defaultPath string = "/solrsearch/select?rows=20&wt=json&q="
)

type response struct {
	Response *docs `json:"response"`
}

type docs struct {
	Docs []artifact `json:"docs"`
}

type artifact struct {
	GroupID    string `json:"g"`
	ArtifactID string `json:"a"`
	Version    string `json:"latestVersion"`
}

func (m *Maven) Find(term string) ([]string, error) {
	var arts []string

	var query string

	terms := Split(term)
	if len(terms) > 1 {
		query = "a:\"" + terms[1] + "\"+AND+g:\"" + terms[0] + "\""
	} else {
		query = term
	}

	var url string
	if m.host != "" {
		url = m.host
	} else {
		url = defaultHost
	}

	res, err := http.Get(url + defaultPath + query)
	if err != nil {
		return arts, err
	}

	defer res.Body.Close()

	var data response
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return arts, err
	}

	for _, art := range data.Response.Docs {
		arts = append(arts, art.GroupID+":"+art.ArtifactID+":"+art.Version)
	}

	return arts, nil
}

func (m *Maven) GetLatestVersion(art string) (string, error) {
	arts, err := m.Find(art)
	if err != nil {
		return "", err
	}

	if len(arts) == 0 {
		return "", fmt.Errorf("No artifacts found")
	}

	return arts[0], nil
}
