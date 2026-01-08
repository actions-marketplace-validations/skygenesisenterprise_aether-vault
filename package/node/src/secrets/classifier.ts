import {
  SecurityLevel,
  SecretCategory,
  SecretMetadata,
  CreateEnhancedSecretRequest,
} from "./security-levels.js";

/**
 * Secret classification engine for fine-grained security management.
 */
export class SecretClassifier {
  /**
   * Classifies a secret based on its content and context.
   */
  public static classify(secret: CreateEnhancedSecretRequest): SecretMetadata {
    const securityLevel = this.determineSecurityLevel(secret);
    const category = this.determineCategory(secret);
    const riskScore = this.calculateRiskScore(secret, securityLevel, category);
    const businessImpact = this.assessBusinessImpact(
      secret,
      category,
      securityLevel,
    );

    return {
      ...secret.metadata,
      securityLevel,
      category,
      riskScore,
      businessImpact,
      requiredPermissions: this.getRequiredPermissions(securityLevel, category),
      rotationPolicy: this.getDefaultRotationPolicy(securityLevel, category),
      retentionPolicy: this.getDefaultRetentionPolicy(securityLevel),
      compliance: this.getComplianceRequirements(securityLevel, category),
      classificationTags: this.generateClassificationTags(
        securityLevel,
        category,
      ),
    };
  }

  /**
   * Determines the security level based on secret characteristics.
   */
  private static determineSecurityLevel(
    secret: CreateEnhancedSecretRequest,
  ): SecurityLevel {
    const { name, value, metadata } = secret;

    // High-value indicators
    const highValuePatterns = [
      /password/i,
      /key/i,
      /secret/i,
      /token/i,
      /credential/i,
      /private/i,
    ];

    // Critical patterns
    const criticalPatterns = [
      /admin/i,
      /root/i,
      /master/i,
      /production/i,
      /prod/i,
      /api.*key/i,
    ];

    const hasHighValuePattern = highValuePatterns.some(
      (pattern) => pattern.test(name) || pattern.test(value),
    );

    const hasCriticalPattern = criticalPatterns.some(
      (pattern) => pattern.test(name) || pattern.test(value),
    );

    // Length and complexity analysis
    const isComplex = value.length > 32 && /[!@#$%^&*(),.?":{}|<>]/.test(value);

    // Explicit security level
    if (metadata.securityLevel) {
      return metadata.securityLevel;
    }

    // Auto-classification logic
    if (hasCriticalPattern || (hasHighValuePattern && isComplex)) {
      return SecurityLevel.TOP_SECRET;
    }

    if (hasHighValuePattern || isComplex) {
      return SecurityLevel.SECRET;
    }

    if (value.length > 16) {
      return SecurityLevel.CONFIDENTIAL;
    }

    if (value.length > 8) {
      return SecurityLevel.INTERNAL;
    }

    return SecurityLevel.PUBLIC;
  }

  /**
   * Determines the secret category based on name and content.
   */
  private static determineCategory(
    secret: CreateEnhancedSecretRequest,
  ): SecretCategory {
    const { name, value } = secret;

    const categoryPatterns = {
      [SecretCategory.API_KEY]: [
        /api.*key/i,
        /token/i,
        /jwt/i,
        /bearer/i,
        /auth/i,
      ],
      [SecretCategory.DATABASE]: [
        /database/i,
        /db/i,
        /mysql/i,
        /postgres/i,
        /mongodb/i,
        /connection.*string/i,
      ],
      [SecretCategory.ENCRYPTION_KEY]: [
        /encryption/i,
        /cipher/i,
        /aes/i,
        /rsa/i,
        /private.*key/i,
        /public.*key/i,
      ],
      [SecretCategory.CERTIFICATE]: [
        /certificate/i,
        /cert/i,
        /ssl/i,
        /tls/i,
        /pem/i,
        /crt/i,
      ],
      [SecretCategory.SERVICE]: [
        /service/i,
        /microservice/i,
        /application/i,
        /app.*secret/i,
      ],
      [SecretCategory.USER_CREDENTIALS]: [
        /user/i,
        /password/i,
        /login/i,
        /auth/i,
      ],
      [SecretCategory.CONFIGURATION]: [
        /config/i,
        /setting/i,
        /env/i,
        /environment/i,
      ],
      [SecretCategory.TEMPORARY]: [/temp/i, /session/i, /nonce/i, /otp/i],
    };

    for (const [category, patterns] of Object.entries(categoryPatterns)) {
      if (
        patterns.some((pattern) => pattern.test(name) || pattern.test(value))
      ) {
        return category as SecretCategory;
      }
    }

    return SecretCategory.CONFIGURATION; // Default category
  }

  /**
   * Calculates risk score (0-100) for the secret.
   */
  private static calculateRiskScore(
    secret: CreateEnhancedSecretRequest,
    securityLevel: SecurityLevel,
    category: SecretCategory,
  ): number {
    let score = 0;

    // Base score by security level
    const levelScores = {
      [SecurityLevel.PUBLIC]: 10,
      [SecurityLevel.INTERNAL]: 25,
      [SecurityLevel.CONFIDENTIAL]: 50,
      [SecurityLevel.SECRET]: 75,
      [SecurityLevel.TOP_SECRET]: 95,
    };

    score += levelScores[securityLevel];

    // Category risk modifiers
    const categoryModifiers = {
      [SecretCategory.API_KEY]: 10,
      [SecretCategory.DATABASE]: 15,
      [SecretCategory.ENCRYPTION_KEY]: 20,
      [SecretCategory.CERTIFICATE]: 5,
      [SecretCategory.SERVICE]: 5,
      [SecretCategory.USER_CREDENTIALS]: 15,
      [SecretCategory.CONFIGURATION]: 0,
      [SecretCategory.TEMPORARY]: -5,
    };

    score += categoryModifiers[category];

    // Value complexity modifier
    if (secret.value.length > 50) score += 5;
    if (/[A-Z]/.test(secret.value)) score += 2;
    if (/[0-9]/.test(secret.value)) score += 2;
    if (/[^A-Za-z0-9]/.test(secret.value)) score += 3;

    // Name sensitivity
    if (/prod/i.test(secret.name) || /production/i.test(secret.name))
      score += 10;
    if (/admin/i.test(secret.name) || /root/i.test(secret.name)) score += 15;

    return Math.min(100, Math.max(0, score));
  }

  /**
   * Assesses business impact level.
   */
  private static assessBusinessImpact(
    secret: CreateEnhancedSecretRequest,
    category: SecretCategory,
    securityLevel: SecurityLevel,
  ): "low" | "medium" | "high" | "critical" {
    // Critical categories
    if (
      [SecretCategory.DATABASE, SecretCategory.ENCRYPTION_KEY].includes(
        category,
      )
    ) {
      return "critical";
    }

    // High security levels
    if (
      [SecurityLevel.SECRET, SecurityLevel.TOP_SECRET].includes(securityLevel)
    ) {
      return "high";
    }

    // Production indicators
    if (/prod/i.test(secret.name) || /production/i.test(secret.name)) {
      return "high";
    }

    // Medium impact
    if ([SecurityLevel.CONFIDENTIAL].includes(securityLevel)) {
      return "medium";
    }

    return "low";
  }

  /**
   * Gets required permissions based on security level and category.
   */
  private static getRequiredPermissions(
    securityLevel: SecurityLevel,
    category: SecretCategory,
  ): string[] {
    const permissions = ["secrets:read"];

    if (securityLevel >= SecurityLevel.CONFIDENTIAL) {
      permissions.push("secrets:decrypt");
    }

    if (securityLevel >= SecurityLevel.SECRET) {
      permissions.push("secrets:audit", "secrets:access_log");
    }

    if (securityLevel === SecurityLevel.TOP_SECRET) {
      permissions.push("secrets:approve", "secrets:monitor");
    }

    // Category-specific permissions
    if (category === SecretCategory.ENCRYPTION_KEY) {
      permissions.push("crypto:use_key");
    }

    if (category === SecretCategory.CERTIFICATE) {
      permissions.push("cert:verify");
    }

    return permissions;
  }

  /**
   * Gets default rotation policy based on security level and category.
   */
  private static getDefaultRotationPolicy(
    securityLevel: SecurityLevel,
    category: SecretCategory,
  ): SecretMetadata["rotationPolicy"] {
    const intervals = {
      [SecurityLevel.PUBLIC]: 0,
      [SecurityLevel.INTERNAL]: 90,
      [SecurityLevel.CONFIDENTIAL]: 60,
      [SecurityLevel.SECRET]: 30,
      [SecurityLevel.TOP_SECRET]: 7,
    };

    const interval = intervals[securityLevel];

    if (interval === 0) {
      return { enabled: false };
    }

    // Category-specific adjustments
    let adjustedInterval = interval;
    if (category === SecretCategory.API_KEY) {
      adjustedInterval = Math.min(interval, 30);
    }

    if (category === SecretCategory.TEMPORARY) {
      adjustedInterval = 1;
    }

    return {
      enabled: true,
      intervalDays: adjustedInterval,
      autoRotate: securityLevel >= SecurityLevel.CONFIDENTIAL,
      notifyBeforeDays: Math.max(7, Math.floor(adjustedInterval / 4)),
    };
  }

  /**
   * Gets default retention policy.
   */
  private static getDefaultRetentionPolicy(
    securityLevel: SecurityLevel,
  ): SecretMetadata["retentionPolicy"] {
    const retainDays = {
      [SecurityLevel.PUBLIC]: 30,
      [SecurityLevel.INTERNAL]: 90,
      [SecurityLevel.CONFIDENTIAL]: 365,
      [SecurityLevel.SECRET]: 2555, // 7 years
      [SecurityLevel.TOP_SECRET]: 3650, // 10 years
    };

    return {
      retainAfterDeletion: securityLevel >= SecurityLevel.CONFIDENTIAL,
      retainDays: retainDays[securityLevel],
      permanentArchive: securityLevel === SecurityLevel.TOP_SECRET,
    };
  }

  /**
   * Gets compliance requirements.
   */
  private static getComplianceRequirements(
    securityLevel: SecurityLevel,
    category: SecretCategory,
  ): SecretMetadata["compliance"] {
    const standards = [];

    if (securityLevel >= SecurityLevel.CONFIDENTIAL) {
      standards.push("ISO27001");
    }

    if (securityLevel >= SecurityLevel.SECRET) {
      standards.push("GDPR", "SOC2");
    }

    if (securityLevel === SecurityLevel.TOP_SECRET) {
      standards.push("PCI-DSS", "HIPAA");
    }

    if (category === SecretCategory.ENCRYPTION_KEY) {
      standards.push("FIPS140-2");
    }

    return {
      standards,
      auditRequired: securityLevel >= SecurityLevel.SECRET,
      encryptionStandard:
        securityLevel >= SecurityLevel.CONFIDENTIAL ? "AES-256" : "AES-128",
    };
  }

  /**
   * Generates classification tags.
   */
  private static generateClassificationTags(
    securityLevel: SecurityLevel,
    category: SecretCategory,
  ): string[] {
    return [
      `level:${securityLevel}`,
      `category:${category}`,
      `auto-classified:${new Date().toISOString()}`,
    ];
  }
}
