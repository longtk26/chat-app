"use client";

import { Conversation } from "@/lib/types";
import { Card, CardContent } from "../ui/card";
import { useRouter, useSearchParams } from "next/navigation";

export type TConversationItemProps = {
    conversation: Conversation;
    currentUserId?: string;
};

const ConversationItem = ({
    conversation,
    currentUserId,
}: TConversationItemProps) => {
    const searchParams = useSearchParams();
    const router = useRouter();
    const conversationName =
        conversation.type === "private"
            ? conversation.users.find((user) => user.id !== currentUserId)
                  ?.username || "Unknown User"
            : conversation.title || "Unnamed Group";
    const conversationIdParam = searchParams.get("conversation_id");
    const isActive = conversationIdParam === conversation.id;

    return (
        <Card
            className={`${isActive ? "border-blue-500 border-2" : ""} hover:cursor-pointer`}
            onClick={() => {
                router.push(`/?conversation_id=${conversation.id}`);
            }}
        >
            <CardContent className="flex flex-row items-center gap-4">
                <div className="w-12 h-12 rounded-full bg-blue-500/10 text-2xl relative">
                    <span className="absolute inset-0 flex items-center justify-center text-blue-500">
                        👨
                    </span>
                </div>
                <h2 className="font-bold">{conversationName}</h2>
            </CardContent>
        </Card>
    );
};

export default ConversationItem;
