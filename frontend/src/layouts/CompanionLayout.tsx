import { Outlet, Link } from 'react-router-dom';

const CompanionLayout = () => {
    return (
        <div className="app-layout" style={{ background: '#2c2f33' }}>
            <aside className="channel-sidebar" style={{ width: '200px', background: '#23272a' }}>
                <div className="server-header">Companion Panel</div>
                <div className="channel-list">
                    <div className="channel-item">Orders</div>
                    <div className="channel-item">Schedule</div>
                    <div className="channel-item">Earnings</div>
                </div>
            </aside>
            <main className="main-content">
                <header className="chat-header">
                    <span>Companion Dashboard</span>
                    <Link to="/">Back to Home</Link>
                </header>
                <div style={{ padding: '20px' }}>
                    <Outlet />
                </div>
            </main>
        </div>
    );
};

export default CompanionLayout;
