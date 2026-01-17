#!/usr/bin/env node

/**
 * Environment Validation Script
 * Validates all required environment variables before build
 */

const fs = require('fs');
const path = require('path');

const REQUIRED_VARS = [
    'VITE_API_BASE_URL',
];

const OPTIONAL_VARS = [
    'VITE_CRYPTO_ENABLED',
    'VITE_CRYPTO_SECRET_KEY',
    'VITE_CRYPTO_IV',
    'VITE_CSP_ENABLED',
    'VITE_TOKEN_STORAGE',
];

const ENV_CONDITIONS = {
    VITE_CRYPTO_SECRET_KEY: (env) => {
        if (env.VITE_CRYPTO_ENABLED === 'true') {
            const key = env.VITE_CRYPTO_SECRET_KEY;
            if (!key || key.length !== 32) {
                return {
                    valid: false,
                    message: 'VITE_CRYPTO_SECRET_KEY must be exactly 32 characters when CRYPTO_ENABLED=true',
                };
            }
        }
        return { valid: true };
    },
    VITE_CRYPTO_IV: (env) => {
        if (env.VITE_CRYPTO_ENABLED === 'true') {
            const iv = env.VITE_CRYPTO_IV;
            if (!iv || iv.length !== 16) {
                return {
                    valid: false,
                    message: 'VITE_CRYPTO_IV must be exactly 16 characters when CRYPTO_ENABLED=true',
                };
            }
        }
        return { valid: true };
    },
};

function validateEnvironment(envFile = '.env.production') {
    const envPath = path.resolve(process.cwd(), envFile);

    if (!fs.existsSync(envPath)) {
        console.error(`❌ Environment file not found: ${envFile}`);
        console.error(`   Please create ${envFile} before building`);
        process.exit(1);
    }

    console.log(`\n🔍 Validating environment: ${envFile}\n`);

    // Parse .env file
    const env = {};
    const envContent = fs.readFileSync(envPath, 'utf-8');
    envContent.split('\n').forEach(line => {
        const [key, ...valueParts] = line.split('=');
        if (key && !key.startsWith('#')) {
            env[key.trim()] = valueParts.join('=').trim();
        }
    });

    // Check required variables
    let hasErrors = false;
    console.log('Required variables:');
    REQUIRED_VARS.forEach(varName => {
        const value = env[varName];
        if (value) {
            console.log(`  ✅ ${varName}`);
        } else {
            console.log(`  ❌ ${varName} - MISSING`);
            hasErrors = true;
        }
    });

    console.log('\nOptional variables:');
    OPTIONAL_VARS.forEach(varName => {
        const value = env[varName];
        if (value) {
            console.log(`  ✅ ${varName}`);
        } else {
            console.log(`  ⚠️  ${varName} - not set (optional)`);
        }
    });

    // Check conditional requirements
    console.log('\nConditional validation:');
    Object.entries(ENV_CONDITIONS).forEach(([varName, validator]) => {
        const result = validator(env);
        if (result.valid) {
            console.log(`  ✅ ${varName} condition`);
        } else {
            console.log(`  ❌ ${varName} - ${result.message}`);
            hasErrors = true;
        }
    });

    // Security warnings
    console.log('\nSecurity checks:');
    if (env.VITE_CRYPTO_ENABLED === 'true') {
        console.log('  ✅ Encryption enabled');
    } else {
        console.log('  ⚠️  Encryption DISABLED (not recommended for production)');
    }

    if (env.VITE_API_BASE_URL?.startsWith('https://')) {
        console.log('  ✅ Using HTTPS');
    } else {
        console.log('  ⚠️  Not using HTTPS (not recommended for production)');
    }

    console.log('\n' + '='.repeat(50));

    if (hasErrors) {
        console.error('\n❌ Environment validation FAILED\n');
        process.exit(1);
    } else {
        console.log('\n✅ Environment validation PASSED\n');
        process.exit(0);
    }
}

// Run validation
const envFile = process.argv[2] || '.env.production';
validateEnvironment(envFile);
