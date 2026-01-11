---
title: Distroless containers, or the art to hide vulnerabilities.
description: A small articles about why distroless containers can be beneficial, but hides vulnerabilities.
tags: [devops, linux, container, distroless, security]
---

## Introduction

Hello, and happy new year! This article is mostly my personal opinion on distroless containers. If you are reading this you might learn something, or not. Since this is just my opinion, it is better to look for additional resource to learn more about the subject. I'll add some references, but remember this article is biased, since it's... well... my personal opinion.

Today, I wanted to talk about distroless containers and static compilation. As a DevOps engineer, these concepts are the best way to ship a container image with the least amount of vulnerabilities, or at least the least amount of _detectable_ vulnerabilities.

## Distroless containers

First, let's talk about distroless containers. For a program to run, you only need a Linux kernel and a C library. Since a container runtime uses the host kernel, there is no need to embed the kernel in the container image. So you only need to embed the C library, which is often glibc or muslc.

To summarize, a distroless container can be defined by **a container with nothing but the application dependencies (like glibc, OpenSSL...) and the application itself.**

To build a distroless container, container builders often use the Linux distribution's OS builder to create a base container image. For example, with `dnf`, you can build an OS image with just glibc:

```dockerfile
# We use Dockerfile to build the image, but you can use anything you want and tar the result.
FROM fedora:latest AS builder
RUN mkdir -p /base \
  && dnf install -y --use-host-config --installroot /base \
  glibc \
  --nodocs \
  --setopt=install_weak_deps=False \
  && dnf --installroot /base clean all


FROM scratch
COPY --from=builder /base /
```

```shell
docker build -t distroless-example .
```

After building it, since Fedora ships bash with glibc, we can run it:

```shell
docker run --rm -it --entrypoint /bin/bash distroless-example bash
```

Since Fedora ships Bash in the containers, it's actually not a good container builder, because Bash can be considered as _bloat_, which increase the container attack surface! And, since we only need the C library to run a program, we could statically link to the program we want to ship.

Statically linking means that the relevant part of the C library will be "embedded" inside the final library, which means the final program can be run on a container without any distribution-specific programs, i.e, a pure distroless container. Here's a simple example:

```c
// main.c
#include <stdio.h>

int main() {
  printf("Hello, world!\n");
  return 0;
}
```

```dockerfile
FROM alpine:latest AS builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /work
COPY . /work
RUN gcc -static -o /work/hello main.c


FROM scratch
COPY --from=builder /work/hello /
ENTRYPOINT [ "/hello" ]
```

```shell
docker build -t distroless-example .
```

And run it:

```shell
docker run --rm -it distroless-example
# Hello, world!
```

And because I didn't include anything besides `hello` in the final base image, I cannot run `bash` or any shell in the container. In fact, the container only contains `hello`. If we try to run `bash` in the container, we will get an error:

```shell
docker run --rm -it --entrypoint /bin/bash distroless-example
# Error: crun: executable file `/bin/bash` not found: No such file or directory: OCI runtime attempted to invoke a command that was not found
```

Our image only containers the program we want to run. It's not only the lightest way to ship a container, but also the most secure one... or is it?

## Considerations

### About dynamic linking, and why static linking hides vulnerabilities

There is a reason why every program on Linux uses a dynamic linking. Dynamic linking does NOT embed the library to the target and instead allows the target to invoke symbols (functions, variables, etc...) from the shared library. The shared library can be loaded and re-used by multiple program. For example, if we dynamically link `hello` and install `ldd` ([manual](https://man7.org/linux/man-pages/man1/ldd.1.html)) to find the linked library:

```dockerfile
FROM alpine:latest AS builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /work
COPY . /work
RUN gcc -o /work/hello main.c


FROM alpine:latest
RUN apk add --no-cache musl-utils
COPY --from=builder /work/hello /
ENTRYPOINT [ "ldd", "/hello" ]
```

```shell
docker build -t distroless-example .
```

and run it:

```shell
docker run --rm -it distroless-example
#        /lib/ld-musl-x86_64.so.1 (0x7fe68f6b6000)
#        libc.musl-x86_64.so.1 => /lib/ld-musl-x86_64.so.1 (0x7fe68f6b6000)
```

We can see that the program is dynamically linked to the muslc library (and it's been found at `/lib/ld-musl-x86_64.so.1`).

Because the library can be re-used, we can **track** runtime libraries. Vulnerability scanners like [Trivy](https://trivy.dev/) will be able to tell from the library version what vulnerabilities are shipped in the final container/OS.

So, **when you are using static linking, you are actually hiding vulnerabilities at runtime**. It can be useful, especially against difficult customers that complains that your container has a bad score against Trivy (a real scenario), but it is technically dangerous since we can't track the runtime dependencies of the program anymore. You'll need to track the dependencies at build-time, which is not an often used practice.

Also, when you are distributing statically linked programs, naive customers might not notice the issues, but experts will. Consider this scenario: If a core dependencies like OpenSSL or glibc has a critical vulnerability and gets statically linked to the final program, how can you make sure the final program is not affected? Can you tell if third-party programs and containers are safe? One way to tell is to have the container to be as transparent as possible.

This is why **signing** and **bill of materials** is important. We would be able to tell the origin of each dependencies, and also the origin of the final program. **This is key to avoid supply chain attacks.**

...But in reality, it is impossible to fully trust a third party since builders can also inject vulnerabilities through the compiler ([Kem Thompson Hack](https://wiki.c2.com/?TheKenThompsonHack)). It's mostly a question of how much trust you can give to the builder.

Lastly, static linking introduces weaknesses in the program, that are resolved in dynamically linked program (see [ASLR Protection for Statically Linked Executables](https://www.leviathansecurity.com/blog/aslr-protection-for-statically-linked-executables)).

To summarize, **static linking is about hiding the attack surface** and **make the container lighter**. Not about building trust, nor making the container more secure, so beware of statically compiled programs and distroless containers!

### No shell? That's not true.

Containers are Linux user namespace, and more precisely a _network/ipc/cgroup/mount/pid/uts/user_ namespace. The name "namespace" has its meaning: IDs are mapped. The user ID is mapped, the PID is mapped, the cgroup is mapped, etc. And by mapping, I really mean "A -> B", like "1 -> 10001". It's the reason why containers are not fully isolated from the host OS and is lighter than virtual machines.

Because this isolation is not perfect, it is actually possible to "walk" into a distroless container. Demonstration:

1. Edit the program to sleep and run it in the background:

   ```c
   // main.c
   #include <stdio.h>
   #include <unistd.h>

   int main() {
     printf("Sleeping indefinitely.\n");

     pause();

     printf("Exited.\n");
     return 0;
   }
   ```

   ```dockerfile
   FROM alpine:latest AS builder
   RUN apk add --no-cache gcc musl-dev
   WORKDIR /work
   COPY . /work
   RUN gcc -static -o /work/hello main.c


   FROM scratch
   COPY --from=builder /work/hello /
   ENTRYPOINT [ "/hello" ]
   ```

   ```shell
   docker build -t distroless-example .
   ```

   And run it:

   ```shell
   docker run --rm -it --name distroless-example -d distroless-example
   ```

2. Now, time to enter the container. You can move to the writable layer of the container:

   ```shell
   # Find the upper directory (writable layer) of the container
   DIR=$(docker inspect -f '{{.GraphDriver.Data.UpperDir}}' distroless-example)

   cd $DIR
   ```

3. At this point, you're on the top layer of the container. There is nothing in it because our program didn't write anything, but you can technically add stuff to it. For example, let's say you want to run `bash` inside the container. Download a static bash, and copy it:

   ```shell
   sudo curl -fsSL https://github.com/robxu9/bash-static/releases/download/5.2.015-1.2.3-2/bash-linux-x86_64 -o ./bash
   sudo chmod +x ./bash
   ```

   Now, you can run it:

   ```shell
   docker exec -it distroless-example /bash
   # bash-5.2#
   ```

   Your shell-less container is no more!

   (By the way, doesn't it seem weird to download a static bash from a random source? Are you able to tell if that bash isn't a malware? Can you trust the compiler that built it and the compiled code?)

Another scenario: You want to run a host program inside the container to debug it. To do that, you can use `nsenter` without the `--mount` flag. This will avoid using the same file-system as the container, and resolve the linked libraries. This can be useful when debugging the network:

```shell
# Move to the writable layer
DIR=$(docker inspect -f '{{.GraphDriver.Data.UpperDir}}' distroless-example)
cd $DIR

# Download curl to simulate the network traffic. You can also include this in the container via Dockerfile.
sudo curl -fsSL https://github.com/moparisthebest/static-curl/releases/download/v8.17.0/curl-amd64 -o ./curl
sudo chmod +x ./curl

# Get the PID of the container
PID=$(docker inspect -f '{{.State.Pid}}' distroless-example)

# Now enter the UTS/Network/PID/IPC/USER namespace of the container and run tcpdump.
# tcpdump is able to run because we are using the host filesystem, which contains all the linked libraries.
sudo nsenter --target $PID --uts --net --pid --user --ipc tcpdump -i any
```

In another window, you can simulate the traffic:

```shell
docker exec -it distroless-example /curl http://google.com
# Don't use HTTPS since we don't have the ca-certificates.crt inside the container.
# You could copy it from the host, though.
```

Your `tcpdump` should show the traffic going through the container network! Pretty cool, huh?

Basically, **distroless container can still be attacked** even if there is no shell. In fact, a lot of attacks doesn't require a shell and will find a way to [execute malicious code inside the container](https://www.cloudflare.com/learning/security/what-is-remote-code-execution/).

### Impure distroless images

Last consideration is about impure distroless images. If you are looking at [Google's distroless project](https://github.com/GoogleContainerTools/distroless), or [Chainguard](https://images.chainguard.dev/), you'll find minimal images. But, they actually contain stuff like glibc, or CA certificates!

Are you able to tell if these can be trusted?

Since you cannot simply shell into these images to audit them, you must rely on external verification methods:

- Look for the Bill of Materials. It is the list of packages used by the container. Chainguard is at least able to provide this (for example: [curl](https://images.chainguard.dev/directory/image/curl/sbom)).
- You can also check the signature of the image with [cosign](https://docs.sigstore.dev/cosign/verifying/verify/).

But, you can only trust as much as your guts! SBOMs and signatures, are only signs of trust, not proofs. *Who audits the auditors?*

## Conclusion

Let's be real about the benefits of a distroless image and statically compiled programs, it's about:

- Reducing the size of the final image and remove bloat.
- Reducing the attack surface of the container tok make life harder for the attacker.
- Lies to your customers and artificially lowering vulnerability scanners scores.

And that's it. It's not about:

- Building trust and providing transparency. A small container can still introduce the same vulnerabilities as a "distrofull" one.
- Avoiding dependencies vulnerabilities. They are still there, just hidden.

However, at most, you can always try to fix these issues by:

- Signing everything you build and sharing a public certificate, so users can verify the image origin.
- Providing a Bill of Materials to track the dependencies of the final program. This also includes statically linked libraries (just be transparent god dammit!).
- Scan the dependencies at build-time, not just at the end of the image build.

However, ultimately, it's the responsibility of the customer to secure a third-party program. They might complain about the security score, the lack of signing, the lack of rootless... without knowing what it really means. "Checkbox compliance" is not enough.Kem Thompson proved that it's impossible to trust third party program no matter how much transparency the builder provides.

Securing third-party software is not a passive task.
No matter how "distroless" an image is, true security relies on **robust sandboxing and strict runtime monitoring**. You cannot scan your way to safety. You must build an environment where the program is restricted by design and damage is limited.
