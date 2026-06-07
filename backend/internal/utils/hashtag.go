package utils

import (
	"regexp"
	"strings"
)

// hashtagRe matches a hashtag the same way the frontend does (#\w+), where \w is
// ASCII [0-9A-Za-z_]. Kept in sync with frontend/src/components/common/RichText.jsx.
var hashtagRe = regexp.MustCompile(`#\w+`)

// ExtractHashtags returns the distinct, lowercased hashtags (with the leading
// '#') found in content, preserving first-seen order. Lowercasing makes the
// trend count case-insensitive so #Golang and #golang are the same tag.
func ExtractHashtags(content string) []string {
	matches := hashtagRe.FindAllString(content, -1)
	if len(matches) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(matches))
	tags := make([]string, 0, len(matches))
	for _, m := range matches {
		tag := strings.ToLower(m)
		if _, dup := seen[tag]; dup {
			continue
		}
		seen[tag] = struct{}{}
		tags = append(tags, tag)
	}
	return tags
}
