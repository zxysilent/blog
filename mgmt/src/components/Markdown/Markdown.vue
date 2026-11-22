<template>
  <!-- 编辑态用 MdEditor：内置工具条、分栏/预览/全屏/目录 -->
  <MdEditor
    v-if="!readonly"
    ref="editorRef"
    :id="EDITOR_ID"
    :model-value="model"
    :theme="theme"
    preview-theme="github"
    :code-theme="codeTheme"
    language="zh-CN"
    :toolbars="toolbars"
    :footers="footers"
    :def-toolbars="(moreToolbars as any)"
    :placeholder="placeholder"
    :style="{ height }"
    @update:model-value="onChange"
    @on-upload-img="onUploadImg"
  />
  <!-- 查看态用 MdPreview：只有渲染结果，无工具条 -->
  <MdPreview
    v-else
    :id="`${EDITOR_ID}-preview`"
    :model-value="model"
    :theme="theme"
    preview-theme="github"
    :code-theme="codeTheme"
    language="zh-CN"
  />
</template>

<script lang="ts" setup>
  import { computed, h, ref, watch } from 'vue';
  import { MdEditor, MdPreview, NormalToolbar } from 'md-editor-v3';
  import { useMessage } from 'naive-ui';
  import { marked } from 'marked';
  import DOMPurify from 'dompurify';
  import { useDesignSettingStore } from '@/store/designSetting';
  import { uploadImage } from '@/api/site';
  // 编辑器样式随本组件一起进入路由懒加载分块，只有打开博文表单才会下载
  import 'md-editor-v3/lib/style.css';

  /**
   * Markdown 编辑器（md-editor-v3）
   *
   * 编辑器只负责「写 + 预览」，落库的 richtext 仍由 marked + DOMPurify 生成，
   * 与历史数据、博客前台渲染保持一致；编辑器自带的渲染结果不参与存储。
   */
  const props = withDefaults(
    defineProps<{
      /** markdown 原文 */
      value: string;
      /** 渲染后的 html，用于回写 richtext */
      html?: string;
      readonly?: boolean;
      /** 编辑器高度 */
      height?: string;
    }>(),
    { html: '', readonly: false, height: '620px' }
  );

  const emit = defineEmits<{
    (e: 'update:value', v: string): void;
    (e: 'update:html', v: string): void;
  }>();

  const EDITOR_ID = 'blog-md-editor';
  const message = useMessage();
  const designStore = useDesignSettingStore();
  const editorRef = ref<any>(null);

  const model = computed(() => props.value || '');
  const theme = computed(() => (designStore.darkTheme ? 'dark' : 'light'));
  const codeTheme = computed(() => (designStore.darkTheme ? 'one-dark' : 'atom'));
  const placeholder = '在这里书写 Markdown 正文...';

  /**
   * 工具条：去掉外链（github）、手动保存（本页用底部按钮提交）与 html 预览；
   * mermaid / katex 未内置渲染依赖，插入后预览看不到效果，也不放出来。
   * 数字 0 指向 #defToolbars 里的第一个自定义按钮（摘要分隔符）。
   */
  const toolbars: any[] = [
    'bold',
    'underline',
    'italic',
    'strikeThrough',
    '-',
    'title',
    'sub',
    'sup',
    '-',
    'quote',
    'unorderedList',
    'orderedList',
    'task',
    '-',
    'codeRow',
    'code',
    'link',
    'image',
    'table',
    '-',
    'revoke',
    'next',
    '-',
    'prettier',
    0,
    '-',
    'pageFullscreen',
    'fullscreen',
    'preview',
    'catalog',
  ];

  const footers: any[] = ['markdownTotal', '=', 'scrollSwitch'];

  /** markdown → 安全 html，作为 richtext 的落库内容 */
  function render(md: string): string {
    if (!md.trim()) return '';
    const raw = marked.parse(md, { async: false }) as string;
    return DOMPurify.sanitize(raw);
  }

  function onChange(v: string) {
    emit('update:value', v);
  }

  // 只跟着 markdown 走：纯富文本来源的记录（markdown 为空）不覆盖 richtext
  watch(
    model,
    (v) => {
      emit('update:html', render(v));
    },
    { immediate: true }
  );

  /** 在光标处插入摘要分隔符 */
  function insertMore() {
    editorRef.value?.insert?.(() => ({ targetValue: '\n<!--more-->\n', select: false }));
  }

  /**
   * 工具条里的「摘要」按钮：v7 的数字索引取的是 defToolbars 属性里的 VNode 数组，
   * 插槽形式不会被索引到，所以这里直接在脚本里构造。
   */
  const moreToolbars = computed(() => [
    h(
      NormalToolbar,
      { title: '插入摘要分隔符 <!--more-->', onClick: insertMore },
      { default: () => h('span', { style: 'padding: 0 4px;font-size: 12px;line-height: 1' }, '摘要') }
    ),
  ]);

  /** 图片上传：复用后台 /api/images，返回可访问地址 */
  async function onUploadImg(files: File[], callback: (urls: string[]) => void) {
    try {
      const urls = await Promise.all(
        files.map(async (file) => {
          const res: any = await uploadImage(file);
          const url = typeof res === 'string' ? res : res?.data || '';
          if (!url) throw new Error('上传失败');
          return url;
        })
      );
      callback(urls);
    } catch (err: any) {
      message.error(err?.message || '图片上传失败');
    }
  }

  defineExpose({ insertMore });
</script>
