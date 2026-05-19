.PHONY: build build-debug


build: main.c
	gcc -std=c99 -Wall main.c mpc.c -lreadline -lm -o lizp

build-debug: main.c
	gcc -std=c99 -Wall -fsanitize=address -g main.c mpc.c -lreadline -lm -o lizp_debug
