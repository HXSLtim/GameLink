import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Trophy, Loader2 } from "lucide-react";
import { useTranslation } from "react-i18next";

import { PageContainer } from "@/components/page-container";
import { Button } from "@/components/ui/button";
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from "@/components/ui/form";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { toast } from "sonner";
import { http } from "@/lib/http";

const formSchema = z.object({
    gameId: z.string().min(1, "verification.skill.select_game"),
    rankId: z.string().min(1, "verification.skill.select_rank"),
    screenshot: z.any().optional(),
});

// Define game ranks mapping with i18n keys
const GAME_RANKS: Record<string, { label: string; value: string }[]> = {
    lol: [
        { label: "verification.skill.lol.challenger", value: "challenger" },
        { label: "verification.skill.lol.grandmaster", value: "grandmaster" },
        { label: "verification.skill.lol.master", value: "master" },
        { label: "verification.skill.lol.diamond", value: "diamond" },
        { label: "verification.skill.lol.platinum", value: "platinum" },
    ],
    wzry: [
        { label: "verification.skill.wzry.glory_king", value: "glory_king" },
        { label: "verification.skill.wzry.king", value: "king" },
        { label: "verification.skill.wzry.star", value: "star" },
        { label: "verification.skill.wzry.diamond", value: "diamond" },
        { label: "verification.skill.wzry.platinum", value: "platinum" },
    ],
    pubg: [
        { label: "verification.skill.pubg.conqueror", value: "conqueror" },
        { label: "verification.skill.pubg.ace", value: "ace" },
        { label: "verification.skill.pubg.crown", value: "crown" },
        { label: "verification.skill.pubg.diamond", value: "diamond" },
        { label: "verification.skill.pubg.platinum", value: "platinum" },
    ],
    apex: [
        { label: "verification.skill.apex.predator", value: "predator" },
        { label: "verification.skill.apex.master", value: "master" },
        { label: "verification.skill.apex.diamond", value: "diamond" },
        { label: "verification.skill.apex.platinum", value: "platinum" },
        { label: "verification.skill.apex.gold", value: "gold" },
    ],
};

export default function SkillAuthPage() {
    const { t } = useTranslation();
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [preview, setPreview] = useState<string | null>(null);

    const form = useForm<z.infer<typeof formSchema>>({
        resolver: zodResolver(formSchema),
    });

    const selectedGame = form.watch("gameId");
    const currentRanks = selectedGame ? GAME_RANKS[selectedGame as keyof typeof GAME_RANKS] || [] : [];

    async function onSubmit(values: z.infer<typeof formSchema>) {
        setIsSubmitting(true);
        try {
            await http.post('/player/verification/skill', {
                gameId: values.gameId,
                rank: values.rankId,
                // In production, upload screenshot first and send URL
                // screenshotUrl: screenshotImageUrl,
            });
            toast.success(t('verification.skill.submit_success'), {
                description: t('verification.skill.submit_success_desc'),
            });
        } catch (error) {
            toast.error(t('verification.skill.submit_failed'), {
                description: error instanceof Error ? error.message : t('verification.skill.try_again')
            });
        } finally {
            setIsSubmitting(false);
        }
    }

    const handleImageChange = (
        e: React.ChangeEvent<HTMLInputElement>
    ) => {
        const file = e.target.files?.[0];
        if (file) {
            const url = URL.createObjectURL(file);
            setPreview(url);
        }
    };

    return (
        <PageContainer>
            <div className="px-4 py-4 md:px-8">
                <h1 className="text-2xl font-bold tracking-tight mb-2">{t('verification.skill.title')}</h1>
                <p className="text-sm text-muted-foreground mb-6">
                    {t('verification.skill.description')}
                </p>

                <div className="max-w-md mx-auto space-y-6">
                    <Alert variant="default" className="bg-blue-50/50 dark:bg-blue-900/20 text-blue-800 dark:text-blue-200 border-blue-200 dark:border-blue-800">
                        <Trophy className="h-4 w-4" />
                        <AlertTitle>{t('verification.skill.alert_title')}</AlertTitle>
                        <AlertDescription>
                            {t('verification.skill.alert_desc')}
                        </AlertDescription>
                    </Alert>

                    <Form {...form}>
                        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
                            <Card>
                                <CardHeader>
                                    <CardTitle>{t('verification.skill.game_info')}</CardTitle>
                                    <CardDescription>
                                        {t('verification.skill.game_info_desc')}
                                    </CardDescription>
                                </CardHeader>
                                <CardContent className="space-y-4">
                                    <FormField
                                        control={form.control}
                                        name="gameId"
                                        render={({ field }) => (
                                            <FormItem>
                                                <FormLabel>{t('verification.skill.game')}</FormLabel>
                                                <Select onValueChange={field.onChange} defaultValue={field.value}>
                                                    <FormControl>
                                                        <SelectTrigger>
                                                            <SelectValue placeholder={t('verification.skill.select_game')} />
                                                        </SelectTrigger>
                                                    </FormControl>
                                                    <SelectContent>
                                                        <SelectItem value="lol">{t('verification.skill.lol.name')}</SelectItem>
                                                        <SelectItem value="wzry">{t('verification.skill.wzry.name')}</SelectItem>
                                                        <SelectItem value="pubg">{t('verification.skill.pubg.name')}</SelectItem>
                                                        <SelectItem value="apex">{t('verification.skill.apex.name')}</SelectItem>
                                                    </SelectContent>
                                                </Select>
                                                <FormMessage />
                                            </FormItem>
                                        )}
                                    />
                                    <FormField
                                        control={form.control}
                                        name="rankId"
                                        render={({ field }) => (
                                            <FormItem>
                                                <FormLabel>{t('verification.skill.rank')}</FormLabel>
                                                <Select onValueChange={field.onChange} defaultValue={field.value} disabled={!selectedGame}>
                                                    <FormControl>
                                                        <SelectTrigger>
                                                            <SelectValue placeholder={selectedGame ? t('verification.skill.select_rank') : t('verification.skill.select_game_first')} />
                                                        </SelectTrigger>
                                                    </FormControl>
                                                    <SelectContent>
                                                        {currentRanks.map((rank) => (
                                                            <SelectItem key={rank.value} value={rank.value}>
                                                                {t(rank.label)}
                                                            </SelectItem>
                                                        ))}
                                                    </SelectContent>
                                                </Select>
                                                <FormMessage />
                                            </FormItem>
                                        )}
                                    />
                                </CardContent>
                            </Card>

                            <Card>
                                <CardHeader>
                                    <CardTitle>{t('verification.skill.screenshot_title')}</CardTitle>
                                    <CardDescription>
                                        {t('verification.skill.screenshot_desc')}
                                    </CardDescription>
                                </CardHeader>
                                <CardContent className="space-y-4">
                                    <div className="space-y-2">
                                        <div
                                            className="relative aspect-video rounded-lg border-2 border-dashed border-muted-foreground/25 hover:border-primary/50 transition-colors flex flex-col items-center justify-center bg-muted/50 cursor-pointer overflow-hidden"
                                            onClick={() => document.getElementById('screenshot-upload')?.click()}
                                        >
                                            {preview ? (
                                                <img src={preview} alt="Screenshot" className="w-full h-full object-cover" />
                                            ) : (
                                                <>
                                                    <Trophy className="h-10 w-10 text-muted-foreground mb-2" />
                                                    <span className="text-sm text-muted-foreground">{t('verification.skill.click_upload')}</span>
                                                </>
                                            )}
                                            <input
                                                id="screenshot-upload"
                                                type="file"
                                                accept="image/*"
                                                className="hidden"
                                                onChange={handleImageChange}
                                            />
                                        </div>
                                    </div>
                                </CardContent>
                            </Card>

                            <Button type="submit" className="w-full" size="lg" disabled={isSubmitting}>
                                {isSubmitting ? (
                                    <>
                                        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                        {t('verification.skill.submitting')}
                                    </>
                                ) : (
                                    t('verification.skill.submit')
                                )}
                            </Button>
                        </form>
                    </Form>
                </div>
            </div>
        </PageContainer>
    );
}
