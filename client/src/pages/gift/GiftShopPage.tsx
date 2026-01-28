import { useState, useEffect, useCallback } from 'react';
import { PageContainer } from '@/components/page-container';
import { FilterBar, type FilterGroup } from '@/components/filter-bar';
import { EmptyState, SectionHeader, LoadingGrid } from '@/components/common';
import { GiftCard } from '@/components/cards';
import { ResponsiveGrid, VStack } from '@/components/layout';
import { Gift, Sparkles, Heart, Crown } from 'lucide-react';
import { itemApi, type ServiceItem } from '@/api/item';
import { toast } from 'sonner';

type GiftCategory = 'all' | 'popular' | 'premium' | 'limited';

const giftCategories: Array<{ id: GiftCategory; label: string; icon?: typeof Gift }> = [
  { id: 'all', label: '全部' },
  { id: 'popular', label: '热门', icon: Heart },
  { id: 'premium', label: '高端', icon: Crown },
  { id: 'limited', label: '限时', icon: Sparkles },
];

export default function GiftShopPage() {
  const [gifts, setGifts] = useState<ServiceItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedCategory, setSelectedCategory] = useState<GiftCategory>('all');

  const loadGifts = useCallback(async () => {
    try {
      setLoading(true);
      const params: Record<string, string> = {
        subCategory: 'gift',
      };
      if (selectedCategory !== 'all') {
        params.category = selectedCategory;
      }
      
      const data = await itemApi.getServiceItems(params);
      setGifts(data);
    } catch (error) {
      console.error('加载礼物失败:', error);
      toast.error('加载礼物失败');
    } finally {
      setLoading(false);
    }
  }, [selectedCategory]);

  useEffect(() => {
    loadGifts();
  }, [loadGifts]);

  const handleFilterChange = (key: string, value: string | number | null) => {
    if (key === 'category') {
      setSelectedCategory((value as GiftCategory) || 'all');
    }
  };

  const filterGroups: FilterGroup[] = [
    {
      key: 'category',
      type: 'badge',
      options: giftCategories.map(c => ({
        id: c.id,
        label: c.label,
        icon: c.icon,
      })),
      selectedId: selectedCategory,
    },
  ];

  const handleSend = (giftId: number) => {
    // TODO: 打开赠送礼物弹窗
    toast.info(`赠送礼物: ${giftId}`);
  };

  const handleSelect = (giftId: number) => {
    // TODO: 查看礼物详情
    toast.info(`查看礼物: ${giftId}`);
  };

  // 分离特色礼物（按排序靠前的作为特色）和普通礼物
  const featuredGifts = gifts.slice(0, 4);
  const regularGifts = gifts.slice(4);

  return (
    <PageContainer>
      <VStack spacing={6} className="py-6">
        <SectionHeader
          title="礼物商店"
          subtitle="为您喜爱的陪玩师送上心意礼物"
          icon={Gift}
          size="lg"
        />

        <FilterBar
          groups={filterGroups}
          onFilterChange={handleFilterChange}
        />

        {loading ? (
          <LoadingGrid count={8} layout="gift" />
        ) : gifts.length === 0 ? (
          <EmptyState
            icon={Gift}
            title="暂无礼物"
            description="当前分类下没有可用的礼物"
            actionLabel="查看全部"
            onAction={() => setSelectedCategory('all')}
          />
        ) : (
          <VStack spacing={8}>
            {/* 特色礼物 */}
            {featuredGifts.length > 0 && (
              <VStack spacing={4}>
                <SectionHeader
                  title="特色礼物"
                  subtitle="限时特惠，不容错过"
                  icon={Sparkles}
                  size="md"
                />
                <ResponsiveGrid preset="cards">
                  {featuredGifts.map((gift) => (
                    <GiftCard
                      key={gift.id}
                      giftId={gift.id}
                      name={gift.name}
                      imageUrl={gift.iconUrl}
                      priceCents={gift.basePriceCents}
                      description={gift.description}
                      variant="featured"
                      onSend={handleSend}
                      onSelect={handleSelect}
                    />
                  ))}
                </ResponsiveGrid>
              </VStack>
            )}

            {/* 全部礼物 */}
            <VStack spacing={4}>
              <SectionHeader
                title="全部礼物"
                size="md"
              />
              <ResponsiveGrid preset="gifts">
                {regularGifts.map((gift) => (
                  <GiftCard
                    key={gift.id}
                    giftId={gift.id}
                    name={gift.name}
                    imageUrl={gift.iconUrl}
                    priceCents={gift.basePriceCents}
                    onSend={handleSend}
                    onSelect={handleSelect}
                  />
                ))}
              </ResponsiveGrid>
            </VStack>
          </VStack>
        )}
      </VStack>
    </PageContainer>
  );
}
