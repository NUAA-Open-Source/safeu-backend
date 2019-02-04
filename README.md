# SafeU Backend

伏羲后端仓库。

> README for developers.

## Programming Language

编程语言为 [`Golang`](https://github.com/golang/go) ，使用 [`govendor`](https://github.com/kardianos/govendor) 作为包管理工具，提交的代码需要使用 `gofmt` 进行格式化。为了开发维护方便，各位同学在编写代码时请注意起名规范及注释编写。该项目中变量命名使用 [Camel-Case](https://zh.wikipedia.org/wiki/%E9%A7%9D%E5%B3%B0%E5%BC%8F%E5%A4%A7%E5%B0%8F%E5%AF%AB)，并建议使用 `JetBrains` 公司的 `GoLand` IDE 进行开发。

## Config

### Database

在 `MySQL / MariaDB` 中新建 `safeu` 数据库：

```sql
CREATE DATABASE safeu;
```

在 `conf/` 下新建 `db.json` 文件，写入以下配置信息：

```json
{
  "Master": {
    "User": "your_user_name",
    "Pass": "your_password",
    "Host": "127.0.0.1",
    "Port": "3306",
    "Database": "safeu",
    "MaxIdleConns": 30,
    "MaxOpenConns": 100,
    "Debug": false
  }
}
```

## Quick Start

```bash
$ git clone https://github.com/Triple-Z/safeu-backend.git
$ cd safeu-backend/
$ go get -u github.com/kardianos/govendor
$ govendor sync
$ go run main.go
```

> 要事先做好数据库的建立和配置。

## Code of Conducts

## API Documentation

API 文档位于 `api/` 文件夹下，采用 [OpenAPI 3.0.0](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.0.md) 标准编写。

Online Documentation: https://app.swaggerhub.com/apis/a2os/safeu

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
