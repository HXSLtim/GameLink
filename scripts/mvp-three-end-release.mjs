import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { spawnSync } from "node:child_process";

const scriptFile = fileURLToPath(import.meta.url);
const scriptDir = path.dirname(scriptFile);
const repoRoot = path.resolve(scriptDir, "..");
const appDir = path.join(repoRoot, "app");
const adminDir = path.join(repoRoot, "admin");
const apiDir = path.join(repoRoot, "api");
const reportDir = path.join(repoRoot, "docs", "reports");
const reportPath = path.join(reportDir, "mvp-three-end-release-report.md");

function hasFlag(flag) {
  return process.argv.includes(flag);
}

function readOption(prefix, fallback) {
  const hit = process.argv.find((arg) => arg.startsWith(`${prefix}=`));
  if (!hit) {
    return fallback;
  }
  return hit.slice(prefix.length + 1).trim();
}

function resolvePowerShell() {
  const candidates = process.platform === "win32" ? ["pwsh", "powershell"] : ["pwsh", "powershell"];
  for (const command of candidates) {
    const probe = spawnSync(command, ["-NoLogo", "-NoProfile", "-Command", "$PSVersionTable.PSVersion.ToString()"], {
      encoding: "utf8",
      shell: false,
    });
    if (probe.status === 0) {
      return command;
    }
  }
  return null;
}

const options = {
  installDeps: hasFlag("--install-deps"),
  withFlowGuard: hasFlag("--with-flow-guard"),
  withWithdrawRegression: hasFlag("--with-withdraw-regression"),
  withFullAcceptance: hasFlag("--with-full-acceptance"),
  strictRelease: hasFlag("--strict-release"),
  dryRun: hasFlag("--dry-run"),
  baseUrl: readOption("--base-url", "http://127.0.0.1:8080/api/v1"),
};

function validateOptions() {
  if (!options.strictRelease) {
    return;
  }

  const requiredFlags = [];
  if (!options.withFlowGuard) {
    requiredFlags.push("--with-flow-guard");
  }
  if (!options.withWithdrawRegression) {
    requiredFlags.push("--with-withdraw-regression");
  }

  if (requiredFlags.length > 0) {
    throw new Error(`strict release requires flags: ${requiredFlags.join(", ")}`);
  }
}

validateOptions();

const powershell = resolvePowerShell();

const steps = [];

if (options.installDeps) {
  steps.push(
    {
      id: "P1",
      gate: "Preflight",
      name: "Install app dependencies",
      cwd: appDir,
      command: "npm",
      args: ["ci"],
    },
    {
      id: "P2",
      gate: "Preflight",
      name: "Install admin dependencies",
      cwd: adminDir,
      command: "npm",
      args: ["ci"],
    },
    {
      id: "P3",
      gate: "Preflight",
      name: "Download api dependencies",
      cwd: apiDir,
      command: "go",
      args: ["mod", "download"],
    }
  );
}

steps.push(
  {
    id: "M1",
    gate: "MVP Core",
    name: "Engineering acceptance",
    cwd: repoRoot,
    command: "node",
    args: ["scripts/engineering-acceptance.mjs"],
  },
  {
    id: "M2",
    gate: "MVP Core",
    name: "Admin type-check",
    cwd: adminDir,
    command: "npm",
    args: ["run", "type-check"],
  },
  {
    id: "M3",
    gate: "MVP Core",
    name: "Admin lint",
    cwd: adminDir,
    command: "npm",
    args: ["run", "lint"],
  },
  {
    id: "M4",
    gate: "MVP Core",
    name: "Admin unit tests",
    cwd: adminDir,
    command: "npm",
    args: ["run", "test:run"],
  },
  {
    id: "M5",
    gate: "MVP Core",
    name: "Admin build",
    cwd: adminDir,
    command: "npm",
    args: ["run", "build"],
  }
);

function addRegressionStep(id, name, scriptName) {
  if (!powershell) {
    steps.push({
      id,
      gate: "Regression",
      name,
      cwd: repoRoot,
      command: "node",
      args: ["-e", `throw new Error('PowerShell not found, cannot run ${scriptName}')`],
    });
    return;
  }

  steps.push({
    id,
    gate: "Regression",
    name,
    cwd: repoRoot,
    command: powershell,
    args: ["-ExecutionPolicy", "Bypass", "-File", path.join(apiDir, "scripts", scriptName), "-BaseUrl", options.baseUrl],
  });
}

if (options.withFlowGuard) {
  addRegressionStep("R1", "Flow guard regression", "run_flow_guard_regression.ps1");
}

if (options.withWithdrawRegression) {
  addRegressionStep("R2", "Withdraw regression", "run_withdraw_flow_regression.ps1");
}

if (options.withFullAcceptance) {
  addRegressionStep("R3", "Full service acceptance", "run_full_service_flow_acceptance.ps1");
}

function runStep(step) {
  const startedAt = Date.now();
  const shouldUseShell = process.platform === "win32" && step.command === "npm";
  const result = shouldUseShell
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
    "# 三端 MVP 发布验收报告",
    "",
    `- 时间: ${new Date().toISOString()}`,
    `- 状态: ${status}`,
    `- 总步骤: ${results.length}`,
    `- 失败步骤: ${failed.length}`,
    "",
    "## 执行参数",
    `- installDeps: ${options.installDeps}`,
    `- withFlowGuard: ${options.withFlowGuard}`,
    `- withWithdrawRegression: ${options.withWithdrawRegression}`,
    `- withFullAcceptance: ${options.withFullAcceptance}`,
    `- strictRelease: ${options.strictRelease}`,
    `- baseUrl: ${options.baseUrl}`,
    "",
    "## 结果总览",
    ...results.map((item) => `- [${item.passed ? "PASS" : "FAIL"}] ${item.id} ${item.gate} / ${item.name} (${item.durationMs}ms)`),
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
        lines.push("```text", detail.slice(0, 12000), "```");
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
  if (options.dryRun) {
    console.log("[MVP-RELEASE] DRY RUN");
    console.log(
      `[MVP-RELEASE] Steps: ${steps
        .map((step) => `${step.id}:${step.name}`)
        .join(" | ")}`
    );
    console.log("[MVP-RELEASE] PASS");
    return;
  }

  const results = [];

  for (const step of steps) {
    console.log(`[MVP-RELEASE] Running ${step.id} ${step.name}...`);
    const result = runStep(step);
    results.push(result);
    console.log(`[MVP-RELEASE] ${step.id} ${result.passed ? "PASS" : "FAIL"} (${result.durationMs}ms)`);
    if (!result.passed) {
      break;
    }
  }

  const report = renderReport(results);
  fs.mkdirSync(reportDir, { recursive: true });
  fs.writeFileSync(reportPath, report.content, "utf8");

  const failedCount = results.filter((item) => !item.passed).length;
  console.log(`[MVP-RELEASE] ${report.status}`);
  console.log(`[MVP-RELEASE] 失败步骤: ${failedCount}`);
  console.log(`[MVP-RELEASE] 报告: ${path.relative(repoRoot, reportPath)}`);

  if (report.status !== "PASS") {
    process.exitCode = 1;
  }
}

main();
