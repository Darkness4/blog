---
title: Go with Portage and Crossdev, for easy static multi-platform compilation of CGO_ENABLED software.
description: Want to statically compile for multi-platform in Go super-easily? Let me introduce Portage, Gentoo's package manager, and Crossdev, Gentoo's solution for cross-compilation.
---

## Table of contents

<div class="toc">

\\{\\{ $.TOC }}

</div>

## Introduction

This article describes how to easily statically compile for multiple platforms. Unfortunately, it does not cover OS X, mainly because there is nothing similar to Portage on OS X. You can use OSXCross, but you'll have to compile the dependency tree manually.

Static compilation has always been difficult in C due to the fact that you have to compile everything statically, or have the statically linked libraries available on your system. This was never a pleasure, and people started jumping on the flatpak/Appimage bandwagon (which is a fair take).

However, there is a certain "sexiness" to just distribute one static binary:

- It is fully contained.
- It doesn't require an OS to run it. You can make empty Docker image with just a binary in it.

This is why, in Go, it's possible to disable `CGO` by setting `CGO_ENABLED=0`, which allows Go to compile statically.

But... what if you need CGO, because of CGO packages like `sqlite` and C libraries like `ffmpeg`? Well, let me introduce you to Portage and Crossdev.

## Understanding static-linked and dynamically-linked.

Before explaining Portage and Crossdev, let me explain how C compilation and linking work.

The compilation stage is where the compiler analyzes the source code and converts it into assembly code. In the same step, the assembler included in the toolchain converts the assembly code into machine code, which is the binary file itself (often in ELF format).

```c
//main.c
int main() { return 0; }
```

```shell
# Compile
gcc -S -fverbose-asm main.c -o main.s
# Assemble
gcc -c main.s -o main.o
# `file` outputs "ELF 64-bit LSB relocatable"
```

An object file contains references to other functions and symbols defined in other object files or libraries. Furthermore, an object file contains no entry point, initialization or finalization code. It is therefore impossible to run it:

```shell
# Inspect symbols
readelf -s main.o

# Symbol table '.symtab' contains 4 entries:
#    Num:    Value          Size Type    Bind   Vis      Ndx Name
#      0: 0000000000000000     0 NOTYPE  LOCAL  DEFAULT  UND
#      1: 0000000000000000     0 FILE    LOCAL  DEFAULT  ABS main.c
#      2: 0000000000000000     0 SECTION LOCAL  DEFAULT    1 .text
#      3: 0000000000000000    11 FUNC    GLOBAL DEFAULT    1 main
```

To resolve the references, add the initialization and finalization code, we need to use a linker like `ld`. This is the linking step:

```shell
# Link
gcc main.o -o main
# Do "gcc main.c -o main" to do everything at once.
# `file` outputs "ELF 64-bit LSB pie executable, dynamically linked, interpreter /lib64/ld-linux-x86-64.so.2"
```

```shell
# Inspect symbols
readelf -s main

# Symbol table '.dynsym' contains 6 entries:
#    Num:    Value          Size Type    Bind   Vis      Ndx Name
#      0: 0000000000000000     0 NOTYPE  LOCAL  DEFAULT  UND
#      1: 0000000000000000     0 FUNC    GLOBAL DEFAULT  UND _[...]@GLIBC_2.34 (2)
#      2: 0000000000000000     0 NOTYPE  WEAK   DEFAULT  UND _ITM_deregisterT[...]
#      3: 0000000000000000     0 NOTYPE  WEAK   DEFAULT  UND __gmon_start__
#      4: 0000000000000000     0 NOTYPE  WEAK   DEFAULT  UND _ITM_registerTMC[...]
#      5: 0000000000000000     0 FUNC    WEAK   DEFAULT  UND [...]@GLIBC_2.2.5 (3)

# Symbol table '.symtab' contains 23 entries:
#    Num:    Value          Size Type    Bind   Vis      Ndx Name
#      0: 0000000000000000     0 NOTYPE  LOCAL  DEFAULT  UND
#      1: 0000000000000000     0 FILE    LOCAL  DEFAULT  ABS main.c
#      2: 0000000000000000     0 FILE    LOCAL  DEFAULT  ABS
#      3: 0000000000003e10     0 OBJECT  LOCAL  DEFAULT   20 _DYNAMIC
#      4: 0000000000002004     0 NOTYPE  LOCAL  DEFAULT   16 __GNU_EH_FRAME_HDR
#      5: 0000000000003fe8     0 OBJECT  LOCAL  DEFAULT   22 _GLOBAL_OFFSET_TABLE_
#      6: 0000000000000000     0 FUNC    GLOBAL DEFAULT  UND __libc_start_mai[...]
#      7: 0000000000000000     0 NOTYPE  WEAK   DEFAULT  UND _ITM_deregisterT[...]
#      8: 0000000000004000     0 NOTYPE  WEAK   DEFAULT   23 data_start
#      9: 0000000000004010     0 NOTYPE  GLOBAL DEFAULT   23 _edata
#     10: 0000000000001160     0 FUNC    GLOBAL HIDDEN    14 _fini
#     11: 0000000000004000     0 NOTYPE  GLOBAL DEFAULT   23 __data_start
#     12: 0000000000000000     0 NOTYPE  WEAK   DEFAULT  UND __gmon_start__
#     13: 0000000000004008     0 OBJECT  GLOBAL HIDDEN    23 __dso_handle
#     14: 0000000000002000     4 OBJECT  GLOBAL DEFAULT   15 _IO_stdin_used
#     15: 0000000000004018     0 NOTYPE  GLOBAL DEFAULT   24 _end
#     16: 0000000000001040    34 FUNC    GLOBAL DEFAULT   13 _start
#     17: 0000000000004010     0 NOTYPE  GLOBAL DEFAULT   24 __bss_start
#     18: 0000000000001155    11 FUNC    GLOBAL DEFAULT   13 main
#     19: 0000000000004010     0 OBJECT  GLOBAL HIDDEN    23 __TMC_END__
#     20: 0000000000000000     0 NOTYPE  WEAK   DEFAULT  UND _ITM_registerTMC[...]
#     21: 0000000000000000     0 FUNC    WEAK   DEFAULT  UND __cxa_finalize@G[...]
#     22: 0000000000001000     0 FUNC    GLOBAL HIDDEN    10 _init
```

You can see the `_init` and `_fini` codes, which are the initialization and finalization codes respectively. You can also see that our executable is dynamically linked to GLIBC 2.34. You can also check the links with `ldd`:

```shell
# Find links
ldd main
#         linux-vdso.so.1 (0x00007fff0237e000)
#         libc.so.6 => /lib64/libc.so.6 (0x00007fd35d245000)
#         /lib64/ld-linux-x86-64.so.2 (0x00007fd35d44d000)
```

Now let's talk about "dynamic" and "static" linking. What we've just done is dynamic linking, i.e. we've linked external libraries (GLIBC) to our executable. The external symbols in our code will use the definition stored in external libraries (`libc.so`). Our binary cannot function without these external libraries. This is why each OSes need a package manager. A package manager allows everyone to use external libraries and find the right version of each one. This is also why we need to compile for every Linux operating system when we want to distribute a dynamically linked library or executable.

However, when this proves too tedious, we can try to establish a "static" link with the executable. The "static" link consists of including the compiled external symbols directly in the final object. The final object will no longer contain references, but the complete definition of the symbols:

```shell
# Static linking (we do everything at once)
gcc -static main.c -o main
# `file` outputs "ELF 64-bit LSB executable, statically linked"
```

```shell
# Inspect symbols
readelf -s main

# Symbol table '.symtab' contains 2015 entries:
#   Num:    Value          Size Type    Bind   Vis      Ndx Name
#     0: 0000000000000000     0 NOTYPE  LOCAL  DEFAULT  UND
#     1: 0000000000000000     0 FILE    LOCAL  DEFAULT  ABS libc_fatal.o
#     2: 0000000000401100     5 FUNC    LOCAL  DEFAULT    6 __libc_message.cold
# ...
```

The object contains no references to external symbols and is therefore no longer linked to external libraries.

```shell
ldd main
#        not a dynamic executable
```

The resulting static object is also larger than the dynamic object.

However, the main reason we prefer to use a statically-linked executable is portability and ease of compilation. Statically compiled objects can run without the need for an operating system. This is an excellent choice when we want to containerize, as it allows us to build "distro-less" container images and drastically reduce the attack surface.

The linking stage is one of the steps often dreaded by programmers, which is why people tend to use pre-compiled runtimes like .NET, Go, Python and Java.

## About Go compilation and CGO

When compiling a Go executable, the Go runtime is included in the final executable. If a Go dependency uses CGO, the final executable will be linked to GLIBC and the external libraries. CGO allows us to include C code in Go using the Foreign Function Interface (FFI).

`github.com/mattn/go-sqlite3` is one of the most popular SQLite 3 driver for Go and uses CGO, making it hard to export a static binary.

For the sake of the example, let's use `math.h`, the math library of the C library.

```go
//main.go
package main

// We've add -lm to link against libm.so. No need to indicate the location of libm.so with -L since it is a system library.
// No need to indicate the location of headers file too with -I.
/*
#cgo LDFLAGS: -lm
#include <math.h>
*/
import "C"

import (
	"fmt"
)

func mySqrt(v float64) float64 {
	return float64(C.sqrt(C.double(v)))
}

func main() {
	fmt.Println(mySqrt(2))
}

```

By add `import "C"`, we are using CGO. `C.sqrt` is defined in the `libm.so` library. If we compile it:

```shell
go build main.go
# `file` outputs "ELF 64-bit LSB executable, dynamically linked, interpreter /lib64/ld-linux-x86-64.so.2"
```

And inspect it:

```shell
readelf -s main | less

# Symbol table '.dynsym' contains 51 entries:
#    Num:    Value          Size Type    Bind   Vis      Ndx Name
#      0: 0000000000000000     0 NOTYPE  LOCAL  DEFAULT  UND
#      1: 0000000000000000     0 FUNC    GLOBAL DEFAULT  UND free@GLIBC_2.2.5 (2)
#      2: 0000000000000000     0 FUNC    GLOBAL DEFAULT  UND [...]@GLIBC_2.3.4 (3)
#      3: 0000000000000000     0 FUNC    GLOBAL DEFAULT  UND _[...]@GLIBC_2.34 (4)
# ...
#
# Symbol table '.symtab' contains 2236 entries:
#    Num:    Value          Size Type    Bind   Vis      Ndx Name
#      0: 0000000000000000     0 NOTYPE  LOCAL  DEFAULT  UND
#      1: 0000000000000000     0 FILE    LOCAL  DEFAULT  ABS go.go
#      2: 0000000000402400     0 FUNC    LOCAL  DEFAULT   13 runtime.text
#      3: 0000000000402400    85 FUNC    LOCAL  DEFAULT   13 internal/abi.Kin[...]
#      4: 0000000000402460    22 FUNC    LOCAL  DEFAULT   13 internal/abi.(*T[...]
#      5: 0000000000402480   180 FUNC    LOCAL  DEFAULT   13 internal/abi.(*T[...]
#      6: 0000000000402540    68 FUNC    LOCAL  DEFAULT   13 internal/abi.(*T[...]
#      7: 00000000004025a0   123 FUNC    LOCAL  DEFAULT   13 internal/abi.Nam[...]
# ...

ldd main
#         linux-vdso.so.1 (0x00007fff685fe000)
#         libm.so.6 => /lib64/libm.so.6 (0x00007f76e20ed000)
#         libc.so.6 => /lib64/libc.so.6 (0x00007f76e1f13000)
#         /lib64/ld-linux-x86-64.so.2 (0x00007f76e21f1000)
```

You can see that the executable is linked with GLIBC and contains the Go runtime symbols.

So how do we compile statically with `CGO_ENABLED`?

With this:

```shell
go build -ldflags '-extldflags "-static"' main.go
```

However, when we compile statically, we need all dependencies to be static libraries. Here, we need `libm.a`, the static version of `libm.so` :

```shell
find /usr -name libm.a
# /usr/lib64/libm.a
```

Since it's a system library, it's easy. But what if you need a whole tree of dependencies that need to be statically compiled? Well, here's Portage.

## Portage: Gentoo's package manager

As you may know, Gentoo is an operating system renowned for the fact that all packages are built at installation time. Some people may say that it is not maintainable (which I may agree), but it is certainly the one of most stable OS of all time.

This is thanks to Portage, Gentoo's package manager. Although it's called a "package manager", it's closer to a build system like [Arch Build System](https://wiki.archlinux.org/title/Arch_build_system), than a package manager like [APT](https://ubuntu.com/server/docs/package-management). Since it's a build system, we can compile and install all the dependencies we want.

For example, let's use `libcurl`:

```go
//main.go
package main

// Use pkg-config instead of LDFLAGS to resolve all the LDFLAGS.
// You can try it on your terminal by running "pkg-config libcurl"
/*
#cgo pkg-config: libcurl
#include <string.h>
#include <curl/curl.h>
*/
import "C"

import (
	"fmt"
)

// strndup is used to safely convert a *C.char string into a Go string.
func strndup(cs *C.char, len int) string {
	return C.GoStringN(cs, C.int(C.strnlen(cs, C.size_t(len))))
}

func curlVersion() string {
	// curl_version returns a static buffer of 300 characters.
	return strndup(C.curl_version(), 300)
}

func main() {
	fmt.Println(curlVersion())
}

```

Just for the test, install `libcurl-dev` (whatever is the name of the cURL development package of your distribution) and test it:

```shell
go build main.go
./main
# libcurl/8.4.0 OpenSSL/3.1.4 zlib/1.2.13 zstd/1.5.5 c-ares/1.19.1 nghttp2/1.57.0
```

Now let's do it with Portage. To customize software, there's usually a configuration stage by using CMake, automake, autoconf... With Portage, we manipulate USE flags to customize a package. To compile static libraries, we need to use the `static-libs` USE flag on our dependencies.

Let's make a Dockerfile to automate the whole process. We use the Docker image `gentoo/stage3:musl` (we use musl libc instead of glibc to avoid any issues with static linking):

```dockerfile
#Dockerfile
FROM gentoo/stage3:musl

# TODO: install Go
# TODO: install dependencies with `static-libs`
# TODO: statically compile our go application
```

Let's install Go:

```dockerfile
FROM gentoo/stage3:musl

# Synchronize with gentoo repository
RUN emerge --sync

# Install Go dependencies
RUN MAKEOPTS="-j$(nproc)" USE="gold" emerge sys-devel/binutils dev-vcs/git

# Install Go (ACCEPT_KEYWORDS="~*" means to allows "unstable" version)
RUN MAKEOPTS="-j$(nproc)" \
  ACCEPT_KEYWORDS="~*" \
  emerge ">=dev-lang/go-1.21.0"

# TODO: install dependencies with `static-libs`
# TODO: statically compile our go application
```

Then install the dependencies. cURL is available as [`net-misc/curl`](https://packages.gentoo.org/packages/net-misc/curl). From the [packages.gentoo.org page](https://packages.gentoo.org/packages/net-misc/curl), you can see the available USE flags and [ebuild files](https://gitweb.gentoo.org/repo/gentoo.git/tree/net-misc/curl), which describes the configuration and compilation steps of the package.

```dockerfile
#...

RUN MAKEOPTS="-j$(nproc)" \
  ACCEPT_KEYWORDS="~*" \
  USE="static-libs" \
  emerge --newuse \
  net-misc/curl \
  sys-libs/zlib \
  net-libs/nghttp2 \
  net-dns/c-ares # We need to add net-libs/nghttp2, sys-libs/zlib and net-dns/c-ares, because the maintainer of net-misc/curl forgot to spread "static-libs" to these packages packages

# TODO: statically compile our go application
```

Then, let's copy our source code and statically compile our application:

```dockerfile
# ...

COPY main.go .

RUN go build -ldflags '-extldflags "-static"' -o main main.go
```

Let's build:

```shell
docker build -t builder .
```

You will see some `undefined references`. This is because since we are statically compiling, we need to statically link with the dependencies of `curl`. Let's edit `main.go`:

```go
//main.go
package main

// HERE
/*
#cgo pkg-config: libcurl libssl libcrypto libnghttp2 libcares zlib
#include <string.h>
#include <curl/curl.h>
*/
import "C"

// ...
```

Let's build again:

```shell
docker build -t builder .
```

Which should work, and let's extract the executable:

```shell
mkdir -p out
docker run --rm \
  -v $(pwd)/out:/out \
  builder \
  cp /main /out/main
```

Then, you can run it:

```shell
out/main
# libcurl/8.4.0 OpenSSL/3.1.4 zlib/1.2.13 zstd/1.5.5 c-ares/1.19.1 nghttp2/1.57.0
```

We've successfully compiled statically for amd64! Now let's build for multiple platforms using Crossdev!

## Crossdev: Gentoo's cross-compilation environment

Let's talk about cross-compiling. Cross-compiling means compiling with a different build platform than the target platform. For example, compiling for `linux/riscv64` on a `linux/amd64` machine is cross-compiling, just as compiling for `windows/amd64` on a `linux/amd64` machine is also cross-compiling.

Cross-compiling can be achieved using either:

- A cross-compiler with its cross-compiling environment
- A virtualized environment (as with Qemu).

Using Qemu is more stable and probably a safer way of cross-compiling, but it is the slower method. Using a cross-compiler is faster because we avoid the virtualization layer.

The aim of Crossdev is to easily install a cross-compiling environment on Gentoo. Its use is as follows for `linux/arm64` :

```dockerfile
#Dockerfile.arm64
# We need to set the platform as the build platform since we are cross-compiling.
FROM gentoo/stage3:musl

RUN emerge --sync
RUN MAKEOPTS="-j$(nproc)" USE="gold" emerge sys-devel/binutils dev-vcs/git
RUN MAKEOPTS="-j$(nproc)" \
  ACCEPT_KEYWORDS="~*" \
  emerge ">=dev-lang/go-1.21.0"

# Setup crossdev
RUN MAKEOPTS="-j$(nproc)" emerge sys-devel/crossdev
RUN mkdir -p /var/db/repos/crossdev/{profiles,metadata} \
  && echo 'crossdev' > /var/db/repos/crossdev/profiles/repo_name \
  && echo 'masters = gentoo' > /var/db/repos/crossdev/metadata/layout.conf \
  && chown -R portage:portage /var/db/repos/crossdev \
  && printf '[crossdev]\nlocation = /var/db/repos/crossdev\npriority = 10\nmasters = gentoo\nauto-sync = no' > /etc/portage/repos.conf \
  && crossdev --target aarch64-unknown-linux-musl

# TODO: install dependencies with `static-libs` for aarch64
```

We don't need to compile Go for other platforms, as it's already a cross-compiler. However, dependencies must be statically compiled in aarch64 :

```dockerfile
# ...

RUN MAKEOPTS="-j$(nproc)" \
  ACCEPT_KEYWORDS="~*" \
  CPU_FLAGS_ARM="v8 vfpv3 neon vfp" \
  USE="static-libs" \
  aarch64-unknown-linux-musl-emerge --newuse \
  net-misc/curl \
  dev-libs/openssl \
  net-libs/nghttp2 \
  sys-libs/zlib \
  net-dns/c-ares
```

Then, we set the target platform at the compilation step:

```dockerfile
# ...

COPY main.go .

RUN GOARCH=arm64 \
  CGO_ENABLED=1 \
  CC="aarch64-unknown-linux-musl-gcc" \
  CXX="aarch64-unknown-linux-musl-g++" \
  PKG_CONFIG="aarch64-unknown-linux-musl-pkg-config" \
  go build -ldflags '-extldflags "-static"' -o main main.go

```

Finally, we can run it:

```shell
docker build -t builder -f Dockerfile.arm64 .
```

Let's check if it compiled for arm64:

```shell
mkdir -p out
docker run --rm \
  -v $(pwd)/out:/out \
  builder \
  cp /main /out/main

file ./out/main
# ELF 64-bit LSB executable, ARM aarch64, statically linked
```

Pretty good, huh?

## Containerization and multi-arch manifests

Remember that a static binary can run on distroless? Let's build a multi-arch container manifest with our static binaries.

First step is to create a multi-arch cross-compilation environment. Earlier, we did a `Dockerfile.arm64` which only targets `arm64` from an `linux/amd64` platform. Let's make one Dockerfile that targets `linux/amd64`, `linux/arm64`, `linux/riscv64` and `windows/amd64` and that supports any linux build platform.

I'm going to use Podman to build my image, you can also use Docker Buildx:

```dockerfile
#Dockerfile.multi
FROM --platform=${BUILDPLATFORM} gentoo/stage3:musl AS builder

RUN emerge --sync
RUN MAKEOPTS="-j$(nproc)" USE="gold" emerge sys-devel/binutils dev-vcs/git
RUN MAKEOPTS="-j$(nproc)" \
  ACCEPT_KEYWORDS="~*" \
  emerge ">=dev-lang/go-1.21.0"

RUN MAKEOPTS="-j$(nproc)" emerge sys-devel/crossdev
RUN mkdir -p /var/db/repos/crossdev/{profiles,metadata} \
  && echo 'crossdev' > /var/db/repos/crossdev/profiles/repo_name \
  && echo 'masters = gentoo' > /var/db/repos/crossdev/metadata/layout.conf \
  && chown -R portage:portage /var/db/repos/crossdev \
  && printf '[crossdev]\nlocation = /var/db/repos/crossdev\npriority = 10\nmasters = gentoo\nauto-sync = no' > /etc/portage/repos.conf \
  && crossdev --target x86_64-unknown-linux-musl \
  && crossdev --target aarch64-unknown-linux-musl \
  && crossdev --target riscv64-unknown-linux-musl

RUN MAKEOPTS="-j$(nproc)" \
  ACCEPT_KEYWORDS="~*" \
  USE="static-libs" \
  x86_64-unknown-linux-musl-emerge --newuse \
  net-misc/curl \
  dev-libs/openssl \
  net-libs/nghttp2 \
  sys-libs/zlib \
  net-dns/c-ares

RUN MAKEOPTS="-j$(nproc)" \
  ACCEPT_KEYWORDS="~*" \
  CPU_FLAGS_ARM="v8 vfpv3 neon vfp" \
  USE="static-libs" \
  aarch64-unknown-linux-musl-emerge --newuse \
  net-misc/curl \
  dev-libs/openssl \
  net-libs/nghttp2 \
  sys-libs/zlib \
  net-dns/c-ares

RUN MAKEOPTS="-j$(nproc)" \
  ACCEPT_KEYWORDS="~*" \
  USE="static-libs" \
  riscv64-unknown-linux-musl-emerge --newuse \
  net-misc/curl \
  dev-libs/openssl \
  net-libs/nghttp2 \
  sys-libs/zlib \
  net-dns/c-ares

COPY main.go .

ARG TARGETOS TARGETARCH

RUN if [ "${TARGETARCH}" = "amd64" ]; then \
  export CC="x86_64-unknown-linux-musl-gcc"; \
  export CXX="x86_64-unknown-linux-musl-g++"; \
  export PKG_CONFIG="x86_64-unknown-linux-musl-pkg-config"; \
  elif [ "${TARGETARCH}" = "arm64" ]; then \
  export CC="aarch64-unknown-linux-musl-gcc"; \
  export CXX="aarch64-unknown-linux-musl-g++"; \
  export PKG_CONFIG="aarch64-unknown-linux-musl-pkg-config"; \
  elif [ "${TARGETARCH}" = "riscv64" ]; then \
  export CC="riscv64-unknown-linux-musl-gcc"; \
  export CXX="riscv64-unknown-linux-musl-g++"; \
  export PKG_CONFIG="riscv64-unknown-linux-musl-pkg-config"; \
  fi; \
  GOOS=${TARGETOS} \
  GOARCH=${TARGETARCH} \
  CGO_ENABLED=1 \
  go build -ldflags '-extldflags "-static"' -o main main.go

FROM scratch

COPY --from=builder /main .

ENTRYPOINT [ "/main" ]
```

`TARGETPLATFORM`, `TARGETARCH`, `TARGETVARIANT`, `TARGETOS`, `BUILDPLATFORM`, `BUILDARCH`, `BUILDOS` and `BUILDVARIANT` are automatically filled by Docker Buildx or Podman.

<center>

```d2
direction: right
builder\n\<build platform\> -> scratch\narm64 -> multi-arch\nmanifest
builder\n\<build platform\> -> scratch\namd64 -> multi-arch\nmanifest
builder\n\<build platform\> -> scratch\nriscv64 -> multi-arch\nmanifest
```

</center>

Let's build:

```shell
podman build \
  --manifest builder:multi \
  --platform linux/amd64,linux/arm64/v8,linux/riscv64 \
  -f Dockerfile.multi .
# OR
docker buildx \
  --platform linux/amd64,linux/arm64/v8,linux/riscv64 \
  -t builder:multi \
  -f Dockerfile.multi .
# You can also push the image by adding --push.
```

And that's it! The manifest contains multiple images for each architecture. By pushing the manifest, there will only be one endpoint for multiple architecture:

```shell
podman manifest push --rm --all builder:multi "docker://ghcr.io/darkness4/multi-arch-example:latest"
```

## Splitting the base image

Sometime, the `builder` image is too big, and you may not wish to build the base image again. We can split the Dockerfile into two:

```dockerfile
#Dockerfile.base
FROM gentoo/stage3:musl

RUN emerge --sync
RUN MAKEOPTS="-j$(nproc)" USE="gold" emerge sys-devel/binutils dev-vcs/git
RUN MAKEOPTS="-j$(nproc)" \
  ACCEPT_KEYWORDS="~*" \
  emerge ">=dev-lang/go-1.21.0"

RUN MAKEOPTS="-j$(nproc)" emerge sys-devel/crossdev
RUN mkdir -p /var/db/repos/crossdev/{profiles,metadata} \
  && echo 'crossdev' > /var/db/repos/crossdev/profiles/repo_name \
  && echo 'masters = gentoo' > /var/db/repos/crossdev/metadata/layout.conf \
  && chown -R portage:portage /var/db/repos/crossdev \
  && printf '[crossdev]\nlocation = /var/db/repos/crossdev\npriority = 10\nmasters = gentoo\nauto-sync = no' > /etc/portage/repos.conf \
  && crossdev --target x86_64-unknown-linux-musl \
  && crossdev --target aarch64-unknown-linux-musl \
  && crossdev --target riscv64-unknown-linux-musl

RUN MAKEOPTS="-j$(nproc)" \
  ACCEPT_KEYWORDS="~*" \
  USE="static-libs" \
  x86_64-unknown-linux-musl-emerge --newuse \
  net-misc/curl \
  dev-libs/openssl \
  net-libs/nghttp2 \
  sys-libs/zlib \
  net-dns/c-ares

RUN MAKEOPTS="-j$(nproc)" \
  ACCEPT_KEYWORDS="~*" \
  CPU_FLAGS_ARM="v8 vfpv3 neon vfp" \
  USE="static-libs" \
  aarch64-unknown-linux-musl-emerge --newuse \
  net-misc/curl \
  dev-libs/openssl \
  net-libs/nghttp2 \
  sys-libs/zlib \
  net-dns/c-ares

RUN MAKEOPTS="-j$(nproc)" \
  ACCEPT_KEYWORDS="~*" \
  USE="static-libs" \
  riscv64-unknown-linux-musl-emerge --newuse \
  net-misc/curl \
  dev-libs/openssl \
  net-libs/nghttp2 \
  sys-libs/zlib \
  net-dns/c-ares

```

Which you can build and push for the **build** platform:

```shell
docker build -t builder:base-amd64 -f Dockerfile.base
docker push builder:base-amd64
```

Then use it in the second Dockerfile:

```dockerfile
FROM --platform=${BUILDPLATFORM} builder:base-amd64 AS builder

COPY main.go .

ARG TARGETOS TARGETARCH

RUN if [ "${TARGETARCH}" = "amd64" ]; then \
  export CC="x86_64-unknown-linux-musl-gcc"; \
  export CXX="x86_64-unknown-linux-musl-g++"; \
  export PKG_CONFIG="x86_64-unknown-linux-musl-pkg-config"; \
  elif [ "${TARGETARCH}" = "arm64" ]; then \
  export CC="aarch64-unknown-linux-musl-gcc"; \
  export CXX="aarch64-unknown-linux-musl-g++"; \
  export PKG_CONFIG="aarch64-unknown-linux-musl-pkg-config"; \
  elif [ "${TARGETARCH}" = "riscv64" ]; then \
  export CC="riscv64-unknown-linux-musl-gcc"; \
  export CXX="riscv64-unknown-linux-musl-g++"; \
  export PKG_CONFIG="riscv64-unknown-linux-musl-pkg-config"; \
  fi; \
  GOOS=${TARGETOS} \
  GOARCH=${TARGETARCH} \
  CGO_ENABLED=1 \
  go build -ldflags '-extldflags "-static"' -o main main.go

FROM scratch

COPY --from=builder /main .

ENTRYPOINT [ "/main" ]
```

```shell
podman build \
  --manifest builder:multi \
  --platform linux/amd64,linux/arm64/v8,linux/riscv64 \
  -f Dockerfile.multi .
# OR
docker buildx \
  --platform linux/amd64,linux/arm64/v8,linux/riscv64 \
  -t builder:multi \
  -f Dockerfile.multi .
```

If you wish to have a multi-platform base image (which I do **not** recommend), you will have to use [`qemu-user-static`](https://github.com/multiarch/qemu-user-static).

## Conclusion

This is how I compiled [fc2-live-dl-go](https://github.com/Darkness4/fc2-live-dl-go). I created several base images with Portage and Crossdev, and used the base images to compile my project.

For OS X, I had to set up my own cross-compiling environment with OSXCross. There's no build system to help me, so I had to download and compile the dependencies manually.

The great thing about splitting the base image is that you can now run the Go compilation step in a CI:

```yaml
build-export-static:
  name: Build and export static Docker
  runs-on: ubuntu-latest

  steps:
    - uses: actions/checkout@v4

    - name: Set up QEMU
      uses: docker/setup-qemu-action@v3

    - name: Set up Docker Context for Buildx
      run: |
        docker context create builders

    - name: Set up Docker Buildx
      id: buildx
      uses: docker/setup-buildx-action@v3
      with:
        version: latest
        endpoint: builders

    - name: Login to GitHub Container Registry
      uses: docker/login-action@v3
      with:
        registry: ghcr.io
        username: ${{ github.actor }}
        password: ${{ secrets.GITHUB_TOKEN }}

    - name: Get the oci compatible version
      if: startsWith(github.ref, 'refs/tags')
      id: get_version
      run: |
        echo "VERSION=$(echo ${GITHUB_REF#refs/*/})" >> $GITHUB_OUTPUT
        echo "OCI_VERSION=$(echo ${GITHUB_REF#refs/*/} | sed 's/+/-/g' | sed -E 's/v(.*)/\1/g' )" >> $GITHUB_OUTPUT

    - name: Build and export dev
      uses: docker/build-push-action@v5
      with:
        file: Dockerfile.static
        platforms: linux/amd64,linux/arm64,linux/riscv64
        push: true
        tags: |
          ghcr.io/example/example:dev
        cache-from: type=gha
        cache-to: type=gha,mode=max

    - name: Build and export
      if: startsWith(github.ref, 'refs/tags')
      uses: docker/build-push-action@v5
      with:
        file: Dockerfile.static
        platforms: linux/amd64,linux/arm64,linux/riscv64
        push: true
        tags: |
          ghcr.io/example/example:latest
          ghcr.io/example/example:${{ steps.get_version.outputs.OCI_VERSION }}
        cache-from: type=gha
        cache-to: type=gha,mode=max
```
