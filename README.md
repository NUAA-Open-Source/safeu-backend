# SafeU Backend

伏羲后端仓库。

> README for developers.

<!-- TOC -->

- [SafeU Backend](#safeu-backend)
  - [Programming Language](#programming-language)
  - [Config](#config)
    - [Databases Config](#databases-config)
      - [MariaDB](#mariadb)
      - [Redis](#redis)
    - [Cloud Services Config](#cloud-services-config)
  - [Security](#security)
    - [CORS](#cors)
    - [CSRF](#csrf)
  - [Error Codes](#error-codes)
  - [Release Mode](#release-mode)
  - [Deploy (Development Environment)](#deploy-development-environment)
    - [Build from Source](#build-from-source)
    - [Run via Docker](#run-via-docker)
    - [Run via Docker Compose (Recommend)](#run-via-docker-compose-recommend)
    - [Scripts](#scripts)
  - [Deploy (Production Environment)](#deploy-production-environment)
    - [Build from Source](#build-from-source-1)
    - [Run via Docker](#run-via-docker-1)
    - [Run via Docker Compose (Recommend)](#run-via-docker-compose-recommend-1)
  - [FaaS](#faas)
  - [Code of Conducts](#code-of-conducts)
  - [API Documentation](#api-documentation)
  - [Development Workflow](#development-workflow)
  - [License](#license)

<!-- /TOC -->

## Programming Language

编程语言为 [`Golang`](https://github.com/golang/go) ，使用 [`govendor`](https://github.com/kardianos/govendor) 作为包管理工具，提交的代码需要使用 `gofmt` 进行格式化。为了开发维护方便，各位同学在编写代码时请注意起名规范及注释编写。该项目中变量命名使用 [Camel-Case](https://zh.wikipedia.org/wiki/%E9%A7%9D%E5%B3%B0%E5%BC%8F%E5%A4%A7%E5%B0%8F%E5%AF%AB)，并建议使用 `JetBrains` 公司的 `GoLand` IDE 进行开发。

## Config

### Databases Config

在 `conf/` 下新建 `db.json` 文件，写入以下配置信息：

```json
{
  "Master": {
    "User": "your_db_user_name",
    "Pass": "your_db_password",
    "Host": "your_database_ip",
    "Port": "your_database_port",
    "Database": "safeu",
    "MaxIdleConns": 30,
    "MaxOpenConns": 100,
    "ConnMaxLifetime": 3600,
    "Debug": false
  },
  "Redis": {
    "Host": "your_redis_host",
    "Port": "your_redis_port",
    "Pass": "your_redis_password"
  }
}
```

> 具体请参照 `conf/db.example.json` 。

#### MariaDB

在 `MySQL / MariaDB` 中新建 `safeu` 数据库：

```sql
CREATE DATABASE safeu;
```

#### Redis

本应用依赖 Redis，需要在 `conf/db.json` 中填好 Redis 相关配置。

### Cloud Services Config

云有关配置。

在 `conf/` 下新建并填写 `cloud.yml` 文件，具体请参照 `conf/cloud.exmaple.yml` 。

## Security

安全设计。

### CORS

Cross-origin Resource Sharing 跨源资源共享。

- 可共享域名由 [CORS_ALLOW_ORIGINS](common/const.go#L65) 决定；
- 可返回的响应头由 [CORS_ALLOW_HEADERS](common/const.go#L72) 决定；
- 允许的 HTTP 方法由 [CORS_ALLOW_METHODS](common/const.go#L80) 决定。

### CSRF

Cross-site Request Forgery 跨域请求伪造防护设计。

先请求 `/csrf` 接口获得 CSRF 口令。当发送 POST 请求时，必须在请求头中加入 `X-CSRF-TOKEN` CSRF 认证头，值为获得到的 CSRF 口令，否则会得到 [10007](ERRORS.md#L15) 错误。

## Error Codes

返回错误码。

错误码对照文档：[ERRORS](ERRORS.md)

## Release Mode

若用于生产环境部署，则建议使用 `RELEASE` 模式：

将 `common/const.go` 中的 `DEBUG` 变量更改为 `false` ，再重新构建容器镜像/运行即可。

## Deploy (Development Environment)

### Build from Source

```bash
$ git clone https://github.com/Triple-Z/safeu-backend.git
$ cd safeu-backend/
$ go get -u github.com/kardianos/govendor
$ govendor sync
$ go run main.go
```

> 要事先做好数据库的建立和配置。

### Run via Docker

通过 Docker 容器来运行应用。

在 `scripts/` 文件夹中有这样两个脚本：`build-docker-images.sh` 和 `run-dev-docker-containers.sh` 。前者用于构建本应用的 Docker 容器镜像，后者用于启动该应用。但要注意，在该这种方式中未含有数据库容器，因此需要额外提供数据库支持，提前编写 `conf/db.json` 数据库配置文件。

```bash
$ cd scripts/
$ ./build-docker-images.sh
$ ./run-dev-docker-containers.sh
```

> 要事先做好数据库的建立和配置。

### Run via Docker Compose (Recommend)

通过 Docker Compose 编排运行应用。

该方法中已自带 MariaDB 数据库，无需额外提供数据库服务。

```bash
$ cd scripts/
$ ./dev-docker-compose.sh up
```

> 需要安装 `docker` 和 `docker-Publicompose` ，可用 `scripts/install-docker-family-on-ubuntuPublic1804.sh` 在 Ubuntu 18.04 中安装 Docker 和 DPubliccker Compose。该脚本为中国网络环境额外定制，Public证安装速度。

### Scripts

所有的脚本文件应在脚本目录下运行。如 `scripts/build-docker-images.sh` 应这样运行：

```bash
$ cd scripts
$./build-docker-images.sh
```

## Deploy (Production Environment)

### Build from Source

与 [Deployment (Development Environment) > Build from Source](#build-from-source) 方法相同。

### Run via Docker

通过 Docker 容器来运行应用。

在 `scripts/` 文件夹中有这样两个脚本：`build-docker-images.sh` 和 `run-production-docker-containers.sh` 。前者用于构建本应用的 Docker 容器镜像，后者用于启动该应用。但要注意，在该这种方式中未含有数据库容器，因此需要额外提供数据库支持，需编写 `conf/db.json` 数据库配置文件及 `conf/cloud.yml` 云配置文件。

```bash
$ cd scripts/
$ ./build-docker-images.sh
$ ./run-production-docker-containers.sh
```

> 要事先做好数据库的建立和配置；
> 
> 若用于生产环境部署，建议使用 `RELEASE` 模式，详情参照 [Release Mode](#release-mode)。

### Run via Docker Compose (Recommend)

通过 Docker Compose 编排运行应用。

该方法中已自带 MariaDB 数据库，无需额外提供数据库服务。

```bash
$ cd scripts/
$ ./prod-docker-compose.sh up
```

> 需要安装 `docker` 和 `docker-compose` ，可用 `scripts/install-docker-family-on-ubuntu-1804.sh` 在 Ubuntu 18.04 中安装 Docker 和 Docker Compose。该脚本为中国网络环境额外定制，保证安装速度。

## FaaS

本项目使用阿里云函数计算以实现文件的压缩功能，函数代码见 `faas/zip-items.py`。若使用同一用户的 OSS 资源，必须基于该用户部署函数计算并且授予 OSS 相关管理权限。

## Code of Conducts

## API Documentation

API 文档位于 `api/` 文件夹下，采用 [OpenAPI 3.0.0](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.0.md) 标准编写。

Online Documentation: https://app.swaggerhub.com/apis-docs/a2os/safeu

> 更多信息见 [SafeU 内部文档](https://docs.google.com/document/d/1UiFHogsqDSqw3SrEAnEMukOoJq3fyxqIIEP-OE7ask0/edit?ts=5c56f70d) **需要权限**。

## Development Workflow

## License

该项目采用 [`Apache 2.0`](LICENSE) 许可证。

```
   Copyright 2019 A2OS SafeU Dev Team

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
```
