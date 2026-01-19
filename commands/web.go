package commands

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/gamelife2025/easytime/pkg/utils"
	"github.com/spf13/cobra"
)

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
	tpl.Parse(web_tpl)
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

var web_tpl = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>{{.Title}} - 时间转换工具</title>
	<style>
		* {
			margin: 0;
			padding: 0;
			box-sizing: border-box;
		}

		body {
			font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
			background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
			min-height: 100vh;
			display: flex;
			justify-content: center;
			align-items: center;
			padding: 20px;
		}

		.container {
			background: white;
			border-radius: 16px;
			box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
			max-width: 500px;
			width: 100%;
			overflow: hidden;
		}

		.header {
			background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
			color: white;
			padding: 40px 30px;
			text-align: center;
		}

		.header h1 {
			font-size: 28px;
			margin-bottom: 8px;
			font-weight: 600;
		}

		.header p {
			font-size: 14px;
			opacity: 0.9;
		}

		.content {
			padding: 30px;
		}

		.form-group {
			margin-bottom: 24px;
		}

		label {
			display: block;
			margin-bottom: 8px;
			font-size: 14px;
			font-weight: 500;
			color: #333;
		}

		input, select {
			width: 100%;
			padding: 12px 16px;
			border: 1px solid #e0e0e0;
			border-radius: 8px;
			font-size: 14px;
			transition: all 0.3s ease;
			font-family: inherit;
		}

		input:focus, select:focus {
			outline: none;
			border-color: #667eea;
			box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
		}

		.button-group {
			display: flex;
			gap: 12px;
		}

		button {
			flex: 1;
			padding: 12px 24px;
			border: none;
			border-radius: 8px;
			font-size: 14px;
			font-weight: 600;
			cursor: pointer;
			transition: all 0.3s ease;
		}

		.btn-convert {
			background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
			color: white;
		}

		.btn-convert:hover {
			transform: translateY(-2px);
			box-shadow: 0 10px 20px rgba(102, 126, 234, 0.3);
		}

		.btn-convert:active {
			transform: translateY(0);
		}

		.btn-now {
			background: #f0f0f0;
			color: #333;
		}

		.btn-now:hover {
			background: #e8e8e8;
		}

		.result {
			background: #f8f9fa;
			border-radius: 8px;
			padding: 20px;
			margin-top: 24px;
			display: none;
		}

		.result.show {
			display: block;
			animation: slideIn 0.3s ease;
		}

		@keyframes slideIn {
			from {
				opacity: 0;
				transform: translateY(-10px);
			}
			to {
				opacity: 1;
				transform: translateY(0);
			}
		}

		.result-item {
			display: grid;
			grid-template-columns: 120px 1fr auto;
			gap: 16px;
			align-items: center;
			padding: 12px 0;
			border-bottom: 1px solid #e0e0e0;
			font-size: 14px;
		}

		.result-item:last-child {
			border-bottom: none;
		}

		.result-label {
			color: #666;
			font-weight: 500;
		}

		.result-value {
			color: #333;
			font-family: 'Monaco', 'Courier New', monospace;
			font-size: 13px;
			word-break: break-all;
			text-align: right;
		}

		.copy-btn {
			background: none;
			border: none;
			color: #667eea;
			cursor: pointer;
			font-size: 12px;
			padding: 4px 8px;
			min-width: 40px;
			text-align: center;
			transition: color 0.3s ease;
			white-space: nowrap;
		}

		.copy-btn:hover {
			color: #764ba2;
		}

		.error {
			background: #fee;
			color: #c33;
			padding: 12px 16px;
			border-radius: 8px;
			margin-top: 16px;
			font-size: 13px;
			display: none;
		}

		.error.show {
			display: block;
		}

		.success-msg {
			position: fixed;
			top: 20px;
			right: 20px;
			background: #4caf50;
			color: white;
			padding: 12px 20px;
			border-radius: 8px;
			font-size: 13px;
			opacity: 0;
			transition: opacity 0.3s ease;
		}

		.success-msg.show {
			opacity: 1;
		}

		.timezone-list {
			max-height: 200px;
			overflow-y: auto;
		}

		.loading {
			display: none;
			text-align: center;
		}

		.spinner {
			display: inline-block;
			width: 16px;
			height: 16px;
			border: 2px solid #f3f3f3;
			border-top: 2px solid #667eea;
			border-radius: 50%;
			animation: spin 0.8s linear infinite;
		}

		@keyframes spin {
			0% { transform: rotate(0deg); }
			100% { transform: rotate(360deg); }
		}

		@media (max-width: 480px) {
			.header {
				padding: 30px 20px;
			}

			.header h1 {
				font-size: 24px;
			}

			.content {
				padding: 20px;
			}

			.result-item {
				grid-template-columns: 1fr;
				gap: 8px;
			}

			.result-label {
				display: block;
			}

			.result-value {
				text-align: left;
				word-break: break-all;
			}

			.copy-btn {
				display: none;
			}
		}
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>⏰ easytime</h1>
		</div>

		<div class="content">
			<form id="convertForm">
				<div class="form-group">
					<label for="input">输入时间或时间戳</label>
					<input 
						type="text" 
						id="input" 
						placeholder="例如：2024-01-29 或 1706486400"
						autocomplete="off"
					>
				</div>

				<div class="form-group">
					<label for="timezone">时区</label>
					<select id="timezone">
						<option value="Local">系统时区</option>
						<option value="UTC">UTC</option>
						<option value="Asia/Shanghai">亚洲/上海 (CST)</option>
						<option value="Asia/Tokyo">亚洲/东京 (JST)</option>
						<option value="Asia/Hong_Kong">亚洲/香港 (HKT)</option>
						<option value="Asia/Singapore">亚洲/新加坡 (SGT)</option>
						<option value="Asia/Bangkok">亚洲/曼谷 (ICT)</option>
						<option value="Australia/Sydney">澳洲/悉尼 (AEDT)</option>
						<option value="America/New_York">美洲/纽约 (EST)</option>
						<option value="America/Los_Angeles">美洲/洛杉矶 (PST)</option>
						<option value="America/Chicago">美洲/芝加哥 (CST)</option>
						<option value="Europe/London">欧洲/伦敦 (GMT)</option>
						<option value="Europe/Paris">欧洲/巴黎 (CET)</option>
						<option value="Europe/Berlin">欧洲/柏林 (CET)</option>
						<option value="Europe/Moscow">欧洲/莫斯科 (MSK)</option>
						<option value="Africa/Cairo">非洲/开罗 (EET)</option>
						<option value="Africa/Johannesburg">非洲/约翰内斯堡 (SAST)</option>
						<option value="Indian/Dubai">印度/迪拜 (GST)</option>
						<option value="Pacific/Auckland">太平洋/奥克兰 (NZDT)</option>
					</select>
				</div>

				<div class="button-group">
					<button type="submit" class="btn-convert">转换</button>
					<button type="button" class="btn-now" id="nowBtn">当前时间</button>
				</div>

				<div class="error" id="error"></div>
			</form>

			<div class="loading" id="loading">
				<div class="spinner"></div>
			</div>

			<div class="result" id="result">
				<div class="result-item">
					<span class="result-label">ISO 8601</span>
					<span class="result-value" id="iso8601"></span>
					<button class="copy-btn" onclick="copyToClipboard('iso8601')">复制</button>
				</div>
				<div class="result-item">
					<span class="result-label">日期</span>
					<span class="result-value" id="date"></span>
					<button class="copy-btn" onclick="copyToClipboard('date')">复制</button>
				</div>
				<div class="result-item">
					<span class="result-label">时间</span>
					<span class="result-value" id="time"></span>
					<button class="copy-btn" onclick="copyToClipboard('time')">复制</button>
				</div>
				<div class="result-item">
					<span class="result-label">时间戳 (秒)</span>
					<span class="result-value" id="timestamp"></span>
					<button class="copy-btn" onclick="copyToClipboard('timestamp')">复制</button>
				</div>
				<div class="result-item">
					<span class="result-label">时间戳 (毫秒)</span>
					<span class="result-value" id="timestamp_ms"></span>
					<button class="copy-btn" onclick="copyToClipboard('timestamp_ms')">复制</button>
				</div>
				<div class="result-item">
					<span class="result-label">时间戳 (微秒)</span>
					<span class="result-value" id="timestamp_us"></span>
					<button class="copy-btn" onclick="copyToClipboard('timestamp_us')">复制</button>
				</div>
				<div class="result-item">
					<span class="result-label">星期</span>
					<span class="result-value" id="weekday"></span>
				</div>
				<div class="result-item">
					<span class="result-label">零点时间戳</span>
					<span class="result-value" id="zero_timestamp"></span>
					<button class="copy-btn" onclick="copyToClipboard('zero_timestamp')">复制</button>
				</div>
				<div class="result-item">
					<span class="result-label">时区</span>
					<span class="result-value" id="timezone_info"></span>
					<button class="copy-btn" onclick="copyToClipboard('timezone_info')">复制</button>
				</div>
			</div>
		</div>
	</div>

	<div class="success-msg" id="successMsg">已复制到剪贴板</div>

	<script>
		const form = document.getElementById('convertForm');
		const input = document.getElementById('input');
		const timezone = document.getElementById('timezone');
		const result = document.getElementById('result');
		const error = document.getElementById('error');
		const loading = document.getElementById('loading');
		const nowBtn = document.getElementById('nowBtn');
		const successMsg = document.getElementById('successMsg');

		form.addEventListener('submit', async (e) => {
			e.preventDefault();
			await convert();
		});

		nowBtn.addEventListener('click', async () => {
			input.value = '';
			await fetchNow();
		});

		input.addEventListener('keypress', (e) => {
			if (e.key === 'Enter') {
				form.dispatchEvent(new Event('submit'));
			}
		});

		async function convert() {
			const inputValue = input.value.trim();
			if (!inputValue) {
				showError('请输入时间或时间戳');
				return;
			}

			showLoading(true);
			hideError();
			result.classList.remove('show');

			try {
				const response = await fetch('/api/convert', {
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
					},
					body: JSON.stringify({
						input: inputValue,
						timezone: timezone.value,
					}),
				});

				const data = await response.json();

				if (data.success) {
					displayResult(data.data);
				} else {
					showError(data.error || '转换失败');
				}
			} catch (err) {
				showError('请求失败: ' + err.message);
			} finally {
				showLoading(false);
			}
		}

		async function fetchNow() {
			showLoading(true);
			hideError();
			result.classList.remove('show');

			try {
				const response = await fetch('/api/now?timezone=' + encodeURIComponent(timezone.value));
				const data = await response.json();

				if (data.success) {
					displayResult(data.data);
				} else {
					showError(data.error || '获取时间失败');
				}
			} catch (err) {
				showError('请求失败: ' + err.message);
			} finally {
				showLoading(false);
			}
		}

		function displayResult(data) {
			document.getElementById('iso8601').textContent = data.iso8601;
			document.getElementById('date').textContent = data.date;
			document.getElementById('time').textContent = data.time;
			document.getElementById('timestamp').textContent = data.timestamp;
			document.getElementById('timestamp_ms').textContent = data.timestamp_ms;
			document.getElementById('timestamp_us').textContent = data.timestamp_us;
			document.getElementById('weekday').textContent = data.weekday;
			document.getElementById('zero_timestamp').textContent = data.zero_timestamp;
			document.getElementById('timezone_info').textContent = data.timezone;

			result.classList.add('show');
		}

		function showError(msg) {
			error.textContent = msg;
			error.classList.add('show');
		}

		function hideError() {
			error.classList.remove('show');
		}

		function showLoading(show) {
			loading.style.display = show ? 'block' : 'none';
		}

		function copyToClipboard(id) {
			const element = document.getElementById(id);
			const text = element.textContent;

			navigator.clipboard.writeText(text).then(() => {
				successMsg.classList.add('show');
				setTimeout(() => {
					successMsg.classList.remove('show');
				}, 2000);
			}).catch(err => {
				console.error('复制失败:', err);
			});
		}

		// Auto focus input
		input.focus();
	</script>
</body>
</html>
`
