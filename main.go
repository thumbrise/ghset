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

// ghset is a CLI tool for declarative GitHub repository settings.
// Describe an existing repo into YAML, spin up a new repo from that YAML.
//
// Usage:
//
//	ghset describe [repo]          — snapshot settings → YAML stdout
//	ghset init [name] [--config f] — create repo from YAML config
package main

import (
	"fmt"
	"os"

	"github.com/thumbrise/ghset/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
