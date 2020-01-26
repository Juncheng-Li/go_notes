# hadoop的学习笔记

## 简介
* 分工叫做map，合并分工叫做reduce - 合起来叫map reduce
* master 和 slave node，master就是name node， slave就是data node
  
### HDFS
* 全称：hadoop distributed file system
* Blocks：默认128MB大小，如果一个block承载了只有5MB的内容，则block占用的空间只有5MB
* 每个分开的block都存在不同的HDFS上
* Name node：HDFS有master-worker pattern，知道所有file的status和metadata。Metadata包括file permission，name和location
  
#### fault tolerance
* replication factor(默认是3，也就是每个block默认复制三份在不同的data nodes上)

#### write mechanism
* 先问一个data node，在按顺序一个一个传过去问
* 返回ack
* 得到并返回这条pipeline
* 根据这条pipeline按顺序写数据（1次）和复制数据（n次）
* 再根据这个pipeline反向传递写入完成的ack
* 但是一个文件的分开的不同的block是同时被写入的（异步的）

#### read mechanism
* client先向name node获取想要的那些block被存放的地址
* client得到地址之后访问core switch然后根据地址得到这些block
* 和write一样，读取的时候每个block读取都是平行的（异步的）

### MapReduce
* 主要成分：Mapper Code，Reducer Code，Driver Code
* Mapper
    1. key是text
    2. value是intWritable
* Reducer
    1. input: Bear, [1, 1], 等等
    2. output: Bear, 2; Car, 3; 等等

## Hadoop
* 一些简单的命令
    1. hadoop fs -cat
    2. hadoop fs -mkdir
    3. hadoop fs -ls
    4. hadoop jar path.jar inputDatasetPath(dir?) outputPath(dir?) 

## Yarn（相当于MapReduce V2）
* App master负责收集所有运行需要的资源，并向Resource manager索取，最后调度container运行任务

## 安装
* core-site.xml：名字和addr
* hdfs-site.xml：replication factor，path（dataNode）， path（nameNode）
* yarn-site.xml: 
* mapred-site.xml: 复制mapered-site.xml.template来创建mapered-site.xml
* bashrc:添加HADOOP_HOME, HADOOP_CONF_DIR, HADOOP_MAPRED_HOME, HADOOP_COMMON_HOME, HADOOP_HDFS_HOME, YARN_HOME, HADOOP_COMMON_LIB_NATIVE_DIR, HADOOP_OPTS
* hadoop-2.10.0/etc/haddop/hadoop-env.sh: 修改JAVA_HOMEJAVA_HOME
* 开启：bin/hadoop namenode -format
* 开启：sbin/start-dfs.sh
* 开启：sbin/start-yarn.sh
* 开启: sbin/mr-jobhistory-deamon.sh start historyserver
* 监控：浏览器：localhost:50070
* 