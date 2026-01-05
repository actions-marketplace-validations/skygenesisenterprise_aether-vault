/**
 * Secret type enumeration.
 * Type definitions for different kinds of secrets.
 */
export type SecretType =
  | "api_key"
  | "password"
  | "certificate"
  | "ssh_key"
  | "database"
  | "token"
  | "oauth"
  | "custom";

/**
 * Secret creation input interface.
 * Data required to create a new secret.
 */
export interface SecretInput {
  /** Secret name */
  name: string;

  /** Secret description */
  description?: string;

  /** Secret value */
  value: string;

  /** Secret type */
  type: SecretType;

  /** Tags for categorization */
  tags?: string;

  /** Optional expiration timestamp */
  expiresAt?: Date;
}

/**
 * Secret update input interface.
 * Data for updating an existing secret.
 */
export interface SecretUpdate {
  /** Updated secret name */
  name?: string;

  /** Updated secret description */
  description?: string;

  /** Updated secret value */
  value?: string;

  /** Updated secret type */
  type?: SecretType;

  /** Updated tags */
  tags?: string;

  /** Updated expiration timestamp */
  expiresAt?: Date;

  /** Updated active status */
  isActive?: boolean;
}

/**
 * Secret interface.
 * Complete secret data model.
 */
export interface Secret {
  /** Secret unique identifier */
  id: string;

  /** Secret name */
  name: string;

  /** Secret description */
  description: string;

  /** Secret value (decrypted) */
  value: string;

  /** Secret type */
  type: SecretType;

  /** Secret tags */
  tags: string;

  /** Whether secret is active */
  isActive: boolean;

  /** Secret creation timestamp */
  createdAt: Date;

  /** Secret last update timestamp */
  updatedAt: Date;

  /** Optional expiration timestamp */
  expiresAt?: Date;

  /** User ID who owns the secret */
  userId: string;

  /** Version number for optimistic locking */
  version: number;
}

/**
 * Secret list response interface.
 * Paginated list of secrets.
 */
export interface SecretListResponse {
  /** Array of secrets */
  secrets: Secret[];

  /** Total number of secrets */
  total: number;

  /** Current page number */
  page: number;

  /** Number of secrets per page */
  pageSize: number;

  /** Total number of pages */
  totalPages: number;
}

/**
 * Authentication credentials interface.
 * Login request data.
 */
export interface AuthCredentials {
  /** User email */
  username: string;

  /** User password */
  password: string;
}

/**
 * Authentication session interface.
 * Session information after successful login.
 */
export interface AuthSession {
  /** JWT access token */
  token: string;

  /** Token expiration timestamp */
  expiresAt: Date;

  /** User information */
  user: UserIdentity;

  /** Token type */
  tokenType: string;

  /** Refresh token (if available) */
  refreshToken?: string;
}

/**
 * User identity interface.
 * Current user information.
 */
export interface UserIdentity {
  /** User unique identifier */
  id: string;

  /** User email */
  email: string;

  /** User first name */
  firstName: string;

  /** User last name */
  lastName: string;

  /** User display name */
  displayName?: string;

  /** User avatar URL */
  avatar?: string;

  /** Account creation timestamp */
  createdAt: Date;

  /** Account last update timestamp */
  updatedAt: Date;

  /** Whether user account is active */
  isActive: boolean;

  /** User roles */
  roles: string[];

  /** User metadata */
  metadata?: Record<string, unknown>;
}

/**
 * Policy interface.
 * Access policy definition.
 */
export interface Policy {
  /** Policy unique identifier */
  id: string;

  /** Policy name */
  name: string;

  /** Policy description */
  description: string;

  /** Policy resource */
  resource: string;

  /** Policy actions */
  actions: string[];

  /** Policy effect */
  effect: "allow" | "deny";

  /** Policy conditions */
  conditions?: Record<string, unknown>;

  /** Policy priority */
  priority: number;

  /** Whether policy is active */
  isActive: boolean;

  /** Policy creation timestamp */
  createdAt: Date;

  /** Policy last update timestamp */
  updatedAt: Date;

  /** User ID who owns the policy */
  userId: string;
}

/**
 * Policy list response interface.
 * List of user policies.
 */
export interface PolicyListResponse {
  /** Array of policies */
  policies: Policy[];

  /** Total number of policies */
  total: number;
}

/**
 * TOTP entry interface.
 * TOTP configuration for 2FA.
 */
export interface TOTPEntry {
  /** TOTP unique identifier */
  id: string;

  /** TOTP name */
  name: string;

  /** TOTP description */
  description?: string;

  /** Whether TOTP is active */
  isActive: boolean;

  /** TOTP creation timestamp */
  createdAt: Date;

  /** TOTP last update timestamp */
  updatedAt: Date;

  /** User ID who owns the TOTP */
  userId: string;

  /** TOTP metadata */
  metadata?: Record<string, unknown>;
}

/**
 * TOTP creation input interface.
 * Data for creating a new TOTP entry.
 */
export interface TOTPCreateInput {
  /** TOTP name */
  name: string;

  /** TOTP description */
  description?: string;

  /** TOTP algorithm (default: SHA1) */
  algorithm?: string;

  /** Number of digits (default: 6) */
  digits?: number;

  /** Time period in seconds (default: 30) */
  period?: number;
}

/**
 * TOTP code interface.
 * Generated TOTP code information.
 */
export interface TOTPCode {
  /** Generated 6-digit code */
  code: string;

  /** Code expiration timestamp */
  expiresAt: Date;

  /** Remaining time in seconds */
  remainingSeconds: number;
}

/**
 * TOTP list response interface.
 * List of user TOTP entries.
 */
export interface TOTPListResponse {
  /** Array of TOTP entries */
  entries: TOTPEntry[];

  /** Total number of entries */
  total: number;
}

/**
 * Audit entry interface.
 * Audit log entry for tracking actions.
 */
export interface AuditEntry {
  /** Audit entry unique identifier */
  id: string;

  /** Action performed */
  action: string;

  /** Resource type */
  resource: string;

  /** Resource identifier */
  resourceId?: string;

  /** User ID who performed action */
  userId?: string;

  /** IP address of request */
  ipAddress: string;

  /** User agent string */
  userAgent: string;

  /** Whether action was successful */
  success: boolean;

  /** Action details */
  details?: Record<string, unknown>;

  /** Action timestamp */
  createdAt: Date;

  /** Request duration in milliseconds */
  duration?: number;

  /** HTTP status code */
  statusCode?: number;
}

/**
 * Audit filter interface.
 * Filtering options for audit logs.
 */
export interface AuditFilter {
  /** Filter by action */
  action?: string;

  /** Filter by resource type */
  resource?: string;

  /** Filter by resource ID */
  resourceId?: string;

  /** Filter by user ID */
  userId?: string;

  /** Filter by success status */
  success?: boolean;

  /** Filter by date range */
  dateFrom?: Date;
  dateTo?: Date;

  /** Pagination parameters */
  page?: number;
  pageSize?: number;

  /** Sort options */
  sortBy?: string;
  sortOrder?: "asc" | "desc";
}

/**
 * Audit list response interface.
 * Paginated list of audit entries.
 */
export interface AuditListResponse {
  /** Array of audit entries */
  entries: AuditEntry[];

  /** Total number of entries */
  total: number;

  /** Current page number */
  page: number;

  /** Number of entries per page */
  pageSize: number;

  /** Total number of pages */
  totalPages: number;
}

/**
 * Health check response interface.
 * System health status.
 */
export interface HealthResponse {
  /** Health status */
  status: "healthy" | "unhealthy";

  /** Current timestamp */
  timestamp: Date;

  /** Service version */
  version: string;

  /** Database connection status */
  database: string;

  /** Additional health details */
  details?: Record<string, unknown>;
}

/**
 * Version response interface.
 * System version information.
 */
export interface VersionResponse {
  /** Current version */
  version: string;

  /** Build timestamp */
  buildTime: string;

  /** Git commit hash */
  gitCommit: string;

  /** Go version */
  goVersion: string;
}

/**
 * Error response interface.
 * Standard error format from API.
 */
export interface ErrorResponse {
  /** Error details */
  error: {
    /** Error code */
    code: string;

    /** Error message */
    message: string;

    /** Additional error details */
    details?: Record<string, unknown>;
  };
}

/**
 * Success response interface.
 * Standard success response format.
 */
export interface SuccessResponse<T = unknown> {
  /** Response data */
  data: T;

  /** Success message */
  message?: string;

  /** Response metadata */
  meta?: Record<string, unknown>;
}
