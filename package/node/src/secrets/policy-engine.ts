import {
  SecurityLevel,
  SecretCategory,
  SecretMetadata,
} from "./security-levels.js";

/**
 * Access control policy types.
 */
export enum PolicyType {
  /** Role-based access control */
  RBAC = "rbac",

  /** Attribute-based access control */
  ABAC = "abac",

  /** Time-based access control */
  TBAC = "tbac",

  /** Location-based access control */
  LBAC = "lbac",

  /** Device-based access control */
  DBAC = "dbac",
}

/**
 * Access control condition types.
 */
export enum ConditionType {
  /** String equality */
  STRING_EQUALS = "StringEquals",

  /** String contains */
  STRING_CONTAINS = "StringContains",

  /** Numeric equality */
  NUMERIC_EQUALS = "NumericEquals",

  /** Numeric greater than */
  NUMERIC_GREATER_THAN = "NumericGreaterThan",

  /** Numeric less than */
  NUMERIC_LESS_THAN = "NumericLessThan",

  /** Boolean equality */
  BOOL = "Bool",

  /** IP address */
  IP_ADDRESS = "IpAddress",

  /** Date/time */
  DATE_TIME = "DateTime",

  /** List contains */
  FOR_ALL_VALUES = "ForAllValues",

  /** Any value matches */
  FOR_ANY_VALUE = "ForAnyValue",
}

/**
 * Access control condition.
 */
export interface AccessCondition {
  /** Condition type */
  type: ConditionType;

  /** Condition key */
  key: string;

  /** Condition values */
  values: string[] | number[] | boolean[];

  /** Whether condition is negated */
  negate?: boolean;

  /** Condition weight for priority */
  weight?: number;
}

/**
 * Access control rule.
 */
export interface AccessRule {
  /** Rule identifier */
  id: string;

  /** Rule name */
  name: string;

  /** Rule description */
  description?: string;

  /** Policy type */
  policyType: PolicyType;

  /** Access conditions (AND logic) */
  conditions: AccessCondition[];

  /** Effect of the rule */
  effect: "allow" | "deny";

  /** Rule priority (lower number = higher priority) */
  priority: number;

  /** Rule status */
  status: "active" | "inactive" | "deprecated";

  /** Time restrictions */
  timeRestrictions?: {
    validFrom?: string;
    validUntil?: string;
    allowedHours?: {
      start: string;
      end: string;
    };
    allowedDays?: number[];
    timezone?: string;
  };

  /** Rule metadata */
  metadata?: Record<string, unknown>;

  /** Created timestamp */
  createdAt: string;

  /** Updated timestamp */
  updatedAt: string;

  /** Created by */
  createdBy: string;

  /** Updated by */
  updatedBy?: string;
}

/**
 * Access control policy.
 */
export interface AccessPolicy {
  /** Policy identifier */
  id: string;

  /** Policy name */
  name: string;

  /** Policy description */
  description?: string;

  /** Policy version */
  version: string;

  /** Associated secret ID or pattern */
  secretId?: string;
  secretPattern?: string;

  /** Security level filter */
  securityLevel?: SecurityLevel[];

  /** Category filter */
  category?: SecretCategory[];

  /** Access rules */
  rules: AccessRule[];

  /** Default effect when no rules match */
  defaultEffect: "allow" | "deny";

  /** Policy status */
  status: "active" | "inactive" | "draft";

  /** Policy metadata */
  metadata?: Record<string, unknown>;

  /** Created timestamp */
  createdAt: string;

  /** Updated timestamp */
  updatedAt: string;

  /** Created by */
  createdBy: string;

  /** Updated by */
  updatedBy?: string;
}

/**
 * Access request context.
 */
export interface AccessContext {
  /** User identifier */
  userId: string;

  /** User roles */
  userRoles: string[];

  /** User attributes */
  userAttributes: Record<string, unknown>;

  /** Request timestamp */
  timestamp: string;

  /** Request source IP */
  sourceIp: string;

  /** User agent */
  userAgent?: string;

  /** Device identifier */
  deviceId?: string;

  /** Geographic location */
  location?: {
    country: string;
    region?: string;
    city?: string;
    coordinates?: {
      latitude: number;
      longitude: number;
    };
  };

  /** Request purpose */
  purpose?: string;

  /** Session identifier */
  sessionId?: string;

  /** Request method */
  requestMethod: string;

  /** Request resource */
  resource: string;
}

/**
 * Access decision result.
 */
export interface AccessDecision {
  /** Whether access is granted */
  granted: boolean;

  /** Decision reason */
  reason: string;

  /** Applied rules */
  appliedRules: {
    ruleId: string;
    ruleName: string;
    effect: "allow" | "deny";
    matched: boolean;
  }[];

  /** Decision timestamp */
  timestamp: string;

  /** Decision duration in milliseconds */
  duration: number;

  /** Additional metadata */
  metadata?: Record<string, unknown>;
}

/**
 * Policy evaluation engine for fine-grained access control.
 */
export class PolicyEngine {
  private policies: Map<string, AccessPolicy> = new Map();

  /**
   * Adds or updates a policy.
   */
  public addPolicy(policy: AccessPolicy): void {
    this.policies.set(policy.id, policy);
  }

  /**
   * Removes a policy.
   */
  public removePolicy(policyId: string): boolean {
    return this.policies.delete(policyId);
  }

  /**
   * Gets a policy by ID.
   */
  public getPolicy(policyId: string): AccessPolicy | undefined {
    return this.policies.get(policyId);
  }

  /**
   * Lists all policies.
   */
  public listPolicies(): AccessPolicy[] {
    return Array.from(this.policies.values());
  }

  /**
   * Evaluates access for a given secret and context.
   */
  public evaluateAccess(
    secretId: string,
    secretMetadata: SecretMetadata,
    context: AccessContext,
  ): AccessDecision {
    const startTime = Date.now();
    const appliedRules: AccessDecision["appliedRules"] = [];

    try {
      // Find applicable policies
      const applicablePolicies = this.findApplicablePolicies(
        secretId,
        secretMetadata,
      );

      if (applicablePolicies.length === 0) {
        return {
          granted: false,
          reason: "No applicable policies found",
          appliedRules,
          timestamp: new Date().toISOString(),
          duration: Date.now() - startTime,
        };
      }

      // Sort policies by priority and rules by priority
      const sortedPolicies = applicablePolicies.sort((a, b) => {
        // Sort by status (active first)
        if (a.status !== b.status) {
          return a.status === "active" ? -1 : 1;
        }
        return 0; // Could add more sorting criteria
      });

      // Evaluate each policy
      for (const policy of sortedPolicies) {
        if (policy.status !== "active") {
          continue;
        }

        // Sort rules by priority
        const sortedRules = policy.rules.sort(
          (a, b) => a.priority - b.priority,
        );

        for (const rule of sortedRules) {
          if (rule.status !== "active") {
            continue;
          }

          const ruleResult = this.evaluateRule(rule, context, secretMetadata);
          appliedRules.push({
            ruleId: rule.id,
            ruleName: rule.name,
            effect: rule.effect,
            matched: ruleResult.matched,
          });

          if (ruleResult.matched) {
            const decision = {
              granted: rule.effect === "allow",
              reason: `Rule "${rule.name}" matched with effect "${rule.effect}"`,
              appliedRules,
              timestamp: new Date().toISOString(),
              duration: Date.now() - startTime,
            };

            return decision;
          }
        }
      }

      // No rules matched, use default effect
      const defaultEffect = sortedPolicies[0]?.defaultEffect || "deny";

      return {
        granted: defaultEffect === "allow",
        reason: `No rules matched, using default effect "${defaultEffect}"`,
        appliedRules,
        timestamp: new Date().toISOString(),
        duration: Date.now() - startTime,
      };
    } catch (error) {
      return {
        granted: false,
        reason: `Policy evaluation failed: ${error instanceof Error ? error.message : "Unknown error"}`,
        appliedRules,
        timestamp: new Date().toISOString(),
        duration: Date.now() - startTime,
      };
    }
  }

  /**
   * Finds policies applicable to a secret.
   */
  private findApplicablePolicies(
    secretId: string,
    secretMetadata: SecretMetadata,
  ): AccessPolicy[] {
    const applicable: AccessPolicy[] = [];

    for (const policy of this.policies.values()) {
      // Check if policy applies to this secret
      if (policy.secretId && policy.secretId !== secretId) {
        continue;
      }

      if (
        policy.secretPattern &&
        !this.matchesPattern(secretId, policy.secretPattern)
      ) {
        continue;
      }

      if (
        policy.securityLevel &&
        !policy.securityLevel.includes(secretMetadata.securityLevel)
      ) {
        continue;
      }

      if (
        policy.category &&
        !policy.category.includes(secretMetadata.category)
      ) {
        continue;
      }

      applicable.push(policy);
    }

    return applicable;
  }

  /**
   * Evaluates a single rule against the context.
   */
  private evaluateRule(
    rule: AccessRule,
    context: AccessContext,
    secretMetadata: SecretMetadata,
  ): { matched: boolean; reason?: string } {
    // Check time restrictions first
    if (rule.timeRestrictions) {
      const now = new Date();

      if (
        rule.timeRestrictions.validFrom &&
        new Date(rule.timeRestrictions.validFrom) > now
      ) {
        return { matched: false, reason: "Rule not yet valid" };
      }

      if (
        rule.timeRestrictions.validUntil &&
        new Date(rule.timeRestrictions.validUntil) < now
      ) {
        return { matched: false, reason: "Rule expired" };
      }

      if (rule.timeRestrictions.allowedHours) {
        const currentHour = now.getHours();
        const [startHour, startMinute] =
          rule.timeRestrictions.allowedHours.start
            .split(":")
            .map((h) => parseInt(h, 10));
        const [endHourH, endMinute] = rule.timeRestrictions.allowedHours.end
          .split(":")
          .map((h) => parseInt(h, 10));

        const currentMinutes = currentHour * 60 + now.getMinutes();
        const startMinutes = startHour * 60 + (startMinute || 0);
        const endMinutes = endHourH * 60 + (endMinute || 0);

        if (currentMinutes < startMinutes || currentMinutes > endMinutes) {
          return { matched: false, reason: "Outside allowed hours" };
        }
      }

      if (rule.timeRestrictions.allowedDays) {
        const currentDay = now.getDay();
        if (!rule.timeRestrictions.allowedDays.includes(currentDay)) {
          return { matched: false, reason: "Outside allowed days" };
        }
      }
    }

    // Evaluate all conditions (AND logic)
    for (const condition of rule.conditions) {
      if (!this.evaluateCondition(condition, context, secretMetadata)) {
        return {
          matched: false,
          reason: `Condition "${condition.key}" not met`,
        };
      }
    }

    return { matched: true };
  }

  /**
   * Evaluates a single condition.
   */
  private evaluateCondition(
    condition: AccessCondition,
    context: AccessContext,
    secretMetadata: SecretMetadata,
  ): boolean {
    let actualValue: unknown;

    // Get the actual value based on the key
    switch (condition.key) {
      case "userId":
        actualValue = context.userId;
        break;
      case "userRoles":
        actualValue = context.userRoles;
        break;
      case "sourceIp":
        actualValue = context.sourceIp;
        break;
      case "securityLevel":
        actualValue = secretMetadata.securityLevel;
        break;
      case "category":
        actualValue = secretMetadata.category;
        break;
      case "ownerId":
        actualValue = secretMetadata.ownerId;
        break;
      default:
        // Check user attributes
        if (condition.key.startsWith("user.")) {
          const attrKey = condition.key.substring(5);
          actualValue = context.userAttributes[attrKey];
        } else {
          actualValue = context.userAttributes[condition.key];
        }
    }

    const result = this.compareValues(
      condition.type,
      actualValue,
      condition.values,
    );
    return condition.negate ? !result : result;
  }

  /**
   * Compares values based on condition type.
   */
  private compareValues(
    type: ConditionType,
    actual: unknown,
    expected: (string | number | boolean)[],
  ): boolean {
    switch (type) {
      case ConditionType.STRING_EQUALS:
        return expected.includes(String(actual));

      case ConditionType.STRING_CONTAINS:
        return expected.some((exp) =>
          String(actual).toLowerCase().includes(String(exp).toLowerCase()),
        );

      case ConditionType.NUMERIC_EQUALS:
        return expected.includes(Number(actual));

      case ConditionType.NUMERIC_GREATER_THAN:
        return Number(actual) > (expected[0] as number);

      case ConditionType.NUMERIC_LESS_THAN:
        return Number(actual) < (expected[0] as number);

      case ConditionType.BOOL:
        return expected.includes(Boolean(actual));

      case ConditionType.IP_ADDRESS:
        return this.checkIpInRange(String(actual), expected as string[]);

      case ConditionType.DATE_TIME:
        return this.checkDateTimeRange(actual as string, expected as string[]);

      case ConditionType.FOR_ALL_VALUES:
        if (Array.isArray(actual)) {
          return (expected as string[]).every((exp) => actual.includes(exp));
        }
        return expected.includes(String(actual));

      case ConditionType.FOR_ANY_VALUE:
        if (Array.isArray(actual)) {
          return (expected as string[]).some((exp) => actual.includes(exp));
        }
        return expected.includes(String(actual));

      default:
        return false;
    }
  }

  /**
   * Checks if IP is in allowed ranges.
   */
  private checkIpInRange(ip: string, ranges: string[]): boolean {
    // Simple IP range checking (could be enhanced with CIDR support)
    return ranges.some((range) => {
      if (range.includes("/")) {
        // CIDR notation - simplified check
        const [network] = range.split("/");
        return ip.startsWith(network.slice(0, -1));
      }
      return ip === range;
    });
  }

  /**
   * Checks if datetime is in allowed ranges.
   */
  private checkDateTimeRange(datetime: string, ranges: string[]): boolean {
    const date = new Date(datetime);
    return ranges.some((range) => {
      const [start, end] = range.split(",");
      return date >= new Date(start) && date <= new Date(end);
    });
  }

  /**
   * Checks if secret ID matches a pattern.
   */
  private matchesPattern(secretId: string, pattern: string): boolean {
    // Simple pattern matching (could be enhanced with regex)
    if (pattern.includes("*")) {
      const regex = new RegExp(pattern.replace(/\*/g, ".*"));
      return regex.test(secretId);
    }
    return secretId === pattern;
  }
}
