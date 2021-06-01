+++
title = "Introducing NanoMDM"
date = "2021-05-30T11:45:00-07:00"
tags = ["mdm", "nanomdm"]
author = "Jesse Peterson"
frontpage = true
+++

I’d like to introduce [NanoMDM](https://github.com/micromdm/nanomdm). From the project’s [README](https://github.com/micromdm/nanomdm/blob/main/README.md):

> NanoMDM is a minimalist Apple MDM server heavily inspired by MicroMDM.

Which is a rather minimalist description itself. To expand a bit: NanoMDM is a fully functional (yet scope-limited) "core" Apple MDM protocol server written in Golang. *Another* small open source MDM server, you ask? Yes! Doesn't that also describe MicroMDM, too? Well, yes it does. So what gives? Why a new project?

I hope to explain how this project came about, why I think it has a place here, and a little about how it differs from MicroMDM. I think this will give some insight into the goals and design of NanoMDM along the way as well. Finally, I hope to recruit your help with this new open source project.

<!--more-->

## Why?

MicroMDM is great. And not just for what it is technologically. The community and pedagogical aspect [that Victor wrote about](https://micromdm.io/blog/wwdc20-v2/) are amazing. These reasons are a big part of why I really enjoy open source software in general and MicroMDM in particular. As well my employer has been using MicroMDM in production for years now. It has and will continue to serve us well.

However MicroMDM is in the midst of changes. As Victor wrote about in the above piece there's some code that hasn't aged well and some design choices that perhaps could use a revisit. Notably MicroMDM v2 development was announced.

Inspired by these developments to take a fresh look at things I started some experimentation with low-level code related to MDM servers this year. This came on the heels of some recent low-level MDM *client* work I did with [mdmb](https://github.com/jessepeterson/mdmb). Initially this MDM server work was just experimentation with request parsing, handling, inspection, and related work. Mostly to get reacquainted with the Apple MDM protocol from the server-side. However in the context of major changes on the horizon in MicroMDM these experiments eventually gave way to thinking about what a different MDM server implementation could look like. For these reasons, and others I'll expand on in a bit, NanoMDM grew into its own project.

Another major motivator for this project (and the experimentation that preceeded it) has been some performance issues we’ve seen with [BoltDB](https://github.com/boltdb/bolt)—MicroMDM’s storage database. BoltDB with MicroMDM has performed great for years on end for us however our organization's growth finally caught up to the limitations of what BoltDB can do for us—it's affected our ability to scale MicroMDM for our fleet of enrolled devices. I’d love to go into these performance and scalability issues further sometime but suffice it to say for now that these issues contributed to the desire to accelerate plans for other storage backend options. Of course MicroMDM has always had plans to to revisit this space but unfortunately project maintainers haven't had the time to commit to this in the past.

## MicroMDM v2

When we think about MDM in general, instead of a singular monolithic technology, MDM is more like a collection of different systems, services, and servers that, when put together, comprise an “MDM server.” For example there’s usually a SCEP service, ABM/DEP API communication & integration, enrollment & profile services, maybe VPP, and of course the "core" MDM protocol that devices enroll with. Treating these services as separate-but-interworking components is the hallmark of MicroMDM. Indeed the “Micro” in MicroMDM has always stood for [microserivces](https://en.wikipedia.org/wiki/Microservices). This is true despite the fact that MicroMDM bundles and distributes these disparate services together — it is actually designed with these components as distinct under the hood.

However, not every organization needs *all* MDM services—or needs them configured in the same ways. An important goal of MicroMDM v2 will be about bringing the customization of these disparate components into a working set of services in an easy-to-configure way. Suffice it to say for now that the Micro in MicroMDM isn’t going anywhere!

Victor wrote in his [first status update on MicroMDM v2](https://micromdm.io/blog/v2dev-status-update-1/) that there has been a [v2dev](https://github.com/micromdm/micromdm/tree/v2dev) development branch started in the GitHub repo. Currently, as far as project code itself goes, there's the beginning of a front-end/UI  with support for user registrations, etc. However, as far as actual support for the business of supporting MDM enrollments by devices, that work hasn't been started just yet.

This presents some choices for us. We could just port the existing MDM protocol code over from v1 to v2 and perhaps be done with it. But for reasons laid out above I think we have the opportunity now to take a fresh look at things. Given this, it is my hope that NanoMDM will be useful for MicroMDM v2’s eventual MDM protocol handling. Maybe just in part, or maybe in whole. Maybe not at all. The future is always difficult to predict.

With that said let me be very clear: MicroMDM is not being supplanted or replaced by NanoMDM. MicroMDM v1 will continue to exist and receive improvements. Indeed PRs and changes have been submitted and merged all the while NanoMDM was initially being put together. NanoMDM is inherently limited in scope. At best, NanoMDM might be a candidate for MDM server code that eventually makes its way into MicroMDM v2 and some may find it useful on its own like me and my organization. But certainly it is no replacement for MicroMDM by itself and was never intended to be.

So, before we get ahead of ourselves, let’s talk a bit about NanoMDM.

## About NanoMDM

There’s a lot to discuss on how and why NanoMDM is put together the way it is but I think a straight-forward way to discuss NanoMDM is to compare and contrast it with MicroMDM. However I want to be very clear there is nothing wrong with MicroMDM or how it’s designed—NanoMDM was *not* designed because MicroMDM was inadequate in some way. It has just taken a different path.

Let’s start with the name: NanoMDM. A silly play off MicroMDM of course but it fits: NanoMDM is *just* an MDM protocol server. No ABM/DEP API access (though, DEP enrollments are supported). No “blueprints” (or other automated MDM commands). No SCEP services. It doesn’t even natively support TLS: you’re expected to reverse-proxy it yourself (for now, anyway). In general it simply does less. It’s only concerned with handling the low-level/"core" aspects of the Apple MDM protocol. That is, more or less, just the device & user enrollment lifecycle, sending APNs push notifications, and queueing & delivering commands. That’s it. I’ll talk more about how this reduced & focused feature set influenced NanoMDM’s architecture below.

One of the goals of NanoMDM is to specifically and directly support horizontal scalability. To that end NanoMDM comes with a MySQL storage backend. It also comes with a “local” filesystem-based storage backend that could theoretically be used to scale horizontally with something like NFS. I wouldn’t recommend that, though!

Of course MicroMDM (v1) was always meant to gain support for SQL and other storage engines but the maintainers just haven’t had the time to commit to this goal over the years. As some folks know there is even a [fork/maintained PR](https://github.com/micromdm/micromdm/pull/558) of MicroMDM with MySQL and PostgreSQL support. However that PR is very large, and represents a non-trivial amount of code to review and to support. The project maintainers haven’t had the time to commit to that and so MicroMDM has continued on with just BoltDB.

NanoMDM also had a goal of using minimal Go dependencies. This is sort of an aside as it speaks more to the project’s development philosophy, but practically speaking there were some other considerations, too. As a candidate for being included in another project (i.e. MicroMDM v2) I wanted to keep our footprint and management overhead low. [Dependency hell](https://en.wikipedia.org/wiki/Dependency_hell) is no fun and projects that pull in a bunch of dependencies exasperate that even with the wonders that [go modules](https://blog.golang.org/using-go-modules) bring. The trade-off of keeping code simple and understandable at the potential expense for some reduced flexibility and/or duplication of effort seems very inline with the Go way. We also use as much Go standard library as we’re able to and where it makes sense. Hopefully this also contributes to shared understanding of the project as we’re using more known idioms. Finally having fewer dependencies contributes to having an easier to understand project overall just by nature of having fewer indirections to follow and a smaller overall codebase. I think we’ve done well with this: currently there are only three dependencies (four if you count the MySQL driver). Naturally this will change as the project evolves but hopefully the underlying goal can be kept to.

## NanoMDM Architecture

It’s probably generous to call NanoMDM’s organization an actual “architecture” given how simple it is. But here, too, we’ll compare and contrast with MicroMDM.

MicroMDM’s original design envisioned that the various components of the MDM protocol could be split off into their own microservices. For example the APNs push notification service is separate from device enrollment service even though both are a part of the “core” MDM protocol. This is in addition to other, further separate concerns related to MDM like SCEP and DEP also being separate. To facilitate these disparate services MicroMDM employs a publish-subscribe event queue. As such MicroMDM generates events for most MDM request processing, marshals (serializes) those events, and sends them to the message queue. Other services then de-queue, unmarshal, and finally process/handle the event. As one example this is how most [“Check-in”](https://developer.apple.com/documentation/devicemanagement/check-in?language=objc) requests from devices are processed: The MDM request is converted into an event, and then the device service (which is separate from the check-in service) listens for the check-in event on the message bus rather than being called directly. This is all great for truly disparate services—communication between services is encapsulated so they can be broken apart, refactored, and scaled (or outright wholly replaced) independently if needed.

However NanoMDM takes a different, simpler approach. Given that NanoMDM’s only concern is the “core” MDM protocol it directly adapts standard Golang HTTP handlers to the NanoMDM service interface. In turn the primary MDM service directly adapts to the storage interfaces. This switch from publish-subscribe to a more request-response paradigm saves a good bit of complexity and simplifies the interfaces and flow of the server. Even simpler, the API endpoints (vs. the MDM endpoints) adapt HTTP handlers directly to the storage layer. One of my hopes is that this simpler design might help spur more contributions from the open source community.

Part of the push for this switch has also been from a few pain points in the past with MicroMDM where pub-sub was used but eventually needed to be (at least in part) request-response. This has complicated some of the services where we had to shim-in direct access to, say, the device database for example. Another specific example will be when we support Bootstrap tokens. The check-in system is largely just a publisher that other subsystems subscribe to to consume. However this will need to, in part, move to a hybrid request-response system because Bootstrap tokens actively return data. You can hear a little more about MicroMDM’s existing architecture in this [talk from 2017](https://www.youtube.com/watch?v=6DBGIDcBKFw).

Another goal with the project’s architecture was to have a clearly delineated persistence layer—called simply ‘storage’ in NanoMDM. The storage layer entirely encapsulates storage and retrieval of data in the MDM server—from enrollment data (devices & users), APNs push data & certs, queued commands & results, etc. If you're familiar with design pattern lingo this sort of resembles the repository pattern. Most of MicroMDM’s systems also share this design—called a ‘store’ in most of MicroMDM’s platform packages. However a key difference here is that each of MicroMDM’s individual services have their own store whereas NanoMDM is organized in such a way as to implement all of the storage, in each backend, for the whole system. While this front-loads the effort of creating a new storage backend (because you more or less have to implement all of it at once) I think the tradeoff is a simpler, overall easier to understand interface. As well the reduced scope of NanoMDM should also reduce this burden.

While this describes the “front” end of NanoMDM (Go HTTP handlers) and the “back” end (storage) I also want to highlight the service layer in the middle. The service (or services) layer is a composable interface that represents MDM client requests and is directly inspired by MicroMDM. All of the things we want to drive from MDM requests happen in the service layer. As mentioned above, actually storing and retrieving data is driven from the service layer. Like MicroMDM there is also a webhook layer which is just another NanoMDM service. The request “dumper,” for debugging, is an example of service middleware, as is the certificate authentication feature. They all share the same interface of a service and are composed and layered together to bring about the server's functionality. Again this is all inspired by MicroMDM.

It’s been hinted at here, but I want to explicitly call it out: a primary difference (perhaps *the* difference) between MicroMDM and NanoMDM is that, while MicroMDM’s implementation of the “core” MDM protocol is componentized, NanoMDM’s is a bit more unified. In other words while MicroMDM envisioned splitting even the core MDM protocol amongst different microservices NanoMDM considers the core MDM protocol feature-set, more or less, a singular concern. Now, technically speaking that’s not the whole story — indeed all of the separate components of the core MDM protocol are modular in NanoMDM’s code, too, and could be split out and used separately. In fact NanoMDM can, right now, be operated in several modes that only handle certain ‘concerns’ such as API (pushing, enqueuing commands) or MDM protocol/enrollment handling including splitting out check-in and command endpoints. However, in general, the choice was made for NanoMDM’s interface design and default out-of-the-box operation to err on a simpler, more unified operation and optimize for a straight-forward request-response design.

The trade-off that was made from hyper-modularity to server simplicity I think is good for NanoMDM’s goals. The MicroMDM server has been distributed as a monolith single binary since its first release—there hasn’t been much call for splitting MicroMDM’s individual core MDM protocol services up—lending some reassurance that this is an acceptable direction to take. As already mentioned I hope a simpler design encourages more code contributors. Further I think by constraining the modularity in NanoMDM will allow for *other* modularity to take place in other components that are more clearly distinct—such as DEP/ABM, VPP, workflow engines, etc. Perhaps this will aid future integration with MicroMDM v2.

## Practically speaking

We’ve covered a bit about the design & internals of NanoMDM. But what about more practical matters like, for instance, just getting it running? NanoMDM is new so its documentation and resources are few just yet—but here’s a few things that hopefully help to get started and/or learn more:

* Check out the [NanoMDM project page](https://github.com/micromdm/nanomdm). The README has an overview of the project, high-level of features, and of some of the operational differences from MicroMDM.
* Like MicroMDM’s server component, NanoMDM is just a single binary. [Download it](https://github.com/micromdm/nanomdm/releases), check out the `-help` flag, try it out!
* The [quickstart guide](https://github.com/micromdm/nanomdm/blob/main/docs/quickstart.md) aims to guide you through an example setup.
* The [micro2nano](https://github.com/micromdm/micro2nano) project provides some tools for migrating from a MicroMDM server to NanoMDM. I hope to write more about these tools!
* If you’re interested in the MySQL backend then take a look at [mysqlscepserver](https://github.com/jessepeterson/mysqlscepserver). It's a simple, somewhat opinionated, SCEP server with a MySQL backend based on MicroMDM’s SCEP server.

I hope to get more time to create more documentation, guides, tools, etc. for this project. Also, your help here is more than welcome, too! Which brings me to:

## Call for participation

I hope to encourage discussion and participation in NanoMDM. There’s few specific things I’d like input on. I've linked to the relevant GitHub Discussions as a place to further discuss these items:

* [The MySQL schema](https://github.com/micromdm/nanomdm/discussions/3), specifically columns types and table relationships
* [PostgreSQL](https://github.com/micromdm/nanomdm/discussions/5) (or, rather, the lack thereof)! There’s a MySQL backend. Care to take a swing at implementing this backend?
* Help design the (currently missing) [enrollment (device) API](https://github.com/micromdm/nanomdm/discussions/6)
* Plenty of issues/considerations in [issue tracker](https://github.com/micromdm/nanomdm/issues)! Or, make a new one!
* Bug reports, bug fixes, documentation fixes and improvements, etc. PRs welcome!

NanoMDM’s scope is more limited than MicroMDM and because of this it may be even more difficult for newcomers to get started with—getting MicroMDM v1 itself going is certainly no walk in the park. Low-level MDM is not for everyone. That said the goal of any open source project is to be useful and hopefully see some adoption. To that end I hope folks find NanoMDM useful and I look forward to collaborating with folks on it.

Come join the fun in #nanomdm on the [MacAdmins slack](https://www.macadmins.org/) and thanks!
