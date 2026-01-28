import { useState, useEffect } from 'react';
import { PageContainer } from '@/components/page-container';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { EmptyState, SectionHeader } from '@/components/common';
import { VStack, HStack } from '@/components/layout';
import { Shield, UserX, Loader2 } from 'lucide-react';
import { blockApi, type BlockedUser } from '@/api/block';
import { toast } from 'sonner';

export default function BlockedUsersPage() {
  const [blockedUsers, setBlockedUsers] = useState<BlockedUser[]>([]);
  const [loading, setLoading] = useState(true);
  const [unblocking, setUnblocking] = useState<number | null>(null);

  useEffect(() => {
    loadBlockedUsers();
  }, []);

  const loadBlockedUsers = async () => {
    try {
      setLoading(true);
      const data = await blockApi.getBlockedUsers();
      setBlockedUsers(data);
    } catch (error) {
      console.error('加载黑名单失败:', error);
      toast.error('加载黑名单失败');
    } finally {
      setLoading(false);
    }
  };

  const handleUnblock = async (userId: number) => {
    try {
      setUnblocking(userId);
      await blockApi.unblockUser(userId);
      setBlockedUsers(prev => prev.filter(u => u.blockedUserId !== userId));
      toast.success('已解除屏蔽');
    } catch (error) {
      console.error('解除屏蔽失败:', error);
      toast.error('解除屏蔽失败');
    } finally {
      setUnblocking(null);
    }
  };

  return (
    <PageContainer>
      <VStack spacing={6} className="py-6">
        <SectionHeader
          title="黑名单管理"
          subtitle="管理您屏蔽的用户"
          icon={Shield}
          size="lg"
        />

        <Card>
          <CardHeader>
            <CardTitle className="text-lg flex items-center gap-2">
              <UserX className="w-5 h-5" />
              已屏蔽用户
            </CardTitle>
          </CardHeader>
          <CardContent>
            {loading ? (
              <div className="flex items-center justify-center py-12">
                <Loader2 className="w-8 h-8 animate-spin text-muted-foreground" />
              </div>
            ) : blockedUsers.length === 0 ? (
              <EmptyState
                icon={Shield}
                title="暂无屏蔽用户"
                description="您还没有屏蔽任何用户"
              />
            ) : (
              <VStack spacing={0} className="divide-y divide-border">
                {blockedUsers.map((user) => (
                  <HStack
                    key={user.blockedUserId}
                    justify="between"
                    align="center"
                    className="py-4"
                  >
                    <HStack spacing={3} align="center">
                      <Avatar className="w-10 h-10">
                        <AvatarImage src={user.blockedUser?.avatarUrl} />
                        <AvatarFallback>
                          {user.blockedUser?.nickname?.charAt(0) || 'U'}
                        </AvatarFallback>
                      </Avatar>
                      <VStack spacing={0.5}>
                        <span className="font-medium">
                          {user.blockedUser?.nickname || '未知用户'}
                        </span>
                        <span className="text-xs text-muted-foreground">
                          屏蔽于 {new Date(user.createdAt).toLocaleDateString()}
                        </span>
                      </VStack>
                    </HStack>
                    
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => handleUnblock(user.blockedUserId)}
                      disabled={unblocking === user.blockedUserId}
                    >
                      {unblocking === user.blockedUserId ? (
                        <Loader2 className="w-4 h-4 animate-spin mr-1" />
                      ) : null}
                      解除屏蔽
                    </Button>
                  </HStack>
                ))}
              </VStack>
            )}
          </CardContent>
        </Card>
      </VStack>
    </PageContainer>
  );
}
