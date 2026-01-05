import { VaultClient } from "../core/client.js";
import { AuditEntry, AuditFilter, AuditListResponse } from "../types/index.js";

/**
 * Client for audit log operations.
 * Provides methods to access and filter audit logs from the Aether Vault API.
 */
export class AuditClient {
  /**
   * Creates a new AuditClient instance.
   *
   * @param client - VaultClient instance for HTTP requests
   */
  constructor(private readonly client: VaultClient) {}

  /**
   * Retrieves audit log entries with optional filtering.
   *
   * @param filter - Optional filter parameters
   * @returns Promise resolving to paginated audit entries
   *
   * @example
   * ```typescript
   * // Get all audit entries
   * const audit = await vault.audit.list();
   *
   * // Filter by specific action
   * const loginAttempts = await vault.audit.list({
   *   action: "login"
   * });
   *
   * // Filter by date range and user
   * const userActions = await vault.audit.list({
   *   userId: "user-123",
   *   dateFrom: new Date("2025-01-01"),
   *   dateTo: new Date("2025-01-31"),
   *   pageSize: 50
   * });
   * ```
   */
  public async list(filter?: AuditFilter): Promise<AuditListResponse> {
    // Convert date objects to ISO strings for API
    const params: Record<string, unknown> = {};

    if (filter) {
      Object.assign(params, {
        action: filter.action,
        resource: filter.resource,
        resource_id: filter.resourceId,
        user_id: filter.userId,
        success: filter.success,
        date_from: filter.dateFrom?.toISOString(),
        date_to: filter.dateTo?.toISOString(),
        page: filter.page,
        page_size: filter.pageSize,
        sort_by: filter.sortBy,
        sort_order: filter.sortOrder,
      });
    }

    const response = await this.client.get<AuditListResponse>(
      "/api/v1/audit/logs",
      params,
    );

    return response;
  }

  /**
   * Retrieves a specific audit entry by ID.
   *
   * @param id - Audit entry ID
   * @returns Promise resolving to audit entry
   *
   * @example
   * ```typescript
   * const entry = await vault.audit.getEntry("audit-123");
   * console.log(entry.action, entry.resource, entry.success);
   * ```
   */
  public async getEntry(id: string): Promise<AuditEntry> {
    const response = await this.client.get<AuditEntry>(
      `/api/v1/audit/logs/${id}`,
    );

    return response;
  }

  /**
   * Retrieves audit entries for a specific user.
   *
   * @param userId - User ID to filter by
   * @param options - Additional filtering options
   * @returns Promise resolving to paginated audit entries
   *
   * @example
   * ```typescript
   * const userAudit = await vault.audit.getUserEntries("user-123", {
   *   page: 1,
   *   pageSize: 25,
   *   dateFrom: new Date("2025-01-01")
   * });
   * ```
   */
  public async getUserEntries(
    userId: string,
    options?: Omit<AuditFilter, "userId">,
  ): Promise<AuditListResponse> {
    return this.list({
      ...options,
      userId,
    });
  }

  /**
   * Retrieves audit entries for a specific resource.
   *
   * @param resource - Resource type (e.g., "secret", "user")
   * @param resourceId - Optional specific resource ID
   * @param options - Additional filtering options
   * @returns Promise resolving to paginated audit entries
   *
   * @example
   * ```typescript
   * // Get all secret access logs
   * const secretLogs = await vault.audit.getResourceEntries("secret");
   *
   * // Get logs for specific secret
   * const specificLogs = await vault.audit.getResourceEntries("secret", "secret-123");
   * ```
   */
  public async getResourceEntries(
    resource: string,
    resourceId?: string,
    options?: Omit<AuditFilter, "resource" | "resourceId">,
  ): Promise<AuditListResponse> {
    const filter: AuditFilter = {
      ...options,
      resource,
    };

    if (resourceId !== undefined) {
      filter.resourceId = resourceId;
    }

    return this.list(filter);
  }

  /**
   * Retrieves failed authentication attempts.
   *
   * @param options - Additional filtering options
   * @returns Promise resolving to paginated audit entries
   *
   * @example
   * ```typescript
   * const failedLogins = await vault.audit.getFailedAuth({
   *   dateFrom: new Date(Date.now() - 24 * 60 * 60 * 1000), // Last 24h
   *   pageSize: 100
   * });
   * ```
   */
  public async getFailedAuth(
    options?: Omit<AuditFilter, "action" | "success">,
  ): Promise<AuditListResponse> {
    return this.list({
      ...options,
      action: "login",
      success: false,
    });
  }

  /**
   * Retrieves secret access logs.
   *
   * @param options - Additional filtering options
   * @returns Promise resolving to paginated audit entries
   *
   * @example
   * ```typescript
   * const secretAccess = await vault.audit.getSecretAccess({
   *   dateFrom: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000), // Last 7 days
   *   sortBy: "createdAt",
   *   sortOrder: "desc"
   * });
   * ```
   */
  public async getSecretAccess(
    options?: Omit<AuditFilter, "resource">,
  ): Promise<AuditListResponse> {
    return this.list({
      ...options,
      resource: "secret",
    });
  }

  /**
   * Retrieves system-level audit logs.
   *
   * @param options - Additional filtering options
   * @returns Promise resolving to paginated audit entries
   *
   * @example
   * ```typescript
   * const systemLogs = await vault.audit.getSystemLogs({
   *   dateFrom: new Date(Date.now() - 24 * 60 * 60 * 1000),
   *   pageSize: 50
   * });
   * ```
   */
  public async getSystemLogs(
    options?: Omit<AuditFilter, "resource">,
  ): Promise<AuditListResponse> {
    return this.list({
      ...options,
      resource: "system",
    });
  }

  /**
   * Exports audit logs to CSV format.
   *
   * @param filter - Optional filter parameters
   * @returns Promise resolving to CSV string
   *
   * @example
   * ```typescript
   * const csvData = await vault.audit.exportToCSV({
   *   dateFrom: new Date("2025-01-01"),
   *   dateTo: new Date("2025-01-31")
   * });
   *
   * // Save to file or process
   * const blob = new Blob([csvData], { type: 'text/csv' });
   * const url = URL.createObjectURL(blob);
   * ```
   */
  public async exportToCSV(filter?: AuditFilter): Promise<string> {
    const params: Record<string, unknown> = {};

    if (filter) {
      Object.assign(params, {
        action: filter.action,
        resource: filter.resource,
        resource_id: filter.resourceId,
        user_id: filter.userId,
        success: filter.success,
        date_from: filter.dateFrom?.toISOString(),
        date_to: filter.dateTo?.toISOString(),
        format: "csv",
      });
    } else {
      params.format = "csv";
    }

    const response = await this.client.get<string>(
      "/api/v1/audit/export",
      params,
    );

    return response;
  }
}
