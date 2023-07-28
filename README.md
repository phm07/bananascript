# BananaScript

BananaScript is a strongly typed, interpreted programming language written in Go.
It was inspired by the Monkey language and its [original interpreter](https://interpreterbook.com/).

## Playground

Check the language out in the [BananaScript Playground](https://bananascript.pauhull.de/).

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

### Type definitions
```
type myNewType := int;

let a: myNewType = 0;  // good
let b: myNewType = ""; // bad
```

### Interfaces
```
fn (int)::sayHello() {
    println("Hi!");
}

type myInterface := iface {
    sayHello: fn() void;
};

fn sayHello(x: myInterface) {
    x.sayHello();
}

123.sayHello();   // good
"123".sayHello(); // bad
```

## Builtins
```
type string := string;
type int := int;
type float := float;
type bool := bool;
type any := iface { };

fn println(any) void;  // Print line to console
fn print(any) void;    // Print to console (no \n)
fn prompt(any) string; // Input prompt
fn min(int, int) int;  // Returns smaller int
fn max(int, int) int;  // Returns bigger int

fn (any)::toString() string; // Returns object's string representation

fn (string)::uppercase() string; // Transforms string to uppercase
fn (string)::lowercase() string; // Transform string to lowercase
fn (string)::length() int;       // Returns string length
fn (string)::parseInt() int;     // Parses int from string

fn (int)::abs() int; // Returns absolute value

fn (float)::abs() float;   // Returns absolute value
fn (float)::ceil() float;  // Rounds value up
fn (float)::floor() float; // Rounds value down
fn (float)::round() float; // Rounds value
```
