package middleware

import (
	"net/http"
	"net/url"
	"strings"
)

func CSRFProtect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only inspect mutating methods
		if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		// 1. Sec-Fetch-Site check (modern browsers)
		if fetchSite := r.Header.Get("Sec-Fetch-Site"); fetchSite != "" {
			if fetchSite == "cross-site" {
				http.Error(w, "Forbidden: cross-origin request rejected", http.StatusForbidden)
				return
			}
		}

		// 2. Origin / Referer check against Host
		origin := r.Header.Get("Origin")
		if origin != "" {
			parsedOrigin, err := url.Parse(origin)
			if err != nil || !strings.EqualFold(parsedOrigin.Host, r.Host) {
				http.Error(w, "Forbidden: origin mismatch", http.StatusForbidden)
				return
			}
		} else if referer := r.Header.Get("Referer"); referer != "" {
			parsedReferer, err := url.Parse(referer)
			if err != nil || !strings.EqualFold(parsedReferer.Host, r.Host) {
				http.Error(w, "Forbidden: referer mismatch", http.StatusForbidden)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
