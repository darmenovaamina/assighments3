package main

import (
	"github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
    "net/http"
    "strconv"
)

var jwtSecretKey = []byte("i_like_cats")

type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Password string `json:"password"`
    Role     string `json:"role"`
}

type Item struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Price float64 `json:"price"`
}

var items = []Item{}
var nextID = 1

var users = []User{
    {ID: 1, Username: "admin", Password: "admin123", Role: "admin"},
}

// Generate JWT token
func generateToken(username, role string) (string, error) {
    claims := jwt.MapClaims{
        "username": username,
        "role":     role,
        "exp":      time.Now().Add(time.Hour * 2).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecretKey)
}

// Registration endpoint
func register(c *gin.Context) {
    var newUser User
    if err := c.ShouldBindJSON(&newUser); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    users = append(users, newUser)
    c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

// Login endpoint
func login(c *gin.Context) {
    var credentials User
    if err := c.ShouldBindJSON(&credentials); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    for _, user := range users {
        if user.Username == credentials.Username && user.Password == credentials.Password {
            token, err := generateToken(user.Username, user.Role)
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
                return
            }
            c.JSON(http.StatusOK, gin.H{"token": token})
            return
        }
    }
    c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
}

func roleMiddleware(role string) gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
            c.Abort()
            return
        }

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return jwtSecretKey, nil
        })

        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok || claims["role"] != role {
            c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: insufficient privileges"})
            c.Abort()
            return
        }

        c.Next()
    }
}

func authMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
            c.Abort()
            return
        }

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return jwtSecretKey, nil
        })

        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        c.Next()
    }
}

func main() {
    r := gin.Default()

	r.Use(cors.Default())

    r.POST("/register", register)
    r.POST("/login", login)
    r.GET("/items", getItems)
    r.GET("/items/:id", getItemByID)

	authorized := r.Group("/")
    authorized.Use(authMiddleware())
    {
        authorized.POST("/items", createItem)
        authorized.PUT("/items/:id", updateItem)
        authorized.DELETE("/items/:id", deleteItem)
    }

	r.GET("/admin", roleMiddleware("admin"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome, Admin"})
	})
	
    r.Run() // Runs on localhost:8080 by default
}

// Create
func createItem(c *gin.Context) {
    var newItem Item
    if err := c.BindJSON(&newItem); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    newItem.ID = nextID
    nextID++
    items = append(items, newItem)
    c.JSON(http.StatusCreated, newItem)
}

// Read all items
func getItems(c *gin.Context) {
    c.JSON(http.StatusOK, items)
}

// Read single item by ID
func getItemByID(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }
    for _, item := range items {
        if item.ID == id {
            c.JSON(http.StatusOK, item)
            return
        }
    }
    c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
}

// Update
func updateItem(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    var updatedItem Item
    if err := c.BindJSON(&updatedItem); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    for i, item := range items {
        if item.ID == id {
            updatedItem.ID = id
            items[i] = updatedItem
            c.JSON(http.StatusOK, updatedItem)
            return
        }
    }
    c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
}

// Delete
func deleteItem(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }
    for i, item := range items {
        if item.ID == id {
            items = append(items[:i], items[i+1:]...)
            c.JSON(http.StatusOK, gin.H{"message": "Item deleted"})
            return
        }
    }
    c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
}
