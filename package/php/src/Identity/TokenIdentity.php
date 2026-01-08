<?php

declare(strict_types=1);

namespace AetherVault\Identity;

final class TokenIdentity implements IdentityInterface
{
    private string $token;

    public function __construct(string $token)
    {
        if (empty($token)) {
            throw new \InvalidArgumentException('Token cannot be empty');
        }
        $this->token = $token;
    }

    public function getCredentials(): array
    {
        return [
            'type' => 'token',
            'token' => $this->token,
        ];
    }

    public function getType(): string
    {
        return 'token';
    }
}