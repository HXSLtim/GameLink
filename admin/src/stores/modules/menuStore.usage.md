# Menu Store Usage Examples

## Overview

The `menuStore` provides complete menu management with permission filtering, caching, and state management using Zustand.

## Features

- **Permission-based filtering**: Automatically filters menus based on user permissions
- **Menu caching**: Persists raw menus and collapsed state to localStorage
- **Nested menu support**: Handles multi-level menu hierarchies
- **Super admin support**: Users with `*` permission see all menus
- **Auto-refresh**: Automatically re-filters menus when permissions change

## Basic Usage

### 1. Fetch and Display Menus

```tsx
import { useMenuStore } from '@/stores';

function MenuSidebar() {
  const { menus, loading, fetchMenus } = useMenuStore();

  useEffect(() => {
    fetchMenus();
  }, []);

  if (loading) return <Spin />;

  return (
    <Menu>
      {menus.map(menu => (
        <Menu.Item key={menu.key}>
          {menu.label}
        </Menu.Item>
      ))}
    </Menu>
  );
}
```

### 2. Menu Collapse State

```tsx
function AppLayout() {
  const { collapsed, toggleCollapsed } = useMenuStore();

  return (
    <Layout>
      <Layout.Sider
        collapsible
        collapsed={collapsed}
        onCollapse={toggleCollapsed}
      >
        <MenuSidebar />
      </Layout.Sider>
      <Layout.Content>{children}</Layout.Content>
    </Layout>
  );
}
```

### 3. Get Menu by Key

```tsx
function Breadcrumb() {
  const { getMenuByKey } = useMenuStore();
  const menu = getMenuByKey('dashboard');

  return (
    <Breadcrumb.Item>{menu?.label}</Breadcrumb.Item>
  );
}
```

### 4. Get Menu Paths

```tsx
// Generate routes from menu configuration
function generateRoutes() {
  const { getMenuPaths } = useMenuStore();
  const paths = getMenuPaths();

  return Object.entries(paths).map(([key, path]) => (
    <Route key={key} path={path} element={<Page />} />
  ));
}
```

### 5. Check Menu Permissions

```tsx
function MenuGuard({ children, menuKey }) {
  const { menus, getMenuByKey } = useMenuStore();
  const menu = getMenuByKey(menuKey);

  // Menu is automatically filtered by permissions,
  // so if it exists, the user has access
  if (!menu) {
    return <NoPermission />;
  }

  return children;
}
```

## State Properties

| Property | Type | Description |
|----------|------|-------------|
| `rawMenus` | `MenuItem[]` | Original menus from API (unfiltered) |
| `menus` | `MenuItem[]` | Filtered menus based on permissions |
| `loading` | `boolean` | Loading state during fetch |
| `collapsed` | `boolean` | Sidebar collapse state |
| `openKeys` | `string[]` | Currently expanded menu keys |
| `selectedKeys` | `string[]` | Currently selected menu keys |

## Actions

| Action | Parameters | Description |
|--------|-----------|-------------|
| `fetchMenus()` | - | Fetch menus from API and filter by permissions |
| `setCollapsed(bool)` | `boolean` | Set sidebar collapsed state |
| `toggleCollapsed()` | - | Toggle sidebar collapsed state |
| `setOpenKeys(keys)` | `string[]` | Set expanded menu keys |
| `setSelectedKeys(keys)` | `string[]` | Set selected menu keys |

## Selectors

| Selector | Parameters | Returns | Description |
|----------|-----------|---------|-------------|
| `getMenuByKey(key)` | `string` | `MenuItem \| undefined` | Find menu by key (supports nested) |
| `getBreadcrumbs(key)` | `string` | `MenuItem[]` | Get breadcrumb items for menu |
| `getMenuPaths()` | - | `Record<string, string>` | Get all menu paths for routing |
| `filterMenusByPermission(menus, perms)` | `MenuItem[], string[]` | `MenuItem[]` | Filter menus by permissions |

## Permission Filtering

The store automatically filters menus based on user permissions:

- **Super admin** (`*` permission): Sees all menus
- **Regular users**: Only sees menus they have permission for
- **Nested menus**: Parent menus without visible children are hidden
- **Auto-refresh**: Menus re-filter when permissions change

## Integration with AdminContext

The menuStore works alongside the existing AdminContext for a smooth migration:

```tsx
// Old way (AdminContext)
const { menus, refreshMenus } = useAdmin();

// New way (menuStore)
const { menus, fetchMenus } = useMenuStore();

// Both work during migration period
```

## Persistence

The store persists to localStorage via Zustand's `persist` middleware:

```ts
{
  name: 'menu-storage',
  partialize: (state) => ({
    rawMenus: state.rawMenus,  // Cache raw menus
    collapsed: state.collapsed, // Cache UI state
  }),
}
```

**Note**: Filtered menus are NOT persisted (recomputed from raw menus + permissions).

## Testing

See `menuStore.test.ts` for complete test coverage:

```bash
npm test -- menuStore.test.ts
```

All 11 tests pass, covering:
- Fetch and filter menus
- Super admin permissions
- Nested menu structures
- Menu collapse state
- Breadcrumbs and paths
- Permission-based filtering

## Migration from AdminContext

To migrate from AdminContext to menuStore:

1. Replace `useAdmin()` with `useMenuStore()`
2. Update property names:
   - `refreshMenus()` → `fetchMenus()`
   - `rawMenus` (same)
   - `menus` (same)
3. Remove `AdminProvider` wrapper (optional)
4. Use `useAuthStore` for permission checks

## Example: Complete Sidebar Component

```tsx
import { useEffect } from 'react';
import { Menu } from 'antd';
import { useMenuStore, useAuthStore } from '@/stores';

function Sidebar() {
  const {
    menus,
    collapsed,
    openKeys,
    selectedKeys,
    fetchMenus,
    setOpenKeys,
    setSelectedKeys,
    toggleCollapsed,
  } = useMenuStore();

  const { user } = useAuthStore();

  useEffect(() => {
    if (user) {
      fetchMenus();
    }
  }, [user]);

  const renderMenuItems = (items: MenuItem[]) => {
    return items
      .filter(item => item.permission === '*' || user?.permissions.includes(item.permission))
      .map(item => {
        if (item.children?.length) {
          return (
            <Menu.SubMenu
              key={item.key}
              icon={item.icon && <Icon type={item.icon} />}
              title={item.label}
            >
              {renderMenuItems(item.children)}
            </Menu.SubMenu>
          );
        }
        return (
          <Menu.Item
            key={item.key}
            icon={item.icon && <Icon type={item.icon} />}
          >
            {item.label}
          </Menu.Item>
        );
      });
  };

  return (
    <Menu
      mode="inline"
      collapsed={collapsed}
      openKeys={openKeys}
      selectedKeys={selectedKeys}
      onOpenChange={setOpenKeys}
      onSelect={({ selectedKeys }) => setSelectedKeys(selectedKeys)}
    >
      {renderMenuItems(menus)}
    </Menu>
  );
}
```
