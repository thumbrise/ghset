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

// Package cmd contains the Cobra command tree for ghset.
package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "ghset",
	Short: "Declarative GitHub repository settings",
	Long: `ghset — describe an existing repo into YAML, spin up a new repo from that YAML.

All GitHub API calls go through gh CLI — no tokens, no OAuth, just subprocess.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.AddCommand(describeCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(applyCmd)
}

// Execute runs the root command. Called from main.
func Execute() error {
	return rootCmd.Execute()
}
