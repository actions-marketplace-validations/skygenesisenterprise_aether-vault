import { VaultClient } from "../core/client.js";
import { HealthResponse, VersionResponse } from "../types/index.js";

/**
 * System client for Aether Vault API.
 * Provides health checks, version information, and system status operations.
 */
export class SystemClient {
  /**
   * Creates a new SystemClient instance.
   *
   * @param client - VaultClient instance for HTTP requests
   */
  constructor(private readonly client: VaultClient) {}

  /**
   * Checks system health status.
   *
   * @returns Promise resolving to health information
   *
   * @example
   * ```typescript
   * const health = await vault.system.health();
   * console.log("System status:", health.status);
   * console.log("Database status:", health.database);
   * console.log("Version:", health.version);
   * ```
   */
  public async health(): Promise<HealthResponse> {
    const response = await this.client.get<{
      status: string;
      timestamp: string;
      version: string;
      database: string;
    }>("/health");

    return {
      status: response.status as "healthy" | "unhealthy",
      timestamp: new Date(response.timestamp),
      version: response.version,
      database: response.database,
    };
  }

  /**
   * Gets system version information.
   *
   * @returns Promise resolving to version details
   *
   * @example
   * ```typescript
   * const version = await vault.system.version();
   * console.log("Aether Vault version:", version.version);
   * console.log("Build time:", version.buildTime);
   * console.log("Git commit:", version.gitCommit);
   * console.log("Go version:", version.goVersion);
   * ```
   */
  public async version(): Promise<VersionResponse> {
    return this.client.get<VersionResponse>("/api/v1/version");
  }

  /**
   * Checks if the system is ready to serve requests.
   *
   * @returns Promise resolving to readiness status
   *
   * @example
   * ```typescript
   * const isReady = await vault.system.ready();
   * if (isReady) {
   *   console.log("System is ready");
   * } else {
   *   console.log("System is still initializing");
   * }
   * ```
   */
  public async ready(): Promise<boolean> {
    try {
      const health = await this.health();
      return health.status === "healthy";
    } catch {
      return false;
    }
  }

  /**
   * Gets system metrics and statistics.
   *
   * @returns Promise resolving to system metrics
   *
   * @example
   * ```typescript
   * const metrics = await vault.system.metrics();
   * console.log("System uptime:", metrics.uptime);
   * console.log("Memory usage:", metrics.memory);
   * console.log("Active connections:", metrics.connections);
   * ```
   */
  public async metrics(): Promise<{
    uptime: number;
    memory: {
      used: number;
      total: number;
      percentage: number;
    };
    connections: {
      active: number;
      total: number;
    };
    requests: {
      total: number;
      perSecond: number;
      averageResponseTime: number;
    };
  }> {
    return this.client.get<any>("/api/v1/system/metrics");
  }

  /**
   * Performs a comprehensive system status check.
   * This combines health, version, and readiness checks.
   *
   * @returns Promise resolving to complete system status
   *
   * @example
   * ```typescript
   * const status = await vault.system.status();
   * console.log("System healthy:", status.healthy);
   * console.log("Version:", status.version);
   * console.log("Ready:", status.ready);
   * console.log("Uptime:", status.uptime);
   * ```
   */
  public async status(): Promise<{
    healthy: boolean;
    ready: boolean;
    version: VersionResponse;
    health: HealthResponse;
    uptime?: number;
    timestamp: Date;
  }> {
    const [health, version, ready] = await Promise.allSettled([
      this.health(),
      this.version(),
      this.ready(),
    ]);

    return {
      healthy:
        health.status === "fulfilled" && health.value.status === "healthy",
      ready: ready.status === "fulfilled" && ready.value,
      version:
        version.status === "fulfilled"
          ? version.value
          : {
              version: "unknown",
              buildTime: "unknown",
              gitCommit: "unknown",
              goVersion: "unknown",
            },
      health:
        health.status === "fulfilled"
          ? health.value
          : {
              status: "unhealthy",
              timestamp: new Date(),
              version: "unknown",
              database: "unknown",
            },
      timestamp: new Date(),
    };
  }
}
