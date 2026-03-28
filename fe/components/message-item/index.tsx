import { cn, getAvatarColor, getInitials } from "@/lib/utils";

export type TMessageItem = {
    type: "me" | "other";
    message: string;
    username?: string;
    time?: string;
};

const MessageItem = ({ type, message, username, time }: TMessageItem) => {
    const isMe = type === "me";

    return (
        <div
            className={cn(
                "flex items-end gap-2",
                isMe ? "flex-row-reverse" : "flex-row",
            )}
        >
            {/* Avatar */}
            {!isMe && (
                <div
                    className={cn(
                        "h-7 w-7 rounded-full flex items-center justify-center text-xs font-semibold text-white flex-shrink-0 mb-1",
                        getAvatarColor(username ?? ""),
                    )}
                >
                    {getInitials(username ?? "?")}
                </div>
            )}

            {/* Bubble + meta */}
            <div
                className={cn(
                    "flex flex-col gap-1 max-w-[70%]",
                    isMe ? "items-end" : "items-start",
                )}
            >
                {!isMe && username && (
                    <span className="text-xs font-medium text-slate-500 px-1">
                        {username}
                    </span>
                )}
                <div
                    className={cn(
                        "px-4 py-2.5 rounded-2xl text-sm leading-relaxed break-words",
                        isMe
                            ? "bg-indigo-600 text-white rounded-br-sm"
                            : "bg-white text-slate-800 border border-slate-200 rounded-bl-sm shadow-sm",
                    )}
                >
                    {message}
                </div>
                {time && (
                    <span className="text-[11px] text-slate-400 px-1">
                        {time}
                    </span>
                )}
            </div>
        </div>
    );
};

export default MessageItem;
