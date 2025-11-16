# GameLink 国际化文档

## 🌍 概述

GameLink 项目支持多语言国际化，为不同地区的用户提供本地化体验。本文档详细说明了国际化实现方案、语言支持范围和扩展指南。

## 📋 支持语言

### 当前支持
- 🇨🇳 **简体中文** (zh-CN) - 默认语言
- 🇺🇸 **英语** (en-US) - 国际化语言
- 🇯🇵 **日语** (ja-JP) - 计划支持
- 🇰🇷 **韩语** (ko-KR) - 计划支持

### 语言代码标准
遵循 [BCP 47](https://tools.ietf.org/html/bcp47) 语言标签标准：
- `zh-CN` - 简体中文（中国大陆）
- `en-US` - 英语（美国）
- `ja-JP` - 日语（日本）
- `ko-KR` - 韩语（韩国）

---

## 🎨 前端国际化

### 技术方案
使用 `react-i18next` 作为国际化解决方案：

```typescript
// i18n/index.ts
import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import Backend from 'i18next-http-backend';
import LanguageDetector from 'i18next-browser-languagedetector';

// 语言资源
import zhCN from './locales/zh-CN.json';
import enUS from './locales/en-US.json';

const resources = {
  'zh-CN': {
    translation: zhCN,
  },
  'en-US': {
    translation: enUS,
  },
};

i18n
  .use(Backend)
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    resources,
    fallbackLng: 'zh-CN',
    debug: process.env.NODE_ENV === 'development',
    interpolation: {
      escapeValue: false,
    },
  });

export default i18n;
```

### 语言资源结构
```
frontend/src/i18n/
├── locales/
│   ├── zh-CN/
│   │   ├── common.json      # 通用词汇
│   │   ├── auth.json        # 认证相关
│   │   ├── order.json       # 订单相关
│   │   ├── payment.json     # 支付相关
│   │   ├── player.json      # 陪玩师相关
│   │   └── admin.json      # 管理后台
│   ├── en-US/
│   │   ├── common.json
│   │   ├── auth.json
│   │   ├── order.json
│   │   ├── payment.json
│   │   ├── player.json
│   │   └── admin.json
│   ├── ja-JP/             # 日语（计划）
│   └── ko-KR/             # 韩语（计划）
├── index.ts                # i18n 配置
└── types.ts               # 类型定义
```

### 语言资源示例

#### 中文 (zh-CN/common.json)
```json
{
  "app": {
    "name": "GameLink",
    "description": "游戏陪玩管理平台"
  },
  "navigation": {
    "home": "首页",
    "games": "游戏",
    "players": "陪玩师",
    "orders": "订单",
    "profile": "个人中心",
    "settings": "设置"
  },
  "actions": {
    "create": "创建",
    "edit": "编辑",
    "delete": "删除",
    "save": "保存",
    "cancel": "取消",
    "confirm": "确认",
    "submit": "提交"
  },
  "status": {
    "pending": "待处理",
    "confirmed": "已确认",
    "in_progress": "进行中",
    "completed": "已完成",
    "cancelled": "已取消"
  }
}
```

#### 英语 (en-US/common.json)
```json
{
  "app": {
    "name": "GameLink",
    "description": "Game Companion Management Platform"
  },
  "navigation": {
    "home": "Home",
    "games": "Games",
    "players": "Players",
    "orders": "Orders",
    "profile": "Profile",
    "settings": "Settings"
  },
  "actions": {
    "create": "Create",
    "edit": "Edit",
    "delete": "Delete",
    "save": "Save",
    "cancel": "Cancel",
    "confirm": "Confirm",
    "submit": "Submit"
  },
  "status": {
    "pending": "Pending",
    "confirmed": "Confirmed",
    "in_progress": "In Progress",
    "completed": "Completed",
    "cancelled": "Cancelled"
  }
}
```

### React 组件中使用
```typescript
import { useTranslation } from 'react-i18next';

export const OrderList: React.FC = () => {
  const { t } = useTranslation('order');

  return (
    <div>
      <h1>{t('title')}</h1>
      <button>{t('actions.create')}</button>
      <span>{t('status.pending')}</span>
    </div>
  );
};

// 使用命名空间
const { t } = useTranslation(['common', 'order']);
const title = t('common:actions.create');
const orderTitle = t('order:title');
```

### 语言切换组件
```typescript
import { useTranslation } from 'react-i18next';

export const LanguageSwitcher: React.FC = () => {
  const { i18n } = useTranslation();

  const languages = [
    { code: 'zh-CN', name: '简体中文', flag: '🇨🇳' },
    { code: 'en-US', name: 'English', flag: '🇺🇸' },
  ];

  const handleLanguageChange = (languageCode: string) => {
    i18n.changeLanguage(languageCode);
    localStorage.setItem('language', languageCode);
  };

  return (
    <div className="language-switcher">
      {languages.map((lang) => (
        <button
          key={lang.code}
          onClick={() => handleLanguageChange(lang.code)}
          className={i18n.language === lang.code ? 'active' : ''}
        >
          <span className="flag">{lang.flag}</span>
          <span className="name">{lang.name}</span>
        </button>
      ))}
    </div>
  );
};
```

---

## 🔧 后端国际化

### API 响应国际化
```go
// 国际化中间件
func I18nMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 从请求头获取语言偏好
        lang := c.GetHeader("Accept-Language")
        if lang == "" {
            lang = "zh-CN" // 默认语言
        }

        // 验证支持的语言
        supportedLangs := []string{"zh-CN", "en-US"}
        if !contains(supportedLangs, lang) {
            lang = "zh-CN"
        }

        c.Set("lang", lang)
        c.Next()
    }
}

// 错误消息国际化
type ErrorMessage struct {
    Code    string                    `json:"code"`
    Message map[string]string          `json:"message"`
    Details map[string]interface{}     `json:"details,omitempty"`
}

func (e *ErrorMessage) Localize(lang string) string {
    if msg, exists := e.Message[lang]; exists {
        return msg
    }
    return e.Message["zh-CN"] // 默认中文
}

// 使用示例
func (h *OrderHandler) CreateOrder(c *gin.Context) {
    lang := c.GetString("lang")
    
    // 国际化错误消息
    errorMsg := &ErrorMessage{
        Code: "PLAYER_NOT_FOUND",
        Message: map[string]string{
            "zh-CN": "陪玩师不存在",
            "en-US": "Player not found",
        },
    }
    
    c.JSON(404, gin.H{
        "error": errorMsg.Localize(lang),
    })
}
```

### 数据库多语言支持
```go
// 多语言字段模型
type Game struct {
    ID          uint64    `gorm:"primaryKey"`
    Name        string    `gorm:"not null"`        // 默认名称（中文）
    NameEn      string    `gorm:"column:name_en"` // 英文名称
    NameJa      string    `gorm:"column:name_ja"` // 日文名称
    NameKo      string    `gorm:"column:name_ko"` // 韩文名称
    Description string    `gorm:"not null"`        // 默认描述
    DescEn      string    `gorm:"column:desc_en"` // 英文描述
    DescJa      string    `gorm:"column:desc_ja"` // 日文描述
    DescKo      string    `gorm:"column:desc_ko"` // 韩文描述
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// 本地化方法
func (g *Game) Localize(lang string) (name, desc string) {
    switch lang {
    case "en-US":
        name = g.NameEn
        desc = g.DescEn
    case "ja-JP":
        name = g.NameJa
        desc = g.DescJa
    case "ko-KR":
        name = g.NameKo
        desc = g.DescKo
    default:
        name = g.Name
        desc = g.Description
    }
    
    if name == "" {
        name = g.Name // 降级到默认语言
    }
    if desc == "" {
        desc = g.Description
    }
    
    return name, desc
}

// 在 Service 中使用
func (s *GameService) GetLocalizedGame(ctx context.Context, id uint64, lang string) (*GameLocalized, error) {
    game, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    name, desc := game.Localize(lang)
    
    return &GameLocalized{
        ID:          game.ID,
        Name:        name,
        Description: desc,
        Original:    game,
    }, nil
}
```

---

## 📱 移动端国际化

### 移动端配置
```typescript
// 移动端 i18n 配置
const mobileI18nConfig = {
  supportedLngs: ['zh-CN', 'en-US'],
  fallbackLng: 'zh-CN',
  debug: __DEV__,
  
  interpolation: {
    escapeValue: false,
  },
  
  detection: {
    order: ['localStorage', 'navigator', 'htmlTag'],
    caches: ['localStorage'],
  },
  
  backend: {
    loadPath: '/assets/locales/{{lng}}/{{ns}}.json',
  },
};
```

### 移动端语言资源
```
mobile/assets/locales/
├── zh-CN/
│   ├── common.json
│   ├── auth.json
│   └── order.json
└── en-US/
    ├── common.json
    ├── auth.json
    └── order.json
```

---

## 🛠️ 开发工具

### 翻译管理脚本
```bash
#!/bin/bash
# scripts/i18n.sh

# 提取待翻译文本
npm run i18n:extract

# 检查缺失翻译
npm run i18n:check

# 合并翻译文件
npm run i18n:merge

# 验证翻译完整性
npm run i18n:validate
```

### 自动翻译工具
```typescript
// scripts/auto-translate.ts
import GoogleTranslate from '@google-cloud/translate';
import fs from 'fs';
import path from 'path';

const translate = new GoogleTranslate();

async function translateFile(sourcePath: string, targetLang: string) {
  const sourceContent = JSON.parse(fs.readFileSync(sourcePath, 'utf8'));
  const translatedContent = {};
  
  for (const [key, value] of Object.entries(sourceContent)) {
    if (typeof value === 'string') {
      const [translation] = await translate.translate(value, targetLang);
      translatedContent[key] = translation;
    } else {
      translatedContent[key] = value;
    }
  }
  
  const targetPath = sourcePath.replace('/zh-CN/', `/${targetLang}/`);
  fs.writeFileSync(targetPath, JSON.stringify(translatedContent, null, 2));
}
```

---

## 📊 翻译覆盖率

### 统计工具
```typescript
// scripts/coverage.ts
interface CoverageReport {
  language: string;
  namespace: string;
  totalKeys: number;
  translatedKeys: number;
  coverage: number;
}

function generateCoverageReport(): CoverageReport[] {
  const languages = ['zh-CN', 'en-US'];
  const namespaces = ['common', 'auth', 'order', 'payment'];
  
  return languages.map(lang => 
    namespaces.map(ns => calculateCoverage(lang, ns))
  ).flat();
}
```

### 覆盖率报告格式
```json
{
  "generated": "2025-11-16T21:30:00Z",
  "languages": {
    "zh-CN": {
      "total": 450,
      "translated": 450,
      "coverage": "100%"
    },
    "en-US": {
      "total": 450,
      "translated": 380,
      "coverage": "84.4%",
      "missing": ["order.player_skills", "payment.refund_reason"]
    }
  }
}
```

---

## 🎯 最佳实践

### 翻译原则
1. **简洁明了** - 避免过长的翻译文本
2. **术语一致** - 建立术语表，确保翻译一致性
3. **文化适配** - 考虑目标语言的文化习惯
4. **技术术语** - 技术词汇保持行业通用译法

### 命名规范
- 使用 camelCase 命名键名
- 按功能模块组织命名空间
- 避免使用特殊字符和空格

### 版本管理
- 翻译文件与代码版本同步
- 使用语义化版本控制
- 记录翻译变更历史

---

## 🚀 扩展指南

### 添加新语言
1. **创建语言目录**
   ```bash
   mkdir -p frontend/src/i18n/locales/ko-KR
   ```

2. **添加语言资源**
   ```json
   // ko-KR/common.json
   {
     "app": {
       "name": "GameLink",
       "description": "게임 동행 관리 플랫폼"
     }
   }
   ```

3. **更新配置**
   ```typescript
   // i18n/index.ts
   import koKR from './locales/ko-KR/common.json';
   
   const resources = {
     'zh-CN': { translation: zhCN },
     'en-US': { translation: enUS },
     'ko-KR': { translation: koKR }, // 新增
   };
   ```

4. **测试验证**
   ```typescript
   // 测试新语言
   describe('Korean Translation', () => {
     it('should display correct app name', () => {
       i18n.changeLanguage('ko-KR');
       expect(t('app.name')).toBe('GameLink');
     });
   });
   ```

### 翻译流程
1. **提取新文本** - 使用脚本提取代码中的待翻译文本
2. **分配翻译任务** - 将翻译文件分配给翻译人员
3. **质量审核** - 审核翻译质量和一致性
4. **集成测试** - 在实际环境中测试翻译效果
5. **发布上线** - 将翻译集成到生产环境

---

## 📞 支持与贡献

### 翻译贡献
欢迎社区贡献翻译：
1. Fork 项目仓库
2. 创建语言分支 `feature/i18n/language-code`
3. 添加翻译文件
4. 提交 Pull Request

### 问题反馈
- 📋 **翻译问题**: [Issues](https://github.com/HXSLtim/GameLink/issues)
- 💬 **翻译讨论**: [Discussions](https://github.com/HXSLtim/GameLink/discussions)

---

## 📄 参考资源

- [react-i18next 官方文档](https://react.i18next.com/)
- [BCP 47 语言标签标准](https://tools.ietf.org/html/bcp47)
- [Unicode CLDR 数据](https://cldr.unicode.org/)
- [Google Cloud Translation API](https://cloud.google.com/translate)

---

**最后更新: 2025-11-16 | 版本: v1.0**
