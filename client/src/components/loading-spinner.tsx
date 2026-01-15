import { Loader2 } from 'lucide-react';
import { cn } from '@/lib/utils';

interface LoadingSpinnerProps {
    className?: string;
    size?: 'sm' | 'md' | 'lg';
    text?: string;
}

const sizeMap = {
    sm: 'w-4 h-4',
    md: 'w-8 h-8',
    lg: 'w-12 h-12',
};

export function LoadingSpinner({ className, size = 'md', text }: LoadingSpinnerProps) {
    return (
        <div className={cn('flex flex-col items-center justify-center gap-3', className)}>
            <Loader2 className={cn('animate-spin text-primary', sizeMap[size])} />
            {text && <p className="text-sm text-muted-foreground">{text}</p>}
        </div>
    );
}

/**
 * Full page loading spinner
 */
export function PageLoading({ text = '加载中...' }: { text?: string }) {
    return (
        <div className="flex h-screen items-center justify-center">
            <LoadingSpinner size="lg" text={text} />
        </div>
    );
}
