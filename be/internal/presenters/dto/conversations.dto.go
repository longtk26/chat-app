package dto

type CreateConversationRequestDto struct {
	Title   string   `json:"title"`
	Type    string   `json:"type"`
	UserIDs []string `json:"user_ids"`
}

type CreateConversationResponseDto struct {
	Message string `json:"message"`
}
