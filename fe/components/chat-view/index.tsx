"use client";

import { Controller, useForm } from "react-hook-form";
import MessageItem from "../message-item";
import { Field, FieldError, FieldGroup, FieldLabel } from "../ui/field";
import { Input } from "../ui/input";
import { Button } from "../ui/button";
import { Textarea } from "../ui/textarea";

export const ChatView = () => {
    const form = useForm<{ message: string }>({
        defaultValues: {
            message: "",
        },
    });

    return (
        <section className="w-full h-screen bg-gray-100 border-l relative">
            <div className="flex items-center gap-4 p-4 w-full bg-white border-b">
                <div className="w-12 h-12 rounded-full bg-blue-500/10 text-2xl relative flex flex-col">
                    <span className="absolute inset-0 flex items-center justify-center text-blue-500">
                        👨
                    </span>
                </div>
                <div className="flex flex-col gap-2 justify-center">
                    <span>User name</span>
                    <span>Active</span>
                </div>
            </div>
            <section className="flex flex-col gap-4 p-4 w-full h-[calc(100%-150px)] overflow-y-auto">
                <MessageItem type="me" />
                <MessageItem type="other" />
            </section>
            <form className="flex justify-between items-center gap-4 p-4 w-full bg-white border-t absolute bottom-0">
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
