// Package server provides a development HTTP server for static sites.
//
// The server serves files from a directory with support for clean URLs,
// where requests to /path/ are served from /path/index.html. Static assets
// such as CSS and images are served with appropriate content types.
//
// # Usage
//
//	srv := server.New(server.Config{
//		Port: 8080,
//		Dir:  "public",
//	})
//	if err := srv.Start(); err != nil {
//		log.Fatal(err)
//	}
//	defer srv.Shutdown(context.Background())
//
// Use [Server.Addr] to retrieve the address the server is listening on,
// which is useful when using port 0 for automatic port assignment in tests.
//
// # Graceful Shutdown
//
// Call [Server.Shutdown] with a context to gracefully stop the server,
// allowing in-flight requests to complete within the context deadline.
package server
