---
created: 2025-09-27T18:57:36Z
last_updated: 2025-09-27T18:57:36Z
version: 1.0
author: Claude Code PM System
---

# Project Style Guide

## Coding Standards

### General Principles
- **Simplicity**: Prefer simple, readable solutions over clever code
- **Consistency**: Follow established patterns throughout the project
- **Documentation**: Code should be self-documenting, with comments for complex logic
- **Minimal Changes**: Implement the most concise solution that changes as little code as possible

### Shell Scripting Standards
- **Shebang**: Always use `#!/bin/bash` for bash scripts
- **Permissions**: Scripts must be executable (`chmod +x`)
- **Error Handling**: Use `set -e` for error handling where appropriate
- **Variables**: Use `${VARIABLE}` syntax for variable expansion
- **Functions**: Define functions with clear, descriptive names

#### Shell Script Example
```bash
#!/bin/bash
set -e

function show_status() {
    local project_name="${1}"
    echo "Status for ${project_name}:"
    # Implementation here
}
```

## Naming Conventions

### Files and Directories
- **Scripts**: kebab-case with `.sh` extension (e.g., `prd-status.sh`)
- **Directories**: lowercase with hyphens (e.g., `ccpm-scripts`)
- **Documentation**: kebab-case with `.md` extension (e.g., `project-overview.md`)
- **Configuration**: lowercase with appropriate extension (e.g., `ccpm.config`)

### Variables and Functions
- **Shell Variables**: UPPER_CASE for environment variables, lower_case for local variables
- **Function Names**: snake_case for shell functions
- **Constants**: UPPER_CASE with underscores

### Commands and Endpoints
- **PM Commands**: Use format `/pm:action` (e.g., `/pm:prd-new`)
- **Context Commands**: Use format `/context:action` (e.g., `/context:create`)
- **Command Parameters**: kebab-case for multi-word parameters

## File Structure Patterns

### Script Organization
```
ccpm/scripts/
├── pm/                    # PM-specific scripts
│   ├── init.sh           # System initialization
│   ├── prd-new.sh        # PRD creation
│   └── status.sh         # Status reporting
├── validation/           # Validation scripts
└── utilities/           # Utility scripts
```

### Documentation Structure
```
.claude/context/
├── progress.md           # Current project status
├── project-structure.md  # Directory organization
├── tech-context.md      # Technology stack
└── [other-context].md   # Additional context files
```

### Configuration Hierarchy
1. **System Defaults**: Built into CCPM scripts
2. **Project Config**: `ccpm/ccpm.config`
3. **Local Settings**: `ccpm/settings.local.json`
4. **Claude Instructions**: `CLAUDE.md`
5. **Runtime Overrides**: Command-line parameters

## Documentation Standards

### Markdown Conventions
- **Headers**: Use ATX-style headers (`#`, `##`, `###`)
- **Lists**: Use `-` for unordered lists, `1.` for ordered lists
- **Code Blocks**: Use triple backticks with language specification
- **Links**: Prefer reference-style links for readability
- **Emphasis**: Use `**bold**` for emphasis, `*italic*` for minor emphasis

### Frontmatter Requirements
All context files must include:
```yaml
---
created: YYYY-MM-DDTHH:MM:SSZ
last_updated: YYYY-MM-DDTHH:MM:SSZ
version: X.Y
author: Claude Code PM System
---
```

### Section Organization
Standard sections for context files:
1. **Purpose/Overview**: What this document covers
2. **Current State**: Present situation
3. **Details**: Comprehensive information
4. **Future Plans**: Planned changes or enhancements

## Comment Style

### Shell Script Comments
```bash
#!/bin/bash
# Script: prd-status.sh
# Purpose: Display status of all PRDs
# Usage: ./prd-status.sh [options]

# Check if PRD directory exists
if [[ ! -d "ccpm/prds" ]]; then
    echo "No PRDs found"
    exit 1
fi

# Process each PRD file
for prd_file in ccpm/prds/*.md; do
    # Extract title from frontmatter
    title=$(grep "^title:" "${prd_file}" | cut -d':' -f2 | xargs)
    echo "PRD: ${title}"
done
```

### Documentation Comments
- **Inline Comments**: Use for clarification of complex logic
- **Section Comments**: Use to separate major sections
- **TODO Comments**: Format as `# TODO: Description of what needs to be done`
- **FIXME Comments**: Format as `# FIXME: Description of issue to fix`

## Configuration Standards

### JSON Configuration
```json
{
    "version": "1.0",
    "project": {
        "name": "project-name",
        "type": "demo"
    },
    "settings": {
        "auto_update": true,
        "validation_level": "strict"
    }
}
```

### YAML Frontmatter
- **Consistent Keys**: Use standard keys across all files
- **ISO Dates**: Use ISO 8601 format for all dates
- **Semantic Versioning**: Use semver for version numbers
- **Clear Values**: Use descriptive, clear values

## Error Handling Patterns

### Script Error Handling
```bash
# Comprehensive error handling
set -euo pipefail

function handle_error() {
    echo "Error on line $1"
    exit 1
}

trap 'handle_error $LINENO' ERR
```

### Validation Patterns
- **Input Validation**: Always validate user inputs
- **File Existence**: Check file existence before operations
- **Permission Checks**: Verify permissions before file operations
- **Graceful Degradation**: Provide fallbacks for non-critical failures

## Testing Conventions

### Script Testing
- **Test Scripts**: Name test scripts with `.test.sh` suffix
- **Test Data**: Use `test-data/` directory for test fixtures
- **Assertions**: Use clear, descriptive test assertions
- **Coverage**: Aim for high test coverage of critical paths

### Validation Testing
- **Automated Validation**: Use `validate.sh` for project validation
- **Continuous Testing**: Test scripts as part of development workflow
- **Integration Testing**: Test full workflows end-to-end

## Best Practices

### Development Workflow
1. **Read First**: Always read existing code before making changes
2. **Small Changes**: Make minimal changes to achieve goals
3. **Test Early**: Test changes as early as possible
4. **Document Changes**: Update documentation with code changes
5. **Follow Patterns**: Use existing patterns and conventions

### Maintenance
- **Regular Updates**: Keep documentation current
- **Cleanup**: Remove unused files and code
- **Refactoring**: Improve code quality continuously
- **Monitoring**: Monitor script performance and reliability