"use client";

import { User } from "@/lib/types";
import { Button } from "../ui/button";
import { MessageCircle } from "lucide-react";
import { conversationMutation } from "@/lib/mutation/conversations.mutation";
import { useAuth } from "@/hooks/useAuth";
import { useRouter } from "next/navigation";
import { cn, getAvatarColor, getInitials } from "@/lib/utils";

export type TUserItemProps = {
    user: User;
};

const UserItem = ({ user }: TUserItemProps) => {
    const {
        createConversation: { mutate },
    } = conversationMutation();
    const { userId } = useAuth();
    const router = useRouter();

    const handleOpenChat = () => {
        mutate(
            {
                user_ids: [user.ID, userId!],
                type: "private",
            },
            {
                onSuccess: (data) => {
                    const conversationId = (data as any).conversation.id;
                    router.push(`?conversation_id=${conversationId}`);
                },
            },
        );
    };

    return (
        <div className="flex items-center justify-between px-4 py-3 hover:bg-slate-50 transition-colors">
            <div className="flex items-center gap-3">
                <div
                    className={cn(
                        "h-10 w-10 rounded-full flex items-center justify-center text-sm font-semibold text-white flex-shrink-0",
                        getAvatarColor(user.Username),
                    )}
                >
                    {getInitials(user.Username)}
                </div>
                <span className="text-sm font-semibold text-slate-800">
                    {user.Username}
                </span>
            </div>
            <Button
                variant="ghost"
                size="icon"
                className="h-8 w-8 text-slate-400 hover:text-indigo-600 hover:bg-indigo-50"
                title="Open chat"
                onClick={handleOpenChat}
            >
                <MessageCircle className="h-4 w-4" />
            </Button>
        </div>
    );
};

export default UserItem;
