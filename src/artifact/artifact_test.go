package artifact

import "testing"

func TestSplit(t *testing.T) {
	art := Split("foo:bar:1.2.3")

	if art[0] != "foo" {
		t.Errorf("For group ID, got %q; expected %q", art[0], "foo")
	}

	if art[1] != "bar" {
		t.Errorf("For artifact ID, got %q; expected %q", art[1], "bar")
	}

	if art[2] != "1.2.3" {
		t.Errorf("For version, got %q; expected %q", art[2], "1.2.3")
	}
}

func TestIsSameArtifact(t *testing.T) {
	art := "foo:bar:1.2.3"
	sameArt := "foo:bar:2.3.4"
	differentArt := "foo:baz:1.2.3"

	if !IsSameArtifact(art, sameArt) {
		t.Errorf("Should be same but different; %q and %q", art, sameArt)
	}

	if IsSameArtifact(art, differentArt) {
		t.Errorf("Should be different but same; %q and %q", art, differentArt)
	}
}

func TestGetLatest(t *testing.T) {
	art := "foo:bar:1.2.3"
	oldArt := "foo:bar:1.2.2"
	newArt := "foo:bar:1.3.0"

	if GetLatest(art, oldArt) != art {
		t.Errorf("for latest artifact, got %q; expected %q", GetLatest(art, oldArt), art)
	}

	if GetLatest(art, newArt) != newArt {
		t.Errorf("for latest artifact, got %q; expected %q", GetLatest(art, oldArt), newArt)
	}
}
