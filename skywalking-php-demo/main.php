<?php

use Swoole\Http\Request;
use Swoole\Http\Response;

$http = new Swoole\Http\Server('0.0.0.0', 10000);

$http->on('request', function (Request $request, Response $response) {
    $method = $request->server['request_method'];
    $path = $request->server['request_uri'];

    switch ($method) {
        case 'GET':
            handleGetRequest($request, $response, $path);
            break;
        case 'POST':
            handlePostRequest($request, $response, $path);
            break;
        default:
            $response->status(400);
            $response->end('Bad Request');
            break;
    }
});

function handleGetRequest(Request $request, Response $response, string $path)
{
    switch ($path) {
        case '/':
            $response->status(200);
            $response->header('Content-Type', 'text/plain');
            $response->end('Hello, World!');
            break;
        case '/ping':
            $response->status(200);
            $response->header('Content-Type', 'text/plain');
            $response->end('Pong');
            break;
        default:
            $response->status(404);
            $response->end('Not Found');
            break;
    }
}

function handlePostRequest(Request $request, Response $response, string $path)
{
    switch ($path) {
        case '/echo':
            $response->status(200);
            $response->header('Content-Type', 'text/plain');
            $response->end($request->rawContent());
            break;
        default:
            $response->status(404);
            $response->end('Not Found');
            break;
    }
}

$http->start();