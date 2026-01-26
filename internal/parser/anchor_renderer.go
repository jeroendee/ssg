package parser

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

// anchorHeadingRenderer wraps heading content in anchor links for direct linking.
type anchorHeadingRenderer struct {
	html.Config
}

// newAnchorHeadingRenderer returns a new renderer that wraps heading content in anchors.
func newAnchorHeadingRenderer(opts ...html.Option) renderer.NodeRenderer {
	r := &anchorHeadingRenderer{
		Config: html.NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

// RegisterFuncs registers the heading renderer function.
func (r *anchorHeadingRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindHeading, r.renderHeading)
}

// renderHeading renders a heading with an anchor link wrapping the content.
func (r *anchorHeadingRenderer) renderHeading(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Heading)
	if entering {
		_, _ = w.WriteString("<h")
		_ = w.WriteByte("123456"[n.Level-1])
		if id, ok := n.AttributeString("id"); ok {
			idBytes := attrToBytes(id)
			_, _ = w.WriteString(` id="`)
			_, _ = w.Write(util.EscapeHTML(idBytes))
			_, _ = w.WriteString(`">`)
			_, _ = w.WriteString(`<a href="#`)
			_, _ = w.Write(util.EscapeHTML(idBytes))
			_, _ = w.WriteString(`">`)
		} else {
			_, _ = w.WriteString(">")
		}
	} else {
		if _, ok := n.AttributeString("id"); ok {
			_, _ = w.WriteString("</a>")
		}
		_, _ = w.WriteString("</h")
		_ = w.WriteByte("123456"[n.Level-1])
		_, _ = w.WriteString(">\n")
	}
	return ast.WalkContinue, nil
}

// attrToBytes converts an attribute value to bytes (handles string or []byte).
func attrToBytes(v any) []byte {
	switch val := v.(type) {
	case string:
		return []byte(val)
	case []byte:
		return val
	default:
		return nil
	}
}
