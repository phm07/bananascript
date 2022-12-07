# BananaScript

BananaScript is a strongly typed, interpreted programming language written in Go.
It was inspired by the Monkey language and its [original interpreter](https://interpreterbook.com/).

## Builds
Pre-built binaries are available [here](https://builds.pauhull.de).

## Language tour

### Hello world
```
println("Hello, world!");
```

### Variables
```
let myString := "Hello, world!";
let myInt: int = 42;
let optionalInt: int? = 0;

myString = "Hi!"; // all variables are mutable
myInt = null; // illegal (null safety)
optionalInt = null; // legal
```

### Functions
```
fn add(a: int, b: int) int {
    return a + b;
}
let ten := add(5, 5);
```

### Loops
```
let i := 0;
while i++ < 5 {
    let j := 0;
    let line := "";
    while j++ < i {
        line = line + "* ";
    }
    println(line);
}
```

### Type extensions
```
fn (int)::fac() int {
    if this <= 1 {
        return 1;
    } else {
        return this * (this - 1).fac();
    }
}

let num := 5.fac(); // 120
```

## Builtin functions
```
fn println(string) void;
fn (string)::toString() string;
fn (int)::toString() string;
fn (bool)::toString() string;
```