import { useMutation } from "@tanstack/react-query";
import { login, register } from "../api/auth";

export const useAuthMutation = () => {
    const { mutate: loginMutate } = useMutation({
        mutationFn: login,
    });
    const { mutate: registerMutate } = useMutation({
        mutationFn: register,
    });

    return {
        login: loginMutate,
        register: registerMutate,
    };
};
