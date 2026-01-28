import { useState, useEffect } from 'react';
import { PageContainer } from '@/components/page-container';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { EmptyState, SectionHeader, StatusBadge } from '@/components/common';
import { VStack, HStack } from '@/components/layout';
import { AlertTriangle, Clock, CheckCircle, XCircle, Loader2, FileText } from 'lucide-react';
import { disputeApi, type Dispute } from '@/api/dispute';
import { toast } from 'sonner';
import { format } from 'date-fns';
import { zhCN } from 'date-fns/locale';

const statusConfig: Record<string, { label: string; variant: 'pending' | 'processing' | 'success' | 'error' | 'default' }> = {
  pending: { label: '待处理', variant: 'pending' },
  processing: { label: '处理中', variant: 'processing' },
  resolved: { label: '已解决', variant: 'success' },
  rejected: { label: '已驳回', variant: 'error' },
  cancelled: { label: '已取消', variant: 'default' },
};

export default function DisputePage() {
  const [disputes, setDisputes] = useState<Dispute[]>([]);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState('all');

  useEffect(() => {
    loadDisputes();
  }, []);

  const loadDisputes = async () => {
    try {
      setLoading(true);
      const data = await disputeApi.getMyDisputes();
      setDisputes(data);
    } catch (error) {
      console.error('加载申诉失败:', error);
      toast.error('加载申诉失败');
    } finally {
      setLoading(false);
    }
  };

  const filteredDisputes = disputes.filter(dispute => {
    if (activeTab === 'all') return true;
    return dispute.status === activeTab;
  });

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'pending': return Clock;
      case 'processing': return Loader2;
      case 'resolved': return CheckCircle;
      case 'rejected': return XCircle;
      default: return FileText;
    }
  };

  return (
    <PageContainer>
      <VStack spacing={6} className="py-6">
        <SectionHeader
          title="我的申诉"
          subtitle="查看和管理您的订单申诉记录"
          icon={AlertTriangle}
          size="lg"
        />

        <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
          <TabsList className="grid w-full grid-cols-4">
            <TabsTrigger value="all">全部</TabsTrigger>
            <TabsTrigger value="pending">待处理</TabsTrigger>
            <TabsTrigger value="processing">处理中</TabsTrigger>
            <TabsTrigger value="resolved">已完成</TabsTrigger>
          </TabsList>

          <TabsContent value={activeTab} className="mt-4">
            {loading ? (
              <div className="flex items-center justify-center py-12">
                <Loader2 className="w-8 h-8 animate-spin text-muted-foreground" />
              </div>
            ) : filteredDisputes.length === 0 ? (
              <EmptyState
                icon={AlertTriangle}
                title="暂无申诉记录"
                description="您还没有提交过任何申诉"
              />
            ) : (
              <VStack spacing={4}>
                {filteredDisputes.map((dispute) => {
                  const config = statusConfig[dispute.status] || statusConfig.pending;
                  const StatusIcon = getStatusIcon(dispute.status);
                  
                  return (
                    <Card key={dispute.id} className="hover:shadow-md transition-shadow">
                      <CardContent className="p-4">
                        <HStack justify="between" align="start">
                          <VStack spacing={2} className="flex-1">
                            <HStack spacing={2} align="center">
                              <StatusBadge status={config.variant} label={config.label} />
                              <span className="text-sm text-muted-foreground">
                                订单号: {dispute.orderId}
                              </span>
                            </HStack>
                            
                            <h4 className="font-medium">{dispute.reason || '申诉'}</h4>
                            
                            {dispute.description && (
                              <p className="text-sm text-muted-foreground line-clamp-2">
                                {dispute.description}
                              </p>
                            )}
                            
                            <HStack spacing={4} className="text-xs text-muted-foreground">
                              <span>
                                提交时间: {format(new Date(dispute.createdAt), 'yyyy-MM-dd HH:mm', { locale: zhCN })}
                              </span>
                              {dispute.resolvedAt && (
                                <span>
                                  处理时间: {format(new Date(dispute.resolvedAt), 'yyyy-MM-dd HH:mm', { locale: zhCN })}
                                </span>
                              )}
                            </HStack>
                          </VStack>
                          
                          <Button variant="outline" size="sm">
                            查看详情
                          </Button>
                        </HStack>
                        
                        {dispute.resolution && (
                          <div className="mt-3 pt-3 border-t border-border">
                            <p className="text-sm">
                              <span className="text-muted-foreground">处理结果: </span>
                              {dispute.resolution}
                            </p>
                          </div>
                        )}
                      </CardContent>
                    </Card>
                  );
                })}
              </VStack>
            )}
          </TabsContent>
        </Tabs>
      </VStack>
    </PageContainer>
  );
}
