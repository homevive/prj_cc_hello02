package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"time"
)

// 日志存储
var (
	logStore   []LogEntry
	logMutex   sync.RWMutex
)

type LogEntry struct {
	Time    string `json:"time"`
	Message string `json:"message"`
	Random  int    `json:"random"`
}

func appendLog(msg string, rnd int) {
	logMutex.Lock()
	defer logMutex.Unlock()
	logStore = append(logStore, LogEntry{
		Time:    time.Now().Format("15:04:05"),
		Message: msg,
		Random:  rnd,
	})
	if len(logStore) > 100 {
		logStore = logStore[len(logStore)-100:]
	}
}

func getLogs() []LogEntry {
	logMutex.RLock()
	defer logMutex.RUnlock()
	result := make([]LogEntry, len(logStore))
	copy(result, logStore)
	return result
}

const html = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>prj_cc_hello02 - 贾氏图腾</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            min-height: 100vh;
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            background: #0a0a12;
            color: #c9d1d9;
            font-family: 'Courier New', 'Source Code Pro', monospace;
            overflow-x: hidden;
        }
        #stars {
            position: fixed;
            top: 0; left: 0;
            width: 100%; height: 100%;
            z-index: 0;
        }

        /* ===== 顶部导航栏 ===== */
        .navbar {
            position: fixed;
            top: 0; left: 0; right: 0;
            z-index: 10;
            display: flex;
            align-items: center;
            justify-content: space-between;
            padding: 0 24px;
            height: 56px;
            background: rgba(22,27,34,0.72);
            border-bottom: 1px solid rgba(255,255,255,0.06);
            backdrop-filter: blur(16px);
            -webkit-backdrop-filter: blur(16px);
        }
        .navbar-left {
            display: flex;
            align-items: center;
            gap: 14px;
        }
        .navbar-logo {
            display: flex;
            align-items: center;
            gap: 8px;
            text-decoration: none;
            color: #e6edf3;
            font-size: 1rem;
            font-weight: 700;
            letter-spacing: 0.02em;
        }
        .navbar-logo .logo-icon {
            width: 30px; height: 30px;
            border-radius: 6px;
            background: linear-gradient(135deg, #ffa94d, #ffd43b);
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 1rem;
            color: #0d1117;
            font-family: "KaiTi", "STKaiti", "楷体", serif;
            font-weight: bold;
        }
        .navbar-search {
            width: 280px;
            height: 32px;
            padding: 0 12px;
            border-radius: 6px;
            border: 1px solid rgba(255,255,255,0.1);
            background: rgba(255,255,255,0.04);
            color: #c9d1d9;
            font-size: 0.8rem;
            font-family: inherit;
            outline: none;
            transition: border-color 0.2s, width 0.3s;
            letter-spacing: 0.02em;
        }
        .navbar-search::placeholder { color: rgba(255,255,255,0.25); }
        .navbar-search:focus {
            border-color: #58a6ff;
            width: 360px;
        }
        .navbar-right {
            display: flex;
            align-items: center;
            gap: 8px;
        }
        .nav-btn {
            padding: 6px 14px;
            border-radius: 6px;
            border: 1px solid rgba(255,255,255,0.1);
            background: transparent;
            color: #c9d1d9;
            font-size: 0.78rem;
            font-family: inherit;
            cursor: pointer;
            letter-spacing: 0.02em;
            transition: all 0.2s;
            white-space: nowrap;
            user-select: none;
        }
        .nav-btn:hover {
            background: rgba(255,255,255,0.06);
            border-color: rgba(255,255,255,0.2);
            color: #e6edf3;
        }
        .nav-btn.active {
            background: rgba(255,255,255,0.08);
            border-color: rgba(255,255,255,0.18);
            color: #fff;
        }
        .nav-btn-accent {
            border-color: rgba(188,140,255,0.35);
            color: #bc8cff;
        }
        .nav-btn-accent:hover {
            background: rgba(188,140,255,0.1);
            border-color: rgba(188,140,255,0.5);
            color: #d4bfff;
        }
        .nav-btn-log {
            border-color: rgba(105,219,124,0.35);
            color: #69db7c;
        }
        .nav-btn-log:hover {
            background: rgba(105,219,124,0.1);
            border-color: rgba(105,219,124,0.5);
            color: #8fe89e;
        }
        .nav-btn-log.on {
            background: rgba(105,219,124,0.15);
            border-color: rgba(105,219,124,0.5);
            color: #8fe89e;
        }
        .nav-btn-icon {
            width: 32px; padding: 6px 0;
            text-align: center;
            font-size: 1rem;
        }
        .nav-avatar {
            width: 28px; height: 28px;
            border-radius: 50%;
            background: linear-gradient(135deg, #58a6ff, #bc8cff);
            cursor: pointer;
            flex-shrink: 0;
        }

        /* ===== 日志面板 ===== */
        .log-panel {
            position: fixed;
            top: 56px;
            right: -420px;
            width: 400px;
            height: calc(100vh - 56px);
            z-index: 9;
            background: rgba(13,17,23,0.94);
            border-left: 1px solid rgba(255,255,255,0.08);
            backdrop-filter: blur(24px);
            -webkit-backdrop-filter: blur(24px);
            display: flex;
            flex-direction: column;
            transition: right 0.35s cubic-bezier(0.16, 1, 0.3, 1);
        }
        .log-panel.open { right: 0; }

        .log-panel-header {
            display: flex;
            align-items: center;
            justify-content: space-between;
            padding: 16px 20px;
            border-bottom: 1px solid rgba(255,255,255,0.06);
            flex-shrink: 0;
        }
        .log-panel-header h3 {
            font-size: 0.9rem;
            font-weight: 600;
            color: #69db7c;
            letter-spacing: 0.04em;
        }
        .log-panel-header .log-count {
            font-size: 0.7rem;
            color: rgba(255,255,255,0.3);
        }
        .log-panel-close {
            width: 26px; height: 26px;
            border-radius: 4px;
            border: 1px solid rgba(255,255,255,0.1);
            background: transparent;
            color: #8b949e;
            font-size: 1rem;
            cursor: pointer;
            line-height: 1;
            transition: all 0.2s;
        }
        .log-panel-close:hover {
            background: rgba(255,255,255,0.06);
            color: #e6edf3;
        }
        .log-list {
            flex: 1;
            overflow-y: auto;
            padding: 12px 20px;
            display: flex;
            flex-direction: column;
            gap: 6px;
        }
        .log-list::-webkit-scrollbar { width: 4px; }
        .log-list::-webkit-scrollbar-track { background: transparent; }
        .log-list::-webkit-scrollbar-thumb { background: rgba(255,255,255,0.08); border-radius: 2px; }

        .log-item {
            display: flex;
            align-items: center;
            gap: 10px;
            padding: 8px 12px;
            border-radius: 6px;
            background: rgba(255,255,255,0.02);
            border: 1px solid rgba(255,255,255,0.04);
            font-size: 0.75rem;
            animation: fadeIn 0.3s ease;
        }
        @keyframes fadeIn { from { opacity: 0; transform: translateY(-6px); } to { opacity: 1; transform: translateY(0); } }
        .log-time {
            color: rgba(255,255,255,0.3);
            flex-shrink: 0;
            font-size: 0.7rem;
        }
        .log-level {
            padding: 2px 6px;
            border-radius: 3px;
            font-size: 0.65rem;
            font-weight: 700;
            flex-shrink: 0;
            background: rgba(88,166,255,0.15);
            color: #58a6ff;
        }
        .log-msg {
            color: #c9d1d9;
            flex: 1;
        }
        .log-random {
            color: #bc8cff;
            font-weight: 700;
            flex-shrink: 0;
        }
        .log-empty {
            text-align: center;
            color: rgba(255,255,255,0.15);
            margin-top: 40px;
            font-size: 0.85rem;
        }

        /* ===== 页眉 ===== */
        .header {
            position: relative;
            z-index: 1;
            margin-bottom: 2rem;
        }
        .header h1 {
            font-size: 3.5rem;
            letter-spacing: 0.3em;
            font-weight: 700;
            background: linear-gradient(135deg, #58a6ff 0%, #bc8cff 50%, #f783ac 100%);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
            text-transform: uppercase;
            user-select: none;
            text-shadow: none;
            filter: drop-shadow(0 0 18px rgba(188,140,255,0.3));
        }

        /* ===== 玻璃卡片 ===== */
        .card {
            position: relative;
            z-index: 1;
            display: flex;
            align-items: center;
            gap: 3rem;
            padding: 2.5rem 3.5rem;
            background: rgba(22,27,34,0.65);
            border: 1px solid rgba(255,255,255,0.08);
            border-radius: 24px;
            backdrop-filter: blur(20px);
            -webkit-backdrop-filter: blur(20px);
            box-shadow: 0 8px 40px rgba(0,0,0,0.4), inset 0 1px 0 rgba(255,255,255,0.04);
        }

        .card-totem { flex-shrink: 0; }
        #totem { width: 280px; height: 280px; }

        .card-divider {
            width: 1px;
            height: 200px;
            background: linear-gradient(
                to bottom,
                transparent,
                rgba(255,255,255,0.12) 20%,
                rgba(255,255,255,0.12) 80%,
                transparent
            );
            flex-shrink: 0;
        }

        .card-info {
            display: flex;
            flex-direction: column;
            gap: 1.2rem;
            min-width: 260px;
        }
        .card-date {
            font-size: 1.3rem;
            color: #8b949e;
            letter-spacing: 0.05em;
        }
        .card-weekday {
            font-size: 1rem;
            color: rgba(255,255,255,0.25);
            margin-top: 0.3rem;
        }
        .card-time-wrap {
            display: flex;
            align-items: baseline;
        }
        .card-time {
            font-size: 3.6rem;
            font-weight: 700;
            color: #e6edf3;
            letter-spacing: 0.04em;
            font-variant-numeric: tabular-nums;
        }
        .card-millis {
            font-size: 1.6rem;
            font-weight: 400;
            color: #bc8cff;
            margin-left: 0.1rem;
        }

        /* ===== Footer ===== */
        .footer {
            position: relative;
            z-index: 1;
            margin-top: 2.5rem;
            font-size: 0.75rem;
            color: rgba(255,255,255,0.12);
            letter-spacing: 0.25em;
        }

        /* ===== 响应式 ===== */
        @media (max-width: 900px) {
            .navbar-search { width: 180px; }
            .navbar-search:focus { width: 220px; }
            .nav-btn { padding: 6px 10px; font-size: 0.72rem; }
            .log-panel { width: 320px; right: -340px; }
        }
        @media (max-width: 750px) {
            .navbar { padding: 0 12px; }
            .navbar-search { display: none; }
            .nav-btn { padding: 4px 8px; font-size: 0.68rem; }
            .navbar-logo span:last-child { display: none; }
            .log-panel { width: 100%; right: -105%; }
            .card {
                flex-direction: column;
                gap: 1.5rem;
                padding: 2rem;
            }
            .card-divider {
                width: 200px;
                height: 1px;
                background: linear-gradient(
                    to right,
                    transparent,
                    rgba(255,255,255,0.12) 20%,
                    rgba(255,255,255,0.12) 80%,
                    transparent
                );
            }
            .card-info {
                align-items: center;
                text-align: center;
            }
            #totem { width: 220px; height: 220px; }
            .header h1 { font-size: 2.2rem; letter-spacing: 0.2em; }
        }
    </style>
</head>
<body>
    <canvas id="stars"></canvas>

    <nav class="navbar">
        <div class="navbar-left">
            <a href="#" class="navbar-logo">
                <span class="logo-icon">贾</span>
                prj_cc_hello02
            </a>
            <input type="text" class="navbar-search" placeholder="Search or type a command...">
        </div>
        <div class="navbar-right">
            <button class="nav-btn active">Dashboard</button>
            <button class="nav-btn">Pull Requests</button>
            <button class="nav-btn">Issues</button>
            <button class="nav-btn">Explore</button>
            <button class="nav-btn nav-btn-icon" title="Notifications">&#128276;</button>
            <button class="nav-btn nav-btn-log" id="btn-log" onclick="toggleLog()">Show Log</button>
            <button class="nav-btn nav-btn-accent">+ New</button>
            <div class="nav-avatar"></div>
        </div>
    </nav>

    <div class="log-panel" id="log-panel">
        <div class="log-panel-header">
            <h3>Server Log</h3>
            <span class="log-count" id="log-count">0 entries</span>
            <button class="log-panel-close" onclick="toggleLog()">&#10005;</button>
        </div>
        <div class="log-list" id="log-list">
            <div class="log-empty">Waiting for logs...</div>
        </div>
    </div>

    <div class="header">
        <h1>Hello World</h1>
    </div>

    <div class="card">
        <div class="card-totem">
            <canvas id="totem" width="280" height="280"></canvas>
        </div>
        <div class="card-divider"></div>
        <div class="card-info">
            <div>
                <div class="card-date" id="date"></div>
                <div class="card-weekday" id="weekday"></div>
            </div>
            <div class="card-time-wrap">
                <span class="card-time" id="time"></span><span class="card-millis" id="millis"></span>
            </div>
        </div>
    </div>

    <div class="footer">Powered by Golang &nbsp;|&nbsp; 贾氏图腾</div>

    <script>
        // ========== 闪烁星星背景 ==========
        const starsCanvas = document.getElementById('stars');
        const starsCtx = starsCanvas.getContext('2d');
        let stars = [];

        function resizeStars() {
            starsCanvas.width = window.innerWidth;
            starsCanvas.height = window.innerHeight;
            initStars();
        }

        function initStars() {
            stars = [];
            const count = Math.floor((starsCanvas.width * starsCanvas.height) / 2000);
            for (let i = 0; i < count; i++) {
                stars.push({
                    x: Math.random() * starsCanvas.width,
                    y: Math.random() * starsCanvas.height,
                    r: Math.random() * 2 + 0.5,
                    baseAlpha: Math.random() * 0.6 + 0.4,
                    speed: Math.random() * 0.02 + 0.005,
                    phase: Math.random() * Math.PI * 2
                });
            }
        }

        function drawStars() {
            starsCtx.clearRect(0, 0, starsCanvas.width, starsCanvas.height);
            for (const s of stars) {
                s.phase += s.speed;
                const alpha = s.baseAlpha + Math.sin(s.phase) * 0.3;
                starsCtx.beginPath();
                starsCtx.arc(s.x, s.y, s.r, 0, Math.PI * 2);
                starsCtx.fillStyle = ` + "`rgba(255,255,255,${alpha})`" + `;
                starsCtx.fill();
            }
            requestAnimationFrame(drawStars);
        }

        window.addEventListener('resize', resizeStars);
        resizeStars();
        drawStars();

        // ========== 贾姓图腾（五彩光环） ==========
        const totemCanvas = document.getElementById('totem');
        const totemCtx = totemCanvas.getContext('2d');
        const W = 280, H = 280, CX = 140, CY = 140;
        const colors = ['#ff4d6a','#ffa94d','#ffd43b','#69db7c','#4dabf7','#b197fc','#f783ac'];
        const particles = [];

        for (let i = 0; i < 220; i++) {
            const ci = Math.floor(Math.random() * colors.length);
            particles.push({
                orbit: 62 + Math.random() * 62,
                angle: Math.random() * Math.PI * 2,
                speed: 0.004 + Math.random() * 0.012,
                r: 1.5 + Math.random() * 2.5,
                color: colors[ci],
                glow: colors[ci]
            });
        }

        function drawTotem() {
            totemCtx.clearRect(0, 0, W, H);

            for (const p of particles) {
                p.angle += p.speed;
                const x = CX + Math.cos(p.angle) * p.orbit;
                const y = CY + Math.sin(p.angle) * p.orbit;

                totemCtx.beginPath();
                totemCtx.arc(x, y, p.r + 4, 0, Math.PI * 2);
                totemCtx.fillStyle = p.glow + '26';
                totemCtx.fill();

                totemCtx.beginPath();
                totemCtx.arc(x, y, p.r, 0, Math.PI * 2);
                totemCtx.fillStyle = p.color;
                totemCtx.fill();
            }

            totemCtx.beginPath();
            totemCtx.arc(CX, CY, 50, 0, Math.PI * 2);
            totemCtx.strokeStyle = '#ffd43b';
            totemCtx.lineWidth = 2;
            totemCtx.setLineDash([5, 10]);
            totemCtx.lineDashOffset = -performance.now() / 250;
            totemCtx.stroke();
            totemCtx.setLineDash([]);

            totemCtx.beginPath();
            totemCtx.arc(CX, CY, 98, 0, Math.PI * 2);
            totemCtx.strokeStyle = 'rgba(255,255,255,0.12)';
            totemCtx.lineWidth = 1;
            totemCtx.stroke();

            totemCtx.save();
            totemCtx.translate(CX, CY);
            totemCtx.rotate(performance.now() / 8000);
            drawPolygon(totemCtx, 0, 0, 110, 8, 'rgba(255,255,255,0.05)', 1);
            totemCtx.restore();

            const grad = totemCtx.createRadialGradient(CX-14, CY-14, 10, CX, CY, 45);
            grad.addColorStop(0, '#2a1a0a');
            grad.addColorStop(0.7, '#1a0f05');
            grad.addColorStop(1, '#0d0803');
            totemCtx.beginPath();
            totemCtx.arc(CX, CY, 45, 0, Math.PI * 2);
            totemCtx.fillStyle = grad;
            totemCtx.fill();
            totemCtx.strokeStyle = '#ffa94d';
            totemCtx.lineWidth = 3;
            totemCtx.stroke();

            for (let i = 0; i < 8; i++) {
                const a = (i / 8) * Math.PI * 2 + performance.now() / 3500;
                const ix = CX + Math.cos(a) * 38;
                const iy = CY + Math.sin(a) * 38;
                totemCtx.beginPath();
                totemCtx.arc(ix, iy, 4, 0, Math.PI * 2);
                totemCtx.fillStyle = colors[i % colors.length];
                totemCtx.fill();
            }

            totemCtx.fillStyle = '#ffd43b';
            totemCtx.font = 'bold 40px "KaiTi", "STKaiti", "楷体", "SimSun", "宋体", serif';
            totemCtx.textAlign = 'center';
            totemCtx.textBaseline = 'middle';
            totemCtx.fillText('贾', CX, CY + 2);

            requestAnimationFrame(drawTotem);
        }

        function drawPolygon(ctx, x, y, r, sides, color, lineWidth) {
            ctx.beginPath();
            for (let i = 0; i < sides; i++) {
                const a = (i / sides) * Math.PI * 2 - Math.PI / 2;
                const px = x + Math.cos(a) * r;
                const py = y + Math.sin(a) * r;
                i === 0 ? ctx.moveTo(px, py) : ctx.lineTo(px, py);
            }
            ctx.closePath();
            ctx.strokeStyle = color;
            ctx.lineWidth = lineWidth;
            ctx.stroke();
        }

        drawTotem();

        // ========== 时钟 ==========
        function refresh() {
            const now = new Date();
            const y = now.getFullYear();
            const m = String(now.getMonth() + 1).padStart(2, '0');
            const d = String(now.getDate()).padStart(2, '0');
            const wd = ['日','一','二','三','四','五','六'][now.getDay()];

            document.getElementById('date').textContent =
                y + '年' + m + '月' + d + '日';
            document.getElementById('weekday').textContent = '星期' + wd;

            const hh = String(now.getHours()).padStart(2, '0');
            const mm = String(now.getMinutes()).padStart(2, '0');
            const ss = String(now.getSeconds()).padStart(2, '0');
            document.getElementById('time').textContent = hh + ':' + mm + ':' + ss;

            document.getElementById('millis').textContent = '.' + String(now.getMilliseconds()).padStart(3, '0');
        }
        setInterval(refresh, 16);
        refresh();

        // ========== 日志面板 ==========
        function toggleLog() {
            const panel = document.getElementById('log-panel');
            const btn = document.getElementById('btn-log');
            panel.classList.toggle('open');
            btn.classList.toggle('on');
            if (panel.classList.contains('open')) {
                btn.textContent = 'Hide Log';
                fetchLogs();
            } else {
                btn.textContent = 'Show Log';
            }
        }

        function fetchLogs() {
            fetch('/api/logs')
                .then(r => r.json())
                .then(data => {
                    const list = document.getElementById('log-list');
                    const count = document.getElementById('log-count');
                    count.textContent = data.length + ' entries';

                    if (data.length === 0) {
                        list.innerHTML = '<div class="log-empty">Waiting for logs...</div>';
                        return;
                    }

                    list.innerHTML = data.map(entry =>
                        '<div class="log-item">' +
                            '<span class="log-time">' + entry.time + '</span>' +
                            '<span class="log-level">INFO</span>' +
                            '<span class="log-msg">' + entry.message + '</span>' +
                            '<span class="log-random">#' + entry.random + '</span>' +
                        '</div>'
                    ).join('');

                    // 自动滚动到底部
                    list.scrollTop = list.scrollHeight;
                });
        }

        // 日志面板打开时每 3 秒刷新
        setInterval(() => {
            const panel = document.getElementById('log-panel');
            if (panel.classList.contains('open')) fetchLogs();
        }, 3000);
    </script>
</body>
</html>`

func portAvailable(port int) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	ln.Close()
	return true
}

func findPort(start int) int {
	for port := start; port < start+100; port++ {
		if portAvailable(port) {
			return port
		}
		msg := fmt.Sprintf("Port %d is in use, trying %d...", port, port+1)
		log.Printf("[INFO] %s\n", msg)
		appendLog(msg, 0)
	}
	msg := fmt.Sprintf("No available port found in range %d-%d", start, start+99)
	log.Fatalf("[FATAL] %s\n", msg)
	return 0
}

func main() {
	port := findPort(8080)

	startMsg := fmt.Sprintf("Server starting on port %d", port)
	log.Printf("[INFO] %s\n", startMsg)
	appendLog(startMsg, port)

	// 每 30 秒输出日志
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			<-ticker.C
			rnd := rand.Intn(90000) + 10000
			msg := fmt.Sprintf("Periodic check completed, random=%d", rnd)
			log.Printf("[INFO] %s\n", msg)
			appendLog(msg, rnd)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, html)
	})

	http.HandleFunc("/api/logs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(getLogs())
	})

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
