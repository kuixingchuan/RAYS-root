# RAYS-root

RAYS 的账号管理与运维平台，基于 Mayfly-Go 改造。

- 管理首页：账号管理、会话撤销、操作审计、运行概览与数据维护。
- 配套 RAYS 后端补丁与配置说明：[integrations/rays](integrations/rays/README.md)。
- 个人图库保持设备内存储，不包含公共数据管理。

以下保留上游项目说明与来源。

# 🌈Dromara mayfly-go

<p align="center">
    <a href="./README_EN.md">English</a> | 
    <a href="https://www.yuque.com/may-fly/mayfly-go">项目文档</a> | 
    <a href="https://space.bilibili.com/484091081/channel/collectiondetail?sid=392854">操作视频</a> | 
</p>

<p align="center">
  <a href="https://gitee.com/dromara/mayfly-go" target="_blank">
    <img src="https://gitee.com/dromara/mayfly-go/badge/star.svg?theme=white" alt="star"/>
    <img src="https://gitee.com/dromara/mayfly-go/badge/fork.svg" alt="fork"/>
  </a>
  <a href="https://github.com/dromara/mayfly-go" target="_blank">
    <img src="https://img.shields.io/github/stars/dromara/mayfly-go.svg?style=social" alt="github star"/>
    <img src="https://img.shields.io/github/forks/dromara/mayfly-go.svg?style=social" alt="github fork"/>
  </a>
  <a href="https://github.com/dromara/mayfly-go" target="_blank">
    <img src="https://gitcode.com/dromara/mayfly-go/star/badge.svg" alt="github star"/>
  </a>
  <a href="https://hub.docker.com/r/mayflygo/mayfly-go/tags" target="_blank">
    <img src="https://img.shields.io/docker/pulls/mayflygo/mayfly-go.svg?label=docker%20pulls&color=fac858" alt="docker pulls"/>
  </a>
  <a href="https://github.com/golang/go" target="_blank">
    <img src="https://img.shields.io/badge/Golang-1.24%2B-yellow.svg" alt="golang"/>
  </a>
  <a href="https://cn.vuejs.org" target="_blank">
    <img src="https://img.shields.io/badge/Vue-3.x-green.svg" alt="vue">
  </a>
</p>

## 前言

Web 版 **统一管理操作平台**，集成了对 Linux 系统的全面操作支持（包括终端管理[终端回放、命令过滤]、文件管理、脚本执行、进程监控及计划任务设置），同时提供了多种数据库（如 MySQL、PostgreSQL、Oracle、SQL Server、达梦、高斯、SQLite 等）的数据操作、数据同步与数据迁移功能。此外，还支持 Redis（单机、哨兵、集群模式）、 MongoDB 、Es 、Kafka的操作管理，并结合工单流程审批功能，为企业提供一站式的运维与管理解决方案。

## 开发语言与主要框架

- 前端：typescript、vue3、element-plus
- 后端：golang、gin、gorm

## 交流及问题反馈加 QQ 群

<a target="_blank" href="https://qm.qq.com/cgi-bin/qm/qr?k=IdJSHW0jTMhmWFHBUS9a83wxtrxDDhFj&jump_from=webapi">119699946</a>

## 演示环境

http://go.mayfly.run
账号/密码：test/test123.

## 系统核心功能截图

#### 首页

![首页](https://foruda.gitee.com/images/1757163736351080323/afb6b330_1240250.png "屏幕截图")

#### 资源管理

![资源树](https://foruda.gitee.com/images/1757163958991119284/83eb2171_1240250.png "屏幕截图")

#### 资源操作

![终端操作](https://foruda.gitee.com/images/1757164093410206293/1c7dda30_1240250.png)

![文件操作](https://foruda.gitee.com/images/1757164149388450531/0542398c_1240250.png)

![文件查看](https://foruda.gitee.com/images/1714378482611638688/7753faf6_1240250.png "屏幕截图")

![sql编辑器](https://foruda.gitee.com/images/1757164386318836686/c3b17a52_1240250.png)

![选表查数据](https://foruda.gitee.com/images/1757164281011401749/5109485f_1240250.png)

![redis操作](https://foruda.gitee.com/images/1757164442298752845/4af1b296_1240250.png)

![mongo操作](https://foruda.gitee.com/images/1714378916425714642/77fc0ed9_1240250.png "屏幕截图")

![es操作](https://foruda.gitee.com/images/1757164553845346963/b5b70381_1240250.png)

![容器操作](https://foruda.gitee.com/images/1757164625186816754/2b195e25_1240250.png)

#### 工单流程审批

![流程审批](https://foruda.gitee.com/images/1714379057627690037/ad136862_1240250.png "屏幕截图")

#### 系统管理

##### 账号管理

![账号管理](https://foruda.gitee.com/images/1714379179491881231/c6d802ae_1240250.png "屏幕截图")

##### 角色管理

![角色管理](https://foruda.gitee.com/images/1714379269408676381/6ac1e85c_1240250.png "屏幕截图")

##### 菜单资源管理

![菜单资源管理](https://foruda.gitee.com/images/1714379321338009940/a00d6a02_1240250.png "屏幕截图")

**其他更多功能&操作指南可查看上述项目文档**

## 💌 支持作者

如果觉得项目不错，或者已经在使用了，希望你可以去 <a target="_blank" href="https://github.com/dromara/mayfly-go">Github</a> 或 <a target="_blank" href="https://gitee.com/dromara/mayfly-go">Gitee</a> 或 <a target="_blank" href="https://gitcode.com/dromara/mayfly-go">Gitcode</a> 帮我点个 ⭐ Star，这将是对我极大的鼓励与支持。

> 喝杯咖啡 ☕️ 或者来杯奶茶 🧋，让作者更有精神，写出更棒的代码！

<img class="no-margin" src="https://foruda.gitee.com/images/1744113367791412282/36a3c23b_1240250.png" alt="微信打赏" width="200" height="200">

> **特别感谢：**
> 赞助金额达 199 元以上，加微信（wx-error），可受邀进入付费交流群，享受更快、更优先的技术支持与交流服务！
