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
	version   string
}

// New creates a new Renderer with embedded templates.
func New() (*Renderer, error) {
	funcMap := template.FuncMap{
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
	}
	tmpl, err := template.New("").Funcs(funcMap).ParseFS(templateFS, "templates/*.html")
	if err != nil {
		return nil, err
	}
	return &Renderer{templates: tmpl}, nil
}

// SetVersion sets the version string to be included in rendered templates.
func (r *Renderer) SetVersion(version string) {
	r.version = version
}

// templateData holds data for template rendering across all page types.
type templateData struct {
	Site         model.Site
	PageTitle    string
	CanonicalURL string
	Summary      string
	IsPost       bool
	OGImage      string
	Version      string
	PageType     string
	Content      template.HTML
	Page         struct {
		Title             string
		Content           template.HTML
		DateAnchors       []string
		CurrentMonthDates []string
		ArchivedYears     []model.YearGroup
		Topics            []model.Topic
	}
}

// blogPostItem represents a post in the blog list.
type blogPostItem struct {
	Title         string
	Slug          string
	DateFormatted string
	WordCount     int
}

// blogListData holds data for blog list template rendering.
type blogListData struct {
	Site         model.Site
	PageTitle    string
	CanonicalURL string
	Summary      string
	IsPost       bool
	OGImage      string
	Version      string
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
	Version       string
	DatePublished string
	Post          struct {
		Title         string
		DateFormatted string
		Content       template.HTML
		WordCount     int
	}
}

// RenderBase renders the base template with site data and content.
func (r *Renderer) RenderBase(site model.Site, content string) (string, error) {
	data := templateData{
		Site:     site,
		Content:  template.HTML(content),
		Version:  r.version,
		PageType: "base",
	}

	var buf bytes.Buffer
	if err := r.templates.ExecuteTemplate(&buf, "base.html", data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// RenderPage renders a static page with the page template.
func (r *Renderer) RenderPage(site model.Site, page model.Page) (string, error) {
	data := templateData{
		Site:         site,
		PageTitle:    page.Title,
		CanonicalURL: site.BaseURL + "/" + page.Slug + "/",
		Summary:      site.Description,
		IsPost:       false,
		OGImage:      ogImageURL(site),
		Version:      r.version,
		PageType:     "page",
	}
	data.Page.Title = page.Title
	data.Page.Content = template.HTML(page.Content)
	data.Page.DateAnchors = page.DateAnchors
	data.Page.CurrentMonthDates = page.CurrentMonthDates
	data.Page.ArchivedYears = page.ArchivedYears
	data.Page.Topics = page.Topics

	var buf bytes.Buffer
	if err := r.templates.ExecuteTemplate(&buf, "base.html", data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// ogImageURL returns the absolute URL for OG image, preferring OGImage with Logo fallback.
func ogImageURL(site model.Site) string {
	if site.OGImage != "" {
		return site.BaseURL + site.OGImage
	}
	if site.Logo != "" {
		return site.BaseURL + site.Logo
	}
	return ""
}

// RenderBlogList renders the blog listing page with all posts.
func (r *Renderer) RenderBlogList(site model.Site, posts []model.Post) (string, error) {
	items := make([]blogPostItem, len(posts))
	for i, p := range posts {
		items[i] = blogPostItem{
			Title:         p.Title,
			Slug:          p.Slug,
			DateFormatted: p.Date.Format("2006-01-02"),
			WordCount:     p.WordCount,
		}
	}

	data := blogListData{
		Site:         site,
		PageTitle:    "Blog",
		CanonicalURL: site.BaseURL + "/blog/",
		Summary:      site.Description,
		IsPost:       false,
		OGImage:      ogImageURL(site),
		Version:      r.version,
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
		Version:       r.version,
		DatePublished: post.Date.Format("2006-01-02"),
	}
	data.Post.Title = post.Title
	data.Post.DateFormatted = post.Date.Format("2006-01-02")
	data.Post.Content = template.HTML(post.Content)
	data.Post.WordCount = post.WordCount

	var buf bytes.Buffer
	if err := r.templates.ExecuteTemplate(&buf, "blog_post.html", data); err != nil {
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
	Version      string
}

// Render404 renders the 404 error page.
func (r *Renderer) Render404(site model.Site) (string, error) {
	data := notFoundData{
		Site:         site,
		PageTitle:    "Page Not Found",
		CanonicalURL: site.BaseURL + "/404/",
		Summary:      site.Description,
		Version:      r.version,
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

// RenderFeed renders an RSS 2.0 feed from FeedItems.
func (r *Renderer) RenderFeed(site model.Site, items []model.FeedItem) (string, error) {
	if len(items) == 0 {
		return "", nil
	}

	// Limit to 20 items
	feedItems := items
	if len(feedItems) > 20 {
		feedItems = feedItems[:20]
	}

	// Build RSS items
	rssItems := make([]rssItem, len(feedItems))
	for i, item := range feedItems {
		rssItems[i] = rssItem{
			Title:       item.FeedTitle(),
			Link:        item.FeedLink(),
			Description: rssCDATA{Content: item.FeedContent()},
			PubDate:     item.FeedDate().Format(time.RFC1123Z),
		}
	}

	feed := rssFeed{
		Version: "2.0",
		Channel: rssChannel{
			Title:         site.Title,
			Link:          site.BaseURL,
			Description:   site.Description,
			LastBuildDate: time.Now().UTC().Format(time.RFC1123Z),
			Items:         rssItems,
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
