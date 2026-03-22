"use client";

import { useConversationsQuery } from "@/lib/query/conversations.query";
import { Conversation } from "@/lib/types";
import { useAuth } from "@/hooks/useAuth";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Avatar } from "@/components/ui/avatar";

export const ConversationList = ({
    onSelectConversation,
}: {
    onSelectConversation: (conversation: Conversation) => void;
}) => {
    const { userId } = useAuth();
    const {
        data: conversationsResponse,
        isLoading,
        isError,
    } = useConversationsQuery(userId!);

    if (!userId) {
        return <div>Please log in.</div>;
    }

    if (isLoading) {
        return <div>Loading conversations...</div>;
    }

    if (isError) {
        return <div>Error loading conversations.</div>;
    }

    const conversations = conversationsResponse?.conversations || [];

    const getConversationTitle = (conv: Conversation) => {
        if (conv.title) return conv.title;
        if (conv.type === "private") {
            const otherUser = conv.users.find((u) => u.id !== userId);
            return otherUser?.username || "Private Chat";
        }
        return conv.users.map((u) => u.username).join(", ");
    };

    const getUserInitials = (name: string) => {
        const names = name.split(" ");
        return names.length > 1 ? names[0][0] + names[1][0] : names[0][0];
    };

    return (
        <Card>
            <CardHeader>
                <CardTitle>Conversations</CardTitle>
            </CardHeader>
            <CardContent className="flex flex-col gap-2">
                {conversations.length === 0 ? (
                    <div className="text-center text-gray-500">
                        No conversations yet. Start a new one!
                    </div>
                ) : (
                    conversations.map((conv) => (
                        <Button
                            key={conv.id}
                            variant="ghost"
                            className="w-full justify-start flex items-center gap-2"
                            onClick={() => onSelectConversation(conv)}
                        >
                            {conv.users[0]?.profileImage ? (
                                <Avatar
                                    src={conv.users[0].profileImage}
                                    alt="User Image"
                                />
                            ) : (
                                <div className="w-8 h-8 bg-gray-300 rounded-full flex items-center justify-center">
                                    {getUserInitials(
                                        getConversationTitle(conv),
                                    )}
                                </div>
                            )}
                            <span>{getConversationTitle(conv)}</span>
                        </Button>
                    ))
                )}
            </CardContent>
        </Card>
    );
};
