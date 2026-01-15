const axios = require('axios');
const fs = require('fs');

const BASE_URL = 'http://localhost:5000/api/v1';

async function debug() {
    console.log('--- DEBUGGING API (Phase 2) ---');
    const ts = Math.floor(Date.now() / 1000);
    const email = `debug_${ts}@example.com`;
    const username = `user_${ts}`;
    const phone = `138${ts.toString().slice(-8)}`; // Fake phone

    // 1. Probe Routes
    const paths = ['/players', '/player', '/player/list', '/home', '/auth/captcha'];
    console.log('1. Probing Routes under /api/v1...');

    for (const p of paths) {
        try {
            const url = `${BASE_URL}${p}`;
            console.log(`   Probing ${url}...`);
            await axios.get(url);
            console.log(`   ✅ ${p} - 200 OK`);
        } catch (e) {
            const status = e.response ? e.response.status : 'ERR';
            const data = e.response ? JSON.stringify(e.response.data).slice(0, 50) : e.message;
            console.log(`   ❌ ${p} - ${status} (${data})`);
        }
    }

    // 2. Try Register with FULL payload
    console.log(`2. Testing Register (FULL PAYLOAD): POST /auth/register`);
    const payload = {
        username,
        email,
        phone,
        password: 'password123',
        name: 'Debug User',
        nickname: 'DebugNick'
    };
    console.log('   Payload:', JSON.stringify(payload));

    try {
        const regRes = await axios.post(`${BASE_URL}/auth/register`, payload);
        console.log('   ✅ Register SUCCESS');
        console.log('   Response:', JSON.stringify(regRes.data));
    } catch (e) {
        console.error('   ❌ Register FAILED');
        if (e.response) {
            console.log('   Status:', e.response.status);
            console.log('   Data:', JSON.stringify(e.response.data, null, 2));
        } else {
            console.error('   Error:', e.message);
        }
    }
}

debug();
