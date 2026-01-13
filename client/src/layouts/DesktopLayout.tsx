import { Outlet, Link, useLocation } from "react-router-dom";
import { cn } from "@/lib/utils";
import {
    LayoutDashboard,
    Gamepad2,
    ShoppingBag,
    MessageSquare,
    User,
    Settings,
    Bell,
    Search
} from "lucide-react";
import { ModeToggle } from "@/components/mode-toggle";
import { LanguageSwitcher } from "@/components/language-switcher";

export default function DesktopLayout() {
    const location = useLocation();

    const navItems = [
        { icon: LayoutDashboard, label: "首页", path: "/" },
        { icon: Gamepad2, label: "陪玩师", path: "/players" },
        { icon: ShoppingBag, label: "订单", path: "/orders" },
        { icon: MessageSquare, label: "聊天", path: "/chat" },
        { icon: User, label: "个人中心", path: "/profile" },
    ];

    const isLinkActive = (path: string) => {
        if (path === "/" && location.pathname === "/") return true;
        if (path !== "/" && location.pathname.startsWith(path)) return true;
        return false;
    };

    return (
        <div className="flex h-screen w-full bg-background text-foreground overflow-hidden font-sans selection:bg-primary/20">
            {/* 侧边栏 */}
            <aside className="w-[240px] flex-shrink-0 bg-sidebar flex flex-col border-r border-border">
                {/* Logo 区域 */}
                <div className="h-14 flex items-center px-4 shadow-sm bg-sidebar">
                    <div className="flex items-center gap-2">
                        <div className="w-8 h-8 rounded-lg bg-primary flex items-center justify-center font-bold text-primary-foreground shadow-lg shadow-primary/20">
                            GL
                        </div>
                        <span className="font-bold text-lg tracking-tight text-sidebar-foreground">GameLink</span>
                    </div>
                </div>

                {/* 导航菜单 */}
                <nav className="flex-1 p-2 space-y-1 overflow-y-auto">
                    <div className="text-xs font-bold text-muted-foreground px-2 py-1.5 uppercase tracking-wider">
                        菜单
                    </div>
                    {navItems.map((item) => {
                        const isActive = isLinkActive(item.path);
                        const Icon = item.icon;

                        return (
                            <Link
                                key={item.path}
                                to={item.path}
                                className={cn(
                                    "flex items-center gap-3 px-3 py-2 rounded-md text-sm font-medium transition-all duration-200 group relative",
                                    isActive
                                        ? "bg-accent text-accent-foreground shadow-sm"
                                        : "text-muted-foreground hover:bg-muted/50 hover:text-foreground"
                                )}
                            >
                                <Icon className={cn("w-5 h-5 transition-colors", isActive ? "text-primary" : "text-muted-foreground group-hover:text-foreground")} />
                                <span>{item.label}</span>
                                {/* 激活指示条 */}
                                {isActive && (
                                    <div className="absolute left-0 top-1/2 -translate-y-1/2 w-1 h-8 bg-primary rounded-r-full" />
                                )}
                            </Link>
                        );
                    })}
                </nav>

                {/* 底部用户信息栏 */}
                <div className="p-3 bg-card/50 mt-auto flex items-center gap-3 border-t border-border">
                    <div className="relative w-9 h-9">
                        <div className="w-full h-full rounded-full bg-muted overflow-hidden">
                            {/* 占位头像 */}
                            <div className="w-full h-full bg-gradient-to-br from-primary to-purple-600 flex items-center justify-center text-xs font-bold text-white">
                                U
                            </div>
                        </div>
                        <div className="absolute bottom-0 right-0 w-3 h-3 bg-green-500 rounded-full border-2 border-background z-50"></div>
                    </div>
                    <div className="flex-1 min-w-0">
                        <div className="text-sm font-medium truncate text-foreground">GameLink 用户</div>
                        <div className="text-xs text-muted-foreground truncate">在线 - 游戏中</div>
                    </div>
                    <button className="p-1.5 hover:bg-muted rounded-md text-muted-foreground hover:text-foreground transition-colors">
                        <Settings className="w-4 h-4" />
                    </button>
                </div>
            </aside>

            {/* 主内容区域 */}
            <main className="flex-1 flex flex-col min-w-0 bg-background relative">
                {/* 顶部工具栏 */}
                <header className="h-14 bg-header border-b border-border flex items-center justify-between px-4 shadow-sm z-10 sticky top-0">
                    <div className="flex items-center gap-4">
                        <LayoutDashboard className="w-5 h-5 text-muted-foreground" />
                        <h2 className="font-semibold text-base text-foreground">
                            {/* 这里可以根据路由动态显示标题 */}
                            仪表盘
                        </h2>
                    </div>

                    <div className="flex items-center gap-3">
                        <div className="relative group">
                            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                                <Search className="h-4 w-4 text-muted-foreground group-focus-within:text-primary transition-colors" />
                            </div>
                            <input
                                type="text"
                                placeholder="搜索全站..."
                                className="bg-muted border-none text-sm rounded-md pl-9 pr-3 py-1.5 w-64 text-foreground placeholder:text-muted-foreground focus:ring-2 focus:ring-primary/50 outline-none transition-all"
                            />
                        </div>
                        <div className="h-6 w-px bg-border mx-1"></div>
                        <LanguageSwitcher />
                        <ModeToggle />
                        <button className="relative p-2 hover:bg-muted rounded-full text-muted-foreground hover:text-foreground transition-colors">
                            <Bell className="w-5 h-5" />
                            <span className="absolute top-1.5 right-1.5 w-2 h-2 bg-destructive rounded-full ring-2 ring-header"></span>
                        </button>
                    </div>
                </header>

                {/* 内容区域 - 由页面自己控制滚动和布局 */}
                <div className="flex-1 flex flex-col min-h-0 overflow-hidden">
                    <Outlet />
                </div>
            </main>
        </div>
    );
}
