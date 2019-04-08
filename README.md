# controlVM-go

# 关键字：
控制；虚拟机；go；VBoxManage；VirtualBox

# 物理环境
* 物理机：MAC OS X
* virtualbox：版本 5.2.16
* 虚拟机：centos 7（四台全是）

# 需求：
1. 实现物理机控制虚拟机，包括开/关虚拟机
2. 在虚拟机无法联网、无法被 ssh 访问时执行虚拟机脚本或 shell 命令，也藉由执行脚本或者命令通过 sed 、awk 等修改文件

# 已实现：
1. 根据配置文件，实现虚拟机开/关
2. 执行虚拟机 shell 命令，只要是在虚拟机能执行的命令，均可以支持

# 开发计划：
* 模块化：将已实现的功能模块化为第三方模块
* BuildCommand：代码部分不够优雅
* 剔除依赖：剔除第三方模块 go-virtualbox 的依赖
* 支持虚拟机脚本直接运行，当前支持间接运行。如运行脚本 /root/run.sh ，在虚拟机中键入 /root/run.sh 回车即可，当前配置文件需写作 bash /root/run.sh  
* 自动化：提供更加自动化的操作方式，修改配置文件决定对虚拟机操作，始终是不够优雅的方式。最好是能够支持物理机命令行 run 代码，命令行中加入各种参数（包括配置文件）

# 帮助文档：
遇到类似错误或者依赖问题： 0x80bb0005 (The guest execution service is not ready (yet))  请访问： https://blog.csdn.net/github_37320188/article/details/89066634
