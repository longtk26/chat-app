import { useMutation, useQueryClient } from "@tanstack/react-query";
import { sendMessage } from "../api/messages";

export const messagesMutation = () => {
    const queryClient = useQueryClient();

    return {
        createMessage: useMutation({
            mutationFn: sendMessage,
            onSuccess: () => {
                queryClient.invalidateQueries({
                    queryKey: ["messages"],
                });
            },
        }),
    };
};
