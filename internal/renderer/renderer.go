package renderer

import (
	"bytes"
	"embed"
	"encoding/xml"
	"html/template"
	"time"

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
	Site         model.Site
	PageTitle    string
	CanonicalURL string
	Summary      string
	IsPost       bool
	OGImage      string
	Content      template.HTML
}

// pageData holds data for page template rendering.
type pageData struct {
	Site         model.Site
	PageTitle    string
	CanonicalURL string
	Summary      string
	IsPost       bool
	OGImage      string
	Page         struct {
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
	Site         model.Site
	PageTitle    string
	CanonicalURL string
	Summary      string
	IsPost       bool
	OGImage      string
	Posts        []blogPostItem
}

// blogPostData holds data for blog post template rendering.
type blogPostData struct {
	Site          model.Site
	PageTitle     string
	CanonicalURL  string
	Summary       string
	IsPost        bool
	OGImage       string
	DatePublished string
	Post          struct {
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
		Site:         site,
		PageTitle:    page.Title,
		CanonicalURL: site.BaseURL + "/" + page.Slug + "/",
		Summary:      site.Description,
		IsPost:       false,
		OGImage:      ogImageURL(site),
	}
	data.Page.Title = page.Title
	data.Page.Content = template.HTML(page.Content)

	var buf bytes.Buffer
	if err := r.templates.ExecuteTemplate(&buf, "page.html", data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// ogImageURL returns the absolute URL for OG image, using site logo as fallback.
func ogImageURL(site model.Site) string {
	if site.Logo == "" {
		return ""
	}
	return site.BaseURL + site.Logo
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
		Site:         site,
		PageTitle:    "Blog",
		CanonicalURL: site.BaseURL + "/blog/",
		Summary:      site.Description,
		IsPost:       false,
		OGImage:      ogImageURL(site),
		Posts:        items,
	}

	var buf bytes.Buffer
	if err := r.templates.ExecuteTemplate(&buf, "blog_list.html", data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// RenderBlogPost renders a single blog post.
func (r *Renderer) RenderBlogPost(site model.Site, post model.Post) (string, error) {
	summary := post.Summary
	if summary == "" {
		summary = site.Description
	}
	data := blogPostData{
		Site:          site,
		PageTitle:     post.Title,
		CanonicalURL:  site.BaseURL + "/blog/" + post.Slug + "/",
		Summary:       summary,
		IsPost:        true,
		OGImage:       ogImageURL(site),
		DatePublished: post.Date.Format("2006-01-02"),
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
	Site         model.Site
	PageTitle    string
	CanonicalURL string
	Summary      string
	IsPost       bool
	OGImage      string
}

// RenderHome renders the homepage.
func (r *Renderer) RenderHome(site model.Site) (string, error) {
	data := homeData{
		Site:         site,
		CanonicalURL: site.BaseURL + "/",
		Summary:      site.Description,
		IsPost:       false,
		OGImage:      ogImageURL(site),
	}

	var buf bytes.Buffer
	if err := r.templates.ExecuteTemplate(&buf, "home.html", data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// notFoundData holds data for 404 page template rendering.
type notFoundData struct {
	Site         model.Site
	PageTitle    string
	CanonicalURL string
	Summary      string
	IsPost       bool
	OGImage      string
}

// Render404 renders the 404 error page.
func (r *Renderer) Render404(site model.Site) (string, error) {
	data := notFoundData{
		Site:         site,
		PageTitle:    "Page Not Found",
		CanonicalURL: site.BaseURL + "/404/",
		Summary:      site.Description,
	}

	var buf bytes.Buffer
	if err := r.templates.ExecuteTemplate(&buf, "404.html", data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// rssChannel represents the RSS channel element.
type rssChannel struct {
	Title         string    `xml:"title"`
	Link          string    `xml:"link"`
	Description   string    `xml:"description"`
	LastBuildDate string    `xml:"lastBuildDate"`
	Items         []rssItem `xml:"item"`
}

// rssItem represents an RSS item element.
type rssItem struct {
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description rssCDATA `xml:"description"`
	PubDate     string   `xml:"pubDate"`
}

// rssCDATA wraps content in CDATA.
type rssCDATA struct {
	Content string `xml:",cdata"`
}

// rssFeed represents the root RSS element.
type rssFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel rssChannel `xml:"channel"`
}

// RenderFeed renders an RSS 2.0 feed for blog posts.
func (r *Renderer) RenderFeed(site model.Site, posts []model.Post) (string, error) {
	if len(posts) == 0 {
		return "", nil
	}

	// Limit to 20 posts
	feedPosts := posts
	if len(feedPosts) > 20 {
		feedPosts = feedPosts[:20]
	}

	// Build items
	items := make([]rssItem, len(feedPosts))
	for i, p := range feedPosts {
		items[i] = rssItem{
			Title:       p.Title,
			Link:        site.BaseURL + "/blog/" + p.Slug + "/",
			Description: rssCDATA{Content: p.Content},
			PubDate:     p.Date.Format(time.RFC1123Z),
		}
	}

	feed := rssFeed{
		Version: "2.0",
		Channel: rssChannel{
			Title:         site.Title,
			Link:          site.BaseURL,
			Description:   site.Description,
			LastBuildDate: time.Now().UTC().Format(time.RFC1123Z),
			Items:         items,
		},
	}

	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")

	enc := xml.NewEncoder(&buf)
	enc.Indent("", "  ")
	if err := enc.Encode(feed); err != nil {
		return "", err
	}

	return buf.String(), nil
}
