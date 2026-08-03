package event

import (
	"ticketBooking/internal/event/dto"
	"ticketBooking/internal/httpresponse"
	"ticketBooking/internal/utils/query"
)

type Service struct {
	repo Repository
}

func NewEventService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateEvent(req *dto.CreateRequest) (*dto.Response, error) {
     event := &Event{
		Title:            req.Title,
		Description:      req.Description,
		Location:         req.Location,
		StartsAt:         req.StartsAt,
		TotalTickets:     req.TotalTickets,
		AvailableTickets: req.TotalTickets,
		Price:            req.Price,
	}

    err := s.repo.Create(event)
	if err != nil {
		return nil, err
	}

	return event.ToResponse(), nil
}


func (s *Service) GetAllEvents(params query.Params) (*httpresponse.PaginatedData, error) {
	events, total, err := s.repo.GetAll(params)
	if err != nil {
		return nil, err
	}

	meta := httpresponse.BuildPaginationMeta(total, params.Page, params.Limit)

	return &httpresponse.PaginatedData{
		Docs:           events,
		PaginationMeta: meta,
	}, nil
}