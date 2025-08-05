# 🛒 Fintks Store API Documentation

## Overview

Fintks Store is a comprehensive e-commerce platform built with Go, GraphQL, PostgreSQL, and AI integration. This document provides detailed information about all available API endpoints, data models, and usage examples.

## Base URL

```
http://localhost:8080
```

## Authentication

The API uses JWT (JSON Web Tokens) for authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## GraphQL Endpoint

```
POST /graphql
```

## Data Models

### User
```graphql
type User {
  id: Int!
  email: String!
  firstName: String!
  lastName: String!
  phone: String
  address: String
  city: String
  country: String
  createdAt: String!
  updatedAt: String!
}
```

### Product
```graphql
type Product {
  id: Int!
  name: String!
  price: Float!
  originalPrice: Float
  categoryId: Int!
  category: Category
  description: String!
  shortDescription: String
  imageUrl: String
  stockQuantity: Int!
  sku: String
  weight: Float
  dimensions: String
  isActive: Boolean!
  isFeatured: Boolean!
  createdAt: String!
  updatedAt: String!
  averageRating: Float
  reviewCount: Int
  isInWishlist: Boolean
  isLiked: Boolean
}
```

### Category
```graphql
type Category {
  id: Int!
  name: String!
  description: String
  imageUrl: String
  parentId: Int
  createdAt: String!
}
```

### CartItem
```graphql
type CartItem {
  id: Int!
  userId: Int!
  productId: Int!
  product: Product
  quantity: Int!
  createdAt: String!
  updatedAt: String!
}
```

### CartSummary
```graphql
type CartSummary {
  items: [CartItem!]!
  totalItems: Int!
  totalPrice: Float!
}
```

### Order
```graphql
type Order {
  id: Int!
  userId: Int!
  user: User
  orderNumber: String!
  status: String!
  totalAmount: Float!
  shippingAddress: String!
  shippingCity: String!
  shippingCountry: String!
  shippingPhone: String!
  paymentMethod: String!
  paymentStatus: String!
  notes: String
  createdAt: String!
  updatedAt: String!
  items: [OrderItem!]!
}
```

### Review
```graphql
type Review {
  id: Int!
  userId: Int!
  user: User
  productId: Int!
  rating: Int!
  title: String
  comment: String
  isVerifiedPurchase: Boolean!
  createdAt: String!
  updatedAt: String!
  likeCount: Int
  isLiked: Boolean
}
```

## Queries

### Authentication

#### Get Current User
```graphql
query {
  me {
    id
    email
    firstName
    lastName
    phone
    address
    city
    country
  }
}
```

### Categories

#### Get All Categories
```graphql
query {
  categories {
    id
    name
    description
    imageUrl
    parentId
    createdAt
  }
}
```

### Products

#### Get All Products
```graphql
query {
  products {
    id
    name
    price
    originalPrice
    imageUrl
    shortDescription
    category {
      name
    }
    stockQuantity
    isFeatured
  }
}
```

#### Get Products with Filters
```graphql
query {
  products(
    categoryId: 1
    search: "iPhone"
    minPrice: 100
    maxPrice: 1000
    isFeatured: true
    limit: 10
  ) {
    id
    name
    price
    imageUrl
    category {
      name
    }
  }
}
```

#### Get Featured Products
```graphql
query {
  featuredProducts {
    id
    name
    price
    imageUrl
    shortDescription
  }
}
```

#### Get Single Product
```graphql
query {
  product(id: 1) {
    id
    name
    price
    originalPrice
    description
    shortDescription
    imageUrl
    stockQuantity
    category {
      name
    }
  }
}
```

### Shopping Cart

#### Get Cart (Requires Authentication)
```graphql
query {
  cart {
    items {
      id
      quantity
      product {
        id
        name
        price
        imageUrl
      }
    }
    totalItems
    totalPrice
  }
}
```

### Wishlist

#### Get Wishlist (Requires Authentication)
```graphql
query {
  wishlist {
    id
    productId
    product {
      id
      name
      price
      imageUrl
    }
    createdAt
  }
}
```

### Orders

#### Get User Orders (Requires Authentication)
```graphql
query {
  orders {
    id
    orderNumber
    status
    totalAmount
    shippingAddress
    shippingCity
    paymentMethod
    paymentStatus
    createdAt
    items {
      productName
      quantity
      totalPrice
    }
  }
}
```

### Reviews

#### Get Product Reviews
```graphql
query {
  productReviews(productId: 1) {
    id
    rating
    title
    comment
    isVerifiedPurchase
    createdAt
    user {
      firstName
      lastName
    }
  }
}
```

## Mutations

### Authentication

#### Register User
```graphql
mutation {
  register(input: {
    email: "newuser@example.com"
    password: "password123"
    firstName: "John"
    lastName: "Doe"
    phone: "+966501234567"
    address: "King Fahd Road"
    city: "Riyadh"
  }) {
    user {
      id
      email
      firstName
      lastName
    }
    token
  }
}
```

#### Login User
```graphql
mutation {
  login(input: {
    email: "customer@fintks.com"
    password: "password123"
  }) {
    user {
      id
      email
      firstName
      lastName
    }
    token
  }
}
```

### Shopping Cart

#### Add to Cart (Requires Authentication)
```graphql
mutation {
  addToCart(input: {
    productId: 1
    quantity: 2
  }) {
    id
    quantity
    product {
      name
      price
    }
  }
}
```

#### Remove from Cart (Requires Authentication)
```graphql
mutation {
  removeFromCart(id: 1)
}
```

### Wishlist

#### Add to Wishlist (Requires Authentication)
```graphql
mutation {
  addToWishlist(productId: 1) {
    id
    productId
    createdAt
  }
}
```

#### Remove from Wishlist (Requires Authentication)
```graphql
mutation {
  removeFromWishlist(productId: 1)
}
```

### Orders

#### Create Order (Requires Authentication)
```graphql
mutation {
  createOrder(input: {
    shippingAddress: "King Fahd Road, Riyadh"
    shippingCity: "Riyadh"
    shippingPhone: "+966501234567"
    paymentMethod: "cash_on_delivery"
    notes: "Please deliver in the morning"
  }) {
    id
    orderNumber
    status
    totalAmount
    shippingAddress
    paymentMethod
    createdAt
  }
}
```

### Reviews

#### Create Review (Requires Authentication)
```graphql
mutation {
  createReview(input: {
    productId: 1
    rating: 5
    title: "Excellent Product!"
    comment: "This product exceeded my expectations. Highly recommended!"
  }) {
    id
    rating
    title
    comment
    isVerifiedPurchase
    createdAt
  }
}
```

### AI Features

#### Translate Text
```graphql
mutation {
  translateText(
    text: "Hello, welcome to our store!"
    from: "English"
    to: "Arabic"
  )
}
```

## Error Handling

The API returns errors in the following format:

```json
{
  "errors": [
    {
      "message": "Error description",
      "locations": [
        {
          "line": 2,
          "column": 3
        }
      ],
      "path": ["fieldName"]
    }
  ],
  "data": null
}
```

### Common Error Messages

- `"user not authenticated"` - User is not logged in
- `"product not found"` - Product with specified ID doesn't exist
- `"insufficient stock"` - Not enough stock available
- `"cart is empty"` - Cannot create order with empty cart
- `"rating must be between 1 and 5"` - Invalid rating value

## Status Codes

- `200` - Success
- `400` - Bad Request (invalid input)
- `401` - Unauthorized (authentication required)
- `403` - Forbidden (insufficient permissions)
- `404` - Not Found
- `500` - Internal Server Error

## Demo Credentials

For testing purposes, you can use these pre-configured accounts:

### Admin Account
- **Email:** admin@fintks.com
- **Password:** password123

### Customer Account
- **Email:** customer@fintks.com
- **Password:** password123

### User Account
- **Email:** user@fintks.com
- **Password:** password123

## Order Status Values

- `pending` - Order created, waiting for processing
- `processing` - Order is being prepared
- `shipped` - Order has been shipped
- `delivered` - Order has been delivered
- `cancelled` - Order has been cancelled

## Payment Methods

- `cash_on_delivery` - Cash on Delivery (COD)

## Payment Status Values

- `pending` - Payment pending
- `paid` - Payment received
- `failed` - Payment failed
- `refunded` - Payment refunded

## Rate Limiting

Currently, there are no rate limits implemented. In production, consider implementing rate limiting to prevent abuse.

## Security Considerations

1. **JWT Tokens**: Tokens expire after 24 hours
2. **Password Hashing**: Passwords are hashed using bcrypt
3. **Input Validation**: All inputs are validated through GraphQL schema
4. **SQL Injection Prevention**: All queries use parameterized statements
5. **CORS**: Configured to allow cross-origin requests

## Testing

### Using cURL

```bash
# Login
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{
    "query": "mutation { login(input: { email: \"customer@fintks.com\", password: \"password123\" }) { user { id email firstName lastName } token } }"
  }'

# Get products (with authentication)
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "query": "{ products { id name price imageUrl } }"
  }'
```

### Using GraphQL Playground

1. Open your browser and go to `http://localhost:8080/graphql`
2. Use the interactive GraphQL playground to test queries and mutations
3. Set the Authorization header in the HTTP Headers section:
   ```json
   {
     "Authorization": "Bearer YOUR_TOKEN_HERE"
   }
   ```

## Frontend Integration

The platform includes a complete frontend interface available at `/app`. This provides a user-friendly way to test all e-commerce features including:

- User registration and login
- Product browsing and search
- Shopping cart management
- Wishlist functionality
- Order creation and tracking

## Deployment

### Local Development
```bash
# Start with Docker Compose
docker-compose up --build

# Or run locally
go run main.go
```

### Production Deployment
```bash
# Build Docker image
docker build -t fintks-store .

# Run with environment variables
docker run -p 8080:8080 \
  -e DATABASE_URL="your-database-url" \
  -e OPENROUTER_API_KEY="your-api-key" \
  fintks-store
```

## Support

For issues and questions:
1. Check the [README.md](README.md) for setup instructions
2. Review the error messages in the API responses
3. Ensure all required environment variables are set
4. Verify database connectivity and schema

---

**Happy coding! 🚀** 