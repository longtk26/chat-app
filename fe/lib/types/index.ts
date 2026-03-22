// User from GET /api/v1/users
export type User = {
    ID: string;
    Username: string;
    Email: string;
};

// User inside a conversation object
export type ConversationParticipant = {
    id: string;
    username: string;
};

export type Conversation = {
    id: string;
    title: string;
    type: "private" | "group";
    users: ConversationParticipant[];
};

export type ListConversationsResponse = {
    conversations: Conversation[];
    total_count: number;
    page: number;
};
