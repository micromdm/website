+++
title = "WWDC20: What's new in managing Apple devices"
date = "2020-06-30T00:00:00+05:00"
tags = ["mdm", "wwdc", "opinion"]
author = "Victor Vrantchan"
frontpage = false
+++

Here are some of my opinions about Mobile Device Management (MDM) related changes announced last week. SimpleMDM did a great job of [summarizing](https://simplemdm.com/mdm-ios-14-macos-11-big-sur/) the announcements, so I won't repeat the same list here. Instead, I'll focus on a few things which are interesting to me. Before I get into that, I have to mention this isn't a well-balanced editorial. It's in my nature to skip over accomplishments and focus on things that still need attention. Which is what you'll read below. But it's important to also state that I enjoyed WWDC20 and I'm excited about many of the technology changes announced. I've also had the chance to meet and work with many Apple employees over the years. My grumpy opinions on software updates and technical debt are not a reflection of how I feel about any individual or even Apple as a whole.

Okay, let's dive in.

## Account Configuration can choose the MDM managed user

Apple claims MDM can [create user accounts](https://developer.apple.com/documentation/devicemanagement/accountconfigurationcommand/command), but it comes with all kinds of limitations. First, the account can *only* be created during the setup assistant stage, while the device is in AwaitingConfiguration state. This leaves users unable to implement a useful workflow since almost everyone has a mixed fleet of "company-owned" and legacy devices. You need a fully implemented fallback workflow in case the automated one fails. There's a number of failure modes that need to be addressed: 

- The device needs to be in [Apple Business Manager](https://www.apple.com/business/it/) (ABM), an inventory system that tracks company-owned devices.  ABM is still not available globally, and there is no 30 day grace period where you can add new purchases in yourself (available for iPhones ).
- There's no enforced activation check on the Mac like there is for phones. Not that I'm advocating for one, but losing controls because a Setup Assistant step fails open is a pain.
- Since you can only rely on this in Setup Assistant, changes made in new releases are not possible to adopt. Even if AccountConfiguration was perfect in macOS 11 and we used it for all *new* devices, it would still take years before all our devices had consistent user creation.

Since the recommended way has these failure modes, most organizations will keep using their own user creation flow. With the caveat that one day Apple will break it. 

Despite the shortcomings, AccountConfiguration is tempting. Admin accounts configured with this command allow for MDM managed [password changes](https://developer.apple.com/documentation/devicemanagement/setautoadminpasswordcommand/command). The MDM could integrate with existing account services to set/change the Mac password by receiving a PBKDF2 derived hash. Neat! Especially for 1:1 environments without traditional directory service binding. 

While I was testing the AccountConfiguration command last year I ended up filing feedback about another artificial limitation of the process. Apple was making a distinction between the "admin" account, and a "primary user". The primary user would be created not by the MDM, but by the end-user typing their username/password into setup assistant, while the MDM was creating a hidden local administrator account for an IT user. There was no way for these two accounts to be the same, even though that's what I wanted. A limitation of not being the *primary* account meant that while the MDM could rotate this admin account password, it could not send profiles or other user channel actions. Well, not any longer!

With the release of macOS 11, there's a way for the MDM to *choose* the account that becomes managed. I [tested](https://twitter.com/wikiwalk/status/1275622118324162561) this workflow and made some [code changes](https://github.com/micromdm/micromdm/pull/679). It works OK now.

Personally, I'm bewildered by much of how this works and why it exists in the first place. Last year Apple announced support for "Managed Apple IDs", which can be federated with major identity providers. This year they also announced support for SCIM, a complementary protocol that federates changes to user accounts across different systems. But I can't use my SSO account to log into my company-owned MacBook.

## Local profile installation is gone

It's not possible to install configuration profiles using the `profiles` command any longer. Profiles now work like iOS; the only two options are MDM or manual installation through System Preferences.

Overall, this is a welcome security change. Malware had the potential to exploit this stuff easily, and this is another class of exploits that Apple eliminated. Apple has been so good at making a "secure by default" platform, that I hesitate to buy hardware from other vendors.

While the profiles change is welcome, it's still a major disruption for many in IT. Similar to DevOps on the server-side, Mac administrators have relied on configuration management tools ([Chef](https://docs.chef.io/resources/osx_profile/), [Puppet](https://github.com/macadmins/puppet-mac_profiles_handler), [Salt](https://github.com/mosen/salt-osx/blob/master/mac-examples/profile.sls)) to generate and install the profiles with code. It's convenient, reliable, and does something which is otherwise hard: being able to dynamically enforce values depending on local system variables. For example, instead of enforcing an immediate screen lock setting, letting the end user choose a value that is acceptable within the range of compliance. Or detecting network changes and only managing something depending on network state. The exact kind of things configuration management is good at.
In addition to the flexibility, organizations that rely on configuration management also tend to have good change management and testing processes. The policies are typically versioned in a system like Git, with linting, code review, and all the bits that come with treating infrastructure as code.

I started working on MicroMDM for two reasons. One was the general curiosity of how things like "zero-touch" provisioning worked. The other was a disappointment with commercial options. I wasn't willing to give up my management workflows, which were driven by Puppet, for the privilege of logging into a web interface and clicking a checkbox or filling out a web form. Commercial vendors offered little as APIs that could support code review, auditing and idempotent management. And often the web forms they did offer missed settings organizations need. I'm sad to say that very little has changed on the commercial front. In open-source, we have examples like [mdmdirector](https://github.com/mdmdirector/mdmdirector)👏 💯, but you can't pay someone money to give you these features.

## Custom app distribution in Apple Business Manager

Apple [announced that there's a new way](https://developer.apple.com/videos/play/wwdc2020/10667/) we can distribute Apps this year. Enterprises can upload custom apps to Apple Business Manager, and these can be installed with MDM. These apps are still, of course, subject to app review by Apple. This new distribution channel is also open for vendors to sell apps to other enterprises.
My understanding of this announcement (which might not be what Apple intended?) is that B2B products, like endpoint security vendors, can now show up here. With a 30% cut for Apple?

It's still too early to tell, as this is just an announcement. I'm excited to try and get an internal app added to this distribution channel. Will we see something like [Nudge](https://github.com/macadmins/nudge) in the App Store one day?

At the same time, having yet another "official" way of distributing apps, makes it more likely that one day there will not be a distribution channel where Apple is not in the middle.

## Managed Software Updates

I've [written](https://micromdm.io/blog/os_update/) about software updates with MDM in the past. It's been a disaster so far. Apple made one change to MDM this year, which allows the MDM to "force" an update. And while that's a necessary *capability*, it's still not something that can be realistically implemented. There is no UX on the user side of the equation. The Mac would restart because the MDM decided, without giving you a chance to quit, without care that a VP is giving a presentation to the board or a doctor is in the middle of a telemedicine call. Apple expects (I guess? There is never official guidance for how these things should be done) that the MDM product will devise an appropriate experience for the user. Having a third party client coordinate the exact moment an MDM server should send the *InstallForceRestart* command is subject to potential race conditions and coordination problems. And of course, nobody has an entitlement to ignore a user's DND settings, while they're about to reboot a device.

What the community has always asked for is for Apple to provide a declarative interface, where the MDM sets

- a minimum OS version/build.
- a deadline when the version is enforced.

In a hypothetical scenario where Apple offered declarative controls to enforce an update, the built-in client would be responsible to deliver the appropriate experience for the end-user. Notifying them well ahead of time about the pending deadline, and potentially offers a safety window when the user can cancel. Like Munki [has done](https://github.com/munki/munki/wiki/Pkginfo-Files#force-install-after-date) for a decade and some vendors have implemented with their agents. For reasons I can only speculate about, Apple has refused to build OS update enforcement correctly.

The one big change in this area is that the way OS updates are done on the Mac is now completely new. Apple replaced the Mac implementation of OS update, with the iOS one. So, maybe I'm wrong, and OS updates are finally reasonably OK to do with MDM on macOS 11? There are no items in the catalogs yet that I can test with. I'll continue testing and filing feedback in this area.

## Apple buys Fleetsmith

The most exciting and surprising bit of WWDC didn't even happen at the conference. [Fleetsmith is now part of Apple](https://blog.fleetsmith.com/fleetsmith-acquired-by-apple/). I've had a chance to meet the Fleetsmith engineering team through various industry events, and they've always impressed me, not to mention how polished their product is. So I'm very happy for that team, and also enthusiastic about what this means for Mac in the enterprise. Mostly. Apple legal did [break existing users](https://twitter.com/wikiwalk/status/1276139303219970053).

Allow me to speculate a bit about what this change means:

*The Good*

As I already mentioned, Fleetsmith has the right engineering talent. They also have the industry experience to understand what companies need and [push for the same kinds of changes I would](https://blog.fleetsmith.com/macos-enterprise-security-roadmap/). So, at least my reaction was thinking I now have an ally on the inside and would see welcome changes in the future.

*The Bad (for the commercial vendors)*

I can't be the only one who's ever speculated/expected to see a MDM product integrated into Apple Business Manager. I think most of us would welcome it, even if the initial offering didn't have everything we need. Right now, the commercial options are mostly disappointing and rely on their sales teams, rather than engineering talent to be in the market. The acquisition of Fleetsmith has got to force some to compete on the merits of their product.

There are also lots of small MDM vendors that are trying to compete with a tiny set of features and come with security risks which prevent them from gaining serious market share. A commercial product from Apple, even one with a limited initial set of features will impact many of these companies.

*The Scary*

I worry that Apple has the hubris to think their offering is enough for everyone. Maybe not at first, but eventually. Imagine a hypothetical scenario where macOS 12-14 comes with a new must-have feature. Let's say Apple allows you to upload a modified APFS snapshot, sign it for you and allow you to boot a Mac from it. But you must use Apple MDM to do it. And slowly the existing management options are replaced by better ones. But they're all exclusive to the Apple service.

I find this scenario scary and yet realistic. There's a lot of money in it. And too often, Apple has been myopic about how companies use their devices. My personal experience has been the opposite, and I don't believe Apple even dogfoods MDM on their own engineering fleet. I hope there's a future where we can continue to orchestrate the devices our companies own the way we want. But hope is not a strategy. And I'm not optimistic.

In conclusion, the more things change, the more things stay the same. Which is another way of saying that changes creep up on you in small increments? And while my overall impression of WWDC20 is that very little has changed in MDM this past week, we will never manage company-owned devices [the way we've always done](https://www.youtube.com/watch?v=4Qv8CE6D_Rc).
