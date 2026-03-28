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
import { cn, getAvatarColor, getInitials } from "@/lib/utils";
import { Send, MessageCircle } from "lucide-react";

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

    const typingTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

    const handleTyping = useCallback(() => {
        if (!conversationId) return;
        socket.emit("typing", conversationId);
        if (typingTimeoutRef.current) clearTimeout(typingTimeoutRef.current);
        typingTimeoutRef.current = setTimeout(() => {
            socket.emit("stop_typing", conversationId);
        }, 2000);
    }, [socket, conversationId]);

    // Empty state — no conversation selected
    if (!conversationId) {
        return (
            <section className="flex flex-col items-center justify-center w-full h-full bg-slate-50 text-center">
                <div className="flex flex-col items-center gap-4 text-slate-400">
                    <div className="h-16 w-16 rounded-2xl bg-indigo-100 flex items-center justify-center">
                        <MessageCircle className="h-8 w-8 text-indigo-400" />
                    </div>
                    <div>
                        <p className="font-semibold text-slate-600">No chat selected</p>
                        <p className="text-sm mt-1">Choose a conversation from the sidebar to start chatting.</p>
                    </div>
                </div>
            </section>
        );
    }

    return (
        <section className="flex flex-col w-full h-full bg-slate-50">
            {/* Header */}
            <div className="flex items-center gap-3 px-5 py-3 bg-white border-b border-slate-200 shadow-sm shrink-0">
                <div
                    className={cn(
                        "h-9 w-9 rounded-full flex items-center justify-center text-sm font-semibold text-white flex-shrink-0",
                        getAvatarColor(otherUser?.username ?? ""),
                    )}
                >
                    {getInitials(otherUser?.username ?? "?")}
                </div>
                <div className="flex flex-col">
                    <span className="font-semibold text-slate-900 text-sm leading-tight">
                        {otherUser?.username ?? "Unknown"}
                    </span>
                    <span
                        className={cn(
                            "text-xs leading-tight transition-colors",
                            isOtherUserTyping
                                ? "text-indigo-500 font-medium"
                                : "text-slate-400",
                        )}
                    >
                        {isOtherUserTyping ? "typing..." : "online"}
                    </span>
                </div>
            </div>

            {/* Messages */}
            <section
                ref={messagesScrollRef}
                onScroll={handleMessagesScroll}
                className="flex flex-col gap-3 px-5 py-4 flex-1 overflow-y-auto"
            >
                {isFetchingNextPage && (
                    <span className="text-xs text-center text-slate-400 py-1">
                        Loading older messages…
                    </span>
                )}

                {isMessagesPending && messages.length === 0 && (
                    <div className="flex flex-col items-center justify-center flex-1 text-slate-400">
                        <span className="text-sm">Loading messages…</span>
                    </div>
                )}

                {!isMessagesPending && messages.length === 0 && (
                    <div className="flex flex-col items-center justify-center flex-1 gap-2 text-slate-400">
                        <MessageCircle className="h-8 w-8 text-slate-300" />
                        <span className="text-sm">No messages yet. Say hello!</span>
                    </div>
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

            {/* Input area */}
            <form
                onSubmit={form.handleSubmit(handleSendMessage)}
                className="flex items-end gap-3 px-5 py-3 bg-white border-t border-slate-200 shrink-0"
            >
                <FieldGroup className="flex-1">
                    <Controller
                        name="message"
                        control={form.control}
                        render={({ field, fieldState }) => (
                            <Field data-invalid={fieldState.invalid}>
                                <Textarea
                                    id="message"
                                    placeholder="Type a message…"
                                    required
                                    rows={1}
                                    className="resize-none min-h-[40px] max-h-32 border-slate-200 focus-visible:ring-indigo-500 rounded-2xl py-2.5 px-4 text-sm leading-relaxed"
                                    {...field}
                                    onChange={(e) => {
                                        field.onChange(e);
                                        handleTyping();
                                    }}
                                    onKeyDown={(e) => {
                                        if (e.key === "Enter" && !e.shiftKey) {
                                            e.preventDefault();
                                            form.handleSubmit(handleSendMessage)();
                                        }
                                    }}
                                />
                                {fieldState.invalid && (
                                    <FieldError errors={[fieldState.error]} />
                                )}
                            </Field>
                        )}
                    />
                </FieldGroup>
                <Button
                    type="submit"
                    size="icon"
                    className="h-10 w-10 rounded-full bg-indigo-600 hover:bg-indigo-700 shrink-0 mb-0.5"
                >
                    <Send className="h-4 w-4" />
                </Button>
            </form>
        </section>
    );
};
