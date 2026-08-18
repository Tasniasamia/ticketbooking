package event

import (
	"ticketBooking/internal/event/dto"
	"ticketBooking/internal/httpresponse"
	"ticketBooking/internal/utils/i18n"
	"ticketBooking/internal/utils/query"
)

type Service struct {
	repo Repository
}

func NewEventService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateEvent(req *dto.CreateRequest) (*dto.RawResponse, error) {
	event := &Event{
		Title:            i18n.LocalizedString(req.Title),
		Description:      i18n.LocalizedString(req.Description),
		Location:         i18n.LocalizedString(req.Location),
		StartsAt:         req.StartsAt,
		TotalTickets:     req.TotalTickets,
		AvailableTickets: req.TotalTickets,
		Price:            req.Price,
		ThumbnailImage:    req.ThumbnailImage,
		Images:             req.Images,
		ManagerID:          req.ManagerID,
		CategoryID:         req.CategoryID,
	}

	if err := s.repo.Create(event); err != nil {
		return nil, err
	}
	return event.ToRawResponse(), nil
}

func (s *Service) GetAllEvents(params query.Params, lang string) (*httpresponse.PaginatedData, error) {
	events, total, err := s.repo.GetAll(params,lang)
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

func (s *Service) GetEventByID(eventId uint, lang string) (*dto.RawResponse, error) {
	event, err := s.repo.GetByID(eventId)
	if err != nil {
		return nil, err
	}
	return event.ToRawResponse(), nil
}

func (s *Service) UpdateEvent(eventId uint, req *dto.UpdateRequest) (*dto.RawResponse, error) {
	event, err := s.repo.GetByID(eventId)
	if err != nil {
		return nil, err
	}

	// Update the event fields with the values from the request
	event.Title = i18n.LocalizedString(req.Title)
	event.Description = i18n.LocalizedString(req.Description)
	event.Location = i18n.LocalizedString(req.Location)
	event.StartsAt = req.StartsAt
	event.Price = req.Price
	event.ThumbnailImage = req.ThumbnailImage
	event.Images = req.Images

	if err := s.repo.Update(event); err != nil {
		return nil, err
	}

	return event.ToRawResponse(), nil
}

func (s *Service) DeleteEvent(eventId uint) error {
	event, err := s.repo.GetByID(eventId)
	if err != nil {
		return err
	}

	return s.repo.Delete(event)
}



