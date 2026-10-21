# Benchmarks

This document records the benchmark methodology and the measured results. Performance
is a release concern, but correctness and a conservative parser are not traded away for
a benchmark number.

## Environment

| Item | Value |
|---|---|
| CPU | 12th Gen Intel Core i5-12450HX |
| OS | Linux amd64 |
| Go | 1.27.2 |
| `GOMAXPROCS` | 12 (default) |
| Parser | `DefaultParser` |
| Benchmark files | `bench_test.go`, `uahttp/bench_test.go` |

Results vary with CPU frequency scaling, compiler version, and background load. Treat
the numbers below as a reproducible local snapshot, not a universal guarantee.

## Command

```sh
go test -run '^$' -bench . -benchmem -count=3 ./...
```

The benchmarks use representative fixtures:

- **Desktop:** Chrome 120 on Windows 10 with `Win64; x64`
- **Mobile:** Chrome 120 on Android 13 with a Pixel 7 Pro and `wv`
- **Bot:** Googlebot 2.1
- **Parallel:** the desktop fixture through `b.RunParallel`

## Baseline

The checked-in reference snapshot is `testdata/bench/baseline.txt`. It covers the core
parser and the `uahttp` integration. It is a local snapshot for comparison, not a
release gate. Generate a fresh result and compare with:

```sh
go test -run '^$' -bench . -benchmem -count=5 ./... > bench-new.txt
go install golang.org/x/perf/cmd/benchstat@latest
"$(go env GOPATH)/bin/benchstat" testdata/bench/baseline.txt bench-new.txt
```

## Results

Representative snapshot from 2026-10-21 with Go 1.27.2:

| Benchmark | Time | Bytes/op | Allocs/op |
|---|---:|---:|---:|
| Desktop Chrome | 1.75 µs | 0 B | 0 |
| Android mobile | 4.77 µs | 176 B | 1 |
| Googlebot | 2.41 µs | 80 B | 1 |
| Desktop, parallel | 0.34 µs | 0 B | 0 |

The mobile path uses one allocation: the lazily built lowercase raw User-Agent used by
Android model/vendor fallback checks. The bot path uses the same lazy copy.

HTTP integration snapshot:

| Benchmark | Time | Bytes/op | Allocs/op |
|---|---:|---:|---:|
| `BenchmarkParseRequest` | 2.90 µs | 256 B | 10 |
| `BenchmarkParseClientHints` | 0.85 µs | 256 B | 10 |

## Targets

| Metric | Target | Result | Status |
|---|---:|---:|---|
| Allocations per parse | ≤4 | 0 desktop, 1 mobile, 1 bot | Met |
| Bytes per parse | ≤512 B | 0–176 B | Met |

The parser performs a bounded token scan through the tokenizer and ordered detectors.
The desktop path is allocation-free. The mobile and bot paths retain one allocation for
the lazily built lowercase raw User-Agent used by Android model/vendor fallback checks.

## Time profile

The parser carries broad detector coverage: OS version extraction, CPU normalization,
bot categories, app/WebView state, TV/console/wearable classification, Android
vendor/model cleanup, client hints, and custom rule dispatch. The desktop path is
already allocation-free, but it still performs a bounded token scan through the
tokenizer and ordered detectors.

The remaining time is spread across:

- tokenization and comment/product scanning
- ordered browser vendor checks
- OS/CPU/app/engine detection
- final version normalization

Further gains would require a more aggressive fast path, cached parse plans for common
families, or a narrower default detector set. Those changes risk correctness and
readability, so they are deferred rather than hidden.

## Benchmark discipline

- Keep `bench_test.go` fixtures stable; changing a fixture invalidates comparisons.
- Run benchmarks on the same host and with the same Go version when comparing changes.
- Use multiple `-count` runs and report a range or median.
- Add a benchmark for any change that touches tokenization, matching, or detector order.
- Never trade away a regression test to improve a benchmark.
- Record the allocation count, not just wall-clock time.
- Do not compare parallel results directly with sequential results.

## Profiling

CPU profile:

```sh
go test -run '^$' -bench '^BenchmarkParseDesktop$' \
  -benchtime=2s -cpuprofile=cpu.out
go tool pprof -top -cum cpu.out
```

Allocation profile:

```sh
go test -run '^$' -bench '^BenchmarkParseMobile$' \
  -benchtime=2s -memprofile=mem.out
go tool pprof -top -alloc_objects mem.out
```

Profiles are local artifacts and are ignored by `.gitignore`.

## Future work

- Reduce the remaining lazy lowercase allocation on Android and bot paths.
- Measure a guarded fast path for common desktop Chromium clients without weakening
  unknown-input behavior.
- Add benchmark fixtures for TV and console clients.
- Track results across releases with `benchstat` once a stable CI benchmark runner is
  available.

Until a faster path is proven and fully regression-tested, the current result is:
**allocation target met, with the remaining time spent in bounded token and detector
scans**.
