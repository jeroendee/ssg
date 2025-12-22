package assets_test

import (
	"strings"
	"testing"

	"github.com/jeroendee/ssg/internal/assets"
)

func TestDefaultStyleCSS(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T, css []byte)
	}{
		{
			name: "returns non-empty content",
			test: func(t *testing.T, css []byte) {
				if len(css) == 0 {
					t.Error("DefaultStyleCSS() returned empty content")
				}
			},
		},
		{
			name: "contains valid CSS markers",
			test: func(t *testing.T, css []byte) {
				content := string(css)
				requiredMarkers := []string{
					":root",
					"--bg-primary",
					"--text-primary",
					"body {",
					"@media (prefers-color-scheme: dark)",
				}
				for _, marker := range requiredMarkers {
					if !strings.Contains(content, marker) {
						t.Errorf("DefaultStyleCSS() missing expected CSS marker: %q", marker)
					}
				}
			},
		},
		{
			name: "contains BearBlog-style comment",
			test: func(t *testing.T, css []byte) {
				content := string(css)
				if !strings.Contains(content, "BearBlog-style CSS") {
					t.Error("DefaultStyleCSS() missing expected header comment")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			css := assets.DefaultStyleCSS()
			tt.test(t, css)
		})
	}
}
