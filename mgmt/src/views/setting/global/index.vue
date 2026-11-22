<template>
  <div>
    <div class="n-layout-page-header">
      <n-card :bordered="false" title="全局配置" />
    </div>

    <n-card :bordered="false" class="mt-4 proCard">
      <n-spin :show="loading">
        <n-empty v-if="!tabs.length" description="暂无配置" class="py-8" />
        <n-tabs
          v-else
          v-model:value="activeTab"
          type="line"
          size="large"
          animated
          @update:value="onTabChange"
        >
          <n-tab-pane v-for="item in tabs" :key="item.id" :name="item.id" :tab="item.tab">
            <!-- 元信息 -->
            <n-space class="mb-3" :size="8" align="center">
              <n-tag size="small" :bordered="false" type="info">{{ item.raw.group }}</n-tag>
              <n-tag size="small" :bordered="false">{{ item.raw.type }}</n-tag>
              <n-text depth="3" class="text-sm">{{ item.raw.describe || '无描述' }}</n-text>
              <n-tag v-if="!item.editable" size="small" :bordered="false" type="warning">
                只读
              </n-tag>
            </n-space>

            <!-- 按 schema 渲染字段，避免把一段裸 JSON 丢给使用者 -->
            <n-form
              ref="formRef"
              :model="draft"
              label-placement="left"
              label-width="130"
              :disabled="readonly"
            >
              <n-grid :cols="24" :x-gap="24">
                <n-form-item-gi
                  v-for="field in item.schema.fields"
                  :key="field.key"
                  :span="field.type === 'textarea' ? 24 : 12"
                  :label="field.label"
                >
                  <!-- 统一用 :value + @update:value，避免在模板里写类型断言 -->
                  <div class="w-full">
                    <n-input
                      v-if="field.type === 'text'"
                      :value="draft[field.key]"
                      :placeholder="field.placeholder"
                      clearable
                      @update:value="(v: string) => (draft[field.key] = v)"
                    />
                    <n-input
                      v-else-if="field.type === 'password'"
                      :value="draft[field.key]"
                      type="password"
                      show-password-on="click"
                      :placeholder="field.placeholder"
                      @update:value="(v: string) => (draft[field.key] = v)"
                    />
                    <n-input
                      v-else-if="field.type === 'textarea'"
                      :value="draft[field.key]"
                      type="textarea"
                      :rows="3"
                      :placeholder="field.placeholder"
                      @update:value="(v: string) => (draft[field.key] = v)"
                    />
                    <n-input-number
                      v-else-if="field.type === 'number'"
                      :value="draft[field.key]"
                      :min="field.min"
                      :max="field.max"
                      :step="field.step || 1"
                      style="width: 100%"
                      @update:value="(v: number | null) => (draft[field.key] = v ?? 0)"
                    >
                      <template v-if="field.suffix" #suffix>{{ field.suffix }}</template>
                    </n-input-number>
                    <n-switch
                      v-else-if="field.type === 'switch'"
                      :value="!!draft[field.key]"
                      @update:value="(v: boolean) => (draft[field.key] = v)"
                    />
                    <!-- 单选：如登录页风格，选项说明跟着选项走 -->
                    <n-radio-group
                      v-else-if="field.type === 'radio'"
                      :value="draft[field.key] || field.options?.[0]?.value"
                      :disabled="readonly"
                      @update:value="(v: string) => (draft[field.key] = v)"
                    >
                      <n-space vertical :size="8">
                        <n-radio v-for="opt in field.options || []" :key="opt.value" :value="opt.value">
                          <div class="flex flex-col">
                            <span>{{ opt.label }}</span>
                            <span v-if="opt.tip" class="option-tip">{{ opt.tip }}</span>
                          </div>
                        </n-radio>
                      </n-space>
                    </n-radio-group>
                    <div v-if="field.tip" class="field-tip">{{ field.tip }}</div>
                  </div>
                </n-form-item-gi>
              </n-grid>
            </n-form>

            <div class="config-actions">
              <n-space>
                <n-button v-if="readonly && item.editable" type="tertiary" @click="openEdit">
                  开启编辑
                </n-button>
                <n-button v-if="!readonly" type="primary" :loading="saving" @click="save">
                  保存
                </n-button>
                <n-button type="success" @click="reload">刷新</n-button>
              </n-space>
              <div class="config-meta">
                <span>修改时间：{{ fmtTime(item.updated) }}</span>
                <span>创建时间：{{ fmtTime(item.created) }}</span>
              </div>
            </div>
          </n-tab-pane>
        </n-tabs>
      </n-spin>
    </n-card>
  </div>
</template>

<script lang="ts" setup>
  /**
   * 全局配置
   *
   * 内置配置（basic/login/blog/auth/upload/sql）的表单化编辑入口，
   * 字段来自 schema.ts，与后端取值区间保持一致。
   * 「配置管理」页只负责列表浏览与自建配置的增删改元信息，两页职责分开。
   */
  import { computed, onMounted, reactive, ref } from 'vue';
  import { NTag, useMessage } from 'naive-ui';
  import { editConfig, getConfigList, setLoginStyle } from '@/api/config';
  import { fmtTime } from '@/api/common';
  import { usePermission } from '@/hooks/usePermission';
  import type { Config } from '@/api/types';
  import { configSchemas, type Field } from './schema';

  const message = useMessage();
  const { hasPermission } = usePermission();

  const loading = ref(false);
  const saving = ref(false);
  const readonly = ref(true);
  const activeTab = ref('');
  const list = ref<Config[]>([]);

  /**
   * 不在本页展示的内置配置
   * layout（布局配置）的编辑入口是右上角「布局配置」抽屉，那里已是可视化开关，
   * 这里再列一份会让人以为两边各自生效
   */
  const HIDDEN_IDS = ['layout'];

  /** 当前编辑副本，切换页签时从 list 重新填充 */
  const draft = reactive<Recordable>({});

  interface TabItem {
    id: string;
    tab: string;
    schema: (typeof configSchemas)[number];
    editable: boolean;
    updated: number;
    created: number;
    raw: Config;
  }

  /** 页签只取 schema 中声明且库里确实存在的配置，顺序跟 schema 一致 */
  const tabs = computed<TabItem[]>(() => {
    const byId = new Map(list.value.map((it) => [it.id, it]));
    const out: TabItem[] = [];
    configSchemas.forEach((sc) => {
      if (HIDDEN_IDS.includes(sc.id)) return;
      const mod = byId.get(sc.id);
      if (!mod) return;
      out.push({
        id: mod.id,
        tab: sc.tab,
        schema: sc,
        editable: mod.editable,
        updated: mod.updated,
        created: mod.created,
        raw: mod,
      });
    });
    return out;
  });

  const current = computed(() => tabs.value.find((it) => it.id === activeTab.value));

  /** 把配置的 JSON 文本解析到 draft，解析失败时保留空对象 */
  function fillDraft(mod: Config) {
    Object.keys(draft).forEach((k) => delete draft[k]);
    try {
      const obj = JSON.parse(mod.value || '{}');
      if (obj && typeof obj === 'object' && !Array.isArray(obj)) {
        Object.assign(draft, obj);
      }
    } catch {
      // 非 JSON 的脏数据：保留空表单，保存时以原值为底合并
    }
  }

  function onTabChange() {
    readonly.value = true;
    if (current.value) fillDraft(current.value.raw);
  }

  function openEdit() {
    readonly.value = false;
  }

  async function load() {
    loading.value = true;
    try {
      const res = await getConfigList({ page: 1, size: 1000 });
      list.value = res.list;
      if (!tabs.value.some((it) => it.id === activeTab.value)) {
        activeTab.value = tabs.value[0]?.id || '';
      }
      if (current.value) fillDraft(current.value.raw);
    } finally {
      loading.value = false;
    }
  }

  function reload() {
    readonly.value = true;
    load();
  }

  /** 未填写的字段按类型取默认值，避免存进空串 */
  function normFieldValue(f: Field, v: any): any {
    if (v !== undefined && v !== null && v !== '') return v;
    if (f.type === 'switch') return false;
    if (f.type === 'radio') return f.options?.[0]?.value ?? '';
    return '';
  }

  /** 解析配置原始 JSON，作为合并基底 */
  function parseRaw(mod: Config): Recordable {
    try {
      const obj = JSON.parse(mod.value || '{}');
      return obj && typeof obj === 'object' && !Array.isArray(obj) ? obj : {};
    } catch {
      return {};
    }
  }

  /** 前端范围校验，与后端保持一致，减少一次无效请求 */
  function checkRange(item: TabItem): string {
    for (const f of item.schema.fields) {
      if (f.type !== 'number') continue;
      const v = Number(draft[f.key]);
      if (Number.isNaN(v)) return `${f.label}必须是数字`;
      if (f.min !== undefined && v < f.min) return `${f.label}不能小于 ${f.min}`;
      if (f.max !== undefined && v > f.max) return `${f.label}不能大于 ${f.max}`;
    }
    return '';
  }

  async function save() {
    const item = current.value;
    if (!item) return;
    if (!item.editable) {
      message.warning('该配置为只读，不允许在线修改');
      return;
    }
    const bad = checkRange(item);
    if (bad) {
      message.warning(bad);
      return;
    }
    // 以原值为底再覆盖表单字段：schema 只声明了部分字段时，
    // 未出现在表单里的键不会被写成空值
    const original = parseRaw(item.raw);
    const obj: Recordable = { ...original };
    item.schema.fields.forEach((f) => {
      obj[f.key] = normFieldValue(f, draft[f.key]);
    });
    saving.value = true;
    try {
      await editConfig(item.id, { value: JSON.stringify(obj) });
      message.success('修改成功');
      readonly.value = true;
      await load();
    } finally {
      saving.value = false;
    }
  }

  onMounted(load);
</script>

<style lang="less" scoped>
  .field-tip {
    margin-top: 4px;
    font-size: 12px;
    line-height: 1.5;
    color: var(--text-color-3);
  }

  .option-tip {
    font-size: 12px;
    line-height: 1.5;
    color: var(--text-color-3);
  }

  .config-actions {
    display: flex;
    flex-wrap: wrap;
    gap: 12px;
    align-items: center;
    justify-content: space-between;
    margin-top: 18px;
  }

  .config-meta {
    display: flex;
    gap: 16px;
    font-size: 12px;
    color: var(--text-color-3);
  }
</style>
