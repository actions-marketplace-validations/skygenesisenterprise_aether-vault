<?php

declare(strict_types=1);

namespace AetherVault\Capability;

use AetherVault\Client\TransportInterface;
use AetherVault\Context\Context;
use AetherVault\Identity\IdentityInterface;

abstract class AbstractCapability
{
    protected TransportInterface $transport;
    protected IdentityInterface $identity;
    protected Context $context;
    protected ?int $ttl = null;
    protected array $constraints = [];

    public function __construct(
        TransportInterface $transport,
        IdentityInterface $identity,
        Context $context
    ) {
        $this->transport = $transport;
        $this->identity = $identity;
        $this->context = $context;
    }

    public function ttl(int $seconds): self
    {
        if ($seconds <= 0 || $seconds > 300) {
            throw new \InvalidArgumentException('TTL must be between 1 and 300 seconds');
        }

        $this->ttl = $seconds;
        return $this;
    }

    public function withConstraint(string $key, string $value): self
    {
        $this->constraints[$key] = $value;
        return $this;
    }

    protected function buildRequest(): array
    {
        return [
            'identity' => $this->identity->getCredentials(),
            'context' => $this->context->toArray(),
            'ttl' => $this->ttl ?? 60,
            'constraints' => $this->constraints,
            'timestamp' => time(),
        ];
    }

    protected function ensureNotExpired(): void
    {
        if ($this->ttl && $this->ttl <= 0) {
            throw new \RuntimeException('Capability has expired');
        }
    }
}