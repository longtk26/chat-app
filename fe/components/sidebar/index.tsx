"use client";

import { useAuth } from "@/hooks/useAuth";
import ConversationItem from "../conversation-item";
import { Button } from "../ui/button";
import { Separator } from "../ui/separator";
import {
    SidebarContent,
    SidebarHeader,
    SidebarProvider,
    SidebarTrigger,
} from "../ui/sidebar";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../ui/tabs";
import UserItem from "../user-item";
import { conversationsQuery } from "@/lib/query/conversations.query";
import { Conversation } from "@/lib/types";
import { usersQuery } from "@/lib/query/users.query";
import { useQuery } from "@tanstack/react-query";
import { useRouter, useSearchParams } from "next/navigation";
import { useEffect } from "react";

const SideBarMessages = ({ children }: { children: React.ReactNode }) => {
    const { token, username, userId } = useAuth();
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

    if (!token || !userId) {
        return <div>Redirecting...</div>;
    }

    if (isFetching || isUsersFetching) {
        return <div>Loading...</div>;
    }

    return (
        <SidebarProvider>
            <section className="w-78 h-full border-r">
                <SidebarHeader className="flex flex-row w-full items-center justify-between px-2">
                    <div className="flex flex-col gap-4">
                        <h1 className="font-bold text-lg">Messages</h1>
                        <p className="capitalize">{username}</p>
                    </div>
                    <Button
                        variant="outline"
                        className="hover:bg-blue-600 hover:text-white"
                        onClick={() => {
                            localStorage.removeItem("auth_token");
                            localStorage.removeItem("username");
                            localStorage.removeItem("user_id");
                            window.location.href = "/login";
                        }}
                    >
                        Logout
                    </Button>
                </SidebarHeader>
                <Separator />
                <SidebarContent>
                    <Tabs defaultValue="chats" className="w-full">
                        <TabsList className="w-full border-b bg-transparent">
                            <TabsTrigger
                                value="chats"
                                className="data-[state=active]:text-white data-[state=active]:bg-blue-600"
                            >
                                Chats
                            </TabsTrigger>
                            <TabsTrigger
                                value="contacts"
                                className="data-[state=active]:text-white data-[state=active]:bg-blue-600"
                            >
                                Contacts
                            </TabsTrigger>
                        </TabsList>
                        <TabsContent
                            value="chats"
                            className="p-4 flex flex-col gap-4"
                        >
                            {conversations.map((conversation: Conversation) => (
                                <ConversationItem
                                    key={conversation.id}
                                    conversation={conversation}
                                    currentUserId={userId}
                                />
                            )) || <p>No conversations found.</p>}
                        </TabsContent>
                        <TabsContent
                            value="contacts"
                            className="p-4 flex flex-col gap-4"
                        >
                            {users.map((user) => (
                                <UserItem key={user.ID} user={user} />
                            )) || <p>No users found.</p>}
                        </TabsContent>
                    </Tabs>
                </SidebarContent>
            </section>
            <SidebarTrigger />
            {children}
        </SidebarProvider>
    );
};

export default SideBarMessages;
