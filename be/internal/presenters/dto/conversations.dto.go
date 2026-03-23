package dto

import "github.com/longtk26/chat-app/pkg/dto"

// Create conversation

type CreateConversationRequestDto struct {
	Title   string   `json:"title"`
	Type    string   `json:"type"`
	UserIDs []string `json:"user_ids"`
}

type CreateConversationResponseDto struct {
	Conversation ConversationDto `json:"conversation"`
	Message      string          `json:"message"`
}

// List conversations

type ListConversationsQueryDto struct {
	dto.PaginationQueryDto
	UserID string `query:"user_id"`
}

type ListConversationsResponseDto struct {
	Conversations []ConversationDto `json:"conversations"`
	dto.PaginationResponseDto
}

type UserDto struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type ConversationDto struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Type          string    `json:"type"`
	Users         []UserDto `json:"users"`
	LastMessageID string    `json:"last_message_id,omitempty"`
}
