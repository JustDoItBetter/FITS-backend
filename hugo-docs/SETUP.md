# Hugo Documentation Setup Guide

Quick guide to get the FITS Backend documentation running.

## Quick Start

### 1. Install Hugo

```bash
# macOS
brew install hugo

# Linux
wget https://github.com/gohugoio/hugo/releases/download/v0.120.0/hugo_extended_0.120.0_linux-amd64.deb
sudo dpkg -i hugo_extended_0.120.0_linux-amd64.deb

# Windows
choco install hugo-extended

# Verify installation
hugo version
```

### 2. Install Theme

```bash
cd hugo-docs

# Create themes directory
mkdir -p themes

# Clone Geekdoc theme
git clone https://github.com/thegeeklab/hugo-geekdoc.git themes/geekdoc
```

Alternative - Download theme release:

```bash
cd hugo-docs/themes
wget https://github.com/thegeeklab/hugo-geekdoc/releases/latest/download/hugo-geekdoc.tar.gz
tar -xzf hugo-geekdoc.tar.gz
```

### 3. Run Development Server

```bash
cd hugo-docs
hugo server -D
```

Open http://localhost:1313

### 4. Build for Production

```bash
cd hugo-docs
hugo --minify
```

Static files will be in `public/` directory.

## Troubleshooting

### Theme Not Found

If you see "theme geekdoc not found":

```bash
# Make sure theme is in the right place
ls themes/geekdoc

# If empty, clone the theme
git clone https://github.com/thegeeklab/hugo-geekdoc.git themes/geekdoc
```

### Hugo Version Too Old

Update Hugo:

```bash
# macOS
brew upgrade hugo

# Linux - download latest from GitHub releases
# https://github.com/gohugoio/hugo/releases
```

### Port Already in Use

Change port:

```bash
hugo server -p 1314
```

## Deployment Options

### GitHub Pages

```yaml
# .github/workflows/deploy-docs.yml
name: Deploy Documentation

on:
  push:
    branches: [main]
    paths:
      - 'hugo-docs/**'

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Hugo
        uses: peaceiris/actions-hugo@v2
        with:
          hugo-version: '0.120.0'
          extended: true
      
      - name: Build
        run: |
          cd hugo-docs
          hugo --minify
      
      - name: Deploy
        uses: peaceiris/actions-gh-pages@v3
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./hugo-docs/public
```

### Netlify

Create `netlify.toml`:

```toml
[build]
  base = "hugo-docs"
  publish = "public"
  command = "hugo --minify"

[build.environment]
  HUGO_VERSION = "0.120.0"
```

### Docker

```dockerfile
FROM klakegg/hugo:ext-alpine AS build

WORKDIR /src
COPY hugo-docs /src

RUN hugo --minify

FROM nginx:alpine
COPY --from=build /src/public /usr/share/nginx/html
```

## Customization

### Change Theme

Edit `config.toml`:

```toml
theme = "your-theme-name"
```

### Update Colors

Create `assets/css/custom.css`:

```css
:root {
  --primary-color: #your-color;
}
```

### Add Logo

Place logo in `static/images/logo.png` and update config:

```toml
[params]
  logo = "/images/logo.png"
```

## Next Steps

- Read the [Hugo Documentation](https://gohugo.io/documentation/)
- Explore the [Geekdoc Theme](https://geekdocs.de/)
- Start adding content in `content/` directory

## Support

- Hugo: https://discourse.gohugo.io/
- Theme: https://github.com/thegeeklab/hugo-geekdoc/issues
