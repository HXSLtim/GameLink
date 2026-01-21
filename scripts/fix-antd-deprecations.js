const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

const adminDir = 'D:\\Desktop\\Code\\GameLink\\admin';

// Files to fix valueStyle
const filesToFix = [
    'src/pages/adminActivity/Rewards.tsx',
    'src/pages/adminCommission/index.tsx',
    'src/pages/adminCoupon/index.tsx',
    'src/pages/adminCoupon/UserCoupon.tsx',
    'src/pages/adminDispute/index.tsx',
    'src/pages/adminOrder/components/OrderDetailTabs.tsx',
    'src/pages/adminPlayer/components/PlayerDetailTabs.tsx',
    'src/pages/adminRecharge/index.tsx',
    'src/pages/adminReferral/Rewards.tsx',
    'src/pages/adminRouting/index.tsx',
    'src/pages/adminService/index.tsx',
    'src/pages/adminSettlement/index.tsx',
    'src/pages/adminSettlement/Players.tsx',
    src/pages/adminTeam/components/MemberCard.tsx,
    'src/pages/adminTeam/index.tsx',
    'src/pages/adminVIP/index.tsx',
    'src/pages/adminWithdraw/index.tsx',
];

let fixedCount = 0;
let errorCount = 0;

filesToFix.forEach(file => {
    const filePath = path.join(adminDir, file);

    if (!fs.existsSync(filePath)) {
        console.log(`⚠️  跳过不存在的文件: ${file}`);
        return;
    }

    try {
        let content = fs.readFileSync(filePath, 'utf8');
        const originalContent = content;

        // Replace valueStyle= to styles={{ label: {
        content = content.replace(
            /valueStyle=\{([^}]+)\}/g,
            'styles={{ label: {$1}'
        );

        if (content !== originalContent) {
            const changes = (content.match(/valueStyle=/g) || []).length;
            fixedCount++;
            fs.writeFileSync(filePath, content, 'utf-8');
            console.log(`✅ 已修复: ${file} (${changes} 处)`);
        } else {
            console.log(`ℹ️  无需修复: ${file}`);
        }
    } catch (err) {
        errorCount++;
        console.error(`❌ 错误: ${file}`);
        console.error(err);
    }
});

console.log(`\n修复完成: ${fixedCount} 个文件, ${errorCount} 个错误`);
