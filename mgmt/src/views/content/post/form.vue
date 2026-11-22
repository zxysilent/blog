<template>
  <div>
    <div class="n-layout-page-header">
      <n-card :bordered="false" :title="isEdit ? '编辑博文' : '新增博文'" />
    </div>

    <n-card :bordered="false" class="mt-4 proCard">
      <n-spin :show="loading">
        <n-form ref="formRef" :model="form" :rules="rules" label-placement="left" label-width="80">
          <n-grid :cols="24" :x-gap="16">
            <n-form-item-gi :span="24" label="标题" path="title">
              <n-input v-model:value="form.title" placeholder="请输入标题" />
            </n-form-item-gi>

            <n-form-item-gi :span="12" label="栏目" path="cate_id">
              <n-tree-select
                v-model:value="form.cate_id"
                :options="cateOptions"
                :disabled="isEdit"
                placeholder="请选择栏目"
                key-field="id"
                label-field="name"
                children-field="children"
                clearable
                @update:value="onCateChange"
              />
            </n-form-item-gi>

            <n-form-item-gi :span="12" label="访问路径" path="path">
              <n-input v-model:value="form.path" placeholder="如 hello-world">
                <template #suffix>.html</template>
              </n-input>
            </n-form-item-gi>

            <!-- 笔记用 pid 自关联成目录树，其余类型层级由栏目决定 -->
            <n-form-item-gi v-if="form.kind === 3" :span="12" label="上级笔记">
              <n-tree-select
                v-model:value="form.pid"
                :options="noteOptions"
                placeholder="留空为顶级笔记"
                key-field="id"
                label-field="title"
                children-field="children"
                clearable
              />
            </n-form-item-gi>

            <n-form-item-gi :span="12" label="标签">
              <n-select
                v-model:value="form.tag_ids"
                :options="tagOptions"
                multiple
                filterable
                placeholder="可多选"
                max-tag-count="responsive"
              />
            </n-form-item-gi>

            <n-form-item-gi :span="12" label="状态">
              <n-radio-group v-model:value="form.status">
                <n-radio-button :value="1">草稿</n-radio-button>
                <n-radio-button :value="2">发布</n-radio-button>
              </n-radio-group>
            </n-form-item-gi>

            <n-form-item-gi :span="12" label="作者">
              <n-input v-model:value="form.author" placeholder="可选" />
            </n-form-item-gi>

            <n-form-item-gi :span="12" label="来源">
              <n-input v-model:value="form.origin" placeholder="可选" />
            </n-form-item-gi>

            <n-form-item-gi :span="24" label="封面图">
              <n-space vertical class="w-full">
                <n-input v-model:value="form.icon" placeholder="图片地址，可留空" clearable />
                <n-space align="center" :size="12">
                  <n-upload
                    :show-file-list="false"
                    accept="image/*"
                    :custom-request="onUploadCover"
                  >
                    <n-button size="small" :loading="uploading">
                      <template #icon>
                        <n-icon><CloudUploadOutlined /></n-icon>
                      </template>
                      上传封面
                    </n-button>
                  </n-upload>
                  <n-image v-if="form.icon" :src="form.icon" width="96" class="rounded" />
                  <span class="text-sm text-gray-400">列表页右上角与文章卡片用它做缩略图</span>
                </n-space>
              </n-space>
            </n-form-item-gi>

            <n-form-item-gi :span="24" label="摘要">
              <n-input
                v-model:value="form.summary"
                type="textarea"
                :rows="2"
                placeholder="列表页展示的摘要，留空时由正文首段截取"
              />
            </n-form-item-gi>

            <n-form-item-gi :span="24" label="选项">
              <n-space>
                <n-checkbox v-model:checked="form.show">前台可见</n-checkbox>
                <n-checkbox v-model:checked="form.istop">置顶</n-checkbox>
                <n-checkbox v-model:checked="form.allow">允许评论</n-checkbox>
              </n-space>
            </n-form-item-gi>

            <n-form-item-gi :span="24" label="正文">
              <Markdown v-model:value="form.markdown" v-model:html="form.richtext" class="w-full" />
            </n-form-item-gi>
          </n-grid>
        </n-form>

        <n-space justify="end" class="mt-4">
          <n-button type="primary" :loading="submitting" @click="submit">
            {{ isEdit ? '保存修改' : '确认新增' }}
          </n-button>
          <n-button @click="goBack">取消返回</n-button>
        </n-space>
      </n-spin>
    </n-card>
  </div>
</template>

<script lang="ts" setup>
  import { computed, onMounted, reactive, ref } from 'vue';
  import { useRoute, useRouter } from 'vue-router';
  import { useMessage } from 'naive-ui';
  import { CloudUploadOutlined } from '@vicons/antd';
  import { Markdown } from '@/components/Markdown';
  import { addPost, editPost, getPost, getPostTags, getPostTree, setPostTags } from '@/api/post';
  import { getCateTree } from '@/api/cate';
  import { getTagList } from '@/api/tag';
  import { uploadImage } from '@/api/site';
  import type { Cate, Post, Tag } from '@/api/types';

  const route = useRoute();
  const router = useRouter();
  const message = useMessage();

  const formRef = ref();
  const loading = ref(false);
  const submitting = ref(false);
  const cateOptions = ref<Cate[]>([]);
  const tagOptions = ref<{ label: string; value: string }[]>([]);
  const noteOptions = ref<Post[]>([]);
  const uploading = ref(false);

  const postId = computed(() => (route.params.id as string) || '');
  const isEdit = computed(() => !!postId.value);

  const form = reactive({
    title: '',
    cate_id: null as string | null,
    pid: null as string | null,
    path: '',
    kind: 1,
    status: 1,
    author: '',
    origin: '',
    icon: '',
    summary: '',
    markdown: '',
    richtext: '',
    show: true,
    istop: false,
    allow: false,
    tag_ids: [] as string[],
  });

  const rules = {
    title: { required: true, message: '请输入标题', trigger: 'blur' },
    path: { required: true, message: '请输入访问路径', trigger: 'blur' },
    cate_id: { required: true, message: '请选择栏目', trigger: ['blur', 'change'] },
  };

  /** 栏目决定内容类型：文章进首页，页面是固定页，笔记建目录树 */
  function onCateChange(val: string | null) {
    const hit = findCate(cateOptions.value, val);
    if (hit) form.kind = hit.kind;
  }

  function findCate(list: Cate[], id: string | null): Cate | undefined {
    if (!id) return undefined;
    for (const it of list) {
      if (it.id === id) return it;
      const sub = findCate(it.children || [], id);
      if (sub) return sub;
    }
    return undefined;
  }

  async function loadOptions() {
    const [cates, tags, notes] = await Promise.all([
      getCateTree({ use: true }),
      getTagList({ page: 1, size: 1000 }),
      // 笔记父级候选：只要笔记类型，后端直接返回树
      getPostTree({ kind: 3 }).catch(() => []),
    ]);
    cateOptions.value = cates;
    tagOptions.value = tags.list.map((it: Tag) => ({ label: it.name, value: it.id }));
    noteOptions.value = notes;
  }

  /** 封面图上传：走后端 /api/images，返回可访问地址 */
  async function onUploadCover({ file, onFinish, onError }: any) {
    uploading.value = true;
    try {
      const res: any = await uploadImage(file.file as File);
      const url = typeof res === 'string' ? res : res?.data || '';
      if (!url) throw new Error('上传失败');
      form.icon = url;
      message.success('上传成功');
      onFinish();
    } catch (err: any) {
      message.error(err?.message || '上传失败');
      onError();
    } finally {
      uploading.value = false;
    }
  }

  async function loadDetail() {
    if (!postId.value) return;
    loading.value = true;
    try {
      const mod: any = await getPost(postId.value);
      if (mod) {
        Object.assign(form, {
          title: mod.title,
          cate_id: mod.cate_id || null,
          pid: mod.pid || null,
          path: mod.path,
          kind: mod.kind,
          status: mod.status,
          author: mod.author,
          origin: mod.origin,
          icon: mod.icon || '',
          summary: mod.summary,
          markdown: mod.markdown,
          richtext: mod.richtext,
          show: mod.show,
          istop: mod.istop,
          allow: mod.allow,
        });
      }
      const list = await getPostTags(postId.value);
      form.tag_ids = list.map((it) => it.id);
    } finally {
      loading.value = false;
    }
  }

  async function submit() {
    try {
      await formRef.value?.validate();
    } catch {
      return;
    }
    submitting.value = true;
    try {
      const payload: Recordable = {
        title: form.title,
        cate_id: form.cate_id,
        pid: form.kind === 3 ? form.pid || '' : '',
        path: form.path,
        status: form.status,
        author: form.author,
        origin: form.origin,
        icon: form.icon,
        summary: form.summary,
        markdown: form.markdown,
        richtext: form.richtext,
        show: form.show,
        istop: form.istop,
        allow: form.allow,
      };
      let id = postId.value;
      if (isEdit.value) {
        // 类型创建后固定，修改接口不接受 kind
        await editPost(id, payload);
        message.success('修改成功');
      } else {
        payload.kind = form.kind;
        const res: any = await addPost(payload);
        id = res?.data || '';
        message.success('新增成功');
      }
      // 标签走独立接口，全量替换
      if (id) {
        await setPostTags(id, form.tag_ids);
      }
      goBack();
    } finally {
      submitting.value = false;
    }
  }

  function goBack() {
    router.push({ name: 'content_post' });
  }

  onMounted(async () => {
    await loadOptions();
    await loadDetail();
  });
</script>
