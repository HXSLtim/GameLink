const axios = require('axios');

const BASE_URL = 'http://localhost:5000/api/v1';

async function verifyAuth() {
    console.log('--- START VERIFY ---');

    // 0. Proxy Check
    try {
        await axios.get(`${BASE_URL}/non-existent-route-${Date.now()}`);
    } catch (e) {
        if (e.response) {
            const contentType = e.response.headers['content-type'];
            if (contentType && contentType.includes('text/html')) {
                console.error('❌ PROXY ERROR: Backend returned HTML. Is the proxy configured correctly?');
                console.error('   Response meant for API was served by Vite frontend.');
                process.exit(1);
            }
            console.log('✅ Proxy Check: Backend returned JSON/Error as expected (not HTML).');
        } else {
            console.error('❌ Proxy Check Failed: Connection Refused or Network Error.');
            process.exit(1);
        }
    }

    const ts = Math.floor(Date.now() / 1000);
    const email = `test_${ts}@example.com`;
    const password = 'password123';
    let loginUsername = email;

    try {
        // 1. Register
        console.log(`1. REG: ${email}`);
        try {
            const regRes = await axios.post(`${BASE_URL}/auth/register`, {
                email,
                password,
                name: 'TestUser'
            });
            console.log('   ✅ REG OK:', regRes.status);
            if (regRes.data?.user?.username) {
                loginUsername = regRes.data.user.username;
            }
        } catch (regErr) {
            console.error('   ❌ REG FAIL:', regErr.response ? JSON.stringify(regErr.response.data) : regErr.message);
            // Don't exit yet, try phone?
        }

        // 2. Login
        console.log(`2. LOGIN: ${loginUsername}`);
        try {
            const loginRes = await axios.post(`${BASE_URL}/auth/login`, {
                username: loginUsername,
                password
            });
            console.log('   ✅ LOGIN OK');
            const { token } = loginRes.data;

            // 3. Me
            console.log(`3. ME`);
            const meRes = await axios.get(`${BASE_URL}/auth/me`, {
                headers: { Authorization: `Bearer ${token}` }
            });
            console.log('   ✅ ME OK:', meRes.data.user?.username || 'User');
            console.log('--- SUCCESS ---');

        } catch (loginErr) {
            console.error('   ❌ LOGIN FAIL:', loginErr.response ? JSON.stringify(loginErr.response.data) : loginErr.message);
            process.exit(1);
        }

    } catch (error) {
        console.error('❌ FATAL:', error.message);
        process.exit(1);
    }
}

verifyAuth();
