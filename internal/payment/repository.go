package payment

import (
	"errors"
	"strings"
	"ticketBooking/internal/utils/query"

	"gorm.io/gorm"
)

var (
	ErrPaymentNotFound       = errors.New("payment not found")
	ErrPaymentMethodNotFound = errors.New("payment method not found")
	ErrPaymentMethodExists   = errors.New("payment method with this code already exists")
)

// PaymentAmountRow — status + currency grouped sum for aggregation.
type PaymentAmountRow struct {
	Status   string  `json:"status"`
	Currency string  `json:"currency"`
	Total    float64 `json:"total"`
}

type Repository interface {
	// Payments
	Create(p *Payment) error
	Update(p *Payment) error
	GetByID(id uint) (*Payment, error)
	GetByTransactionID(tranID string) (*Payment, error)
	GetBySessionID(sessionID string) (*Payment, error)
	GetByBookingID(bookingID uint) (*Payment, error)
	GetByUserID(userID uint) ([]Payment, error)
	GetAllPayments(params query.Params) ([]*Payment, int64, error)
	GetManagerPayments(params query.Params, managerID uint) ([]*Payment, int64, error)
	GetUserPayments(params query.Params, userID uint) ([]*Payment, int64, error)

	AggregateAllPayments(params query.Params) ([]PaymentAmountRow, error)
	AggregateManagerPayments(params query.Params, managerID uint) ([]PaymentAmountRow, error)
	AggregateUserPayments(params query.Params, userID uint) ([]PaymentAmountRow, error)

	// Payment methods
	CreateMethod(m *PaymentMethod) error
	UpdateMethod(m *PaymentMethod) error
	DeleteMethod(id uint) error
	GetMethodByID(id uint) (PaymentMethod, error)
	GetMethodByCode(code string) (PaymentMethod, error)
	ListMethods(enabledOnly bool) ([]PaymentMethod, error)
}

type repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository { return &repository{db: db} }

func (r *repository) Create(p *Payment) error { return r.db.Create(p).Error }
func (r *repository) Update(p *Payment) error { return r.db.Save(p).Error }

func (r *repository) GetByID(id uint) (*Payment, error) {
	var p Payment
	err := r.db.First(&p, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrPaymentNotFound
	}
	return &p, err
}

func (r *repository) GetByTransactionID(tranID string) (*Payment, error) {
	var p Payment
	err := r.db.Where("transaction_id = ?", tranID).Order("created_at DESC").Preload("Bookings").Preload("UserInfo").Preload("EventInfo").Preload("EventInfo.Manager").Preload("EventInfo.Category").First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrPaymentNotFound
	}
	return &p, err
}

func (r *repository) GetBySessionID(sessionID string) (*Payment, error) {
	var p Payment
	err := r.db.Where("gateway_session_id = ?", sessionID).Order("created_at DESC").Preload("Bookings").Preload("UserInfo").Preload("EventInfo").Preload("EventInfo.Manager").Preload("EventInfo.Category").First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrPaymentNotFound
	}
	return &p, err
}

func (r *repository) GetByBookingID(bookingID uint) (*Payment, error) {
	var p Payment
	err := r.db.Where("booking_id = ?", bookingID).Order("created_at DESC").Preload("Bookings").Preload("UserInfo").Preload("EventInfo").Preload("EventInfo.Manager").Preload("EventInfo.Category").First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrPaymentNotFound
	}
	return &p, err
}

func (r *repository) GetByUserID(userID uint) ([]Payment, error) {
	var list []Payment
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Order("created_at DESC").Preload("Bookings").Preload("UserInfo").Preload("EventInfo").Preload("EventInfo.Manager").Preload("EventInfo.Category").Find(&list).Error
	return list, err
}

func (r *repository) GetAllPayments(params query.Params) ([]*Payment, int64, error) {
    var list  []*Payment
	var total int64

	db := r.db.Model(&Payment{})
	db = query.Apply(db, params,[]string{"transaction_id", "status", "created_at"},nil, nil)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = query.Paginate(db, params)

	err := db.Order("created_at DESC").Preload("Bookings").Preload("UserInfo").Preload("EventInfo").Preload("EventInfo.Manager").Preload("EventInfo.Category").Find(&list).Error

	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *repository) GetManagerPayments(params query.Params, managerID uint) ([]*Payment, int64, error){

    var list  []*Payment
	var total int64

	db := r.db.Model(&Payment{}). Joins("JOIN events ON events.id = payments.event_id").
    Where("events.manager_id = ?", managerID)
	
	db = query.Apply(db, params,[]string{"transaction_id", "status", "created_at"},nil, nil)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = query.Paginate(db, params)

	err := db.Preload("Bookings").Preload("UserInfo").Preload("EventInfo").Preload("EventInfo.Manager").Preload("EventInfo.Category").Find(&list).Error

	if err != nil {
		return nil, 0, err
	}

	return list, total, nil

}

func (r *repository) GetUserPayments(params query.Params, userID uint) ([]*Payment, int64, error){
	var list  []*Payment
	var total int64

	db := r.db.Model(&Payment{}).Where("user_id = ?", userID)
	db = query.Apply(db, params,[]string{"transaction_id", "status", "created_at"},nil, nil)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = query.Paginate(db, params)

	err := db.Order("created_at DESC").Preload("Bookings").Preload("UserInfo").Preload("EventInfo").Preload("EventInfo.Manager").Preload("EventInfo.Category").Find(&list).Error

	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func applyAggregateFilters(db *gorm.DB, params query.Params, searchFields []string) *gorm.DB {
	if params.Search != "" && len(searchFields) > 0 {
		var conditions []string
		var args []interface{}
		for _, field := range searchFields {
			conditions = append(conditions, field+" ILIKE ?")
			args = append(args, "%"+params.Search+"%")
		}
		if len(conditions) > 0 {
			db = db.Where(strings.Join(conditions, " OR "), args...)
		}
	}
	if len(params.Filters) > 0 {
		for col, val := range params.Filters {
			db = db.Where(col+" = ?", val)
		}
	}
	return db
}

func (r *repository) AggregateAllPayments(params query.Params) ([]PaymentAmountRow, error) {
	var rows []PaymentAmountRow
	db := r.db.Model(&Payment{}).
		Select("status, currency, COALESCE(SUM(amount), 0) as total")
	db = applyAggregateFilters(db, params, []string{"transaction_id", "status"})
	err := db.Group("status, currency").Scan(&rows).Error
	return rows, err
}

func (r *repository) AggregateManagerPayments(params query.Params, managerID uint) ([]PaymentAmountRow, error) {
	var rows []PaymentAmountRow
	db := r.db.Model(&Payment{}).
		Select("payments.status as status, payments.currency as currency, COALESCE(SUM(payments.amount), 0) as total").
		Joins("JOIN events ON events.id = payments.event_id").
		Where("events.manager_id = ?", managerID)
	db = applyAggregateFilters(db, params, []string{"payments.transaction_id", "payments.status"})
	err := db.Group("payments.status, payments.currency").Scan(&rows).Error
	return rows, err
}

func (r *repository) AggregateUserPayments(params query.Params, userID uint) ([]PaymentAmountRow, error) {
	var rows []PaymentAmountRow
	db := r.db.Model(&Payment{}).
		Select("status, currency, COALESCE(SUM(amount), 0) as total").
		Where("user_id = ?", userID)
	db = applyAggregateFilters(db, params, []string{"transaction_id", "status"})
	err := db.Group("status, currency").Scan(&rows).Error
	return rows, err
}

func (r *repository) CreateMethod(m *PaymentMethod) error {
	return r.db.Create(m).Error
}

func (r *repository) UpdateMethod(m *PaymentMethod) error {
	return r.db.Save(m).Error
}

func (r *repository) DeleteMethod(id uint) error {
	res := r.db.Unscoped().Delete(&PaymentMethod{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrPaymentMethodNotFound
	}
	return nil
}

func (r *repository) GetMethodByID(id uint) (PaymentMethod, error) {
	var m PaymentMethod
	err := r.db.First(&m, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return PaymentMethod{}, ErrPaymentMethodNotFound
	}
	return m, err
}

func (r *repository) GetMethodByCode(code string) (PaymentMethod, error) {
	var m PaymentMethod
	err := r.db.Where("code = ?", code).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return PaymentMethod{}, ErrPaymentMethodNotFound
	}
	return m, err
}

func (r *repository) ListMethods(enabledOnly bool) ([]PaymentMethod, error) {
	var list []PaymentMethod
	q := r.db.Order("id ASC")
	if enabledOnly {
		q = q.Where("enable = ?", true)
	}
	err := q.Find(&list).Error
	return list, err
}
