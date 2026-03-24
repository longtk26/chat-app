"use client";

import { useCallback, useEffect, useMemo, useRef } from "react";
import { Controller, useForm } from "react-hook-form";
import MessageItem from "../message-item";
import { Field, FieldError, FieldGroup } from "../ui/field";
import { Button } from "../ui/button";
import { Textarea } from "../ui/textarea";
import { useSearchParams } from "next/navigation";
import { useInfiniteQuery, useQuery } from "@tanstack/react-query";
import { conversationsQuery } from "@/lib/query/conversations.query";
import { useAuth } from "@/hooks/useAuth";
import { messagesQuery } from "@/lib/query/messages.query";
import { Message } from "@/lib/types";
import { messagesMutation } from "@/lib/mutation/messages.mutation";

const MESSAGE_PAGE_LIMIT = 20;

export const ChatView = () => {
    const searchParams = useSearchParams();
    const form = useForm<{ message: string }>({
        defaultValues: {
            message: "",
        },
    });
    const { userId } = useAuth();
    const conversationId = searchParams.get("conversation_id");

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
            if (seen.has(message.id)) {
                return false;
            }

            seen.add(message.id);
            return true;
        });

        return uniqueNewestToOldest.reverse();
    }, [messagesData]);

    const loadOlderMessages = useCallback(async () => {
        const container = messagesScrollRef.current;
        if (!container || !hasNextPage || isFetchingNextPage) {
            return;
        }

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
            const { scrollTop } = event.currentTarget;

            if (scrollTop <= 80) {
                void loadOlderMessages();
            }
        },
        [loadOlderMessages],
    );

    useEffect(() => {
        hasInitialScrollRef.current = false;
    }, [conversationId]);

    useEffect(() => {
        const container = messagesScrollRef.current;
        if (
            !container ||
            hasInitialScrollRef.current ||
            messages.length === 0
        ) {
            return;
        }

        requestAnimationFrame(() => {
            container.scrollTop = container.scrollHeight;
            hasInitialScrollRef.current = true;
        });
    }, [messages.length]);

    const formatTime = useCallback((time: string) => {
        const parsedDate = new Date(time);
        if (Number.isNaN(parsedDate.getTime())) {
            return "";
        }

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
    };

    return (
        <section className="w-full h-screen bg-gray-100 border-l relative">
            <div className="flex items-center gap-4 p-4 w-full bg-white border-b">
                <div className="w-12 h-12 rounded-full bg-blue-500/10 text-2xl relative flex flex-col">
                    <span className="absolute inset-0 flex items-center justify-center text-blue-500">
                        👨
                    </span>
                </div>
                <div className="flex flex-col gap-2 justify-center">
                    <span>{otherUser?.username}</span>
                    <span>Active</span>
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
                                />
                                {fieldState.invalid && (
                                    <FieldError errors={[fieldState.error]} />
                                )}
                            </Field>
                        )}
                    />
                </FieldGroup>
                <Button type="submit" className="">
                    Send
                </Button>
            </form>
        </section>
    );
};
