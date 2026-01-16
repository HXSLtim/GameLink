
import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { ShieldCheck, Upload, Loader2 } from "lucide-react";
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
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { toast } from "sonner";
import { http } from "@/lib/http";

const formSchema = z.object({
    realName: z.string().min(2, "verification.realname.min_length"),
    idNumber: z.string().regex(/(^\d{15}$)|(^\d{18}$)|(^\d{17}(\d|X|x)$)/, "verification.realname.invalid_id"),
    idCardFront: z.any().optional(),
    idCardBack: z.any().optional(),
});

export default function RealNamePage() {
    const { t } = useTranslation();
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [frontPreview, setFrontPreview] = useState<string | null>(null);
    const [backPreview, setBackPreview] = useState<string | null>(null);

    const form = useForm<z.infer<typeof formSchema>>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            realName: "",
            idNumber: "",
        },
    });

    async function onSubmit(values: z.infer<typeof formSchema>) {
        setIsSubmitting(true);
        try {
            await http.post('/player/verification/realname', {
                realName: values.realName,
                idCard: values.idNumber,
                // In production, upload images first and send URLs
                // idCardFrontUrl: frontImageUrl,
                // idCardBackUrl: backImageUrl,
            });
            toast.success(t('verification.realname.submit_success'), {
                description: t('verification.realname.submit_success_desc'),
            });
        } catch (error) {
            toast.error(t('verification.realname.submit_failed'), {
                description: error instanceof Error ? error.message : t('verification.realname.try_again')
            });
        } finally {
            setIsSubmitting(false);
        }
    }

    const handleImageChange = (
        e: React.ChangeEvent<HTMLInputElement>,
        setPreview: (url: string | null) => void
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
                <h1 className="text-2xl font-bold tracking-tight mb-2">{t('verification.realname.title')}</h1>
                <p className="text-sm text-muted-foreground mb-6">
                    {t('verification.realname.description')}
                </p>

                <div className="max-w-md mx-auto space-y-6">
                    <Alert variant="default" className="bg-blue-50/50 dark:bg-blue-900/20 text-blue-800 dark:text-blue-200 border-blue-200 dark:border-blue-800">
                        <ShieldCheck className="h-4 w-4" />
                        <AlertTitle>{t('verification.realname.security_title')}</AlertTitle>
                        <AlertDescription>
                            {t('verification.realname.security_desc')}
                        </AlertDescription>
                    </Alert>

                    <Form {...form}>
                        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
                            <Card>
                                <CardHeader>
                                    <CardTitle>{t('verification.realname.identity_info')}</CardTitle>
                                    <CardDescription>
                                        {t('verification.realname.identity_info_desc')}
                                    </CardDescription>
                                </CardHeader>
                                <CardContent className="space-y-4">
                                    <FormField
                                        control={form.control}
                                        name="realName"
                                        render={({ field }) => (
                                            <FormItem>
                                                <FormLabel>{t('verification.realname.real_name')}</FormLabel>
                                                <FormControl>
                                                    <Input placeholder={t('verification.realname.real_name_placeholder')} {...field} />
                                                </FormControl>
                                                <FormMessage />
                                            </FormItem>
                                        )}
                                    />
                                    <FormField
                                        control={form.control}
                                        name="idNumber"
                                        render={({ field }) => (
                                            <FormItem>
                                                <FormLabel>{t('verification.realname.id_number')}</FormLabel>
                                                <FormControl>
                                                    <Input placeholder={t('verification.realname.id_number_placeholder')} {...field} />
                                                </FormControl>
                                                <FormMessage />
                                            </FormItem>
                                        )}
                                    />
                                </CardContent>
                            </Card>

                            <Card>
                                <CardHeader>
                                    <CardTitle>{t('verification.realname.id_photos')}</CardTitle>
                                    <CardDescription>
                                        {t('verification.realname.id_photos_desc')}
                                    </CardDescription>
                                </CardHeader>
                                <CardContent className="space-y-4">
                                    <div className="grid grid-cols-2 gap-4">
                                        {/* Front Side */}
                                        <div className="space-y-2">
                                            <FormLabel className="text-xs text-muted-foreground block text-center">{t('verification.realname.photo_front')}</FormLabel>
                                            <div
                                                className="relative aspect-[3/2] rounded-lg border-2 border-dashed border-muted-foreground/25 hover:border-primary/50 transition-colors flex flex-col items-center justify-center bg-muted/50 cursor-pointer overflow-hidden"
                                                onClick={() => document.getElementById('front-upload')?.click()}
                                            >
                                                {frontPreview ? (
                                                    <img src={frontPreview} alt="Front" className="w-full h-full object-cover" />
                                                ) : (
                                                    <>
                                                        <Upload className="h-8 w-8 text-muted-foreground mb-2" />
                                                        <span className="text-xs text-muted-foreground">{t('verification.realname.click_upload')}</span>
                                                    </>
                                                )}
                                                <input
                                                    id="front-upload"
                                                    type="file"
                                                    accept="image/*"
                                                    className="hidden"
                                                    onChange={(e) => handleImageChange(e, setFrontPreview)}
                                                />
                                            </div>
                                        </div>

                                        {/* Back Side */}
                                        <div className="space-y-2">
                                            <FormLabel className="text-xs text-muted-foreground block text-center">{t('verification.realname.photo_back')}</FormLabel>
                                            <div
                                                className="relative aspect-[3/2] rounded-lg border-2 border-dashed border-muted-foreground/25 hover:border-primary/50 transition-colors flex flex-col items-center justify-center bg-muted/50 cursor-pointer overflow-hidden"
                                                onClick={() => document.getElementById('back-upload')?.click()}
                                            >
                                                {backPreview ? (
                                                    <img src={backPreview} alt="Back" className="w-full h-full object-cover" />
                                                ) : (
                                                    <>
                                                        <Upload className="h-8 w-8 text-muted-foreground mb-2" />
                                                        <span className="text-xs text-muted-foreground">{t('verification.realname.click_upload')}</span>
                                                    </>
                                                )}
                                                <input
                                                    id="back-upload"
                                                    type="file"
                                                    accept="image/*"
                                                    className="hidden"
                                                    onChange={(e) => handleImageChange(e, setBackPreview)}
                                                />
                                            </div>
                                        </div>
                                    </div>
                                </CardContent>
                            </Card>

                            <Button type="submit" className="w-full" size="lg" disabled={isSubmitting}>
                                {isSubmitting ? (
                                    <>
                                        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                        {t('verification.realname.submitting')}
                                    </>
                                ) : (
                                    t('verification.realname.submit')
                                )}
                            </Button>
                        </form>
                    </Form>
                </div>
            </div>
        </PageContainer>
    );
}
