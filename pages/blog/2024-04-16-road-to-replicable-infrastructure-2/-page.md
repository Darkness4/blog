---
title: 'Road to replicable infrastructure 2: Kubernetes on immutable OS'
description: Deploy an immutable OS with Kubernetes
tags: []
---

## Table of contents

<div class="toc">

{{% $.TOC %}}

</div>

<hr>

## Introduction

Achieving a replicable infrastructure permits easy maintenance, scaling, and debugging. Most companies will prefer to use a managed Kubernetes (deployed with Terraform) to achieve this. The OS is hidden from the user, and Kubernetes is the only interface to the infrastructure.

However, this is not applicable to all companies. Some companies may prefer to self-host their Kubernetes clusters, and they have to install the OS themselves. The general approach is often:

1. Install the OS, which is often Ubuntu or RHEL.
2. Use Ansible to provision the OS, which is also used to install Kubernetes.

In the [previous article on replicable infrastructure](/blog/2023-09-16-road-to-replicable-infrastructure), we discussed the issues about using Ansible to provision the OS:

- Ansible is not strictly declarative as stateful servers can retain a dirty state. Badly written Ansible playbooks can lead to an undefined behavior.
  - Example: If the playbook creates a file that already exists, it will fail (or worse, it will overwrite the file).
  - Most Ansible playbooks developers must think about idempotency.
- Ansible does not track the state of the server. It is not easy to know what has been done on the server.
- Rolling updates involve risks. Not all nodes are the same: different hardware, different OS, different roles, etc.

To summary, the issue I have with Ansible is: "The initial state of the server is not known.". This is why I prefer a more "pull-based" approach:

1. The server pulls the OS.
2. The OS pulls the post boot scripts, which is also used to install Kubernetes.

With this approach, the initial state of the server is known: the immutable OS. The server is immutable, and the post boot scripts are idempotent.

The benefits of this approach are immediate:

1. Rebooting the server is safe, and the state is known during the whole lifecycle of the server.
2. It is easily scalable.

This behavior is very much similar to a container where the container starts from an immutable image and the entry point is used to configure the container.

## A summary of the first article

### Provisioning the server

In the first article, I used a PXE server to deploy an immutable OS. The OS is packed as a squashfs image, and the server boots from this image by using OverlayFS. The base image is read-only, and the changes are stored in a separate volume. After mounting the volumes, the server runs a systemd service that pulls the post boot scripts from a Git repository, which is used to provision the server.

```d2
shape: sequence_diagram
  Server -> DHCP/PXE.boot: Get IP address, ask for network boot
  DHCP/PXE.boot -> Server: Boot on initrd with kernel and boot parameters
  Server -> TFTP.squashfs: Get squashfs image
  TFTP.squashfs -> Server
  Server -> Server: Mount squashfs image as lowerdir with OverlayFS,\n and mount a writable volume as upperdir
  Server -> Server: Start systemd
  Server -> Git.scripts: Systemd service pulls post boot scripts from Git
  Git.scripts -> Server
  Server -> Server: Run post boot scripts
```

**To simplify the demonstration, I will use a virtual machine to deploy the OS.**

### Building the OS image

To build the OS image, I used Gentoo's Portage to build the OS and Kernel from scratch. Since the objective of this article is about decoupling the OS from Kubernetes, I would love to reuse this approach for the demonstration, but, for the sake of simplicity, I will use Alpine's APK to build the OS. If you want to, you can also use APT or DNF to build the OS.

To summary, the OS image is built with the following steps:

1. Install packages in a different root directory.
2. Configure the OS rapidly (hostname and password).
3. Implement a systemd service to provision the server.
4. Building and packing the kernel and kernels modules.
5. Building the initrd with dracut to enable booting with OverlayFS.

## Our approach

The objective is to create a light base OS image where Kubernetes can be installed. Luckily, **k0s, a lightweight Kubernetes distribution, can be installed on any Linux thanks to its AIO approach**. This is perfect for our demonstration.

We will set up the file-system as follows:

```shell
/etc - ephemeral
/usr - read-only (from the OS image)
/config - "persistent" (only used at boot to configure the server)
/home - persistent
/var - persistent
/opt - persistent
/usr/local - persistent
```

This approach is very similar to most immutable OS like [k3os](https://github.com/rancher/k3os) or SteamOS.

The steps are:

1. Build the OS image and kernel with Alpine's APK.
2. Build the initrd with dracut.
3. Develop the OpenRC service and post boot scripts.
4. Deploy the OS on a virtual machine.
5. Fetch the kubeconfig from the virtual machine.

## Building the OS image

As always, to avoid any side effects, we will use Docker/Podman to run Alpine Linux to build the OS image. Do note that Alpine Linux follows a rolling release:
