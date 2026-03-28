import { ChatView } from "@/components/chat-view";
import SideBarMessages from "@/components/sidebar";
import SocketProvider from "@/provider/socket.provider";

export default function Home() {
    return (
        <main className="h-screen overflow-hidden">
            <SocketProvider>
                <SideBarMessages>
                    <ChatView />
                </SideBarMessages>
            </SocketProvider>
        </main>
    );
}
