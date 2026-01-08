<?php

declare(strict_types=1);

namespace AetherVault\Client;

use AetherVault\Exception\VaultTransportException;
use Psr\Http\Client\ClientInterface;
use Psr\Http\Message\RequestFactoryInterface;
use Psr\Http\Message\StreamFactoryInterface;

final class HttpTransport implements TransportInterface
{
    private string $endpoint;
    private ClientInterface $client;
    private RequestFactoryInterface $requestFactory;
    private StreamFactoryInterface $streamFactory;
    private int $timeout;

    public function __construct(
        string $endpoint,
        ClientInterface $client,
        RequestFactoryInterface $requestFactory,
        StreamFactoryInterface $streamFactory,
        int $timeout = 30
    ) {
        $this->endpoint = rtrim($endpoint, '/');
        $this->client = $client;
        $this->requestFactory = $requestFactory;
        $this->streamFactory = $streamFactory;
        $this->timeout = $timeout;
    }

    public function request(string $method, string $path, array $headers = [], ?string $body = null): array
    {
        try {
            $url = $this->endpoint . $path;
            $request = $this->requestFactory->createRequest($method, $url);

            foreach ($headers as $name => $value) {
                $request = $request->withHeader($name, $value);
            }

            if ($body !== null) {
                $request = $request->withBody($this->streamFactory->createStream($body));
            }

            $response = $this->client->sendRequest($request);
            
            return [
                'status' => $response->getStatusCode(),
                'headers' => $response->getHeaders(),
                'body' => (string) $response->getBody(),
            ];
        } catch (\Throwable $e) {
            throw new VaultTransportException('Transport request failed: ' . $e->getMessage(), 0, $e);
        }
    }
}