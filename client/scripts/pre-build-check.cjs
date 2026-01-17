#!/usr/bin/env node

/**
 * Pre-build Validation Script
 * Runs all checks before production build
 */

const { execSync } = require('child_process');
const path = require('path');

const checks = [
    {
        name: 'TypeScript Compilation',
        command: 'npx tsc --noEmit',
        critical: true,
    },
    {
        name: 'ESLint',
        command: 'npm run lint',
        critical: true,
    },
    {
        name: 'Unit Tests',
        command: 'npm run test:run',
        critical: true,
    },
    {
        name: 'Security Audit',
        command: 'npm run security:audit',
        critical: false,
    },
];

function runCheck(check) {
    console.log(`\n🔍 Running: ${check.name}`);
    console.log('─'.repeat(50));

    try {
        execSync(check.command, {
            stdio: 'inherit',
            cwd: path.resolve(__dirname, '..'),
        });
        console.log(`✅ ${check.name} - PASSED`);
        return true;
    } catch (error) {
        console.error(`❌ ${check.name} - FAILED`);
        if (check.critical) {
            console.error('   This is a CRITICAL check and must pass before building.');
        }
        return false;
    }
}

function main() {
    console.log('\n🚀 Pre-build Checks');
    console.log('='.repeat(50));

    const results = checks.map(check => ({
        name: check.name,
        passed: runCheck(check),
        critical: check.critical,
    }));

    console.log('\n' + '='.repeat(50));
    console.log('📊 Summary');
    console.log('='.repeat(50));

    results.forEach(result => {
        const status = result.passed ? '✅' : '❌';
        const critical = result.critical ? ' [CRITICAL]' : '';
        console.log(`${status} ${result.name}${critical}`);
    });

    const criticalFailed = results.some(r => !r.passed && r.critical);

    console.log('\n' + '='.repeat(50));

    if (criticalFailed) {
        console.error('\n❌ Pre-build checks FAILED\n');
        console.error('Critical checks must pass before building.');
        process.exit(1);
    } else {
        console.log('\n✅ All critical checks passed. Ready to build!\n');
        process.exit(0);
    }
}

main();
