# cuttle 基于fabric-ca的证书颁发

### 项目介绍

```
在开发环境中用cryptogen工具来生成各组织成员的私钥和证书。本项目基于fabric-ca，实现证书的颁发

为了兼容cryptogen（即一键颁发联盟所有证书），本项目支持两种证书颁发方式，均通过配置文件配置实现
```
### 项目依赖


```
go get -u golang.org/x/crypto/sha3
go get -u gopkg.in/yaml.v2
```

### 安装


```
go get -u learnergo/cuttle
```


### 颁发方式
```
- 一键颁发（即cryptogen所实现），需要配置static\crypto-config.yaml文件

- 颁发特定文件，需要配置static\cuttle.yaml文件
```
### 配置介绍

```
static\crypto-config.yaml 文件仿照fabric中crypto-config.yaml文件，但不同的在于每个组织需要制定各自根ca的配置文件，并且在Subject中定义通用Subject属性

static\cuttle.yaml 文件则用于颁发特定证书，配置register和enroll各个细节

static\ca.yaml 文件配置ca信息，包括url和管理员证书
```

### 运行方式


```
main函数中
一键颁发：RunConfig
特定颁发：RunSpeConfig
```

### TODO


```
现在是基于一个根ca，颁发所有相关证书，后续会配置多个ca
```

