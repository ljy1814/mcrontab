### 项目介绍

- etcd
- 调度器
    三个工作.任务事件,空闲休眠,任务执行完毕
    job事件包括:添加,删除,中止任务
    计算本次等待的时间,同时启动当前时间要启动的任务
- http服务器
- 任务管理器
- 监控器
    对任务进行监控,包括正常的任务以及强杀的任务
    将etcd操作传给scheduler
- 执行器
    执行shell命令,返回结果