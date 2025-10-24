#!/bin/bash

# Inventory API Seeder Script
# This script provides easy access to run seeders with common configurations

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
DEFAULT_COUNT=20
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Function to print colored output
print_colored() {
    echo -e "${1}${2}${NC}"
}

# Function to show help
show_help() {
    echo "Inventory API Seeder Script"
    echo "=========================="
    echo
    echo "Usage: $0 [OPTIONS]"
    echo
    echo "Options:"
    echo "  -t, --type TYPE     Type of seed: users, categories, locations, all (default: all)"
    echo "  -c, --count COUNT   Number of records to create (default: 20)"
    echo "  -h, --help          Show this help message"
    echo
    echo "Quick Commands:"
    echo "  $0 --quick-setup    Seed a small dataset for development (10 records each)"
    echo "  $0 --demo-data      Seed demo dataset (50 records each)"
    echo "  $0 --load-test      Seed large dataset for load testing (500 records each)"
    echo
    echo "Examples:"
    echo "  $0                          # Seed all with default count (20)"
    echo "  $0 -t users -c 50          # Seed 50 users"
    echo "  $0 --type=categories --count=30  # Seed 30 categories"
    echo "  $0 --quick-setup            # Quick development setup"
    echo
}

# Function to run seeder
run_seeder() {
    local type=$1
    local count=$2

    print_colored $BLUE "🌱 Running seeder: type=$type, count=$count"

    cd "$PROJECT_ROOT"

    if ! go run cmd/seed/main.go -type="$type" -count="$count"; then
        print_colored $RED "❌ Seeding failed!"
        exit 1
    fi

    print_colored $GREEN "✅ Seeding completed successfully!"
}

# Function for quick setup
quick_setup() {
    print_colored $YELLOW "🚀 Quick Setup: Creating minimal dataset for development..."
    run_seeder "all" 10
}

# Function for demo data
demo_data() {
    print_colored $YELLOW "🎯 Demo Data: Creating demo dataset..."
    run_seeder "all" 50
}

# Function for load test data
load_test() {
    print_colored $YELLOW "🏋️  Load Test: Creating large dataset..."
    print_colored $YELLOW "⚠️  This will create 500 records of each type. Continue? (y/N)"
    read -r response
    if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
        run_seeder "all" 500
    else
        print_colored $YELLOW "Cancelled."
        exit 0
    fi
}

# Parse command line arguments
TYPE="all"
COUNT=$DEFAULT_COUNT

while [[ $# -gt 0 ]]; do
    case $1 in
        -t|--type)
            TYPE="$2"
            shift 2
            ;;
        -c|--count)
            COUNT="$2"
            shift 2
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        --quick-setup)
            quick_setup
            exit 0
            ;;
        --demo-data)
            demo_data
            exit 0
            ;;
        --load-test)
            load_test
            exit 0
            ;;
        --type=*)
            TYPE="${1#*=}"
            shift
            ;;
        --count=*)
            COUNT="${1#*=}"
            shift
            ;;
        *)
            print_colored $RED "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Validate inputs
if [[ ! "$TYPE" =~ ^(users|categories|locations|all)$ ]]; then
    print_colored $RED "Invalid type: $TYPE"
    print_colored $YELLOW "Valid types: users, categories, locations, all"
    exit 1
fi

if ! [[ "$COUNT" =~ ^[0-9]+$ ]] || [ "$COUNT" -le 0 ]; then
    print_colored $RED "Invalid count: $COUNT"
    print_colored $YELLOW "Count must be a positive integer"
    exit 1
fi

# Check if we're in the right directory
if [ ! -f "$PROJECT_ROOT/go.mod" ]; then
    print_colored $RED "Error: go.mod not found. Are you in the correct directory?"
    exit 1
fi

if [ ! -f "$PROJECT_ROOT/cmd/seed/main.go" ]; then
    print_colored $RED "Error: Seeder not found at cmd/seed/main.go"
    exit 1
fi

# Check if .env exists
if [ ! -f "$PROJECT_ROOT/.env" ]; then
    print_colored $YELLOW "⚠️  Warning: .env file not found. Make sure environment variables are set."
fi

# Run the seeder
run_seeder "$TYPE" "$COUNT"
