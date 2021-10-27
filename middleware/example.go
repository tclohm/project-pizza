// MARK: -- middleware! 
// if we needed to pass values between handlers, such as
// the ID of the authenticated user, or a request or a trace
// ID, we can use context.Context attached to the *http.Request
// via the *Request.Context()
// func Middleware(next http.Handler) http.Handler {
// 	// wrap anon func, and cast it to a http.HandlerFunc
// 	// signature matches ServeHTTP(w,r)
// 	return http.HandlerFunc(
// 		func(w http.ResponseWriter, r *http.Request) {
// 			// logic before - reading request values,

// 			// call the 'next' handler in the chain
// 			next.ServeHTTP(w, r)
// 		}
// 	)
// }

// configuredRouter := LoggingMiddleware(OtherMiddleware(YetAnotherMiddleware(router)))
// log.Fatal(http.ListenAndServe(":8000", configuredRouter))

// func NewExampleMiddleware(something string) func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		fn := func(w http.ResponseWriter, r *http.Request) {
// 			// Logic here

// 			// Call the next handler
// 			next.ServeHTTP(w, r))
// 		}

// 		return http.HandlerFunc(fn)
// 	}
// }

// Minimal wrapper for http.ResponseWriter that allows the 
// written HTTP status code to be captured for logging


// type responseWriter struct {
// 	http.ResponseWriter
// 	status int
// 	wroteHeader bool
// }

// func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
// 	return &responseWriter{ResponseWriter: w}
// }

// func (rw *responseWriter) Status() int {
// 	return rw.status
// }

// func (rw *responseWriter) WriteHeader(code int) {
// 	if rw.wroteHeader {
// 		return
// 	}

// 	rw.status = code
// 	rw.ResponseWriter.WriteHeader(code)
// 	rw.wroteHeader = true

// 	return
// }

// // LoggingMiddleware logs the incoming HTTP request & its duration.
// func LoggingMiddleware(logger log.Logger) func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		fn := func(w http.ResponseWriter, r *http.Request) {
// 			defer func() {
// 				if err := recover(); err != nil {
// 					w.WriteHeader(http.StatusInternalServerError)
// 					logger.Log(
// 						"err", err,
// 						"trace", debug.Stack(),
// 					)
// 				}
// 			}()

// 			start := time.Now()
// 			wrapped := wrapResponseWriter(w)
// 			next.ServeHTTP(wrapped, r)
// 			logger.Log(
// 				"status", wrapped.status,
// 				"method", r.Method,
// 				"path", r.URL.EscapedPath(),
// 				"duration", time.Since(start),
// 			)
// 		}
// 		return http.HandlerFunc(fn)
// 	}
// }
