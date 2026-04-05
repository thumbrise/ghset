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
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/thumbrise/ghset/internal/config"
	"github.com/thumbrise/ghset/internal/gh"
)

// Sentinel errors for init command.
var (
	ErrNoSource   = errors.New("no source specified — use --from (file or repo) or pipe from stdin")
	ErrNoRepoName = errors.New("no repo name specified — interactive mode not yet implemented")
)

var initCmd = &cobra.Command{
	Use:   "init [name]",
	Short: "Create repo and apply settings from YAML config",
	Long: `Create a new GitHub repository and apply all settings from a YAML config.

name is the new repo name: my-repo creates under authenticated user,
myorg/my-repo creates under org. Standard gh repo create behavior.
--from accepts a YAML file path or a GitHub repo (owner/name or URL).
No arguments — interactive prompts ask for source and repo name.
Config can also come from stdin (pipe).

Examples:
  ghset init my-new-repo --from opensource.yml
  ghset init my-new-repo --from thumbrise/resilience
  ghset init my-new-repo --from https://github.com/thumbrise/resilience
  cat config.yml | ghset init my-new-repo
  ghset init`,
	Args: cobra.MaximumNArgs(1),
	RunE: runInit,
}

var fromFlag string

func init() {
	initCmd.Flags().StringVarP(&fromFlag, "from", "f", "", "source: YAML file path or GitHub repo (owner/name or URL)")
}

//nolint:cyclop,funlen // simple pipeline
func runInit(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	client, err := gh.NewClient()
	if err != nil {
		return err
	}

	cfg, err := loadFromSource(ctx, client, fromFlag)
	if err != nil {
		return err
	}

	name, err := resolveRepoName(args)
	if err != nil {
		return err
	}

	// 1. Create repo — fatal, can't continue without it.
	fmt.Fprintf(os.Stderr, "==> Creating repo %s...\n", name)

	repoURL, err := client.RepoCreate(ctx, name, cfg.Settings.Visibility)
	if err != nil {
		return fmt.Errorf("creating repo: %w", err)
	}

	// gh repo create returns the URL, extract owner/name for API calls.
	repo, err := parseRepo(repoURL)
	if err != nil {
		// Fallback: if the name already contains owner/, use it directly.
		repo = name
	}

	// Wait for GitHub to fully initialize the repo.
	if err := waitForRepo(ctx, client, repo); err != nil {
		return err
	}

	// Best-effort: apply everything we can, collect warnings.
	var warnings []string

	// 2. Apply settings.
	fmt.Fprintf(os.Stderr, "==> Applying settings...\n")

	if err := applySettings(ctx, client, repo, cfg.Settings); err != nil {
		warnings = append(warnings, fmt.Sprintf("settings: %v", err))
		fmt.Fprintf(os.Stderr, "    ⚠ %v\n", err)
	}

	// 3. Apply security.
	fmt.Fprintf(os.Stderr, "==> Applying security settings...\n")

	warnings = append(warnings, applySecurity(ctx, client, repo, cfg.Security)...)

	// 4. Apply labels.
	fmt.Fprintf(os.Stderr, "==> Applying labels...\n")

	if err := applyLabels(ctx, client, repo, cfg.Labels); err != nil {
		warnings = append(warnings, fmt.Sprintf("labels: %v", err))
		fmt.Fprintf(os.Stderr, "    ⚠ %v\n", err)
	}

	// 5. Apply rulesets.
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

	fmt.Println(repoURL)

	return nil
}

// loadFrom auto-detects whether source is a repo or a file.
// Repo: contains "/" but no ".yml"/".yaml" extension, or starts with "http".
// File: everything else.
func loadFrom(ctx context.Context, client *gh.Client, source string) (config.Repo, error) {
	if isRepoRef(source) {
		repo, err := parseRepo(source)
		if err != nil {
			return config.Repo{}, err
		}

		fmt.Fprintf(os.Stderr, "==> Describing %s...\n", repo)

		return fetchRepo(ctx, client, repo)
	}

	return config.LoadFile(source)
}

// isRepoRef returns true if source looks like a GitHub repo reference.
func isRepoRef(source string) bool {
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		return true
	}

	if strings.HasSuffix(source, ".yml") || strings.HasSuffix(source, ".yaml") {
		return false
	}

	return strings.Contains(source, "/")
}

// resolveRepoName extracts repo name from args or prompts interactively.
func resolveRepoName(args []string) (string, error) {
	if len(args) > 0 {
		return args[0], nil
	}

	//nolint:godox
	// TODO(#5): replace with huh prompt
	return "", ErrNoRepoName
}

// applySettings patches repo settings via PATCH /repos/{owner}/{repo}.
func applySettings(ctx context.Context, client *gh.Client, repo string, s config.Settings) error {
	body := map[string]any{
		"visibility":                  s.Visibility,
		"has_issues":                  s.HasIssues,
		"has_wiki":                    s.HasWiki,
		"has_projects":                s.HasProjects,
		"has_discussions":             s.HasDiscussions,
		"allow_squash_merge":          s.AllowSquashMerge,
		"allow_merge_commit":          s.AllowMergeCommit,
		"allow_rebase_merge":          s.AllowRebaseMerge,
		"delete_branch_on_merge":      s.DeleteBranchOnMerge,
		"web_commit_signoff_required": s.WebCommitSignoffRequired,
	}

	_, err := client.Call(ctx, "PATCH", "repos/"+repo, body)
	if err != nil {
		return fmt.Errorf("applying settings: %w", err)
	}

	return nil
}

// applySecurity patches security settings and toggles dedicated endpoints.
// Returns warnings for features that are not available (e.g. secret scanning on free private repos).
func applySecurity(ctx context.Context, client *gh.Client, repo string, s config.Security) []string {
	var warnings []string

	// Security features via PATCH /repos/{owner}/{repo}.
	body := map[string]any{
		"security_and_analysis": map[string]any{
			"secret_scanning":                       map[string]string{"status": secStatus(s.SecretScanning)},
			"secret_scanning_push_protection":       map[string]string{"status": secStatus(s.SecretScanningPushProtection)},
			"secret_scanning_ai_detection":          map[string]string{"status": secStatus(s.SecretScanningAIDetection)},
			"secret_scanning_non_provider_patterns": map[string]string{"status": secStatus(s.SecretScanningNonProviderPatterns)},
		},
	}

	_, err := client.Call(ctx, "PATCH", "repos/"+repo, body)
	if err != nil {
		w := fmt.Sprintf("security settings: %v", err)
		warnings = append(warnings, w)
		fmt.Fprintf(os.Stderr, "    ⚠ %s\n", w)
	}

	// Vulnerability alerts.
	if s.VulnerabilityAlerts {
		_, err = client.Call(ctx, "PUT", fmt.Sprintf("repos/%s/vulnerability-alerts", repo), nil)
	} else {
		_, err = client.Call(ctx, "DELETE", fmt.Sprintf("repos/%s/vulnerability-alerts", repo), nil)
	}

	if err != nil {
		w := fmt.Sprintf("vulnerability alerts: %v", err)
		warnings = append(warnings, w)
		fmt.Fprintf(os.Stderr, "    ⚠ %s\n", w)
	}

	// Automated security fixes.
	if s.AutomatedSecurityFixes {
		_, err = client.Call(ctx, "PUT", fmt.Sprintf("repos/%s/automated-security-fixes", repo), nil)
	} else {
		_, err = client.Call(ctx, "DELETE", fmt.Sprintf("repos/%s/automated-security-fixes", repo), nil)
	}

	if err != nil {
		w := fmt.Sprintf("automated security fixes: %v", err)
		warnings = append(warnings, w)
		fmt.Fprintf(os.Stderr, "    ⚠ %s\n", w)
	}

	return warnings
}

// applyLabels creates labels from config. Skips labels that already exist.
func applyLabels(ctx context.Context, client *gh.Client, repo string, labels config.Labels) error {
	for _, l := range labels {
		body := map[string]string{
			"name":        l.Name,
			"color":       l.Color,
			"description": l.Description,
		}

		_, err := client.Call(ctx, "POST", fmt.Sprintf("repos/%s/labels", repo), body)
		if err != nil {
			// 422 = label already exists — skip.
			if strings.Contains(err.Error(), "already_exists") || strings.Contains(err.Error(), "422") {
				continue
			}

			return fmt.Errorf("creating label %q: %w", l.Name, err)
		}
	}

	return nil
}

// applyRulesets creates each ruleset via POST /repos/{owner}/{repo}/rulesets.
func applyRulesets(ctx context.Context, client *gh.Client, repo string, rulesets []config.Ruleset) error {
	for _, rs := range rulesets {
		_, err := client.Call(ctx, "POST", fmt.Sprintf("repos/%s/rulesets", repo), rs)
		if err != nil {
			return fmt.Errorf("creating ruleset %q: %w", rs.Name, err)
		}
	}

	return nil
}

var ErrRepoNotReady = errors.New("repo not ready")

// waitForRepo polls until the repo is accessible via API.
func waitForRepo(ctx context.Context, client *gh.Client, repo string) error {
	for range 5 {
		_, err := client.Get(ctx, "repos/"+repo)
		if err == nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
		}
	}

	return fmt.Errorf("%w after 5 seconds: %s", ErrRepoNotReady, repo)
}

func secStatus(enabled bool) string {
	if enabled {
		return "enabled"
	}

	return "disabled"
}
