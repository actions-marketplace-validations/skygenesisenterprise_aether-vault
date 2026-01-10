import { VaultClient } from "../core/client.js";
import {
  NetworkWithConfig,
  CreateNetworkRequest,
  UpdateNetworkRequest,
  NetworkListResponse,
  NetworkFilterParams,
  ProtocolTestRequest,
  ProtocolTestResponse,
  SupportedProtocolsResponse,
  NetworkStats,
  ProtocolType,
} from "../types/network.js";

/**
 * Client for managing networks in Aether Vault.
 * Provides methods to create, read, update, and delete network configurations.
 */
export class NetworkClient {
  /**
   * Creates a new NetworkClient instance.
   *
   * @param httpClient - HTTP client for API communication
   */
  constructor(private readonly httpClient: VaultClient) {}

  /**
   * Lists all networks with optional filtering and pagination.
   *
   * @param params - Optional filter parameters
   * @returns Promise resolving to paginated network list
   *
   * @example
   * ```typescript
   * const networks = await vault.network.list({
   *   page: 1,
   *   pageSize: 20,
   *   type: "https",
   *   status: "active"
   * });
   * ```
   */
  public async list(
    params?: NetworkFilterParams,
  ): Promise<NetworkListResponse> {
    return this.httpClient.get<NetworkListResponse>("/api/v1/network", params);
  }

  /**
   * Gets a network by its ID.
   *
   * @param id - Network ID
   * @returns Promise resolving to the network details with configuration
   *
   * @example
   * ```typescript
   * const network = await vault.network.get(1);
   * console.log(network.name, network.type, network.config?.host);
   * ```
   */
  public async get(id: number): Promise<NetworkWithConfig> {
    return this.httpClient.get<NetworkWithConfig>(`/api/v1/network/${id}`);
  }

  /**
   * Creates a new network.
   *
   * @param network - Network creation data
   * @returns Promise resolving to the created network with configuration
   *
   * @example
   * ```typescript
   * const network = await vault.network.create({
   *   name: "Production API",
   *   type: "https",
   *   config: {
   *     host: "api.example.com",
   *     port: 443,
   *     timeout: 30
   *   }
   * });
   * ```
   */
  public async create(
    network: CreateNetworkRequest,
  ): Promise<NetworkWithConfig> {
    return this.httpClient.post<NetworkWithConfig>("/api/v1/network", network);
  }

  /**
   * Updates an existing network.
   *
   * @param id - Network ID
   * @param updates - Network update data
   * @returns Promise resolving to the updated network with configuration
   *
   * @example
   * ```typescript
   * const updated = await vault.network.update(1, {
   *   name: "Updated API",
   *   config: {
   *     host: "new-api.example.com",
   *     port: 443
   *   }
   * });
   * ```
   */
  public async update(
    id: number,
    updates: UpdateNetworkRequest,
  ): Promise<NetworkWithConfig> {
    return this.httpClient.put<NetworkWithConfig>(
      `/api/v1/network/${id}`,
      updates,
    );
  }

  /**
   * Deletes a network.
   *
   * @param id - Network ID
   * @returns Promise resolving when the network is deleted
   *
   * @example
   * ```typescript
   * await vault.network.delete(1);
   * ```
   */
  public async delete(id: number): Promise<void> {
    return this.httpClient.delete<void>(`/api/v1/network/${id}`);
  }

  /**
   * Gets the list of supported protocols.
   *
   * @returns Promise resolving to supported protocols
   *
   * @example
   * ```typescript
   * const protocols = await vault.network.getSupportedProtocols();
   * console.log("Supported protocols:", protocols.protocols);
   * ```
   */
  public async getSupportedProtocols(): Promise<SupportedProtocolsResponse> {
    return this.httpClient.get<SupportedProtocolsResponse>(
      "/api/v1/network/protocols",
    );
  }

  /**
   * Tests a protocol connection with provided configuration.
   *
   * @param request - Protocol test request
   * @returns Promise resolving to test results
   *
   * @example
   * ```typescript
   * const result = await vault.network.testProtocol({
   *   type: "https",
   *   config: {
   *     host: "api.example.com",
   *     port: 443,
   *     timeout: 10
   *   }
   * });
   *
   * if (result.success) {
   *   console.log("Connection successful, latency:", result.latency);
   * }
   * ```
   */
  public async testProtocol(
    request: ProtocolTestRequest,
  ): Promise<ProtocolTestResponse> {
    return this.httpClient.post<ProtocolTestResponse>(
      "/api/v1/network/test",
      request,
    );
  }

  /**
   * Gets the status of a specific protocol connection.
   *
   * @param id - Network ID
   * @returns Promise resolving to protocol status
   *
   * @example
   * ```typescript
   * const status = await vault.network.getProtocolStatus(1);
   * console.log("Protocol status:", status.status);
   * console.log("Last check:", status.lastCheck);
   * ```
   */
  public async getProtocolStatus(id: number): Promise<{
    protocol: ProtocolType;
    status: string;
    lastCheck: string;
    message?: string;
    latency?: number;
  }> {
    return this.httpClient.get(`/api/v1/network/${id}/status`);
  }

  /**
   * Gets network statistics and metrics.
   *
   * @returns Promise resolving to network statistics
   *
   * @example
   * ```typescript
   * const stats = await vault.network.getStats();
   * console.log("Total networks:", stats.totalNetworks);
   * console.log("Active networks:", stats.activeNetworks);
   * ```
   */
  public async getStats(): Promise<NetworkStats> {
    return this.httpClient.get<NetworkStats>("/api/v1/network/stats");
  }

  /**
   * Searches networks by name.
   *
   * @param query - Search query string
   * @param params - Optional pagination parameters
   * @returns Promise resolving to search results
   *
   * @example
   * ```typescript
   * const results = await vault.network.search("production", {
   *   page: 1,
   *   pageSize: 10
   * });
   * ```
   */
  public async search(
    query: string,
    params?: Omit<NetworkFilterParams, "search">,
  ): Promise<NetworkListResponse> {
    return this.list({
      ...params,
      search: query,
    });
  }

  /**
   * Gets networks filtered by protocol type.
   *
   * @param type - Protocol type to filter by
   * @param params - Optional pagination parameters
   * @returns Promise resolving to filtered network list
   *
   * @example
   * ```typescript
   * const httpsNetworks = await vault.network.getByProtocolType("https");
   * const sshNetworks = await vault.network.getByProtocolType("ssh");
   * ```
   */
  public async getByProtocolType(
    type: ProtocolType,
    params?: Omit<NetworkFilterParams, "type">,
  ): Promise<NetworkListResponse> {
    return this.list({
      ...params,
      type,
    });
  }

  /**
   * Gets networks filtered by status.
   *
   * @param status - Status to filter by
   * @param params - Optional pagination parameters
   * @returns Promise resolving to filtered network list
   *
   * @example
   * ```typescript
   * const activeNetworks = await vault.network.getByStatus("active");
   * const inactiveNetworks = await vault.network.getByStatus("inactive");
   * ```
   */
  public async getByStatus(
    status: string,
    params?: Omit<NetworkFilterParams, "status">,
  ): Promise<NetworkListResponse> {
    return this.list({
      ...params,
      status,
    });
  }

  /**
   * Checks if a network exists by ID.
   *
   * @param id - Network ID
   * @returns Promise resolving to true if the network exists
   *
   * @example
   * ```typescript
   * const exists = await vault.network.exists(1);
   * if (exists) {
   *   console.log("Network exists");
   * }
   * ```
   */
  public async exists(id: number): Promise<boolean> {
    try {
      await this.get(id);
      return true;
    } catch (error) {
      // If we get a 404, the network doesn't exist
      if (
        error &&
        typeof error === "object" &&
        "code" in error &&
        error.code === "NOT_FOUND"
      ) {
        return false;
      }
      // Re-throw other errors
      throw error;
    }
  }

  /**
   * Gets networks created within a specific date range.
   *
   * @param fromDate - Start date (ISO string)
   * @param toDate - End date (ISO string)
   * @param params - Optional pagination parameters
   * @returns Promise resolving to filtered network list
   *
   * @example
   * ```typescript
   * const recentNetworks = await vault.network.getByDateRange(
   *   "2025-01-01T00:00:00.000Z",
   *   "2025-01-31T23:59:59.999Z"
   * );
   * ```
   */
  public async getByDateRange(
    fromDate: string,
    toDate: string,
    params?: Omit<NetworkFilterParams, "createdFrom" | "createdTo">,
  ): Promise<NetworkListResponse> {
    return this.list({
      ...params,
      createdFrom: fromDate,
      createdTo: toDate,
    });
  }

  /**
   * Batch tests multiple networks.
   *
   * @param networkIds - Array of network IDs to test
   * @returns Promise resolving to array of test results
   *
   * @example
   * ```typescript
   * const results = await vault.network.batchTest([1, 2, 3]);
   * results.forEach(result => {
   *   console.log(`Network ${result.networkId}: ${result.success ? 'OK' : 'Failed'}`);
   * });
   * ```
   */
  public async batchTest(networkIds: number[]): Promise<
    Array<{
      networkId: number;
      success: boolean;
      message: string;
      latency?: number;
    }>
  > {
    const results = await Promise.allSettled(
      networkIds.map(async (id) => {
        const status = await this.getProtocolStatus(id);
        return {
          networkId: id,
          success: status.status === "active" || status.status === "healthy",
          message: status.message || status.status,
          latency: status.latency,
        };
      }),
    );

    return results.map((result, index) => {
      if (result.status === "fulfilled") {
        return result.value;
      }
      return {
        networkId: networkIds[index],
        success: false,
        message:
          result.status === "rejected"
            ? String(result.reason)
            : "Unknown error",
      };
    }) as Array<{
      networkId: number;
      success: boolean;
      message: string;
      latency?: number;
    }>;
  }
}
