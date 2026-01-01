# Random Walk

Minimal graphics programming project showing multiple agents performing a random walk using the Ebiten 2D game library for Go.
This was highly inspired by Daniel Hirsch's [video](https://www.youtube.com/watch?v=ErA4U9WqNCE&t=2575s)video. Daniel is a Youtuber I just discovered and his videos are a gem!

## Overview

- Description: A simple visualization where several agents start in the center of a 640×480 window and perform a random walk, leaving colored trails.

## Prerequisites

- Go (1.20+ recommended).
- Ebiten will be fetched automatically from the module when you run or build.

## Run

From the repository root or the `randomwalk` directory:

```bash
cd randomwalk
go run .
```

You can pass an optional single integer argument to set the number of agents:

```bash
go run . 10
```

## Build

```bash
go build -o randomwalk-app
./randomwalk-app 8
```

## Behavior / Usage notes

- Window size: 640×480.
- Agents start at the center (320,240) and leave persistent trails on the canvas.

## Troubleshooting

- If dependencies are missing or build fails, run:

```bash
go mod tidy
```

## Files

- `randomwalk/main.go` — main demo.
- `randomwalk/go.mod` — module file.
