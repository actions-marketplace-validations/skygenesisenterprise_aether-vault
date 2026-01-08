<?php

declare(strict_types=1);

namespace AetherVault\Capability;

use AetherVault\Exception\VaultExpiredCapabilityException;
use AetherVault\Exception\VaultAccessDeniedException;

final class TlsCertificate extends AbstractCapability
{
    private string $domain;

    public function forDomain(string $domain): self
    {
        $this->domain = $domain;
        return $this;
    }

    public function obtain(): CertificateBundle
    {
        $this->ensureNotExpired();

        if (!$this->domain) {
            throw new \InvalidArgumentException('Domain must be specified');
        }

        $request = array_merge($this->buildRequest(), [
            'capability' => 'tls_certificate',
            'domain' => $this->domain,
        ]);

        $response = $this->transport->request('POST', '/v1/capabilities/tls', [], json_encode($request));

        if ($response['status'] === 403) {
            throw new VaultAccessDeniedException('TLS certificate access denied');
        }

        if ($response['status'] === 410) {
            throw new VaultExpiredCapabilityException('TLS certificate capability has expired');
        }

        if ($response['status'] !== 200) {
            throw new \RuntimeException('Failed to obtain TLS certificate');
        }

        $data = json_decode($response['body'], true);
        return new CertificateBundle($data['certificate'], $data['private_key'], $data['ca_chain'], $this->ttl);
    }
}