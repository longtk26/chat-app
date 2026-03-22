import { getUsers } from "@/lib/api/users";

export const usersQuery = {
    listUsers: () => {
        return {
            queryKey: ["users"],
            queryFn: () => getUsers(),
        };
    },
};
