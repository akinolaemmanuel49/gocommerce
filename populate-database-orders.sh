curl -X POST http://localhost:8000/orders \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "67827cb67c8364072523fb4a",
    "items": [
      {
        "productId": "678286807c8364072523fb5c",
        "quantity": 2,
        "price": 15.99
      },
      {
        "productId": "678287157c8364072523fb64",
        "quantity": 1,
        "price": 299.99
      }
    ],
    "totalPrice": 331.97,
    "address": {
      "street": "123 Elm Street",
      "city": "Springfield",
      "state": "IL",
      "zip": "62701",
      "country": "USA"
    }
  }'

curl -X POST http://localhost:8000/orders \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "67827cb67c8364072523fb4a",
    "items": [
      {
        "productId": "678287157c8364072523fb63",
        "quantity": 1,
        "price": 999.99
      },
      {
        "productId": "678286807c8364072523fb60",
        "quantity": 3,
        "price": 16.99
      }
    ],
    "totalPrice": 1050.96,
    "address": {
      "street": "456 Oak Avenue",
      "city": "Metropolis",
      "state": "NY",
      "zip": "10001",
      "country": "USA"
    }
  }'

curl -X POST http://localhost:8000/orders \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "67827cb67c8364072523fb4a",
    "items": [
      {
        "productId": "678287157c8364072523fb65",
        "quantity": 1,
        "price": 1299.99
      },
      {
        "productId": "678286807c8364072523fb5d",
        "quantity": 2,
        "price": 17.99
      }
    ],
    "totalPrice": 1335.97,
    "address": {
      "street": "789 Pine Road",
      "city": "Gotham",
      "state": "NJ",
      "zip": "07001",
      "country": "USA"
    }
  }'

curl -X POST http://localhost:8000/orders \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "67827cb67c8364072523fb4c",
    "items": [
      {
        "productId": "678287157c8364072523fb66",
        "quantity": 1,
        "price": 1999.99
      },
      {
        "productId": "678286807c8364072523fb61",
        "quantity": 2,
        "price": 9.99
      }
    ],
    "totalPrice": 2020.97,
    "address": {
      "street": "321 Birch Boulevard",
      "city": "Smalltown",
      "state": "TX",
      "zip": "75001",
      "country": "USA"
    }
  }'

curl -X POST http://localhost:8000/orders \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "67827cb77c8364072523fb4b",
    "items": [
      {
        "productId": "678287157c8364072523fb67",
        "quantity": 2,
        "price": 199.99
      },
      {
        "productId": "678286807c8364072523fb62",
        "quantity": 1,
        "price": 21.99
      }
    ],
    "totalPrice": 421.97,
    "address": {
      "street": "654 Cedar Street",
      "city": "Lakeside",
      "state": "FL",
      "zip": "33101",
      "country": "USA"
    }
  }'
