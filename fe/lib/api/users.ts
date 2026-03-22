import { User } from "@/lib/types";
import { apiFetch } from "./client";

export const getUsers = async (): Promise<User[]> => {
    return apiFetch("/api/v1/users");
};
