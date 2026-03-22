"use client";

import { useUsersQuery } from "@/lib/query/users.query";
import { User } from "@/lib/types";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Avatar } from "@/components/ui/avatar";
import { useAuth } from "@/hooks/useAuth";

export const UserList = ({
    onSelectUser,
}: {
    onSelectUser: (user: User) => void;
}) => {
    const { data: users, isLoading, isError } = useUsersQuery();
    const { userId } = useAuth();

    if (isLoading) {
        return <div>Loading users...</div>;
    }

    if (isError) {
        return <div>Error loading users.</div>;
    }

    const filteredUsers = users?.filter((user) => user.ID !== userId);

    const getUserInitials = (name: string) => {
        const names = name.split(" ");
        return names.length > 1 ? names[0][0] + names[1][0] : names[0][0];
    };

    return (
        <Card>
            <CardHeader>
                <CardTitle>Users</CardTitle>
            </CardHeader>
            <CardContent className="flex flex-col gap-2">
                {filteredUsers?.map((user) => (
                    <Button
                        key={user.ID}
                        variant="ghost"
                        className="w-full justify-start flex items-center gap-2"
                        onClick={() => onSelectUser(user)}
                    >
                        {user.profileImage ? (
                            <Avatar src={user.profileImage} alt="User Image" />
                        ) : (
                            <div className="w-8 h-8 bg-gray-300 rounded-full flex items-center justify-center">
                                {getUserInitials(user.Username)}
                            </div>
                        )}
                        <span>{user.Username}</span>
                    </Button>
                ))}
            </CardContent>
        </Card>
    );
};
