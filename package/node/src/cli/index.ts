#!/usr/bin/env node

import { initCommand } from "./commands/init.js";

const args = process.argv.slice(2);
const command = args[0];

if (command === "init") {
  const options = {
    force: args.includes("--force") || args.includes("-f"),
    template: "basic",
  };

  const templateIndex = args.findIndex(
    (arg) => arg === "--template" || arg === "-t",
  );
  if (templateIndex !== -1 && args[templateIndex + 1]) {
    options.template = args[templateIndex + 1];
  }

  initCommand(options);
} else if (command === "--version" || command === "-v") {
  console.log("1.0.8");
} else if (command === "--help" || command === "-h") {
  console.log(`
Aether Vault CLI - Initialize and manage vault configuration

Usage:
  aether-vault <command>

Commands:
  init              Initialize a new vault.config.ts file

Options:
  --version, -v     Show version number
  --help, -h        Show this help message

Init Options:
  --force, -f       Overwrite existing configuration file
  --template, -t    Use specific template (basic, production, development)
`);
} else {
  console.error(`Unknown command: ${command}`);
  console.error("Use --help for available commands");
  process.exit(1);
}
