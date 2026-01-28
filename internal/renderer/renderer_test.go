package renderer

import (
	"strings"
	"testing"
	"time"

	"github.com/jeroendee/ssg/internal/model"
)

func TestRenderPage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		site        model.Site
		page        model.Page
		wantTitle   string
		wantContent string
	}{
		{
			name: "renders page with title and content",
			site: model.Site{
				Title:   "Quality Shepherd",
				BaseURL: "https://www.qualityshepherd.nl",
				Navigation: []model.NavItem{
					{Title: "Home", URL: "/"},
					{Title: "About", URL: "/about/"},
				},
			},
			page: model.Page{
				Title:   "About",
				Slug:    "about",
				Content: "<p>Hello! My name is Jeroen.</p>",
			},
			wantTitle:   "About",
			wantContent: "<p>Hello! My name is Jeroen.</p>",
		},
		{
			name: "renders page with special characters in content",
			site: model.Site{
				Title: "Test Site",
			},
			page: model.Page{
				Title:   "Contact",
				Slug:    "contact",
				Content: "<p>Email: test@example.com</p>",
			},
			wantTitle:   "Contact",
			wantContent: "<p>Email: test@example.com</p>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r, err := New()
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}

			got, err := r.RenderPage(tt.site, tt.page)
			if err != nil {
				t.Fatalf("RenderPage() error = %v", err)
			}

			// Check page title appears in output
			if !strings.Contains(got, tt.wantTitle) {
				t.Errorf("RenderPage() missing title %q", tt.wantTitle)
			}

			// Check content is rendered
			if !strings.Contains(got, tt.wantContent) {
				t.Errorf("RenderPage() missing content %q", tt.wantContent)
			}

			// Should have base structure
			if !strings.Contains(got, "<!DOCTYPE html>") {
				t.Error("RenderPage() missing DOCTYPE")
			}
		})
	}
}

func TestRenderBlogList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		site           model.Site
		posts          []model.Post
		wantDates      []string
		wantLinks      []string
		wantWordCounts []string
	}{
		{
			name: "renders blog list with posts",
			site: model.Site{
				Title: "Quality Shepherd",
				Navigation: []model.NavItem{
					{Title: "Blog", URL: "/blog/"},
				},
			},
			posts: []model.Post{
				{
					Page: model.Page{
						Title: "First Post",
						Slug:  "first-post",
					},
					Date: time.Date(2021, 3, 26, 0, 0, 0, 0, time.UTC),
				},
				{
					Page: model.Page{
						Title: "Second Post",
						Slug:  "second-post",
					},
					Date: time.Date(2021, 4, 15, 0, 0, 0, 0, time.UTC),
				},
			},
			wantDates:      []string{"2021-03-26", "2021-04-15"},
			wantLinks:      []string{"First Post", "Second Post"},
			wantWordCounts: []string{},
		},
		{
			name: "renders empty blog list",
			site: model.Site{
				Title: "Test Site",
			},
			posts:          []model.Post{},
			wantDates:      []string{},
			wantLinks:      []string{},
			wantWordCounts: []string{},
		},
		{
			name: "renders blog list with word counts",
			site: model.Site{
				Title: "Test Site",
			},
			posts: []model.Post{
				{
					Page: model.Page{
						Title: "Post with word count",
						Slug:  "post-with-word-count",
					},
					Date:      time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
					WordCount: 150,
				},
				{
					Page: model.Page{
						Title: "Another post",
						Slug:  "another-post",
					},
					Date:      time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC),
					WordCount: 250,
				},
			},
			wantDates:      []string{"2024-01-15", "2024-01-16"},
			wantLinks:      []string{"Post with word count", "Another post"},
			wantWordCounts: []string{"150 words", "250 words"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r, err := New()
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}

			got, err := r.RenderBlogList(tt.site, tt.posts)
			if err != nil {
				t.Fatalf("RenderBlogList() error = %v", err)
			}

			// Check dates are rendered
			for _, date := range tt.wantDates {
				if !strings.Contains(got, date) {
					t.Errorf("RenderBlogList() missing date %q", date)
				}
			}

			// Check post titles/links are rendered
			for _, link := range tt.wantLinks {
				if !strings.Contains(got, link) {
					t.Errorf("RenderBlogList() missing link %q", link)
				}
			}

			// Check word counts are rendered
			for _, wordCount := range tt.wantWordCounts {
				if !strings.Contains(got, wordCount) {
					t.Errorf("RenderBlogList() missing word count %q", wordCount)
				}
			}

			// Should have base structure
			if !strings.Contains(got, "<!DOCTYPE html>") {
				t.Error("RenderBlogList() missing DOCTYPE")
			}
		})
	}
}

func TestRenderBlogPost(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		site        model.Site
		post        model.Post
		wantTitle   string
		wantDate    string
		wantContent string
	}{
		{
			name: "renders blog post with title, date, and content",
			site: model.Site{
				Title: "Quality Shepherd",
				Navigation: []model.NavItem{
					{Title: "Blog", URL: "/blog/"},
				},
			},
			post: model.Post{
				Page: model.Page{
					Title:   "My First Post",
					Slug:    "my-first-post",
					Content: "<p>This is the content of my first post.</p>",
				},
				Date: time.Date(2021, 3, 26, 0, 0, 0, 0, time.UTC),
			},
			wantTitle:   "My First Post",
			wantDate:    "2021-03-26",
			wantContent: "<p>This is the content of my first post.</p>",
		},
		{
			name: "renders post with HTML content",
			site: model.Site{
				Title: "Test Blog",
			},
			post: model.Post{
				Page: model.Page{
					Title:   "Code Post",
					Slug:    "code-post",
					Content: "<pre><code>func main() {}</code></pre>",
				},
				Date: time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC),
			},
			wantTitle:   "Code Post",
			wantDate:    "2023-12-01",
			wantContent: "<pre><code>func main() {}</code></pre>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r, err := New()
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}

			got, err := r.RenderBlogPost(tt.site, tt.post)
			if err != nil {
				t.Fatalf("RenderBlogPost() error = %v", err)
			}

			// Check post title
			if !strings.Contains(got, tt.wantTitle) {
				t.Errorf("RenderBlogPost() missing title %q", tt.wantTitle)
			}

			// Check date is rendered
			if !strings.Contains(got, tt.wantDate) {
				t.Errorf("RenderBlogPost() missing date %q", tt.wantDate)
			}

			// Check content is rendered
			if !strings.Contains(got, tt.wantContent) {
				t.Errorf("RenderBlogPost() missing content %q", tt.wantContent)
			}

			// Should have base structure
			if !strings.Contains(got, "<!DOCTYPE html>") {
				t.Error("RenderBlogPost() missing DOCTYPE")
			}
		})
	}
}

func TestRenderBase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		site       model.Site
		wantTitle  string
		wantNav    []string
		wantHeader bool
		wantFooter bool
	}{
		{
			name: "renders site with navigation",
			site: model.Site{
				Title:   "Quality Shepherd",
				BaseURL: "https://www.qualityshepherd.nl",
				Author:  "Jeroen",
				Navigation: []model.NavItem{
					{Title: "Home", URL: "/"},
					{Title: "About", URL: "/about/"},
					{Title: "Blog", URL: "/blog/"},
				},
			},
			wantTitle:  "Quality Shepherd",
			wantNav:    []string{"Home", "About", "Blog"},
			wantHeader: true,
			wantFooter: true,
		},
		{
			name: "renders site with empty navigation",
			site: model.Site{
				Title:      "Test Site",
				BaseURL:    "https://example.com",
				Author:     "Test",
				Navigation: []model.NavItem{},
			},
			wantTitle:  "Test Site",
			wantNav:    []string{},
			wantHeader: true,
			wantFooter: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r, err := New()
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}

			got, err := r.RenderBase(tt.site, "Test Content")
			if err != nil {
				t.Fatalf("RenderBase() error = %v", err)
			}

			// Check title in output
			if !strings.Contains(got, tt.wantTitle) {
				t.Errorf("RenderBase() missing title %q", tt.wantTitle)
			}

			// Check navigation items
			for _, nav := range tt.wantNav {
				if !strings.Contains(got, nav) {
					t.Errorf("RenderBase() missing nav item %q", nav)
				}
			}

			// Check header presence
			if tt.wantHeader && !strings.Contains(got, "<header") {
				t.Error("RenderBase() missing header element")
			}

			// Check footer presence
			if tt.wantFooter && !strings.Contains(got, "<footer") {
				t.Error("RenderBase() missing footer element")
			}

			// Check valid HTML structure
			if !strings.Contains(got, "<!DOCTYPE html>") {
				t.Error("RenderBase() missing DOCTYPE")
			}
			if !strings.Contains(got, "<html") {
				t.Error("RenderBase() missing html element")
			}
		})
	}
}

func TestRender404(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		site      model.Site
		wantTitle string
		want404   bool
	}{
		{
			name: "renders 404 page with site title",
			site: model.Site{
				Title:   "Quality Shepherd",
				BaseURL: "https://www.qualityshepherd.nl",
				Navigation: []model.NavItem{
					{Title: "Home", URL: "/"},
				},
			},
			wantTitle: "Quality Shepherd",
			want404:   true,
		},
		{
			name: "renders 404 page with not found message",
			site: model.Site{
				Title: "Test Site",
			},
			wantTitle: "Test Site",
			want404:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r, err := New()
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}

			got, err := r.Render404(tt.site)
			if err != nil {
				t.Fatalf("Render404() error = %v", err)
			}

			// Check site title in output
			if !strings.Contains(got, tt.wantTitle) {
				t.Errorf("Render404() missing title %q", tt.wantTitle)
			}

			// Check 404 indicator
			if tt.want404 && !strings.Contains(got, "404") {
				t.Error("Render404() missing 404 indicator")
			}

			// Check valid HTML structure
			if !strings.Contains(got, "<!DOCTYPE html>") {
				t.Error("Render404() missing DOCTYPE")
			}
		})
	}
}

func TestRenderBase_WithLogo(t *testing.T) {
	t.Parallel()

	r, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	site := model.Site{
		Title: "Quality Shepherd",
		Logo:  "/logo.svg",
	}

	got, err := r.RenderBase(site, "Test Content")
	if err != nil {
		t.Fatalf("RenderBase() error = %v", err)
	}

	// Check logo image is rendered
	if !strings.Contains(got, `<img src="/logo.svg"`) {
		t.Error("RenderBase() should render logo image when Logo is set")
	}

	// Check alt attribute contains site title
	if !strings.Contains(got, `alt="Quality Shepherd`) {
		t.Error("RenderBase() logo should have alt attribute with site title")
	}
}

func TestRenderBase_WithoutLogo(t *testing.T) {
	t.Parallel()

	r, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	site := model.Site{
		Title: "Quality Shepherd",
		Logo:  "",
	}

	got, err := r.RenderBase(site, "Test Content")
	if err != nil {
		t.Fatalf("RenderBase() error = %v", err)
	}

	// Check logo image is NOT rendered
	if strings.Contains(got, `<img src="`) {
		t.Error("RenderBase() should not render logo image when Logo is empty")
	}
}

func TestRenderBase_WithFavicon(t *testing.T) {
	t.Parallel()

	r, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	site := model.Site{
		Title:   "Quality Shepherd",
		Favicon: "/favicon.svg",
	}

	got, err := r.RenderBase(site, "Test Content")
	if err != nil {
		t.Fatalf("RenderBase() error = %v", err)
	}

	// Check favicon link is rendered (+ is HTML-escaped to &#43; by html/template)
	if !strings.Contains(got, `<link rel="icon" href="/favicon.svg" type="image/svg&#43;xml">`) {
		t.Errorf("RenderBase() should render favicon link when Favicon is set, got:\n%s", got)
	}
}

func TestRenderBlogPost_WordCount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		site          model.Site
		post          model.Post
		wantWordCount string
	}{
		{
			name: "renders word count in blog post",
			site: model.Site{
				Title: "Test Site",
			},
			post: model.Post{
				Page: model.Page{
					Title:   "Test Post",
					Slug:    "test-post",
					Content: "<p>Test content</p>",
				},
				Date:      time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
				WordCount: 150,
			},
			wantWordCount: "150 words",
		},
		{
			name: "renders zero word count",
			site: model.Site{
				Title: "Test Site",
			},
			post: model.Post{
				Page: model.Page{
					Title:   "Empty Post",
					Slug:    "empty-post",
					Content: "",
				},
				Date:      time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
				WordCount: 0,
			},
			wantWordCount: "0 words",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r, err := New()
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}

			got, err := r.RenderBlogPost(tt.site, tt.post)
			if err != nil {
				t.Fatalf("RenderBlogPost() error = %v", err)
			}

			if !strings.Contains(got, tt.wantWordCount) {
				t.Errorf("RenderBlogPost() missing word count %q in output", tt.wantWordCount)
			}
		})
	}
}

func TestRenderBase_WithoutFavicon(t *testing.T) {
	t.Parallel()

	r, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	site := model.Site{
		Title:   "Quality Shepherd",
		Favicon: "",
	}

	got, err := r.RenderBase(site, "Test Content")
	if err != nil {
		t.Fatalf("RenderBase() error = %v", err)
	}

	// Check favicon link is NOT rendered
	if strings.Contains(got, `rel="icon"`) {
		t.Error("RenderBase() should not render favicon link when Favicon is empty")
	}
}

func TestRenderFeed_EmptyItems(t *testing.T) {
	t.Parallel()

	r, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	site := model.Site{
		Title:       "Test Site",
		Description: "A test site",
		BaseURL:     "https://example.com",
	}

	got, err := r.RenderFeed(site, nil)
	if err != nil {
		t.Fatalf("RenderFeed() error = %v", err)
	}

	if got != "" {
		t.Errorf("RenderFeed() with no items = %q, want empty string", got)
	}
}

// postToFeedItem converts a Post to a FeedItem for testing.
func postToFeedItem(post model.Post, baseURL string) model.FeedItem {
	return model.PostFeedAdapter{Post: &post, BaseURL: baseURL}
}

// postsToFeedItems converts a slice of Posts to FeedItems for testing.
func postsToFeedItems(posts []model.Post, baseURL string) []model.FeedItem {
	items := make([]model.FeedItem, len(posts))
	for i := range posts {
		items[i] = model.PostFeedAdapter{Post: &posts[i], BaseURL: baseURL}
	}
	return items
}

func TestRenderFeed_ValidRSSStructure(t *testing.T) {
	t.Parallel()

	r, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	site := model.Site{
		Title:       "Quality Shepherd",
		Description: "A blog about testing",
		BaseURL:     "https://www.qualityshepherd.nl",
	}
	posts := []model.Post{
		{
			Page: model.Page{
				Title:   "First Post",
				Slug:    "first-post",
				Content: "<p>Hello world</p>",
			},
			Date: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		},
	}

	got, err := r.RenderFeed(site, postsToFeedItems(posts, site.BaseURL))
	if err != nil {
		t.Fatalf("RenderFeed() error = %v", err)
	}

	// Check XML declaration
	if !strings.Contains(got, `<?xml version="1.0" encoding="UTF-8"?>`) {
		t.Error("RenderFeed() missing XML declaration")
	}

	// Check RSS root element
	if !strings.Contains(got, `<rss version="2.0">`) {
		t.Error("RenderFeed() missing RSS root element")
	}

	// Check channel info
	if !strings.Contains(got, "<channel>") {
		t.Error("RenderFeed() missing channel element")
	}
	if !strings.Contains(got, "<title>Quality Shepherd</title>") {
		t.Error("RenderFeed() missing channel title")
	}
	if !strings.Contains(got, "<description>A blog about testing</description>") {
		t.Error("RenderFeed() missing channel description")
	}
	if !strings.Contains(got, "<link>https://www.qualityshepherd.nl</link>") {
		t.Error("RenderFeed() missing channel link")
	}
}

func TestRenderFeed_AbsoluteURLs(t *testing.T) {
	t.Parallel()

	r, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	site := model.Site{
		Title:       "Test Site",
		Description: "A test site",
		BaseURL:     "https://example.com",
	}
	posts := []model.Post{
		{
			Page: model.Page{
				Title:   "Test Post",
				Slug:    "test-post",
				Content: "<p>Content</p>",
			},
			Date: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		},
	}

	got, err := r.RenderFeed(site, postsToFeedItems(posts, site.BaseURL))
	if err != nil {
		t.Fatalf("RenderFeed() error = %v", err)
	}

	// Check absolute URL in item
	if !strings.Contains(got, "<link>https://example.com/blog/test-post/</link>") {
		t.Error("RenderFeed() item should have absolute URL")
	}
}

func TestRenderFeed_RFC822DateFormat(t *testing.T) {
	t.Parallel()

	r, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	site := model.Site{
		Title:       "Test Site",
		Description: "A test site",
		BaseURL:     "https://example.com",
	}
	posts := []model.Post{
		{
			Page: model.Page{
				Title:   "Test Post",
				Slug:    "test-post",
				Content: "<p>Content</p>",
			},
			Date: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		},
	}

	got, err := r.RenderFeed(site, postsToFeedItems(posts, site.BaseURL))
	if err != nil {
		t.Fatalf("RenderFeed() error = %v", err)
	}

	// Check RFC 822 date format in pubDate
	if !strings.Contains(got, "<pubDate>Mon, 15 Jan 2024 10:30:00 +0000</pubDate>") {
		t.Errorf("RenderFeed() should use RFC 822 date format, got:\n%s", got)
	}
}

func TestRenderFeed_CDATAContent(t *testing.T) {
	t.Parallel()

	r, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	site := model.Site{
		Title:       "Test Site",
		Description: "A test site",
		BaseURL:     "https://example.com",
	}
	posts := []model.Post{
		{
			Page: model.Page{
				Title:   "Test Post",
				Slug:    "test-post",
				Content: "<p>HTML content with <strong>tags</strong></p>",
			},
			Date: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		},
	}

	got, err := r.RenderFeed(site, postsToFeedItems(posts, site.BaseURL))
	if err != nil {
		t.Fatalf("RenderFeed() error = %v", err)
	}

	// Check HTML content is wrapped in CDATA
	if !strings.Contains(got, "<![CDATA[<p>HTML content with <strong>tags</strong></p>]]>") {
		t.Errorf("RenderFeed() should wrap HTML content in CDATA, got:\n%s", got)
	}
}

func TestRenderFeed_Max20Items(t *testing.T) {
	t.Parallel()

	r, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	site := model.Site{
		Title:       "Test Site",
		Description: "A test site",
		BaseURL:     "https://example.com",
	}

	// Create 25 posts
	posts := make([]model.Post, 25)
	for i := 0; i < 25; i++ {
		posts[i] = model.Post{
			Page: model.Page{
				Title:   "Post",
				Slug:    "post",
				Content: "<p>Content</p>",
			},
			Date: time.Date(2024, 1, i+1, 0, 0, 0, 0, time.UTC),
		}
	}

	got, err := r.RenderFeed(site, postsToFeedItems(posts, site.BaseURL))
	if err != nil {
		t.Fatalf("RenderFeed() error = %v", err)
	}

	// Count <item> occurrences
	itemCount := strings.Count(got, "<item>")
	if itemCount != 20 {
		t.Errorf("RenderFeed() should limit to 20 items, got %d", itemCount)
	}
}

func TestRenderFeed_MixedItems(t *testing.T) {
	t.Parallel()

	r, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	site := model.Site{
		Title:       "Test Site",
		Description: "A test site",
		BaseURL:     "https://example.com",
	}

	// Create mixed feed items: 1 post and 1 page section
	post := model.Post{
		Page: model.Page{
			Title:   "Blog Post",
			Slug:    "blog-post",
			Content: "<p>Post content</p>",
		},
		Date: time.Date(2026, 1, 27, 0, 0, 0, 0, time.UTC),
	}
	dateSection := model.DateSection{
		PageTitle: "Moments",
		PagePath:  "/moments/",
		Date:      time.Date(2026, 1, 26, 0, 0, 0, 0, time.UTC),
		Anchor:    "2026-01-26",
		Content:   "<p>Section content</p>",
		BaseURL:   "https://example.com",
	}

	items := []model.FeedItem{
		model.PostFeedAdapter{Post: &post, BaseURL: site.BaseURL},
		dateSection,
	}

	got, err := r.RenderFeed(site, items)
	if err != nil {
		t.Fatalf("RenderFeed() error = %v", err)
	}

	// Check post entry
	if !strings.Contains(got, "<title>Blog Post</title>") {
		t.Error("RenderFeed() missing post title")
	}
	if !strings.Contains(got, "<link>https://example.com/blog/blog-post/</link>") {
		t.Error("RenderFeed() missing post link")
	}

	// Check date section entry
	if !strings.Contains(got, "<title>Moments - January 26, 2026</title>") {
		t.Error("RenderFeed() missing date section title")
	}
	if !strings.Contains(got, "<link>https://example.com/moments/#2026-01-26</link>") {
		t.Error("RenderFeed() missing date section link")
	}

	// Check both items present
	itemCount := strings.Count(got, "<item>")
	if itemCount != 2 {
		t.Errorf("RenderFeed() should have 2 items, got %d", itemCount)
	}
}

func TestRenderFeed_DateSectionOnly(t *testing.T) {
	t.Parallel()

	r, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	site := model.Site{
		Title:       "Test Site",
		Description: "A test site",
		BaseURL:     "https://example.com",
	}

	dateSection := model.DateSection{
		PageTitle: "Now",
		PagePath:  "/now/",
		Date:      time.Date(2026, 1, 27, 0, 0, 0, 0, time.UTC),
		Anchor:    "2026-01-27",
		Content:   "<p>What I'm doing now</p>",
		BaseURL:   "https://example.com",
	}

	items := []model.FeedItem{dateSection}

	got, err := r.RenderFeed(site, items)
	if err != nil {
		t.Fatalf("RenderFeed() error = %v", err)
	}

	// Check date section entry
	if !strings.Contains(got, "<title>Now - January 27, 2026</title>") {
		t.Error("RenderFeed() missing date section title")
	}
	if !strings.Contains(got, "<link>https://example.com/now/#2026-01-27</link>") {
		t.Error("RenderFeed() missing date section link")
	}
	if !strings.Contains(got, "<![CDATA[<p>What I'm doing now</p>]]>") {
		t.Error("RenderFeed() missing date section content")
	}
}

func TestRenderPage_CanonicalURL(t *testing.T) {
	t.Parallel()

	r, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	site := model.Site{
		Title:   "Test Site",
		BaseURL: "https://example.com",
	}
	page := model.Page{
		Title: "About",
		Slug:  "about",
	}

	got, err := r.RenderPage(site, page)
	if err != nil {
		t.Fatalf("RenderPage() error = %v", err)
	}

	// Check canonical URL (page = BaseURL + "/" + slug + "/")
	wantCanonical := "https://example.com/about/"
	if !strings.Contains(got, wantCanonical) {
		t.Errorf("RenderPage() should expose canonical URL %q", wantCanonical)
	}
}

func TestRenderBlogList_CanonicalURL(t *testing.T) {
	t.Parallel()

	r, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	site := model.Site{
		Title:   "Test Site",
		BaseURL: "https://example.com",
	}

	got, err := r.RenderBlogList(site, nil)
	if err != nil {
		t.Fatalf("RenderBlogList() error = %v", err)
	}

	// Check canonical URL (blog list = BaseURL + "/blog/")
	wantCanonical := "https://example.com/blog/"
	if !strings.Contains(got, wantCanonical) {
		t.Errorf("RenderBlogList() should expose canonical URL %q", wantCanonical)
	}
}

func TestRenderBlogPost_CanonicalURL(t *testing.T) {
	t.Parallel()

	r, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	site := model.Site{
		Title:   "Test Site",
		BaseURL: "https://example.com",
	}
	post := model.Post{
		Page: model.Page{
			Title: "My Post",
			Slug:  "my-post",
		},
		Date: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
	}

	got, err := r.RenderBlogPost(site, post)
	if err != nil {
		t.Fatalf("RenderBlogPost() error = %v", err)
	}

	// Check canonical URL (blog post = BaseURL + "/blog/" + slug + "/")
	wantCanonical := "https://example.com/blog/my-post/"
	if !strings.Contains(got, wantCanonical) {
		t.Errorf("RenderBlogPost() should expose canonical URL %q", wantCanonical)
	}
}

func TestRender404_CanonicalURL(t *testing.T) {
	t.Parallel()

	r, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	site := model.Site{
		Title:   "Test Site",
		BaseURL: "https://example.com",
	}

	got, err := r.Render404(site)
	if err != nil {
		t.Fatalf("Render404() error = %v", err)
	}

	// Check canonical URL (404 = BaseURL + "/404/")
	wantCanonical := "https://example.com/404/"
	if !strings.Contains(got, wantCanonical) {
		t.Errorf("Render404() should expose canonical URL %q", wantCanonical)
	}
}

func TestRenderBlogPost_MetaDescription(t *testing.T) {
	t.Parallel()

	r, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	site := model.Site{
		Title:       "Test Site",
		Description: "Site description fallback",
		BaseURL:     "https://example.com",
	}
	post := model.Post{
		Page: model.Page{
			Title:   "My Post",
			Slug:    "my-post",
			Content: "<p>Post content</p>",
		},
		Date:    time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		Summary: "This is the post summary",
	}

	got, err := r.RenderBlogPost(site, post)
	if err != nil {
		t.Fatalf("RenderBlogPost() error = %v", err)
	}

	// Blog post should use Summary for meta description
	wantMeta := `<meta name="description" content="This is the post summary">`
	if !strings.Contains(got, wantMeta) {
		t.Errorf("RenderBlogPost() should have meta description with summary, got:\n%s", got)
	}
}

func TestRenderBlogPost_MetaDescription_FallbackToSiteDescription(t *testing.T) {
	t.Parallel()

	r, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	site := model.Site{
		Title:       "Test Site",
		Description: "Site description fallback",
		BaseURL:     "https://example.com",
	}
	post := model.Post{
		Page: model.Page{
			Title:   "My Post",
			Slug:    "my-post",
			Content: "<p>Post content</p>",
		},
		Date:    time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		Summary: "", // No summary
	}

	got, err := r.RenderBlogPost(site, post)
	if err != nil {
		t.Fatalf("RenderBlogPost() error = %v", err)
	}

	// Should fallback to site description when no summary
	wantMeta := `<meta name="description" content="Site description fallback">`
	if !strings.Contains(got, wantMeta) {
		t.Errorf("RenderBlogPost() should fallback to site description, got:\n%s", got)
	}
}

func TestRenderPage_MetaDescription(t *testing.T) {
	t.Parallel()

	r, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	site := model.Site{
		Title:       "Test Site",
		Description: "Site description for pages",
		BaseURL:     "https://example.com",
	}
	page := model.Page{
		Title:   "About",
		Slug:    "about",
		Content: "<p>About content</p>",
	}

	got, err := r.RenderPage(site, page)
	if err != nil {
		t.Fatalf("RenderPage() error = %v", err)
	}

	// Pages should use site description
	wantMeta := `<meta name="description" content="Site description for pages">`
	if !strings.Contains(got, wantMeta) {
		t.Errorf("RenderPage() should have meta description with site description, got:\n%s", got)
	}
}

func TestRenderBlogPost_OpenGraphTags(t *testing.T) {
	t.Parallel()

	r, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	site := model.Site{
		Title:       "Quality Shepherd",
		Description: "A blog about testing",
		BaseURL:     "https://www.qualityshepherd.nl",
		Logo:        "/logo.svg",
	}
	post := model.Post{
		Page: model.Page{
			Title:   "My Test Post",
			Slug:    "my-test-post",
			Content: "<p>Post content</p>",
		},
		Date:    time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		Summary: "This is the post summary",
	}

	got, err := r.RenderBlogPost(site, post)
	if err != nil {
		t.Fatalf("RenderBlogPost() error = %v", err)
	}

	// Check og:title
	if !strings.Contains(got, `<meta property="og:title" content="My Test Post">`) {
		t.Error("RenderBlogPost() missing og:title")
	}

	// Check og:description
	if !strings.Contains(got, `<meta property="og:description" content="This is the post summary">`) {
		t.Error("RenderBlogPost() missing og:description")
	}

	// Check og:url
	if !strings.Contains(got, `<meta property="og:url" content="https://www.qualityshepherd.nl/blog/my-test-post/">`) {
		t.Error("RenderBlogPost() missing og:url")
	}

	// Check og:type (article for blog posts)
	if !strings.Contains(got, `<meta property="og:type" content="article">`) {
		t.Error("RenderBlogPost() should have og:type=article")
	}

	// Check og:site_name
	if !strings.Contains(got, `<meta property="og:site_name" content="Quality Shepherd">`) {
		t.Error("RenderBlogPost() missing og:site_name")
	}

	// Check og:image (uses site logo as fallback)
	if !strings.Contains(got, `<meta property="og:image" content="https://www.qualityshepherd.nl/logo.svg">`) {
		t.Error("RenderBlogPost() missing og:image")
	}
}

func TestRenderPage_OpenGraphTags(t *testing.T) {
	t.Parallel()

	r, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	site := model.Site{
		Title:       "Quality Shepherd",
		Description: "A blog about testing",
		BaseURL:     "https://www.qualityshepherd.nl",
		Logo:        "/logo.svg",
	}
	page := model.Page{
		Title:   "About",
		Slug:    "about",
		Content: "<p>About content</p>",
	}

	got, err := r.RenderPage(site, page)
	if err != nil {
		t.Fatalf("RenderPage() error = %v", err)
	}

	// Check og:title
	if !strings.Contains(got, `<meta property="og:title" content="About">`) {
		t.Error("RenderPage() missing og:title")
	}

	// Check og:type (website for pages)
	if !strings.Contains(got, `<meta property="og:type" content="website">`) {
		t.Error("RenderPage() should have og:type=website")
	}

	// Check og:site_name
	if !strings.Contains(got, `<meta property="og:site_name" content="Quality Shepherd">`) {
		t.Error("RenderPage() missing og:site_name")
	}
}

func TestRenderBlogPost_TwitterCardTags(t *testing.T) {
	t.Parallel()

	r, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	site := model.Site{
		Title:       "Quality Shepherd",
		Description: "A blog about testing",
		BaseURL:     "https://www.qualityshepherd.nl",
		Logo:        "/logo.svg",
	}
	post := model.Post{
		Page: model.Page{
			Title:   "My Test Post",
			Slug:    "my-test-post",
			Content: "<p>Post content</p>",
		},
		Date:    time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		Summary: "This is the post summary",
	}

	got, err := r.RenderBlogPost(site, post)
	if err != nil {
		t.Fatalf("RenderBlogPost() error = %v", err)
	}

	// Check twitter:card
	if !strings.Contains(got, `<meta name="twitter:card" content="summary">`) {
		t.Error("RenderBlogPost() missing twitter:card")
	}

	// Check twitter:title
	if !strings.Contains(got, `<meta name="twitter:title" content="My Test Post">`) {
		t.Error("RenderBlogPost() missing twitter:title")
	}

	// Check twitter:description
	if !strings.Contains(got, `<meta name="twitter:description" content="This is the post summary">`) {
		t.Error("RenderBlogPost() missing twitter:description")
	}

	// Check twitter:image (uses site logo as fallback)
	if !strings.Contains(got, `<meta name="twitter:image" content="https://www.qualityshepherd.nl/logo.svg">`) {
		t.Error("RenderBlogPost() missing twitter:image")
	}
}

func TestRenderPage_TwitterCardTags(t *testing.T) {
	t.Parallel()

	r, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	site := model.Site{
		Title:       "Quality Shepherd",
		Description: "A blog about testing",
		BaseURL:     "https://www.qualityshepherd.nl",
		Logo:        "/logo.svg",
	}
	page := model.Page{
		Title:   "About",
		Slug:    "about",
		Content: "<p>About content</p>",
	}

	got, err := r.RenderPage(site, page)
	if err != nil {
		t.Fatalf("RenderPage() error = %v", err)
	}

	// Check twitter:card
	if !strings.Contains(got, `<meta name="twitter:card" content="summary">`) {
		t.Error("RenderPage() missing twitter:card")
	}

	// Check twitter:title
	if !strings.Contains(got, `<meta name="twitter:title" content="About">`) {
		t.Error("RenderPage() missing twitter:title")
	}
}

func TestRenderPage_JSONLDNoArticleSchema(t *testing.T) {
	t.Parallel()

	r, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	site := model.Site{
		Title:       "Quality Shepherd",
		Description: "A blog about testing",
		BaseURL:     "https://www.qualityshepherd.nl",
	}
	page := model.Page{
		Title:   "About",
		Slug:    "about",
		Content: "<p>About content</p>",
	}

	got, err := r.RenderPage(site, page)
	if err != nil {
		t.Fatalf("RenderPage() error = %v", err)
	}

	// Check WebSite schema exists
	if !strings.Contains(got, `"@type": "WebSite"`) {
		t.Error("RenderPage() missing WebSite @type in JSON-LD")
	}

	// Regular pages should NOT have Article schema
	if strings.Contains(got, `"@type": "Article"`) {
		t.Error("RenderPage() should NOT have Article schema")
	}
}

func TestRenderBlogPost_JSONLDArticleSchema(t *testing.T) {
	t.Parallel()

	r, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	site := model.Site{
		Title:       "Quality Shepherd",
		Description: "A blog about testing",
		BaseURL:     "https://www.qualityshepherd.nl",
		Author:      "Jeroen",
	}
	post := model.Post{
		Page: model.Page{
			Title:   "My Test Post",
			Slug:    "my-test-post",
			Content: "<p>Post content</p>",
		},
		Date:    time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		Summary: "This is the post summary",
	}

	got, err := r.RenderBlogPost(site, post)
	if err != nil {
		t.Fatalf("RenderBlogPost() error = %v", err)
	}

	// Check JSON-LD script tag exists
	if !strings.Contains(got, `<script type="application/ld+json">`) {
		t.Error("RenderBlogPost() missing JSON-LD script tag")
	}

	// Check Article schema type
	if !strings.Contains(got, `"@type": "Article"`) {
		t.Error("RenderBlogPost() missing Article @type in JSON-LD")
	}

	// Check headline
	if !strings.Contains(got, `"headline": "My Test Post"`) {
		t.Error("RenderBlogPost() missing headline in Article JSON-LD")
	}

	// Check datePublished
	if !strings.Contains(got, `"datePublished": "2024-01-15"`) {
		t.Error("RenderBlogPost() missing datePublished in Article JSON-LD")
	}

	// Check author
	if !strings.Contains(got, `"author"`) && !strings.Contains(got, `"Jeroen"`) {
		t.Error("RenderBlogPost() missing author in Article JSON-LD")
	}

	// Check description
	if !strings.Contains(got, `"description": "This is the post summary"`) {
		t.Error("RenderBlogPost() missing description in Article JSON-LD")
	}
}

func TestRenderer_SetVersion(t *testing.T) {
	t.Parallel()

	r, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// SetVersion should store version
	r.SetVersion("v1.2.3")

	site := model.Site{
		Title:   "Test Site",
		BaseURL: "https://example.com",
	}

	// Test version appears in RenderPage
	page := model.Page{
		Title:   "About",
		Slug:    "about",
		Content: "<p>About content</p>",
	}
	got, err := r.RenderPage(site, page)
	if err != nil {
		t.Fatalf("RenderPage() error = %v", err)
	}
	if !strings.Contains(got, "v1.2.3") {
		t.Error("RenderPage() should include version in output")
	}

	// Test version appears in RenderBlogList
	got, err = r.RenderBlogList(site, nil)
	if err != nil {
		t.Fatalf("RenderBlogList() error = %v", err)
	}
	if !strings.Contains(got, "v1.2.3") {
		t.Error("RenderBlogList() should include version in output")
	}

	// Test version appears in RenderBlogPost
	post := model.Post{
		Page: model.Page{
			Title:   "Post",
			Slug:    "post",
			Content: "<p>Content</p>",
		},
		Date: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
	}
	got, err = r.RenderBlogPost(site, post)
	if err != nil {
		t.Fatalf("RenderBlogPost() error = %v", err)
	}
	if !strings.Contains(got, "v1.2.3") {
		t.Error("RenderBlogPost() should include version in output")
	}

	// Test version appears in Render404
	got, err = r.Render404(site)
	if err != nil {
		t.Fatalf("Render404() error = %v", err)
	}
	if !strings.Contains(got, "v1.2.3") {
		t.Error("Render404() should include version in output")
	}

	// Test version appears in RenderBase
	got, err = r.RenderBase(site, "Test content")
	if err != nil {
		t.Fatalf("RenderBase() error = %v", err)
	}
	if !strings.Contains(got, "v1.2.3") {
		t.Error("RenderBase() should include version in output")
	}
}

func TestRenderBlogPost_OGImagePreferredOverLogo(t *testing.T) {
	t.Parallel()

	r, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	site := model.Site{
		Title:   "Test Site",
		BaseURL: "https://example.com",
		Logo:    "/logo.svg",
		OGImage: "/social-image.png",
	}
	post := model.Post{
		Page: model.Page{
			Title:   "Test Post",
			Slug:    "test-post",
			Content: "<p>Content</p>",
		},
		Date: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
	}

	got, err := r.RenderBlogPost(site, post)
	if err != nil {
		t.Fatalf("RenderBlogPost() error = %v", err)
	}

	// OGImage should be preferred over Logo
	wantOGImage := `<meta property="og:image" content="https://example.com/social-image.png">`
	if !strings.Contains(got, wantOGImage) {
		t.Errorf("RenderBlogPost() should use OGImage for og:image, want %q", wantOGImage)
	}

	wantTwitterImage := `<meta name="twitter:image" content="https://example.com/social-image.png">`
	if !strings.Contains(got, wantTwitterImage) {
		t.Errorf("RenderBlogPost() should use OGImage for twitter:image, want %q", wantTwitterImage)
	}

	// Should NOT use the logo
	unwantOGImage := `<meta property="og:image" content="https://example.com/logo.svg">`
	if strings.Contains(got, unwantOGImage) {
		t.Error("RenderBlogPost() should NOT use Logo when OGImage is set")
	}
}

func TestRenderBlogPost_LogoFallbackWhenOGImageEmpty(t *testing.T) {
	t.Parallel()

	r, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	site := model.Site{
		Title:   "Test Site",
		BaseURL: "https://example.com",
		Logo:    "/logo.svg",
		OGImage: "", // Empty - should fall back to Logo
	}
	post := model.Post{
		Page: model.Page{
			Title:   "Test Post",
			Slug:    "test-post",
			Content: "<p>Content</p>",
		},
		Date: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
	}

	got, err := r.RenderBlogPost(site, post)
	if err != nil {
		t.Fatalf("RenderBlogPost() error = %v", err)
	}

	// Should fall back to Logo when OGImage is empty
	wantOGImage := `<meta property="og:image" content="https://example.com/logo.svg">`
	if !strings.Contains(got, wantOGImage) {
		t.Errorf("RenderBlogPost() should fallback to Logo for og:image, want %q", wantOGImage)
	}
}

func TestRenderBlogPost_NoOGImageWhenBothEmpty(t *testing.T) {
	t.Parallel()

	r, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	site := model.Site{
		Title:   "Test Site",
		BaseURL: "https://example.com",
		Logo:    "",
		OGImage: "",
	}
	post := model.Post{
		Page: model.Page{
			Title:   "Test Post",
			Slug:    "test-post",
			Content: "<p>Content</p>",
		},
		Date: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
	}

	got, err := r.RenderBlogPost(site, post)
	if err != nil {
		t.Fatalf("RenderBlogPost() error = %v", err)
	}

	// Should NOT have og:image when both are empty
	if strings.Contains(got, `<meta property="og:image"`) {
		t.Error("RenderBlogPost() should NOT include og:image when both OGImage and Logo are empty")
	}

	if strings.Contains(got, `<meta name="twitter:image"`) {
		t.Error("RenderBlogPost() should NOT include twitter:image when both OGImage and Logo are empty")
	}
}

func TestRenderBase_GoatCounterAnalytics(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		site       model.Site
		wantScript bool
		wantCode   string
	}{
		{
			name: "includes GoatCounter script when configured",
			site: model.Site{
				Title:   "Test Site",
				BaseURL: "https://example.com",
				Analytics: model.Analytics{
					GoatCounter: "aishepherd",
				},
			},
			wantScript: true,
			wantCode:   "aishepherd",
		},
		{
			name: "no GoatCounter script when not configured",
			site: model.Site{
				Title:   "Test Site",
				BaseURL: "https://example.com",
			},
			wantScript: false,
		},
		{
			name: "no GoatCounter script when empty string",
			site: model.Site{
				Title:   "Test Site",
				BaseURL: "https://example.com",
				Analytics: model.Analytics{
					GoatCounter: "",
				},
			},
			wantScript: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r, err := New()
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}

			got, err := r.RenderBase(tt.site, "Test Content")
			if err != nil {
				t.Fatalf("RenderBase() error = %v", err)
			}

			hasGoatCounter := strings.Contains(got, "goatcounter.com/count")
			if tt.wantScript && !hasGoatCounter {
				t.Error("RenderBase() should include GoatCounter script")
			}
			if !tt.wantScript && hasGoatCounter {
				t.Error("RenderBase() should not include GoatCounter script")
			}

			if tt.wantScript && tt.wantCode != "" {
				wantURL := tt.wantCode + ".goatcounter.com"
				if !strings.Contains(got, wantURL) {
					t.Errorf("RenderBase() should include site code %q in GoatCounter URL", tt.wantCode)
				}
			}
		})
	}
}
