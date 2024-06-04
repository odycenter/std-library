# Std-Library
***
### 基础封装库，提供一些工具的封装和实现供上层调用
>#### version 1.0.1

-   ## [std-library](std-library)
    *  ### [array](array) 对数组的操作
    *  ### [auth](auth) 身份验证
		+  #### [google](auth%2Fgoogle) google码
    *  ### [cache](cache) 缓存及淘汰算法实现
		+  #### [lfu](cache%2Flfu) LFU算法
		+  #### [lru](cache%2Flru) LRU算法
    *  ### [compression](compression) 常用压缩算法
		+  #### [gzip](compression%2Fgzip) Gzip压缩算法
    *  ### [containers](containers)  线程安全容器的实现
        +  #### [maps](containers%2Fmaps) safe map
        +  #### [slice](containers%2Fslice) safe slice
    *  ### [crash](crash) 实现TryCache
    *  ### [crypto](crypto) 密码学实现和hash实现
		+  #### [aes](crypto%2Faes)
		+  #### [des](crypto%2Fdes)
		+  #### [ecdsa](crypto%2Fecdsa)
		+  #### [hmac](crypto%2Fhmac)
		+  #### [md5](crypto%2Fmd5)
		+  #### [rsa](crypto%2Frsa)
		+  #### [sha](crypto%2Fsha)
		+  #### [tripleDes](crypto%2FtripleDes)
    *  ### [dbase](dbase) 封装的基于beego的SQL类实现
    *  ### [docs](docs) 文档处理
		+  #### [csv](docs%2Fcsv) 
		+  #### [excel](docs%2Fexcel)
    *  ### [elastic](elastic) ES操作的封装和实现
    *  ### [email](email) EMail操作封装
    *  ### [file](file) 文件操作
    *  ### [fuzzy](fuzzy) 字符串模糊匹配相关操作
    *  ### [grpc](grpc) GRPC封装和实现
    *  ### [imgex](imgex) 图片处理
    * +  #### [resize](imgex%2Fresize)
    *  ### [json](json) json处理
    *  ### ~~[jwt](jwt) JWT token封装~~
    *  ### [kafka](kafka) Kafka的封装和实现
    *  ### [locker](locker) 分布式锁的实现(TODO 基于etcd的实现需要等etcd连接库更新grpc库) 
    *  ### [logs](logs) log封装和实现
    *  ### [mongo](mongo) mongo操作的封装
    *  ### [nacos](nacos) nacos操作的封装
    *  ### [nets](nets) 网络相关方法的封装
    *  ### [numerical](numerical) 数字类操作
    *  ### [operator](operator) 运算符实现
    *  ### [pagex](pagex) 分页操作(多服务于SQL结果的分页)
    *  ### [pyroscope](pyroscope) Pyroscope操作的封装
    *  ### [qbd](qbd) QueryBuilder的实现(只提供有限的SQL注入保护)
    *  ### [rand](rand) 一些常用随机值方法的实现
    *  ### [redis](redis) redis操作的封装 
    *  ### [sdk](sdk) 一些三方SDK的封装
		+  #### [aliyun](sdk%2Faliyun)
            -  ##### [oss](sdk%2Faliyun%2Foss)
		+  #### [aws](sdk%2Faws)
            -  ##### [s3](sdk%2Faws%2Fs3)
		+  #### [google](sdk%2Fgoogle)
            - ##### [ocr](sdk%2Fgoogle%2Focr)
		+  #### [jpush](sdk%2Fjpush)
        +  #### [slack](sdk%2Fslack)
		+  #### [tg](sdk%2Ftg)
    *  ### ~~[searcher](searcher) 一个基于内存的全文搜索引擎~~
    *  ### [stringx](stringx) 字符串操作扩展 
    *  ### [structx](structx) 结构体操作 
    *  ### [tg](tg) 简单的tg发文实现 
    *  ### [timex](timex) 时间操作扩展和封装
    *  ### [unique](unique) 唯一值生成
		+  #### [snowflake](unique%2Fsnowflake)
		+  #### [uuid](unique%2Fuuid)
    *  ### [valid](valid) struct赋值校验验证
    *  ### [websocket](websocket) websocket封装

# 代码文档
### 生成方式
```shell
#启动godocserver 
godoc -goroot ./ -http=:8080
#访问文档地址
http://127.0.0.1:8080/pkg/std-library/
```

# 代码编写及安全性规则

详见:[CODING_RULES_CN.md](CODING_RULES_CN.md)

# 版本待办事项

详见:[TODO.md](TODO.md)

# 已实现的修改

详见:[CHANGE_LOG.md](CHANGE_LOG.md)