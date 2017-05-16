+++
frontpage = true
author = "Jesse Peterson"
tags = ["mdm", "certificate", "apple"]
date = "2017-05-11T13:04:37+02:00"
title = "Understanding MDM Certificates"

+++

MDM requires the use of various digital certificates for its operation. However exactly which certificates and the various ways in which they are generated, acquired, signed, used, exported, imported, and managed within an MDM product may not be so clear. Generally speaking a commercial MDM product or service manages most of the complexity related to these certificates for you but in the case of an open source MDM much of that responsibility will land on you. In this post I hope to bring a better understanding of these certificates with the aim that you'll be managing at least a few of them yourself.

<!--more-->

This post assumes you have at least a basic familiarity of working with certificates, private keys, certificate authorities and the like. To learn more about these and the TLS/SSL protocols have a look at this [TLS & Certificate Survival Guide](http://www.zytrax.com/tech/survival/ssl.html).

# Apple Push Notification service certificate

Perhaps the most important certificate is the Apple Push Notification service certificate (or just "push certificate"). Without going into too much detail about the how the [MDM protocol works](https://developer.apple.com/library/content/documentation/Miscellaneous/Reference/MobileDeviceManagementProtocolRef) exactly (come see our talks this year at [various](https://www.macdevops.ca/) [Mac](http://macadmins.psu.edu/) [Admin](http://www.macsysadmin.se/) conferences!) suffice it to say that the [Apple Push Notification service](https://developer.apple.com/library/content/documentation/NetworkingInternet/Conceptual/RemoteNotificationsPG/APNSOverview.html) (APNs for short) is crucial to how MDM operates. In short: once a device is enrolled in an MDM then push notifications are the only way to have that device actually talk to your MDM on an ongoing basis.

However not just anybody can send these push notifications to Apple (and in turn to your devices). Only services that have acquired a special Apple-signed certificate for sending these push notifications are allowed. While APNs push notifications for MDM are *similar* to "normal" app-based notifications they are slightly different and are acquired in different ways.

## Getting an MDM APNs certificate

There's a few different options for acquiring an MDM push certificate and each comes with its own pros and cons.

### Method A: Roll your own using an MDM Vendor Certificate

An *MDM Vendor Certificate* (or MDM CSR certificate) is a special certificate that can sign *other MDM APNs certificate requests* that can then subsequently be submitted to Apple for signing. This certificate option is only available to members of the [Apple Developer Enterprise Program](https://developer.apple.com/programs/enterprise/) which costs US$300/year. Victor Vrantchan talks about signing up for this account (and this specific certificate option) in [this blog post](/blog/accounts/).

This method is probably the most complicated option (due to the two-step nature of the certificates) and certainly the most expensive. However a developer account is a good resource for any Mac admin anyway, as Victor mentions, and this is probably the most Apple-supported option for hosting your own MDM as you're not beholden to any other entity or software than Apple for your certificates. This is the same method a commercial MDM vendor uses to generate certificates for their customers (and hence explains the two-step signing process — typically an actual MDM user/admin wouldn't be doing the "Vendor" steps).

The high-level overview of getting and using a push certificate with a Vendor certificate is similar to this:

1. Sign up for an Enterprise account and request the `MDM CSR` certificate on the account (see [this blog post](/blog/accounts/))
1. Create a private key and an MDM CSR Vendor Certificate in the Apple portal
2. Once you have the MDM CSR Vendor certificate now generate another, separate "customer" or "end-user" certificate request.
3. This "customer" certificate request needs to be wrapped up and signed into a special format by the MDM Vendor certificate.
4. This signed request must then be uploaded to [identity.apple.com](https://idmsa.apple.com/) where Apple will issue the final actual push certificate that be used in conjunction with the private key for sending APNs MDM push notifications.

This process is also covered here in [Pepijn Bruienne's blog post](http://enterprisemac.bruienne.com/2015/06/06/mdm-azing-setting-up-your-own-mdm-server/). As well for actually working with these certificates [micromdm's certhelper tool](https://github.com/micromdm/tools) has documentation on getting started once all the Apple account details have been taken care of.

Once the final "customer" push certificates are created and exported from Apple (with e.g. [certhelper](https://github.com/micromdm/tools) & identity.apple.com) they can be used directly with MicroMDM. 

### Method B: Export a Profile Manager certificate

*Profile Manager* is Apple's reference MDM product (or proof-of-concept, depending on how jaded you are about it). It's bundled with [macOS Server.app](https://itunes.apple.com/us/app/macos-server/id883878097?mt=12) for US$20 and requires macOS (or a VM running it). It has a neat feature where, with only an Apple ID it can submit a certificate request to Apple, have an MDM certificate signed, and returned back to it in one step. This skips over a bunch of the rigmarole in getting a push certificate.

This technique was probably first documented in 2011 as a part of [David Shuetz Black Hat 2011 presentation "Inside Apple’s MDM Black Box."](https://www.youtube.com/watch?v=OifARLlRMyU). [Page seven of his PDF](https://media.blackhat.com/bh-us-11/Schuetz/BH_US_11_Schuetz_InsideAppleMDM_WP.pdf) documents exporting the PKCS#12 certificate & key from Keychain Access.app once you've [turned on Profile Manager](https://help.apple.com/serverapp/mac/5.3/#/apd05B9B761-D390-4A75-9251-E9AD29A61D0C). The [MicroMDM Quickstart guide](https://github.com/micromdm/micromdm/wiki/Quickstart#getting-an-mdm-push-certificate) also has some documentation on getting at this certificate.

If you just want to try out MicroMDM (or another Open Source MDM solution) this is likely the easiest and quickest way to get an APNs push certificate. The downside is that it is an Apple-proprietary method of getting this certificate that's embedded inside Profile Manager. This probably makes it of questionable legality to use with anything other than Profile Manager. As well there is a nominal cost associated.

Once the Profile Manager certificate is exported as a `.p12` file it can be used directly with MicroMDM.

### Method C: Sign up for mdmcert.download

[mdmcert.download](https://mdmcert.download/) is a service created to issue MDM push certificates to organizations desiring to run open-source MDM solutions. The certificates are free of cost but per Apple only *organizations* (and not individuals) may request a certificate. A part of that requirement is the need to gather some information like business name, email addresses, etc. and of course agree to some terms.

That said it offers what I would call an easier-than-using-MDM Vendor certificate method to get push certificates, but not quite as easy as the Server.app method. And of course it's free.

MicroMDM's [certhelper](https://github.com/micromdm/tools) tool can be used to request an mdmcert.download APNs certificate when following the [mdmcert.download instructions](https://mdmcert.download/instructions). Once certhelper decrypts the encrypted CSR request (and is subsequently uploaded to [identity.apple.com](https://idmsa.apple.com/) and then the certificate downloaded) the push certificates can be directly used with MicroMDM.

## Push certificate gotchas

Once you have an MDM APNs certificate you may then send push notifications to devices that are enrolled in your MDM. But there are a couple of gotchas and things you want to take note of or keep in mind:

* MDM APNs certificates **expire yearly.** This of course means you'll need to renew the certificate with a similar process you followed to get the original certificate. Note that it **must be a renewal** and not a *new* APNs push certificate for a very important reason:
* The APNs Push "topic" (which is embedded in the push certificate and has a `com.apple.mgmt.` prefix for MDM) **can never change for the life of a device's enrollment.** Basically, this means a couple things:
  * You can't just use *any* MDM APNs certificate; they're not interchangeable. When devices enroll into an MDM they are tied to that particular APNs certificate push topic.
  * This is why a *new* certificate can't be used as a renewal — said new certificate would have a different Push topic and would not be able to be used for sending push notifications to your existing enrolled devices. When you renew the certificate on identity.apple.com it must be submitted as a Renewal to a previously created certificate. This doesn't mean you have to use the same private key (which is bad practice), just that it is submitted as a renewal to Apple.
  * If you're using the Server.app extraction method you'll want to keep that instance of Server.app (computer/VM) around so that you can fire it up to get the push certificate with the same topic renewed and exracted again. But hopefully you're only using that method as a test and not necessarily as your production MDM certificate.
* As with any certificate take note of what format the APNs certificate & key are in. Is it two separate PEM-encoded certificate and key files? Is it a single PKCS#12 file encrypted with a password? Knowing these details will be crucial to getting MicroMDM (or any MDM) up and running with them.

# Device Identity certificate

An MDM device enrolls with an identity certificate & key pair. This certificate and private key can either be outright *given* to the device or the device can request that a new certificate be signed on its behalf.

The former is done by embedding a PKCS#12 profile payload in the enrollment profile. This will become the device's identity certificate. The latter is done by the device by itself using the SCEP protocol to request a certificate to be signed. Using the SCEP protocol, for all its faults, is much more secure for a few reasons:

1. The private key is not transferred over the network like the 'embedded' method
2. The private key is not hanging around in an enrollment profile which may be inadvertantly exported, shared, or lost
3. Each device has a unique certificate and private key generated for it to use.

The primary use of this device identity certificate is to authenticate the device to the MDM server whenever an HTTPS connection is made. Depending on how the MDM software is configured and how the enrollment profile are setup the device either performs [TLS/SSL client authentication](https://en.wikipedia.org/wiki/Client_certificate) or it provides a [CMS](https://en.wikipedia.org/wiki/Cryptographic_Message_Syntax) detached signature of the MDM request using its identity key pair. In either case this cryptographically proves that the device is using a certificate that should be known to the MDM server.

In the case of SCEP you likely will never have to touch the identity certificate manually. It automatically enrolls and references the correct certificate. Older versions of MicroMDM and e.g. Project-iMAS required manually providing the device identity certificate & key pair in the enrollment profile but that's no longer needed (nor recommended).

It is also possible to use this certificate to encrypt profiles to a device (and only that specific device in the case of unique per-device certificates). More on that later.

## Device Identity Certificate Authority

When using SCEP the device will be issued a certificate from a Certificate Authority (CA). Similarly with embedded-profile device identity certificates they're probably being generated from a CA, too (though they can be self-signed). This usually implies a certificate chain. While not something you may be hands-on managing for MDM it's brought up here to spark a few questions in your environment:

* Where are your device identity certificates getting issued from (which CA)?
* How do you trust this CA? Who owns/operates it?
* How do you *revoke* certificates for compromised devices or lost certificate/key pairs?
* How is the MDM trusting the provided device certificate? Is it merely trusting that the *issuer* issued this certificate? Or is it verifying the actual device identity certificate itself?

In MicroMDM the SCEP CA is built-in using the [micromdm/scep](https://github.com/micromdm/scep) project. It verifies — on each MDM check-in and MDM command — that the certificate was issued by the built-in SCEP CA. Even so it's good to be knowledgeable about how and where these certificates are issued from and trusted.

# Configuration Profile signing & encrypting certificates

As you probably know [Configuration Profiles](https://developer.apple.com/library/content/featuredarticles/iPhoneConfigurationProfileRef) can be CMS signed and/or encrypted. In order to do those things you must encrypt them to a public key or sign them with a certificate.

Encrypting a Configuration Profile that's destined for a device requires using a public key for which that device has the private key to. This can be known through sending a SCEP profile, a certificate profile payload, or just using the device's enrollment identity already on the device once it's enrolled. When the device has the private key then it can decrypt the encrypted profile and install it. In this way encrypting to the device's specific identity certificate it's possible to encrypt profiles to *just* that device.

Signing a Configuration Profile is also possible. For profile signing to be effective the profile should be signed by a certificate that the device trusts. This can be a certificate in the device's [trusted root store](https://en.wikipedia.org/wiki/Public_key_certificate#Root_programs) (similar to a browser's trusted root store) or it can be a certificate that the device is separately configured to trust. In the case of MDM we likely already have a few certificates that are necessarily trusted that we can use. For example the HTTPS web certificate of the MDM may either be trusted in the system root store or be configured as a part of the enrollment profile (say, for a self-signed certificate, more on that later). This may be used to sign profiles (or packages) too. Often though folks will simply sign profiles with their Apple Developer certificate because the signer of those certificates exists in the trusted root store of the device and will work fine.

Seeing as the trusted certificates on a device may be known to the MDM system (e.g. coming through in the enrollment profile) the possibility exists for the MDM system to sign profiles either on command or in an on-the-fly fashion. MicroMDM itself doesn't yet support profile signing or encrypting but is definitely planned. Pre-encrypted or signed profiles should still work, however.

# Configuration Profile trusted certificates

Configuration Profiles are able to add certificates to the device's trusted store of certificates by using Configuration Profile payloads for certificates. Not only is this valuable in-and-of-itself for normal system administration tasks but also importantly for MDM these certificates can be embedded in the *enrollment profile* such that those certificates will be trusted by the system when the device is enrolled.

This has implications for MDM operation *especially* in the case of using [self-signed certificates](https://en.wikipedia.org/wiki/Self-signed_certificate) for the HTTPS server. By default a self-signed certificate on a normal website would simply not be trusted by a device and it's no different for an MDM server. However if we place the MDM server's self-signed HTTPS certificate inside the enrollment profile then the system will trust it and MDM operation can commence after enrollment. The same goes for profile signing mentioned above.

# HTTPS certificate

The web server portion of MDM requires HTTPS. This of course implies [TLS/SSL](https://en.wikipedia.org/wiki/Transport_Layer_Security) certificates and the related trust issues they bring. That said an MDM server's use of TLS/SSL certificates isn't all that different from your typical web server's TLS/SSL configuration.

This means that the types of a certificates you can use are pretty much the same as you can use for a website. These include:

* A purchased TLS/SSL certificate from a reputable Certificate Authority provider like e.g. [Comodo](https://www.comodo.com/), [GlobalSign](https://www.globalsign.com/en/), [Entrust](https://www.entrust.com/), etc.
* MicroMDM supports [LetsEncrypt](https://letsencrypt.org/) (LE) certificate acquisition for automatic SSL configuration
  * Note per usual LE operation this requires public internet inbound TCP port 443 access to your MDM server and properly working public DNS
* Using a [self-signed](https://en.wikipedia.org/wiki/Self-signed_certificate) (or private CA-signed) certificate

If self-signed or private CA-signed certificates are used then we must add those certificates to our enrollment profile using the appropriate configuration profile payloads as mentioned above. MicroMDM does this for us.

MDM is intended for mobile devices and as such there's a reasonable assumption it will be running in a place that is accessible from the public internet. While this is not a technical requirement — an MDM server can run behind a firewall or in a private network as long as APNs push notificates can be sent and received — a lot of the more interesting features of MDM like, say, Remote Wipe lose a lot of their luster without it.

So by far the easiest TLS/SSL configuration is done with LetsEncrypt. Just make sure your server is publicly accessible on port 443 at a domain name you control. MicroMDM should take care of the rest. That said self-signed & purchased SSL options are supported as well in MicroMDM.

# DEP token certificate

Your MDM server talks to Apple's [Device Enrollment Program](https://www.apple.com/business/dep/) API using OAuth tokens. However before you gain access to these tokens you have to complete a [PKI](https://en.wikipedia.org/wiki/Public_key_infrastructure) process where you upload a certificate (which contains a public key) that Apple will use to encrypt the tokens with. Once you've given Apple your public key you can download the encrypted tokens and subsequently decrypt them and start using them to connect to the DEP API (to e.g. fetch & sync devices, configure a DEP profile, etc.). This process is managed via the [Deployment Portal](https://deploy.apple.com/) once your organization has enrolled in DEP and of course have purchased devices that are registered in your DEP account.

Like with the APNs certificate these DEP tokens (like the VPP tokens) expire yearly and this process must be done again to renew them.

# DEP anchor certificates

Remember how we said we could use self-signed HTTPS web certificates? And  that we need to embed the trust information for the HTTPS certificate in the enrollment profile? For a manual enrollment (that is where to say where you visit the profile website manually) this is usually not a problem because the *initial* certificate prompt that one will get (like with any self-signed website) can be simply overridden by the admin doing the enrollment. However, what about a DEP MDM enrollment where there is no user doing the enrollment?

That's what the DEP anchor certificates are for. There's a special property on the DEP profile just for this scenario. Note that DEP profiles should not be confused with Configuration Profiles — DEP profiles are a completely separate JSON structure. This property is called `anchor_certs` in the DEP profile and it allows you to specify trusted certificates to Apple's DEP API. This property will instruct a device at DEP enrollment time to trust the given certificates when connecting to the MDM server over TLS/SSL.

# Conclusion

My hope is this overview has been helpful in untangling the various certificates used in the MDM protocol and perhaps sheds some light on some of the nuances and gotchas surrounding them.

But even so if you've still got questions or are having trouble working with any of these certificates come join us in the MDM-related [MacAdmins Slack channels](/) where we discuss topics like this.
