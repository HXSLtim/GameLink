import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { PageContainer } from '@/components/page-container';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Switch } from '@/components/ui/switch';
import { Label } from '@/components/ui/label';
import { VStack, HStack } from '@/components/layout';
import { SectionHeader, InfoRow } from '@/components/common';
import {
  Settings,
  User,
  Lock,
  Bell,
  Shield,
  Moon,
  Sun,
  ChevronRight,
  LogOut,
  Trash2,
} from 'lucide-react';
import { useAuthStore } from '@/stores/modules/auth-store';
import { useTheme } from '@/components/theme-provider';
import { toast } from 'sonner';

export default function SettingsPage() {
  const navigate = useNavigate();
  const { logout, user } = useAuthStore();
  const { theme, setTheme } = useTheme();
  
  const [notifications, setNotifications] = useState({
    order: true,
    chat: true,
    system: true,
    marketing: false,
  });

  const handleLogout = async () => {
    try {
      await logout();
      navigate('/login');
      toast.success('已退出登录');
    } catch (error) {
      console.error('退出失败:', error);
      toast.error('退出失败');
    }
  };

  const settingsSections = [
    {
      title: '账户设置',
      items: [
        {
          icon: User,
          label: '编辑个人资料',
          description: '修改头像、昵称等信息',
          onClick: () => navigate('/settings/profile'),
        },
        {
          icon: Lock,
          label: '修改密码',
          description: '更改登录密码',
          onClick: () => navigate('/settings/password'),
        },
        {
          icon: Shield,
          label: '黑名单管理',
          description: '管理已屏蔽的用户',
          onClick: () => navigate('/settings/blocked'),
        },
      ],
    },
  ];

  return (
    <PageContainer>
      <VStack spacing={6} className="py-6">
        <SectionHeader
          title="设置"
          subtitle="管理您的账户和偏好设置"
          icon={Settings}
          size="lg"
        />

        {/* 账户信息卡片 */}
        <Card>
          <CardContent className="p-4">
            <HStack spacing={4} align="center">
              <div className="w-16 h-16 rounded-full bg-primary/10 flex items-center justify-center">
                <User className="w-8 h-8 text-primary" />
              </div>
              <VStack spacing={1} className="flex-1">
                <span className="font-semibold text-lg">{user?.nickname || '用户'}</span>
                <span className="text-sm text-muted-foreground">{user?.phone || user?.email}</span>
              </VStack>
              <Button variant="outline" onClick={() => navigate('/settings/profile')}>
                编辑
              </Button>
            </HStack>
          </CardContent>
        </Card>

        {/* 设置项目 */}
        {settingsSections.map((section) => (
          <Card key={section.title}>
            <CardHeader className="pb-2">
              <CardTitle className="text-base">{section.title}</CardTitle>
            </CardHeader>
            <CardContent className="p-0">
              <VStack spacing={0} className="divide-y divide-border">
                {section.items.map((item) => (
                  <button
                    key={item.label}
                    onClick={item.onClick}
                    className="w-full flex items-center gap-4 p-4 hover:bg-muted/50 transition-colors text-left"
                  >
                    <div className="w-10 h-10 rounded-lg bg-muted flex items-center justify-center shrink-0">
                      <item.icon className="w-5 h-5 text-muted-foreground" />
                    </div>
                    <VStack spacing={0.5} className="flex-1">
                      <span className="font-medium">{item.label}</span>
                      <span className="text-sm text-muted-foreground">{item.description}</span>
                    </VStack>
                    <ChevronRight className="w-5 h-5 text-muted-foreground" />
                  </button>
                ))}
              </VStack>
            </CardContent>
          </Card>
        ))}

        {/* 通知设置 */}
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-base flex items-center gap-2">
              <Bell className="w-4 h-4" />
              通知设置
            </CardTitle>
          </CardHeader>
          <CardContent>
            <VStack spacing={4}>
              <HStack justify="between" align="center">
                <Label htmlFor="order-notify">订单通知</Label>
                <Switch
                  id="order-notify"
                  checked={notifications.order}
                  onCheckedChange={(checked) => setNotifications(prev => ({ ...prev, order: checked }))}
                />
              </HStack>
              <HStack justify="between" align="center">
                <Label htmlFor="chat-notify">消息通知</Label>
                <Switch
                  id="chat-notify"
                  checked={notifications.chat}
                  onCheckedChange={(checked) => setNotifications(prev => ({ ...prev, chat: checked }))}
                />
              </HStack>
              <HStack justify="between" align="center">
                <Label htmlFor="system-notify">系统通知</Label>
                <Switch
                  id="system-notify"
                  checked={notifications.system}
                  onCheckedChange={(checked) => setNotifications(prev => ({ ...prev, system: checked }))}
                />
              </HStack>
              <HStack justify="between" align="center">
                <Label htmlFor="marketing-notify">营销通知</Label>
                <Switch
                  id="marketing-notify"
                  checked={notifications.marketing}
                  onCheckedChange={(checked) => setNotifications(prev => ({ ...prev, marketing: checked }))}
                />
              </HStack>
            </VStack>
          </CardContent>
        </Card>

        {/* 主题设置 */}
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-base flex items-center gap-2">
              {theme === 'dark' ? <Moon className="w-4 h-4" /> : <Sun className="w-4 h-4" />}
              外观设置
            </CardTitle>
          </CardHeader>
          <CardContent>
            <HStack justify="between" align="center">
              <Label>深色模式</Label>
              <Switch
                checked={theme === 'dark'}
                onCheckedChange={(checked) => setTheme(checked ? 'dark' : 'light')}
              />
            </HStack>
          </CardContent>
        </Card>

        {/* 危险操作 */}
        <Card className="border-destructive/30">
          <CardContent className="p-4">
            <VStack spacing={4}>
              <Button
                variant="outline"
                className="w-full justify-start text-destructive hover:text-destructive hover:bg-destructive/10"
                onClick={handleLogout}
              >
                <LogOut className="w-4 h-4 mr-2" />
                退出登录
              </Button>
              <Button
                variant="outline"
                className="w-full justify-start text-destructive hover:text-destructive hover:bg-destructive/10"
                disabled
              >
                <Trash2 className="w-4 h-4 mr-2" />
                注销账号
              </Button>
            </VStack>
          </CardContent>
        </Card>

        {/* 版本信息 */}
        <div className="text-center text-sm text-muted-foreground">
          <p>GameLink v1.0.0</p>
          <p className="mt-1">© 2026 GameLink. All rights reserved.</p>
        </div>
      </VStack>
    </PageContainer>
  );
}
