---
title: Learn software architecture, paradigms and patterns... even the wrong ones.
description: Have you ever wondered whether learning the wrong software architecture is really "wrong"? Personally, I've always asked myself this question, and more often than not I've found my answer on the job.
tags: [software architecture, paradigms, patterns, programming]
---

## Table of contents

<div class="toc">

\\{\\{ $.TOC }}

</div>

<hr>

## Introduction

For once, I want to write about my experience about learning software architecture, paradigms, and patterns. Recently, I've been absorbed by influencers on Youtube (ThePrimeAgen, Jonathan Blow, etc.). It's not like I take any of these influencers as true value, and, in fact, this is why I'm writing this article. My true value has always been my job: "If it works for my team, it may work with you", which is why I'm writing this article.

Most of these Youtubers (at least ThePrimeAgen) don't really like OOP, clean architecture, and things like that. You will often hear "Clean Code is slow", or "OOP adds too much abstraction". But I wanted to convince you to at least try to understand why these concepts suck or work. You should try to understand **why these "concepts" were developed** and **why they failed**, especially if your work in software engineering.

## Software architecture: "frameworks" for beginners

When we talk about software architecture, we are talking about **how to structure your data/state and processes**. A software architecture is a set of practices designed to organize the elements of your software and their relationships in an "orderly" fashion.

This is very similar to a standard, except **it is not**. A software architecture aims to only solve one issue: **reduce complexity of a system**.

Reducing the complexity of a system gives immediate benefits:

- **Maintainability**: Because it is easier to understand, it is easier to debug and troubleshoot.
- **Scalability**: Because its interface are understood, the software can be scaled up through modularity, or simply through development.
- **Robustness**: A simple system mean fewer potential point of failures.

Let's talk a little about [Clean Architecture from Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html), which is the first software architecture I've learned.

Clean Architecture is a "layered" architecture where:

- The entities are the most independent objects. These objects define business rules. It's an object or a method. Example: `Person`
- The use cases are also the most independent objects but depends on at least to entities. Example: Set the `Person`'s name to uppercase, i.e. `PersonNameToUppercase`.
- The controllers, gateways, presenters... i.e. the "data transformers/services" returns data to use cases. Example: `FetchPerson`.
- And the rest (controllers, UI, etc.) sends data to the core layers. Example: a JSON representing a person.

Since it is my first time, I used it everywhere in my personal projects and it worked goddamn well, especially with Dart, Kotlin and Java. It solves these issues:

- Circular dependencies.
- Mixed responsibilities.
- Heavy coupling, which makes the system difficult to extend.

**Notice that these issues happens after a project is well-developed. A small project will never have these issues. These issues are not told by influencers, because they already learned about how to avoid these through best practices, software architecture, patterns, or simply by implicit shared knowledge.**

Many of my friends at engineering schools had at least one dependency cycle, and are always asking this question: **"Where and how can I instantiate my objects?**".

Learning a software architecture makes you learn about **Dependencies Injection**, **Object life cycle**, **Data/Domain/Presentation** separation (which is bullsh-, but I'm going to talk about that later), **Single Responsibility Principle**, **Interfaces** (yes, first timers don't know why we use interfaces) and maybe some **Test Driven Development** (which is also bullsh-) since these are the issues targeted by software architecture.

Simple example of circular dependency:

```kotlin
// Person.kt
class Person(private val name: String, private val address: Address) {
    fun displayInfo() {
        println("Name: $name, Address: ${address.displayAddress()}")
    }
}

// Address.kt
class Address(private val street: String, private val person: Person) {
    fun displayAddress(): String {
        return "Street: $street"
    }
}

fun main() {
    // Attempt to instantiate Person and Address without mutation
    val person = Person("John Doe", /* What do you pass here for Address? */)
    val address = Address("123 Main St", /* What do you pass here for Person? */)

    // This would lead to a compilation error due to unresolved references.
    // It's impossible to create instances without introducing some form of mutation.
}
```

You may say that this isn't real, but this is what I saw with my classmates. And you obviously "solve it" by doing a mutation (`person.address = address`), also known as a lazy initialization, which is not maintainable since there is a time period where `person.address` is undefined.

In software architecture, you learn about **object life-cycle** and **hierarchy**: "Because the life-cycle of `Address` is bound to the life-cycle of `Person` and the relationship between `Person` and `Address` is One-to-Many (one `Address` per `Person`, but a `Address` can have multiple `Person` living there), there is actually no need to reference `Person` inside `Address`.

You need a `Person`'s `Address`? Find the `Person` and return `person.address`. You need the `Persons` living at one particular `Address`, list all the `Person` and filter by `Address`. Too slow? Start indexing by creating a mutable variable.

All these issues seem silly, but they're real. And one way to know how to organize your objects is to go through these architectures. Remember that **you can't look for something you don't know exists**. You need a framework that uses several principles, so that you know later which ones are good and which are not. You can also read a book about patterns and best practices, but I will talk about this later.

Let's go back to Clean Architecture. I used this architecture since the tutorial with [Reso Coder](https://resocoder.com/2019/08/27/flutter-tdd-clean-architecture-course-1-explanation-project-structure/). Then, I used in production... in an enterprise-class project... and... oh boy, oh boy...

Let's explain why a **software architecture is not a standard**. Focus on the fact that a software architecture is a **set of principles** and that it **"attempts" to solve** the problem of complexity.

Did you get it? Assuming your team don't know sh- about Clean Architecture, what you are doing when using Clean Architecture is adding more and more complexity through an obscure choice of principles. Remember that **a person cannot look for something he doesn't know it exists?** Your team don't know what solves Clean Architecture, and why organizing your project benefits to the team.

Even if you did a Proof of Concept to explain to the team, the cost of learning that architecture is non negligible. The moment you use a "specific" software architecture in a team, it adds an incredible amount of complexity. What you need is **YOU** to architecture your software based on **simple** and **logical reasons** with your team and circumstances.

**What you need are the principles, not the framework itself.**

Do a Proof of Concept with the chosen principle, with your team, and take feedback. Eject the useless principles. Try your principles on multiple languages. Some software architecture are designed specifically for OOP, or procedural. For example, Clean Architecture does not work with Go or C.

What principles I've learned from Clean Architecture:

- Layering is useless and adds an arbitrary complexity.
- Use cases are useless as fu-.
- **Data/Domain/Presentation** is just a fancy name for **Input Data/Wanted Data/Output Data**. No reason to follow that separation. However, It is a great idea to define your business logic into immutable objects and pure functions.
- **Dependency Inversion** do work pretty well. You should depend on interfaces, not on implementation. Or more precisely, you should depend on **contracts**, not on the **side effects**. You will see that `interface`s are different depending on the programming languages.
- **DO NOT** organize your file based on architecture (no `data` directory, no `data` module).
- **DO NOT** write any boilerplate **ever**. **Boilerplate** is death. Stick with what you need.
- **KISS** (keep it simple, stupid) and **DRY** (don't repeat yourself) are the best principles. For example, if the object received in a `Service` is exactly the same as the entity (Domain object), then **do not** duplicate, **do not** create a mapper. **YAGNI** (you ain't gonna need it).
- **Dependency injection** is cool. Knowing the life cycle of your **objects** is better. You don't need dependency injection if you know where to instantiate and where to dispose your objects and services.
- ...And maybe more!

These are principles that I've learned, which **you shouldn't care much**, because I hope you **learn the principles that works best with your team**.

## Paradigms: it's simply a point of view

For this topic, I will be short.

A programming paradigm are classes of programming languages telling how state should be handled.

In reality, none of the programming paradigms add value to a project. Whether you are doing some "pure" procedural in C or some "pure" functional in Haskell..., in reality, you are just a fanatic.

That's right. Remember this: a paradigm has its own benefits and drawbacks. Seems obvious, right? In fact, the conclusion is not that obvious.

Why are you not asking this: **Why the fu- are using something that has drawbacks?** Is there no alternative? You are stuck to functional programming?

**A paradigm is a point of view**. That's literally the definition, it's simply a way of classifying programming languages, and you can have the best of both worlds. Sure, a function isn't pure, but it's **fast as fu-**. Sure, there's explicit error handling, but it's **easier to troubleshoot**.

Why the hell do you want to stick to one paradigm? Can't understand the other paradigms? Is this really a skill issue? Isn't better to be a little polyglot, so that you can use any tools? I juggle from C, Bash, C++, Perl, Ruby, Python, Go, Rust, PHP etc... and don't really care about the paradigm. All I care is: **does it work well for my environment**?

## Patterns: weapons of war

This is something you should **never** ignore. Software design patterns are reusable solutions to a commonly occurring problem. Imagine them like a type of weapon: it's a gun, it's an axe, ...

Compared to "Software Architecture" and "Paradigms", these patterns are often used by programmers, because they are considered as **best practices**. These patterns have no philosophical meaning, and doesn't define itself as a "magical way" to solve "everything". Each pattern targets one problem, with specific conditions, like a weapon.

**Patterns are battle-tested common and best practices to respond to one particular issue.**

So learn the [design patterns](https://refactoring.guru/design-patterns). Learn [when to use them](https://refactoring.guru/design-patterns/catalog). You may not use any of them, but **someone will**. You should learn them to be proficient at recognizing them. This will help to reverse engineer open-source projects or frameworks.

To help you start:

- **Creational**
  - [**Factory**](https://refactoring.guru/design-patterns/factory-method) is a object a specific set of object. For example, a `SessionManager` can create `Session` and kill all of them in one function (`dispose` function). This object is very useful when handling multiple life cycles and wish one object to control them all.
  - **[Builder](https://refactoring.guru/design-patterns/builder)** is used heavily in OOP. It is used to "compose" an object. It's basically multiple setters, and each setter has a complex behavior. There are no reason to use this, if the constructor of an object is simple enough.
  - [Singleton](https://refactoring.guru/design-patterns/singleton) can be useful at some point. But it's better to control your object's life cycle instead of using a static variable.
- **Structural**
  - [Adapter](https://refactoring.guru/design-patterns/adapter), you may see one.
- **Behavioral**
  - **[Command](https://refactoring.guru/design-patterns/command)** is used to have one object handling multiple responsibilities. Instead of decoupling into multiple objects, which add boilerplate and repetition, you can instead have multiple "command" objects.
  - [**Observer**](https://refactoring.guru/design-patterns/observer) is used in reactive programming. It's the most common pattern. This is a way to replace "push-based" systems into "pull-based".
  - [**State**](https://refactoring.guru/design-patterns/state), though I prefer **Finite State Machine** (FSM). This is used when the behavior of one object is complex, but finite. This is the best way to design a complex system with its side effects.

You can see that patterns are categorized, but that's very arbitrary. Like algorithm, these patterns should be practiced so that you get a "feeling" when to use them.

## Best practices and standards: the final word

If you are working in a team, the first thing that happen is a "practices conflict". "Formatting rules not respected", "Different IDE, different linting", "Different OS, different tooling"... This will lead to a heated conversation where the team (or team leader) set standards.

However, these best practices are unofficially considered as absolution:

- KISS: Keep it simple, stupid
- DRY: Don't repeat yourself
- YAGNI: You ain't gonna need it

Most arguments will either complement or contradict these practices. For example, Clean Architecture contradicts them all. The Single Responsibility Principle (from the SOLID principles) will complement well KISS and DRY.

Then, the standards of the language will be the final word. If you want Go developers to maintain your project, or to simply understand it, you better follow the Effective Go practices. That goes the same for Java, Rust, ... you can't write Go code like you are writing Java code.

Seems obvious? Well, if you search "Best programming languages", and take the first article, you know that's a nasty review of programming languages. Best-of sh-ty arguments:

- "Con: Steep learning curve for some libraries": Actual skill-issue
- "Con: Dynamic typing": While I do agree, that is not a valid argument since this is a selling point for some languages
- "SQL/HTML/CSS as programming languages"
- "Pro: Cross-platform": All languages are cross-platform, the question is how fast can I set up a cross-compilation environment. Java has 0 seconds, but so does Go. Not a valid argument.
- "Pro: Secure and stable": No language are secure, nor stable. No software produced by any programming language is always secure and stable.
- "Con: Heavily dependent on the X environment": Captain obvious?
- "Pro: Supported by Google": More like a Con. Actually, this argument doesn't make any sense. The maintainer is not important enough to be an argument.
- "Con: No built-in error handling (it's C)": Big skill-issue. Error handling does exist, and it's called non-zero code handling. A better argument would be "No built-in stack trace when panicking", a core dump is insufficient (though, I guess it's great for closed source software).

As you can see, none of them are actual arguments against the language and its standards. I can say better arguments than this sh- article against my favorite language:

- "Visibility is based on the naming of types. Style is strict and non-arguable. Line breaks when the line is too long is weird."
- "Have to re-implement some wheels"
- "Import cycle are not allowed"
- "No UI library"

To summary: **KISS, DRY and YAGNI as your best practices. Other best practices are to be discussed with the team. The programming language standards are absolute.**

## Conclusion

As you can see, **you can't learn by simply reading**. You have to try everything and learn every bit of each concept, preferably with a team.

Each concept has its own benefits and drawbacks, and you should try to only fetch the benefits. Software architecture can be broken into multiple principles, so **take the best principles, and throw away the rest**. Paradigms can be disputable, so take **a programming language that mixes the paradigms with only benefits** or **be polyglot**. **Learn patterns, not because you need it, but because this is actually "unofficially" standard**. And lastly, learn **best practices and standards**, because you have to use them if you want to work with a team.

TL;DR: Learn everything, test everything, take the best, be your own judge.

Congrats, you wasted time on this article. I mean, this article was to criticize influencers. This article is made to influence, so I guess I'm an influencer.

Just because an influencer says something sucks doesn't mean you shouldn't learn it. In fact, it's because it's so popular and sh-ty that you should learn about it, like any popular disease.

## References

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Reso Coder](https://resocoder.com/2019/08/27/flutter-tdd-clean-architecture-course-1-explanation-project-structure/)
- [Design Patterns](https://refactoring.guru/design-patterns)
- [Effective Go](https://golang.org/doc/effective_go.html)
