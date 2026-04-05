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

// Package config defines the YAML config format for ghset.
// Single format: describe writes it, init reads it.
package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

// Repo is the top-level config that describe writes and init reads.
type Repo struct {
	Settings Settings  `json:"settings" yaml:"settings"`
	Security Security  `json:"security" yaml:"security"`
	Labels   Labels    `json:"labels"   yaml:"labels"`
	Rulesets []Ruleset `json:"rulesets" yaml:"rulesets"`
}

// Settings maps to GitHub repository settings (PATCH /repos/{owner}/{repo}).
type Settings struct {
	Visibility               string `json:"visibility"                  yaml:"visibility"`
	HasIssues                bool   `json:"has_issues"                  yaml:"has_issues"`
	HasWiki                  bool   `json:"has_wiki"                    yaml:"has_wiki"`
	HasProjects              bool   `json:"has_projects"                yaml:"has_projects"`
	HasDiscussions           bool   `json:"has_discussions"             yaml:"has_discussions"`
	AllowSquashMerge         bool   `json:"allow_squash_merge"          yaml:"allow_squash_merge"`
	AllowMergeCommit         bool   `json:"allow_merge_commit"          yaml:"allow_merge_commit"`
	AllowRebaseMerge         bool   `json:"allow_rebase_merge"          yaml:"allow_rebase_merge"`
	DeleteBranchOnMerge      bool   `json:"delete_branch_on_merge"      yaml:"delete_branch_on_merge"`
	WebCommitSignoffRequired bool   `json:"web_commit_signoff_required" yaml:"web_commit_signoff_required"`
}

// Security maps to GitHub security settings.
// Some fields come from the repo response, others from dedicated endpoints.
type Security struct {
	SecretScanning                    bool `json:"secret_scanning"                       yaml:"secret_scanning"`
	SecretScanningPushProtection      bool `json:"secret_scanning_push_protection"       yaml:"secret_scanning_push_protection"`
	SecretScanningAIDetection         bool `json:"secret_scanning_ai_detection"          yaml:"secret_scanning_ai_detection"`
	SecretScanningNonProviderPatterns bool `json:"secret_scanning_non_provider_patterns" yaml:"secret_scanning_non_provider_patterns"`
	VulnerabilityAlerts               bool `json:"vulnerability_alerts"                  yaml:"vulnerability_alerts"`
	AutomatedSecurityFixes            bool `json:"automated_security_fixes"              yaml:"automated_security_fixes"`
}

// Labels is a list of labels to ensure on the repo.
// Existing labels with the same name are skipped.
type Labels []Label

// Label is a single GitHub label.
type Label struct {
	Name        string `json:"name"        yaml:"name"`
	Color       string `json:"color"       yaml:"color"`
	Description string `json:"description" yaml:"description"`
}

// Ruleset is a GitHub repository ruleset.
type Ruleset struct {
	Name         string        `json:"name"                    yaml:"name"`
	Target       string        `json:"target"                  yaml:"target"`
	Enforcement  string        `json:"enforcement"             yaml:"enforcement"`
	Conditions   *Conditions   `json:"conditions,omitempty"    yaml:"conditions"`
	BypassActors []BypassActor `json:"bypass_actors,omitempty" yaml:"bypass_actors"`
	Rules        []Rule        `json:"rules"                   yaml:"rules"`
}

// Conditions specifies which refs a ruleset applies to.
type Conditions struct {
	RefName *RefNameCondition `json:"ref_name,omitempty" yaml:"ref_name"`
}

// RefNameCondition filters by ref name patterns.
type RefNameCondition struct {
	Include []string `json:"include" yaml:"include"`
	Exclude []string `json:"exclude" yaml:"exclude"`
}

// BypassActor can bypass the ruleset.
type BypassActor struct {
	ActorID    int    `json:"actor_id"    yaml:"actor_id"`
	ActorType  string `json:"actor_type"  yaml:"actor_type"`
	BypassMode string `json:"bypass_mode" yaml:"bypass_mode"`
}

// Rule is a single rule within a ruleset.
// Parameters is a free-form map because different rule types have different params.
type Rule struct {
	Type       string         `json:"type"                 yaml:"type"`
	Parameters map[string]any `json:"parameters,omitempty" yaml:"parameters"`
}

// LoadFile reads and parses a YAML config from a file path.
func LoadFile(path string) (Repo, error) {
	f, err := os.Open(path) //nolint:gosec // user-provided path is expected
	if err != nil {
		return Repo{}, fmt.Errorf("opening config %s: %w", path, err)
	}
	defer f.Close()

	return Load(f)
}

// Load reads and parses a YAML config from a reader.
func Load(r io.Reader) (Repo, error) {
	var cfg Repo

	dec := yaml.NewDecoder(r)
	dec.KnownFields(true)

	if err := dec.Decode(&cfg); err != nil {
		return Repo{}, fmt.Errorf("parsing config: %w", err)
	}

	return cfg, nil
}

// Marshal serializes a Repo config to YAML bytes.
func Marshal(cfg Repo) ([]byte, error) {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return nil, fmt.Errorf("marshaling config: %w", err)
	}

	return data, nil
}

// MarshalJSON serializes a Repo config to JSON bytes (for gh api --input).
func MarshalJSON(v any) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("marshaling JSON: %w", err)
	}

	return data, nil
}
