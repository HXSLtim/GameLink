
import { useState } from "react";
import { Users, UserPlus, Settings, Crown, Trophy } from "lucide-react";

import { PageContainer } from "@/components/page-container";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { toast } from "sonner";

// 模拟数据
const MOCK_TEAM = {
    id: "1",
    name: "王者无敌战队",
    description: "全员王者，接单必胜！",
    level: 5,
    winRate: "88%",
    members: [
        { id: "1", name: "HADES", avatar: "/avatars/01.png", role: "Leader", rank: "最强王者" },
        { id: "2", name: "Zeus", avatar: "/avatars/02.png", role: "Member", rank: "傲世宗师" },
        { id: "3", name: "Poseidon", avatar: "/avatars/03.png", role: "Member", rank: "超凡大师" },
    ]
};

export default function TeamPage() {
    const [hasTeam, setHasTeam] = useState(false); // 演示切换开关
    const [isCreating, setIsCreating] = useState(false);

    const handleCreateTeam = () => {
        setIsCreating(true);
        // 模拟 API
        setTimeout(() => {
            setIsCreating(false);
            setHasTeam(true);
            toast.success("创建车队成功");
        }, 1500);
    };

    return (
        <PageContainer>
            <div className="px-4 py-4 md:px-8">
                <div className="flex items-center justify-between mb-6">
                    <div className="space-y-1">
                        <h1 className="text-2xl font-bold tracking-tight">我的车队</h1>
                        <p className="text-sm text-muted-foreground">
                            管理您的车队成员和活动
                        </p>
                    </div>
                    {/* 演示切换按钮 */}
                    <Button variant="outline" size="sm" onClick={() => setHasTeam(!hasTeam)}>
                        {hasTeam ? "切换为无车队" : "切换为有车队"}
                    </Button>
                </div>

                {hasTeam ? (
                    <div className="space-y-6">
                        {/* 车队头部信息 */}
                        <Card>
                            <CardContent className="pt-6">
                                <div className="flex flex-col md:flex-row gap-6 items-center md:items-start">
                                    <Avatar className="h-24 w-24">
                                        <AvatarImage src="/team-logo.png" />
                                        <AvatarFallback className="text-2xl">TM</AvatarFallback>
                                    </Avatar>
                                    <div className="flex-1 text-center md:text-left space-y-2">
                                        <div className="flex items-center justify-center md:justify-start gap-2">
                                            <h2 className="text-2xl font-bold">{MOCK_TEAM.name}</h2>
                                            <Badge variant="secondary">Lv.{MOCK_TEAM.level}</Badge>
                                        </div>
                                        <p className="text-muted-foreground">{MOCK_TEAM.description}</p>
                                        <div className="flex items-center justify-center md:justify-start gap-4 pt-2">
                                            <div className="flex items-center gap-1 text-sm text-muted-foreground">
                                                <Users className="h-4 w-4" />
                                                <span>{MOCK_TEAM.members.length} 成员</span>
                                            </div>
                                            <div className="flex items-center gap-1 text-sm text-muted-foreground">
                                                <Trophy className="h-4 w-4" />
                                                <span>胜率 {MOCK_TEAM.winRate}</span>
                                            </div>
                                        </div>
                                    </div>
                                    <div className="flex gap-2">
                                        <Button variant="outline" size="icon" onClick={() => toast.info("战队设置功能开发中")}>
                                            <Settings className="h-4 w-4" />
                                        </Button>
                                        <Button onClick={() => toast.info("邀请功能开发中")}>
                                            <UserPlus className="h-4 w-4 mr-2" />
                                            邀请
                                        </Button>
                                    </div>
                                </div>
                            </CardContent>
                        </Card>

                        {/* 车队成员列表 */}
                        <Card>
                            <CardHeader>
                                <CardTitle>成员列表</CardTitle>
                            </CardHeader>
                            <CardContent>
                                <div className="space-y-4">
                                    {MOCK_TEAM.members.map((member) => (
                                        <div key={member.id} className="flex items-center justify-between p-2 hover:bg-muted/50 rounded-lg transition-colors">
                                            <div className="flex items-center gap-3">
                                                <Avatar>
                                                    <AvatarImage src={member.avatar} />
                                                    <AvatarFallback>{member.name[0]}</AvatarFallback>
                                                </Avatar>
                                                <div>
                                                    <div className="font-medium flex items-center gap-2">
                                                        {member.name}
                                                        {member.role === 'Leader' && (
                                                            <Crown className="h-3 w-3 text-yellow-500" />
                                                        )}
                                                    </div>
                                                    <div className="text-xs text-muted-foreground">{member.rank}</div>
                                                </div>
                                            </div>
                                            <div className="flex items-center gap-2">
                                                <Badge variant={member.role === 'Leader' ? 'default' : 'secondary'}>
                                                    {member.role === 'Leader' ? '队长' : '队员'}
                                                </Badge>
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            </CardContent>
                        </Card>
                    </div>
                ) : (
                    <div className="flex flex-col items-center justify-center py-12 space-y-4 text-center">
                        <div className="p-6 bg-muted rounded-full">
                            <Users className="h-12 w-12 text-muted-foreground" />
                        </div>
                        <div className="space-y-2">
                            <h3 className="text-xl font-semibold">暂无车队</h3>
                            <p className="text-muted-foreground max-w-sm">
                                创建属于你的王者车队，或者加入其他大神的队伍，一起征战峡谷！
                            </p>
                        </div>
                        <div className="flex gap-4 pt-4">
                            <Dialog>
                                <DialogTrigger asChild>
                                    <Button size="lg">创建车队</Button>
                                </DialogTrigger>
                                <DialogContent>
                                    <DialogHeader>
                                        <DialogTitle>创建新车队</DialogTitle>
                                        <DialogDescription>
                                            填写车队基本信息，开启你的战队之旅
                                        </DialogDescription>
                                    </DialogHeader>
                                    <div className="space-y-4 py-4">
                                        <div className="space-y-2">
                                            <Label>车队名称</Label>
                                            <Input placeholder="给车队起个响亮的名字" />
                                        </div>
                                        <div className="space-y-2">
                                            <Label>车队宣言</Label>
                                            <Input placeholder="一句话介绍你的车队" />
                                        </div>
                                    </div>
                                    <DialogFooter>
                                        <Button onClick={handleCreateTeam} disabled={isCreating}>
                                            {isCreating ? "创建中..." : "确认创建"}
                                        </Button>
                                    </DialogFooter>
                                </DialogContent>
                            </Dialog>
                            <Button variant="outline" size="lg">加入车队</Button>
                        </div>
                    </div>
                )}
            </div>
        </PageContainer>
    );
}
