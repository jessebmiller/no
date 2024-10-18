# no

No options, and very few features, static site generator written in Go

# PRs welcome

But I won't merge any options, or probably many features.

# Usage

To generate a site, run the following command in a content directory
or subdirectory

`$ no`

Content directories are identified by being named `_content` and
containins a `_content.tmpl`

Every subdirectory may have any of the following files
`_meta.toml`
`_content.tmpl`
`[name].tmpl`
`[name].md`
`[name].toml`

`_meta.toml`, `_content.tmpl` and a directory named `_content` are
reserved.

For every subdirectory, an HTML page is generated for each [name].md
file

The template used is chosen as follows:

1. If there is a matching `[name].tmpl` file, use that
2. Otherwiseuse the closest sibling or parent `content.tmpl` file

The metadata used is chosen simalarly, but is a union of all metadata
along the path starting with the farthest ancestor `meta.toml` file,
those found in ancestor directories, then keys from closer `meta.toml`
files are added, overriding any conflicts. Finally metadata from
frontmatter is added, again overriding any conflicts.

In the future no may add system metadata like from the filesystem, or
information about other pages


