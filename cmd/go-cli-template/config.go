package main

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"go-cli-template/internal/adapters/editor"
	"go-cli-template/internal/config"
	"go-cli-template/internal/domain"
	"go-cli-template/internal/utils"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "View or edit configuration",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfig(cmd)
		},
	}
	cmd.AddCommand(newConfigInitCmd())
	return cmd
}

func runConfig(cmd *cobra.Command) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	manager := config.NewManager(cwd)
	cfg, err := manager.Load()
	if err != nil {
		return err
	}

	path, err := resolveConfigPath(cwd)
	if err != nil {
		return err
	}

	if !pathExists(path) {
		cfg = domain.DefaultConfig()
		content := renderConfigTemplate(cfg)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			return err
		}
	}

	editorAdapter := editor.New(cfg.Editor)
	if err := editorAdapter.Open(path); err != nil {
		return err
	}
	cmd.Printf("Opened config %s\n", path)
	return nil
}

func resolveConfigPath(cwd string) (string, error) {
	localPath := utils.ConfigPathLocal(cwd)
	if pathExists(localPath) {
		return localPath, nil
	}
	globalPath := utils.ConfigPathGlobal()
	if pathExists(globalPath) {
		return globalPath, nil
	}
	return globalPath, nil
}

func pathExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
