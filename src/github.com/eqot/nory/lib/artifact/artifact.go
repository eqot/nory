package artifact

import (
	"strings"
	"sync"
)

// Split returns []string from string for getting group ID, artifact ID and version
func Split(art string) []string {
	return strings.Split(art, ":")
}

// GetVersion returns artifact version
func GetVersion(art string) string {
	return Split(art)[2]
}

// IsSameArtifact returns boolean if specified artifacts are the same
func IsSameArtifact(art1, art2 string) bool {
	arts1 := Split(art1)
	arts2 := Split(art2)

	if arts1[0] == arts2[0] && arts1[1] == arts2[1] {
		return true
	}

	return false
}

// GetLatest returns the latest artifact from specified artifacts
func GetLatest(art1, art2 string) string {
	if !IsSameArtifact(art1, art2) {
		return ""
	}

	arts1 := Split(art1)
	arts2 := Split(art2)

	if arts1[2] > arts2[2] {
		return art1
	}
	return art2
}

// GetArtifactWithLatestVersion returns artifact with latest version
func GetArtifactWithLatestVersion(art string, ch chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	artifactRepo := &Maven{}

	latestArt, _ := artifactRepo.GetLatestVersion(art)
	if latestArt == "" {
		return
	}

	latestVersion := GetVersion(latestArt)

	ch <- art + ":" + latestVersion
}

// IsUpdatable returns flag if specified artifact can be updatable
func IsUpdatable(art string) bool {
	arts := Split(art)

	return arts[2] < arts[3]
}
