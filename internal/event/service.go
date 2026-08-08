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
		EventImageURL:    req.EventImageURL,
		EventImageId:     req.EventImageId,
	}

	if err := s.repo.Create(event); err != nil {
		return nil, err
	}
	return event.ToRawResponse(), nil
}

func (s *Service) GetAllEvents(params query.Params, lang string) (*httpresponse.PaginatedData, error) {
	events, total, err := s.repo.GetAll(params)
	if err != nil {
		return nil, err
	}

	docs := make([]*dto.Response, 0, len(events))
	for _, e := range events {
		docs = append(docs, e.ToResponse(lang))
	}

	meta := httpresponse.BuildPaginationMeta(total, params.Page, params.Limit)
	return &httpresponse.PaginatedData{
		Docs:           docs,
		PaginationMeta: meta,
	}, nil
}