#!/bin/bash

# Example usage script for Inventory API Seeders
# This script demonstrates various ways to use the seeder system

echo "🌱 Inventory API Seeder Examples"
echo "================================"
echo

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo "❌ Error: Please run this script from the project root directory"
    exit 1
fi

echo "📋 Available seeder commands:"
echo
echo "1. Basic usage examples:"
echo "   go run cmd/seed/main.go                           # Default: all types, 20 records each"
echo "   go run cmd/seed/main.go -type=users -count=50     # 50 users only"
echo "   go run cmd/seed/main.go -type=categories -count=30 # 30 categories (hierarchical)"
echo "   go run cmd/seed/main.go -type=locations -count=40  # 40 locations"
echo

echo "2. Using the shell script (Linux/Mac):"
echo "   ./scripts/seed.sh                                 # Interactive mode"
echo "   ./scripts/seed.sh --quick-setup                   # 10 records each"
echo "   ./scripts/seed.sh --demo-data                     # 50 records each"
echo "   ./scripts/seed.sh -t users -c 100                # 100 users"
echo

echo "3. Using the batch script (Windows):"
echo "   scripts\\seed.bat                                 # Interactive mode"
echo "   scripts\\seed.bat --quick-setup                   # 10 records each"
echo "   scripts\\seed.bat --demo-data                     # 50 records each"
echo "   scripts\\seed.bat -t users -c 100                # 100 users"
echo

echo "4. Using Makefile:"
echo "   make seed-quick                                   # Quick development setup"
echo "   make seed-demo                                    # Demo dataset"
echo "   make seed-users COUNT=100                         # 100 users"
echo "   make seed-all COUNT=200                           # 200 of each type"
echo

echo "🎯 Interactive Demo"
echo "=================="
echo "Would you like to run a quick demo? This will:"
echo "- Create 1 admin user (admin@inventory.com / admin123456)"
echo "- Create 9 additional random users"
echo "- Create 10 categories with parent-child hierarchy"
echo "- Create 10 realistic locations"
echo
read -p "Run demo? (y/N): " response

if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
    echo
    echo "🚀 Running quick demo setup..."

    # Check if .env exists
    if [ ! -f ".env" ]; then
        echo "⚠️  Warning: .env file not found"
        echo "Please ensure your database connection is configured"
        echo
        read -p "Continue anyway? (y/N): " continue_response
        if [[ ! "$continue_response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
            echo "Demo cancelled."
            exit 0
        fi
    fi

    # Run the seeder
    echo "📊 Seeding demo data..."
    if go run cmd/seed/main.go -type=all -count=10; then
        echo
        echo "✅ Demo completed successfully!"
        echo
        echo "📈 What was created:"
        echo "   - 10 users (1 admin + 9 random users)"
        echo "   - 10 categories (~2-3 parents with children)"
        echo "   - 10 locations (mix of predefined and random)"
        echo
        echo "🔑 Admin credentials:"
        echo "   Email: admin@inventory.com"
        echo "   Password: admin123456"
        echo
        echo "💡 Next steps:"
        echo "   1. Start your API server: go run app/main.go"
        echo "   2. Test login with the admin credentials"
        echo "   3. Explore the seeded data through your API endpoints"
        echo
    else
        echo "❌ Demo failed. Please check your database configuration."
    fi
else
    echo "Demo cancelled. You can run the seeders manually using the commands above."
fi

echo
echo "📚 For more information, see:"
echo "   - seeders/README.md"
echo "   - go run cmd/seed/main.go -help"
echo "   - make help (if using Makefile)"
