package artifact

import (
	"encoding/json"
	"net/http"
)

type Maven struct {
}

const url string = "https://search.maven.org/solrsearch/select?rows=20&wt=json&q="

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

	res, err := http.Get(url + term)
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
