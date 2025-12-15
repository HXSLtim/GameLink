/**
 * 权限变更事件工具
 * 用于跨组件/跨标签页通知权限变更
 */

/**
 * 权限变更事件名称
 */
export const PERMISSION_CHANGE_EVENT = 'gamelink:permission-change';

/**
 * 触发权限变更事件
 * 可在任何地方调用以通知权限已变更
 */
export const triggerPermissionChange = () => {
    window.dispatchEvent(new CustomEvent(PERMISSION_CHANGE_EVENT));
    // 同时触发 storage 事件以通知其他标签页
    const timestamp = Date.now().toString();
    localStorage.setItem('permission_change_timestamp', timestamp);
};
