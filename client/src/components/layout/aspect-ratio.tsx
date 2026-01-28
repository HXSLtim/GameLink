import { forwardRef } from 'react';
import { cn } from '@/lib/utils';

type AspectRatioProps = React.HTMLAttributes<HTMLDivElement> & {
    /** 宽高比 */
    ratio?: '1:1' | '4:3' | '16:9' | '21:9' | '3:4' | '9:16' | number;
};

const ratioMap = {
    '1:1': 'aspect-square',
    '4:3': 'aspect-[4/3]',
    '16:9': 'aspect-video',
    '21:9': 'aspect-[21/9]',
    '3:4': 'aspect-[3/4]',
    '9:16': 'aspect-[9/16]',
};

/** 宽高比容器 */
export const AspectRatio = forwardRef<HTMLDivElement, AspectRatioProps>(({
    ratio = '16:9',
    className,
    children,
    ...props
}, ref) => {
    const ratioClass = typeof ratio === 'string' ? ratioMap[ratio] : `aspect-[${ratio}]`;

    return (
        <div
            ref={ref}
            className={cn('relative overflow-hidden', ratioClass, className)}
            {...props}
        >
            {children}
        </div>
    );
});
AspectRatio.displayName = 'AspectRatio';

/** 图片容器 - 预设 object-cover */
export const ImageBox = forwardRef<HTMLDivElement, AspectRatioProps & {
    src?: string;
    alt?: string;
    fallback?: React.ReactNode;
    overlay?: boolean;
}>(({ ratio = '16:9', src, alt, fallback, overlay, className, children, ...props }, ref) => {
    return (
        <AspectRatio ref={ref} ratio={ratio} className={cn('bg-muted', className)} {...props}>
            {src ? (
                <img 
                    src={src} 
                    alt={alt || ''} 
                    className="absolute inset-0 w-full h-full object-cover"
                />
            ) : (
                fallback && (
                    <div className="absolute inset-0 flex items-center justify-center">
                        {fallback}
                    </div>
                )
            )}
            {overlay && (
                <div className="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent" />
            )}
            {children}
        </AspectRatio>
    );
});
ImageBox.displayName = 'ImageBox';
