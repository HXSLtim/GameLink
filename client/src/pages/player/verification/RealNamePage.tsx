
import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { ShieldCheck, Upload, Loader2 } from "lucide-react";

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
    realName: z.string().min(2, "真实姓名至少需要2个字符"),
    idNumber: z.string().regex(/(^\d{15}$)|(^\d{18}$)|(^\d{17}(\d|X|x)$)/, "请输入有效的身份证号"),
    idCardFront: z.any().optional(),
    idCardBack: z.any().optional(),
});

export default function RealNamePage() {
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
            toast.success("提交成功", {
                description: "您的实名认证信息已提交审核",
            });
        } catch (error) {
            toast.error("提交失败", {
                description: error instanceof Error ? error.message : "请稍后重试"
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
                <h1 className="text-2xl font-bold tracking-tight mb-2">实名认证</h1>
                <p className="text-sm text-muted-foreground mb-6">
                    为了保障账号安全，请完成实名认证
                </p>

                <div className="max-w-md mx-auto space-y-6">
                    <Alert variant="default" className="bg-blue-50/50 dark:bg-blue-900/20 text-blue-800 dark:text-blue-200 border-blue-200 dark:border-blue-800">
                        <ShieldCheck className="h-4 w-4" />
                        <AlertTitle>安全提示</AlertTitle>
                        <AlertDescription>
                            您的信息仅用于身份核验，平台将严格保密。
                        </AlertDescription>
                    </Alert>

                    <Form {...form}>
                        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
                            <Card>
                                <CardHeader>
                                    <CardTitle>身份信息</CardTitle>
                                    <CardDescription>
                                        请填写真实的身份信息，与身份证保持一致
                                    </CardDescription>
                                </CardHeader>
                                <CardContent className="space-y-4">
                                    <FormField
                                        control={form.control}
                                        name="realName"
                                        render={({ field }) => (
                                            <FormItem>
                                                <FormLabel>真实姓名</FormLabel>
                                                <FormControl>
                                                    <Input placeholder="请输入真实姓名" {...field} />
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
                                                <FormLabel>身份证号</FormLabel>
                                                <FormControl>
                                                    <Input placeholder="请输入18位身份证号" {...field} />
                                                </FormControl>
                                                <FormMessage />
                                            </FormItem>
                                        )}
                                    />
                                </CardContent>
                            </Card>

                            <Card>
                                <CardHeader>
                                    <CardTitle>证件照片</CardTitle>
                                    <CardDescription>
                                        请上传有效期内的身份证正反面照片
                                    </CardDescription>
                                </CardHeader>
                                <CardContent className="space-y-4">
                                    <div className="grid grid-cols-2 gap-4">
                                        {/* Front Side */}
                                        <div className="space-y-2">
                                            <FormLabel className="text-xs text-muted-foreground block text-center">人像面</FormLabel>
                                            <div
                                                className="relative aspect-[3/2] rounded-lg border-2 border-dashed border-muted-foreground/25 hover:border-primary/50 transition-colors flex flex-col items-center justify-center bg-muted/50 cursor-pointer overflow-hidden"
                                                onClick={() => document.getElementById('front-upload')?.click()}
                                            >
                                                {frontPreview ? (
                                                    <img src={frontPreview} alt="Front" className="w-full h-full object-cover" />
                                                ) : (
                                                    <>
                                                        <Upload className="h-8 w-8 text-muted-foreground mb-2" />
                                                        <span className="text-xs text-muted-foreground">点击上传</span>
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
                                            <FormLabel className="text-xs text-muted-foreground block text-center">国徽面</FormLabel>
                                            <div
                                                className="relative aspect-[3/2] rounded-lg border-2 border-dashed border-muted-foreground/25 hover:border-primary/50 transition-colors flex flex-col items-center justify-center bg-muted/50 cursor-pointer overflow-hidden"
                                                onClick={() => document.getElementById('back-upload')?.click()}
                                            >
                                                {backPreview ? (
                                                    <img src={backPreview} alt="Back" className="w-full h-full object-cover" />
                                                ) : (
                                                    <>
                                                        <Upload className="h-8 w-8 text-muted-foreground mb-2" />
                                                        <span className="text-xs text-muted-foreground">点击上传</span>
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
                                        提交审核中...
                                    </>
                                ) : (
                                    "提交认证"
                                )}
                            </Button>
                        </form>
                    </Form>
                </div>
            </div>
        </PageContainer>
    );
}
