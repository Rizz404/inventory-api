# Inventory API Makefile

.PHONY: help seed seed-users seed-categories seed-locations seed-all seed-quick seed-demo seed-load-test clean-db build-seeder swagger-gen build run dev

# Default target
help:
	@echo "Inventory API Commands"
	@echo "===================="
	@echo ""
	@echo "🚀 Development Commands:"
	@echo "  make dev               - Generate swagger + build + run (development mode)"
	@echo "  make build             - Build the application"
	@echo "  make run               - Build and run the application"
	@echo "  make swagger-gen       - Generate Swagger documentation"
	@echo "  make clean             - Clean build artifacts and docs"
	@echo ""
	@echo "📚 Documentation:"
	@echo "  After running 'make dev' or 'make run', access:"
	@echo "  • Swagger UI: http://localhost:8080/docs/"
	@echo "  • API Base: http://localhost:8080/api/v1/"
	@echo ""
	@echo "🌱 Seeding Commands:"
	@echo "  make seed              - Run seeder with interactive prompts"
	@echo "  make seed-users        - Seed 20 users"
	@echo "  make seed-categories   - Seed 20 categories (with hierarchy)"
	@echo "  make seed-locations    - Seed 20 locations"
	@echo "  make seed-all          - Seed all data types (20 each)"
	@echo ""
	@echo "⚡ Quick Presets:"
	@echo "  make seed-quick        - Quick setup (10 records each)"
	@echo "  make seed-demo         - Demo data (50 records each)"
	@echo "  make seed-load-test    - Load test data (500 records each)"
	@echo ""
	@echo "🔧 Utility Commands:"
	@echo "  make build-seeder      - Build seeder binary"
	@echo "  make clean-db          - Clean database (requires confirmation)"
	@echo "  make install-swagger   - Install Swagger CLI"
	@echo ""
	@echo "📝 Custom Usage:"
	@echo "  make seed TYPE=users COUNT=100"
	@echo "  make seed TYPE=categories COUNT=50"
	@echo ""

# Interactive seeder
seed:
	@echo "Inventory API Seeder"
	@echo "==================="
	@echo "Choose data type to seed:"
	@echo "1) All (users, categories, locations)"
	@echo "2) Users only"
	@echo "3) Categories only"
	@echo "4) Locations only"
	@read -p "Enter choice (1-4): " choice; \
	case $$choice in \
		1) $(MAKE) seed-all ;; \
		2) $(MAKE) seed-users ;; \
		3) $(MAKE) seed-categories ;; \
		4) $(MAKE) seed-locations ;; \
		*) echo "Invalid choice" ;; \
	esac

# Individual seeders with default count
seed-users:
	go run cmd/seed/main.go -type=users -count=$(or $(COUNT),20)

seed-categories:
	go run cmd/seed/main.go -type=categories -count=$(or $(COUNT),20)

seed-locations:
	go run cmd/seed/main.go -type=locations -count=$(or $(COUNT),20)

seed-all:
	go run cmd/seed/main.go -type=all -count=$(or $(COUNT),20)

# Quick presets
seed-quick:
	@echo "🚀 Quick Setup: Creating minimal dataset for development..."
	go run cmd/seed/main.go -type=all -count=10

seed-demo:
	@echo "🎯 Demo Data: Creating demo dataset..."
	go run cmd/seed/main.go -type=all -count=50

seed-load-test:
	@echo "🏋️ Load Test: Creating large dataset..."
	@echo "⚠️  This will create 500 records of each type."
	@read -p "Continue? (y/N): " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		go run cmd/seed/main.go -type=all -count=500; \
	else \
		echo "Cancelled."; \
	fi

# Build seeder binary
build-seeder:
	@echo "Building seeder binary..."
	go build -o bin/seeder cmd/seed/main.go
	@echo "✅ Seeder built at bin/seeder"

# Clean database (dangerous operation)
clean-db:
	@echo "⚠️  WARNING: This will delete ALL data from the database!"
	@echo "This operation cannot be undone."
	@read -p "Type 'DELETE_ALL_DATA' to confirm: " confirm; \
	if [ "$$confirm" = "DELETE_ALL_DATA" ]; then \
		echo "Cleaning database..."; \
		psql $(DSN) -c "TRUNCATE TABLE notifications, scan_logs, issue_reports, maintenance_records, maintenance_schedules, asset_movements, assets, categories_translation, categories, locations_translation, locations, users RESTART IDENTITY CASCADE;"; \
		echo "✅ Database cleaned"; \
	else \
		echo "Cancelled - confirmation text did not match"; \
	fi

# Install dependencies
deps:
	go mod tidy
	go mod download

# Run tests on seeders
test-seeders:
	go test ./seeders/... -v

# Lint seeders
lint-seeders:
	golangci-lint run ./seeders/...

# Format seeder code
fmt-seeders:
	go fmt ./seeders/...
	go fmt ./cmd/seed/...

# Show seeder help
seed-help:
	go run cmd/seed/main.go -help

# Development shortcuts with custom counts
seed-users-100:
	$(MAKE) seed-users COUNT=100

seed-categories-50:
	$(MAKE) seed-categories COUNT=50

seed-locations-75:
	$(MAKE) seed-locations COUNT=75

# Combined operations
setup-dev: seed-quick
	@echo "✅ Development environment set up with test data"

setup-demo: seed-demo
	@echo "✅ Demo environment set up with sample data"

# Check prerequisites
check-env:
	@if [ -z "$(DSN)" ]; then \
		echo "❌ DSN environment variable not set"; \
		echo "Please set DSN in your .env file or environment"; \
		exit 1; \
	fi
	@if [ ! -f ".env" ]; then \
		echo "⚠️  Warning: .env file not found"; \
	fi
	@echo "✅ Environment check passed"

# Run seeder with environment check
seed-safe: check-env
	$(MAKE) seed

# Show current database stats
db-stats:
	@echo "Current Database Statistics:"
	@echo "=========================="
	@psql $(DSN) -c "SELECT 'Users' as table_name, COUNT(*) as count FROM users UNION ALL SELECT 'Categories', COUNT(*) FROM categories UNION ALL SELECT 'Locations', COUNT(*) FROM locations ORDER BY table_name;"

# ==============================================================================
# SWAGGER DOCUMENTATION
# ==============================================================================

# Generate Swagger documentation
swagger-gen:
	@echo "🔄 Generating Swagger documentation..."
	swag init -g app/main.go -o docs --parseDependency --parseInternal
	@echo "✅ Swagger documentation generated successfully"
	@echo "📖 Documentation available at: http://localhost:8080/docs/"

# ==============================================================================
# BUILD & RUN COMMANDS
# ==============================================================================

# Build the application
build:
	@echo "🔨 Building application..."
	go build -o bin/app.exe app/main.go
	@echo "✅ Build completed: bin/app.exe"

# Run the application
run: build
	@echo "🚀 Starting application..."
	./bin/app.exe

# Development mode (with swagger generation)
dev: swagger-gen build
	@echo "🚀 Starting application in development mode..."
	./bin/app.exe

# Install swagger CLI if not available
install-swagger:
	@echo "📦 Installing Swagger CLI..."
	go install github.com/swaggo/swag/cmd/swag@latest
	@echo "✅ Swagger CLI installed"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@if exist bin\\app.exe del bin\\app.exe
	@if exist docs\\docs.go del docs\\docs.go
	@if exist docs\\swagger.json del docs\\swagger.json
	@if exist docs\\swagger.yaml del docs\\swagger.yaml
	@echo "✅ Clean completed"

.DEFAULT_GOAL := help
