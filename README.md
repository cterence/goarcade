# space-invaders

## Profiling

- Use the `-p` flag to start profiling webserver when running the emulator

  ```bash
  ./space-invaders -p [...]
  ```

- Use `go tool pprof` to start profiling

  ```bash
  go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
  ```

- Use `top`, `web` or `png` commands in the repl to explore the results
