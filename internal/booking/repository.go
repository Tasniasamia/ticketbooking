package booking

import (
	"context"
	"errors"
	"ticketBooking/internal/booking/dto"
	"ticketBooking/internal/live"
	"ticketBooking/internal/utils/query"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// HoldTTL = pending booking er jonno ticket koto khon hold thakbe.
// Eta bdlate chaile shudhu ei ek line bodlan (jemon 2 * time.Minute).
const HoldTTL = 5 * time.Minute

// StripeMinExpiry = Stripe checkout session er minimum expiry 30 minute.
// Er kom dile Stripe error dey, tai eta HoldTTL theke alada rakha hoyeche.
const StripeMinExpiry = 30 * time.Minute

var (
	ErrBookingNotFound         = errors.New("booking not found")
	ErrNotEnoughTickets        = errors.New("not enough tickets available")
	ErrBookingAlreadyCancelled = errors.New("booking already cancelled")
	ErrForbiddenBookingAccess  = errors.New("you do not own this booking")
)

type Repository interface {
	Create(booking *Booking) error
	CreateWithReservation(booking *Booking) error
	ReleaseIfPending(bookingID uint) (bool, error)
	ConfirmOrRevive(bookingID uint) error
	ListExpiredPending(before time.Time) ([]Booking, error)
	GetByID(bookingID uint) (*Booking, error)
	GetByUserID(userID uint) ([]Booking, error)
	GetByManagerID(managerID uint) ([]Booking, error)
	GetAll(p query.Params) ([]*Booking, int64, error)
	Update(booking *Booking) error
	UpdateStatus(bookingID uint, status dto.BookingStatus) error
}

type repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository { return &repository{db: db} }

func (r *repository) Create(booking *Booking) error {
	if err := r.db.Create(booking).Error; err != nil {
		return err
	}
	return r.db.
		Preload("UserInfo").
		Preload("EventInfo").
		Preload("EventInfo.Manager").
		Preload("EventInfo.Category").
		First(booking, booking.ID).Error
}

// CreateWithReservation: ticket decrement + booking create ek-i transaction e.
// Majh pothe crash korleo ticket leak hobe na.
func (r *repository) CreateWithReservation(b *Booking) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		res := tx.Exec(
			`UPDATE events SET available_tickets = available_tickets - ? WHERE id = ? AND available_tickets >= ?`,
			b.Quantity, b.EventID, b.Quantity,
		)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return ErrNotEnoughTickets
		}
		return tx.Create(b).Error
	})
	if err == nil {
		live.Default.Publish(b.EventID) // commit hoyar por: browser ke janiye dao
	}
	return err
}

// ReleaseIfPending: sob rollback er EKMATRO path.
// Booking pending thakle-i cancelled kore ar ticket ferot dey.
// Idempotent: webhook, redirect verify, background worker joto bar-i call korun,
// ticket ek bar-i increment hobe.
func (r *repository) ReleaseIfPending(bookingID uint) (bool, error) {
	released := false
	var eventID uint
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var b Booking
		if err := tx.Select("id", "event_id", "quantity").First(&b, bookingID).Error; err != nil {
			return err
		}
		eventID = b.EventID

		res := tx.Model(&Booking{}).
			Where("id = ? AND status = ?", bookingID, dto.BookingPending).
			Update("status", dto.BookingCancelled)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return nil // age-i release/confirm hoye geche
		}
		released = true

		// ei booking er pending payment gulo o cancel
		if err := tx.Exec(
			`UPDATE payments SET status = 'cancelled' WHERE booking_id = ? AND status = 'pending'`,
			bookingID,
		).Error; err != nil {
			return err
		}

		return tx.Exec(
			`UPDATE events SET available_tickets = available_tickets + ? WHERE id = ?`,
			b.Quantity, b.EventID,
		).Error
	})
	if err == nil && released {
		live.Default.Publish(eventID) // ticket ferot gelo: browser ke janiye dao
	}
	return released, err
}

// ConfirmOrRevive: payment success hole booking confirm kore.
// Jodi majhe expire hoye cancelled hoye thake (late payment, jemon SSLCommerz),
// abar ticket reserve kore confirm kore. Ticket na thakle ErrNotEnoughTickets (refund lagbe).
func (r *repository) ConfirmOrRevive(bookingID uint) error {
	revived := false
	var eventID uint
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var b Booking
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&b, bookingID).Error; err != nil {
			return err
		}

		switch b.Status {
		case dto.BookingConfirmed:
			return nil
		case dto.BookingCancelled:
			res := tx.Exec(
				`UPDATE events SET available_tickets = available_tickets - ? WHERE id = ? AND available_tickets >= ?`,
				b.Quantity, b.EventID, b.Quantity,
			)
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected == 0 {
				return ErrNotEnoughTickets
			}
			revived = true
			eventID = b.EventID
		}

		return tx.Model(&Booking{}).Where("id = ?", bookingID).
			Update("status", dto.BookingConfirmed).Error
	})
	if err == nil && revived {
		live.Default.Publish(eventID) // abar ticket atkalo: browser ke janiye dao
	}
	return err
}

func (r *repository) ListExpiredPending(before time.Time) ([]Booking, error) {
	var list []Booking
	err := r.db.
		Where("status = ? AND created_at < ?", dto.BookingPending, before).
		Limit(200).
		Find(&list).Error
	return list, err
}

// StartExpiryWorker: abandoned (payment kora hoyni) booking er ticket ferot dey.
// Multiple instance cholleo safe, karon ReleaseIfPending conditional update.
func StartExpiryWorker(ctx context.Context, repo Repository) {
	go func() {
		t := time.NewTicker(time.Minute)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				// 2 min grace: Stripe webhook er delay er jonno
				list, err := repo.ListExpiredPending(time.Now().Add(-HoldTTL - 2*time.Minute))
				if err != nil {
					continue
				}
				for _, b := range list {
					_, _ = repo.ReleaseIfPending(b.ID)
				}
			}
		}
	}()
}

func (r *repository) GetByID(bookingID uint) (*Booking, error) {
	var b Booking
	err := r.db.First(&b, bookingID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrBookingNotFound
	}
	return &b, err
}

func (r *repository) GetByUserID(userID uint) ([]Booking, error) {
	var bookings []Booking
	err := r.db.Where("user_id = ?", userID).Preload("EventInfo").
		Preload("EventInfo.Manager").
		Preload("EventInfo.Category").
		Preload("UserInfo").Order("created_at DESC").Find(&bookings).Error
	return bookings, err
}

func (r *repository) GetByManagerID(managerID uint) ([]Booking, error) {
	var bookings []Booking
	err := r.db.
		Joins("EventInfo").
		Where("\"EventInfo\".manager_id = ?", managerID).
		Preload("EventInfo").
		Preload("EventInfo.Manager").
		Preload("EventInfo.Category").
		Preload("UserInfo").
		Find(&bookings).Error

	return bookings, err
}

func (r *repository) Update(booking *Booking) error {
	return r.db.Save(booking).Error
}

func (r *repository) UpdateStatus(bookingID uint, status dto.BookingStatus) error {
	return r.db.Model(&Booking{}).Where("id = ?", bookingID).Update("status", status).Error
}

func (r *repository) GetAll(p query.Params) ([]*Booking, int64, error) {
	var bookings []*Booking
	var total int64

	db := r.db.Model(&Booking{})
	db = query.Apply(db, p, []string{"booking_code", "created_at", "status", "total_price"}, nil, nil, nil, nil)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = query.Paginate(db, p)

	err := db.
		Preload("EventInfo").
		Preload("EventInfo.Manager").
		Preload("EventInfo.Category").
		Preload("UserInfo").
		Find(&bookings).Error

	if err != nil {
		return nil, 0, err
	}

	return bookings, total, nil
}
