package service

import (
	"errors"

	"github.com/Kevinmajesta/depublic-backend/internal/entity"
	"github.com/Kevinmajesta/depublic-backend/internal/repository"
	"github.com/google/uuid"
)

type WishlistService interface {
	GetAllWishlist() ([]entity.Wishlist, error)
	GetWishlistByUserId(UserId uuid.UUID) (*entity.Wishlist, error)
	AddWishlist(wishlist *entity.Wishlist) (*entity.Wishlist, error)
	RemoveWishlist(EventId, UserId uuid.UUID) (*entity.Wishlist, error)
}

type wishlistService struct {
	wishlistRepository  repository.WishlistRepository
	repo                repository.EventRepository
	userRepo            repository.UserRepository
	notificationService NotificationService
}

func NewWishlistService(wishlistRepository repository.WishlistRepository, repo repository.EventRepository, userRepo repository.UserRepository, notificationService NotificationService) WishlistService {
	return &wishlistService{wishlistRepository: wishlistRepository, repo: repo, userRepo: userRepo, notificationService: notificationService}
}

func (s *wishlistService) GetAllWishlist() ([]entity.Wishlist, error) {
	wishlists, err := s.wishlistRepository.GetAllWishlist()
	if err != nil {
		return nil, err
	}
	return wishlists, nil
}

func (s *wishlistService) GetWishlistByUserId(UserId uuid.UUID) (*entity.Wishlist, error) {
	eventAdd, err := s.wishlistRepository.CheckEventAdd(UserId)
	if err != nil {
		return nil, err
	}

	//check if user doesn't have an event in wishlist
	if eventAdd == false {
		return nil, errors.New("this user doesn't have any events in wishlist")
	}

	return s.wishlistRepository.GetWishlistByUserId(UserId)
}

func (s *wishlistService) AddWishlist(wishlist *entity.Wishlist) (*entity.Wishlist, error) {

	//check if user exists in db
	checkUser, err := s.userRepo.CheckUser(wishlist.UserId)
	if err != nil {
		return nil, err
	}

	if checkUser == nil {
		return nil, errors.New("users does not exist")
	}

	//check if event exists in db
	checkEvent, err := s.repo.CheckEvent(wishlist.EventId)
	if err != nil {
		return nil, err
	}

	if checkEvent == nil {
		return nil, errors.New("events does not exist")
	}

	// Check if events with the same event_id and user_id are already on the wishlist
	existingWishlist, err := s.wishlistRepository.GetWishlistByEventAndUser(wishlist.EventId, wishlist.UserId)
	if err != nil {
		return nil, err
	}

	if existingWishlist != nil {
		return nil, errors.New("event already added to wishlist")
	}

	notification := &entity.Notification{
		UserID:  wishlist.UserId,
		Type:    "Add To Whislist",
		Message: "Add To Whislist successful for event " + wishlist.EventId.String(),
		IsRead:  false,
	}
	err = s.notificationService.CreateNotification(notification)
	if err != nil {
		return nil, err
	}
	// Jika belum ada, tambahkan ke wishlist
	return s.wishlistRepository.AddWishlist(wishlist)
}

func (s *wishlistService) RemoveWishlist(EventId, UserId uuid.UUID) (*entity.Wishlist, error) {
	// Check if events with the same EventId and UserId are already on the wishlist
	existingWishlist, err := s.wishlistRepository.GetWishlistByEventAndUser(EventId, UserId)
	if err != nil {
		return nil, err
	}

	if existingWishlist == nil {
		return nil, errors.New("wishlist not found for the given event and user")
	}

	// Delete wishlist
	err = s.wishlistRepository.RemoveWishlist(EventId, UserId)
	if err != nil {
		return nil, err
	}

	return existingWishlist, nil
}
