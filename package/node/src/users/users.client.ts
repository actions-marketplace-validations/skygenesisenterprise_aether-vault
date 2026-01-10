import { VaultClient } from "../core/client.js";
import {
  User,
  CreateUserRequest,
  UpdateUserRequest,
  UserListResponse,
  UserFilterParams,
  UserStats,
} from "../types/user.js";

/**
 * Client for managing users in Aether Vault.
 * Provides methods to create, read, update, and delete users.
 */
export class UsersClient {
  /**
   * Creates a new UsersClient instance.
   *
   * @param httpClient - HTTP client for API communication
   */
  constructor(private readonly httpClient: VaultClient) {}

  /**
   * Lists all users with optional filtering and pagination.
   *
   * @param params - Optional filter parameters
   * @returns Promise resolving to paginated user list
   *
   * @example
   * ```typescript
   * const users = await vault.users.list({
   *   page: 1,
   *   pageSize: 20,
   *   isActive: true,
   *   search: "john"
   * });
   * ```
   */
  public async list(params?: UserFilterParams): Promise<UserListResponse> {
    const response = await this.httpClient.get<{ users: User[] }>(
      "/api/v1/users",
      params,
    );

    // Convert response to UserListResponse format
    return {
      users: response.users,
      total: response.users.length,
      page: params?.page || 1,
      pageSize: params?.pageSize || 20,
      totalPages: Math.ceil(response.users.length / (params?.pageSize || 20)),
    };
  }

  /**
   * Gets a user by their ID.
   *
   * @param id - User ID (UUID)
   * @returns Promise resolving to the user details
   *
   * @example
   * ```typescript
   * const user = await vault.users.get("123e4567-e89b-12d3-a456-426614174000");
   * console.log(user.email, user.firstName, user.lastName);
   * ```
   */
  public async get(id: string): Promise<User> {
    return this.httpClient.get<User>(`/api/v1/users/${id}`);
  }

  /**
   * Creates a new user.
   *
   * @param user - User creation data
   * @returns Promise resolving to the created user
   *
   * @example
   * ```typescript
   * const user = await vault.users.create({
   *   email: "john.doe@example.com",
   *   password: "securePassword123",
   *   firstName: "John",
   *   lastName: "Doe"
   * });
   * ```
   */
  public async create(user: CreateUserRequest): Promise<User> {
    return this.httpClient.post<User>("/api/v1/users", user);
  }

  /**
   * Updates an existing user.
   *
   * @param id - User ID (UUID)
   * @param updates - User update data
   * @returns Promise resolving to the updated user
   *
   * @example
   * ```typescript
   * const updated = await vault.users.update("123e4567-e89b-12d3-a456-426614174000", {
   *   firstName: "Johnathan",
   *   lastName: "Smith"
   * });
   * ```
   */
  public async update(id: string, updates: UpdateUserRequest): Promise<User> {
    return this.httpClient.put<User>(`/api/v1/users/${id}`, updates);
  }

  /**
   * Deletes a user.
   *
   * @param id - User ID (UUID)
   * @returns Promise resolving when the user is deleted
   *
   * @example
   * ```typescript
   * await vault.users.delete("123e4567-e89b-12d3-a456-426614174000");
   * ```
   */
  public async delete(id: string): Promise<void> {
    return this.httpClient.delete<void>(`/api/v1/users/${id}`);
  }

  /**
   * Gets user statistics and metrics.
   *
   * @returns Promise resolving to user statistics
   *
   * @example
   * ```typescript
   * const stats = await vault.users.getStats();
   * console.log("Total users:", stats.totalUsers);
   * console.log("Active users:", stats.activeUsers);
   * ```
   */
  public async getStats(): Promise<UserStats> {
    return this.httpClient.get<UserStats>("/api/v1/users/stats");
  }

  /**
   * Searches users by email, first name, or last name.
   *
   * @param query - Search query string
   * @param params - Optional pagination parameters
   * @returns Promise resolving to search results
   *
   * @example
   * ```typescript
   * const results = await vault.users.search("john", {
   *   page: 1,
   *   pageSize: 10
   * });
   * ```
   */
  public async search(
    query: string,
    params?: Omit<UserFilterParams, "search">,
  ): Promise<UserListResponse> {
    return this.list({
      ...params,
      search: query,
    });
  }

  /**
   * Gets users filtered by active status.
   *
   * @param isActive - Whether to get active (true) or inactive (false) users
   * @param params - Optional pagination parameters
   * @returns Promise resolving to filtered user list
   *
   * @example
   * ```typescript
   * const activeUsers = await vault.users.getByActiveStatus(true);
   * const inactiveUsers = await vault.users.getByActiveStatus(false);
   * ```
   */
  public async getByActiveStatus(
    isActive: boolean,
    params?: Omit<UserFilterParams, "isActive">,
  ): Promise<UserListResponse> {
    return this.list({
      ...params,
      isActive,
    });
  }

  /**
   * Checks if a user exists by ID.
   *
   * @param id - User ID (UUID)
   * @returns Promise resolving to true if the user exists
   *
   * @example
   * ```typescript
   * const exists = await vault.users.exists("123e4567-e89b-12d3-a456-426614174000");
   * if (exists) {
   *   console.log("User exists");
   * }
   * ```
   */
  public async exists(id: string): Promise<boolean> {
    try {
      await this.get(id);
      return true;
    } catch (error) {
      // If we get a 404, the user doesn't exist
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
   * Activates a user account.
   *
   * @param id - User ID (UUID)
   * @returns Promise resolving to the updated user
   *
   * @example
   * ```typescript
   * const activated = await vault.users.activate("123e4567-e89b-12d3-a456-426614174000");
   * console.log("User is now active:", activated.isActive);
   * ```
   */
  public async activate(id: string): Promise<User> {
    return this.update(id, { isActive: true });
  }

  /**
   * Deactivates a user account.
   *
   * @param id - User ID (UUID)
   * @returns Promise resolving to the updated user
   *
   * @example
   * ```typescript
   * const deactivated = await vault.users.deactivate("123e4567-e89b-12d3-a456-426614174000");
   * console.log("User is now inactive:", deactivated.isActive);
   * ```
   */
  public async deactivate(id: string): Promise<User> {
    return this.update(id, { isActive: false });
  }

  /**
   * Gets users created within a specific date range.
   *
   * @param fromDate - Start date (ISO string)
   * @param toDate - End date (ISO string)
   * @param params - Optional pagination parameters
   * @returns Promise resolving to filtered user list
   *
   * @example
   * ```typescript
   * const recentUsers = await vault.users.getByDateRange(
   *   "2025-01-01T00:00:00.000Z",
   *   "2025-01-31T23:59:59.999Z"
   * );
   * ```
   */
  public async getByDateRange(
    fromDate: string,
    toDate: string,
    params?: Omit<UserFilterParams, "createdFrom" | "createdTo">,
  ): Promise<UserListResponse> {
    return this.list({
      ...params,
      createdFrom: fromDate,
      createdTo: toDate,
    });
  }
}
