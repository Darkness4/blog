---
title: Gentoo Linux is the best OS for gaming and software development on desktop.
description: "The review about Gentoo Linux after 1 year of intensive usage in gaming and development: it's the best OS in the world."
---

## Table of contents

\\{\\{ $.TOC }}

## Introduction

Since this is a review, let me share my Linux experience before talking about Gentoo Linux.

In the following order, I tried:

- [Ubuntu](https://ubuntu.com/), as my first OS installed by my big brother. Used frequently in 2012. Didn't know anything about the "insides" of Linux. I just installed Linux games through the Ubuntu software center and played flash games. Switched to Windows two years after that. Was a great experience, but couldn't play good games on Linux at the time (even with Wine).
- [Kali Linux](https://www.kali.org) in 2018, installed by myself as a script kiddy (literally) on my new gaming desktop. Worked surprisingly well and learn some stuff about hacking and programming. Still didn't learn the internals of Linux (kernel, bootloader, etc...). Used at school.
- [MX Linux](https://mxlinux.org) in 2020. Was great, but feels off, like the system was bloated (even if Kali was even more bloated).
- [Debian](https://www.debian.org) in the same year, cleaned everything, started learning Linux seriously. Since it was my first year at engineering school, I've learned C programming (with syscalls and stuff), which means I understood what a software is, and what a library is. It was the OS I used the most, but I've has to switch to the unstable branch of Debian.
- [Fedora](https://fedoraproject.org) in 2023 on my new gaming/work desktop which was a computer built by myself. The experience was reaaaally smooth and the best, especially since I've had an NVIDIA card and proprietary drivers for Wi-Fi and Ethernet. It's probably the OS that I recommend for new users.

Let's talk about the last operating system I used and the context in which I find myself in 2023. In 2023, I've started remote working for a startup. I'm doing Software Development for Web3, HPC and DevOps, which means this includes:

- Programming from low level to high level for work.
- Forking projects and installing it to test it locally.
- Lot of windows with lots of SSH connections and documentation alongside.
- Kubernetes at home.

Outside of work, that includes:

- Gaming.
- Testing open-source projects and contributing to it.
- Recreational programming from hackathon to simply start new projects.

Which adds constraints:

- Difficult dependency management with the libraries
- Have to install a lot of software inside `/usr/local`
- Having to compile everything by hand

Which leads me to switch to Gentoo in the same year. This is all so that I can do programming **in peace**: having full control over what is installed and what is running. F- SystemD (why the f- do you handle DNS?!).

At the same time, to judge Gentoo Linux, I've installed Void Linux on my laptop.

## A small review of Void Linux, a SystemD-less binary OS

I used Void Linux for a very short time, but it was enough to tell me it is worse than Arch Linux.

I know some people may say that the packing ecosystem of Void Linux is more "healthy", but the packaging software sucks hard.

One reason: `xbps-src`. If you need a package that is restricted like the Vivaldi web browser, you have to bootstrap yourself. This is simply not a great UX. I mean, what if you forget about it? I'm having trash installed all over the place! This is like my old polluted `Projects/BuildZone` folder that I've had on my old OS, what the hell.

Could be great for Raspberry Pi, but not for a development and gaming desktop.

## The review of Gentoo Linux

### Installation Review

You've come to read my propaganda, great. Let's start with the installation.

The [Gentoo Handbook](https://wiki.gentoo.org/wiki/Handbook:AMD64) tells everything a person need to know to install Gentoo Linux. And none of these steps has caused major issues. I've used the binary kernel the first time, and trimmed it [later](about-kernel-configuration).

You can also use any bootable Linux OS to install Gentoo Linux. I've personally used the Void Linux installation ISO since it's only ~600MB, and it has a package manager so that I can install missing tools.

To summarize my installation:

- Used `cfdisk` to partition my disk and `mkfs.xfs` for the file-system.
- Installed Stage3 OpenRC with Multilib by extracting the `starge3-[...].tar.gz`.
- Configured my `/etc/portage/make.conf`. This is my actual `make.conf`:

  ```shell
  CC="gcc"
  CXX="g++"

  AR="${CHOST}-ar"
  NM="${CHOST}-nm"
  RANLIB="${CHOST}-ranlib"

  COMMON_FLAGS="-march=native -O2 -pipe"
  CFLAGS="${COMMON_FLAGS}"
  CXXFLAGS="${COMMON_FLAGS}"
  FCFLAGS="${COMMON_FLAGS}"
  FFLAGS="${COMMON_FLAGS}"
  LDFLAGS="-Wl,--as-needed"
  RUSTFLAGS="-C target-cpu=native"

  LC_MESSAGES=C
  MAKEOPTS="-j12"

  USE="-telemetry -networkmanager -oss -qt5 -kde -systemd -networkmanager -wayland qt6 bluetooth gtk X pulseaudio dbus pipewire gles gles2 elogind proprietary-codecs pam xorg truetype udisks jpeg png zstd qemu vnc midi alsa appindicator nvenc cuda vulkan x264 x265 srt openssl -webengine mtp ffmpeg openexr -debug vdpau nvidia spell fontconfig theora opus vaapi cups -nls"
  INPUT_DEVICES="libinput"
  VIDEO_CARDS="nvidia"
  ACCEPT_LICENSE="*"
  #GENTOO_MIRRORS="https://mirrors.ircam.fr/pub/gentoo-distfiles/"
  QEMU_SOFTMMU_TARGETS="x86_64 aarch64 s390x ppc64 riscv64"
  QEMU_USER_TARGETS="x86_64 aarch64 s390x ppc64 riscv64"
  L10N="en en-US"
  ```

  Finding the "right" USE flags is quite complicated, it's better to search for examples. Beware of the `-nls` (Native Language Support) which removes translations.

- Installed the binary kernel and built the initramfs (initial ram-based file system, used to boot the OS and load the kernel) with Dracut.
- Installed `iwd` standalone without `netifrc` and `dhcpcd`.
- Installed `sysklogd` `chrony`, rEFInd as bootloader, kicked grub from my old Fedora installation.

You WILL end up with a running Gentoo Linux, no matter what. But, one big issue that I've had: "Now that I have Gentoo Linux running, where the f- is my desktop environment? Where is my display manager?"

You guessed it, time to go outside the handbook, section [Display Manager](https://wiki.gentoo.org/wiki/Display_manager) and [Desktop Environment](https://wiki.gentoo.org/wiki/Desktop_environment).

"Wait... but that's experiential knowledge! How can a new user knows that we need to install a DM, or a DE?! What the f- is X11, or X.org, or even Wayland?!"

That's right, you Google-it: "How to install Desktop Gentoo". Then, [Desktop Environment](https://wiki.gentoo.org/wiki/Desktop_environment), then [GNOME/Guide](https://wiki.gentoo.org/wiki/GNOME/Guide), then [Display Manager](https://wiki.gentoo.org/wiki/Display_manager). You know you won't finish the installation in one day because you forgot about Audio, Screenshot software, Network, Multimedia..., so you switch to Windows for the time being. See the [Recommended applications page](https://wiki.gentoo.org/wiki/Recommended_applications).

Even if the installation of Gentoo Linux is pretty smooth, some steps are somewhat implicit. I wouldn't say it's particularly hard, but I would say it takes more time to install everything that is needed. And, I haven't talked about the USE flags.

### Package Management Review

Most operating systems are often judged by the package manager. Some people will switch from Ubuntu to Arch Linux or Fedora simply because `apt` is slow.

So let's talk about Gentoo's package manager: Portage. Portage is a source-based package manager, which also combines a build system.

Portage uses `rsync` or `git` to synchronize `ebuild` repositories, `ebuild` being the recipe for a package. Portage is capable of resolving build dependencies, runtime dependencies and post-dependencies, allowing unused dependencies to be cleaned up.

Portage uses `USE` flags to customize a package, and `USE` flags can also change the dependencies of the package, meaning you can trim a package to the minimum (for example, `ffmpeg` with only `x264` and no other encoders).

Portage can also use third-party repositories through Overlays, which also includes your own overlay. Meaning, you can change a recipe and install it directly without any extra steps.

Since we are compiling all the software, you can customize the compilation flags inside the `/etc/portage/make.conf` and add `-march=native` to optimize for the CPU. You can also optimize with Link-Time Optimization (LTO) and Profile Guided Optimization (PGO).

Even if Portage is a source-based package manager, it is also able to install binaries, especially for software without any dependencies (Firefox, Vivaldi, the Linux kernel...). If the package cannot not found, there is still flatpak.

Lastly, there are no runtime/linking issues due to the fact it is checked during installation time, i.e., no more `missing *.so` or `undefined reference`.

When using it with a good computer, installing packages are not that long. I've never had to wait at night for everything to compile. And when everything breaks, it is possible to rollback packages and/or debug the compilation steps since the packages are being compiled in the `/var/tmp/portage` directory. But that's even rare, even on the unstable branch `~amd64`.

However, there is one major drawback (for non-programmers), build dependencies are kept on the computer to avoid multiple re-downloads, making the operating system bloated with toolchains. However, since I'm a software developer for HPC and Web3, I'm using these toolchains anyway.

Overall, Portage is probably one of the greatest package manager of all time.

### About Kernel Configuration

After 6 months of usage, I wanted to optimize my kernel and kernel modules. I configured my OS to use a `zstd`-compressed `linux-firmware` and `zstd`-compressed kernel modules.

Then, I've removed everything through the `make menuconfig` of the kernel. I did everything legit without `genkernel`.

![image-20231214174137405](./page.assets/image-20231214174137405.png)

Tricks I used:

- Use <kbd>H</kbd> to show the description of the kernel feature.
- Use <kbd>/</kbd> to search for a kernel feature and its reverse dependencies. Use <kbd>0</kbd> through <kbd>9</kbd> to go to the right page.

Most importantly, I've removed anything related to embedded systems and enterprise-class systems. I've removed unused audio, file-system, Wi-Fi and Ethernet drivers.

Overall, it takes a full day to check everything settings. Starting from scratch is very dangerous, so I recommend to use the kernel config of Fedora, which is the [kernel config of Gentoo used for the binary kernel](https://github.com/projg2/fedora-kernel-config-for-gentoo). You can also use [snippets shared by the community](https://codeberg.org/ranguli/gentoo-popcorn-kernel/). I've managed to reduce the size from 200MB to 87MB of kernel modules, the heaviest being NVIDIA (44MB).

### Stability/Maintenance Review

It's been one year, and nothing happened. Literally. There was a major version change with Python (3.10 to 3.11) and a Pipewire breaking change, but nothing really happened. Most breaking changes are notified through the news:

```shell
eselect news read
```

Besides that, everything has been stable.

### For other stuff review

**Gaming**

Let's talk about gaming with x11. Installing and playing is super-easy through Steam or Wine. However, there is a clear performance drop on Linux compared to Windows.

It's playable, but not as comfy as on Windows.

**Programming**

It's heaven. You can fork dependencies easily, install any version of a library. And you can use [Portage with Crossdev for easy cross-compilation](2023-11-08-go-with-portage-and-crossdev).

**Web browsing**

Yes.

**MIDI and stuff**

The documentation is a quite short. Just don't forget the MIDI driver inside the kernel and ALSA. It works quite well and I was able to use Musescore with my MIDI piano (Kawai ES920).

**OpenRC**

It's fast, it's easy to add our own service. Example with the `iwd`:

```shell
#!/sbin/openrc-run
# Copyright 1999-2020 Gentoo Authors
# Distributed under the terms of the GNU General Public License v2

pidfile="/run/iwd.pid"
command="/usr/libexec/iwd"
command_background="yes"

depend() {
    need dbus
    after bootmisc modules
    before dns dhcpcd net
    keyword -shutdown
}
```

**Movies**

I used `mpv` with `ffmpeg`. My USE flags are `fdk libass x264 x265 chromium nvenc vaapi vpx theora opus abi_x86_32`, which was added depending on my needs, or added automatically through `dispatch-conf`.

### Learning curve... what about it?

Truth to be told, the learning curve is not that hard. Not as hard as compiling the kernel anyway.

The [Gentoo Handbook](https://wiki.gentoo.org/wiki/Handbook:AMD64) tells everything that we need to set up a Gentoo Linux OS. Even using the binary kernel is not that hard.

The installation steps are not hard to follow, and the wiki is filled with enough description and troubleshooting to avoid any errors.

Even maintenance-wise, I haven't seen any disastrous breaking change.

This is certainly not an OS for newbie, i.e, non-C programmers, but this is also certainly not an OS for hardcore Linux users living in their basement with programming socks.

Basically, If you can install Arch Linux, there is no reason to try Gentoo Linux, especially if your computer has a beefy CPU.

## Tackling the myths about Gentoo Linux

### "It takes time to compile, and you compile everyday"

Simply false. On a good computer, it takes maximum one hour. And no, you don't need to compile every time everything, unless you WANT it.

### "It's unmaintainable if you forget about it for one month"

It's the same experience as Arch Linux. Heck, it's safer than Arch Linux because there are two branches (`amd64` and unstable `~amd64`), so you can have the same experience than a standard release distro.

I had to go on vacation for a month, and nothing broke.

The worst experience I had was with Gentoo Prefix (Gentoo as a Linux subsystem), which has to be compiled from step 1. And then everything broke because of a conflict between the host's GLIBC and Gentoo Prefix's GLIBC.

### "It's bloated with dev dependencies"

As a programmer, this is non-issue. Compared with other binary operating systems, Arch Linux is bloated by SystemD, which is much worse. Void Linux is not bloated, but its package manager sucks (tell me the GCC version vs Gentoo or Arch). And I haven't talked about how the kernel is bloated on binary OSes.

There are so much useless kernel modules that we don't need.

### "It's for expert"

There is documentation, and everything is explained. It's not for non-programmer, but come on, it's not like a Linux amateur can't install it and use it every day.

### "It's optimized as hell"

That's the last point which is not in favor for Gentoo. Binaries offered by binary OSes are often more optimized thanks to PGO. PGO takes too much time on Gentoo Linux, therefore it's not worth it to enable it on Gentoo.

On Gentoo, there is one optimization that Gentoo users benefit, which is `-march=native` which allows optimization by using optimized CPU instructions.

But to say it's "optimized as hell" is a generalization. All I can say is that OpenRC is way faster than SystemD.

## Conclusion

Obligatory `fastfetch`:

![image-20231214180934560](./page.assets/image-20231214180934560.png)

That's my review for Gentoo Linux. Just sharing my love for it.

Gentoo Linux give the user full control on what is installed and what is running, without ruining the experience. It's actually addicting how stable it is. And when it crashes, you know it's your fault: either you messed up with the initramfs or kernel (or maybe it's NVIDIA's fault, which I understand). And when it's not your fault, you can patch it, install it through Portage and share a pull request on GitHub.

Installing and using Gentoo Linux makes you understand what's behind a Linux distribution, and you understand what is bloat. You understand that what you've installed are not spyware (vs [Ubuntu](https://www.omgubuntu.co.uk/2018/05/this-is-the-data-ubuntu-collects-about-your-system) and [Fedora](https://fedoraproject.org/wiki/Changes/Telemetry)). Just run `htop`, and you are able to tell which application does what, compared to other binary OSes.

It's the **best OS for desktop that has ever existed**.

If you are reading this article and hesitating of installing it: just do it. You will gain experience and zen. You will **STOP** distro hopping.
