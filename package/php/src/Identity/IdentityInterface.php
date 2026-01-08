<?php

declare(strict_types=1);

namespace AetherVault\Identity;

interface IdentityInterface
{
    public function getCredentials(): array;
    public function getType(): string;
}