"use client";
import { useState, useEffect } from "react";

export const useAuth = () => {
    const [userId, setUserId] = useState<string | null>(null);
    const [token, setToken] = useState<string | null>(null);

    useEffect(() => {
        const storedToken = localStorage.getItem("auth_token");
        if (storedToken) {
            setToken(storedToken);
            const id = storedToken.split("_")[1];
            setUserId(id);
        }
    }, []);

    return { userId, token };
};
