import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { spawnSync } from "node:child_process";

const scriptFile = fileURLToPath(import.meta.url);
const scriptDir = path.dirname(scriptFile);
const repoRoot = path.resolve(scriptDir, "..");
const appDir = path.join(repoRoot, "app");
const apiDir = path.join(repoRoot, "api");
const reportDir = path.join(repoRoot, "docs", "reports");
const reportPath = path.join(reportDir, "engineering-acceptance-report.md");

const steps = [
  {
    id: "A1",
    gate: "Gate A Docs",
    name: "PRD validation",
    cwd: appDir,
    command: "npm",
    args: ["run", "validate:prd"],
  },
  {
    id: "A2",
    gate: "Gate A Docs",
    name: "Business story validation",
    cwd: appDir,
    command: "npm",
    args: ["run", "validate:bizstory"],
  },
  {
    id: "B1",
    gate: "Gate B Frontend",
    name: "Type check",
    cwd: appDir,
    command: "npm",
    args: ["run", "type-check"],
  },
  {
    id: "B2",
    gate: "Gate B Frontend",
    name: "Lint",
    cwd: appDir,
    command: "npm",
    args: ["run", "lint"],
  },
  {
    id: "B3",
    gate: "Gate B Frontend",
    name: "Unit tests",
    cwd: appDir,
    command: "npm",
    args: ["run", "test:run"],
  },
  {
    id: "B4",
    gate: "Gate B Frontend",
    name: "Build",
    cwd: appDir,
    command: "npm",
    args: ["run", "build"],
  },
  {
    id: "C1",
    gate: "Gate C Backend",
    name: "Go test all",
    cwd: apiDir,
    command: "go",
    args: ["test", "./..."],
  },
];

function runStep(step) {
  const startedAt = Date.now();
  const useShellCommand = process.platform === "win32" && step.command === "npm";
  const result = useShellCommand
    ? spawnSync(`${step.command} ${step.args.join(" ")}`, {
        cwd: step.cwd,
        encoding: "utf8",
        shell: true,
        env: process.env,
        maxBuffer: 20 * 1024 * 1024,
      })
    : spawnSync(step.command, step.args, {
        cwd: step.cwd,
        encoding: "utf8",
        shell: false,
        env: process.env,
        maxBuffer: 20 * 1024 * 1024,
      });

  return {
    ...step,
    passed: result.status === 0,
    status: result.status,
    durationMs: Date.now() - startedAt,
    stdout: result.stdout ?? "",
    stderr: result.stderr ?? "",
  };
}

function renderReport(results) {
  const failed = results.filter((item) => !item.passed);
  const status = failed.length === 0 ? "PASS" : "FAIL";

  const lines = [
    "# 工程化验收报告",
    "",
    `- 时间: ${new Date().toISOString()}`,
    `- 状态: ${status}`,
    `- 总步骤: ${results.length}`,
    `- 失败步骤: ${failed.length}`,
    "",
    "## 结果总览",
    ...results.map(
      (item) =>
        `- [${item.passed ? "PASS" : "FAIL"}] ${item.id} ${item.gate} / ${item.name} (${item.durationMs}ms)`
    ),
    "",
    "## 失败详情",
  ];

  if (failed.length === 0) {
    lines.push("- 无");
  } else {
    for (const item of failed) {
      lines.push(`- ${item.id} ${item.name} (exit=${item.status ?? "null"})`);
      const detail = `${item.stdout}\n${item.stderr}`.trim();
      if (detail) {
        lines.push("```text", detail.slice(0, 8000), "```");
      }
    }
  }

  lines.push("", "## 执行命令");
  for (const item of results) {
    lines.push(`- ${item.id}: ${item.command} ${item.args.join(" ")} (cwd=${path.relative(repoRoot, item.cwd)})`);
  }

  lines.push("");
  return { status, content: lines.join("\n") };
}

function main() {
  const results = [];

  for (const step of steps) {
    console.log(`[ACCEPT] Running ${step.id} ${step.name}...`);
    const result = runStep(step);
    results.push(result);
    console.log(`[ACCEPT] ${step.id} ${result.passed ? "PASS" : "FAIL"} (${result.durationMs}ms)`);
  }

  const report = renderReport(results);
  fs.mkdirSync(reportDir, { recursive: true });
  fs.writeFileSync(reportPath, report.content, "utf8");

  const failedCount = results.filter((item) => !item.passed).length;
  console.log(`[ACCEPT] ${report.status}`);
  console.log(`[ACCEPT] 失败步骤: ${failedCount}`);
  console.log(`[ACCEPT] 报告: ${path.relative(repoRoot, reportPath)}`);

  if (report.status !== "PASS") {
    process.exitCode = 1;
  }
}

main();
