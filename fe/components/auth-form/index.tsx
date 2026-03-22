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
                    localStorage.setItem("auth_token", data.access_token);
                    localStorage.setItem("username", data.username);
                    toast.success("Login successful! Redirecting...");
                    router.push("/");
                },
                onError: () => {
                    toast.error(
                        "Something wrong happens. Please check your credentials and try again.",
                    );
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
                    toast.success(
                        "Registration successful! You can now log in.",
                    );
                    router.push("/");
                },
                onError: () => {
                    toast.error("Something wrong happens. Please try again.");
                },
            });
        }
    };

    return (
        <Card className="w-full max-w-sm">
            <CardHeader>
                <CardTitle>
                    {isLogin ? "Welcome back!" : "Create an account"}
                </CardTitle>
                <CardDescription>
                    {isLogin
                        ? "Enter your credentials to access your account."
                        : "Fill in the details to create your account."}
                </CardDescription>
                <CardAction>
                    <Link
                        href={isLogin ? "/register" : "/login"}
                        className="text-sm text-primary underline-offset-4 hover:underline"
                    >
                        {isLogin ? "Register" : "Login"}
                    </Link>
                </CardAction>
            </CardHeader>
            <CardContent>
                <form onSubmit={form.handleSubmit(onSubmit)} className="w-full">
                    <FieldGroup>
                        <div className="flex flex-col gap-6">
                            <Controller
                                name="username"
                                control={form.control}
                                render={({ field, fieldState }) => (
                                    <Field
                                        className="grid gap-2"
                                        data-invalid={fieldState.invalid}
                                    >
                                        <FieldLabel htmlFor="username">
                                            Username
                                        </FieldLabel>
                                        <Input
                                            id="username"
                                            type="text"
                                            placeholder="Enter your username"
                                            required
                                            {...field}
                                        />
                                        {fieldState.invalid && (
                                            <FieldError
                                                errors={[fieldState.error]}
                                            />
                                        )}
                                    </Field>
                                )}
                            />
                            <Controller
                                name="password"
                                control={form.control}
                                render={({ field, fieldState }) => (
                                    <Field
                                        className="grid gap-2"
                                        data-invalid={fieldState.invalid}
                                    >
                                        <div className="flex items-center">
                                            <FieldLabel htmlFor="password">
                                                Password
                                            </FieldLabel>
                                        </div>
                                        <Input
                                            id="password"
                                            type="password"
                                            placeholder="Enter your password"
                                            required
                                            {...field}
                                        />
                                        {fieldState.invalid && (
                                            <FieldError
                                                errors={[fieldState.error]}
                                            />
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
                                            className="grid gap-2"
                                            data-invalid={fieldState.invalid}
                                        >
                                            <div className="flex items-center">
                                                <FieldLabel htmlFor="confirmPassword">
                                                    Confirm Password
                                                </FieldLabel>
                                            </div>
                                            <Input
                                                id="confirmPassword"
                                                type="password"
                                                required
                                                placeholder="Confirm your password"
                                                {...field}
                                            />
                                            {fieldState.invalid && (
                                                <FieldError
                                                    errors={[fieldState.error]}
                                                />
                                            )}
                                        </Field>
                                    )}
                                />
                            ) : null}
                        </div>
                    </FieldGroup>
                    <CardFooter className="flex-col gap-2">
                        <Button type="submit" className="w-full">
                            {isLogin ? "Login" : "Register"}
                        </Button>
                    </CardFooter>
                </form>
            </CardContent>
        </Card>
    );
}
