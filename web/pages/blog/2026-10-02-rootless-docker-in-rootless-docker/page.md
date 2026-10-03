---
title: Setting up Rootless Docker in Rootless Docker
description: Learn how to set up rootless Docker in rootless Docker for secure, isolated CI/CD workloads. Explore why Docker-in-Docker is insecure and how Linux namespaces enable safe nested container execution.
tags:
  [
    devops,
    linux,
    infrastructure,
    docker,
    container,
    rootless,
    security,
    github-actions,
    ci,
  ]
---

## Table of contents

<div class="toc">

{{% $.TOC %}}

</div>

## Introduction

Everything start with this question: **How can I run multiple CI workers on the
same node with a high level of security and isolation?**

Naively, you use ephemeral CI workers inside a container: when the worker
finishes its job, the process ends and the container is restarted with a clean
state. Think about _GitHub Actions Runner Controller (ARC)_: ephemeral CI workers inside a
Kubernetes cluster, automatically scaled to the needs of the workload and
automatically restarted when the job ends.

When you start using GitHub ARC, you encounter now this issue: **How can I use
Docker buildx inside a container?**

Most people would suggest to either:

- Use **Docker in Docker**: The host docker daemon is rootful and launches
  ephemeral CI workers with a restart policy. Then, a second docker daemon is
  launched inside a container, on the host docker. The CI workers use the docker daemon socket of
  the containerized docker to build images.
- Use **Rootless Docker in Docker**: The host docker daemon is rootless (not
  started by root or by the init system) and launches
  ephemeral CI workers with a restart policy. Then, a second docker daemon is
  launched inside a container, on the host docker. The CI workers use the docker daemon socket of
  the containerized docker to build images.
- Use **an alternative to Docker to build images like Buildah**: The host docker
  daemon is rootless and the CI workers use Buildah to build images.

We'll comment on each of those options, then I'll show you how to set up
**Rootless Docker in Rootless Docker**, which is the best option when you want to
use buildx optimizations and buildx github actions.

## Docker in Docker: Easy to set up, but an inscure setup

Here's how to run Docker in Docker:

1. **Start a rootful Docker daemon**: By default, when installing docker, the
   docker daemon runs as root. The daemon is started by systemd and can be
   checked with `systemctl status docker.service` or `systemctl status
dockerd`.

   You might think that you are running rootless, but to be able to run Docker
   commands, you had to join the `docker` group.

   Joining the `docker` groups allows you to run Docker commands without
   `sudo`, but **the daemon is still running as root.**

2. **Start a Docker DinD container**:

   ```shell
   docker run -d \
      --privileged \
      --name docker-dind \
      docker:dind
   ```

   You can already see the vulnerability: it uses `--privileged`. This flag:
   - **Disables Seccomp**: Remove restrictions on system calls like `mount`,
     `pivot_root`, ... which are required to setup a user namespaces and
     container filesystems.
   - **Disables AppArmor/SELinux**: Also remove additional restrictions on
     system calls, linux capabilities and file access.
   - **Grants all Linux Capabilities**: It grants
     [capabilities](https://man7.org/linux/man-pages/man7/capabilities.7.html)
     like `SYS_ADMIN`, `NET_ADMIN`, ...,
     basically bypass permission checks and is in practically equivalent to
     running as a superuser.
   - **Exposes sensitive mounts**:
     - It mounts `/dev` from the host: every devices on the host are exposed,
       including drives, input devices, ... **which are unused by Docker, so no
       reason to expose them.**
     - It mounts `/sys/fs/cgroup` from the host: allows to create new cgroup
       which are used to isolate and track CPU, memory, network, ... of the
       nested container.
     - It remove the masks on `/proc` and `/sys` from the host: allows to read and write on `/proc`.

   So, in other word, when you are running CI jobs (or any multi-tenant
   workload), the nested containers are able to actually elevate their
   privileges, and possibly escape the container isolation.

In fact, there is no point to continue showing you how to run CI actions. I will
show you **how to escape the container isolation**:

```shell
# Enter the dind container to emulate a workload. THis can be the CI worker, or
# anything else.
docker exec -it docker-dind sh

# Mount the host drive
mount /dev/nvme0n1p1 /mnt

# Yay, you've successfully escaped the container isolation!
```

Of course, you may say that "the nested docker daemon will launch secure
containers, so it's not a problem". But in fact it is: the isolation you wanted
from docker in docker isn't actually achieved, and it is in fact a useless
method of isolation.

Here's a possible attack scenario if you were to run a GitHub Actions worker alongside
Docker in Docker. The Github Actions worker is not privileged but the DinD
container is.

Assuming this setup:

```shell
+-------------------------------------------------------------------+
|                        Rootful Docker Host                        |
|                                                                   |
|   +-----------------------+               +-------------------+   |
|   |  GH Action Runner     |               |    DinD Container |   |
|   |  Container            |               |  (Docker-in-Docker)   |
|   |                       |  Docker       |                   |   |
|   |   +---------------+   |  Socket       |   +-----------+   |   |
|   |   | Docker CLI    |---|---------------|---| Docker    |   |   |
|   |   +---------------+   |  /var/run/    |   | Daemon    |   |   |
|   |                       |  docker.sock  |   +-----------+   |   |
|   +-----------------------+               +-------------------+   |
|                                                                   |
+-------------------------------------------------------------------+

```

1. GitHub Actions executes `docker run --rm --privileged alpine /bin/sh -c "mount /dev/nvme0n1p1 /mnt && ls /mnt"`
2. Since the docker daemon is running inside a privileged container, and the
   GitHub Actions worker is executing a privileged command, the attacker will
   actually succeed to mount the host drive.

Entire PoC:

```shell
docker run -d \
  --privileged \
  --name docker-dind \
  -v shared_socket:/var/run \
  docker:dind

docker run --rm \
  -v shared_socket:/var/run \
  docker run --rm --privileged alpine /bin/sh -c "mount /dev/nvme0n1p1 /mnt && ls /mnt"
```

If we assume that `docker run --rm --privileged alpine /bin/sh -c "mount
/dev/nvme0n1p1 /mnt && ls /mnt"` was a CI job, the attacker has successfully
escaped the container isolation.

Therefore, **Docker in Docker is unsafe and unsuitable for isolation**

## Rootless Docker in Docker

### Understanding the differences between Rootful Docker and Rootless Docker

#### Rootful Docker

To create a container, Docker needs to:

- Create a network layer by:
  - Setting up a bridge network interface
  - Setting up a [virtual ethernet
    device](https://man7.org/linux/man-pages/man4/veth.4.html) (`veth`)
  - Wire the `veth` into the container namespace (`ip link set veth-guest netns mycontainer`)
  - Do the classic network interface configuration (IP configuration, routing, ...)
  - Setting up firewall, NAT and port forwarding rules with `iptables`
- Create a namespace with flags (like `pid`, `ipc`, `uts`, `net` ...) using the
  [`clone`](https://man7.org/linux/man-pages/man2/clone.2.html) system call.
  - The process will be running in a "namespace": an abstraction that allows to
    isolate certain aspects of Linux like:
    - Cgroup root directory
    - IPC, POSIX system queues
    - Network devices, stacks, ports and etc.
    - Mount points
    - PID (process IDs)
    - Time (boot and monotonic cloks) (Not used by Docker)
    - User and group IDs (Used by Rootless Docker. It'll be explained later.)
    - UTS (Hostname and NIS domain name)
  - It is possible to test this by using:

    ```shell
    # Docker rootful equivalent
    unshare --uts --ipc --net --pid --fork --mount-proc --cgroup --mount
    ```

- Setup cgroups to track CPU, memory, network, ...
  - It is possible to test this by setting up:

    ```shell
    # As root
    mkdir /sys/fs/cgroup/fake-docker-container

    # Limit memory to 512 MB in bytes = 536870912
    echo "536870912" | sudo tee /sys/fs/cgroup/fake-docker-container/memory.max

    # Limit cpu to 50000us quota per 100000us period
    echo "50000 100000" | sudo tee /sys/fs/cgroup/fake-docker-container/cpu.max

    # Launch your namespace
    # Don't specify --cgroup here because we are using a custom cgroup
    unshare --uts --ipc --net --pid --fork --mount-proc --mount /bin/sh -c "sleep infinity" &
    FAKE_CONTAINER_PID=$!

    # Attach the namespace to the cgroups
    echo "$FAKE_CONTAINER_PID" | sudo tee /sys/fs/cgroup/fake-docker-container/cgroup.procs
    # NB: Docker doesn't use two processes to set up cgroups. We're doing like
    # this because we can't move a namespace to a cgroup before it is created
    # due to `unshare`.
    ```

- Mount the container image file system with a read-only layer and a writable
  layer. This is done by using **OverlayFS**.
  - It is possible to test this by using:

    ```shell
    # As root
    # lower: read-only layer
    # upper: writable layer containing changes (additions, modifications, deletions)
    # work: working directory for overlayfs (not used by the user)
    # merged: the merged directory
    mkdir -p /tmp/overlay-demo/{lower,upper,work,merged}

    # Add some files to the lower (read-only) layer
    echo "Original content from Base Image" > /tmp/overlay-demo/lower/base_file.txt
    echo "File to be deleted" > /tmp/overlay-demo/lower/delete_me.txt

    # Mount the overlay
    mount -t overlay overlay \
      -o lowerdir=/tmp/overlay-demo/lower,upperdir=/tmp/overlay-demo/upper,workdir=/tmp/overlay-demo/work \
      /tmp/overlay-demo/merged

    # Check the content of the merged directory
    ls /tmp/overlay-demo/merged

    # Delete a file and check if it is present on the merged directory and the
    # lower directory
    rm /tmp/overlay-demo/merged/delete_me.txt
    stat /tmp/overlay-demo/merged/delete_me.txt
    stat /tmp/overlay-demo/lower/delete_me.txt

    # If you've setup the namespace from earlier, if you exit the namespace,
    # the mounted overlay will be unmounted.
    ```

- Drop privileges and linux capabilities
- Pivot the root filesystem and unmount the old root filesystem. Docker does
  this basically:

  ```shell
  pivor_root /tmp/overlay-demo/merged /tmp/overlay-demo/merged/old_root
  cd /
  # --lazy: detach the filesystem now, clean up things later
  umount --lazy /tmp/overlay-demo/merged/old_root
  rmdir /old_root
  ```

  (Why `pivot_root` over `chroot` is something you can google yourself. But I'm
  sure you guessed it: for security reasons.)

In all these steps, here's a comprehensive list of all syscalls used in Docker's
container setup, organized by step:

1. Create Namespaces

   ```shell
   unshare()           # Create new namespaces (PID, IPC, UTS, NET, Mount, Cgroup)
   fork()              # Create child process (when using --fork flag)
   ```

2. Setup Cgroups

   ```shell
   mkdir() / mkdirat()         # Create cgroup directories
   open() / openat()           # Open cgroup files for writing
   write()                     # Write limits to cgroup files
   close()                     # Close file descriptors
   statfs()                    # Check cgroup filesystem info
   ```

3. Mount OverlayFS

   ```shell
   mkdir() / mkdirat()         # Create lower, upper, work, merged directories
   mount()                     # Mount overlay filesystem
   open() / openat()           # Create files in lower layer
   write()                     # Write file contents
   ```

4. Manage Files in Overlay

   ```shell
   unlink() / unlinkat()       # Delete files
   stat() / lstat() / fstat()  # Check file status/metadata
   ```

5. Drop Privileges & Capabilities

   ```shell
   setuid() / setuid32()       # Change effective user ID
   setgid() / setgid32()       # Change effective group ID
   setgroups()                 # Set supplementary group IDs
   capset()                    # Drop Linux capabilities
   prctl()                     # Process control (PR_CAPBSET_DROP, etc.)
   seccomp()                   # Setup seccomp filter (optional)
   ```

6. Pivot Root & Cleanup

   ```shell
   pivot_root()        # Pivot to new root filesystem # Requires CAP_SYS_ADMIN + real root
   chdir() / fchdir()  # Change directory
   umount() / umount2() # Unmount old root (with MNT_LAZY)
   rmdir()             # Remove old_root directory
   ```

7. Execute Container Init

   ```shell
   execve()            # Execute container's init process
   ```

(I might have missed some syscalls. If you find any, please let me know.)

#### Rootless docker

##### Leveraging user namespaces to avoid root

When using rootless docker, you run docker as a non-root user (duh!), but this
means some commands will simply not work like `mount` or `pivot_root` which are
vital to the container setup process.

As you saw earlier, we can use `--user` in `unshare` to create a `user`
namespace. In this user namespace, we can map the current user ID to the root user
ID using `--map-root-user`. It is a "fake" root in some sense, but it cannot really call privileged
command like `mount`:

```shell
unshare --user --map-root-user /bin/sh -c "mount --make-private / && touch /tmp/resolv.conf && mount --bind /etc/resolv.conf /tmp/resolv.conf"
# mount: /: permission denied.
```

But in a mount namespace, we can do it:

```shell
unshare --mount --user --map-root-user /bin/sh -c "mount --make-private / && touch /tmp/resolv.conf && mount --bind /etc/resolv.conf /tmp/resolv.conf && cat /tmp/resolv.conf"
# ok
```

However, this is not enough to set up a rootless container: **we need a UID/GID
namespace in which we are allowed to use**. For that, we set up not only a user
namespace but also UID and GID mapping.

This is something that is done once by the root user to provide a specific user
"UID/GID" space. The root user must set up:

```shell
# As root
# Format is <user>:<start-uid>:<count>
echo "my-user:100000:65535" >> /etc/subuid
echo "my-user:100000:65535" >> /etc/subgid
```

The mapping will looks like:

```shell
# Host     Container
100000 --> 0
165535 --> 65535
```

**User namespace with UID/GID mapping allows us to create rootless containers.**

Using user namespaces, we'll be able to replace every rootful component with a
rootless equivalent.

##### Using slirp4netns to connect the external network to the namespace

To set up the network in a rootless environment, we cannot use `veth` and
`bridge` anymore since it requires the root privileges. Instead, we can leverage
TAP devices by using `slirp4netns`. Of course, you can't set up the network
without a user namespace, but basically, the setup looks like this:

```shell
unshare --user --map-root-user \
  --net \
  --mount \
  /bin/sh -c "sleep infinity" &
FAKE_CONTAINER_PID=$!

slirp4netns --configure --mtu=65520 --disable-host-loopback $FAKE_CONTAINER_PID tap0 &

# Enter back the user namespace
nsenter --net --user --preserve-credentials -t $FAKE_CONTAINER_PID

# Check the network status
ip a
# 1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
#     link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
#     inet 127.0.0.1/8 scope host lo
#        valid_lft forever preferred_lft forever
#     inet6 ::1/128 scope host proto kernel_lo
#        valid_lft forever preferred_lft forever
# 2: tap0: <BROADCAST,UP,LOWER_UP> mtu 65520 qdisc fq_codel state UNKNOWN group default qlen 1000
#     link/ether 42:2b:eb:14:1d:a8 brd ff:ff:ff:ff:ff:ff
#     inet 10.0.2.100/24 brd 10.0.2.255 scope global tap0
#        valid_lft forever preferred_lft forever
#     inet6 fe80::402b:ebff:fe14:1da8/64 scope link proto kernel_ll
#        valid_lft forever preferred_lft forever

# Test network access
ping 8.8.8.8 -c 4
# PING 8.8.8.8 (8.8.8.8) 56(84) bytes of data.
# 64 bytes from 8.8.8.8: icmp_seq=1 ttl=255 time=28.0 ms
# 64 bytes from 8.8.8.8: icmp_seq=2 ttl=255 time=18.7 ms
# 64 bytes from 8.8.8.8: icmp_seq=3 ttl=255 time=17.0 ms
# 64 bytes from 8.8.8.8: icmp_seq=4 ttl=255 time=21.9 ms
#
# --- 8.8.8.8 ping statistics ---
# 4 packets transmitted, 4 received, 0% packet loss, time 3001ms
# rtt min/avg/max/mdev = 17.005/21.398/28.034/4.208 ms
```

##### Setting up cgroups

!!!note NOTE

Since I'm running Gentoo Linux with OpenRC, my setup on user cgroups is different from
a SystemD setup. It'll be assumptions for this section.

!!!

SystemD mount user cgroups at `/sys/fs/cgroup/user.slice/user-<UID>.slice` where
the user can create their own cgroups.

You can test using:

```shell
mkdir -p /sys/fs/cgroup/user.slice/user-$(id -u).slice/my-container
echo 100M > /sys/fs/cgroup/user.slice/user-$(id -u).slice/my-container/memory.max
echo $FAKE_CONTAINER_PID | tee /sys/fs/cgroup/user.slice/user-$(id -u).slice/my-container/cgroup.procs
```

###### Setting up the container filesystem

OverlayFS cannot be used without root, so we need to use FUSE (Filesystem in
USErspace). To do that, we can use `fuse-overlayfs`:

```shell
# lower: read-only layer
# upper: writable layer containing changes (additions, modifications, deletions)
# work: working directory for overlayfs (not used by the user)
# merged: the merged directory
mkdir -p /tmp/overlay-demo/fuse/{lower,upper,work,merged}

# Add some files to the lower (read-only) layer
echo "Original content from Base Image" > /tmp/overlay-demo/fuse/lower/base_file.txt
echo "File to be deleted" > /tmp/overlay-demo/fuse/lower/delete_me.txt

# Mount the overlay
fuse-overlayfs \
  -o lowerdir=/tmp/overlay-demo/fuse/lower,upperdir=/tmp/overlay-demo/fuse/upper,workdir=/tmp/overlay-demo/fuse/work \
  /tmp/overlay-demo/fuse/merged

# Check the content of the merged directory
ls /tmp/overlay-demo/fuse/merged

# Delete a file and check if it is present on the merged directory and the
# lower directory
rm /tmp/overlay-demo/fuse/merged/delete_me.txt
stat /tmp/overlay-demo/fuse/merged/delete_me.txt
stat /tmp/overlay-demo/fuse/lower/delete_me.txt

# If you need to unmount: fusermount -u /tmp/overlay-demo/fuse/merged
```

At this point, we've got everything to make a container working. `pivot_root`
works in a user namespace, so we can also switch the root.

### Running rootless docker in docker

Here's how to run Rootless Docker in Docker:

```shell
docker run --rm -v shared_socket:/work alpine chown 1000:1000 /work

docker run -d --name docker-dind \
  --security-opt seccomp=unconfined \
  --security-opt apparmor=unconfined \
  --security-opt systempaths=unconfined \
  --device /dev/fuse \
  --device /dev/net/tun \
  -v shared_socket:/var/run/user/1000 \
  docker:dind-rootless

# Test it
docker run --rm \
  -v shared_socket:/var/run \
  docker run --rm alpine echo "Hello World"
```

To explain some of the flags:

```shell
  --security-opt seccomp=unconfined \
  --security-opt apparmor=unconfined \
```

is required to allow the namespace to run system calls like `mount` or
`pivot_root`. **It does not provide any additional privileges**, but remove
common protections for a container (it is not expected for a container to run
`mount` or `pivot_root`).

```shell
  --security-opt systempaths=unconfined \
```

similar to `seccomp=unconfined` and `apparmor=unconfined`, it allows the
namespace to mount `/sys`, which are required for `cgroups` to work. This is
also a common protection since there are no reasons for a container to interact
with `/sys`..., but this isn't the case for Rootless Docker in Docker.

```shell
  --device /dev/fuse \
```

As we said earlier, for `fuse-overlayfs` to work, we need to mount `/dev/fuse`.

```shell
  --device /dev/net/tun \
```

And for `slirp4netns` to work, we need to mount `/dev/net/tun`.

You can test the PoC from earlier and see that it doesn't work:

```shell
docker run --rm \
  -v shared_socket:/var/run \
  docker run --rm --privileged alpine /bin/sh -c "mount /dev/nvme0n1p1 /mnt && ls /mnt"
# docker: Error response from daemon: failed to create task for container: failed to create shim task: OCI runtime create failed: runc create failed: unable to start container process: error during container init: error mounting "sysfs" to rootfs at "/sys": mount src=sysfs, dst=/sys, dstFd=/proc/thread-self/fd/15, flags=MS_NOSUID|MS_NODEV|MS_NOEXEC: operation not permitted

# And if running without privileges:
docker run --rm \
  -v shared_socket:/var/run \
  docker run --rm alpine /bin/sh -c "mount /dev/nvme0n1p1 /mnt && ls /mnt"
mount: permission denied (are you root?)
```

We have now a proper protection layer for nested containers. But, what if the
attacker was able to escape the nested container and the DinD container? He
would be `root`!

So, time to run Rootless Docker in Rootless Docker.

## Rootless Docker in Docker

If you try to run Rootless Docker in Rootless Docker, you'll get the following error:

```shell
docker run --rm -v shared_socket:/work alpine chown 1000:1000 /work

docker run -d --name docker-dind \
  --security-opt seccomp=unconfined \
  --security-opt apparmor=unconfined \
  --security-opt systempaths=unconfined \
  --device /dev/fuse \
  --device /dev/net/tun \
  -v shared_socket:/var/run/user/1000 \
  docker:dind-rootless

docker logs docker-dind
# [rootlesskit:parent] error: failed to setup UID/GID map: newuidmap 75 [0 1000 1 1 100000 65536] failed: newuidmap: write to uid_map failed: Operation not permitted
```

**Why is that?** Remember, we are pre-allocated a certain range of UIDs and
GIDs, between `0` and `65535`.

The error tells that `rootlesskit` (the process responsible for setting up
rootless in Docker) is unable to setup the UID/GID mapping. Looking at the
parameters sent in that error:

```
newuidmap 75 0 1000 1 1 100000 65536
```

Breaking this down:

| Parameter         | Value    | Meaning                                  |
| ----------------- | -------- | ---------------------------------------- |
| **PID**           | `75`     | Process ID to set up UID mapping for     |
| **ID-inside-ns**  | `0`      | UID 0 inside the user namespace          |
| **ID-outside-ns** | `1000`   | Maps to UID 1000 outside the namespace   |
| **Length**        | `1`      | Length of this mapping (just 1 user)     |
| **ID-inside-ns**  | `1`      | UID 1 inside the namespace               |
| **ID-outside-ns** | `100000` | Maps to UID 100000 outside the namespace |
| **Length**        | `65536`  | Length of this mapping (65,536 users)    |

Or, in other word, `rootlesskit` is trying to map `0` (container-side) to `1000`
(host-side), and `1-65536` (container-side) to `100000-165536` (host-side).

However, we just said earlier that we are allocated between `0` and `65535`, so
`100000-165536` are actually out of range! So, we need to override `/etc/subuid`
and `/etc/subgid` inside the container (or allow more UIDs/GIDs on the host)!

```shell
# rootless is the user in docker:dind-rootless
echo "rootless:30000:20000" > /tmp/subuid
echo "rootless:30000:20000" > /tmp/subgid
# Mapping:
# 0 -> 30000
# 20000 -> 50001
# Above 20000, the container isn't allowed to chmod, chown, su, etc... because
# it is outside of the user namespace (we're hitting the host UIDs/GIDs)

docker run -d --name docker-dind \
  --security-opt seccomp=unconfined \
  --security-opt apparmor=unconfined \
  --security-opt systempaths=unconfined \
  --device /dev/fuse \
  --device /dev/net/tun \
  -v shared_socket:/var/run/user/1000 \
  -v /tmp/subuid:/etc/subuid \
  -v /tmp/subgid:/etc/subgid \
  docker:dind-rootless

# Test it
docker run --rm \
  -v shared_socket:/var/run \
  docker run --rm alpine echo "Hello World"
```

Aaaand it works! Pretty cool right?

## Setting up GitHub Actions runners with Rootless Docker

Assuming an ephemeral GitHub Actions Runner setup, your script should looks
like:

```shell
REPO="org/repo"
GH_RUNNER_NAME="rootless-docker-runner"
LABELS='["self-hosted","linux","x64"]'
ENCODED_JIT_CONFIG=$(curl -fsSL -X "POST" \
    -H "Authorization: Bearer ${GH_TOKEN}" \
    -H "Accept: application/vnd.github+json" \
    "https://api.github.com/repos/${REPO}/actions/runners/generate-jitconfig" -d '{
    "name": "'"$GH_RUNNER_NAME"'",
    "runner_group_id": 1,
    "labels": '"$LABELS"',
    "work_folder": "/home/runner/_work"
  }' | jq -r '.encoded_jit_config')

docker run --rm -v work:/work alpine chown 1001:1001 /work
docker run --rm -v shared_socket:/work alpine chown 1001:1001 /work

docker run --rm --name "runner-container" \
  -v "work:/home/runner/_work" \
  -v "shared_socket:/run/user/1001" \
  -e DOCKER_HOST=unix:///run/user/1001/docker.sock \
  -e TMPDIR=/home/runner/_work/tmp \
  "ghcr.io/actions/actions-runner:2.337.0" \
  bash -c "mkdir -p \$TMPDIR && chmod 755 \$TMPDIR && exec ./run.sh --jitconfig $ENCODED_JIT_CONFIG" &
RUNNER_PID=$!

wait "$RUNNER_PID"
```

To make volume mounts work, the volume `work` will be shared between the runner
and the rootless Docker-in-Docker container.

You can already tell there will be an issue with UID and GID: the runner
is running with UID and GID 1001, but the docker host is running with UID and
GID 1000. So, we need to change the UID and GID of the docker daemon to 1001.

To do this, we need to override `/etc/group` and `/etc/passwd` inside the
container, and set up proper permissions in the working directories.

```shell
cat <<EOF > /tmp/group
root:x:0:root
bin:x:1:root,bin,daemon
daemon:x:2:root,bin,daemon
sys:x:3:root,bin
adm:x:4:root,daemon
tty:x:5:
disk:x:6:root
lp:x:7:lp
kmem:x:9:
wheel:x:10:root
floppy:x:11:root
mail:x:12:mail
news:x:13:news
uucp:x:14:uucp
cron:x:16:cron
audio:x:18:
cdrom:x:19:
dialout:x:20:root
ftp:x:21:
sshd:x:22:
input:x:23:
tape:x:26:root
video:x:27:root
netdev:x:28:
kvm:x:34:kvm
games:x:35:
shadow:x:42:
www-data:x:82:
users:x:100:games
ntp:x:123:
abuild:x:300:
utmp:x:406:
ping:x:999:
nogroup:x:65533:
nobody:x:65534:
docker:x:2375:
dockremap:x:101:dockremap
rootless:x:1001:
EOF

cat <<EOF > /tmp/passwd
root:x:0:0:root:/root:/bin/sh
bin:x:1:1:bin:/bin:/sbin/nologin
daemon:x:2:2:daemon:/sbin:/sbin/nologin
lp:x:4:7:lp:/var/spool/lpd:/sbin/nologin
sync:x:5:0:sync:/sbin:/bin/sync
shutdown:x:6:0:shutdown:/sbin:/sbin/shutdown
halt:x:7:0:halt:/sbin:/sbin/halt
mail:x:8:12:mail:/var/mail:/sbin/nologin
news:x:9:13:news:/usr/lib/news:/sbin/nologin
uucp:x:10:14:uucp:/var/spool/uucppublic:/sbin/nologin
cron:x:16:16:cron:/var/spool/cron:/sbin/nologin
ftp:x:21:21::/var/lib/ftp:/sbin/nologin
sshd:x:22:22:sshd:/dev/null:/sbin/nologin
games:x:35:35:games:/usr/games:/sbin/nologin
ntp:x:123:123:NTP:/var/empty:/sbin/nologin
guest:x:405:100:guest:/dev/null:/sbin/nologin
nobody:x:65534:65534:nobody:/:/sbin/nologin
dockremap:x:100:101::/home/dockremap:/sbin/nologin
rootless:x:1001:1001:Rootless:/home/rootless:/bin/sh
EOF

docker run --rm -v docker:/work alpine chown 1001:1001 /work

docker run -d --name docker-dind \
  --security-opt seccomp=unconfined \
  --security-opt apparmor=unconfined \
  --security-opt systempaths=unconfined \
  --device /dev/fuse \
  --device /dev/net/tun \
  --user 1001:1001 \
  -v shared_socket:/var/run/user/1001 \
  -v work:/home/runner/_work \
  --mount type=volume,src=docker,dst=/home/rootless/.local/share/docker,volume-nocopy \
  -v /tmp/subuid:/etc/subuid \
  -v /tmp/subgid:/etc/subgid \
  -v /tmp/group:/etc/group \
  -v /tmp/passwd:/etc/passwd \
  docker:dind-rootless
```

!!!note NOTE

You need to use `volume-nocopy` to prevent the container from copying the
content and permissions of the `/home/rootless/.local/share/docker` directory.

!!!

And one last thing: if you need to access to the network namespace of the docker
in docker container, you'd want to connect the github runner namespace to the
docker-in-docker container namespace.

To do that, you need to add `--network "container:docker-dind"` to the
`docker run` command.

## The final script

```shell
#!/bin/sh

set -ex
# Setting up variables
: "${REPO?Need REPO}"
GH_TOKEN="$(gh auth token)"

# Clean up procedure
cleanup() {
  docker stop runner-container || true
  docker rm -f runner-container docker-dind || true
  docker volume rm docker work shared_socket || true
  if [ -n "${RUNNER_ID:-}" ]; then
    curl -fsSL -X "DELETE" \
      -H "Authorization: Bearer ${GH_TOKEN}" \
      -H "Accept: application/vnd.github+json" \
      "https://api.github.com/repos/${REPO}/actions/runners/${RUNNER_ID}"
  fi
}
trap cleanup EXIT INT TERM

# Setting up volumes and permissions
docker run --rm -v work:/work alpine chown 1001:1001 /work
docker run --rm -v shared_socket:/work alpine chown 1001:1001 /work
docker run --rm -v docker:/work alpine chown 1001:1001 /work

# Setting up the UID/GID mapping + passwd + group
echo "rootless:30000:20000" >/tmp/subuid
echo "rootless:30000:20000" >/tmp/subgid

cat <<EOF >/tmp/group
root:x:0:root
bin:x:1:root,bin,daemon
daemon:x:2:root,bin,daemon
sys:x:3:root,bin
adm:x:4:root,daemon
tty:x:5:
disk:x:6:root
lp:x:7:lp
kmem:x:9:
wheel:x:10:root
floppy:x:11:root
mail:x:12:mail
news:x:13:news
uucp:x:14:uucp
cron:x:16:cron
audio:x:18:
cdrom:x:19:
dialout:x:20:root
ftp:x:21:
sshd:x:22:
input:x:23:
tape:x:26:root
video:x:27:root
netdev:x:28:
kvm:x:34:kvm
games:x:35:
shadow:x:42:
www-data:x:82:
users:x:100:games
ntp:x:123:
abuild:x:300:
utmp:x:406:
ping:x:999:
nogroup:x:65533:
nobody:x:65534:
docker:x:2375:
dockremap:x:101:dockremap
rootless:x:1001:
EOF

cat <<EOF >/tmp/passwd
root:x:0:0:root:/root:/bin/sh
bin:x:1:1:bin:/bin:/sbin/nologin
daemon:x:2:2:daemon:/sbin:/sbin/nologin
lp:x:4:7:lp:/var/spool/lpd:/sbin/nologin
sync:x:5:0:sync:/sbin:/bin/sync
shutdown:x:6:0:shutdown:/sbin:/sbin/shutdown
halt:x:7:0:halt:/sbin:/sbin/halt
mail:x:8:12:mail:/var/mail:/sbin/nologin
news:x:9:13:news:/usr/lib/news:/sbin/nologin
uucp:x:10:14:uucp:/var/spool/uucppublic:/sbin/nologin
cron:x:16:16:cron:/var/spool/cron:/sbin/nologin
ftp:x:21:21::/var/lib/ftp:/sbin/nologin
sshd:x:22:22:sshd:/dev/null:/sbin/nologin
games:x:35:35:games:/usr/games:/sbin/nologin
ntp:x:123:123:NTP:/var/empty:/sbin/nologin
guest:x:405:100:guest:/dev/null:/sbin/nologin
nobody:x:65534:65534:nobody:/:/sbin/nologin
dockremap:x:100:101::/home/dockremap:/sbin/nologin
rootless:x:1001:1001:Rootless:/home/rootless:/bin/sh
EOF

# Register the runner
JIT_RESPONSE=$(curl -fsSL -X "POST" \
  -H "Authorization: Bearer ${GH_TOKEN}" \
  -H "Accept: application/vnd.github+json" \
  "https://api.github.com/repos/${REPO}/actions/runners/generate-jitconfig" -d '{
    "name": "rootless-docker-runner",
    "runner_group_id": 1,
    "labels": ["self-hosted","linux","x64"],
    "work_folder": "/home/runner/_work"
  }')

RUNNER_ID=$(echo "$JIT_RESPONSE" | jq -r '.runner.id')
ENCODED_JIT_CONFIG=$(echo "$JIT_RESPONSE" | jq -r '.encoded_jit_config')

docker run -d --name docker-dind \
  --security-opt seccomp=unconfined \
  --security-opt apparmor=unconfined \
  --security-opt systempaths=unconfined \
  --device /dev/fuse \
  --device /dev/net/tun \
  --user 1001:1001 \
  -v shared_socket:/var/run/user/1001 \
  -v work:/home/runner/_work \
  --mount type=volume,src=docker,dst=/home/rootless/.local/share/docker,volume-nocopy \
  -v /tmp/subuid:/etc/subuid \
  -v /tmp/subgid:/etc/subgid \
  -v /tmp/group:/etc/group \
  -v /tmp/passwd:/etc/passwd \
  docker:dind-rootless &

sleep 5

docker run --pull=always --rm --name "runner-container" \
  --network "container:docker-dind" \
  -v "work:/home/runner/_work" \
  -v "shared_socket:/run/user/1001" \
  -e DOCKER_HOST=unix:///run/user/1001/docker.sock \
  -e TMPDIR=/home/runner/_work/tmp \
  "ghcr.io/actions/actions-runner:2.337.0" \
  bash -c "mkdir -p \$TMPDIR && chmod 755 \$TMPDIR && exec ./run.sh --jitconfig $ENCODED_JIT_CONFIG" &
RUNNER_PID=$!

wait "$RUNNER_PID"

```

Simulate a GitHub jobs:

```bash
# Enter the container
docker exec -it runner-container bash

# Healthcheck docker
docker ps -a
# CONTAINER ID   IMAGE     COMMAND   CREATED   STATUS    PORTS     NAMES

# Run a server and mount a dir
docker run --rm -d --name nginx -p 80:80 -v /home/runner/_work/test:/test nginx

# Test
curl http://localhost
# <!DOCTYPE html>
# <html>
# <head>
# <title>Welcome to nginx!</title>
# <style>
# html { color-scheme: light dark; }
# body { width: 35em; margin: 0 auto;
# font-family: Tahoma, Verdana, Arial, sans-serif; }
# </style>
# </head>
# <body>
# <h1>Welcome to nginx!</h1>
# <p>If you see this page, nginx is successfully installed and working.
# Further configuration is required for the web server, reverse proxy,
# API gateway, load balancer, content cache, or other features.</p>
#
# <p>For online documentation and support please refer to
# <a href="https://nginx.org/">nginx.org</a>.<br/>
# To engage with the community please visit
# <a href="https://community.nginx.org/">community.nginx.org</a>.<br/>
# For enterprise grade support, professional services, additional
# security features and capabilities please refer to
# <a href="https://f5.com/nginx">f5.com/nginx</a>.</p>
#
# <p><em>Thank you for using nginx.</em></p>
# </body>
# </html>

# Test volume sharing
docker exec -it nginx touch /test/test

stat /home/runner/_work/test/test
```

Pretty good huh? And everything is rootless and 100% secure!

## Conclusion

Rootless Docker in Rootless Docker is not something impossible. A good
understanding of rootless technologies allows you to setup a secure
Docker-in-Docker by leveraging Linux user namespaces, UID/GID mappings, FUSE
overlayFS and `slirp4netns`.

In multi-tenant workload or CI jobs, privilege escalation vulnerabilities plague
standard Docker in Docker configurations, while maintaining full Docker
functionality.

With Rootless Docker in Rootless Docker, we have clear advantages over the
classic Docker-in-Docker and Buildah:

- **True isolation**: Even an attacker would breach all the isolations, it would
  never reach the root user.
- **No security compromises**: `--privileged` is never used and no sensitive
  mounts are used. We're running unprivileged workloads on unprivileged
  environments.
- **Buildx compatibility**: Rootless Docker in Rootless Docker allows you to use
  Docker buildx, which has many optimizations like layer squashing, cache mounts
  and more.
- **Production-Ready**: With a multi-layered isolation, multi-tenant workloads are
  possible and can run on any environment, even in production. Common usecases
  being **CI runners**, but also **AI workloads**! You wouldn't want your Agent
  to escape you containment, would you?

While the setup requires more careful configuration than traditional
Docker-in-Docker, especially around UID/GID mappings and volume
permissions, the investment in security is well worth it, especially when running
untrusted or multi-tenant workloads. For anyone building secure CI/CD
infrastructure, Rootless Docker in Rootless Docker is no longer an optional
optimization but a key part of the solution.
