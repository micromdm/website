+++
title = "Updating to Mojave?"
date = "2018-10-16T13:04:37+02:00"
tags = ["mdm", "dep", "os_update"]
author = "Victor Vrantchan"
frontpage = true
+++

Earlier this year Apple updated the [MDM Protocol Reference](https://developer.apple.com/enterprise/documentation/MDM-Protocol-Reference.pdf) document to add a previously undocumented key to the `AvailableOSUpdates` command. 

| Key                 | Content                                                                       |
| ------------------- | ----------------------------------------------------------------------------- |
| **IsMajorOSUpdate** | Set to true if this is a major OS update (e.g. 10.13.x to 10.14). macOS only. |

This got a few of us in the `#mdmdev` channel on MacAdmins Slack talking. We've tried scheduling OS Updates in the past, and MicroMDM has had support for all these commands for years. The user experience around them is not as great as it could be (more on this some time) so the feature doesn't get much use. But this update was intriguing. The documentation clearly stated that this would be for a _major_ os upgrade, but every time we ran the command, even on older systems, the only available options were point releases (10.13.5 to 10.13.6 for example). No matter what, the update was not listed.  

Mojave was released, but we still didn't see anything. Someone speculated that Apple might have forgotten to add the necessary update to the catalog, and that's why we weren't seeing it. I ended up filing a bug report in early October asking for clarification on the documentation. Was I doing something wrong? Was the documentation mistaken? 
Nothing really worked despite me trying this on 10.12, 10.13 and 10.14 betas. Until today... 

{{<figure src="/os_update/os_update.png" title="" class="screenshot" >}}

Today, a new product key was added to the macOS Software Update catalog. 

```
<key>041-14451</key>
<dict>
    <key>ServerMetadataURL</key>
    <string>http://swcdn.apple.com/[...]/macOSInstallerNotification_GM.smd</string>
    <key>Packages</key>
    <array>
        <dict>
            <key>Digest</key>
            <string>f11446026976c362a451b89d354c2313af19a1a3</string>
            <key>Size</key>
            <integer>1821971</integer>
            <key>MetadataURL</key>
            <string>https://swdist.apple.com/[...]/macOSInstallerNotification_GM.pkm</string>
            <key>URL</key>
            <string>http://swcdn.apple.com/[...]/macOSInstallerNotification_GM.pkg</string>
        </dict>
    </array>
    <key>PostDate</key>
    <date>2018-10-16T20:05:00Z</date>
```

The `.pkg` file includes signed `OSXNotification.bundle` with Info.plist keys like this:
```
<key>FreeSpaceRequired</key>
<integer>12000000000</integer>
<key>HumanReadableName</key>
<string>macOS Mojave</string>
<key>ItemID</key>
<integer>1398502828</integer>
<key>ProductBuildVersion</key>
<string>18A391</string>
```

So this is a configuration update which includes information about the Mojave installer, with both the product version and the Mac App Store ItemID. Neat. Normally you don't see the full OS installer when running `softwareupdate`. It's in the App store and you have to download it like you would download an app. This must be a shim Apple is using to enable a MDM only feature. 

A curious side effect of having the notification bundle come from the Software Update catalog, is that it's repsonsible for this notification, which shows up on user machines regardless of whether you have devices enrolled in MDM:

{{<figure src="/os_update/os_update_prompt.png" title="" class="screenshot" >}}


Let's see what happens when we follow the update steps via MDM. 
Note, to spare you a lot of XML, I'll be referencing a few shorthand commands for interacting with the MDM. These are some simple wrappers that create cURL requests to the MicroMDM API, and map directly to the documented commands in the spec. You can [find them in the repo](https://github.com/micromdm/micromdm/tree/a230831725f173f779cf7e030a93a71b32e8ab13/tools/api), and combined wth [ngrok](https://github.com/micromdm/micromdm/tree/a230831725f173f779cf7e030a93a71b32e8ab13/tools/ngrok), you can test everything yourself without a production MDM setup. 

1. ScheduleOSUpdateScan -- initiates a software update check. 
```
./tools/api/commands/schedule_os_update_scan $udid true
```

2. AvailableOSUpdates -- once the scan is completed, we can see available updates.
```
./tools/api/commands/available_os_updates $udid
```

On a 10.12.6 laptop we see the following:

```
<dict>
    <key>AllowsInstallLater</key>
    <true/>
    <key>AppIdentifiersToClose</key>
    <array/>
    <key>HumanReadableName</key>
    <string>macOS Installer Notification</string>
    <key>HumanReadableNameLocale</key>
    <string>en</string>
    <key>IsConfigDataUpdate</key>
    <true/>
    <key>IsCritical</key>
    <false/>
    <key>IsFirmwareUpdate</key>
    <false/>
    <key>MetadataURL</key>
    <string>http://swcdn.apple.com/[...]/macOSInstallerNotification_GM.smd</string>
    <key>ProductKey</key>
    <string>041-14451</string>
    <key>RestartRequired</key>
    <false/>
    <key>Version</key>
    <string>2.0</string>
</dict>
```

Bingo!

Now we can schedule a OS Update

3. ScheduleOSUpdate for a product key, InstallASAP is the most aggressive install type. 
```
./tools/api/commands/schedule_os_update $udid 041-14451 InstallASAP
```

4. Issue another Scan, followed by AvailableOSUpdates
```
./tools/api/commands/schedule_os_update_scan $udid true
./tools/api/commands/available_os_updates $udid
```

Now we see a new item. 

```
<dict>
    <key>DownloadSize</key>
    <real>12000000000</real>
    <key>HumanReadableName</key>
    <string>macOS Mojave</string>
    <key>HumanReadableNameLocale</key>
    <string>en</string>
    <key>IsConfigDataUpdate</key>
    <false/>
    <key>IsCritical</key>
    <false/>
    <key>IsFirmwareUpdate</key>
    <false/>
    <key>IsMajorOSUpdate</key>
    <true/>
    <key>ProductKey</key>
    <string>_OSX_18A391</string>
    <key>RestartRequired</key>
    <true/>
    <key>Version</key>
    <string>18A391</string>
</dict>
```

Very strange ProductKey but I'll take it. 

```
./tools/api/commands/schedule_os_update $udid _OSX_18A391 InstallASAP
```

```
<key>ErrorChain</key>
<array>
    <dict>
        <key>ErrorCode</key>
        <integer>74</integer>
        <key>ErrorDomain</key>
        <string>MDMClientError</string>
        <key>LocalizedDescription</key>
        <string>Command requires DEP enrollment: ScheduleOSUpdate &lt;MDMClientError:74&gt;</string>
    </dict>
</array>
<key>Status</key>
<string>Error</string>
```

Did you catch that? User Approved MDM is not enough. We need a [DEP enrollment](https://bugreport.apple.com/web/). No problem... a few seconds later. 

```
./tools/api/commands/schedule_os_update $udid _OSX_18A391 NotifyOnly
```

```
<key>ErrorChain</key>
<array>
    <dict>
        <key>ErrorCode</key>
        <integer>12008</integer>
        <key>ErrorDomain</key>
        <string>MCMDMErrorDomain</string>
        <key>LocalizedDescription</key>
        <string>Unsupported InstallAction for this ProductKey</string>
    </dict>
</array>
```

Ah sorry for that detour. I wanted to see if `DownloadOnly` and `NotifyOnly` or `InstallLater` are supported actions. Nope!

But it's ok:

```
./tools/api/commands/schedule_os_update $udid _OSX_18A391 InstallASAP

<key>UpdateResults</key>
<array>
    <dict>
        <key>InstallAction</key>
        <string>InstallASAP</string>
        <key>ProductKey</key>
        <string>_OSX_18A391</string>
        <key>Status</key>
        <string>Installing</string>
    </dict>
</array>
```

Except nothing happens. I can confirm on the Mac there is no network activity. `AvailableOSUpdates` again. We get:

```
<dict>
    <key>DownloadPercentComplete</key>
    <real>0.0</real>
    <key>IsDownloaded</key>
    <false/>
    <key>ProductKey</key>
    <string>_OSX_18A391</string>
    <key>Status</key>
    <string>Idle</string>
</dict>
```

I tried this routine about 20 times and nothing was happening. Absolute madness. And then it occured to me. What if this works like iOS and you need to be plugged in. Sure enough. As soon as I plugged the power adapter in, `softwareupdated` began downloading _something_. `<real>12000000000</real>` bytes later I was prompted with a notification: 

> Update Requested
> A new update was requested to be installed by an administrator. 

And a single availabe option `Restart`. I clicked it, but nothing happened. So I checked `AvailableOSUpdates`:

```
<dict>
	<key>DownloadPercentComplete</key>
	<real>1</real>
	<key>IsDownloaded</key>
	<true/>
	<key>ProductKey</key>
	<string>_OSX_18A391</string>
	<key>Status</key>
	<string>Idle</string>
</dict>
```

Another `ScheduleOSUPdate _OSX_18A391`. Another notification to Restart. Clicking the button does nothing. I decided to look at logs:

```
2018-10-16 21:09:29.624110-0400 0x51e4     Default     0x8000000000002f1f   112    authd: Succeeded authorizing right 'com.apple.ServiceManagement.daemons.modify' by client '/System/Library/PrivateFrameworks/CommerceKit.framework/Versions/A/Resources/storeinstalld' [596] for authorization created by '/System/Library/PrivateFrameworks/CommerceKit.framework/Versions/A/Resources/storeinstalld' [596] (12,0)
2018-10-16 21:09:29.624497-0400 0x51f5     Default     0x8000000000002f1f   596    storeinstalld: (StoreFoundation) [com.apple.commerce.CKLegacy] MajorOSInstallDaemonDelegate: Non-interactive authorization for software install was successful
2018-10-16 21:09:29.624836-0400 0x51f5     Default     0x8000000000002f1f   596    storeinstalld: (StoreFoundation) [com.apple.commerce.CKLegacy] MajorOSInstallDaemonDelegate: Triggering OS installation via worker connection <NSXPCConnection: 0x7fb93ef11300> connection from pid 448 (uid: 0, gid: 0)
2018-10-16 21:09:29.625272-0400 0x4e80     Default     0x8000000000002f1f   719    storeinstallagent: (StoreFoundation) [com.apple.commerce.CKLegacy] MajorOSInstallController: Initialized with IA path: /Applications/Install macOS Mojave.app
2018-10-16 21:09:29.625356-0400 0x4e80     Default     0x8000000000002f1f   719    storeinstallagent: (StoreFoundation) [com.apple.commerce.CKLegacy] MajorOSInstallController: Starting configure
2018-10-16 21:09:29.625902-0400 0x4f83     Default     0x8000000000002f1f   719    storeinstallagent: (StoreFoundation) [com.apple.commerce.CKLegacy] MajorOSInstallController: Setting OSISTarget
2018-10-16 21:09:29.629410-0400 0x4e80     Default     0x0                  719    storeinstallagent: (StoreFoundation) [com.apple.commerce.CKLegacy] MajorOSInstallController: helperToolDied
```

And then, just as I was about to copy these logs to my other computer, my mac rebooted. This was about 2 minutes after clicking on the Restart button. 
I expected to get prompted with a install wizard, similar to going through the `Install macOS Mojave.app` process and having to click next a few times. To my surprise, I got the standard black screen with the Apple logo and a 37 minute progress indicator. 18 minutes after that I was staring at a Mojave login prompt. Hello Dark Mode. 
