package topics_test

import (
	"strings"
	"testing"

	"github.com/jeroendee/ssg/internal/topics"
)

func TestExtract_BasicFrequency(t *testing.T) {
	t.Parallel()
	md := "Claude is great. Claude helps with coding. Anthropic built Claude."
	result := topics.Extract(md)

	found := false
	for _, topic := range result {
		if topic.Word == "claude" && topic.Count == 3 {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected topic 'claude' with count 3, got %v", result)
	}
}

func TestExtract_CaseInsensitive(t *testing.T) {
	t.Parallel()
	md := "LLM and llm and Llm are all the same."
	result := topics.Extract(md)

	found := false
	for _, topic := range result {
		if topic.Word == "llm" && topic.Count == 3 {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected topic 'llm' with count 3, got %v", result)
	}
}

func TestExtract_EmptyInput(t *testing.T) {
	t.Parallel()
	result := topics.Extract("")
	if len(result) != 0 {
		t.Errorf("expected empty slice for empty input, got %v", result)
	}
}

func TestExtract_OnlyStopWords(t *testing.T) {
	t.Parallel()
	md := "the and or is it was the and or is it was"
	result := topics.Extract(md)
	if len(result) != 0 {
		t.Errorf("expected empty slice for stop words only, got %v", result)
	}
}

func TestExtract_StopWordsFiltered(t *testing.T) {
	t.Parallel()
	md := "the the the and and and or or or docker docker docker"
	result := topics.Extract(md)

	for _, topic := range result {
		switch topic.Word {
		case "the", "and", "or":
			t.Errorf("stop word %q should not appear in results", topic.Word)
		}
	}

	found := false
	for _, topic := range result {
		if topic.Word == "docker" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected 'docker' in results, got %v", result)
	}
}

func TestExtract_ShortWordsExcluded(t *testing.T) {
	t.Parallel()
	md := "ai ai ai go go go up up up docker docker docker"
	result := topics.Extract(md)

	for _, topic := range result {
		if len(topic.Word) < 3 {
			t.Errorf("word %q has length %d, expected >= 3", topic.Word, len(topic.Word))
		}
	}
}

func TestExtract_ThreeCharWordIncluded(t *testing.T) {
	t.Parallel()
	md := "llm llm llm"
	result := topics.Extract(md)

	found := false
	for _, topic := range result {
		if topic.Word == "llm" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected 'llm' in results, got %v", result)
	}
}

func TestExtract_MinFrequencyThreshold(t *testing.T) {
	t.Parallel()
	md := "kubernetes docker docker"
	result := topics.Extract(md)

	for _, topic := range result {
		if topic.Word == "kubernetes" {
			t.Error("single-occurrence word 'kubernetes' should not appear in results")
		}
	}

	found := false
	for _, topic := range result {
		if topic.Word == "docker" && topic.Count == 2 {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected 'docker' with count 2, got %v", result)
	}
}

func TestExtract_Top20Cap(t *testing.T) {
	t.Parallel()
	// Create 25 distinct words each appearing 3 times
	var words []string
	wordList := []string{
		"alpha", "bravo", "charlie", "delta", "echo",
		"foxtrot", "golf", "hotel", "india", "juliet",
		"kilo", "lima", "mike", "november", "oscar",
		"papa", "quebec", "romeo", "sierra", "tango",
		"uniform", "victor", "whiskey", "xray", "yankee",
	}
	for _, w := range wordList {
		words = append(words, w, w, w)
	}
	md := strings.Join(words, " ")
	result := topics.Extract(md)

	if len(result) != 20 {
		t.Errorf("expected 20 topics, got %d", len(result))
	}
}

func TestExtract_SortFrequencyDescending(t *testing.T) {
	t.Parallel()
	md := "agent agent agent agent agent agent agent agent agent agent " +
		"claude claude claude claude claude " +
		"docker docker docker"
	result := topics.Extract(md)

	if len(result) < 3 {
		t.Fatalf("expected at least 3 topics, got %d", len(result))
	}
	if result[0].Word != "agent" {
		t.Errorf("result[0].Word = %q, want %q", result[0].Word, "agent")
	}
	if result[1].Word != "claude" {
		t.Errorf("result[1].Word = %q, want %q", result[1].Word, "claude")
	}
	if result[2].Word != "docker" {
		t.Errorf("result[2].Word = %q, want %q", result[2].Word, "docker")
	}
}

func TestExtract_AlphabeticalTiebreaker(t *testing.T) {
	t.Parallel()
	md := "zebra zebra alpha alpha"
	result := topics.Extract(md)

	if len(result) < 2 {
		t.Fatalf("expected at least 2 topics, got %d", len(result))
	}
	if result[0].Word != "alpha" {
		t.Errorf("result[0].Word = %q, want %q (alphabetical tiebreaker)", result[0].Word, "alpha")
	}
	if result[1].Word != "zebra" {
		t.Errorf("result[1].Word = %q, want %q", result[1].Word, "zebra")
	}
}

func TestExtract_LinkTextKeptURLDiscarded(t *testing.T) {
	t.Parallel()
	md := "[Claude](https://anthropic.com) and [Claude](https://anthropic.com) again"
	result := topics.Extract(md)

	for _, topic := range result {
		if topic.Word == "https" || topic.Word == "anthropic" || topic.Word == "com" {
			t.Errorf("URL component %q should not appear in results", topic.Word)
		}
	}

	found := false
	for _, topic := range result {
		if topic.Word == "claude" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected 'claude' from link text, got %v", result)
	}
}

func TestExtract_ImageRefsStripped(t *testing.T) {
	t.Parallel()
	md := "![alt text](assets/image.png) docker docker docker"
	result := topics.Extract(md)

	for _, topic := range result {
		if topic.Word == "assets" || topic.Word == "image" || topic.Word == "png" {
			t.Errorf("image ref component %q should not appear in results", topic.Word)
		}
	}
}

func TestExtract_HTMLEntitiesStripped(t *testing.T) {
	t.Parallel()
	md := "&amp; &quot; docker docker docker"
	result := topics.Extract(md)

	for _, topic := range result {
		if topic.Word == "amp" || topic.Word == "quot" {
			t.Errorf("HTML entity component %q should not appear in results", topic.Word)
		}
	}
}

func TestExtract_InlineCodePreserved(t *testing.T) {
	t.Parallel()
	md := "`openspec` and `openspec` and `openspec` and `openspec` and `openspec`"
	result := topics.Extract(md)

	found := false
	for _, topic := range result {
		if topic.Word == "openspec" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected 'openspec' from inline code, got %v", result)
	}
}

func TestExtract_HyphenatedWords(t *testing.T) {
	t.Parallel()
	md := "pre-push pre-push pre-push"
	result := topics.Extract(md)

	found := false
	for _, topic := range result {
		if topic.Word == "pre-push" && topic.Count == 3 {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected 'pre-push' with count 3, got %v", result)
	}
}
