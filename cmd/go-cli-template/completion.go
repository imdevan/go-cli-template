package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newCompletionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion scripts",
		Long:  completionHelp(),
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCompletion(cmd, args)
		},
	}
	return cmd
}

func runCompletion(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("shell is required (bash, zsh, fish, powershell)")
	}
	shell := strings.ToLower(strings.TrimSpace(args[0]))
	switch shell {
	case "bash":
		return cmd.Root().GenBashCompletion(cmd.OutOrStdout())
	case "zsh":
		return cmd.Root().GenZshCompletion(cmd.OutOrStdout())
	case "fish":
		return cmd.Root().GenFishCompletion(cmd.OutOrStdout(), true)
	case "powershell":
		return cmd.Root().GenPowerShellCompletion(cmd.OutOrStdout())
	default:
		return fmt.Errorf("unsupported shell %q", shell)
	}
}

func completionHelp() string {
	return strings.Join([]string{
		"Examples:",
		"  go-cli-template completion bash > /etc/bash_completion.d/go-cli-template",
		"  go-cli-template completion zsh > ~/.zsh/completion/_go-cli-template",
		"  go-cli-template completion fish > ~/.config/fish/completions/go-cli-template.fish",
		"  go-cli-template completion powershell > go-cli-template.ps1",
	}, "\n")
}
