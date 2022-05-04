# Hex MathRush

MathRush is simple math game which greatly inspired by [this game](https://apps.apple.com/sa/app/1-2-3/id953831664).

It is build using client & server architecture. The client communicate with server using REST API defined in [http_api.md](./docs/http_api.md).

It is intended for showcasing implementation of [Hexagonal Architecture](https://medium.com/@matiasvarela/hexagonal-architecture-in-go-cfd4e436faa3). Hence the name `Hex MathRush`.

## Game Flow

The game flow itself is pretty simple, basically user just need to answer the math problem quickly & correctly to make the score higher. Game will end when either user is too slow answering the question or choose the wrong answer.

Here is the flowchart of the game:

<p align="center">
    <img src="./docs/game_flow.svg" alt="Game Flow" height="500" />
</p>

## How to Run

This app is powered by docker. So make sure to install it before running below command:

```bash
> make run-mem-server
```

Upon success, your console should output message like following:

```bash
2022/04/01 15:53:42 server is listening on :9190...
```
