package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	// Replace with your actual module path:
	"github.com/Gopikrishnaj96/ai-task-manager/db"
)

// A global secret key for signing JWT tokens (don't expose in production!)
var secretKey = []byte("supersecretkey")

// User model matches the 'users' table in PostgreSQL
type User struct {
	ID       int    `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
}

// Task model matches the 'tasks' table in PostgreSQL
type Task struct {
	ID          int    `db:"id" json:"id"`
	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`
	Status      string `db:"status" json:"status"`
	UserID      int    `db:"user_id" json:"user_id"`
}

func main() {
	// Initialize the DB connection
	db.InitDB()

	// Create a Fiber app
	app := fiber.New()

	// Public routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Fiber!")
	})
	app.Post("/signup", signupHandler)
	app.Post("/login", loginHandler)

	// Protected routes (require JWT)
	app.Use(jwtMiddleware)

	app.Post("/tasks", createTaskHandler)
	app.Get("/tasks", getTasksHandler)

	fmt.Println("Starting server on http://localhost:4000")
	if err := app.Listen(":4000"); err != nil {
		log.Fatal(err)
	}
}

// ------------------------------------------------------
// 1. SIGNUP HANDLER
// ------------------------------------------------------
func signupHandler(c *fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Parse JSON request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Hash the password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not hash password"})
	}

	// Insert into the 'users' table
	_, err = db.DB.Exec(`INSERT INTO users (username, password) VALUES ($1, $2)`, req.Username, string(hashed))
	if err != nil {
		// Possibly a duplicate username
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Username already exists"})
	}

	return c.JSON(fiber.Map{"message": "User created successfully"})
}

// ------------------------------------------------------
// 2. LOGIN HANDLER
// ------------------------------------------------------
func loginHandler(c *fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Parse JSON request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Fetch user from the 'users' table
	var user User
	err := db.DB.Get(&user, `SELECT * FROM users WHERE username=$1`, req.Username)
	if err != nil {
		// User not found
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Compare hashed password
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // 24-hour expiration
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
	}

	// Return the token to the client
	return c.JSON(fiber.Map{"token": tokenString})
}

// ------------------------------------------------------
// 3. JWT MIDDLEWARE (Protecting Routes)
// ------------------------------------------------------
func jwtMiddleware(c *fiber.Ctx) error {
	// Look for token in the "Authorization" header
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
	}

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	// Optionally, store claims in c.Locals so handlers can access user info
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Could not parse claims"})
	}

	c.Locals("user_id", claims["user_id"])
	c.Locals("username", claims["username"])

	return c.Next() // proceed to the next handler
}

// ------------------------------------------------------
// 4. TASK HANDLERS
// ------------------------------------------------------
func createTaskHandler(c *fiber.Ctx) error {
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Extract user_id from token claims (stored as float64 in JWT)
	userIDVal := c.Locals("user_id")
	if userIDVal == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "No user ID in token"})
	}
	userID := int(userIDVal.(float64))

	// Insert the task into DB
	_, err := db.DB.Exec(`
        INSERT INTO tasks (title, description, status, user_id)
        VALUES ($1, $2, 'pending', $3)`,
		req.Title, req.Description, userID,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create task"})
	}

	return c.JSON(fiber.Map{"message": "Task created successfully"})
}

func getTasksHandler(c *fiber.Ctx) error {
	// Extract user_id from token claims
	userIDVal := c.Locals("user_id")
	if userIDVal == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "No user ID in token"})
	}
	userID := int(userIDVal.(float64))

	// Fetch tasks for this user
	tasks := []Task{}
	err := db.DB.Select(&tasks, `SELECT * FROM tasks WHERE user_id=$1`, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not fetch tasks"})
	}

	return c.JSON(tasks)
}
