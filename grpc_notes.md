**Protobuf**
* 可以将数据序列化为二进制编码，这会大幅减少需要传输的数据量，从而大幅提高性能
* 通过http2实现异步的请求
* protobuf 代替 json （用protofiles创建gRPC服务，用protocol buffers消息类型来定义方法参数和返回类型）
* 最新版本 proto3 (当前默认proto2, 但是proto3可以支持所有的语言，并且能避免 proto2 客户端与 proto3 服务端交互时出现的兼容性问题）
* 
