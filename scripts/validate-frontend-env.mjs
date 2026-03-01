#!/usr/bin/env node

import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const repoRoot = path.resolve(__dirname, '..');

const mode = process.env.MODE || process.env.NODE_ENV || 'development';

const targets = {
  admin: {
    dir: 'admin',
    required: ['VITE_API_BASE_URL'],
    booleanKeys: [
      'VITE_CRYPTO_ENABLED',
      'VITE_CRYPTO_USE_SIGNATURE',
      'VITE_DEBUG',
      'VITE_ENABLE_WEBSOCKET',
      'VITE_ENABLE_PWA',
    ],
    integerKeys: [
      'VITE_WEBSOCKET_RECONNECT_ATTEMPTS',
      'VITE_WEBSOCKET_RECONNECT_INTERVAL',
      'VITE_WEBSOCKET_HEARTBEAT_INTERVAL',
      'VITE_DEFAULT_PAGE_SIZE',
      'VITE_MAX_PAGE_SIZE',
    ],
  },
  app: {
    dir: 'app',
    required: [],
    booleanKeys: ['VITE_ENABLE_ERROR_MONITORING'],
    integerKeys: [],
  },
};

function parseEnvFile(content) {
  const result = {};
  const lines = content.split(/\r?\n/);

  for (const rawLine of lines) {
    const line = rawLine.trim();
    if (!line || line.startsWith('#')) {
      continue;
    }
    const eqIndex = line.indexOf('=');
    if (eqIndex === -1) {
      continue;
    }

    const key = line.slice(0, eqIndex).trim();
    const value = line.slice(eqIndex + 1).trim();
    if (!key) {
      continue;
    }

    const normalizedValue = value.replace(/^["']|["']$/g, '');
    result[key] = normalizedValue;
  }

  return result;
}

function loadEnv(projectDir) {
  const env = {};
  const candidates = [
    '.env.example',
    '.env',
    '.env.local',
    `.env.${mode}`,
    `.env.${mode}.local`,
  ];

  for (const filename of candidates) {
    const envPath = path.join(projectDir, filename);
    if (!fs.existsSync(envPath)) {
      continue;
    }
    const parsed = parseEnvFile(fs.readFileSync(envPath, 'utf8'));
    Object.assign(env, parsed);
  }

  Object.assign(env, process.env);

  return env;
}

function fail(message) {
  console.error(`❌ ${message}`);
}

function warn(message) {
  console.warn(`⚠️  ${message}`);
}

function ok(message) {
  console.log(`✅ ${message}`);
}

function validateCommon(targetName, targetConfig, env) {
  const errors = [];
  const warnings = [];

  for (const key of targetConfig.required) {
    if (!env[key]) {
      errors.push(`${targetName}: missing required env ${key}`);
    }
  }

  for (const key of targetConfig.booleanKeys) {
    if (!env[key]) {
      continue;
    }
    if (!['true', 'false'].includes(env[key])) {
      errors.push(`${targetName}: ${key} must be "true" or "false", got "${env[key]}"`);
    }
  }

  for (const key of targetConfig.integerKeys) {
    if (!env[key]) {
      continue;
    }
    const parsed = Number(env[key]);
    if (!Number.isInteger(parsed) || parsed < 0) {
      errors.push(`${targetName}: ${key} must be a non-negative integer, got "${env[key]}"`);
    }
  }

  if (targetName === 'admin') {
    const cryptoEnabled = env.VITE_CRYPTO_ENABLED === 'true';

    if (cryptoEnabled) {
      if (!env.VITE_CRYPTO_SECRET_KEY) {
        errors.push('admin: VITE_CRYPTO_SECRET_KEY is required when VITE_CRYPTO_ENABLED=true');
      }
      if (!env.VITE_CRYPTO_IV) {
        errors.push('admin: VITE_CRYPTO_IV is required when VITE_CRYPTO_ENABLED=true');
      }

      if (env.VITE_CRYPTO_SECRET_KEY && ![32, 64].includes(env.VITE_CRYPTO_SECRET_KEY.length)) {
        warnings.push(
          `admin: VITE_CRYPTO_SECRET_KEY length is ${env.VITE_CRYPTO_SECRET_KEY.length}, expected 32 (raw) or 64 (hex)`
        );
      }
      if (env.VITE_CRYPTO_IV && ![16, 32].includes(env.VITE_CRYPTO_IV.length)) {
        warnings.push(
          `admin: VITE_CRYPTO_IV length is ${env.VITE_CRYPTO_IV.length}, expected 16 (raw) or 32 (hex)`
        );
      }
    }
  }

  if (targetName === 'app' && !env.VITE_API_BASE_URL) {
    warnings.push('app: VITE_API_BASE_URL not set, fallback http://localhost:8080/api/v1 will be used');
  }

  if (targetName === 'app') {
    const monitoringEnabled = env.VITE_ENABLE_ERROR_MONITORING === 'true';
    if (monitoringEnabled && !env.VITE_SENTRY_DSN) {
      if (mode === 'production') {
        errors.push('app: VITE_SENTRY_DSN is required in production when VITE_ENABLE_ERROR_MONITORING=true');
      } else {
        warnings.push('app: VITE_ENABLE_ERROR_MONITORING=true but VITE_SENTRY_DSN is empty');
      }
    }
  }

  return { errors, warnings };
}

function resolveTarget() {
  const explicit = process.argv[2];
  if (explicit && targets[explicit]) {
    return explicit;
  }

  const cwdBase = path.basename(process.cwd());
  if (targets[cwdBase]) {
    return cwdBase;
  }

  return null;
}

const targetName = resolveTarget();
if (!targetName) {
  console.error('Usage: node scripts/validate-frontend-env.mjs <admin|app>');
  process.exit(1);
}

const targetConfig = targets[targetName];
const projectDir = path.join(repoRoot, targetConfig.dir);
const env = loadEnv(projectDir);

console.log(`\nFrontend env validation target=${targetName} mode=${mode}`);
console.log(`projectDir=${projectDir}\n`);

const { errors, warnings } = validateCommon(targetName, targetConfig, env);

for (const message of errors) {
  fail(message);
}
for (const message of warnings) {
  warn(message);
}

if (errors.length > 0) {
  console.error('\nEnvironment validation failed.\n');
  process.exit(1);
}

ok('Environment validation passed');
if (warnings.length > 0) {
  console.log(`ℹ️  validation completed with ${warnings.length} warning(s)`);
}
console.log('');
