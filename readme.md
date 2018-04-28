# vue-with-go

「 `go`的server服务, `vue`的前端 」

[![explain](http://llever.com/explain.svg)](https://github.com/chinanf-boy/Source-Explain)

---

## go后端与vue前端例子

### use-使用

- 1. `git clone `

```
go run main.go serve
```

- 2. `go get`

```
go get -u -v github.com/chinanf-boy/vue-with-go-example 
```

``` bash
vue-with-go-example good-serve
//or
vue-with-go-example bad-serve
```


<details>

<summary><h3>有关 good 和 bad serve 有什么区别</h3></summary>

```
go get -u -v github.com/chinanf-boy/vue-with-go-example 
```

一般来说 `go get`之后使用 `serve` 是不行的

>为什么呢❓

我们其实是通过 运行 一个完整的二进制程序 `GOPATH/bin/vue-with-go-example` => `二进制`

当开启服务器, 请求 `clinet/index.html` 怎么会有呢. 不存在的嘛, 这就造成了, 你只能在对应项目下运行

``` bash
vue-with-go-example bad-serve 
```

---

但是, 当然, 有解决办法, 就是把 前端所有的文件 进行 `二进制 `转化处理 _main.go_ [//go:generate fileb0x filebox.json
](https://github.com/UnnoTed/fileb0x)

``` bash
vue-with-go-example good-serve
```

---

filebox.json

``` json
{
    "pkg": "assets",
    "dest": "./assets/",
    "fmt": true,
    "compression": {
        "compress": true,
        "method": ""
    },
    "output": "assets.go", // 通过 assets.ReadFile(path) 意义上路径, 复制 缓存好的 html,js 二进制 给 服务器
    "custom": [{
        "files": ["./client/"],
        "base": "client/",
        "exclude": [
            "less/*.less"
        ]
    }]
}
```

</details>
  
---

例子主要来自 [书签🔖管理-go](https://github.com/RadhiFadlillah/shiori)

> 不可否认, 开启 `shiori serve` 后 对书签🔖的更改, 是如此的容易程序崩溃, 但同时, 也让我们看到 go 与 vue 的结合
