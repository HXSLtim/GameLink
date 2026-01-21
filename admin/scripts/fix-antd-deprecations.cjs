const fs = require('fs');
const path = require('path');

// Files to fix for destroyOnClose -> destroyOnHidden
const destroyOnCloseFiles = [
  'src/pages/adminChat/rooms/index.tsx',
  'src/pages/admin/VIP/components/LevelForm.tsx',
  'src/pages/admin/Team/components/TeamForm.tsx',
  'src/pages/admin/Settlement/components/CompanyForm.tsx',
  'src/pages/admin/Routing/components/RoutingForm.tsx',
  'src/pages/admin/Referral/components/CodeForm.tsx',
  'src/pages/admin/Recharge/components/RefundModal.tsx',
  'src/pages/admin/Recharge/components/OptionForm.tsx',
  'src/pages/admin/GameRank/index.tsx'
];

// Files to fix for direction -> orientation in Space components
const directionFiles = [
  'src/pages/admin/Activity/index.tsx',
  'src/pages/player/Home/index.tsx',
  'src/pages/user/Wallet/index.tsx',
  'src/pages/admin/Routing/Test.tsx',
  'src/pages/admin/Service/index.tsx',
  'src/pages/admin/Routing/index.tsx',
  'src/pages/admin/Coupon/UserCoupon.tsx',
  'src/pages/admin/Coupon/index.tsx',
  'src/pages/admin/Activity/Rewards.tsx',
  'src/pages/admin/Profile/index.tsx',
  'src/pages/user/Ranking/index.tsx',
  'src/pages/sys/setting/index.tsx',
  'src/pages/player/Certification/Rank.tsx',
  'src/pages/player/Certification/Identity.tsx',
  'src/pages/admin/Routing/components/RoutingForm.tsx',
  'src/pages/admin/Review/Reports.tsx',
  'src/pages/admin/Recharge/components/RefundModal.tsx',
  'src/pages/admin/Recharge/components/OptionForm.tsx',
  'src/pages/admin/Notifications/index.tsx',
  'src/pages/admin/Dispute/components/ResolveModal.tsx',
  'src/pages/admin/Dispute/components/DisputeList.tsx',
  'src/pages/admin/Dispute/components/AssignModal.tsx',
  'src/components/ImportModal/index.tsx',
  'src/components/ImportHistoryTable/index.tsx'
];

let fixed = 0;

// Fix destroyOnClose -> destroyOnHidden
console.log('Fixing destroyOnClose -> destroyOnHidden...');
destroyOnCloseFiles.forEach(file => {
  const filePath = path.join(__dirname, '..', file);
  if (fs.existsSync(filePath)) {
    let content = fs.readFileSync(filePath, 'utf8');
    const original = content;
    // Handle destroyOnClose={true} or just destroyOnClose
    content = content.replace(/destroyOnClose\s*=/g, 'destroyOnHidden=');
    content = content.replace(/destroyOnClose:\s*/g, 'destroyOnHidden:');
    content = content.replace(/\bdestroyOnClose(\s?\n)/g, 'destroyOnHidden$1');
    content = content.replace(/\bdestroyOnClose(\s?\/?>)/g, 'destroyOnHidden$1');
    if (content !== original) {
      fs.writeFileSync(filePath, content);
      fixed++;
      console.log('  Fixed:', file);
    }
  }
});

// Fix direction -> orientation in Space components
console.log('\nFixing direction -> orientation in Space components...');
directionFiles.forEach(file => {
  const filePath = path.join(__dirname, '..', file);
  if (fs.existsSync(filePath)) {
    let content = fs.readFileSync(filePath, 'utf8');
    const original = content;
    // Fix direction="horizontal" -> orientation="horizontal"
    content = content.replace(/direction="horizontal"/g, 'orientation="horizontal"');
    content = content.replace(/direction='horizontal'/g, "orientation='horizontal'");
    content = content.replace(/direction="vertical"/g, 'orientation="vertical"');
    content = content.replace(/direction='vertical'/g, "orientation='vertical'");
    content = content.replace(/direction:\s*['"]horizontal['"]/g, "orientation: 'horizontal'");
    content = content.replace(/direction:\s*['"]vertical['"]/g, "orientation: 'vertical'");
    content = content.replace(/direction=\{(['"])(horizontal|vertical)\1\}/g, 'orientation={$2}');
    if (content !== original) {
      fs.writeFileSync(filePath, content);
      fixed++;
      console.log('  Fixed:', file);
    }
  }
});

console.log(`\nTotal fixes applied: ${fixed}`);
