---
created: 2025-09-27T18:57:36Z
last_updated: 2025-09-27T18:57:36Z
version: 1.0
author: Claude Code PM System
---

# System Patterns

## Architectural Patterns

### Command Pattern
- **Implementation**: CCPM uses shell scripts as commands
- **Location**: `ccpm/scripts/pm/*.sh`
- **Benefits**: Modular, extensible command system
- **Usage**: Each PM function is a separate executable script

### Template Pattern
- **Implementation**: Context templates and configuration examples
- **Location**: `ccpm/context/`, `ccpm/settings.json.example`
- **Benefits**: Consistent structure across different project types
- **Usage**: Standardized templates for PRDs, epics, and documentation

### Hook Pattern
- **Implementation**: Git hooks and development lifecycle hooks
- **Location**: `ccpm/hooks/`, `.claude/hooks/`
- **Benefits**: Automated workflow integration
- **Usage**: Triggered actions during git operations and development events

### Configuration Strategy Pattern
- **Implementation**: Multiple configuration layers
- **Locations**:
  - `CLAUDE.md` - Project-specific instructions
  - `ccpm/ccpm.config` - CCPM system configuration
  - `ccpm/settings.local.json` - Local environment settings
- **Benefits**: Flexible, environment-specific configuration

## Design Decisions

### Separation of Concerns
- **PM System**: Isolated in `ccpm/` directory
- **Application Code**: Separated in `src/` directory
- **Configuration**: Centralized but layered
- **Documentation**: Distributed but structured

### Automation First
- **Principle**: Automate repetitive tasks through scripts
- **Implementation**: Comprehensive script library for PM tasks
- **Benefits**: Reduces manual overhead, ensures consistency
- **Examples**: Status reporting, PRD management, epic tracking

### Documentation as Code
- **Principle**: Documentation lives alongside code
- **Implementation**: Markdown files with structured frontmatter
- **Benefits**: Version-controlled, searchable, maintainable
- **Context System**: Automated context generation and management

### Git-Centric Workflow
- **Principle**: Git is the source of truth for project state
- **Implementation**: Git hooks, branch-based workflows
- **Benefits**: Standard developer experience, tool integration
- **Extensions**: GitHub CLI integration for issue management

## Data Flow Patterns

### Command Execution Flow
1. **User Input**: Slash commands in Claude Code
2. **Script Routing**: Commands mapped to shell scripts
3. **Execution**: Bash scripts perform PM operations
4. **Output**: Results returned to user interface
5. **State Update**: Project state updated in git/files

### Context Management Flow
1. **Analysis**: Automated project analysis
2. **Generation**: Context files created/updated
3. **Storage**: Structured storage in `.claude/context/`
4. **Retrieval**: Context loaded for AI assistance
5. **Updates**: Regular context refresh cycles

### Configuration Cascade
1. **Global Defaults**: CCPM system defaults
2. **Project Config**: `ccpm/ccpm.config`
3. **Local Settings**: `ccpm/settings.local.json`
4. **Claude Instructions**: `CLAUDE.md`
5. **Runtime Overrides**: Command-line parameters

## Integration Patterns

### GitHub Integration
- **API Access**: Through GitHub CLI
- **Authentication**: OAuth-based
- **Operations**: Issues, PRs, repository management
- **Extensions**: gh-sub-issue for enhanced functionality

### Claude Code Integration
- **Context System**: Automated context management
- **Command Interface**: Slash command integration
- **Hook System**: Development lifecycle automation
- **AI Assistance**: Context-aware development support

### Shell Integration
- **Script Execution**: Direct bash script invocation
- **Environment**: Standard Unix/Linux shell environment
- **Permissions**: Executable script permissions
- **Error Handling**: Standard shell error codes and output

## Quality Patterns

### Error Handling
- **Graceful Degradation**: Scripts handle missing dependencies
- **Clear Messages**: Descriptive error output
- **Exit Codes**: Standard Unix exit code conventions
- **Logging**: Structured output for debugging

### Testing Strategy
- **Validation Scripts**: Built-in validation (validate.sh)
- **Integration Testing**: Git workflow testing
- **Documentation Testing**: Context accuracy verification
- **System Testing**: End-to-end PM workflow testing

### Maintainability
- **Modular Design**: Clear component boundaries
- **Documentation**: Comprehensive inline and external docs
- **Version Control**: Full git integration
- **Extensibility**: Easy addition of new scripts and features