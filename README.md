# prj_cc_hello02 — 贾氏图腾时钟

基于 Go + Echo v4 的实时时钟 Web 应用，展示贾姓图腾与五彩光环动画，集成 GitHub 仓库 Issue/PR 数据追踪与趋势分析。

当前版本：**v0.2.0-pr6**

---

## 功能

### 前端

- **星空背景** — Canvas 动画实现全屏闪烁星空
- **贾姓图腾** — 五彩粒子环绕「贾」字旋转，带光晕轨迹、符文环、八角星
- **实时时钟** — 精确到毫秒，16ms 刷新间隔，数字等宽对齐
- **玻璃拟态卡片** — `backdrop-filter: blur()` + 半透明边框，图腾与时钟左右并排
- **GitHub 风格导航栏** — Dashboard / Pull Requests / Issues / Explore / +New 按钮
- **日志面板** — 右侧滑出面板，实时显示服务器 INFO 日志（含随机数标记）
- **响应式布局** — 窄屏自动切换为上下堆叠，移动端友好

### 后端

- **定时日志** — 每 30 秒输出 INFO 日志，附带 5 位随机数
- **端口自动探测** — 从 8080 递增查找可用端口，日志中输出过程
- **GitHub 数据同步** — 抓取 `larksuite/cli` 的 Issue 和 PR，存入 SQLite
- **定时同步** — 每 5 小时自动执行一次同步（可配置）
- **手动同步** — 页面按钮点击即刻触发同步，实时反馈结果
- **历史快照** — 每次同步自动记录统计快照，支持趋势分析
- **趋势面板** — Sparkline 柱状图 + 增量变化指标（↑↓ 数字）
- **PR 合并检测** — closed 状态的 PR 会二次查询确认是否已合并

---

## 快速启动

### 环境要求

- Go 1.21+
- GitHub Personal Access Token（需 `repo` 权限，用于读取公开仓库）

### 启动步骤

```bash
# 1. 克隆仓库
git clone https://github.com/homevive/prj_cc_hello02.git
cd prj_cc_hello02

# 2. 设置 GitHub Token
export GITHUB_TOKEN="ghp_xxxxxxxxxxxx"

# 3. 运行
go run main.go

# 或编译后运行
go build -o prj_cc_hello02.exe .
./prj_cc_hello02.exe
```

浏览器打开 `http://localhost:8080`（端口被占用则自动递增）。

### 环境变量

| 变量 | 必需 | 说明 |
|------|:--:|------|
| `GITHUB_TOKEN` | 是 | GitHub Personal Access Token，用于 API 认证 |

---

## API 文档

### `GET /`

返回主页面 HTML。

### `GET /api/logs`

获取最近 100 条服务器日志。

**响应示例：**

```json
[
  {
    "time": "22:58:23",
    "message": "Synced 42 issues",
    "random": 42
  }
]
```

### `GET /api/version`

获取当前版本号。

**响应示例：**

```json
{ "version": "v0.2.0-pr5" }
```

### `GET /api/stats`

获取 `larksuite/cli` 仓库实时统计数据。

**响应示例：**

```json
{
  "issues_open": 25,
  "issues_closed": 6,
  "prs_open": 42,
  "prs_closed": 58,
  "prs_merged": 49,
  "total_items": 131
}
```

### `GET /api/trends`

获取最近 12 条历史快照（按时间倒序）。首次同步仅 1 条，趋势面板需至少 2 条才有增量显示。

**响应示例：**

```json
[
  {
    "id": 2,
    "snapshot_at": "2026-05-18T08:46:41+08:00",
    "issues_open": 25,
    "issues_closed": 6,
    "prs_open": 42,
    "prs_closed": 58,
    "prs_merged": 49,
    "total_items": 131
  }
]
```

### `POST /api/sync`

手动触发一次 GitHub 数据同步（实时抓取、写入数据库、保存快照）。

**响应示例：**

```json
{ "issues": 42, "prs": 58 }
```

---

## 数据库

使用 SQLite 文件数据库 `repo_stats.db`，纯 Go 实现无 CGO 依赖。

### 表 `github_items`

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER PK | GitHub 全局 ID |
| number | INTEGER | Issue/PR 编号 |
| title | TEXT | 标题 |
| state | TEXT | open / closed / merged |
| item_type | TEXT | issue / pr |
| author | TEXT | 提交者 GitHub 用户名 |
| labels | TEXT | 逗号分隔的标签名 |
| url | TEXT | HTML 链接 |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 最后更新时间 |

唯一约束：`(number, item_type)` — 插入冲突时自动更新 state、title、labels、updated_at。

### 表 `snapshots`

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER PK AUTOINCREMENT | 快照序号 |
| snapshot_at | DATETIME | 快照时间 |
| issues_open / issues_closed | INTEGER | Issue 各状态计数 |
| prs_open / prs_closed / prs_merged | INTEGER | PR 各状态计数 |
| total_items | INTEGER | 总条目数 |

每次同步（定时或手动）自动插入一条快照记录，前端趋势面板据此展示 Sparkline 和增量变化。

---

## 项目结构

```
prj_cc_hello02/
├── main.go                 # Echo 服务入口 + 内联 HTML/CSS/JS 前端
├── db.go                   # SQLite 数据库层（建表、CRUD、快照、趋势）
├── github.go               # GitHub REST API 客户端（Issue/PR 抓取）
├── go.mod / go.sum         # Go 模块依赖
├── repo_stats.db           # SQLite 数据库（自动生成）
├── prj_cc_hello02.exe      # 编译产物
└── README.md
```

### 架构说明

```
main.go  ─启动─→ Echo HTTP Server (:8080+)
  │                   │
  ├─ go routine ──→ 每 30s 日志 + 随机数
  ├─ go routine ──→ 首启同步 → 每 5h 自动同步
  │                   │
  └─ GET /api/* ──→ db.go (SQLite) ─→ JSON
       POST /api/sync ─→ github.go ─→ db.go ─→ JSON
```

---

## 技术栈

| 组件 | 技术 |
|------|------|
| 语言 | Go 1.21+ |
| Web 框架 | [Echo v4](https://github.com/labstack/echo) |
| 数据库 | SQLite via [modernc.org/sqlite](https://modernc.org/sqlite) |
| 前端 | HTML5 Canvas + CSS3 + Vanilla JS |
| API | GitHub REST API v3 |

---

## 开发

```bash
# 安装依赖
go mod download

# 运行（开发模式）
go run main.go

# 构建
go build -o prj_cc_hello02.exe .

# 清理
rm -f prj_cc_hello02.exe repo_stats.db
```

Go 代理配置（中国大陆）：

```bash
export GOPROXY=https://goproxy.cn,direct
```

---

## License

MIT
