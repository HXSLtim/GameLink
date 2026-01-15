import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Trophy, Loader2 } from "lucide-react";

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

const formSchema = z.object({
    gameId: z.string().min(1, "请选择游戏"),
    rankId: z.string().min(1, "请选择段位"),
    screenshot: z.any().optional(),
});

// Define game ranks mapping
const GAME_RANKS: Record<string, { label: string; value: string }[]> = {
    lol: [
        { label: "最强王者", value: "challenger" },
        { label: "傲世宗师", value: "grandmaster" },
        { label: "超凡大师", value: "master" },
        { label: "璀璨钻石", value: "diamond" },
        { label: "华贵铂金", value: "platinum" },
    ],
    wzry: [
        { label: "荣耀王者", value: "glory_king" },
        { label: "最强王者", value: "king" },
        { label: "至尊星耀", value: "star" },
        { label: "永恒钻石", value: "diamond" },
        { label: "尊贵铂金", value: "platinum" },
    ],
    pubg: [
        { label: "无敌战神", value: "conqueror" },
        { label: "王牌", value: "ace" },
        { label: "皇冠", value: "crown" },
        { label: "星钻", value: "diamond" },
        { label: "白金", value: "platinum" },
    ],
    apex: [
        { label: "顶尖猎杀者", value: "predator" },
        { label: "大师", value: "master" },
        { label: "钻石", value: "diamond" },
        { label: "白金", value: "platinum" },
        { label: "黄金", value: "gold" },
    ],
};

export default function SkillAuthPage() {
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
            // Simulate API call
            await new Promise((resolve) => setTimeout(resolve, 2000));
            console.log(values);
            toast.success("提交成功", {
                description: "您的技能认证信息已提交审核",
            });
        } catch (error) {
            toast.error("提交失败", {
                description: "请稍后重试"
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
                <h1 className="text-2xl font-bold tracking-tight mb-2">技能认证</h1>
                <p className="text-sm text-muted-foreground mb-6">
                    认证游戏段位，展示您的实力
                </p>

                <div className="max-w-md mx-auto space-y-6">
                    <Alert variant="default" className="bg-blue-50/50 dark:bg-blue-900/20 text-blue-800 dark:text-blue-200 border-blue-200 dark:border-blue-800">
                        <Trophy className="h-4 w-4" />
                        <AlertTitle>认证提示</AlertTitle>
                        <AlertDescription>
                            请上传带有游戏ID的个人主页截图，确保ID清晰可见。
                        </AlertDescription>
                    </Alert>

                    <Form {...form}>
                        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
                            <Card>
                                <CardHeader>
                                    <CardTitle>游戏信息</CardTitle>
                                    <CardDescription>
                                        选择您擅长的游戏及当前段位
                                    </CardDescription>
                                </CardHeader>
                                <CardContent className="space-y-4">
                                    <FormField
                                        control={form.control}
                                        name="gameId"
                                        render={({ field }) => (
                                            <FormItem>
                                                <FormLabel>游戏</FormLabel>
                                                <Select onValueChange={field.onChange} defaultValue={field.value}>
                                                    <FormControl>
                                                        <SelectTrigger>
                                                            <SelectValue placeholder="选择游戏" />
                                                        </SelectTrigger>
                                                    </FormControl>
                                                    <SelectContent>
                                                        <SelectItem value="lol">英雄联盟</SelectItem>
                                                        <SelectItem value="wzry">王者荣耀</SelectItem>
                                                        <SelectItem value="pubg">绝地求生</SelectItem>
                                                        <SelectItem value="apex">Apex Legends</SelectItem>
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
                                                <FormLabel>段位</FormLabel>
                                                <Select onValueChange={field.onChange} defaultValue={field.value} disabled={!selectedGame}>
                                                    <FormControl>
                                                        <SelectTrigger>
                                                            <SelectValue placeholder={selectedGame ? "选择段位" : "请先选择游戏"} />
                                                        </SelectTrigger>
                                                    </FormControl>
                                                    <SelectContent>
                                                        {currentRanks.map((rank) => (
                                                            <SelectItem key={rank.value} value={rank.value}>
                                                                {rank.label}
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
                                    <CardTitle>段位截图</CardTitle>
                                    <CardDescription>
                                        请上传包含游戏ID和段位信息的最新截图
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
                                                    <span className="text-sm text-muted-foreground">点击上传段位截图</span>
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
