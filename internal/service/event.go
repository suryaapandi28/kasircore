package service

import (
	"github.com/google/uuid"
	"github.com/suryaapandi28/kasircore/internal/entity"
	"github.com/suryaapandi28/kasircore/internal/repository"
)

type EventService interface {
	// TODO POST
	AddEvent(event *entity.Events) (*entity.Events, error)
	// TODO UPDATE
	UpdateEvent(event *entity.Events) (*entity.Events, error)
	// UpdateEventByID(eventID uuid.UUID, event *entity.Events) (*entity.Events, error)
	// TODO DELETE
	DeleteEventByID(eventID uuid.UUID) (*entity.Events, error)
	// TODO GET
	GetAllEvent() ([]entity.Events, error)
	GetEventByID(eventID uuid.UUID) (*entity.Events, error)
	SearchEventsByTitle(title string) ([]entity.Events, error)
	// TODO Filtering Events
	FilterEvents(
		categoryID uuid.UUID,
		startDate string,
		endDate string,
		cityEvent string,
		priceMin int,
		priceMax int,
	) ([]entity.Events, error)
	// TODO SORT
	SortEvents(sortBy string) ([]entity.Events, error)
}

type eventService struct {
	eventRepo repository.EventRepository
	// categoryService CategoryService
}

func NewEventService(eventRepo repository.EventRepository) EventService {
	return &eventService{eventRepo: eventRepo}
}

func (s *eventService) AddEvent(event *entity.Events) (*entity.Events, error) {
	return s.eventRepo.AddEvent(event)
}

func (s *eventService) UpdateEvent(event *entity.Events) (*entity.Events, error) {
	return s.eventRepo.UpdateEvent(event)
}

// UpdateEventByID updates an event by ID.
//
//	func (s *eventService) UpdateEventByID(eventID uuid.UUID, event *entity.Events) (*entity.Events, error) {
//		return s.eventRepo.UpdateEventByID(eventID, event)
//	}
func (s *eventService) DeleteEventByID(eventID uuid.UUID) (*entity.Events, error) {
	return s.eventRepo.DeleteEventByID(eventID)
}

func (s *eventService) GetAllEvent() ([]entity.Events, error) {
	return s.eventRepo.GetAllEvent()
}

func (s *eventService) GetEventByID(eventID uuid.UUID) (*entity.Events, error) {
	return s.eventRepo.GetEventByID(eventID)
}

func (s *eventService) SearchEventsByTitle(title string) ([]entity.Events, error) {
	return s.eventRepo.SearchByTitle(title)
}

// Filtering Event
func (s *eventService) FilterEvents(
	categoryID uuid.UUID,
	startDate string,
	endDate string,
	cityEvent string,
	priceMin int,
	priceMax int,
) ([]entity.Events, error) {
	return s.eventRepo.FilterEvents(categoryID, startDate, endDate, cityEvent, priceMin, priceMax)
}

// SORT EVENT
func (s *eventService) SortEvents(sortBy string) ([]entity.Events, error) {
	return s.eventRepo.SortEvents(sortBy)
}
