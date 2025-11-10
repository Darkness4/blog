---
title: Modal dialog with Hyperscript and PicoCSS
description: Small article about a deadly combination.
tags: [dialog, modal, hyperscript, picocss, css, js, html, htmx]
---

## Table of contents

<div class="toc">

{{% $.TOC %}}

</div>

## TL;DR

```html
<button _="on click call #register-dialog.showModal()">Open dialog</button>
<dialog id="register-dialog">
  <article
    _="on click[#register-dialog.open and event.target.matches('dialog')] from elsewhere call #register-dialog.close()"
  >
    <header>
      <button aria-label="Close" rel="prev"></button>
      <p>
        <strong>🗓️ Thank You for Registering!</strong>
      </p>
    </header>
    <p>
      We're excited to have you join us for our upcoming event. Please arrive at
      the museum on time to check in and get started.
    </p>
    <ul>
      <li>Date: Saturday, April 15</li>
      <li>Time: 10:00am - 12:00pm</li>
    </ul>
  </article>
</dialog>
```

<button _="on click call #register-dialog.showModal()">Open dialog</button>

<dialog id="register-dialog">
  <article
    _="on click[#register-dialog.open and event.target.matches('dialog')] from elsewhere call #register-dialog.close()"
  >
    <header>
      <button aria-label="Close" rel="prev"></button>
      <p>
        <strong>🗓️ Thank You for Registering!</strong>
      </p>
    </header>
    <p>
      We're excited to have you join us for our upcoming event. Please arrive at
      the museum on time to check in and get started.
    </p>
    <ul>
      <li>Date: Saturday, April 15</li>
      <li>Time: 10:00am - 12:00pm</li>
    </ul>
  </article>
</dialog>

## How-to create a dialog modal with PicoCSS and vanilla JS

PicoCSS offers a dialog modal directly in the CSS, without any additional JS:

<style>
dialog.example {
  z-index:inherit;
  position:relative;
  min-height:inherit;
  margin:0 calc(var(--pico-spacing) * -1) var(--pico-block-spacing-vertical) calc(var(--pico-spacing) * -1);
  inset:inherit
}
@media (min-width:576px) {
  dialog.example {
    margin:0;
    margin-bottom:var(--pico-block-spacing-vertical)
  }
}
dialog.example>article {
  animation:none
}
</style>

<dialog class="example" open=""><article><header><button aria-label="Close" rel="prev"></button><p><strong>🗓️ Thank You for Registering!</strong></p></header><p>We’re excited to have you join us for our upcoming event. Please arrive at the museum on time to check in and get started.</p><ul><li>Date: Saturday, April 15</li><li>Time: 10:00am - 12:00pm</li></ul></article></dialog>

```html {title="dialog.html"}
<dialog open>
  <article>
    <header>
      <button aria-label="Close" rel="prev"></button>
      <p>
        <strong>🗓️ Thank You for Registering!</strong>
      </p>
    </header>
    <p>
      We're excited to have you join us for our upcoming event. Please arrive at
      the museum on time to check in and get started.
    </p>
    <ul>
      <li>Date: Saturday, April 15</li>
      <li>Time: 10:00am - 12:00pm</li>
    </ul>
  </article>
</dialog>
```

_Credits to [PicoCSS docs](https://picocss.com/docs/modal)_

However, there is a small issue: the example is not interactive!

Thankfully, PicoCSS also offers an example to open and close using JS. However, the give is quite long:

```js {title="dialog.js"}
/*
 * Modal
 *
 * Pico.css - https://picocss.com
 * Copyright 2019-2024 - Licensed under MIT
 */

// Config
const isOpenClass = 'modal-is-open';
const openingClass = 'modal-is-opening';
const closingClass = 'modal-is-closing';
const scrollbarWidthCssVar = '--pico-scrollbar-width';
const animationDuration = 400; // ms
let visibleModal = null;

// Toggle modal
const toggleModal = (event) => {
  event.preventDefault();
  const modal = document.getElementById(event.currentTarget.dataset.target);
  if (!modal) return;
  modal && (modal.open ? closeModal(modal) : openModal(modal));
};

// Open modal
const openModal = (modal) => {
  const { documentElement: html } = document;
  const scrollbarWidth = getScrollbarWidth();
  if (scrollbarWidth) {
    html.style.setProperty(scrollbarWidthCssVar, `${scrollbarWidth}px`);
  }
  html.classList.add(isOpenClass, openingClass);
  setTimeout(() => {
    visibleModal = modal;
    html.classList.remove(openingClass);
  }, animationDuration);
  modal.showModal();
};

// Close modal
const closeModal = (modal) => {
  visibleModal = null;
  const { documentElement: html } = document;
  html.classList.add(closingClass);
  setTimeout(() => {
    html.classList.remove(closingClass, isOpenClass);
    html.style.removeProperty(scrollbarWidthCssVar);
    modal.close();
  }, animationDuration);
};

// Close with a click outside
document.addEventListener('click', (event) => {
  if (visibleModal === null) return;
  const modalContent = visibleModal.querySelector('article');
  const isClickInside = modalContent.contains(event.target);
  !isClickInside && closeModal(visibleModal);
});

// Close with Esc key
document.addEventListener('keydown', (event) => {
  if (event.key === 'Escape' && visibleModal) {
    closeModal(visibleModal);
  }
});

// Get scrollbar width
const getScrollbarWidth = () => {
  const scrollbarWidth =
    window.innerWidth - document.documentElement.clientWidth;
  return scrollbarWidth;
};

// Is scrollbar visible
const isScrollbarVisible = () => {
  return document.body.scrollHeight > screen.height;
};
```

Given the size of the script, I would have extracted this code in a JS file, but I wanted to keep the logic close to the HTML element triggering it.

This is where Hyperscript comes in.

## Hyperscript, the missing piece for homemade SSR

### Difficulties with homemade SSR with HTMX

When doing a homemade SSR server, the server becomes "fat". It is responsible to render the final HTML, instead of an SPA, where the client renders the HTML.

The motivation for SSR comes simply by a wish to go back to the good old days of the web: the server serves a proper page and the client simply renders the page, without any complexity. This has many advantages, with one being simply the performance. Servers are known to be strong, whereas clients are not, especially modern browsers.

When implementing SSR, most of the time, people will prefer full-stack solutions like [NextJS](https://nextjs.org) or [SvelteKit](https://svelte.dev/docs/kit/introduction). But the biggest issue with these frameworks is that they are heavy, and the programming language Typescript doesn't offer enough safety.

This is why, some people will choose to use [HTMX](https://htmx.org/), a solution to do SSR by serving data using HTML. For example, instead of serving:

```json
{
  "results": [
    {
      "title": "My blog article"
    },
    {
      "title": "My second blog article"
    }
  ]
}
```

HTMX would serve a properly rendered HTML fragment:

```html
<section>
  <ul>
    <li>My blog article</li>
    <li>My second blog article</li>
  </ul>
</section>
```

This offers these advantages:

- No JS required, which improves performance
- Locality of behavior
- Agnostic to the programming language
- The server is fat, and the client is [thin](https://en.wikipedia.org/wiki/Thin_client)

But has these drawbacks:

- No front-end framework
  - And if there are, it requires some compilation (React, Vue, etc.)
  - And if they don't require compilation, their performance can be [doubtful](https://krausest.github.io/js-framework-benchmark/current.html) (see alpine, as it is often the most recommended framework to use with HTMX)
- Styling is complicated

### Hyperscript as a simple solution

[Hyperscript](https://hyperscript.org/) is a small library that allows to write simple script alongside the HTML, without disrupting the lisibility of the code.

To install it, simple put in the `index.html`:

```html
<script src="https://unpkg.com/hyperscript.org@0.9.14"></script>
```

Then, you can write:

```html
<button _="on click send hello to <form />">Send</button>
```

...I know. This is a new syntax to learn, and honestly, I still have a hard time with it.

Other drawbacks are:

- No compilation, i.e, no compile-time checks.
- The syntax itself isn't clear due to its "almost-English" nature.
- CSP issues (see [Security - Hyperscript](https://hyperscript.org/docs/#security))

### Dialog Modal with Hyperscript and PicoCSS

Even though Hyperscript has some issues, it is still more clearer than the script above. Here's the implementation of a dialog using Hyperscript:

```html
<button _="on click call #register-dialog.showModal()">Open dialog</button>
<dialog id="register-dialog">
  <article
    _="on click[#register-dialog.open and event.target.matches('dialog')] from elsewhere call #register-dialog.close()"
  >
    <header>
      <button aria-label="Close" rel="prev"></button>
      <p>
        <strong>🗓️ Thank You for Registering!</strong>
      </p>
    </header>
    <p>
      We're excited to have you join us for our upcoming event. Please arrive at
      the museum on time to check in and get started.
    </p>
    <ul>
      <li>Date: Saturday, April 15</li>
      <li>Time: 10:00am - 12:00pm</li>
    </ul>
  </article>
</dialog>
```

<button _="on click call #register-dialog.showModal()">Open dialog</button>

<dialog id="register-dialog">
  <article
    _="on click[#register-dialog.open and event.target.matches('dialog')] from elsewhere call #register-dialog.close()"
  >
    <header>
      <button aria-label="Close" rel="prev" _="on click call #register-dialog.close()"></button>
      <p>
        <strong>🗓️ Thank You for Registering!</strong>
      </p>
    </header>
    <p>
      We're excited to have you join us for our upcoming event. Please arrive at
      the museum on time to check in and get started.
    </p>
    <ul>
      <li>Date: Saturday, April 15</li>
      <li>Time: 10:00am - 12:00pm</li>
    </ul>
  </article>
</dialog>

This is what I'm talking about **locality of behavior**. Thanks to the syntax of Hyperscript, there is no need to search for an element using verbose functions such as `document.getElementById`, and the code is more clear and readable. The only thing I need to example is:

```html
<article
  _="on click[#register-dialog.open and event.target.matches('dialog')] from elsewhere call #register-dialog.close()"
></article>
```

which means:

- I attach an event listener and listen for events only when the dialog is open and the event target is a dialog (the user has clicked on the hitbox of the dialog, this also includes the background).
- If there is an event that is not hitting the `<article>` (`elsewhere` of the `<article>`), I close the dialog.

Pretty cool, huh?

## Conclusion

This is just a small article about this deadly combination. I know you may simply not like Hyperscript, and would criticize it as an alternative to jQuery or vanilla JS, but I think it's better than them since I can write the code alongside the HTML, and it's more readable.

Also, if the syntax of hyperscript is too difficult, hyperscript is also able to call JS code, which makes it a good extension of HTML.
