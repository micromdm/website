.netlify/functions/vanity: netlify-functions/vanity/*.go
	GOBIN=$(PWD)/.netlify/functions go install ./netlify-functions/vanity

netlify: .netlify/functions/vanity
	hugo --gc --minify
	cp _redirects public/_redirects
	
