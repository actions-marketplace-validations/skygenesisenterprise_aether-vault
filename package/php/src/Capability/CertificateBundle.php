<?php

declare(strict_types=1);

namespace AetherVault\Capability;

final class CertificateBundle
{
    private string $certificate;
    private string $privateKey;
    private array $caChain;
    private ?int $expiresAt;

    public function __construct(string $certificate, string $privateKey, array $caChain, ?int $ttl = null)
    {
        $this->certificate = $certificate;
        $this->privateKey = $privateKey;
        $this->caChain = $caChain;
        $this->expiresAt = $ttl ? time() + $ttl : null;
    }

    public function getCertificate(): string
    {
        $this->ensureValid();
        return $this->certificate;
    }

    public function getPrivateKey(): string
    {
        $this->ensureValid();
        return $this->privateKey;
    }

    public function getCaChain(): array
    {
        $this->ensureValid();
        return $this->caChain;
    }

    public function getFullChain(): string
    {
        $this->ensureValid();
        return $this->certificate . "\n" . implode("\n", $this->caChain);
    }

    public function isValid(): bool
    {
        return $this->expiresAt === null || time() < $this->expiresAt;
    }

    private function ensureValid(): void
    {
        if (!$this->isValid()) {
            throw new \RuntimeException('Certificate bundle has expired');
        }
    }
}