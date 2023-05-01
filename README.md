
# Proctree

This program simulates a process tree running in some virtual environment. I plan to use it
in my altnet repository.

```bash 
> go run cmd/app/main.go
```

Most lines you enter create a new `shell` process using the argument as the prompt.

Entering `x` exits the current shell and returns to the parent.

You can list processes with `ps` and kill processes with `kill <pid>`.
Note that killing a process will kill its whole process tree with all of its children.
E.g. `kill 1` exits the program immediately.

