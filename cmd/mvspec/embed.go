package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mvadly/mvspec/internal/config"
)

var outputPath string

func runEmbed(cfg *config.Config) error {
	flag.StringVar(&outputPath, "o", "./mv-docs", "Output directory")
	flag.Parse()

	dir := outputPath
	if dir == "" {
		dir = "./mv-docs"
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	if _, err := os.Stat("mv-spec.json"); err != nil {
		return fmt.Errorf("run mvspec first")
	}

	content := getDocsContent(cfg)
	path := filepath.Join(dir, "docs.go")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return err
	}
	fmt.Printf("Generated %s\n", path)

	if err := writeDefaultFile(filepath.Join(dir, "index.html"), getDefaultIndexHTML()); err != nil {
		return err
	}
	if err := writeDefaultFile(filepath.Join(dir, "styles.css"), getDefaultStyles()); err != nil {
		return err
	}
	if err := writeDefaultFile(filepath.Join(dir, "app.js"), getDefaultAppJS(cfg)); err != nil {
		return err
	}

	fmt.Printf("\nUsage: r.GET(\"/mvdocs/*path\", gin.WrapH(mvdocs.MvHandler()))\n")
	fmt.Printf("Access: http://localhost:<port>/mvdocs\n")

	return nil
}

func writeDefaultFile(path, content string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return err
	}
	fmt.Printf("Generated %s\n", path)
	return nil
}

func getDocsContent(cfg *config.Config) string {
	return "// Package mvdocs provides embedded API documentation handler.\n" +
		"// GENERATED FILE - DO NOT EDIT\n" +
		"// Run 'mvspec embed' to regenerate\n" +
		"\n" +
		"package mvdocs\n" +
		"\n" +
		"import (\n" +
		"\t\"embed\"\n" +
		"\t\"net/http\"\n" +
		"\t\"os\"\n" +
		"\t\"strings\"\n" +
		"\t\"sync\"\n" +
		")\n" +
		"\n" +
		"//go:embed index.html styles.css app.js\n" +
		"var staticFiles embed.FS\n" +
		"\n" +
		"var specOnce sync.Once\n" +
		"var specData []byte\n" +
		"\n" +
		"// MvHandler returns HTTP handler for API documentation.\n" +
		"func MvHandler() http.Handler {\n" +
		"\treturn http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {\n" +
		"\t\tif !isDev(r) {\n" +
		"\t\t\thttp.NotFound(w, r)\n" +
		"\t\t\treturn\n" +
		"\t\t}\n" +
		"\t\tserve(w, r)\n" +
		"\t})\n" +
		"}\n" +
		"\n" +
		"func isDev(r *http.Request) bool {\n" +
		"\tif os.Getenv(\"MVSPEC_DEV_ONLY\") == \"true\" {\n" +
		"\t\treturn true\n" +
		"\t}\n" +
		"\tenv := os.Getenv(\"GO_ENV\")\n" +
		"\tif env == \"\" || env == \"development\" || env == \"local\" {\n" +
		"\t\treturn true\n" +
		"\t}\n" +
		"\treturn strings.HasPrefix(r.RemoteAddr, \"127.\")\n" +
		"}\n" +
		"\n" +
		"func serve(w http.ResponseWriter, r *http.Request) {\n" +
		"\tpath := r.URL.Path\n" +
		"\n" +
		"\tfilePath := strings.TrimPrefix(path, \"/mvdocs\")\n" +
		"\tfilePath = strings.TrimPrefix(filePath, \"/\")\n" +
		"\n" +
		"\tif filePath == \"\" || filePath == \"index.html\" {\n" +
		"\t\tw.Header().Set(\"Content-Type\", \"text/html\")\n" +
		"\t\tserveEmbedded(w, \"index.html\")\n" +
		"\t\treturn\n" +
		"\t}\n" +
		"\n" +
		"\tif filePath == \"mv-spec.json\" {\n" +
		"\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n" +
		"\t\tserveSpec(w)\n" +
		"\t\treturn\n" +
		"\t}\n" +
		"\n" +
		"\tstaticExt := map[string]string{\n" +
		"\t\t\"styles.css\": \"text/css\",\n" +
		"\t\t\"app.js\": \"application/javascript\",\n" +
		"\t\t\"index.html\": \"text/html\",\n" +
		"\t}\n" +
		"\n" +
		"\tfor filename, contentType := range staticExt {\n" +
		"\t\tif filePath == filename {\n" +
		"\t\t\tw.Header().Set(\"Content-Type\", contentType)\n" +
		"\t\t\tserveEmbedded(w, filename)\n" +
		"\t\t\treturn\n" +
		"\t\t}\n" +
		"\t}\n" +
		"\n" +
		"\tw.Header().Set(\"Content-Type\", \"text/html\")\n" +
		"\tserveEmbedded(w, \"index.html\")\n" +
		"}\n" +
		"\n" +
		"func serveEmbedded(w http.ResponseWriter, name string) {\n" +
		"\tdata, err := staticFiles.ReadFile(name)\n" +
		"\tif err != nil {\n" +
		"\t\tdata = []byte(\"<html><body><h1>MVAPI Docs</h1><p>Run mvspec embed first</p></body></html>\")\n" +
		"\t}\n" +
		"\tw.Write(data)\n" +
		"}\n" +
		"\n" +
		"func serveSpec(w http.ResponseWriter) {\n" +
		"\tspecOnce.Do(func() {\n" +
		"\t\tpaths := []string{\"mv-spec.json\", \"./mv-spec.json\", \"../mv-spec.json\"}\n" +
		"\t\tfor _, p := range paths {\n" +
		"\t\t\tif d, err := os.ReadFile(p); err == nil {\n" +
		"\t\t\t\tspecData = d\n" +
		"\t\t\t\treturn\n" +
		"\t\t\t}\n" +
		"\t\t}\n" +
		"\t})\n" +
		"\tif len(specData) == 0 {\n" +
		"\t\tspecData = []byte(\"{\\\"openapi\\\":\\\"3.0.0\\\",\\\"info\\\":{\\\"title\\\":\\\"API\\\"},\\\"paths\\\":{}}\")\n" +
		"\t}\n" +
		"\tw.Write(specData)\n" +
		"}\n" +
		"\n" +
		"func ReloadSpec() {\n" +
		"\tspecOnce = sync.Once{}\n" +
		"\tspecData = nil\n" +
		"}\n"
}

func getDefaultIndexHTML() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>MVAPI Docs</title>
  <link rel="stylesheet" href="styles.css">
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=Fira+Code:wght@400;500&display=swap" rel="stylesheet">
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.16/codemirror.min.css">
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.16/theme/dracula.min.css">
</head>
<body>
  <div class="app">
    <!-- Sidebar -->
    <aside class="sidebar" id="sidebar">
      <div class="sidebar-header">
        <h1 class="logo">MV<span>API</span></h1>
        <button class="env-btn" id="envBtn" title="Environment Variables">⚙</button>
      </div>
      <div class="sidebar-search">
        <input type="text" id="searchInput" placeholder="Search endpoints..." />
      </div>
      <nav class="collections" id="collections"></nav>
      <div class="sidebar-section">
        <h3 class="section-title" id="historyToggle">History <span class="chevron">▾</span></h3>
        <ul class="history-list" id="historyList"></ul>
      </div>
    </aside>

    <!-- Main Content -->
    <main class="main">
      <!-- API Header -->
      <div class="api-header" id="apiHeader">
        <h1 class="api-title" id="apiTitle">API Title</h1>
        <p class="api-description" id="apiDescription">API description</p>
      </div>

      <!-- Request Bar (always top, not affected by toggle) -->
      <div class="request-bar-wrapper">
        <button id="sidebarToggle" class="sidebar-toggle" onclick="toggleSidebar()" title="Toggle Sidebar">
          <svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5">
            <line x1="3" y1="4" x2="13" y2="4"></line>
            <line x1="3" y1="8" x2="13" y2="8"></line>
            <line x1="3" y1="12" x2="13" y2="12"></line>
          </svg>
        </button>
        <select id="methodSelect" class="method-select">
          <option value="GET">GET</option>
          <option value="POST">POST</option>
          <option value="PUT">PUT</option>
          <option value="PATCH">PATCH</option>
          <option value="DELETE">DELETE</option>
          <option value="HEAD">HEAD</option>
          <option value="OPTIONS">OPTIONS</option>
        </select>
        <select id="serverSelect" class="server-select" style="display:none"></select>
        <input type="text" id="urlInput" class="url-input" placeholder="Enter request URL or path..." />
        <button id="sendBtn" class="send-btn">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="22" y1="2" x2="11" y2="13"></line>
            <polygon points="22 2 15 22 11 13 2 9 22 2"></polygon>
          </svg>
        </button>
      </div>

      <!-- Panels (affected by toggle) -->
      <div class="panels-layout" id="panelsLayout">
        <!-- Column 1: Request Panel -->
        <div class="col-request">
          <div class="request-panel">
            <div class="tabs">
              <button class="tab active" data-tab="auth">Auth</button>
              <button class="tab" data-tab="headers">Header</button>
              <button class="tab" data-tab="body">Body</button>
              <button class="tab" data-tab="examples">Example</button>
            </div>
            <div class="tab-content" id="authTab">
              <div class="auth-section">
                <label class="auth-label">Type</label>
                <select id="authType" class="auth-type-select">
                  <option value="none">No Auth</option>
                  <option value="basic">Basic Auth</option>
                  <option value="bearer">Bearer Token</option>
                  <option value="custom">Custom</option>
                </select>
              </div>
              <div id="basicAuthFields" class="auth-fields hidden">
                <div class="auth-section">
                  <label class="auth-label">Username</label>
                  <input type="text" id="basicUsername" class="auth-input" placeholder="Username" />
                </div>
                <div class="auth-section">
                  <label class="auth-label">Password</label>
                  <input type="password" id="basicPassword" class="auth-input" placeholder="Password" />
                </div>
              </div>
              <div id="bearerAuthFields" class="auth-fields hidden">
                <div class="auth-section">
                  <label class="auth-label">Token</label>
                  <input type="text" id="bearerToken" class="auth-input" placeholder="Bearer token" />
                </div>
              </div>
              <div id="customAuthFields" class="auth-fields hidden">
                <div class="auth-section">
                  <label class="auth-label">Header Name</label>
                  <input type="text" id="customAuthKey" class="auth-input" placeholder="Authorization" />
                </div>
                <div class="auth-section">
                  <label class="auth-label">Header Value</label>
                  <input type="text" id="customAuthValue" class="auth-input" placeholder="Token value" />
                </div>
              </div>
            </div>
            <div class="tab-content hidden" id="headersTab">
              <div class="kv-editor" id="paramsEditor">
                <div class="kv-row">
                  <input type="text" placeholder="Key" class="kv-key" />
                  <input type="text" placeholder="Value" class="kv-value" />
                  <button class="kv-remove">×</button>
                </div>
              </div>
              <button class="add-row-btn" data-editor="paramsEditor">+ Add Param</button>
              <div class="kv-editor" id="headersEditor" style="margin-top:12px">
                <div class="kv-row">
                  <input type="text" placeholder="Header Key" class="kv-key" />
                  <input type="text" placeholder="Header Value" class="kv-value" />
                  <button class="kv-remove">×</button>
                </div>
              </div>
              <button class="add-row-btn" data-editor="headersEditor">+ Add Header</button>
            </div>
            <div class="tab-content hidden" id="bodyTab">
              <input type="hidden" id="contentTypeSelect" value="form" />
              <div id="jsonBodyEditor" class="body-editor hidden" placeholder='{ "key": "value" }'></div>
              <div id="formBodyEditor" class="form-body-editor">
                <div class="kv-editor" id="formEditor"></div>
                <button class="add-row-btn" data-editor="formEditor">+ Add Field</button>
              </div>
              <div id="formDataEditor" class="form-body-editor hidden">
                <div class="kv-editor" id="formDataFields"></div>
                <button class="add-row-btn" data-editor="formDataFields">+ Add Field</button>
                <div class="file-input-row">
                  <input type="text" class="kv-key" placeholder="Field name (e.g., file)" />
                  <input type="file" class="file-input" id="fileInput" />
                </div>
              </div>
            </div>
            <div class="tab-content hidden" id="examplesTab">
              <pre id="examplesOutput" class="response-output"></pre>
            </div>
          </div>
        </div>

        <!-- Column 2: Response Panel -->
        <div class="col-response">
          <div class="response-panel">
            <div class="response-meta" id="responseMeta">
              <span class="response-status" id="responseStatus"></span>
              <span class="response-time" id="responseTime"></span>
              <span class="response-size" id="responseSize"></span>
            </div>
            <div class="response-tabs">
              <button class="tab active" data-restab="responseResult">Result</button>
              <button class="tab" data-restab="responseHeaders">Headers</button>
              <button class="tab" data-restab="requestSent">Request</button>
            </div>
            <div class="response-content" id="responseResult">
              <pre id="responseOutput" class="response-output"><code>No response yet. Send a request to see results.</code></pre>
            </div>
            <div class="response-content hidden" id="responseHeaders">
              <pre id="responseHeadersOutput" class="response-output"><code>No response headers</code></pre>
            </div>
            <div class="response-content hidden" id="requestSent">
              <pre id="requestSentOutput" class="response-output"><code>No request headers</code></pre>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>

  <!-- Environment Variables Modal -->
  <div class="modal-overlay hidden" id="envModal">
    <div class="modal">
      <div class="modal-header">
        <h2>Environment Variables</h2>
        <button class="modal-close" id="envClose">×</button>
      </div>
      <div class="modal-body">
        <p class="modal-hint">Use <code>{{variable}}</code> in URLs, headers, and body.</p>
        <div class="kv-editor" id="envEditor">
          <div class="kv-row">
            <input type="text" placeholder="Variable name" class="kv-key" />
            <input type="text" placeholder="Value" class="kv-value" />
            <button class="kv-remove">×</button>
          </div>
        </div>
        <button class="add-row-btn" data-editor="envEditor">+ Add Variable</button>
      </div>
<div class="modal-footer">
          <button class="env-reset-btn" id="envReset">Reset to Defaults</button>
          <button class="send-btn" id="envSave">Save</button>
        </div>
    </div>
  </div>

  <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.16/codemirror.min.js"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.16/mode/javascript/javascript.min.js"></script>
  <script src="app.js"></script>
</body>
</html>`
}

func getDefaultStyles() string {
	return `*,*::before,*::after{box-sizing:border-box;margin:0;padding:0}
:root{
  --bg:#0F172A;--surface:rgba(16,185,129,.06);--surface-hover:rgba(16,185,129,.12);
  --glass:rgba(16,185,129,.08);--glass-border:rgba(16,185,129,.18);
  --primary:#10B981;--primary-glow:rgba(16,185,129,.35);
  --text:#E2E8F0;--text-dim:#94A3B8;--text-muted:#64748B;
  --danger:#EF4444;--warning:#F59E0B;--info:#3B82F6;
  --font:'Inter',system-ui,-apple-system,sans-serif;
  --mono:'Fira Code','Cascadia Code','Consolas',monospace;
  --radius:10px;--radius-sm:6px;
}
html,body{height:100%;background:var(--bg);color:var(--text);font-family:var(--font);font-size:14px;overflow:hidden}
.app{display:flex;height:100vh}

/* Sidebar */
.sidebar{width:280px;min-width:280px;background:rgba(15,23,42,.92);border-right:1px solid var(--glass-border);display:flex;flex-direction:column;backdrop-filter:blur(20px);-webkit-backdrop-filter:blur(20px);transition:width .3s}
.sidebar.hidden{width:0;min-width:0;display:none;overflow:hidden;border:none}
.sidebar-toggle{padding:10px;background:var(--glass);border:1px solid var(--glass-border);border-radius:var(--radius-sm);cursor:pointer;color:var(--text-dim);transition:all .2s;display:flex;align-items:center;justify-content:center}
.sidebar-toggle:hover{background:var(--surface-hover);color:var(--primary)}

.env-btn{background:var(--glass);border:1px solid var(--glass-border);color:var(--text-dim);width:32px;height:32px;border-radius:var(--radius-sm);cursor:pointer;font-size:14px;transition:all .2s;display:flex;align-items:center;justify-content:center}
.env-btn:hover{background:var(--surface-hover);color:var(--primary);border-color:var(--primary)}
.sidebar-header{display:flex;align-items:center;justify-content:space-between;padding:12px 16px}
.logo{font-size:18px;font-weight:700;color:var(--text);margin:0}
.logo span{color:var(--primary)}
.sidebar-search{padding:0 12px 12px}
.sidebar-search input{width:100%;padding:8px 12px;background:var(--glass);border:1px solid var(--glass-border);border-radius:var(--radius-sm);color:var(--text);font-size:13px;outline:none;transition:border-color .2s}
.sidebar-search input:focus{border-color:var(--primary);box-shadow:0 0 0 2px var(--primary-glow)}
.collections{flex:1;overflow-y:auto;padding:4px 8px}
.collections::-webkit-scrollbar{width:4px}
.collections::-webkit-scrollbar-thumb{background:var(--glass-border);border-radius:4px}
.tag-group{margin-bottom:4px}
.tag-name{display:flex;align-items:center;gap:6px;padding:6px 8px;font-size:12px;font-weight:600;color:var(--text-dim);text-transform:uppercase;letter-spacing:.5px;cursor:pointer;border-radius:var(--radius-sm);transition:background .15s}
.tag-name:hover{background:var(--surface)}
.tag-name .chevron{font-size:10px;transition:transform .2s}
.tag-name.collapsed .chevron{transform:rotate(-90deg)}
.endpoint-list{list-style:none}
.tag-name.collapsed+.endpoint-list{display:none}
.endpoint{display:flex;align-items:center;gap:8px;padding:6px 8px 6px 16px;cursor:pointer;border-radius:var(--radius-sm);transition:all .15s;font-size:13px}
.endpoint:hover{background:var(--surface-hover)}
.endpoint.active{background:var(--primary);background:linear-gradient(135deg,rgba(16,185,129,.18),rgba(16,185,129,.08));border:1px solid var(--glass-border)}
.method-badge{font-size:10px;font-weight:700;font-family:var(--mono);padding:2px 6px;border-radius:3px;min-width:42px;text-align:center;text-transform:uppercase}
.method-badge.get{background:rgba(16,185,129,.15);color:#34D399}
.method-badge.post{background:rgba(59,130,246,.15);color:#60A5FA}
.method-badge.put{background:rgba(245,158,11,.15);color:#FBBF24}
.method-badge.patch{background:rgba(168,85,247,.15);color:#C084FC}
.method-badge.delete{background:rgba(239,68,68,.15);color:#F87171}
.method-badge.head{background:rgba(148,163,184,.15);color:#CBD5E1}
.method-badge.options{background:rgba(148,163,184,.15);color:#CBD5E1}
.endpoint-path{color:var(--text-dim);font-family:var(--mono);font-size:12px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.sidebar-section{border-top:1px solid var(--glass-border);padding:8px}
.section-title{display:flex;align-items:center;justify-content:space-between;font-size:12px;font-weight:600;color:var(--text-dim);text-transform:uppercase;letter-spacing:.5px;padding:6px 8px;cursor:pointer;border-radius:var(--radius-sm)}
.section-title:hover{background:var(--surface)}
.history-list{list-style:none;max-height:180px;overflow-y:auto}
.history-list::-webkit-scrollbar{width:4px}
.history-list::-webkit-scrollbar-thumb{background:var(--glass-border);border-radius:4px}
.history-item{display:flex;align-items:center;gap:8px;padding:4px 8px 4px 16px;cursor:pointer;border-radius:var(--radius-sm);font-size:12px;color:var(--text-dim);transition:background .15s}
.history-item:hover{background:var(--surface-hover);color:var(--text)}

/* Main */
.main{flex:1;display:flex;flex-direction:column;overflow:hidden;padding:16px 20px;gap:12px;width:100%;min-width:0}

/* API Header */
.api-header{padding:0 0 12px;border-bottom:1px solid var(--glass-border);margin-bottom:4px}
.api-header.hidden{display:none}
.api-title{font-size:20px;font-weight:700;color:var(--text);margin:0 0 4px}
.api-description{font-size:13px;color:var(--text-dim);margin:0;line-height:1.4}

/* Request Bar Wrapper */
.request-bar-wrapper{display:flex;gap:8px;align-items:center;flex-shrink:0;width:100%}

/* Panels Layout (affected by toggle) */
.panels-layout{display:flex;flex:1;gap:12px;min-height:0;width:100%;min-width:0}
.col-request{flex:1;min-width:0;display:flex;flex-direction:column}
.col-response{min-width:300px;max-width:400px;flex-shrink:0;display:flex;flex-direction:column}
.col-response .response-panel{flex:1;display:flex;flex-direction:column;min-height:0}

/* Vertical layout (stacked) */
.panels-layout.layout-vertical{flex-direction:column}
.panels-layout.layout-vertical .col-request{flex:none;min-height:150px;overflow:hidden}
.panels-layout.layout-vertical .col-response{flex:none;min-height:150px;max-height:40vh}
.panels-layout.layout-vertical .request-panel{height:100vh;overflow-y:auto}

/* Layout Toggle */
.layout-toggle{padding:6px 10px;background:var(--glass);border:1px solid var(--glass-border);border-radius:var(--radius-sm);cursor:pointer;color:var(--text-dim);font-size:16px;transition:all .2s;display:flex;align-items:center;justify-content:center}
.layout-toggle:hover{background:var(--surface-hover);border-color:var(--primary);color:var(--text)}
.layout-icon{stroke:currentColor;fill:none}
#layoutToggle[data-layout="horizontal"] .icon-stacked{display:none}
#layoutToggle[data-layout="vertical"] .icon-side-by-side{display:none}
#layoutToggle:not([data-layout]) .icon-stacked{display:none}

/* Request Bar */
.method-select,.server-select,.auth-type-select{padding:10px 32px 10px 12px;background:var(--glass);border:1px solid var(--glass-border);color:var(--text);font-family:var(--mono);font-size:13px;border-radius:var(--radius-sm);cursor:pointer;outline:none;appearance:none;-webkit-appearance:none;transition:border-color .2s;background-image:url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 24 24' fill='none' stroke='%238b9bb4' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");background-repeat:no-repeat;background-position:right 10px center}
.method-select{color:var(--primary);font-weight:600;min-width:100px}
.method-select:focus,.server-select:focus,.auth-type-select:focus{border-color:var(--primary);box-shadow:0 0 0 2px var(--primary-glow)}
.method-select option,.server-select option,.auth-type-select option{background:var(--bg);color:var(--text)}
.server-select{min-width:150px}
.auth-type-select{width:100%}
.url-input{flex:1;padding:10px 14px;background:var(--glass);border:1px solid var(--glass-border);color:var(--text);font-family:var(--mono);font-size:13px;border-radius:var(--radius-sm);outline:none;transition:border-color .2s}
.url-input:focus{border-color:var(--primary);box-shadow:0 0 0 2px var(--primary-glow)}
.url-input::placeholder{color:var(--text-muted)}
.send-btn{display:flex;align-items:center;gap:6px;padding:10px 24px;background:var(--primary);color:#fff;border:none;border-radius:var(--radius-sm);font-weight:600;font-size:13px;cursor:pointer;transition:all .2s;text-transform:uppercase;letter-spacing:.5px}
.send-btn:hover{box-shadow:0 0 20px var(--primary-glow);transform:translateY(-1px)}
.send-btn:active{transform:translateY(0)}
.send-btn.loading{opacity:.7;pointer-events:none}
.send-btn .spin{animation:spin 1s linear infinite}
@keyframes spin{to{transform:rotate(360deg)}}
.tabs,.response-tabs{display:flex;gap:2px;border-bottom:1px solid var(--glass-border);padding-bottom:0}
.tab{padding:8px 16px;background:none;border:none;border-bottom:2px solid transparent;color:var(--text-dim);font-size:13px;font-weight:500;cursor:pointer;transition:all .2s}
.tab:hover{color:var(--text)}
.tab.active{color:var(--primary);border-bottom-color:var(--primary)}

/* Request Panel */
.request-panel{background:var(--glass);border:1px solid var(--glass-border);border-radius:var(--radius);padding:12px;backdrop-filter:blur(12px);-webkit-backdrop-filter:blur(12px)}
.tab-content{padding-top:10px}
.tab-content.hidden{display:none}
#examplesTab.tab-content{overflow-y:auto;max-height:calc(100vh - 180px)}
.auth-section{margin-bottom:12px}
.auth-label{display:block;font-size:12px;font-weight:500;color:var(--text-dim);margin-bottom:4px}
.auth-input{width:100%;padding:8px 12px;background:rgba(15,23,42,.6);border:1px solid var(--glass-border);border-radius:var(--radius-sm);color:var(--text);font-size:13px;outline:none;transition:border-color .2s}
.auth-input:focus{border-color:var(--primary)}
.auth-fields.hidden{display:none}
.kv-editor{display:flex;flex-direction:column;gap:6px}
.kv-row{display:flex;gap:6px;align-items:center}
.kv-key,.kv-value{flex:1;padding:6px 10px;background:rgba(15,23,42,.6);border:1px solid var(--glass-border);border-radius:var(--radius-sm);color:var(--text);font-size:12px;font-family:var(--mono);outline:none;transition:border-color .2s}
.kv-key:focus,.kv-value:focus{border-color:var(--primary)}
.kv-remove{background:none;border:none;color:var(--text-muted);cursor:pointer;font-size:16px;padding:4px 8px;border-radius:var(--radius-sm);transition:all .15s}
.kv-remove:hover{color:var(--danger);background:rgba(239,68,68,.1)}
.add-row-btn{background:none;border:1px dashed var(--glass-border);color:var(--text-dim);padding:6px 12px;border-radius:var(--radius-sm);cursor:pointer;font-size:12px;margin-top:6px;transition:all .2s;width:100%}
.add-row-btn:hover{border-color:var(--primary);color:var(--primary)}
.body-editor{min-height:150px;height:200px}
.body-editor .CodeMirror{height:100%;font-family:var(--mono);font-size:13px;border-radius:var(--radius-sm);background:rgba(15,23,42,.6);color:var(--text)}
.body-editor .CodeMirror-cursor{border-left:1px solid var(--primary)}
.body-editor .CodeMirror-selected{background:var(--glass-border)}
.body-editor .CodeMirror-gutters{background:rgba(15,23,42,.8);border-right:1px solid var(--glass-border)}
.body-editor .CodeMirror-linenumber{padding:0 8px;color:var(--text-muted)}
.body-editor .cm-string{color:#CE9178}
.body-editor .cm-number{color:#B5CEA8}
.body-editor .cm-key{color:#9CDCFE}
.body-editor .cm-property{color:#9CDCFE}
.body-editor .cm-atom{color:#569CD6}
.body-editor .cm-bool{color:#569CD6}
.body-editor .cm-null{color:#569CD6}

.content-type-selector{display:flex;align-items:center;gap:10px;margin-bottom:10px}
.content-type-select{padding:6px 28px 6px 10px;background:var(--glass);border:1px solid var(--glass-border);color:var(--text);font-size:12px;border-radius:var(--radius-sm);cursor:pointer;appearance:none;-webkit-appearance:none;background-image:url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='10' height='10' viewBox='0 0 24 24' fill='none' stroke='%238b9bb4' stroke-width='2'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");background-repeat:no-repeat;background-position:right 8px center}
.form-body-editor{margin-top:10px}
.form-body-editor.hidden{display:none}
.file-input-row{display:flex;gap:6px;align-items:center;margin-top:10px;padding:8px;background:rgba(15,23,42,.4);border-radius:var(--radius-sm)}
.file-input{flex:1;padding:6px 10px;background:rgba(15,23,42,.6);border:1px solid var(--glass-border);border-radius:var(--radius-sm);color:var(--text);font-size:12px}
.file-input::-webkit-file-upload-button{padding:4px 8px;background:var(--glass);border:1px solid var(--glass-border);color:var(--text);border-radius:var(--radius-sm);cursor:pointer;font-size:11px;margin-right:8px}

/* Response Panel */
.response-panel{flex:1;background:var(--glass);border:1px solid var(--glass-border);border-radius:var(--radius);padding:12px;backdrop-filter:blur(12px);-webkit-backdrop-filter:blur(12px);display:flex;flex-direction:column;overflow:hidden;min-height:0}
.response-meta{display:flex;gap:16px;padding-bottom:8px;font-size:12px;font-family:var(--mono);align-items:center}
.response-status{padding:3px 10px;border-radius:var(--radius-sm);font-weight:600}
.response-status.s2xx{background:rgba(16,185,129,.15);color:#34D399}
.response-status.s3xx{background:rgba(59,130,246,.15);color:#60A5FA}
.response-status.s4xx{background:rgba(245,158,11,.15);color:#FBBF24}
.response-status.s5xx{background:rgba(239,68,68,.15);color:#F87171}
.response-time{color:var(--text-dim)}
.response-size{color:var(--text-dim)}
.response-content{flex:1;overflow-y:auto;min-height:0;display:flex;flex-direction:column}
.response-content.hidden{display:none}
.response-content::-webkit-scrollbar{width:6px}
.response-content::-webkit-scrollbar-track{background:transparent}
.response-content::-webkit-scrollbar-thumb{background:var(--glass-border);border-radius:4px}
.response-output{flex:1;margin:0;padding:12px;font-family:var(--mono);font-size:13px;line-height:1.6;white-space:pre-wrap;word-break:break-word;color:var(--text-dim);max-height:100%;overflow-y:auto}
.response-example{margin-bottom:16px;padding:12px;background:var(--glass);border-radius:var(--radius-sm);border:1px solid var(--glass-border)}
.response-code{display:inline-block;padding:2px 8px;border-radius:var(--radius-sm);font-weight:600;font-size:12px;margin-right:8px}
.response-code.s2xx{background:rgba(16,185,129,.15);color:#34D399}
.response-code.s3xx{background:rgba(59,130,246,.15);color:#60A5FA}
.response-code.s4xx{background:rgba(245,158,11,.15);color:#FBBF24}
.response-code.s5xx{background:rgba(239,68,68,.15);color:#F87171}
.response-desc{color:var(--text-dim);font-size:12px;margin-bottom:8px}
.response-header-row{display:flex;align-items:center;gap:8px;margin-bottom:8px}
.response-desc-text{flex:1;color:var(--text-dim);font-size:12px}
.try-btn{padding:4px 12px;background:linear-gradient(135deg,rgba(16,185,129,0.8),rgba(16,185,129,0.6));color:#fff;border:none;border-radius:var(--radius-sm);cursor:pointer;font-size:11px;font-weight:500}
.try-btn:hover{opacity:0.9}
.example-json{margin:8px 0 0;padding:8px;background:rgba(0,0,0,.2);border-radius:var(--radius-sm);font-size:12px;white-space:pre-wrap}
.example-label{font-weight:600;margin:12px 0 4px;color:var(--text)}
.result-section{margin-bottom:16px;padding-bottom:16px;border-bottom:1px solid var(--glass-border)}
.result-section:last-child{border-bottom:none;margin-bottom:0;padding-bottom:0}
.curl-command{background:rgba(0,0,0,.4);padding:10px;border-radius:var(--radius-sm);font-size:12px;overflow-x:auto}

/* JSON Highlighting */
.json-key{color:#34D399}
.json-string{color:#FBBF24}
.json-number{color:#60A5FA}
.json-boolean{color:#C084FC}
.json-null{color:#F87171}

/* Modal */
.modal-overlay{position:fixed;inset:0;background:rgba(0,0,0,.6);display:flex;align-items:center;justify-content:center;z-index:100;backdrop-filter:blur(4px)}
.modal-overlay.hidden{display:none}
.modal{background:rgba(15,23,42,.95);border:1px solid var(--glass-border);border-radius:var(--radius);padding:0;width:500px;max-width:90vw;max-height:80vh;display:flex;flex-direction:column;backdrop-filter:blur(20px);box-shadow:0 25px 50px rgba(0,0,0,.4)}
.modal-header{display:flex;align-items:center;justify-content:space-between;padding:16px 20px;border-bottom:1px solid var(--glass-border)}
.modal-header h2{font-size:16px;font-weight:600;color:var(--text)}
.modal-close{background:none;border:none;color:var(--text-dim);font-size:20px;cursor:pointer;padding:4px 8px;border-radius:var(--radius-sm)}
.modal-close:hover{color:var(--danger)}
.modal-body{padding:16px 20px;overflow-y:auto;flex:1}
.modal-hint{font-size:12px;color:var(--text-dim);margin-bottom:12px}
.modal-hint code{background:var(--glass);padding:2px 6px;border-radius:3px;font-family:var(--mono);color:var(--primary)}
.modal-footer{padding:12px 20px;border-top:1px solid var(--glass-border);display:flex;justify-content:space-between}
.env-reset-btn{background:transparent;border:1px solid var(--glass-border);color:var(--text-dim);padding:8px 16px;border-radius:var(--radius-sm);cursor:pointer;font-size:13px;transition:all .2s}
.env-reset-btn:hover{background:var(--surface-hover);color:var(--danger);border-color:var(--danger)}

/* Responsive */
@media(max-width:1200px){
  .col-response{min-width:250px;max-width:300px}
}
@media(max-width:992px){
  .sidebar{width:240px;min-width:240px}
  .col-response{min-width:200px;max-width:250px}
}
@media(max-width:768px){
  .sidebar{width:200px;min-width:200px}
  .method-badge{min-width:36px;font-size:9px}
  .col-response{min-width:180px;max-width:200px}
  .request-bar-wrapper{flex-wrap:wrap}
  .server-select{min-width:120px}
  .url-input{min-width:100px}
  .sidebar-toggle,#sidebarToggle{width:36px;height:36px;padding:8px}
  .send-btn{padding:10px 16px}
  .main{padding:12px;gap:8px}
}
@media(max-width:600px){
  .app{flex-direction:column}
  .sidebar{width:100%;min-width:100%;max-height:35vh;border-right:none;border-bottom:1px solid var(--glass-border)}
  .main{flex-direction:column;padding:8px;gap:8px}
  .request-bar-wrapper{flex-wrap:wrap;gap:8px}
  .sidebar-toggle,#sidebarToggle{width:100%;height:40px;padding:8px;order:1}
  .method-select,.server-select,.url-input,#sendBtn{flex:1 1 100%;min-width:100%;width:100%}
  .method-select{order:2;flex:0 0 80px}
  .server-select{order:3}
  .url-input{order:4}
  #sendBtn{order:5;justify-content:center;width:100%}
  .logo{font-size:16px}
  .sidebar-header{padding:10px 12px}
  .sidebar-toggle,#sidebarToggle{width:36px;height:36px;padding:8px}
  .sendBtn{padding:8px 16px}
  .col-response{min-width:100%;max-width:100%}
  .panels-layout{flex-direction:column}
  .col-request,.col-response{min-height:200px}
  .tabs,.response-tabs{overflow-x:auto;flex-wrap:nowrap;-webkit-overflow-scrolling:touch}
  .tab{white-space:nowrap;padding:8px 12px;font-size:12px}
  .request-panel,.response-panel{padding:8px}
  .auth-input,.kv-key,.kv-value{width:100%;min-width:0}
  .kv-row{flex-wrap:wrap}
  .kv-remove{flex-shrink:0}
  body{font-size:13px}
  .api-header{padding:0 0 8px;margin-bottom:4px}
  .api-title{font-size:18px}
  .api-description{font-size:12px}
}
@media(max-width:480px){
  .main{padding:8px;gap:6px}
  .panels-layout{gap:6px}
  .sidebar{width:100%;min-width:100%;max-height:30vh}
  .method-select{order:2;flex:0 0 70px}
  .send-btn{padding:10px 16px}
  .endpoint{padding:6px 8px 6px 12px;font-size:12px}
  .tag-name{font-size:11px}
  .modal{width:95%;max-width:none;margin:10px}
}`
}

func getDefaultAppJS(cfg *config.Config) string {
	// Build default env vars from config
	var defaultEnvVars string
	if len(cfg.EnvVars) > 0 {
		envMap := make(map[string]string)
		for _, e := range cfg.EnvVars {
			envMap[e.Name] = e.Value
		}
		b, _ := json.Marshal(envMap)
		defaultEnvVars = string(b)
	} else {
		defaultEnvVars = "{}"
	}

	return `(function(){
  "use strict";

  // --- State ---
  let spec = null;
  let currentEntry = null;
  let defaultEnvVars = ` + defaultEnvVars + `;
  let envVars = JSON.parse(localStorage.getItem("mvapi_env") || JSON.stringify(defaultEnvVars));
  let history = JSON.parse(localStorage.getItem("mvapi_history") || "[]");
  let bodyEditorInstance = null;

  // --- DOM refs ---
  const $ = (s) => document.querySelector(s);
  const $$ = (s) => document.querySelectorAll(s);

  const collectionsEl   = $("#collections");
  const searchInput     = $("#searchInput");
  const methodSelect    = $("#methodSelect");
  const serverSelect    = $("#serverSelect");
  const urlInput        = $("#urlInput");
  const sendBtn         = $("#sendBtn");
  const bodyEditor      = $("#bodyEditor");
  const responseOutput  = $("#responseOutput");
  const responseStatus  = $("#responseStatus");
  const responseTime    = $("#responseTime");
  const responseSize    = $("#responseSize");
  const examplesOutput = $("#examplesOutput");
  const historyList     = $("#historyList");
  const envBtn          = $("#envBtn");
  const envModal        = $("#envModal");
  const envClose        = $("#envClose");
  const envSave         = $("#envSave");
  const envReset        = $("#envReset");
  const envEditorEl     = $("#envEditor");
  const responseHeadersOutput = $("#responseHeadersOutput");
  const requestSentOutput = $("#requestSentOutput");

  // --- Init ---
  loadSpec();
  renderHistory();
  setupTabs();
  setupKVEditors();
  setupEnvModal();
  setupAuth();
  setupBodyEditor();

  sendBtn.addEventListener("click", sendRequest);
  urlInput.addEventListener("keydown", (e) => { if(e.key==="Enter") sendRequest(); });
  searchInput.addEventListener("input", filterCollections);

  // --- Load OpenAPI Spec ---
  function loadSpec() {
    const apiHeader = document.getElementById('apiHeader');
    const apiTitle = document.getElementById('apiTitle');
    const apiDescription = document.getElementById('apiDescription');
    
    fetch("mv-spec.json")
      .then((r) => r.json())
      .then((data) => { 
        spec = data; 
        
        // Populate API header
        if (spec.info) {
          apiTitle.textContent = spec.info.title || 'API';
          apiDescription.textContent = spec.info.description || '';
          apiDescription.style.display = spec.info.description ? 'block' : 'none';
          apiHeader.classList.remove('hidden');
        } else {
          apiHeader.classList.add('hidden');
        }
        
        renderCollections(); 
        renderServerDropdown();
      })
      .catch(() => {
        document.getElementById('apiHeader').classList.add('hidden');
        collectionsEl.innerHTML = '<p style="padding:12px;color:var(--text-dim);font-size:12px">Could not load mv-spec.json</p>';
      });
  }

  // --- Render Server Dropdown ---
  function renderServerDropdown() {
    if (spec.servers && spec.servers.length > 1) {
      serverSelect.innerHTML = "";
      spec.servers.forEach((server) => {
        const option = document.createElement('option');
        option.value = server.url;
        option.textContent = server.url;
        serverSelect.appendChild(option);
      });
      serverSelect.style.display = 'block';
    } else {
      serverSelect.style.display = 'none';
    }
  }

  // Handle server change
  serverSelect.addEventListener('change', function() {
    if (currentEntry) {
      urlInput.value = serverSelect.value + currentEntry.path;
    }
  });

  // --- Render Collections ---
  function renderCollections(filter) {
    if (!spec || !spec.paths) return;
    collectionsEl.innerHTML = "";
    const tagged = {};
    const untagged = [];

    for (const [path, methods] of Object.entries(spec.paths)) {
      for (const [method, op] of Object.entries(methods)) {
        if (typeof op !== "object" || !op) continue;
        const entry = { method: method.toUpperCase(), path, summary: op.summary || path, op };
        if (op.tags && op.tags.length > 0) {
          for (const tag of op.tags) {
            if (!tagged[tag]) tagged[tag] = [];
            tagged[tag].push(entry);
          }
        } else {
          untagged.push(entry);
        }
      }
    }

    const filterLower = (filter || "").toLowerCase();

    function matchesFilter(e) {
      if (!filterLower) return true;
      return e.path.toLowerCase().includes(filterLower) ||
             e.method.toLowerCase().includes(filterLower) ||
             (e.summary && e.summary.toLowerCase().includes(filterLower));
    }

    for (const [tag, entries] of Object.entries(tagged)) {
      const filtered = entries.filter(matchesFilter);
      if (filtered.length === 0) continue;
      renderTagGroup(tag, filtered);
    }
    if (untagged.length > 0) {
      const filtered = untagged.filter(matchesFilter);
      if (filtered.length > 0) renderTagGroup("Other", filtered);
    }
  }

  function renderTagGroup(tag, entries) {
    const group = document.createElement("div");
    group.className = "tag-group";
    const header = document.createElement("div");
    header.className = "tag-name";
    header.innerHTML = tag + ' <span class="chevron">▾</span>';
    header.addEventListener("click", () => header.classList.toggle("collapsed"));
    group.appendChild(header);

    const list = document.createElement("ul");
    list.className = "endpoint-list";
    for (const entry of entries) {
      const li = document.createElement("li");
      li.className = "endpoint";
      li.innerHTML = '<span class="method-badge ' + entry.method.toLowerCase() + '">' +
        entry.method + '</span><span class="endpoint-path" title="' + escapeAttr(entry.path) + '">' +
        escapeHTML(entry.summary) + '</span>';
      li.addEventListener("click", () => selectEndpoint(entry));
      list.appendChild(li);
    }
    group.appendChild(list);
    collectionsEl.appendChild(group);
  }

  function selectEndpoint(entry) {
    currentEntry = entry;
    $$(".endpoint.active").forEach((el) => el.classList.remove("active"));
    event.currentTarget.classList.add("active");
    methodSelect.value = entry.method;
    
    // Use selected server or first server from spec
    let basePath = "";
    if (serverSelect.style.display !== 'none' && serverSelect.value) {
      basePath = serverSelect.value;
    } else {
      basePath = (spec.servers && spec.servers[0] && spec.servers[0].url) || "";
    }
    urlInput.value = basePath + entry.path;

    // Populate params from path parameters
    const paramsEditor = $("#paramsEditor");
    paramsEditor.innerHTML = "";
    if (entry.op.parameters) {
      for (const p of entry.op.parameters) {
        if (p.in === "query") addKVRow(paramsEditor, p.name, "", p.description || "");
      }
    }
    if (paramsEditor.children.length === 0) addKVRow(paramsEditor, "", "");

    // Populate headers
    const headersEditor = $("#headersEditor");
    headersEditor.innerHTML = "";
    // Content-Type will be set based on body type in sendRequest, don't add here
    if (entry.op.parameters) {
      for (const p of entry.op.parameters) {
        if (p.in === "header") addKVRow(headersEditor, p.name, "");
      }
    }

    // Clear body
    setBodyEditorValue("");
    
    // Detect content-type from endpoint and set body editors
    const jsonBodyEditor = document.getElementById('jsonBodyEditor');
    const formBodyEditor = document.getElementById('formBodyEditor');
    const formDataEditor = document.getElementById('formDataEditor');
    const contentTypeInput = document.getElementById('contentTypeSelect');
    
    jsonBodyEditor.classList.add('hidden');
    formBodyEditor.classList.add('hidden');
    formDataEditor.classList.add('hidden');
    
    let detectedType = 'form'; // default
    
    if (entry.op.requestBody && entry.op.requestBody.content) {
      const contentKeys = Object.keys(entry.op.requestBody.content);
      const hasJson = contentKeys.includes('application/json');
      const hasMultipart = contentKeys.some(k => k.includes('multipart'));
      const hasFormUrlEncoded = contentKeys.includes('application/x-www-form-urlencoded');
      
      if (hasJson) {
        detectedType = 'json';
        jsonBodyEditor.classList.remove('hidden');
        const jsonContent = entry.op.requestBody.content["application/json"];
        if (jsonContent && jsonContent.schema) {
          setBodyEditorValue(buildExampleBody(jsonContent.schema));
        }
      } else if (hasMultipart) {
        detectedType = 'form-data';
        formDataEditor.classList.remove('hidden');
        const formDataEl = document.getElementById('formDataFields');
        formDataEl.innerHTML = '';
        
        // Get schema from spec and populate form fields
        if (entry.op.requestBody && entry.op.requestBody.content && entry.op.requestBody.content["multipart/form-data"]) {
          const schemaRef = entry.op.requestBody.content["multipart/form-data"].schema;
          if (schemaRef && schemaRef.$ref) {
            const refName = schemaRef.$ref.split('/').pop();
            const schema = spec.components && spec.components.schemas && spec.components.schemas[refName];
            if (schema && schema.properties) {
              // Check for file inputs (format: binary) and create file input rows for each
              for (const [propName, propSchema] of Object.entries(schema.properties)) {
                const isFile = propSchema.format === 'binary' || propSchema.type === 'file';
                if (isFile) {
                  // Create file input row for this file field
                  const fileInputRow = document.createElement('div');
                  fileInputRow.className = 'file-input-row';
                  fileInputRow.innerHTML = '<input type="text" class="kv-key" placeholder="Field name" value="' + propName + '"><input type="file" class="file-input" />';
                  formDataEl.parentNode.insertBefore(fileInputRow, formDataEl.nextSibling);
                  continue;
                }
                addKVRow(formDataEl, propName, '');
              }
              if (Object.keys(schema.properties).length === 0) {
                addKVRow(formDataEl, '', '');
              }
            } else {
              addKVRow(formDataEl, '', '');
            }
          } else {
            addKVRow(formDataEl, '', '');
          }
        } else {
          addKVRow(formDataEl, '', '');
        }
      } else if (hasFormUrlEncoded) {
        detectedType = 'form';
        formBodyEditor.classList.remove('hidden');
        const formEl = document.getElementById('formEditor');
        formEl.innerHTML = '';
        addKVRow(formEl, '', '');
      } else {
        // Fallback: show form editor
        formBodyEditor.classList.remove('hidden');
        const formEl = document.getElementById('formEditor');
        formEl.innerHTML = '';
        addKVRow(formEl, '', '');
      }
    } else {
      // No requestBody, show form editor by default
      formBodyEditor.classList.remove('hidden');
      const formEl = document.getElementById('formEditor');
      formEl.innerHTML = '';
      addKVRow(formEl, '', '');
    }
    
    contentTypeInput.value = detectedType;
    
    if (bodyEditorInstance) {
      bodyEditorInstance.refresh();
    }

    // Populate examples (request + response together)
    examplesOutput.innerHTML = "";
    if (entry.op.responses) {
      let examplesHTML = "";
      for (const [code, resp] of Object.entries(entry.op.responses)) {
        const statusClass = code.startsWith("2") ? "s2xx" : code.startsWith("4") ? "s4xx" : code.startsWith("5") ? "s5xx" : "s3xx";
        const description = resp.description || "";
        
        examplesHTML += '<div class="response-example">';
        examplesHTML += '<div class="response-header-row">'; 
        examplesHTML += '<span class="response-code ' + statusClass + '">' + code + '</span>';
        examplesHTML += '<span class="response-desc-text">' + description + '</span>';
        if (resp.requestExample !== undefined) {
          examplesHTML += '<button class="try-btn" onclick="tryRequestExample(' + code + ')">Try</button>';
        }
        examplesHTML += '</div>';

        if (resp.requestExample !== undefined) {
          examplesHTML += '<div class="example-label">Request:</div>';
          examplesHTML += '<pre class="example-json">' + syntaxHighlight(JSON.stringify(resp.requestExample, null, 2)) + '</pre>';
        }

        if (resp.example !== undefined) {
          examplesHTML += '<div class="example-label">Response:</div>';
          examplesHTML += '<pre class="example-json">' + syntaxHighlight(JSON.stringify(resp.example, null, 2)) + '</pre>';
        }
        examplesHTML += '</div>';
      }
      examplesOutput.innerHTML = examplesHTML || '<div style="color:var(--text-dim);padding:12px">No examples defined</div>';
    }
  }

  function buildExampleBody(schema) {
    if (schema["$ref"]) {
      const refName = schema["$ref"].split("/").pop();
      // Try direct lookup first
      if (spec.components && spec.components.schemas && spec.components.schemas[refName]) {
        return buildExampleBody(spec.components.schemas[refName]);
      }
      // Try without package prefix (e.g., "dto.LoginDTO" -> "LoginDTO")
      const parts = refName.split(".");
      if (parts.length > 1) {
        const shortName = parts[parts.length - 1];
        if (spec.components && spec.components.schemas && spec.components.schemas[shortName]) {
          return buildExampleBody(spec.components.schemas[shortName]);
        }
      }
      return "{}";
    }
    if (schema.example) return JSON.stringify(schema.example, null, 2);
    if (schema.properties) {
      const obj = {};
      for (const [key, prop] of Object.entries(schema.properties)) {
        if (prop.example !== undefined) obj[key] = prop.example;
        else if (prop.type === "string") obj[key] = "";
        else if (prop.type === "integer" || prop.type === "number") obj[key] = 0;
        else if (prop.type === "boolean") obj[key] = false;
        else if (prop.type === "array") obj[key] = [];
        else if (prop.type === "object") obj[key] = {};
        else obj[key] = null;
      }
      return JSON.stringify(obj, null, 2);
    }
    return "{}";
  }

  // --- Try Request Example ---
  function tryRequestExample(code) {
    if (!currentEntry || !currentEntry.op.responses || !currentEntry.op.responses[code]) return;
    const resp = currentEntry.op.responses[code];
    if (resp.requestExample !== undefined) {
      setBodyEditorValue(JSON.stringify(resp.requestExample, null, 2));
      // Switch to Body tab
      document.querySelector('.request-panel .tabs .tab[data-tab="body"]').click();
    }
  }
  window.tryRequestExample = tryRequestExample;

  // --- Toggle Sidebar ---
  function toggleSidebar() {
    const sidebar = document.getElementById('sidebar');
    sidebar.classList.toggle('hidden');
    localStorage.setItem('mvapi_sidebar', sidebar.classList.contains('hidden') ? 'hidden' : 'visible');
  }
  window.toggleSidebar = toggleSidebar;

  // Load saved sidebar preference
  const savedSidebar = localStorage.getItem('mvapi_sidebar');
  if (savedSidebar === 'hidden') {
    document.getElementById('sidebar').classList.add('hidden');
  }

  // --- Send Request ---
  function sendRequest() {
    const method = methodSelect.value;
    let url = substituteEnv(urlInput.value.trim());
    if (!url) return;

    // Build query params
    const params = getKVPairs("paramsEditor");
    let finalUrl = url;
    if (params.length > 0 && params[0].key) {
      const qs = params.filter(p => p.key).map((p) => encodeURIComponent(p.key) + "=" + encodeURIComponent(p.value)).join("&");
      finalUrl += (url.includes("?") ? "&" : "?") + qs;
    }

    // Build headers
    const headerPairs = getKVPairs("headersEditor");
    const headers = {};
    for (const h of headerPairs) {
      if (h.key) headers[substituteEnv(h.key)] = substituteEnv(h.value);
    }

    // Add auth header
    const authType = document.getElementById('authType').value;
    if (authType === 'basic') {
      const username = document.getElementById('basicUsername').value;
      const password = document.getElementById('basicPassword').value;
      if (username) {
        const encoded = btoa(username + ':' + password);
        headers['Authorization'] = 'Basic ' + encoded;
      }
    } else if (authType === 'bearer') {
      const token = document.getElementById('bearerToken').value;
      if (token) {
        headers['Authorization'] = 'Bearer ' + token;
      }
    } else if (authType === 'custom') {
      const key = document.getElementById('customAuthKey').value;
      const value = document.getElementById('customAuthValue').value;
      if (key && value) {
        headers[key] = value;
      }
    }

    // Handle body based on content type
    const contentType = document.getElementById('contentTypeSelect').value;
    let body = null;
    
    if (method !== "GET" && method !== "HEAD") {
      if (contentType === 'json') {
        body = substituteEnv(getBodyEditorValue().trim());
        if (body) headers['Content-Type'] = 'application/json';
      } else if (contentType === 'form') {
        const formPairs = getKVPairs('formEditor');
        const formData = formPairs.filter(p => p.key).map(p => encodeURIComponent(p.key) + "=" + encodeURIComponent(p.value)).join("&");
        if (formData) {
          body = formData;
          headers['Content-Type'] = 'application/x-www-form-urlencoded';
        }
      } else if (contentType === 'form-data') {
        const formData = new FormData();
        const formPairs = getKVPairs('formDataFields');
        formPairs.forEach(p => {
          if (p.key) formData.append(p.key, p.value);
        });
        // Add file if selected
        const fileInput = document.getElementById('fileInput');
        if (fileInput && fileInput.files.length > 0) {
          formData.append(fileInput.files[0].name, fileInput.files[0]);
        }
        body = formData;
        delete headers['Content-Type']; // Let browser set multipart boundary
      }
    }

    // Build curl command
    let curl = "curl -X " + method;
    for (const [key, value] of Object.entries(headers)) {
      curl += " -H \"" + key + ": " + value + "\"";
    }
    let bodyStr = "";
    if (body && method !== "GET" && method !== "HEAD") {
      if (typeof body === 'string') {
        bodyStr = body;
      } else if (body instanceof FormData) {
        for (const [key, value] of body.entries()) {
          bodyStr += (bodyStr ? "&" : "") + key + "=" + value;
        }
      }
      if (bodyStr) curl += " -d '" + bodyStr.replace(/'/g, "'\\''") + "'";
    }
    curl += " " + finalUrl;

    const opts = { method, headers };
    if (method !== "GET" && method !== "HEAD") {
      if (body) opts.body = body;
    }

    sendBtn.classList.add("loading");
    sendBtn.innerHTML = '<svg class="spin" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10" stroke-opacity="0.25"/><path d="M12 2a10 10 0 0 1 10 10"/></svg>';
    const start = performance.now();

    fetch(url, opts)
      .then(async (res) => {
        const elapsed = Math.round(performance.now() - start);
        const text = await res.text();
        const size = new Blob([text]).size;
        showResponse(res.status, res.statusText, elapsed, size, text, res.headers, curl);
        addHistory(method, urlInput.value.trim(), res.status);
      })
      .catch((err) => {
        const elapsed = Math.round(performance.now() - start);
        showResponse(0, "Error", elapsed, 0, err.message, null, curl);
      })
      .finally(() => {
        sendBtn.classList.remove("loading");
        sendBtn.innerHTML = '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="22" y1="2" x2="11" y2="13"></line><polygon points="22 2 15 22 11 13 2 9 22 2"></polygon></svg>';
      });
  }

  function showResponse(status, statusText, time, size, body, headers, curl) {
    const statusClass = status >= 500 ? "s5xx" : status >= 400 ? "s4xx" : status >= 300 ? "s3xx" : status >= 200 ? "s2xx" : "";
    responseStatus.textContent = status ? status + " " + statusText : "Error";
    responseStatus.className = "response-status " + statusClass;
    responseTime.textContent = time + " ms";
    responseSize.textContent = formatSize(size);

    // Build result content with curl, response, and headers
    let resultHTML = '<div class="result-section">';
    resultHTML += '<div class="example-label">Curl:</div>';
    resultHTML += '<pre class="example-json curl-command">' + escapeHTML(curl) + '</pre>';
    resultHTML += '</div>';

    // Add Request Headers section
    const headerPairs = getKVPairs("headersEditor");
    if (headerPairs.length > 0) {
      resultHTML += '<div class="result-section">';
      resultHTML += '<div class="example-label">Request Headers:</div>';
      let reqHeadersText = "";
      for (const h of headerPairs) {
        if (h.key) reqHeadersText += substituteEnv(h.key) + ": " + substituteEnv(h.value) + "\n";
      }
      resultHTML += '<pre class="example-json">' + escapeHTML(reqHeadersText) + '</pre>';
      resultHTML += '</div>';
    }

    // Add Response Headers section
    if (headers) {
      resultHTML += '<div class="result-section">';
      resultHTML += '<div class="example-label">Response Headers:</div>';
      let respHeadersText = "";
      headers.forEach((v, k) => { respHeadersText += k + ": " + v + "\n"; });
      resultHTML += '<pre class="example-json">' + escapeHTML(respHeadersText) + '</pre>';
      resultHTML += '</div>';
    }

    // Add Response Body
    resultHTML += '<div class="result-section">';
    resultHTML += '<div class="example-label">Response:</div>';
    // Try to format as JSON
    let formatted;
    try {
      const parsed = JSON.parse(body);
      formatted = syntaxHighlight(JSON.stringify(parsed, null, 2));
    } catch(e) {
      formatted = escapeHTML(body);
    }
    resultHTML += '<pre class="example-json">' + formatted + '</pre>';
    resultHTML += '</div>';

    responseOutput.innerHTML = resultHTML;

    // Populate response headers (for Headers tab)
    if (headers) {
      let headersText = "";
      headers.forEach((v, k) => { headersText += k + ": " + v + "\n"; });
      responseHeadersOutput.innerHTML = '<code>' + escapeHTML(headersText) + '</code>';
    } else {
      responseHeadersOutput.innerHTML = '<code>No response headers</code>';
    }

    // Populate request headers (for Request tab) - reuse headerPairs from above
    if (headerPairs.length > 0) {
      let requestHeadersText = "";
      for (const h of headerPairs) {
        if (h.key) requestHeadersText += substituteEnv(h.key) + ": " + substituteEnv(h.value) + "\n";
      }
      requestSentOutput.innerHTML = '<code>' + escapeHTML(requestHeadersText) + '</code>';
    } else {
      requestSentOutput.innerHTML = '<code>No request headers</code>';
    }
  }

  function syntaxHighlight(json) {
    return json.replace(/&/g,"&amp;").replace(/</g,"&lt;").replace(/>/g,"&gt;")
      .replace(/("(\\u[a-fA-F0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?)/g, function(match) {
        if (/:$/.test(match)) return '<span class="json-key">' + match + '</span>';
        return '<span class="json-string">' + match + '</span>';
      })
      .replace(/\b(-?\d+(\.\d+)?([eE][+-]?\d+)?)\b/g, '<span class="json-number">$1</span>')
      .replace(/\b(true|false)\b/g, '<span class="json-boolean">$1</span>')
      .replace(/\bnull\b/g, '<span class="json-null">null</span>');
  }

  // --- History ---
  function addHistory(method, url, status) {
    history.unshift({ method, url, status, time: Date.now() });
    if (history.length > 50) history = history.slice(0, 50);
    localStorage.setItem("mvapi_history", JSON.stringify(history));
    renderHistory();
  }

  function renderHistory() {
    historyList.innerHTML = "";
    for (const h of history.slice(0, 20)) {
      const li = document.createElement("li");
      li.className = "history-item";
      li.innerHTML = '<span class="method-badge ' + h.method.toLowerCase() + '">' +
        h.method + '</span><span>' + escapeHTML(truncate(h.url, 30)) + '</span>';
      li.addEventListener("click", () => {
        methodSelect.value = h.method;
        urlInput.value = h.url;
      });
      historyList.appendChild(li);
    }
  }

  // --- Tabs ---
  function setupTabs() {
    $$(".tabs .tab").forEach((tab) => {
      tab.addEventListener("click", () => {
        tab.parentElement.querySelectorAll(".tab").forEach((t) => t.classList.remove("active"));
        tab.classList.add("active");
        const target = tab.dataset.tab;
        const panel = tab.closest(".request-panel");
        panel.querySelectorAll(".tab-content").forEach((c) => c.classList.add("hidden"));
        panel.querySelector("#" + target + "Tab").classList.remove("hidden");
        if (target === "body" && bodyEditorInstance) {
          setTimeout(() => bodyEditorInstance.refresh(), 50);
        }
      });
    });
    $$(".response-tabs .tab").forEach((tab) => {
      tab.addEventListener("click", () => {
        tab.parentElement.querySelectorAll(".tab").forEach((t) => t.classList.remove("active"));
        tab.classList.add("active");
        const target = tab.dataset.restab;
        const panel = tab.closest(".response-panel");
        panel.querySelectorAll(".response-content").forEach((c) => c.classList.add("hidden"));
        panel.querySelector("#" + target).classList.remove("hidden");
      });
    });
  }

  // --- KV Editors ---
  function setupKVEditors() {
    $$(".add-row-btn").forEach((btn) => {
      btn.addEventListener("click", () => {
        const editor = $("#" + btn.dataset.editor);
        addKVRow(editor, "", "");
      });
    });
    document.addEventListener("click", (e) => {
      if (e.target.classList.contains("kv-remove")) {
        const row = e.target.closest(".kv-row");
        const editor = row.parentElement;
        if (editor.children.length > 1) row.remove();
        else { row.querySelector(".kv-key").value = ""; row.querySelector(".kv-value").value = ""; }
      }
    });
  }

  function addKVRow(editor, key, value, placeholder) {
    const row = document.createElement("div");
    row.className = "kv-row";
    row.innerHTML = '<input type="text" placeholder="' + (placeholder || "Key") + '" class="kv-key" value="' + escapeAttr(key) + '" />' +
      '<input type="text" placeholder="Value" class="kv-value" value="' + escapeAttr(value) + '" />' +
      '<button class="kv-remove">×</button>';
    editor.appendChild(row);
  }

  function getKVPairs(editorId) {
    const rows = $$("#" + editorId + " .kv-row");
    const pairs = [];
    rows.forEach((row) => {
      const k = row.querySelector(".kv-key").value.trim();
      const v = row.querySelector(".kv-value").value.trim();
      if (k) pairs.push({ key: k, value: v });
    });
    return pairs;
  }

// --- Auth ---
  function setupAuth() {
    const authTypeSelect = document.getElementById('authType');
    const basicFields = document.getElementById('basicAuthFields');
    const bearerFields = document.getElementById('bearerAuthFields');
    const customFields = document.getElementById('customAuthFields');

    // Load saved auth
    const savedAuth = JSON.parse(localStorage.getItem('mvapi_auth') || '{}');
    if (savedAuth.type) {
      authTypeSelect.value = savedAuth.type;
    }
    if (savedAuth.username) document.getElementById('basicUsername').value = savedAuth.username;
    if (savedAuth.password) document.getElementById('basicPassword').value = savedAuth.password;
    if (savedAuth.token) document.getElementById('bearerToken').value = savedAuth.token;
    if (savedAuth.customKey) document.getElementById('customAuthKey').value = savedAuth.customKey;
    if (savedAuth.customValue) document.getElementById('customAuthValue').value = savedAuth.customValue;
    updateAuthFields();

    // Handle type change
    authTypeSelect.addEventListener('change', updateAuthFields);
  }

  function updateAuthFields() {
    const authType = document.getElementById('authType').value;
    document.getElementById('basicAuthFields').classList.toggle('hidden', authType !== 'basic');
    document.getElementById('bearerAuthFields').classList.toggle('hidden', authType !== 'bearer');
    document.getElementById('customAuthFields').classList.toggle('hidden', authType !== 'custom');

    // Save to localStorage
    const auth = {
      type: authType,
      username: document.getElementById('basicUsername').value,
      password: document.getElementById('basicPassword').value,
      token: document.getElementById('bearerToken').value,
      customKey: document.getElementById('customAuthKey').value,
      customValue: document.getElementById('customAuthValue').value
    };
    localStorage.setItem('mvapi_auth', JSON.stringify(auth));
  }

  // Add event listeners to auth inputs to save on change
  document.addEventListener('change', function(e) {
    if (e.target.id === 'authType' || 
        e.target.id === 'basicUsername' || 
        e.target.id === 'basicPassword' ||
        e.target.id === 'bearerToken' ||
        e.target.id === 'customAuthKey' ||
        e.target.id === 'customAuthValue') {
      updateAuthFields();
    }
  });

// --- Body Editor (CodeMirror) ---
  function setupBodyEditor() {
    const bodyEditorEl = document.getElementById('jsonBodyEditor');
    if (!bodyEditorEl) return;
    
    bodyEditorInstance = CodeMirror(bodyEditorEl, {
      mode: "application/json",
      theme: "dracula",
      lineNumbers: true,
      lineWrapping: true,
      indentUnit: 2,
      tabSize: 2,
      styleActiveLine: true,
      matchBrackets: true,
      placeholder: '{ "key": "value" }'
});
    
    bodyEditorInstance.on('change', function() {
      // Value is automatically available via bodyEditorInstance.getValue()
    });
  }

  // Helper to get body editor value
  function getBodyEditorValue() {
    if (bodyEditorInstance) {
      return bodyEditorInstance.getValue();
    }
    return "";
  }
  
  // Helper to set body editor value
  function setBodyEditorValue(value) {
    if (bodyEditorInstance) {
      bodyEditorInstance.setValue(value || "");
    }
  }

// --- Environment Variables ---
  function setupEnvModal() {
    envBtn.addEventListener("click", () => {
      renderEnvEditor();
      envModal.classList.remove("hidden");
    });
    envClose.addEventListener("click", () => envModal.classList.add("hidden"));
    envModal.addEventListener("click", (e) => { if (e.target === envModal) envModal.classList.add("hidden"); });
    envSave.addEventListener("click", saveEnv);
    envReset.addEventListener("click", resetEnv);
  }

  function resetEnv() {
    localStorage.removeItem('mvapi_env');
    envVars = JSON.parse(JSON.stringify(defaultEnvVars));
    renderEnvEditor();
    envModal.classList.add("hidden");
  }

  function renderEnvEditor() {
    envEditorEl.innerHTML = "";
    const entries = Object.entries(envVars);
    if (entries.length === 0) {
      addKVRow(envEditorEl, "", "");
    } else {
      for (const [k, v] of entries) addKVRow(envEditorEl, k, v);
    }
  }

  function saveEnv() {
    envVars = {};
    const pairs = getKVPairs("envEditor");
    for (const p of pairs) envVars[p.key] = p.value;
    localStorage.setItem("mvapi_env", JSON.stringify(envVars));
    envModal.classList.add("hidden");
  }

  function substituteEnv(str) {
    return str.replace(/\{\{(\w+)\}\}/g, (_, name) => envVars[name] !== undefined ? envVars[name] : "{{" + name + "}}");
  }

  // --- Filter ---
  function filterCollections() {
    renderCollections(searchInput.value);
  }

  // --- Utilities ---
  function escapeHTML(s) { return s.replace(/&/g,"&amp;").replace(/</g,"&lt;").replace(/>/g,"&gt;").replace(/"/g,"&quot;"); }
  function escapeAttr(s) { return String(s).replace(/&/g,"&amp;").replace(/"/g,"&quot;").replace(/</g,"&lt;").replace(/>/g,"&gt;"); }
  function truncate(s, n) { return s.length > n ? s.substring(0, n) + "..." : s; }
  function formatSize(bytes) {
    if (bytes === 0) return "0 B";
    if (bytes < 1024) return bytes + " B";
    return (bytes / 1024).toFixed(1) + " KB";
  }
})();`
}
