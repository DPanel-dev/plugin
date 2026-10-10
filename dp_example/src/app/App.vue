<script setup>
import { onMounted, ref } from "vue";

const config = ref({});
const loading = ref(true);
const error = ref("");

const fields = [
  { key: "DISPLAY_TITLE", label: "展示标题" },
  { key: "DISPLAY_MESSAGE", label: "展示内容" },
  { key: "DISPLAY_COUNT", label: "示例数字" },
  { key: "DISPLAY_MODE", label: "显示模式" },
  { key: "DISPLAY_TAGS", label: "标签" },
  { key: "DISPLAY_COLOR", label: "主题颜色" },
  { key: "DISPLAY_NOTE", label: "外观备注" },
  { key: "CONTACTS", label: "联系人" },
  { key: "SERVER_PORT", label: "示例服务端口" },
];

onMounted(async () => {
  try {
    const apiUrl = `${new URL(import.meta.url).origin}/api/test`;
    const response = await fetch(apiUrl, { cache: "no-store" });
    if (!response.ok) throw new Error(`HTTP ${response.status}`);
    config.value = await response.json();
  } catch (cause) {
    error.value = `读取配置失败：${cause.message}`;
  } finally {
    loading.value = false;
  }
});
</script>

<template>
  <main class="card">
    <p class="eyebrow">DPanel 插件示例</p>
    <h1>配置测试</h1>
    <p class="description">以下内容由后端的 /api/test 接口读取容器环境变量后返回。</p>
    <p v-if="loading">正在加载配置…</p>
    <dl v-else-if="!error">
      <div v-for="field in fields" :key="field.key">
        <dt>{{ field.label }}</dt>
        <dd>{{ config[field.key] || "（未填写）" }}</dd>
      </div>
    </dl>
    <p v-if="error" class="error" role="alert">{{ error }}</p>
  </main>
</template>
