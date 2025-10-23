# Quick Start - FITS Backend Documentation

Get the documentation running in 3 commands.

## Prerequisites

- Hugo Extended 0.110.0+ ([Install Guide](https://gohugo.io/installation/))

## Setup

```bash
# 1. Navigate to docs directory
cd hugo-docs

# 2. Download theme
git clone https://github.com/thegeeklab/hugo-geekdoc.git themes/geekdoc

# 3. Start server
hugo server -D
```

Open http://localhost:1313

## Build for Production

```bash
hugo --minify
```

Output: `public/` directory

## Documentation Structure

```
📚 FITS Backend Documentation
│
├── 🏗️  Project Overview
│   ├── Architecture          # System design & components
│   └── System Design         # Design principles & patterns
│
├── 🚀 Getting Started
│   ├── Installation          # Setup instructions
│   └── Quick Start           # First API requests
│
├── 💻 Development
│   ├── Local Setup           # Dev environment
│   ├── Testing Strategy      # Test guidelines
│   └── Contribution          # How to contribute
│
├── 🔌 API Reference
│   └── Authentication        # JWT & security
│
└── 🏭 Infrastructure
    ├── Deployment            # Production setup
    └── Security              # Security checklist
```

## Quick Links

- **Homepage**: http://localhost:1313
- **Swagger UI**: http://localhost:8080/docs (when API running)
- **Architecture**: /project-overview/architecture/
- **Quick Start Guide**: /getting-started/quick-start/

## Troubleshooting

### Theme not found

```bash
# Download theme
cd hugo-docs
git clone https://github.com/thegeeklab/hugo-geekdoc.git themes/geekdoc
```

### Port in use

```bash
hugo server -p 1314
```

## Deployment

### GitHub Pages

Add to `.github/workflows/docs.yml`:

```yaml
- name: Build docs
  run: cd hugo-docs && hugo --minify
- name: Deploy
  uses: peaceiris/actions-gh-pages@v3
  with:
    publish_dir: ./hugo-docs/public
```

### Netlify

Push to GitHub, connect in Netlify:
- Build: `cd hugo-docs && hugo --minify`
- Publish: `hugo-docs/public`

## Need Help?

- Hugo Docs: https://gohugo.io/documentation/
- Theme Docs: https://geekdocs.de/
- Project Issues: https://github.com/JustDoItBetter/FITS-backend/issues

Happy documenting! 📖
