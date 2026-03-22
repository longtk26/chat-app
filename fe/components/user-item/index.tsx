import { User } from "@/lib/types";
import { Card, CardContent } from "../ui/card";
import { Button } from "../ui/button";
import { MessageCircle } from "lucide-react";

export type TUserItemProps = {
    user: User;
};

const UserItem = ({ user }: TUserItemProps) => {
    return (
        <Card>
            <CardContent className="flex flex-row items-center gap-4">
                <div className="flex flex-row items-center gap-4">
                    <div className="w-12 h-12 rounded-full bg-blue-500/10 text-2xl relative">
                        <span className="absolute inset-0 flex items-center justify-center text-blue-500">
                            👨
                        </span>
                    </div>
                    <h2 className="font-bold">{user.Username}</h2>
                </div>
                <Button
                    variant="outline"
                    size="sm"
                    className="flex items-center gap-2 hover:bg-blue-600 hover:text-white"
                >
                    <MessageCircle />
                </Button>
            </CardContent>
        </Card>
    );
};

export default UserItem;
