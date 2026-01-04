/**
 * Permission Utils Tests
 *
 * Coverage Target: 90%+
 *
 * Test Scenarios:
 * 1. PermissionStore class (set, get, clear, subscribe)
 * 2. hasPermission function (single, array, any/all modes)
 * 3. hasAllPermissions function
 * 4. hasAnyPermission function
 * 5. filterActionsByPermission function
 * 6. PermissionParser utility (getModule, getResource, getAction, build, matches)
 * 7. Super admin (*) wildcard handling
 * 8. Edge cases (empty permissions, missing permissions)
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import * as fc from 'fast-check';
import {
    permissionStore,
    hasPermission,
    hasAllPermissions,
    hasAnyPermission,
    filterActionsByPermission,
    PermissionParser,
} from './permission';

describe('permission utils', () => {
    beforeEach(() => {
        // Clear permissions before each test
        permissionStore.clearPermissions();
    });

    describe('PermissionStore', () => {
        describe('setPermissions and getPermissions', () => {
            it('should set and get permissions', () => {
                const permissions = ['admin.users.read', 'admin.users.write'];
                permissionStore.setPermissions(permissions);

                expect(permissionStore.getPermissions()).toEqual(permissions);
            });

            it('should handle empty permission list', () => {
                permissionStore.setPermissions([]);

                expect(permissionStore.getPermissions()).toEqual([]);
            });

            it('should handle duplicate permissions', () => {
                const permissions = ['admin.users.read', 'admin.users.read'];
                permissionStore.setPermissions(permissions);

                expect(permissionStore.getPermissions()).toEqual(permissions);
            });
        });

        describe('clearPermissions', () => {
            it('should clear all permissions', () => {
                permissionStore.setPermissions(['admin.users.read']);
                permissionStore.clearPermissions();

                expect(permissionStore.getPermissions()).toEqual([]);
            });

            it('should handle clearing empty permissions', () => {
                permissionStore.clearPermissions();

                expect(permissionStore.getPermissions()).toEqual([]);
            });
        });

        describe('subscribe', () => {
            it('should notify subscribers on permission change', () => {
                const listener = vi.fn();
                permissionStore.subscribe(listener);

                permissionStore.setPermissions(['admin.users.read']);

                expect(listener).toHaveBeenCalledWith(['admin.users.read']);
            });

            it('should notify multiple subscribers', () => {
                const listener1 = vi.fn();
                const listener2 = vi.fn();
                permissionStore.subscribe(listener1);
                permissionStore.subscribe(listener2);

                permissionStore.setPermissions(['admin.users.read']);

                expect(listener1).toHaveBeenCalledWith(['admin.users.read']);
                expect(listener2).toHaveBeenCalledWith(['admin.users.read']);
            });

            it('should unsubscribe correctly', () => {
                const listener = vi.fn();
                const unsubscribe = permissionStore.subscribe(listener);

                unsubscribe();
                permissionStore.setPermissions(['admin.users.read']);

                expect(listener).not.toHaveBeenCalled();
            });

            it('should notify on clearPermissions', () => {
                const listener = vi.fn();
                permissionStore.setPermissions(['admin.users.read']);
                permissionStore.subscribe(listener);

                permissionStore.clearPermissions();

                expect(listener).toHaveBeenCalledWith([]);
            });
        });

        describe('hasPermission', () => {
            it('should return false for empty permission list', () => {
                expect(permissionStore.hasPermission('admin.users.read')).toBe(false);
            });

            it('should return true for matching permission', () => {
                permissionStore.setPermissions(['admin.users.read']);

                expect(permissionStore.hasPermission('admin.users.read')).toBe(true);
            });

            it('should return false for non-matching permission', () => {
                permissionStore.setPermissions(['admin.users.read']);

                expect(permissionStore.hasPermission('admin.users.write')).toBe(false);
            });

            it('should return true for super admin wildcard', () => {
                permissionStore.setPermissions(['*']);

                expect(permissionStore.hasPermission('admin.users.read')).toBe(true);
                expect(permissionStore.hasPermission('any.permission')).toBe(true);
            });

            describe('array permissions with mode', () => {
                beforeEach(() => {
                    permissionStore.setPermissions(['admin.users.read', 'admin.users.write']);
                });

                it('should return true in "any" mode when at least one matches', () => {
                    expect(
                        permissionStore.hasPermission(['admin.users.read', 'admin.users.delete'], 'any')
                    ).toBe(true);
                });

                it('should return false in "any" mode when none match', () => {
                    expect(
                        permissionStore.hasPermission(['admin.users.delete', 'admin.users.update'], 'any')
                    ).toBe(false);
                });

                it('should return true in "all" mode when all match', () => {
                    expect(
                        permissionStore.hasPermission(['admin.users.read', 'admin.users.write'], 'all')
                    ).toBe(true);
                });

                it('should return false in "all" mode when one does not match', () => {
                    expect(
                        permissionStore.hasPermission(['admin.users.read', 'admin.users.delete'], 'all')
                    ).toBe(false);
                });

                it('should default to "any" mode', () => {
                    expect(
                        permissionStore.hasPermission(['admin.users.read', 'admin.users.delete'])
                    ).toBe(true);
                });
            });
        });
    });

    describe('hasPermission function', () => {
        it('should use permissionStore for checking', () => {
            permissionStore.setPermissions(['admin.users.read']);

            expect(hasPermission('admin.users.read')).toBe(true);
            expect(hasPermission('admin.users.write')).toBe(false);
        });

        it('should pass through array permissions', () => {
            permissionStore.setPermissions(['admin.users.read']);

            expect(hasPermission(['admin.users.read', 'admin.users.write'], 'any')).toBe(true);
            expect(hasPermission(['admin.users.read', 'admin.users.write'], 'all')).toBe(false);
        });

        it('should pass through wildcard', () => {
            permissionStore.setPermissions(['*']);

            expect(hasPermission('any.permission')).toBe(true);
        });
    });

    describe('hasAllPermissions function', () => {
        it('should return true when user has all permissions', () => {
            permissionStore.setPermissions(['admin.users.read', 'admin.users.write']);

            expect(hasAllPermissions(['admin.users.read', 'admin.users.write'])).toBe(true);
        });

        it('should return false when user is missing any permission', () => {
            permissionStore.setPermissions(['admin.users.read']);

            expect(hasAllPermissions(['admin.users.read', 'admin.users.write'])).toBe(false);
        });

        it('should return true for super admin', () => {
            permissionStore.setPermissions(['*']);

            expect(hasAllPermissions(['admin.users.read', 'admin.users.write'])).toBe(true);
        });

        it('should return true for empty array', () => {
            expect(hasAllPermissions([])).toBe(true);
        });
    });

    describe('hasAnyPermission function', () => {
        it('should return true when user has at least one permission', () => {
            permissionStore.setPermissions(['admin.users.read']);

            expect(hasAnyPermission(['admin.users.read', 'admin.users.write'])).toBe(true);
        });

        it('should return false when user has none of the permissions', () => {
            permissionStore.setPermissions(['admin.orders.read']);

            expect(hasAnyPermission(['admin.users.read', 'admin.users.write'])).toBe(false);
        });

        it('should return true for super admin', () => {
            permissionStore.setPermissions(['*']);

            expect(hasAnyPermission(['admin.users.read', 'admin.users.write'])).toBe(true);
        });

        it('should return false for empty array', () => {
            expect(hasAnyPermission([])).toBe(false);
        });
    });

    describe('filterActionsByPermission function', () => {
        interface Action {
            key: string;
            label: string;
            permission?: string | string[];
        }

        it('should filter actions based on permissions', () => {
            permissionStore.setPermissions(['admin.users.read', 'admin.users.write']);

            const actions: Action[] = [
                { key: 'read', label: 'View', permission: 'admin.users.read' },
                { key: 'write', label: 'Edit', permission: 'admin.users.write' },
                { key: 'delete', label: 'Delete', permission: 'admin.users.delete' },
            ];

            const filtered = filterActionsByPermission(actions);

            expect(filtered).toHaveLength(2);
            expect(filtered.map(a => a.key)).toEqual(['read', 'write']);
        });

        it('should include actions without permission requirement', () => {
            permissionStore.setPermissions(['admin.users.read']);

            const actions: Action[] = [
                { key: 'read', label: 'View', permission: 'admin.users.read' },
                { key: 'export', label: 'Export' }, // No permission
                { key: 'delete', label: 'Delete', permission: 'admin.users.delete' },
            ];

            const filtered = filterActionsByPermission(actions);

            expect(filtered).toHaveLength(2);
            expect(filtered.map(a => a.key)).toEqual(['read', 'export']);
        });

        it('should handle array permissions in actions', () => {
            permissionStore.setPermissions(['admin.users.read']);

            const actions: Action[] = [
                { key: 'read', label: 'View', permission: ['admin.users.read', 'admin.users.write'] },
                { key: 'delete', label: 'Delete', permission: ['admin.users.delete', 'admin.users.update'] },
            ];

            const filtered = filterActionsByPermission(actions);

            expect(filtered).toHaveLength(1);
            expect(filtered[0].key).toBe('read');
        });

        it('should use custom permission key', () => {
            permissionStore.setPermissions(['read']);

            const actions = [
                { key: 'read', label: 'View', perm: 'read' },
                { key: 'write', label: 'Edit', perm: 'write' },
            ] as unknown as Action[];

            const filtered = filterActionsByPermission(actions, 'perm' as keyof Action);

            expect(filtered).toHaveLength(1);
            expect(filtered[0].key).toBe('read');
        });

        it('should return all actions for super admin', () => {
            permissionStore.setPermissions(['*']);

            const actions: Action[] = [
                { key: 'read', label: 'View', permission: 'admin.users.read' },
                { key: 'write', label: 'Edit', permission: 'admin.users.write' },
                { key: 'delete', label: 'Delete', permission: 'admin.users.delete' },
            ];

            const filtered = filterActionsByPermission(actions);

            expect(filtered).toHaveLength(3);
        });
    });

    describe('PermissionParser', () => {
        describe('getModule', () => {
            it('should extract module from permission code', () => {
                expect(PermissionParser.getModule('admin.users.read')).toBe('admin');
                expect(PermissionParser.getModule('player.games.create')).toBe('player');
            });

            it('should handle incomplete codes', () => {
                expect(PermissionParser.getModule('admin')).toBe('admin');
                expect(PermissionParser.getModule('')).toBe('');
            });
        });

        describe('getResource', () => {
            it('should extract resource from permission code', () => {
                expect(PermissionParser.getResource('admin.users.read')).toBe('users');
                expect(PermissionParser.getResource('player.games.create')).toBe('games');
            });

            it('should handle incomplete codes', () => {
                expect(PermissionParser.getResource('admin.users')).toBe('users');
                expect(PermissionParser.getResource('admin')).toBe('');
            });
        });

        describe('getAction', () => {
            it('should extract action from permission code', () => {
                expect(PermissionParser.getAction('admin.users.read')).toBe('read');
                expect(PermissionParser.getAction('player.games.create')).toBe('create');
            });

            it('should handle incomplete codes', () => {
                expect(PermissionParser.getAction('admin.users.read.extra')).toBe('read');
                expect(PermissionParser.getAction('admin.users')).toBe('');
            });
        });

        describe('build', () => {
            it('should build permission code from parts', () => {
                expect(PermissionParser.build('admin', 'users', 'read')).toBe('admin.users.read');
                expect(PermissionParser.build('player', 'games', 'create')).toBe('player.games.create');
            });

            it('should handle empty parts', () => {
                expect(PermissionParser.build('', '', '')).toBe('..');
            });
        });

        describe('matches', () => {
            it('should match exact patterns', () => {
                expect(PermissionParser.matches('admin.users.read', 'admin.users.read')).toBe(true);
                expect(PermissionParser.matches('admin.users.read', 'admin.users.write')).toBe(false);
            });

            it('should match wildcard patterns', () => {
                expect(PermissionParser.matches('admin.users.read', 'admin.users.*')).toBe(true);
                expect(PermissionParser.matches('admin.users.write', 'admin.users.*')).toBe(true);
                expect(PermissionParser.matches('admin.users.read', 'admin.*.*')).toBe(true);
                expect(PermissionParser.matches('player.users.read', 'admin.*.*')).toBe(false);
            });

            it('should handle partial wildcards', () => {
                expect(PermissionParser.matches('admin.users.read', 'admin.users.re*')).toBe(false);
                expect(PermissionParser.matches('admin.users.read', 'admin.*ead')).toBe(false);
            });

            it('should require matching number of parts', () => {
                expect(PermissionParser.matches('admin.users', 'admin.users.read')).toBe(false);
                expect(PermissionParser.matches('admin.users.read', 'admin.users')).toBe(false);
            });
        });
    });

    describe('property-based tests', () => {
        describe('hasPermission properties', () => {
            /**
             * Property: Empty permission list should always fail for non-empty requirements
             */
            it('should return false for empty permissions with non-empty requirements', () => {
                fc.assert(
                    fc.property(
                        fc.array(fc.stringMatching(/^[a-z]+\.[a-z]+\.[a-z]+$/), { minLength: 1 }),
                        (requiredPermissions) => {
                            permissionStore.setPermissions([]);
                            return permissionStore.hasPermission(requiredPermissions, 'any') === false &&
                                   permissionStore.hasPermission(requiredPermissions, 'all') === false;
                        }
                    ),
                    { numRuns: 20 }
                );
            });

            /**
             * Property: Super admin wildcard should match any permission
             */
            it('should return true for super admin with any requirement', () => {
                fc.assert(
                    fc.property(
                        fc.array(fc.stringMatching(/^[a-z]+\.[a-z]+\.[a-z]+$/)),
                        fc.constantFrom('any' as const, 'all' as const),
                        (requiredPermissions, mode) => {
                            permissionStore.setPermissions(['*']);
                            return permissionStore.hasPermission(requiredPermissions, mode) === true;
                        }
                    ),
                    { numRuns: 20 }
                );
            });

            /**
             * Property: "any" mode should pass if at least one permission matches
             */
            it('should pass in any mode with at least one matching permission', () => {
                fc.assert(
                    fc.property(
                        fc.array(fc.stringMatching(/^[a-z]+\.[a-z]+\.[a-z]+$/), { minLength: 2, maxLength: 5 }),
                        fc.integer({ min: 0 }),
                        (permissions, selectIndex) => {
                            const actualIndex = selectIndex % permissions.length;
                            permissionStore.setPermissions([permissions[actualIndex]]);
                            return permissionStore.hasPermission(permissions, 'any') === true;
                        }
                    ),
                    { numRuns: 20 }
                );
            });

            /**
             * Property: "all" mode should fail if any permission is missing
             */
            it('should fail in all mode when missing any permission', () => {
                fc.assert(
                    fc.property(
                        fc.array(fc.stringMatching(/^[a-z]+\.[a-z]+\.[a-z]+$/), { minLength: 2, maxLength: 5 }),
                        fc.integer({ min: 0 }),
                        (permissions, removeIndex) => {
                            const actualIndex = removeIndex % permissions.length;
                            permissionStore.setPermissions(permissions.filter((_, i) => i !== actualIndex));
                            return permissionStore.hasPermission(permissions, 'all') === false;
                        }
                    ),
                    { numRuns: 20 }
                );
            });
        });

        describe('PermissionParser properties', () => {
            /**
             * Property: build(getModule(x), getResource(x), getAction(x)) == x
             */
            it('should be reversible for complete permission codes', () => {
                fc.assert(
                    fc.property(
                        fc.tuple(
                            fc.stringMatching(/^[a-z]+$/),
                            fc.stringMatching(/^[a-z]+$/),
                            fc.stringMatching(/^[a-z]+$/)
                        ),
                        ([module, resource, action]) => {
                            const code = PermissionParser.build(module, resource, action);
                            return PermissionParser.getModule(code) === module &&
                                   PermissionParser.getResource(code) === resource &&
                                   PermissionParser.getAction(code) === action;
                        }
                    ),
                    { numRuns: 20 }
                );
            });

            /**
             * Property: matches(x, x) should always be true
             */
            it('should match exact same code', () => {
                fc.assert(
                    fc.property(fc.stringMatching(/^[a-z]+\.[a-z]+\.[a-z]+$/), (code) => {
                        return PermissionParser.matches(code, code) === true;
                    }),
                    { numRuns: 20 }
                );
            });

            /**
             * Property: matches(x, '*.*.*') should always be true for 3-part codes
             */
            it('should match wildcard pattern', () => {
                fc.assert(
                    fc.property(fc.stringMatching(/^[a-z]+\.[a-z]+\.[a-z]+$/), (code) => {
                        return PermissionParser.matches(code, '*.*.*') === true;
                    }),
                    { numRuns: 20 }
                );
            });
        });

        describe('filterActionsByPermission properties', () => {
            /**
             * Property: Filtered actions should always be subset of original
             */
            it('should only return actions from original list', () => {
                fc.assert(
                    fc.property(
                        fc.array(fc.record({
                            key: fc.string(),
                            permission: fc.option(fc.stringMatching(/^[a-z]+\.[a-z]+\.[a-z]+$/)),
                        })),
                        fc.array(fc.stringMatching(/^[a-z]+\.[a-z]+\.[a-z]+$/)),
                        (actions, userPermissions) => {
                            permissionStore.setPermissions(userPermissions);
                            const filtered = filterActionsByPermission(actions);

                            // Every filtered action should be in the original list
                            return filtered.every(action => actions.includes(action));
                        }
                    ),
                    { numRuns: 20 }
                );
            });

            /**
             * Property: Super admin should get all actions
             */
            it('should return all actions for super admin', () => {
                fc.assert(
                    fc.property(
                        fc.array(fc.record({
                            key: fc.string(),
                            permission: fc.option(fc.string()),
                        })),
                        (actions) => {
                            permissionStore.setPermissions(['*']);
                            const filtered = filterActionsByPermission(actions);
                            return filtered.length === actions.length;
                        }
                    ),
                    { numRuns: 20 }
                );
            });
        });
    });

    describe('edge cases', () => {
        it('should handle empty strings in permission checks', () => {
            permissionStore.setPermissions(['admin.users.read']);

            expect(hasPermission('')).toBe(false);
            expect(hasPermission(['', 'admin.users.read'], 'any')).toBe(true);
        });

        it('should handle very long permission codes', () => {
            const longCode = 'a'.repeat(1000);
            permissionStore.setPermissions([longCode]);

            expect(hasPermission(longCode)).toBe(true);
        });

        it('should handle special characters in permission codes', () => {
            const permissions = ['admin.users.read', 'admin.users.write', 'admin.users.delete'];
            permissionStore.setPermissions(permissions);

            expect(hasPermission('admin.users.read')).toBe(true);
            expect(hasPermission('admin.users.write')).toBe(true);
            expect(hasPermission('admin.users.delete')).toBe(true);
        });

        it('should handle undefined/null in filterActionsByPermission', () => {
            permissionStore.setPermissions(['admin.users.read']);

            const actions = [
                { key: 'read', label: 'View', permission: undefined },
                { key: 'export', label: 'Export', permission: null },
            ] as unknown as { key: string; label: string; permission: string | null | undefined }[];

            const filtered = filterActionsByPermission(actions);

            expect(filtered).toHaveLength(2);
        });
    });

    describe('subscription edge cases', () => {
        it('should handle duplicate subscriptions', () => {
            const listener = vi.fn();
            permissionStore.subscribe(listener);
            permissionStore.subscribe(listener);

            permissionStore.setPermissions(['admin.users.read']);

            expect(listener).toHaveBeenCalledTimes(2);
        });

        it('should handle unsubscribe during notification', () => {
            let unsubscribe: () => void;
            const listener1 = vi.fn(() => unsubscribe());
            const listener2 = vi.fn();

            unsubscribe = permissionStore.subscribe(listener1);
            permissionStore.subscribe(listener2);

            permissionStore.setPermissions(['admin.users.read']);

            expect(listener1).toHaveBeenCalledTimes(1);
            expect(listener2).toHaveBeenCalledTimes(1);
        });
    });
});
