<p align="center">
  <img width="150" src="doc/img/logo.png">
</p>

<p align="center">
  <a href="https://github.com/apioak/apioak-admin">
    <img src="https://img.shields.io/badge/apioak--admin-v0.6.0-blue" alt="apioak-admin">
  </a>
  <a href="https://github.com/golang/go">
    <img src="https://img.shields.io/badge/GO-v1.16-blue" alt="go-1.16">
  </a>
  <a href="https://github.com/gin-gonic/gin">
    <img src="https://img.shields.io/badge/gin-v1.7.2-blue" alt="gin-1.7.2">
  </a>
</p>

[简体中文](README_CN.md) | [English](README.md)

## Introduction
`apioak-admin` is the control plane backend project of `apioak` gateway, based on <a target="_blank" href="https://github.com/golang/go">Go 1.16</a> and <a target="_blank" href="https://github.com/gin-gonic/gin">Gin 1.7.2</a> development, the project matches the data surface project <a target="_blank" href="https ://github.com/apioak/apioak">apioak</a>.
The project aims to simplify the use of `apioak`, optimize the user's operation, and achieve a minimal operation to complete the launch and release of a complete service configuration.

## Quick start
For the convenience of use, the front-end and back-end projects are merged and packaged as out-of-the-box executable files, which only need to be downloaded in [Releases](https://github.com/apioak/apioak-admin/releases) Compress the package and decompress it, then configure the `config/app.yaml` configuration file in the corresponding directory and execute the executable file to complete the deployment of the project. Just access the contents of the `server` configuration item in the `config/app.yaml` configuration file.

## Self-compiled
```
go build -o apioak-admin main.go
```

## Rely
For the system dependencies necessary to install `apioak-admin` on different operating systems (`MySQL >= 5.7 or MariaDB >= 10.2`, etc.), please refer to: [Dependency Installation Documentation](doc/zh_CN/install-dependencies.md ).

## Configuration
- Import the database configuration file to `MySQL` or `MariaDB`, the data table configuration file path `/{path}/config/apioak.sql`.

- Create a `config` directory in the directory where the `apioak-admin` executable file generated after compiling the command is located, and copy the configuration file `app_example.yaml` under the `apioak-admin` project to this directory, and change the name to `app.yaml`, and then configure in that configuration file.
  > - `database`: database connection information.
  > - `token`: User login to issue `token` configuration information.
  > - `server`: Information about accessing the service after starting the service.
  > - `apioak`: Data plane configuration synchronization connection information.
  > - `logger`: Record log configuration information.
  > - `validator`: The language of parameter verification information. zh:Chinese (default) / en:English

## Run
```
./apioak-admin
```
