#include <stdio.h>
#include <stdlib.h>

int main()
{
    char *command = NULL;
    while (1)
    {
        printf("wshell> ");
        size_t command_size = 30;

        getline(&command, &command_size, stdin);

        printf("%s", command);
    }

    free(command);
}