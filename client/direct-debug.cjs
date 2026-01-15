const axios = require('axios');
const fs = require('fs');

// Target the backend directly
const BACKEND_URL = 'http://localhost:8080';

async function scan() {
    const log = [];
    const logFn = (msg) => {
        console.log(msg);
        log.push(msg);
    };

    logFn(`--- DIRECT BACKEND SCAN: ${BACKEND_URL} ---`);

    // 1. Path Scan
    const paths = [
        '/api/v1/public/players',
        '/api/v1/auth/me', // protected, should be 401
        '/api/v1/players',
        '/api/v1/player/list',
        '/api/v1/users',
        '/api/v1/auth/login', // Expect 405 Method Not Allowed types if GET
        '/health',
        '/api/health',
        '/ping'
    ];

    for (const p of paths) {
        try {
            logFn(`GET ${p}`);
            const res = await axios.get(`${BACKEND_URL}${p}`);
            logFn(`   ✅ ${res.status} - Data: ${JSON.stringify(res.data).slice(0, 50)}`);
        } catch (e) {
            const status = e.response ? e.response.status : 'ERR';
            const data = e.response ? JSON.stringify(e.response.data) : e.message;
            logFn(`   ❌ ${status} - ${data.slice(0, 100)}`);
        }
    }

    // 2. Register Retry
    const ts = Math.floor(Date.now() / 1000);
    const payload = {
        username: `user_${ts}`,
        password: 'password123',
        email: `test_${ts}@qq.com`, // Try a common domain
        phone: `138${ts.toString().slice(-8)}`,
        name: 'DirectDebug',
        nickname: 'DirectNick'
    };
    logFn(`\nPOST /api/v1/auth/register`);
    logFn(`Payload: ${JSON.stringify(payload)}`);

    try {
        const res = await axios.post(`${BACKEND_URL}/api/v1/auth/register`, payload);
        logFn(`   ✅ ${res.status} - SUCCESS`);
        logFn(`   ${JSON.stringify(res.data)}`);
    } catch (e) {
        logFn(`   ❌ FAILED`);
        if (e.response) logFn(`   Status: ${e.response.status}`);
        if (e.response) logFn(`   Data: ${JSON.stringify(e.response.data, null, 2)}`);
        else logFn(`   Error: ${e.message}`);
    }

    fs.writeFileSync('direct_debug_log.txt', log.join('\n'));
}

scan();
