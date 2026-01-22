package commands

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/gamelife2025/easytime/pkg/utils"
	"github.com/spf13/cobra"
)

//go:embed web.html
var webHTML string

var webPort *int

func init() {
	webPort = webCmd.PersistentFlags().Int("port", 80, " --port 80")
	rootCmd.AddCommand(webCmd)
}

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "以web端运行",
	RunE: func(cmd *cobra.Command, args []string) error {
		return http.ListenAndServe(fmt.Sprintf(":%d", *webPort), &HTTPServer{})
	},
}

type HTTPServer struct {
}

type ConvertResponse struct {
	Success bool        `json:"success"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type TimeInfo struct {
	ISO8601       string `json:"iso8601"`
	Timestamp     int64  `json:"timestamp"`
	TimestampMs   int64  `json:"timestamp_ms"`
	TimestampUs   int64  `json:"timestamp_us"`
	Date          string `json:"date"`
	Time          string `json:"time"`
	Weekday       string `json:"weekday"`
	ZeroTimestamp int64  `json:"zero_timestamp"`
	Timezone      string `json:"timezone"`
}

func (server *HTTPServer) ServeHTTP(write http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/api/convert" {
		server.handleConvert(write, req)
		return
	}

	if req.URL.Path == "/api/now" {
		server.handleNow(write, req)
		return
	}

	// Serve HTML
	tpl := template.New("main")
	tpl.Parse(webHTML)
	tpl.Execute(write, map[string]interface{}{
		"Title": "easytime",
	})
}

func (server *HTTPServer) handleConvert(write http.ResponseWriter, req *http.Request) {
	write.Header().Set("Content-Type", "application/json")

	if req.Method != http.MethodPost {
		resp := ConvertResponse{Success: false, Error: "Method not allowed"}
		json.NewEncoder(write).Encode(resp)
		return
	}

	var input struct {
		Input    string `json:"input"`
		Timezone string `json:"timezone"`
	}

	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		resp := ConvertResponse{Success: false, Error: "Invalid request"}
		json.NewEncoder(write).Encode(resp)
		return
	}

	input.Input = strings.TrimSpace(input.Input)
	if input.Input == "" {
		resp := ConvertResponse{Success: false, Error: "Input cannot be empty"}
		json.NewEncoder(write).Encode(resp)
		return
	}

	// Load timezone
	var loc *time.Location = time.Local
	if input.Timezone != "" && input.Timezone != "Local" {
		var err error
		loc, err = time.LoadLocation(input.Timezone)
		if err != nil {
			resp := ConvertResponse{Success: false, Error: "Invalid timezone"}
			json.NewEncoder(write).Encode(resp)
			return
		}
	}
	var t time.Time
	var err error

	t, err = utils.Get(input.Input)
	if err != nil {
		resp := ConvertResponse{Success: false, Error: "Invalid input format"}
		json.NewEncoder(write).Encode(resp)
		return
	}
	t.In(loc)
	info := buildTimeInfo(t)
	resp := ConvertResponse{Success: true, Data: info}
	json.NewEncoder(write).Encode(resp)
}

func (server *HTTPServer) handleNow(write http.ResponseWriter, req *http.Request) {
	write.Header().Set("Content-Type", "application/json")

	timezone := req.URL.Query().Get("timezone")
	var loc *time.Location = time.Local
	if timezone != "" && timezone != "Local" {
		var err error
		loc, err = time.LoadLocation(timezone)
		if err != nil {
			resp := ConvertResponse{Success: false, Error: "Invalid timezone"}
			json.NewEncoder(write).Encode(resp)
			return
		}
	}

	t := time.Now().In(loc)
	info := buildTimeInfo(t)
	resp := ConvertResponse{Success: true, Data: info}
	json.NewEncoder(write).Encode(resp)
}

func buildTimeInfo(t time.Time) TimeInfo {
	weekdays := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	return TimeInfo{
		ISO8601:       t.Format(time.RFC3339),
		Timestamp:     t.Unix(),
		TimestampMs:   t.UnixMilli(),
		TimestampUs:   t.UnixMicro(),
		Date:          t.Format("2006-01-02"),
		Time:          t.Format("15:04:05"),
		Weekday:       weekdays[t.Weekday()],
		ZeroTimestamp: time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix(),
		Timezone:      t.Location().String(),
	}
}
