"use client";

import { Conversation } from "@/lib/types";
import { cn, getAvatarColor, getInitials } from "@/lib/utils";
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
        <button
            type="button"
            onClick={() => router.push(`/?conversation_id=${conversation.id}`)}
            className={cn(
                "flex w-full items-center gap-3 px-4 py-3 text-left transition-colors",
                isActive
                    ? "bg-indigo-50 border-r-2 border-indigo-600"
                    : "hover:bg-slate-50 border-r-2 border-transparent",
            )}
        >
            <div
                className={cn(
                    "h-10 w-10 rounded-full flex items-center justify-center text-sm font-semibold text-white flex-shrink-0",
                    getAvatarColor(conversationName),
                )}
            >
                {getInitials(conversationName)}
            </div>
            <div className="flex flex-col min-w-0">
                <span
                    className={cn(
                        "text-sm font-semibold truncate",
                        isActive ? "text-indigo-700" : "text-slate-800",
                    )}
                >
                    {conversationName}
                </span>
                <span className="text-xs text-slate-400 truncate">
                    Tap to open chat
                </span>
            </div>
        </button>
    );
};

export default ConversationItem;
