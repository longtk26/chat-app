"use client";
import { useState, useEffect } from "react";
import { getCookie } from "@/utils/cookies";

export const useAuth = () => {
    const [userId, setUserId] = useState<string | null>(null);
    const [token, setToken] = useState<string | null>(null);
    const [username, setUsername] = useState<string | null>(null);
    const [isLoaded, setIsLoaded] = useState(false);

    useEffect(() => {
        const storedToken = getCookie("auth_token");
        const userName = getCookie("username");
        if (storedToken) {
            setToken(storedToken);
            const id = storedToken.split("_")[1];
            setUserId(id);
        }
        if (userName) {
            setUsername(userName);
        }
        setIsLoaded(true);
    }, []);

    return { userId, token, username, isLoaded };
};
