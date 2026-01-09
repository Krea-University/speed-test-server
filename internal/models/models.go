// Package models defines database models and operations
package models

import (
	"time"

	"github.com/google/uuid"
)

// SpeedTest represents a speed test record in the database
type SpeedTest struct {
	ID                  string    `json:"id" db:"id"`
	ClientIP            string    `json:"client_ip" db:"client_ip"`
	UserAgent           *string   `json:"user_agent,omitempty" db:"user_agent"`
	TestType            string    `json:"test_type" db:"test_type"`
	DownloadSpeedMbps   *float64  `json:"download_speed_mbps,omitempty" db:"download_speed_mbps"`
	UploadSpeedMbps     *float64  `json:"upload_speed_mbps,omitempty" db:"upload_speed_mbps"`
	PingLatencyMs       *float64  `json:"ping_latency_ms,omitempty" db:"ping_latency_ms"`
	JitterMs            *float64  `json:"jitter_ms,omitempty" db:"jitter_ms"`
	DownloadSizeBytes   *int64    `json:"download_size_bytes,omitempty" db:"download_size_bytes"`
	UploadSizeBytes     *int64    `json:"upload_size_bytes,omitempty" db:"upload_size_bytes"`
	TestDurationSeconds *float64  `json:"test_duration_seconds,omitempty" db:"test_duration_seconds"`
	ISP                 *string   `json:"isp,omitempty" db:"isp"`
	Country             *string   `json:"country,omitempty" db:"country"`
	Region              *string   `json:"region,omitempty" db:"region"`
	City                *string   `json:"city,omitempty" db:"city"`
	ServerName          string    `json:"server_name" db:"server_name"`
	ServerCountry       string    `json:"server_country" db:"server_country"`
	ServerCity          string    `json:"server_city" db:"server_city"`
	Sponsor             string    `json:"sponsor" db:"sponsor"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

// OoklaCompatibleResponse represents an Ookla-compatible speed test response
type OoklaCompatibleResponse struct {
	Type      string              `json:"type"`
	Timestamp time.Time           `json:"timestamp"`
	Ping      *OoklaPingResult    `json:"ping,omitempty"`
	Download  *OoklaSpeedResult   `json:"download,omitempty"`
	Upload    *OoklaSpeedResult   `json:"upload,omitempty"`
	Interface *OoklaInterfaceInfo `json:"interface,omitempty"`
	Server    *OoklaServerInfo    `json:"server,omitempty"`
	Result    *OoklaResultInfo    `json:"result,omitempty"`
}

// OoklaPingResult represents ping results in Ookla format
type OoklaPingResult struct {
	Jitter  float64 `json:"jitter"`
	Latency float64 `json:"latency"`
	Low     float64 `json:"low"`
	High    float64 `json:"high"`
}

// OoklaSpeedResult represents speed results in Ookla format
type OoklaSpeedResult struct {
	Bandwidth int     `json:"bandwidth"` // bits per second
	Bytes     int64   `json:"bytes"`
	Elapsed   int     `json:"elapsed"` // milliseconds
	Latency   float64 `json:"latency"`
}

// OoklaInterfaceInfo represents network interface info
type OoklaInterfaceInfo struct {
	InternalIP string `json:"internalIp"`
	Name       string `json:"name"`
	MacAddr    string `json:"macAddr"`
	IsVpn      bool   `json:"isVpn"`
	ExternalIP string `json:"externalIp"`
}

// OoklaServerInfo represents server information
type OoklaServerInfo struct {
	ID       int     `json:"id"`
	Host     string  `json:"host"`
	Port     int     `json:"port"`
	Name     string  `json:"name"`
	Location string  `json:"location"`
	Country  string  `json:"country"`
	CC       string  `json:"cc"`
	Sponsor  string  `json:"sponsor"`
	Distance float64 `json:"distance"`
	Latency  float64 `json:"latency"`
}

// OoklaResultInfo represents result metadata
type OoklaResultInfo struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// APIKey represents an API key record
type APIKey struct {
	ID                 string     `json:"id" db:"id"`
	KeyHash            string     `json:"-" db:"key_hash"`
	Name               string     `json:"name" db:"name"`
	Description        *string    `json:"description,omitempty" db:"description"`
	RateLimitPerMinute int        `json:"rate_limit_per_minute" db:"rate_limit_per_minute"`
	IsActive           bool       `json:"is_active" db:"is_active"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	LastUsedAt         *time.Time `json:"last_used_at,omitempty" db:"last_used_at"`
}

// RateLimit represents rate limiting data
type RateLimit struct {
	ID           string    `json:"id" db:"id"`
	Identifier   string    `json:"identifier" db:"identifier"`
	Endpoint     string    `json:"endpoint" db:"endpoint"`
	RequestCount int       `json:"request_count" db:"request_count"`
	WindowStart  time.Time `json:"window_start" db:"window_start"`
}

// RateLimitWhitelist represents whitelisted IPs/ranges
type RateLimitWhitelist struct {
	ID          string    `json:"id" db:"id"`
	IPAddress   *string   `json:"ip_address,omitempty" db:"ip_address"`
	IPRange     *string   `json:"ip_range,omitempty" db:"ip_range"`
	Description *string   `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	IsActive    bool      `json:"is_active" db:"is_active"`
}

// TestSession represents an active speed test session
type TestSession struct {
	Token             string    `json:"token"`
	ClientIP          string    `json:"client_ip"`
	UserAgent         string    `json:"user_agent"`
	CreatedAt         time.Time `json:"created_at"`
	ExpiresAt         time.Time `json:"expires_at"`
	PingLatencyMs     *float64  `json:"ping_latency_ms,omitempty"`
	JitterMs          *float64  `json:"jitter_ms,omitempty"`
	DownloadSpeedMbps *float64  `json:"download_speed_mbps,omitempty"`
	UploadSpeedMbps   *float64  `json:"upload_speed_mbps,omitempty"`
	DownloadBytes     *int64    `json:"download_bytes,omitempty"`
	UploadBytes       *int64    `json:"upload_bytes,omitempty"`
	ISP               *string   `json:"isp,omitempty"`
	Country           *string   `json:"country,omitempty"`
	Region            *string   `json:"region,omitempty"`
	City              *string   `json:"city,omitempty"`
}

// NewTestSession creates a new test session with a unique token
func NewTestSession(clientIP, userAgent string) *TestSession {
	return &TestSession{
		Token:     uuid.New().String(),
		ClientIP:  clientIP,
		UserAgent: userAgent,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(15 * time.Minute), // 15 minute expiry
	}
}

// NewSpeedTest creates a new speed test record
func NewSpeedTest(clientIP, testType string) *SpeedTest {
	return &SpeedTest{
		ID:            uuid.New().String(),
		ClientIP:      clientIP,
		TestType:      testType,
		ServerName:    "Krea Speed Test Server",
		ServerCountry: "IN",
		ServerCity:    "Sri City",
		Sponsor:       "Krea University",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// ToOoklaFormat converts a SpeedTest to Ookla-compatible format
func (st *SpeedTest) ToOoklaFormat() *OoklaCompatibleResponse {
	response := &OoklaCompatibleResponse{
		Type:      "result",
		Timestamp: st.CreatedAt,
		Server: &OoklaServerInfo{
			ID:       1,
			Host:     "speed.krea.edu.in",
			Port:     8080,
			Name:     st.ServerName,
			Location: st.ServerCity,
			Country:  st.ServerCountry,
			CC:       st.ServerCountry,
			Sponsor:  st.Sponsor,
			Distance: 0,
		},
		Result: &OoklaResultInfo{
			ID:  st.ID,
			URL: "https://speed.krea.edu.in/result/" + st.ID,
		},
	}

	// Add ping results if available
	if st.PingLatencyMs != nil {
		response.Ping = &OoklaPingResult{
			Latency: *st.PingLatencyMs,
			Jitter:  0, // Default if not available
			Low:     *st.PingLatencyMs,
			High:    *st.PingLatencyMs,
		}
		if st.JitterMs != nil {
			response.Ping.Jitter = *st.JitterMs
		}
		response.Server.Latency = *st.PingLatencyMs
	}

	// Add download results if available
	if st.DownloadSpeedMbps != nil {
		bandwidth := int(*st.DownloadSpeedMbps * 1000000) // Convert Mbps to bps
		response.Download = &OoklaSpeedResult{
			Bandwidth: bandwidth,
			Elapsed:   int(*st.TestDurationSeconds * 1000), // Convert to ms
		}
		if st.DownloadSizeBytes != nil {
			response.Download.Bytes = *st.DownloadSizeBytes
		}
		if st.PingLatencyMs != nil {
			response.Download.Latency = *st.PingLatencyMs
		}
	}

	// Add upload results if available
	if st.UploadSpeedMbps != nil {
		bandwidth := int(*st.UploadSpeedMbps * 1000000) // Convert Mbps to bps
		response.Upload = &OoklaSpeedResult{
			Bandwidth: bandwidth,
			Elapsed:   int(*st.TestDurationSeconds * 1000), // Convert to ms
		}
		if st.UploadSizeBytes != nil {
			response.Upload.Bytes = *st.UploadSizeBytes
		}
		if st.PingLatencyMs != nil {
			response.Upload.Latency = *st.PingLatencyMs
		}
	}

	return response
}
