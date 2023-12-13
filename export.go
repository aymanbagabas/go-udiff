package udiff

// UnifiedDiff is a unified diff.
type UnifiedDiff = unified

// ToUnifiedDiff takes a file contents and a sequence of edits, and calculates
// a unified diff that represents those edits.
func ToUnifiedDiff(fromName, toName string, content string, edits []Edit, contextLines int) (UnifiedDiff, error) {
	return toUnified(fromName, toName, content, edits, contextLines)
}
