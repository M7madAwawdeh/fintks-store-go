package graph

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	Phone        string    `json:"phone"`
	Address      string    `json:"address"`
	City         string    `json:"city"`
	Country      string    `json:"country"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// Category represents a product category
type Category struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ImageURL    string    `json:"imageUrl"`
	ParentID    *int      `json:"parentId"`
	CreatedAt   time.Time `json:"createdAt"`
}

// Product represents a product in the catalog
type Product struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	Price             float64   `json:"price"`
	OriginalPrice     *float64  `json:"originalPrice"`
	CategoryID        int       `json:"categoryId"`
	Category          *Category `json:"category"`
	Description       string    `json:"description"`
	ShortDescription  string    `json:"shortDescription"`
	ImageURL          string    `json:"imageUrl"`
	StockQuantity     int       `json:"stockQuantity"`
	SKU               string    `json:"sku"`
	Weight            *float64  `json:"weight"`
	Dimensions        string    `json:"dimensions"`
	IsActive          bool      `json:"isActive"`
	IsFeatured        bool      `json:"isFeatured"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
	AverageRating     float64   `json:"averageRating"`
	ReviewCount       int       `json:"reviewCount"`
	IsInWishlist      bool      `json:"isInWishlist"`
	IsLiked           bool      `json:"isLiked"`
}

// CartItem represents an item in the shopping cart
type CartItem struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	ProductID int       `json:"productId"`
	Product   *Product  `json:"product"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// WishlistItem represents an item in the wishlist
type WishlistItem struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	ProductID int       `json:"productId"`
	Product   *Product  `json:"product"`
	CreatedAt time.Time `json:"createdAt"`
}

// Order represents a customer order
type Order struct {
	ID             int       `json:"id"`
	UserID         int       `json:"userId"`
	User           *User     `json:"user"`
	OrderNumber    string    `json:"orderNumber"`
	Status         string    `json:"status"`
	TotalAmount    float64   `json:"totalAmount"`
	ShippingAddress string   `json:"shippingAddress"`
	ShippingCity   string    `json:"shippingCity"`
	ShippingCountry string   `json:"shippingCountry"`
	ShippingPhone  string    `json:"shippingPhone"`
	PaymentMethod  string    `json:"paymentMethod"`
	PaymentStatus  string    `json:"paymentStatus"`
	Notes          string    `json:"notes"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Items          []*OrderItem `json:"items"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID           int       `json:"id"`
	OrderID      int       `json:"orderId"`
	ProductID    int       `json:"productId"`
	ProductName  string    `json:"productName"`
	ProductPrice float64   `json:"productPrice"`
	Quantity     int       `json:"quantity"`
	TotalPrice   float64   `json:"totalPrice"`
	CreatedAt    time.Time `json:"createdAt"`
}

// Review represents a product review
type Review struct {
	ID                int       `json:"id"`
	UserID            int       `json:"userId"`
	User              *User     `json:"user"`
	ProductID         int       `json:"productId"`
	Rating            int       `json:"rating"`
	Title             string    `json:"title"`
	Comment           string    `json:"comment"`
	IsVerifiedPurchase bool     `json:"isVerifiedPurchase"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
	LikeCount         int       `json:"likeCount"`
	IsLiked           bool      `json:"isLiked"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

// CartSummary represents cart summary information
type CartSummary struct {
	Items      []*CartItem `json:"items"`
	TotalItems int         `json:"totalItems"`
	TotalPrice float64     `json:"totalPrice"`
}

// SearchResult represents search results
type SearchResult struct {
	Products   []*Product `json:"products"`
	Categories []*Category `json:"categories"`
	Total      int        `json:"total"`
}

// Pagination represents pagination information
type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
} 