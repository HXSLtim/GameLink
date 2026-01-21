import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { ArrowLeft } from 'lucide-react';
import { PageContainer, PageHeader } from '@/components/page-container';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';
import { Card, CardContent } from '@/components/ui/card';
import { GameSelector, type Game } from '@/components/game';
import { useLFGStore } from '@/stores';
import type { LFGRequestType } from '@/stores/modules/lfg-store';

export default function CreateLFGPage() {
    const { t } = useTranslation();
    const navigate = useNavigate();
    const { createRequest, isLoading } = useLFGStore();

    const [formData, setFormData] = useState({
        gameId: 0,
        gameName: '',
        requestType: 'find_player' as LFGRequestType,
        title: '',
        description: '',
        requiredPlayers: '1',
        minRank: '',
        maxPriceCents: '',
        expireMinutes: '30',
    });

    const [errors, setErrors] = useState<Record<string, string>>({});

    const validate = () => {
        const newErrors: Record<string, string> = {};
        if (!formData.gameId) {
            newErrors.gameId = t('lfg.form.errors.gameRequired');
        }
        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!validate()) return;

        try {
            await createRequest({
                gameId: formData.gameId,
                requestType: formData.requestType,
                title: formData.title.trim() || undefined,
                description: formData.description.trim() || undefined,
                requiredPlayers: Number(formData.requiredPlayers) || 1,
                minRank: formData.minRank || undefined,
                maxPriceCents: formData.maxPriceCents
                    ? Number(formData.maxPriceCents) * 100
                    : undefined,
                expireMinutes: Number(formData.expireMinutes) || 30,
            });
            navigate('/lfg');
        } catch (error) {
            console.error('Failed to create LFG request:', error);
        }
    };

    return (
        <PageContainer>
            <div className="flex items-center gap-4 mb-6">
                <Button variant="ghost" size="icon" onClick={() => navigate('/lfg')}>
                    <ArrowLeft className="h-5 w-5" />
                </Button>
                <PageHeader
                    title={t('lfg.form.title')}
                    description={t('lfg.form.description')}
                />
            </div>

            <Card className="max-w-2xl">
                <CardContent className="pt-6">
                    <form onSubmit={handleSubmit} className="space-y-6">
                        {/* Request Type */}
                        <div className="space-y-2">
                            <Label>{t('lfg.form.fields.type')}</Label>
                            <Select
                                value={formData.requestType}
                                onValueChange={(value) =>
                                    setFormData({ ...formData, requestType: value as LFGRequestType })
                                }
                            >
                                <SelectTrigger>
                                    <SelectValue />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectItem value="find_player">
                                        {t('lfg.type.findPlayer')}
                                    </SelectItem>
                                    <SelectItem value="find_team">
                                        {t('lfg.type.findTeam')}
                                    </SelectItem>
                                </SelectContent>
                            </Select>
                            <p className="text-sm text-muted-foreground">
                                {formData.requestType === 'find_player'
                                    ? t('lfg.form.fields.typeDescriptionFindPlayer')
                                    : t('lfg.form.fields.typeDescriptionFindTeam')}
                            </p>
                        </div>

                        {/* Game Selection */}
                        <div className="space-y-2">
                            <Label>{t('lfg.form.fields.game')}</Label>
                            <GameSelector
                                value={formData.gameId || undefined}
                                onChange={(gameId: number, game: Game) => {
                                    setFormData({
                                        ...formData,
                                        gameId,
                                        gameName: game.name,
                                    });
                                    if (errors.gameId) {
                                        setErrors({ ...errors, gameId: '' });
                                    }
                                }}
                                placeholder={t('lfg.form.fields.gamePlaceholder')}
                                error={errors.gameId}
                            />
                            {errors.gameId && (
                                <p className="text-sm text-destructive">{errors.gameId}</p>
                            )}
                        </div>

                        {/* Title */}
                        <div className="space-y-2">
                            <Label htmlFor="title">{t('lfg.form.fields.title')}</Label>
                            <Input
                                id="title"
                                value={formData.title}
                                onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                                placeholder={t('lfg.form.fields.titlePlaceholder')}
                                maxLength={64}
                            />
                        </div>

                        {/* Description */}
                        <div className="space-y-2">
                            <Label htmlFor="description">{t('lfg.form.fields.description')}</Label>
                            <Textarea
                                id="description"
                                value={formData.description}
                                onChange={(e) =>
                                    setFormData({ ...formData, description: e.target.value })
                                }
                                placeholder={t('lfg.form.fields.descriptionPlaceholder')}
                                rows={3}
                                maxLength={256}
                            />
                        </div>

                        {/* Required Players */}
                        <div className="space-y-2">
                            <Label>{t('lfg.form.fields.requiredPlayers')}</Label>
                            <Select
                                value={formData.requiredPlayers}
                                onValueChange={(value) =>
                                    setFormData({ ...formData, requiredPlayers: value })
                                }
                            >
                                <SelectTrigger>
                                    <SelectValue />
                                </SelectTrigger>
                                <SelectContent>
                                    {[1, 2, 3, 4, 5].map((n) => (
                                        <SelectItem key={n} value={String(n)}>
                                            {n} {t('lfg.form.fields.people')}
                                        </SelectItem>
                                    ))}
                                </SelectContent>
                            </Select>
                        </div>

                        {/* Min Rank */}
                        <div className="space-y-2">
                            <Label htmlFor="minRank">{t('lfg.form.fields.minRank')}</Label>
                            <Input
                                id="minRank"
                                value={formData.minRank}
                                onChange={(e) => setFormData({ ...formData, minRank: e.target.value })}
                                placeholder={t('lfg.form.fields.minRankPlaceholder')}
                                maxLength={32}
                            />
                        </div>

                        {/* Max Price */}
                        <div className="space-y-2">
                            <Label htmlFor="maxPrice">{t('lfg.form.fields.maxPrice')}</Label>
                            <Input
                                id="maxPrice"
                                type="number"
                                value={formData.maxPriceCents}
                                onChange={(e) =>
                                    setFormData({ ...formData, maxPriceCents: e.target.value })
                                }
                                placeholder={t('lfg.form.fields.maxPricePlaceholder')}
                            />
                            <p className="text-sm text-muted-foreground">
                                {t('lfg.form.fields.maxPriceHint')}
                            </p>
                        </div>

                        {/* Expire Time */}
                        <div className="space-y-2">
                            <Label>{t('lfg.form.fields.expireTime')}</Label>
                            <Select
                                value={formData.expireMinutes}
                                onValueChange={(value) =>
                                    setFormData({ ...formData, expireMinutes: value })
                                }
                            >
                                <SelectTrigger>
                                    <SelectValue />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectItem value="15">15 {t('lfg.form.fields.minutes')}</SelectItem>
                                    <SelectItem value="30">30 {t('lfg.form.fields.minutes')}</SelectItem>
                                    <SelectItem value="60">60 {t('lfg.form.fields.minutes')}</SelectItem>
                                    <SelectItem value="120">120 {t('lfg.form.fields.minutes')}</SelectItem>
                                </SelectContent>
                            </Select>
                        </div>

                        {/* Submit */}
                        <div className="flex gap-3 pt-4">
                            <Button
                                type="button"
                                variant="outline"
                                onClick={() => navigate('/lfg')}
                                className="flex-1"
                            >
                                {t('common.cancel')}
                            </Button>
                            <Button type="submit" disabled={isLoading} className="flex-1">
                                {isLoading ? t('lfg.form.creating') : t('lfg.form.submit')}
                            </Button>
                        </div>
                    </form>
                </CardContent>
            </Card>
        </PageContainer>
    );
}
