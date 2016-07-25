package artifact

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	fileNameArtID         = "./maven_test_http_rxjava.txt"
	fileNameGroupAndArtID = "./maven_test_http_io.reactivex_rxjava.txt"
	expectedArt           = "io.reactivex:rxjava:1.1.8"
)

var testHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var filename string
	if r.URL.String() == "/solrsearch/select?rows=20&wt=json&q=rxjava" {
		filename = fileNameArtID
	} else {
		filename = fileNameGroupAndArtID
	}

	res, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Failed to load test file; %q\n", filename)
	}

	w.Write(res)
})

func TestFind(t *testing.T) {
	ts := httptest.NewServer(testHandler)
	defer ts.Close()

	artifactRepo := &Maven{host: ts.URL}

	arts, err := artifactRepo.Find("rxjava")
	if err != nil {
		t.Errorf("Error %q", err)
	}

	if arts[0] != expectedArt {
		t.Errorf("For artifact, got %q; expected %q", arts[0], expectedArt)
	}

	arts, err = artifactRepo.Find("io.reactivex:rxjava")
	if err != nil {
		t.Errorf("Error %q", err)
	}

	if arts[0] != expectedArt {
		t.Errorf("For artifact, got %q; expected %q", arts[0], expectedArt)
	}
}

func TestGetLatestVersion(t *testing.T) {
	ts := httptest.NewServer(testHandler)
	defer ts.Close()

	artifactRepo := &Maven{host: ts.URL}

	art, err := artifactRepo.GetLatestVersion("rxjava")
	if err != nil {
		t.Errorf("Error %q", err)
	}

	if art != expectedArt {
		t.Errorf("For latest version, got %q; expected %q", art, expectedArt)
	}
}
