package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func SkipMiddleware(paths []string, middleware gin.HandlerFunc) gin.HandlerFunc {
	skipPaths := make(map[string]struct{}, len(paths))
	skipPrefixes := make([]string, 0)
	for _, p := range paths {
		if strings.HasSuffix(p, "*any") {
			skipPrefixes = append(skipPrefixes, strings.TrimSuffix(p, "*any"))
			continue
		}
		skipPaths[p] = struct{}{}
	}

	matches := func(path, fullPath string) bool {
		if _, ok := skipPaths[path]; ok {
			return true
		}
		if fullPath != "" {
			if _, ok := skipPaths[fullPath]; ok {
				return true
			}
		}
		for _, prefix := range skipPrefixes {
			if strings.HasPrefix(path, prefix) || (fullPath != "" && strings.HasPrefix(fullPath, prefix)) {
				return true
			}
		}
		return false
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path
		full := c.FullPath()
		if matches(path, full) {
			c.Next()
			return
		}
		middleware(c)
	}
}
