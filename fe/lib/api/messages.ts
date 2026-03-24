import { ListMessagesResponse } from "@/lib/types";
import { apiFetch } from "./client";

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

export const sendMessage = async (payload: {
    senderId: string;
    conversationId: string;
    content: string;
}) => {
    const payloadToSend = {
        sender_id: payload.senderId,
        conversation_id: payload.conversationId,
        content: payload.content,
    };
    return apiFetch("/api/v1/messages", {
        method: "POST",
        body: JSON.stringify(payloadToSend),
    });
};
