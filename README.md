# DPanel 插件规范说明

插件是 DPanel 的功能扩展方式。它在系统中主要有两个作用：一是承载并非所有用户都需要的功能，让主程序保持精简；二是接入第三方应用，为其提供适合的可视化配置界面。

本规范约定了插件如何安装和运行、如何收集配置、如何接入页面，以及如何在权限范围内调用 DPanel API。你可以结合示例 `dp_example` 插件更好的了解插件的定义及使用。

## 目录组织

插件以目录的形式存在，一个目录表示一个独立完整的插件。目录名必须与 `plugin.name` 一致，格式为 `前缀_插件名`。前缀标识发布者，两部分使用小写字母、数字和连字符。
你也可以通过 Zip 或是 Git 仓库的形式组织多个插件，每个插件占用仓库根目录下的一个一级目录。

```text
plugin-repository/
├── dp_rclone/
│   ├── manifest.yaml
│   └── docker-compose.yaml
└── dp_example/
    ├── manifest.yaml
    ├── docker-compose.yaml
    ├── asset/
    └── server/
```

每一个插件中必须包含 `manifest.yaml` 和 `docker-compose.yaml` 两个文件。其中 `manifest.yaml` 是基本的插件的配置描述文件，`docker-compose.yaml` 则是最终插件运行的容器启动编排文件。
插件中可以根据自己的需求决定是否提供前端文件以及后端源码。前端文件可以通过 `manifest.yaml` 配置在系统中通过微前端展现出来，后端由需要自行完善 `Dockerfile` 在 `docker-compose.yaml` 中编译生成镜像。

## 配置

## manifest.yaml

```yaml
version: "1.0"
plugin:
  name: dp_rclone
  version: "1.0.0"
  defaultLang: zh
  title:
    zh: Rclone 同步
    en: Rclone Sync
  bootstrap: compose
  
view:
  type: wujie
  entry: "http://{{DP_PLUGIN_HOST_IP}}:{{SERVER_PORT}}/"

config:
  - name: SYNC_PATH
    label:
      zh: 本地同步目录
      en: Local directory
    type: text
    required: true
```
