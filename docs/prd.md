# Bible Versions Multi-Instance Support - Epic Documentation

## Epic Overview

**Goal:** Add multiple Bible versions (starting with NKJV) to the Vietnamese Christian education website using Docusaurus's built-in multi-instance feature, allowing users to compare different Bible translations.

**Status:** Ready for Implementation
**Priority:** Medium
**Estimated Duration:** 3-4 days (3 stories)
**Risk Level:** Low

## Business Value

### User Benefits
- **Compare Translations:** Users can view the same Bible chapter in different versions
- **Study Enhancement:** Multiple translations provide deeper understanding of scriptures
- **Accessibility:** NKJV serves English-speaking users while maintaining Vietnamese content

### Technical Benefits
- **Extensible Architecture:** Easy to add more Bible versions in the future
- **Native Docusaurus Features:** Leverages tested, maintained multi-instance functionality
- **Maintainability:** Uses Docusaurus best practices and standard patterns

## Current System Context

### Existing Architecture
- **Framework:** Docusaurus 3.9.1 with TypeScript
- **Content:** Vietnamese Bible (VI1934) in `VI1934/` directory
- **Navigation:** Auto-generated sidebar from directory structure
- **Language:** Vietnamese (vi) as default locale
- **Current Route:** Single instance at root path ("/")

### Current Configuration
```typescript
// Current single instance setup
docs: {
    sidebarPath: "./sidebars.ts",
    path: "VI1934",
    routeBasePath: "/",
    sidebarCollapsible: true,
}
```

## Technical Approach

### Docusaurus Multi-Instance Feature
Leveraging Docusaurus's built-in multi-instance docs capability:
- **Multiple `@docusaurus/plugin-content-docs` instances** with unique IDs
- **Separate routing** for each Bible version
- **Built-in version switching** using `docsPluginId` navbar items
- **CLI commands** for version management (`docs:version:vi1934`, `docs:version:nkjv`)

### Architecture Changes
```
Current:                    Target:
┌─────────────────┐         ┌─────────────────┐
│ Single Instance │         │ Multiple         │
│ VI1934 only     │         │ Instances       │
│ Route: "/"      │         │ VI1934: "/"     │
└─────────────────┘         │ NKJV: "/nkjv"    │
                           └─────────────────┘
```

## User Stories

### Story 1: Docusaurus Multi-Instance Plugin Configuration
**As a website administrator, I want to configure multiple Docusaurus docs plugin instances with unique IDs, So that I can leverage Docusaurus's built-in multi-instance and version switching features.**

**Acceptance Criteria:**
- [ ] Configure VI1934 as first docs instance with `id: "vi1934"`
- [ ] Configure NKJV as second docs instance with `id: "nkjv"`
- [ ] Both instances use separate sidebar configurations
- [ ] Existing VI1934 functionality preserved at root route
- [ ] NKJV accessible at "/nkjv" route
- [ ] CLI commands work for both instances

### Story 2: NKJV Multi-Instance Content Structure
**As a content manager, I want to create NKJV Bible content using Docusaurus multi-instance patterns, So that the NKJV version integrates seamlessly with Docusaurus's version management system.**

**Acceptance Criteria:**
- [ ] Create NKJV directory structure matching Docusaurus multi-instance requirements
- [ ] Sample NKJV content (at least 2 books) with proper frontmatter
- [ ] Category files and navigation structure following Docusaurus patterns
- [ ] NKJV content structure compatible with Docusaurus version management
- [ ] Sidebar generation works correctly for NKJV instance
- [ ] Version isolation maintained between VI1934 and NKJV

### Story 3: Built-in Version Switching Integration
**As a website user, I want to use Docusaurus's built-in version switching to navigate between Bible versions, So that I can easily compare different translations using native Docusaurus UI components.**

**Acceptance Criteria:**
- [ ] Add version dropdown to navbar using `docsPluginId` items
- [ ] Version switching maintains current book/chapter context when possible
- [ ] Dropdown displays in Vietnamese language following existing UI patterns
- [ ] Existing navbar items and functionality preserved
- [ ] Version dropdown uses Docusaurus's built-in switching components
- [ ] Mobile responsiveness maintained for version switching

## Implementation Details

### Story 1: Plugin Configuration Changes

#### docusaurus.config.ts Modifications
```typescript
presets: [
    [
        "classic", {
        docs: [
            {
                id: "vi1934",
                path: "VI1934",
                routeBasePath: "/",
                sidebarPath: "./sidebars-vi1934.ts",
                versionName: "VI1934",
                disableVersioning: true,
            },
            {
                id: "nkjv",
                path: "NKJV",
                routeBasePath: "/nkjv",
                sidebarPath: "./sidebars-nkjv.ts",
                versionName: "NKJV",
                disableVersioning: true,
            }
        ],
        // ... rest of config
    }
],
```

#### Files to Create/Modify
- `docusaurus.config.ts` - Update to multi-instance configuration
- `sidebars-vi1934.ts` - Rename existing `sidebars.ts`
- `sidebars-nkjv.ts` - New NKJV sidebar configuration

### Story 2: Content Structure

#### Directory Structure
```
NKJV/
├── intro.md              # NKJV introduction
├── 1co/                  # 1 Corinthians
│   ├── 1.md             # Chapter 1
│   ├── _category_.json  # Category metadata
├── 1gi/                  # 1 John
│   ├── 1.md             # Chapter 1
│   └── _category_.json  # Category metadata
```

#### Sample NKJV Content (NKJV/1co/1.md)
```markdown
---
title: 1 Corinthians 1
book: 1 Corinthians
chapter: 1
translation: NKJV
language: EN
slug: /nkjv/1co/1
sidebar_position: 1
---

# Chapter 1

## Greeting

¹ Paul, called to be an apostle of Jesus Christ through the will of God...
```

### Story 3: Version Switching

#### Navbar Configuration
```typescript
navbar: {
    items: [
        {
            type: "docsVersionDropdown",
            position: "left",
            dropdownItemsBefore: [
                {
                    type: "html",
                    value: '<div class="dropdown-header">Phiên Bản Kinh Thánh</div>',
                },
            ],
        },
    ],
},
```

#### Localization Support
- Add Vietnamese translation for version dropdown labels
- Custom CSS styling for enhanced UX (optional)

## Risk Assessment

### Risk Matrix
| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Breaking existing VI1934 functionality | Low | High | Maintain existing route config, test thoroughly |
| Multi-instance configuration errors | Medium | Medium | Follow Docusaurus documentation exactly |
| NKJV content structure issues | Low | Medium | Mirror VI1934 patterns exactly |
| Version switching UX problems | Low | Low | Use Docusaurus built-in components |

### Rollback Plan
1. **Revert docusaurus.config.ts** to single instance configuration
2. **Remove NKJV directory** if content issues arise
3. **Restore original sidebar.ts** if sidebar problems occur

## Testing Strategy

### Unit Tests
- [ ] Configuration validation
- [ ] Route generation
- [ ] Sidebar rendering

### Integration Tests
- [ ] Multi-instance build process
- [ ] Version switching functionality
- [ ] Content accessibility

### User Acceptance Testing
- [ ] Navigate VI1934 content (regression test)
- [ ] Navigate NKJV content
- [ ] Switch between versions
- [ ] Test mobile responsiveness
- [ ] Verify Vietnamese language support

## Success Criteria

### Must-Have
- [ ] Both VI1934 and NKJV versions accessible
- [ ] Version switching works seamlessly
- [ ] No regression in existing functionality
- [ ] Mobile responsive design
- [ ] Vietnamese language support maintained

### Should-Have
- [ ] Context preservation when switching versions
- [ ] Enhanced version dropdown styling
- [ ] Comprehensive NKJV content sample

### Could-Have
- [ ] Smart book/chapter matching between versions
- [ ] Additional Bible versions beyond NKJV
- [ ] Advanced comparison features

## Future Enhancements

### Phase 2: Additional Bible Versions
- KJV (King James Version)
- NIV (New International Version)
- ESV (English Standard Version)

### Phase 3: Advanced Features
- Side-by-side version comparison
- Cross-version search
- Reading plan integration
- Audio support for additional versions

## Dependencies

### External Dependencies
- Docusaurus 3.9.1+ multi-instance support
- Node.js 22.x
- TypeScript 4.x+

### Internal Dependencies
- Existing VI1934 content structure
- Current theme configuration
- Vietnamese language setup

## Deployment Notes

### Build Commands
```bash
# Standard build
npm run build

# Type checking
npm run typecheck

# Development testing
npm start
```

### Deployment Considerations
- Increased build time due to multiple instances
- Larger bundle size with additional content
- No breaking changes to existing deployment pipeline

## Timeline and Milestones

### Week 1
- **Day 1:** Story 1 implementation and testing
- **Day 2:** Story 2 implementation and testing
- **Day 3:** Story 3 implementation and testing
- **Day 4:** Integration testing and bug fixes
- **Day 5:** User acceptance testing and documentation

### Milestone Reviews
- **End of Story 1:** Multi-instance configuration complete
- **End of Story 2:** NKJV content structure ready
- **End of Story 3:** Version switching functional
- **End of Epic:** Full feature ready for production

## Approval

**Product Owner:** _________________________
**Technical Lead:** _________________________
**Date:** _________________________

---

*This epic documentation was generated using the Brownfield Epic creation process and is ready for implementation.*