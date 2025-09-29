# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Vietnamese Christian education website built with Docusaurus that provides Bible study resources. The site serves as a comprehensive platform for studying the Bible in Vietnamese, featuring the 1925 Vietnamese Bible translation and various theological materials.

## Development Commands

```bash
# Development server (runs on http://localhost:3000)
npm start

# Build production site
npm run build

# Type checking
npm run typecheck

# Serve build locally for testing
npm run serve

# Clear build cache
npm run clear

# Docusaurus utilities
npm run swizzle    # Customize Docusaurus components
npm run write-translations
npm run write-heading-ids
```

## Architecture

### Docusaurus Configuration
- **Framework**: Docusaurus 3.9.1 with TypeScript
- **Content Structure**: Uses docs plugin with sidebar navigation
- **Language**: Vietnamese (vi) as default locale
- **Base Path**: Serves from root (`/`) not `/docs`

### Content Organization
- **Main Content**: Located in `VI1934/` directory
- **Sidebar**: Auto-generated from directory structure (`sidebars.ts`)
- **Introduction**: `VI1934/intro.md` serves as the homepage (slug: `/`)
- **Book Structure**: Each book has its own directory (e.g., `1co/`, `1gi/`, `1su/`) with numbered chapters

### Key Components
- **docusaurus.config.ts**: Main configuration with Vietnamese language settings
- **sidebars.ts**: Auto-generated sidebar configuration
- **Theme**: Classic theme with custom CSS in `src/css/custom.css`
- **Comments**: Disqus integration for user engagement
- **SEO**: Optimized for Vietnamese content

### Content Structure Pattern
```
VI1934/
├── intro.md              # Homepage content
├── [book-code]/          # e.g., 1co (1 Corinthians), 1gi (1 John)
│   ├── 1.md             # Chapter 1
│   ├── 2.md             # Chapter 2
│   └── _category_.json  # Category configuration
```

### Deployment
- **Platform**: Vercel (configured in docusaurus.config.ts)
- **Production URL**: https://phankhoi.vercel.app
- **Build Command**: `npm run build`
- **Output Directory**: `build/`

## Content Guidelines

### File Naming
- Use standard Bible book abbreviations (e.g., `1co`, `1gi`, `1su`)
- Chapter files are numbered (e.g., `1.md`, `2.md`)
- Each book directory must have `_category_.json`

### Frontmatter
- Introduction page uses `slug: /` and `sidebar_position: -1`
- Standard Docusaurus frontmatter for other pages

### Language Considerations
- All content should be in Vietnamese
- UTF-8 encoding is required
- Use proper Vietnamese typography and diacritics

## Development Notes

### Node.js Version

- Requires Node.js 22.x (specified in package.json engines)

### TypeScript

- Type checking available with `npm run typecheck`
- Full TypeScript support with Docusaurus type definitions

### Customization

- Theme customization through `src/css/custom.css`
- Component customization available via `npm run swizzle`
- Navigation and footer links configured in `docusaurus.config.ts`
- 