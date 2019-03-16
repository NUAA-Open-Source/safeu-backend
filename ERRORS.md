# Error Code Documentation

## System Errors

1xxxx 为系统级错误。

| Error Code | Error Message (en) | Error Message (zh) |
| :----: | :----: | :----: |
| 10001 | System error | 系统错误 |
| 10002 | Service unavailable | 服务暂停 |
| 10003 | Parameter error | 参数错误 |
| 10004 | Parameter value invalid | 参数非法 |
| 10005 | Missing required parameter | 缺少参数 |
| 10006 | Resource unavailable | 资源不存在 |
| 10007 | CSRF token mismatch | CSRF 认证失败 |

## Application Errors

2xxxx 为应用级错误。

### General Errors

| Error Code | Error Message (en) | Error Message (zh) |
| :----: | :----: | :----: |
| 20000 | General error | 通用应用错误（以返回错误信息为准） |

### Upload Errors

201xx 为上传相关错误。

### Update Errors

202xx 为更新相关错误。

| Error Code | Error Message (en) | Error Message (zh) |
| :----: | :----: | :----: |
| 20201 | Can't find user token | 无法找到 Token |

### Download Errors

203xx 为下载相关错误。

| Error Code | Error Message (en) | Error Message (zh) |
| :----: | :----: | :----: |
| 20301 | Missing token in header | 请求头缺少 Token |
| 20302 | Token used | Token 已被使用 |
| 20303 | Token expired | Token 已过期 |
| 20304 | Token revoked | Token 不合法 |
| 20305 | Can't get the download link | 无法获取下载链接 |
| 20306 | The retrieve code mismatch auth | 提取码无法对应auth |
| 20307 | The retrieve code repeat | 提取码重复 |
### Delete Errors

204xx 为删除相关错误。

### Validate Errors

205xx 为认证相关错误。

| Error Code | Error Message (en) | Error Message (zh) |
| :----: | :----: | :----: |
| 20501 | Incorrect password | 密码错误 |
