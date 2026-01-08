/**
 * Export all enhanced secrets management modules
 */

// Core interfaces and types
export * from "./security-levels.js";

// Classification engine
export * from "./classifier.js";

// Access control and policy engine
export * from "./policy-engine.js";

// Multi-level encryption
export * from "./encryption.js";

// Enhanced client with fine-grained management
export * from "./enhanced-client.js";

// Legacy client for backward compatibility
export { SecretsClient } from "./secrets.client.js";
