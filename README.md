# DPanel 插件开发规范

本规范定义了如何通过 manifest.yaml 和 docker-compose.yaml 文件构建 DPanel 面板插件。

---

## 文件结构

```
plugin-name/
├── manifest.yaml           # 插件配置清单（必需）
├── docker-compose.yaml     # 容器编排配置（必需）
└── README.md              # 插件说明文档（可选）
```

---

## 一、manifest.yaml 配置规范

### 1.1 基本结构

| 字段 | 类型 | 必需 | 说明 |
|------|------|------|------|
| `version` | string | ✅ | 规范版本，固定为 "1.0" |
| `plugin` | object | ✅ | 插件元数据 |
| `view` | object | ❌ | 视图配置 |
| `config` | array | ❌ | 配置项列表 |

### 1.2 plugin 插件元数据

| 字段 | 类型 | 必需 | 说明 |
|------|------|------|------|
| `name` | string | ✅ | 插件唯一标识符 |
| `version` | string | ✅ | 版本号 |
| `defaultLang` | string | ✅ | 默认语言，如 zh-CN、en-US |
| `title` | object | ✅ | 多语言标题 |
| `bootstrap` | string | ✅ | 启动方式：compose 或 shell |

### 1.3 view 视图配置

| 字段 | 类型 | 必需 | 说明 |
|------|------|------|------|
| `type` | string | ✅ | 前端模式：wujie、iframe、none |
| `entry` | string | ✅ | 访问地址，支持系统变量替换 |
| `props` | object | ❌ | 传递给前端的属性 |

**系统变量**：
- `{{DP_PLUGIN_NAME}}` - 插件容器名称
- `{{DP_PLUGIN_HOST_IP}}` - DPanel 主机 IP
- `{{SERVER_PORT}}` - 插件服务端口

### 1.4 config 配置项

#### 1.4.1 基本字段属性

| 字段 | 类型 | 必需 | 说明 |
|------|------|------|------|
| `name` | string | ✅ | 环境变量名/字段标识符 |
| `label` | string/object | ✅ | 显示标签 |
| `type` | string | ✅ | 字段类型 |
| `description` | string | ❌ | 配置项说明 |
| `required` | boolean | ❌ | 是否必填 |
| `default` | any | ❌ | 默认值 |
| `placeholder` | string | ❌ | 占位符 |

#### 1.4.2 国际化规范

**label 字段**：
- 简单字符串：`label: "配置项名称"`
- 多语言对象：`label: {zh-CN: "中文名", en-US: "English Name"}`

**description 字段**：
- 使用简单字符串（默认为中文）

**select 的 values**：
- label 使用简单字符串

#### 1.4.3 字段类型

| 类型 | 说明 |
|------|------|
| `text` | 单行文本输入 |
| `textarea` | 多行文本输入 |
| `number` | 数字输入 |
| `select` | 单选下拉框 |
| `group` | 字段分组 |

#### 1.4.4 select 类型

**额外字段**：
- `values` - 选项数组，每项包含 `label` 和 `value`
- `multiple` - 是否多选，值为 "true" 或 "false"

#### 1.4.5 group 类型

**说明**：group 是一种特殊类型，通过 `fields` 数组引用其他字段的名称。被 group 引用的字段不会单独显示，而是作为 group 的一部分。未被任何 group 引用的字段在根表单显示。

**字段属性**：
| 字段 | 类型 | 必需 | 说明 |
|------|------|------|------|
| `name` | string | ✅ | 分组唯一标识符 |
| `label` | string/object | ✅ | 分组标题 |
| `type` | string | ✅ | 固定为 "group" |
| `multiple` | boolean | ❌ | 是否支持多例，默认 false |
| `min` | number | ❌ | 最少实例数（multiple=true 时） |
| `max` | number | ❌ | 最多实例数（multiple=true 时，0=无限制） |
| `default` | number | ❌ | 默认实例数（multiple=true 时） |
| `collapsed` | boolean | ❌ | 是否默认折叠（multiple=false 时） |
| `fields` | array | ✅ | 引用的字段名称数组 |
| `description` | string | ❌ | 分组说明 |

**显示规则**：
- 被 group 引用的字段：在 group 内部显示
- 未被任何 group 引用的字段：在根表单显示
- group（multiple=false）：显示为折叠面板
- group（multiple=true）：显示为可添加/删除的列表

---

## 二、docker-compose.yaml 模板规范

### 2.1 模板语法

使用 Go Template 语法，支持变量替换、条件判断、循环。

#### 2.1.1 变量引用

**用户配置字段**：`{{.字段名称}}`

**系统变量**：`${DP_PLUGIN_NAME}`、`${DP_PLUGIN_HOST_IP}`、`${SERVER_PORT}`

#### 2.1.2 条件判断

```yaml
#{{- if .条件 }}
内容
#{{- end }}
```

#### 2.1.3 循环

```yaml
#{{- range $key, $val := .列表 }}
内容
#{{- end }}
```

#### 2.1.4 注释

使用 `#{{` 和 `}}` 作为模板注释，不会出现在最终文件中。

### 2.2 常用模板函数

| 函数 | 说明 |
|------|------|
| `split` | 分割字符串 |
| `splitList` | 分割为列表 |
| `index` | 获取数组元素 |
| `trim` | 去除空格 |
| `replace` | 字符串替换 |
| `eq` | 相等比较 |
| `ne` | 不等比较 |

### 2.3 访问规则

**根表单字段**：直接通过字段名访问
```yaml
{{.FIELD_NAME}}
```

**single group（multiple=false）的字段**：直接通过字段名访问
```yaml
{{.FIELD_NAME}}
```

**multiple group（multiple=true）的字段**：通过 range 循环访问
```yaml
#{{- range $index, $item := .GROUP_NAME }}
{{ $item.FIELD_NAME }}
#{{- end }}
```

### 2.4 扩展配置

**x-dpanel-service**：DPanel 特定配置

**image_tar - 本地镜像导入配置**

用于指定不同架构的本地镜像 tar 包路径，DPanel 会根据宿主机架构自动导入对应的镜像。

| 字段 | 说明 |
|------|------|
| `amd64` | AMD64 架构的镜像 tar 包路径 |
| `arm64` | ARM64 架构的镜像 tar 包路径 |
| `arm` | ARMv7 架构的镜像 tar 包路径 |

**路径说明**：
- 路径相对于插件的 asset 目录
- 例如：`asset/plugin/dpanel-plugin-explorer/image-amd64.tar`
- 支持的架构：amd64、arm64、arm

**external - 外部配置**

| 字段 | 说明 |
|------|------|
| `volumes` | 外部挂载卷数组 |
| `volumes_from` | 继承容器的卷数组 |

**使用场景**：
- 插件需要使用自定义镜像
- 镜像未发布到公共仓库
- 需要支持多架构
- 镜像打包在插件内部

---

## 三、最佳实践

### 3.1 命名规范

- **插件名称**：小写字母、数字、连字符
- **字段名称**：大写字母、下划线
- **容器名称**：使用 `${DP_PLUGIN_NAME}`

### 3.2 配置设计

- 重要字段使用多语言 label
- 提供合理的默认值
- 添加清晰的 description
- 必填项标记 `required: true`
- 合理使用 group 组织相关字段

### 3.3 模板编写

- 使用 `#{{-` 注释模板逻辑
- 条件渲染使用 `if` 判断
- 多值配置使用 `range` 循环
- 注意变量作用域

---

## 四、参考资源

- **DPanel 项目**：https://github.com/donknap/dpanel
- **Go Template**：https://pkg.go.dev/text/template
- **Docker Compose**：https://docs.docker.com/compose/

---

## 五、配置示例

本仓库中的 manifest.yaml 和 docker-compose.yaml 提供了完整的配置示例。
