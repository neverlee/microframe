# MicroFrame

MicroFrame是基于[go-micro](https://github.com/micro/go-micro)结合go的plugin（需要go1.8或以上）开发的一个微服务框架。

[go-micro](https://github.com/micro/go-micro)本身已经做了非常好的抽象和插件化。MicroFrame没有直接采用go-micro，而是在它的基础上重新开发有下面一些原因。
1. 对纯粹的业务开发屏蔽掉掉所有服务本身的代码（当时也有考虑过使用生成代码）
2. go-micro使用命令行参数设置启动状态，MicroFrame则依据配置文件
3. 统一微服务节点结构
4. 统计提供日志收集，状态监控，PV统计，性能统计，请求跟踪调试的服务。
5. 所有节点可以通过配置，直接使用以上这些公共服务

## 微服务架构设计
[MicroFrame](http://github.com/neverlee/microframe/raw/master/docs/MicroFrame.png)

### 节点设计
1. 每个微服务节点为一个MicroFrame框架服务实例（包括API、WEB、SRV、Worker、数据库代理、配置中心节点、日志收集节点）
2. MicroFrame通过配置选择加载需要公共插件和业务插件，而成为不同的功能节点
3. 对于一个新的业务功能，我们只需要为其开发相应的微服务Handle业务插件即可

### MicroFrame框架运行流程
1. 框架读取配置文件
2. 框架根据配置文件加载插件（当内部插件与外部插件同名时，会加载外部插件）
3. 调用已经加载插件的NewPlugin函数新建插件结构体（插件的NewPlugin调用无序）
4. 对插件按插件本身阶段类型以及在配置文件中的位置进行排序
5. 顺序依次调用插件的Preinit函数（wrap操作需要在这一步增加）
6. 框架初始化自己的Service
7. 顺序依次调用插件的Init函数（这一步注册业务Handle）
8. 顺序依次调用插件的Start函数
9. 逆序依次调用插件的Stop函数
10. 逆序依次调用插件的Uninit函数

## 微服务技术选型示例
* 服务发现 consul（优先考虑）
* 服务保护 hystrix
* 服务通信 protobuf（优先）, jsonrpc（兼容）
* 框架 [microframe（基于go-micro）](https://github.com/neverlee/microframe)
* 部署方案 kubernetes，docker swarm
* 监控报警 influxdb, telegraf, hrafana
* 日志收集 fluentd（ELK）
* 配置文件格式 yaml（优先），toml
* 日志格式 [xclog](https://github.com/neverlee/xclog/go)

