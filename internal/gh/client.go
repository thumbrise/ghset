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

// Package gh wraps the gh CLI as a subprocess for GitHub API calls.
// Auth is handled entirely by gh — no tokens in ghset.
package gh

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
)

// ErrGHNotFound is returned when the gh CLI is not installed.
var ErrGHNotFound = errors.New("gh CLI not found — install from https://cli.github.com/")

// Client calls the GitHub API via gh CLI subprocess.
type Client struct{}

// NewClient creates a Client. Verifies that gh is installed.
func NewClient() (*Client, error) {
	if _, err := exec.LookPath("gh"); err != nil {
		return nil, ErrGHNotFound
	}

	return &Client{}, nil
}

// Get calls gh api {endpoint} and returns the response body.
func (c *Client) Get(ctx context.Context, endpoint string) ([]byte, error) {
	return c.run(ctx, []string{"api", endpoint})
}

// GetStatus calls gh api {endpoint} and returns the exit code.
// Used for endpoints that signal via HTTP status (204 vs 404).
func (c *Client) GetStatus(ctx context.Context, endpoint string) (int, error) {
	cmd := exec.CommandContext(ctx, "gh", "api", endpoint) //nolint:gosec // endpoint from code, not user

	var stderr bytes.Buffer

	cmd.Stderr = &stderr

	err := cmd.Run()
	if err == nil {
		return 0, nil
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode(), nil
	}

	return -1, fmt.Errorf("gh api %s: %w", endpoint, err)
}

// Call calls gh api -X {method} {endpoint} with JSON body on stdin.
func (c *Client) Call(ctx context.Context, method, endpoint string, body any) ([]byte, error) {
	if body == nil {
		return c.run(ctx, []string{"api", "-X", method, endpoint})
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshaling body: %w", err)
	}

	return c.runWithStdin(ctx, []string{"api", "-X", method, endpoint, "--input", "-"}, jsonData)
}

// RepoCreate creates a new repository via gh repo create.
func (c *Client) RepoCreate(ctx context.Context, name string, private bool) (string, error) {
	args := []string{"repo", "create", name, "--add-readme"}

	if private {
		args = append(args, "--private")
	} else {
		args = append(args, "--public")
	}

	out, err := c.run(ctx, args)
	if err != nil {
		return "", err
	}

	return string(bytes.TrimSpace(out)), nil
}

// run executes a gh command and returns stdout.
func (c *Client) run(ctx context.Context, args []string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "gh", args...) //nolint:gosec // args from code

	var stdout, stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("gh %v: %s: %w", args, stderr.String(), err)
	}

	return stdout.Bytes(), nil
}

// runWithStdin executes a gh command with data piped to stdin.
func (c *Client) runWithStdin(ctx context.Context, args []string, stdin []byte) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "gh", args...) //nolint:gosec // args from code
	cmd.Stdin = bytes.NewReader(stdin)

	var stdout, stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("gh %v: %s: %w", args, stderr.String(), err)
	}

	return stdout.Bytes(), nil
}
