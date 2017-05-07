+++
frontpage = true
author = "Victor Vrantchan"
tags = ["mdm", "dep","enterprise", "account", "apple"]
date = "2017-05-06T13:04:37+02:00"
title = "The business side of MDM - Do you know your DUNS Number?"

+++

Getting started with MDM is hard enough, but the toughest part is not technical, it's getting approval to create all the accounts with Apple in order to begin. I recently went through the whole process for my company and thought it would be useful to others if I blogged about it. Below is infromation about what accounts you might need, and how you might go about acquiring them.

<!--more-->

# Have your company's trust

When you sign up on behalf of your organization with Apple, the form will ask wether you're trusted to make such a decision on behalf of your company. The form will also ask for the contact person of your superior to verify. Apple will most surely call you and the contact you've specified during sign up to verify. I was contacted both during DEP approval and when opening an Enterprise Developer account.

My advice is to coordinate with the person whose name you write down, and make sure they're available that week. It will really help expedite the approval process if both you and your superior are around and able to take phone calls. In both cases, my conversation with Apple only took a few minutes. 


# Have a DUNS number

The Developer program as well as DEP and VPP accounts require that you register with a DUNS Number. While many organizations already have one, not everyone does, and it might take a while to register for one. 

Apple publishes a [help document](https://developer.apple.com/support/D-U-N-S/) and [lookup tool](https://developer.apple.com/enroll/duns-lookup/) in case you don't know yours. 

If you're just applying now, keep this timeline in mind:

> While expediting your D-U-N-S Number creation process may enable you to receive your number sooner, it will still take up to 14 business days for D&B to provide updated information to Apple. You will not be able to use your number to enroll until this step has been completed.

# Create your Apple IDs

Apple will require you to create an AppleID for each account you create. If you're creating accounts on behalf of a company, it's important to remmember you might not be the only one accessing it, and that one day you might leave and someone else will inherit that account. With that in mind, I went the route of creating a group in Gmail which has several aliases -- one for each administrative Apple ID.

# Sign up for an Enterprise Developer Account

An organizational developer account with Apple is important for any Mac Administrator. If you manage macs in your enterprise, you likely create and distribute scripts and internal applications. And when you do that, the software should always be signed. Beyond the ability to sign software, you also gain access to documentation which is not otherwise acccessible.

If you're looking to either develop an MDM, or manage one yourself (as opposed to signing up with a vendor), you will need an enterpise developer to get access to the MDM vendor certificate. It is called the `MDM CSR` option. Curiously, this certificate option is only available under iOS certificates, not macOS. 

If you already have an enterprise account, but don't see the the `MDM CSR` option, it's because Apple requires special approval. The main administrator of your enterprise account must email Apple to ask for the option to be enabled. 

If you're creating your enterprise account today, Apple will ask you if you're an MDM vendor. Select yes, even if you don't intend to use MDM beyond your own company's needs. During the phone call for approval, Apple inquired about our intended use, and I specified that it was for internal use. They did not object to this use.

{{<figure src="/images/mdm_vendor.png" class="screenshot" >}}

You can begin your enterprise developer enrollment [here](https://developer.apple.com/programs/enterprise/). This account also comes with a yearly cost of $299.

# Purchase your devices direct from Apple or through an authorized reseller 

If you'd like to take advantage of the Device Enrollment Program(DEP), you must first buy your devices through an approved channel, like an authorized reseller or direct from Apple.

https://ecommerce.apple.com

# Sign up for DEP and VPP 

Just like the developer account, DEP registration requires a DUNS Number and a phonecall from Apple for approval. The AppleID you sign up with must also be new/not associated with other accounts. You'll also be required to enable 2FA for each of these iCloud accounts. 

You can begin the process at [deploy.apple.com](https://deploy.apple.com), unless you're an EDU organization, in which case sign up with Apple School Manager at [school.apple.com](https://school.apple.com/). Apple School Manager is a version of DEP reserved for educational institutions.


