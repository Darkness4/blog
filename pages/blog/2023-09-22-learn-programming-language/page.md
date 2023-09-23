---
title: Learning your first programming language
description: About the programming language to learn first in 2023. It is a may be a filler post, but my opinion are here.
---

_Wow. It's one of those posts again!_

Yes, I'm writing this blog post to fill up my blog. Also, when searching "first programming language to learn", it's always a "top 5" with no example and just bullet points.

Most of them recommend Python, in which I certainly do **not** recommend.

There is a lot of criteria to consider when choosing a programming language, especially the first.

## Criteria

Before they even start, in reality, someone never wants to program first. Instead, they have an idea in mind, and want to make it possible.

In software programming, the first-time programmer often need to learn these concepts:

- Syntax
- Types
- Operators
- Control Flow
- Functions
- Error handling
- Debugging
- Algorithm
- Data structures
- IO
- Code style and Best practices

That's a lot, and most of the time, the novice programmer doesn't want to learn it all. However, learning as much as possible will certainly benefit the novice programmer's future work.

This is why, my criteria for choosing a programming language are the following:

1. **Ease of Learning**: Being able to learn most of the concepts in a minimum of time.
2. **Availability**: Not every programming language is accessible. Setting up a programming environment should not be difficult.
3. **Goal**: It is possible to program in every language. But some are certainly "domain"-specific.
4. **Trend**: Learning a dying programming language will simply demotivate the programmer.

I'm also taking account that the person won't be learning one programming language.

To be a good developer is to diversify since the future is uncertain. **Do not over-specialize, be a generalist.**

I've used all the languages mentioned in this article at some point in my life as a software engineer. I used develop software for Embedded Systems, Machine Learning, DevOps, Games, VR, Backend, Front-end (Web, Qt, Flutter, Android). It's always interesting to learn the syntax, rationale, behavior and paradigm of a programming language to understand how the programming ecosystem will evolve.

## Your first programming language for embedded systems: Arduino (C++)

Are we already starting with something hard? Not really, let me explain.

A good gift to someone to start programming is often an Arduino with an LCD shield.

While C++ is often recognized as a complex language, Arduino has abstracted most complex parts and propose a simple syntax:

```cpp
void setup() {
  // put your setup code here, to run once:
}

void loop() {
  // put your main code here, to run repeatedly:
}
```

When opening the Arduino IDE, you are already exposed to Arduino's syntax:

- `setup` is a function that only runs once
- `loop` is a function that run repeatedly

The IDE has a simple `compile` button and `upload` button to verify and deploy on an actual Arduino.

The IDE also includes examples, one being `DigitalReadSerial` which read the signal of a button and send the result via USB:

```cpp
// digital pin 2 has a pushbutton attached to it. Give it a name:
int pushButton = 2;

// the setup routine runs once when you press reset:
void setup() {
  // initialize serial communication at 9600 bits per second:
  Serial.begin(9600);
  // make the pushbutton's pin an input:
  pinMode(pushButton, INPUT);
}

// the loop routine runs over and over again forever:
void loop() {
  // read the input pin:
  int buttonState = digitalRead(pushButton);
  // print out the state of the button:
  Serial.println(buttonState);
  delay(1);  // delay in between reads for stability
}
```

The code is commented, and the syntax is easy to understand:

```cpp
// Decleare a variable
<variable-type> <variable-type> = <variable-value>;  // int pushButton = 2;

// Declare a function
<return-type> <function-identifier>(<function-parameters>) { // void setup() {
    <function-body>
}

// Use a function
<function-identifier>(<function-parameters>); // digitalRead(pushButton);
```

Many benefits to use Arduino as a first programming language:

- The syntax is **explicit**.
- The integrated development environment is **friendly**.
- The execution is **simple.**
- The community is **alive**.
- The language can be generalized in **software development** with C++.
- Pointers are not shown and Arduino allows global variables instead.

Overalls, Arduino is one great way to start programming!

## Your first programming language for software

No Arduino? Well, It's time to search through the many programming languages available.

I will be quick, these are languages that programmers often start to learn:

- Python
- C
- C++
- Rust
- Go
- JavaScript
- Typescript
- Java
- C#
- PHP
- Kotlin
- Objective-C
- Lua
- Ruby

That's a lot, and all of them can make software (to a certain level). As a senior developer, I would recommend a programming language in which you will be able to refactor and maintain in the future. **It's almost certain that the first project of a programmer is ugly, filled with bad practices.**

I will immediately remove languages that are too domain-specific and that are hard to make a Terminal User Interface or Graphical User Interface. After all, one concept to learn is input/output, meaning being able to print and read from the user inputs.

This removes **Java**, **C#**, **PHP**, **Kotlin**, **Objective-C**.

We can also remove language that a "dying": **Lua** and **Ruby**. Note that they aren't really dying, but shine only with certain conditions. **Lua** is extremely great to write plugins and mods (Garry's mod for example, or Roblox). **Ruby** shines in backend development, but **Python** and **Go** has proven to be more "friendly" to backend development.

I will now talk about the rest.

### Python, JavaScript and Typescript: Simple, but full of the bad practices, like the scripting languages they are.

To understand why I would not recommend Python, JavaScript and Typescript. Let me cite the Wikipedia definition of Python:

> **Python** is a [high-level](https://en.wikipedia.org/wiki/High-level_programming_language), [general-purpose programming language](https://en.wikipedia.org/wiki/General-purpose_programming_language). Its design philosophy emphasizes [code readability](https://en.wikipedia.org/wiki/Code_readability) with the use of [significant indentation](https://en.wikipedia.org/wiki/Off-side_rule).
>
> Python is [dynamically typed](https://en.wikipedia.org/wiki/Type_system#DYNAMIC) and [garbage-collected](<https://en.wikipedia.org/wiki/Garbage_collection_(computer_science)>). It supports multiple [programming paradigms](https://en.wikipedia.org/wiki/Programming_paradigm), including [structured](https://en.wikipedia.org/wiki/Structured_programming) (particularly [procedural](https://en.wikipedia.org/wiki/Procedural_programming)), [object-oriented](https://en.wikipedia.org/wiki/Object-oriented_programming) and [functional programming](https://en.wikipedia.org/wiki/Functional_programming). It is often described as a "batteries included" language due to its comprehensive [standard library](https://en.wikipedia.org/wiki/Standard_library).

All of these features abstracts a **lot** of programming concepts that are important as a software engineer.

- Let's start by **high-level**: this means a lot of computer and programming "details" are heavily abstracted like the compilation step and memory management. While I do also think that **memory management** (memory allocation and freeing) is not required, I do think that missing the **compilation step** is a grave mistake. When using an **interpreted** (vs compiled) programming languages, most errors happens at runtime. To avoid most errors, you need to install tools like the static analyzer **pylint** (which is not included in the REPL/interactive mode).

- **Dynamically typed**: This means that types are NOT strict. For example:

  ```python3
  # <variable-identifier> = <variable-value>
  a = 2
  a = "test"
  print(a) # Prints "test"
  ```

  This works. However, this behavior is problematic: **what is the type of "a" at a certain point in time**? Again, this causes more issues at runtime than a compiled language. Typescript also **has this issue**:

  ```typescript
  // let <variable-identifier> = <variable-value>
  let a = 2;
  a = 'test' as unknown as number;
  console.log(a); // Prints "test"
  ```

- **Garbage-collected**: This means that the language manages memory allocation for you. To be honest, this isn't really a problem, especially since we all have buffed computers.

- **Multiple programming paradigms**: This means there are multiple ways to achieve a solution. While you may say: "Good! That mean I'm free to do whatever I want, with the syntax that I prefer.", this also means: "I cannot read other people solution because I don't understand the syntax". This feature adds **complexity** and **implicitness** to the language.

- **Indents** vs **Brackets**: In reality, you don't care about indents (tab, spaces...) and brackets. You should let your IDE format your code, so that it is **standard**.

That's why I would never recommend Python, JavaScript or Typescript. These are **scripting** languages like **Bash** and **Perl**. Learning them can be easy at first, but like any scripting language, you're much more likely to shoot yourself in the foot when you scale up.

To avoid shooting yourself in the foot, you have to learn best practices like the PEP8 and PEP257 and have a multitude of tools to lint and format.

To summarize, Python, JavaScript and Typescript suffer from:

- **Implicitness**
- **Error-prone** and likely at **runtime**
- Syntax is simple, but the number of solutions adds in **complexity**
- Prone to **bad practices**

And I haven't talked about Typescript `let`, `var` `const`, `===`... and the packaging ecosystem which is also dreadful (`pnpm` or `yarn` or `bun`? `pip` or `conda`? `left-pad`???).

Overall, to be avoided at all costs. Especially, for learning programming concepts.

### C: simple to write, hard to bootstrap, compile and debug

Compared to what we saw, C is a low-level compiled programming language with a static type system. C is **not** designed for Object-Oriented Programming, but it is still possible to do OOP by using structures, header files and composition.

C syntax is quite explicit and simple:

```C
#include <stdio.h>

int main() {
    int count = 1; // Initialize the counter to 1

    while (count <= 10) { // Continue looping while count is less than or equal to 10
        printf("Count: %d\n", count); // Display the current count
        count++; // Increment the counter
    }

    return 0; // Exit the program
}
```

But one concept can be difficult to learn: Pointers and Memory Management:

```c
#include <stdio.h>
#include <stdlib.h>

int main() {
    // Declare a pointer to an integer
    int *ptr;

    // Allocate memory for an integer dynamically
    ptr = (int *)malloc(sizeof(int));

    if (ptr == NULL) {
        printf("Memory allocation failed\n");
        return 1;
    }

    // Assign a value to the dynamically allocated memory
    *ptr = 42;

    // Access and print the value using the pointer
    printf("Value: %d\n", *ptr);

    // Deallocate the memory to prevent memory leaks
    free(ptr);

    // After freeing, the pointer should be set to NULL
    ptr = NULL;

    return 0;
}
```

To summary, a pointer holds a memory address. `malloc(sizeof(int))` allocates memory for a `integer` (natural number) and returns the memory address. To change/get the value of the allocated memory, you must invoke `*ptr`. After doing `malloc`, you should always call `free` at the end.

Since C is too "simple" and low-level, managing the memory can be quite dreadful. There is no safety for the pointer: when reading the value of a pointer (also known as **dereferencing**), this may read garbage because the memory wasn't allocated, or the pointer is NULL. Reading a NULL pointer also crash the program with a **segmentation fault (core dump)**, with no log nor stack-trace. You need to use a **debugger** to be able to check where the program crashed. Also, fetching the `core dump` is not trivial.

There are other drawbacks too:

- Not friendly with Windows (but you should **switch to Linux** anyway)

- Dreadful linking step: While you can learn how to compile with C (compilation = convert code into binaries), linking third-party libraries is quite difficult. You need to learn how to use a build system like Make or CMake to be able to link safely, which **may adds a configure step**:

  ```makefile
  # Makefile for compiling the main.c program with math library

  # Compiler and flags
  CC = gcc
  CFLAGS = -Wall -Wextra -g
  LDFLAGS = -lm

  # Target executable
  TARGET = main

  all: $(TARGET)

  # Linking step
  $(TARGET): main.o
      $(CC) $(LDFLAGS) -o $(TARGET) main.o

  # Compile step
  main.o: main.c
      $(CC) $(CFLAGS) -c main.c

  clean:
      rm -f $(TARGET)
  ```

- Dreadful static compilation and export: The program is hard to share. You need the libraries versions to match between computers if compiled as a dynamically-linked executable (default behavior). If compiled as a static executable, you need all the C library to be statically compiled. However, **glibc is not statically compilable**, and you need to use **muslc**.

- Dreadful ecosystem due the last two points.

- Macros which are used for compile-time operations. This is often abused due to the lack of "checks" to the macros.

I would recommend learning **C** as the second programming language to understand how memory allocation works. Otherwise, there are better alternatives.

### C++: same as C, but with a better standard library and easier to use for OOP

Like the title said. It has the same drawbacks as C.

At least it easier to use a dynamically-sized list with C++:

```c
// Example with C
// structure definition
typedef struct {
    int *data;
    size_t size;
    size_t capacity;
} IntList;

void initIntList(IntList *list, size_t initialCapacity) {
    list->data = (int *)malloc(initialCapacity * sizeof(int));
    if (list->data == NULL) {
        fprintf(stderr, "Memory allocation failed\n");
        exit(1);
    }
    list->size = 0;
    list->capacity = initialCapacity;
}

void pushBackInt(IntList *list, int value) {
    if (list->size == list->capacity) {
        list->capacity *= 2;
        list->data = (int *)realloc(list->data, list->capacity * sizeof(int));
        if (list->data == NULL) {
            fprintf(stderr, "Memory allocation failed\n");
            exit(1);
        }
    }
    list->data[list->size++] = value;
}

void freeIntList(IntList *list) {
    free(list->data);
    list->data = NULL;
    list->size = list->capacity = 0;
}

int main() {
    IntList intList;
    initIntList(&intList, 2);

    pushBackInt(&intList, 42);
    pushBackInt(&intList, 56);
    pushBackInt(&intList, 78);

    for (size_t i = 0; i < intList.size; ++i) {
        printf("%d ", intList.data[i]);
    }
    printf("\n");

    freeIntList(&intList);

    return 0;
}
```

```c++
// Example with C++
#include <iostream>
#include <vector>

int main() {
    std::vector<int> intVector;
    // The return type here is a little complex. This means:
    // <namespace>::<object name><[template type]>
    // C++ is able to template classes and functions, meaning it is able to "generalize" a class or a function to other types.

    intVector.push_back(42);
    intVector.push_back(56);
    intVector.push_back(78);

    for (size_t i = 0; i < intVector.size(); ++i) {
        std::cout << intVector[i] << " "; // This is printing.
        // Basically, cout is a input stream/flow. By doing "<<", you adds the value to the stream.
        // Different, but actually quite clever.
    }
    std::cout << std::endl;

    return 0;
}
```

A lot of data structure is included in the standard library, making a great language for competitive programming.

However, C++ adds a few new drawbacks:

- New syntax which may be more complex. For example, there is too many ways to iterate:

  ```c++
  // 1. Using a for loop with vector
  for (int i = 0; i < 5; i++) {
      std::cout << vec[i] << " ";
  }
  std::cout << std::endl;

  // 2. Using a range-based for loop with vector
  for (int num : vec) {
      std::cout << num << " ";
  }
  std::cout << std::endl;

  // 3. Using an iterator with vector
  for (std::vector<int>::iterator it = vec.begin(); it != vec.end(); ++it) {
      std::cout << *it << " ";
  }

  // 4. Using for_each with a lambda from the algorithm library
  std::for_each(std::begin(arr), std::end(arr), [](int num) {
      std::cout << num << " ";
  }); // A lambda (anonymous first-class function) is used here.
  // The syntax is the following:
  // [](<function-parameters>) -> <return-type> {
  // The identifier is [], which means anonymous.
  std::cout << std::endl;
  ```

- Functions overloading, which adds complexity to the language by making multiple functions that hold the same name. Even the auto-completion is hard to read.

  ```c++
  #include <iostream>

  // Function overloading
  int add(int a, int b) {
      return a + b;
  }

  int add(int a, int b, int c) {
      return a + b + c;
  }

  int add(int a, int b, int c, int d) {
      return a + b + c + d;
  }

  int main() {
      // Calling functions with different parameter lists
      int sum1 = add(2, 3);
      int sum2 = add(2, 3, 4);
      int sum3 = add(2, 3, 4, 5);

      std::cout << "Sum 1: " << sum1 << std::endl;
      std::cout << "Sum 2: " << sum2 << std::endl;
      std::cout << "Sum 3: " << sum3 << std::endl;

      return 0;
  }
  ```

- Abusive templates in some cases, making it hard to debug or even read/write code.

Overalls, C++ is a no-go to learn as a first programming language.

### Rust: an alternative to C++, but a lot more "safe"

Rust is a compiled low-level programming language with type safety, null safety, memory safety, concurrency and multi-paradigm. It borrows concepts from functional programming like immutability, first-class citizen functions. Basically, it solves all the problems with C.

Some people say that Rust has a steep learning curve, I would argue that the steepness only come when using functional programming features and memory management.

Syntax-wise, it is quite simple:

```rust
// Function that adds two numbers and returns the result
fn add_numbers(a: i32, b: i32) -> i32 { // fn <function-identifier>(<function-parameters>) -> <return-type> {
    let sum = a + b; // Variable declaration and addition
    return sum
}

fn main() {
    // Variable declaration with explicit type annotation
    let x: i32 = 5;
    let y: i32 = 7;

    // Function call and assignment of the result to a variable
    let result = add_numbers(x, y);

    // Printing the result to the console
    println!("The sum of {} and {} is {}.", x, y, result);
}
```

Rust also has multiple ways to iterate a vector (dynamically-sized array):

```rust
// Method 1: Using a for loop
println!("Using a for loop:");
for num in &numbers {
    println!("{}", num);
}

// Method 2: Using a while loop with an iterator
println!("Using a while loop with an iterator:");
let mut iterator = numbers.iter();
while let Some(num) = iterator.next() {
    println!("{}", num);
}

// Method 3: Using `iter().enumerate()` to get both index and value
println!("Using `iter().enumerate()` to get both index and value:");
for (index, num) in numbers.iter().enumerate() {
    println!("Index: {}, Value: {}", index, num);
}

// Method 4: Using the `iter()` method and `for_each()`
println!("Using the `iter()` method and `for_each()`:");
numbers.iter().for_each(|&num| {
    println!("{}", num);
});
```

Rust is still complex due to the richness of the ecosystem. This could cause issues when reading the code of a person.

Rust solves these issues compared to C++:

- Memory safety by using an ownership system instead of malloc or smart pointers.
- Null-safety by using `Option.None` and `Option.Some(T)`
- Package management: There is [cargo](https://doc.rust-lang.org/cargo/), package manager that can be used to compile and pack your program, and manage dependencies. C++ has no official package manager.

The most difficult part is memory management, due to its borrower checker. Although the borrow checker enables memory safety, learning it before you know "why it's needed" can be quite "paradoxical" (maybe). It's probably a good idea to learn C and Rust at the same time.

Another criticism of Rust is the "implicit" nature of macros. This point is also debatable.

Overalls, Rust is pretty good to start with, until you hit the memory management wall combined with concurrency.

### Go: THE language that I recommend

Go is a compiled high-level programming language with static typing and garbage collection.

Immediately, one drawback is implicit memory management and hardware abstraction. However, this issue only comes if we WANT to manage memory and handle hardware at low-level.

Go's syntax is **simple and explicit** and actually shines because of it. Almost all program written in Go will have the same solution. This is to iterate a slice (dynamically-sized view of an array):

```go
package main

import "fmt"

func main() { // func <function-identifier>(<function-parameters) (return-types) {
	// Create a slice of integers
	numbers := []int{1, 2, 3, 4, 5}

	// Iterate over the slice using a for loop
	fmt.Println("Iterating over the slice:")
	for index, value := range numbers {
		fmt.Printf("Index: %d, Value: %d\n", index, value)
	}
}
```

While the language has no `;`, the indentation is also not important compared to Python. [Go determines a statement by reading ahead](https://www.youtube.com/watch?v=09-pyU0MLCY) and insert the semicolon when two statements are not compatible.

Since this is the language I recommend, let me list all the benefits of Go:

- **Explicit** and **simple** syntax.

- **Easy install and setting up** of Go in your favorite IDE

- **Easy package management**, directly included in Go via `go get`.

- **Easy formatting** with `goimports` or `gofmt`.

- **Easy concurrency** with goroutines and standard patterns like the Context API.

- **Standard concurrent stream with channels**:

  ```go
  package main

  import (
      "fmt"
  )

  func main() {
      // Create a channel
      ch := make(chan int)

      // Start a goroutine
      go func() {
          ch <- 42 // Send a value to the channel
      }()

      // Receive and print the value from the channel
      value := <-ch
      fmt.Println("Received:", value)
  }
  ```

- **Easy directory structure** via modules and packages:

  ```shell
  .
  ├── cmd # <your-module>/cmd
  │   └── cmd.go
  ├── main.go
  ├── go.mod # <your-module>
  └── go.sum
  ```

- **Easy and explicit error handling** via error as a value:

  ```go
  // Divide function returns a result and an error
  func Divide(a, b float64) (float64, error) {
      if b == 0 {
          return 0, errors.New("division by zero")
      }
      return a / b, nil
  }
  ```

- **Easy testing and benchmarking**:

  ```go
  // mymath_test.go
  package mymath_test

  import (
      "testing"
      "mymodule/mypackage/mymath"
  )

  func TestAdd(t *testing.T) {
      result := mymath.Add(2, 3)
      expected := 5

      if result != expected {
          t.Errorf("Add(2, 3) = %d; expected %d", result, expected)
      }
  }

  func BenchmarkAdd(b *testing.B) {
      // Benchmark the Add function with inputs 2 and 3
      for i := 0; i < b.N; i++ {
          _ = mymath.Add(2, 3)
      }
  }
  ```

- **Easy static and multi-platform compilation**:

  ```shell
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./executable ./main.go # static executable for linux
  CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./executable ./main.go # static executable for windows
  ```

- And more...

Go shines by being simple and explicit, making it easy to read by other people. It enables to follows best practices automatically and write idiomatic Go code. Therefore, it is a great language for open-source.

Even if Go is simple, it might be too simple. Go is missing some features that some other modern programming language has:

- No default parameters. But there is one way to achieve this, is to use the optional functions via variable-length parameters list:

  ```go
  // Options is a struct for specifying optional parameters
  type Options struct {
      Param1 string
      Param2 int
  }

  // DefaultOptions returns an Options struct with default values
  func DefaultOptions() Options {
      return Options{
          Param1: "default_param1",
          Param2: 42,
      }
  }

  func MyFunction(requiredParam string, options ...func(*Options)) {
      // Create an Options struct with default values
      opt := DefaultOptions()

      // Apply the provided options
      for _, option := range options {
          option(&opt)
      }

      // Use the options in the function
      fmt.Println("Required Parameter:", requiredParam)
      fmt.Println("Optional Parameter 1:", opt.Param1)
      fmt.Println("Optional Parameter 2:", opt.Param2)
  }
  ```

  Rust also has this issue.

- No compile-time solution besides code generation. No customizable build-system.

- Visibility by naming: Majuscule for public, minuscule for internal. This is actually debatable. In my opinion, it isn't that bad and you write code faster.

- No nil safety combined with Interfaces as pointers, making it "panic"-prone:

  ```go
  // MyInterface is an example interface
  // An interface is a "unimplemented" type used to define public behavior (e.g. public methods).
  type MyInterface interface {
      SayHello()
  }
  
  func main() {
      var myVar MyInterface // MyInterface is nil, not MyInterface{} because it is an unimplemented type
  
      // This will cause a runtime panic because myVar is nil
      myInterface.SayHello()
  }
  ```

And, that's it. Some may say that Go is too simple. I do believe that **explicitness** and **simplicity** are the most important criteria to select the first programming language to learn all the programming concepts.

Overall, Go is the language I recommend.

### Bonus: Zig: THE next language that I recommend

Zig is a compiled low-level programming language with static typing that can be used to replace C.

Zig has already [an article explaining why Zig is needed](https://ziglang.org/learn/why_zig_rust_d_cpp/). To summarize, it want to be even more explicit and simple:

- **No hidden control flow**

- **No hidden allocations**

- **Optional standard library**

- **A Package manager** (vs C/C++)

- **Easy static compilation** (vs C/C++/Rust)

- **Build system included** (vs C/C++/Rust/Go)

- **comptime** statements to declare compile time variable and functions:

  ```zig
  fn multiply(a: i64, b: i64) i64 {
      return a * b;
  }
  
  pub fn main() void {
      const len = comptime multiply(4, 5);
      const my_static_array: [len]u8 = undefined;
  }
  ```

Zig is still not production-ready, so we have to wait for concurrency. Zig is also lacking in linting/formatting due to the zig language server slow development.

Overall, if Zig is stable, you should learn it, because it may be the next C/C++/Go/Rust. It is already used in production with [Bun](https://github.com/oven-sh/bun).

## One last point: To OOP or not

Object-Oriented Programming or Functional Programming are programming paradigms that can be used to architecture your code.

Some talk about Clean Code, Clean Architecture, others talk about functions purity.

In reality, as I said at the beginning, you should never over-specialize. Always take these principles with a pinch of salt. You can try them, but you have to accept that there may be other ways.

A computer is naturally "procedural" so functional programming may not be adapter. Same for OOP:

- There more case where composition is better than inheritance
- Over-architecture is useless if your coworker don't know about it

## Conclusion

Based on explicitness and simplicity, here's my recommendation on which programming language you should learn:

1. Zig (when stable)
2. Go
3. Rust + C (together would be great)

With these programming languages, you can learn most of the programming concepts rapidly and naturally, without learning bad practices.

Also, this is my recommendations for domain-specific programming languages:

- Web front-end: TypeScript
- GUI: C++ with Qt, or Python with Qt. (Or Typescript with Electron, please don't). Dart with Flutter. Kotlin with JavaFX.
- Android: Kotlin or C++.
- Embedded Systems: C, C++, Zig, Rust or Arduino (if possible).
- High Performance Computing: C, C++, Zig, Rust.
- Games: C++ with Unreal, C# with Godot (F- Unity). For the choice of the game engine, I would recommend to start with any of them and let your passion go crazy.

I also recommend you to learn about:

- Auto-completion
- Auto-formatting
- Linting
- Snippets

Which should be included in your favorite IDE. These are the tools used to make programming more fun.

Although I've given you some recommendations on these languages, the final criterion is often the same: "Do you have enough resources (time, money, ...) to learn and use it?". **Like any tool, the real cost is not the tool itself, but the consequences that come from it.**

*If you have more or less time and are used to programming syntax, you can learn a large number of programming languages with [Learn X in Y minutes](https://learnxinyminutes.com).*

