package topics

import (
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/jeroendee/ssg/internal/model"
)

// Extract extracts recurring subject words from markdown content and returns
// them as frequency-counted topics, sorted by count descending with
// alphabetical tiebreaker. Returns at most 18 topics. Words must be at
// least 3 characters long and appear at least 2 times.
func Extract(markdown string) []model.Topic {
	if markdown == "" {
		return nil
	}

	// Strip markdown syntax
	text := stripMarkdown(markdown)

	// Tokenize
	words := tokenize(text)

	// Count frequencies
	freq := make(map[string]int)
	for _, w := range words {
		w = strings.ToLower(w)
		if len(w) < 3 {
			continue
		}
		if stopWords[w] {
			continue
		}
		freq[w]++
	}

	// Filter by minimum frequency and build result
	var result []model.Topic
	for word, count := range freq {
		if count < 2 {
			continue
		}
		result = append(result, model.Topic{Word: word, Count: count})
	}

	// Sort: frequency descending, alphabetical tiebreaker
	sort.Slice(result, func(i, j int) bool {
		if result[i].Count != result[j].Count {
			return result[i].Count > result[j].Count
		}
		return result[i].Word < result[j].Word
	})

	// Cap at 18
	if len(result) > 18 {
		result = result[:18]
	}

	return result
}

// Regex patterns for markdown stripping.
var (
	// Match image references: ![alt](path)
	reImage = regexp.MustCompile(`!\[[^\]]*\]\([^)]*\)`)
	// Match links: [text](url) — capture text, discard URL
	reLink = regexp.MustCompile(`\[([^\]]*)\]\([^)]*\)`)
	// Match HTML entities: &amp; &quot; etc.
	reHTMLEntity = regexp.MustCompile(`&[a-zA-Z]+;`)
	// Match inline code backticks (keep content)
	reInlineCode = regexp.MustCompile("`([^`]*)`")
)

// stripMarkdown removes markdown syntax artifacts before tokenization.
func stripMarkdown(md string) string {
	// Strip image refs first (before links, since ![...] contains [...])
	text := reImage.ReplaceAllString(md, "")
	// Replace links with just the text
	text = reLink.ReplaceAllString(text, "$1")
	// Strip HTML entities
	text = reHTMLEntity.ReplaceAllString(text, "")
	// Strip inline code backticks but keep content
	text = reInlineCode.ReplaceAllString(text, "$1")
	return text
}

// tokenize splits text into words, preserving hyphens within words.
func tokenize(text string) []string {
	var words []string
	var current strings.Builder

	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			current.WriteRune(r)
		} else if r == '-' && current.Len() > 0 {
			// Peek: hyphen might be part of a compound word
			current.WriteRune(r)
		} else {
			if current.Len() > 0 {
				w := strings.TrimRight(current.String(), "-")
				if w != "" {
					words = append(words, w)
				}
				current.Reset()
			}
		}
	}
	if current.Len() > 0 {
		w := strings.TrimRight(current.String(), "-")
		if w != "" {
			words = append(words, w)
		}
	}

	return words
}

// stopWords contains common English words to exclude from topic extraction.
var stopWords = map[string]bool{
	// Articles
	"the": true, "a": true, "an": true,
	// Prepositions
	"in": true, "on": true, "at": true, "to": true, "for": true,
	"with": true, "from": true, "by": true, "of": true, "about": true,
	"into": true, "through": true, "during": true, "before": true,
	"after": true, "above": true, "below": true, "between": true,
	"under": true, "over": true, "out": true, "off": true, "up": true,
	"down": true, "upon": true, "along": true, "across": true, "via": true,
	// Pronouns
	"he": true, "she": true, "it": true, "they": true, "we": true,
	"you": true, "his": true, "her": true, "its": true, "their": true,
	"our": true, "your": true, "him": true, "them": true, "who": true,
	"whom": true, "whose": true, "which": true, "that": true,
	"this": true, "these": true, "those": true, "what": true,
	"myself": true, "yourself": true, "himself": true, "herself": true,
	"itself": true, "ourselves": true, "themselves": true,
	// Common verbs
	"is": true, "are": true, "was": true, "were": true, "be": true,
	"been": true, "being": true, "has": true, "have": true, "had": true,
	"having": true, "do": true, "does": true, "did": true, "doing": true,
	"will": true, "would": true, "shall": true, "should": true,
	"may": true, "might": true, "must": true, "can": true, "could": true,
	"am": true, "get": true, "got": true, "gets": true, "make": true,
	"made": true, "let": true, "say": true, "said": true, "know": true,
	"think": true, "take": true, "come": true, "see": true, "want": true,
	"use": true, "used": true, "using": true, "find": true, "give": true,
	"tell": true, "work": true, "call": true, "try": true, "ask": true,
	"need": true, "seem": true, "feel": true, "leave": true, "put": true,
	"keep": true, "set": true, "run": true, "move": true, "go": true,
	"went": true, "gone": true, "going": true,
	// Conjunctions
	"and": true, "but": true, "or": true, "nor": true, "so": true,
	"yet": true, "both": true, "either": true, "neither": true,
	"not": true, "only": true, "own": true, "same": true,
	// Common adverbs
	"also": true, "just": true, "then": true, "than": true,
	"now": true, "here": true, "there": true, "when": true,
	"where": true, "why": true, "how": true, "all": true,
	"each": true, "every": true, "any": true, "few": true,
	"more": true, "most": true, "other": true, "some": true,
	"such": true, "no": true, "very": true, "too": true,
	"quite": true, "enough": true, "well": true, "back": true,
	"still": true, "even": true, "never": true, "always": true,
	"often": true, "ever": true, "much": true, "many": true,
	// Other common words
	"like": true, "one": true, "two": true,
	"new": true, "old": true, "first": true, "last": true,
	"long": true, "great": true, "little": true, "right": true,
	"big": true, "high": true, "small": true, "large": true,
	"next": true, "early": true, "young": true, "important": true,
	"public": true, "bad": true, "different": true, "able": true,
	"way": true, "day": true, "time": true, "year": true,
	"people": true, "part": true, "place": true, "case": true,
	"thing": true, "man": true, "world": true, "life": true,
	"hand": true, "point": true, "end": true, "another": true,
	"again": true, "don": true, "article": true, "post": true, "dev": true,
	"based": true,
}
