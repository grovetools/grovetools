# grovetools

A Grove ecosystem.

## Development Setup

### Go Module Configuration

These repositories are private. Configure Go to skip the public proxy and checksum database:

```bash
go env -w GOPRIVATE="github.com/grovetools/*"
```

### Building

Build all packages:

```bash
make build
```

### Directory Structure

| Directory | Binary | Description |
|-----------|--------|-------------|
| agentlogs | aglogs | Agent log viewer |
| core | - | Shared library |
| cx | cx | Context/rules management |
| docgen | docgen | Documentation generator |
| flow | flow | LLM job orchestration |
| grove | grove | Meta CLI tool |
| grove-anthropic | grove-anthropic | Anthropic API client |
| grove-gemini | grove-gemini | Gemini API client |
| grove.nvim | grove-nvim | Neovim plugin |
| hooks | hooks | Git/editor hooks |
| nb | nb | Notebook/notes management |
| notify | notify | Notifications |
| nav | nav | Tmux session navigator |
| project-tmpl-go | - | Go project template |
| skills | skills | Skill definitions |
| tend | tend | E2E test framework |
