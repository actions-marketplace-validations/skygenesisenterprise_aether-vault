import { VaultClient } from "../core/client.js";
import { Policy, PolicyListResponse } from "../types/index.js";

/**
 * Client for policy management operations.
 * Provides methods to manage access policies in Aether Vault.
 */
export class PolicyClient {
  /**
   * Creates a new PolicyClient instance.
   *
   * @param client - VaultClient instance for HTTP requests
   */
  constructor(private readonly client: VaultClient) {}

  /**
   * Retrieves all policies for the current user.
   *
   * @returns Promise resolving to list of policies
   *
   * @example
   * ```typescript
   * const policies = await vault.policies.list();
   * console.log(`Found ${policies.length} policies`);
   * ```
   */
  public async list(): Promise<PolicyListResponse> {
    const response =
      await this.client.get<PolicyListResponse>("/api/v1/policies");

    return response;
  }

  /**
   * Retrieves a specific policy by ID.
   *
   * @param id - Policy ID
   * @returns Promise resolving to policy
   *
   * @example
   * ```typescript
   * const policy = await vault.policies.get("policy-123");
   * console.log(policy.name, policy.effect, policy.actions);
   * ```
   */
  public async get(id: string): Promise<Policy> {
    const response = await this.client.get<Policy>(`/api/v1/policies/${id}`);

    return response;
  }

  /**
   * Creates a new policy.
   *
   * @param policy - Policy data to create
   * @returns Promise resolving to created policy
   *
   * @example
   * ```typescript
   * const newPolicy = await vault.policies.create({
   *   name: "Allow secret access",
   *   description: "Allows read access to user secrets",
   *   resource: "secret",
   *   actions: ["read"],
   *   effect: "allow",
   *   priority: 100
   * });
   * ```
   */
  public async create(
    policy: Omit<Policy, "id" | "createdAt" | "updatedAt" | "userId">,
  ): Promise<Policy> {
    const response = await this.client.post<Policy>("/api/v1/policies", policy);

    return response;
  }

  /**
   * Updates an existing policy.
   *
   * @param id - Policy ID to update
   * @param updates - Policy updates
   * @returns Promise resolving to updated policy
   *
   * @example
   * ```typescript
   * const updatedPolicy = await vault.policies.update("policy-123", {
   *   description: "Updated policy description",
   *   isActive: false
   * });
   * ```
   */
  public async update(
    id: string,
    updates: Partial<Omit<Policy, "id" | "createdAt" | "updatedAt" | "userId">>,
  ): Promise<Policy> {
    const response = await this.client.put<Policy>(
      `/api/v1/policies/${id}`,
      updates,
    );

    return response;
  }

  /**
   * Deletes a policy.
   *
   * @param id - Policy ID to delete
   * @returns Promise resolving when policy is deleted
   *
   * @example
   * ```typescript
   * await vault.policies.delete("policy-123");
   * console.log("Policy deleted successfully");
   * ```
   */
  public async delete(id: string): Promise<void> {
    await this.client.delete<void>(`/api/v1/policies/${id}`);
  }

  /**
   * Enables a policy.
   *
   * @param id - Policy ID to enable
   * @returns Promise resolving to updated policy
   *
   * @example
   * ```typescript
   * const policy = await vault.policies.enable("policy-123");
   * console.log("Policy is now active:", policy.isActive);
   * ```
   */
  public async enable(id: string): Promise<Policy> {
    return this.update(id, { isActive: true });
  }

  /**
   * Disables a policy.
   *
   * @param id - Policy ID to disable
   * @returns Promise resolving to updated policy
   *
   * @example
   * ```typescript
   * const policy = await vault.policies.disable("policy-123");
   * console.log("Policy is now inactive:", policy.isActive);
   * ```
   */
  public async disable(id: string): Promise<Policy> {
    return this.update(id, { isActive: false });
  }

  /**
   * Retrieves policies by resource type.
   *
   * @param resource - Resource type to filter by
   * @returns Promise resolving to filtered policies
   *
   * @example
   * ```typescript
   * const secretPolicies = await vault.policies.getByResource("secret");
   * const userPolicies = await vault.policies.getByResource("user");
   * ```
   */
  public async getByResource(resource: string): Promise<Policy[]> {
    const response = await this.client.get<PolicyListResponse>(
      "/api/v1/policies",
      { resource },
    );

    return response.policies;
  }

  /**
   * Retrieves active policies only.
   *
   * @returns Promise resolving to active policies
   *
   * @example
   * ```typescript
   * const activePolicies = await vault.policies.getActive();
   * console.log(`Found ${activePolicies.length} active policies`);
   * ```
   */
  public async getActive(): Promise<Policy[]> {
    const response = await this.client.get<PolicyListResponse>(
      "/api/v1/policies",
      { isActive: true },
    );

    return response.policies;
  }

  /**
   * Retrieves policies by effect type.
   *
   * @param effect - Effect type to filter by ("allow" or "deny")
   * @returns Promise resolving to filtered policies
   *
   * @example
   * ```typescript
   * const allowPolicies = await vault.policies.getByEffect("allow");
   * const denyPolicies = await vault.policies.getByEffect("deny");
   * ```
   */
  public async getByEffect(effect: "allow" | "deny"): Promise<Policy[]> {
    const response = await this.client.get<PolicyListResponse>(
      "/api/v1/policies",
      { effect },
    );

    return response.policies;
  }

  /**
   * Evaluates if a specific action is allowed for a resource.
   * This is a utility method that checks all applicable policies.
   *
   * @param resource - Resource type
   * @param action - Action to perform
   * @param context - Additional context for evaluation
   * @returns Promise resolving to evaluation result
   *
   * @example
   * ```typescript
   * const canRead = await vault.policies.evaluate("secret", "read", {
   *   userId: "user-123",
   *   resourceId: "secret-456"
   * });
   *
   * if (canRead.allowed) {
   *   // Proceed with action
   * } else {
   *   // Action is denied by policy
   *   console.log("Access denied:", canRead.reason);
   * }
   * ```
   */
  public async evaluate(
    resource: string,
    action: string,
    context?: Record<string, unknown>,
  ): Promise<{ allowed: boolean; reason?: string; matchedPolicies: string[] }> {
    const response = await this.client.post<{
      allowed: boolean;
      reason?: string;
      matchedPolicies: string[];
    }>("/api/v1/policies/evaluate", {
      resource,
      action,
      context,
    });

    return response;
  }
}
