curl -X POST http://localhost:8000/products \
-H "Content-Type: application/json" \
-d '{
    "name": "Classic Leather Jacket",
    "description": "A timeless leather jacket made from premium materials.",
    "price": 199.99,
    "images": [
        "https://example.com/images/leather-jacket-1.jpg",
        "https://example.com/images/leather-jacket-2.jpg"
    ],
    "categoryId": "67827cb87c8364072523fb55",
    "brand": "UrbanStyle"
}'

curl -X POST http://localhost:8000/products \
-H "Content-Type: application/json" \
-d '{
    "name": "Slim Fit Jeans",
    "description": "Stylish and comfortable slim fit jeans.",
    "price": 49.99,
    "images": [
        "https://example.com/images/jeans-1.jpg",
        "https://example.com/images/jeans-2.jpg"
    ],
    "categoryId": "67827cb87c8364072523fb55",
    "brand": "DenimTrend"
}'

curl -X POST http://localhost:8000/products \
-H "Content-Type: application/json" \
-d '{
    "name": "Cotton Crew Neck T-Shirt",
    "description": "Soft and breathable cotton t-shirt available in multiple colors.",
    "price": 19.99,
    "images": [
        "https://example.com/images/tshirt-1.jpg",
        "https://example.com/images/tshirt-2.jpg"
    ],
    "categoryId": "67827cb87c8364072523fb55",
    "brand": "ComfortWear"
}'

curl -X POST http://localhost:8000/products \
-H "Content-Type: application/json" \
-d '{
    "name": "Designer Sneakers",
    "description": "Trendy and comfortable sneakers for all-day wear.",
    "price": 89.99,
    "images": [
        "https://example.com/images/sneakers-1.jpg",
        "https://example.com/images/sneakers-2.jpg"
    ],
    "categoryId": "67827cb87c8364072523fb55",
    "brand": "StepUp"
}'

curl -X POST http://localhost:8000/products \
-H "Content-Type: application/json" \
-d '{
    "name": "Wool Overcoat",
    "description": "Elegant wool overcoat to keep you warm and stylish.",
    "price": 249.99,
    "images": [
        "https://example.com/images/overcoat-1.jpg",
        "https://example.com/images/overcoat-2.jpg"
    ],
    "categoryId": "67827cb87c8364072523fb55",
    "brand": "LuxuryWear"
}'

curl -X POST http://localhost:8000/products \
-H "Content-Type: application/json" \
-d '{
    "name": "Silk Scarf",
    "description": "Luxurious silk scarf with elegant patterns.",
    "price": 39.99,
    "images": [
        "https://example.com/images/scarf-1.jpg",
        "https://example.com/images/scarf-2.jpg"
    ],
    "categoryId": "67827cb87c8364072523fb55",
    "brand": "Elegance"
}'

curl -X POST http://localhost:8000/products \
-H "Content-Type: application/json" \
-d '{
    "name": "Casual Baseball Cap",
    "description": "Adjustable baseball cap for everyday casual wear.",
    "price": 14.99,
    "images": [
        "https://example.com/images/cap-1.jpg",
        "https://example.com/images/cap-2.jpg"
    ],
    "categoryId": "67827cb87c8364072523fb55",
    "brand": "HeadStyle"
}'

curl -X POST http://localhost:8000/products \
-H "Content-Type: application/json" \
-d '{
    "name": "Leather Belt",
    "description": "Durable leather belt with a classic buckle design.",
    "price": 29.99,
    "images": [
        "https://example.com/images/belt-1.jpg",
        "https://example.com/images/belt-2.jpg"
    ],
    "categoryId": "67827cb87c8364072523fb55",
    "brand": "Accessorize"
}'

curl -X POST http://localhost:8000/products \
-H "Content-Type: application/json" \
-d '{
    "name": "Chiffon Maxi Dress",
    "description": "Elegant and flowy chiffon maxi dress perfect for events.",
    "price": 99.99,
    "images": [
        "https://example.com/images/dress-1.jpg",
        "https://example.com/images/dress-2.jpg"
    ],
    "categoryId": "67827cb87c8364072523fb55",
    "brand": "GlamourWear"
}'

curl -X POST http://localhost:8000/products \
-H "Content-Type: application/json" \
-d '{
    "name": "Ankle Boots",
    "description": "Stylish ankle boots with a sturdy sole for all-day comfort.",
    "price": 79.99,
    "images": [
        "https://example.com/images/boots-1.jpg",
        "https://example.com/images/boots-2.jpg"
    ],
    "categoryId": "67827cb87c8364072523fb55",
    "brand": "UrbanWalk"
}'
