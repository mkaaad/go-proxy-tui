# mihomo RESTful API 接口文档

> 本文档基于 mihomo (`MetaCubeX/mihomo`, `Meta` 分支) 源码 `hub/route/` 中的路由定义整理,覆盖外部控制器 (External Controller) 提供的全部 RESTful 接口。

## 1. 概述

mihomo 内置了一个完整的 RESTful API 服务(即「外部控制器」),供外部程序、Web 面板(如 metacubexd / yacd)和脚本对内核进行监控与控制,包括:

- 查看/修改运行配置
- 查询与切换代理、策略组
- 查询与动态禁用规则
- 查看/关闭活跃连接
- 实时流量、内存、日志推送(SSE 与 WebSocket)
- DNS 查询、缓存刷新、代理 Provider 管理、内核重启与升级等

### 1.1 启用方式

在 mihomo 配置文件的 `general` 段中开启 `external-controller` 即可:

```yaml
external-controller: 127.0.0.1:9090   # HTTP 监听地址
secret: "your-secret"                 # 鉴权令牌(建议设置)
external-controller-tls: 127.0.0.1:9443   # 可选,HTTPS
external-controller-unix: mihomo.sock     # 可选,Unix Socket(Linux/macOS)
external-controller-pipe: \\.\pipe\mihomo # 可选,Windows 命名管道
external-doh-server: /dns-query           # 可选,DoH 服务路径
```

- 访问 Unix Socket / 命名管道时**不校验** `secret`。
- 若 `external-controller` 绑定到非回环地址,必须设置 `secret` 或使用 TLS。

## 2. 通用约定

### 2.1 Base URL

```
http://<external-controller-addr>/
```

本文示例均以 `http://127.0.0.1:9090` 为例。

### 2.2 鉴权

设置了 `secret` 后,所有 HTTP 请求需要在请求头携带 Bearer Token:

```
Authorization: Bearer <secret>
```

WebSocket 请求(浏览器场景)无法自定义请求头,改为通过查询参数传 token:

```
ws://127.0.0.1:9090/traffic?token=<secret>
```

未鉴权或鉴权失败时返回 `401`:

```json
{ "message": "Unauthorized" }
```

### 2.3 请求与响应

- 请求体均为 `application/json`。
- 成功时根据接口返回 JSON、`204 No Content` 或流式数据。
- 失败时返回非 2xx 状态码 + JSON 错误体:

```json
{ "message": "Body invalid" }
```

常见错误状态码见 [第 12 节](#12-错误码)。

## 3. 接口索引

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/` | 探活,返回 `{"hello":"mihomo"}` |
| GET | `/version` | 获取内核版本 |
| GET | `/logs` | 实时日志流(SSE/WS) |
| GET | `/traffic` | 实时流量(SSE/WS) |
| GET | `/memory` | 实时内存占用(SSE/WS) |
| GET | `/configs` | 获取当前运行配置 |
| PUT | `/configs` | 重载配置(支持 path / payload) |
| PATCH | `/configs` | 部分更新运行配置 |
| POST | `/configs/geo` | 更新 GEO 数据库(GeoIP/GeoSite/MMDB/ASN) |
| GET | `/proxies` | 获取全部代理 |
| GET | `/proxies/:name` | 获取单个代理/策略组详情 |
| PUT | `/proxies/:name` | 切换选择器(Selector)选中的代理 |
| DELETE | `/proxies/:name` | 重置 URLTest 组的固定选择(取消锁定) |
| GET | `/proxies/:name/delay` | 测试单个代理延迟 |
| GET | `/group` | 获取全部策略组 |
| GET | `/group/:name` | 获取单个策略组 |
| GET | `/group/:name/delay` | 测试策略组全部节点延迟 |
| GET | `/rules` | 获取全部规则 |
| PATCH | `/rules/disable` | 按索引动态禁用/启用规则 |
| GET | `/connections` | 获取活跃连接快照(SSE/WS) |
| DELETE | `/connections` | 关闭全部连接 |
| DELETE | `/connections/:id` | 关闭指定连接 |
| GET | `/providers/proxies` | 获取全部代理 Provider |
| GET | `/providers/proxies/:name` | 获取单个代理 Provider |
| PUT | `/providers/proxies/:name` | 手动更新代理 Provider |
| GET | `/providers/proxies/:name/healthcheck` | 触发 Provider 健康检查 |
| GET | `/providers/proxies/:name/:proxy` | 获取 Provider 内指定代理 |
| GET | `/providers/proxies/:name/:proxy/healthcheck` | 测试 Provider 内代理延迟 |
| GET | `/providers/rules` | 获取全部规则 Provider |
| PUT | `/providers/rules/:name` | 手动更新规则 Provider |
| GET | `/dns/query` | 通过内核 DNS 解析域名 |
| POST | `/cache/fakeip/flush` | 清空 FakeIP 池 |
| POST | `/cache/dns/flush` | 清空 DNS 缓存 |
| GET | `/storage/:key` | 读取 key-value 存储 |
| PUT | `/storage/:key` | 写入 key-value 存储(JSON,≤1MB) |
| DELETE | `/storage/:key` | 删除 key-value 存储 |
| POST | `/restart` | 重启内核 |
| POST | `/upgrade` | 升级内核二进制 |
| POST | `/upgrade/geo` | 更新 GEO 数据库 |
| POST | `/upgrade/ui` | 更新 external-ui 面板资源 |
| PUT | `/debug/gc` | 手动触发 GC(需 log-level=debug) |
| GET | `/debug/pprof/*` | Go pprof 性能分析(需 log-level=debug) |
| GET/WS | `/ui` | external-ui 静态面板资源(需配置) |

## 4. 系统与版本

### 4.1 GET `/` — 探活

返回固定内容,可用于健康检查。

```bash
curl -s http://127.0.0.1:9090/
```

```json
{ "hello": "mihomo" }
```

### 4.2 GET `/version` — 版本信息

```bash
curl -s http://127.0.0.1:9090/version
```

```json
{
  "meta": true,
  "version": "v1.19.0"
}
```

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `meta` | bool | 是否为 mihomo(Meta 内核) |
| `version` | string | 版本号 |

## 5. 配置管理

### 5.1 GET `/configs` — 获取运行配置

返回当前运行时通用配置。

```bash
curl -s http://127.0.0.1:9090/configs
```

```json
{
  "port": 7890,
  "socks-port": 7891,
  "redir-port": 0,
  "tproxy-port": 0,
  "mixed-port": 0,
  "tun": {},
  "tuic-server": {},
  "ss-config": "",
  "vmess-config": "",
  "authentication": [],
  "skip-auth-prefixes": [],
  "lan-allowed-ips": [],
  "lan-disallowed-ips": [],
  "allow-lan": false,
  "bind-address": "*",
  "inbound-tfo": false,
  "inbound-mptcp": false,
  "mode": "rule",
  "unified-delay": true,
  "log-level": "info",
  "ipv6": true,
  "interface-name": "",
  "routing-mark": 0,
  "geox-url": {},
  "geo-auto-update": false,
  "geo-update-interval": 24,
  "geodata-mode": false,
  "geodata-loader": "memconservative",
  "geosite-matcher": "mph",
  "tcp-concurrent": false,
  "find-process-mode": "off",
  "sniffing": false,
  "global-ua": "",
  "etag-support": true,
  "keep-alive-idle": 30,
  "keep-alive-interval": 30,
  "disable-keep-alive": false
}
```

主要字段说明:

| 字段 | 说明 |
| --- | --- |
| `port` / `socks-port` / `redir-port` / `tproxy-port` / `mixed-port` | 各入站端口,0 表示未启用 |
| `tun` | TUN 模式配置对象 |
| `allow-lan` | 是否允许局域网访问 |
| `bind-address` | 入站监听地址 |
| `mode` | 运行模式:`rule` / `global` / `direct` |
| `log-level` | 日志级别:`debug` / `info` / `warning` / `error` |
| `ipv6` | 是否启用 IPv6 |
| `unified-delay` | 是否使用统一延迟 |
| `sniffing` | 是否启用流量嗅探 |
| `find-process-mode` | 进程查找模式:`off` / `strict` / `always` |
| `routing-mark` | 路由标记 |
| `geox-url` | GEO 数据库下载地址 |

### 5.2 PUT `/configs` — 重载配置

从配置文件路径或直接传入配置内容重载内核配置。

- 请求体(二选一):

```json
{ "path": "/etc/mihomo/config.yaml" }
```

```json
{ "payload": "mixed-port: 7893\nmode: rule\n..." }
```

- `path` 必须为绝对路径,且需位于 mihomo 工作目录或 `SAFE_PATHS` 白名单内。
- 可选查询参数 `?force=true`,强制应用配置(跳过部分校验)。

```bash
curl -X PUT http://127.0.0.1:9090/configs?force=true \
  -H "Content-Type: application/json" \
  -d '{"path": "/etc/mihomo/config.yaml"}'
```

成功返回 `204 No Content`。

### 5.3 PATCH `/configs` — 部分更新配置

只修改传入的字段,未传字段保持不变。可用于动态切换模式、端口、开关等。

```bash
curl -X PATCH http://127.0.0.1:9090/configs \
  -H "Content-Type: application/json" \
  -d '{"mode": "global", "allow-lan": true}'
```

可更新的顶层字段:

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `port` | int | HTTP 代理端口 |
| `socks-port` | int | SOCKS5 代理端口 |
| `redir-port` | int | 透明代理(重定向)端口 |
| `tproxy-port` | int | TPROXY 端口 |
| `mixed-port` | int | 混合代理端口 |
| `tun` | object | TUN 配置(见下) |
| `tuic-server` | object | TUIC 服务端配置 |
| `ss-config` | string | Shadowsocks 入站配置 |
| `vmess-config` | string | Vmess 入站配置 |
| `allow-lan` | bool | 是否允许局域网 |
| `skip-auth-prefixes` | string[] | 免鉴权 CIDR 前缀 |
| `lan-allowed-ips` | string[] | 允许访问的局域网 IP |
| `lan-disallowed-ips` | string[] | 禁止访问的局域网 IP |
| `bind-address` | string | 入站监听地址 |
| `mode` | string | 模式:`rule` / `global` / `direct` |
| `log-level` | string | 日志级别 |
| `ipv6` | bool | 是否启用 IPv6 |
| `sniffing` | bool | 是否启用嗅探 |
| `tcp-concurrent` | bool | TCP 并发连接 |
| `find-process-mode` | string | 进程查找模式 |
| `interface-name` | string | 出站接口名 |

`tun` 对象常用字段:

```json
{
  "enable": true,
  "device": "Mihomo",
  "stack": "mixed",
  "dns-hijack": ["any:53"],
  "auto-route": true,
  "auto-detect-interface": true,
  "mtu": 9000,
  "strict-route": false,
  "route-address": [],
  "route-exclude-address": [],
  "endpoint-independent-nat": false,
  "udp-timeout": 120,
  "file-descriptor": 0
}
```

成功返回 `204 No Content`。

> 本项目 `internal/client/mihomo/client.go` 中 `SwitchModes` 即通过该接口切换模式:
> ```go
> m.apiClient.Patch("/configs", map[string]string{"mode": mode})
> ```

### 5.4 POST `/configs/geo` — 更新 GEO 数据库

同步更新 GeoIP / GeoLite2 MMDB / ASN / GeoSite 数据库,成功返回 `204 No Content`。

```bash
curl -X POST http://127.0.0.1:9090/configs/geo
```

## 6. 代理

### 6.1 GET `/proxies` — 获取全部代理

返回所有代理与策略组,key 为代理名。

```bash
curl -s http://127.0.0.1:9090/proxies
```

```json
{
  "proxies": {
    "PROXY": {
      "type": "Selector",
      "name": "PROXY",
      "history": [],
      "all": ["🇺🇸 US-01", "🇯🇵 JP-01", "DIRECT"],
      "now": "🇺🇸 US-01",
      "udp": true,
      "xudp": false,
      "tfo": false,
      "mptcp": false,
      "alive": true,
      "delay": 128
    },
    "🇺🇸 US-01": {
      "type": "Vless",
      "name": "🇺🇸 US-01",
      "history": [
        { "time": "2026-08-21T18:00:00+08:00", "delay": 128 }
      ],
      "udp": true,
      "xudp": false,
      "tfo": false,
      "mptcp": false,
      "alive": true,
      "delay": 128
    }
  }
}
```

字段说明:

| 字段 | 说明 |
| --- | --- |
| `type` | 代理类型:`Selector` / `URLTest` / `Fallback` / `LoadBalance` / `Direct` / `Reject` / `Vmess` / `Vless` / `Trojan` / `Shadowsocks` / `Hysteria2` / `Tuic` / `Wireguard` 等 |
| `name` | 代理名 |
| `history` | 延迟测试历史记录 |
| `all` | 仅策略组:包含的全部可选代理 |
| `now` | 仅选择器:当前选中的代理 |
| `udp` / `xudp` / `tfo` / `mptcp` | 协议特性支持情况 |
| `alive` | 节点是否可用(最近延迟测试结果) |
| `delay` | 最近一次延迟(ms) |

### 6.2 GET `/proxies/:name` — 获取单个代理

```bash
curl -s http://127.0.0.1:9090/proxies/PROXY
```

返回该代理/策略组的完整对象(结构同 `GET /proxies` 中的单项)。

### 6.3 PUT `/proxies/:name` — 切换选择器选中的代理

请求体指定目标代理名(必须是该选择器 `all` 列表中的成员)。

```bash
curl -X PUT http://127.0.0.1:9090/proxies/PROXY \
  -H "Content-Type: application/json" \
  -d '{"name": "🇯🇵 JP-01"}'
```

- 仅 `Selector` 类型支持,否则返回 `400`。
- 成功返回 `204 No Content`,选择会持久化到缓存文件,重启后保留。

### 6.4 DELETE `/proxies/:name` — 取消固定选择

对 `URLTest` 等可「固定选择」的非 Selector 策略组生效:清空手动固定,恢复自动测速选择。

```bash
curl -X DELETE http://127.0.0.1:9090/proxies/🚀 节点选择
```

成功返回 `204 No Content`。

### 6.5 GET `/proxies/:name/delay` — 测试代理延迟

| 查询参数 | 必填 | 说明 |
| --- | --- | --- |
| `url` | 是 | 测速目标 URL |
| `timeout` | 是 | 超时时间(ms) |
| `expected` | 否 | 期望的状态码/区间,如 `200`、`200,300`、`200-299` |

```bash
curl -s "http://127.0.0.1:9090/proxies/🇺🇸%20US-01/delay?url=http://www.gstatic.com/generate_204&timeout=5000"
```

```json
{ "delay": 128 }
```

- 超时返回 `504`;测速失败/不可用返回 `503`。

## 7. 策略组

策略组指 `Selector` / `URLTest` / `Fallback` / `LoadBalance` 类型的代理集合。

### 7.1 GET `/group` — 获取全部策略组

仅返回策略组,结构为 `{"proxies": [...]}`(内容同代理对象)。

```bash
curl -s http://127.0.0.1:9090/group
```

```json
{
  "proxies": [
    {
      "type": "Selector",
      "name": "PROXY",
      "all": ["🇺🇸 US-01", "🇯🇵 JP-01", "DIRECT"],
      "now": "🇺🇸 US-01",
      "history": [],
      "udp": true
    }
  ]
}
```

### 7.2 GET `/group/:name` — 获取单个策略组

```bash
curl -s http://127.0.0.1:9090/group/PROXY
```

非策略组返回 `404`。

### 7.3 GET `/group/:name/delay` — 测试组内全部节点延迟

| 查询参数 | 必填 | 说明 |
| --- | --- | --- |
| `url` | 是 | 测速目标 URL |
| `timeout` | 是 | 超时时间(ms) |
| `expected` | 否 | 期望的状态码/区间 |

返回 `{"代理名": 延迟ms}` 的映射:

```bash
curl -s "http://127.0.0.1:9090/group/PROXY/delay?url=http://www.gstatic.com/generate_204&timeout=5000"
```

```json
{
  "🇺🇸 US-01": 128,
  "🇯🇵 JP-01": 156,
  "DIRECT": 32
}
```

> 对 URLTest 组调用时,若此前手动固定过节点,会先清除固定再测速。

## 8. 规则

### 8.1 GET `/rules` — 获取全部规则

```bash
curl -s http://127.0.0.1:9090/rules
```

```json
{
  "rules": [
    {
      "type": "DOMAIN-SUFFIX",
      "payload": "google.com",
      "proxy": "PROXY",
      "size": 0,
      "index": 0,
      "extra": {
        "disabled": false,
        "hitCount": 0,
        "hitAt": "0001-01-01T00:00:00Z",
        "missCount": 0,
        "missAt": "0001-01-01T00:00:00Z"
      }
    }
  ]
}
```

| 字段 | 说明 |
| --- | --- |
| `type` | 规则类型:`DOMAIN` / `DOMAIN-SUFFIX` / `DOMAIN-KEYWORD` / `IP-CIDR` / `GEOIP` / `GEOSITE` / `MATCH` / `PROCESS-NAME` 等 |
| `payload` | 规则内容 |
| `proxy` | 命中后使用的出站 |
| `size` | 仅 GEOIP/GEOSITE 规则:记录条数 |
| `index` | 规则下标(0 起,配置文件中顺序) |
| `extra` | 附加统计:是否被动态禁用、命中/未命中次数与时间 |

### 8.2 PATCH `/rules/disable` — 动态禁用/启用规则

请求体为「规则下标 → 是否禁用」的映射:

```bash
# 禁用下标为 3 的规则
curl -X PATCH http://127.0.0.1:9090/rules/disable \
  -H "Content-Type: application/json" \
  -d '{"3": true}'

# 重新启用
curl -X PATCH http://127.0.0.1:9090/rules/disable \
  -H "Content-Type: application/json" \
  -d '{"3": false}'
```

- 下标越界的条目会被忽略。
- 该操作为运行时生效,重启后重置。
- 成功返回 `204 No Content`。

## 9. 连接管理

### 9.1 GET `/connections` — 获取活跃连接

支持两种方式:

- **HTTP**(一次性快照):

```bash
curl -s http://127.0.0.1:9090/connections
```

- **WebSocket**(周期性推送快照,`interval` 指定推送间隔 ms,默认 1000):

```
ws://127.0.0.1:9090/connections?token=<secret>&interval=2000
```

响应快照结构:

```json
{
  "downloadTotal": 1024000,
  "uploadTotal": 512000,
  "connections": [
    {
      "id": "b4a3f2e1-...",
      "metadata": {
        "network": "tcp",
        "type": "HTTP",
        "sourceIP": "127.0.0.1",
        "destinationIP": "142.250.72.14",
        "sourcePort": "53210",
        "destinationPort": "443",
        "host": "www.google.com",
        "sniffHost": "www.google.com",
        "dnsMode": "normal",
        "processPath": "/usr/bin/curl",
        "specialProxy": "",
        "specialRules": ""
      },
      "upload": 0,
      "download": 0,
      "start": "2026-08-21T18:00:00.123456789+08:00",
      "chains": ["PROXY", "🇺🇸 US-01"],
      "providerChains": ["subscription-a"],
      "rule": "DomainSuffix",
      "rulePayload": "google.com"
    }
  ],
  "memory": 52428800
}
```

字段说明:

| 字段 | 说明 |
| --- | --- |
| `downloadTotal` / `uploadTotal` | 累计下载/上传字节数 |
| `connections[]` | 活跃连接数组 |
| `id` | 连接 UUID |
| `metadata` | 连接元信息(见下) |
| `upload` / `download` | 该连接已传输字节数 |
| `start` | 连接建立时间 |
| `chains` | 出站链路(策略组 → 实际节点) |
| `providerChains` | 命中的 Provider 链路 |
| `rule` / `rulePayload` | 命中规则类型与内容 |
| `memory` | 当前内存占用(字节) |

`metadata` 字段:

| 字段 | 说明 |
| --- | --- |
| `network` | `tcp` / `udp` |
| `type` | 连接类型:`HTTP` / `SOCKS` / `CONNECT` 等 |
| `sourceIP` / `sourcePort` | 源地址/端口 |
| `destinationIP` / `destinationPort` | 目标地址/端口 |
| `host` | 请求 Host |
| `sniffHost` | 嗅探到的域名 |
| `dnsMode` | DNS 模式 |
| `processPath` | 发起连接的进程路径(需开启 `find-process-mode`) |

### 9.2 DELETE `/connections` — 关闭全部连接

```bash
curl -X DELETE http://127.0.0.1:9090/connections
```

成功返回 `204 No Content`。

### 9.3 DELETE `/connections/:id` — 关闭指定连接

```bash
curl -X DELETE http://127.0.0.1:9090/connections/b4a3f2e1-...
```

成功返回 `204 No Content`(即使连接已不存在)。

## 10. Provider

### 10.1 代理 Provider

#### GET `/providers/proxies` — 全部代理 Provider

```bash
curl -s http://127.0.0.1:9090/providers/proxies
```

```json
{
  "providers": {
    "subscription-a": {
      "name": "subscription-a",
      "type": "Proxy",
      "vehicleType": "HTTP",
      "updatedAt": "2026-08-21T10:00:00+08:00",
      "proxies": [
        { "type": "Vless", "name": "🇺🇸 US-01", "udp": true }
      ],
      "testUrl": "http://www.gstatic.com/generate_204",
      "subscriptionInfo": {
        "upload": 0,
        "download": 1073741824,
        "total": 10737418240,
        "expire": 1798560000
      }
    }
  }
}
```

#### GET `/providers/proxies/:name` — 单个 Provider

```bash
curl -s http://127.0.0.1:9090/providers/proxies/subscription-a
```

#### PUT `/providers/proxies/:name` — 手动更新 Provider

拉取订阅/读取本地文件刷新 Provider,成功返回 `204 No Content`。

```bash
curl -X PUT http://127.0.0.1:9090/providers/proxies/subscription-a
```

#### GET `/providers/proxies/:name/healthcheck` — 触发健康检查

```bash
curl -s http://127.0.0.1:9090/providers/proxies/subscription-a/healthcheck
```

成功返回 `204 No Content`。

#### GET `/providers/proxies/:name/:proxy` — Provider 内单个代理

```bash
curl -s "http://127.0.0.1:9090/providers/proxies/subscription-a/🇺🇸%20US-01"
```

#### GET `/providers/proxies/:name/:proxy/healthcheck` — 测试 Provider 内代理延迟

参数同 `GET /proxies/:name/delay`(`url` / `timeout` / `expected`):

```bash
curl -s "http://127.0.0.1:9090/providers/proxies/subscription-a/🇺🇸%20US-01/healthcheck?url=http://www.gstatic.com/generate_204&timeout=5000"
```

### 10.2 规则 Provider

#### GET `/providers/rules` — 全部规则 Provider

```bash
curl -s http://127.0.0.1:9090/providers/rules
```

```json
{
  "providers": {
    "rule-provider-a": {
      "name": "rule-provider-a",
      "type": "Rule",
      "vehicleType": "HTTP",
      "updatedAt": "2026-08-21T10:00:00+08:00",
      "ruleCount": 100
    }
  }
}
```

#### PUT `/providers/rules/:name` — 更新规则 Provider

```bash
curl -X PUT http://127.0.0.1:9090/providers/rules/rule-provider-a
```

成功返回 `204 No Content`。

## 11. 实时数据流(SSE / WebSocket)

`/traffic`、`/memory`、`/logs` 三个接口支持两种消费方式:

1. **HTTP + `text/event-stream`**:以 chunked 流持续输出 JSON,每行一条记录。
2. **WebSocket**:通过 `Upgrade: websocket` 请求头或 `?token=` 建立 WS 连接后持续推送。

### 11.1 GET `/traffic` — 实时流量

每秒推送一次当前速率与累计流量:

```bash
curl -N http://127.0.0.1:9090/traffic
```

```json
{ "up": 1024, "down": 4096, "upTotal": 512000, "downTotal": 1024000 }
```

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `up` / `down` | int64 | 当前上行/下行速率(字节/秒) |
| `upTotal` / `downTotal` | int64 | 累计上行/下行字节数 |

WebSocket:

```
ws://127.0.0.1:9090/traffic?token=<secret>
```

### 11.2 GET `/memory` — 实时内存

每秒推送一次堆内存占用,首次推送为 `0`:

```bash
curl -N http://127.0.0.1:9090/memory
```

```json
{ "inuse": 52428800, "oslimit": 0 }
```

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `inuse` | uint64 | 当前 Go 堆内存使用(字节) |
| `oslimit` | uint64 | 系统内存上限(当前恒为 0,预留字段) |

### 11.3 GET `/logs` — 实时日志

| 查询参数 | 说明 |
| --- | --- |
| `level` | 最低日志级别,默认 `info`;可选 `debug` / `info` / `warning` / `error` |
| `format` | 设为 `structured` 时输出结构化日志 |

```bash
curl -N "http://127.0.0.1:9090/logs?level=info"
```

默认格式(每行一条):

```json
{ "type": "info", "payload": "[TCP] dial 1.2.3.4:443 success" }
```

结构化格式(`?format=structured`):

```json
{
  "time": "18:00:01",
  "level": "info",
  "message": "[TCP] dial 1.2.3.4:443 success",
  "fields": []
}
```

> `type` 取值:`debug` / `info` / `warning`(结构化输出时转为 `warn`)/ `error`。
> 注意 `warning` 级别在查询参数中为 `warning`,结构化输出中为 `warn`。

## 12. DNS 与缓存

### 12.1 GET `/dns/query` — 内核 DNS 解析

使用内核内置 resolver 查询域名。

| 查询参数 | 必填 | 默认 | 说明 |
| --- | --- | --- | --- |
| `name` | 是 | - | 要解析的域名 |
| `type` | 否 | `A` | 查询类型:`A` / `AAAA` / `CNAME` / `TXT` / `MX` 等 |

```bash
curl -s "http://127.0.0.1:9090/dns/query?name=www.google.com&type=A"
```

```json
{
  "Status": 0,
  "Question": [{ "name": "www.google.com.", "type": 1 }],
  "TC": false,
  "RD": true,
  "RA": true,
  "AD": false,
  "CD": false,
  "Answer": [
    { "name": "www.google.com.", "type": 1, "TTL": 300, "data": "142.250.72.14" }
  ]
}
```

- `Status` 为 DNS 返回码(0 = NOERROR)。
- `Answer` / `Authority` / `Additional` 分别对应应答、授权、附加区段。
- DNS 未启用时返回 `500`。

### 12.2 POST `/cache/fakeip/flush` — 清空 FakeIP 池

```bash
curl -X POST http://127.0.0.1:9090/cache/fakeip/flush
```

成功返回 `204 No Content`。

### 12.3 POST `/cache/dns/flush` — 清空 DNS 缓存

```bash
curl -X POST http://127.0.0.1:9090/cache/dns/flush
```

成功返回 `204 No Content`。

## 13. 存储(Key-Value)

用于面板/外部程序在 mihomo 缓存文件中保存自定义数据。

### 13.1 GET `/storage/:key` — 读取

```bash
curl -s http://127.0.0.1:9090/storage/my-key
```

返回存储的 JSON 值;不存在时返回 `null`。

### 13.2 PUT `/storage/:key` — 写入

请求体为任意合法 JSON,大小上限 1MB:

```bash
curl -X PUT http://127.0.0.1:9090/storage/my-key \
  -H "Content-Type: application/json" \
  -d '{"theme": "dark"}'
```

成功返回 `204 No Content`;非法 JSON 返回 `400`,超限返回 `413`。

### 13.3 DELETE `/storage/:key` — 删除

```bash
curl -X DELETE http://127.0.0.1:9090/storage/my-key
```

成功返回 `204 No Content`。

## 14. 重启与升级

> `embedMode` 下(以库方式嵌入时)`/restart` 及 `/upgrade`(内核/GEO)不可用。

### 14.1 POST `/restart` — 重启内核

以相同参数重启当前进程,返回 `{"status":"ok"}` 后进程立即重启:

```bash
curl -X POST http://127.0.0.1:9090/restart
```

```json
{ "status": "ok" }
```

### 14.2 POST `/upgrade` — 升级内核

| 查询参数 | 说明 |
| --- | --- |
| `channel` | 升级通道(如 `alpha` / `beta` / `stable`) |
| `force` | 设为 `true` 时强制升级(忽略版本比较) |

```bash
curl -X POST "http://127.0.0.1:9090/upgrade?channel=stable&force=true"
```

下载并替换二进制后自动重启,返回 `{"status":"ok"}`。

### 14.3 POST `/upgrade/geo` — 更新 GEO 数据库

同 `POST /configs/geo`:

```bash
curl -X POST http://127.0.0.1:9090/upgrade/geo
```

### 14.4 POST `/upgrade/ui` — 更新面板资源

下载并解压 `external-ui` 指向的面板资源:

```bash
curl -X POST http://127.0.0.1:9090/upgrade/ui
```

返回 `{"status":"ok"}`。

## 15. 调试接口

仅在 `log-level: debug` 时挂载 `/debug`。

### 15.1 PUT `/debug/gc` — 手动触发 GC

调用 Go 运行时 `debug.FreeOSMemory()`,立即回收内存:

```bash
curl -X PUT http://127.0.0.1:9090/debug/gc
```

### 15.2 GET `/debug/pprof/*` — 性能分析

标准 Go pprof 接口:

```bash
curl -s http://127.0.0.1:9090/debug/pprof/
curl -s http://127.0.0.1:9090/debug/pprof/heap > heap.out
curl -s http://127.0.0.1:9090/debug/pprof/goroutine?debug=1
```

可使用 `go tool pprof` 进行分析。

## 16. 静态面板与 DoH

### 16.1 `/ui` — 内置面板

配置了 `external-ui` 后,可通过 API 端口直接访问面板:

```
http://127.0.0.1:9090/ui/
```

### 16.2 DoH 服务

配置了 `external-doh-server`(如 `/dns-query`)后,该路径提供 RFC 8484 DoH 服务:

```
https://127.0.0.1:9443/dns-query?dns=<base64url 编码的 DNS 报文>   # GET
POST /dns-query  (Content-Type: application/dns-message)            # POST
```

## 17. 错误码

| 状态码 | 含义 | 常见场景 |
| --- | --- | --- |
| `400 Bad Request` | 请求体非法 / 参数错误 | JSON 解析失败、缺少必填参数、目标不是 Selector |
| `401 Unauthorized` | 未通过鉴权 | `secret` 缺失或错误 |
| `403 Forbidden` | 无权限 | 访问安全路径外的文件等 |
| `404 Not Found` | 资源不存在 | 代理/策略组/Provider 名称不存在 |
| `413 Payload Too Large` | 请求体过大 | storage 写入超过 1MB |
| `500 Internal Server Error` | 服务端错误 | DNS 未启用、升级失败、GEO 更新失败 |
| `503 Service Unavailable` | 服务暂不可用 | 延迟测速失败、Provider 更新失败 |
| `504 Gateway Timeout` | 超时 | 延迟测速超时 |

错误响应体统一为:

```json
{ "message": "错误描述" }
```

## 18. 常用操作速查

```bash
# 探活
curl http://127.0.0.1:9090/

# 版本
curl http://127.0.0.1:9090/version

# 切换模式为 global
curl -X PATCH http://127.0.0.1:9090/configs -d '{"mode":"global"}'

# 获取代理列表
curl http://127.0.0.1:9090/proxies

# 切换选择器选中节点
curl -X PUT http://127.0.0.1:9090/proxies/PROXY -d '{"name":"节点名"}'

# 测试节点延迟
curl "http://127.0.0.1:9090/proxies/节点名/delay?url=http://www.gstatic.com/generate_204&timeout=5000"

# 测试整组延迟
curl "http://127.0.0.1:9090/group/PROXY/delay?url=http://www.gstatic.com/generate_204&timeout=5000"

# 查看活跃连接
curl http://127.0.0.1:9090/connections

# 关闭全部连接
curl -X DELETE http://127.0.0.1:9090/connections

# 清空 DNS 缓存
curl -X POST http://127.0.0.1:9090/cache/dns/flush

# 刷新代理订阅
curl -X PUT http://127.0.0.1:9090/providers/proxies/订阅名

# 实时流量(SSE)
curl -N http://127.0.0.1:9090/traffic

# 实时流量(WebSocket)
wscat -c "ws://127.0.0.1:9090/traffic?token=<secret>"
```

带鉴权的示例(所有请求均可加 `-H "Authorization: Bearer <secret>"`):

```bash
curl -s http://127.0.0.1:9090/version \
  -H "Authorization: Bearer your-secret"
```

## 19. 与本地客户端的对应关系

本项目 `internal/client/mihomo/` 中的客户端封装与本文档接口的对应关系:

| 本地方法 | 接口 |
| --- | --- |
| `Ping()` | `GET /version` |
| `SwitchModes(mode)` | `PATCH /configs` `{"mode": mode}` |
| `GetModes()` | 硬编码 `["Rule", "Global", "Direct"]`(对应 `mode` 可选值) |

后续如需扩展(连接管理、节点测速、订阅刷新等),可直接参照上文对应接口在 `internal/client/mihomo/` 中补充实现。
