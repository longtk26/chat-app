import { Conversation } from "@/lib/types";
import { Card, CardContent } from "../ui/card";

export type TConversationItemProps = {
    conversation: Conversation;
    currentUserId?: string;
};

const ConversationItem = ({
    conversation,
    currentUserId,
}: TConversationItemProps) => {
    const conversationName =
        conversation.type === "private"
            ? conversation.users.find((user) => user.id !== currentUserId)
                  ?.username || "Unknown User"
            : conversation.title || "Unnamed Group";

    return (
        <Card>
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
