package middleware

import (
	"gochat/internal/shared/collections"
	"strings"

	"github.com/gin-gonic/gin"
)

func SkipMiddleware(paths []string, middleware gin.HandlerFunc) gin.HandlerFunc {
	trie := buildTrie(paths)

	matches := func(path string) bool {
		if path == "" {
			return false
		}
		return trie.Match(path)
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path
		full := c.FullPath()
		if matches(path) || matches(full) {
			c.Next()
			return
		}
		middleware(c)
	}
}

func buildTrie(paths []string) *collections.Trie {
	trie := collections.NewTrie()
	for _, p := range paths {
		if strings.HasSuffix(p, "*any") {
			trie.Insert(strings.TrimSuffix(p, "*any"), true)
			continue
		}
		trie.Insert(p, false)
	}
	return trie
}
