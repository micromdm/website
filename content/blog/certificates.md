+++
frontpage = true
author = "Jesse Peterson"
tags = ["mdm", "certificate", "apple"]
date = "2017-05-11T13:04:37+02:00"
title = "Understanding MDM Certificates"
lastmod = "2019-12-19T13:03:00+08:00"

+++

Apple Mobile Device Management (MDM) requires the use of various digital certificates for its operation. But exactly which certificates and the various ways in which they are generated, acquired, signed, used, exported, imported, and managed within an MDM product may not be so clear. Generally speaking a commercial MDM product or service manages most of the complexity related to these certificates for you but in the case of an open source MDM much of that responsibility will land on you. In this post I hope to bring a better understanding of these certificates with the aim that you'll be managing at least a few of them yourself.

<!--more-->

This post assumes you have at least a basic familiarity of working with certificates, private keys, certificate authorities and the like. To learn more about these and the TLS/SSL protocols have a look at this [TLS & Certificate Survival Guide](http://www.zytrax.com/tech/survival/ssl.html).

# Apple Push Notification service certificate

Perhaps the most important certificate is the Apple Push Notification service (APNs) certificate (or just "push certificate"). Without going into too much detail about the how the [MDM protocol works](https://developer.apple.com/business/documentation/MDM-Protocol-Reference.pdf) exactly suffice it to say that the [Apple Push Notification service](https://developer.apple.com/library/content/documentation/NetworkingInternet/Conceptual/RemoteNotificationsPG/APNSOverview.html) (APNs for short) is crucial to how MDM operates. Push notifications are the only way to have a device talk to your MDM after initial enrollment.

However not just anybody can send these push notifications to Apple (and in turn to your devices). Only services that have acquired a special Apple-signed certificate for sending these push notifications are allowed. While APNs push notifications for MDM are *similar* to "normal" app-based APNs push notifications they are slightly different and are acquired in different vastly ways.

## Getting an MDM APNs certificate

There are a few different options for acquiring an MDM push certificate and each has its own pros and cons.

### Method A: Roll your own using an MDM Vendor Certificate

An *MDM Vendor Certificate* (or MDM CSR certificate) is a special certificate that can sign *other MDM APNs certificate requests*. These requests can then be submitted to Apple for signing to get the MDM push certificate. This certificate option is only available to members of the [Apple Developer Enterprise Program](https://developer.apple.com/programs/enterprise/) which has a cost of US$300/year. Victor Vrantchan talks about signing up for this account (and this specific certificate option) in [this blog post](/blog/accounts/).

This method is probably the most complicated option (due to the two-step nature of the certificates) and certainly the most expensive. However a developer account is a good resource for any Mac admin anyway, as Victor mentions, and this is probably the most Apple-supported option for hosting your own MDM as you're not beholden to any other entity or software than Apple for your certificates. This is the same method a commercial MDM vendor uses to generate certificates for their customers (and hence explains the two-step signing process — typically an actual MDM user/admin wouldn't be doing the "Vendor" steps).

The high-level overview of getting and using an MDM push certificate with a Vendor certificate is similar to this:

{{<figure src="/certificates/mdm_csr_process.png" title="Diagram showing certificate signing workflow for MDM Vendor Certificates" class="graphic" >}}

1. Sign up for an Enterprise account and request the *MDM CSR* certificate on the account (see [this blog post](/blog/accounts/))
2. Create a private key and an MDM CSR Vendor Certificate in the Apple portal
3. Once you have the MDM CSR Vendor certificate now generate another, separate "customer" or "end-user" certificate request.
4. This "customer" certificate request needs to be wrapped up and signed into a special format by the MDM Vendor certificate.
5. This signed request must then be uploaded to [identity.apple.com](https://identity.apple.com/) where Apple will issue the final actual push certificate that is used in conjunction with the private key for sending APNs MDM push notifications.

To see a walk-through on this process take a look at [this presentation on MicroMDM](https://www.youtube.com/watch?v=WGKT-PyHz6I&t=26m55s) starting at 26m55s. This process is also covered here in [Pepijn Bruienne's blog post](http://enterprisemac.bruienne.com/2015/06/06/mdm-azing-setting-up-your-own-mdm-server/).  As well for actually working with these certificates [micromdm's mdmctl mdmcert tool](https://github.com/micromdm/micromdm/wiki/Generating-MicroMDM-MDM-Certificates) has documentation on getting started once all the Apple account details have been taken care of.

Once the final "customer" push certificates are created and downloaded from Apple (with e.g. *mdmctl mdmcert* & [identity.apple.com](https://identity.apple.com/)) they can be used directly with MicroMDM. 

### Method B: Export a Profile Manager certificate

*Profile Manager* is Apple's reference MDM product (or proof-of-concept, depending on how jaded you are about it). It's bundled with [macOS Server.app](https://itunes.apple.com/us/app/macos-server/id883878097?mt=12) for US$20 and requires macOS (or a VM running it). It has a neat feature where with only an Apple ID it can submit a certificate request to Apple, have an MDM certificate signed, and returned back to it in one step. This skips over a bunch of the rigmarole in getting a push certificate.

This technique was probably first documented in 2011 as a part of [David Shuetz Black Hat 2011 presentation "Inside Apple’s MDM Black Box."](https://www.youtube.com/watch?v=OifARLlRMyU). [Page seven of his PDF](https://media.blackhat.com/bh-us-11/Schuetz/BH_US_11_Schuetz_InsideAppleMDM_WP.pdf) documents exporting the PKCS#12 certificate & key from Keychain Access.app once you've [turned on Profile Manager](https://help.apple.com/serverapp/mac/5.3/#/apd05B9B761-D390-4A75-9251-E9AD29A61D0C). The [MicroMDM Quickstart guide](https://github.com/micromdm/micromdm/wiki/Quickstart#getting-an-mdm-push-certificate) also has some documentation on getting at this certificate.

If you just want to try out MicroMDM (or another Open Source MDM solution) this is likely the easiest and quickest way to get an APNs push certificate. The downside is that it is an Apple-proprietary method of getting this certificate that's embedded inside Profile Manager. This probably makes it of questionable legality to use with anything other than Profile Manager. As well there is a nominal cost associated.

Once the Profile Manager certificate is exported as a `.p12` file it can be used directly with MicroMDM.

### Method C: Sign up for mdmcert.download

[mdmcert.download](https://mdmcert.download/) is a service created to issue MDM push certificates to organizations wishing to run open-source MDM solutions. The certificates are free of cost but per Apple, only *organizations* (and not individuals) may agree to request a certificate. Apple also requires gathering some information like business name, email addresses, etc. That may not be something you're willing to share or legally able to do for your organization.

That said it offers a method to get push certificates that's easier than method A, above, but isn't quite as easy as the Server.app method. As mentioned above it's free, too.

MicroMDM's [mdmctl mdmcert.download option](https://github.com/micromdm/micromdm/wiki/mdmcert.download) can be used to request an mdmcert.download APNs certificate when following the [mdmcert.download instructions](https://mdmcert.download/instructions). Once mdmctl decrypts the encrypted CSR request (which is then subsequently uploaded to [identity.apple.com](https://identity.apple.com/) in order to retrieve the certificate) the push certificates you download can be directly used with MicroMDM.

## Push certificate gotchas

Once you have an MDM APNs certificate you have the ability to send push notifications to devices that are enrolled in your MDM. But there are a couple of caveats that you want to keep in mind:

* MDM APNs certificates **expire yearly.** This of course means you'll need to renew the certificate with a similar process you followed to get the original certificate. Note that it **must be a renewal** and not a *new* APNs push certificate for a very important reason:
* The APNs Push "topic" (which is embedded in the push certificate and has a `com.apple.mgmt.` prefix for MDM) **can never change during the life of a device's enrollment.** Basically, this means a couple things:
  * You can't just use *any* MDM APNs certificate; they're not interchangeable. When devices enroll into an MDM they are tied to that particular APNs certificate push topic.
  * This is why a *new* certificate can't be used as a renewal — said new certificate would have a different Push topic and would not be able to be used for sending push notifications to your existing enrolled devices. When you renew the certificate on [identity.apple.com](https://identity.apple.com/) it must be submitted as a Renewal to a previously created certificate. This doesn't mean you have to use the same private key (which is bad practice), just that it is submitted as a renewal to Apple.
  * If you're using the Server.app extraction method you'll want to keep that instance of Server.app (computer/VM) around so that you can fire it up to get the push certificate with the same topic renewed and exracted again. But hopefully you're only using that method as a test and not necessarily as your production MDM certificate.
* MDM APNs certificates are tied to an Apple ID — the account that is used to issue the APNs Push certificate on [identity.apple.com](https://identity.apple.com/). Because this is the only account that can issue renewals for your existing enrollments you'll want to keep track of its login credentials. Gentle reminder to always use 2FA as well.
* As with any certificate take note of what format the APNs certificate & key are in. Is it two separate PEM-encoded certificate and key files? Is it a single PKCS#12 file encrypted with a password? Knowing these details will be crucial to getting MicroMDM (or any MDM) up and running with them.

# Device Identity certificate

An MDM device enrolls into MDM with an identity certificate & key pair. This certificate and private key can either be a) outright *given* to the device or b) the device can request that a new certificate be signed on its behalf.

The former is done by embedding a PKCS#12 profile payload in the enrollment profile. This will become the device's identity certificate. The latter is done by the device by itself using the SCEP protocol to request a certificate to be signed. Using the SCEP protocol, for all its faults, is much more secure for a few reasons:

1. The private key is not transferred over the network. Using the 'embedded' method it is.
2. The private key is not in an enrollment profile on disk in a Download folder which may be inadvertantly exported, shared, or lost. If the identity is compromised you can essentially spoof the device connecting to the MDM server.
3. Each device has a unique certificate and private key generated for it to use.

The primary use of this device identity certificate is to authenticate the device to the MDM server whenever an HTTPS connection is made. Depending on how the MDM software and enrollment profile are configured the device either performs [TLS/SSL client authentication](https://en.wikipedia.org/wiki/Client_certificate) or it provides a [CMS](https://en.wikipedia.org/wiki/Cryptographic_Message_Syntax) detached signature of the MDM request using its identity key pair. In either case, this cryptographically "proves" that the device is using a certificate that *should* be already known to the MDM server and only belonging to that device.

In the case of SCEP you likely will never have to touch the identity certificate manually. It automatically enrolls and references the correct certificate. Older versions of MicroMDM and e.g. Project-iMAS required manually providing the device identity certificate & key pair in the enrollment profile but that's no longer needed, or recommended.

It is also possible to use this certificate to encrypt profiles to a device (and only that specific device in the case of unique per-device certificates). More on that later.

## Device Identity Certificate Authority

When using SCEP the device will be issued a certificate from a Certificate Authority (CA). While embedded-profile device identity certificates can also be issued from a CA (this is what [Commandment MDM](https://github.com/jessepeterson/commandment) does for example) they're likely just self-signed. Having a CA implies a certificate chain and associated trust concerns they bring. While not something you may be hands-on managing for MDM these are a few questions and concerns you'll want to think about when deploying your MDM:

* Where are your device identity certificates getting issued from, which CA?
* How do you trust this CA? Who owns/operates it?
* How do you *revoke* certificates for compromised devices or lost certificate & private key pairs?
* How is the MDM trusting the provided device certificate? Is it merely trusting that the *issuer* issued this certificate? Or is it verifying the contents of the device identity certificate itself?

In MicroMDM the SCEP CA is built-in using the [micromdm/scep](https://github.com/micromdm/scep) project. On each MDM check-in and MDM command it verifies that the certificate was issued by the built-in SCEP CA. Even with the aid of well-designed and straightforward SCEP systems that work with minimal configuration, it is good to understand where device certificates are issued from, how they actually get issued, and how to verify they are trusted. MicroMDM also [validates that the enrolled certificate matches the UDID of the device](https://github.com/micromdm/micromdm/pull/429) so that devices can only access their own command queue.

# Configuration Profile signing & encrypting certificates

Apple's [Configuration Profiles](https://developer.apple.com/business/documentation/Configuration-Profile-Reference.pdf) can be CMS signed and/or encrypted. To do that, you must encrypt them with a public key or sign them with a certificate, respectively.

Encrypting a Configuration Profile requires using a public key that the device has the private key to. This can be known through sending a SCEP profile, a certificate profile payload, or just using the device's enrollment identity already on the device once it's enrolled. When the device has the corresponding private key, it can decrypt the encrypted profile and install it. In this way, you can use a device's specific identify certificate to encrypt a profile so that *only the target device* can decrypt them.

Signing a Configuration Profile is also possible. For profile signing to be effective the profile should be signed by a certificate that the device trusts. This can be a certificate in the device's [trusted root store](https://en.wikipedia.org/wiki/Public_key_certificate#Root_programs) (similar to a browser's trusted root store) or it can be a certificate that the device is separately configured to trust. In the case of MDM we likely already have a few certificates that are necessarily trusted that we can use. For example the HTTPS web certificate of the MDM may either be trusted in the system root store or be configured as a part of the enrollment profile (say, for a self-signed certificate, more on that later). This may be used to sign profiles (or packages). Often folks will simply sign profiles with their Apple Developer certificate because the signer of those certificates exists in the trusted root store of the device and will just work without further changes or trust management.

Seeing as the trusted certificates on a device may be known to the MDM system (e.g. coming through in the enrollment profile) the possibility exists for the MDM system to sign profiles either on command or in an on-the-fly fashion. MicroMDM now has support profile for [signing a profile](https://github.com/micromdm/micromdm/releases/tag/v1.6.0) but not encrypting.. Pre-encrypted or signed profiles should also work.

# Configuration Profile trusted certificates

Configuration Profiles are able to add certificates to the trusted store of certificates of a device by using Configuration Profile payloads for certificates. Not only is this valuable in-and-of-itself for normal system administration tasks but also importantly for MDM these certificates can be embedded in the *enrollment profile* such that those certificates will be trusted by the system when the device is enrolled.

This has implications for MDM operation *especially* in the case of using [self-signed certificates](https://en.wikipedia.org/wiki/Self-signed_certificate) for the HTTPS server. By default a self-signed certificate on a normal website would simply not be trusted by a device and it's no different for an MDM server. However if we place the MDM server's self-signed HTTPS certificate inside the enrollment profile then the system will trust it and MDM operation can commence after enrollment. The same goes for profile signing mentioned above. However, this may not work on iOS devices.

# HTTPS certificate

The web server portion of MDM requires HTTPS. This of course implies [TLS/SSL](https://en.wikipedia.org/wiki/Transport_Layer_Security) certificates and the related trust issues they bring. That said an MDM server's use of TLS/SSL certificates isn't all that different from your typical web server's TLS/SSL configuration.

This means that the types of a certificates you can use are pretty much the same as you can use for a website. These include:

* A purchased TLS/SSL certificate from a reputable Certificate Authority provider like e.g. [Comodo](https://www.comodo.com/), [GlobalSign](https://www.globalsign.com/en/), [Entrust](https://www.entrust.com/), etc.
* MicroMDM supports [LetsEncrypt](https://letsencrypt.org/) (LE) certificate acquisition for automatic SSL configuration
  * Note per usual LE operation this requires public internet inbound TCP port 443 access to your MDM server and properly working public DNS
  * MicroMDM's use of Let's Encrypt is currently broken, but there are [workarounds](https://github.com/micromdm/micromdm/wiki/Generating-LetsEncrypt-Certs-with-the-Lego-Client)
* Using a [self-signed](https://en.wikipedia.org/wiki/Self-signed_certificate) (or private CA-signed) certificate

If self-signed or private CA-signed certificates are used then we must add those certificates to our enrollment profile using the appropriate configuration profile payloads as mentioned above. MicroMDM does this for us.

MDM is intended for mobile devices and as such there's a reasonable assumption it will be running in a place that is accessible from the public internet. While this is not a technical requirement — an MDM server can run behind a firewall or in a private network as long as APNs push notifications can be sent and received — a lot of the more interesting features of MDM like, say, Remote Wipe lose a lot of their luster without it.

So by far the easiest TLS/SSL configuration is done with LetsEncrypt. Just make sure your server is publicly accessible on port 443 at a domain name you control. MicroMDM should take care of the rest. That said self-signed & purchased SSL options are supported as well in MicroMDM.

# DEP token certificate

Your MDM server talks to Apple's [Device Enrollment Program](https://www.apple.com/business/dep/) API using OAuth tokens. However before you gain access to these tokens you have to complete a [PKI](https://en.wikipedia.org/wiki/Public_key_infrastructure) process where you upload a certificate (which contains a public key) that Apple will use to encrypt the tokens with. Once you've given Apple your public key you can download the encrypted tokens and subsequently decrypt them and start using them to connect to the DEP API (to e.g. fetch & sync devices, configure a DEP profile, etc.). This process is managed via the [Apple Business Manager](https://support.apple.com/guide/apple-business-manager/welcome/web) (ABM) interface once your organization has enrolled in ABM and of course have purchased devices that are registered in your ABM account. If you're a school there is the very similar Apple School Manager (ASM).

Like with the APNs certificate these DEP tokens (like the VPP tokens) expire yearly and this process must be done again to renew them.

# DEP anchor certificates

Remember how we said we could use self-signed HTTPS web certificates? And  that we need to embed the trust information for the HTTPS certificate in the enrollment profile? For a manual enrollment (that is where to say where you visit the profile website manually) this is usually not a problem because the *initial* certificate prompt that one will get (like with any self-signed website) can be simply overridden by the admin doing the enrollment. However, what about a DEP MDM enrollment where there is no user doing the enrollment?

That's what the DEP anchor certificates are for. There's a special property on the DEP profile just for this scenario. Note that DEP profiles should not be confused with Configuration Profiles — DEP profiles are a completely separate JSON structure. This property is called `anchor_certs` in the DEP profile and it allows you to specify trusted certificates to Apple's DEP API. This property will instruct a device at DEP enrollment time to trust the given certificates when connecting to the MDM server over TLS/SSL.

# Conclusion

My hope is this overview has been helpful in untangling the various certificates used in the MDM protocol and perhaps sheds some light on some of the nuances and gotchas surrounding them.

But even so if you've still got questions or are having trouble working with any of these certificates come join us in the MDM-related [MacAdmins Slack channels](/) where we discuss topics like this.
