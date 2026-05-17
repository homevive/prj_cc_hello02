# prj_cc_hello02 - 贾氏图腾时钟

基于 Golang 的实时时钟 Web 应用，展示贾姓图腾与五彩光环动画。

## 功能

- GitHub 风格顶部导航栏
- 闪烁星空背景（Canvas）
- 贾姓图腾，五彩粒子环绕旋转
- 实时数字时钟（精确到毫秒）
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

## 技术栈

- Go 标准库 `net/http`
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
