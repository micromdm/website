+++
title = "About Software Updates in Big Sur"
date = "2020-12-17T00:00:00+05:00"
tags = ["mdm", "softwareupdate", "opinion"]
author = "Victor Vrantchan"
frontpage = false
+++

If you work with Apple in some capacity, you know that they're not very likely to admit mistakes. I'm not aware of Apple publishing postmortems after outages or providing details about known issues. So it's up to the developer and admin communities at large to help each other learn about outages and potential causes. With the release of macOS 11.1 this week I've been debugging a new issue with the Mac software update process, one which will affect most enterprise users. I decided I should write about it and let everyone know what the issue is, workarounds, and how to avoid it. I've also been paying close attention to some recent changes to the software update mechanism, so I'll try to mention them below as well.

## A macOS 11 bug prevents upgrades to the next version

After macOS 11.1 was released earlier this week, many users [started reporting](https://twitter.com/Contains_ENG/status/1339399335298166785) that they were not able to see the software update. Others reported that they saw it, but were not able to download it. I personally experienced this issue when 11.1 was in beta and filed a case with Apple about it, but then moved on to other problems. When the final release came out this week, there were widespread reports of it on the MacAdmins Slack. Reading through system logs, I as well as other admins were able to find what appears to be the root cause. **Under certain conditions, macOS 11.0.1 and macOS 11.1 hosts are requesting the update server send the 11.0.1 update, instead of requesting the next available one. The server rejects this update as it's already either installed or older.** This somehow corrupts the state of the software update process, and the update is no longer visible as an option in System Preferences. 

{{<figure src="/big_sur_softwareupdate/log.jpg" title="" class="screenshot" >}}

### Workarounds

- The update is visible again immediately after the restart. But it's unclear if the update can always be installed after it's visible since the condition that made it fail to download the first time can be triggered again.
- Removing the MDM enrollment profile causes the update to be seen again. I and several others tested this extensively, and it's definitely the enrollment profile, not some other policy managed by MDM. This is obviously a terrible workaround and I'd hesitate to mention it to users as it could cause security agents to stop working, and countless approval dialogs we're so familiar with. Not to mention some of you have the MDM enrollment profile as non-removable or users who are not administrators, so they don't have permission to unenroll. Getting them back enrolled in the MDM might prove to be a challenge too.
- Distributing the full 11.1, and eventually the 11.2 installers.

### Can Apple fix this bug without manual intervention?

Apple is well aware this is a problem now, so I am confident the issue will be addressed in 11.2. Unfortunately, it is a client-side issue affecting both 11.0.1 and 11.1 clients. So there's not much Apple *can* do to provide a fix. One potential solution I see is for Apple to detect the wrong request and instead of rejecting the download, offer the right file archive instead. But this is a complex system and it's unclear if the server-side changes alone are enough. 

Something else Apple could try is to side load a patch through another software update. There's background configuration and malware removal tools that are likely capable of fixing the issue on the system. But it's an ugly hack and one I'd personally stay clear of even if the option was available. 

In my opinion, the most likely outcome is that the bug will be fixed in 11.2, but clients that have already upgraded to Big Sur (or any M1 macs you might've bought) will have to work around the problem themselves. If we're lucky Apple will publish a support article and that will be that. 

## What else you need to know

Big Sur has changed the software update mechanism entirely, especially on Apple Silicon macs. It's a long time coming and Apple spoke about some of this at WWDC in the MDM and IT sessions, so it shouldn't be entirely surprising. But a lot was left unsaid for us to discover on our own. 

- Combo update packages are [no longer published](https://eclecticlight.co/2020/12/17/apple-has-stopped-providing-standalone-installers-for-macos-updates/) on the Apple website. This might surprise many of you who rely on distributing them. The main reason is that the entire format of the updates has changed, and it's also no longer possible to install updates without the Mac having access to the internet. The updates must come from Apple, and they must be managed by `softwareupdate` and related processes on the system.
- On Apple Silicon, third party processes are no longer able to script the `softwareupdate` command. Running the software update command as a root process now prompts for the administrator password, who also needs to be a secure token user. There are only two possible options for OS updates to be applied; either the human user of the device itself or the MDM process [which has special permission](https://support.apple.com/guide/deployment-reference-macos/using-secure-and-bootstrap-tokens-apdff2cf769b/1/web/1.0). It might also be possible for the Mac to update itself with the auto-update mechanism, but there are too many bugs right now to observe how well that works. My personal experience is that it doesn't and I'm greeted with a password prompt to enter the next day...
- Specifying a custom URL for the `softwareupdate` process to pull updates from is no longer possible. Apple advertised this deprecation for all of last year, so it should surprise no one but it hurts. [Reposado](https://github.com/wdas/reposado) was one of several great tools that made it possible to have unstable/beta/stable tracks within an organization.
- —ignore is gone as a flag. It [was gone](https://mrmacintosh.com/10-15-5-2020-003-updates-changes-to-softwareupdate-ignore/) in 10.15.5 briefly too, so you likely know about this one. This time it's gone for good and never coming back. An MDM server can delay the client from seeing OS updates for up to 90 days only, but that is the absolute maximum going forward. Even for a future major release like macOS 12. If this is important to you, file feedback for Apple to provide a second deferral option, specific to major version numbers.

I work with the MDM protocol a lot day-to-day and have been [testing](https://micromdm.io/blog/os_update/) software updates for a while. I was even [optimistic](https://micromdm.io/blog/wwdc20-what-s-new-in-managing/) about what it would look like in Big Sur back in June. Apple had promised it's an entirely new implementation, closer to what is available on iOS and that everything would work better than before. But we're not off to a good start, and all the concerns I had for several years now are back. The [design](https://developer.apple.com/documentation/devicemanagement/commands_and_queries) of OS updates in MDM is brittle, requiring multiple remote procedure calls to accomplish something that was previously done by a few lines of shell scripting. And that would be bad on its own, but the bad design is coupled with a buggy implementation; there are many known issues, besides the one I described above. Unless something drastically changes, we're likely to see many months, if not years of software update bugs that are entirely out of our control. 

Apple was never particularly great at building systems for the enterprise. But, until recently, the Unix components were available for software developers and system administrators to work with. The end result ended up being that if you made an investment within your organization, macOS was a great experience for end-users. It took a lot of work, but it was all achievable. Today, Apple is no better at developing systems that are required in the enterprise. But the ball is entirely in their court. Apple still makes great tools for consumers, and there will be a demand from employees to provide them with Apple gear. But, if Apple can't start acting on our collective feedback, the experience of using a Mac in the workplace will quickly become unbearable to most.

