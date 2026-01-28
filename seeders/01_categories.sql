-- Seeder SQL for Categories
-- This file populates the categories table with sample data

-- Clear existing data (optional, comment out if you want to preserve existing data)
-- TRUNCATE TABLE categories RESTART IDENTITY CASCADE;

INSERT INTO categories (name, description) VALUES
('Electronics', 'Electronic devices and gadgets'),
('Food & Beverage', 'Food and drink products'),
('Clothing', 'Apparel and fashion items'),
('Books & Stationery', 'Books, notebooks, and office supplies'),
('Health & Beauty', 'Personal care and beauty products'),
('Sports & Outdoor', 'Sports equipment and outdoor gear'),
('Home & Garden', 'Home decor and gardening supplies'),
('Toys & Games', 'Toys, games, and entertainment');

-- Verify data
SELECT * FROM categories;
