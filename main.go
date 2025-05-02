package main

import (
	"e-commerce_with_golang/config"
	"e-commerce_with_golang/database"
	"e-commerce_with_golang/models"
	"e-commerce_with_golang/repositories"
	"e-commerce_with_golang/routes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("development-secret-key-123")

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func main() {
	cfg := config.LoadConfig()
	database.ConnectDB(cfg)

	// Run migrations
	if err := database.Migrate(database.DB); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Initialize repositories
	productRepo := repositories.NewProductRepository(database.DB)
	reviewRepo := repositories.NewReviewRepository(database.DB)
	userRepo := repositories.NewUserRepository(database.DB)

	// Setup router with recovery middleware
	router := gin.Default()

	// Register route
	router.POST("/register", func(c *gin.Context) {
		var req RegisterRequest

		// Bind and validate request
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("Register request validation failed: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request format",
				"details": err.Error(),
			})
			return
		}

		// Check if user already exists
		existingUser, err := userRepo.GetUserByEmail(req.Email)
		if err == nil && existingUser != nil {
			c.JSON(http.StatusConflict, gin.H{
				"error": "User already exists",
			})
			return
		}

		// Create new user
		newUser := &models.User{
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
		}
		if err := userRepo.CreateUser(newUser); err != nil {
			log.Printf("User creation failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create user",
			})
			return
		}

		// Generate JWT token with secure claims
		expirationTime := time.Now().Add(24 * time.Hour)
		claims := jwt.MapClaims{
			"user_id": newUser.ID,
			"email":   newUser.Email,
			"exp":     expirationTime.Unix(),
			"iat":     time.Now().Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			log.Printf("Token generation failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate authentication token",
			})
			return
		}

		c.SetSameSite(http.SameSiteStrictMode)
		c.SetCookie(
			"auth_token",
			tokenString,
			int(expirationTime.Sub(time.Now()).Seconds()),
			"/",
			"",    // domain
			false, // secure (set to true in production with HTTPS)
			true,  // httpOnly
		)

		// Successful response
		c.JSON(http.StatusOK, gin.H{
			"token":      tokenString,
			"expires_at": expirationTime.Format(time.RFC3339),
			"user_id":    newUser.ID,
			"message":    "Registration successful",
		})
	})

	router.POST("/login", func(c *gin.Context) {
		var req models.LoginRequest

		// Bind and validate request
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("Login request validation failed: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request format",
				"details": err.Error(),
			})
			return
		}

		log.Printf("Login attempt for email: %s", req.Email)

		// Get user by email
		user, err := userRepo.GetUserByEmail(req.Email)
		if err != nil {
			log.Printf("Database error during user lookup: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid credentials",
				"message": "Authentication failed",
			})
			return
		}

		if user == nil || user.ID == 0 {
			log.Printf("No user found for email: %s", req.Email)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid credentials",
				"message": "Authentication failed",
			})
			return
		}

		log.Printf("Found user ID: %d, Password hash length: %d",
			user.ID, len(user.Password))

		start := time.Now()
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
		elapsed := time.Since(start)
		log.Printf("Password comparison took %v", elapsed)

		if err != nil {
			log.Printf("Password mismatch for user %d: %v", user.ID, err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid credentials",
				"message": "Authentication failed",
			})
			return
		}

		expirationTime := time.Now().Add(24 * time.Hour)
		claims := jwt.MapClaims{
			"user_id": user.ID,
			"email":   user.Email,
			"exp":     expirationTime.Unix(),
			"iat":     time.Now().Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			log.Printf("Token generation failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate authentication token",
			})
			return
		}

		c.SetSameSite(http.SameSiteStrictMode)
		c.SetCookie(
			"auth_token",
			tokenString,
			int(expirationTime.Sub(time.Now()).Seconds()),
			"/",
			"",   
			false, 
			true,  
		)
		c.JSON(http.StatusOK, gin.H{
			"token":      tokenString,
			"expires_at": expirationTime.Format(time.RFC3339),
			"user_id":    user.ID,
			"message":    "Login successful",
		})
	})

	routes.SetupRoutes(router, productRepo, reviewRepo, authCheck())

	port := "8080"
	log.Printf("Server running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func authCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to get token from header first
		authHeader := c.GetHeader("Authorization")
		tokenString := ""

		if authHeader != "" {
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				tokenString = authHeader[7:]
			}
		} else {
			if cookie, err := c.Cookie("auth_token"); err == nil {
				tokenString = cookie
			}
		}

		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization token required",
			})
			return
		}

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtKey, nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid token",
				"details": err.Error(),
			})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("userID", claims["user_id"])
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token claims",
			})
		}
	}
}
