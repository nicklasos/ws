package main

import (
	"os"
	"testing"
)

func TestIsSpam(t *testing.T) {
	r := isSpam("Hello, how are you doing?")
	if r {
		t.Error("Should not be detected as spam")
	}

	str := `foo bar https://lol.gr
		foo bar www.le
		www.vbucks.net
		vbucks.te`

	r = isSpam(str)
	if !r {
		t.Error("Should be detected as spam")
	}

	str = `foo bar vbucks.te`

	r = isSpam(str)
	if !r {
		t.Error("Should be detected as spam")
	}
}

func TestSpamExclude(t *testing.T) {
	os.Setenv("SPAM_LINKS_EXCLUDE", `https:\/\/site\.foo[^\s]*`)

	r := isSpam(`foo bar https://site.foo/some/page`)
	if r {
		t.Error("Excluded link should not be detected as spam")
	}

	r = isSpam(`foo bar www.localhost.net https://site.foo/some/page`)
	if !r {
		t.Error("Should be not excludede from spam")
	}
}

func TestSpamExcludeWithoutEnvParam(t *testing.T) {
	r := isSpam(`foo bar www.localhost.net`)
	if !r {
		t.Error("Should be marked as spam")
	}
}
