package assets

import _ "embed"

//go:embed default_style.css
var defaultStyleCSS []byte

// DefaultStyleCSS returns the embedded default CSS stylesheet as a byte slice.
// The stylesheet includes BearBlog-style responsive design with Solarized color scheme.
func DefaultStyleCSS() []byte {
	return defaultStyleCSS
}
