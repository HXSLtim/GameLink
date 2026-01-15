import { PageContainer } from '@/components/page-container';
import { Button } from '@/components/ui/button';
import { ArrowLeft } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

export default function PrivacyPage() {
    const navigate = useNavigate();
    const { t } = useTranslation();

    return (
        <PageContainer>
            <div className="max-w-4xl mx-auto py-8 px-4">
                <div className="mb-6 flex items-center gap-4">
                    <Button variant="ghost" size="icon" onClick={() => navigate(-1)}>
                        <ArrowLeft className="h-5 w-5" />
                    </Button>
                    <h1 className="text-3xl font-bold">{t('legal.privacy_title', { defaultValue: 'Privacy Policy' })}</h1>
                </div>

                <div className="prose prose-sm dark:prose-invert max-w-none bg-card p-8 rounded-xl border shadow-sm">
                    <p className="lead">Last updated: January 1, 2024</p>

                    <h3>1. Information We Collect</h3>
                    <p>We collect information you provide directly to us, such as when you create an account, update your profile, or communicate with us.</p>

                    <h3>2. How We Use Your Information</h3>
                    <p>We use the information we collect to provide, maintain, and improve our services, to process your transactions, and to communicate with you.</p>

                    <h3>3. Information Sharing</h3>
                    <p>We do not share your personal information with third parties except as described in this policy or with your consent.</p>

                    <h3>4. Data Security</h3>
                    <p>We take reasonable measures to help protect information about you from loss, theft, misuse and unauthorized access, disclosure, alteration and destruction.</p>

                    <p>We use cookies and similar technologies to collect information about your activity, browser, and device.</p>
                </div>
            </div>
        </PageContainer>
    );
}
