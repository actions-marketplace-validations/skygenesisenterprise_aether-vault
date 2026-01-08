import { SecurityLevel } from "./security-levels.js";
import {
  createHash,
  randomBytes,
  createCipheriv,
  createDecipheriv,
  scrypt,
} from "crypto";

/**
 * Encryption algorithms supported by the vault.
 */
export enum EncryptionAlgorithm {
  /** AES-128-CBC */
  AES_128_CBC = "aes-128-cbc",

  /** AES-256-CBC */
  AES_256_CBC = "aes-256-cbc",

  /** AES-128-GCM */
  AES_128_GCM = "aes-128-gcm",

  /** AES-256-GCM */
  AES_256_GCM = "aes-256-gcm",

  /** ChaCha20-Poly1305 */
  CHACHA20_POLY1305 = "chacha20-poly1305",
}

/**
 * Key derivation functions.
 */
export enum KeyDerivationFunction {
  /** PBKDF2 with SHA-256 */
  PBKDF2_SHA256 = "pbkdf2-sha256",

  /** scrypt */
  SCRYPT = "scrypt",

  /** Argon2id */
  ARGON2ID = "argon2id",
}

/**
 * Encryption metadata.
 */
export interface EncryptionMetadata {
  /** Encryption algorithm used */
  algorithm: EncryptionAlgorithm;

  /** Key identifier */
  keyId: string;

  /** Key version */
  keyVersion: number;

  /** Initialization vector */
  iv: string;

  /** Authentication tag (for AEAD ciphers) */
  tag?: string;

  /** Key derivation function */
  kdf?: KeyDerivationFunction;

  /** Salt for key derivation */
  salt?: string;

  /** Iterations for PBKDF2 */
  iterations?: number;

  /** Memory cost for scrypt */
  memoryCost?: number;

  /** Time cost for scrypt */
  timeCost?: number;

  /** Parallelism for scrypt */
  parallelism?: number;

  /** Encryption timestamp */
  encryptedAt: string;

  /** Checksum for integrity verification */
  checksum: string;
}

/**
 * Encryption key metadata.
 */
export interface EncryptionKey {
  /** Key identifier */
  id: string;

  /** Key version */
  version: number;

  /** Security level this key is designed for */
  securityLevel: SecurityLevel;

  /** Key algorithm */
  algorithm: EncryptionAlgorithm;

  /** Key size in bits */
  keySize: number;

  /** Key usage */
  usage: "encryption" | "decryption" | "both";

  /** Key status */
  status: "active" | "deprecated" | "revoked";

  /** Key creation timestamp */
  createdAt: string;

  /** Key expiration timestamp */
  expiresAt?: string;

  /** Key rotation schedule */
  rotationSchedule?: {
    intervalDays: number;
    nextRotation: string;
  };

  /** Key metadata */
  metadata?: Record<string, unknown>;
}

/**
 * Multi-level encryption engine.
 */
export class MultiLevelEncryption {
  private keys: Map<string, EncryptionKey> = new Map();
  private masterKey: Buffer;

  constructor(masterKey?: string | Buffer) {
    if (masterKey) {
      this.masterKey =
        typeof masterKey === "string"
          ? Buffer.from(masterKey, "hex")
          : masterKey;
    } else {
      // Generate a random master key if none provided
      this.masterKey = randomBytes(32);
    }
  }

  /**
   * Derives encryption key for a given security level.
   */
  private async deriveKey(
    keyId: string,
    securityLevel: SecurityLevel,
    salt?: Buffer,
  ): Promise<Buffer> {
    const saltBuffer = salt || randomBytes(32);
    const keySize = this.getKeySizeForSecurityLevel(securityLevel);

    // Use scrypt for key derivation
    return new Promise((resolve, reject) => {
      scrypt(
        this.masterKey,
        `${keyId}:${securityLevel}:${saltBuffer.toString("hex")}`,
        keySize,
        (err, derivedKey) => {
          if (err) reject(err);
          else resolve(derivedKey);
        },
      );
    });
  }

  /**
   * Gets appropriate encryption algorithm for security level.
   */
  private getAlgorithmForSecurityLevel(
    securityLevel: SecurityLevel,
  ): EncryptionAlgorithm {
    switch (securityLevel) {
      case SecurityLevel.PUBLIC:
        return EncryptionAlgorithm.AES_128_CBC;
      case SecurityLevel.INTERNAL:
        return EncryptionAlgorithm.AES_256_CBC;
      case SecurityLevel.CONFIDENTIAL:
        return EncryptionAlgorithm.AES_128_GCM;
      case SecurityLevel.SECRET:
        return EncryptionAlgorithm.AES_256_GCM;
      case SecurityLevel.TOP_SECRET:
        return EncryptionAlgorithm.CHACHA20_POLY1305;
      default:
        return EncryptionAlgorithm.AES_256_GCM;
    }
  }

  /**
   * Gets key size for security level.
   */
  private getKeySizeForSecurityLevel(securityLevel: SecurityLevel): number {
    switch (securityLevel) {
      case SecurityLevel.PUBLIC:
        return 16; // 128 bits
      case SecurityLevel.INTERNAL:
      case SecurityLevel.CONFIDENTIAL:
      case SecurityLevel.SECRET:
      case SecurityLevel.TOP_SECRET:
        return 32; // 256 bits
      default:
        return 32;
    }
  }

  /**
   * Encrypts data based on security level.
   */
  public async encrypt(
    data: string,
    keyId: string,
    securityLevel: SecurityLevel,
    additionalData?: string,
  ): Promise<{ encrypted: string; metadata: EncryptionMetadata }> {
    const algorithm = this.getAlgorithmForSecurityLevel(securityLevel);
    const salt = randomBytes(32);
    const key = await this.deriveKey(keyId, securityLevel, salt);
    const iv = randomBytes(16);

    let encrypted: Buffer;
    let tag: Buffer | undefined;
    const timestamp = new Date().toISOString();

    try {
      if (algorithm.includes("gcm")) {
        // AEAD mode (GCM)
        const cipher = createCipheriv(algorithm, key, iv) as any;

        if (additionalData) {
          (cipher as any).setAAD(Buffer.from(additionalData));
        }

        encrypted = Buffer.concat([
          cipher.update(data, "utf8"),
          cipher.final(),
        ]);
        tag = (cipher as any).getAuthTag();
      } else if (algorithm === EncryptionAlgorithm.CHACHA20_POLY1305) {
        // ChaCha20-Poly1305 (simplified implementation)
        const cipher = createCipheriv(algorithm, key, iv) as any;

        if (additionalData) {
          cipher.setAAD(Buffer.from(additionalData));
        }

        encrypted = Buffer.concat([
          cipher.update(data, "utf8"),
          cipher.final(),
        ]);
        tag = (cipher as any).getAuthTag();
      } else {
        // CBC mode
        const cipher = createCipheriv(algorithm, key, iv);
        encrypted = Buffer.concat([
          cipher.update(data, "utf8"),
          cipher.final(),
        ]);
      }

      // Calculate checksum
      const checksum = createHash("sha256")
        .update(encrypted)
        .update(key)
        .update(iv)
        .digest("hex");

      const metadata: EncryptionMetadata = {
        algorithm,
        keyId,
        keyVersion: 1, // Simplified - should track actual key version
        iv: iv.toString("hex"),
        ...(tag && { tag: tag.toString("hex") }),
        kdf: KeyDerivationFunction.SCRYPT,
        salt: salt.toString("hex"),
        encryptedAt: timestamp,
        checksum,
      };

      return {
        encrypted: encrypted.toString("hex"),
        metadata,
      };
    } catch (error) {
      throw new Error(
        `Encryption failed: ${error instanceof Error ? error.message : "Unknown error"}`,
      );
    }
  }

  /**
   * Decrypts data using provided metadata.
   */
  public async decrypt(
    encryptedData: string,
    metadata: EncryptionMetadata,
    additionalData?: string,
  ): Promise<string> {
    try {
      const encrypted = Buffer.from(encryptedData, "hex");
      const iv = Buffer.from(metadata.iv, "hex");
      const salt = metadata.salt
        ? Buffer.from(metadata.salt, "hex")
        : undefined;

      // Determine security level from algorithm
      const securityLevel = this.inferSecurityLevelFromAlgorithm(
        metadata.algorithm,
      );
      const key = await this.deriveKey(metadata.keyId, securityLevel, salt);

      let decrypted: Buffer;

      if (metadata.algorithm.includes("gcm")) {
        // AEAD mode (GCM)
        const decipher = createDecipheriv(metadata.algorithm, key, iv) as any;

        if (metadata.tag) {
          decipher.setAuthTag(Buffer.from(metadata.tag, "hex"));
        }

        if (additionalData) {
          decipher.setAAD(Buffer.from(additionalData));
        }

        decrypted = Buffer.concat([
          decipher.update(encrypted),
          decipher.final(),
        ]);
      } else if (metadata.algorithm === EncryptionAlgorithm.CHACHA20_POLY1305) {
        // ChaCha20-Poly1305
        const decipher = createDecipheriv(metadata.algorithm, key, iv) as any;

        if (metadata.tag) {
          decipher.setAuthTag(Buffer.from(metadata.tag, "hex"));
        }

        if (additionalData) {
          decipher.setAAD(Buffer.from(additionalData), {
            plaintextLength: encrypted.length,
          });
        }

        decrypted = Buffer.concat([
          decipher.update(encrypted),
          decipher.final(),
        ]);
      } else {
        // CBC mode
        const decipher = createDecipheriv(metadata.algorithm, key, iv);
        decrypted = Buffer.concat([
          decipher.update(encrypted),
          decipher.final(),
        ]);
      }

      // Verify integrity with checksum
      const calculatedChecksum = createHash("sha256")
        .update(encrypted)
        .update(key)
        .update(iv)
        .digest("hex");

      if (calculatedChecksum !== metadata.checksum) {
        throw new Error(
          "Integrity verification failed - data may be corrupted",
        );
      }

      return decrypted.toString("utf8");
    } catch (error) {
      throw new Error(
        `Decryption failed: ${error instanceof Error ? error.message : "Unknown error"}`,
      );
    }
  }

  /**
   * Infers security level from encryption algorithm.
   */
  private inferSecurityLevelFromAlgorithm(
    algorithm: EncryptionAlgorithm,
  ): SecurityLevel {
    switch (algorithm) {
      case EncryptionAlgorithm.AES_128_CBC:
        return SecurityLevel.PUBLIC;
      case EncryptionAlgorithm.AES_256_CBC:
        return SecurityLevel.INTERNAL;
      case EncryptionAlgorithm.AES_128_GCM:
        return SecurityLevel.CONFIDENTIAL;
      case EncryptionAlgorithm.AES_256_GCM:
        return SecurityLevel.SECRET;
      case EncryptionAlgorithm.CHACHA20_POLY1305:
        return SecurityLevel.TOP_SECRET;
      default:
        return SecurityLevel.INTERNAL;
    }
  }

  /**
   * Generates a new encryption key.
   */
  public generateKey(
    keyId: string,
    securityLevel: SecurityLevel,
    algorithm?: EncryptionAlgorithm,
  ): EncryptionKey {
    const selectedAlgorithm =
      algorithm || this.getAlgorithmForSecurityLevel(securityLevel);
    const keySize = this.getKeySizeForSecurityLevel(securityLevel);

    return {
      id: keyId,
      version: 1,
      securityLevel,
      algorithm: selectedAlgorithm,
      keySize: keySize * 8, // Convert to bits
      usage: "both",
      status: "active",
      createdAt: new Date().toISOString(),
      rotationSchedule: {
        intervalDays: this.getRotationInterval(securityLevel),
        nextRotation: new Date(
          Date.now() +
            this.getRotationInterval(securityLevel) * 24 * 60 * 60 * 1000,
        ).toISOString(),
      },
    };
  }

  /**
   * Gets rotation interval for security level.
   */
  private getRotationInterval(securityLevel: SecurityLevel): number {
    switch (securityLevel) {
      case SecurityLevel.PUBLIC:
        return 365; // 1 year
      case SecurityLevel.INTERNAL:
        return 180; // 6 months
      case SecurityLevel.CONFIDENTIAL:
        return 90; // 3 months
      case SecurityLevel.SECRET:
        return 30; // 1 month
      case SecurityLevel.TOP_SECRET:
        return 7; // 1 week
      default:
        return 90;
    }
  }

  /**
   * Registers an encryption key.
   */
  public registerKey(key: EncryptionKey): void {
    this.keys.set(key.id, key);
  }

  /**
   * Gets an encryption key by ID.
   */
  public getKey(keyId: string): EncryptionKey | undefined {
    return this.keys.get(keyId);
  }

  /**
   * Lists all registered keys.
   */
  public listKeys(): EncryptionKey[] {
    return Array.from(this.keys.values());
  }

  /**
   * Rotates an encryption key.
   */
  public async rotateKey(keyId: string): Promise<EncryptionKey> {
    const existingKey = this.keys.get(keyId);
    if (!existingKey) {
      throw new Error(`Key ${keyId} not found`);
    }

    // Mark existing key as deprecated
    existingKey.status = "deprecated";

    // Create new key version
    const newKey: EncryptionKey = {
      ...existingKey,
      version: existingKey.version + 1,
      status: "active",
      createdAt: new Date().toISOString(),
      rotationSchedule: {
        intervalDays:
          existingKey.rotationSchedule?.intervalDays ||
          this.getRotationInterval(existingKey.securityLevel),
        nextRotation: new Date(
          Date.now() +
            (existingKey.rotationSchedule?.intervalDays ||
              this.getRotationInterval(existingKey.securityLevel)) *
              24 *
              60 *
              60 *
              1000,
        ).toISOString(),
      },
    };

    this.keys.set(keyId, newKey);
    return newKey;
  }

  /**
   * Validates encryption metadata.
   */
  public validateMetadata(metadata: EncryptionMetadata): boolean {
    try {
      // Check required fields
      if (
        !metadata.algorithm ||
        !metadata.keyId ||
        !metadata.iv ||
        !metadata.checksum
      ) {
        return false;
      }

      // Check algorithm validity
      if (!Object.values(EncryptionAlgorithm).includes(metadata.algorithm)) {
        return false;
      }

      // Check IV length (should be 16 bytes for most algorithms)
      if (Buffer.from(metadata.iv, "hex").length !== 16) {
        return false;
      }

      // Check checksum format (64 character hex string for SHA-256)
      if (
        metadata.checksum.length !== 64 ||
        !/^[a-f0-9]{64}$/i.test(metadata.checksum)
      ) {
        return false;
      }

      // Check tag for AEAD ciphers
      if (
        metadata.algorithm.includes("gcm") ||
        metadata.algorithm === EncryptionAlgorithm.CHACHA20_POLY1305
      ) {
        if (!metadata.tag || Buffer.from(metadata.tag, "hex").length !== 16) {
          return false;
        }
      }

      return true;
    } catch {
      return false;
    }
  }
}
