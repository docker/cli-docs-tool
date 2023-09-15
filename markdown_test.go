// Copyright 2017 cli-docs-tool authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package clidocstool

import "testing"

func TestCleanupMarkDown(t *testing.T) {
	tests := []struct {
		doc, in, expected string
	}{
		{
			doc: "whitespace around sections",
			in: `

	## Section start

Some lines.
And more lines.

`,
			expected: `## Section start

Some lines.
And more lines.`,
		},
		{
			doc: "lines with inline tabs",
			in: `## Some	Heading

A line with tabs		in it.
Tabs	should be replaced by spaces`,
			expected: `## Some    Heading

A line with tabs        in it.
Tabs    should be replaced by spaces`,
		},
		{
			doc: "lines with trailing spaces",
			in: `## Some Heading with spaces                  
       
This is a line.              
    This is an indented line        

### Some other heading         

Last line.`,
			expected: `## Some Heading with spaces

This is a line.
    This is an indented line

### Some other heading

Last line.`,
		},
		{
			doc: "lines with trailing tabs",
			in: `## Some Heading with tabs				
		
This is a line.		
	This is an indented line		

### Some other heading 	

Last line.`,
			expected: `## Some Heading with tabs

This is a line.
    This is an indented line

### Some other heading

Last line.`,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.doc, func(t *testing.T) {
			out, _ := cleanupMarkDown(tc.in)
			if out != tc.expected {
				t.Fatalf("\nexpected:\n%q\nactual:\n%q\n", tc.expected, out)
			}
		})
	}
}

func TestGetCustomAnchors(t *testing.T) {
	tests := []struct {
		in, id, expected string
	}{
		{
			in: `# Heading 1 {#heading1}`,
			id: "heading1",
		},
		{
			in: `## Heading 2 {#heading2 }`,
			id: "heading2",
		},
		{
			in: `### Heading 3 { #heading3 }`,
			id: "heading3",
		},
		{
			in: `#### Heading 4  {#heading4}`,
			id: "heading4",
		},
		{
			in: `##### Heading 5 {id=heading5}`,
			id: "heading5",
		},
		{
			in: `###### Heading 6 {id="heading6"}`,
			id: "heading6",
		},
		{
			in: `## Heading 7 { id="heading7" }`,
			id: "heading7",
		},
		{
			in: `## Heading 8 {id="heading8" }`,
			id: "heading8",
		},
		{
			in: `## Heading 9 {id="heading9"}`,
			id: "heading9",
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.in, func(t *testing.T) {
			id := getCustomAnchor(tc.in)
			if id != tc.id {
				t.Fatalf("expected: %s, actual:   %s\n", tc.id, id)
			}
		})
	}
}
