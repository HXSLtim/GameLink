/**
 * Menu Permission Utility Tests
 *
 * Property-based tests for menu permission filtering and route protection
 * **Feature: rbac-button-level-permission, Property 14: 菜单权限过滤**
 * **Feature: rbac-button-level-permission, Property 15: 路由权限保护**
 * **Validates: Requirements 8.1, 8.3**
 */
import { describe, it } from 'vitest';
import * as fc from 'fast-check';
import type { Menu } from '@/api/admin';
import {
    hasPermission,
    filterMenusByPermission,
    hasRoutePermission,
    getFirstAccessiblePath,
} from './menuPermission';

// Arbitrary generators for property testing
const permissionCodeArb = fc.stringMatching(/^[a-z]+\.[a-z]+\.[a-z]+$/);
const permissionListArb = fc.array(permissionCodeArb, { minLength: 0, maxLength: 10 });

// Menu generator
const menuArb: fc.Arbitrary<Menu> = fc.record({
    id: fc.integer({ min: 1, max: 1000 }),
    name: fc.string({ minLength: 1, maxLength: 20 }),
    path: fc.stringMatching(/^\/[a-z]+(?:\/[a-z]+)*$/),
    component: fc.string({ minLength: 1, maxLength: 50 }),
    parentId: fc.option(fc.integer({ min: 1, max: 1000 }), { nil: null }),
    order: fc.integer({ min: 0, max: 100 }),
    visible: fc.boolean(),
    type: fc.constantFrom('menu' as const, 'page' as const, 'button' as const),
    permission: fc.option(permissionCodeArb, { nil: '' }),
    icon: fc.option(fc.string({ minLength: 1, maxLength: 20 }), { nil: undefined }),
    redirect: fc.option(fc.string({ minLength: 1, maxLength: 50 }), { nil: undefined }),
    description: fc.option(fc.string({ minLength: 1, maxLength: 100 }), { nil: undefined }),
    children: fc.constant(undefined),
});

// Menu with children generator
const menuWithChildrenArb: fc.Arbitrary<Menu> = fc.record({
    id: fc.integer({ min: 1, max: 1000 }),
    name: fc.string({ minLength: 1, maxLength: 20 }),
    path: fc.stringMatching(/^\/[a-z]+$/),
    component: fc.string({ minLength: 1, maxLength: 50 }),
    parentId: fc.constant(null),
    order: fc.integer({ min: 0, max: 100 }),
    visible: fc.constant(true),
    type: fc.constant('menu' as const),
    permission: fc.constant(''),
    icon: fc.option(fc.string({ minLength: 1, maxLength: 20 }), { nil: undefined }),
    redirect: fc.option(fc.string({ minLength: 1, maxLength: 50 }), { nil: undefined }),
    description: fc.option(fc.string({ minLength: 1, maxLength: 100 }), { nil: undefined }),
    children: fc.array(menuArb, { minLength: 1, maxLength: 5 }),
});

describe('hasPermission', () => {
    /**
     * Property: Super admin (*) should always have permission
     */
    it('should return true for super admin regardless of required permission', () => {
        fc.assert(
            fc.property(permissionCodeArb, (permission) => {
                const userPermissions = ['*'];
                return hasPermission(userPermissions, permission) === true;
            }),
            { numRuns: 100 }
        );
    });

    /**
     * Property: Empty permission requirement should always return true
     */
    it('should return true when permission is empty', () => {
        fc.assert(
            fc.property(permissionListArb, (userPermissions) => {
                return hasPermission(userPermissions, '') === true;
            }),
            { numRuns: 100 }
        );
    });

    /**
     * Property: User with the exact permission should have access
     */
    it('should return true when user has the exact permission', () => {
        fc.assert(
            fc.property(permissionCodeArb, permissionListArb, (permission, extraPermissions) => {
                const userPermissions = [permission, ...extraPermissions];
                return hasPermission(userPermissions, permission) === true;
            }),
            { numRuns: 100 }
        );
    });

    /**
     * Property: User without the permission should not have access
     */
    it('should return false when user lacks the required permission', () => {
        fc.assert(
            fc.property(
                permissionCodeArb,
                fc.array(permissionCodeArb, { minLength: 0, maxLength: 5 }),
                (requiredPermission, userPermissions) => {
                    // Ensure user doesn't have the required permission or super admin
                    const filteredPermissions = userPermissions.filter(
                        p => p !== requiredPermission && p !== '*'
                    );
                    return hasPermission(filteredPermissions, requiredPermission) === false;
                }
            ),
            { numRuns: 100 }
        );
    });
});

describe('filterMenusByPermission', () => {
    /**
     * **Feature: rbac-button-level-permission, Property 14: 菜单权限过滤**
     * **Validates: Requirements 8.1, 8.2**
     * 
     * Property: Super admin should see all visible menus (bypasses permission check)
     * Note: Hidden menus (visible=false) are still filtered out for everyone,
     * including super admin. Permission bypass only affects permission checks.
     */
    it('should return all menus for super admin (permission bypass)', () => {
        fc.assert(
            fc.property(
                fc.array(menuArb, { minLength: 1, maxLength: 10 }),
                (menus) => {
                    const userPermissions = ['*'];
                    const filteredMenus = filterMenusByPermission(menus, userPermissions);
                    
                    // Super admin bypasses permission check, but hidden menus are still filtered
                    // Count visible menus (visible !== false)
                    const visibleMenus = menus.filter(m => m.visible !== false);
                    return filteredMenus.length === visibleMenus.length;
                }
            ),
            { numRuns: 100 }
        );
    });

    /**
     * **Feature: rbac-button-level-permission, Property 14: 菜单权限过滤**
     * **Validates: Requirements 8.1**
     * 
     * Property: Filtered menus should only contain items user has permission for
     */
    it('should only include menus user has permission for', () => {
        fc.assert(
            fc.property(
                fc.array(menuArb, { minLength: 1, maxLength: 10 }),
                permissionListArb,
                (menus, userPermissions) => {
                    // Ensure no super admin
                    const filteredUserPermissions = userPermissions.filter(p => p !== '*');
                    const filteredMenus = filterMenusByPermission(menus, filteredUserPermissions);
                    
                    // Every filtered menu should either have no permission requirement
                    // or user should have the required permission
                    return filteredMenus.every(menu => {
                        if (!menu.permission) return true;
                        return filteredUserPermissions.includes(menu.permission);
                    });
                }
            ),
            { numRuns: 100 }
        );
    });

    /**
     * **Feature: rbac-button-level-permission, Property 14: 菜单权限过滤**
     * **Validates: Requirements 8.1**
     * 
     * Property: Hidden menus (visible=false) should never appear in filtered results
     */
    it('should exclude hidden menus regardless of permissions', () => {
        fc.assert(
            fc.property(
                fc.array(menuArb, { minLength: 1, maxLength: 10 }),
                permissionListArb,
                (menus, userPermissions) => {
                    const filteredMenus = filterMenusByPermission(menus, userPermissions);
                    
                    // No hidden menus should appear
                    return filteredMenus.every(menu => menu.visible !== false);
                }
            ),
            { numRuns: 100 }
        );
    });

    /**
     * **Feature: rbac-button-level-permission, Property 14: 菜单权限过滤**
     * **Validates: Requirements 8.2**
     * 
     * Property: Parent menu should be hidden if no children are accessible
     */
    it('should hide parent menu when no children are accessible', () => {
        fc.assert(
            fc.property(
                menuWithChildrenArb,
                (parentMenu) => {
                    // Create a parent with children that all require permissions user doesn't have
                    const childrenWithPermissions = parentMenu.children?.map((child, i) => ({
                        ...child,
                        visible: true,
                        permission: `unique.permission.${i}`,
                    })) || [];
                    
                    const menuWithRestrictedChildren: Menu = {
                        ...parentMenu,
                        children: childrenWithPermissions,
                    };
                    
                    // User has no permissions
                    const userPermissions: string[] = [];
                    const filteredMenus = filterMenusByPermission([menuWithRestrictedChildren], userPermissions);
                    
                    // Parent should be hidden because no children are accessible
                    return filteredMenus.length === 0;
                }
            ),
            { numRuns: 100 }
        );
    });

    /**
     * **Feature: rbac-button-level-permission, Property 14: 菜单权限过滤**
     * **Validates: Requirements 8.1**
     * 
     * Property: Menus without permission requirement should be accessible
     */
    it('should include menus without permission requirement', () => {
        fc.assert(
            fc.property(
                fc.array(
                    fc.record({
                        id: fc.integer({ min: 1, max: 1000 }),
                        name: fc.string({ minLength: 1, maxLength: 20 }),
                        path: fc.stringMatching(/^\/[a-z]+$/),
                        component: fc.string({ minLength: 1, maxLength: 50 }),
                        parentId: fc.constant(null),
                        order: fc.integer({ min: 0, max: 100 }),
                        visible: fc.constant(true),
                        type: fc.constant('menu' as const),
                        permission: fc.constant(''), // No permission required
                        children: fc.constant(undefined),
                    }),
                    { minLength: 1, maxLength: 5 }
                ),
                (menus) => {
                    // User has no permissions
                    const userPermissions: string[] = [];
                    const filteredMenus = filterMenusByPermission(menus, userPermissions);
                    
                    // All menus without permission should be accessible
                    return filteredMenus.length === menus.length;
                }
            ),
            { numRuns: 100 }
        );
    });
});

describe('hasRoutePermission', () => {
    /**
     * **Feature: rbac-button-level-permission, Property 15: 路由权限保护**
     * **Validates: Requirements 8.3**
     * 
     * Property: Super admin should have access to all routes
     */
    it('should return true for super admin regardless of route', () => {
        fc.assert(
            fc.property(
                fc.array(menuArb, { minLength: 1, maxLength: 10 }),
                fc.stringMatching(/^\/[a-z]+(?:\/[a-z]+)*$/),
                (menus, path) => {
                    const userPermissions = ['*'];
                    return hasRoutePermission(menus, path, userPermissions) === true;
                }
            ),
            { numRuns: 100 }
        );
    });

    /**
     * **Feature: rbac-button-level-permission, Property 15: 路由权限保护**
     * **Validates: Requirements 8.3**
     * 
     * Property: Routes not in menu should be accessible (public pages)
     */
    it('should return true for routes not in menu', () => {
        fc.assert(
            fc.property(
                fc.array(menuArb, { minLength: 0, maxLength: 5 }),
                permissionListArb,
                (menus, userPermissions) => {
                    // Use a path that definitely won't be in the menus
                    const uniquePath = '/unique/path/not/in/menu';
                    return hasRoutePermission(menus, uniquePath, userPermissions) === true;
                }
            ),
            { numRuns: 100 }
        );
    });

    /**
     * **Feature: rbac-button-level-permission, Property 15: 路由权限保护**
     * **Validates: Requirements 8.3**
     * 
     * Property: User with permission should access protected route
     */
    it('should return true when user has permission for the route', () => {
        fc.assert(
            fc.property(
                permissionCodeArb,
                fc.stringMatching(/^\/[a-z]+$/),
                (permission, path) => {
                    const menu: Menu = {
                        id: 1,
                        name: 'Test Menu',
                        path: path,
                        component: 'TestComponent',
                        parentId: null,
                        order: 0,
                        visible: true,
                        type: 'menu',
                        permission: permission,
                    };
                    
                    const userPermissions = [permission];
                    return hasRoutePermission([menu], path, userPermissions) === true;
                }
            ),
            { numRuns: 100 }
        );
    });

    /**
     * **Feature: rbac-button-level-permission, Property 15: 路由权限保护**
     * **Validates: Requirements 8.3**
     * 
     * Property: User without permission should not access protected route
     */
    it('should return false when user lacks permission for the route', () => {
        fc.assert(
            fc.property(
                permissionCodeArb,
                fc.stringMatching(/^\/[a-z]+$/),
                fc.array(permissionCodeArb, { minLength: 0, maxLength: 5 }),
                (requiredPermission, path, userPermissions) => {
                    const menu: Menu = {
                        id: 1,
                        name: 'Test Menu',
                        path: path,
                        component: 'TestComponent',
                        parentId: null,
                        order: 0,
                        visible: true,
                        type: 'menu',
                        permission: requiredPermission,
                    };
                    
                    // Ensure user doesn't have the required permission or super admin
                    const filteredPermissions = userPermissions.filter(
                        p => p !== requiredPermission && p !== '*'
                    );
                    
                    return hasRoutePermission([menu], path, filteredPermissions) === false;
                }
            ),
            { numRuns: 100 }
        );
    });

    /**
     * **Feature: rbac-button-level-permission, Property 15: 路由权限保护**
     * **Validates: Requirements 8.3**
     * 
     * Property: Routes without permission requirement should be accessible
     */
    it('should return true for routes without permission requirement', () => {
        fc.assert(
            fc.property(
                fc.stringMatching(/^\/[a-z]+$/),
                (path) => {
                    const menu: Menu = {
                        id: 1,
                        name: 'Public Menu',
                        path: path,
                        component: 'PublicComponent',
                        parentId: null,
                        order: 0,
                        visible: true,
                        type: 'menu',
                        permission: '', // No permission required
                    };
                    
                    // Even with no permissions, should have access
                    return hasRoutePermission([menu], path, []) === true;
                }
            ),
            { numRuns: 100 }
        );
    });
});

describe('getFirstAccessiblePath', () => {
    /**
     * Property: Should return null when no menus are accessible
     */
    it('should return null when no menus are accessible', () => {
        fc.assert(
            fc.property(
                fc.array(
                    fc.record({
                        id: fc.integer({ min: 1, max: 1000 }),
                        name: fc.string({ minLength: 1, maxLength: 20 }),
                        path: fc.stringMatching(/^\/[a-z]+$/),
                        component: fc.string({ minLength: 1, maxLength: 50 }),
                        parentId: fc.constant(null),
                        order: fc.integer({ min: 0, max: 100 }),
                        visible: fc.constant(true),
                        type: fc.constant('menu' as const),
                        permission: permissionCodeArb, // All require permission
                        children: fc.constant(undefined),
                    }),
                    { minLength: 1, maxLength: 5 }
                ),
                (menus) => {
                    // User has no permissions
                    const userPermissions: string[] = [];
                    const firstPath = getFirstAccessiblePath(menus, userPermissions);
                    
                    return firstPath === null;
                }
            ),
            { numRuns: 100 }
        );
    });

    /**
     * Property: Should return a valid path when menus are accessible
     */
    it('should return a valid path when menus are accessible', () => {
        fc.assert(
            fc.property(
                fc.array(
                    fc.record({
                        id: fc.integer({ min: 1, max: 1000 }),
                        name: fc.string({ minLength: 1, maxLength: 20 }),
                        path: fc.stringMatching(/^\/[a-z]+$/),
                        component: fc.string({ minLength: 1, maxLength: 50 }),
                        parentId: fc.constant(null),
                        order: fc.integer({ min: 0, max: 100 }),
                        visible: fc.constant(true),
                        type: fc.constant('menu' as const),
                        permission: fc.constant(''), // No permission required
                        children: fc.constant(undefined),
                    }),
                    { minLength: 1, maxLength: 5 }
                ),
                (menus) => {
                    const userPermissions: string[] = [];
                    const firstPath = getFirstAccessiblePath(menus, userPermissions);
                    
                    // Should return one of the menu paths
                    return firstPath !== null && menus.some(m => m.path === firstPath);
                }
            ),
            { numRuns: 100 }
        );
    });

    /**
     * Property: Super admin should always get a path if menus exist
     */
    it('should return a path for super admin when visible menus exist', () => {
        fc.assert(
            fc.property(
                fc.array(
                    fc.record({
                        id: fc.integer({ min: 1, max: 1000 }),
                        name: fc.string({ minLength: 1, maxLength: 20 }),
                        path: fc.stringMatching(/^\/[a-z]+$/),
                        component: fc.string({ minLength: 1, maxLength: 50 }),
                        parentId: fc.constant(null),
                        order: fc.integer({ min: 0, max: 100 }),
                        visible: fc.constant(true),
                        type: fc.constant('menu' as const),
                        permission: permissionCodeArb,
                        children: fc.constant(undefined),
                    }),
                    { minLength: 1, maxLength: 5 }
                ),
                (menus) => {
                    const userPermissions = ['*'];
                    const firstPath = getFirstAccessiblePath(menus, userPermissions);
                    
                    // Super admin should get a path
                    return firstPath !== null;
                }
            ),
            { numRuns: 100 }
        );
    });
});
