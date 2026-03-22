import { getConversations } from "@/lib/api/conversations";

export const conversationsQuery = {
    listConversations: (userId: string) => ({
        queryKey: ["conversations", userId],
        queryFn: () => getConversations(userId),
        enabled: !!userId,
    }),
};
