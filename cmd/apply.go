// Copyright 2026 thumbrise
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/thumbrise/ghset/internal/config"
	"github.com/thumbrise/ghset/internal/gh"
)

var applyCmd = &cobra.Command{
	Use:   "apply [repo] [--from source]",
	Short: "Apply settings to an existing repo",
	Long: `Apply settings from a YAML config to an existing GitHub repository.

repo accepts owner/name or full URL https://github.com/owner/name.
--from accepts a YAML file path or a GitHub repo (owner/name or URL).
Config can also come from stdin (pipe).

Examples:
  ghset apply thumbrise/ghset --from opensource.yml
  ghset apply thumbrise/ghset --from thumbrise/resilience
  cat config.yml | ghset apply thumbrise/ghset`,
	Args: cobra.MaximumNArgs(1),
	RunE: runApply,
}

var applyFromFlag string

func init() {
	applyCmd.Flags().StringVarP(&applyFromFlag, "from", "f", "", "source: YAML file path or GitHub repo (owner/name or URL)")
}

func runApply(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	client, err := gh.NewClient()
	if err != nil {
		return err
	}

	cfg, err := loadFromSource(ctx, client, applyFromFlag)
	if err != nil {
		return err
	}

	repo, err := resolveApplyRepo(args)
	if err != nil {
		return err
	}

	// Best-effort: apply everything we can, collect warnings.
	var warnings []string

	fmt.Fprintf(os.Stderr, "==> Applying settings to %s...\n", repo)

	if err := applySettings(ctx, client, repo, cfg.Settings); err != nil {
		warnings = append(warnings, fmt.Sprintf("settings: %v", err))
		fmt.Fprintf(os.Stderr, "    ⚠ %v\n", err)
	}

	fmt.Fprintf(os.Stderr, "==> Applying security settings...\n")

	warnings = append(warnings, applySecurity(ctx, client, repo, cfg.Security)...)

	fmt.Fprintf(os.Stderr, "==> Applying labels...\n")

	if err := applyLabels(ctx, client, repo, cfg.Labels); err != nil {
		warnings = append(warnings, fmt.Sprintf("labels: %v", err))
		fmt.Fprintf(os.Stderr, "    ⚠ %v\n", err)
	}

	fmt.Fprintf(os.Stderr, "==> Applying rulesets...\n")

	if err := applyRulesets(ctx, client, repo, cfg.Rulesets); err != nil {
		warnings = append(warnings, fmt.Sprintf("rulesets: %v", err))
		fmt.Fprintf(os.Stderr, "    ⚠ %v\n", err)
	}

	if len(warnings) > 0 {
		fmt.Fprintf(os.Stderr, "==> Done! (%d warning(s))\n", len(warnings))
	} else {
		fmt.Fprintf(os.Stderr, "==> Done!\n")
	}

	return nil
}

// resolveApplyRepo extracts owner/name from args or prompts interactively.
func resolveApplyRepo(args []string) (string, error) {
	if len(args) > 0 {
		return parseRepo(args[0])
	}

	//nolint:godox
	// TODO(#5): replace with huh prompt
	return "", ErrNoRepo
}

// loadFromSource loads config from --from flag or stdin. Shared by init and apply.
func loadFromSource(ctx context.Context, client *gh.Client, from string) (config.Repo, error) {
	if from != "" {
		return loadFrom(ctx, client, from)
	}

	// Check if stdin has data (pipe).
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		return config.Load(os.Stdin)
	}

	//nolint:godox
	// TODO(#5): replace with huh prompt
	return config.Repo{}, ErrNoSource
}
