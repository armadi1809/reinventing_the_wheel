.PHONY: clean run

run: wshell
	./wshell

clean:
	rm -f wshell

wshell: main.c
	gcc -Wall main.c -o wshell