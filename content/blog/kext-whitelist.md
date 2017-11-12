+++
author = "Victor Vrantchan"
tags = ["mdm", "kext"]
date = "2017-11-13T13:04:37+02:00"
title = "Surviving the Kextpocalyse, Round 0"
frontpage = true

+++

There's an understandable sense of urgency in the MacAdmin community around MDM and Kernel extensions these days. If you've been [paying attention](http://www.richard-purves.com/2017/11/09/mdm-and-the-kextpocalypse-2/) you know that in order to be able to reliably deploy kernel extensions in an enterprise environment, DEP and MDM are becoming a requirement. Apple is likely not going to stop there, and both MDM and DEP will become a requirement for Mac management over the next year or two. But it all feels extremely rushed, and not every enterprise will be ready. Even Apple is not ready yet -- I recently came back from a vacation abroad, where I noticed I couldn't even log in to the DEP portal, because DEP was not available in the country. 

If there's one thing we sysadmins are good at, it's thinking on our feet under pressure and coming up with _alternative_ solutions to problems in our path. Sure, MicroMDM might be a good solution longterm, but you're not ready to deal with the certificates, and deploy it yet. But you need a solution for your upcoming kext problems as early as today. I have one!

# The New Kext Rules

- User will be prompted to approve a kernel extension if the Mac is not enrolled in MDM.
- If there's an MDM, all kext installs are whitelisted. 
- You'll be grandfathered in to new rules if the enrollment profile is already there.

With the above rules in mind, we can start working on an MDM, but short term our only problem is that we need the Mac to be considered enrolled. Turns out, this is actually somewhat trivial. Allow me to explain.

# MDM Enrollment Rules

Rolling out a new MDM might be hard, but enrolling into one isn't. When I was [just starting to play with MDM](https://groob.io/posts/mdm-experiments/), there was no public documentation of the spec, yet I was able to have a test device enrolled into an nginx server in just a few minutes. There's only two things you need to enroll into MDM before 10.13.2 rolls out: An HTTP Server that responds with HTTP `200 OK` to any request coming from the MDM. You won't be able to use that to install profiles, or do other things MDMs allow, but you'll be considered enrolled, which is what we're aiming for here.

# Hackery

So far we've established that we want enroll our entire Mac fleet into an MDM for the purpose of being grandfathered in to the default Kernel Extension whitelisting rule. Later, when 10.13.2+ ships with new Macs, or when we re-image(is that still a thing???) we'll remove the MDM enrollment and enroll in a proper MDM. 
Now we need an enrollment profile. We'll use Profile Manager to create and export it, then a text editor to tweak it, and finally Munki/Puppet/Chef/ARD/Custom Pkg to deploy it to our Mac fleet.

1) Configure Profile Manger and export the default enrollment profile.

2) Create a Device Identity Certificate
By default the Profile Manger Enrollment Profile uses SCEP, but we'll get rid of that, and add a .p12 file instead. Export the certificate you created as a P12 file.

3) Use Apple Configurator to tweak the Enrollment Profile.
- Remove SCEP
- Remove any existing entries from Certificates
- Add your .p12 file as an entry to Certificates

Now open the enrollment profile in a text editor. There's a few more changes to do there.
- Find ServerURL and CheckinURL keys and change those to a URL that returns 200 OK no matter the request. An `nginx` or `caddy` instance with Let's Encrypt certs will do. `facebook.com` also works, although I don't recommend/endorse that.

- Find the `IdentityCertificateUUID` key, and make sure that the UUID there is identitical to the `PayloadUUID` for the .p12 payload. if it's not put the `IdentityCertificateUUID` uuid in that field.

Your Enrollment Profile is now ready. Feel free to deploy it to your fleet of <= 10.13.1 devices and you should be ready to survive the 10.13.2 rollout... for now.

I added a sample enrollment profile which should work for anyone who wants to try the above steps or double check their own: https://gist.github.com/groob/c54b3907498de18221f5c93d56083a54

# Conclusion

I wrote up the above workflow because I found it amusing that the requirement to be considered "Enrolled" are so light, but Apple is enforcing it everywhere. 
Long term, I _do not endorse_, anything I've described above. If you're not yet using an MDM, take a look at setting one up. And set a DEP account as well. The days where you could ignore MDM are long gone, and any enterprise that does not have these services enabled over the next couple months will suffer.  




