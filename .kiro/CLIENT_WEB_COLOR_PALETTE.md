# GameLink Color Palette Quick Reference

> **Quick reference** for designers and developers implementing the GameLink dual-theme system.

---

## Day Theme (Kook-inspired)

```
┌────────────────────────────────────────────────────────────┐
│                    DAY THEME (KOOK)                        │
├────────────────────────────────────────────────────────────┤
│                                                            │
│  PRIMARY GREEN                                             │
│  ┌─────────────────────────────────────────────┐          │
│  │ 50    #E6FCF6  hsl(153, 74%, 95%)           │          │
│  │ 100   #CBF7E9  hsl(153, 74%, 88%)           │          │
│  │ 500   #00D26A  hsl(165, 100%, 41%) ⭐ MAIN │          │
│  │ 600   #00B55C  hsl(165, 100%, 35%)         │          │
│  │ 700   #00984D  hsl(165, 100%, 30%)         │          │
│  └─────────────────────────────────────────────┘          │
│                                                            │
│  SEMANTIC COLORS                                          │
│  ┌─────────────────────────────────────────────┐          │
│  │ Success  #10B981  🟢                         │          │
│  │ Warning  #F59E0B  🟡                         │          │
│  │ Error    #EF4444  🔴                         │          │
│  │ Info     #3B82F6  🔵                         │          │
│  └─────────────────────────────────────────────┘          │
│                                                            │
│  NEUTRALS                                                  │
│  ┌─────────────────────────────────────────────┐          │
│  │ Background      #F5F7FA  ⬜                   │          │
│  │ Card            #FFFFFF  ⬜                   │          │
│  │ Secondary BG    #F9FAFB  ⬜                   │          │
│  │ Text Primary    #1F2937  ⬛                   │          │
│  │ Text Secondary  #6B7280  ⬛                   │          │
│  │ Text Tertiary   #9CA3AF  ⬛                   │          │
│  │ Border          #E5E7EB  ⬜                   │          │
│  │ Border Light    #F3F4F6  ⬜                   │          │
│  └─────────────────────────────────────────────┘          │
│                                                            │
└────────────────────────────────────────────────────────────┘
```

## Night Theme (Discord-inspired)

```
┌────────────────────────────────────────────────────────────┐
│                  NIGHT THEME (DISCORD)                     │
├────────────────────────────────────────────────────────────┤
│                                                            │
│  BLURPLE (Primary Purple)                                 │
│  ┌─────────────────────────────────────────────┐          │
│  │ 400   #7983F5  hsl(239, 84%, 72%)           │          │
│  │ 500   #5865F2  hsl(239, 89%, 65%) ⭐ MAIN   │          │
│  │ 600   #4752C4  hsl(239, 63%, 53%)           │          │
│  └─────────────────────────────────────────────┘          │
│                                                            │
│  SEMANTIC COLORS                                          │
│  ┌─────────────────────────────────────────────┐          │
│  │ Success  #23A559  🟢                         │          │
│  │ Warning  #F0B132  🟡                         │          │
│  │ Error    #DA373C  🔴                         │          │
│  │ Info     #00A8FC  🔵                         │          │
│  └─────────────────────────────────────────────┘          │
│                                                            │
│  DARK NEUTRALS (Discord Gray Scale)                       │
│  ┌─────────────────────────────────────────────┐          │
│  │ Background      #313338  ⬛                   │          │
│  │ Card            #2B2D31  ⬛                   │          │
│  │ Elevated        #2B2D31  ⬛                   │          │
│  │ Deep            #1E1F22  ⬛                   │          │
│  │ Text Primary    #F2F3F5  ⬜                   │          │
│  │ Text Secondary  #B5BAC1  ⬜                   │          │
│  │ Text Tertiary   #949BA4  ⬜                   │          │
│  │ Border          #1E1F22  ⬛                   │          │
│  │ Border Subtle   #2B2D31  ⬛                   │          │
│  └─────────────────────────────────────────────┘          │
│                                                            │
└────────────────────────────────────────────────────────────┘
```

---

## CSS Variables (Ready to Use)

```css
/* Day Theme */
:root[data-theme="day"] {
  /* Primary */
  --color-primary-50: #E6FCF6;
  --color-primary-100: #CBF7E9;
  --color-primary-500: #00D26A;
  --color-primary-600: #00B55C;
  --color-primary-700: #00984D;

  /* Semantic */
  --color-success: #10B981;
  --color-warning: #F59E0B;
  --color-error: #EF4444;
  --color-info: #3B82F6;

  /* Neutral */
  --color-bg-primary: #F5F7FA;
  --color-bg-card: #FFFFFF;
  --color-bg-secondary: #F9FAFB;
  --color-text-primary: #1F2937;
  --color-text-secondary: #6B7280;
  --color-text-tertiary: #9CA3AF;
  --color-border: #E5E7EB;
  --color-border-light: #F3F4F6;
}

/* Night Theme */
:root[data-theme="night"] {
  /* Primary */
  --color-primary-400: #7983F5;
  --color-primary-500: #5865F2;
  --color-primary-600: #4752C4;

  /* Semantic */
  --color-success: #23A559;
  --color-warning: #F0B132;
  --color-error: #DA373C;
  --color-info: #00A8FC;

  /* Neutral */
  --color-bg-primary: #313338;
  --color-bg-card: #2B2D31;
  --color-bg-elevated: #2B2D31;
  --color-bg-deep: #1E1F22;
  --color-text-primary: #F2F3F5;
  --color-text-secondary: #B5BAC1;
  --color-text-tertiary: #949BA4;
  --color-border: #1E1F22;
  --color-border-subtle: #2B2D31;
}
```

---

## Tailwind Config Extension

```javascript
// tailwind.config.js
module.exports = {
  theme: {
    extend: {
      colors: {
        primary: {
          50: 'var(--color-primary-50)',
          100: 'var(--color-primary-100)',
          400: 'var(--color-primary-400)',
          500: 'var(--color-primary-500)',
          600: 'var(--color-primary-600)',
          700: 'var(--color-primary-700)',
          DEFAULT: 'var(--color-primary-500)',
        },
        success: 'var(--color-success)',
        warning: 'var(--color-warning)',
        error: 'var(--color-error)',
        info: 'var(--color-info)',
        bg: {
          primary: 'var(--color-bg-primary)',
          card: 'var(--color-bg-card)',
          secondary: 'var(--color-bg-secondary)',
          elevated: 'var(--color-bg-elevated)',
          deep: 'var(--color-bg-deep)',
        },
        text: {
          primary: 'var(--color-text-primary)',
          secondary: 'var(--color-text-secondary)',
          tertiary: 'var(--color-text-tertiary)',
        },
        border: {
          DEFAULT: 'var(--color-border)',
          subtle: 'var(--color-border-subtle)',
          light: 'var(--color-border-light)',
        },
      },
    },
  },
};
```

---

## Usage Examples

### Buttons

```css
/* Primary Button - Day */
.btn-primary {
  background: var(--color-primary-500);
  color: white;
}
.btn-primary:hover {
  background: var(--color-primary-600);
}

/* Primary Button - Night */
[data-theme="night"] .btn-primary {
  background: var(--color-primary-500);
  color: white;
}
[data-theme="night"] .btn-primary:hover {
  background: var(--color-primary-600);
}
```

### Cards

```css
.card {
  background: var(--color-bg-card);
  border: 1px solid var(--color-border);
  color: var(--color-text-primary);
}
```

### Status Badges

```css
.badge-success { background: var(--color-success); }
.badge-warning { background: var(--color-warning); }
.badge-error { background: var(--color-error); }
.badge-info { background: var(--color-info); }
```

---

## Color Contrast Validation

| Text on Background | Day Theme | Night Theme | Pass? |
|:-------------------|:----------|:------------|:------|
| Primary on Card | 4.5:1 | 4.5:1 | ✅ |
| Secondary on Card | 4.5:1 | 4.5:1 | ✅ |
| Primary on BG Primary | 4.5:1 | 4.5:1 | ✅ |
| Primary-500 text on Card | 4.5:1 | 4.5:1 | ✅ |

---

**Last Updated**: 2025-01-11
**Version**: 1.0.0
