import { Secret, SecretFilterParams } from "./secrets.client.js";

/**
 * Security levels for secret classification and handling.
 */
export enum SecurityLevel {
  /** Public - no encryption needed */
  PUBLIC = "public",

  /** Internal - basic encryption */
  INTERNAL = "internal",

  /** Confidential - standard encryption + access control */
  CONFIDENTIAL = "confidential",

  /** Secret - enhanced encryption + strict access control */
  SECRET = "secret",

  /** Top Secret - maximum encryption + audit + rotation */
  TOP_SECRET = "top_secret",
}

/**
 * Secret categories for classification.
 */
export enum SecretCategory {
  /** API keys and tokens */
  API_KEY = "api_key",

  /** Database credentials */
  DATABASE = "database",

  /** Encryption keys */
  ENCRYPTION_KEY = "encryption_key",

  /** Certificates */
  CERTIFICATE = "certificate",

  /** Service credentials */
  SERVICE = "service",

  /** User credentials */
  USER_CREDENTIALS = "user_credentials",

  /** Configuration values */
  CONFIGURATION = "configuration",

  /** Temporary tokens */
  TEMPORARY = "temporary",
}

/**
 * Secret metadata for fine-grained management.
 */
export interface SecretMetadata {
  /** Security classification level */
  securityLevel: SecurityLevel;

  /** Secret category */
  category: SecretCategory;

  /** Owner identifier */
  ownerId: string;

  /** Access control list */
  acl?: string[];

  /** Required permissions for access */
  requiredPermissions?: string[];

  /** Geographic restrictions */
  geoRestrictions?: string[];

  /** Time-based access restrictions */
  timeRestrictions?: {
    startHour?: number;
    endHour?: number;
    daysOfWeek?: number[];
    timezone?: string;
  };

  /** IP address restrictions */
  ipRestrictions?: string[];

  /** Device restrictions */
  deviceRestrictions?: string[];

  /** Rotation policy */
  rotationPolicy: {
    enabled: boolean;
    intervalDays?: number;
    autoRotate?: boolean;
    notifyBeforeDays?: number;
  };

  /** Retention policy */
  retentionPolicy: {
    retainAfterDeletion: boolean;
    retainDays?: number;
    permanentArchive?: boolean;
  };

  /** Compliance requirements */
  compliance: {
    standards?: string[];
    auditRequired?: boolean;
    encryptionStandard?: string;
  };

  /** Custom classification tags */
  classificationTags?: string[];

  /** Risk score (0-100) */
  riskScore?: number;

  /** Business impact level */
  businessImpact?: "low" | "medium" | "high" | "critical";
}

/**
 * Enhanced secret interface with fine-grained metadata.
 */
export interface EnhancedSecret extends Omit<Secret, "metadata"> {
  /** Enhanced metadata with security classification */
  metadata: SecretMetadata & Record<string, unknown>;

  /** Secret version information */
  version: {
    current: number;
    total: number;
    lastRotated?: string;
    nextRotation?: string;
  };

  /** Access statistics */
  accessStats: {
    totalAccess: number;
    lastAccess?: string;
    uniqueUsers: number;
    failedAttempts: number;
  };

  /** Secret status */
  status: "active" | "suspended" | "deprecated" | "compromised";

  /** Secret checksum for integrity verification */
  checksum?: string;

  /** Encryption information */
  encryption: {
    algorithm: string;
    keyId: string;
    iv?: string;
    strength: number;
  };
}

/**
 * Secret creation request with enhanced metadata.
 */
export interface CreateEnhancedSecretRequest {
  /** Secret name/key */
  name: string;

  /** Secret value */
  value: string;

  /** Enhanced security metadata */
  metadata: SecretMetadata;

  /** Secret description */
  description?: string;

  /** Initial tags */
  tags?: string[];

  /** Expiration timestamp */
  expiresAt?: string;

  /** Whether to enable versioning */
  enableVersioning?: boolean;

  /** Initial access control settings */
  initialAccess?: {
    allowedUsers?: string[];
    allowedRoles?: string[];
    accessDuration?: number; // in hours
  };
}

/**
 * Secret update request with enhanced options.
 */
export interface UpdateEnhancedSecretRequest {
  /** New secret value */
  value?: string;

  /** Updated metadata */
  metadata?: Partial<SecretMetadata>;

  /** New description */
  description?: string;

  /** Updated tags */
  tags?: string[];

  /** New expiration */
  expiresAt?: string;

  /** Status change */
  status?: EnhancedSecret["status"];

  /** Force rotation */
  forceRotation?: boolean;

  /** Access control updates */
  accessControl?: {
    addUsers?: string[];
    removeUsers?: string[];
    addRoles?: string[];
    removeRoles?: string[];
  };
}

/**
 * Fine-grained filter parameters for enhanced secret search.
 */
export interface EnhancedSecretFilterParams extends SecretFilterParams {
  /** Filter by security level */
  securityLevel?: SecurityLevel[];

  /** Filter by category */
  category?: SecretCategory[];

  /** Filter by owner */
  ownerId?: string[];

  /** Filter by status */
  status?: EnhancedSecret["status"][];

  /** Filter by risk score range */
  riskScoreRange?: {
    min?: number;
    max?: number;
  };

  /** Filter by business impact */
  businessImpact?: ("low" | "medium" | "high" | "critical")[];

  /** Filter by compliance standards */
  complianceStandards?: string[];

  /** Filter by access frequency */
  accessFrequency?: {
    min?: number;
    max?: number;
  };

  /** Filter by rotation status */
  rotationStatus?: "enabled" | "disabled" | "overdue";

  /** Filter by encryption algorithm */
  encryptionAlgorithm?: string[];

  /** Include/exclude compromised secrets */
  excludeCompromised?: boolean;

  /** Filter by last access time */
  lastAccessAfter?: string;
  lastAccessBefore?: string;
}
