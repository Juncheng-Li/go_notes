**架构拓补**
web管理系统
服务节点（node服务节点，N个业务节点）- node节点：统一管理，提供起停、发布、监控等功能，同时接受业务服务节点上报过来的心跳
registry
Patch发布管理, config配置中心, log, stat, property业务属性-内存，队列大小, notify
