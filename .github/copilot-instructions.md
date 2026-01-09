// Copilot Guidelines - Inventory API

// 1) Go Conventions & Best Practices
// - Follow standard Go project layout (cmd, internal, pkg, etc.)
// - Use meaningful package names (lowercase, no underscores)
// - Naming: camelCase locals, PascalCase exports, avoid getters prefix

// 2) Project Structure Rules
// - domain: pure business logic, no dependencies on frameworks
// - internal: implementation details, not exposed
// - services: business logic orchestration, call repositories
// - repositories: data access layer only, return domain without response prefix types
// - handlers: thin layer, validation + service calls + response
// - mappers: convert between domain & DB models, one-way transformations

// 3) Database & GORM
// - Always use transactions for multi-table operations
// - Preload relationships explicitly, avoid N+1 queries
// - Use db.Model(&Type{}) instead of raw table names
// - Soft delete: gorm.DeletedAt field, use Unscoped() when needed
// - Migrations: never modify, create new ones for changes

// 4) Testing (when explicitly requested)
// - Use testify/assert and testify/mock
// - Test file naming: \*\_test.go in same package
// - Table-driven tests for multiple scenarios
// - Mock external dependencies (DB, APIs, clients)

// 5) Performance
// - Use pointers for large structs (>64 bytes)
// - Defer only when necessary (defers have overhead)
// - Batch operations when processing multiple records
// - Use goroutines for independent tasks, add context cancellation

// 6) Comments (Better Comments format)
// TODO: | FIXME: | ! warning | ? question | \* important note

// 7) Response: brief & to the point
// Only mention what changed/added/removed, no lengthy explanations

// 8) Docs: NO .md files unless explicitly requested
// Keep inline comments 1-2 lines max, code should be self-explanatory

// 9) Terminal: use modern CLI tools
// Files: eza, fd, rg, bat, sd | Git: lazygit, gh, delta
// Nav: z (zoxide), fzf, yazi | Dev: glow, jq, tldr, micro
// Monitor: btm, procs, dust, duf
// ❌ Avoid: dir, findstr, find, grep, cat, manual cd

// 10) Deployment & Docker
// - Backend deployed on AWS using Docker & Docker Compose
// - ❌ NEVER run docker/docker-compose commands in terminal
// - Only provide docker commands as instructions/notes
// - User will manually execute Docker operations

// Terminal workflows:
// fd -e go | rg "TODO" // find TODOs
// eza --tree -L 3 internal/ // show structure
// z inventory // jump to project
// lazygit // interactive git
// go mod tidy && go mod verify // clean deps
// air // hot reload dev server
