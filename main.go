package main

import (
	"ai-catalog/auth"
	"ai-catalog/graph"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/graphql-go/graphql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type GraphQLRequest struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

// AuthMiddleware extracts JWT token and adds user to context
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			token, err := auth.ExtractTokenFromHeader(authHeader)
			if err == nil {
				user, err := auth.GetUserFromToken(token)
				if err == nil {
					// Convert auth.User to graph.User for context
					graphUser := &graph.User{
						ID:    user.ID,
						Email: user.Email,
					}
					// Add user to context
					ctx := context.WithValue(r.Context(), "user", graphUser)
					r = r.WithContext(ctx)
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}

// Connect establishes a connection to PostgreSQL database
func Connect() (*sql.DB, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:password@localhost:5432/ai_catalog?sslmode=disable"
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	// Create all tables
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %v", err)
	}

	fmt.Println("Database connected successfully")
	return db, nil
}

// createTables creates all necessary tables by executing init.sql
func createTables(db *sql.DB) error {
	initSQL, err := os.ReadFile("init.sql")
	if err != nil {
		return fmt.Errorf("failed to read init.sql: %v", err)
	}
	_, err = db.Exec(string(initSQL))
	if err != nil {
		return fmt.Errorf("failed to execute init.sql: %v", err)
	}
	return nil
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Connect to database
	db, err := Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Set database connection in graph package
	graph.SetDB(db)

	// Create GraphQL schema
	schema, err := graph.Schema()
	if err != nil {
		log.Fatal("Failed to create GraphQL schema:", err)
	}

	// Create router
	router := mux.NewRouter()

	// Apply authentication middleware
	router.Use(AuthMiddleware)

	// Serve static files
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Frontend interface
	router.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	// GraphQL handler
	router.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req GraphQLRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		params := graphql.Params{
			Schema:         *schema,
			RequestString:  req.Query,
			OperationName:  req.OperationName,
			VariableValues: req.Variables,
			Context:        r.Context(),
		}

		result := graphql.Do(params)
		if len(result.Errors) > 0 {
			log.Printf("GraphQL errors: %v", result.Errors)
		}

		json.NewEncoder(w).Encode(result)
	})

	// Homepage handler
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		html := `
<!DOCTYPE html>
<html>
<head>
    <title>🛒 Fintks Store - E-commerce Platform</title>
    <style>
        body { 
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; 
            margin: 0; 
            padding: 0; 
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
        }
        .container { 
            max-width: 1200px; 
            margin: 0 auto; 
            padding: 40px 20px;
        }
        .header {
            text-align: center;
            color: white;
            margin-bottom: 50px;
        }
        .header h1 {
            font-size: 3rem;
            margin-bottom: 10px;
            text-shadow: 2px 2px 4px rgba(0,0,0,0.3);
        }
        .header p {
            font-size: 1.2rem;
            opacity: 0.9;
        }
        .features {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 30px;
            margin-bottom: 50px;
        }
        .feature {
            background: white;
            padding: 30px;
            border-radius: 15px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.1);
            transition: transform 0.3s ease;
        }
        .feature:hover {
            transform: translateY(-5px);
        }
        .feature h3 {
            color: #333;
            margin-bottom: 15px;
            font-size: 1.5rem;
        }
        .feature p {
            color: #666;
            line-height: 1.6;
        }
        .endpoint { 
            background: #f8f9fa; 
            padding: 25px; 
            margin: 20px 0; 
            border-radius: 10px;
            border-left: 4px solid #667eea;
        }
        .endpoint h2 {
            color: #333;
            margin-top: 0;
        }
        code { 
            background: #e9ecef; 
            padding: 3px 6px; 
            border-radius: 4px;
            font-family: 'Courier New', monospace;
        }
        .demo-credentials {
            background: #fff3cd;
            border: 1px solid #ffeaa7;
            padding: 20px;
            border-radius: 10px;
            margin: 30px 0;
        }
        .demo-credentials h3 {
            color: #856404;
            margin-top: 0;
        }
        .demo-credentials ul {
            color: #856404;
            margin: 10px 0;
        }
        .btn {
            display: inline-block;
            background: #667eea;
            color: white;
            padding: 12px 24px;
            text-decoration: none;
            border-radius: 6px;
            margin: 10px 5px;
            transition: background 0.3s ease;
        }
        .btn:hover {
            background: #5a6fd8;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🛒 Fintks Store</h1>
            <p>Modern E-commerce Platform with AI-Powered Features</p>
        </div>

        <div class="features">
            <div class="feature">
                <h3>🔐 User Authentication</h3>
                <p>Secure JWT-based authentication with user registration, login, and profile management.</p>
            </div>
            <div class="feature">
                <h3>🛍️ Product Catalog</h3>
                <p>Comprehensive product management with categories, search, filtering, and AI-generated descriptions.</p>
            </div>
            <div class="feature">
                <h3>🛒 Shopping Cart</h3>
                <p>Full-featured shopping cart with quantity management and real-time updates.</p>
            </div>
            <div class="feature">
                <h3>💳 Checkout Process</h3>
                <p>Complete checkout flow with cash on delivery payment method and order tracking.</p>
            </div>
            <div class="feature">
                <h3>⭐ Reviews & Ratings</h3>
                <p>Customer reviews with ratings, likes, and verified purchase badges.</p>
            </div>
            <div class="feature">
                <h3>🤖 AI Integration</h3>
                <p>AI-powered product descriptions in Arabic and multi-language translation services.</p>
            </div>
        </div>

        <div class="demo-credentials">
            <h3>🧪 Demo Credentials</h3>
            <ul>
                <li><strong>Admin:</strong> admin@fintks.com / password123</li>
                <li><strong>Customer:</strong> customer@fintks.com / password123</li>
                <li><strong>User:</strong> user@fintks.com / password123</li>
            </ul>
        </div>

        <div class="endpoint">
            <h2>📊 GraphQL Endpoint</h2>
            <p><code>POST /graphql</code></p>
            <p>Use this endpoint to interact with the GraphQL API.</p>
            <a href="/graphql" class="btn" target="_blank">Try GraphQL Playground</a>
        </div>

        <div class="endpoint">
            <h2>🖥️ Interactive Frontend</h2>
            <p>Test the e-commerce features with our interactive web interface.</p>
            <a href="/app" class="btn">Launch E-commerce App</a>
        </div>

        <div class="endpoint">
            <h2>🔐 Authentication Examples</h2>
            <h4>Register User:</h4>
            <pre><code>mutation {
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
}</code></pre>

            <h4>Login:</h4>
            <pre><code>mutation {
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
}</code></pre>
        </div>

        <div class="endpoint">
            <h2>🛍️ Product Queries</h2>
            <h4>Get All Products:</h4>
            <pre><code>{
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
}</code></pre>

            <h4>Get Featured Products:</h4>
            <pre><code>{
  products(isFeatured: true) {
    id
    name
    price
    imageUrl
  }
}</code></pre>

            <h4>Search Products:</h4>
            <pre><code>{
  products(search: "iPhone") {
    id
    name
    price
    imageUrl
  }
}</code></pre>
        </div>

        <div class="endpoint">
            <h2>🤖 AI Features</h2>
            <h4>Translate Text:</h4>
            <pre><code>mutation {
  translateText(
    text: "Hello, welcome to our store!"
    from: "English"
    to: "Arabic"
  )
}</code></pre>
        </div>

        <div class="endpoint">
            <h2>📱 Technologies Used</h2>
            <ul>
                <li><strong>Backend:</strong> Go 1.22 with GraphQL</li>
                <li><strong>Database:</strong> PostgreSQL with comprehensive schema</li>
                <li><strong>Authentication:</strong> JWT with bcrypt password hashing</li>
                <li><strong>AI Services:</strong> OpenRouter API integration</li>
                <li><strong>Payment:</strong> Cash on Delivery (COD)</li>
                <li><strong>Deployment:</strong> Docker & Google Cloud Platform ready</li>
            </ul>
        </div>
    </div>
</body>
</html>`
		fmt.Fprint(w, html)
	})

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Fintks Store starting on port %s", port)
	log.Printf("📊 GraphQL endpoint: http://localhost:%s/graphql", port)
	log.Printf("🖥️ Frontend app: http://localhost:%s/app", port)
	log.Printf("🏠 Homepage: http://localhost:%s", port)
	
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal("Failed to start server:", err)
	}
} 