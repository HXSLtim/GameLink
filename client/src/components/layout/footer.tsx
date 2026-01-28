import { forwardRef } from 'react';
import { cn } from '@/lib/utils';
import { HStack, VStack } from './flex';

type FooterProps = React.HTMLAttributes<HTMLElement> & {
    /** 是否显示简洁版本 */
    compact?: boolean;
};

const currentYear = new Date().getFullYear();

/** 页脚组件 */
export const Footer = forwardRef<HTMLElement, FooterProps>(({
    compact = false,
    className,
    ...props
}, ref) => {
    if (compact) {
        return (
            <footer
                ref={ref}
                className={cn(
                    "py-4 text-center text-xs text-muted-foreground/60",
                    className
                )}
                {...props}
            >
                <p>© {currentYear} GameLink. All rights reserved.</p>
            </footer>
        );
    }

    return (
        <footer
            ref={ref}
            className={cn(
                "mt-auto pt-6 sm:pt-8 border-t border-border/30",
                className
            )}
            {...props}
        >
            <VStack spacing={4} align="center" className="text-center text-xs sm:text-sm text-muted-foreground">
                {/* 链接 */}
                <HStack spacing={6} className="flex-wrap justify-center">
                    <a href="/about" className="hover:text-foreground transition-colors">关于我们</a>
                    <a href="/terms" className="hover:text-foreground transition-colors">服务条款</a>
                    <a href="/privacy" className="hover:text-foreground transition-colors">隐私政策</a>
                    <a href="/contact" className="hover:text-foreground transition-colors">联系我们</a>
                    <a href="/help" className="hover:text-foreground transition-colors">帮助中心</a>
                </HStack>

                {/* 版权信息 */}
                <VStack spacing={1}>
                    <p>© {currentYear} GameLink 游戏陪玩平台. All rights reserved.</p>
                    <p className="text-muted-foreground/60">
                        <a 
                            href="https://beian.miit.gov.cn/" 
                            target="_blank" 
                            rel="noopener noreferrer"
                            className="hover:text-foreground transition-colors"
                        >
                            京ICP备XXXXXXXX号-1
                        </a>
                    </p>
                </VStack>

                {/* 免责声明 */}
                <p className="text-muted-foreground/50 max-w-md">
                    本平台仅提供信息撮合服务，不对服务内容承担责任
                </p>
            </VStack>
        </footer>
    );
});
Footer.displayName = 'Footer';
