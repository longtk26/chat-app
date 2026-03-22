import { useMutation, useQueryClient } from "@tanstack/react-query";
import { createConversation } from "../api/conversations";
import { useAuth } from "@/hooks/useAuth";

export const useConversationMutation = () => {
    const queryClient = useQueryClient();
    const { userId } = useAuth();

    const { mutate: createConversationMutate, ...rest } = useMutation({
        mutationFn: createConversation,
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["conversations", userId] });
        },
    });

    return {
        createConversation: createConversationMutate,
        ...rest,
    };
};
