/*
** File: constants.go
** Description: Shared sort/order tokens and SQL LIKE helpers for repository
** search queries.
 */

package repositories

import "strings"

const (
	sortCreatedAt = "created_at"
	orderAsc      = "asc"
	orderDesc     = "desc"
)

// likeWildcardReplacer escapes the SQL LIKE / ILIKE metacharacters so that
// caller-supplied text is matched literally. PostgreSQL's default escape
// character is `\` (when standard_conforming_strings is on, which is the
// default since PostgreSQL 9.1). Order matters — escape `\` first so the
// inserted backslashes from the next two rules are not double-escaped.
var likeWildcardReplacer = strings.NewReplacer(
	`\`, `\\`,
	`%`, `\%`,
	`_`, `\_`,
)

// escapeLike returns q wrapped as a contains-pattern (`%q%`) with the
// wildcards in q escaped so it matches literally.
func escapeLike(q string) string {
	return "%" + likeWildcardReplacer.Replace(q) + "%"
}
