import { Outlet, Link, useLocation } from "react-router-dom";
import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { useOrderStore, useAuthStore } from "@/stores";
import { cn } from "@/lib/utils";
import {
    LayoutDashboard,
    Gamepad2,
    ShoppingBag,
    MessageSquare,
    User,
    Settings,
    Bell,
    Search,
    Users,
    Target,
    Wallet,
    // Star import removed as it is unused
    ArrowLeftRight,
    BadgeCheck,
    // TrendingUp import removed as it is unused
} from "lucide-react";
import { ModeToggle } from "@/components/mode-toggle";
import { LanguageSwitcher } from "@/components/language-switcher";

import { BottomNav } from "@/components/bottom-nav";

export default function DesktopLayout() {
    const location = useLocation();
    const { t } = useTranslation();
    const {
        user,
        viewMode,
        isPlayer,
        switchToPlayerMode,
        switchToUserMode
    } = useAuthStore();

    // 用户视图导航
    const userNavItems = [
        { icon: LayoutDashboard, label: t('nav.home'), path: "/" },
        { icon: Gamepad2, label: t('nav.players'), path: "/players" },
        { icon: Users, label: t('nav.rooms'), path: "/rooms" },
        { icon: Target, label: t('nav.lfg'), path: "/lfg" },
        { icon: ShoppingBag, label: t('nav.orders'), path: "/orders" },
        { icon: MessageSquare, label: t('nav.chat'), path: "/chat" },
        { icon: User, label: t('nav.profile'), path: "/profile" },
    ];

    // 陪玩视图导航
    const playerNavItems = [
        { icon: LayoutDashboard, label: t('nav.playerDashboard'), path: "/player/dashboard" },
        { icon: ShoppingBag, label: t('nav.playerOrders'), path: "/player/orders" },
        { icon: Wallet, label: t('nav.earnings'), path: "/player/earnings" },
        { icon: Users, label: t('nav.team'), path: "/player/team" },
        { icon: MessageSquare, label: t('nav.chat'), path: "/chat" },
        { icon: BadgeCheck, label: t('nav.verification'), path: "/player/verification/realname" },
        { icon: User, label: t('nav.playerProfile'), path: "/player/profile/edit" },
    ];

    // 根据当前视图模式选择导航项
    const navItems = viewMode === 'player' ? playerNavItems : userNavItems;

    const isLinkActive = (path: string) => {
        if (path === "/" && location.pathname === "/") return true;
        if (path !== "/" && location.pathname.startsWith(path)) return true;
        return false;
    };

    const { subscribeToOrderUpdates, unsubscribeFromOrderUpdates } = useOrderStore();

    useEffect(() => {
        // Initialize WebSocket subscription for real-time order updates
        subscribeToOrderUpdates();
        return () => {
            unsubscribeFromOrderUpdates();
        };
    }, [subscribeToOrderUpdates, unsubscribeFromOrderUpdates]);

    return (
        <div className="flex h-screen w-full bg-background text-foreground overflow-hidden font-sans selection:bg-primary/20">
            {/* Sidebar (Desktop Only) */}
            <aside className="hidden md:flex w-[240px] flex-shrink-0 flex-col border-r border-border/40 bg-background/60 backdrop-blur-xl z-20 transition-all duration-300">
                {/* Logo */}
                <div className="h-16 flex items-center px-6 border-b border-border/40">
                    <div className="flex items-center gap-3">
                        <div className="w-9 h-9 rounded-xl bg-gradient-to-br from-primary to-primary/60 flex items-center justify-center font-bold text-primary-foreground shadow-lg shadow-primary/20 ring-1 ring-white/10">
                            GL
                        </div>
                        <span className="font-bold text-xl tracking-tight bg-clip-text text-transparent bg-gradient-to-r from-foreground to-foreground/70">
                            {t('app.name')}
                        </span>
                    </div>
                </div>

                {/* Nav Menu */}
                <nav className="flex-1 p-4 space-y-1.5 overflow-y-auto">
                    <div className="text-xs font-bold text-muted-foreground px-3 py-2 opacity-70">
                        {t('nav.menu')}
                    </div>
                    {navItems.map((item) => {
                        const isActive = isLinkActive(item.path);
                        const Icon = item.icon;

                        return (
                            <Link
                                key={item.path}
                                to={item.path}
                                className={cn(
                                    "flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition-all duration-300 group relative overflow-hidden",
                                    isActive
                                        ? "bg-primary/10 text-primary shadow-sm"
                                        : "text-muted-foreground hover:bg-muted/60 hover:text-foreground"
                                )}
                            >
                                {isActive && (
                                    <div className="absolute inset-0 bg-gradient-to-r from-primary/10 to-transparent opacity-50" />
                                )}
                                <Icon className={cn("w-5 h-5 flex-shrink-0 transition-transform duration-300", isActive ? "scale-110" : "group-hover:scale-110")} />
                                <span className={cn("relative z-10 truncate", isActive && "font-semibold")}>{item.label}</span>

                                {isActive && (
                                    <div className="absolute left-0 top-1/2 -translate-y-1/2 w-1 h-6 bg-primary rounded-r-full shadow-[0_0_8px_rgba(var(--primary),0.5)]" />
                                )}
                            </Link>
                        );
                    })}
                </nav>

                {/* User Info Footer */}
                <div className="p-4 bg-gradient-to-t from-background/90 to-transparent mt-auto border-t border-border/40">
                    {/* View Mode Switcher - 只有陪玩用户才显示 */}
                    {isPlayer && (
                        <button
                            onClick={() => viewMode === 'player' ? switchToUserMode() : switchToPlayerMode()}
                            className="w-full mb-3 p-3 rounded-xl bg-gradient-to-r from-primary/10 to-purple-500/10 border border-primary/20 hover:border-primary/40 transition-all group flex items-center justify-between"
                        >
                            <div className="flex items-center gap-2">
                                <ArrowLeftRight className="w-4 h-4 text-primary" />
                                <span className="text-sm font-medium text-foreground">
                                    {viewMode === 'player' ? t('nav.switchToUser') : t('nav.switchToPlayer')}
                                </span>
                            </div>
                            <div className={cn(
                                "px-2 py-0.5 rounded-full text-xs font-medium",
                                viewMode === 'player'
                                    ? "bg-purple-500/20 text-purple-400"
                                    : "bg-primary/20 text-primary"
                            )}>
                                {viewMode === 'player' ? t('nav.playerMode') : t('nav.userMode')}
                            </div>
                        </button>
                    )}

                    <div className="p-3 rounded-2xl bg-card/40 border border-border/50 backdrop-blur-sm shadow-sm flex items-center gap-3 hover:bg-card/60 transition-colors cursor-pointer group">
                        <div className="relative w-10 h-10">
                            <div className="w-full h-full rounded-full bg-muted overflow-hidden ring-2 ring-background group-hover:ring-primary/20 transition-all">
                                {user?.avatar ? (
                                    <img src={user.avatar} alt={user.name || user.username} className="w-full h-full object-cover" />
                                ) : (
                                    <div className="w-full h-full bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-xs font-bold text-white">
                                        {(user?.name || user?.username || 'U').charAt(0).toUpperCase()}
                                    </div>
                                )}
                            </div>
                            <div className="absolute bottom-0 right-0 w-3 h-3 bg-green-500 rounded-full border-2 border-background z-50 shadow-sm"></div>
                        </div>
                        <div className="flex-1 min-w-0">
                            <div className="text-sm font-semibold truncate text-foreground group-hover:text-primary transition-colors">
                                {user?.name || user?.username || t('nav.guest')}
                            </div>
                            <div className="text-xs text-muted-foreground truncate">
                                {viewMode === 'player' ? t('nav.playerMode') : t('nav.online_status')}
                            </div>
                        </div>
                        <Link to="/profile" className="p-1.5 hover:bg-muted rounded-full text-muted-foreground hover:text-foreground transition-colors">
                            <Settings className="w-4 h-4" />
                        </Link>
                    </div>
                </div>
            </aside>

            {/* Main Content */}
            <main className="flex-1 flex flex-col min-w-0 bg-muted/10 relative">
                {/* Header */}
                <header className="h-16 flex items-center justify-between px-6 sticky top-0 z-40 backdrop-blur-md bg-background/70 border-b border-border/40 transition-all">
                    <div className="flex items-center gap-4">
                        <div className="p-2 bg-background/50 rounded-lg border border-border/50 shadow-sm">
                            <LayoutDashboard className="w-5 h-5 text-muted-foreground" />
                        </div>
                        <h2 className="font-semibold text-lg text-foreground/90 tracking-tight">
                            {t('nav.dashboard')}
                        </h2>
                    </div>

                    <div className="flex items-center gap-4">
                        <div className="relative group hidden md:block">
                            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                                <Search className="h-4 w-4 text-muted-foreground group-focus-within:text-primary transition-colors" />
                            </div>
                            <input
                                type="text"
                                placeholder={t('nav.search_placeholder')}
                                className="bg-background/50 border border-border/40 text-sm rounded-full pl-10 pr-4 py-2 w-64 text-foreground placeholder:text-muted-foreground focus:ring-2 focus:ring-primary/20 focus:border-primary/50 outline-none transition-all shadow-sm focus:w-72"
                            />
                        </div>
                        <div className="h-6 w-px bg-border/60 mx-1"></div>
                        <div className="flex items-center gap-2">
                            <LanguageSwitcher />
                            <ModeToggle />
                        </div>
                        <button className="relative p-2.5 bg-background/50 hover:bg-background border border-border/40 hover:border-border rounded-full text-muted-foreground hover:text-foreground transition-all shadow-sm">
                            <Bell className="w-5 h-5" />
                            <span className="absolute top-2 right-2.5 w-2 h-2 bg-red-500 rounded-full ring-2 ring-background animate-pulse"></span>
                        </button>
                    </div>
                </header>

                <div className="flex-1 flex flex-col min-h-0 overflow-hidden relative">
                    {/* Background decorations */}
                    <div className="absolute inset-0 pointer-events-none overflow-hidden">
                        <div className="absolute top-[-20%] right-[-10%] w-[500px] h-[500px] bg-primary/5 rounded-full blur-3xl opacity-50" />
                        <div className="absolute bottom-[-20%] left-[-10%] w-[500px] h-[500px] bg-purple-500/5 rounded-full blur-3xl opacity-50" />
                    </div>
                    {/* Added pb-20 for mobile bottom nav spacing */}
                    <div className="flex-1 relative z-10 overflow-y-scroll pb-20 md:pb-0">
                        <Outlet />
                    </div>
                </div>
            </main>

            {/* Mobile Bottom Nav */}
            <BottomNav />
        </div>
    );
}
