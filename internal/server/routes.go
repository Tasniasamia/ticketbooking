

// package server

// import (
// 	"ticketBooking/internal/blog"
// 	"ticketBooking/internal/booking"
// 	"ticketBooking/internal/comment"
// 	"ticketBooking/internal/live"
// 	"ticketBooking/internal/messaging"

// 	"ticketBooking/internal/config"
// 	"ticketBooking/internal/currency"
// 	"ticketBooking/internal/event"
// 	"ticketBooking/internal/eventCategory"
// 	"ticketBooking/internal/language"
// 	"ticketBooking/internal/media"
// 	"ticketBooking/internal/payment"
// 	"ticketBooking/internal/settings"
// 	"ticketBooking/internal/translation"
// 	"ticketBooking/internal/user"

// 	"github.com/labstack/echo/v5"
// 	"gorm.io/gorm"
// )

// type RouteRegistrar func(api *echo.Group, db *gorm.DB, config config.Config)

// func RegisterAllRoutes(api *echo.Group, db *gorm.DB, config config.Config, hub *messaging.Hub) {
// 	// Public live channel: ticket sonkha bodlale browser ke push kore (WebSocket).
// 	api.GET("/live/events", live.Default.Serve)

// 	registrars := []RouteRegistrar{
// 		//user.UserRegisterRoutes,
// 		event.EventRegisterRoutes,
// 		language.LanguageRegisterRoutes,
// 		translation.TranslationRegisterRoutes,
// 		media.MediaRegisterRoutes,
// 		currency.CurrencyRegisterRoutes,
// 		settings.SettingsRegisterRoutes,
// 		booking.BookingRegisterRoutes,
// 		eventCategory.EventCategoryRegisterRoutes,
// 		blog.BlogRegisterRoutes,
// 		comment.CommentRegisterRoutes,
// 	}

// 	for _, register := range registrars {
// 		register(api, db, config)
// 	}
// 	messaging.RegisterRoutes(api, db, config, hub)
// 	payment.PaymentRegisterRoutes(api, db, config, hub)

// 	// 2. user service-এ সেই msgService পাস করো
// 	user.UserRegisterRoutes(api, db, config) // ← এখানে

// }


package server

import (
	"ticketBooking/internal/blog"
	"ticketBooking/internal/booking"
	"ticketBooking/internal/comment"
	"ticketBooking/internal/live"
	"ticketBooking/internal/messaging"
	messageDTO "ticketBooking/internal/messaging/dto"

	"ticketBooking/internal/config"
	"ticketBooking/internal/currency"
	"ticketBooking/internal/event"
	"ticketBooking/internal/eventCategory"
	"ticketBooking/internal/language"
	"ticketBooking/internal/media"
	"ticketBooking/internal/payment"
	"ticketBooking/internal/settings"
	"ticketBooking/internal/translation"
	"ticketBooking/internal/user"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type RouteRegistrar func(api *echo.Group, db *gorm.DB, config config.Config)

// messaging → user.ConversationStarter (package server, cycle নেই)
type messagingConvBridge struct {
	svc *messaging.Service
}

func (b *messagingConvBridge) Start(senderID, receiverID uint) error {
	_, err := b.svc.StartConversation(senderID, &messageDTO.StartConversationRequest{
		ReceiverID: receiverID,
	})
	return err
}

func RegisterAllRoutes(api *echo.Group, db *gorm.DB, config config.Config, hub *messaging.Hub) {
	api.GET("/live/events", live.Default.Serve)

	registrars := []RouteRegistrar{
		event.EventRegisterRoutes,
		language.LanguageRegisterRoutes,
		translation.TranslationRegisterRoutes,
		media.MediaRegisterRoutes,
		currency.CurrencyRegisterRoutes,
		settings.SettingsRegisterRoutes,
		booking.BookingRegisterRoutes,
		eventCategory.EventCategoryRegisterRoutes,
		blog.BlogRegisterRoutes,
		comment.CommentRegisterRoutes,
	}

	for _, register := range registrars {
		register(api, db, config)
	}

	messaging.RegisterRoutes(api, db, config, hub)
	payment.PaymentRegisterRoutes(api, db, config, hub)

	// verify-otp → admin conversation (user package messaging import করে না)
	msgSvc := messaging.NewService(messaging.NewRepository(db), hub)
	user.UserRegisterRoutes(api, db, config, &messagingConvBridge{svc: msgSvc})
}