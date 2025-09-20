package syntax

import (
	"fmt"
	"image/color"
	"log"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/styles"
)

// ThemeManager manages syntax highlighting themes
type ThemeManager struct {
	currentTheme string
	themes       map[string]*chroma.Style
	customThemes map[string]*ThemeConfig
}

// ThemeConfig represents a custom theme configuration
type ThemeConfig struct {
	Name        string            `json:"name"`
	Background  string            `json:"background"`
	Foreground  string            `json:"foreground"`
	Selection   string            `json:"selection"`
	LineNumber  string            `json:"lineNumber"`
	TokenColors map[string]string `json:"tokenColors"`
}

// NewThemeManager creates a new theme manager
func NewThemeManager() *ThemeManager {
	tm := &ThemeManager{
		currentTheme: "github",
		themes:       make(map[string]*chroma.Style),
		customThemes: make(map[string]*ThemeConfig),
	}
	
	// Load default themes
	tm.loadDefaultThemes()
	
	return tm
}

// loadDefaultThemes loads the default available themes
func (tm *ThemeManager) loadDefaultThemes() {
	// Cache commonly used themes
	defaultThemes := []string{
		"github",
		"monokai",
		"dracula",
		"solarized-dark",
		"solarized-light",
		"vs",
		"xcode",
		"atom-one-dark",
		"atom-one-light",
	}

	for _, themeName := range defaultThemes {
		if style := styles.Get(themeName); style != nil {
			tm.themes[themeName] = style
		}
	}
}

// GetTheme returns the current theme
func (tm *ThemeManager) GetTheme() *chroma.Style {
	if theme, exists := tm.themes[tm.currentTheme]; exists {
		return theme
	}
	
	// Try to get from chroma styles
	if style := styles.Get(tm.currentTheme); style != nil {
		tm.themes[tm.currentTheme] = style
		return style
	}
	
	// Fall back to default
	log.Printf("Theme '%s' not found, using fallback", tm.currentTheme)
	return styles.Fallback
}

// SetTheme sets the current theme
func (tm *ThemeManager) SetTheme(themeName string) error {
	// Check if theme exists
	if _, exists := tm.themes[themeName]; exists {
		tm.currentTheme = themeName
		return nil
	}
	
	// Try to load from chroma
	if style := styles.Get(themeName); style != nil {
		tm.themes[themeName] = style
		tm.currentTheme = themeName
		return nil
	}
	
	return fmt.Errorf("theme '%s' not found", themeName)
}

// GetCurrentThemeName returns the name of the current theme
func (tm *ThemeManager) GetCurrentThemeName() string {
	return tm.currentTheme
}

// GetAvailableThemes returns a list of available theme names
func (tm *ThemeManager) GetAvailableThemes() []string {
	var themes []string
	
	// Add cached themes
	for name := range tm.themes {
		themes = append(themes, name)
	}
	
	// Add all available chroma themes
	for name := range styles.Registry {
		found := false
		for _, existing := range themes {
			if existing == name {
				found = true
				break
			}
		}
		if !found {
			themes = append(themes, name)
		}
	}
	
	return themes
}

// GetThemeInfo returns information about a theme
func (tm *ThemeManager) GetThemeInfo(themeName string) *ThemeInfo {
	style := styles.Get(themeName)
	if style == nil {
		return nil
	}

	info := &ThemeInfo{
		Name:        themeName,
		Description: fmt.Sprintf("Syntax highlighting theme: %s", themeName),
		IsDark:      tm.isThemeDark(style),
	}

	// Get background color if available
	if bg := style.Get(chroma.Background); bg.Background.IsSet() {
		info.BackgroundColor = chromaToHex(bg.Background)
	}

	return info
}

// ThemeInfo contains information about a theme
type ThemeInfo struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	IsDark          bool   `json:"isDark"`
	BackgroundColor string `json:"backgroundColor"`
}

// isThemeDark attempts to determine if a theme is dark
func (tm *ThemeManager) isThemeDark(style *chroma.Style) bool {
	bg := style.Get(chroma.Background)
	if !bg.Background.IsSet() {
		return false
	}

	// Convert to RGB and check brightness
	r := uint8((bg.Background >> 16) & 0xFF)
	g := uint8((bg.Background >> 8) & 0xFF)
	b := uint8(bg.Background & 0xFF)
	
	// Calculate perceived brightness
	brightness := (0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 255.0
	
	return brightness < 0.5
}

// chromaToHex converts chroma.Colour to hex string
func chromaToHex(c chroma.Colour) string {
	if !c.IsSet() {
		return ""
	}
	return fmt.Sprintf("#%06x", uint32(c)&0xFFFFFF)
}

// GetTokenColor returns the color for a specific token type in the current theme
func (tm *ThemeManager) GetTokenColor(tokenType chroma.TokenType) color.Color {
	style := tm.GetTheme()
	if style == nil {
		return color.Black
	}

	tokenStyle := style.Get(tokenType)
	if tokenStyle.Colour.IsSet() {
		return chromaToRGBA(tokenStyle.Colour)
	}

	// Fall back to default color
	return color.Black
}

// GetBackgroundColor returns the background color of the current theme
func (tm *ThemeManager) GetBackgroundColor() color.Color {
	style := tm.GetTheme()
	if style == nil {
		return color.White
	}

	bg := style.Get(chroma.Background)
	if bg.Background.IsSet() {
		return chromaToRGBA(bg.Background)
	}

	return color.White
}

// GetForegroundColor returns the foreground color of the current theme
func (tm *ThemeManager) GetForegroundColor() color.Color {
	style := tm.GetTheme()
	if style == nil {
		return color.Black
	}

	// Try to get text color
	text := style.Get(chroma.Text)
	if text.Colour.IsSet() {
		return chromaToRGBA(text.Colour)
	}

	// Fall back to black for light themes, white for dark themes
	if tm.isThemeDark(style) {
		return color.White
	}
	return color.Black
}

// RegisterCustomTheme registers a custom theme
func (tm *ThemeManager) RegisterCustomTheme(config *ThemeConfig) error {
	if config.Name == "" {
		return fmt.Errorf("theme name cannot be empty")
	}

	tm.customThemes[config.Name] = config
	
	// TODO: Convert ThemeConfig to chroma.Style and register it
	// This would require more complex implementation to create a chroma.Style from config
	
	return nil
}

// GetCustomThemes returns all registered custom themes
func (tm *ThemeManager) GetCustomThemes() map[string]*ThemeConfig {
	return tm.customThemes
}

// IsDarkTheme returns true if the current theme is dark
func (tm *ThemeManager) IsDarkTheme() bool {
	style := tm.GetTheme()
	if style == nil {
		return false
	}
	return tm.isThemeDark(style)
}

// GetRecommendedThemes returns recommended themes for light and dark modes
func (tm *ThemeManager) GetRecommendedThemes() (lightThemes, darkThemes []string) {
	lightThemes = []string{
		"github",
		"vs",
		"xcode",
		"atom-one-light",
		"solarized-light",
	}
	
	darkThemes = []string{
		"monokai",
		"dracula",
		"atom-one-dark",
		"solarized-dark",
		"vim",
	}
	
	return lightThemes, darkThemes
}