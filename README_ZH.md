切换新的git库本地需要做的事情：
1. 修改本用户下的 profile文件
  vi ~/.profile
  在文件最后增加如下两行：
  export GO111MODULE=on
  export GOPROXY=https://goproxy.cn,direct
  然后执行配置保存后，执行source ~/.profile，使修改配置即时生效
2.代码工程不要放到gopath对应的路径下