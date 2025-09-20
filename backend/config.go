package backend

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Configuration holds all application settings
type Configuration struct {
	Editor EditorConfig `json:"editor"`
	UI     UIConfig     `json:"ui"`
	Syntax SyntaxConfig `json:"syntax"`
}

// EditorConfig holds editor-specific settings
type EditorConfig struct {
	FontSize        int    `json:"fontSize"`
	TabSize         int    `json:"tabSize"`
	WordWrap        bool   `json:"wordWrap"`
	ShowLineNumbers bool   `json:"showLineNumbers"`
	AutoIndent      bool   `json:"autoIndent"`
	InsertSpaces    bool   `json:"insertSpaces"`
	TrimWhitespace  bool   `json:"trimWhitespace"`
}

// UIConfig holds user interface settings
type UIConfig struct {
	Theme        string `json:"theme"`        // "light" | "dark" | "auto"
	WindowWidth  int    `json:"windowWidth"`
	WindowHeight int    `json:"windowHeight"`
	WindowX      int    `json:"windowX"`
	WindowY      int    `json:"windowY"`
	Maximized    bool   `json:"maximized"`
}

// SyntaxConfig holds syntax highlighting settings
type SyntaxConfig struct {
	DefaultTheme string          `json:"defaultTheme"`
	Languages    map[string]bool `json:"languages"`
	Enabled      bool            `json:"enabled"`
}

// ConfigManager manages application configuration
type ConfigManager struct {
	config     *Configuration
	configPath string
}

// NewConfigManager creates a new configuration manager
func NewConfigManager() *ConfigManager {
	configDir, err := getConfigDir()
	if err != nil {
		configDir = "."
	}
	
	configPath := filepath.Join(configDir, "goeditor.json")
	
	return &ConfigManager{
		config:     DefaultConfiguration(),
		configPath: configPath,
	}
}

// DefaultConfiguration returns the default configuration
func DefaultConfiguration() *Configuration {
	return &Configuration{
		Editor: EditorConfig{
			FontSize:        14,
			TabSize:         4,
			WordWrap:        false,
			ShowLineNumbers: true,
			AutoIndent:      true,
			InsertSpaces:    true,
			TrimWhitespace:  true,
		},
		UI: UIConfig{
			Theme:        "light",
			WindowWidth:  800,
			WindowHeight: 600,
			WindowX:      100,
			WindowY:      100,
			Maximized:    false,
		},
		Syntax: SyntaxConfig{
			DefaultTheme: "github",
			Languages: map[string]bool{
				"go":         true,
				"javascript": true,
				"typescript": true,
				"python":     true,
				"java":       true,
				"html":       true,
				"css":        true,
				"json":       true,
				"xml":        true,
				"yaml":       true,
			},
			Enabled: true,
		},
	}
}

// Load loads configuration from file
func (cm *ConfigManager) Load() error {
	if _, err := os.Stat(cm.configPath); os.IsNotExist(err) {
		// Config file doesn't exist, use defaults
		return nil
	}

	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var config Configuration
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	// Merge with defaults to ensure all fields are present
	cm.mergeWithDefaults(&config)
	cm.config = &config

	return nil
}

// Save saves configuration to file
func (cm *ConfigManager) Save() error {
	// Ensure config directory exists
	configDir := filepath.Dir(cm.configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(cm.config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(cm.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetConfig returns the current configuration
func (cm *ConfigManager) GetConfig() *Configuration {
	return cm.config
}

// UpdateEditorConfig updates editor configuration
func (cm *ConfigManager) UpdateEditorConfig(config EditorConfig) {
	cm.config.Editor = config
}

// UpdateUIConfig updates UI configuration
func (cm *ConfigManager) UpdateUIConfig(config UIConfig) {
	cm.config.UI = config
}

// UpdateSyntaxConfig updates syntax configuration
func (cm *ConfigManager) UpdateSyntaxConfig(config SyntaxConfig) {
	cm.config.Syntax = config
}

// GetEditorConfig returns editor configuration
func (cm *ConfigManager) GetEditorConfig() EditorConfig {
	return cm.config.Editor
}

// GetUIConfig returns UI configuration
func (cm *ConfigManager) GetUIConfig() UIConfig {
	return cm.config.UI
}

// GetSyntaxConfig returns syntax configuration
func (cm *ConfigManager) GetSyntaxConfig() SyntaxConfig {
	return cm.config.Syntax
}

// mergeWithDefaults ensures all configuration fields have values
func (cm *ConfigManager) mergeWithDefaults(config *Configuration) {
	defaults := DefaultConfiguration()

	// Merge editor config
	if config.Editor.FontSize == 0 {
		config.Editor.FontSize = defaults.Editor.FontSize
	}
	if config.Editor.TabSize == 0 {
		config.Editor.TabSize = defaults.Editor.TabSize
	}

	// Merge UI config
	if config.UI.Theme == "" {
		config.UI.Theme = defaults.UI.Theme
	}
	if config.UI.WindowWidth == 0 {
		config.UI.WindowWidth = defaults.UI.WindowWidth
	}
	if config.UI.WindowHeight == 0 {
		config.UI.WindowHeight = defaults.UI.WindowHeight
	}

	// Merge syntax config
	if config.Syntax.DefaultTheme == "" {
		config.Syntax.DefaultTheme = defaults.Syntax.DefaultTheme
	}
	if config.Syntax.Languages == nil {
		config.Syntax.Languages = defaults.Syntax.Languages
	}
}

// getConfigDir returns the configuration directory path
func getConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".config", "goeditor")
	return configDir, nil
}