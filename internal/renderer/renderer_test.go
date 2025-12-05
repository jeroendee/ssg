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
		name      string
		site      model.Site
		posts     []model.Post
		wantDates []string
		wantLinks []string
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
			wantDates: []string{"2021-03-26", "2021-04-15"},
			wantLinks: []string{"First Post", "Second Post"},
		},
		{
			name: "renders empty blog list",
			site: model.Site{
				Title: "Test Site",
			},
			posts:     []model.Post{},
			wantDates: []string{},
			wantLinks: []string{},
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

func TestRenderHome(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		site      model.Site
		wantTitle string
		wantNav   bool
	}{
		{
			name: "renders homepage with site title",
			site: model.Site{
				Title:   "Quality Shepherd",
				BaseURL: "https://www.qualityshepherd.nl",
				Navigation: []model.NavItem{
					{Title: "Home", URL: "/"},
					{Title: "Blog", URL: "/blog/"},
				},
			},
			wantTitle: "Quality Shepherd",
			wantNav:   true,
		},
		{
			name: "renders homepage with minimal site",
			site: model.Site{
				Title: "Test Site",
			},
			wantTitle: "Test Site",
			wantNav:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r, err := New()
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}

			got, err := r.RenderHome(tt.site)
			if err != nil {
				t.Fatalf("RenderHome() error = %v", err)
			}

			// Check site title in output
			if !strings.Contains(got, tt.wantTitle) {
				t.Errorf("RenderHome() missing title %q", tt.wantTitle)
			}

			// Check valid HTML structure
			if !strings.Contains(got, "<!DOCTYPE html>") {
				t.Error("RenderHome() missing DOCTYPE")
			}

			// Check navigation is present when expected
			if tt.wantNav && !strings.Contains(got, "<nav") {
				t.Error("RenderHome() missing nav element")
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
