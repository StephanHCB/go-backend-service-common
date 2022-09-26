package swaggerctl

import (
	"context"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	auwebswaggerui "github.com/StephanHCB/go-autumn-web-swagger-ui"
	"github.com/StephanHCB/go-backend-service-common/web/util/media"
	"github.com/go-chi/chi/v5"
	"github.com/go-http-utils/headers"
	"net/http"
	"os"
	"path/filepath"
)

type SwaggerCtlImpl struct{}

func (c *SwaggerCtlImpl) WireUp(ctx context.Context, router chi.Router) {
	// 	serve swagger-ui and openapi spec json (which needs to be in the file system of your container)
	c.AddStaticHttpFilesystemRoute(router, auwebswaggerui.Assets, "/swagger-ui")
	c.AddStaticSingleFileRoute(router, "docs", "openapi-v3-spec.json", "/")
	c.AddRedirect(router, "/v3/api-docs", "/openapi-v3-spec.json")
}

// serve static files from an http.FileSystem instance via a chi router
//
// inspired by https://github.com/go-chi/chi/blob/master/_examples/fileserver/main.go
//
// parameters:
// - uriPath under which sub-url they should be served, if you leave out trailing slash, also adds a redirect
//   example: "/swagger-ui"
func (c *SwaggerCtlImpl) AddStaticHttpFilesystemRoute(server chi.Router, fs http.FileSystem, uriPath string) {
	strippedFs := http.StripPrefix(uriPath, http.FileServer(fs))

	if hasNoTrailingSlash(uriPath) {
		server.Get(uriPath, http.RedirectHandler(uriPath+"/", 301).ServeHTTP)
		uriPath += "/"
	}
	uriPath += "*"

	server.Get(uriPath, func(w http.ResponseWriter, r *http.Request) {
		strippedFs.ServeHTTP(w, r)
	})
}

// serve a single static file via a chi router
//
// parameters:
// - relativeFilesystemPath where to find the file(s) to serve relative to the current working directory
//   example: "docs"
//   note: do NOT add a trailing slash
// - filename which exact file to serve. This will be added to the route, so only exactly this file is made available
//   example: "swagger.json"
// - uriPath under which path the file should be served
//   example: "/"
//   Note: unfortunately it is not possible to use a different filename, as this is a direct filesystem directory server.
//   That's why we have the redirect.
func (c *SwaggerCtlImpl) AddStaticSingleFileRoute(server chi.Router, relativeFilesystemPath string, filename string, uriPath string) {
	workDir, _ := os.Getwd()
	filePath := filepath.Join(workDir, relativeFilesystemPath, filename)

	contents, err := os.ReadFile(filePath)
	if err != nil {
		aulogging.Logger.NoCtx().Error().WithErr(err).Printf("failed to read file %s - skipping: %s", filePath, err.Error())
		return
	}

	if hasNoTrailingSlash(uriPath) {
		uriPath += "/"
	}

	server.Get(uriPath+filename, func(w http.ResponseWriter, r *http.Request) {
		// this stops browsers from caching our swagger.json
		w.Header().Set(headers.CacheControl, "no-cache")
		w.Header().Set(headers.ContentType, media.ContentTypeApplicationJson)
		_, _ = w.Write(contents)
	})
}

func hasNoTrailingSlash(path string) bool {
	return path != "/" && path[len(path)-1] != '/'
}

func (c *SwaggerCtlImpl) AddRedirect(server chi.Router, source string, target string) {
	server.Get(source, func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target, http.StatusMovedPermanently)
	})
}
