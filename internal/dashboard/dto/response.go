package dto

import (
	"time"

	"ticketBooking/internal/utils/i18n"
)

// BaseCurrency: dashboard এর সব টাকার অংক এই currency তে আসে (currencies.rate_to_bdt দিয়ে convert করা)।
// Frontend এর <Price from="BDT" /> user এর select করা currency তে বদলে দেখায়।
const BaseCurrency = "BDT"

// Stat: একটা সংখ্যা + আগের ২৮ দিনের তুলনায় কত % বদলেছে (আগে কিছু না থাকলে nil)।
type Stat struct {
	Value     float64  `json:"value"`
	ChangePct *float64 `json:"change_pct"`
}

type RevenuePoint struct {
	Week   string  `json:"week"` // YYYY-MM-DD (সোমবার)
	Amount float64 `json:"amount"`
}

// UpcomingEvent: admin আর manager দুজনেরই upcoming list এ এটা যায়।
type UpcomingEvent struct {
	ID             uint                 `json:"id"`
	Title          i18n.LocalizedString `json:"title"`
	Location       i18n.LocalizedString `json:"location"`
	ImageURL       string               `json:"image_url"`
	StartsAt       time.Time            `json:"starts_at"`
	Status         string               `json:"status"`
	TotalSeats     int                  `json:"total_seats"`
	AvailableSeats int                  `json:"available_seats"`
	Booked         int                  `json:"booked"`      // total - available
	PaidAmount     float64              `json:"paid_amount"` // এই event এ এখন পর্যন্ত paid (BDT)
}

type RecentBooking struct {
	BookingCode string               `json:"booking_code"`
	Customer    string               `json:"customer"`
	Event       i18n.LocalizedString `json:"event"`
	Amount      float64              `json:"amount"` // booking এর currency তেই (Currency দেখো)
	Currency    string               `json:"currency"`
	Status      string               `json:"status"` // paid | pending | failed | cancelled
	CreatedAt   time.Time            `json:"created_at"`
}

// Transaction: একটা paid (success) payment।
type Transaction struct {
	ID          uint                 `json:"id"`
	BookingCode string               `json:"booking_code"`
	EventTitle  i18n.LocalizedString `json:"event_title"`
	Amount      float64              `json:"amount"`
	Currency    string               `json:"currency"` // payment এর নিজের currency
	PaidAt      time.Time            `json:"paid_at"`
}

/* ---------------- Admin ---------------- */

type AdminDashboard struct {
	Currency        string          `json:"currency"`
	Revenue         Stat            `json:"revenue"`
	TicketsSold     Stat            `json:"tickets_sold"`
	ActiveEvents    Stat            `json:"active_events"`
	NewUsers        Stat            `json:"new_users"`
	TotalEvents     int             `json:"total_events"`   // upcoming event (canceled বাদে)
	TotalTickets    int             `json:"total_tickets"`  // ওই upcoming event গুলোর মোট ticket
	TotalUsers      int             `json:"total_users"`    // role = user
	TotalManagers   int             `json:"total_managers"` // role = manager
	RevenueSeries   []RevenuePoint  `json:"revenue_series"` // শেষ ১২ সপ্তাহ
	RevenueTotal    float64         `json:"revenue_total"`
	ThisWeekRevenue float64         `json:"this_week_revenue"`
	UpcomingEvents  []UpcomingEvent `json:"upcoming_events"` // ৪টা
	RecentBookings  []RecentBooking `json:"recent_bookings"` // ৬টা
}

/* ---------------- Manager ---------------- */

type ManagerStats struct {
	TotalEvents       int     `json:"total_events"`       // manager এর বানানো সব event
	UpcomingEvents    int     `json:"upcoming_events"`    // কয়টা upcoming
	TotalPaid         float64 `json:"total_paid"`         // সব event এ মোট paid
	UpcomingPaid      float64 `json:"upcoming_paid"`      // upcoming event গুলোতে paid
	UpcomingPotential float64 `json:"upcoming_potential"` // upcoming এর সব ticket বিক্রি হলে মোট (total_tickets × price)
	SeatsTotal        int     `json:"seats_total"`        // upcoming event গুলোর মোট seat
	SeatsAvailable    int     `json:"seats_available"`    // এখনো খালি seat
	BookingsCount     int     `json:"bookings_count"`     // upcoming event গুলোতে কয়টা booking
}

type ManagerDashboard struct {
	Currency       string          `json:"currency"`
	Stats          ManagerStats    `json:"stats"`
	UpcomingEvents []UpcomingEvent `json:"upcoming_events"` // ৪টা
	Transactions   []Transaction   `json:"transactions"`    // সাম্প্রতিক paid, ৫টা
}

/* ---------------- User ---------------- */

type UserStats struct {
	TotalPaid      float64 `json:"total_paid"`
	PaidEvents     int     `json:"paid_events"`     // কয়টা event এর জন্য payment করেছে
	UnpaidUpcoming int     `json:"unpaid_upcoming"` // upcoming event এ booking আছে, payment বাকি
}

type UserEvent struct {
	ID          uint                 `json:"id"`
	Title       i18n.LocalizedString `json:"title"`
	Location    i18n.LocalizedString `json:"location"`
	ImageURL    string               `json:"image_url"`
	StartsAt    time.Time            `json:"starts_at"`
	BookingCode string               `json:"booking_code"`
	Quantity    int                  `json:"quantity"`
}

type UnpaidBooking struct {
	BookingID   uint                 `json:"booking_id"`
	BookingCode string               `json:"booking_code"`
	EventID     uint                 `json:"event_id"`
	EventTitle  i18n.LocalizedString `json:"event_title"`
	Location    i18n.LocalizedString `json:"location"`
	StartsAt    time.Time            `json:"starts_at"`
	Quantity    int                  `json:"quantity"`
	Amount      float64              `json:"amount"`
	Currency    string               `json:"currency"`
}

type UserDashboard struct {
	Currency       string          `json:"currency"`
	Stats          UserStats       `json:"stats"`
	UpcomingEvents []UserEvent     `json:"upcoming_events"` // paid + upcoming
	Payments       []Transaction   `json:"payments"`        // সাম্প্রতিক paid, ৫টা
	UnpaidUpcoming []UnpaidBooking `json:"unpaid_upcoming"` // payment বাকি
}
