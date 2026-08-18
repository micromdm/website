.netlify/functions/vanity: netlify-functions/vanity/*.go
	CGO_ENABLED=0 GOBIN=$(PWD)/.netlify/functions go install ./netlify-functions/vanity

netlify: .netlify/functions/vanity
	hugo --gc --minify
	cp _redirects public/_redirects
	
