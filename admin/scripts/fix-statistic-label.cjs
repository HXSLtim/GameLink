const fs = require('fs');
const path = require('path');

/**
 * Fix Ant Design 6.0 Statistic component issues
 * - The 'label' prop was removed from Statistic component
 * - We need to wrap Statistic with a div/label or use the title prop pattern
 */

// Files with Statistic label issues (from build output)
const filesWithIssues = [
  'src/pages/admin/Order/components/OrderDetailTabs.tsx',
  'src/pages/admin/Player/components/PlayerDetailTabs.tsx',
  'src/pages/admin/Recharge/index.tsx',
  'src/pages/admin/Referral/index.tsx',
  'src/pages/admin/Routing/index.tsx',
  'src/pages/admin/Routing/Test.tsx',
  'src/pages/admin/Service/index.tsx',
  'src/pages/admin/Settlement/index.tsx',
  'src/pages/admin/Settlement/Players.tsx',
  'src/pages/admin/SettlementCompany/index.tsx',
  'src/pages/admin/Team/components/MemberCard.tsx',
  'src/pages/admin/Team/index.tsx',
  'src/pages/admin/VIP/index.tsx',
  'src/pages/admin/Withdraw/index.tsx',
  'src/pages/adminChat/records/index.tsx',
  'src/pages/adminChat/rooms/index.tsx',
  'src/pages/player/Earnings/index.tsx',
  'src/pages/player/Home/index.tsx',
  'src/pages/user/Wallet/index.tsx'
];

let fixed = 0;

console.log('Fixing Statistic component label prop issues...\n');

filesWithIssues.forEach(file => {
  const filePath = path.join(__dirname, '..', file);
  if (!fs.existsSync(filePath)) {
    console.log(`⚠️  Skipped (not found): ${file}`);
    return;
  }

  let content = fs.readFileSync(filePath, 'utf8');
  const original = content;

  // Pattern 1: Fix styles={{ label: ... }} in Statistic
  // In Ant Design 6.0, label prop was removed. We need to move it outside.
  content = content.replace(
    /(<Statistic[^>]*?)styles=\{\{\s*label:\s*([^}]+)\s*\}\}([^>]*>)/g,
    (match, prefix, labelStyle, suffix) => {
      // Check if the label is already provided separately
      return `${prefix}${suffix}`;
    }
  );

  // Pattern 2: Replace Statistic with label prop with a wrapped version
  // <Statistic label="xxx" ... />  ->  <div><div>xxx</div><Statistic ... /></div>
  content = content.replace(
    /<Statistic\s+([^>]*?)label=\{([^}]+)\}([^>]*)\s*\/>/g,
    (match, beforeProps, labelValue, afterProps) => {
      return `<div className="statistic-item"><div className="statistic-label">${labelValue}</div><Statistic ${beforeProps}${afterProps} /></div>`;
    }
  );

  // Pattern 3: Handle label with string literal
  // <Statistic label="text" ... />  ->  <div><div>text</div><Statistic ... /></div>
  content = content.replace(
    /<Statistic\s+([^>]*?)label="([^"]+)"([^>]*)\s*\/>/g,
    (match, beforeProps, labelText, afterProps) => {
      return `<div className="statistic-item"><div className="statistic-label">${labelText}</div><Statistic ${beforeProps}${afterProps} /></div>`;
    }
  );

  // Pattern 4: Handle closing tag version
  content = content.replace(
    /<Statistic\s+([^>]*?)label=\{([^}]+)\}([^>]*)>(\s*<\/Statistic>)/g,
    (match, beforeProps, labelValue, afterProps, closingTag) => {
      return `<div className="statistic-item"><div className="statistic-label">${labelValue}</div><Statistic ${beforeProps}${afterProps}>${closingTag}</div>`;
    }
  );

  content = content.replace(
    /<Statistic\s+([^>]*?)label="([^"]+)"([^>]*)>(\s*<\/Statistic>)/g,
    (match, beforeProps, labelText, afterProps, closingTag) => {
      return `<div className="statistic-item"><div className="statistic-label">${labelText}</div><Statistic ${beforeProps}${afterProps}>${closingTag}</div>`;
    }
  );

  if (content !== original) {
    fs.writeFileSync(filePath, content);
    fixed++;
    console.log(`✅ Fixed: ${file}`);
  } else {
    console.log(`ℹ️  No changes needed: ${file}`);
  }
});

console.log(`\n📊 Summary:`);
console.log(`   Files processed: ${filesWithIssues.length}`);
console.log(`   Files fixed: ${fixed}`);
console.log(`   Files skipped: ${filesWithIssues.length - fixed}`);

// Also add CSS suggestion
console.log(`\n💡 CSS Tip: Add this to your global styles for the new structure:`);
console.log(`
.statistic-item {
  display: inline-block;
  text-align: center;
}
.statistic-label {
  color: rgba(0, 0, 0, 0.45);
  font-size: 14px;
  margin-bottom: 4px;
}
`);
