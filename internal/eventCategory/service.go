package eventCategory;

import (
	"ticketBooking/internal/eventCategory/dto"
	"ticketBooking/internal/httpresponse"
	"ticketBooking/internal/utils/i18n"
	"ticketBooking/internal/utils/query"
)

type Service struct {
	repo Repository
}

func NewEventCategoryService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateEventCategory(req *dto.CreateEventCategoryRequest) (*dto.RawResponse, error) {
	eventCategory := &EventCategory{
		Name:            i18n.LocalizedString(req.Name),
		Description:      i18n.LocalizedString(req.Description),
		EventCategoryImageURL:    req.EventCategoryImageURL,
		EventCategoryImageId:     req.EventCategoryImageId,

	}

	if err := s.repo.Create(eventCategory); err != nil {
		return nil, err
	}
	return eventCategory.ToRawResponse(), nil
}

func (s *Service) GetAllEventCategories(params query.Params, lang string) (*httpresponse.PaginatedData, error) {
	eventCategories, total, err := s.repo.GetAll(params,lang)
	if err != nil {
		return nil, err
	}

	docs := make([]*dto.Response, 0, len(eventCategories))
	for _, e := range eventCategories {
		docs = append(docs, e.ToResponse(lang))
	}

	meta := httpresponse.BuildPaginationMeta(total, params.Page, params.Limit)
	return &httpresponse.PaginatedData{
		Docs:           docs,
		PaginationMeta: meta,
	}, nil
}


func (s *Service) GetEventCategoryByID(eventId uint, lang string) (*dto.Response, error) {
	eventCategory, err := s.repo.GetByID(eventId)
	if err != nil {
		return nil, err
	}
	return eventCategory.ToResponse(lang), nil
}

func (s *Service) UpdateEventCategory(eventId uint, req *dto.UpdateEventCategoryRequest) (*dto.RawResponse, error) {
	eventCategory, err := s.repo.GetByID(eventId)
	if err != nil {
		return nil, err
	}

	eventCategory.Name = i18n.LocalizedString(req.Name)
	eventCategory.Description = i18n.LocalizedString(req.Description)
	eventCategory.EventCategoryImageURL = req.EventCategoryImageURL
	eventCategory.EventCategoryImageId = req.EventCategoryImageId

	if err := s.repo.Update(eventCategory); err != nil {
		return nil, err
	}

	return eventCategory.ToRawResponse(), nil
}

func (s *Service) DeleteEventCategory(eventId uint) error {
	return s.repo.Delete(eventId)
}