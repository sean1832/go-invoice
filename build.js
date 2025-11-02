#!/usr/bin/env node

/**
 * Cross-platform build script for go-invoice
 * Builds frontend, copies to backend, and builds Go binary
 *
 * Usage:
 *   node build.js           - Full build
 *   node build.js copy      - Copy frontend to backend only
 *   node build.js frontend  - Build frontend only
 *   node build.js backend   - Build backend only
 */

const { execSync } = require("child_process");
const fs = require("fs");
const path = require("path");

// Colors for terminal output
const colors = {
  reset: "\x1b[0m",
  bright: "\x1b[1m",
  green: "\x1b[32m",
  red: "\x1b[31m",
  yellow: "\x1b[33m",
  cyan: "\x1b[36m",
};

function log(message, color = colors.reset) {
  console.log(`${color}${message}${colors.reset}`);
}

function logStep(message) {
  log(`\n${colors.bright}${colors.cyan}► ${message}${colors.reset}`);
}

function logSuccess(message) {
  log(`${colors.green}✓ ${message}${colors.reset}`);
}

function logError(message) {
  log(`${colors.red}✗ ${message}${colors.reset}`);
}

/**
 * Execute a command and stream output
 */
function exec(command, cwd = process.cwd()) {
  try {
    execSync(command, {
      cwd,
      stdio: "inherit",
      shell: true,
    });
    return true;
  } catch (error) {
    return false;
  }
}

/**
 * Recursively copy a directory
 */
function copyDir(src, dest) {
  fs.mkdirSync(dest, { recursive: true });
  const entries = fs.readdirSync(src, { withFileTypes: true });

  for (let entry of entries) {
    const srcPath = path.join(src, entry.name);
    const destPath = path.join(dest, entry.name);

    if (entry.isDirectory()) {
      copyDir(srcPath, destPath);
    } else {
      fs.copyFileSync(srcPath, destPath);
    }
  }
}

/**
 * Build frontend
 */
function buildFrontend(frontendDir) {
  logStep("Building frontend...");
  if (!exec("npm run build", frontendDir)) {
    logError("Frontend build failed!");
    process.exit(1);
  }
  logSuccess("Frontend built successfully");
}

/**
 * Copy frontend to backend
 */
function copyFrontend(buildDir, distDir) {
  logStep("Copying frontend to backend...");
  try {
    if (!fs.existsSync(buildDir)) {
      logError(`Build directory not found: ${buildDir}`);
      process.exit(1);
    }

    if (fs.existsSync(distDir)) {
      fs.rmSync(distDir, { recursive: true, force: true });
    }

    copyDir(buildDir, distDir);
    logSuccess("Frontend copied to backend");
  } catch (error) {
    logError(`Copy failed: ${error.message}`);
    process.exit(1);
  }
}

/**
 * Build backend
 */
function buildBackend(backendDir) {
  logStep("Building Go backend...");
  const isWindows = process.platform === "win32";
  const binaryName = isWindows ? "go-invoice.exe" : "go-invoice";
  const buildCmd = `go build -o bin/${binaryName} .`;

  if (!exec(buildCmd, backendDir)) {
    logError("Go build failed!");
    process.exit(1);
  }
  logSuccess("Backend built successfully");

  const binaryPath = path.join(backendDir, "bin", binaryName);
  return binaryPath;
}

/**
 * Main build process
 */
function build() {
  const rootDir = __dirname;
  const frontendDir = path.join(rootDir, "frontend");
  const backendDir = path.join(rootDir, "backend");
  const buildDir = path.join(frontendDir, "build");
  const distDir = path.join(backendDir, "internal", "ui", "dist");

  const command = process.argv[2] || "all";

  log(`${colors.bright}${colors.green}Building go-invoice${colors.reset}\n`);

  switch (command) {
    case "frontend":
      buildFrontend(frontendDir);
      break;

    case "copy":
      copyFrontend(buildDir, distDir);
      break;

    case "backend":
      const binaryPath = buildBackend(backendDir);
      log(`\n${colors.bright}${colors.green}✓ Backend build complete!${colors.reset}`);
      log(`\nRun the application:`);
      log(`  ${colors.cyan}${binaryPath}${colors.reset}\n`);
      break;

    case "all":
    default:
      buildFrontend(frontendDir);
      copyFrontend(buildDir, distDir);
      const fullBinaryPath = buildBackend(backendDir);

      log(`\n${colors.bright}${colors.green}✓ Build complete!${colors.reset}`);
      log(`\nRun the application:`);
      log(`  ${colors.cyan}${fullBinaryPath}${colors.reset}\n`);
      break;
  }
}

// Run build
build();
