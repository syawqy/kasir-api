-- Seeder SQL for Products
-- This file populates the products table with sample data
-- Run 01_categories.sql BEFORE running this file

-- Clear existing data (optional, comment out if you want to preserve existing data)
-- TRUNCATE TABLE products RESTART IDENTITY CASCADE;

-- Electronics (category_id = 1)
INSERT INTO products (name, price, stock, category_id) VALUES
('Laptop HP ProBook', 12000000, 15, 1),
('Mouse Logitech Wireless', 250000, 50, 1),
('Keyboard Mechanical RGB', 750000, 30, 1),
('Monitor 24 inch Full HD', 2500000, 20, 1),
('Webcam HD 1080p', 450000, 25, 1);

-- Food & Beverage (category_id = 2)
INSERT INTO products (name, price, stock, category_id) VALUES
('Mineral Water 600ml', 3000, 200, 2),
('Instant Noodles', 5000, 150, 2),
('Coffee Arabica Premium', 45000, 80, 2),
('Green Tea Box 25 sachets', 35000, 60, 2),
('Chocolate Bar', 15000, 100, 2);

-- Clothing (category_id = 3)
INSERT INTO products (name, price, stock, category_id) VALUES
('T-Shirt Cotton Basic', 75000, 100, 3),
('Jeans Denim Regular Fit', 250000, 60, 3),
('Sneakers Sport', 450000, 40, 3),
('Hoodie Fleece', 180000, 50, 3),
('Baseball Cap', 65000, 80, 3);

-- Books & Stationery (category_id = 4)
INSERT INTO products (name, price, stock, category_id) VALUES
('Notebook A5 80 Pages', 12000, 150, 4),
('Ballpoint Pen Pack of 10', 25000, 200, 4),
('Novel Fiction Bestseller', 85000, 45, 4),
('Marker Highlighter Set', 35000, 90, 4),
('Sticky Notes Colorful', 18000, 120, 4);

-- Health & Beauty (category_id = 5)
INSERT INTO products (name, price, stock, category_id) VALUES
('Face Wash Cleanser', 45000, 80, 5),
('Shampoo Anti-Dandruff', 32000, 100, 5),
('Moisturizer SPF 30', 95000, 60, 5),
('Hand Sanitizer 100ml', 15000, 150, 5),
('Lip Balm Natural', 22000, 90, 5);

-- Sports & Outdoor (category_id = 6)
INSERT INTO products (name, price, stock, category_id) VALUES
('Yoga Mat TPE', 120000, 40, 6),
('Dumbbells Set 5kg', 180000, 30, 6),
('Water Bottle 1L', 45000, 80, 6),
('Resistance Bands Set', 85000, 50, 6),
('Jump Rope Speed', 35000, 70, 6);

-- Home & Garden (category_id = 7)
INSERT INTO products (name, price, stock, category_id) VALUES
('Ceramic Plant Pot', 55000, 60, 7),
('LED Light Bulb 12W', 28000, 100, 7),
('Wall Clock Modern', 125000, 35, 7),
('Scented Candle Lavender', 45000, 70, 7),
('Doormat Anti-Slip', 68000, 50, 7);

-- Toys & Games (category_id = 8)
INSERT INTO products (name, price, stock, category_id) VALUES
('Puzzle 1000 Pieces', 95000, 40, 8),
('Board Game Strategy', 250000, 25, 8),
('Action Figure Set', 180000, 35, 8),
('Remote Control Car', 350000, 20, 8),
('Building Blocks 500pcs', 125000, 45, 8);

-- Verify data
SELECT p.id, p.name, p.price, p.stock, c.name as category 
FROM products p 
LEFT JOIN categories c ON p.category_id = c.id 
ORDER BY p.category_id, p.id;
