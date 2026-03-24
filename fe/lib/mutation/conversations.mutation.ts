import { useMutation, useQueryClient } from "@tanstack/react-query";
import { createConversation } from "../api/conversations";
import { useAuth } from "@/hooks/useAuth";

export const conversationMutation = () => {
    const queryClient = useQueryClient();
    const { userId } = useAuth();

    return {
        createConversation: useMutation({
            mutationFn: createConversation,
            onSuccess: () => {
                queryClient.invalidateQueries({
                    queryKey: ["conversations", userId],
                });
            },
        }),
    };
};
