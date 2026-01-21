/**
 * PermissionGuard Component Property-Based Tests
 * 
 * **Feature: rbac-button-level-permission, Property 6: 前端权限检查一致性**
 * **Validates: Requirements 3.1, 3.4**
 * 
 * Property: For any user with a set of permissions and any permission requirement
 * (single or multiple with any/all mode), the PermissionGuard component should
 * render children if and only if the user satisfies the permission requirement.
 */
import { describe, it, vi, beforeEach } from 'vitest';
import { render } from '@testing-library/react';
import * as fc from 'fast-check';

// Mock the AdminContext
const mockHasPermission = vi.fn();
const mockAdminContext = {
    menus: [],
    permissions: [] as string[],
    loading: false,
    refreshMenus: vi.fn(),
    hasPermission: mockHasPermission,
    hasAllPermissions: vi.fn(),
    hasAnyPermission: vi.fn(),
    isSuperAdmin: false,
};

vi.mock('@/context/useAdmin', () => ({
    useAdmin: () => mockAdminContext,
}));

// Import after mocking
import { PermissionGuard, PermissionButton } from './PermissionGuard';

/**
 * Pure permission checking logic - extracted for property testing
 * This mirrors the ACTUAL logic in usePermission.ts checkPermission function
 *
 * IMPORTANT: The actual usePermission hook determines super admin status
 * by checking if permissions include '*'. There is NO separate isSuperAdmin field
 * used by the checkPermission function in the actual implementation.
 */
function checkPermission(
    userPermissions: string[],
    requiredPermissions: string[],
    mode: 'any' | 'all'
): boolean {
    const permissionList = requiredPermissions;

    // If no valid permissions required, allow
    const hasValidPermission = permissionList.some(p => p && p.length > 0);
    if (!hasValidPermission) {
        return true;
    }

    // Check if user has any permissions
    if (!userPermissions.length) {
        return false;
    }

    // Super admin has all permissions (determined by '*' in permissions)
    if (userPermissions.includes('*')) {
        return true;
    }

    if (mode === 'all') {
        return permissionList.every(p => userPermissions.includes(p));
    } else {
        return permissionList.some(p => userPermissions.includes(p));
    }
}

// Arbitrary generators for property testing
// Use more constrained generators to avoid edge cases
const permissionSegmentArb = fc.stringMatching(/^[a-z]{2,10}$/); // 2-10 lowercase letters
const permissionCodeArb = fc.tuple(
    permissionSegmentArb,
    permissionSegmentArb,
    permissionSegmentArb
).map(parts => `${parts[0]}.${parts[1]}.${parts[2]}`);

const permissionListArb = fc.array(permissionCodeArb, { minLength: 1, maxLength: 5 });
const userPermissionListArb = fc.array(permissionCodeArb, { minLength: 0, maxLength: 5 });
const modeArb = fc.constantFrom('any' as const, 'all' as const);

describe('PermissionGuard Property Tests', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        mockAdminContext.loading = false;
        mockAdminContext.isSuperAdmin = false;
        mockAdminContext.permissions = [];

        // Set up mockHasPermission to use the actual permission checking logic
        mockHasPermission.mockImplementation((permission: string | string[], mode: 'any' | 'all' = 'any') => {
            return checkPermission(mockAdminContext.permissions, Array.isArray(permission) ? permission : [permission], mode);
        });
    });

    /**
     * **Feature: rbac-button-level-permission, Property 6: 前端权限检查一致性**
     * **Validates: Requirements 3.1, 3.4**
     *
     * Property: For any user permissions and required permissions with any mode,
     * the component renders children iff the permission check passes.
     */
    it('Property 6: should render children if and only if user has required permissions', () => {
        fc.assert(
            fc.property(
                userPermissionListArb,
                permissionListArb,
                modeArb,
                (userPermissions, requiredPermissions, mode) => {
                    // Setup mock context
                    mockAdminContext.permissions = userPermissions;
                    mockAdminContext.isSuperAdmin = false;

                    // Calculate expected result using pure function
                    const expectedHasPermission = checkPermission(
                        userPermissions,
                        requiredPermissions,
                        mode
                    );

                    const testId = `test-child-${Math.random()}`;
                    const { container } = render(
                        <PermissionGuard
                            permission={requiredPermissions}
                            mode={mode}
                        >
                            <div data-testid={testId}>Protected Content</div>
                        </PermissionGuard>
                    );

                    // Check if children are rendered
                    const childElement = container.querySelector(`[data-testid="${testId}"]`);
                    const isRendered = childElement !== null;

                    return isRendered === expectedHasPermission;
                }
            ),
            { numRuns: 100 }
        );
    });

    /**
     * Property: Super admin should always see protected content
     * Super admin is determined by having '*' in permissions
     */
    it('Property 6a: super admin should always have access regardless of permissions', () => {
        fc.assert(
            fc.property(
                permissionListArb,
                modeArb,
                (requiredPermissions, mode) => {
                    // Setup as super admin by adding '*' to permissions
                    mockAdminContext.permissions = ['*'];
                    mockAdminContext.isSuperAdmin = true;

                    const testId = `super-admin-${Math.random()}`;
                    const { container } = render(
                        <PermissionGuard
                            permission={requiredPermissions.length > 0 ? requiredPermissions : ['any.permission.code']}
                            mode={mode}
                        >
                            <div data-testid={testId}>Super Admin Content</div>
                        </PermissionGuard>
                    );

                    const childElement = container.querySelector(`[data-testid="${testId}"]`);
                    return childElement !== null;
                }
            ),
            { numRuns: 100 }
        );
    });

    /**
     * Property: 'all' mode requires all permissions
     */
    it('Property 6b: all mode should require every permission in the list', () => {
        fc.assert(
            fc.property(
                fc.array(permissionCodeArb, { minLength: 2, maxLength: 5 }),
                fc.integer({ min: 0 }),
                (requiredPermissions, removeIndex) => {
                    // User has all but one permission
                    const actualRemoveIndex = removeIndex % requiredPermissions.length;
                    const userPermissions = requiredPermissions.filter((_, i) => i !== actualRemoveIndex);

                    mockAdminContext.permissions = userPermissions;
                    mockAdminContext.isSuperAdmin = false;

                    const testId = `all-mode-${Math.random()}`;
                    const { container } = render(
                        <PermissionGuard permission={requiredPermissions} mode="all">
                            <div data-testid={testId}>All Mode Content</div>
                        </PermissionGuard>
                    );

                    const childElement = container.querySelector(`[data-testid="${testId}"]`);
                    // Should NOT render because user is missing one permission
                    return childElement === null;
                }
            ),
            { numRuns: 100 }
        );
    });

    /**
     * Property: 'any' mode requires at least one permission
     */
    it('Property 6c: any mode should pass if user has at least one required permission', () => {
        fc.assert(
            fc.property(
                fc.array(permissionCodeArb, { minLength: 2, maxLength: 5 }),
                fc.integer({ min: 0 }),
                (requiredPermissions, selectIndex) => {
                    // User has exactly one of the required permissions
                    const actualSelectIndex = selectIndex % requiredPermissions.length;
                    const userPermissions = [requiredPermissions[actualSelectIndex]];

                    mockAdminContext.permissions = userPermissions;
                    mockAdminContext.isSuperAdmin = false;

                    const testId = `any-mode-${Math.random()}`;
                    const { container } = render(
                        <PermissionGuard permission={requiredPermissions} mode="any">
                            <div data-testid={testId}>Any Mode Content</div>
                        </PermissionGuard>
                    );

                    const childElement = container.querySelector(`[data-testid="${testId}"]`);
                    // Should render because user has at least one permission
                    return childElement !== null;
                }
            ),
            { numRuns: 100 }
        );
    });

    /**
     * Property: Loading state should not render children
     */
    it('Property 6d: loading state should not render children (avoid permission flicker)', () => {
        fc.assert(
            fc.property(
                permissionListArb,
                (requiredPermissions) => {
                    mockAdminContext.loading = true;
                    mockAdminContext.permissions = ['*']; // Even with all permissions
                    mockAdminContext.isSuperAdmin = true;

                    const testId = `loading-${Math.random()}`;
                    const loadingId = `loading-content-${Math.random()}`;
                    const { container } = render(
                        <PermissionGuard
                            permission={requiredPermissions.length > 0 ? requiredPermissions : ['any.permission.code']}
                            loading={<div data-testid={loadingId}>Loading...</div>}
                        >
                            <div data-testid={testId}>Content</div>
                        </PermissionGuard>
                    );

                    const childElement = container.querySelector(`[data-testid="${testId}"]`);
                    const loadingElement = container.querySelector(`[data-testid="${loadingId}"]`);
                    // Loading content should be shown, children should not
                    return childElement === null && loadingElement !== null;
                }
            ),
            { numRuns: 100 }
        );
    });

    /**
     * Property: Disabled mode should render disabled children when no permission
     */
    it('Property 6e: disabled mode should render disabled element when user lacks permission', () => {
        fc.assert(
            fc.property(
                permissionCodeArb,
                (requiredPermission) => {
                    mockAdminContext.permissions = [];
                    mockAdminContext.isSuperAdmin = false;
                    mockAdminContext.loading = false;

                    const { container } = render(
                        <PermissionGuard
                            permission={requiredPermission}
                            disabled={true}
                        >
                            <button>Action</button>
                        </PermissionGuard>
                    );

                    // In disabled mode, the button should be rendered but disabled
                    const button = container.querySelector('button');
                    return button !== null && button.disabled === true;
                }
            ),
            { numRuns: 100 }
        );
    });

    /**
     * Property: Fallback should be rendered when no permission and not in disabled mode
     */
    it('Property 6f: fallback content should render when user lacks permission', () => {
        fc.assert(
            fc.property(
                permissionCodeArb,
                (requiredPermission) => {
                    mockAdminContext.permissions = [];
                    mockAdminContext.isSuperAdmin = false;
                    mockAdminContext.loading = false;

                    const fallbackId = `fallback-${Math.random()}`;
                    const childId = `child-${Math.random()}`;
                    const { container } = render(
                        <PermissionGuard
                            permission={requiredPermission}
                            fallback={<div data-testid={fallbackId}>No Access</div>}
                        >
                            <div data-testid={childId}>Protected</div>
                        </PermissionGuard>
                    );

                    const fallback = container.querySelector(`[data-testid="${fallbackId}"]`);
                    const child = container.querySelector(`[data-testid="${childId}"]`);

                    // Fallback should be rendered, child should not
                    return fallback !== null && child === null;
                }
            ),
            { numRuns: 100 }
        );
    });
});

describe('PermissionButton Property Tests', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        mockAdminContext.loading = false;
        mockAdminContext.isSuperAdmin = false;
        mockAdminContext.permissions = [];

        // Set up mockHasPermission to use the actual permission checking logic
        mockHasPermission.mockImplementation((permission: string | string[], mode: 'any' | 'all' = 'any') => {
            return checkPermission(mockAdminContext.permissions, Array.isArray(permission) ? permission : [permission], mode);
        });
    });

    /**
     * Property: PermissionButton should be disabled (not hidden) when disableOnNoPermission is true
     * Note: By default, PermissionButton hides the button when there's no permission.
     * Set disableOnNoPermission to true to show it as disabled instead.
     * Note: Children must be a button element for the disabled prop to be applied.
     */
    it('Property 6g: PermissionButton should be disabled when disableOnNoPermission is true and user lacks permission', () => {
        fc.assert(
            fc.property(
                permissionCodeArb,
                (requiredPermission) => {
                    mockAdminContext.permissions = [];
                    mockAdminContext.isSuperAdmin = false;
                    mockAdminContext.loading = false;

                    const { container } = render(
                        <PermissionButton
                            permission={requiredPermission}
                            disableOnNoPermission
                        >
                            <button>Click Me</button>
                        </PermissionButton>
                    );

                    const button = container.querySelector('button');
                    return button !== null && button.disabled === true;
                }
            ),
            { numRuns: 100 }
        );
    });

    /**
     * Property: PermissionButton should be hidden when hideOnNoPermission is true
     */
    it('Property 6h: PermissionButton should be hidden when hideOnNoPermission is true and user lacks permission', () => {
        fc.assert(
            fc.property(
                permissionCodeArb,
                (requiredPermission) => {
                    mockAdminContext.permissions = [];
                    mockAdminContext.isSuperAdmin = false;
                    mockAdminContext.loading = false;

                    const { container } = render(
                        <PermissionButton
                            permission={requiredPermission}
                            hideOnNoPermission
                        >
                            Click Me
                        </PermissionButton>
                    );

                    const button = container.querySelector('button');
                    return button === null;
                }
            ),
            { numRuns: 100 }
        );
    });

    /**
     * Property: PermissionButton should be enabled when user has permission
     * Note: Children must be a button element for the disabled prop to be applied.
     */
    it('Property 6i: PermissionButton should be enabled when user has permission', () => {
        fc.assert(
            fc.property(
                permissionCodeArb,
                (requiredPermission) => {
                    mockAdminContext.permissions = [requiredPermission];
                    mockAdminContext.isSuperAdmin = false;
                    mockAdminContext.loading = false;

                    const { container } = render(
                        <PermissionButton permission={requiredPermission}>
                            <button>Click Me</button>
                        </PermissionButton>
                    );

                    const button = container.querySelector('button');
                    return button !== null && button.disabled === false;
                }
            ),
            { numRuns: 100 }
        );
    });
});
