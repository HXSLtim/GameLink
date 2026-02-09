/**
 * GL 组件库 - 基于 uv-ui 的二次封装
 * 
 * 命名规范：gl-xxx
 * 特点：
 * - 主题感知（使用 CSS 变量）
 * - 响应式设计（PC/移动端适配）
 * - 统一的设计语言
 */

// 基础组件
export { default as GlButton } from './Button/index.vue'
export { default as GlTag } from './Tag/index.vue'
export { default as GlAvatar } from './Avatar/index.vue'
export { default as GlCard } from './Card/index.vue'
export { default as GlEmpty } from './Empty/index.vue'
export { default as GlInput } from './Input/index.vue'
export { default as GlSwitch } from './Switch/index.vue'

// 后续扩展：
// export { default as GlModal } from './Modal/index.vue'
// export { default as GlToast } from './Toast/index.vue'
// export { default as GlPopup } from './Popup/index.vue'
