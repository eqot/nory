package artifact

type ArtifactRepo interface {
	Find(term string) ([]string, error)
	GetLatestVersion(art string) string
}
