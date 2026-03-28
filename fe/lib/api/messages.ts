import { ListMessagesResponse, Message } from "@/lib/types";
import { apiFetch } from "./client";

export type SendMessageVariables = {
    senderId: string;
    conversationId: string;
    content: string;
};

export type SendMessageResponse = {
    message: Message;
};

type GetMessagesParams = {
    conversationId: string;
    limit?: number;
    cursorTime?: string;
    lastMessageId?: string;
};

export const getMessages = async ({
    conversationId,
    limit = 20,
    cursorTime,
    lastMessageId,
}: GetMessagesParams): Promise<ListMessagesResponse> => {
    const query = new URLSearchParams();

    query.set("conversation_id", conversationId);
    query.set("limit", String(limit));

    if (cursorTime) {
        query.set("cursor_time", cursorTime);
    }

    if (lastMessageId) {
        query.set("last_message_id", lastMessageId);
    }

    return apiFetch(`/api/v1/messages?${query.toString()}`);
};

export const sendMessage = async (payload: SendMessageVariables): Promise<SendMessageResponse> => {
    const payloadToSend = {
        sender_id: payload.senderId,
        conversation_id: payload.conversationId,
        content: payload.content,
    };
    return apiFetch<SendMessageResponse>("/api/v1/messages", {
        method: "POST",
        body: JSON.stringify(payloadToSend),
    });
};
