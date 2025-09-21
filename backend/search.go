package backend

import (
	"regexp"
	"strings"
	"unicode"
)

// SearchManager handles search and replace operations
type SearchManager struct {
	currentPattern string
	matches        []Match
	currentIndex   int
	options        SearchOptions
	lastSearchText string
}

// Match represents a search match result
type Match struct {
	Start    Position `json:"start"`
	End      Position `json:"end"`
	Text     string   `json:"text"`
}

// SearchOptions defines search behavior options
type SearchOptions struct {
	CaseSensitive     bool `json:"caseSensitive"`
	WholeWord         bool `json:"wholeWord"`
	RegularExpression bool `json:"regularExpression"`
	WrapAround        bool `json:"wrapAround"`
}

// ReplaceOptions defines replace behavior options
type ReplaceOptions struct {
	SearchOptions
	ReplaceAll bool `json:"replaceAll"`
}

// NewSearchManager creates a new search manager
func NewSearchManager() *SearchManager {
	return &SearchManager{
		matches:      make([]Match, 0),
		currentIndex: -1,
		options: SearchOptions{
			CaseSensitive:     false,
			WholeWord:         false,
			RegularExpression: false,
			WrapAround:        true,
		},
	}
}

// SetOptions updates the search options
func (sm *SearchManager) SetOptions(options SearchOptions) {
	sm.options = options
	// Clear matches if options changed to force re-search
	if sm.currentPattern != "" {
		sm.matches = make([]Match, 0)
		sm.currentIndex = -1
	}
}

// GetOptions returns the current search options
func (sm *SearchManager) GetOptions() SearchOptions {
	return sm.options
}

// Find searches for a pattern in the given text and returns all matches
func (sm *SearchManager) Find(text, pattern string) []Match {
	if pattern == "" {
		sm.matches = make([]Match, 0)
		sm.currentIndex = -1
		return sm.matches
	}

	sm.currentPattern = pattern
	sm.lastSearchText = text
	sm.matches = make([]Match, 0)
	sm.currentIndex = -1

	if sm.options.RegularExpression {
		sm.findRegex(text, pattern)
	} else {
		sm.findLiteral(text, pattern)
	}

	return sm.matches
}

// findLiteral performs literal string search
func (sm *SearchManager) findLiteral(text, pattern string) {
	searchText := text
	searchPattern := pattern

	// Handle case sensitivity
	if !sm.options.CaseSensitive {
		searchText = strings.ToLower(text)
		searchPattern = strings.ToLower(pattern)
	}

	lines := strings.Split(text, "\n")
	searchLines := strings.Split(searchText, "\n")

	for lineNum, line := range searchLines {
		startCol := 0
		for {
			index := strings.Index(line[startCol:], searchPattern)
			if index == -1 {
				break
			}

			actualIndex := startCol + index
			
			// Check whole word option
			if sm.options.WholeWord && !sm.isWholeWord(line, actualIndex, len(searchPattern)) {
				startCol = actualIndex + 1
				continue
			}

			// Create match using original text (preserve case)
			originalLine := lines[lineNum]
			matchText := originalLine[actualIndex : actualIndex+len(searchPattern)]
			
			match := Match{
				Start: Position{Line: lineNum + 1, Column: actualIndex + 1},
				End:   Position{Line: lineNum + 1, Column: actualIndex + len(searchPattern) + 1},
				Text:  matchText,
			}
			sm.matches = append(sm.matches, match)

			startCol = actualIndex + 1
		}
	}
}

// findRegex performs regular expression search
func (sm *SearchManager) findRegex(text, pattern string) {
	var flags string
	if !sm.options.CaseSensitive {
		flags = "(?i)"
	}
	
	if sm.options.WholeWord {
		pattern = `\b` + pattern + `\b`
	}

	regex, err := regexp.Compile(flags + pattern)
	if err != nil {
		// Invalid regex, return empty matches
		return
	}

	lines := strings.Split(text, "\n")
	for lineNum, line := range lines {
		matches := regex.FindAllStringIndex(line, -1)
		for _, match := range matches {
			start, end := match[0], match[1]
			matchText := line[start:end]
			
			matchObj := Match{
				Start: Position{Line: lineNum + 1, Column: start + 1},
				End:   Position{Line: lineNum + 1, Column: end + 1},
				Text:  matchText,
			}
			sm.matches = append(sm.matches, matchObj)
		}
	}
}

// isWholeWord checks if the match at the given position is a whole word
func (sm *SearchManager) isWholeWord(line string, start, length int) bool {
	end := start + length

	// Check character before match
	if start > 0 {
		prevChar := rune(line[start-1])
		if unicode.IsLetter(prevChar) || unicode.IsDigit(prevChar) || prevChar == '_' {
			return false
		}
	}

	// Check character after match
	if end < len(line) {
		nextChar := rune(line[end])
		if unicode.IsLetter(nextChar) || unicode.IsDigit(nextChar) || nextChar == '_' {
			return false
		}
	}

	return true
}

// GetMatches returns all current matches
func (sm *SearchManager) GetMatches() []Match {
	return sm.matches
}

// GetMatchCount returns the number of matches found
func (sm *SearchManager) GetMatchCount() int {
	return len(sm.matches)
}

// GetCurrentIndex returns the current match index
func (sm *SearchManager) GetCurrentIndex() int {
	return sm.currentIndex
}

// NextMatch moves to the next match and returns it
func (sm *SearchManager) NextMatch() *Match {
	if len(sm.matches) == 0 {
		return nil
	}

	sm.currentIndex++
	if sm.currentIndex >= len(sm.matches) {
		if sm.options.WrapAround {
			sm.currentIndex = 0
		} else {
			sm.currentIndex = len(sm.matches) - 1
			return nil
		}
	}

	return &sm.matches[sm.currentIndex]
}

// PreviousMatch moves to the previous match and returns it
func (sm *SearchManager) PreviousMatch() *Match {
	if len(sm.matches) == 0 {
		return nil
	}

	sm.currentIndex--
	if sm.currentIndex < 0 {
		if sm.options.WrapAround {
			sm.currentIndex = len(sm.matches) - 1
		} else {
			sm.currentIndex = 0
			return nil
		}
	}

	return &sm.matches[sm.currentIndex]
}

// SetCurrentMatch sets the current match by index
func (sm *SearchManager) SetCurrentMatch(index int) *Match {
	if index < 0 || index >= len(sm.matches) {
		return nil
	}

	sm.currentIndex = index
	return &sm.matches[sm.currentIndex]
}

// GetCurrentMatch returns the current match
func (sm *SearchManager) GetCurrentMatch() *Match {
	if sm.currentIndex < 0 || sm.currentIndex >= len(sm.matches) {
		return nil
	}
	return &sm.matches[sm.currentIndex]
}

// Replace performs text replacement
func (sm *SearchManager) Replace(text, pattern, replacement string, options ReplaceOptions) (string, int) {
	sm.SetOptions(options.SearchOptions)
	
	if options.ReplaceAll {
		return sm.replaceAll(text, pattern, replacement)
	} else {
		return sm.replaceCurrent(text, pattern, replacement)
	}
}

// replaceAll replaces all occurrences of the pattern
func (sm *SearchManager) replaceAll(text, pattern, replacement string) (string, int) {
	// Find all matches first
	matches := sm.Find(text, pattern)
	if len(matches) == 0 {
		return text, 0
	}

	// Replace from end to beginning to maintain position accuracy
	result := text
	replacedCount := 0

	for i := len(matches) - 1; i >= 0; i-- {
		match := matches[i]
		
		// Convert positions to string indices
		lines := strings.Split(result, "\n")
		if match.Start.Line-1 >= len(lines) {
			continue
		}
		
		line := lines[match.Start.Line-1]
		if match.Start.Column-1 >= len(line) || match.End.Column-1 > len(line) {
			continue
		}

		// Calculate absolute position in text
		absoluteStart := 0
		for j := 0; j < match.Start.Line-1; j++ {
			absoluteStart += len(lines[j]) + 1 // +1 for newline
		}
		absoluteStart += match.Start.Column - 1

		absoluteEnd := absoluteStart + (match.End.Column - match.Start.Column)

		// Perform replacement
		if absoluteEnd <= len(result) {
			actualReplacement := replacement
			if sm.options.RegularExpression {
				// For regex, we might need to handle capture groups
				actualReplacement = sm.processRegexReplacement(match.Text, pattern, replacement)
			}
			
			result = result[:absoluteStart] + actualReplacement + result[absoluteEnd:]
			replacedCount++
		}
	}

	return result, replacedCount
}

// replaceCurrent replaces only the current match
func (sm *SearchManager) replaceCurrent(text, pattern, replacement string) (string, int) {
	currentMatch := sm.GetCurrentMatch()
	if currentMatch == nil {
		// No current match, find first match
		matches := sm.Find(text, pattern)
		if len(matches) == 0 {
			return text, 0
		}
		sm.currentIndex = 0
		currentMatch = &matches[0]
	}

	// Calculate absolute position
	lines := strings.Split(text, "\n")
	if currentMatch.Start.Line-1 >= len(lines) {
		return text, 0
	}

	absoluteStart := 0
	for i := 0; i < currentMatch.Start.Line-1; i++ {
		absoluteStart += len(lines[i]) + 1 // +1 for newline
	}
	absoluteStart += currentMatch.Start.Column - 1

	absoluteEnd := absoluteStart + (currentMatch.End.Column - currentMatch.Start.Column)

	if absoluteEnd > len(text) {
		return text, 0
	}

	// Perform replacement
	actualReplacement := replacement
	if sm.options.RegularExpression {
		actualReplacement = sm.processRegexReplacement(currentMatch.Text, pattern, replacement)
	}

	result := text[:absoluteStart] + actualReplacement + text[absoluteEnd:]
	
	// Update matches after replacement
	sm.Find(result, pattern)
	
	return result, 1
}

// processRegexReplacement handles regex replacement with capture groups
func (sm *SearchManager) processRegexReplacement(matchText, pattern, replacement string) string {
	var flags string
	if !sm.options.CaseSensitive {
		flags = "(?i)"
	}
	
	if sm.options.WholeWord {
		pattern = `\b` + pattern + `\b`
	}

	regex, err := regexp.Compile(flags + pattern)
	if err != nil {
		return replacement
	}

	return regex.ReplaceAllString(matchText, replacement)
}

// Clear clears all search state
func (sm *SearchManager) Clear() {
	sm.currentPattern = ""
	sm.matches = make([]Match, 0)
	sm.currentIndex = -1
	sm.lastSearchText = ""
}

// HasMatches returns true if there are any matches
func (sm *SearchManager) HasMatches() bool {
	return len(sm.matches) > 0
}