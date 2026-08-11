package currency

import (
	"errors"
	"math"
	"strings"
	"ticketBooking/internal/httpresponse"
	"ticketBooking/internal/utils/query"
	"ticketBooking/internal/currency/dto"
)

type Service interface {
	// Convert amount from fromCode → toCode (single public convert API)
	Convert(amount float64, fromCode, toCode string) (float64, error)

	GetDefault() (*Currency, error)
	GetDefaultCode() string
	GetByCode(code string) (*Currency, error)
	ListAll(params query.Params) (*httpresponse.PaginatedData, error)
	ListEnabled(params query.Params) (*httpresponse.PaginatedData, error) 

	Create(req dto.CreateCurrencyRequest) (*Currency, error)
	Update(id uint, req dto.UpdateCurrencyRequest) (*Currency, error)
	SetDefault(code string) (*Currency, error)
	Delete(id uint) error

	SeedDefaults() error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// Convert: amount in fromCode → toCode
// Base pivot = BDT (RateToBDT)
// amount_in_BDT = amount * from.RateToBDT
// result        = amount_in_BDT / to.RateToBDT
func (s *service) Convert(amount float64, fromCode, toCode string) (float64, error) {
	fromCode = strings.ToUpper(strings.TrimSpace(fromCode))
	toCode = strings.ToUpper(strings.TrimSpace(toCode))

	if fromCode == toCode {
		return round2(amount), nil
	}

	from, err := s.repo.GetByCode(fromCode)
	if err != nil {
		return 0, err
	}
	to, err := s.repo.GetByCode(toCode)
	if err != nil {
		return 0, err
	}

	if from.RateToBDT <= 0 || to.RateToBDT <= 0 {
		return 0, errors.New("invalid exchange rate")
	}

	amountInBDT := amount * from.RateToBDT
	result := amountInBDT / to.RateToBDT
	return round2(result), nil
}

func (s *service) GetDefault() (*Currency, error) {
	return s.repo.GetDefault()
}

func (s *service) GetDefaultCode() string {
	c, err := s.repo.GetDefault()
	if err != nil {
		return "BDT"
	}
	return c.Code
}

func (s *service) GetByCode(code string) (*Currency, error) {
	return s.repo.GetByCodeAnyStatus(code)
}

// func (s *service) ListAll(params query.Params, lang string) (*httpresponse.PaginatedData, error) {
// 		currencies, total, err := s.repo.GetAll(params)
// 			if err != nil {
// 		return nil, err
// 	}

// 	docs := make([]*dto.Response, 0, len(currencies))
// 	for _, c := range currencies {
// 		docs = append(docs, c.ToResponse())
// 	}

// 	meta := httpresponse.BuildPaginationMeta(total, params.Page, params.Limit)
// 	return &httpresponse.PaginatedData{
// 		Docs:           docs,
// 		PaginationMeta: meta,
// 	}, nil

// }



func (s *service) ListAll(params query.Params) (*httpresponse.PaginatedData, error) {
	events, total, err := s.repo.GetAll(params)
	if err != nil {
		return nil, err
	}

	docs := make([]*dto.Response, 0, len(events))
	for _, e := range events {
		docs = append(docs, e.ToResponse())
	}

	meta := httpresponse.BuildPaginationMeta(total, params.Page, params.Limit)
	return &httpresponse.PaginatedData{
		Docs:           docs,
		PaginationMeta: meta,
	}, nil
}




func (s *service) ListEnabled(params query.Params) (*httpresponse.PaginatedData, error) {
	events, total, err := s.repo.GetAllEnabled(params)
	if err != nil {
		return nil, err
	}

	docs := make([]*dto.Response, 0, len(events))
	for _, e := range events {
		docs = append(docs, e.ToResponse())
	}

	meta := httpresponse.BuildPaginationMeta(total, params.Page, params.Limit)
	return &httpresponse.PaginatedData{
		Docs:           docs,
		PaginationMeta: meta,
	}, nil
}



// func (s *service) ListEnabled() ([]Currency, error) {
// 	return s.repo.GetAllEnabled()


// }






func (s *service) Create(req dto.CreateCurrencyRequest) (*Currency, error) {
	code := strings.ToUpper(strings.TrimSpace(req.Code))
	if code == "" {
		return nil, errors.New("code is required")
	}
	if req.RateToBDT <= 0 {
		return nil, errors.New("rate_to_bdt must be greater than 0")
	}

	// already exists?
	if _, err := s.repo.GetByCodeAnyStatus(code); err == nil {
		return nil, errors.New("currency already exists")
	}

	status := req.Status
	if status == "" {
		status = "enable"
	}

	c := &Currency{
		Code:      code,
		Name:      req.Name,
		Symbol:    req.Symbol,
		RateToBDT: req.RateToBDT,
		Status:    status,
		IsDefault: false,
	}

	if req.IsDefault {
		if err := s.repo.ClearDefault(); err != nil {
			return nil, err
		}
		c.IsDefault = true
	}

	if err := s.repo.Create(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *service) Update(id uint, req dto.UpdateCurrencyRequest) (*Currency, error) {
	c, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		c.Name = *req.Name
	}
	if req.Symbol != nil {
		c.Symbol = *req.Symbol
	}
	if req.RateToBDT != nil {
		if *req.RateToBDT <= 0 {
			return nil, errors.New("rate_to_bdt must be greater than 0")
		}
		c.RateToBDT = *req.RateToBDT
	}
	if req.Status != nil {
		c.Status = *req.Status
	}

	if err := s.repo.Update(c); err != nil {
		return nil, err
	}
	return c, nil
}

// SetDefault — single click: clear old default, set new one
func (s *service) SetDefault(code string) (*Currency, error) {
	c, err := s.repo.GetByCodeAnyStatus(code)
	if err != nil {
		return nil, err
	}
	if c.Status != "enable" {
		return nil, errors.New("cannot set disabled currency as default")
	}

	if err := s.repo.ClearDefault(); err != nil {
		return nil, err
	}

	c.IsDefault = true
	if err := s.repo.Update(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *service) Delete(id uint) error {
	c, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if c.IsDefault {
		return errors.New("cannot delete default currency; set another default first")
	}
	return s.repo.Delete(id)
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

// SeedDefaults — only inserts missing currencies, does not overwrite admin rates
func (s *service) SeedDefaults() error {
	defaults := []Currency{
		{Code: "BDT", Name: "Bangladeshi Taka", Symbol: "৳", RateToBDT: 1, IsDefault: true, Status: "enable"},
		{Code: "USD", Name: "US Dollar", Symbol: "$", RateToBDT: 120, IsDefault: false, Status: "enable"},
		{Code: "EUR", Name: "Euro", Symbol: "€", RateToBDT: 130, IsDefault: false, Status: "enable"},
		{Code: "INR", Name: "Indian Rupee", Symbol: "₹", RateToBDT: 1.45, IsDefault: false, Status: "enable"},
		{Code: "GBP", Name: "British Pound", Symbol: "£", RateToBDT: 155, IsDefault: false, Status: "enable"},
	}

	for i := range defaults {
		existing, err := s.repo.GetByCodeAnyStatus(defaults[i].Code)
		if err != nil {
			// not found → create
			if err := s.repo.Create(&defaults[i]); err != nil {
				return err
			}
			continue
		}
		// already exists → only update name/symbol if empty, never touch rate/default
		_ = existing
	}
	return nil
}
