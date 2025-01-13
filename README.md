# E-Commerce Platform

---

## **Todo**
- [x] Implement all CRUD routes and expose them
- [ ] Update deletion logic. Add "$set" to all queries.
- [ ] Document implemented API routes using POSTMAN
    - [x] Orders
    - [ ] Users
    - [ ] Products
    - [ ] Categories
- [ ] Implement messaging for orders, whenever the status of an order changes, a message is sent to the queue, which is then sent to the email address of the user that owns the order.
- [ ] Write tests

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
