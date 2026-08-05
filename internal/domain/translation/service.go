package translation

import (
	"errors"

	"ticketBooking/internal/httpresponse"
	"ticketBooking/internal/domain/translation/dto"
	"ticketBooking/internal/utils/i18n"
	"ticketBooking/internal/utils/query"

	"gorm.io/gorm"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(req *dto.CreateRequest) (*dto.Response, error) {
	existing, err := s.repo.GetByKey(req.Key)
	if err == nil && existing != nil {
		return nil, errors.New("translation key already exists")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	values := i18n.LocalizedString{}
	if req.Values != nil {
		values = i18n.LocalizedString(req.Values)
	}

	t := &Translation{
		Key:    req.Key,
		Values: values,
	}

	if err := s.repo.Create(t); err != nil {
		return nil, err
	}
	return t.ToResponse(), nil
}

func (s *Service) UpdateByKey(key string, req *dto.UpdateRequest) (*dto.Response, error) {
	t, err := s.repo.GetByKey(key)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("translation key not found")
		}
		return nil, err
	}

	if t.Values == nil {
		t.Values = make(i18n.LocalizedString)
	}

	for lang, val := range req.Values {
		t.Values[lang] = val
	}

	if err := s.repo.Update(t); err != nil {
		return nil, err
	}
	return t.ToResponse(), nil
}

func (s *Service) BulkUpdate(req *dto.BulkUpdateRequest) ([]*dto.Response, error) {
	results := make([]*dto.Response, 0, len(req.Items))

	for _, item := range req.Items {
		t, err := s.repo.GetByKey(item.Key)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				values := i18n.LocalizedString(item.Values)
				t = &Translation{Key: item.Key, Values: values}
				if err := s.repo.Create(t); err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		} else {
			if t.Values == nil {
				t.Values = make(i18n.LocalizedString)
			}
			for lang, val := range item.Values {
				t.Values[lang] = val
			}
			if err := s.repo.Update(t); err != nil {
				return nil, err
			}
		}
		results = append(results, t.ToResponse())
	}
	return results, nil
}

func (s *Service) GetAll(params query.Params) (*httpresponse.PaginatedData, error) {
	list, total, err := s.repo.GetAll(params)
	if err != nil {
		return nil, err
	}

	docs := make([]*dto.Response, 0, len(list))
	for _, t := range list {
		docs = append(docs, t.ToResponse())
	}

	meta := httpresponse.BuildPaginationMeta(total, params.Page, params.Limit)
	return &httpresponse.PaginatedData{
		Docs:           docs,
		PaginationMeta: meta,
	}, nil
}

func (s *Service) GetGroup() (dto.GroupResponse, error) {
	list, err := s.repo.GetAllSimple()
	if err != nil {
		return nil, err
	}

	group := make(dto.GroupResponse)
	for _, t := range list {
		vals := make(map[string]string)
		for k, v := range t.Values {
			vals[k] = v
		}
		group[t.Key] = vals
	}
	return group, nil
}


func (s *Service) GetByKey(key string) (*dto.Response, error) {
	t, err := s.repo.GetByKey(key)
	if err != nil {
		return nil, err
	}
	return t.ToResponse(), nil
}

func (s *Service) DeleteByKey(key string) error {
	_, err := s.repo.GetByKey(key)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("translation key not found")
		}
		return err
	}
	return s.repo.DeleteByKey(key)
}