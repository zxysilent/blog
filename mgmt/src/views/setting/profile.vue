<template>
  <div>
    <div class="n-layout-page-header">
      <n-card :bordered="false" title="个人中心" />
    </div>

    <n-grid class="mt-4" :cols="24" :x-gap="16" item-responsive>
      <n-grid-item span="24 900:10">
        <n-card :bordered="false" title="账号信息" class="proCard">
          <n-descriptions :column="1" label-placement="left">
            <n-descriptions-item label="账号">{{ info.num || '-' }}</n-descriptions-item>
            <n-descriptions-item label="角色">{{ info.role?.name || '-' }}</n-descriptions-item>
            <n-descriptions-item label="状态">
              <n-tag :type="info.status === 1 ? 'success' : 'error'" size="small">
                {{ info.status === 1 ? '正常' : '禁用' }}
              </n-tag>
            </n-descriptions-item>
            <n-descriptions-item label="创建时间">{{ fmtTime(info.created) }}</n-descriptions-item>
          </n-descriptions>
        </n-card>
      </n-grid-item>

      <n-grid-item span="24 900:14">
        <n-card :bordered="false" class="proCard">
          <n-tabs type="line" animated>
            <n-tab-pane name="profile" tab="基本资料">
              <n-form
                ref="profileRef"
                :model="profile"
                :rules="profileRules"
                label-placement="left"
                label-width="80"
              >
                <n-form-item label="名称" path="name">
                  <n-input v-model:value="profile.name" placeholder="显示名称" />
                </n-form-item>
                <n-form-item label="头像">
                  <n-input v-model:value="profile.avatar" placeholder="图片地址，可留空" />
                </n-form-item>
                <n-form-item label="手机">
                  <n-input v-model:value="profile.phone" placeholder="可选" />
                </n-form-item>
                <n-form-item label="邮箱">
                  <n-input v-model:value="profile.email" placeholder="可选" />
                </n-form-item>
                <n-form-item>
                  <n-button type="primary" :loading="savingProfile" @click="saveProfile">
                    保存资料
                  </n-button>
                </n-form-item>
              </n-form>
            </n-tab-pane>

            <n-tab-pane name="password" tab="修改密码">
              <n-form
                ref="passwdRef"
                :model="passwd"
                :rules="passwdRules"
                label-placement="left"
                label-width="80"
              >
                <n-form-item label="原密码" path="opasswd">
                  <n-input
                    v-model:value="passwd.opasswd"
                    type="password"
                    show-password-on="click"
                    placeholder="请输入当前密码"
                  />
                </n-form-item>
                <n-form-item label="新密码" path="npasswd">
                  <n-input
                    v-model:value="passwd.npasswd"
                    type="password"
                    show-password-on="click"
                    placeholder="至少 6 位"
                  />
                </n-form-item>
                <n-form-item label="确认密码" path="confirm">
                  <n-input
                    v-model:value="passwd.confirm"
                    type="password"
                    show-password-on="click"
                    placeholder="再次输入新密码"
                  />
                </n-form-item>
                <n-form-item>
                  <n-button type="primary" :loading="savingPasswd" @click="savePasswd">
                    修改密码
                  </n-button>
                </n-form-item>
              </n-form>
            </n-tab-pane>
          </n-tabs>
        </n-card>
      </n-grid-item>
    </n-grid>
  </div>
</template>

<script lang="ts" setup>
  import { onMounted, reactive, ref } from 'vue';
  import { useMessage } from 'naive-ui';
  import { changePassword, getUserInfo, updateProfile } from '@/api/auth';
  import { fmtTime } from '@/api/common';
  import type { User } from '@/api/types';
  import { useUserStore } from '@/store/user';

  const message = useMessage();
  const userStore = useUserStore();

  const info = ref<Partial<User>>({});
  const profileRef = ref();
  const passwdRef = ref();
  const savingProfile = ref(false);
  const savingPasswd = ref(false);

  const profile = reactive({ name: '', avatar: '', phone: '', email: '' });
  const passwd = reactive({ opasswd: '', npasswd: '', confirm: '' });

  const profileRules = {
    name: { required: true, message: '请输入名称', trigger: 'blur' },
  };

  // 新密码至少 6 位由后端把关，前端只做即时反馈
  const passwdRules = {
    opasswd: { required: true, message: '请输入当前密码', trigger: 'blur' },
    npasswd: {
      required: true,
      trigger: 'blur',
      validator: (_rule, value: string) => {
        if (!value || value.length < 6) return new Error('新密码至少 6 位');
        return true;
      },
    },
    confirm: {
      required: true,
      trigger: 'blur',
      validator: (_rule, value: string) => {
        if (value !== passwd.npasswd) return new Error('两次输入的新密码不一致');
        return true;
      },
    },
  };

  async function load() {
    const data = await getUserInfo();
    info.value = data || {};
    profile.name = data?.name || '';
    profile.avatar = data?.avatar || '';
    profile.phone = data?.phone || '';
    profile.email = data?.email || '';
  }

  async function saveProfile() {
    try {
      await profileRef.value?.validate();
    } catch {
      return;
    }
    savingProfile.value = true;
    try {
      await updateProfile({ ...profile });
      // 同步刷新 store 与页面数据，头部昵称会立即更新
      await userStore.getInfo();
      await load();
      message.success('资料已保存');
    } finally {
      savingProfile.value = false;
    }
  }

  async function savePasswd() {
    try {
      await passwdRef.value?.validate();
    } catch {
      return;
    }
    savingPasswd.value = true;
    try {
      await changePassword({ opasswd: passwd.opasswd, npasswd: passwd.npasswd });
      message.success('密码已修改，下次登录请使用新密码');
      passwd.opasswd = '';
      passwd.npasswd = '';
      passwd.confirm = '';
    } finally {
      savingPasswd.value = false;
    }
  }

  onMounted(() => {
    load().catch(() => undefined);
  });
</script>
