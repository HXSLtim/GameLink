import { describe, it, expect, beforeEach, vi } from 'vitest';
import { getItem, setItem, removeItem, clear, isAvailable } from './storage';

describe('storage utils', () => {
  beforeEach(() => {
    localStorage.clear();
    vi.clearAllMocks();
  });

  describe('getItem', () => {
    it('should get string value', () => {
      localStorage.setItem('test', 'value');
      expect(getItem('test')).toBe('value');
    });

    it('should get JSON value', () => {
      const obj = { name: 'test', value: 123 };
      localStorage.setItem('test', JSON.stringify(obj));
      expect(getItem('test')).toEqual(obj);
    });

    it('should return default value when key not exists', () => {
      expect(getItem('nonexistent', 'default')).toBe('default');
    });

    it('should return null when key not exists and no default', () => {
      expect(getItem('nonexistent')).toBeNull();
    });

    it('should handle localStorage errors gracefully', () => {
      const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
      vi.spyOn(Storage.prototype, 'getItem').mockImplementation(() => {
        throw new Error('Storage error');
      });

      expect(getItem('test', 'default')).toBe('default');
      expect(consoleSpy).toHaveBeenCalled();

      consoleSpy.mockRestore();
    });
  });

  describe('setItem', () => {
    it('should set string value', () => {
      expect(setItem('test', 'value')).toBe(true);
      expect(localStorage.getItem('test')).toBe('value');
    });

    it('should set object value as JSON', () => {
      const obj = { name: 'test', value: 123 };
      expect(setItem('test', obj)).toBe(true);
      expect(localStorage.getItem('test')).toBe(JSON.stringify(obj));
    });

    it('should handle quota exceeded error', () => {
      const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
      const warnSpy = vi.spyOn(console, 'warn').mockImplementation(() => {});

      vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => {
        const error = new Error('Quota exceeded');
        error.name = 'QuotaExceededError';
        throw error;
      });

      expect(setItem('test', 'value')).toBe(false);
      expect(consoleSpy).toHaveBeenCalled();
      expect(warnSpy).toHaveBeenCalledWith('[Storage] localStorage quota exceeded');

      consoleSpy.mockRestore();
      warnSpy.mockRestore();
    });
  });

  describe('removeItem', () => {
    it('should remove item', () => {
      localStorage.setItem('test', 'value');
      expect(removeItem('test')).toBe(true);
      expect(localStorage.getItem('test')).toBeNull();
    });

    it('should handle errors gracefully', () => {
      const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
      vi.spyOn(Storage.prototype, 'removeItem').mockImplementation(() => {
        throw new Error('Storage error');
      });

      expect(removeItem('test')).toBe(false);
      expect(consoleSpy).toHaveBeenCalled();

      consoleSpy.mockRestore();
    });
  });

  describe('clear', () => {
    it('should clear all items', () => {
      localStorage.setItem('test1', 'value1');
      localStorage.setItem('test2', 'value2');
      expect(clear()).toBe(true);
      expect(localStorage.length).toBe(0);
    });
  });

  describe('isAvailable', () => {
    it('should return true when localStorage is available', () => {
      expect(isAvailable()).toBe(true);
    });

    it('should return false when localStorage is not available', () => {
      vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => {
        throw new Error('Storage not available');
      });

      expect(isAvailable()).toBe(false);
    });
  });
});

