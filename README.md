## Running

go build && ./tyrande.exe

## OpenCV 

Download opencv and mingw, following installation instructions on github.


Use 64 bit vc11 libs and dlls.
Replace the version string in the include statements.
Remove the .dll in the include statements
Put the files in: C:\mingw\mingw64\x86_64-w64-mingw32\lib
Comment out method bodies in cvaux.go


Make sure the script compiles.
```
cd Documents\code\go\src\github.com\lazywei\go-opencv\samples && go run hellocv.go
```

### Profiling

Instructions taken from here: https://medium.com/@hackintoshrao/daily-code-optimization-using-benchmarks-and-profiling-in-golang-gophercon-india-2016-talk-874c8b4dc3c5#.wlwu2sxi8

Run profiling command:

```
go test -run=^$ -bench=. -cpuprofile=cpu.out
go tool pprof tyrande.test.exe cpu.out
web
```

```
go tool pprof tyrande.exe cpu.out
web
```

See top 20 functions: 

```
top20
```
