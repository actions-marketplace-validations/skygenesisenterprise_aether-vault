<?php

declare(strict_types=1);

namespace AetherVault\Client;

interface TransportInterface
{
    public function request(string $method, string $path, array $headers = [], ?string $body = null): array;
}