/**
 * Network management types for Aether Vault SDK.
 * Based on server Go models from server/src/model/network.go and server/src/controllers/network.go
 */

/**
 * Protocol type enumeration.
 * Supported network protocols for network connections.
 */
export type ProtocolType =
  | "http"
  | "https"
  | "ssh"
  | "ftp"
  | "sftp"
  | "webdav"
  | "smb"
  | "nfs"
  | "rsync"
  | "git"
  | "custom";

/**
 * Protocol configuration interface.
 * Configuration details for a specific protocol connection.
 */
export interface ProtocolConfig {
  /** Host address or hostname */
  host: string;

  /** Port number */
  port: number;

  /** Username for authentication (optional) */
  username?: string;

  /** Password for authentication (optional) */
  password?: string;

  /** Private key for SSH authentication (optional) */
  privateKey?: string;

  /** Certificate for HTTPS authentication (optional) */
  certificate?: string;

  /** Connection timeout in seconds (optional) */
  timeout?: number;

  /** Custom headers for HTTP/HTTPS (optional) */
  headers?: Record<string, string>;

  /** Additional protocol-specific options (optional) */
  options?: Record<string, unknown>;
}

/**
 * Network interface representing a vault network.
 * Matches the Go Network struct with proper TypeScript types.
 */
export interface Network {
  /** Unique network identifier */
  id: number;

  /** Network name (unique) */
  name: string;

  /** Network protocol type */
  type: ProtocolType;

  /** Network status */
  status: string;

  /** Network creation timestamp */
  createdAt: string;

  /** Network last update timestamp */
  updatedAt: string;
}

/**
 * Network with configuration interface.
 * Network object that includes protocol configuration.
 */
export interface NetworkWithConfig extends Network {
  /** Protocol configuration details */
  config?: ProtocolConfig;
}

/**
 * Network creation request interface.
 * Used for creating new networks via the API.
 */
export interface CreateNetworkRequest {
  /** Network name (must be unique) */
  name: string;

  /** Network protocol type */
  type: ProtocolType;

  /** Protocol configuration */
  config?: ProtocolConfig;
}

/**
 * Network update request interface.
 * Used for updating existing networks via the API.
 */
export interface UpdateNetworkRequest {
  /** New network name (optional) */
  name?: string;

  /** New protocol type (optional) */
  type?: ProtocolType;

  /** New protocol configuration (optional) */
  config?: ProtocolConfig;
}

/**
 * Network list response interface.
 * Response containing multiple networks with configuration.
 */
export interface NetworkListResponse {
  /** Array of networks */
  networks: NetworkWithConfig[];

  /** Total number of networks */
  total: number;

  /** Current page number */
  page: number;

  /** Number of networks per page */
  pageSize: number;

  /** Total number of pages */
  totalPages: number;
}

/**
 * Network filter parameters interface.
 * Filtering and pagination options for network list requests.
 */
export interface NetworkFilterParams extends Record<string, unknown> {
  /** Page number (default: 1) */
  page?: number;

  /** Number of items per page (default: 20) */
  pageSize?: number;

  /** Sort field */
  sortBy?: string;

  /** Sort direction */
  sortOrder?: "asc" | "desc";

  /** Filter by protocol type */
  type?: ProtocolType;

  /** Filter by status */
  status?: string;

  /** Search in network name */
  search?: string;

  /** Filter by creation date from */
  createdFrom?: string;

  /** Filter by creation date to */
  createdTo?: string;
}

/**
 * Protocol status interface.
 * Status information for a specific protocol connection.
 */
export interface ProtocolStatus {
  /** Protocol type */
  protocol: ProtocolType;

  /** Connection status */
  status: string;

  /** Last status check timestamp */
  lastCheck: string;

  /** Status message (optional) */
  message?: string;

  /** Connection latency in milliseconds (optional) */
  latency?: number;
}

/**
 * Protocol test request interface.
 * Request for testing a protocol connection.
 */
export interface ProtocolTestRequest {
  /** Protocol type to test */
  type: ProtocolType;

  /** Protocol configuration for testing */
  config: ProtocolConfig;
}

/**
 * Protocol test response interface.
 * Response from protocol connection test.
 */
export interface ProtocolTestResponse {
  /** Whether the test was successful */
  success: boolean;

  /** Test result message */
  message: string;

  /** Test latency in milliseconds (optional) */
  latency?: number;

  /** Additional test details (optional) */
  details?: Record<string, unknown>;
}

/**
 * Supported protocols response interface.
 * Response containing list of supported protocols.
 */
export interface SupportedProtocolsResponse {
  /** Array of supported protocols */
  protocols: ProtocolType[];

  /** Total number of supported protocols */
  count: number;
}

/**
 * Network statistics interface.
 * Aggregated statistics about networks.
 */
export interface NetworkStats {
  /** Total number of networks */
  totalNetworks: number;

  /** Number of active networks */
  activeNetworks: number;

  /** Number of inactive networks */
  inactiveNetworks: number;

  /** Networks by protocol type */
  networksByProtocol: Record<ProtocolType, number>;

  /** Networks created in the last 24 hours */
  networksCreated24h: number;

  /** Networks created in the last 7 days */
  networksCreated7d: number;

  /** Networks created in the last 30 days */
  networksCreated30d: number;
}
