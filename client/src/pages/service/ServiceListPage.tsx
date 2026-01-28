import { useState, useEffect, useCallback } from 'react';
import { useSearchParams } from 'react-router-dom';
import { PageContainer } from '@/components/page-container';
import { FilterBar, type FilterGroup } from '@/components/filter-bar';
import { EmptyState, SectionHeader, LoadingGrid } from '@/components/common';
import { ServiceCard } from '@/components/cards';
import { ResponsiveGrid, VStack } from '@/components/layout';
import { Gamepad2, Package } from 'lucide-react';
import { itemApi, type ServiceItem } from '@/api/item';
import { gameApi, type Game } from '@/api/game';
import { toast } from 'sonner';

type ServiceType = 'all' | 'accompany' | 'coaching' | 'boost';

const serviceTypes: Array<{ id: ServiceType; label: string }> = [
  { id: 'all', label: '全部' },
  { id: 'accompany', label: '陪玩' },
  { id: 'coaching', label: '教学' },
  { id: 'boost', label: '代练' },
];

export default function ServiceListPage() {
  const [searchParams, setSearchParams] = useSearchParams();
  
  const [services, setServices] = useState<ServiceItem[]>([]);
  const [games, setGames] = useState<Game[]>([]);
  const [loading, setLoading] = useState(true);
  const [gamesLoading, setGamesLoading] = useState(true);
  
  // 筛选状态
  const [selectedType, setSelectedType] = useState<ServiceType>(
    (searchParams.get('type') as ServiceType) || 'all'
  );
  const [selectedGameId, setSelectedGameId] = useState<number | null>(
    searchParams.get('gameId') ? Number(searchParams.get('gameId')) : null
  );

  // 加载游戏列表
  useEffect(() => {
    const loadGames = async () => {
      try {
        setGamesLoading(true);
        const data = await gameApi.getGames();
        setGames(data);
      } catch (error) {
        console.error('加载游戏列表失败:', error);
      } finally {
        setGamesLoading(false);
      }
    };
    loadGames();
  }, []);

  // 加载服务项目
  const loadServices = useCallback(async () => {
    try {
      setLoading(true);
      const params: Record<string, string | number> = {
        subCategory: 'service',
      };
      if (selectedGameId) {
        params.gameId = selectedGameId;
      }
      if (selectedType !== 'all') {
        params.type = selectedType;
      }
      
      const data = await itemApi.getServiceItems(params);
      setServices(data);
    } catch (error) {
      console.error('加载服务项目失败:', error);
      toast.error('加载服务项目失败');
    } finally {
      setLoading(false);
    }
  }, [selectedType, selectedGameId]);

  useEffect(() => {
    loadServices();
  }, [loadServices]);

  // 更新 URL 参数
  useEffect(() => {
    const params = new URLSearchParams();
    if (selectedType !== 'all') {
      params.set('type', selectedType);
    }
    if (selectedGameId) {
      params.set('gameId', String(selectedGameId));
    }
    setSearchParams(params, { replace: true });
  }, [selectedType, selectedGameId, setSearchParams]);

  const handleFilterChange = (key: string, value: string | number | null) => {
    if (key === 'type') {
      setSelectedType((value as ServiceType) || 'all');
    } else if (key === 'game') {
      setSelectedGameId(value as number | null);
    }
  };

  const filterGroups: FilterGroup[] = [
    {
      key: 'type',
      type: 'badge',
      options: serviceTypes,
      selectedId: selectedType,
    },
    {
      key: 'game',
      type: 'dropdown',
      dropdownOptions: games.map(g => ({ id: g.id, label: g.name })),
      selectedId: selectedGameId,
      placeholder: '选择游戏',
      icon: Gamepad2,
      loading: gamesLoading,
    },
  ];

  const handleBuy = (serviceId: number) => {
    // TODO: 跳转到下单页面
    toast.info(`购买服务: ${serviceId}`);
  };

  const handleSelect = (serviceId: number) => {
    // TODO: 跳转到服务详情页面
    toast.info(`查看服务: ${serviceId}`);
  };

  return (
    <PageContainer>
      <VStack spacing={6} className="py-6">
        <SectionHeader
          title="服务项目"
          subtitle="浏览各类游戏服务，找到适合您的陪玩或代练"
          icon={Package}
          size="lg"
        />

        <FilterBar
          groups={filterGroups}
          onFilterChange={handleFilterChange}
        />

        {loading ? (
          <LoadingGrid count={6} layout="service" />
        ) : services.length === 0 ? (
          <EmptyState
            icon={Package}
            title="暂无服务项目"
            description="当前筛选条件下没有找到服务项目"
            actionLabel="清除筛选"
            onAction={() => {
              setSelectedType('all');
              setSelectedGameId(null);
            }}
          />
        ) : (
          <ResponsiveGrid preset="services">
            {services.map((service) => (
              <ServiceCard
                key={service.id}
                serviceId={service.id}
                name={service.name}
                description={service.description}
                coverUrl={service.iconUrl}
                priceCents={service.basePriceCents}
                unit={`/${service.serviceHours || 1}小时`}
                requiredPlayers={service.requiredPlayers}
                durationMinutes={(service.serviceHours || 1) * 60}
                category={service.category}
                onBuy={handleBuy}
                onSelect={handleSelect}
              />
            ))}
          </ResponsiveGrid>
        )}
      </VStack>
    </PageContainer>
  );
}
