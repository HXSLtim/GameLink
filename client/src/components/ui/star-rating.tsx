import { Star } from "lucide-react";
import { cn } from "@/lib/utils";
import { useState } from "react";

interface StarRatingProps {
    max?: number;
    value?: number;
    onChange?: (value: number) => void;
    readOnly?: boolean;
    className?: string;
    size?: "sm" | "md" | "lg";
}

export function StarRating({
    max = 5,
    value = 0,
    onChange,
    readOnly = false,
    className,
    size = "md"
}: StarRatingProps) {
    const [hoverValue, setHoverValue] = useState<number | null>(null);

    const handleMouseEnter = (index: number) => {
        if (!readOnly) {
            setHoverValue(index + 1);
        }
    };

    const handleMouseLeave = () => {
        if (!readOnly) {
            setHoverValue(null);
        }
    };

    const handleClick = (index: number) => {
        if (!readOnly && onChange) {
            onChange(index + 1);
        }
    };

    const starSizes = {
        sm: "h-3 w-3",
        md: "h-5 w-5",
        lg: "h-8 w-8"
    };

    return (
        <div className={cn("flex items-center space-x-1", className)} onMouseLeave={handleMouseLeave}>
            {Array.from({ length: max }).map((_, i) => {
                const filled = (hoverValue !== null ? hoverValue : value) > i;
                return (
                    <Star
                        key={i}
                        className={cn(
                            starSizes[size],
                            "transition-all duration-200",
                            filled ? "fill-yellow-400 text-yellow-400" : "fill-muted text-muted-foreground/30",
                            !readOnly && "cursor-pointer hover:scale-110"
                        )}
                        onMouseEnter={() => handleMouseEnter(i)}
                        onClick={() => handleClick(i)}
                    />
                );
            })}
        </div>
    );
}
