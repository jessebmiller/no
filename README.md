# no

No options, and very few features, static site generator written in Go

# PRs welcome

But I won't merge any options, or probably many features.

# Usage

To generate a site, run the following command in a content directory or subdirectory

`$ no`

Content directories are identified by being named `content` and containint a `content.tmpl` file

Every subdirectory may have any of the following files
`meta.toml`
`content.tmpl`
`[name].tmpl`
`[name].md`
`[name].toml`

For every subdirectory, an HTML page is generated for each [name].md file

The template used is chosen as follows:

1. If there is a matching `[name].tmpl` file, use that
2. Otherwise use the closest sibling or parent `content.tmpl` file

## How no builds pages

First no builds a content tree with the content and metadata for each Markdown file.
For each piece of content, the selected template generated with the following.

content       // this content
metadata      // this content's metadata
fullSite      // the full tree
subSite       // the subtree with this directory as the root

## Where does the metadata come from?

TBD, haven't written anything yet.
probably frontmatter
probably the union (closest wins for repeated keys) of all `meta.toml` files in the path and a matching `[name].toml` file
potentially filesystem or git metadata
potentially metadata encoded in the path