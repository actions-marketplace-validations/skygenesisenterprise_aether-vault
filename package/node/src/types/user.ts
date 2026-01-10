/**
 * User management types for Aether Vault SDK.
 * Based on server Go models from server/src/model/user.go and server/src/controllers/user.go
 */

/**
 * User interface representing a vault user.
 * Matches the Go User struct with proper TypeScript types.
 */
export interface User {
  /** Unique user identifier (UUID) */
  id: string;

  /** User email address (unique) */
  email: string;

  /** User first name */
  firstName: string;

  /** User last name */
  lastName: string;

  /** Whether the user account is active */
  isActive: boolean;

  /** Account creation timestamp */
  createdAt: string;

  /** Account last update timestamp */
  updatedAt: string;
}

/**
 * User creation request interface.
 * Used for creating new users via the API.
 */
export interface CreateUserRequest {
  /** User email address (must be unique) */
  email: string;

  /** User password (minimum 8 characters) */
  password: string;

  /** User first name */
  firstName: string;

  /** User last name */
  lastName: string;
}

/**
 * User update request interface.
 * Used for updating existing users via the API.
 */
export interface UpdateUserRequest {
  /** New first name (optional) */
  firstName?: string;

  /** New last name (optional) */
  lastName?: string;

  /** New active status (optional) */
  isActive?: boolean;
}

/**
 * User list response interface.
 * Paginated response containing multiple users.
 */
export interface UserListResponse {
  /** Array of users */
  users: User[];

  /** Total number of users */
  total: number;

  /** Current page number */
  page: number;

  /** Number of users per page */
  pageSize: number;

  /** Total number of pages */
  totalPages: number;
}

/**
 * User filter parameters interface.
 * Filtering and pagination options for user list requests.
 */
export interface UserFilterParams extends Record<string, unknown> {
  /** Page number (default: 1) */
  page?: number;

  /** Number of items per page (default: 20) */
  pageSize?: number;

  /** Sort field */
  sortBy?: string;

  /** Sort direction */
  sortOrder?: "asc" | "desc";

  /** Filter by active status */
  isActive?: boolean;

  /** Search in email, first name, or last name */
  search?: string;

  /** Filter by creation date from */
  createdFrom?: string;

  /** Filter by creation date to */
  createdTo?: string;
}

/**
 * User statistics interface.
 * Aggregated statistics about users.
 */
export interface UserStats {
  /** Total number of users */
  totalUsers: number;

  /** Number of active users */
  activeUsers: number;

  /** Number of inactive users */
  inactiveUsers: number;

  /** Users created in the last 24 hours */
  usersCreated24h: number;

  /** Users created in the last 7 days */
  usersCreated7d: number;

  /** Users created in the last 30 days */
  usersCreated30d: number;
}
