# Database Seeders

This directory contains SQL files to populate the database with sample data for development and testing purposes.

## Prerequisites

- PostgreSQL database must be running
- Database and tables must be created (run migrations first)
- You need to have proper database credentials

## Running Seeders

### Option 1: Using psql Command Line

```bash
# Navigate to the kasir-api directory
cd d:\codes\go\kasir-api

# Connect to your database and run seeders in order
psql -h localhost -U postgres -d kasir -f seeders/01_categories.sql
psql -h localhost -U postgres -d kasir -f seeders/02_products.sql
```

Replace the connection parameters:
- `-h localhost`: Your database host
- `-U postgres`: Your database username
- `-d kasir`: Your database name

### Option 2: Using PostgreSQL Client (pgAdmin, DBeaver, etc.)

1. Open your PostgreSQL client
2. Connect to your `kasir` database
3. Open and execute `01_categories.sql` first
4. Then open and execute `02_products.sql`

### Option 3: Copy-Paste into psql Interactive Mode

```bash
# Connect to the database
psql -h localhost -U postgres -d kasir

# Then copy and paste the contents of each SQL file in order
```

## Execution Order

**IMPORTANT**: Always run the seeders in this order:

1. `01_categories.sql` - Creates 8 sample categories
2. `02_products.sql` - Creates 40 sample products linked to categories

The order is important because products have foreign key references to categories.

## Sample Data

### Categories (8 total)
- Electronics
- Food & Beverage
- Clothing
- Books & Stationery
- Health & Beauty
- Sports & Outdoor
- Home & Garden
- Toys & Games

### Products (40 total)
- 5 products per category
- Realistic Indonesian prices (in Rupiah)
- Various stock quantities

## Clearing Data

If you want to clear existing data before re-running seeders, uncomment the `TRUNCATE` statements at the top of each SQL file:

```sql
TRUNCATE TABLE categories RESTART IDENTITY CASCADE;
```

⚠️ **Warning**: This will delete ALL data in the tables and reset the ID sequences!

## Verification

Each seeder file includes a `SELECT` statement at the end to verify the data was inserted correctly. Check the output to confirm successful seeding.
