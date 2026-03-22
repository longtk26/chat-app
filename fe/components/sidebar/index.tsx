"use client";

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

const SideBarMessages = ({ children }: { children: React.ReactNode }) => {
    return (
        <SidebarProvider>
            <section className="w-78 h-full border-r">
                <SidebarHeader className="flex flex-row w-full items-center justify-between px-2">
                    <div className="flex flex-col gap-4">
                        <h1 className="font-bold text-lg">Messages</h1>
                        <p>Your name</p>
                    </div>
                    <Button
                        variant="outline"
                        className="hover:bg-blue-600 hover:text-white"
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
                            <ConversationItem />
                            <ConversationItem />
                            <ConversationItem />
                        </TabsContent>
                        <TabsContent
                            value="contacts"
                            className="p-4 flex flex-col gap-4"
                        >
                            <UserItem />
                            <UserItem />
                            <UserItem />
                            <UserItem />
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
