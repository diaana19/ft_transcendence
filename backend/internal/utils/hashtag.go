package utils

import (
	"regexp"
	"strings"
)

var hashtagRe = regexp.MustCompile(`#\w+`)

// NormalizeHashtag puts the tag in lower case and adds the # when it is missing.
func NormalizeHashtag(raw string) string {
	raw = strings.ToLower(strings.TrimSpace(raw))
	if !strings.HasPrefix(raw, "#") {
		raw = "#" + raw
	}
	return hashtagRe.FindString(raw)
}

// ExtractHashtags returns all the unique hashtags found in the content.
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
