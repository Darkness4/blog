---
title: A first try on Zig and C interop
description: Trying Zig with C libraries for the first time.
tags: [zig, c, ffmpeg, av1, ffi]
---

## Table of contents

<div class="toc">

\\{\\{ $.TOC }}

</div>

## Introduction

!!!warning WARNING

Zig is still in development, and the language API is not stable. The code in this article may not work in future versions of Zig.

The version of Zig used in this article is `0.13.0-dev.380+b32aa99b8`.

!!!

Lately, I've been programming in Zig for the Advent of Code 2023. What I've learned about the language is that it works perfectly well for low-level programming, but lacked in some areas (my biggest gripe is the lack of HTTP2).

However, Zig claims to be able to interop with C, which is something that I've been wanting to try for a long time. As you know, I used to use CGO in Go, which permits me to use complex C libraries with a high-level layer in Go.

In this article, I will try to demonstrate my experience with Zig and C interop. I will use a simple example: a C library that transcode a video into AV1 format, and a Zig program that uses this library.

## The concept

Transcoding a video seems quite a complex task, but thanks to the `libavcodec` and `libavformat` C libraries, it is quite "easy" to do (or at least, to understand).

Transcoding a video follows these steps:

1. Open the input video file.
2. Demux the input video file: read packets from the input video file.
3. Decode the packets: decode the packets into frames by passing them to the decoder.
4. Encode the frames: encode the frames into packets by passing them to the encoder.
5. Mux the packets: write the packets to the output video file.

```d2
direction: right

Open -> Demux: file
Demux -> Decode: packet
Decode -> Encode: frame
Encode -> Mux: packet
```

The aim of this article is to demonstrate how to use a C library without the need to make a Zig wrapper around it. This is a common practice in Go, where you can use CGO to call C functions directly.

## Zig and tricks

I'm taking account that people reading this article are not familiar with Zig, so I will explain some concepts that I've learned while programming in Zig.

### Memory allocators

Zig does not have a garbage collector, so you have to manage memory "yourself". By yourself, I mean that you have to allocate and deallocate memory manually.

Compared to C, you are free to choose the type of allocator you want to use. The standard library provides a `std.mem.Allocator` interface that you can implement to create your own allocator.

In this article, I will use two allocators:

- `std.heap.GeneralPurposeAllocator`: a simple allocator. We could have used the `std.heap.page_allocator` or the `std.heap.c_allocator`, but for the sake of being simple, I will use the `GeneralPurposeAllocator`.
- `std.heap.ArenaAllocator`: an allocator that groups allocations in arenas. This is useful when you want to deallocate a group of allocations at once, like the arguments of the program.

In Zig, the `GeneralPurposeAllocator` can be used to detect memory leaks:

```zig
var gpa = std.heap.GeneralPurposeAllocator(.{}){};
const gpa_allocator = gpa.allocator();

pub fn main() !void {
    // Detect memory leaks
    defer std.debug.assert(gpa.deinit() == .ok);
    // Your code here
}
```

Oh yeah, Zig has a `defer` statement, which is quite similar to Go's `defer`, but it is scoped to the curly braces, compared to Go's `defer`, which is scoped to the function. And the `.ok` is an enum value.

About the `gpa` declaration, I'm calling a function with an "empty" struct as an argument, which returns a type. Do note I'm quoting "empty", because Zig has default values for structs, which is quite useful.

To instanciate the type `std.heap.GeneralPurposeAllocator(.{})`, we add the curly braces after the type name:

```zig
const my_struct: type = struct {
    a: u8,
    b: u8,
};

var object = my_struct{};
var object: my_struct = .{}; // Equivalent
```

### Easy memory management with `defer`

The best feature of Zig is the `defer` statement because it completes the "flow" of the function. Similar to Go, `defer` can be used to clean up resources at the end of the function:

```zig
func do() !void {
    var array = try allocator.alloc(i64, array_size);
    defer allocator.free(array);

    // Your code here
}
```

Compared to C:

```c
int do() {
    int ret = 0;
    int *array = malloc(array_size * sizeof(int));
    if (array == NULL) {
        ret = 1;
        goto end;
    }

    // Your code here

end:
    if (array != NULL) free(array);
    return ret;
}
```

But, one issue that I've had is that since `defer` is scoped to the curly braces, it's almost unusable in `if` statements:

```zig
var data: ?[]u8 = null;
if (first_file) {
    try allocator.alloc(u8, data_size);
    defer allocator.free(data);
} // Is cleared here
```

Zig has `errdefer`, which is "almost" what I what, but only triggers when an error occurs:

```zig
var data: ?[]u8 = null;
if (first_file) {
    try allocator.alloc(u8, data_size);
    errdefer allocator.free(data);
}

return; // errdefer is not triggered: no error returned.
```

It would be nice to have an actual equivalent of Go's `defer` in Zig.

### Error handling

Zig has an error handling similar to Go, but slightly more primitive. To compare:

**C:**

```c
int my_function() {
    if (error_condition) {
        return -1;
    }
    return 0;
}
```

**Go:**

```go
func myFunction() error {
    if errorCondition {
        return errors.New("my error")
    }
    return nil
}
```

**Zig:**

```zig
fn my_function() !void {
    if (error_condition) {
        return error.MyError;
    }
}
```

In Zig, errors are not implementations of an error interface like in Go, but are enums. And errors can have subsets and supersets, which is somewhat confusing at first, but quite powerful:

```zig
const std = @import("std");

const FileOpenError = error{
    AccessDenied,
    OutOfMemory,
    FileNotFound,
};

const AllocationError = error{
    OutOfMemory,
};

test "coerce subset to superset" {
    const err = subset_to_superset(AllocationError.OutOfMemory);
    try std.testing.expect(err == FileOpenError.OutOfMemory);
}

fn subset_to_superset(err: AllocationError) FileOpenError {
    return err;
}
```

`error` is the base superset.

As you can see, there is some flexibility in error handling in Zig, but has one downside: error does not have any value (no message, no custom data). However, Zig errors remember the stack trace, which is quite useful for debugging and could help avoid the need to pass custom data in the error.

### Zig's quick error handling

In Go, we often have this pattern:

```go
func myFunction() error {
    if err := someFunction(); err != nil {
        return err
    }
    return nil
}
```

In Zig, we can use the `try` keyword to return an error immediately:

```zig
fn my_function() !void {
    try some_function();
}
```

This helps to streamline the error handling in Zig:

```zig
fn complex_function() !void {
    try handle_data(try fetch_data());
}
```

### C interop

Zig has its own C compiler and own "C-translator". To include a C library in Zig, you have to use the `@cImport` directive:

```zig
const c = @cImport({
    @cInclude("libavcodec/avcodec.h");
    @cInclude("libavformat/avformat.h");
    @cInclude("libavutil/avutil.h");
});

pub fn main() !void {
    c.call_some_c_function();
}
```

Which is quite similar to Go:

```go
/*
#cgo pkg-config: libavformat libavcodec libavutil
#include <libavcodec/avcodec.h>
#include <libavformat/avformat.h>
#include <libavutil/avutil.h>
*/
import "C"

func main() {
    C.call_some_c_function()
}
```

About `pkg-config`, Zig has its own build system. You can create a `build.zig` and pass the C libraries you want to link to:

```zig
pub fn build(b: *std.Build) void {
    const target = b.standardTargetOptions(.{});
    const optimize = b.standardOptimizeOption(.{});

    const exe = b.addExecutable(.{
        .name = "av1-transcoder",
        .root_source_file = b.path("src/main.zig"),
        .target = target,
        .optimize = optimize,
        .link_libc = true,
    });

    exe.addIncludePath(.{
        .src_path = .{ .owner = b, .sub_path = "src" },
    });
    exe.linkSystemLibrary2("libavcodec", .{ .preferred_link_mode = .static });
    exe.linkSystemLibrary2("libavutil", .{ .preferred_link_mode = .static });
    exe.linkSystemLibrary2("libavformat", .{ .preferred_link_mode = .static });
    exe.linkSystemLibrary2("swresample", .{ .preferred_link_mode = .static });
    exe.linkSystemLibrary2("SvtAv1Dec", .{ .preferred_link_mode = .static });
    exe.linkSystemLibrary2("SvtAv1Enc", .{ .preferred_link_mode = .static });

    b.installArtifact(exe);
}
```

Your IDE/LSP won't be able to detect the C symbols at first, but after compiling the project, it will be able to detect them.

But one difference is certain: Zig has less "bridges" between Zig and C, which makes the code more readable than Go's:

```zig
// In the .zig-cache directory, there is the translation of the C library to Zig
pub extern fn avformat_open_input(ps: [*c][*c]AVFormatContext, url: [*c]const u8, fmt: [*c]const AVInputFormat, options: [*c]?*AVDictionary) c_int;

// Usage
fn test_avformat_open_input(input [*:0]const u8) c_int {
    var ifmt_ctx: ?[*]c.AVFormatContext = null;
    var ret = c.avformat_open_input(&ifmt_ctx, input_file, null, null);
}
```

```go
// This is stored in the $HOME/.cache/go-build directory. The translation is not readable.
//go:cgo_unsafe_args
func _Cfunc_avformat_open_input(p0 **_Ctype_struct_AVFormatContext, p1 *_Ctype_char, p2 *_Ctype_struct_AVInputFormat, p3 **_Ctype_struct_AVDictionary) (r1 _Ctype_int) {
	_cgo_runtime_cgocall(_cgo_fa42a779fc4c_Cfunc_avformat_open_input, uintptr(unsafe.Pointer(&p0)))
	if _Cgo_always_false {
		_Cgo_use(p0)
		_Cgo_use(p1)
		_Cgo_use(p2)
		_Cgo_use(p3)
	}
	return
}

// Usage
func testAVFormatOpenInput(input string) C.int {
	var ifmt *C.AVFormatContext = nil
	return C.avformat_open_input(&ifmt, C.CString(input), nil, nil)
}
```

Oh wait, again! Here, Zig has multiple features that enhance the safety of the code:

- `[*:0]const u8` indicates a slice (`[]const u8`) that is sentinel-terminated (`[:0]const u8`) and is the pointer (`[*]const u8`, which makes `[*:0]const u8`). To summarize, this is a C string. Zig strings does not need to be manipulated with a pointer.
- `?[*]c.AVFormatContext` is a nullable pointer to a `c.AVFormatContext` struct. Zig has basic null safety. Since C does not have any null safety, you may see instead `[*c]c.AVFormatContext` which is a C pointer to a `c.AVFormatContext` struct and can be null.

Both languages suffer from one major issue: the comments are not passed to the translation, which means that deprecation notices or warnings are not passed to the Zig/Go code.

Overall, Zig has a slightly better C interop than Go.

### Limitations of the C interop

Zig has some limitations when it comes to C interop:

- Some macros are translated to Zig, but not all of them. You may have to write the Zig equivalent of the macro:

  ```zig
  pub const av_err2str = @compileError("unable to translate C expr: expected ')' instead got '['");
  ```

  The worst issue that I had is due to the strictness of Zig's type system. Some macros do not translate well between `u64`, `c_int`, `usize`... This is quite a pain, because FFmpeg (`libavutil`) uses a macro to define errors at compile time.

- `const`-hell. Zig is able to handle `const` at pointer ("pointer's value is immutable") and struct level ("struct is immutable"). However, Zig is quite picky when passing a `const` pointer to a C function (developper's fault):

  ```zig
  c.av_guess_frame_rate(@constCast(ifmt_ctx), @constCast(in_stream), null);
  ```

  Technically, `av_guess_frame_rate` accepts a const pointer because it does not modify the pointer.

### Visibility

Zig visibility is scoped to the file, which is more similar to Python or C than Go. This forces you to mostly develop in a single file, which is quite a pain when you have a large project.

To export a function, you have to use the `pub` keyword:

```zig
pub fn my_function() void {
    // Your code here
}
```

To import a function, you have, well..., to import it:

```zig
const my_module = @import("my_module.zig");

pub fn main() void {
    my_module.my_function();
}
```

Somewhat, I prefer Go's visibility, which is scoped to the package and allows you to separate responsibilities easily.

I mean, just look at the `std` package in Zig. It's quite a mess. (Example: [general_purpose_allocator.zig](https://github.com/ziglang/zig/blob/master/lib/std/heap/general_purpose_allocator.zig)). Tests are also in the same file, which, I guess, it's fine.

The reason why I think that Zig is more messy than C and Python is because C's header indicates explicitly what is exported and what is not. And, Python hasn't really a visibility system, but it's quite easy to understand what is exported simply by looking at the variables and functions names.

Overall, Zig's visilibilty tries to be the best of both worlds: everything private in one file like in C, but without the hassle of header files, with the sacrifice of having one messy file. I hope there will be some styling guidelines in the future.

### Syntax

Zig's syntax has a lot of quality of life improvements compared to C and Go. I won't go too much into detail, but here are some examples:

- Dereferencing and null-safety chaining:

  ```zig
  my_ptr.?.*.another_ptr.?.*.nullable_struct.?.ptr.*
  ```

  Go equivalent:

  ```go
  if myPtr != nil && myPtr.anotherPtr != nil && myPtr.anotherPtr.nullableStruct != nil {
      myPtr.anotherPtr.nullableStruct.ptr // Implicit dereference
  }
  ```

  C equivalent:

  ```c
  if (my_ptr != NULL && my_ptr->another_ptr != NULL && my_ptr->another_ptr->nullable_struct != NULL) {
      *(my_ptr->another_ptr->nullable_struct->ptr);
  }
  ```

- For loops can uses ranges and zip (simultanous iteration)
- Etc...

## Developing the AV1 transcoder

### Developing the Remuxer

I had some issue with developing with `libav*` libraries due to the lack of resources. But after a while, I was able to do it in Zig without any wrapper.

To transcode a video, you must first think about remuxing the video: Read packets and write packets into a new container. Remuxing is relatively easy to do, since it's all about concatenating packets.

The steps are:

1. Open the input video file.
2. Open the output video file.
3. Demux the input video file: read packets (`av_read_frame`) from the input video file in a while loop.
4. Process the packet: rescale the timestamps of the packet, and fix eventual discontinuities.
5. Mux the packets: write the packets to the output video file (`av_interleaved_write_frame`)
6. Close the input and output video files. Clean up everything.

The [example](https://www.ffmpeg.org/doxygen/trunk/remuxing_8c-example.html) given by the FFmpeg documentation is quite accurate (minus the mpegts discontinuities fix)

I won't go too much into details, but here are the sexy stuff which was improved by Zig:

- No more `goto`, no more dangling `ret`. Zig `defer` is powerful enough to handle most the memory management:

  ```zig
  // Using Zig's allocator
  var stream_mapping = try allocator.alloc(i64, stream_mapping_size);
  defer allocator.free(stream_mapping);

  // Using C's (libav) allocator
  var enc_ctx = c.avcodec_alloc_context3(enc);
  if (enc_ctx == null) {
    // Error handling
  }
  defer c.avcodec_free_context(&enc_ctx);
  ```

- Slightly object-oriented, to bind multiple lifecycles into one:

  ```zig
  const Context = struct {
    stream_mapping: []i64,
    dts_offset: []i64,
    // ...
    allocator: std.mem.Allocator,

    fn init(allocator: Allocator, size: usize) !Context {
      return .{
        .stream_mapping = try allocator.alloc(i64, size),
        .dts_offset = try allocator.alloc(i64, size),
        // ...
      };
    }

    fn deinit(self: *Context) void {
      self.allocator.free(self.stream_mapping);
      self.allocator.free(self.dts_offset);
      // ...
    }
  };
  ```

Now, let's develop the transcoding side of the program.

### Developing the Transcoder

Transcoding a video adds 6 steps:

- Initializing the decoder.
- Initializing the encoder.
- Add a `while` loop to decode frames.
- Add a `while` loop to encode frames.
- Flush the encoder.
- Flush the decoder.

The steps include also fixing the timestamps and the frame rate.

You can pretty much use the [example](https://ffmpeg.org/doxygen/trunk/transcode_8c-example.html) given by the FFmpeg documentation (minus the filters).

The code looks like this:

```zig
// .. First loop, in the remuxer
while (true) {
    ret = c.av_read_frame(ifmt_ctx, pkt);
    if (ret < 0) {
        // No more packets
        break;
    }
    defer c.av_packet_unref(pkt);

    const in_stream_index = @as(usize, @intCast(pkt.stream_index));

    // Packet is blacklisted
    if (in_stream_index >= stream_mapping_size or stream_mapping[in_stream_index] < 0) {
        continue;
    }
    const out_stream_index = @as(usize, @intCast(stream_mapping[in_stream_index]));
    pkt.stream_index = @as(c_int, @intCast(out_stream_index));

    const stream_ctx = stream_ctxs[out_stream_index];

    try stream_ctx.fix_discontinuity_ts(pkt);

    // Input to decoder timebase
    try stream_ctx.transcode_write_frame(pkt);
} // while packets.

// ...

// Second loop, in the decoder (stream_ctx)
fn transcode_write_frame(self: StreamContext, pkt: ?*c.AVPacket) !void {
    // Send packet to decoder
    try self.decoder.send_packet(pkt);

    while (true) {
        // Fetch decoded frame from decoded packet
        const frame = self.decoder.receive_frame() catch |e| switch (e) {
            AVError.EAGAIN => return,
            AVError.EOF => return,
            else => return e,
        };
        defer c.av_frame_unref(frame);

        frame.*.pts = frame.*.best_effort_timestamp;

        if (frame.*.pts != c.AV_NOPTS_VALUE) {
            frame.*.pts = c.av_rescale_q(frame.*.pts, self.decoder.dec_ctx.?.*.pkt_timebase, self.encoder.enc_ctx.?.*.time_base);
        }

        try self.encode_write_frame(frame);
    }
}

// Third loop, in the encoder (also in stream_ctx)
fn encode_write_frame(self: StreamContext, dec_frame: ?*c.AVFrame) !void {
    self.encoder.unref_pkt();

    try self.encoder.send_frame(dec_frame);

    while (true) {
        // Read encoded data from the encoder.
        var pkt = self.encoder.receive_packet() catch |e| switch (e) {
            AVError.EAGAIN => return,
            AVError.EOF => return,
            else => return e,
        };

        // Remux the packet
        pkt.stream_index = @as(c_int, @intCast(self.stream_index));

        // Encoder to output timebase
        c.av_packet_rescale_ts(pkt, self.encoder.enc_ctx.?.*.time_base, self.out_stream.*.time_base);

        try self.fix_monotonic_ts(pkt);

        // Write packet
        const ret = c.av_interleaved_write_frame(self.ofmt_ctx, pkt);
        if (ret < 0) {
            err.print("av_interleaved_write_frame", ret);
            return ret_to_error(ret);
        }
    }
}
```

Or in simple words:

1. Read a packet from the input video file by calling `av_read_frame`.
2. Send packet to the decoder by calling `avcodec_send_packet`.
3. Receive a frame from the decoder by calling `avcodec_receive_frame`.
4. Send frame to the encoder by calling `avcodec_send_frame`.
5. Receive a packet from the encoder by calling `avcodec_receive_packet`.
6. Write the packet to the output video file by calling `av_interleaved_write_frame`.

And that's it! You have a video transcoder in Zig. (You'll also need to fix the timestamps and discontinuities, but that's another story).

## Last part: the build system

The `build.zig` file that I've given earlier is quite enough to build the project. Zig automatically statically links the libraries, making the executable portable.

Oh, and you'll need to fork SVT-AV1 to enable the `static` flag:

```shell
# svt-av1-9999.ebuild, using Gentoo's Portage to build the package

multilib_src_configure() {
  append-ldflags -Wl,-z,noexecstack

  local mycmakeargs=(
    -DBUILD_TESTING=OFF
    -DCMAKE_OUTPUT_DIRECTORY="${BUILD_DIR}"
    -DBUILD_SHARED_LIBS="$(usex static-libs OFF ON)" # Enable static libraries based on the USE flag "static-libs"
  )

  [[ ${ABI} != amd64 ]] && mycmakeargs+=(-DCOMPILE_C_ONLY=ON)

  cmake_src_configure
}
```

However, I have one MAJOR issue: the libc is not statically linked, which means I'm unable to create a distroless Docker image. It would have been nice to have a `prefered_link_mode` for the libc.

When disabling `link_libc`, the executable does not compile since symbols are missing (even with `linkSystemLibrary2("c", .{ .preferred_link_mode = .static })`). Normally, Zig automatically links the libc statically, but it seems that this isn't the case here.

To force the static linking, you can enable `.linkage = .static` in the `addExecutable` function. And instead of using `linkSystemLibrary2`, you can use `addObjectFile`. The issue with this technique is that you have to do everything manually, and you cannot use `pkg-config` to find the missing includes and libraries:

```zig
pub fn build(b: *std.Build) void {
    const target = b.standardTargetOptions(.{});
    const optimize = b.standardOptimizeOption(.{});

    const exe = b.addExecutable(.{
        .name = "av1-transcoder",
        .root_source_file = b.path("src/main.zig"),
        .target = target,
        .optimize = optimize,
        .linkage = .static,
        .link_libc = true,
    });

    exe.addIncludePath(.{
        .src_path = .{ .owner = b, .sub_path = "src" },
    });
    exe.addIncludePath(.{
        .src_path = .{ .owner = b, .sub_path = "/usr/include" },
    });
    exe.addObjectFile(.{ .src_path = .{
        .owner = b,
        .sub_path = "/usr/lib/libavcodec.a",
    } });
    exe.addObjectFile(.{ .src_path = .{
        .owner = b,
        .sub_path = "/usr/lib/libavutil.a",
    } });
    exe.addObjectFile(.{ .src_path = .{
        .owner = b,
        .sub_path = "/usr/lib/libavformat.a",
    } });
    exe.addObjectFile(.{ .src_path = .{
        .owner = b,
        .sub_path = "/usr/lib/libswresample.a",
    } });
    exe.addObjectFile(.{ .src_path = .{
        .owner = b,
        .sub_path = "/usr/lib/libSvtAv1Dec.a",
    } });
    exe.addObjectFile(.{ .src_path = .{
        .owner = b,
        .sub_path = "/usr/lib/libSvtAv1Enc.a",
    } });

```

At this point, the generated artifact is a static executable:

```shell
$ ldd ./zig-out/bin/av1-transcoder
ldd: ./zig-out/bin/av1-transcoder: Not a valid dynamic program
```

Yay!

Small issue: Because the paths are hardcoded, it will be quite difficult to cross-compile the project at the moment.

## Conclusion

Zig C interop is almost impeccable, and at least way better than Go's. Symbols are directly translated to Zig, and the memory management is quite easy to handle. The Zig API plugs well with the C API, making C code slightly safer.

However, the build system, while quite powerful, lacks of flexibility around the libc linking and `pkg-config` support. Perhaps sticking to a Makefile would be better for now.

Overall, Zig presents some potential for low-level programming, or at least for dynamic libraries development. The features and syntax complement very well with C. But, I would still recommend C++ if you want to develop stable and production-ready software. While C++ is complex due to the richness of the language, C++ offers its kind of safety (smart pointers) which can help you avoid memory leaks and dangling pointers.

Lastly, because Zig is still in development, Zig lacks of high-level libraries and frameworks, which limits the use of Zig in production (no gRPC).

So, to conclude, I will use Zig for competitive programming and for basic API like my AV1 transcoder bot. But for production, I will stick to Go.

## References

- [Git Repository](https://github.com/Darkness4/av1-transcoder-bot)
- [FFmpeg documentation](https://ffmpeg.org/doxygen/trunk/index.html)
- [Zig documentation](https://ziglang.org/documentation/master/)
