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
function buildBackend(backendDir, version = "dev", targetOS = null, targetArch = "amd64") {
  // Normalize Inputs
  // If targetOS is not provided, use the host OS (development mode)
  // Map friendly names (macos) to Go standards (darwin)
  const hostOS =
    process.platform === "win32" ? "windows" : process.platform === "darwin" ? "darwin" : "linux";
  let goos = targetOS ? targetOS.toLowerCase() : hostOS;

  if (goos === "macos") goos = "darwin";
  if (goos === "win") goos = "windows";

  // Define Output Name
  // Windows binaries require .exe extension
  const ext = goos === "windows" ? ".exe" : "";
  const binaryName = `go-invoice-${goos}-${targetArch}${ext}`;
  const outputPath = path.join(backendDir, "bin", binaryName);

  logStep(`Building Backend [${version}] for ${goos}/${targetArch}...`);

  // Configure Build Environment
  // CGO_ENABLED=0 is critical for portability (Alpine support) and easy cross-compilation
  const buildEnv = {
    ...process.env,
    CGO_ENABLED: "0",
    GOOS: goos,
    GOARCH: targetArch,
  };

  // Linker Flags
  // -s -w: Strip debug symbols for non-dev builds to reduce size
  let ldflags = `-X 'main.Version=${version}'`;
  if (version !== "dev") {
    ldflags += " -s -w";
  }

  const buildCmd = `go build -ldflags "${ldflags}" -o "${outputPath}" .`;

  try {
    execSync(buildCmd, {
      cwd: backendDir,
      env: buildEnv,
      stdio: "inherit", // Stream compiler output to terminal
    });
    logSuccess(`Artifact created: bin/${binaryName}`);
  } catch (error) {
    logError(`Build failed for ${goos}/${targetArch}`);
    process.exit(1);
  }

  return outputPath;
}

/**
 * Get version from package.json or environment variable
 */
function getVersion(rootDir) {
  // Check for VERSION environment variable (used in CI/CD)
  if (process.env.VERSION) {
    return process.env.VERSION;
  }

  // Read from root package.json
  const packageJsonPath = path.join(rootDir, "package.json");
  try {
    const packageJson = JSON.parse(fs.readFileSync(packageJsonPath, "utf8"));
    return packageJson.version || "dev";
  } catch (error) {
    logError(`Failed to read version from package.json: ${error.message}`);
    return "dev";
  }
}

/**
 * Helper to parse named flags from process.argv
 * e.g., --platform=linux becomes { platform: "linux" }
 */
function parseFlags() {
  const flags = {
    // Defaults: null indicates "detect host environment" later
    platform: null,
    arch: "amd64",
    version: null,
  };

  // Skip "node" and "build.js"
  const args = process.argv.slice(2);
  const command = args[0] && !args[0].startsWith("-") ? args[0] : "all";

  args.forEach((arg) => {
    if (arg.startsWith("--platform=")) flags.platform = arg.split("=")[1];
    if (arg.startsWith("--arch=")) flags.arch = arg.split("=")[1];
    if (arg.startsWith("--version=")) flags.version = arg.split("=")[1];
  });

  return { command, ...flags };
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

  // Parse CLI flags
  const { command, platform, arch, version: flagVersion } = parseFlags();

  const version = flagVersion || getVersion(rootDir);

  log(`${colors.bright}${colors.green}Building go-invoice [${command}]${colors.reset}\n`);

  switch (command) {
    case "frontend":
      buildFrontend(frontendDir);
      break;

    case "copy":
      copyFrontend(buildDir, distDir);
      break;

    case "backend":
      const binaryPath = buildBackend(backendDir, version, platform, arch);
      log(`\n${colors.bright}${colors.green}✓ Backend build complete!${colors.reset}`);
      log(`\nRun the application:`);
      log(`  ${colors.cyan}${binaryPath}${colors.reset}\n`);
      break;

    case "all":
    default:
      buildFrontend(frontendDir);
      copyFrontend(buildDir, distDir);
      const fullBinaryPath = buildBackend(backendDir, version);

      log(`\n${colors.bright}${colors.green}✓ Build complete!${colors.reset}`);
      log(`\nRun the application:`);
      log(`  ${colors.cyan}${fullBinaryPath}${colors.reset}\n`);
      break;
  }
}

// Run build
build();
