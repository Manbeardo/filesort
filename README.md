# filesort

Example solutions for https://www.reddit.com/r/golang/comments/1jcnqfi/how_the_hell_do_i_make_this_go_program_faster/

How to use:

- download the example dataset  
  ```curl -OL https://github.com/brannondorsey/naive-hashcat/releases/download/data/rockyou.txt```
- run the benchmarks  
  ```go test -bench=. ./...```

Example output:

```
goos: linux
goarch: amd64
pkg: github.com/Manbeardo/filesort
cpu: Intel(R) Core(TM) Ultra 7 265K
BenchmarkBuffered-20                   1        7192006209 ns/op
BenchmarkNaive-20                      1        6696147768 ns/op
BenchmarkNoMap-20                      1        2385805439 ns/op
BenchmarkNoStrAlloc-20                 1        3060783581 ns/op
BenchmarkStreamed-20                   1        1040885412 ns/op
PASS
ok      github.com/Manbeardo/filesort   20.705s
```
