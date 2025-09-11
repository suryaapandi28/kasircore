package service

import (
	"errors"
	"strconv"

	"github.com/google/uuid"
	"github.com/suryaapandi28/kasircore/internal/entity"
	"github.com/suryaapandi28/kasircore/internal/repository"
)

type CartService interface {
	GetAllCart() ([]entity.Carts, error)
	GetCartByUserId(UserId uuid.UUID) (*entity.Carts, error)
	AddToCart(UserId, EventId uuid.UUID) (*entity.Carts, error)
	RemoveCart(CartId uuid.UUID) (bool, error)
	UpdateQuantityAdd(UserId, EventId uuid.UUID) error
	UpdateQuantityLess(UserId, EventId uuid.UUID) error
}

type cartService struct {
	cartRepository      repository.CartRepository
	repo                repository.EventRepository
	userRepo            repository.UserRepository
	notificationService NotificationService
}

func NewCartService(cartRepository repository.CartRepository, repo repository.EventRepository, userRepo repository.UserRepository, notificationService NotificationService) CartService {
	return &cartService{cartRepository: cartRepository, repo: repo, userRepo: userRepo, notificationService: notificationService}
}

func (s *cartService) GetAllCart() ([]entity.Carts, error) {
	carts, err := s.cartRepository.GetAllCart()
	if err != nil {
		return nil, err
	}
	return carts, nil
}

func (s *cartService) GetCartByUserId(UserId uuid.UUID) (*entity.Carts, error) {

	eventAdd, err := s.cartRepository.CheckEventAdd(UserId)
	if err != nil {
		return nil, err
	}

	//check if user doesn't have an event in carts
	if eventAdd == false {
		return nil, errors.New("this user doesn't have any events in cart")
	}

	return s.cartRepository.GetCartByUserId(UserId)
}

func (s *cartService) AddToCart(UserId, EventId uuid.UUID) (*entity.Carts, error) {

	//check if user exists in db
	checkUser, err := s.userRepo.CheckUser(UserId)
	if err != nil {
		return nil, err
	}

	if checkUser == nil {
		return nil, errors.New("users does not exist")
	}

	//check if event exists in db
	checkEvent, err := s.repo.CheckEvent(EventId)
	if err != nil {
		return nil, err
	}

	if checkEvent == nil {
		return nil, errors.New("events does not exist")
	}

	//check if the user already has a quantity of an event in his cart.
	exist, err := s.cartRepository.CheckIfEventAlreadyAdded(UserId, EventId)
	if err != nil {
		return nil, err
	}

	//if it exist cannot add again event
	if exist {
		return nil, errors.New("you've already added this event!")
	}

	//check max add 1 event by user
	eventAdd, err := s.cartRepository.CheckEventAdd(UserId)
	if err != nil {
		return nil, err
	}

	if eventAdd {
		return nil, errors.New("max add of cart one event!")
	}

	// Check available quantity for the event
	qtyEvent, err := s.repo.CheckQtyEvent(EventId)
	if err != nil {
		return nil, err
	}

	// If event is out of stock
	if qtyEvent < 1 {
		return nil, errors.New("out of stock")
	}

	// Get the price of the event
	priceEvent, err := s.repo.CheckPriceEvent(EventId)
	if err != nil {
		return nil, err
	}

	// Get the date of the event
	dateEvent, err := s.repo.CheckDateEvent(EventId)
	if err != nil {
		return nil, err
	}

	// Calculate the total price
	pricetot := priceEvent * 1
	totalPrice := strconv.Itoa(pricetot)

	// Create the new cart entry
	cart := &entity.Carts{
		Cart_id:     uuid.New().String(),
		User_id:     UserId.String(),
		Event_id:    EventId.String(),
		Qty:         "1",
		Ticket_date: dateEvent,
		Price:       totalPrice,
		Auditable:   entity.NewAuditable(),
	}

	// Add the new cart entry to the repository
	err = s.cartRepository.CreateCart(cart)
	if err != nil {
		return nil, err
	}

	// Decrease the stock in the events table
	err = s.repo.DecreaseEventStock(EventId, 1)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(cart.User_id)
	if err != nil {
		return nil, err
	}

	notification := &entity.Notification{
		UserID:  userID,
		Type:    "Add To Cart",
		Message: "Add To Cart successful for event " + cart.Event_id,
		IsRead:  false,
	}
	err = s.notificationService.CreateNotification(notification)
	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (s *cartService) UpdateQuantityAdd(UserId, EventId uuid.UUID) error {

	//check if user exists in db
	checkUser, err := s.userRepo.CheckUser(UserId)
	if err != nil {
		return err
	}

	if checkUser == nil {
		return errors.New("users does not exist")
	}

	//check if event exists in db
	checkEvent, err := s.repo.CheckEvent(EventId)
	if err != nil {
		return err
	}

	if checkEvent == nil {
		return errors.New("events does not exist")
	}

	// Get the last quantity in the cart
	lastQtyCart, err := s.cartRepository.GetUserTotalQtyInCart(UserId, EventId)
	if err != nil {
		return err
	}

	// If the quantity in the cart is already 5 or more, return an error
	if lastQtyCart >= 5 {
		return errors.New("you cannot update quantity more than 5")
	}

	// Increase the quantity in the cart
	err = s.cartRepository.UpdateQuantityAdd(UserId, EventId)
	if err != nil {
		return err
	}

	// Get the price of the related event
	price, err := s.repo.CheckPriceEvent(EventId)
	if err != nil {
		return err
	}

	// Calculate the new total price after increasing the quantity
	newTotalPrice := price * int(lastQtyCart+1)

	// Update the total price in the cart
	err = s.cartRepository.UpdateTotalPrice(UserId, EventId, newTotalPrice)
	if err != nil {
		return err
	}

	// Decrease the stock of the event that was increased
	err = s.repo.DecreaseEventStock(EventId, 1)
	if err != nil {
		return err
	}

	return nil
}

func (s *cartService) UpdateQuantityLess(UserId, EventId uuid.UUID) error {

	//check if user exists in db
	checkUser, err := s.userRepo.CheckUser(UserId)
	if err != nil {
		return err
	}

	if checkUser == nil {
		return errors.New("users does not exist")
	}

	//check if event exists in db
	checkEvent, err := s.repo.CheckEvent(EventId)
	if err != nil {
		return err
	}

	if checkEvent == nil {
		return errors.New("events does not exist")
	}

	// retrieve the last quantity of users and event
	lastQtyCart, err := s.cartRepository.GetUserTotalQtyInCart(UserId, EventId)
	if err != nil {
		return err
	}

	// check if the last quantity is less than 1
	if lastQtyCart <= 1 {
		return errors.New("you cannot update quantity less than 1")
	}

	// subtract quantity from the cart
	err = s.cartRepository.UpdateQuantityLess(UserId, EventId)
	if err != nil {
		return err
	}

	// retrieve a price from containing the event
	price, err := s.repo.CheckPriceEvent(EventId)
	if err != nil {
		return err
	}

	// calculate the new amount after subtracting quantities
	newTotalPrice := price * int(lastQtyCart-1)

	// update amount of carts
	err = s.cartRepository.UpdateTotalPrice(UserId, EventId, newTotalPrice)
	if err != nil {
		return err
	}

	// adding back stock that has been reduced
	err = s.repo.IncreaseEventStock(EventId, 1)
	if err != nil {
		return err
	}

	return nil
}

func (s *cartService) RemoveCart(CartId uuid.UUID) (bool, error) {
	// Ambil informasi tentang cart berdasarkan cart_id
	cart, err := s.cartRepository.FindCartById(CartId)
	if err != nil {
		return false, err
	}

	// Hapus entri dari keranjang
	return s.cartRepository.RemoveCart(cart)
}
