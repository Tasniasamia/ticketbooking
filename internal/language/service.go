package language

import (
	"errors"
	"ticketBooking/internal/language/dto"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(req *dto.CreateRequest) (*dto.Response, error) {
	existing, _ := s.repo.GetByCode(req.Code)
	if existing != nil {
		return nil, errors.New("language code already exists")
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	if req.IsDefault {
		_ = s.repo.ClearDefault()
	}

	lang := &Language{
		Name:      req.Name,
		Code:      req.Code,
		RTL:       req.RTL,
		IsActive:  isActive,
		IsDefault: req.IsDefault,
		Flag:      req.Flag,
	}

	if err := s.repo.Create(lang); err != nil {
		return nil, err
	}

	return toResponse(lang), nil
}

func (s *Service) GetAll() ([]*dto.Response, error) {
	list, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	result := make([]*dto.Response, 0, len(list))
	for _, l := range list {
		result = append(result, toResponse(l))
	}
	return result, nil
}

func (s *Service) GetByID(id uint) (*dto.Response, error) {
	lang, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return toResponse(lang), nil
}

func (s *Service) Update(id uint, req *dto.UpdateRequest) (*dto.Response, error) {
	lang, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		lang.Name = req.Name
	}
	if req.Code != "" {
		lang.Code = req.Code
	}
	if req.RTL != nil {
		lang.RTL = *req.RTL
	}
	if req.IsActive != nil {
		lang.IsActive = *req.IsActive
	}
	if req.Flag != "" {
		lang.Flag = req.Flag
	}
	if req.IsDefault != nil && *req.IsDefault {
		_ = s.repo.ClearDefault()
		lang.IsDefault = true
	}

	if err := s.repo.Update(lang); err != nil {
		return nil, err
	}
	return toResponse(lang), nil
}

func (s *Service) Delete(id uint) error {
	lang, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if lang.IsDefault {
		return errors.New("cannot delete default language")
	}
	return s.repo.Delete(id)
}

func toResponse(l *Language) *dto.Response {
	return &dto.Response{
		ID:        l.ID,
		Name:      l.Name,
		Code:      l.Code,
		RTL:       l.RTL,
		IsActive:  l.IsActive,
		IsDefault: l.IsDefault,
		Flag:      l.Flag,
	}
}