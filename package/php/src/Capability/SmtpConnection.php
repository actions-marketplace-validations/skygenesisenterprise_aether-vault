<?php

declare(strict_types=1);

namespace AetherVault\Capability;

final class SmtpConnection
{
    private array $credentials;
    private ?int $expiresAt;

    public function __construct(array $credentials, ?int $ttl = null)
    {
        $this->credentials = $credentials;
        $this->expiresAt = $ttl ? time() + $ttl : null;
    }

    public function getHost(): string
    {
        $this->ensureValid();
        return $this->credentials['host'];
    }

    public function getPort(): int
    {
        $this->ensureValid();
        return $this->credentials['port'];
    }

    public function getUsername(): string
    {
        $this->ensureValid();
        return $this->credentials['username'];
    }

    public function getPassword(): string
    {
        $this->ensureValid();
        return $this->credentials['password'];
    }

    public function getEncryption(): ?string
    {
        $this->ensureValid();
        return $this->credentials['encryption'] ?? null;
    }

    public function isValid(): bool
    {
        return $this->expiresAt === null || time() < $this->expiresAt;
    }

    private function ensureValid(): void
    {
        if (!$this->isValid()) {
            throw new \RuntimeException('SMTP connection has expired');
        }
    }
}