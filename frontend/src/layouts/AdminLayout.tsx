import { Outlet, Link } from 'react-router-dom';

const AdminLayout = () => {
    return (
        <div className="app-layout" style={{ background: '#1a1a1a' }}>
            <aside className="channel-sidebar" style={{ width: '240px', background: '#111' }}>
                <div className="server-header" style={{ color: '#ff5555' }}>Admin Console</div>
                <div className="channel-list">
                    <div className="channel-category">Management</div>
                    <div className="channel-item">Users</div>
                    <div className="channel-item">Companions</div>
                    <div className="channel-category">Finance</div>
                    <div className="channel-item">Transactions</div>
                    <div className="channel-item">Payouts</div>
                </div>
            </aside>
            <main className="main-content">
                <header className="chat-header" style={{ background: '#222' }}>
                    <span>System Administration</span>
                    <Link to="/">Exit</Link>
                </header>
                <div style={{ padding: '20px' }}>
                    <Outlet />
                </div>
            </main>
        </div>
    );
};

export default AdminLayout;
