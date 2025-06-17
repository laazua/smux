### SMUX

- 说明
1. go 标准库net套接字编程
2. 只做了简单的粘包处理
3. 添加了tls双向认证
4. 修改certfile.sh脚本中[ alt_names ]块设置服务端地址为服务的端地址
5. 如果要动态扩容server节点, 则需要保证所有server节点和client节点基于相同的ca证书来进行证书签发