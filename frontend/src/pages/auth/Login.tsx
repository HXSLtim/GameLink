import { useNavigate, useLocation } from 'react-router-dom';

const Login = () => {
    const navigate = useNavigate();
    const location = useLocation();
    const from = location.state?.from?.pathname || '/';

    const handleLogin = (role: string) => {
        // Mock login logic - in real app, this would set auth token and user info
        localStorage.setItem('user_role', role);

        if (role === 'ADMIN') navigate('/admin');
        else if (role === 'COMPANION') navigate('/companion');
        else navigate('/');
    };

    return (
        <div style={{
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            justifyContent: 'center',
            height: '100vh',
            gap: '20px'
        }}>
            <h1>Login</h1>
            <div style={{ display: 'flex', gap: '10px' }}>
                <button onClick={() => handleLogin('USER')} className="cta-button">Login as User</button>
                <button onClick={() => handleLogin('COMPANION')} className="cta-button">Login as Companion</button>
                <button onClick={() => handleLogin('ADMIN')} className="cta-button">Login as Admin</button>
            </div>
        </div>
    );
};

export default Login;
