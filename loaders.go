package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"github.com/yuin/goldmark"
	"go.abhg.dev/goldmark/frontmatter"
	"html/template"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func LoadPage(path string) (Page, error) {
	source, err := ioutil.ReadFile(path)
	if err != nil {
		return Page{}, err
	}
	md := goldmark.New(
		goldmark.WithExtensions(&frontmatter.Extender{}),
	)

	var buf bytes.Buffer
	if err := md.Convert(source, &buf); err != nil {
		return Page{}, err
	}

	var meta Metadata

	return Page{Content: buf, Meta: meta}, nil
}

func LoadMetadata(path string) (Metadata, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var m Metadata
	err = toml.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func LoadSiteTree(rootTreePath string) (SiteNode, error) {
	if filepath.Base(rootTreePath) != CONTENT_FOLDER {
		return SiteNode{}, errors.New(fmt.Sprintf(
			"ERROR: %s is not a content root, expected folder to be named %s",
			rootTreePath,
			CONTENT_FOLDER,
		))
	}
	defaultTemplatePath := filepath.Join(rootTreePath, CONTENT_FOLDER_DEFAULT_TEMPLATE)
	if _, err := os.Stat(defaultTemplatePath); errors.Is(err, os.ErrNotExist) {
		return SiteNode{}, errors.New(fmt.Sprintf(
			"ERROR: %s is not a content root, expected it to contain %s",
			rootTreePath,
			CONTENT_FOLDER_DEFAULT_TEMPLATE,
		))
	}

	parentTree := NewSiteNode()
	tree := &parentTree
	err := filepath.WalkDir(rootTreePath, func(path string, d fs.DirEntry, err error) error {
		if d.Name() == CONTENT_FOLDER {
			return nil
		}
		dirName := filepath.Base(filepath.Dir(path))

		var workingTree *SiteNode
		existingTree, ok := tree.Children[dirName]
		if !ok {
			t := NewSiteNode()
			tree.Children[dirName] = t
			workingTree = &t
		} else {
			workingTree = &existingTree
		}

		if d.IsDir() {
			workingTree.Children[d.Name()] = NewSiteNode()
			tree = workingTree
			return nil
		}

		ext := filepath.Ext(d.Name())
		name := strings.TrimSuffix(d.Name(), ext)
		switch ext {
		case ".toml":
			m, err := LoadMetadata(path)
			if err != nil {
				return err
			}
			workingTree.Meta[name] = m
		case ".md":
			p, err := LoadPage(path)
			if err != nil {
				return err
			}
			workingTree.Pages[name] = p
		case ".tmpl":
			t, err := template.ParseFiles(path)
			if err != nil {
				return err
			}
			workingTree.Templates[name] = t
		default:
			fmt.Println("No ignoring", d.Name(), ext == "")
		}

		return nil
	})
	if err != nil {
		return SiteNode{}, err
	}
	return parentTree.Children[CONTENT_FOLDER], nil
}
