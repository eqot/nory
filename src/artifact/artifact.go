package artifact

type Artifact interface {
	Find(term string) ([]string, error)
}
