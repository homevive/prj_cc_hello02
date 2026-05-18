# prj_cc_hello02 - 贾氏图腾时钟

基于 Golang + Echo 框架的实时时钟 Web 应用，展示贾姓图腾与五彩光环动画，集成 GitHub 仓库数据追踪。

## 功能

- GitHub 风格顶部导航栏
- 闪烁星空背景（Canvas）
- 贾姓图腾，五彩粒子环绕旋转
- 实时数字时钟（精确到毫秒）
- 玻璃拟态卡片布局，响应式适配
- 实时日志面板（每 30 秒自动输出）
- **larksuite/cli 仓库 Issue/PR 追踪**
- **定时同步（每 5 小时）+ 手动同步按钮**
- **历史快照趋势图（Sparkline）+ 增量变化指标**
- 端口自动探测（从 8080 递增）

## 快速启动

```bash
# 设置 GitHub Token（用于读取 larksuite/cli 数据）
export GITHUB_TOKEN="your_github_token"

# 源码运行
go run main.go

# 或编译后运行
go build -o prj_cc_hello02.exe .
./prj_cc_hello02.exe
```

浏览器打开 `http://localhost:8080`。

## API

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/` | 主页面 |
| GET | `/api/logs` | 获取最近 100 条服务器日志 |
| GET | `/api/version` | 获取当前版本号 |
| GET | `/api/stats` | 获取 larksuite/cli 实时统计 |
| GET | `/api/trends` | 获取最近 12 条历史快照 |
| POST | `/api/sync` | 手动触发 GitHub 数据同步 |

## 数据库

使用 SQLite (`repo_stats.db`)，包含两张表：

- **github_items** — 存储 Issue 和 PR 数据（自动去重更新）
- **snapshots** — 每次同步时记录统计快照，用于趋势分析

## 技术栈

- Go + Echo v4 Web 框架
- HTML5 Canvas 动画
- CSS 玻璃拟态 + 响应式布局
- SQLite（modernc.org/sqlite，纯 Go 实现）
- GitHub REST API

## 项目结构

```
prj_cc_hello02/
├── go.mod / go.sum        # 依赖管理
├── main.go                 # 服务入口 + 前端 HTML/CSS/JS
├── db.go                   # SQLite 数据库层
├── github.go               # GitHub API 客户端
├── prj_cc_hello02.exe      # 编译产物
└── README.md
```
