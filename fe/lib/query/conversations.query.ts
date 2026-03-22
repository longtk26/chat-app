import { getConversations } from "@/lib/api/conversations";
import { useQuery } from "@tanstack/react-query";

export const useConversationsQuery = (userId: string) => {
    return useQuery({
        queryKey: ["conversations", userId],
        queryFn: () => getConversations(userId),
        enabled: !!userId,
    });
};
