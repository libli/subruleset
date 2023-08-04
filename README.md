# subconverter 支持 no-resolve

准备配置文件：
```bash
touch config.yaml
```

内容如下：
```yaml
token: "自己随便设置"
urls:
  naxi: http://127.0.0.1:25500/sub?target=clash&expand=false&url=[机场url]&config=[配置文件url]
```

其中key: naxi 自定义，为最终url的path部分。建议和你机场名一致。

value: subconverter转换的url，其中target=clash表示转换为clash配置，expand=false表示不展开规则，url=[机场url]表示机场订阅地址，config=[配置文件url]表示配置文件地址。

启动：
```bash
docker run --name=subruleset -d --restart=unless-stopped \
  -v /data/subruleset/config.yaml:/app/config.yaml \
  -p 8080:8080 \
  libli/subconverter:latest
```
