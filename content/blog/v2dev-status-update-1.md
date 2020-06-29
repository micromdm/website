+++
title = "Dev Log: My first status update on v2"
date = "2020-06-29T00:00:00+05:00"
tags = ["v2", "v2dev", "status"]
author = "Victor Vrantchan"
frontpage = true
+++

After announcing the kickoff of a big rewrite last week, I wasn't sure how to actually get things going. There's a lot I want to get done, but I also want to do it in a way that tells a story. So it's important not to rush too fast ahead. Then Saturday, fueled by a fresh pot of [good coffee](https://parlorcoffee.com/collections/all-products/products/prospect) and WWDC excitement, things really took off. Now I have a bunch to share.

## v2dev working branch

I created a new [`v2dev`](https://github.com/micromdm/micromdm/tree/v2dev) git branch which will eventually become the default. I kept the git history up to this point but removed all files except for the LICENSE, and CODE_OF_CONDUCT. Both are essential and unchanged.

Next, I added a new [CONTRIBUTING](https://github.com/micromdm/micromdm/blob/v2dev/CONTRIBUTING.md) document. I tried to make it welcoming to new contributors. If you spot an opportunity to improve it, [let](https://github.com/micromdm/micromdm/issues/new) [me](twitter.com/wikiwalk) [know](https://github.com/micromdm/micromdm/pulls).

## website/ in the main repository

The `micromdm/website` repo had the code for [micromdm.io](http://micromdm.io). I [moved](https://github.com/micromdm/micromdm/tree/v2dev/website) it into the v2dev branch instead, [preserving](https://github.community/t/combining-repositories/2060/2) the original history. I think having fewer repositories will be easier to maintain. Breaking things out into multiple repositories makes sense for large teams, or for projects that appeal to a new audience ([micromdm/scep](github.com/micromdm/scep) is a good example). 

## Architecture Decision Records

Architecture Decision Records (ADR) are a way of documenting decisions to capture the context leading up to, and the consequences of adopting the decision. ADRs are adopted in [some engineering organizations](https://engineering.atspotify.com/2020/04/14/when-should-i-write-an-architecture-decision-record/) and I've advocated for their use on some teams I work with.

Joining a project that's been around for a while is confusing. Questions like "Why is this project using X?", or "Why did you not do Y?" come up a lot. Talking to the engineers (if they're still around) the answer is usually. "Well that API didn't exist three years ago.", or "We tried to do X but it didn't work. We had a deadline, so we did Y instead." ADRs help to explain why someone made a specific choice at a specific point in time. 

I created [`docs/architecture/decisions/`](https://github.com/micromdm/micromdm/tree/038fd359aeba5f4b6a2273a2662cd89375dd0e13/docs/architecture/decisions) as a place to capture these decisions. For now, the only file is a template.

```markdown
# Title Which Captures The Decision

## Status

What is the status, such as proposed, accepted, rejected, deprecated, superseded, etc.?

## Context

What is the issue that we're seeing that is motivating this decision or change?

## Decision

What is the change that we're proposing and/or doing?

## Consequences

What becomes easier or more difficult to do because of this change?
```

I'll try and publish proposals and their outcomes on the blog.

## Go module for v2

In Go, every import outside of the standard library is specified [with a URL](https://golang.org/cmd/go/#hdr-Remote_import_paths). To import the webhooks module in MicroMDM use [`github.com/micromdm/micromdm/workflow/webhook`](https://pkg.go.dev/github.com/micromdm/micromdm/workflow/webhook?tab=doc). To hash passwords use [`golang.org/x/crypto/bcrypt`](https://pkg.go.dev/golang.org/x/crypto/bcrypt). For MicroMDM, I wanted to reduce the coupling with GitHub, and also give the imports a bit of branding.

```bash
go mod init micromdm.io/v2
```

To depend on a library you would run:

```bash
go get -u micromdm.io/v2/pkg/useful/library@v2dev
```

The `@v2dev` portion is temporary and points to the development git branch. Specifying it won't be necessary when I change the default branch on GitHub.

## Hosting on Netlify

Up to this point, the website was updated by pushing the compiled [Hugo](https://gohugo.io/) output to a cloud storage bucket. To support micromdm.io/v2 in the [import path](https://golang.org/cmd/go/#hdr-Remote_import_paths), I needed something that would serve the go-import `meta` tag 

```bash
<meta name="go-import" 
    content="micromdm.io/v2/pkg/log git https://github.com/micromdm/micromdm">
```

I signed up for Netlify because they make it extremely easy to host a static site and add a bit of on-demand logic with lambda [functions](https://www.netlify.com/products/functions/). I've never used Netlify before, but within an hour I had the entire website migrated and running the vanity URL service.  

The `website/_redirects` files configures Netlify to redirect HTTP requests for `v2/*` to the lambda service.

```bash
/v2/* /.netlify/functions/vanity 200!
```

The function is a Go binary, compiled with this make target as part of the deploy process:

```bash
.netlify/functions/vanity: netlify-functions/vanity/*.go
        GOBIN=$(PWD)/.netlify/functions go install ./netlify-functions/vanity
```

And [here is the code](https://github.com/micromdm/micromdm/blob/038fd359aeba5f4b6a2273a2662cd89375dd0e13/website/netlify-functions/vanity/vanity.go) it executes.

I'm not sure I'll stick with Netlify long term, or if I'm going to end up moving to a cloud provider with more features later. Some of the projects I have in mind also require a database. For now, this was the cheapest (it's free!) set up and requires no maintenance.
