// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package udiff_test

import (
	"strings"
	"testing"

	diff "github.com/aymanbagabas/go-udiff"
	"github.com/aymanbagabas/go-udiff/difftest"
)

func TestApplyUnifiedContextLines(t *testing.T) {
	// Issue #40: ApplyUnified fails on diffs with context lines
	cases := []struct {
		name, old, new string
	}{
		{"three_lines_middle_change", "a\nb\nc\n", "a\nB\nc\n"},
		{"one_line_change", "b\n", "B\n"},
		{"multiple_context_lines", "x\ny\nz\n1\n2\n3\n", "x\ny\nZ\n1\n2\n3\n"},
		{"longer_file", "line 1\nline 2\nline 3\nline 4\nline 5\n", "line 1\nline 2\nline MODIFIED\nline 4\nline 5\n"},
		{"first_line_change", "1\n2\n3\n4\n", "FIRST\n2\n3\n4\n"},
		{"last_line_change", "1\n2\n3\n4\n", "1\n2\n3\nLAST\n"},
		{"multi_hunk", "1\n2\n3\n4\n5\n6\n7\n8\n9\n10\n", "1\nTWO\n3\n4\n5\n6\n7\n8\nNINE\n10\n"},
		{"delete_line", "1\n2\n3\n4\n", "1\n3\n4\n"},
		{"insert_line", "1\n3\n4\n", "1\n2\n3\n4\n"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := diff.Unified("a/f", "b/f", tc.old, tc.new)
			got, err := diff.ApplyUnified(d, tc.old)
			if err != nil {
				t.Fatalf("ApplyUnified failed on diff:\n%s\nerr: %v", d, err)
			}
			if got != tc.new {
				t.Errorf("got %q, want %q", got, tc.new)
			}
		})
	}
}

func TestApplyUnifiedRoundTrip(t *testing.T) {
	for _, tc := range difftest.TestCases {
		t.Run(tc.Name, func(t *testing.T) {
			if !strings.HasSuffix(tc.In, "\n") || !strings.HasSuffix(tc.Out, "\n") {
				// ApplyUnified is line-based for newline-terminated files
				return
			}
			d := diff.Unified(difftest.FileA, difftest.FileB, tc.In, tc.Out)
			if d == "" {
				return
			}
			got, err := diff.ApplyUnified(d, tc.In)
			if err != nil {
				t.Fatalf("ApplyUnified failed: %v", err)
			}
			if got != tc.Out {
				t.Errorf("got %q, want %q", got, tc.Out)
			}
		})
	}
}
