
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

## tips
```
Issuer 和 Subject 的区别
	•	Issuer：表示证书的颁发者，即签署该证书的 CA。
	•	Subject：表示证书的持有者，即证书所代表的实体。

  
```