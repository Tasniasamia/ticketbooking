package payment

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrPaymentNotFound       = errors.New("payment not found")
	ErrPaymentMethodNotFound = errors.New("payment method not found")
	ErrPaymentMethodExists   = errors.New("payment method with this code already exists")
)

type Repository interface {
	// Payments
	Create(p *Payment) error
	Update(p *Payment) error
	GetByID(id uint) (*Payment, error)
	GetByTransactionID(tranID string) (*Payment, error)
	GetBySessionID(sessionID string) (*Payment, error)
	GetByBookingID(bookingID uint) (*Payment, error)
	GetByUserID(userID uint) ([]Payment, error)

	// Payment methods
	CreateMethod(m *PaymentMethod) error
	UpdateMethod(m *PaymentMethod) error
	DeleteMethod(id uint) error
	GetMethodByID(id uint) (*PaymentMethod, error)
	GetMethodByCode(code string) (*PaymentMethod, error)
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
	err := r.db.Where("transaction_id = ?", tranID).First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrPaymentNotFound
	}
	return &p, err
}

func (r *repository) GetBySessionID(sessionID string) (*Payment, error) {
	var p Payment
	err := r.db.Where("gateway_session_id = ?", sessionID).First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrPaymentNotFound
	}
	return &p, err
}

func (r *repository) GetByBookingID(bookingID uint) (*Payment, error) {
	var p Payment
	err := r.db.Where("booking_id = ?", bookingID).First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrPaymentNotFound
	}
	return &p, err
}

func (r *repository) GetByUserID(userID uint) ([]Payment, error) {
	var list []Payment
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&list).Error
	return list, err
}

// ---- Payment methods ----

func (r *repository) CreateMethod(m *PaymentMethod) error {
	return r.db.Create(m).Error
}

func (r *repository) UpdateMethod(m *PaymentMethod) error {
	return r.db.Save(m).Error
}

func (r *repository) DeleteMethod(id uint) error {
	res := r.db.Delete(&PaymentMethod{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrPaymentMethodNotFound
	}
	return nil
}

func (r *repository) GetMethodByID(id uint) (*PaymentMethod, error) {
	var m PaymentMethod
	err := r.db.First(&m, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrPaymentMethodNotFound
	}
	return &m, err
}

func (r *repository) GetMethodByCode(code string) (*PaymentMethod, error) {
	var m PaymentMethod
	err := r.db.Where("code = ?", code).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrPaymentMethodNotFound
	}
	return &m, err
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
