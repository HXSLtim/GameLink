# 国际化实现指南

## 当前状态

已创建基础的 i18n 架构，包括：

1. **基础工具** (`src/i18n/index.ts`)
   - `getCurrentLocale()` - 获取当前语言
   - `setLocale()` - 设置语言
   - `t()` - 翻译函数（占位）

2. **翻译文件**
   - `src/i18n/locales/zh-CN.ts` - 中文翻译
   - `src/i18n/locales/en-US.ts` - 英文翻译

3. **支持的语言**
   - `zh-CN` - 简体中文
   - `en-US` - 美国英语

## 下一步实施（建议）

### 方案 A：使用 react-i18next（推荐）

```bash
npm install react-i18next i18next
```

**优点**：
- 功能完整，社区成熟
- 支持复数、插值、命名空间等高级功能
- 性能优化良好
- 提供 React Hooks

**示例配置**：

```typescript
// src/i18n/config.ts
import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import { zhCN } from './locales/zh-CN';
import { enUS } from './locales/en-US';

i18n
  .use(initReactI18next)
  .init({
    resources: {
      'zh-CN': { translation: zhCN },
      'en-US': { translation: enUS },
    },
    lng: 'zh-CN',
    fallbackLng: 'zh-CN',
    interpolation: {
      escapeValue: false,
    },
  });

export default i18n;
```

**使用方式**：

```typescript
import { useTranslation } from 'react-i18next';

function MyComponent() {
  const { t } = useTranslation();
  
  return (
    <div>
      <h1>{t('common.confirm')}</h1>
      <button>{t('auth.login')}</button>
    </div>
  );
}
```

### 方案 B：使用 Arco Design 的国际化

Arco Design 自带国际化支持：

```typescript
import { ConfigProvider } from '@arco-design/web-react';
import zhCN from '@arco-design/web-react/es/locale/zh-CN';
import enUS from '@arco-design/web-react/es/locale/en-US';

function App() {
  const [locale, setLocale] = useState(zhCN);
  
  return (
    <ConfigProvider locale={locale}>
      {/* 应用内容 */}
    </ConfigProvider>
  );
}
```

### 方案 C：混合方案（推荐用于本项目）

结合 react-i18next（业务文本）和 Arco Design 国际化（UI 组件）：

```typescript
// main.tsx
import i18n from './i18n/config';
import { I18nextProvider } from 'react-i18next';
import { ConfigProvider } from '@arco-design/web-react';
import zhCN from '@arco-design/web-react/es/locale/zh-CN';

root.render(
  <I18nextProvider i18n={i18n}>
    <ConfigProvider locale={zhCN}>
      <App />
    </ConfigProvider>
  </I18nextProvider>
);
```

## 实施步骤

### 1. 安装依赖（如选择方案 A 或 C）

```bash
npm install react-i18next i18next
```

### 2. 配置 i18n

创建 `src/i18n/config.ts`（参考上面示例）

### 3. 在 main.tsx 中集成

```typescript
import './i18n/config'; // 导入配置

// 其余代码...
```

### 4. 创建语言切换组件

```typescript
// src/components/LanguageSwitcher/LanguageSwitcher.tsx
import { Select } from '@arco-design/web-react';
import { useTranslation } from 'react-i18next';

export const LanguageSwitcher: React.FC = () => {
  const { i18n } = useTranslation();
  
  return (
    <Select
      value={i18n.language}
      onChange={(value) => i18n.changeLanguage(value)}
      options={[
        { label: '简体中文', value: 'zh-CN' },
        { label: 'English', value: 'en-US' },
      ]}
    />
  );
};
```

### 5. 替换硬编码文本

逐步将硬编码的中文文本替换为 `t()` 函数调用：

```typescript
// Before
<h1>用户管理</h1>

// After
<h1>{t('menu.users')}</h1>
```

## 翻译键命名规范

- 使用点号分隔命名空间：`namespace.key`
- 命名空间示例：
  - `common.*` - 通用文本
  - `auth.*` - 认证相关
  - `menu.*` - 菜单项
  - `error.*` - 错误消息
  - `validation.*` - 表单验证
  - `page.{pageName}.*` - 页面特定文本

## 类型安全

使用 TypeScript 确保翻译键的类型安全：

```typescript
// src/i18n/types.ts
import { TFunction } from 'i18next';
import type { TranslationKeys } from './locales/zh-CN';

// 类型安全的 t 函数
export type TypeSafeTranslate = TFunction<'translation', undefined, TranslationKeys>;
```

## 注意事项

1. **性能优化**
   - 使用命名空间避免加载所有翻译
   - 考虑翻译文件的懒加载

2. **日期和数字格式化**
   - 使用 `Intl` API 或 `dayjs` 处理日期
   - 使用 `Intl.NumberFormat` 处理数字

3. **复数和性别**
   - i18next 支持复数规则
   - 需要为不同语言配置不同规则

4. **测试**
   - 为每种语言创建截图测试
   - 验证文本长度不会破坏布局

5. **持续集成**
   - 使用工具检测遗漏的翻译
   - 自动化翻译流程（可选）

## 工具推荐

- **i18n Ally** (VSCode 插件) - 翻译管理
- **BabelEdit** - 翻译编辑器
- **Crowdin** / **Lokalise** - 协作翻译平台

## 示例代码

完整示例见：`src/i18n/` 目录

## 参考资料

- [react-i18next 文档](https://react.i18next.com/)
- [i18next 文档](https://www.i18next.com/)
- [Arco Design 国际化](https://arco.design/react/docs/i18n)

