---
title: Presenting my home Raspberry Pi cluster
description: Presenting my home Raspberry Pi Kubernetes cluster which is hosting this blog.
tags:
  [raspberry pi, hpc, kubernetes, cluster, home, monitoring, storage, devops]
---

## Table of contents

<div class="toc">

\\{\\{ $.TOC }}

</div>

## Introduction

I've finally finished upscaling my home Raspberry Pi cluster to 3 nodes to be able to run my monitoring services. I'm pretty proud of it, so I wanted to share it with you all. Here's a diagram of the setup:

```d2
modem -> router: Wi-Fi  {
  style: {
    stroke-dash: 3
    stroke: green
  }
}
router -> node1 {
  style: {
    stroke: green
  }
}
router -> node2 {
  style: {
    stroke: green
  }
}
router -> node3 {
  style: {
    stroke: green
  }
}
router -> storage1 {
  style: {
    stroke: green
  }
}
storage1 -> node1 {
  style.animated: true
}
storage1 -> node2 {
  style.animated: true
}
storage1 -> node3 {
  style.animated: true
}

node1 -- power supply {
  style: {
    stroke: red
  }
}
node2 -- power supply {
  style: {
    stroke: red
  }
}
node3 -- power supply {
  style: {
    stroke: red
  }
}
storage1 -- power supply {
  style: {
    stroke: red
  }
}
router -- power supply {
  style: {
    stroke: red
  }
}
```

The specifications of the nodes are as follows:

- Router: Raspberry Pi 3 Model B Plus Rev 1.3. It runs Raspbian with an ARP proxy called parprouted. More details later.
- Nodes: Raspberry Pi 4 Model B Rev 1.4 8GB or Rev 1.5 4B. They run on k3os and have a 32GB SD card.
- Storage: ODroid XU4 with a 1TB SSD.

The motivation for this project is very simple: I wanted to self host my services: DNS, Blog, Password Manager, etc... I also wanted to learn more about Kubernetes and how to manage a cluster.

The architecture is very similar to an HPC cluster, but on a smaller scale. The storage is separated from the worker nodes, and the network layer (interconnect) is compact and efficient. To make the cluster portable (due to my home), I've set up a router which acts as an ARP proxy to route the traffic toward the internet. In other words:

- The storage and nodes are interconnected with one switch, using simply the L2 layer, permitting a high-speed connection between the nodes and the storage.
- The storage and nodes are exposed to the internet via the router, a much slower connection, but enough for my needs. The bottleneck at my home is actually the lack of Wi-Fi 6E, which would permit a faster connection between the router and the modem.

Now that you have an overview of the cluster, let's dive into the details. I will show my setup, the software I use, and how I manage the cluster.

## Building the cluster

### Storage

The storage is an ODroid XU4 with a 1TB SSD. The choice wasn't one: a friend gave it to me for free. While a lot of people would have used a preconfigured NAS, I wanted to avoid any abstraction layer, which would permits me to scale up/explore new technologies, the long term goal being to set up a parallel filesystem like Lustre or Ceph.

The storage node has 2 storage: the rootfs and the data. The rootfs is a simple 32GB SD card, and the data is a 1TB SSD. The SD card is formatted with ext4, and the SSD is formatted with btrfs.

Why not ZFS? Cause btrfs is integrated into the Linux kernel, and I wanted to avoid any kernel module. I've also used btrfs in the past as a main filesystem, and had good experiences with it.

Volumes are exported via NFS, CIFS (yuck) and FTP. The reason: for max compability with my devices. Android works well with FTP and CIFS. Linux works well with NFS.

Why not Ceph? Cause I don't have enough nodes to set up a proper Ceph cluster, nor I have the memory to run it. Having run BeeGFS and Longhorn in the past, I know that a parallel filesystem is a lot of work to maintain and needs a lot of resources.

Why not iSCSI? Cause NFS is more convenient. I'm actually more interested in setting up a parallel filesystem than a block storage.

What's now and the future? ODroid XU4 is running on ARMv7, which is a bit old and may lack in support in the future. The next step would be to scale up to an SBC with ARMv8, maybe the ODroid M1 or HC4, which are specifically made for storage. To add one more reason, I want to install Kubernetes on the storage node, which would allow me to run Rook (easy Ceph for Kubernetes). The HC4 is quite old, but offers 2 SATA ports, while the M1 is more recent and offers more performance. The benchmarks given by Hardkernel:

- ODroid HC4 (SATA): Random Write 1M: 427.137
- ODroid M1 (NVMe): Random Write 1M: 1163.95

To be decided.

### Power supply

The router and storage is powered directly on the wall. The Raspberry Pi 4, using USB-C, is using an Anker power supply USB-A hub.

Why not plugging everything to the hub? For the storage, it's for security reasons. For the router, the monitoring system alerted me that the power supply was not enough. Plus, since the Rapsberry Pi 3 is using micro-USB, the cable is sh-t and unreliable.

### Router

At the beginning, the router wasn't part of the project. It was a random router found by my brother running on OpenWRT. The "relay" worked pretty well to relay L2 packets (ARP, DHCP, etc...), but the Wi-Fi was terrible. Since I had a Raspberry Pi 3 lying around, I decided to use it as a router. The Raspberry Pi 3 had a better antenna (could at least run on 5 GHz), but the ARP proxy could relay the DHCP packets, only the ARP packets. Therefore, I've set up static IPs on the nodes and storage, and the router is now acting as a simple router.

Quite disappointing to be honest since on HPC clusters, the DHCP server is able to network boot the nodes, which would have been the next step.

The discussion about the setup is [here](https://raspberrypi.stackexchange.com/questions/88954/workaround-for-a-wifi-bridge-on-a-raspberry-pi-with-proxy-arp). Around the internet, some people said it's an issue with the firmware, which sucks.

Anyway, it's working well now and I'm waiting to upgrade to a Wi-Fi 6E network.

### Nodes

Time for the most interesting part: the worker nodes. Actually, one of the nodes is actually a controller node too, but due to the lack of resources, I've decided to run both roles on the same node. (To be honest, I have some Raspberry pi 3 lying around which could have served me as controller nodes, but I feared my USB power supply hub wouldn't be able to power everything).

Each node is running on k3os, an immutable OS made for Kubernetes. By immutable, I mean that the OS is read-only and the user data is stored in a separate writable partition. This pattern is now very popular thanks to the Steam Deck and various consumer OSes (Fedora Silverblue, Ubuntu Core, etc...). In past articles, I've talked about how to setup an immutable OS thanks to OverlayFS. On k3os, the setup is much more simpler:

- The OS (squashfs) is read-only and mounted on `/usr` with the `ro` flag.
- `/etc` is ephemeral, meaning that it's stored in a tmpfs and is lost at reboot.
- `/boot` is the EFI partition, stored on the SD card.
- `/k3os/system` is the boot configuration files, read-only, stored on the SD card.
- The user data (ext4) is mounted on `/` with the `rw` flag, making `/home`, `/var`, `/opt`, `/usr/local` persistent.

Or in other words:

```shell
$ findmnt

TARGET                SOURCE                 FSTYPE      OPTIONS
/                     /dev/mmcblk0p2[/k3os/data]
│                                            ext4        rw,relatime
├─/usr                /dev/loop0             squashfs    ro,relatime
│ ├─/usr/lib/modules  /dev/mmcblk0p2[/lib/modules]
│ │                                          ext4        rw,relatime
│ └─/usr/lib/firmware /dev/mmcblk0p2[/lib/firmware]
│                                            ext4        rw,relatime
├─/etc                none                   tmpfs       rw,relatime
├─/proc               none                   proc        rw,relatime
├─/run                tmpfs                  tmpfs       rw,nosuid,nodev,noexec,relatime,size=383296k,mode=755
├─/tmp                tmpfs                  tmpfs       rw,nosuid,nodev,noexec,relatime,size=383296k
├─/dev                dev                    devtmpfs    rw,nosuid,relatime,size=10240k,nr_inodes=437696,mode=755
├─/sys                sysfs                  sysfs       rw,nosuid,nodev,noexec,relatime
├─/boot               /dev/mmcblk0p2[/boot]  ext4        rw,relatime
└─/k3os/system        /dev/mmcblk0p2[/k3os/system]
                                             ext4        ro,relatime
```

Today, k3os is discontinued, but, technically, the techniques used in k3os are still valid and doesn't mean that k3os is dead. To summarize on how k3os works:

- K3os is using a custom init program to setup the immutable OS. To be more precise:
  - The Raspberry PI bootloader loads the OS according the `boot/cmdline.txt` and the `boot/config.txt` files.
  - K3Os runs `/sbin/init.preinit` (which just runs checks), and then run a "custom" `/sbin/init`, which is runs scripts in the `/usr/libexec` directory:
    - First, the `bootstrap` script runs `k3os config --initrd`, which setups the OS (mount filesystems, etc...).
    - Second, the `mode` script, which runs the `mode-local` script, which setup ssh and the k3os mode.
    - Third, the `boot` script, which runs the `config --boot` command, which setups the services (DNS, WiFi, SSH, k3s, etc...).
    - Lastly, the `/sbin/init` is executed, which is the real init system from Busybox, which executes OpenRC.

As you can see, it seems quite complicated, but the tools used are pretty much native to Linux. Thanks to this setup, k3os is versioned and stored in the path `/k3os/system`, and the user data is stored in the path `/k3os/data`.

Not only that, but this makes easy to upgrade the hardware as I just have to move the SD card to the new hardware, with nothing to change.

## Setting up Kubernetes and the services

### Core services

K3os comes with preconfigured services, such as:

- Traefik ingress controller
- CoreDNS, the DNS server
- Flannel, the CNI plugin
- KlipperLB, the load balancer for bare-metal

The issue with this is that Traefik and CoreDNS are not configured to my needs. Traefik is in heavy development and CoreDNS needs custom rules. Therefore, I've decided to remove them and install them manually.

We also need to add `cert-manager` and `sealed-secrets` to handle certificates and encrypted secrets.

### Monitoring services

At the beginning, I tried to use Prometheus but was way too heavy and hogs the resources. I discovered VictoriaMetrics, which is a lightweight version of Prometheus. Nowadays, VictoriaMetrics is pretty stable and works well with Grafana. I decided to also not install Alertmanager, as Grafana has its own.

Node-exporters are installed on each machine to monitor the resources. Alerts are configured in Grafana.

These include:

- CPU usage
- Memory usage
- Disk usage
- Health of the nodes
- Network usage
- Temperature of the nodes
- Power supply voltage

And this is where I detected that one of my Raspberry Pi was under-powered, and decided to plug it directly to the wall.

Containers are also monitored thanks to the built-in monitoring system of Kubernetes:

- Node status
- Disk usage
- Memory usage
- Pod status
- Deployment status
- Job status
- Container status
- Volume status

Thanks to this, Lens (a GUI for Kubernetes) is able to quickly show me the status of the cluster.

The monitoring system is now on point, with only two things missing: logs and traces. I'm currently looking into Loki and Jaeger/Tempo to complete the monitoring system.

### Application Services

Without going into much details, here are the services I'm running:

- Various web services from my projects (blog, pilot projects, etc...)
- Password manager VaultWarden
- Note-taking database Joplin
- Filebrowser, to browse the files on the storage (there is a small bottleneck here, but this is mostly for reading than writing)
- My bots to archive and remux videos
- My Gotify notification server
- The CoreDNS server which is accessible from the local network. A custom `hosts` file is used to adblock some domains.

Everything is deployed via Kustomize or Helm, with the complement of `cert-manager` and `sealed-secrets`.

The deployment of services is push-based, meaning it's not entirely automated, but still tracked in Git.

### Storage services

Even if Kubernetes is able to host many services, there is an obvious bottleneck: NFS. NFS is unstable and has some "syncing/locking" issues. When running qBittorrent, the software is IO blocked and hog a lot the CPU. Therefore, I've decided to run qBittorrent on the storage node.

More precisely, qBittorrent is running in a home-made container: a user namespace. To avoid any security issues, the container is running rootless, in a user namespace, plugged to a VPN. Obviously, there is a reduction in bandwidth due to the VPN (Wireguard uses UDP encapsulation), but it's a good trade-off for security.

Overall, the storage had never failed me, and I'm pretty happy with it. The only issue is the lack of a parallel filesystem, which would permit me to scale up the storage.

## Conclusion

Not very much to say. The cluster is running well, and I'm happy with it. Thanks to this architecture, I'm able to move around the whole cluster in my house and do easy maintenance, something that I wouldn't do with a homemade NAS. Feels pretty great to just unplug the power strip, place it in a small box, and move it to another room.

<center>

<img src="/blog/2024-06-19-home-raspi/page.assets/signal-2024-06-19-201128_003.jpeg" alt="signal-2024-06-19-201128_003" style="zoom:25%;" />

_The cluster on the shelf, wirelessly connected to the modem. There is a 120 mm fan behind the cluster._

<img src="/blog/2024-06-19-home-raspi/page.assets/signal-2024-06-19-201128_002.jpeg" alt="signal-2024-06-19-201128_002" style="zoom:25%;" />

_The cluster is quite small. Not as small as [Jeff Geerling](https://www.jeffgeerling.com/blog/2024/turing-rk1-2x-faster-18x-pricier-pi-5) 😭. At least, it's eco-friendly... right?_

</center>

Theoretical maximum power draw is:

- 1 Rpi 3B+: 12.5W
- 3 Rpi 4B: 45W
- 1 ODroid XU4: 20W
- 1 Switch 1G: ~13W

Total: 90.5W for 1To of storage, 12 cores, 16GB of RAM.

Obviously, the power draw is certainly lower. But I don't have the necessary tools to measure it.

Other than that, the next steps would be to upgrade the storage node and the router, and to add more nodes to the cluster. I'm also looking into setting up a parallel filesystem, but this is a long-term goal. The lack of proper backup could also be an issue (I'm mean, I'm using btrfs, but none of its features). A proper NAS like TrueNAS offers multiple type of RAID, here I have only one disk.

And about K3os, while I'm quite happy about it, I'm a little saddened that it's discontinued, especially when its core is robust and works well. At this point Talos Linux seems to be the best alternative, or I would do my own solution with Gentoo Linux, but I'm not sure if it's worth the switch.
