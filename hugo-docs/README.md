# FITS Backend Documentation

This directory contains the Hugo-based documentation for FITS Backend.

## Prerequisites

- Hugo Extended v0.110.0 or higher
- Git

## Installation

### Install Hugo

#### macOS

```bash
brew install hugo
```

#### Linux (Ubuntu/Debian)

```bash
# Download latest Hugo extended
wget https://github.com/gohugoio/hugo/releases/download/v0.120.0/hugo_extended_0.120.0_linux-amd64.deb
sudo dpkg -i hugo_extended_0.120.0_linux-amd64.deb
```

#### Windows (Chocolatey)

```powershell
choco install hugo-extended
```

### Install Theme

This documentation uses the Geekdoc theme:

```bash
cd hugo-docs

# Initialize Hugo module
hugo mod init github.com/JustDoItBetter/FITS-backend

# Download theme
mkdir -p themes
cd themes
git clone https://github.com/thegeeklab/hugo-geekdoc.git geekdoc
```

## Development

### Run Local Server

```bash
cd hugo-docs
hugo server -D
```

Open http://localhost:1313 in your browser.

The server will automatically reload when you make changes.

### Build Static Site

```bash
hugo --minify
```

Output will be in the `public/` directory.

## Project Structure

```
hugo-docs/
├── config.toml              # Hugo configuration
├── content/                 # Markdown content
│   ├── _index.md           # Homepage
│   ├── project-overview/   # Architecture & design
│   ├── getting-started/    # Installation & quick start
│   ├── development/        # Dev guides
│   ├── api/                # API documentation
│   └── infrastructure/     # Deployment & security
├── static/                 # Static assets
│   └── images/            # Diagrams & images
├── layouts/               # Custom templates
├── themes/               # Hugo themes
└── public/              # Generated site (gitignored)
```

## Adding Content

### Create New Page

```bash
hugo new project-overview/new-page.md
```

### Front Matter Template

```yaml
---
title: "Page Title"
description: "Brief description"
weight: 1  # Order in menu
date: 2025-10-23
---

# Page Content

Your markdown content here.
```

## Deploying

### GitHub Pages

```bash
# Build site
hugo --minify

# Deploy to gh-pages branch
git subtree push --prefix hugo-docs/public origin gh-pages
```

### Netlify

1. Push to GitHub
2. Connect repository in Netlify
3. Configure build:
   - Build command: `cd hugo-docs && hugo --minify`
   - Publish directory: `hugo-docs/public`

### Custom Server

```bash
# Build site
hugo --minify

# Copy public/ directory to web server
rsync -avz public/ user@server:/var/www/docs/
```

## Maintenance

### Update Theme

```bash
cd themes/geekdoc
git pull origin main
```

### Update Content

1. Edit markdown files in `content/`
2. Test locally with `hugo server`
3. Commit and push changes
4. Deploy

## Links

- [Hugo Documentation](https://gohugo.io/documentation/)
- [Geekdoc Theme](https://geekdocs.de/)
- [Markdown Guide](https://www.markdownguide.org/)

## Contributing

See [Contribution Guidelines](/development/contribution-guidelines/) for how to contribute to the documentation.
