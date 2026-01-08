<?php

declare(strict_types=1);

namespace AetherVault\Audit;

final class AuditTrail
{
    private array $events = [];

    public function log(string $action, array $context = []): void
    {
        $this->events[] = [
            'timestamp' => time(),
            'action' => $action,
            'context' => $context,
        ];
    }

    public function getEvents(): array
    {
        return $this->events;
    }

    public function clear(): void
    {
        $this->events = [];
    }
}