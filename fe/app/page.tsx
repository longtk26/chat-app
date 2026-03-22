import { ChatView } from "@/components/chat-view";
import SideBarMessages from "@/components/sidebar";

export default function Home() {
    return (
        <main className="flex">
            <SideBarMessages>
                <ChatView />
            </SideBarMessages>
        </main>
    );
}
