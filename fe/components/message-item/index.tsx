export type TMessageItem = {
    type: "me" | "other";
    message: string;
    username?: string;
    time?: string;
};

const MessageItem = ({ type, message, username, time }: TMessageItem) => {
    return (
        <div
            className={`bg-transparent flex items-center ${type === "me" ? "flex-row-reverse" : ""}`}
        >
            <div
                className={`flex ${type === "me" ? "" : "flex-row-reverse"} gap-2`}
            >
                <div className="flex flex-col gap-2">
                    <span>{username}</span>
                    <div
                        className={`bg-blue-500 text-white rounded-lg px-4 py-2 ${type === "me" ? "rounded-br-none" : "rounded-bl-none"}`}
                    >
                        {message}
                    </div>
                    <span>{time}</span>
                </div>
                <div className="w-12 h-12 rounded-full bg-blue-500/10 text-2xl relative flex">
                    <span className="absolute inset-0 flex items-center justify-center text-blue-500">
                        👨
                    </span>
                </div>
            </div>
        </div>
    );
};

export default MessageItem;
