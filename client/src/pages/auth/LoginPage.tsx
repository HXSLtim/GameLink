import { useState } from 'react';
import { useNavigate, useLocation, Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useAuthStore } from '@/stores';
import { Button } from '@/components/ui/button';
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Checkbox } from "@/components/ui/checkbox"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Loader2, Gamepad2, Ghost, Rocket, Mail, Smartphone } from "lucide-react"
import { LanguageSwitcher } from '@/components/language-switcher';
import { ModeToggle } from '@/components/mode-toggle';

// Storage key for remembered username only (NEVER store passwords)
const REMEMBER_USERNAME_KEY = 'gamelink_remembered_username';

interface LocationState {
    from?: { pathname: string };
}

export default function LoginPage() {
    const navigate = useNavigate();
    const location = useLocation();
    const { login, register, loading, error: storeError } = useAuthStore();
    const { t } = useTranslation();

    const [isRegister, setIsRegister] = useState(false);
    const [validationErr, setValidationErr] = useState<string | null>(null);

    // Load remembered username only (NEVER store passwords)
    const getSavedUsername = (): string | null => {
        try {
            return localStorage.getItem(REMEMBER_USERNAME_KEY);
        } catch {
            return null;
        }
    };

    const savedUsername = getSavedUsername();

    // Login fields
    const [username, setUsername] = useState(savedUsername || '');
    const [password, setPassword] = useState('');
    const [rememberMe, setRememberMe] = useState(!!savedUsername);

    // Register additional fields
    const [regMethod, setRegMethod] = useState<'email' | 'phone'>('email');
    const [email, setEmail] = useState('');
    const [phone, setPhone] = useState('');
    const [nickname, setNickname] = useState('');

    // Get return url from location state or default to home
    const from = (location.state as LocationState)?.from?.pathname || '/';

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setValidationErr(null);

        try {
            if (isRegister) {
                // Validation: Password length
                if (password.length < 6) {
                    setValidationErr(t('auth.password_min_length', { defaultValue: "Password must be at least 6 characters" }));
                    return;
                }

                if (regMethod === 'email' && !email) {
                    setValidationErr(t('auth.email_required', { defaultValue: 'Email is required' }));
                    return;
                }
                if (regMethod === 'phone' && !phone) {
                    setValidationErr(t('auth.phone_required', { defaultValue: 'Phone number is required' }));
                    return;
                }

                await register({
                    phone: regMethod === 'phone' ? phone : undefined,
                    email: regMethod === 'email' ? email : undefined,
                    password,
                    name: nickname || (regMethod === 'email' ? email.split('@')[0] : `User${phone.slice(-4)}`)
                });
            } else {
                await login({ username, password });

                // Save or clear remembered username only (NEVER store passwords)
                try {
                    if (rememberMe) {
                        localStorage.setItem(REMEMBER_USERNAME_KEY, username);
                    } else {
                        localStorage.removeItem(REMEMBER_USERNAME_KEY);
                    }
                } catch {
                    // Ignore localStorage errors (quota exceeded, disabled, etc.)
                }
            }
            navigate(from, { replace: true });
        } catch {
            // Error handled in store
        }
    };

    const toggleMode = () => {
        setIsRegister(!isRegister);
        setValidationErr(null);
    };

    const error = validationErr || storeError;

    return (
        <div className="w-full h-screen min-h-[600px] lg:grid lg:grid-cols-2 relative transition-colors duration-300">
            {/* Language Switcher & Mode Toggle - Absolute positioned */}
            <div className="absolute top-4 right-4 z-50 flex items-center gap-2">
                <LanguageSwitcher />
                <ModeToggle />
            </div>

            {/* Left Column - Visuals */}
            <div className="hidden lg:flex flex-col relative bg-muted text-white dark:border-r overflow-hidden">
                <div className="absolute inset-0 bg-zinc-900" />
                {/* Abstract Gradient Background */}
                <div className="absolute inset-0 bg-gradient-to-br from-violet-600/30 via-zinc-900 to-indigo-900/30" />
                <div className="absolute inset-0 bg-[url('https://images.unsplash.com/photo-1542751371-adc38448a05e?q=80&w=2670&auto=format&fit=crop')] bg-cover bg-center opacity-10 mix-blend-overlay" />

                {/* Branding Content */}
                <div className="relative z-20 flex items-center text-lg font-medium p-10">
                    <Gamepad2 className="mr-2 h-6 w-6" />
                    {t('app.name')}
                </div>

                <div className="relative z-20 mt-auto p-10 space-y-4">
                    <blockquote className="space-y-2">
                        <p className="text-xl font-semibold leading-relaxed">
                            {t('auth.slogan')}
                        </p>
                    </blockquote>
                    <div className="flex gap-4 pt-4 text-sm text-muted-foreground/60">
                        <div className="flex items-center gap-1"><Ghost className="h-4 w-4" /> {t('auth.anonymous')}</div>
                        <div className="flex items-center gap-1"><Rocket className="h-4 w-4" /> {t('auth.secure')}</div>
                    </div>
                </div>
            </div>

            {/* Right Column - Form */}
            <div className="flex items-center justify-center py-12 px-4 sm:px-8 bg-background">
                <div className="mx-auto grid w-full max-w-[380px] gap-6">
                    <div className="flex flex-col space-y-2 text-center">
                        <h1 className="text-3xl font-bold tracking-tight">
                            {isRegister ? t('auth.create_account') : t('auth.welcome_back')}
                        </h1>
                        <p className="text-sm text-muted-foreground">
                            {isRegister ? t('auth.enter_credentials') : t('auth.enter_credentials')}
                        </p>
                    </div>

                    <form onSubmit={handleSubmit} className="grid gap-4">
                        {error && (
                            <div className="p-3 text-sm font-medium text-destructive bg-destructive/10 rounded-md border border-destructive/20 animate-in fade-in slide-in-from-top-1">
                                {error}
                            </div>
                        )}

                        {/* Login Mode */}
                        {!isRegister && (
                            <>
                                <div className="grid gap-2">
                                    <Label htmlFor="username">{t('auth.username')}</Label>
                                    <Input
                                        id="username"
                                        placeholder={t('auth.username')}
                                        type="text"
                                        autoCapitalize="none"
                                        autoComplete="username"
                                        autoCorrect="off"
                                        value={username}
                                        onChange={(e) => setUsername(e.target.value)}
                                        disabled={loading}
                                        className="bg-background/50 backdrop-blur-sm"
                                        required={!isRegister}
                                    />
                                </div>
                                <div className="grid gap-2">
                                    <div className="flex items-center justify-between">
                                        <Label htmlFor="password">{t('auth.password')}</Label>
                                        <Link
                                            to="/forgot-password"
                                            className="text-sm text-primary underline-offset-4 hover:underline"
                                        >
                                            {t('auth.forgot_password')}
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
                                <div className="flex items-center space-x-2">
                                    <Checkbox
                                        id="remember"
                                        checked={rememberMe}
                                        onCheckedChange={(checked) => setRememberMe(checked as boolean)}
                                    />
                                    <Label htmlFor="remember" className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                                        {t('auth.remember_me', { defaultValue: 'Remember me' })}
                                    </Label>
                                </div>
                            </>
                        )}

                        {/* Register Mode */}
                        {isRegister && (
                            <Tabs defaultValue="email" value={regMethod} onValueChange={(v) => setRegMethod(v as 'email' | 'phone')} className="w-full">
                                <TabsList className="grid w-full grid-cols-2 mb-4">
                                    <TabsTrigger value="email">
                                        <Mail className="h-4 w-4 mr-2" />
                                        {t('auth.email', { defaultValue: 'Email' })}
                                    </TabsTrigger>
                                    <TabsTrigger value="phone">
                                        <Smartphone className="h-4 w-4 mr-2" />
                                        {t('auth.phone', { defaultValue: 'Phone' })}
                                    </TabsTrigger>
                                </TabsList>

                                <div className="space-y-4">
                                    <div className="grid gap-2">
                                        <Label htmlFor="nickname">{t('auth.nickname')} <span className="text-destructive">*</span></Label>
                                        <Input
                                            id="nickname"
                                            placeholder={t('auth.nickname')}
                                            type="text"
                                            value={nickname}
                                            onChange={(e) => setNickname(e.target.value)}
                                            disabled={loading}
                                            className="bg-background/50 backdrop-blur-sm"
                                            required={isRegister}
                                        />
                                    </div>

                                    <TabsContent value="email" className="space-y-4 m-0">
                                        <div className="grid gap-2 animate-in fade-in slide-in-from-top-2">
                                            <Label htmlFor="email">{t('auth.email')} <span className="text-destructive">*</span></Label>
                                            <Input
                                                id="email"
                                                placeholder="name@example.com"
                                                type="email"
                                                autoComplete="email"
                                                value={email}
                                                onChange={(e) => setEmail(e.target.value)}
                                                disabled={loading}
                                                className="bg-background/50 backdrop-blur-sm"
                                                required={regMethod === 'email'}
                                            />
                                        </div>
                                    </TabsContent>

                                    <TabsContent value="phone" className="space-y-4 m-0">
                                        <div className="grid gap-2 animate-in fade-in slide-in-from-top-2">
                                            <Label htmlFor="phone">{t('auth.phone')} <span className="text-destructive">*</span></Label>
                                            <Input
                                                id="phone"
                                                placeholder="13800000000"
                                                type="tel"
                                                value={phone}
                                                onChange={(e) => setPhone(e.target.value)}
                                                disabled={loading}
                                                className="bg-background/50 backdrop-blur-sm"
                                                required={regMethod === 'phone'}
                                            />
                                        </div>
                                    </TabsContent>

                                    <div className="grid gap-2">
                                        <Label htmlFor="password">{t('auth.password')} <span className="text-destructive">*</span></Label>
                                        <Input
                                            id="password"
                                            placeholder="••••••••"
                                            type="password"
                                            autoComplete="new-password"
                                            value={password}
                                            onChange={(e) => setPassword(e.target.value)}
                                            disabled={loading}
                                            className="bg-background/50 backdrop-blur-sm"
                                            required
                                        />
                                        <p className="text-[10px] text-muted-foreground">{t('auth.password_min_length')}</p>
                                    </div>
                                </div>
                            </Tabs>
                        )}

                        <Button disabled={loading} className="w-full mt-2 bg-gradient-to-r from-violet-600 to-indigo-600 hover:from-violet-700 hover:to-indigo-700 transition-all shadow-lg hover:shadow-violet-500/25">
                            {loading ? (
                                <>
                                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                    {isRegister ? t('auth.create_account') : t('auth.sign_in')}...
                                </>
                            ) : (
                                isRegister ? t('auth.create_account') : t('auth.sign_in')
                            )}
                        </Button>
                    </form>

                    <div className="relative">
                        <div className="absolute inset-0 flex items-center">
                            <span className="w-full border-t" />
                        </div>
                        <div className="relative flex justify-center text-xs uppercase">
                            <span className="bg-background px-2 text-muted-foreground">
                                {t('auth.continue_with')}
                            </span>
                        </div>
                    </div>

                    <Button variant="outline" type="button" disabled={loading} className="w-full" onClick={toggleMode}>
                        {isRegister ? t('auth.already_have_account') : t('auth.create_account')}
                    </Button>

                    <p className="px-8 text-center text-sm text-muted-foreground">
                        {t('auth.terms_agreement')}{" "}
                        <Link to="/terms" className="underline underline-offset-4 hover:text-primary">
                            {t('auth.terms')}
                        </Link>{" "}
                        {t('auth.and')}{" "}
                        <Link to="/privacy" className="underline underline-offset-4 hover:text-primary">
                            {t('auth.privacy')}
                        </Link>
                        .
                    </p>
                </div>
            </div>
        </div>
    );
}
