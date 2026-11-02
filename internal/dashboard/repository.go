package dashboard

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"

	"ticketBooking/internal/dashboard/dto"
	"ticketBooking/internal/event"
	"ticketBooking/internal/utils/i18n"
)

// Table / column নাম (project এর entity অনুযায়ী):
//   events   (gorm.Model, title/location jsonb, starts_at, total_tickets, available_tickets, price, currency, manager_id, status, thumbnail_image)
//   bookings (id, user_id, event_id, quantity, total_price, status[pending|confirmed|cancelled], booking_code, created_at)
//   payments (id, booking_id, amount, status[pending|success|failed|cancelled], paid_at, created_at)
//   users    (id, name, role, created_at)
//   currencies (code, rate_to_bdt)
//
// টাকার অংক সবসময় event এর currency তে ধরা হয় (payment ওই দামেই হয়), তারপর BDT তে convert করে sum করা হয়।

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db} }

const (
	payOK       = "success"
	evCanceled  = "canceled"  // event status
	bkCancelled = "cancelled" // booking status
	paidAtExpr  = "COALESCE(p.paid_at, p.created_at)"
)

/* ---------------- currency helpers ---------------- */

type moneyRow struct {
	Currency string  `gorm:"column:currency"`
	Total    float64 `gorm:"column:total"`
}

// rates: currencies table থেকে code → rate_to_bdt। পড়া না গেলে সব ১ ধরে (crash করে না)।
func (r *Repository) rates() map[string]float64 {
	var rows []struct {
		Code string  `gorm:"column:code"`
		Rate float64 `gorm:"column:rate_to_bdt"`
	}
	m := map[string]float64{}
	if err := r.db.Raw(`SELECT code, rate_to_bdt FROM currencies`).Scan(&rows).Error; err != nil {
		log.Println("dashboard: currency rates unavailable, using 1:1 →", err)
		return m
	}
	for _, row := range rows {
		m[strings.ToUpper(row.Code)] = row.Rate
	}
	return m
}

func toBase(rates map[string]float64, amount float64, currency string) float64 {
	rate := rates[strings.ToUpper(currency)]
	if rate <= 0 {
		rate = 1
	}
	return amount * rate
}

func sumBase(rates map[string]float64, rows []moneyRow) float64 {
	var t float64
	for _, row := range rows {
		t += toBase(rates, row.Total, row.Currency)
	}
	return t
}

/* ---------------- small utils ---------------- */

func pct(cur, prev float64) *float64 {
	if prev == 0 {
		return nil
	}
	v := (cur - prev) / prev * 100
	return &v
}

func stat(cur, prev float64) dto.Stat { return dto.Stat{Value: cur, ChangePct: pct(cur, prev)} }

func weekStart(t time.Time) time.Time {
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return t.AddDate(0, 0, -((int(t.Weekday()) + 6) % 7)) // সোমবার = শুরু
}

// thumbURL: media.MediaImage এর field নামে নির্ভর না করে JSON দিয়ে url বের করে।
func thumbURL(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	var o struct {
		URL string `json:"url"`
	}
	_ = json.Unmarshal(b, &o)
	return o.URL
}

/* ---------------- shared: upcoming events ---------------- */

// upcomingEvents: managerID == 0 হলে সবার (শুধু approved), নইলে ওই manager এর (canceled বাদে সব)।
func (r *Repository) upcomingEvents(managerID uint, limit int, rates map[string]float64) ([]dto.UpcomingEvent, error) {
	q := r.db.Model(&event.Event{}).Where("starts_at >= ?", time.Now())
	if managerID != 0 {
		q = q.Where("manager_id = ? AND status <> ?", managerID, evCanceled)
	} else {
		q = q.Where("status = ?", "approved")
	}
	var events []event.Event
	if err := q.Order("starts_at ASC").Limit(limit).Find(&events).Error; err != nil {
		return nil, err
	}

	out := make([]dto.UpcomingEvent, 0, len(events))
	if len(events) == 0 {
		return out, nil
	}
	ids := make([]uint, len(events))
	for i, e := range events {
		ids[i] = e.ID
	}

	// প্রতি event এ কত paid হয়েছে (event এর নিজের currency তে)
	var paid []struct {
		EventID uint    `gorm:"column:event_id"`
		Total   float64 `gorm:"column:total"`
	}
	if err := r.db.Raw(`
		SELECT b.event_id AS event_id, COALESCE(SUM(p.amount),0) AS total
		FROM payments p JOIN bookings b ON b.id = p.booking_id
		WHERE p.status = ? AND b.event_id IN ? GROUP BY b.event_id`, payOK, ids).Scan(&paid).Error; err != nil {
		return nil, err
	}
	paidBy := map[uint]float64{}
	for _, row := range paid {
		paidBy[row.EventID] = row.Total
	}

	for _, e := range events {
		out = append(out, dto.UpcomingEvent{
			ID:             e.ID,
			Title:          e.Title,
			Location:       e.Location,
			ImageURL:       thumbURL(e.ThumbnailImage),
			StartsAt:       e.StartsAt,
			Status:         string(e.Status),
			TotalSeats:     e.TotalTickets,
			AvailableSeats: max(e.AvailableTickets, 0),
			Booked:         max(e.TotalTickets-e.AvailableTickets, 0),
			PaidAmount:     toBase(rates, paidBy[e.ID], e.Currency),
		})
	}
	return out, nil
}

/* ---------------- admin ---------------- */

func (r *Repository) paidRevenue(rates map[string]float64, from, to time.Time) (float64, error) {
	var rows []moneyRow
	err := r.db.Raw(`
		SELECT e.currency AS currency, COALESCE(SUM(p.amount),0) AS total
		FROM payments p JOIN bookings b ON b.id = p.booking_id JOIN events e ON e.id = b.event_id
		WHERE p.status = ? AND `+paidAtExpr+` >= ? AND `+paidAtExpr+` < ?
		GROUP BY e.currency`, payOK, from, to).Scan(&rows).Error
	return sumBase(rates, rows), err
}

func (r *Repository) count(q string, args ...any) (float64, error) {
	var v float64
	err := r.db.Raw(q, args...).Scan(&v).Error
	return v, err
}

func (r *Repository) Admin() (*dto.AdminDashboard, error) {
	now := time.Now()
	cur0, prev0 := now.AddDate(0, 0, -28), now.AddDate(0, 0, -56)
	rates := r.rates()
	out := &dto.AdminDashboard{Currency: dto.BaseCurrency}

	revCur, err := r.paidRevenue(rates, cur0, now)
	if err != nil {
		return nil, err
	}
	revPrev, err := r.paidRevenue(rates, prev0, cur0)
	if err != nil {
		return nil, err
	}
	out.Revenue = stat(revCur, revPrev)

	tixQ := `SELECT COALESCE(SUM(b.quantity),0) FROM bookings b
		WHERE b.created_at >= ? AND b.created_at < ?
		AND EXISTS (SELECT 1 FROM payments p WHERE p.booking_id = b.id AND p.status = ?)`
	tixCur, err := r.count(tixQ, cur0, now, payOK)
	if err != nil {
		return nil, err
	}
	tixPrev, err := r.count(tixQ, prev0, cur0, payOK)
	if err != nil {
		return nil, err
	}
	out.TicketsSold = stat(tixCur, tixPrev)

	usrQ := `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL AND role = 'user' AND created_at >= ? AND created_at < ?`
	usrCur, err := r.count(usrQ, cur0, now)
	if err != nil {
		return nil, err
	}
	usrPrev, err := r.count(usrQ, prev0, cur0)
	if err != nil {
		return nil, err
	}
	out.NewUsers = stat(usrCur, usrPrev)

	// Active events = এখন approved + upcoming; % = আগের ২৮ দিনের তুলনায় নতুন approved event কত
	active, err := r.count(`SELECT COUNT(*) FROM events WHERE deleted_at IS NULL AND status = 'approved' AND starts_at >= ?`, now)
	if err != nil {
		return nil, err
	}
	evQ := `SELECT COUNT(*) FROM events WHERE deleted_at IS NULL AND status = 'approved' AND created_at >= ? AND created_at < ?`
	evCur, _ := r.count(evQ, cur0, now)
	evPrev, _ := r.count(evQ, prev0, cur0)
	out.ActiveEvents = dto.Stat{Value: active, ChangePct: pct(evCur, evPrev)}

	// Total events / total tickets: upcoming (canceled বাদে) event গুলোর জন্য
	var upc struct {
		Events  int `gorm:"column:events"`
		Tickets int `gorm:"column:tickets"`
	}
	if err := r.db.Raw(`
		SELECT COUNT(*) AS events, COALESCE(SUM(total_tickets),0) AS tickets
		FROM events WHERE deleted_at IS NULL AND starts_at >= ? AND status <> ?`, now, evCanceled).Scan(&upc).Error; err != nil {
		return nil, err
	}
	out.TotalEvents, out.TotalTickets = upc.Events, upc.Tickets

	// Total users / managers
	var ppl []struct {
		Role  string `gorm:"column:role"`
		Total int    `gorm:"column:total"`
	}
	if err := r.db.Raw(`SELECT role, COUNT(*) AS total FROM users WHERE deleted_at IS NULL GROUP BY role`).Scan(&ppl).Error; err != nil {
		return nil, err
	}
	for _, x := range ppl {
		switch x.Role {
		case "user":
			out.TotalUsers = x.Total
		case "manager":
			out.TotalManagers = x.Total
		}
	}

	// শেষ ১২ সপ্তাহের revenue (দিন ধরে group করে Go তে সপ্তাহে বসানো)
	start := weekStart(now).AddDate(0, 0, -7*11)
	var days []struct {
		Day      time.Time `gorm:"column:day"`
		Currency string    `gorm:"column:currency"`
		Total    float64   `gorm:"column:total"`
	}
	if err := r.db.Raw(`
		SELECT date_trunc('day', `+paidAtExpr+`) AS day, e.currency AS currency, COALESCE(SUM(p.amount),0) AS total
		FROM payments p JOIN bookings b ON b.id = p.booking_id JOIN events e ON e.id = b.event_id
		WHERE p.status = ? AND `+paidAtExpr+` >= ?
		GROUP BY 1, 2`, payOK, start).Scan(&days).Error; err != nil {
		return nil, err
	}
	byWeek := map[string]float64{}
	for _, d := range days {
		k := weekStart(d.Day.In(now.Location())).Format("2006-01-02")
		byWeek[k] += toBase(rates, d.Total, d.Currency)
	}
	for i := 0; i < 12; i++ {
		k := start.AddDate(0, 0, 7*i).Format("2006-01-02")
		out.RevenueSeries = append(out.RevenueSeries, dto.RevenuePoint{Week: k, Amount: byWeek[k]})
		out.RevenueTotal += byWeek[k]
	}
	out.ThisWeekRevenue = out.RevenueSeries[11].Amount

	if out.UpcomingEvents, err = r.upcomingEvents(0, 4, rates); err != nil {
		return nil, err
	}

	var rows []struct {
		BookingCode string               `gorm:"column:booking_code"`
		Customer    string               `gorm:"column:customer"`
		Event       i18n.LocalizedString `gorm:"column:event"`
		Amount      float64              `gorm:"column:amount"`
		Currency    string               `gorm:"column:currency"`
		Status      string               `gorm:"column:status"`
		CreatedAt   time.Time            `gorm:"column:created_at"`
	}
	if err := r.db.Raw(`
		SELECT b.booking_code AS booking_code, COALESCE(u.name,'') AS customer, e.title AS event,
		  b.total_price AS amount, e.currency AS currency, b.created_at AS created_at,
		  CASE
		    WHEN EXISTS (SELECT 1 FROM payments p WHERE p.booking_id = b.id AND p.status = ?) THEN 'paid'
		    WHEN b.status = ? THEN 'cancelled'
		    WHEN EXISTS (SELECT 1 FROM payments p WHERE p.booking_id = b.id AND p.status = 'failed') THEN 'failed'
		    ELSE 'pending'
		  END AS status
		FROM bookings b JOIN events e ON e.id = b.event_id LEFT JOIN users u ON u.id = b.user_id
		ORDER BY b.created_at DESC LIMIT 6`, payOK, bkCancelled).Scan(&rows).Error; err != nil {
		return nil, err
	}
	out.RecentBookings = make([]dto.RecentBooking, 0, len(rows))
	for _, x := range rows {
		out.RecentBookings = append(out.RecentBookings, dto.RecentBooking{
			BookingCode: x.BookingCode, Customer: x.Customer, Event: x.Event, Amount: x.Amount,
			Currency: x.Currency, Status: x.Status, CreatedAt: x.CreatedAt,
		})
	}
	return out, nil
}

/* ---------------- shared: paid transactions ---------------- */

type txRow struct {
	ID          uint                 `gorm:"column:id"`
	BookingCode string               `gorm:"column:booking_code"`
	EventTitle  i18n.LocalizedString `gorm:"column:event_title"`
	Amount      float64              `gorm:"column:amount"`
	Currency    string               `gorm:"column:currency"`
	PaidAt      time.Time            `gorm:"column:paid_at"`
}

// recentPaid: scope = "e.manager_id = ?" বা "b.user_id = ?"
func (r *Repository) recentPaid(scope string, id uint, limit int) ([]dto.Transaction, error) {
	var rows []txRow
	err := r.db.Raw(`
		SELECT p.id AS id, b.booking_code AS booking_code, e.title AS event_title,
		  p.amount AS amount, e.currency AS currency, `+paidAtExpr+` AS paid_at
		FROM payments p JOIN bookings b ON b.id = p.booking_id JOIN events e ON e.id = b.event_id
		WHERE p.status = ? AND `+scope+`
		ORDER BY `+paidAtExpr+` DESC LIMIT ?`, payOK, id, limit).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]dto.Transaction, 0, len(rows))
	for _, x := range rows {
		out = append(out, dto.Transaction{
			ID: x.ID, BookingCode: x.BookingCode, EventTitle: x.EventTitle,
			Amount: x.Amount, Currency: x.Currency, PaidAt: x.PaidAt,
		})
	}
	return out, nil
}

/* ---------------- manager ---------------- */

func (r *Repository) Manager(managerID uint) (*dto.ManagerDashboard, error) {
	now := time.Now()
	rates := r.rates()
	out := &dto.ManagerDashboard{Currency: dto.BaseCurrency}
	s := &out.Stats

	// event সংখ্যা + seat
	var ev struct {
		TotalEvents    int `gorm:"column:total_events"`
		UpcomingEvents int `gorm:"column:upcoming_events"`
		SeatsTotal     int `gorm:"column:seats_total"`
		SeatsAvailable int `gorm:"column:seats_available"`
	}
	if err := r.db.Raw(`
		SELECT COUNT(*) AS total_events,
		  COUNT(*) FILTER (WHERE starts_at >= @now AND status <> @canceled) AS upcoming_events,
		  COALESCE(SUM(total_tickets) FILTER (WHERE starts_at >= @now AND status <> @canceled),0) AS seats_total,
		  COALESCE(SUM(available_tickets) FILTER (WHERE starts_at >= @now AND status <> @canceled),0) AS seats_available
		FROM events WHERE manager_id = @m AND deleted_at IS NULL`,
		map[string]any{"now": now, "canceled": evCanceled, "m": managerID}).Scan(&ev).Error; err != nil {
		return nil, err
	}
	s.TotalEvents, s.UpcomingEvents = ev.TotalEvents, ev.UpcomingEvents
	s.SeatsTotal, s.SeatsAvailable = ev.SeatsTotal, max(ev.SeatsAvailable, 0)

	// upcoming event গুলোর সব ticket বিক্রি হলে কত আসত
	var pot []moneyRow
	if err := r.db.Raw(`
		SELECT currency, COALESCE(SUM(total_tickets * price),0) AS total FROM events
		WHERE manager_id = ? AND deleted_at IS NULL AND starts_at >= ? AND status <> ?
		GROUP BY currency`, managerID, now, evCanceled).Scan(&pot).Error; err != nil {
		return nil, err
	}
	s.UpcomingPotential = sumBase(rates, pot)

	// মোট paid (সব event) আর upcoming event এ paid
	paidQ := `
		SELECT e.currency AS currency, COALESCE(SUM(p.amount),0) AS total
		FROM payments p JOIN bookings b ON b.id = p.booking_id JOIN events e ON e.id = b.event_id
		WHERE p.status = ? AND e.manager_id = ? %s GROUP BY e.currency`
	var all, up []moneyRow
	if err := r.db.Raw(strings.Replace(paidQ, "%s", "", 1), payOK, managerID).Scan(&all).Error; err != nil {
		return nil, err
	}
	if err := r.db.Raw(strings.Replace(paidQ, "%s", "AND e.starts_at >= ?", 1), payOK, managerID, now).Scan(&up).Error; err != nil {
		return nil, err
	}
	s.TotalPaid, s.UpcomingPaid = sumBase(rates, all), sumBase(rates, up)

	// upcoming event গুলোতে কয়টা booking
	var bc int
	if err := r.db.Raw(`
		SELECT COUNT(*) FROM bookings b JOIN events e ON e.id = b.event_id
		WHERE e.manager_id = ? AND e.deleted_at IS NULL AND e.starts_at >= ? AND e.status <> ? AND b.status <> ?`,
		managerID, now, evCanceled, bkCancelled).Scan(&bc).Error; err != nil {
		return nil, err
	}
	s.BookingsCount = bc

	var err error
	if out.UpcomingEvents, err = r.upcomingEvents(managerID, 4, rates); err != nil {
		return nil, err
	}
	out.Transactions, err = r.recentPaid("e.manager_id = ?", managerID, 5)
	return out, err
}

/* ---------------- user ---------------- */

func (r *Repository) User(userID uint) (*dto.UserDashboard, error) {
	now := time.Now()
	rates := r.rates()
	out := &dto.UserDashboard{Currency: dto.BaseCurrency}

	var paid []moneyRow
	if err := r.db.Raw(`
		SELECT e.currency AS currency, COALESCE(SUM(p.amount),0) AS total
		FROM payments p JOIN bookings b ON b.id = p.booking_id JOIN events e ON e.id = b.event_id
		WHERE p.status = ? AND b.user_id = ? GROUP BY e.currency`, payOK, userID).Scan(&paid).Error; err != nil {
		return nil, err
	}
	out.Stats.TotalPaid = sumBase(rates, paid)

	if err := r.db.Raw(`
		SELECT COUNT(DISTINCT b.event_id) FROM payments p JOIN bookings b ON b.id = p.booking_id
		WHERE p.status = ? AND b.user_id = ?`, payOK, userID).Scan(&out.Stats.PaidEvents).Error; err != nil {
		return nil, err
	}

	// ---- যেসব upcoming event এ payment হয়ে গেছে ----
	var pe []struct {
		EventID     uint   `gorm:"column:event_id"`
		BookingCode string `gorm:"column:booking_code"`
		Quantity    int    `gorm:"column:quantity"`
	}
	if err := r.db.Raw(`
		SELECT b.event_id AS event_id, b.booking_code AS booking_code, b.quantity AS quantity
		FROM bookings b JOIN events e ON e.id = b.event_id
		WHERE b.user_id = ? AND e.deleted_at IS NULL AND e.starts_at >= ?
		  AND EXISTS (SELECT 1 FROM payments p WHERE p.booking_id = b.id AND p.status = ?)
		ORDER BY e.starts_at ASC, b.id ASC LIMIT 30`, userID, now, payOK).Scan(&pe).Error; err != nil {
		return nil, err
	}
	out.UpcomingEvents = []dto.UserEvent{}
	if len(pe) > 0 {
		ids := make([]uint, 0, len(pe))
		for _, x := range pe {
			ids = append(ids, x.EventID)
		}
		var events []event.Event
		if err := r.db.Where("id IN ?", ids).Find(&events).Error; err != nil {
			return nil, err
		}
		byID := map[uint]event.Event{}
		for _, e := range events {
			byID[e.ID] = e
		}
		seen := map[uint]int{} // event id → out এর index (একই event এর একাধিক booking এর ticket যোগ হবে)
		for _, x := range pe {
			e, found := byID[x.EventID]
			if !found {
				continue
			}
			if i, dup := seen[x.EventID]; dup {
				out.UpcomingEvents[i].Quantity += x.Quantity
				continue
			}
			if len(out.UpcomingEvents) == 4 {
				continue
			}
			seen[x.EventID] = len(out.UpcomingEvents)
			out.UpcomingEvents = append(out.UpcomingEvents, dto.UserEvent{
				ID: e.ID, Title: e.Title, Location: e.Location, ImageURL: thumbURL(e.ThumbnailImage),
				StartsAt: e.StartsAt, BookingCode: x.BookingCode, Quantity: x.Quantity,
			})
		}
	}

	var err error
	if out.Payments, err = r.recentPaid("b.user_id = ?", userID, 5); err != nil {
		return nil, err
	}

	// ---- upcoming event এ booking আছে কিন্তু payment হয়নি ----
	var rows []struct {
		BookingID   uint                 `gorm:"column:booking_id"`
		BookingCode string               `gorm:"column:booking_code"`
		EventID     uint                 `gorm:"column:event_id"`
		EventTitle  i18n.LocalizedString `gorm:"column:event_title"`
		Location    i18n.LocalizedString `gorm:"column:location"`
		StartsAt    time.Time            `gorm:"column:starts_at"`
		Quantity    int                  `gorm:"column:quantity"`
		Amount      float64              `gorm:"column:amount"`
		Currency    string               `gorm:"column:currency"`
	}
	if err := r.db.Raw(`
		SELECT b.id AS booking_id, b.booking_code AS booking_code, e.id AS event_id, e.title AS event_title,
		  e.location AS location, e.starts_at AS starts_at, b.quantity AS quantity,
		  b.total_price AS amount, e.currency AS currency
		FROM bookings b JOIN events e ON e.id = b.event_id
		WHERE b.user_id = ? AND b.status <> ? AND e.deleted_at IS NULL AND e.starts_at >= ?
		  AND NOT EXISTS (SELECT 1 FROM payments p WHERE p.booking_id = b.id AND p.status = ?)
		ORDER BY e.starts_at ASC, b.id ASC LIMIT 50`, userID, bkCancelled, now, payOK).Scan(&rows).Error; err != nil {
		return nil, err
	}
	out.Stats.UnpaidUpcoming = len(rows)
	out.UnpaidUpcoming = make([]dto.UnpaidBooking, 0, 10)
	for i, x := range rows {
		if i == 10 {
			break
		}
		out.UnpaidUpcoming = append(out.UnpaidUpcoming, dto.UnpaidBooking{
			BookingID: x.BookingID, BookingCode: x.BookingCode, EventID: x.EventID, EventTitle: x.EventTitle,
			Location: x.Location, StartsAt: x.StartsAt, Quantity: x.Quantity, Amount: x.Amount, Currency: x.Currency,
		})
	}
	return out, nil
}
