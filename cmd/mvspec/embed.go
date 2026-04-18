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
</head>
<body>
  <div class="app">
    <!-- Sidebar -->
    <aside class="sidebar">
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
      <!-- Request Bar (always top, not affected by toggle) -->
      <div class="request-bar-wrapper">
        <button id="layoutToggle" class="layout-toggle" onclick="toggleLayout()" title="Toggle Layout">
          <svg class="layout-icon icon-side-by-side" width="16" height="16" viewBox="0 0 16 16" fill="none">
            <rect x="1" y="3" width="5" height="10" rx="1" stroke="currentColor" stroke-width="1.5"/>
            <rect x="10" y="3" width="5" height="10" rx="1" stroke="currentColor" stroke-width="1.5"/>
          </svg>
          <svg class="layout-icon icon-stacked" width="16" height="16" viewBox="0 0 16 16" fill="none">
            <rect x="3" y="1" width="10" height="5" rx="1" stroke="currentColor" stroke-width="1.5"/>
            <rect x="3" y="10" width="10" height="5" rx="1" stroke="currentColor" stroke-width="1.5"/>
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
        <button id="sendBtn" class="send-btn">Send</button>
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
              <textarea id="bodyEditor" class="body-editor" placeholder='{ "key": "value" }'></textarea>
            </div>
            <div class="tab-content hidden" id="examplesTab">
              <pre id="requestExamplesOutput" class="response-output"></pre>
              <pre id="responseExamplesOutput" class="response-output" style="margin-top:16px;padding-top:16px;border-top:1px solid var(--glass-border)"></pre>
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
        <button class="send-btn" id="envSave">Save</button>
      </div>
    </div>
  </div>

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
.sidebar{width:280px;min-width:280px;background:rgba(15,23,42,.92);border-right:1px solid var(--glass-border);display:flex;flex-direction:column;backdrop-filter:blur(20px);-webkit-backdrop-filter:blur(20px)}
.sidebar-header{display:flex;align-items:center;justify-content:space-between;padding:16px 16px 12px}
.logo{font-size:18px;font-weight:700;color:var(--primary);letter-spacing:1px}
.logo span{color:var(--text)}
.env-btn{background:var(--glass);border:1px solid var(--glass-border);color:var(--text-dim);width:32px;height:32px;border-radius:var(--radius-sm);cursor:pointer;font-size:14px;transition:all .2s}
.env-btn:hover{background:var(--surface-hover);color:var(--primary);border-color:var(--primary)}
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
.main{flex:1;display:flex;flex-direction:column;overflow:hidden;padding:16px 20px;gap:12px}

/* Request Bar Wrapper */
.request-bar-wrapper{display:flex;gap:8px;align-items:center;flex-shrink:0}

/* Panels Layout (affected by toggle) */
.panels-layout{display:flex;flex:1;gap:12px;min-height:0;width:100%}
.col-request{flex:1;min-width:0;display:flex;flex-direction:column}
.col-response{min-width:300px;max-width:400px;flex-shrink:0;display:flex;flex-direction:column}
.col-response .response-panel{flex:1;display:flex;flex-direction:column;min-height:0}

/* Vertical layout (stacked) */
.panels-layout.layout-vertical{flex-direction:column}
.panels-layout.layout-vertical .col-request{flex:none;min-height:150px;max-height:40vh;overflow:hidden}
.panels-layout.layout-vertical .col-response{flex:none;min-height:150px;max-height:40vh}
.panels-layout.layout-vertical .request-panel{height:100%;overflow-y:auto}

/* Layout Toggle */
.layout-toggle{padding:6px 10px;background:var(--glass);border:1px solid var(--glass-border);border-radius:var(--radius-sm);cursor:pointer;color:var(--text-dim);font-size:16px;transition:all .2s;display:flex;align-items:center;justify-content:center}
.layout-toggle:hover{background:var(--surface-hover);border-color:var(--primary);color:var(--text)}
.layout-icon{stroke:currentColor;fill:none}
#layoutToggle[data-layout="horizontal"] .icon-stacked{display:none}
#layoutToggle[data-layout="vertical"] .icon-side-by-side{display:none}
#layoutToggle:not([data-layout]) .icon-stacked{display:none}

/* Request Bar */
.method-select{padding:8px 12px;background:var(--glass);border:1px solid var(--glass-border);color:var(--primary);font-family:var(--mono);font-weight:600;font-size:13px;border-radius:var(--radius-sm);cursor:pointer;outline:none;appearance:none;-webkit-appearance:none;min-width:100px;transition:border-color .2s}
.method-select:focus{border-color:var(--primary);box-shadow:0 0 0 2px var(--primary-glow)}
.method-select option{background:var(--bg);color:var(--text)}
.server-select{padding:8px 12px;background:var(--glass);border:1px solid var(--glass-border);color:var(--text-dim);font-family:var(--mono);font-size:13px;border-radius:var(--radius-sm);cursor:pointer;outline:none;min-width:150px}
.server-select:focus{border-color:var(--primary)}
.url-input{flex:1;padding:8px 14px;background:var(--glass);border:1px solid var(--glass-border);color:var(--text);font-family:var(--mono);font-size:13px;border-radius:var(--radius-sm);outline:none;transition:border-color .2s}
.url-input:focus{border-color:var(--primary);box-shadow:0 0 0 2px var(--primary-glow)}
.send-btn{padding:8px 24px;background:var(--primary);color:#fff;border:none;border-radius:var(--radius-sm);font-weight:600;font-size:13px;cursor:pointer;transition:all .2s;text-transform:uppercase;letter-spacing:.5px}
.send-btn:hover{box-shadow:0 0 20px var(--primary-glow);transform:translateY(-1px)}
.send-btn:active{transform:translateY(0)}
.send-btn.loading{opacity:.7;pointer-events:none}

/* Tabs */
.tabs,.response-tabs{display:flex;gap:2px;border-bottom:1px solid var(--glass-border);padding-bottom:0}
.tab{padding:8px 16px;background:none;border:none;border-bottom:2px solid transparent;color:var(--text-dim);font-size:13px;font-weight:500;cursor:pointer;transition:all .2s}
.tab:hover{color:var(--text)}
.tab.active{color:var(--primary);border-bottom-color:var(--primary)}

/* Request Panel */
.request-panel{background:var(--glass);border:1px solid var(--glass-border);border-radius:var(--radius);padding:12px;backdrop-filter:blur(12px);-webkit-backdrop-filter:blur(12px)}
.tab-content{padding-top:10px}
.tab-content.hidden{display:none}
.auth-section{margin-bottom:12px}
.auth-label{display:block;font-size:12px;font-weight:500;color:var(--text-dim);margin-bottom:4px}
.auth-type-select{width:100%;padding:8px 12px;background:var(--glass);border:1px solid var(--glass-border);color:var(--text);font-size:13px;border-radius:var(--radius-sm);cursor:pointer;outline:none}
.auth-type-select:focus{border-color:var(--primary)}
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
.body-editor{width:100%;min-height:120px;padding:10px;background:rgba(15,23,42,.6);border:1px solid var(--glass-border);border-radius:var(--radius-sm);color:var(--text);font-family:var(--mono);font-size:13px;resize:vertical;outline:none;transition:border-color .2s}
.body-editor:focus{border-color:var(--primary)}

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
.modal-footer{padding:12px 20px;border-top:1px solid var(--glass-border);display:flex;justify-content:flex-end}

/* Responsive */
@media(max-width:768px){
  .sidebar{width:220px;min-width:220px}
  .method-badge{min-width:36px;font-size:9px}
}
@media(max-width:600px){
  .app{flex-direction:column}
  .sidebar{width:100%;min-width:100%;max-height:40vh;border-right:none;border-bottom:1px solid var(--glass-border)}
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
  const responseExamplesOutput = $("#responseExamplesOutput");
  const requestExamplesOutput = $("#requestExamplesOutput");
  const historyList     = $("#historyList");
  const envBtn          = $("#envBtn");
  const envModal        = $("#envModal");
  const envClose        = $("#envClose");
  const envSave         = $("#envSave");
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

  sendBtn.addEventListener("click", sendRequest);
  urlInput.addEventListener("keydown", (e) => { if(e.key==="Enter") sendRequest(); });
  searchInput.addEventListener("input", filterCollections);

  // --- Load OpenAPI Spec ---
  function loadSpec() {
    fetch("mv-spec.json")
      .then((r) => r.json())
      .then((data) => { 
        spec = data; 
        renderCollections(); 
        renderServerDropdown();
      })
      .catch(() => {
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
    addKVRow(headersEditor, "Content-Type", "application/json");
    if (entry.op.parameters) {
      for (const p of entry.op.parameters) {
        if (p.in === "header") addKVRow(headersEditor, p.name, "");
      }
    }

    // Clear body
    bodyEditor.value = "";
    if (entry.op.requestBody && entry.op.requestBody.content) {
      const jsonContent = entry.op.requestBody.content["application/json"];
      if (jsonContent && jsonContent.schema) {
        bodyEditor.value = buildExampleBody(jsonContent.schema);
      }
    }

    // Populate response examples
    responseExamplesOutput.innerHTML = "";
    requestExamplesOutput.innerHTML = "";
    if (entry.op.responses) {
      let examplesHTML = "";
      let requestExamplesHTML = "";
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

        // Request example
        if (resp.requestExample !== undefined) {
          examplesHTML += '<div class="example-label">Request:</div>';
          examplesHTML += '<pre class="example-json">' + syntaxHighlight(JSON.stringify(resp.requestExample, null, 2)) + '</pre>';
          requestExamplesHTML += '<div class="response-example">';
          requestExamplesHTML += '<div class="response-header-row">';
          requestExamplesHTML += '<span class="response-code ' + statusClass + '">' + code + '</span>';
          requestExamplesHTML += '<span class="response-desc-text">' + description + '</span>';
          requestExamplesHTML += '<button class="try-btn" onclick="tryRequestExample(' + code + ')">Try</button>';
          requestExamplesHTML += '</div>';
          requestExamplesHTML += '<div class="example-label">Request:</div>';
          requestExamplesHTML += '<pre class="example-json">' + syntaxHighlight(JSON.stringify(resp.requestExample, null, 2)) + '</pre>';
        }

        // Response example
        if (resp.example !== undefined) {
          examplesHTML += '<div class="example-label">Response:</div>';
          examplesHTML += '<pre class="example-json">' + syntaxHighlight(JSON.stringify(resp.example, null, 2)) + '</pre>';
          if (resp.requestExample !== undefined) {
            requestExamplesHTML += '<div class="example-label">Response:</div>';
            requestExamplesHTML += '<pre class="example-json">' + syntaxHighlight(JSON.stringify(resp.example, null, 2)) + '</pre>';
          }
        }
        examplesHTML += '</div>';
        if (resp.requestExample !== undefined) {
          requestExamplesHTML += '</div>';
        }
      }
      responseExamplesOutput.innerHTML = examplesHTML;
      requestExamplesOutput.innerHTML = requestExamplesHTML || '<div style="color:var(--text-dim);padding:12px">No examples defined</div>';
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
      bodyEditor.value = JSON.stringify(resp.requestExample, null, 2);
      // Switch to Body tab
      document.querySelector('.request-panel .tabs .tab[data-tab="body"]').click();
    }
  }
  window.tryRequestExample = tryRequestExample;

  // --- Toggle Layout ---
  function toggleLayout() {
    const container = document.getElementById('panelsLayout');
    const isVertical = container.classList.toggle('layout-vertical');
    const btn = document.getElementById('layoutToggle');
    btn.setAttribute('data-layout', isVertical ? 'vertical' : 'horizontal');
    localStorage.setItem('mvapi_panels_layout', isVertical ? 'vertical' : 'horizontal');
  }
  window.toggleLayout = toggleLayout;

  // Load saved layout preference
  const savedLayout = localStorage.getItem('mvapi_panels_layout');
  const btn = document.getElementById('layoutToggle');
  if (savedLayout === 'vertical') {
    document.getElementById('panelsLayout').classList.add('layout-vertical');
    btn.setAttribute('data-layout', 'vertical');
  } else {
    btn.setAttribute('data-layout', 'horizontal');
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

    // Build curl command
    let curl = "curl -X " + method;
    for (const [key, value] of Object.entries(headers)) {
      curl += " -H \"" + key + ": " + value + "\"";
    }
    const body = substituteEnv(bodyEditor.value.trim());
    if (body && method !== "GET" && method !== "HEAD") {
      curl += " -d '" + body.replace(/'/g, "'\\''") + "'";
    }
    curl += " " + finalUrl;

    const opts = { method, headers };
    if (method !== "GET" && method !== "HEAD") {
      if (body) opts.body = body;
    }

    sendBtn.classList.add("loading");
    sendBtn.textContent = "Sending...";
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
        sendBtn.textContent = "Send";
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

  // --- Environment Variables ---
  function setupEnvModal() {
    envBtn.addEventListener("click", () => {
      renderEnvEditor();
      envModal.classList.remove("hidden");
    });
    envClose.addEventListener("click", () => envModal.classList.add("hidden"));
    envModal.addEventListener("click", (e) => { if (e.target === envModal) envModal.classList.add("hidden"); });
    envSave.addEventListener("click", saveEnv);
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
