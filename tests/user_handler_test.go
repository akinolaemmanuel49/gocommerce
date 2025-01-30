package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/akinolaemmanuel49/gocommerce/configs"
	auth_middlewares "github.com/akinolaemmanuel49/gocommerce/internal/auth/middlewares"
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

var mongoURI = configs.GetMongoDBURI()

const (
	dbName         = "GoCommerceTest"
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
	fmt.Println("LOUD!!!")
	fmt.Println(mongoURI)
	// Setup the test database
	testDB := setupUserTest(t)
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
	testDB := setupUserTest(t)
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

	createdAt := time.Date(2025, time.January, 23, 12, 53, 47, 406000000, time.UTC)
	updatedAt := time.Date(2025, time.January, 23, 12, 53, 47, 406000000, time.UTC)

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
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
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
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
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
		{
			name:        "Read Own Customer User",
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
					CreatedAt: createdAt,
					UpdatedAt: updatedAt,
					IsDeleted: false,
				},
			},
			expectedCode: http.StatusOK,
		},
		{
			name:        "Read Own Admin User",
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
					CreatedAt: createdAt,
					UpdatedAt: updatedAt,
					IsDeleted: false,
				},
			},
			expectedCode: http.StatusOK,
		},
		{
			name:        "Read Other Admin User",
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
					CreatedAt: createdAt,
					UpdatedAt: updatedAt,
					IsDeleted: false,
				},
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "Read Other Customer User",
			userID:       "679203704b42eafa5d57d30b",
			query:        "?id=",
			otherUserID:  "679203704b42eafa5d57d30a",
			role:         "customer",
			expectedCode: http.StatusForbidden,
		},
		{
			name:         "Read Own Non Existent User",
			userID:       "679203704b42eafa5d57d30c",
			query:        "",
			otherUserID:  "",
			role:         "customer",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "Read Other Non Existent User",
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

			req = GetRequestAuthenticated(t, jwtSecretKey, tt.userID, tt.role, req)

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
				// user.CreatedAt = createdAt
				// user.UpdatedAt = updatedAt
				assert.Equal(t, tt.expectedUser, user)
			}
		})
	}
}

func TestUserHandler_ReadAll(t *testing.T) {
	// Setup the test database
	testDB := setupUserTest(t)

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
	router.Handle("/users/all", authMiddleware(http.HandlerFunc(userHandler.ReadAll))).Methods(http.MethodGet)

	createdAt := time.Date(2025, time.January, 23, 12, 53, 47, 406000000, time.UTC)
	updatedAt := time.Date(2025, time.January, 23, 12, 53, 47, 406000000, time.UTC)

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
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
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
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
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
		role         string
		expectedCode int
	}{
		{
			name:         "Read All Users Admin",
			userID:       insertedAdminUser.ID,
			role:         insertedAdminUser.Role,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Forbidden",
			userID:       insertedCustomerUser.ID,
			role:         insertedCustomerUser.Role,
			expectedCode: http.StatusForbidden,
		},
		{
			name:         "Not Found",
			userID:       "679203704b42eafa5d57d30c",
			role:         "admin",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "Unauthorized",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create the HTTP request
			req, _ := http.NewRequest(http.MethodGet, "/users/all", nil)

			req = GetRequestAuthenticated(t, jwtSecretKey, tt.userID, tt.role, req)
			if tt.name == "Unauthorized" {
				req.Header.Del("Authorization")
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Serve the request using the router
			router.ServeHTTP(rr, req)

			// Assert the response status code
			assert.Equal(t, tt.expectedCode, rr.Code)
		})
	}
}

func TestUserHandler_Update(t *testing.T) {
	// Setup the test database
	testDB := setupUserTest(t)

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
	router.Handle("/users", authMiddleware(http.HandlerFunc(userHandler.Update))).Methods(http.MethodPut)

	createdAt := time.Date(2025, time.January, 23, 12, 53, 47, 406000000, time.UTC)
	updatedAt := time.Date(2025, time.January, 23, 12, 53, 47, 406000000, time.UTC)

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
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
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
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
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
		payload      map[string]interface{}
		expectedUser models.User
		expectedCode int
	}{
		{
			name:   "Update Own Admin User",
			userID: insertedAdminUser.ID,
			role:   insertedAdminUser.Role,
			payload: map[string]interface{}{
				"firstName": "Valeria",
				"lastName":  "Smart",
				"phone":     "91 321 654",
				"address": models.Address{
					Street:  "21 Baker St.",
					City:    "Manchester",
					State:   "Grand",
					Zip:     "112233",
					Country: "Ujigan",
				},
			},
			expectedUser: models.User{
				ID:    "679203704b42eafa5d57d30a",
				Email: "valentinadoe@example.com",
				// PasswordHash: "$2a$10$tbP6YLiT1A8rWwzdNthKAugBvmc5zF8GSF6QDdewDFh9pfWpqcvgW",
				FirstName: "Valeria",
				LastName:  "Smart",
				Address: models.Address{
					Street:  "21 Baker St.",
					City:    "Manchester",
					State:   "Grand",
					Zip:     "112233",
					Country: "Ujigan",
				},
				Phone: "91 321 654",
				Role:  "admin",
				CommonFields: models.CommonFields{
					CreatedAt: createdAt,
					UpdatedAt: updatedAt,
					IsDeleted: false,
				},
			},
			expectedCode: http.StatusOK,
		}, {
			name:   "Update Own Customer User",
			userID: insertedCustomerUser.ID,
			role:   insertedCustomerUser.Role,
			payload: map[string]interface{}{
				"firstName": "Brady",
				"lastName":  "Smart",
				"phone":     "91 321 654",
				"address": models.Address{
					Street:  "21 Baker St.",
					City:    "Manchester",
					State:   "Grand",
					Zip:     "112233",
					Country: "Ujigan",
				},
			},
			expectedUser: models.User{
				ID:    "679203704b42eafa5d57d30b",
				Email: "brandondoe@example.com",
				// PasswordHash: "$2a$10$tbP6YLiT1A8rWwzdNthKAugBvmc5zF8GSF6QDdewDFh9pfWpqcvgW",
				FirstName: "Brady",
				LastName:  "Smart",
				Address: models.Address{
					Street:  "21 Baker St.",
					City:    "Manchester",
					State:   "Grand",
					Zip:     "112233",
					Country: "Ujigan",
				},
				Phone: "91 321 654",
				Role:  "customer",
				CommonFields: models.CommonFields{
					CreatedAt: createdAt,
					UpdatedAt: updatedAt,
				},
			},
			expectedCode: http.StatusOK,
		}, {
			name:        "Update Other Admin User",
			userID:      insertedAdminUser.ID,
			role:        insertedAdminUser.Role,
			query:       "?id=",
			otherUserID: "679203704b42eafa5d57d30b",
			payload: map[string]interface{}{
				"firstName": "Brandon",
				"lastName":  "Doe",
				"phone":     "91 321 654",
				"address": models.Address{
					Street:  "21 Baker St.",
					City:    "Manchester",
					State:   "Grand",
					Zip:     "112233",
					Country: "Ujigan",
				},
			},
			expectedUser: models.User{
				ID:    "679203704b42eafa5d57d30b",
				Email: "brandondoe@example.com",
				// PasswordHash: "$2a$10$tbP6YLiT1A8rWwzdNthKAugBvmc5zF8GSF6QDdewDFh9pfWpqcvgW",
				FirstName: "Brandon",
				LastName:  "Doe",
				Address: models.Address{
					Street:  "21 Baker St.",
					City:    "Manchester",
					State:   "Grand",
					Zip:     "112233",
					Country: "Ujigan",
				},
				Phone: "91 321 654",
				Role:  "customer",
				CommonFields: models.CommonFields{
					CreatedAt: createdAt,
					UpdatedAt: updatedAt,
				},
			},
			expectedCode: http.StatusOK,
		}, {
			name:         "Forbidden",
			userID:       insertedCustomerUser.ID,
			role:         insertedCustomerUser.Role,
			query:        "?id=",
			otherUserID:  "679203704b42eafa5d57d30a",
			expectedCode: http.StatusForbidden,
		}, {
			name:         "Not Found",
			userID:       insertedAdminUser.ID,
			role:         insertedAdminUser.Role,
			query:        "?id=",
			otherUserID:  "679203704b42eafa5d57d30c",
			expectedCode: http.StatusNotFound,
		}, {
			name:         "Unauthorized",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshall the payload to JSON
			body, _ := json.Marshal(tt.payload)

			// Create the HTTP request
			req, _ := http.NewRequest(http.MethodPut, "/users"+tt.query+tt.otherUserID, bytes.NewReader(body))
			req = GetRequestAuthenticated(t, jwtSecretKey, tt.userID, tt.role, req)
			if tt.name == "Unauthorized" {
				req.Header.Del("Authorization")
			}

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Serve the request using the router
			router.ServeHTTP(rr, req)

			// Assert the response status code
			assert.Equal(t, tt.expectedCode, rr.Code)

			var updatedUser models.User
			json.NewDecoder(rr.Body).Decode(&updatedUser)
			if tt.expectedUser != (models.User{}) {
				updatedUser.CreatedAt = createdAt
				updatedUser.UpdatedAt = updatedAt
				assert.Equal(t, tt.expectedUser, updatedUser)
			}
		})
	}
}

func TestUserHandler_Delete(t *testing.T) {
	// Setup the test database
	testDB := setupUserTest(t)

	defer testDB.TeardownTestDatabase()

	userHandler := spawnUserHandler(testDB.Database)

	// Generate JWT token
	jwtSecretKey := []byte("jwt-secret-key")

	// Seed the database with the mock users
	adminObjectID, _ := primitive.ObjectIDFromHex("679203704b42eafa5d57d30a")
	customerAObjectID, _ := primitive.ObjectIDFromHex("679203704b42eafa5d57d30b")
	customerBObjectID, _ := primitive.ObjectIDFromHex("679203704b42eafa5d57d30c")

	// Set up the router
	router := mux.NewRouter()
	authMiddleware := auth_middlewares.AuthMiddleware(jwtSecretKey)
	router.Handle("/users", authMiddleware(http.HandlerFunc(userHandler.Update))).Methods(http.MethodPut)

	createdAt := time.Date(2025, time.January, 23, 12, 53, 47, 406000000, time.UTC)
	updatedAt := time.Date(2025, time.January, 23, 12, 53, 47, 406000000, time.UTC)

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
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			IsDeleted: false,
		},
	}

	mockCustomerUserA := userDbIn{
		ID:           customerAObjectID,
		Email:        "brandondoe@example.com",
		PasswordHash: "$2a$10$tbP6YLiT1A8rWwzdNthKAugBvmc5zF8GSF6QDdewDFh9pfWpqcvgW",
		FirstName:    "Brandon",
		LastName:     "Doe",
		Address:      models.Address{},
		Phone:        "91 123 456",
		Role:         "customer",
		CommonFields: models.CommonFields{
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			IsDeleted: false,
		},
	}

	mockCustomerUserB := userDbIn{
		ID:           customerBObjectID,
		Email:        "janicedoe@example.com",
		PasswordHash: "$2a$10$tbP6YLiT1A8rWwzdNthKAugBvmc5zF8GSF6QDdewDFh9pfWpqcvgW",
		FirstName:    "Janice",
		LastName:     "Doe",
		Address:      models.Address{},
		Phone:        "91 123 456",
		Role:         "customer",
		CommonFields: models.CommonFields{
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			IsDeleted: false,
		},
	}

	// Insert mock users
	_, err := testDB.Database.Collection("users").InsertOne(context.TODO(), mockAdminUser)
	if err != nil {
		t.Fatalf("Failed to insert mock user: %v", err)
	}

	_, err = testDB.Database.Collection("users").InsertOne(context.TODO(), mockCustomerUserA)
	if err != nil {
		t.Fatalf("Failed to insert mock user: %v", err)
	}

	_, err = testDB.Database.Collection("users").InsertOne(context.TODO(), mockCustomerUserB)
	if err != nil {
		t.Fatalf("Failed to insert mock user: %v", err)
	}
	var insertedAdminUser models.User
	err = testDB.Database.Collection("users").FindOne(context.TODO(), bson.M{"_id": adminObjectID}).Decode(&insertedAdminUser)
	if err != nil {
		t.Fatalf("Failed to retrieve inserted user: %v", err)
	}

	var insertedCustomerUserA models.User
	err = testDB.Database.Collection("users").FindOne(context.TODO(), bson.M{"_id": customerAObjectID}).Decode(&insertedCustomerUserA)
	if err != nil {
		t.Fatalf("Failed to retrieve inserted user: %v", err)
	}

	var insertedCustomerUserB models.User
	err = testDB.Database.Collection("users").FindOne(context.TODO(), bson.M{"_id": customerBObjectID}).Decode(&insertedCustomerUserB)
	if err != nil {
		t.Fatalf("Failed to retrieve inserted user: %v", err)
	}

	tests := []struct {
		name         string
		userID       string
		query        string
		otherUserID  string
		role         string
		payload      map[string]interface{}
		expectedUser models.User
		expectedCode int
	}{
		{
			name:         "Delete Own Customer User",
			userID:       insertedCustomerUserA.ID,
			role:         insertedCustomerUserA.Role,
			expectedCode: http.StatusOK,
		}, {
			name:         "Forbidden",
			userID:       insertedCustomerUserB.ID,
			role:         insertedCustomerUserB.Role,
			query:        "?id=",
			otherUserID:  "679203704b42eafa5d57d30a",
			expectedCode: http.StatusForbidden,
		}, {
			name:         "Delete Other Admin User",
			userID:       insertedAdminUser.ID,
			role:         insertedAdminUser.Role,
			query:        "?id=",
			otherUserID:  "679203704b42eafa5d57d30c",
			expectedCode: http.StatusOK,
		}, {
			name:         "Delete Own Admin User",
			userID:       insertedAdminUser.ID,
			role:         insertedAdminUser.Role,
			expectedCode: http.StatusOK,
		}, {
			name:         "Unauthorized",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshall the payload to JSON
			body, _ := json.Marshal(tt.payload)

			// Create the HTTP request
			req, _ := http.NewRequest(http.MethodPut, "/users"+tt.query+tt.otherUserID, bytes.NewReader(body))
			req = GetRequestAuthenticated(t, jwtSecretKey, tt.userID, tt.role, req)
			if tt.name == "Unauthorized" {
				req.Header.Del("Authorization")
			}

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Serve the request using the router
			router.ServeHTTP(rr, req)

			// Assert the response status code
			assert.Equal(t, tt.expectedCode, rr.Code)

			var updatedUser models.User
			json.NewDecoder(rr.Body).Decode(&updatedUser)
			if tt.expectedUser != (models.User{}) {
				updatedUser.CreatedAt = createdAt
				updatedUser.UpdatedAt = updatedAt
				assert.Equal(t, tt.expectedUser, updatedUser)
			}
		})
	}
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

func setupUserTest(t *testing.T) *TestDatabase {
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
