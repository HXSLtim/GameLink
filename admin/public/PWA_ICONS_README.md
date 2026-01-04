# PWA 图标说明

## 当前状态

✅ 已创建 SVG 图标: `icon.svg`
⚠️ 需要生成 PNG 图标以支持完整的 PWA 功能

## 需要的图标尺寸

PWA 需要以下 PNG 图标：

- `icon-192x192.png` (192x192 像素)
- `icon-512x512.png` (512x512 像素)

## 方法 1: 使用在线工具 (推荐)

### 方案 A: SVG to PNG Converter

1. 访问 https://cloudconvert.com/svg-to-png
2. 上传 `icon.svg`
3. 设置输出尺寸为 192x192，下载并重命名为 `icon-192x192.png`
4. 重复步骤，设置输出尺寸为 512x512，下载并重命名为 `icon-512x512.png`
5. 将两个 PNG 文件放入 `admin/public/` 目录

### 方案 B: RealFaviconGenerator (更全面)

1. 访问 https://realfavicongenerator.net/
2. 上传 `icon.svg`
3. 配置 PWA 设置：
   - iOS 设备: 选择 "Apple touch icon"
   - Android 设备: 选择 "Android Chrome"
   - Windows 设备: 选择 "Windows 8/10 tiles"
4. 下载生成的包
5. 将 `icon-192x192.png` 和 `icon-512x512.png` 复制到 `admin/public/` 目录

## 方法 2: 使用命令行工具

### 使用 ImageMagick (需要安装)

```bash
# Windows (使用 Chocolatey 安装)
choco install imagemagick

# 生成 192x192 图标
magick convert -background none -resize 192x192 icon.svg icon-192x192.png

# 生成 512x512 图标
magick convert -background none -resize 512x512 icon.svg icon-512x512.png
```

### 使用 Sharp (Node.js)

```bash
# 安装 sharp
npm install -g sharp-cli

# 生成图标
sharp input.svg -resize 192 192 --flatten background="{r:0,g:0,b:0,a:0}" icon-192x192.png
sharp input.svg -resize 512 512 --flatten background="{r:0,g:0,b:0,a:0}" icon-512x512.png
```

## 方法 3: 使用 Figma 或 Sketch

### Figma:

1. 在 Figma 中导入 `icon.svg`
2. 创建 192x192 和 512x512 的框架
3. 导出为 PNG (2x 或 3x)
4. 重命名并放入 `admin/public/` 目录

### Sketch:

1. 在 Sketch 中导入 `icon.svg`
2. 使用 Export 功能设置尺寸
3. 导出为 PNG
4. 重命名并放入 `admin/public/` 目录

## 验证图标

生成图标后，可以访问以下地址验证 PWA 配置：

- Chrome DevTools: Application → Manifest
- Lighthouse: 运行 PWA 审计
- 在线验证: https://www.pwabuilder.com/

## 当前临时方案

在生成 PNG 图标之前，PWA 功能仍然可以工作，但图标可能显示为默认浏览器图标。

建议尽快生成 PNG 图标以获得完整的 PWA 体验。

## 图标设计说明

当前 SVG 图标特点：

- **背景色**: #1890ff (Ant Design 蓝色)
- **图案**: 游戏手柄 + GL 文字
- **圆角**: 100px (大圆角)
- **风格**: 扁平化，现代简约

如需修改图标设计，编辑 `icon.svg` 文件并重新生成 PNG 即可。
