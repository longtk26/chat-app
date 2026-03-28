"use client";

import { useAuth } from "@/hooks/useAuth";
import ConversationItem from "../conversation-item";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../ui/tabs";
import UserItem from "../user-item";
import { conversationsQuery } from "@/lib/query/conversations.query";
import { Conversation } from "@/lib/types";
import { usersQuery } from "@/lib/query/users.query";
import { useQuery } from "@tanstack/react-query";
import { useRouter, useSearchParams } from "next/navigation";
import { useEffect } from "react";
import { deleteCookie } from "@/utils/cookies";
import { LogOut, MessageCircle } from "lucide-react";
import { Button } from "../ui/button";
import { cn, getAvatarColor, getInitials } from "@/lib/utils";

const SideBarMessages = ({ children }: { children: React.ReactNode }) => {
    const { token, username, userId, isLoaded } = useAuth();
    const searchParams = useSearchParams();
    const router = useRouter();

    const { data: conversationsData, isFetching } = useQuery(
        conversationsQuery.listConversations(userId!),
    );
    const { data: usersData, isFetching: isUsersFetching } = useQuery(
        usersQuery.listUsers(),
    );

    const conversations = conversationsData?.conversations || [];
    const users = usersData || [];
    const firstConversationId =
        conversations.length > 0 ? conversations[0].id : null;
    const conversationIdParam = searchParams.get("conversation_id");
    const redirectConversationId = conversationIdParam || firstConversationId;

    useEffect(() => {
        if (redirectConversationId) {
            router.replace(`/?conversation_id=${redirectConversationId}`);
        }
    }, [redirectConversationId]);

    useEffect(() => {
        if (isLoaded && (!token || !userId)) {
            router.push("/login");
        }
    }, [isLoaded, token, userId, router]);

    const handleLogout = () => {
        deleteCookie("auth_token");
        deleteCookie("username");
        window.location.href = "/login";
    };

    if (!isLoaded || isFetching || isUsersFetching) {
        return (
            <div className="flex h-screen w-full">
                <aside className="w-72 flex-shrink-0 border-r border-slate-200 bg-white flex flex-col">
                    <div className="flex items-center gap-3 px-4 py-4 border-b border-slate-100">
                        <div className="h-9 w-9 rounded-full bg-slate-200 animate-pulse" />
                        <div className="h-4 w-24 rounded bg-slate-200 animate-pulse" />
                    </div>
                    <div className="flex flex-col gap-3 p-4">
                        {[...Array(5)].map((_, i) => (
                            <div key={i} className="flex items-center gap-3">
                                <div className="h-10 w-10 rounded-full bg-slate-200 animate-pulse flex-shrink-0" />
                                <div className="h-4 flex-1 rounded bg-slate-200 animate-pulse" />
                            </div>
                        ))}
                    </div>
                </aside>
                <main className="flex-1">{children}</main>
            </div>
        );
    }

    return (
        <div className="flex h-screen w-full overflow-hidden">
            {/* Sidebar */}
            <aside className="w-72 flex-shrink-0 border-r border-slate-200 bg-white flex flex-col h-full">
                {/* User header */}
                <div className="flex items-center justify-between px-4 py-3 border-b border-slate-100">
                    <div className="flex items-center gap-3">
                        <div
                            className={cn(
                                "h-9 w-9 rounded-full flex items-center justify-center text-sm font-semibold text-white flex-shrink-0",
                                getAvatarColor(username ?? ""),
                            )}
                        >
                            {getInitials(username ?? "?")}
                        </div>
                        <span className="font-semibold text-slate-800 text-sm truncate capitalize">
                            {username}
                        </span>
                    </div>
                    <Button
                        variant="ghost"
                        size="icon"
                        className="h-8 w-8 text-slate-400 hover:text-rose-500 hover:bg-rose-50"
                        title="Logout"
                        onClick={handleLogout}
                    >
                        <LogOut className="h-4 w-4" />
                    </Button>
                </div>

                {/* App name */}
                <div className="flex items-center gap-2 px-4 py-2 border-b border-slate-100">
                    <MessageCircle className="h-4 w-4 text-indigo-600" />
                    <span className="text-xs font-semibold text-indigo-600 tracking-wide uppercase">
                        Chats
                    </span>
                </div>

                {/* Tabs */}
                <Tabs defaultValue="chats" className="flex flex-col flex-1 overflow-hidden">
                    <TabsList className="flex w-full rounded-none border-b border-slate-100 bg-transparent p-0 h-auto shrink-0">
                        <TabsTrigger
                            value="chats"
                            className="flex-1 rounded-none border-b-2 border-transparent py-2.5 text-xs font-semibold text-slate-500 data-[state=active]:border-indigo-600 data-[state=active]:text-indigo-600 data-[state=active]:bg-transparent data-[state=active]:shadow-none"
                        >
                            Chats
                        </TabsTrigger>
                        <TabsTrigger
                            value="contacts"
                            className="flex-1 rounded-none border-b-2 border-transparent py-2.5 text-xs font-semibold text-slate-500 data-[state=active]:border-indigo-600 data-[state=active]:text-indigo-600 data-[state=active]:bg-transparent data-[state=active]:shadow-none"
                        >
                            Contacts
                        </TabsTrigger>
                    </TabsList>

                    <TabsContent
                        value="chats"
                        className="flex-1 overflow-y-auto m-0 p-0"
                    >
                        {conversations.length === 0 ? (
                            <p className="px-4 py-8 text-center text-sm text-slate-400">
                                No conversations yet.
                            </p>
                        ) : (
                            conversations.map((conversation: Conversation) => (
                                <ConversationItem
                                    key={conversation.id}
                                    conversation={conversation}
                                    currentUserId={userId ?? undefined}
                                />
                            ))
                        )}
                    </TabsContent>

                    <TabsContent
                        value="contacts"
                        className="flex-1 overflow-y-auto m-0 p-0"
                    >
                        {users.length === 0 ? (
                            <p className="px-4 py-8 text-center text-sm text-slate-400">
                                No contacts found.
                            </p>
                        ) : (
                            users.map((user) => (
                                <UserItem key={user.ID} user={user} />
                            ))
                        )}
                    </TabsContent>
                </Tabs>
            </aside>

            {/* Main content */}
            <main className="flex-1 overflow-hidden">{children}</main>
        </div>
    );
};

export default SideBarMessages;
