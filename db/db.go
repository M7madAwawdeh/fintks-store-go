package db

import (
	"ai-catalog/auth"
	"ai-catalog/graph"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// Connect establishes a connection to PostgreSQL database
func Connect() error {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:password@localhost:5432/ai_catalog?sslmode=disable"
	}

	var err error
	DB, err = sql.Open("postgres", databaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	// Create all tables
	if err := createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %v", err)
	}

	fmt.Println("Database connected successfully")
	return nil
}

// createTables creates all necessary tables by executing init.sql
func createTables() error {
	initSQL, err := os.ReadFile("init.sql")
	if err != nil {
		return fmt.Errorf("failed to read init.sql: %v", err)
	}
	_, err = DB.Exec(string(initSQL))
	if err != nil {
		return fmt.Errorf("failed to execute init.sql: %v", err)
	}
	return nil
}

// GetUserByEmail retrieves a user by email
func GetUserByEmail(email string) (*graph.User, error) {
	query := `
		SELECT id, email, first_name, last_name, phone, address, city, country, created_at, updated_at
		FROM users WHERE email = $1
	`
	user := &graph.User{}
	err := DB.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.FirstName, &user.LastName,
		&user.Phone, &user.Address, &user.City, &user.Country,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByID retrieves a user by ID
func GetUserByID(id int) (*graph.User, error) {
	query := `
		SELECT id, email, first_name, last_name, phone, address, city, country, created_at, updated_at
		FROM users WHERE id = $1
	`
	user := &graph.User{}
	err := DB.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.FirstName, &user.LastName,
		&user.Phone, &user.Address, &user.City, &user.Country,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// CreateUser creates a new user
func CreateUser(email, password, firstName, lastName, phone, address, city string) (*graph.User, error) {
	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	query := `
		INSERT INTO users (email, password_hash, first_name, last_name, phone, address, city)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, email, first_name, last_name, phone, address, city, country, created_at, updated_at
	`
	user := &graph.User{}
	err = DB.QueryRow(query, email, passwordHash, firstName, lastName, phone, address, city).Scan(
		&user.ID, &user.Email, &user.FirstName, &user.LastName,
		&user.Phone, &user.Address, &user.City, &user.Country,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// AuthenticateUser authenticates a user with email and password
func AuthenticateUser(email, password string) (*graph.User, error) {
	query := `SELECT password_hash FROM users WHERE email = $1`
	var passwordHash string
	err := DB.QueryRow(query, email).Scan(&passwordHash)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if !auth.CheckPassword(password, passwordHash) {
		return nil, fmt.Errorf("invalid password")
	}

	return GetUserByEmail(email)
}

// UpdateUser updates user information
func UpdateUser(id int, firstName, lastName, phone, address, city string) (*graph.User, error) {
	query := `
		UPDATE users 
		SET first_name = $2, last_name = $3, phone = $4, address = $5, city = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, email, first_name, last_name, phone, address, city, country, created_at, updated_at
	`
	user := &graph.User{}
	err := DB.QueryRow(query, id, firstName, lastName, phone, address, city).Scan(
		&user.ID, &user.Email, &user.FirstName, &user.LastName,
		&user.Phone, &user.Address, &user.City, &user.Country,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetCategories retrieves all categories
func GetCategories() ([]graph.Category, error) {
	query := `SELECT id, name, description, image_url, parent_id, created_at FROM categories ORDER BY name`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []graph.Category
	for rows.Next() {
		var category graph.Category
		err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.ImageURL, &category.ParentID, &category.CreatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

// GetCategory retrieves a category by ID
func GetCategory(id int) (*graph.Category, error) {
	query := `SELECT id, name, description, image_url, parent_id, created_at FROM categories WHERE id = $1`
	category := &graph.Category{}
	err := DB.QueryRow(query, id).Scan(&category.ID, &category.Name, &category.Description, &category.ImageURL, &category.ParentID, &category.CreatedAt)
	if err != nil {
		return nil, err
	}
	return category, nil
}

// GetProducts retrieves products with optional filters
func GetProducts(categoryID *int, search *string, minPrice, maxPrice *float64, isFeatured *bool, page, limit int) ([]graph.Product, error) {
	baseQuery := `
		SELECT p.id, p.name, p.price, p.original_price, p.category_id, p.description, p.short_description,
		       p.image_url, p.stock_quantity, p.sku, p.weight, p.dimensions, p.is_active, p.is_featured,
		       p.created_at, p.updated_at,
		       c.id, c.name, c.description, c.image_url, c.parent_id, c.created_at
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.is_active = true
	`
	
	var conditions []string
	var args []interface{}
	argCount := 1

	if categoryID != nil {
		conditions = append(conditions, fmt.Sprintf("p.category_id = $%d", argCount))
		args = append(args, *categoryID)
		argCount++
	}

	if search != nil && *search != "" {
		conditions = append(conditions, fmt.Sprintf("(p.name ILIKE $%d OR p.description ILIKE $%d)", argCount, argCount))
		args = append(args, "%"+*search+"%")
		argCount++
	}

	if minPrice != nil {
		conditions = append(conditions, fmt.Sprintf("p.price >= $%d", argCount))
		args = append(args, *minPrice)
		argCount++
	}

	if maxPrice != nil {
		conditions = append(conditions, fmt.Sprintf("p.price <= $%d", argCount))
		args = append(args, *maxPrice)
		argCount++
	}

	if isFeatured != nil {
		conditions = append(conditions, fmt.Sprintf("p.is_featured = $%d", argCount))
		args = append(args, *isFeatured)
		argCount++
	}

	if len(conditions) > 0 {
		baseQuery += " AND " + conditions[0]
		for i := 1; i < len(conditions); i++ {
			baseQuery += " AND " + conditions[i]
		}
	}

	baseQuery += " ORDER BY p.created_at DESC"

	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount, argCount+1)
		args = append(args, limit, offset)
	}

	rows, err := DB.Query(baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []graph.Product
	for rows.Next() {
		var product graph.Product
		var category graph.Category
		err := rows.Scan(
			&product.ID, &product.Name, &product.Price, &product.OriginalPrice, &product.CategoryID,
			&product.Description, &product.ShortDescription, &product.ImageURL, &product.StockQuantity,
			&product.SKU, &product.Weight, &product.Dimensions, &product.IsActive, &product.IsFeatured,
			&product.CreatedAt, &product.UpdatedAt,
			&category.ID, &category.Name, &category.Description, &category.ImageURL, &category.ParentID, &category.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		product.Category = &category
		products = append(products, product)
	}
	return products, nil
}

// GetProduct retrieves a product by ID
func GetProduct(id int) (*graph.Product, error) {
	query := `
		SELECT p.id, p.name, p.price, p.original_price, p.category_id, p.description, p.short_description,
		       p.image_url, p.stock_quantity, p.sku, p.weight, p.dimensions, p.is_active, p.is_featured,
		       p.created_at, p.updated_at,
		       c.id, c.name, c.description, c.image_url, c.parent_id, c.created_at
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id = $1 AND p.is_active = true
	`
	
	var product graph.Product
	var category graph.Category
	err := DB.QueryRow(query, id).Scan(
		&product.ID, &product.Name, &product.Price, &product.OriginalPrice, &product.CategoryID,
		&product.Description, &product.ShortDescription, &product.ImageURL, &product.StockQuantity,
		&product.SKU, &product.Weight, &product.Dimensions, &product.IsActive, &product.IsFeatured,
		&product.CreatedAt, &product.UpdatedAt,
		&category.ID, &category.Name, &category.Description, &category.ImageURL, &category.ParentID, &category.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	product.Category = &category
	return &product, nil
}

// GetFeaturedProducts retrieves featured products
func GetFeaturedProducts() ([]graph.Product, error) {
	isFeatured := true
	return GetProducts(nil, nil, nil, nil, &isFeatured, 0, 0)
}

// SearchProducts searches products by query
func SearchProducts(query string) (*graph.SearchResult, error) {
	products, err := GetProducts(nil, &query, nil, nil, nil, 0, 0)
	if err != nil {
		return nil, err
	}

	return &graph.SearchResult{
		Products: products,
		Count:    len(products),
		Query:    query,
	}, nil
}

// GetCartItems retrieves cart items for a user
func GetCartItems(userID int) ([]graph.CartItem, error) {
	query := `
		SELECT ci.id, ci.user_id, ci.product_id, ci.quantity, ci.created_at,
		       p.id, p.name, p.price, p.original_price, p.category_id, p.description, p.short_description,
		       p.image_url, p.stock_quantity, p.sku, p.weight, p.dimensions, p.is_active, p.is_featured,
		       p.created_at, p.updated_at,
		       c.id, c.name, c.description, c.image_url, c.parent_id, c.created_at
		FROM cart ci
		JOIN products p ON ci.product_id = p.id
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE ci.user_id = $1
		ORDER BY ci.created_at DESC
	`
	
	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cartItems []graph.CartItem
	for rows.Next() {
		var cartItem graph.CartItem
		var product graph.Product
		var category graph.Category
		err := rows.Scan(
			&cartItem.ID, &cartItem.UserID, &cartItem.ProductID, &cartItem.Quantity, &cartItem.CreatedAt,
			&product.ID, &product.Name, &product.Price, &product.OriginalPrice, &product.CategoryID,
			&product.Description, &product.ShortDescription, &product.ImageURL, &product.StockQuantity,
			&product.SKU, &product.Weight, &product.Dimensions, &product.IsActive, &product.IsFeatured,
			&product.CreatedAt, &product.UpdatedAt,
			&category.ID, &category.Name, &category.Description, &category.ImageURL, &category.ParentID, &category.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		product.Category = &category
		cartItem.Product = &product
		cartItems = append(cartItems, cartItem)
	}
	return cartItems, nil
}

// AddToCart adds a product to cart
func AddToCart(userID, productID, quantity int) (*graph.CartItem, error) {
	// Check if product exists and has enough stock
	product, err := GetProduct(productID)
	if err != nil {
		return nil, fmt.Errorf("product not found")
	}
	if product.StockQuantity < quantity {
		return nil, fmt.Errorf("insufficient stock")
	}

	// Check if item already exists in cart
	var existingID int
	checkQuery := `SELECT id FROM cart WHERE user_id = $1 AND product_id = $2`
	err = DB.QueryRow(checkQuery, userID, productID).Scan(&existingID)
	
	if err == sql.ErrNoRows {
		// Insert new cart item
		insertQuery := `
			INSERT INTO cart (user_id, product_id, quantity)
			VALUES ($1, $2, $3)
			RETURNING id, user_id, product_id, quantity, created_at
		`
		cartItem := &graph.CartItem{}
		err = DB.QueryRow(insertQuery, userID, productID, quantity).Scan(
			&cartItem.ID, &cartItem.UserID, &cartItem.ProductID, &cartItem.Quantity, &cartItem.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		cartItem.Product = product
		return cartItem, nil
	} else if err != nil {
		return nil, err
	} else {
		// Update existing cart item
		updateQuery := `
			UPDATE cart SET quantity = quantity + $3, created_at = CURRENT_TIMESTAMP
			WHERE id = $1
			RETURNING id, user_id, product_id, quantity, created_at
		`
		cartItem := &graph.CartItem{}
		err = DB.QueryRow(updateQuery, existingID, userID, quantity).Scan(
			&cartItem.ID, &cartItem.UserID, &cartItem.ProductID, &cartItem.Quantity, &cartItem.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		cartItem.Product = product
		return cartItem, nil
	}
}

// UpdateCartItem updates cart item quantity
func UpdateCartItem(id, quantity int) (*graph.CartItem, error) {
	query := `
		UPDATE cart SET quantity = $2, created_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, user_id, product_id, quantity, created_at
	`
	cartItem := &graph.CartItem{}
	err := DB.QueryRow(query, id, quantity).Scan(
		&cartItem.ID, &cartItem.UserID, &cartItem.ProductID, &cartItem.Quantity, &cartItem.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Get product details
	product, err := GetProduct(cartItem.ProductID)
	if err != nil {
		return nil, err
	}
	cartItem.Product = product
	return cartItem, nil
}

// RemoveFromCart removes an item from cart
func RemoveFromCart(id int) error {
	query := `DELETE FROM cart WHERE id = $1`
	_, err := DB.Exec(query, id)
	return err
}

// ClearCart clears all items from user's cart
func ClearCart(userID int) error {
	query := `DELETE FROM cart WHERE user_id = $1`
	_, err := DB.Exec(query, userID)
	return err
}

// GetWishlistItems retrieves wishlist items for a user
func GetWishlistItems(userID int) ([]graph.WishlistItem, error) {
	query := `
		SELECT wi.id, wi.user_id, wi.product_id, wi.created_at,
		       p.id, p.name, p.price, p.original_price, p.category_id, p.description, p.short_description,
		       p.image_url, p.stock_quantity, p.sku, p.weight, p.dimensions, p.is_active, p.is_featured,
		       p.created_at, p.updated_at,
		       c.id, c.name, c.description, c.image_url, c.parent_id, c.created_at
		FROM wishlist wi
		JOIN products p ON wi.product_id = p.id
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE wi.user_id = $1
		ORDER BY wi.created_at DESC
	`
	
	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wishlistItems []graph.WishlistItem
	for rows.Next() {
		var wishlistItem graph.WishlistItem
		var product graph.Product
		var category graph.Category
		err := rows.Scan(
			&wishlistItem.ID, &wishlistItem.UserID, &wishlistItem.ProductID, &wishlistItem.CreatedAt,
			&product.ID, &product.Name, &product.Price, &product.OriginalPrice, &product.CategoryID,
			&product.Description, &product.ShortDescription, &product.ImageURL, &product.StockQuantity,
			&product.SKU, &product.Weight, &product.Dimensions, &product.IsActive, &product.IsFeatured,
			&product.CreatedAt, &product.UpdatedAt,
			&category.ID, &category.Name, &category.Description, &category.ImageURL, &category.ParentID, &category.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		product.Category = &category
		wishlistItem.Product = &product
		wishlistItems = append(wishlistItems, wishlistItem)
	}
	return wishlistItems, nil
}

// AddToWishlist adds a product to wishlist
func AddToWishlist(userID, productID int) (*graph.WishlistItem, error) {
	// Check if product exists
	product, err := GetProduct(productID)
	if err != nil {
		return nil, fmt.Errorf("product not found")
	}

	// Check if already in wishlist
	var existingID int
	checkQuery := `SELECT id FROM wishlist WHERE user_id = $1 AND product_id = $2`
	err = DB.QueryRow(checkQuery, userID, productID).Scan(&existingID)
	if err == nil {
		return nil, fmt.Errorf("product already in wishlist")
	}

	// Add to wishlist
	insertQuery := `
		INSERT INTO wishlist (user_id, product_id)
		VALUES ($1, $2)
		RETURNING id, user_id, product_id, created_at
	`
	wishlistItem := &graph.WishlistItem{}
	err = DB.QueryRow(insertQuery, userID, productID).Scan(
		&wishlistItem.ID, &wishlistItem.UserID, &wishlistItem.ProductID, &wishlistItem.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	wishlistItem.Product = product
	return wishlistItem, nil
}

// RemoveFromWishlist removes a product from wishlist
func RemoveFromWishlist(userID, productID int) error {
	query := `DELETE FROM wishlist WHERE user_id = $1 AND product_id = $2`
	_, err := DB.Exec(query, userID, productID)
	return err
}

// GetOrders retrieves orders for a user
func GetOrders(userID int) ([]graph.Order, error) {
	query := `
		SELECT o.id, o.user_id, o.order_number, o.total_amount, o.status, o.payment_method,
		       o.payment_status, o.shipping_address, o.shipping_city, o.shipping_phone, o.notes,
		       o.created_at, o.updated_at
		FROM orders o
		WHERE o.user_id = $1
		ORDER BY o.created_at DESC
	`
	
	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []graph.Order
	for rows.Next() {
		var order graph.Order
		err := rows.Scan(
			&order.ID, &order.UserID, &order.OrderNumber, &order.TotalAmount, &order.Status,
			&order.PaymentMethod, &order.PaymentStatus, &order.ShippingAddress, &order.ShippingCity,
			&order.ShippingPhone, &order.Notes, &order.CreatedAt, &order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

// GetOrder retrieves an order by ID
func GetOrder(id int) (*graph.Order, error) {
	query := `
		SELECT o.id, o.user_id, o.order_number, o.total_amount, o.status, o.payment_method,
		       o.payment_status, o.shipping_address, o.shipping_city, o.shipping_phone, o.notes,
		       o.created_at, o.updated_at
		FROM orders o
		WHERE o.id = $1
	`
	
	order := &graph.Order{}
	err := DB.QueryRow(query, id).Scan(
		&order.ID, &order.UserID, &order.OrderNumber, &order.TotalAmount, &order.Status,
		&order.PaymentMethod, &order.PaymentStatus, &order.ShippingAddress, &order.ShippingCity,
		&order.ShippingPhone, &order.Notes, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Get order items
	orderItems, err := GetOrderItems(id)
	if err != nil {
		return nil, err
	}
	order.Items = orderItems

	return order, nil
}

// GetOrderItems retrieves items for an order
func GetOrderItems(orderID int) ([]graph.OrderItem, error) {
	query := `
		SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, oi.price, oi.created_at,
		       p.id, p.name, p.price, p.original_price, p.category_id, p.description, p.short_description,
		       p.image_url, p.stock_quantity, p.sku, p.weight, p.dimensions, p.is_active, p.is_featured,
		       p.created_at, p.updated_at,
		       c.id, c.name, c.description, c.image_url, c.parent_id, c.created_at
		FROM order_items oi
		JOIN products p ON oi.product_id = p.id
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE oi.order_id = $1
	`
	
	rows, err := DB.Query(query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderItems []graph.OrderItem
	for rows.Next() {
		var orderItem graph.OrderItem
		var product graph.Product
		var category graph.Category
		err := rows.Scan(
			&orderItem.ID, &orderItem.OrderID, &orderItem.ProductID, &orderItem.Quantity, &orderItem.Price, &orderItem.CreatedAt,
			&product.ID, &product.Name, &product.Price, &product.OriginalPrice, &product.CategoryID,
			&product.Description, &product.ShortDescription, &product.ImageURL, &product.StockQuantity,
			&product.SKU, &product.Weight, &product.Dimensions, &product.IsActive, &product.IsFeatured,
			&product.CreatedAt, &product.UpdatedAt,
			&category.ID, &category.Name, &category.Description, &category.ImageURL, &category.ParentID, &category.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		product.Category = &category
		orderItem.Product = &product
		orderItems = append(orderItems, orderItem)
	}
	return orderItems, nil
}

// CreateOrder creates a new order
func CreateOrder(userID int, orderNumber, shippingAddress, shippingCity, shippingPhone, paymentMethod, notes string, totalAmount float64) (*graph.Order, error) {
	tx, err := DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Create order
	orderQuery := `
		INSERT INTO orders (user_id, order_number, total_amount, status, payment_method, payment_status,
		                   shipping_address, shipping_city, shipping_phone, notes)
		VALUES ($1, $2, $3, 'pending', $4, 'pending', $5, $6, $7, $8)
		RETURNING id, user_id, order_number, total_amount, status, payment_method, payment_status,
		          shipping_address, shipping_city, shipping_phone, notes, created_at, updated_at
	`
	
	order := &graph.Order{}
	err = tx.QueryRow(orderQuery, userID, orderNumber, totalAmount, paymentMethod,
		shippingAddress, shippingCity, shippingPhone, notes).Scan(
		&order.ID, &order.UserID, &order.OrderNumber, &order.TotalAmount, &order.Status,
		&order.PaymentMethod, &order.PaymentStatus, &order.ShippingAddress, &order.ShippingCity,
		&order.ShippingPhone, &order.Notes, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Get cart items
	cartItems, err := GetCartItems(userID)
	if err != nil {
		return nil, err
	}

	// Create order items and update product stock
	for _, cartItem := range cartItems {
		// Create order item
		orderItemQuery := `
			INSERT INTO order_items (order_id, product_id, quantity, price)
			VALUES ($1, $2, $3, $4)
		`
		_, err = tx.Exec(orderItemQuery, order.ID, cartItem.ProductID, cartItem.Quantity, cartItem.Product.Price)
		if err != nil {
			return nil, err
		}

		// Update product stock
		updateStockQuery := `
			UPDATE products SET stock_quantity = stock_quantity - $2, updated_at = CURRENT_TIMESTAMP
			WHERE id = $1
		`
		_, err = tx.Exec(updateStockQuery, cartItem.ProductID, cartItem.Quantity)
		if err != nil {
			return nil, err
		}
	}

	// Clear cart
	_, err = tx.Exec("DELETE FROM cart WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	// Get order items for response
	orderItems, err := GetOrderItems(order.ID)
	if err != nil {
		return nil, err
	}
	order.Items = orderItems

	return order, nil
}

// UpdateOrderStatus updates order status
func UpdateOrderStatus(id int, status string) (*graph.Order, error) {
	query := `
		UPDATE orders SET status = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, user_id, order_number, total_amount, status, payment_method, payment_status,
		          shipping_address, shipping_city, shipping_phone, notes, created_at, updated_at
	`
	
	order := &graph.Order{}
	err := DB.QueryRow(query, id, status).Scan(
		&order.ID, &order.UserID, &order.OrderNumber, &order.TotalAmount, &order.Status,
		&order.PaymentMethod, &order.PaymentStatus, &order.ShippingAddress, &order.ShippingCity,
		&order.ShippingPhone, &order.Notes, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Get order items
	orderItems, err := GetOrderItems(order.ID)
	if err != nil {
		return nil, err
	}
	order.Items = orderItems

	return order, nil
}

// GetProductReviews retrieves reviews for a product
func GetProductReviews(productID int) ([]graph.Review, error) {
	query := `
		SELECT r.id, r.user_id, r.product_id, r.rating, r.title, r.comment, r.is_verified_purchase,
		       r.created_at, r.updated_at,
		       u.id, u.email, u.first_name, u.last_name, u.phone, u.address, u.city, u.country, u.created_at, u.updated_at
		FROM reviews r
		JOIN users u ON r.user_id = u.id
		WHERE r.product_id = $1
		ORDER BY r.created_at DESC
	`
	
	rows, err := DB.Query(query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []graph.Review
	for rows.Next() {
		var review graph.Review
		var user graph.User
		err := rows.Scan(
			&review.ID, &review.UserID, &review.ProductID, &review.Rating, &review.Title, &review.Comment,
			&review.IsVerifiedPurchase, &review.CreatedAt, &review.UpdatedAt,
			&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Phone, &user.Address, &user.City, &user.Country, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		review.User = &user
		reviews = append(reviews, review)
	}
	return reviews, nil
}

// CreateReview creates a new review
func CreateReview(userID, productID, rating int, title, comment string, isVerifiedPurchase bool) (*graph.Review, error) {
	query := `
		INSERT INTO reviews (user_id, product_id, rating, title, comment, is_verified_purchase)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, product_id, rating, title, comment, is_verified_purchase, created_at, updated_at
	`
	
	review := &graph.Review{}
	err := DB.QueryRow(query, userID, productID, rating, title, comment, isVerifiedPurchase).Scan(
		&review.ID, &review.UserID, &review.ProductID, &review.Rating, &review.Title, &review.Comment,
		&review.IsVerifiedPurchase, &review.CreatedAt, &review.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Get user details
	user, err := GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	review.User = user

	return review, nil
}

// CheckUserPurchasedProduct checks if user has purchased the product
func CheckUserPurchasedProduct(userID, productID int) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM order_items oi
			JOIN orders o ON oi.order_id = o.id
			WHERE o.user_id = $1 AND oi.product_id = $2 AND o.status = 'delivered'
		)
	`
	var exists bool
	err := DB.QueryRow(query, userID, productID).Scan(&exists)
	return exists, err
} 