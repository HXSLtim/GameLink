import { useState } from 'react';
import { useNavigate, useLocation, Link } from 'react-router-dom';
import { useAuthStore } from '@/stores';
import { Button } from '@/components/ui/button';
import { LanguageSwitcher } from '@/components/language-switcher';
import { ModeToggle } from '@/components/mode-toggle';
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Loader2, Gamepad2, Ghost, Rocket } from "lucide-react"

export default function LoginPage() {
    const navigate = useNavigate();
    const location = useLocation();
    const { login, loading, error } = useAuthStore();

    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');

    // Get return url from location state or default to home
    const from = (location.state as any)?.from?.pathname || '/';

    const handleLogin = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            await login({ username, password });
            navigate(from, { replace: true });
        } catch (err) {
            // Error handled in store
        }
    };

    return (
        <div className="w-full h-screen min-h-[600px] lg:grid lg:grid-cols-2 relative transition-colors duration-300">
            {/* Language Switcher & Mode Toggle - Absolute positioned */}
            <div className="absolute top-4 right-4 z-50 flex items-center gap-2">
                <LanguageSwitcher />
                <ModeToggle />
            </div>

            {/* Left Column - Visuals */}
            <div className="hidden lg:flex flex-col relative bg-muted text-white dark:border-r">
                <div className="absolute inset-0 bg-zinc-900" />
                {/* Abstract Gradient Background */}
                <div className="absolute inset-0 bg-gradient-to-br from-violet-600/30 via-zinc-900 to-indigo-900/30" />
                <div className="absolute inset-0 bg-[url('https://images.unsplash.com/photo-1542751371-adc38448a05e?q=80&w=2670&auto=format&fit=crop')] bg-cover bg-center opacity-10 mix-blend-overlay" />

                {/* Branding Content */}
                <div className="relative z-20 flex items-center text-lg font-medium p-10">
                    <Gamepad2 className="mr-2 h-6 w-6" />
                    GameLink
                </div>

                <div className="relative z-20 mt-auto p-10 space-y-4">
                    <blockquote className="space-y-2">
                        <p className="text-xl font-semibold leading-relaxed">
                            "Connect with pro players, elevate your skills, and experience gaming like never before."
                        </p>
                    </blockquote>
                    <div className="flex gap-4 pt-4 text-sm text-muted-foreground/60">
                        <div className="flex items-center gap-1"><Ghost className="h-4 w-4" /> Anonymous</div>
                        <div className="flex items-center gap-1"><Rocket className="h-4 w-4" /> Fast & Secure</div>
                    </div>
                </div>
            </div>

            {/* Right Column - Form */}
            <div className="flex items-center justify-center py-12 px-4 sm:px-8">
                <div className="mx-auto grid w-full max-w-[380px] gap-6">
                    <div className="flex flex-col space-y-2 text-center">
                        <h1 className="text-3xl font-bold tracking-tight">Welcome Back</h1>
                        <p className="text-sm text-muted-foreground">
                            Enter your credentials to access your account
                        </p>
                    </div>

                    <form onSubmit={handleLogin} className="grid gap-4">
                        {error && (
                            <div className="p-3 text-sm font-medium text-destructive bg-destructive/10 rounded-md border border-destructive/20 animate-in fade-in slide-in-from-top-1">
                                {error}
                            </div>
                        )}
                        <div className="grid gap-2">
                            <Label htmlFor="username">Username / Email</Label>
                            <Input
                                id="username"
                                placeholder="name@example.com"
                                type="text"
                                autoCapitalize="none"
                                autoComplete="username"
                                autoCorrect="off"
                                value={username}
                                onChange={(e) => setUsername(e.target.value)}
                                disabled={loading}
                                className="bg-background/50 backdrop-blur-sm"
                                required
                            />
                        </div>
                        <div className="grid gap-2">
                            <div className="flex items-center justify-between">
                                <Label htmlFor="password">Password</Label>
                                <Link
                                    to="/forgot-password"
                                    className="text-sm text-primary underline-offset-4 hover:underline"
                                >
                                    Forgot password?
                                </Link>
                            </div>
                            <Input
                                id="password"
                                placeholder="••••••••"
                                type="password"
                                autoComplete="current-password"
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                                disabled={loading}
                                className="bg-background/50 backdrop-blur-sm"
                                required
                            />
                        </div>
                        <Button disabled={loading} className="w-full bg-gradient-to-r from-violet-600 to-indigo-600 hover:from-violet-700 hover:to-indigo-700 transition-all shadow-lg hover:shadow-violet-500/25">
                            {loading ? (
                                <>
                                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                    Signing In...
                                </>
                            ) : (
                                'Sign In'
                            )}
                        </Button>
                    </form>

                    <div className="relative">
                        <div className="absolute inset-0 flex items-center">
                            <span className="w-full border-t" />
                        </div>
                        <div className="relative flex justify-center text-xs uppercase">
                            <span className="bg-background px-2 text-muted-foreground">
                                Or continue with
                            </span>
                        </div>
                    </div>

                    <Button variant="outline" type="button" disabled={loading} className="w-full">
                        Create an account
                    </Button>

                    <p className="px-8 text-center text-sm text-muted-foreground">
                        By clicking continue, you agree to our{" "}
                        <Link to="/terms" className="underline underline-offset-4 hover:text-primary">
                            Terms
                        </Link>{" "}
                        and{" "}
                        <Link to="/privacy" className="underline underline-offset-4 hover:text-primary">
                            Privacy Policy
                        </Link>
                        .
                    </p>
                </div>
            </div>
        </div>
    );
}
