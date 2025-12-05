# space-invaders

## Test results

### TST8080.COM
<!-- TST8080.COM -->
<!-- /TST8080.COM -->

### 8080PRE.COM

<!-- 8080PRE.COM -->
<!-- /8080PRE.COM -->

### CPUTEST.COM

<!-- CPUTEST.COM -->
<!-- /CPUTEST.COM -->

### 8080EXM.COM

<!-- 8080EXM.COM -->
<!-- /8080EXM.COM -->

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
