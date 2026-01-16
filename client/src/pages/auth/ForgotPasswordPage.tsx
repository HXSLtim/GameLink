import { useState, useEffect } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Button } from '@/components/ui/button';
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { ArrowLeft, Loader2, Mail, CheckCircle2, KeyRound, ShieldCheck } from "lucide-react";
import { toast } from 'sonner';
import { http } from '@/lib/http';

export default function ForgotPasswordPage() {
    const navigate = useNavigate();
    const { t } = useTranslation();

    // Steps: 0 = Email, 1 = Verify Code, 2 = Reset Password, 3 = Success
    const [step, setStep] = useState(0);
    const [email, setEmail] = useState('');
    const [code, setCode] = useState('');
    const [newPassword, setNewPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');

    const [loading, setLoading] = useState(false);
    const [countdown, setCountdown] = useState(0);

    useEffect(() => {
        let timer: NodeJS.Timeout;
        if (countdown > 0) {
            timer = setInterval(() => setCountdown(c => c - 1), 1000);
        }
        return () => clearInterval(timer);
    }, [countdown]);

    const handleSendCode = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!email) return;

        setLoading(true);
        try {
            await http.post('/public/verification/send', {
                target: email,
                type: 'email',
                scene: 'reset_password' // Optional, dependent on backend generic service
            });
            toast.success(t('auth.code_sent', { defaultValue: 'Verification code sent' }));
            setStep(1);
            setCountdown(60);
        } catch (error) {
            console.error(error);
            toast.error(error instanceof Error ? error.message : t('auth.error_generic', { defaultValue: 'Failed to send code' }));
        } finally {
            setLoading(false);
        }
    };

    const handleVerifyCode = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!code || code.length !== 6) return;

        setLoading(true);
        try {
            await http.post('/public/verification/verify', {
                target: email,
                code: code,
                type: 'email',
                scene: 'reset_password'
            });
            toast.success(t('auth.code_verified', { defaultValue: 'Verified successfully' }));
            setStep(2);
        } catch (error) {
            console.error(error);
            toast.error(error instanceof Error ? error.message : t('auth.invalid_code', { defaultValue: 'Invalid verification code' }));
        } finally {
            setLoading(false);
        }
    };

    const handleResetPassword = async (e: React.FormEvent) => {
        e.preventDefault();
        if (newPassword !== confirmPassword) {
            toast.error(t('settings.password_mismatch', { defaultValue: 'Passwords do not match' }));
            return;
        }

        setLoading(true);
        try {
            await http.post('/public/auth/reset-password', {
                email,
                code,
                newPassword,
            });

            setStep(3);
            toast.success(t('auth.password_reset_success', { defaultValue: 'Password reset successfully' }));
        } catch (error) {
            toast.error(error instanceof Error ? error.message : t('auth.error_generic'));
        } finally {
            setLoading(false);
        }
    };

    const renderStepContent = () => {
        switch (step) {
            case 0: // Input Email
                return (
                    <form onSubmit={handleSendCode} className="space-y-4 animate-in fade-in slide-in-from-right-4 duration-300">
                        <div className="space-y-2">
                            <Label htmlFor="email">{t('auth.email', { defaultValue: 'Email' })}</Label>
                            <div className="relative">
                                <Mail className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                                <Input
                                    id="email"
                                    placeholder="name@example.com"
                                    type="email"
                                    className="pl-9"
                                    value={email}
                                    onChange={(e) => setEmail(e.target.value)}
                                    required
                                    disabled={loading}
                                />
                            </div>
                        </div>
                        <Button type="submit" className="w-full" disabled={loading}>
                            {loading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                            {t('auth.send_code', { defaultValue: 'Send Verification Code' })}
                        </Button>
                    </form>
                );

            case 1: // Verify Code
                return (
                    <form onSubmit={handleVerifyCode} className="space-y-4 animate-in fade-in slide-in-from-right-4 duration-300">
                        <div className="text-center space-y-2 mb-6">
                            <div className="text-sm text-muted-foreground">
                                {t('auth.code_sent_desc', { defaultValue: 'Enter the 6-digit code sent to' })} <br />
                                <span className="font-medium text-foreground">{email}</span>
                            </div>
                        </div>
                        <div className="space-y-2">
                            <Label htmlFor="code" className="sr-only">{t('auth.code', { defaultValue: 'Verification Code' })}</Label>
                            <div className="relative flex justify-center">
                                <ShieldCheck className="absolute left-10 top-3 h-4 w-4 text-muted-foreground hidden sm:block" />
                                <Input
                                    id="code"
                                    placeholder="000000"
                                    value={code}
                                    onChange={(e) => setCode(e.target.value.replace(/[^0-9]/g, '').slice(0, 6))}
                                    className="text-center text-2xl tracking-[0.5em] font-mono h-12 w-full max-w-[240px]"
                                    required
                                    maxLength={6}
                                    disabled={loading}
                                    autoFocus
                                />
                            </div>
                        </div>
                        <Button type="submit" className="w-full" disabled={loading || code.length !== 6}>
                            {loading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                            {t('auth.verify_code', { defaultValue: 'Verify Code' })}
                        </Button>
                        <div className="text-center">
                            <Button
                                type="button"
                                variant="link"
                                size="sm"
                                disabled={countdown > 0 || loading}
                                onClick={handleSendCode}
                                className="text-muted-foreground"
                            >
                                {countdown > 0 ? `${t('auth.resend_in', { defaultValue: 'Resend in' })} ${countdown}s` : t('auth.resend_code', { defaultValue: 'Resend Code' })}
                            </Button>
                        </div>
                    </form>
                );

            case 2: // Reset Password
                return (
                    <form onSubmit={handleResetPassword} className="space-y-4 animate-in fade-in slide-in-from-right-4 duration-300">
                        <div className="space-y-2">
                            <Label htmlFor="new-password">{t('settings.new_password', { defaultValue: 'New Password' })}</Label>
                            <div className="relative">
                                <KeyRound className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                                <Input
                                    id="new-password"
                                    type="password"
                                    className="pl-9"
                                    value={newPassword}
                                    onChange={(e) => setNewPassword(e.target.value)}
                                    required
                                    disabled={loading}
                                    minLength={6}
                                />
                            </div>
                        </div>
                        <div className="space-y-2">
                            <Label htmlFor="confirm-password">{t('settings.confirm_password', { defaultValue: 'Confirm Password' })}</Label>
                            <div className="relative">
                                <KeyRound className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                                <Input
                                    id="confirm-password"
                                    type="password"
                                    className="pl-9"
                                    value={confirmPassword}
                                    onChange={(e) => setConfirmPassword(e.target.value)}
                                    required
                                    disabled={loading}
                                    minLength={6}
                                />
                            </div>
                        </div>
                        <Button type="submit" className="w-full" disabled={loading}>
                            {loading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                            {t('auth.reset_password', { defaultValue: 'Reset Password' })}
                        </Button>
                    </form>
                );

            case 3: // Success
                return (
                    <div className="flex flex-col items-center justify-center py-6 space-y-4 animate-in fade-in zoom-in duration-500">
                        <div className="h-20 w-20 bg-green-500/10 text-green-500 rounded-full flex items-center justify-center">
                            <CheckCircle2 className="h-10 w-10" />
                        </div>
                        <div className="text-center space-y-2">
                            <h3 className="font-semibold text-xl">{t('auth.password_reset_success_title', { defaultValue: 'Password Reset!' })}</h3>
                            <p className="text-sm text-muted-foreground max-w-xs mx-auto">
                                {t('auth.password_reset_success_desc', { defaultValue: 'Your password has been successfully reset. You can now login with your new password.' })}
                            </p>
                        </div>
                        <Button className="w-full mt-4" onClick={() => navigate('/login')}>
                            {t('auth.back_to_login', { defaultValue: 'Back to Login' })}
                        </Button>
                    </div>
                );
        }
    };

    return (
        <div className="flex items-center justify-center min-h-screen bg-muted/30 p-4">
            <Card className="w-full max-w-md border-white/10 shadow-xl bg-background/60 backdrop-blur-xl transition-all duration-300">
                <CardHeader className="space-y-1">
                    <div className="flex items-center gap-2 mb-2">
                        {step < 3 && (
                            <Button
                                variant="ghost"
                                size="icon"
                                className="h-8 w-8 -ml-2 rounded-full"
                                onClick={() => step > 0 ? setStep(s => s - 1) : navigate('/login')}
                            >
                                <ArrowLeft className="h-4 w-4" />
                            </Button>
                        )}
                        <CardTitle className="text-2xl font-bold">
                            {step === 0 && t('auth.forgot_password', { defaultValue: 'Forgot password?' })}
                            {step === 1 && t('auth.verify_email', { defaultValue: 'Verify Email' })}
                            {step === 2 && t('auth.reset_password', { defaultValue: 'Reset Password' })}
                            {step === 3 && t('auth.success', { defaultValue: 'Success' })}
                        </CardTitle>
                    </div>
                    {step < 3 && (
                        <CardDescription>
                            {step === 0 && t('auth.forgot_password_desc', { defaultValue: 'Enter your email address and we\'ll send you a verification code.' })}
                            {step === 1 && t('auth.enter_code_desc', { defaultValue: 'Please enter the verification code sent to your email.' })}
                            {step === 2 && t('auth.set_new_password_desc', { defaultValue: 'Create a new password for your account.' })}
                        </CardDescription>
                    )}
                </CardHeader>
                <CardContent>
                    {renderStepContent()}
                </CardContent>
                {step === 0 && (
                    <CardFooter className="flex justify-center border-t p-4">
                        <div className="text-sm text-muted-foreground">
                            {t('auth.remember_password', { defaultValue: 'Remember your password?' })}{" "}
                            <Link to="/login" className="underline underline-offset-4 hover:text-primary">
                                {t('auth.sign_in', { defaultValue: 'Sign in' })}
                            </Link>
                        </div>
                    </CardFooter>
                )}
            </Card>
        </div>
    );
}
