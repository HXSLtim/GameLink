import { PageContainer } from '@/components/page-container';

import { Button } from '@/components/ui/button';
import { ArrowLeft } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

export default function TermsPage() {
    const navigate = useNavigate();
    const { t } = useTranslation();

    return (
        <PageContainer>
            <div className="max-w-4xl mx-auto py-8 px-4">
                <div className="mb-6 flex items-center gap-4">
                    <Button variant="ghost" size="icon" onClick={() => navigate(-1)}>
                        <ArrowLeft className="h-5 w-5" />
                    </Button>
                    <h1 className="text-3xl font-bold">{t('legal.terms_title', { defaultValue: 'Terms of Service' })}</h1>
                </div>

                <div className="prose prose-sm dark:prose-invert max-w-none bg-card p-8 rounded-xl border shadow-sm">
                    <p className="lead">Last updated: January 1, 2024</p>

                    <h3>1. Acceptance of Terms</h3>
                    <p>By accessing and using GameLink, you accept and agree to be bound by the terms and provision of this agreement.</p>

                    <h3>2. Description of Service</h3>
                    <p>GameLink provides a platform for gamers to connect, play together, and offer coaching services ("Services").</p>

                    <h3>3. User Conduct</h3>
                    <p>You agree to use the Service only for lawful purposes. You are prohibited from posting or transmitting any unlawful, threatening, libelous, defamatory, obscene, scandalous, inflammatory, pornographic, or profane material.</p>

                    <h3>4. Payment and Refunds</h3>
                    <p>Payments for services are processed securely. Refunds are subject to our refund policy and are handled on a case-by-case basis.</p>

                    <h3>5. Account Security</h3>
                    <p>You are responsible for maintaining the confidentiality of your login credentials and for all activities that occur under your account.</p>

                    <p>We reserve the right to terminate or suspend access to our Service immediately, without prior notice or liability, for any reason whatsoever, including without limitation if you breach the Terms.</p>
                </div>
            </div>
        </PageContainer>
    );
}
