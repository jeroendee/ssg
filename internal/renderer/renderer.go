package renderer

import (
	"bytes"
	"embed"
	"html/template"

	"github.com/jeroendee/ssg/internal/model"
)

//go:embed templates/*.html
var templateFS embed.FS

// Renderer renders HTML templates with site data.
type Renderer struct {
	templates *template.Template
}

// New creates a new Renderer with embedded templates.
func New() (*Renderer, error) {
	tmpl, err := template.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		return nil, err
	}
	return &Renderer{templates: tmpl}, nil
}

// baseData holds data for base template rendering.
type baseData struct {
	Site    model.Site
	Content template.HTML
}

// pageData holds data for page template rendering.
type pageData struct {
	Site model.Site
	Page struct {
		Title   string
		Content template.HTML
	}
}

// blogPostItem represents a post in the blog list.
type blogPostItem struct {
	Title         string
	Slug          string
	DateFormatted string
}

// blogListData holds data for blog list template rendering.
type blogListData struct {
	Site  model.Site
	Posts []blogPostItem
}

// blogPostData holds data for blog post template rendering.
type blogPostData struct {
	Site model.Site
	Post struct {
		Title         string
		DateFormatted string
		Content       template.HTML
	}
}

// RenderBase renders the base template with site data and content.
func (r *Renderer) RenderBase(site model.Site, content string) (string, error) {
	data := baseData{
		Site:    site,
		Content: template.HTML(content),
	}

	var buf bytes.Buffer
	if err := r.templates.ExecuteTemplate(&buf, "base.html", data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// RenderPage renders a static page with the page template.
func (r *Renderer) RenderPage(site model.Site, page model.Page) (string, error) {
	data := pageData{
		Site: site,
	}
	data.Page.Title = page.Title
	data.Page.Content = template.HTML(page.Content)

	var buf bytes.Buffer
	if err := r.templates.ExecuteTemplate(&buf, "page.html", data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// RenderBlogList renders the blog listing page with all posts.
func (r *Renderer) RenderBlogList(site model.Site, posts []model.Post) (string, error) {
	items := make([]blogPostItem, len(posts))
	for i, p := range posts {
		items[i] = blogPostItem{
			Title:         p.Title,
			Slug:          p.Slug,
			DateFormatted: p.Date.Format("2006-01-02"),
		}
	}

	data := blogListData{
		Site:  site,
		Posts: items,
	}

	var buf bytes.Buffer
	if err := r.templates.ExecuteTemplate(&buf, "blog_list.html", data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// RenderBlogPost renders a single blog post.
func (r *Renderer) RenderBlogPost(site model.Site, post model.Post) (string, error) {
	data := blogPostData{
		Site: site,
	}
	data.Post.Title = post.Title
	data.Post.DateFormatted = post.Date.Format("2006-01-02")
	data.Post.Content = template.HTML(post.Content)

	var buf bytes.Buffer
	if err := r.templates.ExecuteTemplate(&buf, "blog_post.html", data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// homeData holds data for home page template rendering.
type homeData struct {
	Site model.Site
}

// RenderHome renders the homepage.
func (r *Renderer) RenderHome(site model.Site) (string, error) {
	data := homeData{
		Site: site,
	}

	var buf bytes.Buffer
	if err := r.templates.ExecuteTemplate(&buf, "home.html", data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// notFoundData holds data for 404 page template rendering.
type notFoundData struct {
	Site model.Site
}

// Render404 renders the 404 error page.
func (r *Renderer) Render404(site model.Site) (string, error) {
	data := notFoundData{
		Site: site,
	}

	var buf bytes.Buffer
	if err := r.templates.ExecuteTemplate(&buf, "404.html", data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
