//go:build build

package index

import (
	"bytes"
	"strings"

	"github.com/yuin/goldmark/ast"
)

// Header represents a single header in the hierarchy
type Header struct {
	Level    int
	Text     string
	Anchor   string
	Content  string
	Children []*Header
}

// headerInfo is used to track position of headers so we can extract their content
type headerInfo struct {
	header *Header
	node   *ast.Heading
	start  int // Start position after this heading
	end    int // End position (start of next heading or EOF)
}

// extractHeaders walks the AST and extracts all headers with their content
func extractHeaders(n ast.Node, source []byte) []*Header {
	var headers []*Header
	var headingInfos []headerInfo

	// First pass: collect all headings and their positions
	ast.Walk(n, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if heading, ok := node.(*ast.Heading); ok {
			text := extractHeaderText(heading, source)

			// Extract anchor (id attribute)
			t, _ := heading.AttributeString("id")
			anchor := string(t.([]byte))

			header := &Header{
				Level:    heading.Level,
				Text:     text,
				Anchor:   anchor,
				Content:  "",
				Children: []*Header{},
			}

			headers = append(headers, header)

			// Track position info - content starts after this heading ends
			// Use the stop position of the last line of the heading
			lastLine := heading.Lines().At(heading.Lines().Len() - 1)
			headingInfos = append(headingInfos, headerInfo{
				header: header,
				node:   heading,
				start:  lastLine.Stop,
			})
		}
		return ast.WalkContinue, nil
	})

	// TODO: Content parsing is disabled for now.
	// // Second pass: determine content ranges and extract content
	// for i := range headingInfos {
	// 	// Find the end position (start of next heading at same or higher level)
	// 	endPos := len(source)
	// 	currentLevel := headingInfos[i].header.Level

	// 	for j := i + 1; j < len(headingInfos); j++ {
	// 		nextLevel := headingInfos[j].header.Level
	// 		// Stop at next heading of same or higher level (lower level number = higher in hierarchy)
	// 		if nextLevel <= currentLevel {
	// 			// Get the start of the next heading's first line
	// 			endPos = headingInfos[j].node.Lines().At(0).Start
	// 			break
	// 		}
	// 	}

	// 	headingInfos[i].end = endPos

	// 	// Extract content between start and end positions
	// 	contentBytes := source[headingInfos[i].start:headingInfos[i].end]
	// 	headingInfos[i].header.Content = cleanContent(contentBytes)
	// }

	return headers
}

// cleanContent removes leading/trailing whitespace and normalizes the content
func cleanContent(content []byte) string {
	// Trim whitespace
	trimmed := bytes.TrimSpace(content)

	// Convert to string and normalize newlines
	str := string(trimmed)

	// Optional: collapse multiple newlines into double newlines
	str = strings.ReplaceAll(str, "\n\n\n", "\n\n")

	return str
}

// extractHeaderText extracts the text content from a heading node
func extractHeaderText(heading *ast.Heading, source []byte) string {
	var buf bytes.Buffer
	for child := heading.FirstChild(); child != nil; child = child.NextSibling() {
		if textNode, ok := child.(*ast.Text); ok {
			buf.Write(textNode.Segment.Value(source))
		}
	}
	return buf.String()
}

// buildHierarchy converts a flat list of headers into a hierarchical structure
func buildHierarchy(headers []*Header) []*Header {
	if len(headers) == 0 {
		return []*Header{}
	}

	var root []*Header
	var stack []*Header

	for _, header := range headers {
		// Pop stack until we find a parent with lower level
		for len(stack) > 0 && stack[len(stack)-1].Level >= header.Level {
			stack = stack[:len(stack)-1]
		}

		if len(stack) == 0 {
			// Top-level header
			root = append(root, header)
		} else {
			// Add as child to the parent
			parent := stack[len(stack)-1]
			parent.Children = append(parent.Children, header)
		}

		// Push current header onto stack
		stack = append(stack, header)
	}

	return root
}
