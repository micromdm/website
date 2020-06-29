// Binary vanity serves go-import and go-source URLs for Go module imports.
package main

import (
	"bytes"
	"html/template"
	"net/http"
	"net/url"
	"strings"

	"micromdm.io/v2/pkg/log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// See "go help importpath" or https://golang.org/cmd/go/#hdr-Remote_import_paths
var vanityTemplate = template.Must(template.New("vanity.html").Parse(`
<head>
  <meta name="go-import" content="{{ .Permalink }} git {{ .Repo }}">
  <meta name="go-source" content="{{ .Permalink }} {{ .Repo }} {{ .Repo }}/tree/{{ .Tree }}{/dir} {{ .Repo }}/blob/{{ .Tree }}{/dir}/{file}#L{line}">
</head>
`))

type server struct {
	logger log.Logger
	myURL  *url.URL
}

func (srv server) handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	log.Info(srv.logger).Log("request_path", request.Path)

	pURL, err := srv.myURL.Parse(request.Path)
	if err != nil {
		log.Info(srv.logger).Log("msg", "parsing relative to  base url", "err", err)
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "oops",
		}, nil
	}

	var vanity = struct {
		Permalink string
		Repo      string
		Tree      string
	}{
		Permalink: strings.TrimPrefix(pURL.String(), "https://"),
		Repo:      "https://github.com/micromdm/micromdm",
		Tree:      "v2dev",
	}

	log.Info(srv.logger).Log(
		"request_path", request.Path,
		"permalink", vanity.Permalink,
	)

	buf := new(bytes.Buffer)
	if err := vanityTemplate.Execute(buf, vanity); err != nil {
		log.Info(srv.logger).Log("msg", "execute vanityTemplate", err, "err")
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "oops",
		}, nil
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "text/html; charset=UTF-8"},
		Body:       buf.String(),
	}, nil
}

func main() {
	myURL, _ := url.Parse("https://micromdm.io/")
	srv := server{
		logger: log.New(),
		myURL:  myURL,
	}
	lambda.Start(srv.handler)
}
