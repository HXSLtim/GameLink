package middleware

import (
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// MetricsAuthConfig configures the metrics endpoint authentication
type MetricsAuthConfig struct {
	// Enabled enables or disables authentication for the metrics endpoint
	Enabled bool

	// AllowedCIDRs is a list of CIDR blocks allowed to access metrics (e.g., "10.0.0.0/8", "192.168.1.0/24")
	// If empty, only localhost (127.0.0.0/8 and ::1/128) is allowed
	AllowedCIDRs []string

	// TrustedProxies is a list of trusted proxy IPs for X-Forwarded-For header parsing
	TrustedProxies []string
}

// DefaultMetricsAuthConfig returns the default metrics authentication configuration
// By default, only localhost can access metrics in production
func DefaultMetricsAuthConfig() MetricsAuthConfig {
	return MetricsAuthConfig{
		Enabled:        true,
		AllowedCIDRs:   []string{}, // Empty means localhost only
		TrustedProxies: []string{"127.0.0.1", "::1"},
	}
}

// MetricsAuth creates a middleware that restricts access to the metrics endpoint
// based on IP address. This is suitable for Prometheus scraping in production.
//
// Usage:
//   router.GET("/metrics", middleware.MetricsAuth(middleware.DefaultMetricsAuthConfig()), handler)
func MetricsAuth(config MetricsAuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.Enabled {
			// Authentication disabled - allow all access (dev/staging only)
			c.Next()
			return
		}

		// Get client IP
		clientIP := getMetricsClientIP(c)

		// Check if IP is allowed
		if !isIPAllowed(clientIP, config.AllowedCIDRs) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"code":    http.StatusForbidden,
				"message": "access to metrics endpoint is forbidden",
			})
			return
		}

		c.Next()
	}
}

// getMetricsClientIP extracts the client IP from the request, considering X-Forwarded-For header
// Renamed from getClientIP to avoid conflict with rateLimiter.getClientIP
func getMetricsClientIP(c *gin.Context) string {
	// Check X-Forwarded-For header for reverse proxy scenarios
	xForwardedFor := c.GetHeader("X-Forwarded-For")
	if xForwardedFor != "" {
		// X-Forwarded-For can contain multiple IPs, the first one is the client
		ips := strings.Split(xForwardedFor, ",")
		if len(ips) > 0 {
			clientIP := strings.TrimSpace(ips[0])
			if clientIP != "" {
				return clientIP
			}
		}
	}

	// Check X-Real-IP header
	xRealIP := c.GetHeader("X-Real-IP")
	if xRealIP != "" {
		return xRealIP
	}

	// Fall back to RemoteAddr
	return c.ClientIP()
}

// isIPAllowed checks if the given IP is allowed based on the CIDR list
func isIPAllowed(clientIP string, allowedCIDRs []string) bool {
	ip := net.ParseIP(clientIP)
	if ip == nil {
		return false
	}

	// If no CIDRs specified, only allow localhost
	if len(allowedCIDRs) == 0 {
		return isLocalhost(ip)
	}

	// Check against each CIDR block
	for _, cidr := range allowedCIDRs {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if ipNet.Contains(ip) {
			return true
		}
	}

	// Also allow localhost even if CIDRs are specified
	return isLocalhost(ip)
}

// isLocalhost checks if the IP is a loopback address
func isLocalhost(ip net.IP) bool {
	if ip.IsLoopback() {
		return true
	}

	// Check for IPv4 loopback (127.0.0.0/8)
	if ip4 := ip.To4(); ip4 != nil && ip4[0] == 127 {
		return true
	}

	return false
}

// AdminOnly creates a middleware that requires admin authentication
// This can be used as an alternative to IP-based authentication
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user is authenticated
		userID, exists := c.Get("user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    http.StatusUnauthorized,
				"message": "authentication required",
			})
			return
		}

		// Check if user has admin role
		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"code":    http.StatusForbidden,
				"message": "admin access required",
			})
			return
		}

		// Check for admin or superAdmin role
		roleStr, ok := role.(string)
		if !ok || (roleStr != "admin" && roleStr != "superAdmin") {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"code":    http.StatusForbidden,
				"message": "admin access required",
			})
			return
		}

		// Set admin ID for logging
		c.Set("admin_id", userID)
		c.Next()
	}
}
