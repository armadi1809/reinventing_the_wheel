#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/wait.h>

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
    cmd->redirect_input[0] = -1;
    cmd->redirect_input[1] = -1;
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

void close_all_pipes(int n_pipes, int (*pipes)[2])
{
    for (int i = 0; i < n_pipes; ++i)
    {
        close(pipes[i][0]);
        close(pipes[i][1]);
    }
}

pid_t run_with_redirection(Command *cmd, int n_pipes, int pipes[][2])
{

    if (cmd->args[0] == NULL)
        return -1;
    if (strcmp(cmd->args[0], "exit") == 0)
        exit(EXIT_SUCCESS);
    if (strcmp(cmd->args[0], "cd") == 0)
    {
        if (cmd->args[1] == NULL)
        {
            fprintf(stderr, "cd: expected argument\n");
        }
        else
        {
            if (chdir(cmd->args[1]) != 0)
            {
                perror("cd");
            }
        }
        return -1;
    }
    pid_t pid = fork();
    if (pid < 0)
    {
        perror("fork");
        exit(EXIT_FAILURE);
    }
    if (pid == 0)
    {
        if (cmd->redirect_input[0] != -1)
        {
            dup2(cmd->redirect_input[0], STDIN_FILENO);
            close(cmd->redirect_input[0]);
        }
        if (cmd->redirect_input[1] != -1)
        {
            dup2(cmd->redirect_input[1], STDOUT_FILENO);
            close(cmd->redirect_input[1]);
        }

        close_all_pipes(n_pipes, pipes);
        execvp(cmd->program_name, cmd->args);
        perror("command execution failed");
        exit(EXIT_FAILURE);
    }
    return pid;
}

int main()
{
    char *command = NULL;
    Pipeline *pipeline;

    while (1)
    {
        printf("wshell> ");
        size_t command_size = 30;

        getline(&command, &command_size, stdin);

        pipeline = parsePipeline(command);

        int num_pipes = pipeline->num_commands - 1;
        int pipes[num_pipes][2];

        for (int i = 1; i < pipeline->num_commands; ++i)
        {
            pipe(pipes[i - 1]);
            pipeline->commands[i]->redirect_input[STDIN_FILENO] = pipes[i - 1][0];
            pipeline->commands[i - 1]->redirect_input[STDOUT_FILENO] = pipes[i - 1][1];
        }

        for (int i = 0; i < pipeline->num_commands; ++i)
        {
            run_with_redirection(pipeline->commands[i], num_pipes, pipes);
        }

        close_all_pipes(num_pipes, pipes);

        for (int i = 0; i < pipeline->num_commands; ++i)
        {
            wait(NULL);
        }
    }

    free(command);
    free(pipeline);
    return 0;
}