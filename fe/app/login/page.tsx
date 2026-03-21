import { AuthForm } from "@/components/auth-form";

const LoginPage = () => {
    return (
        <main className="flex min-h-screen items-center justify-center">
            <AuthForm type="login" />
        </main>
    );
};

export default LoginPage;
