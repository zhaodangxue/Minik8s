# Go Build Tag

go语言中选择文件进行编译的机制

在go文件的开头（package之前）添加一行注释，形如

```
//go:build ${标签名称}
```

标签名称是例如dev, release, mambaout，这样的字符串

在编译时，使用

```
go build --tags mambaout target.go
```

就可以在编译时对文件进行筛选，具体而言，如果文件没定义tag，那么必定被选中加入编译，定义了tags，若匹配则加入编译，不匹配则不加入编译
