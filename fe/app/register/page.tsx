import { AuthForm } from "@/components/auth-form";

const RegisterPage = () => {
    return (
        <main className="flex min-h-screen items-center justify-center bg-gradient-to-br from-slate-50 via-white to-indigo-50">
            <AuthForm type="register" />
        </main>
    );
};

export default RegisterPage;
