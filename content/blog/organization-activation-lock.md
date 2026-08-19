+++
title = "Organization-linked Activation Lock & Bypass"
date = "2026-08-19T04:27:17Z"
tags = ["micromdm", "nanomdm", "nanodep", "bypass", "activation", "lock"]
author = "Jesse Peterson"
frontpage = true
+++

To quote our friends over at Iru (formerly Kandji):

> Activation Lock is a theft-deterrent feature found in [iOS and iPadOS devices](https://support.apple.com/HT201365) and modern [Mac computers](https://support.apple.com/HT208987) (with the Apple T2 Security chip and Apple silicon).

NanoMDM and NanoDEP now have direct support for organization-linked Activation Lock, the flavor of Activation Lock that is *not* tied to a user's Apple Account (i.e. not "Find My"). NanoDEP (as of [v0.7.0](https://github.com/micromdm/nanodep/releases/tag/v0.7.0)) can **lock** a device to your organization using the Apple DEP API, and NanoMDM (as of [v0.9.0](https://github.com/micromdm/nanomdm/releases/tag/v0.9.0)) can **unlock** a device using Apple's Escrow Key Unlock endpoint. Both sides revolve around the humble Activation Lock Bypass Code. Below we'll walk through both, step by step.

<!--more-->

I encourage you to read [Iru's excellent write-up on Activation Lock](https://www.iru.com/blog/how-to-manage-activation-lock) by Mike Boylan to ground yourself in what activation lock is and how it works! Also, while it does contain some proprietary tool specific information, [SimpleMDM's take on managing Activation Lock for enterprise environments](https://simplemdm.com/blog/how-to-manage-activation-lock-for-enterprise-environments/) is also a good read to get some background, too.

## A quick primer

Before we get our hands dirty, two concepts worth having in the back of your head.

First, on a *supervised* device, "Activation Lock" and "Find My" are not the same thing. Apple gives you [two paths](https://support.apple.com/guide/deployment/activation-lock-depf4ab94ef1/web):

* **User-linked Activation Lock** is the Find My-adjacent one: the device gets tied to whichever Apple Account signs in and turns on Find My. Managing it is a per-device dance between your MDM and the device itself.
* **Organization-linked Activation Lock** is the one this post is about. Here you don't touch the device at all — because the source of truth for Activation Lock is Apple's *servers*, your MDM simply tells Apple's servers "lock this serial number." The device becomes Activation Locked *to your organization*, with no Apple Account, no Find My, no iCloud involved. It's a pure server-side lock, originally introduced way back in iOS 9.3 for Shared iPad.

Second, the thing that makes both lock and unlock work is the **Activation Lock Bypass Code** (or "ALBC"). A single bypass code has three forms:

* **Code** — the dash-separated string a human can read (and type), like `3UM43-PUYVY-QYD1-UVCC-HEHJ-FKA4`.
* **Raw** — the underlying 16 bytes, hex-encoded.
* **Hash** — a hex-encoded PBKDF2-derived hash of that value.

These are *not* all interchangeable: the **hash is one-way.** Given the `code` (or `raw`) you can always derive the `hash`, but you can never go backwards from `hash` to the `code` you'd need to unlock with. So when you generate a bypass code, store the **code** (or `raw`) — that's the form that can later become whatever you need.

To that point: when you *lock*, Apple wants the **hash**; when you *unlock*, Apple wants the **code**. We'll point this out again below.

## The setup

Both walks below assume you've already stood up the two Nano-suite pieces, per their [quickstart guides](https://github.com/micromdm/nanodep/blob/main/docs/quickstart.md):

* A running **`depserver`** (NanoDEP), configured with your DEP tokens and reachable at a known URL with an API key.
* A running **NanoMDM** server with your APNs push certificate uploaded.

The NanoDEP `tools/` directory ships a handful of ready-made shell scripts for exactly this kind of work. They just need a few environment variables:

```bash
export BASE_URL='http://[::1]:9001'   # your depserver
export APIKEY=supersecret              # should match depserver's -api flag
export DEP_NAME=mdmserver1             # the DEP "name" you configured
```

With those set, we're off.

## Activation Locking a device

Let's say a device has gone missing, or you just want to make sure the next person to walk off with a Mac doesn't get a free computer. You want to lock it to the org.

### 1. Make a bypass code

A bypass code is the key that will later let *you* unlock this device. Generate one now and, crucially, **save it**:

```bash
$ ./bypasscode
b7608a5fd1fa8500c6095ee152ea4caf  raw
PXH8M-QYJZA-2H1J-H9CV-HN5U-KDN7  code
6789c28a5a596746af3e7e4d18facb19c1161e4053ccca02a3e65ce838938c0e  hash
```

([`bypasscode`](https://github.com/micromdm/nanodep/blob/main/docs/operations-guide.md#bypasscode) is a stand-alone helper shipped with NanoDEP; if you prefer, the same thing is available over HTTP at [`GET /v1/bypasscode`](https://github.com/micromdm/nanodep/blob/main/docs/operations-guide.md#activation-lock-bypass-code) on your `depserver`.) Keep the `code` somewhere safe — you'll need it to unlock later. For locking, you'll use the `hash`.

### 2. Lock the device

Hand the serial number and the `hash` to the DEP API. The [`dep-activation-lock.sh`](https://github.com/micromdm/nanodep/blob/main/docs/operations-guide.md#dep-activation-locksh) script does the heavy lifting:

```bash
$ ./dep-activation-lock.sh C8TJ500QF1MN 6789c28a5a596746af3e7e4d18facb19c1161e4053ccca02a3e65ce838938c0e "This Mac belongs to Acme Inc."
```

That's it. Apple's servers now have that serial marked as Activation Locked to your organization. You can optionally leave off the third argument (the "lost message"), but it's a nice touch — it shows on the device's lock screen, so a well-meaning finder sees who to call. You can also omit the hash, but then the lock ties to whoever created the MDM server token, and you lose the easy unlock path.

## Activation Unlocking a device

The device has come back home (or a user returned from leave, or you're reprovisioning). Now you want to clear the lock.

### 1. Dig up the device details and your saved code

You'll need the device's `serial`, its `productType` (like `MacBookPro17,1` — this is model info your MDM will already have), and the **code** you saved back in step 1 of locking. Note the shift: this time it's the human-readable `code`, *not* the `hash`.

### 2. Tell NanoMDM to unlock it

NanoMDM's [`/v1/escrowkeyunlock`](https://github.com/micromdm/nanomdm/blob/main/docs/operations-guide.md#escrow-key-unlock) endpoint wraps Apple's Escrow Key Unlock call. Fire it a single request:

```bash
curl \
	-u nanomdm:nanomdm \
	-d "topic=com.apple.mgmt.External.f3abfeac-1f27-4c8e-8a63-dd17555d35d9" \
	-d "serial=C8TJ500QF1MN" \
	-d "productType=MacBookPro17,1" \
	-d "escrowKey=3UM43-PUYVY-QYD1-UVCC-HEHJ-FKA4" \
	-d "orgName=Acme Inc" \
	-d "guid=12346" \
	'http://[::1]:9000/v1/escrowkeyunlock'
```

And the lock is cleared. The `topic` is your APNs push certificate topic (NanoMDM uses it to authenticate to Apple via mutual TLS), and `orgName`/`guid` are just audit strings you get to choose.

### No server handy? Type it directly

If the device is sitting in front of you and you don't have (or don't want) your MDM server in the loop, the bypass code works by hand too. On an iPhone or iPad, enter the **code** in the Apple ID password field and leave the username blank. On a Mac, open Recovery Assistant and choose "Activate with MDM Key," then enter it there.

### No code at all? Turn it off in ABM or ASM

If the device is organization-owned and lives in your Apple Business Manager or Apple School Manager instance, you can sidestep the bypass code entirely: sign in with a role that has [Manage Devices](https://support.apple.com/guide/apple-business-manager/turn-off-activation-lock-axm812df1dd8/1/web/1) privileges, select the device, open the menu, and choose **Turn Off Activation Lock**. Handy for lost codes and emergencies.

## For integrators: what's under the hood

The walkthrough above leans on NanoDEP's shell scripts, but everything they do is also exposed as libraries and endpoints for you to build on.

* **NanoDEP** adds a [`godep.Client.ActivationLock()`](https://pkg.go.dev/github.com/micromdm/nanodep/godep#Client.ActivationLock) method wrapping Apple's [Enable activation lock on a remote device](https://developer.apple.com/documentation/devicemanagement/activation-lock-devices) endpoint:

  ```go
  resp, err := client.ActivationLock(ctx, name, serial, escrowKeyHash, "This is an example lost message.")
  ```

  The `GET /v1/bypasscode` endpoint (and the `bypasscode` binary) generate or parse codes, returning all three forms. `tools/dep-activation-lock.sh` is just a thin `curl`+`jq` wrapper around the `depserver` reverse proxy — and that reverse proxy is directly available to you too, at `POST /proxy/{name}/device/activationlock`, if you'd rather drive Apple's endpoint yourself without the Go library.

* **NanoMDM** adds the `POST /v1/escrowkeyunlock` endpoint, which accepts the same form fields as the curl example above (`serial`, `productType`, `escrowKey`, `orgName`, `guid`, plus `imei`/`imei2`/`meid` for cellular devices) and forwards them to Apple's `escrowKeyUnlock`, authenticating with your stored APNs certificate. For Go integrators the same operation is exposed directly via [`escrowkeyunlock.DoEscrowKeyUnlock()`](https://pkg.go.dev/github.com/micromdm/nanomdm/http/escrowkeyunlock#DoEscrowKeyUnlock).

Apple's [Creating and Using Bypass Codes](https://developer.apple.com/documentation/devicemanagement/creating-and-using-bypass-codes) documents the code format for the curious: the human-readable form is a base-32-ish encoding of the raw bytes, and the hash is a PBKDF2 derivation of those bytes. None of that is anything you need to implement yourself, though — the NanoDEP `bypasscode` tool and `/v1/bypasscode` endpoint handle it for you, and both are built on the [`albc`](https://pkg.go.dev/github.com/micromdm/nanodep/albc) package, which is available to your own Go code if you need the code↔hash↔raw conversions.

## A few more notes

* **The bypass code still matters, but losing it is no longer a dead end.** Organization-linked locking is server-side, so there's no device-derived escrow key and no 14-day retrieval window like user-linked Activation Lock. Happily, since WWDC24 (see below) administrators with the appropriate role can [turn off Activation Lock](https://developer.apple.com/videos/play/wwdc2024/10143) for organization-owned devices directly in Apple Business Manager or Apple School Manager — select the device, open the menu, and choose "Turn Off Activation Lock." That's a handy fallback for lost codes, general unlocking, and emergencies, alongside AppleCare for devices that aren't listed there. The DEP server API that NanoDEP and NanoMDM use is very much still around, too.
* **Only one lock at a time.** If a device somehow has both a user-linked and an organization-linked lock, the first one wins. On a supervised, organization-linked device you generally want Find My left off so there's no ambiguity about who holds the device.
* **This is a building block.** Both projects deliberately expose these as low-level primitives — endpoints, a Go method, and scripts — rather than an opinionated "activation lock module." That's very much the Nano-suite way: wire them together however your own tooling needs.

Finally, keep an eye on Apple's WWDC announcements (as always). [WWDC24](https://developer.apple.com/videos/play/wwdc2024/10143) brought the ability to turn off Activation Lock for organization-owned devices from within Apple Business Manager and Apple School Manager — the manual fallback above. What hasn't landed (as far as I know) is programmatic Activation Lock and unlock through the [Apple School and Business Manager API](https://developer.apple.com/documentation/apple-school-and-business-manager-api), which would be a natural fit for [NanoAXM](https://github.com/micromdm/nanoaxm) and something we'd love to add support for. I have no insider information here, but if you'd like to see it too, please [file feedback](https://feedbackassistant.apple.com/) with Apple.
