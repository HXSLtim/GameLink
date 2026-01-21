const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

/**
 * Fix all Statistic styles.label issues in the codebase
 */

// Find all files with Statistic styles.label pattern
let files = [];
try {
  const result = execSync('grep -rl "styles.*label" src --include=*.tsx 2>/dev/null || true', {
    encoding: 'utf8',
    cwd: path.join(__dirname, '..')
  });
  files = result.split('\n').filter(f => f && f.trim());
} catch (e) {
  // Fallback: manually list remaining files
  files = [
    'src/pages/admin/Activity/index.tsx',
    'src/pages/admin/Activity/Rewards.tsx',
    'src/pages/admin/Commission/index.tsx',
    'src/pages/admin/Coupon/index.tsx',
    'src/pages/admin/Coupon/UserCoupon.tsx',
    'src/pages/admin/Dispute/index.tsx'
  ];
}

console.log('Found files with styles.label:', files.length);

let fixed = 0;
files.forEach(file => {
  const filePath = path.join(__dirname, '..', file);
  if (!fs.existsSync(filePath)) {
    console.log('⚠️  Not found:', file);
    return;
  }

  let content = fs.readFileSync(filePath, 'utf8');
  const orig = content;

  // Remove entire styles prop if it only contains label
  content = content.replace(/styles=\{\{\s*label:\s*\{[^}]*\}\s*\}\}\s*/g, '');

  // Remove label property from styles objects
  content = content.replace(/,\s*label:\s*\{[^}]+\}\s*(?=\})/g, '');

  // Handle multi-line label styles
  content = content.replace(/,\s*label:\s*\{[\s\S]*?\n\s*\}/g, '');

  if (content !== orig) {
    fs.writeFileSync(filePath, content);
    fixed++;
    console.log('✅ Fixed:', file);
  }
});

console.log(`\n✅ Total fixed: ${fixed} files`);
