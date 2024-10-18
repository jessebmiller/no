package main

import (
	"html/template"
	"os"
	"errors"
	"fmt"
	"path/filepath"
	"bufio"
)

// GenerateSite generates a static site in BuildDir from contentPath
func GenerateSite(contentPath string, buildDir string) error {
	siteTree, err := LoadSiteTree(contentPath)
	if err != nil {
		return err
	}

	templates := NewOverlayMap(make(map[string]*template.Template))
	metadata := NewOverlayMap(make(map[string]Metadata))

	buildPath := filepath.Join(contentPath, buildDir)
	os.RemoveAll(buildPath)
	os.MkdirAll(buildPath, 0774)
	return GenerateSiteNode(siteTree, buildDir, contentPath, templates, metadata)
}

// GenerateSiteNode recursively generates a content node
// by generating all pages in this node then calling itself on all the
// children updating the templates and metadata overlay maps as it
// goes
func GenerateSiteNode(
	node SiteNode,
	nodePath string,
	contentRootPath string,
	parentTemplates OverlayMap[*template.Template],
	parentMetadata OverlayMap[Metadata],
) error {
	templates := parentTemplates.Push(node.Templates)
	metadata := parentMetadata.Push(node.Meta)
	for name, page := range node.Pages {
		tmpl, err := templates.Try(name, DEFAULT_TEMPLATE_NAME)
		if err != nil {
			return errors.New(fmt.Sprintf(
				"ERROR: Could not find template for %s",
				name,
			))
		}
		nodeMeta, _ := metadata.Try(name, DEFAULT_METADATA_NAME)
		pageMeta := NewOverlayMap[any](nodeMeta).Push(page.Meta)
		content := page.Content

		path := filepath.Join(contentRootPath, nodePath, fmt.Sprintf("%s.html", name))
		os.MkdirAll(filepath.Dir(path), 0774)
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		w := bufio.NewWriter(f)
		tmpl.Execute(w, TemplateData{pageMeta, content})
		w.Flush()
	}

	for pathPart, child := range node.Children {
		fmt.Println(nodePath, pathPart)
		return GenerateSiteNode(
			child,
			filepath.Join(nodePath, pathPart),
			contentRootPath,
			templates,
			metadata,
		)
	}
	return nil
}
