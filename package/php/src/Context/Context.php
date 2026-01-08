<?php

declare(strict_types=1);

namespace AetherVault\Context;

final class Context
{
    private string $hostname;
    private string $appName;
    private string $environment;
    private array $metadata;

    private function __construct(string $hostname, string $appName, string $environment, array $metadata = [])
    {
        $this->hostname = $hostname;
        $this->appName = $appName;
        $this->environment = $environment;
        $this->metadata = $metadata;
    }

    public static function auto(): self
    {
        $hostname = gethostname() ?: 'unknown';
        $appName = $_ENV['APP_NAME'] ?? $_SERVER['APP_NAME'] ?? 'aether-app';
        $environment = $_ENV['APP_ENV'] ?? $_SERVER['APP_ENV'] ?? 'development';

        return new self($hostname, $appName, $environment);
    }

    public static function create(string $hostname, string $appName, string $environment, array $metadata = []): self
    {
        return new self($hostname, $appName, $environment, $metadata);
    }

    public function withMetadata(string $key, string $value): self
    {
        $new = clone $this;
        $new->metadata[$key] = $value;
        return $new;
    }

    public function getHostname(): string
    {
        return $this->hostname;
    }

    public function getAppName(): string
    {
        return $this->appName;
    }

    public function getEnvironment(): string
    {
        return $this->environment;
    }

    public function getMetadata(): array
    {
        return $this->metadata;
    }

    public function toArray(): array
    {
        return [
            'hostname' => $this->hostname,
            'app_name' => $this->appName,
            'environment' => $this->environment,
            'metadata' => $this->metadata,
        ];
    }
}