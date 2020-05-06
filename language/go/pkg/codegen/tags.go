package codegen

import (
	"fmt"
	"sort"
	"strings"

	"github.com/anz-bank/sysl/unsorted/codegen/golang"
)

// TODO: Make this more robust (e.g., What about tags containing newlines?)
func Tag(tags map[string]string) *golang.BasicLit {
	keys := make([]string, 0, len(tags))
	for key := range tags {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf(`%s:%q`, key, tags[key]))
	}
	tag := strings.Join(parts, " ")
	if strings.Contains(tag, "`") {
		return golang.String(tag)
	}
	return &golang.BasicLit{Token: golang.Token{Text: "`" + tag + "`"}}
}
