#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

typedef struct
{
    char *program_name;
    int redirect_input[2];
    // flexible array member has to be the last member of struct for compiler to calculate the size of struct correctly
    char *args[];
} Command;

typedef struct
{
    int num_commands;
    Command *commands[];
} Pipeline;

char *next_token(char **command)
{
    char *token;

    while ((token = strsep(command, " \t\n")) && !*token)
        ;

    return token;
}

void fork_and_exeute_with_io(int stdin_fd, int stdout_fd, char *args[])
{
    pid_t pid = fork();
    if (pid < 0)
    {
        perror("fork");
        exit(EXIT_FAILURE);
    }
    if (pid == 0)
    {
        dup2(stdin_fd, STDIN_FILENO);
        dup2(stdout_fd, STDOUT_FILENO);
        execvp(args[0], args);
        perror("command execution failed");
        exit(EXIT_FAILURE);
    }
    else
    {
        wait(NULL);
    }
}

Command *parseCommand(char *command_str)
{
    char *token;
    char *args[100];
    int i = 0;

    while ((token = next_token(&command_str)) != NULL)
    {
        args[i++] = token;
    }
    args[i] = NULL;

    Command *cmd = malloc(sizeof(Command) + (i + 1) * sizeof(char *));
    cmd->program_name = args[0];
    for (int j = 0; j <= i; j++)
    {
        cmd->args[j] = args[j];
    }
    return cmd;
}

Pipeline *parsePipeline(char *line)
{
    char *command_str;
    int num_commands = 0;
    Command *commands[100];

    while ((command_str = strsep(&line, "|")) != NULL)
    {
        commands[num_commands++] = parseCommand(command_str);
    }

    Pipeline *pipeline = malloc(sizeof(Pipeline) + num_commands * sizeof(Command *));
    pipeline->num_commands = num_commands;
    for (int i = 0; i < num_commands; i++)
    {
        pipeline->commands[i] = commands[i];
    }
    return pipeline;
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