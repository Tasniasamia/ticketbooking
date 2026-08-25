package messaging

import (
	"errors"

	"time"

	"gorm.io/gorm"
)

var (
	ErrConversationNotFound = errors.New("conversation not found")
	ErrMessageNotFound      = errors.New("message not found")
	ErrForbidden            = errors.New("you are not allowed to access this conversation")
	ErrCannotMessage        = errors.New("you are not allowed to message this user")
)

type Repository interface {
	// Conversation
	FindOrCreateConversation(userA, userB uint, eventID *uint) (*Conversation, error)
	GetConversationByID(id uint) (*Conversation, error)
	GetConversationsForUser(userID uint) ([]Conversation, error)
	UpdateLastMessageAt(convID uint, t time.Time) error

	// Message
	CreateMessage(msg *Message) error
	GetMessages(convID uint, limit, offset int) ([]Message, int64, error)
	MarkMessagesRead(convID uint, userID uint, messageIDs []uint) (int64, error)
	UnreadCount(convID uint, userID uint) (int64, error)
	GetLastMessage(convID uint) (*Message, error)

	// Permission helpers
	CanMessage(senderID, receiverID uint, senderRole, receiverRole string) (bool, error)
	GetUserRole(userID uint) (string, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindOrCreateConversation(userA, userB uint, eventID *uint) (*Conversation, error) {
	u1, u2 := NormalizePair(userA, userB)

	var conv Conversation
	err := r.db.Preload("User1").Preload("User2").Where("user1_id = ? AND user2_id = ?", u1, u2).First(&conv).Error
	if err == nil {
		return &conv, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	conv = Conversation{
		User1ID: u1,
		User2ID: u2,
		EventID: eventID,
	}
	if err := r.db.Create(&conv).Error; err != nil {
		return nil, err
	}
		if err := r.db.
		Preload("User1").
		Preload("User2").
		First(&conv, conv.ID).Error; err != nil {
		return nil, err}
	

	return &conv, nil
}

func (r *repository) GetConversationByID(id uint) (*Conversation, error) {
	var conv Conversation
	err := r.db.
		Preload("User1").
		Preload("User2").
		First(&conv, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrConversationNotFound
	}
	return &conv, err
}

func (r *repository) GetConversationsForUser(userID uint) ([]Conversation, error) {
	var list []Conversation
	err := r.db.
		Where("user1_id = ? OR user2_id = ?", userID, userID).
		Preload("User1").
		Preload("User2").
		Order("COALESCE(last_message_at, created_at) DESC").
		Find(&list).Error
	return list, err
}

func (r *repository) UpdateLastMessageAt(convID uint, t time.Time) error {
	return r.db.Model(&Conversation{}).Where("id = ?", convID).
		Update("last_message_at", t).Error
}

func (r *repository) CreateMessage(msg *Message) error {
	if err := r.db.Create(msg).Error; err != nil {
		return err
	}
	// reload with sender
	return r.db.Preload("Sender").First(msg, msg.ID).Error
}

func (r *repository) GetMessages(convID uint, limit, offset int) ([]Message, int64, error) {
	var total int64
	r.db.Model(&Message{}).Where("conversation_id = ?", convID).Count(&total)

	var msgs []Message
	err := r.db.
		Where("conversation_id = ?", convID).
		Preload("Sender").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&msgs).Error
	return msgs, total, err
}

func (r *repository) MarkMessagesRead(convID uint, userID uint, messageIDs []uint) (int64, error) {
	// only mark messages that are NOT sent by this user
	res := r.db.Model(&Message{}).
		Where("conversation_id = ? AND sender_id != ? AND id IN ?", convID, userID, messageIDs).
		Update("is_read", true)
	return res.RowsAffected, res.Error
}

func (r *repository) UnreadCount(convID uint, userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&Message{}).
		Where("conversation_id = ? AND sender_id != ? AND is_read = false", convID, userID).
		Count(&count).Error
	return count, err
}

func (r *repository) GetLastMessage(convID uint) (*Message, error) {
	var msg Message
	err := r.db.
		Where("conversation_id = ?", convID).
		Preload("Sender").
		Order("created_at DESC").
		First(&msg).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &msg, err
}

// ---------- Permission logic ----------

func (r *repository) GetUserRole(userID uint) (string, error) {
	var role string
	err := r.db.Table("users").Select("role").Where("id = ?", userID).Scan(&role).Error
	return role, err
}

func (r *repository) CanMessage(senderID, receiverID uint, senderRole, receiverRole string) (bool, error) {
	if senderID == receiverID {
		return false, nil
	}

	// Admin can talk to anyone, anyone can talk to Admin
	if senderRole == "admin" || receiverRole == "admin" {
		return true, nil
	}

	// Manager → User: only if the user booked any of manager's events
	if senderRole == "manager" && receiverRole == "user" {
		return r.userBookedManagerEvent(receiverID, senderID)
	}

	// User → Manager: only if user booked any event of that manager
	if senderRole == "user" && receiverRole == "manager" {
		return r.userBookedManagerEvent(senderID, receiverID)
	}

	// User ↔ User: only if they share at least one common event booking
	if senderRole == "user" && receiverRole == "user" {
		return r.shareCommonEvent(senderID, receiverID)
	}

	// Manager ↔ Manager not allowed (unless you want later)
	return false, nil
}

// user booked any event owned by managerID
func (r *repository) userBookedManagerEvent(userID, managerID uint) (bool, error) {
	var count int64
	err := r.db.Table("bookings").
		Joins("JOIN events ON events.id = bookings.event_id").
		Where("bookings.user_id = ? AND events.manager_id = ? AND bookings.deleted_at IS NULL", userID, managerID).
		Count(&count).Error
	return count > 0, err
}

// both users have at least one booking on the same event
func (r *repository) shareCommonEvent(userA, userB uint) (bool, error) {
	var count int64
	err := r.db.Raw(`
		SELECT COUNT(*) FROM (
			SELECT event_id FROM bookings WHERE user_id = ? AND deleted_at IS NULL
			INTERSECT
			SELECT event_id FROM bookings WHERE user_id = ? AND deleted_at IS NULL
		) AS common
	`, userA, userB).Scan(&count).Error
	return count > 0, err
}
