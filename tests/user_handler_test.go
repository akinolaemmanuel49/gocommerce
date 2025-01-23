package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	auth_middlewares "github.com/akinolaemmanuel49/gocommerce/internal/auth/middlewares"
	"github.com/akinolaemmanuel49/gocommerce/internal/auth/utils"
	"github.com/akinolaemmanuel49/gocommerce/internal/handlers"
	"github.com/akinolaemmanuel49/gocommerce/internal/models"
	"github.com/akinolaemmanuel49/gocommerce/internal/repositories"
	"github.com/akinolaemmanuel49/gocommerce/internal/services"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	mongoURI       = "mongodb://localhost:27017"
	dbName         = "test_gocommerce"
	seedFile       = "../tests/test_data_users.json"
	collectionName = "users"
)

type userDbIn struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty"`
	Email               string             `bson:"email,omitempty"`
	PasswordHash        string             `bson:"passwordHash,omitempty"`
	FirstName           string             `bson:"firstName,omitempty"`
	LastName            string             `bson:"lastName,omitempty"`
	Address             models.Address     `bson:"address,omitempty"`
	Phone               string             `bson:"phone,omitempty"`
	Role                string             `bson:"role,omitempty"`
	models.CommonFields `bson:"inline"`
}

func TestUserHandler_Create(t *testing.T) {
	// Setup the test database
	testDB := setupTest(t)
	defer testDB.TeardownTestDatabase()

	userHandler := spawnUserHandler(testDB.Database)

	// Set up the router
	router := mux.NewRouter()
	router.HandleFunc("/users", userHandler.Create).Methods(http.MethodPost)

	// Define test cases
	tests := []struct {
		name         string
		payload      map[string]string
		expectedCode int
		expectedRole string // use empty string for failure cases
	}{
		{
			name: "Create Customer",
			payload: map[string]string{
				"email":     "samueldoe@example.com",
				"password":  "password",
				"firstName": "Samuel",
				"lastName":  "Doe",
				"role":      "customer",
			},
			expectedCode: http.StatusCreated,
			expectedRole: "customer",
		},
		{
			name: "Create Admin",
			payload: map[string]string{
				"email":     "jeremiahdoe@example.com",
				"password":  "password",
				"firstName": "Jeremiah",
				"lastName":  "Doe",
				"role":      "admin",
			},
			expectedCode: http.StatusCreated,
			expectedRole: "admin",
		},
		{
			name: "Bad Request",
			payload: map[string]string{
				"email":     "",
				"password":  "password",
				"firstName": "Isaac",
				"lastName":  "Doe",
				"role":      "customer",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Conflict",
			payload: map[string]string{
				"email":     "jeremiahdoe@example.com",
				"password":  "password",
				"firstName": "Jeremiah",
				"lastName":  "Doe",
				"role":      "customer",
			},
			expectedCode: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal the payload to JSON
			body, _ := json.Marshal(tt.payload)

			// Create the HTTP request
			req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Serve the request using the router
			router.ServeHTTP(rr, req)

			// Assert the response status code
			assert.Equal(t, tt.expectedCode, rr.Code)

			// If the expected role is set, validate the response body
			if tt.expectedRole != "" {
				var createdUser models.User
				json.NewDecoder(rr.Body).Decode(&createdUser)
				assert.Equal(t, tt.expectedRole, createdUser.Role)
			}
		})
	}
}

func TestUserHandler_Read(t *testing.T) {
	// Setup the test database
	testDB := setupTest(t)
	defer testDB.TeardownTestDatabase()

	userHandler := spawnUserHandler(testDB.Database)

	// Generate JWT token
	jwtSecretKey := []byte("jwt-secret-key")

	// Seed the database with the mock users
	customerObjectID, _ := primitive.ObjectIDFromHex("679203704b42eafa5d57d30b")
	adminObjectID, _ := primitive.ObjectIDFromHex("679203704b42eafa5d57d30a")

	// Set up the router
	router := mux.NewRouter()
	authMiddleware := auth_middlewares.AuthMiddleware(jwtSecretKey)
	router.Handle("/users", authMiddleware(http.HandlerFunc(userHandler.Read))).Methods(http.MethodGet)

	// Mock customer user
	mockAdminUser := userDbIn{
		ID:           adminObjectID,
		Email:        "valentinadoe@example.com",
		PasswordHash: "$2a$10$tbP6YLiT1A8rWwzdNthKAugBvmc5zF8GSF6QDdewDFh9pfWpqcvgW",
		FirstName:    "Valentina",
		LastName:     "Doe",
		Address:      models.Address{},
		Phone:        "91 123 456",
		Role:         "admin",
		CommonFields: models.CommonFields{
			CreatedAt: time.Date(2025, time.January, 23, 12, 53, 47, 406286520, time.Local),
			UpdatedAt: time.Date(2025, time.January, 23, 12, 53, 47, 406286705, time.Local),
			IsDeleted: false,
		},
	}

	mockCustomerUser := userDbIn{
		ID:           customerObjectID,
		Email:        "brandondoe@example.com",
		PasswordHash: "$2a$10$tbP6YLiT1A8rWwzdNthKAugBvmc5zF8GSF6QDdewDFh9pfWpqcvgW",
		FirstName:    "Brandon",
		LastName:     "Doe",
		Address:      models.Address{},
		Phone:        "91 123 456",
		Role:         "customer",
		CommonFields: models.CommonFields{
			CreatedAt: time.Date(2025, time.January, 23, 12, 53, 47, 406286520, time.Local),
			UpdatedAt: time.Date(2025, time.January, 23, 12, 53, 47, 406286705, time.Local),
			IsDeleted: false,
		},
	}

	// Insert mock users
	_, err := testDB.Database.Collection("users").InsertOne(context.TODO(), mockAdminUser)
	if err != nil {
		t.Fatalf("Failed to insert mock user: %v", err)
	}

	_, err = testDB.Database.Collection("users").InsertOne(context.TODO(), mockCustomerUser)
	if err != nil {
		t.Fatalf("Failed to insert mock user: %v", err)
	}

	var insertedAdminUser models.User
	err = testDB.Database.Collection("users").FindOne(context.TODO(), bson.M{"_id": adminObjectID}).Decode(&insertedAdminUser)
	if err != nil {
		t.Fatalf("Failed to retrieve inserted user: %v", err)
	}

	var insertedCustomerUser models.User
	err = testDB.Database.Collection("users").FindOne(context.TODO(), bson.M{"_id": customerObjectID}).Decode(&insertedCustomerUser)
	if err != nil {
		t.Fatalf("Failed to retrieve inserted user: %v", err)
	}

	tests := []struct {
		name         string
		userID       string
		query        string
		otherUserID  string
		role         string
		expectedUser models.User
		expectedCode int
	}{
		{name: "Read Own Customer User",
			userID:      "679203704b42eafa5d57d30b",
			query:       "",
			otherUserID: "",
			role:        "customer",
			expectedUser: models.User{
				ID:        "679203704b42eafa5d57d30b",
				Email:     "brandondoe@example.com",
				FirstName: "Brandon",
				LastName:  "Doe",
				Address:   models.Address{},
				Phone:     "91 123 456",
				Role:      "customer",
				CommonFields: models.CommonFields{
					CreatedAt: insertedCustomerUser.CreatedAt,
					UpdatedAt: insertedCustomerUser.UpdatedAt,
					IsDeleted: false,
				},
			},
			expectedCode: http.StatusOK,
		},
		{name: "Read Own Admin User",
			userID:      "679203704b42eafa5d57d30a",
			query:       "",
			otherUserID: "",
			role:        "admin",
			expectedUser: models.User{
				ID:        "679203704b42eafa5d57d30a",
				Email:     "valentinadoe@example.com",
				FirstName: "Valentina",
				LastName:  "Doe",
				Address:   models.Address{},
				Phone:     "91 123 456",
				Role:      "admin",
				CommonFields: models.CommonFields{
					CreatedAt: insertedCustomerUser.CreatedAt,
					UpdatedAt: insertedCustomerUser.UpdatedAt,
					IsDeleted: false,
				},
			},
			expectedCode: http.StatusOK,
		},
		{name: "Read Other Admin User",
			userID:      "679203704b42eafa5d57d30a",
			query:       "?id=",
			otherUserID: "679203704b42eafa5d57d30b",
			role:        "admin",
			expectedUser: models.User{
				ID:        "679203704b42eafa5d57d30b",
				Email:     "brandondoe@example.com",
				FirstName: "Brandon",
				LastName:  "Doe",
				Address:   models.Address{},
				Phone:     "91 123 456",
				Role:      "customer",
				CommonFields: models.CommonFields{
					CreatedAt: insertedCustomerUser.CreatedAt,
					UpdatedAt: insertedCustomerUser.UpdatedAt,
					IsDeleted: false,
				},
			},
			expectedCode: http.StatusOK,
		},
		{name: "Read Other Customer User",
			userID:       "679203704b42eafa5d57d30b",
			query:        "?id=",
			otherUserID:  "679203704b42eafa5d57d30a",
			role:         "customer",
			expectedCode: http.StatusForbidden,
		},
		{name: "Read Own Non Existent User",
			userID:       "679203704b42eafa5d57d30c",
			query:        "",
			otherUserID:  "",
			role:         "customer",
			expectedCode: http.StatusNotFound,
		},
		{name: "Read Other Non Existent User",
			userID:       "679203704b42eafa5d57d30a",
			query:        "?id=",
			otherUserID:  "679203704b42eafa5d57d30c",
			role:         "admin",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create the HTTP request
			req, _ := http.NewRequest(http.MethodGet, "/users"+tt.query+tt.otherUserID, nil)
			req.Header.Set("Content-Type", "application/json")

			req = getRequestAuthenticated(t, jwtSecretKey, tt.userID, tt.role, req)

			// Create response recorder
			rr := httptest.NewRecorder()

			// Serve the request using the router
			router.ServeHTTP(rr, req)

			var user models.User
			if err := json.NewDecoder(rr.Body).Decode(&user); err != nil {
				t.Fatalf("Failed to decode response body: %v", err)
			}

			// Assert the response status code
			assert.Equal(t, tt.expectedCode, rr.Code)
			if tt.expectedUser != (models.User{}) {
				// Hacky, but fixes the problem with time fields
				user.CreatedAt = insertedCustomerUser.CreatedAt
				user.UpdatedAt = insertedCustomerUser.UpdatedAt
				assert.Equal(t, tt.expectedUser, user)
			}
		})
	}
}

func TestUserHandler_ReadAll(t *testing.T) {
	// Setup the test database
	testDB := setupTest(t)

	defer testDB.TeardownTestDatabase()

	userHandler := spawnUserHandler(testDB.Database)

	// Generate JWT token
	jwtSecretKey := []byte("jwt-secret-key")

	// Seed the database with the mock users
	customerObjectID, _ := primitive.ObjectIDFromHex("679203704b42eafa5d57d30b")
	adminObjectID, _ := primitive.ObjectIDFromHex("679203704b42eafa5d57d30a")

	// Mock customer user
	mockAdminUser := userDbIn{
		ID:           adminObjectID,
		Email:        "valentinadoe@example.com",
		PasswordHash: "$2a$10$tbP6YLiT1A8rWwzdNthKAugBvmc5zF8GSF6QDdewDFh9pfWpqcvgW",
		FirstName:    "Valentina",
		LastName:     "Doe",
		Address:      models.Address{},
		Phone:        "91 123 456",
		Role:         "admin",
		CommonFields: models.CommonFields{
			CreatedAt: time.Date(2025, time.January, 23, 12, 53, 47, 406286520, time.Local),
			UpdatedAt: time.Date(2025, time.January, 23, 12, 53, 47, 406286705, time.Local),
			IsDeleted: false,
		},
	}

	mockCustomerUser := userDbIn{
		ID:           customerObjectID,
		Email:        "brandondoe@example.com",
		PasswordHash: "$2a$10$tbP6YLiT1A8rWwzdNthKAugBvmc5zF8GSF6QDdewDFh9pfWpqcvgW",
		FirstName:    "Brandon",
		LastName:     "Doe",
		Address:      models.Address{},
		Phone:        "91 123 456",
		Role:         "customer",
		CommonFields: models.CommonFields{
			CreatedAt: time.Date(2025, time.January, 23, 12, 53, 47, 406286520, time.Local),
			UpdatedAt: time.Date(2025, time.January, 23, 12, 53, 47, 406286705, time.Local),
			IsDeleted: false,
		},
	}

	// Insert mock users
	_, err := testDB.Database.Collection("users").InsertOne(context.TODO(), mockAdminUser)
	if err != nil {
		t.Fatalf("Failed to insert mock user: %v", err)
	}

	_, err = testDB.Database.Collection("users").InsertOne(context.TODO(), mockCustomerUser)
	if err != nil {
		t.Fatalf("Failed to insert mock user: %v", err)
	}

	var insertedAdminUser models.User
	err = testDB.Database.Collection("users").FindOne(context.TODO(), bson.M{"_id": adminObjectID}).Decode(&insertedAdminUser)
	if err != nil {
		t.Fatalf("Failed to retrieve inserted user: %v", err)
	}

	var insertedCustomerUser models.User
	err = testDB.Database.Collection("users").FindOne(context.TODO(), bson.M{"_id": customerObjectID}).Decode(&insertedCustomerUser)
	if err != nil {
		t.Fatalf("Failed to retrieve inserted user: %v", err)
	}

	// Success Case: Read All Users Admin
	t.Run("Read All Users Admin", func(t *testing.T) {
		// Create a mock HTTP request
		req, err := http.NewRequest(http.MethodGet, "/users/all", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		// JWT token for user Valentina Doe
		validTokenString, _ := utils.GenerateJWT(jwtSecretKey, "679203704b42eafa5d57d30a", "admin")

		type contextKey string
		const UserClaimsKey contextKey = "userClaims"

		req.Header.Set("Authorization", "Bearer "+validTokenString)
		req.Header.Set("Content-Type", "application/json")

		// Parse and validate the token
		claims, err := utils.ParseJWT(jwtSecretKey, validTokenString)
		if err != nil {
			t.Fatalf("Failed to extract claims from token")
		}

		ctx := context.WithValue(req.Context(), UserClaimsKey, claims)
		req = req.WithContext(ctx)

		// Create response recorder
		rr := httptest.NewRecorder()

		// Call the handler's ReadAll method
		authMiddleware := auth_middlewares.AuthMiddleware(jwtSecretKey)
		chain := authMiddleware(http.HandlerFunc(userHandler.ReadAll))
		chain.ServeHTTP(rr, req)

		// Make assertions
		assert.Equal(t, req.Header.Get("Authorization"), "Bearer "+validTokenString)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	// Failure Case: Read All Users Customer
	t.Run("Read All Users Customer", func(t *testing.T) {
		// Create a mock HTTP request
		req, err := http.NewRequest(http.MethodGet, "/users/all", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		// JWT token for user Valentina Doe
		validTokenString, _ := utils.GenerateJWT(jwtSecretKey, "679203704b42eafa5d57d30b", "customer")

		type contextKey string
		const UserClaimsKey contextKey = "userClaims"

		req.Header.Set("Authorization", "Bearer "+validTokenString)
		req.Header.Set("Content-Type", "application/json")

		// Parse and validate the token
		claims, err := utils.ParseJWT(jwtSecretKey, validTokenString)
		if err != nil {
			t.Fatalf("Failed to extract claims from token")
		}

		ctx := context.WithValue(req.Context(), UserClaimsKey, claims)
		req = req.WithContext(ctx)

		// Create response recorder
		rr := httptest.NewRecorder()

		// Call the handler's ReadAll method
		authMiddleware := auth_middlewares.AuthMiddleware(jwtSecretKey)
		chain := authMiddleware(http.HandlerFunc(userHandler.ReadAll))
		chain.ServeHTTP(rr, req)

		// Make assertions
		assert.Equal(t, req.Header.Get("Authorization"), "Bearer "+validTokenString)
		assert.Equal(t, http.StatusForbidden, rr.Code)
	})

	// Failure Case: Unauthorized
	t.Run("Unauthorized", func(t *testing.T) {
		// Create a mock HTTP request
		req, err := http.NewRequest(http.MethodGet, "/users/all", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Call the handler's Read method
		authMiddleware := auth_middlewares.AuthMiddleware(jwtSecretKey)
		chain := authMiddleware(http.HandlerFunc(userHandler.ReadAll))
		chain.ServeHTTP(rr, req)

		// Make assertions
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})
}

func TestUserHandler_Update(t *testing.T) {
	// Setup the test database
	testDB := setupTest(t)

	defer testDB.TeardownTestDatabase()

	userHandler := spawnUserHandler(testDB.Database)

	// Success Case: Update Customer
	t.Run("Create Customer", func(t *testing.T) {
		// Create a mock HTTP request
		payload := map[string]string{
			"email":     "samueldoe@example.com",
			"password":  "password",
			"firstName": "Samuel",
			"lastName":  "Doe",
			"role":      "customer", // role is customer
		}
		body, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Call the handler's Create method
		http.HandlerFunc(userHandler.Create).ServeHTTP(rr, req)

		// Make assertions
		assert.Equal(t, http.StatusCreated, rr.Code)

		var createdUser models.User
		json.NewDecoder(rr.Body).Decode(&createdUser)
		assert.Equal(t, "customer", createdUser.Role)
	})

	t.Run("Create Admin", func(t *testing.T) {
		// Create a mock HTTP request
		payload := map[string]string{
			"email":     "jeremiahdoe@example.com",
			"password":  "password",
			"firstName": "Jeremiah",
			"lastName":  "Doe",
			"role":      "admin", // role is admin
		}
		body, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Call the handler's Create method
		http.HandlerFunc(userHandler.Create).ServeHTTP(rr, req)

		// Make assertions
		assert.Equal(t, http.StatusCreated, rr.Code)

		var createdUser models.User
		json.NewDecoder(rr.Body).Decode(&createdUser)
		assert.Equal(t, "admin", createdUser.Role)
	})

	// Failure Case: Email Empty String
	t.Run("Email Empty String", func(t *testing.T) {
		// Create a mock HTTP request
		payload := map[string]string{
			"email":     "", // email is empty string
			"password":  "password",
			"firstName": "Isaac",
			"lastName":  "Doe",
			"role":      "customer",
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Call the handler's Create method
		http.HandlerFunc(userHandler.Create).ServeHTTP(rr, req)

		// Assert the response
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	// Failure Case: Email Invalid
	t.Run("Email Invalid", func(t *testing.T) {
		// Create a mock HTTP request
		payload := map[string]string{
			"email":     "isaacdoeexample.com", // email is missing an '@'
			"password":  "password",
			"firstName": "Isaac",
			"lastName":  "Doe",
			"role":      "customer",
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Call the handler's Create method
		http.HandlerFunc(userHandler.Create).ServeHTTP(rr, req)

		// Assert the response
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	// Failure Case: Password Empty String
	t.Run("Password Empty String", func(t *testing.T) {
		// Create a mock HTTP request
		payload := map[string]string{
			"email":     "isaacdoe@example.com",
			"password":  "", // password is empty string
			"firstName": "Isaac",
			"lastName":  "Doe",
			"role":      "customer",
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Call the handler's Create method
		http.HandlerFunc(userHandler.Create).ServeHTTP(rr, req)

		// Assert the response
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	// Failure Case: Password Length Validation
	t.Run("Password Length Validation", func(t *testing.T) {
		// Create a mock HTTP request
		payload := map[string]string{
			"email":     "isaacdoe@example.com",
			"password":  "pass", // password length is less than 8
			"firstName": "Isaac",
			"lastName":  "Doe",
			"role":      "customer",
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Call the handler's Create method
		http.HandlerFunc(userHandler.Create).ServeHTTP(rr, req)

		// Assert the response
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	// Failure Case: First Name Empty String
	t.Run("First Name Empty String", func(t *testing.T) {
		// Create a mock HTTP request
		payload := map[string]string{
			"email":     "isaacdoe@example.com",
			"password":  "password",
			"firstName": "", // firstName is empty string
			"lastName":  "Doe",
			"role":      "customer",
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Call the handler's Create method
		http.HandlerFunc(userHandler.Create).ServeHTTP(rr, req)

		// Assert the response
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	// Failure Case: Last Name Empty String
	t.Run("Last Name Empty String", func(t *testing.T) {
		// Create a mock HTTP request
		payload := map[string]string{
			"email":     "isaacdoe@example.com",
			"password":  "password",
			"firstName": "Isaac",
			"lastName":  "", // lastName is empty string
			"role":      "customer",
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Call the handler's Create method
		http.HandlerFunc(userHandler.Create).ServeHTTP(rr, req)

		// Assert the response
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	// Failure Case: Role Empty String
	t.Run("Role Empty String", func(t *testing.T) {
		// Create a mock HTTP request
		payload := map[string]string{
			"email":     "isaacdoe@example.com",
			"password":  "password",
			"firstName": "Isaac",
			"lastName":  "Doe",
			"role":      "", // role is empty string
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Call the handler's Create method
		http.HandlerFunc(userHandler.Create).ServeHTTP(rr, req)

		// Assert the response
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	// Failure Case: Email Conflict
	t.Run("Email Conflict", func(t *testing.T) {
		// Create a mock HTTP request
		payload := map[string]string{
			"email":     "jeremiahdoe@example.com", // email already in use
			"password":  "password",
			"firstName": "Jeremiah",
			"lastName":  "Doe",
			"role":      "customer",
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Call the handler's Create method
		http.HandlerFunc(userHandler.Create).ServeHTTP(rr, req)

		// Assert the response
		assert.Equal(t, http.StatusConflict, rr.Code)
	})
}

// Spawn user handler
func spawnUserHandler(db *mongo.Database) *handlers.UserHandler {
	userService := spawnUserService(db)
	logger := log.New(io.Discard, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger := log.New(io.Discard, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	return handlers.NewUserHandler(userService, logger, errorLogger)
}

// Spawn user service
func spawnUserService(db *mongo.Database) *services.UserService {
	userRepository := repositories.NewUserRepository(db)
	return services.NewUserService(userRepository)
}

func setupTest(t *testing.T) *TestDatabase {
	testDB, err := SetupTestDatabase(mongoURI, dbName)
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}

	// Seed the database
	if err := testDB.SeedDatabase(collectionName, seedFile); err != nil {
		t.Fatalf("Failed to seed database: %v", err)
	}
	return testDB
}

func getRequestAuthenticated(t *testing.T, jwtSecretKey []byte, userID string, role string, req *http.Request) *http.Request {
	validTokenString, _ := utils.GenerateJWT(jwtSecretKey, userID, role)

	type contextKey string
	const UserClaimsKey contextKey = "userClaims"

	req.Header.Set("Authorization", "Bearer "+validTokenString)
	req.Header.Set("Content-Type", "application/json")

	// Parse and validate the token
	claims, err := utils.ParseJWT(jwtSecretKey, validTokenString)
	if err != nil {
		t.Fatalf("Failed to extract claims from token")
	}

	ctx := context.WithValue(req.Context(), UserClaimsKey, claims)
	req = req.WithContext(ctx)
	return req
}
