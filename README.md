Hypothetically, if someone made a video game that outlined enemies in a specific color, you could write a bot that scraped the screen in real time instead of looking at memory. I wonder if any games do that?

This is an academic proof of concept, it is not meant for real use by anyone. Matching the performance of a memory scraping bot with a CV one is one hell of an uphill battle. The purpose of this project is to see how well I could do, not to create a tool for cheaters.

At last update the bot can handle reading the screen with high precision. It tracks the nearest centerpoint for an enemy when the "aim" button is pressed, but has some problems with aiming displacement. 

All that said, you're not going to find specific instructions on setting this project up here. I'm sticking to my "I didnt mean harm" argument above, but I've also seen a handful of cease-and-desist letters sent on similar projects. 

## Running

```
go build && ./tyrande.exe
```

## OpenCV 

Download opencv and mingw, following installation instructions on github.


- Use 64 bit vc11 libs and dlls.
- Replace the version string in the include statements.
- Remove the .dll in the include statements
- Put the files in: C:\mingw\mingw64\x86_64-w64-mingw32\lib
- Comment out method bodies in cvaux.go


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
