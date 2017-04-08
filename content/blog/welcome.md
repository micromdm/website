+++
author = "Victor Vrantchan"
tags = ["mdm"]
date = "2015-08-24T13:04:37+02:00"
title = "What's next for MicroMDM"
frontpage = true

+++

I first [wrote](https://groob.io/posts/mdm-experiments/) about experimenting with MDM at the end of 2015. Since then, Apple has made the MDM specification [public](https://developer.apple.com/library/content/documentation/Miscellaneous/Reference/MobileDeviceManagementProtocolRef/3-MDM_Protocol/MDM_Protocol.html), many administrators are looking to swap imaging devices for a DEP workflow, and even commercial vendors [are taking notice](https://simplemdm.com/2017/03/07/deploy-munki-apple-dep-mdm/) of the [needs](http://blog.eriknicolasgomez.com/2017/03/08/Custom-DEP-Part-1-An-Introduction/) of our community. One thing has become increasingly clear - MDM will play a critical role in the future of managing Apple devices in the enterprise. And with the renewed interest in MDM from the macadmin community, it's only fair to be asked -- What is the future of the MicroMDM project?   <!--more-->

# Growing the community

Last summer I met Jesse Petterson at the Penn State Macadmins conference. Jesse is one of the [pioneers](https://github.com/jessepeterson/commandment) in developing an open source MDM project, and has helped me figure out a lot of the initial hurdles to get MicroMDM up and running. We've had a lot of opportunities to chat about [what features we'd like to see](https://github.com/micromdm/micromdm/issues/110) in an MDM server, and how to make the project easy to get started for new users. 

Jesse and I will be speaking about MDM and DEP at a number of Mac conferences this year, starting with [macdevops::YVR](https://www.macdevops.ca/speakers/) in June. 

If you're interested in the project, follow [@micromdm_io](https://twitter.com/micromdm_io) on Twitter or join the discussion on [Slack](https://macadmins.herokuapp.com/).

# Focus on user needs

Long term, we're looking to make MicroMDM the solution of choice for Apple device management, focusing on usability, extensibility and security. To achieve this goal, MicroMDM development will focus on actual user needs. The immediate focus will be building our integration with [DEP](https://deploy.apple.com) and allowing administrators to provision macOS devices. We're looking at what to build next. If you have a specific need that would help you adopt MicroMDM, consider opening an [issue](https://github.com/micromdm/micromdm/issues/new) or sending an email to [hello@micromdm.io](mailto:hello@micromdm.io).

