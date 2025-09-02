# Inventory API Seeders

This package provides comprehensive data seeding functionality for the Inventory API. The seeders can generate realistic test data for users, categories (with hierarchical structure), and locations.

## Features

- **🔧 Flexible Command System**: Run seeders individually or all together
- **📊 Configurable Record Count**: Set custom number of records (default: 20)
- **🏗️ Hierarchical Categories**: Smart parent-child category creation
- **🌍 Multilingual Support**: Seeds data with Indonesian and English translations
- **🎯 Realistic Data**: Uses predefined templates and faker library for realistic data

## Quick Start

### Prerequisites

Make sure you have:
- Go 1.21 or higher
- PostgreSQL database with migrations applied
- Environment variables configured (see main application)

### Installation

1. Install the required dependency:
```bash
go mod tidy
```

2. Make sure your `.env` file is configured with database connection:
```env
DSN=postgres://username:password@localhost:5432/inventory_db?sslmode=disable
```

## Usage

### Command Line Interface

Navigate to your project root and run:

```bash
# Seed all data with default count (20 records each)
go run cmd/seed/main.go

# Seed all data with custom count
go run cmd/seed/main.go -type=all -count=50

# Seed specific data types
go run cmd/seed/main.go -type=users -count=30
go run cmd/seed/main.go -type=categories -count=25
go run cmd/seed/main.go -type=locations -count=40
```

### Command Options

| Option   | Description                 | Default | Valid Values                              |
| -------- | --------------------------- | ------- | ----------------------------------------- |
| `-type`  | Type of data to seed        | `all`   | `users`, `categories`, `locations`, `all` |
| `-count` | Number of records to create | `20`    | Any positive integer                      |
| `-help`  | Show help message           | `false` | -                                         |

### Examples

```bash
# Seed 100 users only
go run cmd/seed/main.go -type=users -count=100

# Seed 30 categories (will create ~7 parents with ~3 children each)
go run cmd/seed/main.go -type=categories -count=30

# Seed 50 locations
go run cmd/seed/main.go -type=locations -count=50

# Seed everything with 200 records each
go run cmd/seed/main.go -type=all -count=200

# Show help
go run cmd/seed/main.go -help
```

## Data Details

### Users
- **Always creates 1 admin user**: `admin@inventory.com` / `admin123456`
- **Role distribution**: 10% Admin, 30% Staff, 60% Employee
- **Includes**: Realistic names, emails, employee IDs, avatar URLs
- **Languages**: Random preference between English and Indonesian
- **Status**: 90% active users

### Categories
- **Hierarchical structure**: Creates parent categories first, then children
- **Parent count**: ~25% of total count (minimum 3)
- **Predefined parents**: Electronics, Furniture, Vehicles, Office Supplies, Tools, Safety
- **Children**: Realistic subcategories under each parent
- **Multilingual**: Both English and Indonesian names/descriptions

### Locations
- **Predefined locations**: Realistic office locations (lobby, meeting rooms, warehouses)
- **Building structure**: Multiple buildings with floor information
- **Geographic data**: Coordinates around Jakarta area
- **Multilingual**: English and Indonesian location names

## Architecture

### File Structure
```
cmd/seed/
├── main.go              # CLI entry point
seeders/
├── seeder_manager.go    # Main seeder coordinator
├── user_seeder.go       # User data seeding
├── category_seeder.go   # Category data seeding (with hierarchy)
└── location_seeder.go   # Location data seeding
```

### Key Components

1. **SeederManager**: Orchestrates all seeders and manages execution order
2. **Individual Seeders**: Specialized seeders for each entity type
3. **CLI Interface**: Command-line interface for easy execution

## Advanced Configuration

### Category Hierarchy Logic

For categories with total count N:
- **Parent count**: N ÷ 4 (minimum 3)
- **Children count**: N - parent_count
- **Distribution**: Children are distributed evenly among parents

Example with count=24:
- 6 parent categories
- 18 child categories (3 per parent)

### Predefined Data Templates

The seeders include realistic predefined templates:

**Category Templates**:
- Electronics → Computers, Phones, Printers
- Furniture → Chairs, Desks, Storage
- Vehicles → Cars, Motorcycles

**Location Templates**:
- HQ_LOBBY, HQ_IT_ROOM, HQ_MEETING_A
- WAREHOUSE_A, WAREHOUSE_B
- OFFICE_FL1, OFFICE_FL2

## Error Handling

- **Graceful failure**: If individual records fail, seeding continues
- **Progress reporting**: Shows successful vs failed creations
- **Database validation**: Respects unique constraints and foreign keys
- **Rollback safe**: Each seeder is independent

## Performance Notes

- **Batch processing**: Creates records in sequence for data integrity
- **Progress indicators**: Shows progress every 10 records for users
- **Memory efficient**: Doesn't load all data into memory at once

## Troubleshooting

### Common Issues

1. **Database connection failed**
   ```
   Solution: Check your DSN environment variable
   ```

2. **Unique constraint violations**
   ```
   Solution: Seeders handle duplicates gracefully, check output for warnings
   ```

3. **Foreign key violations**
   ```
   Solution: Ensure database migrations are applied correctly
   ```

### Debug Mode

Add debug output by modifying the seeder functions or check the database directly:

```sql
-- Check seeded data
SELECT COUNT(*) FROM users;
SELECT COUNT(*) FROM categories WHERE parent_id IS NULL; -- Parent categories
SELECT COUNT(*) FROM categories WHERE parent_id IS NOT NULL; -- Child categories
SELECT COUNT(*) FROM locations;
```

## Contributing

When adding new seeders:

1. Create a new seeder file in `seeders/` package
2. Implement the seeder interface with realistic data
3. Add the seeder to `SeederManager`
4. Update the CLI interface to support the new type
5. Add documentation and examples

## License

This seeder system is part of the Inventory API project and follows the same license terms.
