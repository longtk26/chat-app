"use client";
import { useState, useEffect } from "react";

export const useAuth = () => {
    const [userId, setUserId] = useState<string | null>(null);
    const [token, setToken] = useState<string | null>(null);
    const [username, setUsername] = useState<string | null>(null);
 
    useEffect(() => {
        const storedToken = localStorage.getItem("auth_token");
        const userName = localStorage.getItem("username");
        if (storedToken) {
            setToken(storedToken);
            const id = storedToken.split("_")[1];
            setUserId(id);
        }
        if (userName) {
            setUsername(userName);
        }
    }, []);

    return { userId, token, username };
};
