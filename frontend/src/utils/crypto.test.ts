import { describe, it, expect } from 'vitest';
import { CryptoUtil } from './crypto';

describe('CryptoUtil', () => {
  describe('encrypt and decrypt', () => {
    it('should encrypt and decrypt string correctly', () => {
      const original = 'Hello, World!';
      const encrypted = CryptoUtil.encrypt(original);
      const decrypted = CryptoUtil.decrypt<string>(encrypted);

      expect(decrypted).toBe(original);
      expect(encrypted).not.toBe(original);
    });

    it('should encrypt and decrypt object correctly', () => {
      const original = { username: 'admin', password: '123456' };
      const encrypted = CryptoUtil.encrypt(original);
      const decrypted = CryptoUtil.decrypt<typeof original>(encrypted);

      expect(decrypted).toEqual(original);
    });

    it('should encrypt and decrypt array correctly', () => {
      const original = [1, 2, 3, 4, 5];
      const encrypted = CryptoUtil.encrypt(original);
      const decrypted = CryptoUtil.decrypt<number[]>(encrypted);

      expect(decrypted).toEqual(original);
    });
  });

  describe('encryptFields and decryptFields', () => {
    it('should encrypt only specified fields', () => {
      const original = {
        username: 'admin',
        password: '123456',
        email: 'admin@example.com',
      };

      const encrypted = CryptoUtil.encryptFields(original, ['password', 'email']);

      expect(encrypted.username).toBe('admin');
      expect(encrypted.password).not.toBe('123456');
      expect(encrypted.email).not.toBe('admin@example.com');
    });

    it('should decrypt only specified fields', () => {
      const original = {
        username: 'admin',
        password: '123456',
        email: 'admin@example.com',
      };

      const encrypted = CryptoUtil.encryptFields(original, ['password', 'email']);
      const decrypted = CryptoUtil.decryptFields(encrypted, ['password', 'email']);

      expect(decrypted.username).toBe('admin');
      expect(decrypted.password).toBe('123456');
      expect(decrypted.email).toBe('admin@example.com');
    });
  });

  describe('hash functions', () => {
    it('should generate MD5 hash', () => {
      const data = 'Hello, World!';
      const hash = CryptoUtil.md5(data);

      expect(hash).toBeTruthy();
      expect(hash.length).toBe(32); // MD5 = 32 hex characters
    });

    it('should generate SHA256 hash', () => {
      const data = 'Hello, World!';
      const hash = CryptoUtil.sha256(data);

      expect(hash).toBeTruthy();
      expect(hash.length).toBe(64); // SHA-256 = 64 hex characters
    });

    it('should generate consistent hash for same input', () => {
      const data = 'Hello, World!';
      const hash1 = CryptoUtil.sha256(data);
      const hash2 = CryptoUtil.sha256(data);

      expect(hash1).toBe(hash2);
    });
  });

  describe('generateSignature', () => {
    it('should generate signature', () => {
      const data = { test: 'data' };
      const timestamp = Date.now();
      const signature = CryptoUtil.generateSignature(data, timestamp);

      expect(signature).toBeTruthy();
      expect(signature.length).toBe(64); // SHA-256
    });

    it('should generate consistent signature for same input', () => {
      const data = { test: 'data' };
      const timestamp = 1234567890;
      const signature1 = CryptoUtil.generateSignature(data, timestamp);
      const signature2 = CryptoUtil.generateSignature(data, timestamp);

      expect(signature1).toBe(signature2);
    });

    it('should generate different signatures for different data', () => {
      const timestamp = Date.now();
      const signature1 = CryptoUtil.generateSignature({ a: 1 }, timestamp);
      const signature2 = CryptoUtil.generateSignature({ a: 2 }, timestamp);

      expect(signature1).not.toBe(signature2);
    });
  });
});

