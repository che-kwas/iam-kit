[![CI](https://github.com/che-kwas/iam-kit/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/che-kwas/iam-kit/actions/workflows/ci.yaml)

# IAM Kit

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

### 通用错误码

[error_code_base.md](error_code_base.md)

## TODO

- test
- metrics
- profiling
- version
- validate
- README
