+++
title = "Five years behind, Five years ahead"
date = "2020-06-22T00:00:00+05:00"
tags = ["mdm", "wwdc", "v2"]
author = "Victor Vrantchan"
frontpage = true
+++

Apple MDM is over 10 years old today. Initially, MDM was a protocol to manage iPhones, but later expanded to the growing range of Apple devices. MicroMDM started about five years ago, when a few of us in the MacAdmin community became curious about what the MDM could do. We were intrigued by the "zero touch" deployment, which was brand new at the time. Lots of things were different. Almost nobody was using MDM and Apple Business Manager for the Mac. There were some vendors (notably JAMF) that had support, but the rest of us, including commercial vendors were making do with our own agents/scripts and imaging. [DeployStudio](https://www.deploystudio.com), [Imagr](https://github.com/imagr) and [Restor](http://github.com/google/restor) were the common tools for setting up a new Mac. The MDM documentation itself was not accessible, except to registered vendors. I [wrote](https://groob.io/posts/mdm-experiments/) a quick post in 2015 about experimenting with MDM. About six months later, the MicroMDM project was officially announced. 

Over the years, MicroMDM attracted some interest from different areas of the industry, which shows the project filling use-cases that aren't otherwise filled by other solutions. 
- A few organizations with larger fleets adopted MicroMDM to manage their Mac devices. 
- MicroMDM is the backend for multiple commercial startups. These are not necessarily the commercial products you use to manage your devices. Instead, custom workflows in education, healthcare and hospitality are being deferred to dedicated MDM services, and MicroMDM is a part of that story. 
- MicroMDM's [SCEP server](https://github.com/micromdm/scep), a project created to support PKI in configuration profiles, became one of the more popular open source SCEP implementations, with many users unrelated to Apple MDM.
- As a lightweight project, un-encumbered by the usual business logic, it's the project of choice for many to validate the Apple protocol and try out new features. Every year after WWDC, MicroMDM is one of the first projects out there to implement any newly announced features.
- The project has over 1k stars on GitHub. More importantly, every release for the last few years has new, external contributors, most of whom are new to both Go and Apple MDM. 

While I'm proud of the achievements highlighted above, and more, I have to also dwell on some things which give me anxiety as a maintainer.
- The project is still difficult to approach. A large motivation for building the project for me was to learn Go programming, and it shows. There's a general lack of documentation for developers, and some choices for writing code that haven't aged well. 
- The project [is not recommended](https://github.com/micromdm/micromdm/blob/main/docs/user-guide/introduction.md#not-a-product) for many who'd like to use it. It remains one of the few maintained open source MDM solutions, but it's not what I'd recommend as a low-cost replacement of Profile Manager. I can't tell someone to use MicroMDM without also telling them they'll eventually have to write a bunch of code on their own to make it useful.

When I reflect on MicroMDM up to now, the aspect I'm most excited by, is the pedagogical one. MicroMDM inspired at least a few people to look at Go as an option. It also made the MDM features more accessible to those who wanted to see exactly how each feature works under the hood. So if I plan on maintaining MicroMDM going forward, and remain happy doing it, that's the area I should be spending more time on.

Todays is the first day of WWDC 2020. Today is also the day we begin working on the eventual v2 of MicroMDM. With the next version, I want to keep what makes MicroMDM so unique, but also make it a Mac management solution I could recommend to any organization, large or small. I also want to focus on expanding the number of contributors and power users by keeping a development blog, regular office hours and finding new opportunities to teach.

Let's [see](https://developer.apple.com/wwdc20/) what's new in managing Apple devices this year, then come together to build something new.
