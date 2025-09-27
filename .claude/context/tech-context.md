---
created: 2025-09-27T18:57:36Z
last_updated: 2025-09-27T18:57:36Z
version: 1.0
author: Claude Code PM System
---

# Technology Context

## Core Technologies

### Project Management Stack
- **CCPM**: Claude Code Project Management system
- **GitHub CLI**: v2.79.0 (2025-09-08) - Repository and issue management
- **Git**: Version control with hook integration
- **gh-sub-issue**: GitHub CLI extension for enhanced issue management

### Development Environment
- **Platform**: macOS (Darwin 24.6.0)
- **Shell**: Bash scripting environment
- **Date**: 2025-09-27 (Current session)

### Language Detection
- **Primary Language**: Not yet determined (src/ directory empty)
- **Scripting**: Bash shell scripts for automation
- **Documentation**: Markdown format
- **Configuration**: JSON format

## Dependencies

### System Dependencies
- **Required**: GitHub CLI (gh) - ✅ Installed and authenticated
- **Required**: Git - ✅ Available and configured
- **Optional**: Various language-specific tools (TBD based on project choice)

### CCPM Dependencies
- **gh-sub-issue**: ✅ Installed
- **Git repository**: ✅ Initialized
- **Shell environment**: ✅ Bash-compatible

### Project Dependencies
- **No package.json**: Node.js not currently used
- **No requirements.txt**: Python not currently used
- **No Cargo.toml**: Rust not currently used
- **No go.mod**: Go not currently used
- **Technology Stack**: To be determined based on project requirements

## Development Tools

### Installed Tools
- **GitHub CLI**: Full authentication and API access
- **Git**: Repository management and version control
- **CCPM Scripts**: 15+ automation scripts for project management
- **Claude Code**: AI-powered development environment

### Available Scripts
Located in `ccpm/scripts/pm/`:
- `init.sh` - System initialization
- `help.sh` - Command documentation
- `status.sh` - Project status reporting
- `standup.sh` - Daily standup automation
- `prd-*.sh` - PRD management
- `epic-*.sh` - Epic management
- `validate.sh` - Validation and checks
- And more...

### Hook System
- **Git Hooks**: Available in `ccpm/hooks/` and `.claude/hooks/`
- **Automation**: bash-worktree-fix.sh for git workflow enhancement
- **Integration**: Claude Code lifecycle hooks

## Configuration Management

### CCPM Configuration
- **Main Config**: `ccpm/ccpm.config`
- **Settings**: `ccpm/settings.local.json`
- **Example Settings**: `ccpm/settings.json.example`

### Claude Code Configuration
- **CLAUDE.md**: Project-specific instructions and testing commands
- **Context System**: Automated context management in `.claude/context/`

## Technology Decisions

### Architectural Choices
- **Documentation-First**: Markdown-based documentation system
- **Automation-Heavy**: Extensive shell scripting for PM tasks
- **Git-Centric**: Git-based workflow with hook integration
- **Modular Design**: Clear separation between PM system and application code

### Standards and Conventions
- **Shell Scripts**: Bash-compatible, executable permissions
- **Documentation**: Markdown format with frontmatter
- **Configuration**: JSON format for structured data
- **Naming**: Kebab-case for scripts, camelCase for JSON

## Version Information
- **CCPM**: Freshly installed (latest)
- **GitHub CLI**: 2.79.0 (2025-09-08)
- **Platform**: macOS Darwin 24.6.0
- **Project**: Initial setup phase, version TBD