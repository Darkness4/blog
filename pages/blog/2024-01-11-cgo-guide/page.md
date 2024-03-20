---
title: Using C libraries in Go with CGO
description: Simple guide and recommendations about CGO. For documentation purposes.
tags: [go, cgo, ffi, c]
---

## Table of contents

<div class="toc">

\\{\\{ $.TOC }}

</div>

## Introduction

I've made an [article about using Go with Portage and Crossdev for easy static compilation of CGO_ENABLED software](/blog/2023-11-08-go-with-portage-and-crossdev), but I didn't make an article about how to use CGO properly.

You see... I've seen comments in my projects that recommend using third-party libraries that implement full CGO bindings, without knowing the implications of doing so. While CGO is a powerful feature, there are major drawbacks to using it.

This article tries to explain what is FFI and how to use it properly.

## Quick reminder about C compilation

When compiling C, two major steps happens:

1. **Compilation**: Converts code into binary objects.
2. **Linking**: Combines all objects based on the `#include` directives and makes the final executable.

There are more steps involved into that, but these two steps are the most important one for understanding CGO.

Let's talk quickly about the first step. The compilation step happens in every compiled language (obviously), meaning Rust, CPP, C, Zig... If the language is compiled, there is a high chance that it can be exported into a library and called by another language.

Let's talk about the second step. The linking steps requires **two type of files:**

- **Compiled libraries and/or objects** containing symbols (functions, constants, ...).
- **Header files** declaring the existing symbols in the library and/or objects.

The header files are "interfaces" files which contains only symbols declaration:

```c
// library/library.h
#ifndef LIBRARY_H
#define LIBRARY_H

int sum(int a, int b);

#endif // LIBRARY_H
```

In C, when doing `#include "library.h"`, the line is simply replaced by the content of the header. The linker will resolve the references (`int sum(int a, int b)`) and link to the right library (`liblibrary.so` for example).

To avoid double includes, the `#ifndef LIBRARY_H`, `#define LIBRARY_H` and `#endif // LIBRARY_H` are used to define the symbols only once. These are called "include guards".

Example:

```c
#include "library.h"
#include "library.h"
```

Will only declare once:

```c
int sum(int a, int b);
```

When using FFI, these two steps are still necessary. Your library must be **compiled**, and a **header file** must be exposed. Thankfully, Go handle these two steps automatically. Go will use the header file to resolve references and will use the compiled libraries to link symbols declared in the header file to their implementation.

## Why use CGO instead of pure Go ?

While Go is simple and allows simple implementation, there are times when we need implementations that are _blazingly_ fast. For example, you wouldn't want to re-implement a video encoder in Go because there exists implementations that are faster in assembly.

The idea is to develop a simple layer with Go, and use complex implementations for the specific use cases. Let's take a real-world example: You want to develop a website that transcode a video with `libx264`. If we were to do it purely in C, there is a C API that we can use for `libx264`, but the C HTTP server is difficult to use and might be unsafe. If we were to do it purely in Go, we have a simple standard `net/http` server, but there are no `libx264` in pure Go.

The solution become evident: we need to combine Go and C.

If you've studied software architecture at all, you should know that this is easy to do, since Go acts as a "presentation" layer and C as a "business logic" layer. Go depends on C, but C does not depend on Go (Go calls C functions, but C doesn't need to call Go functions). What's more, since we only need the `int transcode(video_path)` function, the interface between these two layers is very thin, making it easy to convert.

So, let's talk about FFI, which is the practice used to convert functions in different languages.

## Foreign Function Interface (FFI)

A foreign function interface allows a programming language to use functions written in another language. Most of the time, that interface is represented by a C header `.h` file, which allows referencing functions defined inside a library file.

In the context of Go, CGO allows Go to call C functions and use C libraries, and vice-versa. However, due to the nature of the memory management in C, Go has some constraints:

1. **Data types in Go or C involves marshaling**. Not all types in Go can be converted in C.
2. **Memory allocated in C must be managed manually**. Any `malloc` must be `free`d, even if Go includes a garbage collector.
3. **Special "rules" in Go** are needed to configure the C compiler.

These points are already well documented in the [Go documentation](https://pkg.go.dev/cmd/cgo). To avoid repetition and obsolete recommendations, please read the documentation.

But, to summarize, if you have a `library.c`, `library.h` and `library.go`, with `library.go` converting C functions into pure Go functions, it would look like that:

```c
// library/library.c
int sum(int a, int b) { return a + b; }

```

```c
// library/library.h
#ifndef LIBRARY_H
#define LIBRARY_H

int sum(int a, int b);

#endif // LIBRARY_H

```

```go
// library/library.go
package library

// #include "library.h"
import "C"

func Sum(a int, b int) int {
	return int(C.sum(C.int(a), C.int(b)))
}

```

_You'll see that auto-completion doesn't work if the C header is incorrect. You need to click on `regenerate cgo definitions` each time you edit C code._

Data types are converted via `C.int` and `int`. There is no memory management here, but we would use `C.free` from the `stdlib` library if there is any `malloc`. Special rules are indicated above the `import "C"`. Comments declared above `import "C"` have actual effects.

Go has its own C compiler, so it compiles the C code and includes it in the final binary automatically.

To includes external system libraries, we recommend to use `pkg-config` as it is used by many libraries on Linux. `pkg-config` is a utility used to generate `CFLAGS` and `LDFLAGS` easily. For example, `/usr/lib64/pkgconfig/libcrypto.pc` contains:

```shell
prefix=/usr
exec_prefix=${prefix}
libdir=${exec_prefix}/lib64
includedir=${prefix}/include
enginesdir=${libdir}/engines-3
modulesdir=${libdir}/ossl-modules

Name: OpenSSL-libcrypto
Description: OpenSSL cryptography library
Version: 3.1.4
Libs: -L${libdir} -lcrypto
Libs.private: -ldl -pthread
Cflags: -I${includedir}
```

`Libs` are `LDFLAGS` and `Cflags` are `CFLAGS`. Normally, you would use `CFLAGS=-I<path/to/files.h>`, where the path to header files is `/usr/include` on Linux distributions. And, you would use `LDFLAGS=-L<path/to/libraries.so> -l<name of library>`, where the path to library files is `/usr/lib64` on Linux distributions. The name of the library is based on the file name. For example, `libcrypto.so` can be used by adding the `-lcrypto` flag (so, without the `lib` prefix and `.so` suffix).

`pkg-config` can be easily used in CGO. Example with OpenSSL:

```go
package crypto

// #cgo pkg-config: libcrypto
// #include <openssl/crypto.h>
import "C"

func Version() string {
    // OpenSSL_version gives a static null-terminated string. It's safe to call C.GoString and not freeing the pointer.
	return C.GoString(C.OpenSSL_version(C.OPENSSL_VERSION))
}

```

_Be sure to install the OpenSSL development headers and `pkg-config`._

The name of the library used in `pkg-config` depends on the file name stored inside the `/usr/lib64/pkgconfig/` directory. For example, `libcrypto.pc` can be used by calling `libcrypto`, without `.pc` suffix.

As a best practice, it is recommended to write wrapper functions (functions that just convert from Go to C), instead of using CGO everywhere. A wrapper function is also called a **binding**.

## About Go libraries offering CGO bindings

Now that you're familiar with FFI, you should know that some libraries offer complete CGO bindings. However, there's one big problem: unused bindings add unnecessary coupling between Go and C.

This is harmful because the coupling will **force the users and developers to use a specific version of a library instead of a version range.** For example, [go-astiav](https://github.com/asticode/go-astiav) is a library which provides Go functions to interact with the FFmpeg library. Because of the bindings, the unused functions force:

- The developers to use FFmpeg `5.1.2`. This makes the program unmaintainable in case FFmpeg is upgraded.
- The users to use FFmpeg `5.1.2` (if dynamically compiled). This makes the resulting program non-GPL compliant (impossible to change the distribution of FFmpeg).

Basically, _go-astiav_ is unusable due to its extreme coupling between Go and the library.

_"But wait! Why some of the most popular library ([go-sqlite3](https://github.com/mattn/go-sqlite3) for example) are distributed if it's dangerous to use?"_

The reason is quite simple: `libsqlite3` exposes a stable API (stable header files), which means it can be used over a wider version range despite the coupling. This allows users and developers to use any version of `libsqlite3`.

However, in general, when possible, it's better to **write the necessary C functions** and **minimize the interface between Go and C**.

_"Wait... I have to write C code?"_

Sadly, it's not only the most efficient way, but the safest way. Remember that you are using FFI because you want to connect the C business logic to the Go layer, not write C code in Go (even if it might be tempting).

For example, you want to use the math library to compute the roots of a quadratic operation:

$$
\displaylines{
ax^2 + bx + c = 0 \\
x = \frac{{-b \pm \sqrt{{b^2 - 4ac}}}}{{2a}}
}
$$

**Using the CGO-ed `libm`**:

```go
// quadratic.go
package quadratic

// #cgo LDFLAGS: -lm
// #include <complex.h>
// #include <math.h>
import "C"

func FindRoots2(a float64, b float64, c float64) (r1 complex128, r2 complex128) {
	discriminant := b*b - 4*a*c

	r1 = (-complex(b, 0) + complex128(C.csqrt(C.complex(complex(discriminant, 0))))) / complex(2*a, 0)
	r2 = (-complex(b, 0) - complex128(C.csqrt(C.complex(complex(discriminant, 0))))) / complex(2*a, 0)
	return r1, r2
}

```

**Using C, then CGO**:

```c
// quadratic.h
#ifndef QUADRATIC_H
#define QUADRATIC_H

#include <complex.h>

struct Roots {
  double complex root1;
  double complex root2;
};

struct Roots findRoots(double a, double b, double c);

#endif
```

```go
// quadratic.c
#include "quadratic.h"

struct Roots findRoots(double a, double b, double c) {
  struct Roots roots;

  double complex discriminant = b * b - 4 * a * c;

  roots.root1 = (-b + csqrt(discriminant)) / (2 * a);
  roots.root2 = (-b - csqrt(discriminant)) / (2 * a);

  return roots;
}

```

```go
// quadratic.go
package quadratic

// #cgo LDFLAGS: -lm
// #include "quadratic.h"
import "C"

func FindRoots(a float64, b float64, c float64) (complex128, complex128) {
	r := C.findRoots(C.double(a), C.double(b), C.double(c))
	return complex128(r.root1), complex128(r.root2)
}
```

You can see that when using CGO directly with `libm`, you are fighting against CGO's type-system. Not only that, but the IDE doesn't give any clues when there is a compilation bug, or doesn't even auto-complete. It's hard to debug, hard to write, and you are fighting against CGO. There are also issues with certain macros where it is simply not usable.

Using the second method, you can write C code with autocompletion thanks to `clangd` and have all the features of C, including macros. The wrapper function is small and easy to maintain. **The interface between Go and C is small**, which decouples Go and C.

## Conclusion

To summarize this article, to be able to maintain a Go program using C libraries, you must:

- Avoid using libraries that give all CGO wrappers, as unused bindings will cause unnecessary coupling between the program and the C library version.
- Write wrappers manually to limit the interface between Go and C, making the program easier to maintain.
- Separate CGO-related code from the rest of the Go code.
- Use `pkg-config` whenever possible.
- And the basics of C:
  - Don't forget error handling in C (return code).
  - Don't forget memory management.
  - Don't forget thread safety (mutex/semaphore).

The smaller the interface between Go and C, the more compatible it is with C libraries.

## References

- [CGo documentation](https://pkg.go.dev/cmd/cgo)
