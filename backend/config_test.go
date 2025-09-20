package backend

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfiguration(t *testing.T) {
	config := DefaultConfiguration()
	
	// Test editor defaults
	if config.Editor.FontSize != 14 {
		t.Errorf("Expected FontSize 14, got %d", config.Editor.FontSize)
	}
	
	if config.Editor.TabSize != 4 {
		t.Errorf("Expected TabSize 4, got %d", config.Editor.TabSize)
	}
	
	if !config.Editor.ShowLineNumbers {
		t.Error("Expected ShowLineNumbers to be true")
	}
	
	// Test UI defaults
	if config.UI.Theme != "light" {
		t.Errorf("Expected Theme 'light', got '%s'", config.UI.Theme)
	}
	
	if config.UI.WindowWidth != 800 {
		t.Errorf("Expected WindowWidth 800, got %d", config.UI.WindowWidth)
	}
	
	// Test syntax defaults
	if config.Syntax.DefaultTheme != "github" {
		t.Errorf("Expected DefaultTheme 'github', got '%s'", config.Syntax.DefaultTheme)
	}
	
	if !config.Syntax.Enabled {
		t.Error("Expected Syntax.Enabled to be true")
	}
}

func TestConfigManagerSaveLoad(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "goeditor_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create config manager with custom path
	cm := &ConfigManager{
		config:     DefaultConfiguration(),
		configPath: filepath.Join(tempDir, "test_config.json"),
	}
	
	// Modify some settings
	cm.config.Editor.FontSize = 16
	cm.config.UI.Theme = "dark"
	cm.config.Syntax.DefaultTheme = "monokai"
	
	// Save configuration
	if err := cm.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}
	
	// Create new config manager and load
	cm2 := &ConfigManager{
		config:     DefaultConfiguration(),
		configPath: filepath.Join(tempDir, "test_config.json"),
	}
	
	if err := cm2.Load(); err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	
	// Verify loaded values
	if cm2.config.Editor.FontSize != 16 {
		t.Errorf("Expected FontSize 16, got %d", cm2.config.Editor.FontSize)
	}
	
	if cm2.config.UI.Theme != "dark" {
		t.Errorf("Expected Theme 'dark', got '%s'", cm2.config.UI.Theme)
	}
	
	if cm2.config.Syntax.DefaultTheme != "monokai" {
		t.Errorf("Expected DefaultTheme 'monokai', got '%s'", cm2.config.Syntax.DefaultTheme)
	}
}

func TestConfigManagerMergeWithDefaults(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "goeditor_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	configPath := filepath.Join(tempDir, "partial_config.json")
	
	// Create partial config file
	partialConfig := map[string]interface{}{
		"editor": map[string]interface{}{
			"fontSize": 18,
		},
		"ui": map[string]interface{}{
			"theme": "dark",
		},
	}
	
	data, err := json.MarshalIndent(partialConfig, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal partial config: %v", err)
	}
	
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		t.Fatalf("Failed to write partial config: %v", err)
	}
	
	// Load with config manager
	cm := &ConfigManager{
		config:     DefaultConfiguration(),
		configPath: configPath,
	}
	
	if err := cm.Load(); err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	
	// Check that partial values are loaded
	if cm.config.Editor.FontSize != 18 {
		t.Errorf("Expected FontSize 18, got %d", cm.config.Editor.FontSize)
	}
	
	if cm.config.UI.Theme != "dark" {
		t.Errorf("Expected Theme 'dark', got '%s'", cm.config.UI.Theme)
	}
	
	// Check that missing values use defaults
	if cm.config.Editor.TabSize != 4 {
		t.Errorf("Expected TabSize 4 (default), got %d", cm.config.Editor.TabSize)
	}
	
	if cm.config.UI.WindowWidth != 800 {
		t.Errorf("Expected WindowWidth 800 (default), got %d", cm.config.UI.WindowWidth)
	}
}