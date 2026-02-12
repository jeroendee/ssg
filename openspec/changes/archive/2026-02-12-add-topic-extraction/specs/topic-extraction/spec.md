## ADDED Requirements

### Requirement: Extract topics from markdown content

The system SHALL extract recurring subject words from raw markdown text and return them as frequency-counted topics.

#### Scenario: Basic topic extraction

- **WHEN** extracting topics from markdown containing "Claude is great. Claude helps with coding. Anthropic built Claude."
- **THEN** the result SHALL include "claude" with count 3
- **AND** the result SHALL include "anthropic" with count 1 (if minimum threshold is met)

#### Scenario: Case-insensitive counting

- **WHEN** extracting topics from markdown containing "LLM" and "llm" and "Llm"
- **THEN** the system SHALL count all variants as a single topic "llm" with count 3

#### Scenario: Empty markdown

- **WHEN** extracting topics from empty markdown
- **THEN** the result SHALL be an empty slice

#### Scenario: Markdown with only stop words

- **WHEN** extracting topics from markdown containing only common English words like "the and or is it"
- **THEN** the result SHALL be an empty slice

### Requirement: Filter stop words

The system SHALL exclude common English stop words from topic results.

#### Scenario: Common articles filtered

- **WHEN** extracting topics from markdown containing "the", "a", "an"
- **THEN** those words SHALL NOT appear in the results

#### Scenario: Common prepositions filtered

- **WHEN** extracting topics from markdown containing "in", "on", "at", "to", "for", "with", "from"
- **THEN** those words SHALL NOT appear in the results

#### Scenario: Common pronouns filtered

- **WHEN** extracting topics from markdown containing "he", "she", "it", "they", "we", "you"
- **THEN** those words SHALL NOT appear in the results

#### Scenario: Common verbs filtered

- **WHEN** extracting topics from markdown containing "is", "are", "was", "were", "has", "have", "had", "been", "will", "can"
- **THEN** those words SHALL NOT appear in the results

### Requirement: Enforce minimum word length

The system SHALL exclude words shorter than 3 characters.

#### Scenario: Short words excluded

- **WHEN** extracting topics from markdown containing "ai", "go", "up", "do"
- **THEN** those 2-character words SHALL NOT appear in the results

#### Scenario: Three-character words included

- **WHEN** extracting topics from markdown containing "llm" repeated 3 times
- **THEN** "llm" SHALL appear in the results

### Requirement: Enforce minimum frequency threshold

The system SHALL exclude words that appear fewer than 2 times.

#### Scenario: Single-occurrence word excluded

- **WHEN** extracting topics from markdown where "kubernetes" appears exactly once
- **THEN** "kubernetes" SHALL NOT appear in the results

#### Scenario: Two-occurrence word included

- **WHEN** extracting topics from markdown where "docker" appears exactly 2 times
- **THEN** "docker" SHALL appear in the results with count 2

### Requirement: Limit to top 20 topics

The system SHALL return at most 20 topics.

#### Scenario: More than 20 qualifying words

- **WHEN** extracting topics from markdown with more than 20 distinct qualifying words
- **THEN** the result SHALL contain exactly 20 topics
- **AND** the 20 topics SHALL be those with the highest frequencies

#### Scenario: Fewer than 20 qualifying words

- **WHEN** extracting topics from markdown with 5 distinct qualifying words
- **THEN** the result SHALL contain exactly 5 topics

### Requirement: Sort topics by frequency descending

The system SHALL return topics ordered from most frequent to least frequent.

#### Scenario: Frequency ordering

- **WHEN** extracting topics from markdown where "agent" appears 10 times, "claude" appears 5 times, and "docker" appears 3 times
- **THEN** the result SHALL be ordered: agent, claude, docker

#### Scenario: Equal frequency ordering

- **WHEN** two topics have the same frequency
- **THEN** the system SHALL use alphabetical order as tiebreaker

### Requirement: Strip markdown syntax before tokenization

The system SHALL remove markdown formatting artifacts before counting words.

#### Scenario: Link text extracted, URLs discarded

- **WHEN** extracting topics from markdown containing `[Claude](https://anthropic.com)`
- **THEN** "claude" SHALL be counted
- **AND** "https", "anthropic", "com" from the URL SHALL NOT be counted

#### Scenario: Image references stripped

- **WHEN** extracting topics from markdown containing `![alt text](assets/image.png)`
- **THEN** "assets", "image", "png" SHALL NOT be counted

#### Scenario: HTML entities stripped

- **WHEN** extracting topics from markdown containing `&amp;` or `&quot;`
- **THEN** "amp", "quot" SHALL NOT be counted

#### Scenario: Inline code preserved as words

- **WHEN** extracting topics from markdown containing `` `openspec` `` repeated 5 times
- **THEN** "openspec" SHALL appear in the results

### Requirement: Preserve hyphenated words

The system SHALL treat hyphenated terms as single words.

#### Scenario: Hyphenated word counted as one

- **WHEN** extracting topics from markdown containing "pre-push" repeated 3 times
- **THEN** "pre-push" SHALL appear as a single topic with count 3
