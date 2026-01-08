<?php

declare(strict_types=1);

namespace AetherVault\Capability;

final class DatabaseConnection
{
    private string $connectionString;
    private ?int $expiresAt;

    public function __construct(string $connectionString, ?int $ttl = null)
    {
        $this->connectionString = $connectionString;
        $this->expiresAt = $ttl ? time() + $ttl : null;
    }

    public function getConnectionString(): string
    {
        $this->ensureValid();
        return $this->connectionString;
    }

    public function isValid(): bool
    {
        return $this->expiresAt === null || time() < $this->expiresAt;
    }

    private function ensureValid(): void
    {
        if (!$this->isValid()) {
            throw new \RuntimeException('Database connection has expired');
        }
    }
}