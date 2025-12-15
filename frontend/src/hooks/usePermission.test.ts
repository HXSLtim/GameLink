/**
 * usePermission Hook Tests
 *
 * Tests for permission checking hooks
 * **Validates: Requirements 3.1, 3.4**
 */
import { describe, it, vi, beforeEach } from 'vitest';
import { renderHook } from '@testing-library/react';
import * as fc from 'fast-check';

// Mock the AdminContext
const mockAdminContext = {
  menus: [],
  permissions: [] as string[],
  loading: false,
  refreshMenus: vi.fn(),
  hasPermission: vi.fn(),
  hasAllPermissions: vi.fn(),
  hasAnyPermission: vi.fn(),
  isSuperAdmin: false,
};

vi.mock('@/context/useAdmin', () => ({
  useAdmin: () => mockAdminContext,
}));

// Import after mocking
import {
  usePermission,
  usePermissions,
  useHasPermission,
  usePermissionChecker,
  usePermissionCheckerWithLoading,
  checkPermission,
} from './usePermission';

// Arbitrary generators for property testing
const permissionCodeArb = fc.stringMatching(/^[a-z]+\.[a-z]+\.[a-z]+$/);
const permissionListArb = fc.array(permissionCodeArb, { minLength: 0, maxLength: 10 });
const modeArb = fc.constantFrom('any' as const, 'all' as const);

describe('checkPermission (pure function)', () => {
  /**
   * Property: Empty permission list should always return true
   */
  it('should return true when no permissions are required', () => {
    fc.assert(
      fc.property(permissionListArb, (userPermissions) => {
        return checkPermission(userPermissions, [], 'any') === true;
        return checkPermission(userPermissions, [], 'all') === true;
      }),
      { numRuns: 100 }
    );
  });

  /**
   * Property: Super admin (*) should always have permission
   */
  it('should return true for super admin regardless of required permissions', () => {
    fc.assert(
      fc.property(permissionListArb, modeArb, (requiredPermissions, mode) => {
        const userPermissions = ['*'];
        return checkPermission(userPermissions, requiredPermissions, mode) === true;
      }),
      { numRuns: 100 }
    );
  });

  /**
   * Property: User with no permissions should fail non-empty requirements
   */
  it('should return false when user has no permissions and requirements exist', () => {
    fc.assert(
      fc.property(
        fc.array(permissionCodeArb, { minLength: 1, maxLength: 5 }),
        modeArb,
        (requiredPermissions, mode) => {
          return checkPermission([], requiredPermissions, mode) === false;
        }
      ),
      { numRuns: 100 }
    );
  });

  /**
   * Property: 'any' mode should pass if at least one permission matches
   */
  it('should return true in any mode when user has at least one required permission', () => {
    fc.assert(
      fc.property(
        fc.array(permissionCodeArb, { minLength: 2, maxLength: 5 }),
        fc.integer({ min: 0 }),
        (requiredPermissions, selectIndex) => {
          const actualIndex = selectIndex % requiredPermissions.length;
          const userPermissions = [requiredPermissions[actualIndex]];
          return checkPermission(userPermissions, requiredPermissions, 'any') === true;
        }
      ),
      { numRuns: 100 }
    );
  });

  /**
   * Property: 'all' mode should fail if any permission is missing
   */
  it('should return false in all mode when user is missing any required permission', () => {
    fc.assert(
      fc.property(
        fc.array(permissionCodeArb, { minLength: 2, maxLength: 5 }),
        fc.integer({ min: 0 }),
        (requiredPermissions, removeIndex) => {
          const actualIndex = removeIndex % requiredPermissions.length;
          const userPermissions = requiredPermissions.filter((_, i) => i !== actualIndex);
          return checkPermission(userPermissions, requiredPermissions, 'all') === false;
        }
      ),
      { numRuns: 100 }
    );
  });

  /**
   * Property: 'all' mode should pass when user has all required permissions
   */
  it('should return true in all mode when user has all required permissions', () => {
    fc.assert(
      fc.property(permissionListArb, (requiredPermissions) => {
        // User has all required permissions plus some extra
        const userPermissions = [...requiredPermissions, 'extra.permission.code'];
        return checkPermission(userPermissions, requiredPermissions, 'all') === true;
      }),
      { numRuns: 100 }
    );
  });
});


describe('usePermission Hook', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockAdminContext.loading = false;
    mockAdminContext.permissions = [];
    mockAdminContext.isSuperAdmin = false;
  });

  /**
   * Property: Hook should return correct hasPermission based on user permissions
   * **Validates: Requirements 3.1**
   */
  it('should return hasPermission=true when user has the required permission', () => {
    fc.assert(
      fc.property(permissionCodeArb, (permission) => {
        mockAdminContext.permissions = [permission];
        mockAdminContext.loading = false;

        const { result } = renderHook(() => usePermission(permission));

        return result.current.hasPermission === true && result.current.loading === false;
      }),
      { numRuns: 100 }
    );
  });

  /**
   * Property: Hook should return hasPermission=false when user lacks permission
   */
  it('should return hasPermission=false when user lacks the required permission', () => {
    fc.assert(
      fc.property(permissionCodeArb, (permission) => {
        mockAdminContext.permissions = [];
        mockAdminContext.loading = false;

        const { result } = renderHook(() => usePermission(permission));

        return result.current.hasPermission === false && result.current.loading === false;
      }),
      { numRuns: 100 }
    );
  });

  /**
   * Property: Hook should return hasPermission=false during loading
   */
  it('should return hasPermission=false when loading', () => {
    fc.assert(
      fc.property(permissionCodeArb, (permission) => {
        mockAdminContext.permissions = [permission]; // Even with permission
        mockAdminContext.loading = true;

        const { result } = renderHook(() => usePermission(permission));

        return result.current.hasPermission === false && result.current.loading === true;
      }),
      { numRuns: 100 }
    );
  });

  /**
   * Property: Super admin should always have permission
   */
  it('should return hasPermission=true for super admin', () => {
    fc.assert(
      fc.property(permissionCodeArb, (permission) => {
        mockAdminContext.permissions = ['*'];
        mockAdminContext.loading = false;

        const { result } = renderHook(() => usePermission(permission));

        return result.current.hasPermission === true;
      }),
      { numRuns: 100 }
    );
  });

  /**
   * Property: Multiple permissions with 'any' mode
   */
  it('should return true in any mode when user has at least one permission', () => {
    fc.assert(
      fc.property(
        fc.array(permissionCodeArb, { minLength: 2, maxLength: 5 }),
        fc.integer({ min: 0 }),
        (permissions, selectIndex) => {
          const actualIndex = selectIndex % permissions.length;
          mockAdminContext.permissions = [permissions[actualIndex]];
          mockAdminContext.loading = false;

          const { result } = renderHook(() => usePermission(permissions, 'any'));

          return result.current.hasPermission === true;
        }
      ),
      { numRuns: 100 }
    );
  });

  /**
   * Property: Multiple permissions with 'all' mode
   */
  it('should return false in all mode when user is missing any permission', () => {
    fc.assert(
      fc.property(
        fc.array(permissionCodeArb, { minLength: 2, maxLength: 5 }),
        fc.integer({ min: 0 }),
        (permissions, removeIndex) => {
          const actualIndex = removeIndex % permissions.length;
          mockAdminContext.permissions = permissions.filter((_, i) => i !== actualIndex);
          mockAdminContext.loading = false;

          const { result } = renderHook(() => usePermission(permissions, 'all'));

          return result.current.hasPermission === false;
        }
      ),
      { numRuns: 100 }
    );
  });
});

describe('usePermissions Hook', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockAdminContext.loading = false;
    mockAdminContext.permissions = [];
    mockAdminContext.isSuperAdmin = false;
  });

  /**
   * Property: Batch permission check should return correct results for each permission
   * **Validates: Requirements 3.4**
   */
  it('should return correct boolean for each permission in the map', () => {
    fc.assert(
      fc.property(
        fc.array(permissionCodeArb, { minLength: 1, maxLength: 5 }),
        fc.array(fc.boolean(), { minLength: 1, maxLength: 5 }),
        (permissions, hasPermissions) => {
          // Ensure arrays have same length
          const len = Math.min(permissions.length, hasPermissions.length);
          const perms = permissions.slice(0, len);
          const has = hasPermissions.slice(0, len);

          // Build permission map
          const permissionMap: Record<string, string> = {};
          perms.forEach((p, i) => {
            permissionMap[`perm${i}`] = p;
          });

          // Set user permissions based on hasPermissions array
          mockAdminContext.permissions = perms.filter((_, i) => has[i]);
          mockAdminContext.loading = false;

          const { result } = renderHook(() => usePermissions(permissionMap));

          // Verify each permission result
          return perms.every((_, i) => {
            const key = `perm${i}` as keyof typeof result.current;
            return result.current[key] === has[i];
          });
        }
      ),
      { numRuns: 100 }
    );
  });

  /**
   * Property: Super admin should have all permissions in batch check
   */
  it('should return true for all permissions when user is super admin', () => {
    fc.assert(
      fc.property(
        fc.array(permissionCodeArb, { minLength: 1, maxLength: 5 }),
        (permissions) => {
          const permissionMap: Record<string, string> = {};
          permissions.forEach((p, i) => {
            permissionMap[`perm${i}`] = p;
          });

          mockAdminContext.permissions = ['*'];
          mockAdminContext.loading = false;

          const { result } = renderHook(() => usePermissions(permissionMap));

          return permissions.every((_, i) => {
            const key = `perm${i}` as keyof typeof result.current;
            return result.current[key] === true;
          });
        }
      ),
      { numRuns: 100 }
    );
  });

  /**
   * Property: Loading state should return false for all permissions
   */
  it('should return false for all permissions when loading', () => {
    fc.assert(
      fc.property(
        fc.array(permissionCodeArb, { minLength: 1, maxLength: 5 }),
        (permissions) => {
          const permissionMap: Record<string, string> = {};
          permissions.forEach((p, i) => {
            permissionMap[`perm${i}`] = p;
          });

          mockAdminContext.permissions = permissions; // Even with all permissions
          mockAdminContext.loading = true;

          const { result } = renderHook(() => usePermissions(permissionMap));

          return (
            result.current.loading === true &&
            permissions.every((_, i) => {
              const key = `perm${i}` as keyof typeof result.current;
              return result.current[key] === false;
            })
          );
        }
      ),
      { numRuns: 100 }
    );
  });
});


describe('useHasPermission Hook', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockAdminContext.loading = false;
    mockAdminContext.permissions = [];
    mockAdminContext.isSuperAdmin = false;
  });

  /**
   * Property: Should return boolean directly
   * **Validates: Requirements 3.1**
   */
  it('should return true when user has permission', () => {
    fc.assert(
      fc.property(permissionCodeArb, (permission) => {
        mockAdminContext.permissions = [permission];
        mockAdminContext.loading = false;

        const { result } = renderHook(() => useHasPermission(permission));

        return result.current === true;
      }),
      { numRuns: 100 }
    );
  });

  it('should return false when user lacks permission', () => {
    fc.assert(
      fc.property(permissionCodeArb, (permission) => {
        mockAdminContext.permissions = [];
        mockAdminContext.loading = false;

        const { result } = renderHook(() => useHasPermission(permission));

        return result.current === false;
      }),
      { numRuns: 100 }
    );
  });
});

describe('usePermissionChecker Hook', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockAdminContext.loading = false;
    mockAdminContext.permissions = [];
    mockAdminContext.isSuperAdmin = false;
  });

  /**
   * Property: Checker function should return correct results for dynamic checks
   * **Validates: Requirements 3.4**
   */
  it('should return a function that correctly checks permissions', () => {
    fc.assert(
      fc.property(
        permissionListArb,
        permissionCodeArb,
        (userPermissions, checkPermission) => {
          mockAdminContext.permissions = userPermissions;
          mockAdminContext.loading = false;

          const { result } = renderHook(() => usePermissionChecker());
          const checker = result.current;

          const expected = userPermissions.includes(checkPermission) || userPermissions.includes('*');
          return checker(checkPermission) === expected;
        }
      ),
      { numRuns: 100 }
    );
  });

  /**
   * Property: Checker should support multiple permissions with any mode
   */
  it('should check multiple permissions with any mode', () => {
    fc.assert(
      fc.property(
        fc.array(permissionCodeArb, { minLength: 2, maxLength: 5 }),
        fc.integer({ min: 0 }),
        (permissions, selectIndex) => {
          const actualIndex = selectIndex % permissions.length;
          mockAdminContext.permissions = [permissions[actualIndex]];
          mockAdminContext.loading = false;

          const { result } = renderHook(() => usePermissionChecker());
          const checker = result.current;

          return checker(permissions, 'any') === true;
        }
      ),
      { numRuns: 100 }
    );
  });

  /**
   * Property: Checker should support multiple permissions with all mode
   */
  it('should check multiple permissions with all mode', () => {
    fc.assert(
      fc.property(
        fc.array(permissionCodeArb, { minLength: 2, maxLength: 5 }),
        (permissions) => {
          mockAdminContext.permissions = permissions;
          mockAdminContext.loading = false;

          const { result } = renderHook(() => usePermissionChecker());
          const checker = result.current;

          return checker(permissions, 'all') === true;
        }
      ),
      { numRuns: 100 }
    );
  });

  /**
   * Property: Super admin checker should always return true
   */
  it('should return true for super admin regardless of permission', () => {
    fc.assert(
      fc.property(permissionCodeArb, modeArb, (permission, mode) => {
        mockAdminContext.permissions = ['*'];
        mockAdminContext.loading = false;

        const { result } = renderHook(() => usePermissionChecker());
        const checker = result.current;

        return checker(permission, mode) === true;
      }),
      { numRuns: 100 }
    );
  });
});

describe('usePermissionCheckerWithLoading Hook', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockAdminContext.loading = false;
    mockAdminContext.permissions = [];
    mockAdminContext.isSuperAdmin = false;
  });

  /**
   * Property: Should return loading state and checker function
   * **Validates: Requirements 3.4**
   */
  it('should return loading state correctly', () => {
    fc.assert(
      fc.property(fc.boolean(), (isLoading) => {
        mockAdminContext.loading = isLoading;
        mockAdminContext.permissions = ['some.permission.code'];

        const { result } = renderHook(() => usePermissionCheckerWithLoading());

        return result.current.loading === isLoading;
      }),
      { numRuns: 100 }
    );
  });

  /**
   * Property: Checker should return false during loading
   */
  it('should return false from checker when loading', () => {
    fc.assert(
      fc.property(permissionCodeArb, (permission) => {
        mockAdminContext.loading = true;
        mockAdminContext.permissions = [permission]; // Even with permission

        const { result } = renderHook(() => usePermissionCheckerWithLoading());
        const { check, loading } = result.current;

        return loading === true && check(permission) === false;
      }),
      { numRuns: 100 }
    );
  });

  /**
   * Property: Checker should work correctly when not loading
   */
  it('should return correct result from checker when not loading', () => {
    fc.assert(
      fc.property(
        permissionListArb,
        permissionCodeArb,
        (userPermissions, checkPerm) => {
          mockAdminContext.loading = false;
          mockAdminContext.permissions = userPermissions;

          const { result } = renderHook(() => usePermissionCheckerWithLoading());
          const { check, loading } = result.current;

          const expected = userPermissions.includes(checkPerm) || userPermissions.includes('*');
          return loading === false && check(checkPerm) === expected;
        }
      ),
      { numRuns: 100 }
    );
  });
});
