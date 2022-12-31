<p align="center">
  <img width="150" src="doc/img/logo.png">
</p>

<p align="center">
  <a href="https://github.com/apioak/apioak-admin">
    <img src="https://img.shields.io/badge/apioak--admin-v0.6.0-blue" alt="apioak-admin">
  </a>
  <a href="https://github.com/golang/go">
    <img src="https://img.shields.io/badge/Go-v1.16-blue" alt="go-1.16">
  </a>
  <a href="https://github.com/gin-gonic/gin">
    <img src="https://img.shields.io/badge/Gin-v1.7.2-blue" alt="gin-1.7.2">
  </a>
</p>

[简体中文](README_CN.md) | [English](README.md)

## 简介
`apioak-admin` 是`apioak`网关的控制面后端项目，基于 <a target="_blank" href="https://github.com/golang/go">Go 1.16</a> 和 <a target="_blank" href="https://github.com/gin-gonic/gin">Gin 1.7.2</a> 开发，项目配合数据面的项目 <a target="_blank" href="https://github.com/apioak/apioak">apioak</a> 一起使用。
该项目旨在简化`apioak`的上手使用，优化用户的操作，已达到极简的操作即可完成一个完整服务配置的上线与发布。

## 编译
```
go build -o apioak-admin main.go
```

## 依赖
在不同的操作系统上安装 `apioak-admin` 所必需的系统依赖（`MySQL >= 5.7 或 MariaDB >= 10.2`等），请参见：[依赖安装文档](doc/zh_CN/install-dependencies.md)。

## 配置
- 导入数据库配置文件到 `MySQL` 或 `MariaDB` 中，数据表配置文件路径 `/path/config/apioak.sql`。

- 在上面编译命令后生成的 `apioak` 可执行文件的所在目录创建 `config` 目录，同时将 `apioak-admin` 项目下的配置文件 `app_example.yaml` 复制到该目录下，并更改名称为 `app.yaml` ，然后在该配置文件中配置 `database` 项的数据库连接信息，`server` 项启动服务后访问服务的信息， `apioak` 项为数据面的数据面部署的`admin-api`访问信息。

## 运行
直接执行可执行文件即可完成项目启动
```
./apioak
```

## 其他
项目中为了方便使用，有可执行文件已经包含了本项目的前后端项目，只需要将对应的可执行文件下载下来，然后增加可执行文件目录下的 `config/app.yaml` 配置文件可以完成前后端项目的启动，直接访问 `config/app.yaml` 配置文件中的 `server` 项即可进行操作使用。









