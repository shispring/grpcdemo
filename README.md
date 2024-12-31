
# grpcdemo

## CA 证书制作：
### 1、生成.key  私钥文件
```
openssl genrsa -out ca.key 4096
```

### 2、生成证书
```
#生成.csr 证书签名请求文件
openssl req -new -key ca.key -out ca.csr  -subj "/C=CN/ST=Beijing/L=Beijing/O=demo/OU=IT Department/CN=www.demo.cn"

#自签名生成.crt 证书文件
openssl req -new -x509 -days 3650 -key ca.key -out ca.crt -config /opt/homebrew/etc/openssl@3/openssl.cnf \
  -extensions v3_ca \
  -subj "/C=CN/ST=Beijing/L=Beijing/O=demo/OU=IT Department/CN=www.demo.cn"
```

或者合成一步
```
openssl req -new -x509 -days 3650 -key ca.key -out ca.crt \
  -subj "/C=CN/ST=Beijing/L=Beijing/O=demo/OU=IT Department/CN=www.demo.cn" \
  -config /opt/homebrew/etc/openssl@3/openssl.cnf \
  -extensions v3_ca
```

---

## 服务端证书

### 1、生成 .key 私钥文件
```
openssl genrsa -out server.key 2048
```

### 2、生成 .csr 证书签名请求文件
```
#查看openssl.conf
openssl version -d

#拷贝一份新的配置
cp /opt/homebrew/etc/openssl@3/openssl.cnf server_openssl.cnf

#添加SAN配置
echo "[SAN]" >> server_openssl.cnf
echo "subjectAltName=DNS:*.demo.cn,DNS:www.demo.cn" >> server_openssl.cnf

#生成 .csr
openssl req -new -key server.key -out server.csr -config server_openssl.cnf -reqexts SAN \
  -subj "/C=CN/ST=Beijing/L=Beijing/O=demo/OU=IT Department/CN=www.demo.cn"

```

### 3、签名生成.crt 证书文件
```
openssl x509 -req -days 3650 -in server.csr -out server.crt \
  -CA ca.crt -CAkey ca.key -CAcreateserial \
  -extensions SAN -extfile server_openssl.cnf
```

---

## 查看证书内容
```
openssl x509 -in ./ca.crt -text -noout
openssl x509 -in ./server.crt -text -noout
```


## 验证
```
go test -v ./test/server_test.go 
go test -v ./test/client_test.go 
```

## 使用grpcurl调试
```
#使用了TLS模式调试步骤
grpcurl -v \
  -cacert ca.crt \
  -authority www.demo.cn \
  127.0.0.1:8972 list

grpcurl -v \
  -cacert ca.crt \
  -authority www.demo.cn \
  127.0.0.1:8972 list helloworld.Greeter

grpcurl -v \
  -cacert ca.crt \
  -authority www.demo.cn \
  127.0.0.1:8972 describe helloworld.Greeter.SayHello


grpcurl -v \
  -cacert ca.crt \
  -authority www.demo.cn \
  127.0.0.1:8972 describe helloworld.HelloRequest



grpcurl \
  -cert server.crt \
  -key server.key \
  -cacert ca.crt \
  -authority www.demo.cn \
  -d '{"name": "China"}' \
  127.0.0.1:8972 helloworld.Greeter/SayHello

或者

grpcurl \
  -cert server.crt \
  -key server.key \
  -cacert ca.crt \
  -authority www.demo.cn \
  -d '{"name": "China"}' \
  127.0.0.1:8972 helloworld.Greeter.SayHello


1.-cert server.crt
指定客户端使用的证书文件，这里是服务端证书文件 server.crt。
2.-key server.key
指定客户端使用的私钥文件，与 server.crt 配套使用。
3.-cacert ca.crt
指定 CA 根证书文件，用于验证服务端证书的合法性。
4.-authority www.demo.cn
指定服务端证书的 Common Name 或 subjectAltName，这里需要与服务端证书中的 CN 或 SAN 匹配（例如 www.demo.cn）。
5.-d '{"name": "China"}'
指定请求数据，这里是一个 JSON 格式的请求体，调用服务端的 SayHello 方法时需要传递此数据。 
```

---

### **关于 `-subj` 的使用与含义**

`-subj` 参数用于在命令行中直接指定证书的主题信息，避免交互式输入。这是非交互式生成证书时非常常用的方式。

#### **`-subj` 参数的格式**
`-subj` 参数的值是一个以 `/` 分隔的字符串，每个字段的含义如下：

| 字段名 | 全称                   | 含义                                                                                      | 示例               |
|--------|------------------------|-------------------------------------------------------------------------------------------|--------------------|
| `C`    | Country Name           | 国家代码，需为两位 ISO 3166 国家代码                                                      | `C=CN`            |
| `ST`   | State or Province Name | 州或省的名称（可选）                                                                       | `ST=Beijing`      |
| `L`    | Locality Name          | 地区或城市名称                                                                            | `L=Beijing`       |
| `O`    | Organization Name      | 组织名称（通常是公司名称）                                                                 | `O=demo`           |
| `OU`   | Organizational Unit    | 组织部门名称（可选，通常是团队或部门名称）                                                 | `OU=IT Department`|
| `CN`   | Common Name            | 通用名称，通常是域名或主机名                                                               | `CN=www.demo.cn`   |
| `emailAddress` | Email Address | 邮箱地址（可选）                                                                          | `emailAddress=admin@demo.cn` |

#### **示例**
```bash
-subj "/C=CN/ST=Beijing/L=Beijing/O=demo/OU=IT Department/CN=www.demo.cn"
```

---

### **注意事项**
1. **字段的匹配**：
   - `CN`（Common Name）字段必须与您希望匹配的域名一致。例如，`CN=www.demo.cn` 表示该证书适用于 `www.demo.cn`。
   - 如果需要支持多个域名，请在扩展字段（如 `subjectAltName`）中添加更多域名，而不是只依赖 `CN`。

2. **格式正确性**：
   - 每个字段必须以 `/` 开头，不能有多余的空格。
   - 如果某些字段为空，可以省略，但必须确保其他字段正确填写。

3. **兼容性问题**：
   - `CN` 字段在现代证书中已逐渐被 `subjectAltName` 替代，因此推荐在配置文件中明确添加 `SAN` 扩展。

---

### **Issuer 的含义**

在读取证书时，`Issuer` 表示签发该证书的实体信息，即证书的颁发者（CA）。以下是常见的字段及其含义：

| 字段名 | 含义                                   | 示例                |
|--------|----------------------------------------|---------------------|
| `C`    | 签发者所在国家代码                     | `C=CN`             |
| `ST`   | 签发者所在州或省份                     | `ST=Beijing`       |
| `L`    | 签发者所在城市                         | `L=Beijing`        |
| `O`    | 签发者组织名称                         | `O=demo Root CA`    |
| `OU`   | 签发者组织部门                         | `OU=Certificate Authority` |
| `CN`   | 签发者的通用名称（通常是 CA 的名称）    | `CN=demo Root CA`   |

---

#### **Issuer 和 Subject 的区别**
- **`Issuer`**：表示证书的颁发者，即签署该证书的 CA。
- **`Subject`**：表示证书的持有者，即证书所代表的实体。

---

### **如何查看证书的 Issuer 信息**
使用以下命令查看证书的详细信息：
```bash
openssl x509 -in server.crt -text -noout
```

---

### **总结**
1. `-subj` 提供了一种简洁、非交互式的方式指定证书的主题信息。
2. 证书的 `Issuer` 表示签发证书的 CA 信息，通常需要与信任链中的根证书匹配。
3. 在现代证书中，`subjectAltName` 的配置非常重要，`CN` 仅作为备用字段。