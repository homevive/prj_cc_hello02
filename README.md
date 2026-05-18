# prj_cc_hello02 — 贾氏图腾时钟

基于 Go + Echo v4 的实时时钟 Web 应用。展示贾姓图腾与五彩光环动画，集成 `larksuite/cli` 仓库 Issue/PR 数据追踪、历史快照与趋势分析。

当前版本：**v0.2.0-pr8**

---

## 需求总览

本项目从一个简单的「Hello World + 时钟」起步，逐步演进为集实时动画、日志系统、GitHub 数据追踪、SQLite 存储、趋势分析于一体的全栈 Web 应用。以下是所有需求的完整汇总。

---

## 功能清单

### 前端

| 功能 | 说明 |
|------|------|
| 闪烁星空背景 | Canvas 全屏星空动画，每颗星星独立闪烁频率和相位 |
| 贾姓图腾 | 220 颗五彩粒子环绕「贾」字旋转，带光晕轨迹、金色虚线符文环、八角星外圈 |
| 实时数字时钟 | 精确到毫秒，16ms 刷新间隔（约 60fps），等宽数字字体 |
| Hello World 标题 | 顶部渐变色彩标题，宽字距（0.3em），紫色辉光滤镜 |
| 玻璃拟态卡片 | `backdrop-filter: blur(20px)` + 半透明边框，图腾（280px）与时钟左右并排 |
| 日期与星期 | 卡片内显示年月日 + 星期几 |
| 导航栏 | GitHub 风格：项目 Logo（贾字图标）+ 搜索框 + Dashboard/PRs/Issues/Explore/+New 按钮 + 通知铃铛 + 头像 |
| 版本号 | 显示在导航栏项目名称右侧，通过 `/api/version` 获取 |
| 日志面板 | 右侧滑出面板，3 秒自动刷新，显示最近 100 条 INFO 日志（时间 + 消息 + 随机数标记） |
| 统计面板（Overview） | larksuite/cli 仓库统计卡片（Issues Open/Closed、PRs Open/Merged、Total、PRs Closed） |
| 趋势折线图（Overview） | Canvas 四色折线图（Issues Open 绿、Issues Closed 灰、PRs Open 蓝、PRs Merged 紫），横轴为月/日 时:00，带渐变填充 + 图例 + 网格线 + Y 轴刻度，每 60 秒自动刷新，支持 HiDPI |
| 增量变化指标（Overview） | 折线图下方显示最近两次快照的差值（↑ +N / ↓ -N / 0），彩色标识 |
| Tab 切换 | Overview（总览 + 趋势）/ Details（明细查询）双视图切换 |
| 明细查询面板（Details） | 类型下拉筛选（All/Issues/PRs）+ 起止日期选择器 + 分页表格（#，Title，Type，State，Author，Labels，Updated）|
| 明细分页 | 首/末/上/下翻页按钮 + 页码/总页数/总条数显示，每页 20 条 |
| Sync Now 按钮 | 手动触发 GitHub 数据同步，同步按钮禁用 + "Syncing..." 状态反馈，完成后自动刷新统计和趋势 |
| 响应式布局 | 窄屏（≤750px）自动切换为上下堆叠，导航栏按钮缩略，搜索框隐藏 |

### 后端

| 功能 | 说明 |
|------|------|
| Echo v4 框架 | Web 层使用 Echo v4 框架，含 RequestLogger + Recover + CORS 中间件 |
| 端口自动探测 | 从 8080 递增检查端口可用性，找到第一个空闲端口后启动，过程输出 INFO 日志 |
| 30 秒定时日志 | goroutine 每 30 秒输出 INFO 日志，附带 10000-99999 五位随机数 |
| 线程安全日志存储 | `sync.RWMutex` 保护的内存日志队列，最多保留 100 条 |
| GitHub 数据抓取 | REST API v3 抓取 `larksuite/cli` Issues（100 条/页，过滤 PR）+ Pull Requests（100 条/页）|
| PR 合并检测 | closed 状态的 PR 二次调用 `/pulls/:number/merge` 确认是否 merged（HTTP 204 = 已合并）|
| SQLite 存储 | `github_items` 表存储 Issue/PR 明细，`snapshots` 表存储历史统计快照 |
| 唯一约束去重 | `(number, item_type)` 联合唯一约束，INSERT ON CONFLICT 自动更新 state/title/labels/updated_at |
| 5 小时定时同步 | 首次启动立即同步，之后每 5 小时自动同步一次 |
| 手动同步 | `POST /api/sync` 端点，页面按钮触发，返回 issues/prs 数量 |
| 历史快照 | 每次同步（定时或手动）自动调用 `saveSnapshot()` 写入 snapshots 表 |
| 趋势查询 | `getTrends(12)` 返回最近 12 条快照（倒序），供前端折线图使用 |
| 明细查询 | `GET /api/items` 支持 type/start/end/page/size 参数，动态构建 WHERE 子句，返回分页结果 + 总数 |

---

## 快速启动

### 环境要求

- Go 1.21+
- GitHub Personal Access Token（需 `repo` 权限，读取公开仓库）

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

### Go 代理配置（中国大陆）

```bash
export GOPROXY=https://goproxy.cn,direct
```

---

## API 文档

### `GET /`
返回主页面 HTML。

### `GET /api/logs`
获取最近 100 条服务器日志。

```json
[
  { "time": "22:58:23", "message": "Synced 42 issues", "random": 42 }
]
```

### `GET /api/version`
获取当前版本号。

```json
{ "version": "v0.2.0-pr8" }
```

### `GET /api/stats`
获取 `larksuite/cli` 实时统计。

```json
{
  "issues_open": 25, "issues_closed": 6,
  "prs_open": 42, "prs_closed": 58, "prs_merged": 49,
  "total_items": 131
}
```

### `GET /api/trends`
获取最近 12 条历史快照（按时间倒序）。首次同步仅 1 条，需至少 2 条才有趋势显示。

```json
[
  {
    "id": 2,
    "snapshot_at": "2026-05-18T08:46:41+08:00",
    "issues_open": 25, "issues_closed": 6,
    "prs_open": 42, "prs_closed": 58, "prs_merged": 49,
    "total_items": 131
  }
]
```

### `POST /api/sync`
手动触发一次 GitHub 数据同步（抓取 → 写入 → 快照 → 统计）。

```json
{ "issues": 42, "prs": 58 }
```

### `GET /api/items`
分页查询 Issue/PR 明细，支持按类型和日期范围筛选。

**参数：**

| 参数 | 说明 | 示例 |
|------|------|------|
| type | 类型 | `issue` / `pr` / 空 = 全部 |
| start | 起始日期 | `2026-05-01` |
| end | 截止日期 | `2026-05-18` |
| page | 页码（默认 1） | `1` |
| size | 每页条数（1-100，默认 20） | `20` |

**响应：**

```json
{
  "items": [
    {
      "id": 4463229004,
      "number": 929,
      "title": "form-questions-create: attachment 默认只接收图片",
      "state": "open",
      "item_type": "issue",
      "author": "pipi-HST",
      "labels": "bug,domain/base",
      "url": "https://github.com/larksuite/cli/issues/929",
      "created_at": "2026-05-17T10:54:38Z",
      "updated_at": "2026-05-17T13:10:47Z"
    }
  ],
  "total": 31,
  "page": 1,
  "size": 20
}
```

---

## 数据库

使用 SQLite (`repo_stats.db`)，纯 Go 实现无 CGO 依赖。

### 表 `github_items` — Issue/PR 明细

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER PK | GitHub 全局 ID |
| number | INTEGER | Issue/PR 编号 |
| title | TEXT | 标题 |
| state | TEXT | `open` / `closed` / `merged` |
| item_type | TEXT | `issue` / `pr` |
| author | TEXT | GitHub 用户名 |
| labels | TEXT | 逗号分隔标签 |
| url | TEXT | HTML 链接 |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 最后更新时间 |

- 唯一约束：`(number, item_type)`
- 冲突策略：`ON CONFLICT DO UPDATE` (state, title, labels, updated_at)
- 索引：`idx_item_type`, `idx_state`

### 表 `snapshots` — 历史统计快照

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER PK AUTOINCREMENT | 自增序号 |
| snapshot_at | DATETIME | 快照时间（本地时间） |
| issues_open | INTEGER | 开放 Issue 数 |
| issues_closed | INTEGER | 已关闭 Issue 数 |
| prs_open | INTEGER | 开放 PR 数 |
| prs_closed | INTEGER | 已关闭 PR 数（含 merged） |
| prs_merged | INTEGER | 已合并 PR 数 |
| total_items | INTEGER | 总条目数 |

每次同步（定时或手动）自动插入一条记录。

---

## 项目结构

```
prj_cc_hello02/
├── main.go                 # Echo 服务入口 + 内联 HTML/CSS/JS 前端
├── db.go                   # SQLite 数据库层（建表、CRUD、快照、趋势、筛选查询）
├── github.go               # GitHub REST API 客户端（Issue/PR 抓取 + 合并检测）
├── go.mod / go.sum         # Go 模块依赖
├── repo_stats.db           # SQLite 数据库文件（自动生成）
├── prj_cc_hello02.exe      # 编译产物
└── README.md
```

### 架构流程

```
main.go
  │
  ├─ findPort(8080) ─── 端口探测
  ├─ initDB() ───────── SQLite 初始化
  ├─ go syncAll() ──── 首次同步 + 每 5h 定时同步
  │     └─ syncIssues()/syncPRs() → upsertItem() → saveSnapshot() → logStats()
  ├─ go 30s ticker ─── 随机数日志
  │
  └─ Echo HTTP Server (:8080+)
        ├─ GET /             → HTML
        ├─ GET /api/logs     → getLogs()
        ├─ GET /api/version  → version
        ├─ GET /api/stats    → getStats()
        ├─ GET /api/trends   → getTrends(12)
        ├─ POST /api/sync    → syncIssues/syncPRs/saveSnapshot
        └─ GET /api/items    → countItemsFiltered/queryItemsFiltered
```

---

## 技术栈

| 组件 | 技术 | 版本 |
|------|------|------|
| 语言 | Go | 1.21+ |
| Web 框架 | Echo | v4.15.2 |
| 数据库驱动 | modernc.org/sqlite | v1.50.1 |
| 前端 | HTML5 Canvas + CSS3 + Vanilla JS | — |
| 外部 API | GitHub REST API v3 | — |
| 认证 | Personal Access Token (Bearer) | — |

---

## 版本历史

| 版本 | PR | 说明 |
|------|:--:|------|
| `v0.2.0-pr8` | #4 | 完善 README：汇总所有需求 |
| `v0.2.0-pr7` | #4 | 明细查询面板（Tab + 日期筛选 + 分页表）；趋势折线图（Canvas 四色线 + 图例） |
| `v0.2.0-pr6` | #4 | 完善 README 文档（API 示例、数据库 Schema、架构图） |
| `v0.2.0-pr5` | #4 | Sync Now 按钮 + POST /api/sync 端点 + README 更新 |
| `v0.2.0-pr4` | #4 | 5 小时定时同步 + snapshots 表 + /api/trends + Sparkline + 增量指标 |
| `v0.2.0-pr3` | #3 | Echo v4 重构 + larksuite/cli 统计面板 |
| — | #2 | 版本号显示（PR 编号驱动） |
| — | #1 | 初始版本：Hello World + 贾氏图腾时钟 + 星空背景 + 导航栏 + 日志面板 |

---

## 开发

```bash
# 安装依赖
go mod download

# 运行（开发模式）
go run main.go

# 构建
go build -o prj_cc_hello02.exe .

# 清理编译产物和数据库
rm -f prj_cc_hello02.exe repo_stats.db
```

---

## License

MIT
