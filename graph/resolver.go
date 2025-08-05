package graph

import (
	"ai-catalog/auth"
	"ai-catalog/handlers"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/graphql-go/graphql"
	_ "github.com/lib/pq"
)

// Global database connection
var DB *sql.DB

// SetDB sets the database connection for the resolver
func SetDB(db *sql.DB) {
	DB = db
}

// GraphQL Types
var UserType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id":        &graphql.Field{Type: graphql.Int},
		"email":     &graphql.Field{Type: graphql.String},
		"firstName": &graphql.Field{Type: graphql.String},
		"lastName":  &graphql.Field{Type: graphql.String},
		"phone":     &graphql.Field{Type: graphql.String},
		"address":   &graphql.Field{Type: graphql.String},
		"city":      &graphql.Field{Type: graphql.String},
		"country":   &graphql.Field{Type: graphql.String},
		"createdAt": &graphql.Field{Type: graphql.String},
		"updatedAt": &graphql.Field{Type: graphql.String},
	},
})

var CategoryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Category",
	Fields: graphql.Fields{
		"id":          &graphql.Field{Type: graphql.Int},
		"name":        &graphql.Field{Type: graphql.String},
		"description": &graphql.Field{Type: graphql.String},
		"imageUrl":    &graphql.Field{Type: graphql.String},
		"parentId":    &graphql.Field{Type: graphql.Int},
		"createdAt":   &graphql.Field{Type: graphql.String},
	},
})

var ProductType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Product",
	Fields: graphql.Fields{
		"id":               &graphql.Field{Type: graphql.Int},
		"name":             &graphql.Field{Type: graphql.String},
		"price":            &graphql.Field{Type: graphql.Float},
		"originalPrice":    &graphql.Field{Type: graphql.Float},
		"categoryId":       &graphql.Field{Type: graphql.Int},
		"category":         &graphql.Field{Type: CategoryType},
		"description":      &graphql.Field{Type: graphql.String},
		"shortDescription": &graphql.Field{Type: graphql.String},
		"imageUrl":         &graphql.Field{Type: graphql.String},
		"stockQuantity":    &graphql.Field{Type: graphql.Int},
		"sku":              &graphql.Field{Type: graphql.String},
		"weight":           &graphql.Field{Type: graphql.Float},
		"dimensions":       &graphql.Field{Type: graphql.String},
		"isActive":         &graphql.Field{Type: graphql.Boolean},
		"isFeatured":       &graphql.Field{Type: graphql.Boolean},
		"createdAt":        &graphql.Field{Type: graphql.String},
		"updatedAt":        &graphql.Field{Type: graphql.String},
		"averageRating":    &graphql.Field{Type: graphql.Float},
		"reviewCount":      &graphql.Field{Type: graphql.Int},
		"isInWishlist":     &graphql.Field{Type: graphql.Boolean},
		"isLiked":          &graphql.Field{Type: graphql.Boolean},
	},
})

var CartItemType = graphql.NewObject(graphql.ObjectConfig{
	Name: "CartItem",
	Fields: graphql.Fields{
		"id":        &graphql.Field{Type: graphql.Int},
		"userId":    &graphql.Field{Type: graphql.Int},
		"productId": &graphql.Field{Type: graphql.Int},
		"product":   &graphql.Field{Type: ProductType},
		"quantity":  &graphql.Field{Type: graphql.Int},
		"createdAt": &graphql.Field{Type: graphql.String},
		"updatedAt": &graphql.Field{Type: graphql.String},
	},
})

var CartSummaryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "CartSummary",
	Fields: graphql.Fields{
		"items":      &graphql.Field{Type: graphql.NewList(CartItemType)},
		"totalItems": &graphql.Field{Type: graphql.Int},
		"totalPrice": &graphql.Field{Type: graphql.Float},
	},
})

var WishlistItemType = graphql.NewObject(graphql.ObjectConfig{
	Name: "WishlistItem",
	Fields: graphql.Fields{
		"id":        &graphql.Field{Type: graphql.Int},
		"userId":    &graphql.Field{Type: graphql.Int},
		"productId": &graphql.Field{Type: graphql.Int},
		"product":   &graphql.Field{Type: ProductType},
		"createdAt": &graphql.Field{Type: graphql.String},
	},
})

var OrderItemType = graphql.NewObject(graphql.ObjectConfig{
	Name: "OrderItem",
	Fields: graphql.Fields{
		"id":           &graphql.Field{Type: graphql.Int},
		"orderId":      &graphql.Field{Type: graphql.Int},
		"productId":    &graphql.Field{Type: graphql.Int},
		"productName":  &graphql.Field{Type: graphql.String},
		"productPrice": &graphql.Field{Type: graphql.Float},
		"quantity":     &graphql.Field{Type: graphql.Int},
		"totalPrice":   &graphql.Field{Type: graphql.Float},
		"createdAt":    &graphql.Field{Type: graphql.String},
	},
})

var OrderType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Order",
	Fields: graphql.Fields{
		"id":              &graphql.Field{Type: graphql.Int},
		"userId":          &graphql.Field{Type: graphql.Int},
		"user":            &graphql.Field{Type: UserType},
		"orderNumber":     &graphql.Field{Type: graphql.String},
		"status":          &graphql.Field{Type: graphql.String},
		"totalAmount":     &graphql.Field{Type: graphql.Float},
		"shippingAddress": &graphql.Field{Type: graphql.String},
		"shippingCity":    &graphql.Field{Type: graphql.String},
		"shippingCountry": &graphql.Field{Type: graphql.String},
		"shippingPhone":   &graphql.Field{Type: graphql.String},
		"paymentMethod":   &graphql.Field{Type: graphql.String},
		"paymentStatus":   &graphql.Field{Type: graphql.String},
		"notes":           &graphql.Field{Type: graphql.String},
		"createdAt":       &graphql.Field{Type: graphql.String},
		"updatedAt":       &graphql.Field{Type: graphql.String},
		"items":           &graphql.Field{Type: graphql.NewList(OrderItemType)},
	},
})

var ReviewType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Review",
	Fields: graphql.Fields{
		"id":                &graphql.Field{Type: graphql.Int},
		"userId":            &graphql.Field{Type: graphql.Int},
		"user":              &graphql.Field{Type: UserType},
		"productId":         &graphql.Field{Type: graphql.Int},
		"rating":            &graphql.Field{Type: graphql.Int},
		"title":             &graphql.Field{Type: graphql.String},
		"comment":           &graphql.Field{Type: graphql.String},
		"isVerifiedPurchase": &graphql.Field{Type: graphql.Boolean},
		"createdAt":         &graphql.Field{Type: graphql.String},
		"updatedAt":         &graphql.Field{Type: graphql.String},
		"likeCount":         &graphql.Field{Type: graphql.Int},
		"isLiked":           &graphql.Field{Type: graphql.Boolean},
	},
})

var AuthResponseType = graphql.NewObject(graphql.ObjectConfig{
	Name: "AuthResponse",
	Fields: graphql.Fields{
		"user":  &graphql.Field{Type: UserType},
		"token": &graphql.Field{Type: graphql.String},
	},
})

// Root Query
var RootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"me": &graphql.Field{
			Type: UserType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if user, ok := p.Context.Value("user").(*User); ok {
					return user, nil
				}
				return nil, nil
			},
		},
		"categories": &graphql.Field{
			Type: graphql.NewList(CategoryType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				rows, err := DB.Query("SELECT id, name, description, image_url, parent_id, created_at FROM categories ORDER BY name")
				if err != nil {
					return nil, err
				}
				defer rows.Close()

				var categories []Category
				for rows.Next() {
					var c Category
					err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.ImageURL, &c.ParentID, &c.CreatedAt)
					if err != nil {
						return nil, err
					}
					categories = append(categories, c)
				}
				return categories, nil
			},
		},
		"products": &graphql.Field{
			Type: graphql.NewList(ProductType),
			Args: graphql.FieldConfigArgument{
				"categoryId":  &graphql.ArgumentConfig{Type: graphql.Int},
				"search":      &graphql.ArgumentConfig{Type: graphql.String},
				"minPrice":    &graphql.ArgumentConfig{Type: graphql.Float},
				"maxPrice":    &graphql.ArgumentConfig{Type: graphql.Float},
				"isFeatured":  &graphql.ArgumentConfig{Type: graphql.Boolean},
				"page":        &graphql.ArgumentConfig{Type: graphql.Int},
				"limit":       &graphql.ArgumentConfig{Type: graphql.Int},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				query := "SELECT id, name, price, original_price, category_id, description, short_description, image_url, stock_quantity, sku, weight, dimensions, is_active, is_featured, created_at, updated_at FROM products WHERE is_active = true"
				var args []interface{}
				argCount := 1

				if categoryID, ok := p.Args["categoryId"].(int); ok {
					query += fmt.Sprintf(" AND category_id = $%d", argCount)
					args = append(args, categoryID)
					argCount++
				}

				if search, ok := p.Args["search"].(string); ok && search != "" {
					query += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", argCount, argCount)
					args = append(args, "%"+search+"%")
					argCount++
				}

				if minPrice, ok := p.Args["minPrice"].(float64); ok {
					query += fmt.Sprintf(" AND price >= $%d", argCount)
					args = append(args, minPrice)
					argCount++
				}

				if maxPrice, ok := p.Args["maxPrice"].(float64); ok {
					query += fmt.Sprintf(" AND price <= $%d", argCount)
					args = append(args, maxPrice)
					argCount++
				}

				if isFeatured, ok := p.Args["isFeatured"].(bool); ok {
					query += fmt.Sprintf(" AND is_featured = $%d", argCount)
					args = append(args, isFeatured)
					argCount++
				}

				query += " ORDER BY created_at DESC"

				// Add pagination
				if limit, ok := p.Args["limit"].(int); ok && limit > 0 {
					query += fmt.Sprintf(" LIMIT $%d", argCount)
					args = append(args, limit)
					argCount++
				}

				rows, err := DB.Query(query, args...)
				if err != nil {
					return nil, err
				}
				defer rows.Close()

				var products []Product
				for rows.Next() {
					var p Product
					err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.OriginalPrice, &p.CategoryID, &p.Description, &p.ShortDescription, &p.ImageURL, &p.StockQuantity, &p.SKU, &p.Weight, &p.Dimensions, &p.IsActive, &p.IsFeatured, &p.CreatedAt, &p.UpdatedAt)
					if err != nil {
						return nil, err
					}
					products = append(products, p)
				}
				return products, nil
			},
		},
		"product": &graphql.Field{
			Type: ProductType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id := p.Args["id"].(int)
				var product Product
				err := DB.QueryRow("SELECT id, name, price, original_price, category_id, description, short_description, image_url, stock_quantity, sku, weight, dimensions, is_active, is_featured, created_at, updated_at FROM products WHERE id = $1", id).Scan(&product.ID, &product.Name, &product.Price, &product.OriginalPrice, &product.CategoryID, &product.Description, &product.ShortDescription, &product.ImageURL, &product.StockQuantity, &product.SKU, &product.Weight, &product.Dimensions, &product.IsActive, &product.IsFeatured, &product.CreatedAt, &product.UpdatedAt)
				if err != nil {
					return nil, err
				}
				return product, nil
			},
		},
		"featuredProducts": &graphql.Field{
			Type: graphql.NewList(ProductType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				rows, err := DB.Query("SELECT id, name, price, original_price, category_id, description, short_description, image_url, stock_quantity, sku, weight, dimensions, is_active, is_featured, created_at, updated_at FROM products WHERE is_active = true AND is_featured = true ORDER BY created_at DESC LIMIT 10")
				if err != nil {
					return nil, err
				}
				defer rows.Close()

				var products []Product
				for rows.Next() {
					var p Product
					err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.OriginalPrice, &p.CategoryID, &p.Description, &p.ShortDescription, &p.ImageURL, &p.StockQuantity, &p.SKU, &p.Weight, &p.Dimensions, &p.IsActive, &p.IsFeatured, &p.CreatedAt, &p.UpdatedAt)
					if err != nil {
						return nil, err
					}
					products = append(products, p)
				}
				return products, nil
			},
		},
		"cart": &graphql.Field{
			Type: CartSummaryType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				user, ok := p.Context.Value("user").(*User)
				if !ok {
					return nil, fmt.Errorf("user not authenticated")
				}

				rows, err := DB.Query(`
					SELECT c.id, c.user_id, c.product_id, c.quantity, c.created_at, c.updated_at,
						   p.id, p.name, p.price, p.original_price, p.category_id, p.description, p.short_description, p.image_url, p.stock_quantity, p.sku, p.weight, p.dimensions, p.is_active, p.is_featured, p.created_at, p.updated_at
					FROM cart c
					JOIN products p ON c.product_id = p.id
					WHERE c.user_id = $1
				`, user.ID)
				if err != nil {
					return nil, err
				}
				defer rows.Close()

				var items []*CartItem
				var totalPrice float64
				var totalItems int
				for rows.Next() {
					var item CartItem
					var product Product
					err := rows.Scan(
						&item.ID, &item.UserID, &item.ProductID, &item.Quantity, &item.CreatedAt, &item.UpdatedAt,
						&product.ID, &product.Name, &product.Price, &product.OriginalPrice, &product.CategoryID, &product.Description, &product.ShortDescription, &product.ImageURL, &product.StockQuantity, &product.SKU, &product.Weight, &product.Dimensions, &product.IsActive, &product.IsFeatured, &product.CreatedAt, &product.UpdatedAt,
					)
					if err != nil {
						return nil, err
					}
					item.Product = &product
					items = append(items, &item)
					totalPrice += product.Price * float64(item.Quantity)
					totalItems += item.Quantity
				}

				return CartSummary{
					Items:      items,
					TotalItems: totalItems,
					TotalPrice: totalPrice,
				}, nil
			},
		},
		"wishlist": &graphql.Field{
			Type: graphql.NewList(WishlistItemType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				user, ok := p.Context.Value("user").(*User)
				if !ok {
					return nil, fmt.Errorf("user not authenticated")
				}

				rows, err := DB.Query(`
					SELECT w.id, w.user_id, w.product_id, w.created_at,
						   p.id, p.name, p.price, p.original_price, p.category_id, p.description, p.short_description, p.image_url, p.stock_quantity, p.sku, p.weight, p.dimensions, p.is_active, p.is_featured, p.created_at, p.updated_at
					FROM wishlist w
					JOIN products p ON w.product_id = p.id
					WHERE w.user_id = $1
					ORDER BY w.created_at DESC
				`, user.ID)
				if err != nil {
					return nil, err
				}
				defer rows.Close()

				var items []WishlistItem
				for rows.Next() {
					var item WishlistItem
					var product Product
					err := rows.Scan(
						&item.ID, &item.UserID, &item.ProductID, &item.CreatedAt,
						&product.ID, &product.Name, &product.Price, &product.OriginalPrice, &product.CategoryID, &product.Description, &product.ShortDescription, &product.ImageURL, &product.StockQuantity, &product.SKU, &product.Weight, &product.Dimensions, &product.IsActive, &product.IsFeatured, &product.CreatedAt, &product.UpdatedAt,
					)
					if err != nil {
						return nil, err
					}
					item.Product = &product
					items = append(items, item)
				}

				return items, nil
			},
		},
		"orders": &graphql.Field{
			Type: graphql.NewList(OrderType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				user, ok := p.Context.Value("user").(*User)
				if !ok {
					return nil, fmt.Errorf("user not authenticated")
				}

				rows, err := DB.Query(`
					SELECT id, user_id, order_number, status, total_amount, shipping_address, shipping_city, shipping_country, shipping_phone, payment_method, payment_status, notes, created_at, updated_at
					FROM orders
					WHERE user_id = $1
					ORDER BY created_at DESC
				`, user.ID)
				if err != nil {
					return nil, err
				}
				defer rows.Close()

				var orders []Order
				for rows.Next() {
					var order Order
					err := rows.Scan(&order.ID, &order.UserID, &order.OrderNumber, &order.Status, &order.TotalAmount, &order.ShippingAddress, &order.ShippingCity, &order.ShippingCountry, &order.ShippingPhone, &order.PaymentMethod, &order.PaymentStatus, &order.Notes, &order.CreatedAt, &order.UpdatedAt)
					if err != nil {
						return nil, err
					}
					orders = append(orders, order)
				}

				return orders, nil
			},
		},
		"productReviews": &graphql.Field{
			Type: graphql.NewList(ReviewType),
			Args: graphql.FieldConfigArgument{
				"productId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				productID := p.Args["productId"].(int)

				rows, err := DB.Query(`
					SELECT r.id, r.user_id, r.product_id, r.rating, r.title, r.comment, r.is_verified_purchase, r.created_at, r.updated_at,
						   u.id, u.email, u.first_name, u.last_name, u.phone, u.address, u.city, u.country, u.created_at, u.updated_at
					FROM reviews r
					JOIN users u ON r.user_id = u.id
					WHERE r.product_id = $1
					ORDER BY r.created_at DESC
				`, productID)
				if err != nil {
					return nil, err
				}
				defer rows.Close()

				var reviews []Review
				for rows.Next() {
					var review Review
					var user User
					err := rows.Scan(
						&review.ID, &review.UserID, &review.ProductID, &review.Rating, &review.Title, &review.Comment, &review.IsVerifiedPurchase, &review.CreatedAt, &review.UpdatedAt,
						&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Phone, &user.Address, &user.City, &user.Country, &user.CreatedAt, &user.UpdatedAt,
					)
					if err != nil {
						return nil, err
					}
					review.User = &user
					reviews = append(reviews, review)
				}

				return reviews, nil
			},
		},
	},
})

// Root Mutation
var RootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"register": &graphql.Field{
			Type: AuthResponseType,
			Args: graphql.FieldConfigArgument{
				"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.NewInputObject(graphql.InputObjectConfig{
					Name: "RegisterInput",
					Fields: graphql.InputObjectConfigFieldMap{
						"email":     &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
						"password":  &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
						"firstName": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
						"lastName":  &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
						"phone":     &graphql.InputObjectFieldConfig{Type: graphql.String},
						"address":   &graphql.InputObjectFieldConfig{Type: graphql.String},
						"city":      &graphql.InputObjectFieldConfig{Type: graphql.String},
					},
				}))},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				input := p.Args["input"].(map[string]interface{})
				
				user, err := CreateUser(
					input["email"].(string),
					input["password"].(string),
					input["firstName"].(string),
					input["lastName"].(string),
					input["phone"].(string),
					input["address"].(string),
					input["city"].(string),
				)
				if err != nil {
					return nil, err
				}

				// Convert graph.User to auth.User
				authUser := &auth.User{
					ID:    user.ID,
					Email: user.Email,
				}
				token, err := auth.GenerateToken(authUser)
				if err != nil {
					return nil, err
				}

				return AuthResponse{
					User:  user,
					Token: token,
				}, nil
			},
		},
		"login": &graphql.Field{
			Type: AuthResponseType,
			Args: graphql.FieldConfigArgument{
				"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.NewInputObject(graphql.InputObjectConfig{
					Name: "LoginInput",
					Fields: graphql.InputObjectConfigFieldMap{
						"email":    &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
						"password": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
					},
				}))},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				input := p.Args["input"].(map[string]interface{})
				
				user, err := AuthenticateUser(input["email"].(string), input["password"].(string))
				if err != nil {
					return nil, err
				}

				// Convert graph.User to auth.User
				authUser := &auth.User{
					ID:    user.ID,
					Email: user.Email,
				}
				token, err := auth.GenerateToken(authUser)
				if err != nil {
					return nil, err
				}

				return AuthResponse{
					User:  user,
					Token: token,
				}, nil
			},
		},
		"addToCart": &graphql.Field{
			Type: CartItemType,
			Args: graphql.FieldConfigArgument{
				"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.NewInputObject(graphql.InputObjectConfig{
					Name: "AddToCartInput",
					Fields: graphql.InputObjectConfigFieldMap{
						"productId": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Int)},
						"quantity":  &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Int)},
					},
				}))},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				user, ok := p.Context.Value("user").(*User)
				if !ok {
					return nil, fmt.Errorf("user not authenticated")
				}

				input := p.Args["input"].(map[string]interface{})
				productID := input["productId"].(int)
				quantity := input["quantity"].(int)

				// Check if product exists and has stock
				var stockQuantity int
				err := DB.QueryRow("SELECT stock_quantity FROM products WHERE id = $1 AND is_active = true", productID).Scan(&stockQuantity)
				if err != nil {
					return nil, fmt.Errorf("product not found")
				}

				if stockQuantity < quantity {
					return nil, fmt.Errorf("insufficient stock")
				}

				// Add to cart (upsert)
				var cartItem CartItem
				query := `
					INSERT INTO cart (user_id, product_id, quantity)
					VALUES ($1, $2, $3)
					ON CONFLICT (user_id, product_id)
					DO UPDATE SET quantity = cart.quantity + $3, updated_at = CURRENT_TIMESTAMP
					RETURNING id, user_id, product_id, quantity, created_at, updated_at
				`
				err = DB.QueryRow(query, user.ID, productID, quantity).Scan(
					&cartItem.ID, &cartItem.UserID, &cartItem.ProductID, &cartItem.Quantity, &cartItem.CreatedAt, &cartItem.UpdatedAt,
				)
				if err != nil {
					return nil, err
				}

				return cartItem, nil
			},
		},
		"removeFromCart": &graphql.Field{
			Type: graphql.Boolean,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				user, ok := p.Context.Value("user").(*User)
				if !ok {
					return nil, fmt.Errorf("user not authenticated")
				}

				id := p.Args["id"].(int)
				result, err := DB.Exec("DELETE FROM cart WHERE id = $1 AND user_id = $2", id, user.ID)
				if err != nil {
					return nil, err
				}

				rowsAffected, _ := result.RowsAffected()
				return rowsAffected > 0, nil
			},
		},
		"addToWishlist": &graphql.Field{
			Type: WishlistItemType,
			Args: graphql.FieldConfigArgument{
				"productId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				user, ok := p.Context.Value("user").(*User)
				if !ok {
					return nil, fmt.Errorf("user not authenticated")
				}

				productID := p.Args["productId"].(int)

				// Check if product exists
				var productExists bool
				err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE id = $1 AND is_active = true)", productID).Scan(&productExists)
				if err != nil || !productExists {
					return nil, fmt.Errorf("product not found")
				}

				var wishlistItem WishlistItem
				query := `
					INSERT INTO wishlist (user_id, product_id)
					VALUES ($1, $2)
					ON CONFLICT (user_id, product_id) DO NOTHING
					RETURNING id, user_id, product_id, created_at
				`
				err = DB.QueryRow(query, user.ID, productID).Scan(
					&wishlistItem.ID, &wishlistItem.UserID, &wishlistItem.ProductID, &wishlistItem.CreatedAt,
				)
				if err != nil {
					return nil, err
				}

				return wishlistItem, nil
			},
		},
		"createOrder": &graphql.Field{
			Type: OrderType,
			Args: graphql.FieldConfigArgument{
				"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.NewInputObject(graphql.InputObjectConfig{
					Name: "CreateOrderInput",
					Fields: graphql.InputObjectConfigFieldMap{
						"shippingAddress": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
						"shippingCity":    &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
						"shippingPhone":   &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
						"paymentMethod":   &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
						"notes":           &graphql.InputObjectFieldConfig{Type: graphql.String},
					},
				}))},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				user, ok := p.Context.Value("user").(*User)
				if !ok {
					return nil, fmt.Errorf("user not authenticated")
				}

				input := p.Args["input"].(map[string]interface{})

				// Start transaction
				tx, err := DB.Begin()
				if err != nil {
					return nil, err
				}
				defer tx.Rollback()

				// Get cart items
				cartRows, err := tx.Query(`
					SELECT c.id, c.product_id, c.quantity, p.name, p.price
					FROM cart c
					JOIN products p ON c.product_id = p.id
					WHERE c.user_id = $1
				`, user.ID)
				if err != nil {
					return nil, err
				}
				defer cartRows.Close()

				var cartItems []struct {
					ID       int
					ProductID int
					Quantity int
					Name     string
					Price    float64
				}
				var totalAmount float64

				for cartRows.Next() {
					var item struct {
						ID       int
						ProductID int
						Quantity int
						Name     string
						Price    float64
					}
					err := cartRows.Scan(&item.ID, &item.ProductID, &item.Quantity, &item.Name, &item.Price)
					if err != nil {
						return nil, err
					}
					cartItems = append(cartItems, item)
					totalAmount += item.Price * float64(item.Quantity)
				}

				if len(cartItems) == 0 {
					return nil, fmt.Errorf("cart is empty")
				}

				// Generate order number
				orderNumber := fmt.Sprintf("ORD-%d-%s", time.Now().Year(), strconv.FormatInt(time.Now().Unix(), 10))

				// Create order
				var order Order
				err = tx.QueryRow(`
					INSERT INTO orders (user_id, order_number, status, total_amount, shipping_address, shipping_city, shipping_phone, payment_method, notes)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
					RETURNING id, user_id, order_number, status, total_amount, shipping_address, shipping_city, shipping_country, shipping_phone, payment_method, payment_status, notes, created_at, updated_at
				`, user.ID, orderNumber, "pending", totalAmount, input["shippingAddress"], input["shippingCity"], input["shippingPhone"], input["paymentMethod"], input["notes"]).Scan(
					&order.ID, &order.UserID, &order.OrderNumber, &order.Status, &order.TotalAmount, &order.ShippingAddress, &order.ShippingCity, &order.ShippingCountry, &order.ShippingPhone, &order.PaymentMethod, &order.PaymentStatus, &order.Notes, &order.CreatedAt, &order.UpdatedAt,
				)
				if err != nil {
					return nil, err
				}

				// Create order items
				for _, item := range cartItems {
					_, err = tx.Exec(`
						INSERT INTO order_items (order_id, product_id, product_name, product_price, quantity, total_price)
						VALUES ($1, $2, $3, $4, $5, $6)
					`, order.ID, item.ProductID, item.Name, item.Price, item.Quantity, item.Price*float64(item.Quantity))
					if err != nil {
						return nil, err
					}

					// Update product stock
					_, err = tx.Exec("UPDATE products SET stock_quantity = stock_quantity - $1 WHERE id = $2", item.Quantity, item.ProductID)
					if err != nil {
						return nil, err
					}
				}

				// Clear cart
				_, err = tx.Exec("DELETE FROM cart WHERE user_id = $1", user.ID)
				if err != nil {
					return nil, err
				}

				// Commit transaction
				if err = tx.Commit(); err != nil {
					return nil, err
				}

				return order, nil
			},
		},
		"createReview": &graphql.Field{
			Type: ReviewType,
			Args: graphql.FieldConfigArgument{
				"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.NewInputObject(graphql.InputObjectConfig{
					Name: "CreateReviewInput",
					Fields: graphql.InputObjectConfigFieldMap{
						"productId": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Int)},
						"rating":    &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Int)},
						"title":     &graphql.InputObjectFieldConfig{Type: graphql.String},
						"comment":   &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
					},
				}))},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				user, ok := p.Context.Value("user").(*User)
				if !ok {
					return nil, fmt.Errorf("user not authenticated")
				}

				input := p.Args["input"].(map[string]interface{})
				productID := input["productId"].(int)
				rating := input["rating"].(int)
				title := input["title"].(string)
				comment := input["comment"].(string)

				// Validate rating
				if rating < 1 || rating > 5 {
					return nil, fmt.Errorf("rating must be between 1 and 5")
				}

				// Check if user has purchased the product
				var hasPurchased bool
				err := DB.QueryRow(`
					SELECT EXISTS(
						SELECT 1 FROM order_items oi
						JOIN orders o ON oi.order_id = o.id
						WHERE o.user_id = $1 AND oi.product_id = $2 AND o.status = 'delivered'
					)
				`, user.ID, productID).Scan(&hasPurchased)
				if err != nil {
					return nil, err
				}

				var review Review
				query := `
					INSERT INTO reviews (user_id, product_id, rating, title, comment, is_verified_purchase)
					VALUES ($1, $2, $3, $4, $5, $6)
					RETURNING id, user_id, product_id, rating, title, comment, is_verified_purchase, created_at, updated_at
				`
				err = DB.QueryRow(query, user.ID, productID, rating, title, comment, hasPurchased).Scan(
					&review.ID, &review.UserID, &review.ProductID, &review.Rating, &review.Title, &review.Comment, &review.IsVerifiedPurchase, &review.CreatedAt, &review.UpdatedAt,
				)
				if err != nil {
					return nil, err
				}

				return review, nil
			},
		},
		"translateText": &graphql.Field{
			Type: graphql.String,
			Args: graphql.FieldConfigArgument{
				"text": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"from": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"to":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				text := p.Args["text"].(string)
				from := p.Args["from"].(string)
				to := p.Args["to"].(string)

				translatedText, err := handlers.TranslateText(text, from, to)
				if err != nil {
					return nil, fmt.Errorf("translation failed: %v", err)
				}

				return translatedText, nil
			},
		},
	},
})

// Schema creates the GraphQL schema
func Schema() (*graphql.Schema, error) {
	schemaConfig := graphql.SchemaConfig{
		Query:    RootQuery,
		Mutation: RootMutation,
	}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		return nil, err
	}
	return &schema, nil
} 

// Database functions moved from db package to avoid circular dependency

// GetUserByEmail retrieves a user by email
func GetUserByEmail(email string) (*User, error) {
	query := `
		SELECT id, email, first_name, last_name, phone, address, city, country, created_at, updated_at
		FROM users WHERE email = $1
	`
	user := &User{}
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
func GetUserByID(id int) (*User, error) {
	query := `
		SELECT id, email, first_name, last_name, phone, address, city, country, created_at, updated_at
		FROM users WHERE id = $1
	`
	user := &User{}
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
func CreateUser(email, password, firstName, lastName, phone, address, city string) (*User, error) {
	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	query := `
		INSERT INTO users (email, password_hash, first_name, last_name, phone, address, city)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, email, first_name, last_name, phone, address, city, country, created_at, updated_at
	`
	user := &User{}
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
func AuthenticateUser(email, password string) (*User, error) {
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