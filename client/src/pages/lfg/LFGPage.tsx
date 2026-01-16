import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Plus, Search, Filter, Users } from 'lucide-react';
import { PageContainer, PageHeader } from '@/components/page-container';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';
import { LFGRequestCard, LFGRequestCardSkeleton } from '@/components/lfg';
import { useLFGStore, useAuthStore } from '@/stores';
import type { LFGRequestType } from '@/stores/modules/lfg-store';

export default function LFGPage() {
    const { t } = useTranslation();
    const navigate = useNavigate();
    const { user } = useAuthStore();
    const {
        requests,
        myRequests,
        activeRequest,
        isLoading,
        fetchPendingRequests,
        fetchMyRequests,
        fetchActiveRequest,
        acceptRequest,
        cancelRequest,
    } = useLFGStore();

    const [typeFilter, setTypeFilter] = useState<LFGRequestType | 'all'>('all');
    const [searchQuery, setSearchQuery] = useState('');
    const [acceptingId, setAcceptingId] = useState<number | null>(null);

    useEffect(() => {
        fetchPendingRequests();
        fetchMyRequests();
        fetchActiveRequest();
    }, [fetchPendingRequests, fetchMyRequests, fetchActiveRequest]);

    const handleCreateRequest = () => {
        navigate('/lfg/create');
    };

    const handleAccept = async (requestId: number) => {
        setAcceptingId(requestId);
        try {
            const room = await acceptRequest(requestId);
            if (room?.id) {
                navigate(`/rooms/${room.id}`);
            }
        } finally {
            setAcceptingId(null);
        }
    };

    const handleCancel = async (requestId: number) => {
        await cancelRequest(requestId);
        fetchActiveRequest();
    };

    const filteredRequests = requests.filter((request) => {
        if (typeFilter !== 'all' && request.requestType !== typeFilter) {
            return false;
        }
        if (searchQuery) {
            const query = searchQuery.toLowerCase();
            return (
                request.title?.toLowerCase().includes(query) ||
                request.description?.toLowerCase().includes(query) ||
                request.gameName?.toLowerCase().includes(query) ||
                request.userNickname?.toLowerCase().includes(query)
            );
        }
        return true;
    });

    return (
        <PageContainer>
            <PageHeader
                title={t('lfg.title')}
                description={t('lfg.description')}
            />

            {/* Active Request Banner */}
            {activeRequest && (
                <div className="mb-6 p-4 rounded-lg bg-primary/10 border border-primary/20">
                    <div className="flex items-center justify-between">
                        <div>
                            <h3 className="font-medium text-foreground">
                                {t('lfg.activeRequest')}
                            </h3>
                            <p className="text-sm text-muted-foreground">
                                {activeRequest.title || t('lfg.untitled')}
                            </p>
                        </div>
                        <Button
                            variant="outline"
                            size="sm"
                            onClick={() => handleCancel(activeRequest.id)}
                        >
                            {t('lfg.cancel')}
                        </Button>
                    </div>
                </div>
            )}

            {/* Actions Bar */}
            <div className="flex flex-col sm:flex-row gap-3 mb-6">
                {/* Search */}
                <div className="relative flex-1">
                    <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                    <Input
                        placeholder={t('lfg.searchPlaceholder')}
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        className="pl-9"
                    />
                </div>

                {/* Type Filter */}
                <Select
                    value={typeFilter}
                    onValueChange={(value) => setTypeFilter(value as LFGRequestType | 'all')}
                >
                    <SelectTrigger className="w-full sm:w-[180px]">
                        <Filter className="h-4 w-4 mr-2" />
                        <SelectValue placeholder={t('lfg.filterType')} />
                    </SelectTrigger>
                    <SelectContent>
                        <SelectItem value="all">{t('lfg.type.all')}</SelectItem>
                        <SelectItem value="find_player">{t('lfg.type.findPlayer')}</SelectItem>
                        <SelectItem value="find_team">{t('lfg.type.findTeam')}</SelectItem>
                    </SelectContent>
                </Select>

                {/* Create Button */}
                <Button
                    onClick={handleCreateRequest}
                    className="gap-2"
                    disabled={!!activeRequest}
                >
                    <Plus className="h-4 w-4" />
                    {t('lfg.create')}
                </Button>
            </div>

            {/* Tabs */}
            <Tabs defaultValue="browse" className="space-y-4">
                <TabsList>
                    <TabsTrigger value="browse">{t('lfg.tabs.browse')}</TabsTrigger>
                    <TabsTrigger value="my">{t('lfg.tabs.my')}</TabsTrigger>
                </TabsList>

                {/* Browse Tab */}
                <TabsContent value="browse" className="space-y-4">
                    {isLoading ? (
                        <div className="space-y-4">
                            {[1, 2, 3].map((i) => (
                                <LFGRequestCardSkeleton key={i} />
                            ))}
                        </div>
                    ) : filteredRequests.length === 0 ? (
                        <div className="text-center py-12">
                            <Users className="h-12 w-12 mx-auto mb-4 text-muted-foreground" />
                            <h3 className="text-lg font-medium text-foreground mb-2">
                                {t('lfg.empty')}
                            </h3>
                            <p className="text-sm text-muted-foreground mb-4">
                                {t('lfg.emptyDescription')}
                            </p>
                            <Button onClick={handleCreateRequest} variant="outline">
                                <Plus className="h-4 w-4 mr-2" />
                                {t('lfg.createFirst')}
                            </Button>
                        </div>
                    ) : (
                        <div className="space-y-4">
                            {filteredRequests.map((request) => (
                                <LFGRequestCard
                                    key={request.id}
                                    request={request}
                                    currentUserId={Number(user?.id)}
                                    onAccept={handleAccept}
                                    onCancel={handleCancel}
                                    isAccepting={acceptingId === request.id}
                                />
                            ))}
                        </div>
                    )}
                </TabsContent>

                {/* My Requests Tab */}
                <TabsContent value="my" className="space-y-4">
                    {myRequests.length === 0 ? (
                        <div className="text-center py-12">
                            <Users className="h-12 w-12 mx-auto mb-4 text-muted-foreground" />
                            <h3 className="text-lg font-medium text-foreground mb-2">
                                {t('lfg.myEmpty')}
                            </h3>
                            <p className="text-sm text-muted-foreground">
                                {t('lfg.myEmptyDescription')}
                            </p>
                        </div>
                    ) : (
                        <div className="space-y-4">
                            {myRequests.map((request) => (
                                <LFGRequestCard
                                    key={request.id}
                                    request={request}
                                    currentUserId={Number(user?.id)}
                                    onCancel={handleCancel}
                                />
                            ))}
                        </div>
                    )}
                </TabsContent>
            </Tabs>
        </PageContainer>
    );
}
