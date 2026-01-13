import React from "react";
import { cn } from "@/lib/utils";
import { Separator } from "@/components/ui/separator";
import { Skeleton } from "@/components/ui/skeleton";

interface PageContainerProps extends React.HTMLAttributes<HTMLDivElement> {
    children: React.ReactNode;
    scrollable?: boolean;
}

export function PageContainer({
    children,
    scrollable = true,
    className,
    ...props
}: PageContainerProps) {
    return (
        <div
            className={cn(
                "h-full flex-1 flex flex-col min-h-0", // min-h-0 is crucial for nested flex scrolling
                scrollable && "overflow-y-auto",
                className
            )}
            {...props}
        >
            {children}
        </div>
    );
}

interface PageHeaderProps {
    title: string;
    description?: string;
    className?: string;
    action?: React.ReactNode;
}

export function PageHeader({
    title,
    description,
    className,
    action,
}: PageHeaderProps) {
    return (
        <div className={cn("px-4 py-4 md:px-8", className)}>
            <div className="flex items-center justify-between">
                <div className="space-y-1">
                    <h1 className="text-2xl font-bold tracking-tight">{title}</h1>
                    {description && (
                        <p className="text-sm text-muted-foreground">{description}</p>
                    )}
                </div>
                {action && <div>{action}</div>}
            </div>
            <Separator className="my-4" />
        </div>
    );
}

export function PageHeaderSkeleton() {
    return (
        <div className="px-4 py-6 md:px-8 space-y-4">
            <div className="space-y-2">
                <Skeleton className="h-8 w-[200px]" />
                <Skeleton className="h-4 w-[300px]" />
            </div>
            <Separator />
        </div>
    );
}
