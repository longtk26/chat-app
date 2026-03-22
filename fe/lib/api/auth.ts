import { apiFetch } from "./client";

export const login = async ({
    username,
    password,
}: {
    username: string;
    password: string;
}) => {
    return apiFetch<{ access_token: string; username: string }>(
        "/api/v1/auth/login",
        {
            method: "POST",
            body: JSON.stringify({ username, password }),
        },
    );
};

export const register = async ({
    username,
    password,
}: {
    username: string;
    password: string;
}) => {
    return apiFetch("/api/v1/auth/register", {
        method: "POST",
        body: JSON.stringify({ username, password }),
    });
};
