# space-invaders

## Test results

### TST8080.COM
<!-- TST8080.COM -->
```txt
panic: failed to init sdl: No available video device

goroutine 1 [running, locked to thread]:
github.com/cterence/space-invaders/internal/arcade/ui.(*UI).Init(0xc000190000)
	/home/runner/work/space-invaders/space-invaders/internal/arcade/ui/ui.go:30 +0x265
github.com/cterence/space-invaders/internal/arcade.Run({0xc0000122e8?, 0xc0000b3a30?}, {0xc00003cac0, 0x1, 0xe?}, {0xc0000b3a00, 0x5, 0x0?})
	/home/runner/work/space-invaders/space-invaders/internal/arcade/arcade.go:94 +0x21e
main.main.func2({0x9d8508, 0xc00009d5f0}, 0xc00009d5f0?)
	/home/runner/work/space-invaders/space-invaders/main.go:74 +0x17b
github.com/urfave/cli/v3.(*Command).run(0xc0000ca2c8, {0x9d8508, 0xc00009d5f0}, {0xc000024040, 0x2, 0x2})
	/home/runner/go/pkg/mod/github.com/urfave/cli/v3@v3.6.1/command_run.go:354 +0x28ec
github.com/urfave/cli/v3.(*Command).Run(...)
	/home/runner/go/pkg/mod/github.com/urfave/cli/v3@v3.6.1/command_run.go:94
main.main()
	/home/runner/work/space-invaders/space-invaders/main.go:97 +0x6fb
```
<!-- /TST8080.COM -->

### 8080PRE.COM

<!-- 8080PRE.COM -->
```txt
panic: failed to init sdl: No available video device

goroutine 1 [running, locked to thread]:
github.com/cterence/space-invaders/internal/arcade/ui.(*UI).Init(0xc000110000)
	/home/runner/work/space-invaders/space-invaders/internal/arcade/ui/ui.go:30 +0x265
github.com/cterence/space-invaders/internal/arcade.Run({0xc0000122e8?, 0xc0000b3a30?}, {0xc00003cac0, 0x1, 0xe?}, {0xc0000b3a00, 0x5, 0x0?})
	/home/runner/work/space-invaders/space-invaders/internal/arcade/arcade.go:94 +0x21e
main.main.func2({0x9d8508, 0xc00009d5f0}, 0xc00009d5f0?)
	/home/runner/work/space-invaders/space-invaders/main.go:74 +0x17b
github.com/urfave/cli/v3.(*Command).run(0xc0000ca2c8, {0x9d8508, 0xc00009d5f0}, {0xc000024040, 0x2, 0x2})
	/home/runner/go/pkg/mod/github.com/urfave/cli/v3@v3.6.1/command_run.go:354 +0x28ec
github.com/urfave/cli/v3.(*Command).Run(...)
	/home/runner/go/pkg/mod/github.com/urfave/cli/v3@v3.6.1/command_run.go:94
main.main()
	/home/runner/work/space-invaders/space-invaders/main.go:97 +0x6fb
```
<!-- /8080PRE.COM -->

### CPUTEST.COM

<!-- CPUTEST.COM -->
```txt
panic: failed to init sdl: No available video device

goroutine 1 [running, locked to thread]:
github.com/cterence/space-invaders/internal/arcade/ui.(*UI).Init(0xc00019e000)
	/home/runner/work/space-invaders/space-invaders/internal/arcade/ui/ui.go:30 +0x265
github.com/cterence/space-invaders/internal/arcade.Run({0xc0001242d0?, 0xc00013fa30?}, {0xc00011ea70, 0x1, 0xe?}, {0xc00013fa00, 0x5, 0x0?})
	/home/runner/work/space-invaders/space-invaders/internal/arcade/arcade.go:94 +0x21e
main.main.func2({0x9d8508, 0xc00011d5f0}, 0xc00011d5f0?)
	/home/runner/work/space-invaders/space-invaders/main.go:74 +0x17b
github.com/urfave/cli/v3.(*Command).run(0xc00015e2c8, {0x9d8508, 0xc00011d5f0}, {0xc00012a000, 0x2, 0x2})
	/home/runner/go/pkg/mod/github.com/urfave/cli/v3@v3.6.1/command_run.go:354 +0x28ec
github.com/urfave/cli/v3.(*Command).Run(...)
	/home/runner/go/pkg/mod/github.com/urfave/cli/v3@v3.6.1/command_run.go:94
main.main()
	/home/runner/work/space-invaders/space-invaders/main.go:97 +0x6fb
```
<!-- /CPUTEST.COM -->

### 8080EXM.COM

<!-- 8080EXM.COM -->
```txt
panic: failed to init sdl: No available video device

goroutine 1 [running, locked to thread]:
github.com/cterence/space-invaders/internal/arcade/ui.(*UI).Init(0xc00019e000)
	/home/runner/work/space-invaders/space-invaders/internal/arcade/ui/ui.go:30 +0x265
github.com/cterence/space-invaders/internal/arcade.Run({0xc0001262d0?, 0xc00013fa30?}, {0xc000120a70, 0x1, 0xe?}, {0xc00013fa00, 0x5, 0x0?})
	/home/runner/work/space-invaders/space-invaders/internal/arcade/arcade.go:94 +0x21e
main.main.func2({0x9d8508, 0xc00011f5f0}, 0xc00011f5f0?)
	/home/runner/work/space-invaders/space-invaders/main.go:74 +0x17b
github.com/urfave/cli/v3.(*Command).run(0xc00015e2c8, {0x9d8508, 0xc00011f5f0}, {0xc00012a000, 0x2, 0x2})
	/home/runner/go/pkg/mod/github.com/urfave/cli/v3@v3.6.1/command_run.go:354 +0x28ec
github.com/urfave/cli/v3.(*Command).Run(...)
	/home/runner/go/pkg/mod/github.com/urfave/cli/v3@v3.6.1/command_run.go:94
main.main()
	/home/runner/work/space-invaders/space-invaders/main.go:97 +0x6fb
```
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
