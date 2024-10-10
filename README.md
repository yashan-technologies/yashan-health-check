# YHC(Yashan Health Check) | 崖山深度健康检查工具

YHC(崖山深度健康检查工具)，是一款针对崖山数据库的深度检查工具。旨在提供更为专业的数据库检查报告及建议，为客户数据保障提供全面灵活的解决方案。

## 用户定位

- 崖山DBA

## 产品定位

- 轻量的独立工具
- 开箱即用

## 场景建议

- 性能监控和优化
- 故障排查和问题定位
- 安全审计和合规性检查
- 数据库升级和迁移前检查
- 数据库日常巡检
- 其他任何想要快速检查相关信息时
- 部署后检查校验

## 核心功能

### 服务器信息检查

- 服务器的基本信息(操作系统、硬件配置、防火墙等)
- 服务器的负载情况(网络流量、CPU占用、I/O负载、内存、磁盘容量等)
- 服务器的系统日志
- ...

### 数据库信息检查

- 数据库的基本信息（版本、实例信息、主备信息等）
- 数据库的文件检查（数据文件、控制文件、备份文件等）
- 数据库的对象检查（表空间、表、索引、约束、序列等）
- 数据库的负载检查（会话、事务、等待事件、锁等待、缓存池命中率等）
- 数据库的安全检查（密码强度、用户权限、登录配置、默认表空间等）
- 数据库日志分析（alert.log、run.log、redo日志等）
- ...

### 检查结果告警与建议

- 支持自定义指标告警表达式
- 灵活可配的阈值、告警级别以及告警建议

### 丰富可扩展的检查指标

- 提供了丰富的默认指标
- 支持自定义检查指标（bash、sql）

### 灵活可配的检查策略

- 支持自定义时间周期，绝对灵活的时间周期选择
- 支持自定义路径数据检查，目录或文件均可批量选择
- ...

### 丰富健全的数据管理

- 支持自定义检查数据的存放路径
- 多种检查数据展示形式(html、docx)
- ...
## 使用方法

### 工具帮助信息

```
bash # ./yhcctl -h
Usage: yhcctl <command>

Yhcctl is used to manage the yashan health check.

Flags:
  -h, --help                          Show context-sensitive help.
  -v, --version                       Show version.
  -c, --config="./config/yhc.toml"    Configuration file.

Commands:
  check            The check command is used to yashan health check.

  after-install    The after-install command is used to verify the installation of Yashandb after it has been installed.

Run "yhcctl <command> --help" for more information on a command.
```

### 最佳实践

```shell
# 标准健康检查
./yhcctl check

# 部署后检查校验
./yhcctl after-install
```

>更多使用方法详见产品文档 (工具包路径/docs/yhc.pdf)