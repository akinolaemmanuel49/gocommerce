package services

import (
	"context"
	"time"

	"github.com/akinolaemmanuel49/gocommerce/common/errors"
	"github.com/akinolaemmanuel49/gocommerce/internal/models"
	"github.com/akinolaemmanuel49/gocommerce/internal/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewCartService creates a new instance of CartService
func NewCartService(cartRepository *repositories.CartRepository, userService UserService, productService ProductService) *CartService {
	return &CartService{
		cartRepository: cartRepository,
	}
}

// CreateCart creates a new instance of a cart in the database
func (s *CartService) CreateCart(ctx context.Context, newCart *models.CreateCart, userID string) (*models.Cart, error) {
	// Check for valid user
	_, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Create new cart
	cart := models.NewCart(newCart, userID)

	// Save cart to database
	result, err := s.cartRepository.Insert(ctx, cart)
	if err != nil {
		return nil, err
	}

	// Cast result to models.Cart
	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, err
	}
	cart.ID = objectID.Hex() // Convert ObjectID to string

	return cart, nil
}

// AddProductToCart adds a product to a cart
func (s *CartService) AddProductToCart(ctx context.Context, userID string, cartID string, productID string, quantity int) (*models.Cart, error) {
	if quantity <= 0 {
		return nil, errors.NewValidationError("Quantity", "quantity must be greater than zero")
	}

	// Retrieve cart by ID
	cart, err := s.cartRepository.FindByID(ctx, cartID)
	if err != nil {
		return nil, err
	}

	if userID != cart.UserID {
		return nil, errors.NewForbiddenError("You do not have permission to modify this resource")
	}

	// Retrieve product by ID
	product, err := s.productRepository.FindByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	// Check if product exists in the cart
	productFound := false
	for i, item := range cart.Items {
		if item.Product.ID == productID {
			cart.Items[i].Quantity += quantity
			cart.TotalPrice += product.Price * float64(quantity)
			productFound = true
			break
		}
	}

	// If product not found, add it as a new item
	if !productFound {
		cart.Items = append(cart.Items, models.CartItem{
			Product:  *product,
			Quantity: quantity,
		})
		cart.TotalPrice += product.Price * float64(quantity)
	}

	// Update cart in the database
	updateFields := bson.M{
		"$set": bson.M{
			"items":      cart.Items,
			"totalPrice": cart.TotalPrice,
		},
	}
	_, err = s.cartRepository.Update(ctx, cart.ID, updateFields)
	if err != nil {
		return nil, err
	}

	return cart, nil
}

// RemoveProductFromCart removes a product or reduces its quantity in the cart
func (s *CartService) RemoveProductFromCart(ctx context.Context, userID string, cartID string, productID string, quantity int) (*models.Cart, error) {
	if quantity <= 0 {
		return nil, errors.NewValidationError("Quantity", "Quantity must be greater than zero")
	}

	// Retrieve the cart by ID
	cart, err := s.cartRepository.FindByID(ctx, cartID)
	if err != nil {
		return nil, err
	}

	if userID != cart.UserID {
		return nil, errors.NewForbiddenError("You do not have permission to modify this resource")
	}
	// Check if the product exists in the cart
	productFound := false
	for i, item := range cart.Items {
		if item.Product.ID == productID {
			productFound = true

			// Reduce the quantity or remove the product
			if item.Quantity <= quantity {
				// Remove the product entirely
				cart.TotalPrice -= item.Product.Price * float64(item.Quantity)
				cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			} else {
				// Reduce the quantity
				cart.Items[i].Quantity -= quantity
				cart.TotalPrice -= item.Product.Price * float64(quantity)
			}
			break
		}
	}

	// If the product is not found, return an error
	if !productFound {
		return nil, errors.NewNotFoundError("Cart", "Product", "Product not found in the cart")
	}

	// Update the cart in the database
	updateFields := bson.M{
		"$set": bson.M{
			"items":      cart.Items,
			"totalPrice": cart.TotalPrice,
		},
	}
	_, err = s.cartRepository.Update(ctx, cart.ID, updateFields)
	if err != nil {
		return nil, err
	}

	return cart, nil
}

// RetrieveCartByID retrieves a cart by its ID
func (s *CartService) RetrieveCartByID(ctx context.Context, cartID string) (*models.Cart, error) {
	// Retrieve cart by ID
	cart, err := s.cartRepository.FindByID(ctx, cartID)
	if err != nil {
		return nil, err
	}

	return cart, nil
}

// RetrieveAllCarts retrieves all carts based on filters and implements cursor-based pagination
func (s *CartService) RetrieveAllCarts(ctx context.Context, filter map[string]interface{}, lastID string, limit int) ([]models.Cart, string, error) {
	carts, nextCursor, err := s.cartRepository.FindAll(ctx, filter, lastID, limit)
	if err != nil {
		return nil, "", err
	}

	return carts, nextCursor, nil
}

// DeleteCartByID deletes a cart
func (s *CartService) DeleteCartByID(ctx context.Context, userID string, ID string) error {
	// Check for existing cart
	existingCart, err := s.cartRepository.FindByID(ctx, ID)
	if userID != existingCart.UserID {
		return errors.NewForbiddenError("You do not have permission to modify this resource")
	}

	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	if err == mongo.ErrNoDocuments {
		return nil
	}

	if existingCart != nil {
		// Apply transformation, set cart IsDeleted field to true
		order := &models.Order{
			CommonFields: models.CommonFields{
				IsDeleted: true,
				UpdatedAt: time.Now(),
			},
		}

		deleted := bson.M{"$set": order}

		_, err = s.cartRepository.Update(ctx, ID, deleted)
		if err != nil {
			return err
		}
	}

	return nil
}
