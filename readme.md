### 致谢
- webdav部分逻辑来自github.com/hacdias/webdav
- 图形界面采用walk

### TLS证书
- TLS_CERT = "C:\\nginx-1.18.0\\conf\\nolva.pem"

### TLS私钥
- TLS_KEY = "C:\\nginx-1.18.0\\conf\\nolva.key"

### webdav.yaml
```yaml
#  Ip   string `yaml:"ip,omitempty"`
#  Port uint16 `yaml:"port"`
#
#  Auth bool   `yaml:"auth"`
#  User string `yaml:"user,omitempty"`
#  Pass string `yaml:"pass,omitempty"`
#
#  Scope  string `yaml:"scope"`
#  Modify bool   `yaml:"modify"`
#
#  Tls  bool   `yaml:"tls"`
#  Cert string `yaml:"cert,omitempty"`
#  Key  string `yaml:"key,omitempty"`

log: webdav.log

# 默认服务器及密码，用于根目录
default:
  user: nolva
  auth: true
  pass: "{bcrypt}..."
  port: 8082
  scope: .
  tls: true
  modify: true


# 其他服务器
#server:
#  -
#    user: admin
#    auth: true
#    pass: admin
#    port: 8080
#    scope: .
#    tls: true
#    modify: false
```
