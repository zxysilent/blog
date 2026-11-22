<template>
  <div>
    <div class="n-layout-page-header">
      <n-card :bordered="false" title="内容预览" />
    </div>

    <n-card :bordered="false" class="mt-4 proCard">
      <n-spin :show="loading">
        <n-empty v-if="!mod" description="内容不存在或已被删除" class="py-10" />

        <template v-else>
          <!-- 元信息 -->
          <div class="flex flex-wrap items-center gap-2 mb-3">
            <n-tag size="small" :type="mod.status === 2 ? 'success' : 'default'">
              {{ mod.status === 2 ? '已发布' : '草稿' }}
            </n-tag>
            <n-tag size="small" type="info">{{ kindText }}</n-tag>
            <n-tag v-if="mod.istop" size="small" type="error" :bordered="false">置顶</n-tag>
            <n-tag v-if="!mod.show" size="small" :bordered="false">前台隐藏</n-tag>
            <n-tag v-if="mod.inner" size="small" type="warning" :bordered="false">内置</n-tag>
          </div>

          <h1 class="text-2xl font-bold mb-3">{{ mod.title }}</h1>

          <n-image
            v-if="mod.icon"
            :src="mod.icon"
            class="mb-4 rounded"
            :img-props="{ style: 'max-width:100%;border-radius:6px' }"
          />

          <n-space :size="16" class="mb-4 text-sm text-gray-500">
            <span>栏目：{{ cateText }}</span>
            <span v-if="mod.author">作者：{{ mod.author }}</span>
            <span v-if="mod.origin">来源：{{ mod.origin }}</span>
            <span>浏览：{{ mod.hits }}</span>
            <span>更新：{{ fmtTime(mod.updated) }}</span>
            <span>创建：{{ fmtTime(mod.created) }}</span>
          </n-space>

          <n-space :size="8" class="mb-4" align="center">
            <span class="text-sm text-gray-500">访问路径</span>
            <n-tag :bordered="false" type="default">{{ mod.path }}</n-tag>
            <n-button size="tiny" quaternary @click="copyPath">复制</n-button>
          </n-space>

          <n-space v-if="mod.tags?.length" :size="8" class="mb-4">
            <n-tag v-for="tg in mod.tags" :key="tg.id" size="small" :bordered="false">
              #{{ tg.name }}
            </n-tag>
          </n-space>

          <n-divider />

          <div v-if="mod.summary" class="mb-4 p-3 rounded bg-gray-50 text-gray-600">
            摘要：{{ mod.summary }}
          </div>

          <!-- 正文：服务端渲染结果，再经 DOMPurify 白名单过滤 -->
          <div class="markdown-body" v-html="safeHtml"></div>
          <n-empty v-if="!safeHtml" description="该内容还没有正文" class="py-8" />

          <n-divider />

          <!-- 上下篇 -->
          <n-space :size="12">
            <n-button quaternary :disabled="!prevPost" @click="gotoPost(prevPost)">
              ← {{ prevPost?.title || '没有上一篇' }}
            </n-button>
            <n-button quaternary :disabled="!nextPost" @click="gotoPost(nextPost)">
              {{ nextPost?.title || '没有下一篇' }} →
            </n-button>
          </n-space>

          <n-space justify="end" class="mt-6">
            <n-button type="primary" @click="goEdit">编辑内容</n-button>
            <n-button @click="goBack">返回列表</n-button>
          </n-space>
        </template>
      </n-spin>
    </n-card>
  </div>
</template>

<script lang="ts" setup>
  import { computed, onMounted, ref, watch } from 'vue';
  import { useRoute, useRouter } from 'vue-router';
  import { useMessage } from 'naive-ui';
  import DOMPurify from 'dompurify';
  import { marked } from 'marked';
  import { getPost, getPostNeighbors } from '@/api/post';
  import { getCateTree } from '@/api/cate';
  import { fmtTime } from '@/api/common';
  import type { Cate, Post } from '@/api/types';

  const route = useRoute();
  const router = useRouter();
  const message = useMessage();

  const loading = ref(false);
  const mod = ref<Post | null>(null);
  const cateOptions = ref<Cate[]>([]);
  const neighbors = ref<{ prev?: Post; next?: Post }>({});

  // 后端找不到相邻内容时返回的是空结构（id 为空），这里统一成 undefined
  const prevPost = computed(() => (neighbors.value.prev?.id ? neighbors.value.prev : undefined));
  const nextPost = computed(() => (neighbors.value.next?.id ? neighbors.value.next : undefined));

  const postId = computed(() => (route.params.id as string) || '');
  /** 预览权限与列表一致：打开详情即视为读过 */
  const kindMap: Record<number, string> = { 1: '文章', 2: '页面', 3: '笔记' };
  const kindText = computed(() => kindMap[mod.value?.kind ?? 1] || '-');

  /**
   * 正文：优先用服务端渲染的 richtext；
   * 部分内容（如接口直录）只有 markdown，这里做一次本地渲染兵底。
   * 不管哪种来源，输出前都过一遍 DOMPurify。
   */
  const safeHtml = computed(() => {
    const mod_ = mod.value;
    if (!mod_) return '';
    const rich = mod_.richtext || '';
    if (rich) return DOMPurify.sanitize(rich);
    const md = mod_.markdown || '';
    return md ? DOMPurify.sanitize(marked.parse(md) as string) : '';
  });

  const cateText = computed(() => {
    const id = mod.value?.cate_id || '';
    if (!id) return '-';
    return findCate(cateOptions.value, id)?.name || '-';
  });

  function findCate(list: Cate[], id: string): Cate | undefined {
    for (const it of list) {
      if (it.id === id) return it;
      const sub = findCate(it.children || [], id);
      if (sub) return sub;
    }
    return undefined;
  }

  async function load() {
    if (!postId.value) return;
    loading.value = true;
    try {
      const [detail, near] = await Promise.all([
        getPost(postId.value),
        getPostNeighbors(postId.value).catch(() => ({})),
      ]);
      mod.value = detail || null;
      neighbors.value = (near as any) || {};
    } finally {
      loading.value = false;
    }
  }

  async function copyPath() {
    try {
      await navigator.clipboard.writeText(mod.value?.path || '');
      message.success('路径已复制');
    } catch {
      message.warning('浏览器不允许复制，请手动选择');
    }
  }

  function gotoPost(item?: Post) {
    if (!item?.id) return;
    router.push({ name: 'content_post_detail', params: { id: item.id } });
  }

  function goEdit() {
    router.push({ name: 'content_post_edit', params: { id: postId.value } });
  }

  function goBack() {
    router.push({ name: 'content_post' });
  }

  watch(postId, load);

  onMounted(async () => {
    cateOptions.value = await getCateTree({});
    await load();
  });
</script>
