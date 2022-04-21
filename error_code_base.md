## 错误码列表

| Identifier | Code | HTTP Code | Description |
| ---------- | ---- | --------- | ----------- |
| ErrSuccess | 100001 | 200 | OK |
| ErrUnknown | 100002 | 500 | Internal server error |
| ErrBadParams | 100003 | 400 | Bad request parameters |
| ErrNotFound | 100004 | 404 | Not found |
| ErrPasswordInvalid | 100101 | 401 | Password invalid |
| ErrHeaderInvalid | 100102 | 401 | Authorization header invalid |
| ErrSignatureInvalid | 100103 | 401 | Signature invalid |
| ErrTokenInvalid | 100104 | 401 | Token invalid |
| ErrTokenExpired | 100105 | 401 | Token expired |
| ErrPermissionDenied | 100106 | 403 | Permission denied |
| ErrDatabase | 100201 | 500 | Database error |

