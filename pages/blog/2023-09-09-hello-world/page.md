---
title: Hello world!
description: The very first article. About the motivations of developing this blog from scratch with Go and HTMX, and why I want to write articles on this blog.
---

_This article is pretty personal, so if you skip it, I don't mind it._

This article written in Markdown is rendered with an HTTP server in Go, and rendered server-side with HTMX.

Try switching articles to see how fast they load!

## Why ?

Let's say I don't like any of the existing solutions for writing a blog.

My requirements are as follows:

- Server-side rendering, at least for routing.
- Light as hell. I want to run this site on a Raspberry Pi.
- Markdown.

With only these requirements, I haven't found a framework that meets my needs:

- CMS: Nope, nope, nope. This is a technical blog, we **code**.
- [Docusaurus](https://docusaurus.io): SSR, but f- React. It isn't even compatible with the latest React version. Talk about dependencies hell.
- [SvelteKit](https://kit.svelte.dev) or [SveltePress](https://sveltepress.site): SSR, quite "light" (at least it is simple). I have started with this, honestly. But after reading seeing the [motherfuckingwebsite](https://motherfuckingwebsite.com), I had an illumination.
- [Hugo](https://gohugo.io): Static, which is great. But, I want the snappy SSR like Docusaurus or SvelteKit.

## HTMX and Go

[HTMX](https://htmx.org) is a JavaScript library which enables to avoid writing JS by abstracting common HTML document manipulation and, more importantly, by using hypermedia as medium instead of JSON.

**Example**

```html
<script src="https://unpkg.com/htmx.org@1.9.5"></script>
<!-- have a button POST a click via AJAX -->
<button hx-post="/clicked" hx-swap="outerHTML">Click Me</button>
```

On click, the server respond:

```html
<p>HTMX works!</p>
```

By coupling HTMX and Go, we can create a simple HTTP server with state and interactivity!

However, I don't want to talk about the technology behind the blog for this first article. Just remember there is only Go and HTMX behind it. No node/pnpm/js. No TailwindCSS bullsh- too, this is just PicoCSS with hard-coded `style`.

## Okay, so what is this blog ?

This is a **personal** and **technical** blog.

By personal, I mean: I am the god of this blog. I can write about anything. This blog reflects me.

By technical, I mean: "Behind the scene of a solution". This blog will document and discuss **conceptualization**, **development**, **implementation** and **deployment** of solutions. It is NOT about writing a tutorial, there won't be any "steps-by-steps". This is about the "what", "why", "how", "what's good or bad" and "what's next" of a solution. RTFM if you need to. Or just learn to reverse engineer.

TL;DR: **I** document my **findings**.

## So, what's next ?

I am planning to write about:

- The conceptualization, development and implementation of this blog.
- About dracut live images and the deployment of stateless images to server.
- Infra., software, ... you name it.
