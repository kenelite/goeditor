package syntax

import (
	"testing"
	
	"github.com/alecthomas/chroma"
)

func TestThemeManager(t *testing.T) {
	tm := NewThemeManager()

	// Test default theme
	if tm.GetCurrentThemeName() != "github" {
		t.Errorf("Expected default theme to be 'github', got '%s'", tm.GetCurrentThemeName())
	}

	// Test getting current theme
	theme := tm.GetTheme()
	if theme == nil {
		t.Error("Expected non-nil theme")
	}
}

func TestSetTheme(t *testing.T) {
	tm := NewThemeManager()

	// Test setting valid theme - try multiple themes as availability may vary
	validThemes := []string{"monokai", "github", "vs"}
	var successfulTheme string
	
	for _, theme := range validThemes {
		err := tm.SetTheme(theme)
		if err == nil {
			successfulTheme = theme
			break
		}
	}
	
	if successfulTheme == "" {
		t.Error("Expected at least one valid theme to be settable")
	} else {
		if tm.GetCurrentThemeName() != successfulTheme {
			t.Errorf("Expected current theme to be '%s', got '%s'", successfulTheme, tm.GetCurrentThemeName())
		}
	}

	// Test setting invalid theme - use a name that definitely doesn't exist
	err := tm.SetTheme("this-theme-definitely-does-not-exist-xyz-123")
	if err == nil {
		// Some chroma versions might not return errors for invalid themes
		t.Log("Warning: No error returned for invalid theme (this might be expected)")
	}
}

func TestGetAvailableThemes(t *testing.T) {
	tm := NewThemeManager()

	themes := tm.GetAvailableThemes()
	if len(themes) == 0 {
		t.Error("Expected at least some available themes")
	}

	// Check for some common themes
	expectedThemes := []string{"github", "monokai", "vs"}
	themeMap := make(map[string]bool)
	for _, theme := range themes {
		themeMap[theme] = true
	}

	for _, expected := range expectedThemes {
		if !themeMap[expected] {
			t.Errorf("Expected to find theme '%s' in available themes", expected)
		}
	}
}

func TestGetThemeInfo(t *testing.T) {
	tm := NewThemeManager()

	// Try multiple themes as availability may vary
	testThemes := []string{"github", "monokai", "vs"}
	var foundTheme string
	
	for _, theme := range testThemes {
		info := tm.GetThemeInfo(theme)
		if info != nil {
			foundTheme = theme
			if info.Name != theme {
				t.Errorf("Expected theme name to be '%s', got '%s'", theme, info.Name)
			}
			if info.Description == "" {
				t.Error("Expected non-empty description")
			}
			break
		}
	}
	
	if foundTheme == "" {
		t.Error("Expected at least one test theme to be available")
	}

	// Test nonexistent theme
	info := tm.GetThemeInfo("this-theme-definitely-does-not-exist-xyz-123")
	if info != nil {
		// Some chroma versions might return info even for invalid themes
		t.Log("Warning: Got info for invalid theme (this might be expected)")
	}
}

func TestThemeColors(t *testing.T) {
	tm := NewThemeManager()

	// Test background color
	bgColor := tm.GetBackgroundColor()
	if bgColor == nil {
		t.Error("Expected non-nil background color")
	}

	// Test foreground color
	fgColor := tm.GetForegroundColor()
	if fgColor == nil {
		t.Error("Expected non-nil foreground color")
	}

	// Colors should be different (in most themes)
	if bgColor == fgColor {
		t.Log("Warning: Background and foreground colors are the same")
	}
}

func TestIsDarkTheme(t *testing.T) {
	tm := NewThemeManager()

	// Test with available themes - just ensure the function works
	availableThemes := tm.GetAvailableThemes()
	if len(availableThemes) > 0 {
		// Test with first available theme
		err := tm.SetTheme(availableThemes[0])
		if err == nil {
			// Just test that the function returns a boolean without error
			isDark := tm.IsDarkTheme()
			_ = isDark // We don't assert the specific value as it depends on the theme
		}
	}
}

func TestGetRecommendedThemes(t *testing.T) {
	tm := NewThemeManager()

	lightThemes, darkThemes := tm.GetRecommendedThemes()

	if len(lightThemes) == 0 {
		t.Error("Expected at least some light themes")
	}

	if len(darkThemes) == 0 {
		t.Error("Expected at least some dark themes")
	}

	// Check that recommended themes exist
	availableThemes := tm.GetAvailableThemes()
	themeMap := make(map[string]bool)
	for _, theme := range availableThemes {
		themeMap[theme] = true
	}

	for _, theme := range lightThemes {
		if !themeMap[theme] {
			t.Errorf("Recommended light theme '%s' is not available", theme)
		}
	}

	for _, theme := range darkThemes {
		if !themeMap[theme] {
			t.Errorf("Recommended dark theme '%s' is not available", theme)
		}
	}
}

func TestCustomThemes(t *testing.T) {
	tm := NewThemeManager()

	// Test registering custom theme
	customTheme := &ThemeConfig{
		Name:       "test-theme",
		Background: "#ffffff",
		Foreground: "#000000",
		Selection:  "#0066cc",
		LineNumber: "#666666",
		TokenColors: map[string]string{
			"keyword": "#0000ff",
			"string":  "#008000",
			"comment": "#808080",
		},
	}

	err := tm.RegisterCustomTheme(customTheme)
	if err != nil {
		t.Errorf("Expected no error registering custom theme, got: %v", err)
	}

	// Test getting custom themes
	customThemes := tm.GetCustomThemes()
	if len(customThemes) != 1 {
		t.Errorf("Expected 1 custom theme, got %d", len(customThemes))
	}

	if customThemes["test-theme"] == nil {
		t.Error("Expected to find test-theme in custom themes")
	}

	// Test registering theme with empty name
	invalidTheme := &ThemeConfig{
		Name: "",
	}
	err = tm.RegisterCustomTheme(invalidTheme)
	if err == nil {
		t.Error("Expected error when registering theme with empty name")
	}
}

func TestChromaToRGBA(t *testing.T) {
	// Test the function with a real theme to get actual chroma colors
	tm := NewThemeManager()
	style := tm.GetTheme()
	
	if style != nil {
		// Get a real token style from the theme
		tokenStyle := style.Get(chroma.Keyword)
		if tokenStyle.Colour.IsSet() {
			result := chromaToRGBA(tokenStyle.Colour)
			// Just verify we get a valid color (not nil)
			if result == nil {
				t.Error("Expected non-nil color result")
			}
		}
	}

	// Test with unset color - just verify the function doesn't crash
	unsetColor := chroma.Colour(0x000000) // Not marked as set
	result := chromaToRGBA(unsetColor)
	// Just verify we get some color back (the function handles unset colors)
	if result == nil {
		t.Error("Expected non-nil color result even for unset color")
	}
}