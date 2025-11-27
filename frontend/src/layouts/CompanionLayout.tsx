import { Outlet, Link } from 'react-router-dom';
import { ThemeToggle } from '@/components';

const CompanionLayout = () => {
    return (
        <div className="app-layout" style={{ background: 'var(--bg-primary)', color: 'var(--text-normal)', minHeight: '100vh', display: 'flex' }}>
            <aside className="channel-sidebar" style={{ width: '200px', background: 'var(--bg-secondary)' }}>
                <div className="server-header" style={{ padding: '16px', fontWeight: 'bold', borderBottom: '1px solid var(--background-modifier-active)' }}>Companion Panel</div>
                <div className="channel-list" style={{ padding: '8px' }}>
                    <div className="channel-item" style={{ padding: '8px', cursor: 'pointer', borderRadius: '4px', color: 'var(--text-muted)' }}>Orders</div>
                    <div className="channel-item" style={{ padding: '8px', cursor: 'pointer', borderRadius: '4px', color: 'var(--text-muted)' }}>Schedule</div>
                    <div className="channel-item" style={{ padding: '8px', cursor: 'pointer', borderRadius: '4px', color: 'var(--text-muted)' }}>Earnings</div>
                </div>
            </aside>
            <main className="main-content" style={{ flex: 1, display: 'flex', flexDirection: 'column' }}>
                <header className="chat-header" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '0 16px', height: '48px', borderBottom: '1px solid var(--background-modifier-active)', background: 'var(--bg-primary)' }}>
                    <span style={{ fontWeight: 'bold' }}>Companion Dashboard</span>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
                        <ThemeToggle />
                        <Link to="/" style={{ color: 'var(--text-link)' }}>Back to Home</Link>
                    </div>
                </header>
                <div style={{ padding: '20px', flex: 1, overflow: 'auto' }}>
                    <Outlet />
                </div>
            </main>
        </div>
    );
};

export default CompanionLayout;
