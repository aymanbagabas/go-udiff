// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package udiff_test

import (
	"fmt"
	"strings"
	"testing"

	diff "github.com/aymanbagabas/go-udiff"
	"github.com/aymanbagabas/go-udiff/difftest"
)

func TestApplyUnifiedContextLines(t *testing.T) {
	tests := []struct {
		name string
		old  string
		new  string
	}{
		{
			name: "issue 40 three-line file with middle change",
			old:  "a\nb\nc\n",
			new:  "a\nB\nc\n",
		},
		{
			name: "single line file without context",
			old:  "b\n",
			new:  "B\n",
		},
		{
			name: "change at start of file with following context",
			old:  "first\nsecond\nthird\nfourth\n",
			new:  "FIRST\nsecond\nthird\nfourth\n",
		},
		{
			name: "change at end of file with preceding context",
			old:  "first\nsecond\nthird\nfourth\n",
			new:  "first\nsecond\nthird\nFOURTH\n",
		},
		{
			name: "deletion in middle with context",
			old:  "line1\nline2\nline3\nline4\nline5\n",
			new:  "line1\nline2\nline4\nline5\n",
		},
		{
			name: "insertion in middle with context",
			old:  "line1\nline2\nline4\nline5\n",
			new:  "line1\nline2\nline3\nline4\nline5\n",
		},
		{
			name: "empty context lines",
			old:  "line1\n\nline3\nline4\n",
			new:  "line1\n\nLINE3\nline4\n",
		},
		{
			name: "multiple hunks separated by unchanged lines",
			old: strings.Join([]string{
				"line 1",
				"line 2",
				"line 3",
				"line 4",
				"line 5",
				"line 6",
				"line 7",
				"line 8",
				"line 9",
				"line 10",
				"line 11",
				"line 12",
				"line 13",
				"line 14",
				"line 15",
				"",
			}, "\n"),
			new: strings.Join([]string{
				"line 1",
				"MODIFIED 2",
				"line 3",
				"line 4",
				"line 5",
				"line 6",
				"line 7",
				"line 8",
				"line 9",
				"line 10",
				"line 11",
				"line 12",
				"line 13",
				"MODIFIED 14",
				"line 15",
				"",
			}, "\n"),
		},
		{
			name: "identical files",
			old:  "alpha\nbeta\ngamma\n",
			new:  "alpha\nbeta\ngamma\n",
		},
		{
			name: "completely replaced content",
			old:  "old1\nold2\nold3\n",
			new:  "new1\nnew2\nnew3\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ud := diff.Unified("a/file", "b/file", tc.old, tc.new)
			got, err := diff.ApplyUnified(ud, tc.old)
			if err != nil {
				t.Fatalf("ApplyUnified failed: %v\nUnified diff:\n%s", err, ud)
			}
			if got != tc.new {
				t.Errorf("ApplyUnified mismatch:\ngot:\n%q\nwant:\n%q\nUnified diff:\n%s", got, tc.new, ud)
			}
		})
	}
}

func TestApplyUnifiedVariousContextSizes(t *testing.T) {
	old := "1\n2\n3\n4\n5\n6\n7\n8\n9\n10\n"
	new := "1\n2\n3\nFOUR\n5\n6\n7\nEIGHT\n9\n10\n"
	edits := diff.Strings(old, new)

	for _, ctx := range []int{0, 1, 2, 3, 5} {
		t.Run(fmt.Sprintf("context_%d", ctx), func(t *testing.T) {
			ud, err := diff.ToUnified("a/file", "b/file", old, edits, ctx)
			if err != nil {
				t.Fatalf("ToUnified failed: %v", err)
			}
			got, err := diff.ApplyUnified(ud, old)
			if err != nil {
				t.Fatalf("ApplyUnified failed with context %d: %v\nDiff:\n%s", ctx, err, ud)
			}
			if got != new {
				t.Fatalf("ApplyUnified mismatch with context %d: got %q, want %q", ctx, got, new)
			}
		})
	}
}

func TestApplyUnifiedDifftestCases(t *testing.T) {
	for _, tc := range difftest.TestCases {
		if !strings.HasSuffix(tc.In, "\n") || !strings.HasSuffix(tc.Out, "\n") {
			// ApplyUnified operates on line-delimited content.
			continue
		}
		t.Run(tc.Name, func(t *testing.T) {
			ud := diff.Unified("a", "b", tc.In, tc.Out)
			got, err := diff.ApplyUnified(ud, tc.In)
			if err != nil {
				t.Fatalf("ApplyUnified failed: %v\nUnified diff:\n%s", err, ud)
			}
			if got != tc.Out {
				t.Errorf("ApplyUnified mismatch:\ngot:\n%q\nwant:\n%q\nUnified diff:\n%s", got, tc.Out, ud)
			}
		})
	}
}

func TestApplyUnifiedErrors(t *testing.T) {
	t.Run("invalid hunk header", func(t *testing.T) {
		invalidDiff := "@@ invalid @@\n"
		_, err := diff.ApplyUnified(invalidDiff, "a\n")
		if err == nil {
			t.Fatal("expected error for invalid hunk header, got nil")
		}
	})

	t.Run("unexpected prefix in diff", func(t *testing.T) {
		invalidDiff := "--- a\n+++ b\n@@ -1,1 +1,1 @@\n?unknown\n"
		_, err := diff.ApplyUnified(invalidDiff, "a\n")
		if err == nil {
			t.Fatal("expected error for unexpected prefix, got nil")
		}
	})
}
