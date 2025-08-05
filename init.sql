-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    address TEXT,
    city VARCHAR(100),
    country VARCHAR(100) DEFAULT 'Saudi Arabia',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create categories table
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    image_url VARCHAR(500),
    parent_id INTEGER REFERENCES categories(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Update products table with more e-commerce fields
DROP TABLE IF EXISTS products CASCADE;
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    original_price DECIMAL(10,2),
    category_id INTEGER REFERENCES categories(id),
    description TEXT,
    short_description VARCHAR(500),
    image_url VARCHAR(500),
    stock_quantity INTEGER DEFAULT 0,
    sku VARCHAR(100) UNIQUE,
    weight DECIMAL(8,2),
    dimensions VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    is_featured BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create product images table
CREATE TABLE IF NOT EXISTS product_images (
    id SERIAL PRIMARY KEY,
    product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
    image_url VARCHAR(500) NOT NULL,
    alt_text VARCHAR(255),
    is_primary BOOLEAN DEFAULT false,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create cart table
CREATE TABLE IF NOT EXISTS cart (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, product_id)
);

-- Create wishlist table
CREATE TABLE IF NOT EXISTS wishlist (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, product_id)
);

-- Create orders table
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    order_number VARCHAR(50) UNIQUE NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    total_amount DECIMAL(10,2) NOT NULL,
    shipping_address TEXT NOT NULL,
    shipping_city VARCHAR(100) NOT NULL,
    shipping_country VARCHAR(100) DEFAULT 'Saudi Arabia',
    shipping_phone VARCHAR(20) NOT NULL,
    payment_method VARCHAR(50) DEFAULT 'cash_on_delivery',
    payment_status VARCHAR(50) DEFAULT 'pending',
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create order items table
CREATE TABLE IF NOT EXISTS order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES orders(id) ON DELETE CASCADE,
    product_id INTEGER REFERENCES products(id),
    product_name VARCHAR(255) NOT NULL,
    product_price DECIMAL(10,2) NOT NULL,
    quantity INTEGER NOT NULL,
    total_price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create reviews table
CREATE TABLE IF NOT EXISTS reviews (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    title VARCHAR(255),
    comment TEXT,
    is_verified_purchase BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create review likes table
CREATE TABLE IF NOT EXISTS review_likes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    review_id INTEGER REFERENCES reviews(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, review_id)
);

-- Create product likes table
CREATE TABLE IF NOT EXISTS product_likes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, product_id)
);

-- Insert sample categories
INSERT INTO categories (name, description, image_url) VALUES
    ('Electronics', 'Latest electronic devices and gadgets', 'https://images.unsplash.com/photo-1498049794561-7780e7231661?w=400'),
    ('Fashion', 'Trendy clothing and accessories', 'https://images.unsplash.com/photo-1445205170230-053b83016050?w=400'),
    ('Home & Garden', 'Everything for your home and garden', 'https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=400'),
    ('Sports', 'Sports equipment and activewear', 'https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?w=400'),
    ('Books', 'Books for all ages and interests', 'https://images.unsplash.com/photo-1544947950-fa07a98d237f?w=400'),
    ('Beauty', 'Beauty and personal care products', 'https://images.unsplash.com/photo-1596462502278-27bfdc403348?w=400'),
    ('Toys', 'Toys and games for children', 'https://images.unsplash.com/photo-1566576912321-d58ddd7a6088?w=400'),
    ('Automotive', 'Car accessories and maintenance', 'https://images.unsplash.com/photo-1549317661-bd32c8ce0db2?w=400')
ON CONFLICT (id) DO NOTHING;

-- Insert sample products with categories
INSERT INTO products (name, price, original_price, category_id, description, short_description, image_url, stock_quantity, sku, weight, is_featured) VALUES
    ('iPhone 15 Pro', 999.99, 1099.99, 1, 'Latest iPhone with advanced camera system and A17 Pro chip. Features titanium design, 48MP camera, and all-day battery life.', 'Premium smartphone with cutting-edge technology', 'https://images.unsplash.com/photo-1592750475338-74b7b21085ab?w=400', 50, 'IPH15PRO-001', 0.187, true),
    ('MacBook Air M2', 1199.99, 1299.99, 1, 'Ultra-thin laptop with powerful M2 chip and all-day battery life. Perfect for work and creativity.', 'Lightweight laptop with exceptional performance', 'https://images.unsplash.com/photo-1517336714731-489689fd1ca8?w=400', 30, 'MBA-M2-001', 1.24, true),
    ('Sony WH-1000XM5', 349.99, 399.99, 1, 'Premium noise-cancelling headphones with exceptional sound quality and 30-hour battery life.', 'Best-in-class noise cancellation', 'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=400', 25, 'SONY-WH5-001', 0.25, false),
    ('Nike Air Max 270', 129.99, 149.99, 2, 'Comfortable running shoes with Air Max technology for maximum cushioning and style.', 'Comfortable and stylish running shoes', 'https://images.unsplash.com/photo-1542291026-7eec264c27ff?w=400', 100, 'NIKE-AM270-001', 0.85, true),
    ('Adidas Ultraboost 22', 179.99, 199.99, 2, 'High-performance running shoes with responsive cushioning and energy return technology.', 'Professional running shoes', 'https://images.unsplash.com/photo-1608231387042-66d1773070a5?w=400', 75, 'ADIDAS-UB22-001', 0.9, false),
    ('Samsung 4K Smart TV', 799.99, 899.99, 1, '55-inch 4K Ultra HD Smart TV with HDR and built-in streaming apps.', 'Crystal clear 4K entertainment', 'https://images.unsplash.com/photo-1593359677879-a4bb92f829d1?w=400', 20, 'SAMSUNG-4K-001', 15.5, true),
    ('Coffee Maker', 89.99, 99.99, 3, 'Programmable coffee maker with 12-cup capacity and auto-shutoff feature.', 'Perfect morning coffee every time', 'https://images.unsplash.com/photo-1517668808822-9ebb02f2a0e6?w=400', 45, 'COFFEE-001', 2.1, false),
    ('Yoga Mat', 29.99, 39.99, 4, 'Non-slip yoga mat with carrying strap, perfect for home workouts and studio sessions.', 'Premium non-slip yoga mat', 'https://images.unsplash.com/photo-1544367567-0f2fcb009e0b?w=400', 200, 'YOGA-MAT-001', 0.8, false),
    ('Wireless Earbuds', 79.99, 99.99, 1, 'True wireless earbuds with noise cancellation and 24-hour battery life.', 'Crystal clear wireless audio', 'https://images.unsplash.com/photo-1590658268037-6bf12165a8df?w=400', 150, 'EARBUDS-001', 0.05, true),
    ('Smart Watch', 299.99, 349.99, 1, 'Fitness tracking smartwatch with heart rate monitor and GPS.', 'Track your fitness goals', 'https://images.unsplash.com/photo-1523275335684-37898b6baf30?w=400', 60, 'SMARTWATCH-001', 0.12, true)
ON CONFLICT (id) DO NOTHING;

-- Insert sample users (password: password123)
INSERT INTO users (email, password_hash, first_name, last_name, phone, address, city) VALUES
    ('admin@fintks.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Admin', 'User', '+966501234567', 'King Fahd Road, Riyadh', 'Riyadh'),
    ('customer@fintks.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Ahmed', 'Al-Saud', '+966507654321', 'Prince Sultan Street, Jeddah', 'Jeddah'),
    ('user@fintks.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Fatima', 'Al-Zahra', '+966508765432', 'King Abdullah Road, Dammam', 'Dammam')
ON CONFLICT (id) DO NOTHING;

-- Insert sample reviews
INSERT INTO reviews (user_id, product_id, rating, title, comment, is_verified_purchase) VALUES
    (2, 1, 5, 'Excellent Phone!', 'The iPhone 15 Pro is amazing. The camera quality is outstanding and the battery life is great.', true),
    (3, 1, 4, 'Great but expensive', 'Very good phone with excellent features, but quite expensive.', true),
    (2, 2, 5, 'Perfect for work', 'The MacBook Air M2 is perfect for my work needs. Fast and reliable.', true),
    (3, 4, 4, 'Comfortable shoes', 'Very comfortable running shoes. Good for daily use.', true),
    (2, 6, 5, 'Amazing TV', 'The picture quality is incredible. Highly recommended!', true)
ON CONFLICT (id) DO NOTHING;

-- Insert sample cart items
INSERT INTO cart (user_id, product_id, quantity) VALUES
    (2, 1, 1),
    (2, 4, 2),
    (3, 2, 1),
    (3, 8, 1)
ON CONFLICT (user_id, product_id) DO NOTHING;

-- Insert sample wishlist items
INSERT INTO wishlist (user_id, product_id) VALUES
    (2, 3),
    (2, 9),
    (3, 1),
    (3, 6)
ON CONFLICT (user_id, product_id) DO NOTHING;

-- Insert sample orders
INSERT INTO orders (user_id, order_number, status, total_amount, shipping_address, shipping_city, shipping_phone, payment_method, payment_status) VALUES
    (2, 'ORD-2024-001', 'delivered', 1129.98, 'King Fahd Road, Riyadh', 'Riyadh', '+966501234567', 'cash_on_delivery', 'paid'),
    (3, 'ORD-2024-002', 'processing', 89.99, 'Prince Sultan Street, Jeddah', 'Jeddah', '+966507654321', 'cash_on_delivery', 'pending')
ON CONFLICT (id) DO NOTHING;

-- Insert sample order items
INSERT INTO order_items (order_id, product_id, product_name, product_price, quantity, total_price) VALUES
    (1, 1, 'iPhone 15 Pro', 999.99, 1, 999.99),
    (1, 4, 'Nike Air Max 270', 129.99, 1, 129.99),
    (2, 7, 'Coffee Maker', 89.99, 1, 89.99)
ON CONFLICT (id) DO NOTHING; 