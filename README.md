# https-proxy

### build

```
sh build.sh

# for go 1.9 version build
docker-compose run app sh -x build.sh
```

### CHANGELOG

##### 1.2.0 / 2019-09-06

[CHANGE] 
- 修改程序中client的复用和管理
- 修改原有log到logrus库，并设置对应debug信息输出
[ENHANCEMENT] 
[BUGFIX]


##### 2.0.0 /2020-07-01
[CHANGE]
- 使用krakend来实现路由管理
