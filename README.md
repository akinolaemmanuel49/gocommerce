# E-Commerce Platform

---

## **Todo**
- [x] Implement all CRUD routes and expose them
- [x] Update deletion logic. Add "$set" to all queries.
- [x] Document implemented and planned API routes using POSTMAN
    - [x] Orders
    - [x] Users
    - [x] Products
    - [x] Categories
    - [x] Carts
    - [x] Auth
    - [x] Health
- [x] Implement carts repository.
- [x] Implement carts service.
- [x] Implement carts handler.
- [x] Refactor orders feature.
- [x] Implement authentication feature (JWT Authentication).
- [x] Implement auto set shipping address to owner address if present.
- [x] Implement messaging for orders, whenever the status of an order changes, a message is sent to the queue.
- [ ] Add route protection to concerned routes.
    - [x] Users
    - [x] Orders
    - [x] Products
    - [x] Categories
    - [x] Carts
- [x] Document implemented and planned API routes using SWAGGER
    - [x] Orders
    - [x] Users
    - [x] Products
    - [x] Categories
    - [x] Carts
    - [x] Auth
    - [x] Health
- [ ] Write tests
    - [x] Users
    - [ ] Products
    - [ ] Categories
    - [ ] Carts
    - [ ] Auth
- [ ] Implement email messaging for order creation, status changes and cancellation.

---

## **Features**

- [ ] **Product Catalog**: Search, filter, and browse through products.
- [ ] **User Accounts**: User registration, login, and profile management.
- [ ] **Shopping Cart**: Add, remove, and modify items in a cart.
- [ ] **Order Management**: Place orders, track status, and view order history.
- [ ] **Admin Dashboard**: Manage inventory and view analytics.
- [ ] **Payment Integration**: Seamlessly integrate popular gateways like Stripe or PayPal.
- [ ] **Notifications**: Real-time or scheduled notifications for order updates.

---

## **Tech Stack**

- **Backend**: Go
- **Database**: MongoDB
- **Task Queue**: RabbitMQ (for asynchronous task processing)
