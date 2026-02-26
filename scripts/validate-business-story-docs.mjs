import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const scriptFile = fileURLToPath(import.meta.url);
const scriptDir = path.dirname(scriptFile);
const repoRoot = path.resolve(scriptDir, "..");

const files = {
  prd: path.join(repoRoot, "docs", "PRD.md"),
  comprehensive: path.join(repoRoot, "docs", "PRD_COMPREHENSIVE.md"),
  governance: path.join(repoRoot, "docs", "PRD_GOVERNANCE.md"),
};

const reportDir = path.join(repoRoot, "docs", "reports");
const reportPath = path.join(reportDir, "business-story-validation-report.md");

const hardErrors = [];
const warnings = [];

function readFileSafe(filePath) {
  if (!fs.existsSync(filePath)) {
    hardErrors.push(`缺少文件: ${path.relative(repoRoot, filePath)}`);
    return "";
  }

  return fs.readFileSync(filePath, "utf8");
}

function checkContains(text, regex, message, level = "hard") {
  if (regex.test(text)) {
    return;
  }

  if (level === "hard") {
    hardErrors.push(message);
    return;
  }

  warnings.push(message);
}

function checkNoMatch(text, regex, message, level = "hard") {
  if (!regex.test(text)) {
    return;
  }

  if (level === "hard") {
    hardErrors.push(message);
    return;
  }

  warnings.push(message);
}

function countStoryIds(text) {
  const matches = text.match(/US\d{3}/g);
  return matches ? matches : [];
}

function validateStoryRows(text) {
  const lines = text.split(/\r?\n/);
  const storyLines = lines.filter((line) => /\|\s*US\d{3}\s*\|/.test(line));

  for (const line of storyLines) {
    const cols = line
      .split("|")
      .map((item) => item.trim())
      .filter((item) => item.length > 0);

    if (cols.length < 3) {
      hardErrors.push(`用户故事行字段不足: ${line.trim()}`);
      continue;
    }

    const acceptance = cols[cols.length - 1];
    if (!acceptance || acceptance === "-") {
      hardErrors.push(`用户故事缺少验收标准: ${line.trim()}`);
    }
  }
}

function checkDuplicateStoryIds(ids) {
  const seen = new Set();
  const duplicated = new Set();

  for (const id of ids) {
    if (seen.has(id)) {
      duplicated.add(id);
    }

    seen.add(id);
  }

  if (duplicated.size > 0) {
    warnings.push(`用户故事编号重复: ${Array.from(duplicated).join(", ")}`);
  }
}

function main() {
  const prd = readFileSafe(files.prd);
  const comprehensive = readFileSafe(files.comprehensive);
  const governance = readFileSafe(files.governance);

  checkContains(prd, /^##\s*二、功能架构/m, "PRD.md 缺少“二、功能架构”章节");
  checkContains(comprehensive, /^##\s*四、主要业务流程/m, "PRD_COMPREHENSIVE.md 缺少“四、主要业务流程”章节");
  checkContains(comprehensive, /^##\s*七、用户故事与优先级/m, "PRD_COMPREHENSIVE.md 缺少“七、用户故事与优先级”章节");
  checkContains(governance, /单一事实源/, "PRD_GOVERNANCE.md 缺少单一事实源约束");

  const storyIds = countStoryIds(comprehensive);
  if (storyIds.length < 5) {
    hardErrors.push(`用户故事数量不足（当前 ${storyIds.length}，要求至少 5）`);
  }

  validateStoryRows(comprehensive);

  checkContains(
    comprehensive,
    /```mermaid|stateDiagram|flowchart/i,
    "PRD_COMPREHENSIVE.md 业务流程缺少流程图标记（mermaid/stateDiagram/flowchart）"
  );

  checkNoMatch(prd, /app-react/i, "PRD.md 存在过时路径 app-react");
  checkNoMatch(comprehensive, /app-react/i, "PRD_COMPREHENSIVE.md 存在过时路径 app-react");

  checkDuplicateStoryIds(storyIds);

  checkContains(comprehensive, /####\s*P0/i, "用户故事缺少 P0 优先级分段", "warn");
  checkContains(comprehensive, /####\s*P1/i, "用户故事缺少 P1 优先级分段", "warn");
  checkContains(comprehensive, /####\s*P2/i, "用户故事缺少 P2 优先级分段", "warn");
  checkContains(comprehensive, /Given|When|Then/i, "验收标准中缺少 Given/When/Then 结构化描述（建议）", "warn");

  fs.mkdirSync(reportDir, { recursive: true });

  const status = hardErrors.length === 0 ? "PASS" : "FAIL";
  const report = [
    "# 业务逻辑与用户故事文档校验报告",
    "",
    `- 时间: ${new Date().toISOString()}`,
    `- 状态: ${status}`,
    `- 阻断问题: ${hardErrors.length}`,
    `- 告警问题: ${warnings.length}`,
    `- 用户故事数量: ${storyIds.length}`,
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

  console.log(`[BIZ-STORY-CHECK] ${status}`);
  console.log(`[BIZ-STORY-CHECK] 阻断问题: ${hardErrors.length}`);
  console.log(`[BIZ-STORY-CHECK] 告警问题: ${warnings.length}`);
  console.log(`[BIZ-STORY-CHECK] 用户故事数量: ${storyIds.length}`);
  console.log(`[BIZ-STORY-CHECK] 报告: ${path.relative(repoRoot, reportPath)}`);

  if (hardErrors.length > 0) {
    process.exitCode = 1;
  }
}

main();
