const fs = require('fs');
const path = require('path');

/**
 * Simple fix for Ant Design 6.0 Statistic component
 * Removes the deprecated 'label' property from styles prop
 */

const filesToFix = [
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

console.log('Fixing Statistic styles.label issues...\n');

filesToFix.forEach(file => {
  const filePath = path.join(__dirname, '..', file);
  if (!fs.existsSync(filePath)) {
    console.log(`⚠️  Skipped: ${file} (not found)`);
    return;
  }

  let content = fs.readFileSync(filePath, 'utf8');
  const original = content;

  // Simply remove the label property from styles
  // styles={{ label: { color: '...' }, ... }}  ->  styles={{ ... }}
  content = content.replace(/styles=\{\{\s*label:\s*\{[^}]*\}\s*,?\s*/g, 'styles={{');
  content = content.replace(/,\s*label:\s*\{[^}]*\}\s*\}\}/g, ' }}');
  content = content.replace(/label:\s*\{[^}]*\},?\s*/g, '');

  if (content !== original) {
    fs.writeFileSync(filePath, content);
    fixed++;
    console.log(`✅ Fixed: ${file}`);
  } else {
    console.log(`ℹ️  No changes: ${file}`);
  }
});

console.log(`\n✅ Fixed ${fixed} files`);
