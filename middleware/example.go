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