package main

import (
	"bytes"
	"html/template"
)

type Metadata = map[string]any

// Page struct is the content and metadata from a markdown file and some frontmatter
type Page struct {
	Content  bytes.Buffer
	Meta Metadata
}

// SiteNode struct holds all the elements of a content directory
type SiteNode struct {
	Templates map[string]*template.Template
	Meta  map[string]Metadata
	Pages     map[string]Page
	Children  map[string]SiteNode
}

// NewSiteNode creates a new SiteNode with no nil maps
func NewSiteNode() SiteNode {
	return SiteNode{
		make(map[string]*template.Template),
		make(map[string]Metadata),
		make(map[string]Page),
		make(map[string]SiteNode),
	}
}

type TemplateData struct {
	Meta OverlayMap[any]
	Content  bytes.Buffer
}
