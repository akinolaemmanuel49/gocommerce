curl -X GET http://localhost:8080/products -H "Content-Type: application/json"

curl -X GET http://localhost:8080/users -H "Content-Type: application/json"

curl -X GET http://localhost:8080/orders -H "Content-Type: application/json"

curl -X GET http://localhost:8080/categories -H "Content-Type: application/json"

curl -X POST http://localhost:8080/products \
-H "Content-Type: application/json" \
-d '{
  "name": "Sample Product",
  "description": "A detailed description of the sample product",
  "price": 49.99,
  "images": ["image1.jpg", "image2.jpg"],
  "categoryId": "sample-category-id",
  "brand": "Sample Brand"
}'

curl -X POST http://localhost:8080/users \
-H "Content-Type: application/json" \
-d '{
  "email": "johndoe@example.com",
  "password": "password",
  "firstName": "John",
  "lastName": "Doe",
  "role": "customer"
}'

curl -X POST http://localhost:8080/orders \
-H "Content-Type: application/json" \
-d '{
  "userId": "676d92ed2ebcfc217fade2b5",
  "items": [
    { "productId": "product-id-1", "quantity": 2, "price": 19.99 },
    { "productId": "product-id-2", "quantity": 1, "price": 49.99 }
  ],
  "totalPrice": 89.97,
  "status": "pending",
  "shippingAddress": {
    "street": "123 Main St",
    "city": "Sample City",
    "state": "Sample State",
    "zip": "12345",
    "country": "Sample Country"
  }
}'

curl -X POST http://localhost:8080/categories \
-H "Content-Type: application/json" \
-d '{
  "name": "Electronics",
  "description": "All kinds of electronic gadgets and devices"
}'
