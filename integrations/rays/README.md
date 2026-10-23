# RAYS 后端配套接入

管理平台需要配套 RAYS Spring Boot 管理接口。backend.patch 包含本轮后端实现、迁移和测试，不包含个人数据、运行凭据及其他清理改动。

在 RAYS 业务仓库根目录先运行 `git apply --check /path/to/backend.patch`，确认后再应用。补丁基于本轮修改前的 RAYS 版本；已完成这些改动的工作区不要重复应用。

配置和验证范围见随附说明。
