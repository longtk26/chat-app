package dto

import baseDto "github.com/longtk26/chat-app/pkg/dto"

type ListMessagesQueryDto struct {
	baseDto.CursorPaginationQueryDto
	ConversationID string `query:"conversation_id"`
	LastMessageID  string `query:"last_message_id"`
}

type MessageDto struct {
	ID             string `json:"id"`
	SenderID       string `json:"sender_id"`
	SenderName     string `json:"sender_name"`
	ConversationID string `json:"conversation_id"`
	Content        string `json:"content"`
	SentAt         string `json:"sent_at"`
	UpdatedAt      string `json:"updated_at"`
}

type ListMessagesResponseDto struct {
	Messages []MessageDto `json:"messages"`
	baseDto.CursorPaginationResponseDto
}

type SendMessageRequestDto struct {
	SenderID       string `json:"sender_id"`
	ConversationID string `json:"conversation_id"`
	Content        string `json:"content"`
}

type SendMessageResponseDto struct {
	Message MessageDto `json:"message"`
}

type UpdateMessageRequestDto struct {
	Content string `json:"content"`
}

type UpdateMessageResponseDto struct {
	Message MessageDto `json:"message"`
}
