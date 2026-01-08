<?php

declare(strict_types=1);

namespace AetherVault\Capability;

use AetherVault\Exception\VaultExpiredCapabilityException;
use AetherVault\Exception\VaultAccessDeniedException;

final class SmtpAccess extends AbstractCapability
{
    private bool $readOnly = false;

    public function readOnly(): self
    {
        $this->readOnly = true;
        return $this;
    }

    public function connect(): SmtpConnection
    {
        $this->ensureNotExpired();

        $request = array_merge($this->buildRequest(), [
            'capability' => 'smtp',
            'read_only' => $this->readOnly,
        ]);

        $response = $this->transport->request('POST', '/v1/capabilities/smtp', [], json_encode($request));

        if ($response['status'] === 403) {
            throw new VaultAccessDeniedException('SMTP access denied');
        }

        if ($response['status'] === 410) {
            throw new VaultExpiredCapabilityException('SMTP capability has expired');
        }

        if ($response['status'] !== 200) {
            throw new \RuntimeException('Failed to obtain SMTP capability');
        }

        $data = json_decode($response['body'], true);
        return new SmtpConnection($data['credentials'], $this->ttl);
    }
}