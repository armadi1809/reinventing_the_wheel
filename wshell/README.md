# wshell

Minimal Unix-like shell written in C as a learning exercise. The goal is to understand the mechanics behind command parsing, `fork`, `exec`, `wait`, pipes, and a couple of built-in shell commands by implementing a small shell from scratch.

## Overview

- Description: reads a command line from standard input, splits it into commands separated by `|`, creates the required pipes, and runs each stage with `fork` + `execvp`.
- Built-ins: `cd` and `exit`.

## Prerequisites

- A C compiler such as `gcc`.
- A Unix-like environment (tested layout assumes macOS or Linux).
- `make` for the provided Makefile.

## Build

From the `wshell` directory:

```bash
make
```

This compiles `main.c` into the `wshell` executable with:

```bash
gcc -Wall main.c -o wshell
```

## Run

```bash
cd wshell
make run
```

Or run the executable directly after building:

```bash
./wshell
```

You will get a prompt like:

```text
wshell>
```

## Supported usage

Simple commands:

```bash
pwd
ls
echo hello
```

Pipelines:

```bash
ls | wc
cat main.c | grep fork
ps aux | grep ssh
```

Built-ins:

```bash
cd ..
exit
```

## Current limitations

This is intentionally small and does not try to be a full POSIX shell. In its current form it does not support:

- quoting or escaping rules
- environment variable expansion
- command history
- output redirection like `>` or input redirection like `<`
- background jobs with `&`
- advanced parsing such as `&&`, `||`, subshells, or command substitution

## Files

- `main.c` — parser, built-in handling, pipe setup, and process execution.
- `Makefile` — build and run helpers.

## Clean

Remove the compiled binary with:

```bash
make clean
```
