package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// 日志存储
var (
	logStore []LogEntry
	logMutex sync.RWMutex
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

const version = "v0.2.0-pr7"

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
        .nav-version {
            padding: 2px 7px;
            border-radius: 3px;
            border: 1px solid rgba(255,255,255,0.08);
            font-size: 0.65rem;
            color: rgba(255,255,255,0.3);
            letter-spacing: 0.04em;
            user-select: none;
            flex-shrink: 0;
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
        .log-time { color: rgba(255,255,255,0.3); flex-shrink: 0; font-size: 0.7rem; }
        .log-level {
            padding: 2px 6px; border-radius: 3px;
            font-size: 0.65rem; font-weight: 700; flex-shrink: 0;
            background: rgba(88,166,255,0.15); color: #58a6ff;
        }
        .log-msg { color: #c9d1d9; flex: 1; }
        .log-random { color: #bc8cff; font-weight: 700; flex-shrink: 0; }
        .log-empty {
            text-align: center; color: rgba(255,255,255,0.15);
            margin-top: 40px; font-size: 0.85rem;
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
            width: 1px; height: 200px;
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
        .card-date { font-size: 1.3rem; color: #8b949e; letter-spacing: 0.05em; }
        .card-weekday { font-size: 1rem; color: rgba(255,255,255,0.25); margin-top: 0.3rem; }
        .card-time-wrap { display: flex; align-items: baseline; }
        .card-time {
            font-size: 3.6rem; font-weight: 700; color: #e6edf3;
            letter-spacing: 0.04em; font-variant-numeric: tabular-nums;
        }
        .card-millis { font-size: 1.6rem; font-weight: 400; color: #bc8cff; margin-left: 0.1rem; }

        /* ===== GitHub 统计面板 ===== */
        .stats-panel {
            position: relative;
            z-index: 1;
            margin-top: 2rem;
            width: 100%;
            max-width: 620px;
            background: rgba(22,27,34,0.65);
            border: 1px solid rgba(255,255,255,0.08);
            border-radius: 16px;
            backdrop-filter: blur(16px);
            -webkit-backdrop-filter: blur(16px);
            padding: 1.5rem 2rem;
        }
        .stats-panel h3 {
            font-size: 0.85rem;
            color: #8b949e;
            letter-spacing: 0.06em;
        }
        .stats-tabs {
            display: flex;
            gap: 0;
            margin-bottom: 1rem;
            border-bottom: 1px solid rgba(255,255,255,0.06);
        }
        .stats-tab {
            padding: 8px 18px;
            border: none;
            background: transparent;
            color: rgba(255,255,255,0.3);
            font-size: 0.78rem;
            font-family: inherit;
            cursor: pointer;
            letter-spacing: 0.04em;
            border-bottom: 2px solid transparent;
            transition: all 0.2s;
            margin-bottom: -1px;
        }
        .stats-tab:hover { color: #c9d1d9; }
        .stats-tab.active {
            color: #58a6ff;
            border-bottom-color: #58a6ff;
        }
        .stats-tab-content { display: none; }
        .stats-tab-content.active { display: block; }
        .stats-panel h3 a {
            color: #58a6ff;
            text-decoration: none;
        }
        .stats-header {
            display: flex;
            align-items: center;
            justify-content: space-between;
            margin-bottom: 1rem;
        }
        .btn-sync {
            padding: 5px 14px;
            border-radius: 6px;
            border: 1px solid rgba(88,166,255,0.35);
            background: rgba(88,166,255,0.08);
            color: #58a6ff;
            font-size: 0.7rem;
            font-family: inherit;
            cursor: pointer;
            letter-spacing: 0.04em;
            transition: all 0.2s;
            white-space: nowrap;
        }
        .btn-sync:hover {
            background: rgba(88,166,255,0.16);
            border-color: rgba(88,166,255,0.55);
        }
        .btn-sync:disabled {
            opacity: 0.4;
            cursor: not-allowed;
        }
        .stats-grid {
            display: grid;
            grid-template-columns: repeat(4, 1fr);
            gap: 1rem;
        }
        .stat-card {
            text-align: center;
            padding: 0.8rem 0.5rem;
            border-radius: 8px;
            background: rgba(255,255,255,0.02);
            border: 1px solid rgba(255,255,255,0.04);
        }
        .stat-card .stat-num {
            font-size: 1.8rem;
            font-weight: 700;
            letter-spacing: 0.02em;
        }
        .stat-card .stat-label {
            font-size: 0.7rem;
            color: rgba(255,255,255,0.3);
            margin-top: 0.3rem;
            letter-spacing: 0.04em;
        }
        .stat-open .stat-num { color: #69db7c; }
        .stat-closed .stat-num { color: #8b949e; }
        .stat-merged .stat-num { color: #bc8cff; }
        .stat-total .stat-num { color: #58a6ff; }
        .stats-sub {
            margin-top: 1.2rem;
            display: flex;
            gap: 2rem;
            font-size: 0.72rem;
            color: rgba(255,255,255,0.25);
        }
        .stats-sync {
            margin-top: 0.8rem;
            font-size: 0.68rem;
            color: rgba(255,255,255,0.15);
            letter-spacing: 0.04em;
        }
        .trend-section {
            margin-top: 1.2rem;
            padding-top: 1rem;
            border-top: 1px solid rgba(255,255,255,0.06);
        }
        .trend-section h4 {
            font-size: 0.72rem;
            color: rgba(255,255,255,0.25);
            letter-spacing: 0.06em;
            margin-bottom: 0.7rem;
        }
        #trend-canvas {
            width: 100%;
            height: 200px;
            border-radius: 6px;
            background: rgba(0,0,0,0.2);
        }
        .trend-deltas {
            display: flex;
            gap: 1.2rem;
            flex-wrap: wrap;
            font-size: 0.7rem;
        }
        .trend-delta {
            display: flex;
            align-items: center;
            gap: 4px;
            letter-spacing: 0.03em;
        }
        .trend-delta.up { color: #69db7c; }
        .trend-delta.down { color: #f85149; }
        .trend-delta.flat { color: rgba(255,255,255,0.25); }
        .trend-no-data {
            font-size: 0.7rem;
            color: rgba(255,255,255,0.15);
            text-align: center;
            padding: 0.6rem 0;
        }

        /* ===== 明细面板 ===== */
        .detail-filter {
            display: flex;
            gap: 10px;
            margin-bottom: 1rem;
            flex-wrap: wrap;
            align-items: center;
        }
        .detail-filter select,
        .detail-filter input[type="date"] {
            padding: 6px 10px;
            border-radius: 6px;
            border: 1px solid rgba(255,255,255,0.1);
            background: rgba(0,0,0,0.3);
            color: #c9d1d9;
            font-size: 0.74rem;
            font-family: inherit;
            outline: none;
        }
        .detail-filter select:focus,
        .detail-filter input[type="date"]:focus {
            border-color: #58a6ff;
        }
        .detail-filter input[type="date"]::-webkit-calendar-picker-indicator {
            filter: invert(0.7);
            cursor: pointer;
        }
        .detail-filter label {
            font-size: 0.7rem;
            color: rgba(255,255,255,0.3);
            letter-spacing: 0.04em;
        }
        .btn-filter {
            padding: 6px 16px;
            border-radius: 6px;
            border: 1px solid rgba(88,166,255,0.35);
            background: rgba(88,166,255,0.1);
            color: #58a6ff;
            font-size: 0.74rem;
            font-family: inherit;
            cursor: pointer;
            letter-spacing: 0.04em;
            transition: all 0.2s;
        }
        .btn-filter:hover {
            background: rgba(88,166,255,0.18);
            border-color: rgba(88,166,255,0.5);
        }
        .detail-table-wrap {
            overflow-x: auto;
            max-height: 360px;
            overflow-y: auto;
            border-radius: 6px;
            border: 1px solid rgba(255,255,255,0.06);
        }
        .detail-table-wrap::-webkit-scrollbar { width: 4px; height: 4px; }
        .detail-table-wrap::-webkit-scrollbar-track { background: transparent; }
        .detail-table-wrap::-webkit-scrollbar-thumb { background: rgba(255,255,255,0.08); border-radius: 2px; }
        .detail-table {
            width: 100%;
            border-collapse: collapse;
            font-size: 0.72rem;
        }
        .detail-table th {
            position: sticky;
            top: 0;
            background: rgba(22,27,34,0.95);
            padding: 10px 12px;
            text-align: left;
            color: rgba(255,255,255,0.35);
            letter-spacing: 0.04em;
            border-bottom: 1px solid rgba(255,255,255,0.08);
            white-space: nowrap;
            font-weight: 600;
        }
        .detail-table td {
            padding: 8px 12px;
            border-bottom: 1px solid rgba(255,255,255,0.03);
            color: #c9d1d9;
            white-space: nowrap;
        }
        .detail-table tr:hover td {
            background: rgba(255,255,255,0.02);
        }
        .detail-table .col-title {
            max-width: 280px;
            overflow: hidden;
            text-overflow: ellipsis;
        }
        .detail-table .col-title a {
            color: #58a6ff;
            text-decoration: none;
        }
        .detail-table .col-title a:hover { text-decoration: underline; }
        .badge {
            padding: 2px 7px;
            border-radius: 10px;
            font-size: 0.65rem;
            letter-spacing: 0.03em;
        }
        .badge-open { background: rgba(105,219,124,0.12); color: #69db7c; }
        .badge-closed { background: rgba(139,148,158,0.12); color: #8b949e; }
        .badge-merged { background: rgba(188,140,255,0.12); color: #bc8cff; }
        .badge-issue { background: rgba(88,166,255,0.1); color: #58a6ff; }
        .badge-pr { background: rgba(255,169,77,0.1); color: #ffa94d; }

        .detail-pagination {
            display: flex;
            align-items: center;
            justify-content: center;
            gap: 8px;
            margin-top: 0.8rem;
            font-size: 0.72rem;
        }
        .detail-pagination button {
            padding: 5px 12px;
            border-radius: 4px;
            border: 1px solid rgba(255,255,255,0.1);
            background: transparent;
            color: #c9d1d9;
            font-size: 0.7rem;
            font-family: inherit;
            cursor: pointer;
            transition: all 0.2s;
        }
        .detail-pagination button:hover:not(:disabled) {
            background: rgba(255,255,255,0.06);
            border-color: rgba(255,255,255,0.2);
        }
        .detail-pagination button:disabled {
            opacity: 0.3;
            cursor: not-allowed;
        }
        .detail-pagination .page-info {
            color: rgba(255,255,255,0.3);
        }
        .detail-pagination .page-info strong {
            color: #58a6ff;
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
            .stats-grid { grid-template-columns: repeat(2, 1fr); }
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
                width: 200px; height: 1px;
                background: linear-gradient(
                    to right,
                    transparent,
                    rgba(255,255,255,0.12) 20%,
                    rgba(255,255,255,0.12) 80%,
                    transparent
                );
            }
            .card-info { align-items: center; text-align: center; }
            #totem { width: 220px; height: 220px; }
            .header h1 { font-size: 2.2rem; letter-spacing: 0.2em; }
            .stats-panel { max-width: 100%; border-radius: 0; }
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
            <span class="nav-version" id="version"></span>
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

    <div class="stats-panel">
        <div class="stats-header">
            <h3><a href="https://github.com/larksuite/cli" target="_blank">larksuite/cli</a> Repository Stats</h3>
            <button class="btn-sync" id="btn-sync" onclick="manualSync()">Sync Now</button>
        </div>
        <div class="stats-tabs">
            <button class="stats-tab active" onclick="switchTab('overview')">Overview</button>
            <button class="stats-tab" onclick="switchTab('details')">Details</button>
        </div>
        <div class="stats-tab-content active" id="tab-overview">
        <div class="stats-grid">
            <div class="stat-card stat-open">
                <div class="stat-num" id="stat-issues-open">-</div>
                <div class="stat-label">Issues Open</div>
            </div>
            <div class="stat-card stat-closed">
                <div class="stat-num" id="stat-issues-closed">-</div>
                <div class="stat-label">Issues Closed</div>
            </div>
            <div class="stat-card stat-open">
                <div class="stat-num" id="stat-prs-open">-</div>
                <div class="stat-label">PRs Open</div>
            </div>
            <div class="stat-card stat-merged">
                <div class="stat-num" id="stat-prs-merged">-</div>
                <div class="stat-label">PRs Merged</div>
            </div>
        </div>
        <div class="stats-sub">
            <span>Total Items: <strong id="stat-total">-</strong></span>
            <span>PRs Closed: <strong id="stat-prs-closed">-</strong></span>
        </div>
        <div class="stats-sync" id="stat-sync">Syncing...</div>
        <div class="trend-section">
            <h4>Trend (12 snapshots)</h4>
            <canvas id="trend-canvas" width="600" height="200"></canvas>
            <div class="trend-deltas" id="trend-deltas"></div>
            <div class="trend-no-data" id="trend-no-data" style="display:none">Collecting snapshots...</div>
        </div>
        </div>
        <div class="stats-tab-content" id="tab-details">
            <div class="detail-filter">
                <select id="filter-type">
                    <option value="">All Types</option>
                    <option value="issue">Issues</option>
                    <option value="pr">Pull Requests</option>
                </select>
                <label>From</label>
                <input type="date" id="filter-start">
                <label>To</label>
                <input type="date" id="filter-end">
                <button class="btn-filter" onclick="fetchDetails(1)">Search</button>
            </div>
            <div class="detail-table-wrap">
                <table class="detail-table">
                    <thead>
                        <tr>
                            <th>#</th>
                            <th>Title</th>
                            <th>Type</th>
                            <th>State</th>
                            <th>Author</th>
                            <th>Labels</th>
                            <th>Updated</th>
                        </tr>
                    </thead>
                    <tbody id="detail-tbody">
                        <tr><td colspan="7" style="text-align:center;color:rgba(255,255,255,0.15);padding:2rem;">Select date range and click Search</td></tr>
                    </tbody>
                </table>
            </div>
            <div class="detail-pagination" id="detail-pager"></div>
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
                    list.scrollTop = list.scrollHeight;
                });
        }

        setInterval(() => {
            const panel = document.getElementById('log-panel');
            if (panel.classList.contains('open')) fetchLogs();
        }, 3000);

        // ========== 版本号 ==========
        fetch('/api/version')
            .then(r => r.json())
            .then(data => { document.getElementById('version').textContent = data.version; });

        // ========== GitHub 统计 ==========
        function fetchStats() {
            fetch('/api/stats')
                .then(r => r.json())
                .then(data => {
                    document.getElementById('stat-issues-open').textContent = data.issues_open;
                    document.getElementById('stat-issues-closed').textContent = data.issues_closed;
                    document.getElementById('stat-prs-open').textContent = data.prs_open;
                    document.getElementById('stat-prs-merged').textContent = data.prs_merged;
                    document.getElementById('stat-total').textContent = data.total_items;
                    document.getElementById('stat-prs-closed').textContent = data.prs_closed;
                    document.getElementById('stat-sync').textContent =
                        'Last synced: ' + new Date().toLocaleTimeString();
                })
                .catch(() => {
                    document.getElementById('stat-sync').textContent = 'Sync in progress...';
                });
        }
        fetchStats();
        setInterval(fetchStats, 30000);

        function manualSync() {
            const btn = document.getElementById('btn-sync');
            btn.disabled = true;
            btn.textContent = 'Syncing...';
            fetch('/api/sync', { method: 'POST' })
                .then(r => r.json())
                .then(data => {
                    if (data.error) {
                        document.getElementById('stat-sync').textContent = 'Sync failed: ' + data.error;
                    } else {
                        document.getElementById('stat-sync').textContent =
                            'Last synced: ' + new Date().toLocaleTimeString() +
                            ' (issues: ' + data.issues + ', PRs: ' + data.prs + ')';
                        fetchStats();
                        fetchTrends();
                    }
                })
                .catch(err => {
                    document.getElementById('stat-sync').textContent = 'Sync error: ' + err;
                })
                .finally(() => {
                    btn.disabled = false;
                    btn.textContent = 'Sync Now';
                });
        }

        function drawLineChart(canvas, snapshots) {
            const ctx = canvas.getContext('2d');
            const W = canvas.width, H = canvas.height;
            const pad = { top: 12, right: 16, bottom: 32, left: 48 };
            const pw = W - pad.left - pad.right;
            const ph = H - pad.top - pad.bottom;
            ctx.clearRect(0, 0, W, H);

            // Background
            ctx.fillStyle = 'rgba(13,17,23,0.5)';
            ctx.fillRect(0, 0, W, H);

            const keys = ['issues_open', 'issues_closed', 'prs_open', 'prs_merged'];
            const colors = { issues_open: '#69db7c', issues_closed: '#8b949e', prs_open: '#58a6ff', prs_merged: '#bc8cff' };
            const names = { issues_open: 'Issues Open', issues_closed: 'Issues Closed', prs_open: 'PRs Open', prs_merged: 'PRs Merged' };

            const reversed = [...snapshots].reverse();
            let maxVal = 1;
            reversed.forEach(s => keys.forEach(k => { if (s[k] > maxVal) maxVal = s[k]; }));
            maxVal = Math.ceil(maxVal * 1.15);

            // Grid
            const gridLines = 5;
            ctx.strokeStyle = 'rgba(255,255,255,0.04)';
            ctx.lineWidth = 1;
            for (let i = 0; i <= gridLines; i++) {
                const y = pad.top + (ph / gridLines) * i;
                ctx.beginPath();
                ctx.moveTo(pad.left, y);
                ctx.lineTo(W - pad.right, y);
                ctx.stroke();
                ctx.fillStyle = 'rgba(255,255,255,0.2)';
                ctx.font = '9px "Courier New", monospace';
                ctx.textAlign = 'right';
                ctx.fillText(Math.round(maxVal - (maxVal / gridLines) * i), pad.left - 6, y + 3);
            }

            // X-axis labels
            ctx.fillStyle = 'rgba(255,255,255,0.2)';
            ctx.textAlign = 'center';
            const step = Math.max(1, Math.floor(reversed.length / 6));
            reversed.forEach((s, i) => {
                if (i % step === 0 || i === reversed.length - 1) {
                    const x = pad.left + (pw / Math.max(1, reversed.length - 1)) * i;
                    const d = new Date(s.snapshot_at);
                    ctx.fillText(
                        (d.getMonth()+1) + '/' + d.getDate() + ' ' +
                        String(d.getHours()).padStart(2,'0') + ':00',
                        x, H - 8
                    );
                }
            });

            // Draw lines
            keys.forEach(key => {
                ctx.beginPath();
                ctx.strokeStyle = colors[key];
                ctx.lineWidth = 2;
                ctx.lineJoin = 'round';
                reversed.forEach((s, i) => {
                    const x = pad.left + (pw / Math.max(1, reversed.length - 1)) * i;
                    const y = pad.top + ph - (s[key] / maxVal) * ph;
                    if (i === 0) ctx.moveTo(x, y); else ctx.lineTo(x, y);
                });
                ctx.stroke();

                // Gradient fill
                const lastX = pad.left + pw;
                ctx.lineTo(lastX, pad.top + ph);
                ctx.lineTo(pad.left, pad.top + ph);
                ctx.closePath();
                const grad = ctx.createLinearGradient(0, pad.top, 0, pad.top + ph);
                grad.addColorStop(0, colors[key] + '30');
                grad.addColorStop(1, colors[key] + '00');
                ctx.fillStyle = grad;
                ctx.fill();
            });

            // Legend
            let lx = pad.left;
            ctx.font = '11px "Courier New", monospace';
            keys.forEach(key => {
                const tw = ctx.measureText(names[key]).width + 16;
                ctx.fillStyle = colors[key] + '20';
                ctx.fillRect(lx, 2, tw, 14);
                ctx.strokeStyle = colors[key];
                ctx.lineWidth = 1.5;
                ctx.strokeRect(lx, 2, tw, 14);
                ctx.fillStyle = colors[key];
                ctx.textAlign = 'center';
                ctx.fillText(names[key], lx + tw / 2, 13);
                lx += tw + 6;
            });
        }

        function fetchTrends() {
            fetch('/api/trends')
                .then(r => r.json())
                .then(data => {
                    const canvas = document.getElementById('trend-canvas');
                    const deltas = document.getElementById('trend-deltas');
                    const noData = document.getElementById('trend-no-data');
                    if (!data || data.length < 2) {
                        canvas.style.display = 'none';
                        deltas.innerHTML = '';
                        noData.style.display = 'block';
                        return;
                    }
                    canvas.style.display = 'block';
                    noData.style.display = 'none';
                    const dpr = window.devicePixelRatio || 1;
                    const rect = canvas.getBoundingClientRect();
                    canvas.width = rect.width * dpr;
                    canvas.height = 200 * dpr;
                    canvas.getContext('2d').scale(dpr, dpr);
                    canvas.style.width = rect.width + 'px';
                    canvas.style.height = '200px';
                    drawLineChart(canvas, data);

                    const latest = data[0], prev = data[1];
                    function deltaEl(label, val) {
                        let cls = 'flat', sign = '';
                        if (val > 0) { cls = 'up'; sign = '+'; }
                        else if (val < 0) { cls = 'down'; }
                        return '<span class="trend-delta ' + cls + '">' + label + ': <strong>' + sign + val + '</strong></span>';
                    }
                    deltas.innerHTML = deltaEl('Issues Open', latest.issues_open - prev.issues_open) +
                        deltaEl('Issues Closed', latest.issues_closed - prev.issues_closed) +
                        deltaEl('PRs Open', latest.prs_open - prev.prs_open) +
                        deltaEl('PRs Merged', latest.prs_merged - prev.prs_merged);
                });
        }
        fetchTrends();
        setInterval(fetchTrends, 60000);
        window.addEventListener('resize', fetchTrends);

        function switchTab(name) {
            document.querySelectorAll('.stats-tab').forEach(t => t.classList.remove('active'));
            document.querySelectorAll('.stats-tab-content').forEach(c => c.classList.remove('active'));
            document.querySelector('.stats-tab[onclick*="' + name + '"]').classList.add('active');
            document.getElementById('tab-' + name).classList.add('active');
            if (name === 'overview') { fetchTrends(); }
        }

        let detailPage = 1, detailTotal = 0, detailSize = 20;
        function fetchDetails(page) {
            detailPage = page || 1;
            const type = document.getElementById('filter-type').value;
            const start = document.getElementById('filter-start').value;
            const end = document.getElementById('filter-end').value;
            const params = new URLSearchParams({ type, start, end, page: detailPage, size: detailSize });
            fetch('/api/items?' + params)
                .then(r => r.json())
                .then(data => {
                    detailTotal = data.total;
                    const tbody = document.getElementById('detail-tbody');
                    if (!data.items || data.items.length === 0) {
                        tbody.innerHTML = '<tr><td colspan="7" style="text-align:center;color:rgba(255,255,255,0.15);padding:2rem;">No results found</td></tr>';
                        document.getElementById('detail-pager').innerHTML = '';
                        return;
                    }
                    function badge(val, cls) { return '<span class="badge ' + cls + '">' + val + '</span>'; }
                    tbody.innerHTML = data.items.map(item => {
                        const d = new Date(item.updated_at);
                        const ds = d.getFullYear() + '-' + String(d.getMonth()+1).padStart(2,'0') + '-' + String(d.getDate()).padStart(2,'0');
                        return '<tr>' +
                            '<td><a href="' + item.url + '" target="_blank" style="color:#58a6ff;text-decoration:none">#' + item.number + '</a></td>' +
                            '<td class="col-title"><a href="' + item.url + '" target="_blank">' + item.title + '</a></td>' +
                            '<td>' + badge(item.item_type === 'issue' ? 'Issue' : 'PR', 'badge-' + item.item_type) + '</td>' +
                            '<td>' + badge(item.state, 'badge-' + item.state) + '</td>' +
                            '<td>' + (item.author || '-') + '</td>' +
                            '<td style="max-width:140px;overflow:hidden;text-overflow:ellipsis">' + (item.labels || '-') + '</td>' +
                            '<td>' + ds + '</td>' +
                            '</tr>';
                    }).join('');

                    const totalPages = Math.ceil(detailTotal / detailSize);
                    const pager = document.getElementById('detail-pager');
                    pager.innerHTML =
                        '<button onclick="fetchDetails(1)" ' + (detailPage <= 1 ? 'disabled' : '') + '>First</button>' +
                        '<button onclick="fetchDetails(' + (detailPage-1) + ')" ' + (detailPage <= 1 ? 'disabled' : '') + '>Prev</button>' +
                        '<span class="page-info">Page <strong>' + detailPage + '</strong> / ' + totalPages + ' (' + detailTotal + ' items)</span>' +
                        '<button onclick="fetchDetails(' + (detailPage+1) + ')" ' + (detailPage >= totalPages ? 'disabled' : '') + '>Next</button>' +
                        '<button onclick="fetchDetails(' + totalPages + ')" ' + (detailPage >= totalPages ? 'disabled' : '') + '>Last</button>';
                });
        }
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

	// 初始化 SQLite
	if err := initDB("repo_stats.db"); err != nil {
		log.Fatalf("[FATAL] DB init failed: %v\n", err)
	}
	log.Printf("[INFO] Database initialized\n")
	appendLog("Database initialized", 0)

	// 后台同步 GitHub 数据（首次 + 每 5 小时）
	go func() {
		syncAll()
		ticker := time.NewTicker(5 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			syncAll()
		}
	}()

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

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// 中间件
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log.Printf("[INFO] %s %s %d\n", v.Method, v.URI, v.Status)
			return nil
		},
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// 路由
	e.GET("/", func(c echo.Context) error {
		return c.HTML(200, html)
	})

	e.GET("/api/logs", func(c echo.Context) error {
		return c.JSON(200, getLogs())
	})

	e.GET("/api/version", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"version": version})
	})

	e.GET("/api/stats", func(c echo.Context) error {
		s, err := getStats()
		if err != nil {
			return c.JSON(500, map[string]string{"error": err.Error()})
		}
		return c.JSON(200, s)
	})

	e.GET("/api/trends", func(c echo.Context) error {
		snapshots, err := getTrends(12)
		if err != nil {
			return c.JSON(500, map[string]string{"error": err.Error()})
		}
		if snapshots == nil {
			snapshots = []Snapshot{}
		}
		return c.JSON(200, snapshots)
	})

	e.POST("/api/sync", func(c echo.Context) error {
		issueCount, err := syncIssues()
		if err != nil {
			log.Printf("[INFO] Manual sync issues error: %v\n", err)
		}
		prCount, err := syncPRs()
		if err != nil {
			log.Printf("[INFO] Manual sync PRs error: %v\n", err)
		}
		if err := saveSnapshot(); err != nil {
			log.Printf("[INFO] Manual sync save snapshot failed: %v\n", err)
		}
		logStats()
		return c.JSON(200, map[string]int{
			"issues": issueCount,
			"prs":    prCount,
		})
	})

	e.GET("/api/items", func(c echo.Context) error {
		itemType := c.QueryParam("type")
		startDate := c.QueryParam("start")
		endDate := c.QueryParam("end")
		page, _ := strconv.Atoi(c.QueryParam("page"))
		size, _ := strconv.Atoi(c.QueryParam("size"))
		if page < 1 {
			page = 1
		}
		if size < 1 || size > 100 {
			size = 20
		}

		total, err := countItemsFiltered(itemType, startDate, endDate)
		if err != nil {
			return c.JSON(500, map[string]string{"error": err.Error()})
		}
		items, err := queryItemsFiltered(itemType, startDate, endDate, size, (page-1)*size)
		if err != nil {
			return c.JSON(500, map[string]string{"error": err.Error()})
		}
		if items == nil {
			items = []GitHubItem{}
		}
		return c.JSON(200, map[string]interface{}{
			"items": items,
			"total": total,
			"page":  page,
			"size":  size,
		})
	})

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("http://localhost%s\n", addr)
	log.Fatal(e.Start(addr))
}
