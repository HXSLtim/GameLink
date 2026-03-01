import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const scriptFile = fileURLToPath(import.meta.url);
const scriptDir = path.dirname(scriptFile);
const repoRoot = path.resolve(scriptDir, "..");

const mvpReportPath = path.join(repoRoot, "docs", "reports", "mvp-three-end-release-report.md");
const engineeringReportPath = path.join(repoRoot, "docs", "reports", "engineering-acceptance-report.md");

function assertFileExists(filePath) {
  if (!fs.existsSync(filePath)) {
    throw new Error(`报告不存在: ${path.relative(repoRoot, filePath)}`);
  }
}

function assertLineContains(content, text, fileLabel) {
  if (!content.includes(text)) {
    throw new Error(`${fileLabel} 缺少关键行: ${text}`);
  }
}

function extractNumber(content, pattern, fallbackLabel) {
  const match = content.match(pattern);
  if (!match) {
    throw new Error(`无法解析 ${fallbackLabel}`);
  }

  return Number(match[1]);
}

function main() {
  assertFileExists(mvpReportPath);
  assertFileExists(engineeringReportPath);

  const mvpReport = fs.readFileSync(mvpReportPath, "utf8");
  const engineeringReport = fs.readFileSync(engineeringReportPath, "utf8");

  assertLineContains(mvpReport, "- 状态: PASS", "mvp-three-end-release-report");
  assertLineContains(engineeringReport, "- 状态: PASS", "engineering-acceptance-report");

  const mvpFailed = extractNumber(mvpReport, /- 失败步骤: (\d+)/, "mvp failed count");
  const engineeringFailed = extractNumber(engineeringReport, /- 失败步骤: (\d+)/, "engineering failed count");

  if (mvpFailed !== 0) {
    throw new Error(`三端发布报告存在失败步骤: ${mvpFailed}`);
  }

  if (engineeringFailed !== 0) {
    throw new Error(`工程验收报告存在失败步骤: ${engineeringFailed}`);
  }

  const mvpStepCount = extractNumber(mvpReport, /- 总步骤: (\d+)/, "mvp total step count");
  if (mvpStepCount < 5) {
    throw new Error(`三端发布总步骤异常（过少）: ${mvpStepCount}`);
  }

  console.log("[MVP-REPORT-VALIDATION] PASS");
  console.log(`[MVP-REPORT-VALIDATION] mvp report: ${path.relative(repoRoot, mvpReportPath)}`);
  console.log(`[MVP-REPORT-VALIDATION] engineering report: ${path.relative(repoRoot, engineeringReportPath)}`);
}

main();
