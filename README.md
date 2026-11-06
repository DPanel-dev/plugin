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

## 镜像元数据协议

主程序、应用商店源和插件源共用 `com.dpanel.*` 字段命名。Annotation 只包含展示和检测更新所需的版本信息；镜像 config 的 label 保存这些字段的同名同值副本及完整元数据。不另设 `namespace` 字段，也不使用 `org.opencontainers.*` 作为本协议字段。

| 字段 | Annotation | Label | 含义 |
| --- | --- | --- | --- |
| `com.dpanel.version` | 必填 | 必填 | 展示版本，不参与更新比较 |
| `com.dpanel.release` | 必填 | 必填 | 发布修订标识，用于判断内容是否变化 |
| `com.dpanel.image.createdAt` | 不写入 | 必填 | 本次镜像打包时间，使用 UTC RFC 3339，例如 `2026-11-06T08:30:00Z` |
| `com.dpanel.store.root` | 不写入 | 商店源必填 | 镜像内的文件提取目录，解压时剥离此前缀 |
| `com.dpanel.store.source` | 不写入 | 商店源必填 | 内容来源仓库地址，不包含凭据、分支片段、末尾斜线或 `.git` 后缀 |
| `com.dpanel.store.commit` | 不写入 | Git 商店源必填 | 实际打包的完整 commit ID，不使用短 commit 或打包工具仓库的 commit |

插件源属于商店源，使用相同的 `store.*` 字段；主程序镜像不要求填写商店字段。版本描述整个发布包，不取其中某个应用或插件的版本，也不把构建工具或管理面板自身的版本当作所构建内容的版本。

### 版本与更新判断

- `version` 是非空展示字符串，发布标签可去掉开头的 `v`，没有发布版本的滚动仓库使用 `rolling`。它不参与排序或更新判断。
- `release` 是非空修订字符串。Git 仓库快照默认使用实际打包的完整 commit；主程序可以使用发布批次标识。`store.commit` 只记录源码提交，即使默认与 `release` 相同，也不替代它的更新比较职责。
- 只在同一个产品或源包、发布通道、变体和目标平台内比较 `release`。相同且本地已有成功同步的可用内容时，可以跳过文件层下载；不同表示目标发布有变化，不按字符串大小推断新旧。仅 `version`、打包时间或镜像 digest 变化不触发内容更新。
- 同一比较范围内，一个 `release` 必须对应相同的有效内容。目录选择、打包规则、基础镜像或依赖变化导致有效内容变化时，必须改变 `release`，即使源码 commit 没有变化。
- 远端查询失败、版本字段缺失或同名字段冲突时，结果为无法确认或无效产物，不能表示已是最新。本地内容缺失或上次同步失败时，不能仅凭相同 `release` 跳过同步。
- 检查后按取得的 manifest digest 锁定下载对象，避免可变 tag 在检查与下载之间指向另一版。只有内容下载、校验及正式发布成功后，才能记录本次成功同步的修订；失败保留上次成功状态。

### 元数据存放位置

两个版本 annotation 写入 image manifest 的 `annotations`。存在多平台 image index 时，在各目标平台的 manifest descriptor 上同步写入对应值，不能只写 index 顶层。Label 写入 image config 的 `config.Labels`；相同字段由同一份构建元数据生成，值必须一致。

检测版本只需要读取 index 或 manifest。读取完整 label 时，再按 manifest 中的 config digest 获取 config JSON，无须下载镜像文件层。跨 registry 分发应复制同一产物并保留 digest 和元数据，不在每个 registry 分别重建。

### 商店源目录

`store.root` 使用规范化的 POSIX 绝对目录路径，`/` 表示镜像根目录。消费者只提取该目录下的内容，并剥离整个目录前缀；不自动猜测第一层目录。例如：

```text
com.dpanel.store.root=/payload
镜像内：/payload/apps/demo/docker-compose.yml
解压后：<目标目录>/apps/demo/docker-compose.yml
```

root 必须存在且为实际目录，不能通过符号链接跳转；拒绝包含 `..`、反斜杠等非法路径及越出提取边界的归档条目或链接。root 外的内容不写入目标目录。该字段表示最终镜像内的位置，不是宿主机路径；若构建时选择 Git 仓库子目录，应以最终镜像布局填写。

### 接入与兼容

以上是发布者与消费者的协议约定，不表示所有消费端已经实现版本检测和 root 剥离。旧插件消费者仍依赖 `com.dpanel.plugin.source` 与 `com.dpanel.plugin.root` 时，发布者暂时保留与 `store.source/root` 同值的旧 label；这些字段仅用于迁移兼容，待消费端完成切换并验证后移除，不放入 annotation。旧消费者只支持 `root=/` 时，不能宣称其支持其他目录。

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

### 插件描述

`plugin.description` 按语言填写 Markdown 正文，使用 YAML 的 `|` 保留换行：

```yaml
plugin:
  description:
    zh: |
      ## 插件用途

      - **功能**：说明插件提供的能力。
      - **使用方式**：说明运行前的配置和使用步骤。
```

也可以使用 `plugin.descriptionFile` 按语言引用 Markdown 文件。路径相对于该插件的 `manifest.yaml` 所在目录：

```yaml
plugin:
  descriptionFile:
    zh: README.md
  description:
    en: |
      ## About this plugin

      Describe its features and usage here.
```

两个字段都可省略，也可以分别为不同语言使用正文和文件；同一语言不能同时出现在两个字段中。字段值中的普通字符串不会自动被识别为文件路径。

引用的文件须为插件目录内的 UTF-8 `.md` 普通文件，大小不超过 1 MiB；不支持远程 URL、绝对路径、越界路径或文件符号链接。文件须随插件仓库一起交付。引用缺失、无效或同一语言重复定义会使本次来源同步失败，保留之前成功同步的索引。

描述在来源同步时读取并保存，修改后需要重新同步插件源。详情页按当前界面语言显示，未提供对应语言时回退到 `plugin.defaultLang` 的描述；两者均未提供时不显示描述。两种来源均按 Markdown 渲染并清理不安全 HTML。文件中的相对图片和附件链接不会自动映射到插件目录，需使用可访问的完整 URL。
