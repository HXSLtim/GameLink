import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Plus, Search, Filter, Gamepad2 } from 'lucide-react';
import { PageContainer, PageHeader } from '@/components/page-container';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';
import { RoomCard, RoomCardSkeleton } from '@/components/room';
import { useRoomStore, useAuthStore } from '@/stores';
import type { ChatGroupStatus } from '@/stores/modules/room-store';

export default function RoomListPage() {
    const { t } = useTranslation();
    const navigate = useNavigate();
    const { user } = useAuthStore();
    const { rooms, isLoading, pagination, fetchRooms, joinRoom } = useRoomStore();

    const [statusFilter, setStatusFilter] = useState<ChatGroupStatus | 'all'>('all');
    const [searchQuery, setSearchQuery] = useState('');
    const [joiningRoomId, setJoiningRoomId] = useState<number | null>(null);

    useEffect(() => {
        fetchRooms({
            page: 1,
            pageSize: 20,
            status: statusFilter === 'all' ? undefined : statusFilter,
        });
    }, [fetchRooms, statusFilter]);

    const handleCreateRoom = () => {
        navigate('/rooms/create');
    };

    const handleJoinRoom = async (roomId: number, password?: string) => {
        setJoiningRoomId(roomId);
        try {
            await joinRoom(roomId, password);
        } finally {
            setJoiningRoomId(null);
        }
    };

    const filteredRooms = rooms.filter((room) => {
        if (searchQuery) {
            const query = searchQuery.toLowerCase();
            return (
                room.name.toLowerCase().includes(query) ||
                room.gameName?.toLowerCase().includes(query) ||
                room.description?.toLowerCase().includes(query)
            );
        }
        return true;
    });

    return (
        <PageContainer>
            <PageHeader
                title={t('room.title')}
                description={t('room.description')}
            />

            {/* Actions Bar */}
            <div className="flex flex-col sm:flex-row gap-3 mb-6">
                {/* Search */}
                <div className="relative flex-1">
                    <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                    <Input
                        placeholder={t('room.searchPlaceholder')}
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        className="pl-9"
                    />
                </div>

                {/* Status Filter */}
                <Select
                    value={statusFilter}
                    onValueChange={(value) => setStatusFilter(value as ChatGroupStatus | 'all')}
                >
                    <SelectTrigger className="w-full sm:w-[180px]">
                        <Filter className="h-4 w-4 mr-2" />
                        <SelectValue placeholder={t('room.filterStatus')} />
                    </SelectTrigger>
                    <SelectContent>
                        <SelectItem value="all">{t('room.status.all')}</SelectItem>
                        <SelectItem value="waiting">{t('room.status.waiting')}</SelectItem>
                        <SelectItem value="ready">{t('room.status.ready')}</SelectItem>
                        <SelectItem value="in_game">{t('room.status.in_game')}</SelectItem>
                        <SelectItem value="finished">{t('room.status.finished')}</SelectItem>
                    </SelectContent>
                </Select>

                {/* Create Button */}
                <Button onClick={handleCreateRoom} className="gap-2">
                    <Plus className="h-4 w-4" />
                    {t('room.create')}
                </Button>
            </div>

            {/* Room Grid */}
            {isLoading ? (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {[1, 2, 3, 4, 5, 6].map((i) => (
                        <RoomCardSkeleton key={i} />
                    ))}
                </div>
            ) : filteredRooms.length === 0 ? (
                <div className="text-center py-12">
                    <Gamepad2 className="h-12 w-12 mx-auto mb-4 text-muted-foreground" />
                    <h3 className="text-lg font-medium text-foreground mb-2">
                        {t('room.empty')}
                    </h3>
                    <p className="text-sm text-muted-foreground mb-4">
                        {t('room.emptyDescription')}
                    </p>
                    <Button onClick={handleCreateRoom} variant="outline">
                        <Plus className="h-4 w-4 mr-2" />
                        {t('room.createFirst')}
                    </Button>
                </div>
            ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {filteredRooms.map((room) => (
                        <RoomCard
                            key={room.id}
                            room={room}
                            currentUserId={Number(user?.id)}
                            onJoin={handleJoinRoom}
                            isJoining={joiningRoomId === room.id}
                        />
                    ))}
                </div>
            )}

            {/* Pagination info */}
            {!isLoading && filteredRooms.length > 0 && (
                <div className="mt-6 text-center text-sm text-muted-foreground">
                    {t('room.showing', {
                        count: filteredRooms.length,
                        total: pagination.total,
                    })}
                </div>
            )}
        </PageContainer>
    );
}
