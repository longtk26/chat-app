import { ListConversationsResponse } from "@/lib/types";
import { apiFetch } from "./client";

export const getConversations = async (
    userId: string,
): Promise<ListConversationsResponse> => {
    return apiFetch(`/api/v1/conversations?user_id=${userId}`);
};

export const createConversation = async (payload: {
    title?: string;
    type: "private" | "group";
    user_ids: string[];
}) => {
    return apiFetch("/api/v1/conversations", {
        method: "POST",
        body: JSON.stringify(payload),
    });
};
