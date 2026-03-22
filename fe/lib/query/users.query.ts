import { getUsers } from "@/lib/api/users";
import { useQuery } from "@tanstack/react-query";

export const useUsersQuery = () => {
    return useQuery({
        queryKey: ["users"],
        queryFn: getUsers,
    });
};
