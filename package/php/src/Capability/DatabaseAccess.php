<?php

declare(strict_types=1);

namespace AetherVault\Capability;

use AetherVault\Exception\VaultExpiredCapabilityException;
use AetherVault\Exception\VaultAccessDeniedException;

final class DatabaseAccess extends AbstractCapability
{
    private string $type;
    private bool $readOnly = false;

    public function postgres(): self
    {
        $this->type = 'postgres';
        return $this;
    }

    public function mysql(): self
    {
        $this->type = 'mysql';
        return $this;
    }

    public function readOnly(): self
    {
        $this->readOnly = true;
        return $this;
    }

    public function connect(): DatabaseConnection
    {
        $this->ensureNotExpired();

        if (!$this->type) {
            throw new \InvalidArgumentException('Database type must be specified');
        }

        $request = array_merge($this->buildRequest(), [
            'capability' => 'database',
            'type' => $this->type,
            'read_only' => $this->readOnly,
        ]);

        $response = $this->transport->request('POST', '/v1/capabilities/database', [], json_encode($request));

        if ($response['status'] === 403) {
            throw new VaultAccessDeniedException('Database access denied');
        }

        if ($response['status'] === 410) {
            throw new VaultExpiredCapabilityException('Database capability has expired');
        }

        if ($response['status'] !== 200) {
            throw new \RuntimeException('Failed to obtain database capability');
        }

        $data = json_decode($response['body'], true);
        return new DatabaseConnection($data['connection_string'], $this->ttl);
    }
}