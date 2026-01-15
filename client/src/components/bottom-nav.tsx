import { Link, useLocation } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { cn } from "@/lib/utils";
import {
    LayoutDashboard,
    Gamepad2,
    ShoppingBag,
    MessageSquare,
    User
} from "lucide-react";

export function BottomNav() {
    const { t } = useTranslation();
    const location = useLocation();

    const navItems = [
        { icon: LayoutDashboard, label: t('nav.home'), path: "/" },
        { icon: Gamepad2, label: t('nav.players'), path: "/players" },
        { icon: ShoppingBag, label: t('nav.orders'), path: "/orders" },
        { icon: MessageSquare, label: t('nav.chat'), path: "/chat" },
        { icon: User, label: t('nav.profile'), path: "/profile" },
    ];

    const isLinkActive = (path: string) => {
        if (path === "/" && location.pathname === "/") return true;
        if (path !== "/" && location.pathname.startsWith(path)) return true;
        return false;
    };

    return (
        <div className="md:hidden fixed bottom-0 left-0 right-0 h-16 bg-background/80 backdrop-blur-xl border-t border-border/40 z-50 px-6 flex items-center justify-between shadow-negative-lg">
            {navItems.map((item) => {
                const isActive = isLinkActive(item.path);
                const Icon = item.icon;

                return (
                    <Link
                        key={item.path}
                        to={item.path}
                        className={cn(
                            "flex flex-col items-center justify-center gap-1 transition-all duration-300 relative",
                            isActive ? "text-primary" : "text-muted-foreground hover:text-foreground"
                        )}
                    >
                        {isActive && (
                            <div className="absolute -top-3 left-1/2 -translate-x-1/2 w-8 h-1 bg-primary rounded-b-full shadow-[0_2px_8px_rgba(var(--primary),0.5)]" />
                        )}
                        <Icon className={cn("w-6 h-6 transition-transform", isActive && "scale-110")} />
                        <span className="text-[10px] font-medium">{item.label}</span>
                    </Link>
                );
            })}
        </div>
    );
}
