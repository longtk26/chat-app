import { InfiniteData, useMutation, useQueryClient } from "@tanstack/react-query";
import { sendMessage, SendMessageResponse, SendMessageVariables } from "../api/messages";
import { ListMessagesResponse } from "../types";

export const messagesMutation = () => {
    const queryClient = useQueryClient();

    return {
        createMessage: useMutation<SendMessageResponse, Error, SendMessageVariables>({
            mutationFn: sendMessage,
            onSuccess: (data, variables) => {
                // Directly add the sent message to the cache instead of invalidating.
                // invalidateQueries would trigger a refetch that can race with incoming
                // socket events, causing the refetch to overwrite messages received via socket.
                queryClient.setQueryData<InfiniteData<ListMessagesResponse>>(
                    ["messages", variables.conversationId, 20],
                    (oldData) => {
                        if (!oldData) return oldData;
                        const exists = oldData.pages.some((page) =>
                            page.messages.some((m) => m.id === data.message.id),
                        );
                        if (exists) return oldData;
                        return {
                            ...oldData,
                            pages: [
                                {
                                    ...oldData.pages[0],
                                    messages: [data.message, ...oldData.pages[0].messages],
                                },
                                ...oldData.pages.slice(1),
                            ],
                        };
                    },
                );
            },
        }),
    };
};
