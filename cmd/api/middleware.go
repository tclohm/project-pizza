package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {
	// client struct to hold the rate limiter and last seen time for each client
	type client struct {
		limiter *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu sync.Mutex
		clients = make(map[string]*client)
	)

	// launch background goroutine to remove old entries from the client map
	go func() {
		for {
			time.Sleep(time.Minute)
			// lock the mutex to prevent any rate limiter checks from happening while cleanup takes place
			mu.Lock()

			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3 * time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// extract clients IP
		if app.config.limiter.enabled {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

			// Lock the mutex to prevent this code from being executed concurrently
			mu.Lock()

			// if the ip address doesn't exist, initialize a new rate limiter and add the ip address
			if _, found := clients[ip]; !found {
				clients[ip] = &client{
					limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst),
				}
			}

			// update the last seen time for the client
			clients[ip].lastSeen = time.Now()

			// call allow for the current ip address. If the request isn't allowed, unlock the mutex
			// and send a 429 Too Many Requests
			// response, just like before
			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				app.rateLimitExceededResponse(w, r)
				return
			}

			// unlock mutex before calling the next handler in the chain
			mu.Unlock()
		}
		

		next.ServeHTTP(w, r)
	})
}

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")

		w.Header().Add("Vary", "Access-Control-Request-Method")

		origin := r.Header.Get("Origin")

		if origin != "" {
			for i := range app.config.cors.trustedOrigins {
				if origin == app.config.cors.trustedOrigins[i] {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					
					// check request has http method options
					if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
						w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
						//w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

						w.WriteHeader(http.StatusOK)
						return
					}

					break
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}