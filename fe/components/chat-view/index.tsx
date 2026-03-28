"use client";

import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { Controller, useForm } from "react-hook-form";
import MessageItem from "../message-item";
import { Field, FieldError, FieldGroup } from "../ui/field";
import { Button } from "../ui/button";
import { Textarea } from "../ui/textarea";
import { useSearchParams } from "next/navigation";
import {
    InfiniteData,
    useInfiniteQuery,
    useQuery,
    useQueryClient,
} from "@tanstack/react-query";
import { conversationsQuery } from "@/lib/query/conversations.query";
import { useAuth } from "@/hooks/useAuth";
import { messagesQuery } from "@/lib/query/messages.query";
import { ListMessagesResponse, Message } from "@/lib/types";
import { messagesMutation } from "@/lib/mutation/messages.mutation";
import { useSocket } from "@/provider/socket.provider";

const MESSAGE_PAGE_LIMIT = 20;

export const ChatView = () => {
    const searchParams = useSearchParams();
    const form = useForm<{ message: string }>({
        defaultValues: { message: "" },
    });
    const { userId } = useAuth();
    const conversationId = searchParams.get("conversation_id");
    const socket = useSocket();
    const queryClient = useQueryClient();

    const [isOtherUserTyping, setIsOtherUserTyping] = useState(false);

    const { data: conversation } = useQuery(
        conversationsQuery.getConversation(conversationId!),
    );

    const {
        data: messagesData,
        fetchNextPage,
        hasNextPage,
        isFetchingNextPage,
        isPending: isMessagesPending,
    } = useInfiniteQuery(
        messagesQuery.listMessages(conversationId ?? "", MESSAGE_PAGE_LIMIT),
    );

    const messagesScrollRef = useRef<HTMLElement | null>(null);
    const hasInitialScrollRef = useRef(false);

    const users = conversation?.users ?? [];
    const otherUser = users.find((user) => user.id !== userId);

    const messages = useMemo(() => {
        const pages = messagesData?.pages ?? [];
        const newestToOldest = pages.flatMap((page) => page.messages);
        const seen = new Set<string>();

        const uniqueNewestToOldest = newestToOldest.filter((message) => {
            if (seen.has(message.id)) return false;
            seen.add(message.id);
            return true;
        });

        return uniqueNewestToOldest.reverse();
    }, [messagesData]);

    // Join/leave conversation room on mount and when conversation changes
    useEffect(() => {
        if (!conversationId) return;
        socket.emit("join_conversation", conversationId);
        return () => {
            socket.emit("leave_conversation", conversationId);
        };
    }, [socket, conversationId]);

    // Handle incoming real-time events
    useEffect(() => {
        if (!conversationId) return;

        const handleNewMessage = (data: unknown) => {
            const message = data as Message;
            queryClient.setQueryData<InfiniteData<ListMessagesResponse>>(
                ["messages", conversationId, MESSAGE_PAGE_LIMIT],
                (oldData) => {
                    if (!oldData) return oldData;
                    // Deduplicate: skip if message already exists
                    const exists = oldData.pages.some((page) =>
                        page.messages.some((m) => m.id === message.id),
                    );
                    if (exists) return oldData;
                    return {
                        ...oldData,
                        pages: [
                            {
                                ...oldData.pages[0],
                                messages: [message, ...oldData.pages[0].messages],
                            },
                            ...oldData.pages.slice(1),
                        ],
                    };
                },
            );
        };

        const handleUserTyping = (data: unknown) => {
            const { user_id } = data as { user_id: string };
            if (user_id !== userId) {
                setIsOtherUserTyping(true);
            }
        };

        const handleUserStopTyping = () => {
            setIsOtherUserTyping(false);
        };

        socket.on("new_message", handleNewMessage);
        socket.on("user_typing", handleUserTyping);
        socket.on("user_stop_typing", handleUserStopTyping);

        return () => {
            socket.off("new_message", handleNewMessage);
            socket.off("user_typing", handleUserTyping);
            socket.off("user_stop_typing", handleUserStopTyping);
        };
    }, [socket, conversationId, userId, queryClient]);

    const loadOlderMessages = useCallback(async () => {
        const container = messagesScrollRef.current;
        if (!container || !hasNextPage || isFetchingNextPage) return;

        const prevScrollTop = container.scrollTop;
        const prevScrollHeight = container.scrollHeight;

        await fetchNextPage();

        requestAnimationFrame(() => {
            const nextScrollHeight = container.scrollHeight;
            container.scrollTop =
                nextScrollHeight - prevScrollHeight + prevScrollTop;
        });
    }, [fetchNextPage, hasNextPage, isFetchingNextPage]);

    const handleMessagesScroll = useCallback(
        (event: React.UIEvent<HTMLElement>) => {
            if (event.currentTarget.scrollTop <= 80) {
                void loadOlderMessages();
            }
        },
        [loadOlderMessages],
    );

    useEffect(() => {
        hasInitialScrollRef.current = false;
        setIsOtherUserTyping(false);
    }, [conversationId]);

    useEffect(() => {
        const container = messagesScrollRef.current;
        if (!container || messages.length === 0) return;

        if (!hasInitialScrollRef.current) {
            requestAnimationFrame(() => {
                container.scrollTop = container.scrollHeight;
                hasInitialScrollRef.current = true;
            });
            return;
        }

        // Auto-scroll to bottom when new messages arrive if already near the bottom
        const distanceFromBottom =
            container.scrollHeight - container.scrollTop - container.clientHeight;
        if (distanceFromBottom < 100) {
            requestAnimationFrame(() => {
                container.scrollTop = container.scrollHeight;
            });
        }
    }, [messages.length]);

    const formatTime = useCallback((time: string) => {
        const parsedDate = new Date(time);
        if (Number.isNaN(parsedDate.getTime())) return "";
        return parsedDate.toLocaleTimeString([], {
            hour: "2-digit",
            minute: "2-digit",
        });
    }, []);

    const {
        createMessage: { mutate: createMessage },
    } = messagesMutation();

    const handleSendMessage = (data: { message: string }) => {
        if (!conversationId) return;
        createMessage({
            conversationId,
            content: data.message,
            senderId: userId!,
        });
        form.reset();
    };

    // Typing indicator: emit "typing" on input, then "stop_typing" after 2 s of silence
    const typingTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

    const handleTyping = useCallback(() => {
        if (!conversationId) return;
        socket.emit("typing", conversationId);
        if (typingTimeoutRef.current) clearTimeout(typingTimeoutRef.current);
        typingTimeoutRef.current = setTimeout(() => {
            socket.emit("stop_typing", conversationId);
        }, 2000);
    }, [socket, conversationId]);

    return (
        <section className="w-full h-screen bg-gray-100 border-l relative">
            <div className="flex items-center gap-4 p-4 w-full bg-white border-b">
                <div className="w-12 h-12 rounded-full bg-blue-500/10 text-2xl relative flex flex-col">
                    <span className="absolute inset-0 flex items-center justify-center text-blue-500">
                        👨
                    </span>
                </div>
                <div className="flex flex-col gap-1 justify-center">
                    <span className="font-medium">{otherUser?.username}</span>
                    <span className="text-xs text-gray-500">
                        {isOtherUserTyping ? "Typing..." : "Active"}
                    </span>
                </div>
            </div>

            <section
                ref={messagesScrollRef}
                onScroll={handleMessagesScroll}
                className="flex flex-col gap-4 p-4 w-full h-[calc(100%-150px)] overflow-y-auto"
            >
                {isFetchingNextPage && (
                    <span className="text-xs text-center text-gray-500">
                        Loading older messages...
                    </span>
                )}

                {isMessagesPending && messages.length === 0 && (
                    <span className="text-sm text-gray-500 text-center">
                        Loading messages...
                    </span>
                )}

                {!isMessagesPending && messages.length === 0 && (
                    <span className="text-sm text-gray-500 text-center">
                        No messages yet
                    </span>
                )}

                {messages.map((message: Message) => (
                    <MessageItem
                        key={message.id}
                        type={message.sender_id === userId ? "me" : "other"}
                        username={message.sender_name}
                        message={message.content}
                        time={formatTime(message.sent_at)}
                    />
                ))}
            </section>

            <form
                onSubmit={form.handleSubmit(handleSendMessage)}
                className="flex justify-between items-center gap-4 p-4 w-full bg-white border-t absolute bottom-0"
            >
                <FieldGroup>
                    <Controller
                        name="message"
                        control={form.control}
                        render={({ field, fieldState }) => (
                            <Field data-invalid={fieldState.invalid}>
                                <Textarea
                                    id="message"
                                    placeholder="Enter your message"
                                    required
                                    {...field}
                                    onChange={(e) => {
                                        field.onChange(e);
                                        handleTyping();
                                    }}
                                />
                                {fieldState.invalid && (
                                    <FieldError errors={[fieldState.error]} />
                                )}
                            </Field>
                        )}
                    />
                </FieldGroup>
                <Button type="submit">Send</Button>
            </form>
        </section>
    );
};
