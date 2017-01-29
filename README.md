# win32
go win32 API + generator

## build

* go
* gcc (如果需要 cgo 实现回调)

对于自定义的 module
```
% cd <module>
% go generate -tags genapi
% go build -i
```

## API 描述
在 `/*` 和 `*/` 编写 **module**.dll 下 API 的 go 函数声明
