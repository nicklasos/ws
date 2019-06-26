package main

import (
	"os"
	"regexp"
)

var re = regexp.MustCompile(`(?m)[-a-zA-Z0-9@:%_\+.~#?&//=]{2,256}\.[a-z]{2,4}\b(\/[-a-zA-Z0-9@:%_\+.~#?&//=]*)?`)
var rEx *regexp.Regexp

func isSpam(str string) bool {
	if rEx == nil && os.Getenv("SPAM_LINKS_EXCLUDE") != "" {
		rEx = regexp.MustCompile(os.Getenv("SPAM_LINKS_EXCLUDE"))
	}

	founded := re.FindAllString(str, -1)

	if len(founded) != 0 {
		if os.Getenv("SPAM_LINKS_EXCLUDE") != "" {
			exclude := rEx.FindAllString(str, -1)

			// Do not exclude if there more than one link
			if len(exclude) != 0 && len(founded) == 1 {
				return false
			}
		}

		return true
	}

	return false
}
