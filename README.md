# fileFinder (文件搜索器)
![](https://img.shields.io/badge/author-TheSevenSky-blue) ![](https://img.shields.io/badge/build-passing-yellow) ![](https://img.shields.io/badge/Release-Development-red)
![](https://camo.githubusercontent.com/8ab5e05ff609c4a280640cef9c5beeb9bc1953881e9daba2d6235b5989381557/68747470733a2f2f676f646f632e6f72672f6769746875622e636f6d2f6a657373656475666669656c642f6c617a796769743f7374617475732e737667)

一个非常快速好用的基于go语言编写的文件搜索工具
<hr/>

## 快速开始

```git
git clone https://github.com/sta-golang/fileFinder.git
cd fileFinder
go build
./fileFinder vi /user/bin
```
上面命令演示了拉取下来并且完成编译后运行
<br/>
最后一行指令是从/user/bin 目录下查找所有的名字里带vi的文件

<br/>
当然上面得到的结果是忽略大小写的,如果需要标准匹配大小写就需要第三个参数
<br/>

```shell
./fileFinder vi /user/bin false
```


## 作者信息
毕业于XUPT 软件科技协会实验室成员


