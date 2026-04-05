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
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thumbrise/ghset/internal/config"
	"github.com/thumbrise/ghset/internal/gh"
)

// Sentinel errors for describe command.
var (
	ErrNoRepo      = errors.New("no repo specified — interactive mode not yet implemented")
	ErrInvalidRepo = errors.New("invalid repo format")
)

var describeCmd = &cobra.Command{
	Use:   "describe [repo]",
	Short: "Snapshot repo settings into YAML",
	Long: `Snapshot GitHub repository settings into YAML and print to stdout.

repo accepts owner/name or full URL https://github.com/owner/name.
No argument — interactive prompt asks for it.

Examples:
  ghset describe thumbrise/ghset > config.yml
  ghset describe https://github.com/thumbrise/ghset
  ghset describe`,
	Args: cobra.MaximumNArgs(1),
	RunE: runDescribe,
}

func runDescribe(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	client, err := gh.NewClient()
	if err != nil {
		return err
	}

	repo, err := resolveRepo(args)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "==> Describing %s...\n", repo)

	cfg, err := fetchRepo(ctx, client, repo)
	if err != nil {
		return err
	}

	data, err := config.Marshal(cfg)
	if err != nil {
		return err
	}

	fmt.Print(string(data))

	return nil
}

// resolveRepo extracts owner/name from args or prompts interactively.
func resolveRepo(args []string) (string, error) {
	if len(args) > 0 {
		return parseRepo(args[0])
	}

	//nolint:godox
	// TODO(#5): replace with huh prompt
	return "", ErrNoRepo
}

// parseRepo normalizes a repo argument to owner/name format.
// Accepts: "owner/name" or "https://github.com/owner/name".
func parseRepo(raw string) (string, error) {
	// Full URL.
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		u, err := url.Parse(raw)
		if err != nil {
			return "", fmt.Errorf("%w: %w", ErrInvalidRepo, err)
		}

		parts := strings.Split(strings.Trim(u.Path, "/"), "/")
		if len(parts) < 2 {
			return "", fmt.Errorf("%w: URL %q must contain owner/name", ErrInvalidRepo, raw)
		}

		return parts[0] + "/" + parts[1], nil
	}

	// owner/name.
	if strings.Count(raw, "/") == 1 {
		return raw, nil
	}

	return "", fmt.Errorf("%w: %q must be owner/name or a GitHub URL", ErrInvalidRepo, raw)
}

// fetchRepo fetches all repo settings from GitHub API and maps to config.Repo.
func fetchRepo(ctx context.Context, client *gh.Client, repo string) (config.Repo, error) {
	settings, security, err := fetchSettings(ctx, client, repo)
	if err != nil {
		return config.Repo{}, err
	}

	labels, err := fetchLabels(ctx, client, repo)
	if err != nil {
		return config.Repo{}, err
	}

	rulesets, err := fetchRulesets(ctx, client, repo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "    ⚠ rulesets: %v\n", err)
	}

	vulnAlerts, err := fetchVulnerabilityAlerts(ctx, client, repo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "    ⚠ vulnerability alerts: %v\n", err)
	}

	autoFix := fetchAutomatedSecurityFixes(ctx, client, repo)

	security.VulnerabilityAlerts = vulnAlerts
	security.AutomatedSecurityFixes = autoFix

	return config.Repo{
		Settings: settings,
		Security: security,
		Labels:   labels,
		Rulesets: rulesets,
	}, nil
}

// fetchSettings fetches repo settings and security from GET /repos/{owner}/{repo}.
func fetchSettings(ctx context.Context, client *gh.Client, repo string) (config.Settings, config.Security, error) {
	data, err := client.Get(ctx, "repos/"+repo)
	if err != nil {
		return config.Settings{}, config.Security{}, fmt.Errorf("fetching repo settings: %w", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return config.Settings{}, config.Security{}, fmt.Errorf("parsing repo response: %w", err)
	}

	settings := config.Settings{
		Visibility:               str(raw, "visibility"),
		HasIssues:                boolean(raw, "has_issues"),
		HasWiki:                  boolean(raw, "has_wiki"),
		HasProjects:              boolean(raw, "has_projects"),
		HasDiscussions:           boolean(raw, "has_discussions"),
		AllowSquashMerge:         boolean(raw, "allow_squash_merge"),
		AllowMergeCommit:         boolean(raw, "allow_merge_commit"),
		AllowRebaseMerge:         boolean(raw, "allow_rebase_merge"),
		DeleteBranchOnMerge:      boolean(raw, "delete_branch_on_merge"),
		WebCommitSignoffRequired: boolean(raw, "web_commit_signoff_required"),
	}

	security := config.Security{
		SecretScanning:                    securityFeature(raw, "secret_scanning"),
		SecretScanningPushProtection:      securityFeature(raw, "secret_scanning_push_protection"),
		SecretScanningAIDetection:         securityFeature(raw, "secret_scanning_ai_detection"),
		SecretScanningNonProviderPatterns: securityFeature(raw, "secret_scanning_non_provider_patterns"),
	}

	return settings, security, nil
}

// fetchLabels fetches all labels from GET /repos/{owner}/{repo}/labels.
func fetchLabels(ctx context.Context, client *gh.Client, repo string) ([]config.Label, error) {
	data, err := client.GetPaginated(ctx, fmt.Sprintf("repos/%s/labels?per_page=100", repo))
	if err != nil {
		return nil, fmt.Errorf("fetching labels: %w", err)
	}

	var raw []map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parsing labels: %w", err)
	}

	labels := make([]config.Label, 0, len(raw))
	for _, l := range raw {
		labels = append(labels, config.Label{
			Name:        str(l, "name"),
			Color:       str(l, "color"),
			Description: str(l, "description"),
		})
	}

	return labels, nil
}

// fetchRulesets fetches all rulesets with full details.
func fetchRulesets(ctx context.Context, client *gh.Client, repo string) ([]config.Ruleset, error) {
	data, err := client.GetPaginated(ctx, fmt.Sprintf("repos/%s/rulesets?includes=all", repo))
	if err != nil {
		return nil, fmt.Errorf("fetching rulesets: %w", err)
	}

	var list []map[string]any
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, fmt.Errorf("parsing rulesets list: %w", err)
	}

	rulesets := make([]config.Ruleset, 0, len(list))

	for _, item := range list {
		id := int(num(item, "id"))

		detail, err := client.Get(ctx, fmt.Sprintf("repos/%s/rulesets/%d", repo, id))
		if err != nil {
			return nil, fmt.Errorf("fetching ruleset %d: %w", id, err)
		}

		var rs config.Ruleset
		if err := json.Unmarshal(detail, &rs); err != nil {
			return nil, fmt.Errorf("parsing ruleset %d: %w", id, err)
		}

		rulesets = append(rulesets, rs)
	}

	return rulesets, nil
}

// fetchVulnerabilityAlerts checks if vulnerability alerts are enabled.
// GET /repos/{owner}/{repo}/vulnerability-alerts → exit 0 = enabled.
func fetchVulnerabilityAlerts(ctx context.Context, client *gh.Client, repo string) (bool, error) {
	code, err := client.GetStatus(ctx, fmt.Sprintf("repos/%s/vulnerability-alerts", repo))
	if err != nil {
		return false, fmt.Errorf("checking vulnerability alerts: %w", err)
	}

	return code == 0, nil
}

// fetchAutomatedSecurityFixes checks if automated security fixes are enabled.
// Returns false if the endpoint is not supported (404) or response is malformed.
func fetchAutomatedSecurityFixes(ctx context.Context, client *gh.Client, repo string) bool {
	data, err := client.Get(ctx, fmt.Sprintf("repos/%s/automated-security-fixes", repo))
	if err != nil {
		return false
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return false
	}

	return boolean(raw, "enabled")
}

// --- JSON helpers ---

func str(m map[string]any, key string) string {
	v, _ := m[key].(string)

	return v
}

func boolean(m map[string]any, key string) bool {
	v, _ := m[key].(bool)

	return v
}

func num(m map[string]any, key string) float64 {
	v, _ := m[key].(float64)

	return v
}

// securityFeature extracts a security_and_analysis feature status.
// GitHub returns: {"security_and_analysis": {"secret_scanning": {"status": "enabled"}}}.
func securityFeature(raw map[string]any, feature string) bool {
	sa, ok := raw["security_and_analysis"].(map[string]any)
	if !ok {
		return false
	}

	feat, ok := sa[feature].(map[string]any)
	if !ok {
		return false
	}

	return str(feat, "status") == "enabled"
}
