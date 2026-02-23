#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
char *next_token(char **command)
{
    char *token;

    while ((token = strsep(command, " \t\n")) && !*token)
        ;

    return token;
}

int main()
{
    char *command = NULL;
    while (1)
    {
        char *token;
        char *args[100];
        int i = 0;
        printf("wshell> ");
        size_t command_size = 30;

        getline(&command, &command_size, stdin);
        while ((token = next_token(&command)) != NULL)
        {
            args[i++] = token;
        }
        args[i] = NULL;

        if (args[0] == NULL)
            continue;
        if (strcmp(args[0], "exit") == 0)
            break;
        if (strcmp(args[0], "cd") == 0)
        {
            if (args[1] == NULL)
            {
                fprintf(stderr, "cd: expected argument\n");
            }
            else
            {
                if (chdir(args[1]) != 0)
                {
                    perror("cd");
                }
            }
            continue;
        }
        pid_t pid = fork();
        if (pid < 0)
        {
            perror("fork");
            exit(EXIT_FAILURE);
        }
        if (pid == 0)
        {
            execvp(args[0], args);
            perror("command execution failed");
            exit(EXIT_FAILURE);
        }
        else
        {
            wait(NULL);
        }
    }

    free(command);
}