// package messaging

// import (
// 	"errors"
// 	"strings"
// 	"ticketBooking/internal/messaging/dto"
// 	"time"
// )

// type Service struct {
// 	repo Repository
// 	hub  *Hub
// }

// func NewService(repo Repository, hub *Hub) *Service {
// 	return &Service{repo: repo, hub: hub}
// }

// func (s *Service) StartConversation(senderID uint, req *dto.StartConversationRequest) (*dto.ConversationResponse, error) {
// 	senderRole, err := s.repo.GetUserRole(senderID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	receiverRole, err := s.repo.GetUserRole(req.ReceiverID)
// 	if err != nil {
// 		return nil, errors.New("receiver not found")
// 	}
// 	ok, err := s.repo.CanMessage(senderID, req.ReceiverID, senderRole, receiverRole)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if !ok {
// 		return nil, ErrCannotMessage
// 	}
// 	conv, err := s.repo.FindOrCreateConversation(senderID, req.ReceiverID, req.EventID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return s.buildConversationResponse(conv, senderID)
// }

// func (s *Service) SendMessage(senderID uint, req *dto.SendMessageRequest) (*dto.MessageResponse, error) {
// 	senderRole, err := s.repo.GetUserRole(senderID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	receiverRole, err := s.repo.GetUserRole(req.ReceiverID)
// 	if err != nil {
// 		return nil, errors.New("receiver not found")
// 	}
// 	ok, err := s.repo.CanMessage(senderID, req.ReceiverID, senderRole, receiverRole)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if !ok {
// 		return nil, ErrCannotMessage
// 	}

// 	content := strings.TrimSpace(req.Content)
// 	hasFile := strings.TrimSpace(req.AttachmentURL) != ""
// 	if content == "" && !hasFile {
// 		return nil, errors.New("content or attachment is required")
// 	}

// 	conv, err := s.repo.FindOrCreateConversation(senderID, req.ReceiverID, req.EventID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	msg := &Message{
// 		ConversationID:    conv.ID,
// 		SenderID:          senderID,
// 		Content:           content,
// 		AttachmentURL:     req.AttachmentURL,
// 		AttachmentType:    req.AttachmentType,
// 		AttachmentName:    req.AttachmentName,
// 		AttachmentMediaID: req.AttachmentMediaID,
// 		IsRead:            false,
// 	}
// 	if err := s.repo.CreateMessage(msg); err != nil {
// 		return nil, err
// 	}
// 	_ = s.repo.UpdateLastMessageAt(conv.ID, time.Now())

// 	// reload with sender
// 	if loaded, e := s.repo.GetMessageByID(msg.ID); e == nil {
// 		msg = loaded
// 	}
// 	res := msg.ToResponse()

// 	if s.hub != nil {
// 		out := dto.WSOutgoing{Type: "message", Payload: res}
// 		s.hub.SendToUser(req.ReceiverID, out)
// 		s.hub.SendToUser(senderID, out)
// 	}
// 	return res, nil
// }

// func (s *Service) EditMessage(userID, messageID uint, content string) (*dto.MessageResponse, error) {
// 	msg, err := s.repo.GetMessageByID(messageID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if msg.SenderID != userID {
// 		return nil, ErrForbidden
// 	}
// 	content = strings.TrimSpace(content)
// 	if content == "" {
// 		return nil, errors.New("content is required")
// 	}
// 	msg.Content = content
// 	msg.IsEdited = true
// 	if err := s.repo.UpdateMessage(msg); err != nil {
// 		return nil, err
// 	}
// 	res := msg.ToResponse()
// 	if s.hub != nil {
// 		conv, _ := s.repo.GetConversationByID(msg.ConversationID)
// 		out := dto.WSOutgoing{Type: "message_updated", Payload: res}
// 		if conv != nil {
// 			s.hub.SendToUser(conv.User1ID, out)
// 			s.hub.SendToUser(conv.User2ID, out)
// 		}
// 	}
// 	return res, nil
// }

// func (s *Service) DeleteMessage(userID, messageID uint) error {
// 	msg, err := s.repo.GetMessageByID(messageID)
// 	if err != nil {
// 		return err
// 	}
// 	if msg.SenderID != userID {
// 		return ErrForbidden
// 	}
// 	convID := msg.ConversationID
// 	if err := s.repo.DeleteMessage(messageID); err != nil {
// 		return err
// 	}
// 	if s.hub != nil {
// 		conv, _ := s.repo.GetConversationByID(convID)
// 		out := dto.WSOutgoing{
// 			Type: "message_deleted",
// 			Payload: map[string]any{
// 				"id":              messageID,
// 				"conversation_id": convID,
// 			},
// 		}
// 		if conv != nil {
// 			s.hub.SendToUser(conv.User1ID, out)
// 			s.hub.SendToUser(conv.User2ID, out)
// 		}
// 	}
// 	return nil
// }

// func (s *Service) GetMyConversations(userID uint) ([]*dto.ConversationResponse, error) {
// 	list, err := s.repo.GetConversationsForUser(userID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	result := make([]*dto.ConversationResponse, 0, len(list))
// 	for i := range list {
// 		resp, err := s.buildConversationResponse(&list[i], userID)
// 		if err != nil {
// 			continue
// 		}
// 		result = append(result, resp)
// 	}
// 	return result, nil
// }

// func (s *Service) GetMessages(userID, convID uint, page, pageSize int) ([]*dto.MessageResponse, int64, error) {
// 	conv, err := s.repo.GetConversationByID(convID)
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	if conv.User1ID != userID && conv.User2ID != userID {
// 		return nil, 0, ErrForbidden
// 	}
// 	if page < 1 {
// 		page = 1
// 	}
// 	if pageSize < 1 || pageSize > 100 {
// 		pageSize = 30
// 	}
// 	offset := (page - 1) * pageSize
// 	msgs, total, err := s.repo.GetMessages(convID, pageSize, offset)
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
// 		msgs[i], msgs[j] = msgs[j], msgs[i]
// 	}
// 	res := make([]*dto.MessageResponse, 0, len(msgs))
// 	for i := range msgs {
// 		res = append(res, msgs[i].ToResponse())
// 	}
// 	return res, total, nil
// }

// func (s *Service) MarkRead(userID, convID uint, messageIDs []uint) error {
// 	conv, err := s.repo.GetConversationByID(convID)
// 	if err != nil {
// 		return err
// 	}
// 	if conv.User1ID != userID && conv.User2ID != userID {
// 		return ErrForbidden
// 	}
// 	_, err = s.repo.MarkMessagesRead(convID, userID, messageIDs)
// 	if err != nil {
// 		return err
// 	}
// 	if s.hub != nil {
// 		otherID := conv.User1ID
// 		if otherID == userID {
// 			otherID = conv.User2ID
// 		}
// 		s.hub.SendToUser(otherID, dto.WSOutgoing{
// 			Type: "read",
// 			Payload: map[string]any{
// 				"conversation_id": convID,
// 				"message_ids":     messageIDs,
// 				"read_by":         userID,
// 			},
// 		})
// 	}
// 	return nil
// }

// func (s *Service) buildConversationResponse(conv *Conversation, currentUserID uint) (*dto.ConversationResponse, error) {
// 	var participant *dto.UserBrief
// 	if conv.User1ID == currentUserID {
// 		if conv.User2.ID != 0 {
// 			participant = &dto.UserBrief{
// 				ID: conv.User2.ID, Name: conv.User2.Name, Email: conv.User2.Email,
// 				Role: string(conv.User2.Role), ProfileImage: conv.User2.ProfileImage,
// 				ProfileImageId: conv.User2.ProfileImageId,
// 			}
// 		}
// 	} else if conv.User1.ID != 0 {
// 		participant = &dto.UserBrief{
// 			ID: conv.User1.ID, Name: conv.User1.Name, Email: conv.User1.Email,
// 			Role: string(conv.User1.Role), ProfileImage: conv.User1.ProfileImage,
// 			ProfileImageId: conv.User1.ProfileImageId,
// 		}
// 	}
// 	unread, _ := s.repo.UnreadCount(conv.ID, currentUserID)
// 	last, _ := s.repo.GetLastMessage(conv.ID)
// 	var lastResp *dto.MessageResponse
// 	if last != nil {
// 		lastResp = last.ToResponse()
// 	}
// 	updated := conv.UpdatedAt
// 	if conv.LastMessageAt != nil {
// 		updated = *conv.LastMessageAt
// 	}
// 	return &dto.ConversationResponse{
// 		ID:          conv.ID,
// 		Participant: participant,
// 		EventID:     conv.EventID,
// 		LastMessage: lastResp,
// 		UnreadCount: unread,
// 		UpdatedAt:   updated,
// 		CreatedAt:   conv.CreatedAt,
// 	}, nil
// }

// func (s *Service) HandleWSMessage(senderID uint, payload dto.WSMessagePayload) (*dto.MessageResponse, error) {
// 	return s.SendMessage(senderID, &dto.SendMessageRequest{
// 		ReceiverID: payload.ReceiverID,
// 		Content:    payload.Content,
// 		EventID:    payload.EventID,
// 	})
// }

// func (s *Service) HandleWSRead(userID uint, payload dto.WSReadPayload) error {
// 	return s.MarkRead(userID, payload.ConversationID, payload.MessageIDs)
// }

package messaging

import (
	"errors"
	"strings"
	"ticketBooking/internal/messaging/dto"
	"time"
)

type Service struct {
	repo Repository
	hub  *Hub
}

func NewService(repo Repository, hub *Hub) *Service {
	return &Service{repo: repo, hub: hub}
}

func (s *Service) StartConversation(senderID uint, req *dto.StartConversationRequest) (*dto.ConversationResponse, error) {
	senderRole, err := s.repo.GetUserRole(senderID)
	if err != nil {
		return nil, err
	}
	receiverRole, err := s.repo.GetUserRole(req.ReceiverID)
	if err != nil {
		return nil, errors.New("receiver not found")
	}
	ok, err := s.repo.CanMessage(senderID, req.ReceiverID, senderRole, receiverRole)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrCannotMessage
	}
	conv, err := s.repo.FindOrCreateConversation(senderID, req.ReceiverID, req.EventID)
	if err != nil {
		return nil, err
	}
	return s.buildConversationResponse(conv, senderID)
}

func (s *Service) SendMessage(senderID uint, req *dto.SendMessageRequest) (*dto.MessageResponse, error) {
	senderRole, err := s.repo.GetUserRole(senderID)
	if err != nil {
		return nil, err
	}
	receiverRole, err := s.repo.GetUserRole(req.ReceiverID)
	if err != nil {
		return nil, errors.New("receiver not found")
	}
	ok, err := s.repo.CanMessage(senderID, req.ReceiverID, senderRole, receiverRole)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrCannotMessage
	}

	content := strings.TrimSpace(req.Content)
	hasFile := strings.TrimSpace(req.AttachmentURL) != ""
	if content == "" && !hasFile {
		return nil, errors.New("content or attachment is required")
	}

	conv, err := s.repo.FindOrCreateConversation(senderID, req.ReceiverID, req.EventID)
	if err != nil {
		return nil, err
	}

	msg := &Message{
		ConversationID:    conv.ID,
		SenderID:          senderID,
		Content:           content,
		AttachmentURL:     req.AttachmentURL,
		AttachmentType:    req.AttachmentType,
		AttachmentName:    req.AttachmentName,
		AttachmentMediaID: req.AttachmentMediaID,
		IsRead:            false,
	}
	if err := s.repo.CreateMessage(msg); err != nil {
		return nil, err
	}
	_ = s.repo.UpdateLastMessageAt(conv.ID, time.Now())

	// reload with sender
	if loaded, e := s.repo.GetMessageByID(msg.ID); e == nil {
		msg = loaded
	}
	res := msg.ToResponse()

	if s.hub != nil {
		out := dto.WSOutgoing{Type: "message", Payload: res}
		s.hub.SendToUser(req.ReceiverID, out)
		s.hub.SendToUser(senderID, out)
	}
	return res, nil
}

func (s *Service) EditMessage(userID, messageID uint, content string) (*dto.MessageResponse, error) {
	msg, err := s.repo.GetMessageByID(messageID)
	if err != nil {
		return nil, err
	}
	if msg.SenderID != userID {
		return nil, ErrForbidden
	}
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, errors.New("content is required")
	}
	msg.Content = content
	msg.IsEdited = true
	if err := s.repo.UpdateMessage(msg); err != nil {
		return nil, err
	}
	res := msg.ToResponse()
	if s.hub != nil {
		conv, _ := s.repo.GetConversationByID(msg.ConversationID)
		out := dto.WSOutgoing{Type: "message_updated", Payload: res}
		if conv != nil {
			s.hub.SendToUser(conv.User1ID, out)
			s.hub.SendToUser(conv.User2ID, out)
		}
	}
	return res, nil
}

func (s *Service) DeleteMessage(userID, messageID uint) error {
	msg, err := s.repo.GetMessageByID(messageID)
	if err != nil {
		return err
	}
	if msg.SenderID != userID {
		return ErrForbidden
	}
	convID := msg.ConversationID
	if err := s.repo.DeleteMessage(messageID); err != nil {
		return err
	}
	if s.hub != nil {
		conv, _ := s.repo.GetConversationByID(convID)
		out := dto.WSOutgoing{
			Type: "message_deleted",
			Payload: map[string]any{
				"id":              messageID,
				"conversation_id": convID,
			},
		}
		if conv != nil {
			s.hub.SendToUser(conv.User1ID, out)
			s.hub.SendToUser(conv.User2ID, out)
		}
	}
	return nil
}

func (s *Service) GetMyConversations(userID uint) ([]*dto.ConversationResponse, error) {
	list, err := s.repo.GetConversationsForUser(userID)
	if err != nil {
		return nil, err
	}
	result := make([]*dto.ConversationResponse, 0, len(list))
	for i := range list {
		resp, err := s.buildConversationResponse(&list[i], userID)
		if err != nil {
			continue
		}
		result = append(result, resp)
	}
	return result, nil
}

func (s *Service) GetMessages(userID, convID uint, beforeID uint, pageSize int) ([]*dto.MessageResponse, int64, bool, error) {
	conv, err := s.repo.GetConversationByID(convID)
	if err != nil {
		return nil, 0, false, err
	}
	if conv.User1ID != userID && conv.User2ID != userID {
		return nil, 0, false, ErrForbidden
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 30
	}
	msgs, total, hasMore, err := s.repo.GetMessages(convID, pageSize, beforeID)
	if err != nil {
		return nil, 0, false, err
	}
	// repo returns newest first; UI wants oldest -> newest
	for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
		msgs[i], msgs[j] = msgs[j], msgs[i]
	}
	res := make([]*dto.MessageResponse, 0, len(msgs))
	for i := range msgs {
		res = append(res, msgs[i].ToResponse())
	}
	return res, total, hasMore, nil
}

func (s *Service) MarkRead(userID, convID uint, messageIDs []uint) error {
	conv, err := s.repo.GetConversationByID(convID)
	if err != nil {
		return err
	}
	if conv.User1ID != userID && conv.User2ID != userID {
		return ErrForbidden
	}
	_, err = s.repo.MarkMessagesRead(convID, userID, messageIDs)
	if err != nil {
		return err
	}
	if s.hub != nil {
		otherID := conv.User1ID
		if otherID == userID {
			otherID = conv.User2ID
		}
		s.hub.SendToUser(otherID, dto.WSOutgoing{
			Type: "read",
			Payload: map[string]any{
				"conversation_id": convID,
				"message_ids":     messageIDs,
				"read_by":         userID,
			},
		})
	}
	return nil
}

func (s *Service) buildConversationResponse(conv *Conversation, currentUserID uint) (*dto.ConversationResponse, error) {
	var participant *dto.UserBrief
	if conv.User1ID == currentUserID {
		if conv.User2.ID != 0 {
			participant = &dto.UserBrief{
				ID: conv.User2.ID, Name: conv.User2.Name, Email: conv.User2.Email,
				Role: string(conv.User2.Role), ProfileImage: conv.User2.ProfileImage,
				ProfileImageId: conv.User2.ProfileImageId,
			}
		}
	} else if conv.User1.ID != 0 {
		participant = &dto.UserBrief{
			ID: conv.User1.ID, Name: conv.User1.Name, Email: conv.User1.Email,
			Role: string(conv.User1.Role), ProfileImage: conv.User1.ProfileImage,
			ProfileImageId: conv.User1.ProfileImageId,
		}
	}
	unread, _ := s.repo.UnreadCount(conv.ID, currentUserID)
	last, _ := s.repo.GetLastMessage(conv.ID)
	var lastResp *dto.MessageResponse
	if last != nil {
		lastResp = last.ToResponse()
	}
	updated := conv.UpdatedAt
	if conv.LastMessageAt != nil {
		updated = *conv.LastMessageAt
	}
	return &dto.ConversationResponse{
		ID:          conv.ID,
		Participant: participant,
		EventID:     conv.EventID,
		LastMessage: lastResp,
		UnreadCount: unread,
		UpdatedAt:   updated,
		CreatedAt:   conv.CreatedAt,
	}, nil
}

func (s *Service) HandleWSMessage(senderID uint, payload dto.WSMessagePayload) (*dto.MessageResponse, error) {
	return s.SendMessage(senderID, &dto.SendMessageRequest{
		ReceiverID: payload.ReceiverID,
		Content:    payload.Content,
		EventID:    payload.EventID,
	})
}

func (s *Service) HandleWSRead(userID uint, payload dto.WSReadPayload) error {
	return s.MarkRead(userID, payload.ConversationID, payload.MessageIDs)
}
