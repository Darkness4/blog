---
title: Setting up Yubikey GPG with LUKS and Dracut
description: Did you know that Dracut natively supports LUKS with Yubikey?
tags: [devops, linux, infrastructure, dracut, luks, gitops]
---

## Table of contents

<div class="toc">

{{% $.TOC %}}

</div>

## Introduction

Full disk encryption with LUKS is the most common setup to protect your data at
rest. However, there is drawback: storing the secret is not convenient. Either
you have to remember a secret password, or you have to store a secret key
somewhere.

To avoid entering their password twice, most people will setup a secret key, but
store it alongside the initramfs, in the boot/EFI partition! This is a big
vulnerability as the key can be easily stolen if the attacker make the simple
effort to read the boot/EFI partition. It defeats the purpose of LUKS.

In this article, I will show you how to properly setup LUKS with a Yubikey (or
any GPG smartcard) and use it with Dracut.

## The in-depth explanation

### The boot process

Let's talk quickly about the Linux boot process. Linux boot process is:

1. Load the EFI Partition.
2. Load Grand Unified Bootload (GRUB).
3. Load the Kernel (/boot/vmlinuz-<kernel version>).
4. Load the initramfs (/boot/initramfs-<kernel version>.img).
5. Call the init system (/sbin/init).
6. Boot SystemD and the services.

The "play" of LUKS is around the step 4 and 5.

### The initramfs and Dracut

The initramfs is a lightweight Linux OS that is invoked during the boot process
to setup:

- The kernel modules
- The firmware and microcode loading
- The root filesystem mounting (and various other mounts)

Or, it can also do... nothing. Dracut is the software used to configure the
initramfs and instruct what to do to boot the system. Some people uses for
network booting, diskless booting, ... and in our case, decrypting the root
filesystem before mounting it.

### LUKS, GPG and Yubikey

Before going too much in details about Dracut modules, let's just explain what
are the roles of LUKS, GPG and Yubikey.

LUKS encrypts the root filesystem with a symmetric key. That symmetric key is
generated using a cryptographic pseudo-random number generator (CPRNG).

To store the symmetric key, we can use PGP (Pretty Good Privacy) to encrypt the
symmetric key with a GnuPG public (asymmetric) key. And the private key
dedicated to decrypt the encrypted key is stored on the Yubikey, the secure
medium.

```d2 {layout="elk", title="LUKS, GPG and Yubikey"}
CPRNG: {
  label: "CPRNG\n(Cryptographic PRNG)"
  shape: oval
}

symmetric_key: {
  label: "Symmetric Key"
  shape: document
}

root_fs: {
  label: "Root Filesystem"
  shape: cylinder
}

luks: {
  label: "LUKS"
  shape: rectangle
}

encrypted_key: {
  label: "Encrypted Symmetric Key"
  shape: document
}

gpg_public: {
  label: "GnuPG Public Key\n(asymmetric)"
  shape: document
}

gpg_private: {
  label: "GnuPG Private Key\n(asymmetric)"
  shape: document
}

yubikey: {
  label: "Yubikey\n(secure medium)"
  shape: hexagon
}

CPRNG -> symmetric_key: generates
symmetric_key -> luks: used by
luks -> root_fs: encrypts

symmetric_key -> gpg_public: encrypted using
gpg_public -> encrypted_key: generates

gpg_private -> yubikey : stored
encrypted_key -> gpg_private: decrypted using
gpg_private -> symmetric_key: recovers
```

To summarize:

- The **LUKS encrypted key** is stored in the boot partition.
- The **GPG public key** is embedded in the initramfs.
- The **GPG private key** is stored in the Yubikey.

No secret is persisted and exposed.

To do the "play" indicated above (LUKS key recovery and root filesystem decryption), Dracut modules are used.

### Dracut modules, the final pieces

The dracut modules dedicated LUKS aren't very complex and some people has
already remade them even though the solution already exists (probably because it
isn't very much documented). The flow is simple:

**Build phase**

- Embed cryptsetup and linux kernel modules required to mount the root
  filesystem.
- Embed gpg and the gpg-agent. Make sure gpg has been compiled with `smartcard`
  and `libusb` support.
  - Also embed `scdaemon` to manage smartcards.
  - Also embed the user's GPG public key, to identify which private key to use.
- Load the kernel modules to also display the prompt.

**Runtime phase**

After loading and parsing all the kernel command line parameters, and upon reaching the
`crypt` module `cmdline` hook, the following happens:

- Is the key found?
  - If yes:
    - Is smartcard supported? - If yes, call `gpg --card-status` to healthcheck then
      prompt the user to enter the PIN of the Yubikey (using `gpg --pinentry-mode=loopback --decrypt`).
    - Else, prompt the user to enter the GPG passphrase (using `gpg --decrypt`)
  - Else:
    - Prompt the user to enter the LUKS passphrase (using `cryptsetup luksOpen`)

Upon decrypting the root filesystem, `cryptsetup` will mount the root filesystem
to `/dev/mapper/<name>`, in which the user has set the `root=` parameter to this
path, permitting the system to boot.

That's it! All of this is explained at:

- [Full Disk Encryption From Scratch - Gentoo wiki](https://wiki.gentoo.org/wiki/Full_Disk_Encryption_From_Scratch)
- [dracut/modules.d/73crypt-gpg/README - dracut-ng/dracut](https://github.com/dracut-ng/dracut/blob/abfb29003bf3f7462caa3c07f9361bb8131c49f5/modules.d/73crypt-gpg/README)

## The guide

### System preparation

Assuming you have booted on the USB drive, the following sections will be about
how to setup LUKS. You can skip it if you already have LUKS setup.

#### Disk preparation

We'll setup the disk as follows:

```shell
Device Location     Label        Mountpoint  Size   Filesystem  Usage
/dev/nvme0n1
 ├── /dev/nvme0n1p1 [EFI]        /boot/efi   1 GB   fat32       Bootloader
 └── /dev/nvme0n1p2 [BOOTX]      /boot       1 GB   ext4        Bootloader support files, kernel and initramfs
 └── /dev/nvme0n1p3 [ROOT]       (root)      ->END  luks        Encrypted root device, mapped to the name 'cryptroot'
      └──  /dev/mapper/cryptroot /           ->END  ext4        Root filesystem
```

#### Disk partition setup

##### Enter fdisk and setup GPT label

```shell {title="root@usblinux:/"}
fdisk /dev/nvme0n1
```

```shell {title="root@usblinux:/ # fdisk"}
Welcome to fdisk (util-linux 2.38.1).
Changes will remain in memory only, until you decide to write them.
Be careful before using the write command.

Device does not contain a recognized partition table.
Created a new DOS disklabel with disk identifier 0x81391dbc.

Command (m for help): g

Created a new GPT disklabel (GUID: 3E57DFCE-CDD9-6F42-8418-F0B6B4A08294).
```

##### Create the EFI System Partition (ESP)

```shell {title="root@usblinux:/ # fdisk"}
Command (m for help): n

Partition number (1-128, default 1): ↵
First sector (2048-1953525134, default 2048): ↵
Last sector, +/-sectors or +/-size{K,M,G,T,P} (2048-1953525134, default 1953523711): +1G

Created a new partition 1 of type 'Linux filesystem' and of size 1 GiB.
```

Add the **EFI System** property:

```shell {title="root@usblinux:/ # fdisk"}
Command (m for help): t

Selected partition 1
Partition type or alias (type L to list all): 1
Changed type of partition 'Linux filesystem' to 'EFI System'.
```

##### Create the linux boot partition

```shell {title="root@usblinux:/ # fdisk"}
Command (m for help): n

Partition number (2-128, default 2): ↵
First sector (2099200-1953525134, default 2099200): ↵
Last sector, +/-sectors or +/-size{K,M,G,T,P} (2099200-1953525134, default 1953523711): +1G

Created a new partition 2 of type 'Linux filesystem' and of size 1 GiB.
```

Add the **Linux Extended Boot** property:

```shell {title="root@usblinux:/ # fdisk"}
Command (m for help): t

Partition number (1-2, default 2): ↵
Partition type or alias (type L to list all): 142
Changed type of partition 'Linux filesystem' to 'Linux Extended Boot'.
```

##### Create the LUKS partition

```shell {title="root@usblinux:/ # fdisk"}
Command (m for help): n

Partition number (3-128, default 3): ↵
First sector (2099200-1953525134, default 2099200): ↵
Last sector, +/-sectors or +/-size{K,M,G,T,P} (2099200-1953525134, default 1953523711): ↵

Created a new partition 3 of type 'Linux filesystem' and of size 930.5 GiB.
```

Add the **Linux Root (x86-64)** property:

```shell {title="root@usblinux:/ # fdisk"}
Command (m for help): t

Partition number (1-3, default 3): ↵
Partition type or alias (type L to list all): 23
Changed type of partition 'Linux filesystem' to 'Linux Root (x86-64)'.
```

##### Update the partition table

```shell {title="root@usblinux:/ # fdisk"}
Command (m for help): w
The partition table has been altered.
Calling ioctl() to re-read partition table.
Syncing disks.
```

#### Basic LUKS setup

##### Set up phassphrase encrypted volume

As an initial setup, it's better to setup a passphrase encrypted volume as a
fallback, as losing a key file would be catastrophic. Since we are setting up
secret keys, we recommend setting up a **long passphrase**.

```shell {title="root@usblinux:/"}
cryptsetup luksFormat --key-size 512 /dev/nvme0n1p3
```

```shell {title="root@usblinux:/ # cryptsetup"}
WARNING!
========
This will overwrite data on /dev/nvme0n1p3 irrevocably.

Are you sure? (Type 'yes' in capital letters):
YES
Enter passphrase for /dev/nvme0n1p3:
```

#### Format the Filesystems

```shell {title="root@usblinux:/"}
mkfs.vfat -F32 /dev/nvme0n1p1
mkfs.ext4 -L boot /dev/nvme0n1p2
cryptsetup luksOpen /dev/nvme0n1p3 cryptroot
mkfs.ext4 -L rootfs /dev/mapper/cryptroot
```

For the next steps, mount as follows:

```shell {title="root@usblinux:/"}
mount /dev/mapper/cryptroot /mnt/rootfs
mount /dev/nvme0n1p2 /mnt/rootfs/boot
mount /dev/nvme0n1p1 /mnt/rootfs/boot/efi

mount --bind /dev /mnt/rootfs/dev
mount --bind /proc /mnt/rootfs/proc
mount --bind /sys /mnt/rootfs/sys
```

### Set up the encrypted LUKS key

#### Set up GPG keys on the Yubikey

Check if the YubiKey is detected:

```shell {title="root@usblinux:/"}
gpg --card-status
```

Then enter the generation menu:

```shell {title="root@usblinux:/"}
gpg --card-edit
```

Generate a new key:

```shell {title="root@usblinux:/ # gpg --card-edit"}
Command> generate

Make off-card backup of encryption key? (Y/n) n

gpg: 3 Admin PIN attempts remaining before card is permanently locked

Admin PIN

PIN

Please specify how long the key should be valid.
         0 = key does not expire
         <n>  = key expires in n days
         <n>w = key expires in n weeks
         <n>m = key expires in n months
         <n>y = key expires in n years
Key is valid for? (0) 0

Key does not expire at all.
Is this correct? (y/N) y

You need a user ID to identify your key; the software constructs the user ID
from the Real Name, Comment and Email Address in this form:
"Heinrich Heine (Der Dichter) <heinrichh@duesseldorf.de>"

Real name: John Doe
Email address: john@example.com
Comment: tester
You selected this USER-ID:
    "John Doe (tester) <john@example.com>"

Change (N)ame, (C)omment, (E)mail or (O)kay/(Q)uit? O
```

At this point, the key is stored on the Yubikey. Exit the generation menu with
`exit`.

#### Create a keyfile and add it to the LUKS volume

Create the key file:

```shell {title="root@usblinux:/"}
dd bs=8388608 count=1 if=/dev/urandom of=/run/luks.key

# 1+0 records in
# 1+0 records out
# 8388608 bytes (8.4 MB, 8.0 MiB) copied, 0.0211386 s, 397 MB/
```

Register the key file with the LUKS volume:

```shell {title="root@usblinux:/"}
cryptsetup luksAddKey /dev/nvme0n1p3 /run/luks.key
```

#### Encrypt the key file with the Yubikey using GPG

You must encrypt the key file with the Yubikey **and with a passphrase**. It
will allow you to recover the key file if the Yubikey is lost.

```shell {title="root@usblinux:/"}
gpg --encrypt \
  --recipient john@example.com \
  --cipher-algo aes256 \
  --armor \
  --symmetric \
  --output /run/luks.key.gpg \
  /run/luks.key
```

Move the encrypted key file to the EFI partition:

```shell {title="root@usblinux:/"}
mv /run/luks.key.gpg /mnt/rootfs/boot/luks.key.gpg
```

Since we worked in the `/run` directory, no cleartext data has been persisted.
We are persisting the encrypted key file in the EFI partition.

During the boot process, in the initramfs, we will load the encrypted key file,
the private key from the Yubikey and decrypt it to use it to unlock the LUKS
volume.

The dracut module responsible for LUKS uses FIFOs to avoid leaking the key file.

### Setting up dracut and the initramfs

To use the correct private key, the GPG agent needs to load a public key. Export
the public key from the Yubikey:

```shell {title="root@usblinux:/"}
gpg --export -a john@example.com > /mnt/rootfs/etc/dracut.conf.d/crypt-public-key.gpg
```

Add the following content to the `/mnt/rootfs/etc/dracut.conf.d/crypt.conf`
file:

```shell {title="root@usblinux:/mnt/rootfs/etc/dracut.conf.d/crypt.conf"}
add_dracutmodules+=" crypt crypt-gpg dm rootfs-block "
# The leading and trailing space are important.
```

Since we use GRUB, to simplify the initramfs generation, we embed the kernel
command line parameters in the initramfs. Get the block device UUID using `blkid`:

```shell {title="root@usblinux:/"}
blkid -s UUID

# /dev/nvme0n1p1: UUID="029D-E51D"
# /dev/nvme0n1p2: UUID="bb2c5038-021d-4189-97aa-82f9c51ebfa1"
# /dev/nvme0n1p3: UUID="9ddc623a-aad5-4739-94ba-62898308437a"
```

Then append the following parameters to the
`/mnt/rootfs/etc/dracut.conf.d/crypt.conf` file:

```shell {title="root@usblinux:/mnt/rootfs/etc/dracut.conf.d/crypt.conf"}
kernel_cmdline+=" root=/dev/mapper/cryptroot rd.luks.uuid=9ddc623a-aad5-4739-94ba-62898308437a rd.luks.name=9ddc623a-aad5-4739-94ba-62898308437a=cryptroot rd.luks.key=/luks.key.gpg:UUID=bb2c5038-021d-4189-97aa-82f9c51ebfa1 "
# The leading and trailing space are important.
```

!!!warning WARNING

Due to an issue with Plymouth, the "eye-candy" splash screen, the GPG prompt
might not show up.

In this case, you need to disable Plymouth:

```shell {title="root@usblinux:/mnt/rootfs/etc/dracut.conf.d/crypt.conf"}
kernel_cmdline+=" rd.plymouth=0 "
```

!!!

At this point, you need to rebuild the initramfs. Mount bind the directories and chroot:

```shell {title="root@usblinux:/"}
chroot /mnt/rootfs
```

```shell {title="root@linux:/"}
# Look for the kernel version before typing the command:
dracut -fv --kver <kernel version>
# Ex: dracut -fv --kver 5.15.0-1-amd64
```

We assume the initramfs has been built in the `/boot` directory. Check the
`/boot/initramfs-<kernel version>.img` file:

```shell {title="root@linux:/"}
ls -lah /boot/

# -rw-------  1 root root  50M Aug 17 21:05 initramfs-5.15.0-1-amd64.img
```

Since the kernel/initramfs hasn't changed, we shouldn't need to update the GRUB
configuration.

### Booting the system

Upon boot, after GRUB, you should be greeted by a prompt:

```shell {title="root@initramfs:/"}
PIN (OpenPGP card <serial number>) [1/3]:
```

If the Yubikey hasn't been inserted, the prompt will look like this:

```shell {title="root@initramfs:/"}
gpg: selecting card failed: No such device
gpg: OpenPGP card not available: No such device
Password (/luks.key.gpg on /dev/nvme0n1p2 for /dev/nvme0n1p3) [1/3]:
```

And if you fail to enter the correct passphrase 3 times in a row, the prompt will look like this:

```shell {title="root@initramfs:/"}
Wrong password
Nothing to read on input.
Enter passphrase for /dev/nvme0n1p3:
```

## Conclusion

Hope you've learned something. LUKS decryption using Yubikey isn't very well
known, so I wrote this article in hope to spread `crypt-gpg` usefulness.

I thank **Alexander Moch** for this wonderful
[article (Using a YubiKey to unlock LUKS and Root on ZFS with native encryption)](https://www.alexmoch.com/blog/gentoo-yubikey-encrypted-root-on-zfs/),
which started my journey into setting up LUKS with Yubikey, until I've learned I
could just use `crypt-gpg` and don't need to do custom Dracut modules.

Anyway, all I can say is that LUKS is cool, and with Yubikey, even cooler.
