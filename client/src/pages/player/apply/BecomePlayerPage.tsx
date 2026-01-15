import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { PageContainer } from '@/components/page-container';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { CheckCircle2, Upload, ChevronRight, ChevronLeft, Gamepad2, ShieldCheck, UserCircle2 } from 'lucide-react';
import { toast } from 'sonner';
import { useNavigate } from 'react-router-dom';

// Define game ranks mapping
const GAME_RANKS: Record<string, { label: string; value: string }[]> = {
    lol: [
        { label: "Challenger", value: "challenger" },
        { label: "Grandmaster", value: "grandmaster" },
        { label: "Master", value: "master" },
        { label: "Diamond", value: "diamond" },
        { label: "Platinum", value: "platinum" },
    ],
    apex: [
        { label: "Apex Predator", value: "predator" },
        { label: "Master", value: "master" },
        { label: "Diamond", value: "diamond" },
        { label: "Platinum", value: "platinum" },
    ],
    valorant: [
        { label: "Radiant", value: "radiant" },
        { label: "Immortal", value: "immortal" },
        { label: "Ascendant", value: "ascendant" },
        { label: "Diamond", value: "diamond" },
    ],
    genshin: [
        { label: "Adventure Rank 60", value: "ar60" },
        { label: "Adventure Rank 55-59", value: "ar55" },
        { label: "Spiral Abyss 12-3", value: "abyss12" },
    ],
};

export default function BecomePlayerPage() {
    const { t } = useTranslation();
    const navigate = useNavigate();
    const [step, setStep] = useState(1);
    const [loading, setLoading] = useState(false);

    // Form State
    const [formData, setFormData] = useState({
        gameId: '',
        rank: '',
        rankImage: null as File | null,
        description: '',
        price: '',
    });

    const handleNext = () => {
        if (step === 1 && (!formData.gameId || !formData.rank)) {
            toast.error(t('player.apply.error_step1', { defaultValue: 'Please select a game and rank.' }));
            return;
        }
        if (step === 2 && !formData.rankImage) {
            toast.error(t('player.apply.error_step2', { defaultValue: 'Please upload rank proof.' }));
            return;
        }
        setStep(step + 1);
    };

    const handleBack = () => {
        setStep(step - 1);
    };

    const handleSubmit = async () => {
        if (!formData.description || !formData.price) {
            toast.error(t('player.apply.error_step3', { defaultValue: 'Please fill in all details.' }));
            return;
        }

        setLoading(true);
        // Simulate API call
        setTimeout(() => {
            setLoading(false);
            toast.success(t('player.apply.success', { defaultValue: 'Application submitted! We will review it shortly.' }));
            navigate('/player/dashboard');
        }, 1500);
    };

    const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (e.target.files && e.target.files[0]) {
            setFormData({ ...formData, rankImage: e.target.files[0] });
        }
    };

    const steps = [
        { id: 1, title: t('player.apply.step1_title', { defaultValue: 'Select Game' }), icon: Gamepad2 },
        { id: 2, title: t('player.apply.step2_title', { defaultValue: 'Rank Proof' }), icon: ShieldCheck },
        { id: 3, title: t('player.apply.step3_title', { defaultValue: 'Profile Info' }), icon: UserCircle2 },
    ];

    return (
        <PageContainer>
            <div className="max-w-3xl mx-auto py-12 px-4">
                <div className="text-center mb-10">
                    <h1 className="text-3xl font-bold tracking-tight mb-2">{t('player.apply.title', { defaultValue: 'Become a Pro Player' })}</h1>
                    <p className="text-muted-foreground">{t('player.apply.subtitle', { defaultValue: 'Share your skills and earn money playing games.' })}</p>
                </div>

                {/* Stepper */}
                <div className="flex justify-between items-center mb-12 relative">
                    <div className="absolute top-1/2 left-0 w-full h-0.5 bg-muted -z-10 transform -translate-y-1/2" />
                    {steps.map((s) => (
                        <div key={s.id} className="flex flex-col items-center bg-background px-4 z-10">
                            <div className={`w-10 h-10 rounded-full flex items-center justify-center border-2 transition-colors ${step >= s.id ? 'border-primary bg-primary text-primary-foreground' : 'border-muted text-muted-foreground'}`}>
                                <s.icon className="h-5 w-5" />
                            </div>
                            <span className={`text-sm mt-2 font-medium ${step >= s.id ? 'text-primary' : 'text-muted-foreground'}`}>{s.title}</span>
                        </div>
                    ))}
                </div>

                {/* Form Content */}
                <Card className="border-2">
                    <CardHeader>
                        <CardTitle>{steps[step - 1].title}</CardTitle>
                        <CardDescription>{t('player.apply.step_desc', { defaultValue: 'Fill in the information below.' })}</CardDescription>
                    </CardHeader>
                    <CardContent className="space-y-6">


                        {step === 1 && (
                            <div className="space-y-4 animate-in fade-in slide-in-from-right-4">
                                <div className="space-y-2">
                                    <Label>{t('player.apply.game', { defaultValue: 'Game' })}</Label>
                                    <Select onValueChange={(v: string) => setFormData({ ...formData, gameId: v, rank: '' })} value={formData.gameId}>
                                        <SelectTrigger>
                                            <SelectValue placeholder={t('player.apply.select_game', { defaultValue: 'Select a game' })} />
                                        </SelectTrigger>
                                        <SelectContent>
                                            <SelectItem value="lol">League of Legends</SelectItem>
                                            <SelectItem value="apex">Apex Legends</SelectItem>
                                            <SelectItem value="valorant">Valorant</SelectItem>
                                            <SelectItem value="genshin">Genshin Impact</SelectItem>
                                        </SelectContent>
                                    </Select>
                                </div>
                                <div className="space-y-2">
                                    <Label>{t('player.apply.rank', { defaultValue: 'Current Rank' })}</Label>
                                    <Select
                                        onValueChange={(v: string) => setFormData({ ...formData, rank: v })}
                                        value={formData.rank}
                                        disabled={!formData.gameId}
                                    >
                                        <SelectTrigger>
                                            <SelectValue placeholder={t('player.apply.select_rank', { defaultValue: 'Select your rank' })} />
                                        </SelectTrigger>
                                        <SelectContent>
                                            {(formData.gameId && GAME_RANKS[formData.gameId] ? GAME_RANKS[formData.gameId] : []).map((rank) => (
                                                <SelectItem key={rank.value} value={rank.value}>
                                                    {rank.label}
                                                </SelectItem>
                                            ))}
                                        </SelectContent>
                                    </Select>
                                </div>
                            </div>
                        )}

                        {step === 2 && (
                            <div className="space-y-4 animate-in fade-in slide-in-from-right-4">
                                <Label>{t('player.apply.upload_proof', { defaultValue: 'Upload Screenshot of Rank' })}</Label>
                                <div className="border-2 border-dashed rounded-xl p-10 flex flex-col items-center justify-center text-center cursor-pointer hover:bg-muted/50 transition-colors relative">
                                    <input
                                        type="file"
                                        accept="image/*"
                                        className="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
                                        onChange={handleFileChange}
                                    />
                                    {formData.rankImage ? (
                                        <div className="flex flex-col items-center text-primary">
                                            <CheckCircle2 className="h-10 w-10 mb-2" />
                                            <p className="font-medium">{formData.rankImage.name}</p>
                                        </div>
                                    ) : (
                                        <>
                                            <Upload className="h-10 w-10 text-muted-foreground mb-4" />
                                            <p className="font-medium">{t('player.apply.drop_file', { defaultValue: 'Click to upload or drag and drop' })}</p>
                                            <p className="text-sm text-muted-foreground mt-1">PNG, JPG up to 5MB</p>
                                        </>
                                    )}
                                </div>
                            </div>
                        )}

                        {step === 3 && (
                            <div className="space-y-4 animate-in fade-in slide-in-from-right-4">
                                <div className="space-y-2">
                                    <Label>{t('player.apply.price', { defaultValue: 'Hourly Rate (CNY)' })}</Label>
                                    <Input
                                        type="number"
                                        placeholder="Min 10"
                                        value={formData.price}
                                        onChange={(e) => setFormData({ ...formData, price: e.target.value })}
                                    />
                                </div>
                                <div className="space-y-2">
                                    <Label>{t('player.apply.bio', { defaultValue: 'Self Introduction' })}</Label>
                                    <Textarea
                                        placeholder={t('player.apply.bio_placeholder', { defaultValue: 'Describe your playstyle, availability, and why people should hire you...' })}
                                        className="min-h-[120px]"
                                        value={formData.description}
                                        onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                                    />
                                </div>
                            </div>
                        )}

                        <div className="flex justify-between pt-4">
                            <Button variant="ghost" onClick={handleBack} disabled={step === 1}>
                                <ChevronLeft className="mr-2 h-4 w-4" /> {t('common.back', { defaultValue: 'Back' })}
                            </Button>

                            {step < 3 ? (
                                <Button onClick={handleNext}>
                                    {t('common.next', { defaultValue: 'Next' })} <ChevronRight className="ml-2 h-4 w-4" />
                                </Button>
                            ) : (
                                <Button onClick={handleSubmit} disabled={loading}>
                                    {loading ? t('common.submitting', { defaultValue: 'Submitting...' }) : t('common.submit', { defaultValue: 'Submit Application' })}
                                </Button>
                            )}
                        </div>
                    </CardContent>
                </Card>
            </div>
        </PageContainer>
    );
}
