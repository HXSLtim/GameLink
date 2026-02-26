import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const scriptFile = fileURLToPath(import.meta.url);
const scriptDir = path.dirname(scriptFile);
const repoRoot = path.resolve(scriptDir, "..");

const files = {
  prd: path.join(repoRoot, "docs", "PRD.md"),
  governance: path.join(repoRoot, "docs", "PRD_GOVERNANCE.md"),
  roadmap: path.join(repoRoot, "docs", "PRODUCT_ROADMAP.md"),
  progress: path.join(repoRoot, "docs", "PROGRESS.md"),
  comprehensive: path.join(repoRoot, "docs", "PRD_COMPREHENSIVE.md"),
};

const reportDir = path.join(repoRoot, "docs", "reports");
const reportPath = path.join(reportDir, "prd-validation-report.md");

const hardErrors = [];
const warnings = [];

function readFileSafe(filePath) {
  if (!fs.existsSync(filePath)) {
    hardErrors.push(`缺少文件: ${path.relative(repoRoot, filePath)}`);
    return "";
  }

  return fs.readFileSync(filePath, "utf8");
}

function has(text, pattern) {
  return pattern.test(text);
}

function checkContains(text, regex, description, level = "hard") {
  if (!has(text, regex)) {
    if (level === "hard") {
      hardErrors.push(description);
      return;
    }

    warnings.push(description);
  }
}

function checkNoMatch(text, regex, description, level = "hard") {
  if (has(text, regex)) {
    if (level === "hard") {
      hardErrors.push(description);
      return;
    }

    warnings.push(description);
  }
}

function parseLastUpdated(text) {
  const match = text.match(/\*\*\s*最后更新\s*\*\*\s*[:：]\s*(\d{4}-\d{2}-\d{2})|最后更新\s*[:：]\s*(\d{4}-\d{2}-\d{2})/);
  if (!match) {
    return null;
  }

  return match[1] ?? match[2] ?? null;
}

function daysBetween(isoDate, now = new Date()) {
  const date = new Date(`${isoDate}T00:00:00`);
  if (Number.isNaN(date.getTime())) {
    return null;
  }

  const diff = now.getTime() - date.getTime();
  return Math.floor(diff / (1000 * 60 * 60 * 24));
}

function pushStalenessWarning(docName, text, thresholdDays) {
  const updated = parseLastUpdated(text);
  if (!updated) {
    warnings.push(`${docName} 缺少“最后更新”字段`);
    return;
  }

  const age = daysBetween(updated);
  if (age === null) {
    warnings.push(`${docName} 的“最后更新”格式无效: ${updated}`);
    return;
  }

  if (age > thresholdDays) {
    warnings.push(`${docName} 可能过期，已 ${age} 天未更新（阈值 ${thresholdDays}）`);
  }
}

function main() {
  const prd = readFileSafe(files.prd);
  const governance = readFileSafe(files.governance);
  const roadmap = readFileSafe(files.roadmap);
  const progress = readFileSafe(files.progress);
  const comprehensive = readFileSafe(files.comprehensive);

  checkContains(prd, /SSOT/i, "PRD.md 未声明 SSOT（单一事实源）");
  checkContains(prd, /PRD_GOVERNANCE\.md/, "PRD.md 未链接治理规则文档");

  checkContains(governance, /单一事实源/, "PRD_GOVERNANCE.md 缺少“单一事实源”章节");
  checkContains(governance, /更新规则/, "PRD_GOVERNANCE.md 缺少“更新规则”章节");
  checkContains(governance, /Owner/, "PRD_GOVERNANCE.md 缺少 Owner 机制");

  checkContains(roadmap, /只维护里程碑|仅维护“?时间与里程碑/, "PRODUCT_ROADMAP.md 未声明只维护里程碑/时间规划");
  checkContains(progress, /仅记录交付进度与阻塞项|仅记录交付进度/, "PROGRESS.md 未声明只维护进度/阻塞项");
  checkContains(comprehensive, /主需求事实源请以\s*`docs\/PRD\.md`\s*为准/, "PRD_COMPREHENSIVE.md 未声明从属于 PRD SSOT");

  const coreDocs = [
    ["docs/PRD.md", prd],
    ["docs/PRD_GOVERNANCE.md", governance],
    ["docs/PRODUCT_ROADMAP.md", roadmap],
    ["docs/PROGRESS.md", progress],
    ["docs/PRD_COMPREHENSIVE.md", comprehensive],
  ];

  for (const [name, content] of coreDocs) {
    checkNoMatch(content, /app-react/i, `${name} 存在过时路径 app-react`);
  }

  for (const [name, content] of coreDocs) {
    const lines = content.split(/\r?\n/);
    for (const line of lines) {
      if (!/(小程序\/H5|uni-app)/i.test(line)) {
        continue;
      }

      if (/版本历史|v\d+\.\d+/i.test(line)) {
        continue;
      }

      warnings.push(`${name} 可能存在历史技术栈残留: ${line.trim()}`);
    }
  }

  pushStalenessWarning("docs/PRD.md", prd, 30);
  pushStalenessWarning("docs/PRODUCT_ROADMAP.md", roadmap, 30);
  pushStalenessWarning("docs/PROGRESS.md", progress, 30);

  fs.mkdirSync(reportDir, { recursive: true });

  const status = hardErrors.length === 0 ? "PASS" : "FAIL";
  const report = [
    "# PRD 文档校验报告",
    "",
    `- 时间: ${new Date().toISOString()}`,
    `- 状态: ${status}`,
    `- 阻断问题: ${hardErrors.length}`,
    `- 告警问题: ${warnings.length}`,
    "",
    "## 阻断问题（P0）",
    ...(hardErrors.length === 0 ? ["- 无"] : hardErrors.map((item) => `- ${item}`)),
    "",
    "## 告警问题（P1）",
    ...(warnings.length === 0 ? ["- 无"] : warnings.map((item) => `- ${item}`)),
    "",
    "## 校验范围",
    ...Object.values(files).map((item) => `- ${path.relative(repoRoot, item)}`),
    "",
  ].join("\n");

  fs.writeFileSync(reportPath, report, "utf8");

  console.log(`[PRD-CHECK] ${status}`);
  console.log(`[PRD-CHECK] 阻断问题: ${hardErrors.length}`);
  console.log(`[PRD-CHECK] 告警问题: ${warnings.length}`);
  console.log(`[PRD-CHECK] 报告: ${path.relative(repoRoot, reportPath)}`);

  if (hardErrors.length > 0) {
    process.exitCode = 1;
  }
}

main();
