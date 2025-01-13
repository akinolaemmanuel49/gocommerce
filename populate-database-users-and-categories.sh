!#/bin/bash

# Create users
curl -X POST http://localhost:8000/users \
-H "Content-Type: application/json" \
-d '{"email": "john.smith@example.com", "password": "password", "firstName": "John", "lastName": "Smith", "role": "customer"}'

curl -X POST http://localhost:8000/users \
-H "Content-Type: application/json" \
-d '{"email": "emma.jones@example.com", "password": "password", "firstName": "Emma", "lastName": "Jones", "role": "admin"}'

curl -X POST http://localhost:8000/users \
-H "Content-Type: application/json" \
-d '{"email": "lisa.brown@example.com", "password": "password", "firstName": "Lisa", "lastName": "Brown", "role": "customer"}'

curl -X POST http://localhost:8000/users \
-H "Content-Type: application/json" \
-d '{"email": "michael.johnson@example.com", "password": "password", "firstName": "Michael", "lastName": "Johnson", "role": "customer"}'

curl -X POST http://localhost:8000/users \
-H "Content-Type: application/json" \
-d '{"email": "sarah.davis@example.com", "password": "password", "firstName": "Sarah", "lastName": "Davis", "role": "admin"}'

curl -X POST http://localhost:8000/users \
-H "Content-Type: application/json" \
-d '{"email": "robert.martin@example.com", "password": "password", "firstName": "Robert", "lastName": "Martin", "role": "customer"}'

curl -X POST http://localhost:8000/users \
-H "Content-Type: application/json" \
-d '{"email": "karen.moore@example.com", "password": "password", "firstName": "Karen", "lastName": "Moore", "role": "customer"}'

curl -X POST http://localhost:8000/users \
-H "Content-Type: application/json" \
-d '{"email": "david.lee@example.com", "password": "password", "firstName": "David", "lastName": "Lee", "role": "admin"}'

curl -X POST http://localhost:8000/users \
-H "Content-Type: application/json" \
-d '{"email": "amy.white@example.com", "password": "password", "firstName": "Amy", "lastName": "White", "role": "customer"}'

curl -X POST http://localhost:8000/users \
-H "Content-Type: application/json" \
-d '{"email": "james.clark@example.com", "password": "password", "firstName": "James", "lastName": "Clark", "role": "customer"}'

# Create categories
curl -X POST http://localhost:8000/categories \
-H "Content-Type: application/json" \
-d '{
    "name": "Electronics",
    "description": "Latest gadgets and electronic devices.",
    "image": "https://example.com/images/electronics.jpg"
}'

curl -X POST http://localhost:8000/categories \
-H "Content-Type: application/json" \
-d '{
    "name": "Fashion",
    "description": "Trendy clothing and accessories.",
    "image": "https://example.com/images/fashion.jpg"
}'

curl -X POST http://localhost:8000/categories \
-H "Content-Type: application/json" \
-d '{
    "name": "Home & Kitchen",
    "description": "Appliances, furniture, and kitchenware.",
    "image": "https://example.com/images/home-kitchen.jpg"
}'

curl -X POST http://localhost:8000/categories \
-H "Content-Type: application/json" \
-d '{
    "name": "Sports & Outdoors",
    "description": "Gear and equipment for outdoor activities.",
    "image": "https://example.com/images/sports-outdoors.jpg"
}'

curl -X POST http://localhost:8000/categories \
-H "Content-Type: application/json" \
-d '{
    "name": "Books",
    "description": "A wide selection of books across genres.",
    "image": "https://example.com/images/books.jpg"
}'

# Create products

# Create orders