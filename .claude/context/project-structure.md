---
created: 2025-09-27T18:57:36Z
last_updated: 2025-09-27T18:57:36Z
version: 1.0
author: Claude Code PM System
---

# Project Structure

## Root Directory Organization
```
ccpm-demo/
├── .claude/                    # Claude Code configuration
│   ├── context/               # Project context documentation
│   └── hooks/                 # Git and development hooks
├── .git/                      # Git repository metadata
├── ccpm/                      # CCPM system files (untracked)
│   ├── agents/               # Specialized agent configurations
│   ├── commands/             # Custom command definitions
│   ├── context/              # CCPM context templates
│   ├── epics/                # Epic management files
│   ├── hooks/                # Project-specific hooks
│   ├── prds/                 # Product Requirements Documents
│   ├── rules/                # Development rules and guidelines
│   └── scripts/              # Automation scripts
├── doc/                       # Documentation directory
├── install/                   # Installation scripts
├── src/                       # Source code directory (empty)
├── zh-docs/                   # Chinese documentation
└── [Various .md files]        # Project documentation
```

## Key Directories

### `/ccpm/` - Project Management System
- **Purpose**: Complete project management infrastructure
- **Contents**: Scripts, configurations, templates for development workflow
- **Key Subdirectories**:
  - `scripts/pm/` - Project management automation (15+ scripts)
  - `rules/` - Development guidelines and standards
  - `agents/` - Specialized AI agent configurations
  - `hooks/` - Git and development lifecycle hooks

### `/.claude/` - Claude Code Configuration
- **Purpose**: Claude Code IDE integration and context management
- **Contents**: Hooks, context files, and IDE-specific configurations
- **Key Files**:
  - `context/` - Project context documentation (this directory)
  - `hooks/` - Development workflow automation

### `/src/` - Source Code
- **Purpose**: Main application source code
- **Status**: Currently empty, ready for development
- **Organization**: TBD based on technology stack choice

### Documentation Structure
- **Root Level**: High-level project documentation (README, LICENSE, etc.)
- **`/doc/`**: Detailed technical documentation
- **`/zh-docs/`**: Chinese language documentation
- **CCPM Documentation**: Integrated within ccpm/ structure

## File Naming Patterns
- **Configuration Files**: `.json`, `.config` extensions
- **Documentation**: `.md` markdown format
- **Scripts**: `.sh` shell scripts with executable permissions
- **Templates**: Located in respective subdirectories with clear naming

## Module Organization
- **Monolithic Structure**: Single repository for all components
- **Separation of Concerns**: Clear division between PM system, source code, and documentation
- **Extensible Design**: Easy to add new PM scripts and configurations

## Development Environment Structure
- **Git Integration**: Full git repository with hook support
- **PM Integration**: Comprehensive project management workflow
- **Documentation-First**: Emphasis on clear documentation and context
- **Automation-Ready**: Scripts and hooks for common development tasks