package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"

	"go-cli-template/internal/domain"
	"go-cli-template/internal/utils"
)

// ManagerImpl loads and saves configuration files.
type ManagerImpl struct {
	cwd string
}

// NewManager returns a config manager rooted at the provided cwd.
func NewManager(cwd string) *ManagerImpl {
	return &ManagerImpl{cwd: cwd}
}

// LoadWithOverride loads config from a specific path, layered on defaults.
func (m *ManagerImpl) LoadWithOverride(path string) (domain.Config, error) {
	config := domain.DefaultConfig()
	if strings.TrimSpace(path) == "" {
		return m.Load()
	}
	partial, err := readConfig(path)
	if err != nil {
		return domain.Config{}, err
	}
	if partial != nil {
		applyPartial(&config, partial)
	}
	return config, nil
}

// Load reads config with precedence: defaults < global < local.
func (m *ManagerImpl) Load() (domain.Config, error) {
	config := domain.DefaultConfig()

	globalPath := utils.ConfigPathGlobal()
	if partial, err := readConfig(globalPath); err != nil {
		return domain.Config{}, err
	} else if partial != nil {
		applyPartial(&config, partial)
	}

	localPath := utils.ConfigPathLocal(m.cwd)
	if partial, err := readConfig(localPath); err != nil {
		return domain.Config{}, err
	} else if partial != nil {
		applyPartial(&config, partial)
	}

	return config, nil
}

// Save persists config to the global config path.
func (m *ManagerImpl) Save(config domain.Config) error {
	path := utils.ConfigPathGlobal()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := toml.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// Exists reports whether a local or global config file exists.
func (m *ManagerImpl) Exists() (bool, error) {
	globalPath := utils.ConfigPathGlobal()
	if exists, err := fileExists(globalPath); err != nil {
		return false, err
	} else if exists {
		return true, nil
	}
	localPath := utils.ConfigPathLocal(m.cwd)
	return fileExists(localPath)
}

type partialConfig struct {
	Editor               *string `toml:"editor"`
	Primary              *string `toml:"primary"`
	Secondary            *string `toml:"secondary"`
	Headings             *string `toml:"headings"`
	Text                 *string `toml:"text"`
	TextHighlight        *string `toml:"text_highlight"`
	DescriptionHighlight *string `toml:"description_highlight"`
	Tags                 *string `toml:"tags"`
	Flags                *string `toml:"flags"`
	Muted                *string `toml:"muted"`
	Accent               *string `toml:"accent"`
	Border               *string `toml:"border"`
	InteractiveDefault   *bool   `toml:"interactive_default"`
	ListSpacing          *string `toml:"list_spacing"`
}

func readConfig(path string) (*partialConfig, error) {
	if exists, err := fileExists(path); err != nil || !exists {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var partial partialConfig
	if err := toml.Unmarshal(data, &partial); err != nil {
		return nil, err
	}
	return &partial, nil
}

func applyPartial(config *domain.Config, partial *partialConfig) {
	if partial.Editor != nil {
		config.Editor = *partial.Editor
	}
	if partial.Primary != nil {
		config.Primary = *partial.Primary
	}
	if partial.Secondary != nil {
		config.Secondary = *partial.Secondary
	}
	if partial.Headings != nil {
		config.Headings = *partial.Headings
	}
	if partial.Text != nil {
		config.Text = *partial.Text
	}
	if partial.TextHighlight != nil {
		config.TextHighlight = *partial.TextHighlight
	}
	if partial.DescriptionHighlight != nil {
		config.DescriptionHighlight = *partial.DescriptionHighlight
	}
	if partial.Tags != nil {
		config.Tags = *partial.Tags
	}
	if partial.Flags != nil {
		config.Flags = *partial.Flags
	}
	if partial.Muted != nil {
		config.Muted = *partial.Muted
	}
	if partial.Accent != nil {
		config.Accent = *partial.Accent
	}
	if partial.Border != nil {
		config.Border = *partial.Border
	}
	if partial.InteractiveDefault != nil {
		config.InteractiveDefault = *partial.InteractiveDefault
	}
	if partial.ListSpacing != nil {
		config.ListSpacing = *partial.ListSpacing
	}
}

func expandPath(value string) string {
	expanded := os.ExpandEnv(value)
	if expanded == "" {
		return expanded
	}
	if expanded == "~" {
		if home, err := os.UserHomeDir(); err == nil {
			return home
		}
		return expanded
	}
	if strings.HasPrefix(expanded, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, strings.TrimPrefix(expanded, "~/"))
		}
	}
	return expanded
}

func fileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return !info.IsDir(), nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}
