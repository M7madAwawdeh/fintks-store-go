# 🛒 Fintks Store - Modern E-commerce Platform

A comprehensive, AI-powered e-commerce platform built with Go, GraphQL, PostgreSQL, and OpenRouter AI integration. This project demonstrates how to build a full-featured online store with authentication, shopping cart, checkout process, reviews, and AI-powered features.

## 🚀 Features

### 🔐 Authentication & User Management
- **JWT-based Authentication** - Secure login/register with bcrypt password hashing
- **User Profiles** - Complete user management with address and contact information
- **Role-based Access** - Different user roles (admin, customer, user)

### 🛍️ Product Management
- **Product Catalog** - Comprehensive product database with categories
- **Search & Filtering** - Advanced search with category and price filtering
- **Product Images** - Multiple image support with primary image designation
- **Stock Management** - Real-time inventory tracking
- **Featured Products** - Highlight special products

### 🛒 Shopping Experience
- **Shopping Cart** - Full cart functionality with quantity management
- **Wishlist** - Save products for later
- **Product Reviews** - Customer reviews with ratings and likes
- **Product Likes** - Social features for product engagement

### 💳 Checkout & Orders
- **Cash on Delivery** - COD payment method (popular in Saudi Arabia)
- **Order Management** - Complete order lifecycle tracking
- **Order History** - User order history and status tracking
- **Shipping Information** - Address and delivery management

### 🤖 AI-Powered Features
- **Arabic Product Descriptions** - AI-generated marketing descriptions in Arabic
- **Multi-language Translation** - Real-time text translation between languages
- **Smart Product Recommendations** - AI-powered product suggestions

### 📱 Modern Tech Stack
- **GraphQL API** - Type-safe, efficient API with real-time updates
- **PostgreSQL Database** - Robust, scalable data storage
- **Docker Support** - Easy deployment and scaling
- **Cloud Ready** - Designed for Google Cloud Platform deployment

## 🛠️ Technologies

- **Backend**: Go 1.22
- **API**: GraphQL (github.com/graphql-go/graphql)
- **Database**: PostgreSQL 15
- **Authentication**: JWT with bcrypt
- **AI Services**: OpenRouter API
- **Containerization**: Docker & Docker Compose
- **Cloud Platform**: Google Cloud Platform (Cloud Run, Cloud Build)

## 📋 Prerequisites

- Go 1.22 or later
- Docker and Docker Compose
- PostgreSQL (if running locally)
- OpenRouter API key (for AI features)

## 🏃‍♂️ Quick Start

### Option 1: Using Docker Compose (Recommended)

1. **Clone the repository**
   ```bash
   git clone https://github.com/M7madAwawdeh/fintks-store-go.git
   cd fintks-store
   ```

2. **Set up environment variables**
   ```bash
   # Copy the .env file
   cp .env.example .env
   
   # Edit the .env file with your settings
   OPENROUTER_API_KEY="your_openrouter_api_key_here"
   ```

3. **Start the services**
   ```bash
   docker-compose up --build
   ```

4. **Access the application**
   - Homepage: http://localhost:8080
   - GraphQL endpoint: http://localhost:8080/graphql

### Option 2: Local Development

1. **Install dependencies**
   ```bash
   go mod download
   ```

2. **Set up PostgreSQL**
   ```bash
   # Create database
   createdb fintks_store
   
   # Run initialization script
   psql -d fintks_store -f init.sql
   ```

3. **Set environment variables**
   ```bash
   export DATABASE_URL="postgres://username:password@localhost:5432/fintks_store?sslmode=disable"
   export OPENROUTER_API_KEY="your_openrouter_api_key_here"
   ```

4. **Run the application**
   ```bash
   go run main.go
   ```

## 📊 API Usage

### 🔐 Authentication

**Register a new user:**
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

**Login:**
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

### 🛍️ Product Queries

**Get all products:**
```graphql
{
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
  }
}
```

**Get featured products:**
```graphql
{
  products(isFeatured: true) {
    id
    name
    price
    imageUrl
  }
}
```

**Search products:**
```graphql
{
  products(search: "iPhone") {
    id
    name
    price
    imageUrl
  }
}
```

**Get products by category:**
```graphql
{
  products(categoryId: 1) {
    id
    name
    price
    imageUrl
  }
}
```

### 🛒 Shopping Cart

**Add to cart:**
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

**View cart:**
```graphql
{
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

### 🤖 AI Features

**Translate text:**
```graphql
mutation {
  translateText(
    text: "Hello, welcome to our store!"
    from: "English"
    to: "Arabic"
  )
}
```

## 🧪 Demo Credentials

The application comes with pre-configured demo accounts:

- **Admin**: admin@fintks.com / password123
- **Customer**: customer@fintks.com / password123  
- **User**: user@fintks.com / password123

## 🐳 Docker Commands

```bash
# Build the image
docker build -t fintks-store .

# Run with Docker Compose
docker-compose up

# Run in background
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Rebuild and restart
docker-compose up --build
```

## ☁️ Google Cloud Platform Deployment

### Prerequisites
- Google Cloud SDK installed and authenticated
- Cloud Run, Cloud Build, and Artifact Registry APIs enabled

### Deployment Steps

1. **Set your project ID**
   ```bash
   export PROJECT_ID="your-gcp-project-id"
   ```

2. **Build and push the image**
   ```bash
   gcloud builds submit --tag gcr.io/$PROJECT_ID/fintks-store
   ```

3. **Deploy to Cloud Run**
   ```bash
   gcloud run deploy fintks-store \
     --image gcr.io/$PROJECT_ID/fintks-store \
     --platform managed \
     --port 8080 \
     --allow-unauthenticated \
     --set-env-vars="OPENROUTER_API_KEY=$OPENROUTER_API_KEY" \
     --set-env-vars="DATABASE_URL=$DATABASE_URL"
   ```

## 🔧 Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `postgres://postgres:password@localhost:5432/fintks_store?sslmode=disable` |
| `OPENROUTER_API_KEY` | OpenRouter API key for AI services | Required |
| `PORT` | Server port | `8080` |
| `JWT_SECRET` | JWT signing secret | `your-secret-key` |

## 📁 Project Structure

```
fintks-store/
├── main.go              # HTTP server with authentication middleware
├── go.mod               # Go module dependencies
├── .env                 # Environment variables
├── graph/
│   ├── model.go         # All data models (User, Product, Order, etc.)
│   ├── schema.graphqls  # GraphQL schema definition
│   └── resolver.go      # GraphQL resolvers
├── auth/
│   └── auth.go         # JWT authentication utilities
├── db/
│   └── db.go           # Database connection and user management
├── handlers/
│   ├── ai.go           # AI description generation
│   └── lang.go         # Translation services
├── Dockerfile          # Multi-stage Docker build
├── docker-compose.yml  # Multi-service orchestration
├── init.sql           # Database initialization with demo data
└── README.md          # Project documentation
```

## 🎯 Learning Goals

This project demonstrates:

1. **Full E-commerce Implementation** - Complete online store with all essential features
2. **Modern Go Development** - Using Go modules, proper package structure, and best practices
3. **GraphQL API Design** - Building a comprehensive, type-safe API
4. **Database Design** - Complex relational database with proper relationships
5. **Authentication & Security** - JWT-based auth with password hashing
6. **AI Service Integration** - External API integration for AI features
7. **Containerization** - Docker multi-stage builds and service orchestration
8. **Cloud Deployment** - Deploying to Google Cloud Platform
9. **Payment Integration** - Cash on Delivery implementation
10. **User Experience** - Shopping cart, wishlist, reviews, and social features

## 🛡️ Security Features

- **Password Hashing** - bcrypt for secure password storage
- **JWT Authentication** - Stateless authentication with token expiration
- **Input Validation** - GraphQL schema validation
- **SQL Injection Prevention** - Parameterized queries
- **CORS Configuration** - Proper cross-origin resource sharing

## 🚀 Performance Features

- **Database Indexing** - Optimized queries with proper indexes
- **Connection Pooling** - Efficient database connection management
- **GraphQL Efficiency** - Only fetch requested data
- **Docker Optimization** - Multi-stage builds for smaller images

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

If you encounter any issues or have questions:

1. Check the [Issues](../../issues) page
2. Review the documentation above
3. Create a new issue with detailed information

## 🌟 Features Roadmap

- [ ] Payment gateway integration (Stripe, PayPal)
- [ ] Email notifications
- [ ] SMS notifications
- [ ] Admin dashboard
- [ ] Analytics and reporting
- [ ] Mobile app API
- [ ] Multi-language support
- [ ] Advanced search with Elasticsearch
- [ ] Real-time chat support
- [ ] Loyalty program
- [ ] Coupon system
- [ ] Inventory alerts
- [ ] Order tracking
- [ ] Return/refund management

---


**Happy shopping! 🛒✨** 
