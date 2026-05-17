# prj_cc_hello02 - 贾氏图腾时钟

基于 Golang + Marten 框架的实时时钟 Web 应用，展示贾姓图腾与五彩光环动画。

## 功能

- GitHub 风格顶部导航栏（搜索框 + 导航按钮）
- 贾姓图腾，五彩粒子环绕旋转（Canvas 动画）
- 实时数字时钟（精确到毫秒）
- 闪烁星空背景
- 右侧滑出日志面板（服务端每 30s 生成随机数日志）
- 端口自动探测（8080 起，被占用则 +1）
- 版本号展示（基于 GitHub PR 编号）
- 玻璃拟态卡片布局，响应式适配

## 快速启动

```bash
# 源码运行
go run main.go

# 或编译后运行
go build -o prj_cc_hello02.exe .
./prj_cc_hello02.exe
```

浏览器打开 `http://localhost:8080`。

## API

| 端点 | 说明 |
|------|------|
| `GET /` | 主页面（HTML） |
| `GET /api/logs` | 服务端日志列表（JSON） |
| `GET /api/version` | 版本信息（JSON） |

## 技术栈

- [Marten](https://github.com/gomarten/marten) — 零依赖轻量 Web 框架
- HTML5 Canvas 动画
- CSS 玻璃拟态 + 响应式布局

## 项目结构

```
prj_cc_hello02/
├── go.mod
├── main.go              # 服务入口，内联 HTML/CSS/JS
├── prj_cc_hello02.exe   # 编译产物
└── README.md
```
