"use client";

import { Button } from "@/components/ui/button";
import {
    Card,
    CardAction,
    CardContent,
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Field, FieldError, FieldGroup, FieldLabel } from "../ui/field";
import { Controller, useForm } from "react-hook-form";
import Link from "next/link";
import { useAuthMutation } from "@/lib/mutation/auth.mutation";
import { toast } from "sonner";
import { useRouter } from "next/navigation";
import { setCookie } from "@/utils/cookies";
import { MessageCircle } from "lucide-react";

export type TAuthFormProps = {
    type: "login" | "register";
};

export type TAuthFormValues = {
    username: string;
    password: string;
    confirmPassword?: string;
};

export function AuthForm({ type }: TAuthFormProps) {
    const isLogin = type === "login";
    const form = useForm<TAuthFormValues>({
        defaultValues: {
            username: "",
            password: "",
            confirmPassword: "",
        },
    });
    const { login, register } = useAuthMutation();
    const router = useRouter();

    const onSubmit = (payload: TAuthFormValues) => {
        if (isLogin) {
            login(payload, {
                onSuccess: (data) => {
                    setCookie("auth_token", data.access_token);
                    setCookie("username", data.username);
                    toast.success("Logged in successfully!");
                    router.push("/");
                },
                onError: () => {
                    toast.error("Invalid credentials. Please try again.");
                },
            });
        } else {
            const { confirmPassword, ...registerPayload } = payload;
            if (confirmPassword !== payload.password) {
                form.setError("confirmPassword", {
                    type: "manual",
                    message: "Passwords do not match",
                });
                return;
            }
            register(registerPayload, {
                onSuccess: () => {
                    toast.success("Account created! You can now log in.");
                    router.push("/login");
                },
                onError: () => {
                    toast.error("Something went wrong. Please try again.");
                },
            });
        }
    };

    return (
        <Card className="w-full max-w-sm shadow-xl border border-slate-100">
            <CardHeader className="pb-4 text-center">
                <div className="mx-auto mb-3 flex h-12 w-12 items-center justify-center rounded-2xl bg-indigo-600">
                    <MessageCircle className="h-6 w-6 text-white" />
                </div>
                <CardTitle className="text-xl font-bold text-slate-900">
                    {isLogin ? "Welcome back" : "Create account"}
                </CardTitle>
                <CardDescription className="text-slate-500 text-sm">
                    {isLogin
                        ? "Sign in to continue to your chats."
                        : "Fill in the details below to get started."}
                </CardDescription>
                <CardAction>
                    <Link
                        href={isLogin ? "/register" : "/login"}
                        className="text-sm font-medium text-indigo-600 hover:text-indigo-700 hover:underline underline-offset-4"
                    >
                        {isLogin ? "Create account" : "Sign in instead"}
                    </Link>
                </CardAction>
            </CardHeader>
            <CardContent>
                <form onSubmit={form.handleSubmit(onSubmit)} className="w-full">
                    <FieldGroup>
                        <div className="flex flex-col gap-4">
                            <Controller
                                name="username"
                                control={form.control}
                                render={({ field, fieldState }) => (
                                    <Field
                                        className="grid gap-1.5"
                                        data-invalid={fieldState.invalid}
                                    >
                                        <FieldLabel
                                            htmlFor="username"
                                            className="text-sm font-medium text-slate-700"
                                        >
                                            Username
                                        </FieldLabel>
                                        <Input
                                            id="username"
                                            type="text"
                                            placeholder="Enter your username"
                                            required
                                            className="h-10 border-slate-200 focus-visible:ring-indigo-500"
                                            {...field}
                                        />
                                        {fieldState.invalid && (
                                            <FieldError errors={[fieldState.error]} />
                                        )}
                                    </Field>
                                )}
                            />
                            <Controller
                                name="password"
                                control={form.control}
                                render={({ field, fieldState }) => (
                                    <Field
                                        className="grid gap-1.5"
                                        data-invalid={fieldState.invalid}
                                    >
                                        <FieldLabel
                                            htmlFor="password"
                                            className="text-sm font-medium text-slate-700"
                                        >
                                            Password
                                        </FieldLabel>
                                        <Input
                                            id="password"
                                            type="password"
                                            placeholder="Enter your password"
                                            required
                                            className="h-10 border-slate-200 focus-visible:ring-indigo-500"
                                            {...field}
                                        />
                                        {fieldState.invalid && (
                                            <FieldError errors={[fieldState.error]} />
                                        )}
                                    </Field>
                                )}
                            />
                            {!isLogin ? (
                                <Controller
                                    name="confirmPassword"
                                    control={form.control}
                                    render={({ field, fieldState }) => (
                                        <Field
                                            className="grid gap-1.5"
                                            data-invalid={fieldState.invalid}
                                        >
                                            <FieldLabel
                                                htmlFor="confirmPassword"
                                                className="text-sm font-medium text-slate-700"
                                            >
                                                Confirm Password
                                            </FieldLabel>
                                            <Input
                                                id="confirmPassword"
                                                type="password"
                                                placeholder="Confirm your password"
                                                required
                                                className="h-10 border-slate-200 focus-visible:ring-indigo-500"
                                                {...field}
                                            />
                                            {fieldState.invalid && (
                                                <FieldError errors={[fieldState.error]} />
                                            )}
                                        </Field>
                                    )}
                                />
                            ) : null}
                        </div>
                    </FieldGroup>
                    <CardFooter className="flex-col gap-2 px-0 pt-6 pb-0">
                        <Button
                            type="submit"
                            className="w-full h-10 bg-indigo-600 hover:bg-indigo-700 text-white font-medium"
                        >
                            {isLogin ? "Sign in" : "Create account"}
                        </Button>
                    </CardFooter>
                </form>
            </CardContent>
        </Card>
    );
}
