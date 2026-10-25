// package event

// import (
// 	"errors"
// 	"strings"
// 	"ticketBooking/internal/event/dto"
// 	"ticketBooking/internal/httpresponse"
// 	"ticketBooking/internal/user"
// 	"ticketBooking/internal/utils/i18n"
// 	"ticketBooking/internal/utils/query"
// )

// type Service struct {
// 	repo Repository
// 	userRepo user.Repository
// }

// func NewEventService(repo Repository,userRepo user.Repository) *Service {
// 	return &Service{repo: repo, userRepo: userRepo}
// }

// func (s *Service) CreateEvent(req *dto.CreateRequest) (*dto.RawResponse, error) {
// 	currency := strings.ToUpper(strings.TrimSpace(req.Currency))
// 	if currency == "" {
// 		currency = "BDT"
// 	}

// 	event := &Event{
// 		Title:            i18n.LocalizedString(req.Title),
// 		Description:      i18n.LocalizedString(req.Description),
// 		Location:         i18n.LocalizedString(req.Location),
// 		StartsAt:         req.StartsAt,
// 		TotalTickets:     req.TotalTickets,
// 		AvailableTickets: req.TotalTickets,
// 		Price:            req.Price,
// 		Currency:         currency,
// 		ThumbnailImage:    req.ThumbnailImage,
// 		Images:             req.Images,
// 		ManagerID:          req.ManagerID,
// 		CategoryID:         req.CategoryID,
// 	}
   
//    if _,err :=s.userRepo.GetUserActiveById(req.ManagerID); err != nil {
//        return nil, errors.New("Invalid manager ID");
//    }

	
// 	if err := s.repo.Create(event); err != nil {
// 		return nil, err
// 	}
// 	return event.ToRawResponse(), nil
// }

// func (s *Service) GetAllEvents(params query.Params, lang string) (*httpresponse.PaginatedData, error) {
// 	events, total, err := s.repo.GetAll(params,lang)
// 	if err != nil {
// 		return nil, err
// 	}

// 	docs := make([]*dto.RawResponse, 0, len(events))
// 	for _, e := range events {
// 		docs = append(docs, e.ToRawResponse())
// 	}

// 	meta := httpresponse.BuildPaginationMeta(total, params.Page, params.Limit)
// 	return &httpresponse.PaginatedData{
// 		Docs:           docs,
// 		PaginationMeta: meta,
// 	}, nil
// }


// func (s *Service) GetEventByID(eventId uint, lang string) (*dto.RawResponse, error) {
// 	event, err := s.repo.GetByID(eventId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return event.ToRawResponse(), nil
// }
// func (s *Service) GetEventByIDAdmin(eventId uint) (*dto.RawResponse, error) {
// 	event, err := s.repo.GetByIDAdmin(eventId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return event.ToRawResponse(), nil
// }
// func (s *Service) UpdateEvent(eventId uint, req *dto.UpdateRequest) (*dto.RawResponse, error) {
// 	event, err := s.repo.GetByID(eventId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Update the event fields with the values from the request
// 	event.Title = i18n.LocalizedString(req.Title)
// 	event.Description = i18n.LocalizedString(req.Description)
// 	event.Location = i18n.LocalizedString(req.Location)
// 	event.StartsAt = req.StartsAt
// 	event.Price = req.Price
// 	if c := strings.ToUpper(strings.TrimSpace(req.Currency)); c != "" {
// 		event.Currency = c
// 	}
// 	event.ThumbnailImage = req.ThumbnailImage
// 	event.Images = req.Images

// 	if err := s.repo.Update(event); err != nil {
// 		return nil, err
// 	}

// 	return event.ToRawResponse(), nil
// }

// func (s *Service) DeleteEvent(eventId uint) error {
// 	event, err := s.repo.GetByID(eventId)
// 	if err != nil {
// 		return err
// 	}

// 	return s.repo.Delete(event)
// }


// func (s *Service) UpdateEventStatus(req *dto.UpdateEventStatusRequest) error {
// 	event, err := s.repo.GetByID(req.EventId)
// 	if err != nil {
// 		return err
// 	}

// 	event.Status = req.Status

// 	return s.repo.Update(event)
// }


package event

import (
	"errors"
	"strings"
	"ticketBooking/internal/event/dto"
	"ticketBooking/internal/httpresponse"
	"ticketBooking/internal/user"
	"ticketBooking/internal/utils/i18n"
	"ticketBooking/internal/utils/query"
)

// ErrForbidden: manager অন্যের event এ হাত দিতে চাইলে।
var ErrForbidden = errors.New("you can only manage your own events")

const roleAdmin = "admin"

type Service struct {
	repo     Repository
	userRepo user.Repository
}

func NewEventService(repo Repository, userRepo user.Repository) *Service {
	return &Service{repo: repo, userRepo: userRepo}
}

// CreateEvent: managerID আসে token থেকে (client এর পাঠানো কিছু বিশ্বাস করা হয় না)।
// নতুন event সবসময় pending, admin approve না করা পর্যন্ত public এ আসে না।
func (s *Service) CreateEvent(req *dto.CreateRequest, managerID uint) (*dto.RawResponse, error) {
	currency := strings.ToUpper(strings.TrimSpace(req.Currency))
	if currency == "" {
		currency = "BDT"
	}

	if _, err := s.userRepo.GetUserActiveById(managerID); err != nil {
		return nil, errors.New("Invalid manager ID")
	}

	event := &Event{
		Title:            i18n.LocalizedString(req.Title),
		Description:      i18n.LocalizedString(req.Description),
		Location:         i18n.LocalizedString(req.Location),
		StartsAt:         req.StartsAt,
		TotalTickets:     req.TotalTickets,
		AvailableTickets: req.TotalTickets,
		Price:            req.Price,
		Currency:         currency,
		ThumbnailImage:   req.ThumbnailImage,
		Images:           req.Images,
		ManagerID:        managerID,
		CategoryID:       req.CategoryID,
		Status:           dto.Pending,
	}

	if err := s.repo.Create(event); err != nil {
		return nil, err
	}
	return event.ToRawResponse(), nil
}

func (s *Service) GetAllEvents(params query.Params, lang string) (*httpresponse.PaginatedData, error) {
	events, total, err := s.repo.GetAll(params, lang)
	if err != nil {
		return nil, err
	}

	docs := make([]*dto.RawResponse, 0, len(events))
	for _, e := range events {
		docs = append(docs, e.ToRawResponse())
	}

	meta := httpresponse.BuildPaginationMeta(total, params.Page, params.Limit)
	return &httpresponse.PaginatedData{
		Docs:           docs,
		PaginationMeta: meta,
	}, nil
}

// Public: শুধু approved event।
func (s *Service) GetEventByID(eventId uint, lang string) (*dto.RawResponse, error) {
	event, err := s.repo.GetByID(eventId)
	if err != nil {
		return nil, err
	}
	return event.ToRawResponse(), nil
}

// Admin: যেকোনো status এর event।
func (s *Service) GetEventByIDAdmin(eventId uint) (*dto.RawResponse, error) {
	event, err := s.repo.GetByIDAdmin(eventId)
	if err != nil {
		return nil, err
	}
	return event.ToRawResponse(), nil
}

// owns: admin সবকিছু পারে, manager শুধু নিজের event।
func owns(event *Event, actorID uint, actorRole string) bool {
	return actorRole == roleAdmin || event.ManagerID == actorID
}

func (s *Service) UpdateEvent(eventId uint, req *dto.UpdateRequest, actorID uint, actorRole string) (*dto.RawResponse, error) {
	// GetByID শুধু approved আনে, তাই pending event edit করতে GetByIDAdmin লাগে।
	event, err := s.repo.GetByIDAdmin(eventId)
	if err != nil {
		return nil, err
	}
	if !owns(event, actorID, actorRole) {
		return nil, ErrForbidden
	}

	event.Title = i18n.LocalizedString(req.Title)
	event.Description = i18n.LocalizedString(req.Description)
	event.Location = i18n.LocalizedString(req.Location)
	event.StartsAt = req.StartsAt
	event.Price = req.Price
	if c := strings.ToUpper(strings.TrimSpace(req.Currency)); c != "" {
		event.Currency = c
	}
	if req.CategoryID != 0 {
		event.CategoryID = req.CategoryID
	}
	event.ThumbnailImage = req.ThumbnailImage
	event.Images = req.Images

	if err := s.repo.Update(event); err != nil {
		return nil, err
	}

	// নতুন category সহ ফেরত দিতে আবার আনি
	fresh, err := s.repo.GetByIDAdmin(eventId)
	if err != nil {
		return nil, err
	}
	return fresh.ToRawResponse(), nil
}

func (s *Service) DeleteEvent(eventId uint, actorID uint, actorRole string) error {
	event, err := s.repo.GetByIDAdmin(eventId)
	if err != nil {
		return err
	}
	if !owns(event, actorID, actorRole) {
		return ErrForbidden
	}
	return s.repo.Delete(event)
}

// UpdateEventStatus: শুধু admin (route এ AdminMiddleware)। GetByID এখানে চলবে না,
// কারণ সেটা approved ছাড়া পায় না, ফলে pending event কখনো approve হতো না।
func (s *Service) UpdateEventStatus(req *dto.UpdateEventStatusRequest) error {
	event, err := s.repo.GetByIDAdmin(req.EventId)
	if err != nil {
		return err
	}

	event.Status = req.Status
	return s.repo.Update(event)
}
