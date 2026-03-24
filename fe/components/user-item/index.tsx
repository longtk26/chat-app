"use client";

import { User } from "@/lib/types";
import { Card, CardContent } from "../ui/card";
import { Button } from "../ui/button";
import { MessageCircle } from "lucide-react";
import { conversationMutation } from "@/lib/mutation/conversations.mutation";
import { useAuth } from "@/hooks/useAuth";
import { useRouter } from "next/navigation";

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
                    console.log("Conversation created:", data);
                    const conversationId = (data as any).conversation.id;
                    router.push(`?conversation_id=${conversationId}`);
                },
            },
        );
    };

    return (
        <Card>
            <CardContent className="flex flex-row items-center gap-4">
                <div className="flex flex-row items-center gap-4">
                    <div className="w-12 h-12 rounded-full bg-blue-500/10 text-2xl relative">
                        <span className="absolute inset-0 flex items-center justify-center text-blue-500">
                            👨
                        </span>
                    </div>
                    <h2 className="font-bold">{user.Username}</h2>
                </div>
                <Button
                    variant="outline"
                    size="sm"
                    className="flex items-center gap-2 hover:bg-blue-600 hover:text-white"
                    onClick={handleOpenChat}
                >
                    <MessageCircle />
                </Button>
            </CardContent>
        </Card>
    );
};

export default UserItem;
