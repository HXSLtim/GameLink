import { Outlet, Link } from 'react-router-dom';
import '../App.css'; // Re-use existing styles for now or create new ones

const ClientLayout = () => {
    return (
        <div className="app-layout">
            <nav className="server-sidebar">
                <div className="server-icon home">🏠</div>
                <div className="server-separator" />
                {/* User specific navigation */}
            </nav>
            <div style={{ flex: 1, display: 'flex', flexDirection: 'column' }}>
                <header className="chat-header" style={{ justifyContent: 'space-between' }}>
                    <span>GameLink Client</span>
                    <div>
                        <Link to="/login" style={{ marginRight: '10px' }}>Login</Link>
                        <Link to="/companion">Switch to Companion</Link>
                    </div>
                </header>
                <main className="main-content" style={{ position: 'relative' }}>
                    <Outlet />
                </main>
            </div>
        </div>
    );
};

export default ClientLayout;
