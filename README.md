# 说明

我自己是使用TTRSS，同时配合`mercury-fulltext`插件和`mercury-parser-api`服务进行全文抓取。

在看帖子[2023 年，我为什么选择 Miniflux 作为 RSS 主力工具 - V2EX](https://www.v2ex.com/t/963837?r=AboutRSS)的时候，OP提到miniflux的全文抓取是他认为最好的。

尝试了下miniflux，认为它不满足自己的要求。它的图片代理功能(PROXY_OPTION=all)是全局的，无法针对某些Feed单独设置，而我又有很多源是国内的，通过服务代理一方面减慢了速度，另一方面也浪费了带宽。

因此突发奇想，要不干脆把miniflux的全文抓取逻辑代码抠出来，伪装成`mercuyr-parser-api`的接口，这样我在TTRSS中不就可以用了。

该仓库是在[Release Miniflux 2.0.47 · miniflux/v2](https://github.com/miniflux/v2/releases/tag/2.0.47)的基础上，把全文抓取的相关的代码做了一部分精简后整理出来的，接口仿`mercury-parser-api`。

个人使用了后，发现它对一些微信公众号的文章抓取不到位，在一些有图片的情况下无法抓取到，而`mercury-parser-api`是没有问题的。因此最终还是换回了`mercury-parser-api`。

# 使用

编译后，启动，默认监听在`3002`端口。

```bash
go build .
```

请求。

```bash
curl http://127.0.0.1:3002/parser?url=www.baidu.com
```

响应。

```bash
{
    "content": "" // 内容
}
```