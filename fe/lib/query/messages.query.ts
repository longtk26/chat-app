import { getMessages } from "@/lib/api/messages";

type MessagesPageParam = {
    cursorTime?: string;
    lastMessageId?: string;
};

export const messagesQuery = {
    listMessages: (conversationId: string, limit = 20) => ({
        queryKey: ["messages", conversationId, limit],
        queryFn: ({ pageParam }: { pageParam: MessagesPageParam }) =>
            getMessages({
                conversationId,
                limit,
                cursorTime: pageParam?.cursorTime,
                lastMessageId: pageParam?.lastMessageId,
            }),
        enabled: !!conversationId,
        initialPageParam: {} as MessagesPageParam,
        getNextPageParam: (lastPage: {
            has_more: boolean;
            next_cursor_time?: string;
            messages: Array<{ id: string }>;
        }) => {
            if (!lastPage.has_more || lastPage.messages.length === 0) {
                return undefined;
            }

            const oldestMessage =
                lastPage.messages[lastPage.messages.length - 1];

            return {
                cursorTime: lastPage.next_cursor_time,
                lastMessageId: oldestMessage.id,
            } satisfies MessagesPageParam;
        },
    }),
};
