# IAM Kit
[![CI](https://github.com/che-kwas/iam-kit/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/che-kwas/iam-kit/actions/workflows/ci.yaml)

## 错误码

### 格式

错误码为`6位数字`，每2位为1组，分别代表`服务`(1-2位) + `模块`(3-4位) + `错误编号`(5-6位)。示例：

```sh
# 10（通用错误）+ 01（认证/授权模块）+ 05（token过期）
100105
```

### 错误码段

1. 通用：`10XXYY`
2. iam-apiserver：`11XXYY`
3. iam-authz-server：`12XXYY`
4. iam-pump：`13XXYY`

## TODO

1. init最佳实践
2. config最佳实践
  - load config独立出来，不只是server需要load
3. errors最佳实践（参考kratos公众号文章）
4. logger最佳实践
  - withValues
  - WithContext
  - FromContext
  - global singleton
  - L()
5. mysql logger
6. test
7. policy audit
8. migrate
9. metrics
10. profiling
11. version
12. list selector
13. server可以只创建http/grpc server
