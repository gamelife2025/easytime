# CI/CD 工作流说明

## 概述

本项目配置了 GitHub Actions 工作流，支持自动化编译、测试和发布。

## 工作流说明

### 1. Test 工作流 (`test.yml`)

在以下情况触发：
- 推送到 `main` 或 `develop` 分支
- 创建或更新 Pull Request 到 `main` 或 `develop` 分支

执行的任务：
- **Test**: 运行单元测试并生成代码覆盖率报告
- **Build Check**: 检查各平台的编译情况 (Linux/Windows/macOS)
- **Lint**: 运行代码检查工具

### 2. Build and Release 工作流 (`build.yml`)

#### 触发方式：

**方式一：标签发布**（推荐）
```bash
git tag v1.0.0
git push origin v1.0.0
```

**方式二：手动触发**
在 GitHub 仓库页面 → Actions → Build and Release → Run workflow

#### 编译平台：

| 操作系统 | 架构 | 可用 |
|---------|------|------|
| Linux | x86_64, arm64, 386, ARM | ✓ |
| Windows | x86_64, 386, ARM64 | ✓ |
| macOS | x86_64, ARM64 | ✓ |

#### 发布流程：

1. 编译所有平台的可执行文件
2. 上传为 artifacts
3. 创建 GitHub Release
4. 自动上传二进制文件到 Release 页面

## 本地编译

### 编译当前平台
```bash
make build
```

### 编译所有平台
```bash
make build-all
```

编译结果将保存到 `build/` 目录

### 运行测试
```bash
make test
```

### 清理编译文件
```bash
make clean
```

## 使用发布的二进制文件

标签发布后，在 GitHub Releases 页面可下载预编译的可执行文件。

文件命名规则：`easytime-{OS}-{ARCH}[.exe]`

例如：
- `easytime-linux-amd64` - Linux x86_64 版本
- `easytime-windows-amd64.exe` - Windows x86_64 版本
- `easytime-darwin-arm64` - macOS ARM64 版本

## 代码检查规则

项目配置了 `.golangci.yml`，包含以下检查：
- 代码格式（gofmt）
- 导入优化（goimports）
- 错误处理（errcheck）
- 类型检查（typecheck）
- 代码简化（gosimple）
- 等等

## 环境版本

- **Go**: 1.20+
- **GitHub Actions**: 最新

## 故障排除

### 编译失败
检查 Actions 日志中的错误信息，确保代码符合 Go 语言规范。

### 发布失败
确保：
1. 标签格式为 `v*` (如 `v1.0.0`)
2. 所有编译任务已完成
3. GitHub Token 权限正确

### 测试失败
运行本地测试：
```bash
go test -v ./...
```

## 更多参考

- [GitHub Actions 文档](https://docs.github.com/en/actions)
- [Go 交叉编译指南](https://golang.org/doc/install/source#environment)
