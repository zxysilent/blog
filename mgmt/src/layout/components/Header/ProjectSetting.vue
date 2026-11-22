<template>
  <n-drawer v-model:show="isDrawer" :width="width" :placement="placement">
    <n-drawer-content :title="drawerTitle" :native-scrollbar="false">
      <div class="drawer">
        <n-divider title-placement="center">{{ '系统主题' }}</n-divider>

        <div class="drawer-setting-item align-items-top">
          <span
            class="theme-item"
            v-for="(item, index) in appThemeList"
            :key="index"
            :style="{ 'background-color': item }"
            @click="togTheme(item)"
          >
            <n-icon size="12" v-if="item === designStore.appTheme">
              <CheckOutlined />
            </n-icon>
          </span>
        </div>

        <n-divider title-placement="center">{{ '导航栏模式' }}</n-divider>

        <div class="drawer-setting-item align-items-top">
          <div class="drawer-setting-item-style align-items-top">
            <n-tooltip placement="top">
              <template #trigger>
                <img
                  src="~@/assets/images/nav-theme-dark.svg"
                  @click="togNavMode('vertical')"
                  :alt="'左侧菜单模式'"
                />
              </template>
              <span>{{ '左侧菜单模式' }}</span>
            </n-tooltip>
            <n-badge dot color="#19be6b" v-show="settingStore.navMode === 'vertical'" />
          </div>

          <div class="drawer-setting-item-style">
            <n-tooltip placement="top">
              <template #trigger>
                <img
                  src="~@/assets/images/nav-horizontal.svg"
                  :alt="'顶部菜单模式'"
                  @click="togNavMode('horizontal')"
                />
              </template>
              <span>{{ '顶部菜单模式' }}</span>
            </n-tooltip>
            <n-badge dot color="#19be6b" v-show="settingStore.navMode === 'horizontal'" />
          </div>

          <div class="drawer-setting-item-style">
            <n-tooltip placement="top">
              <template #trigger>
                <img
                  src="~@/assets/images/nav-horizontal-mix.svg"
                  @click="togNavMode('horizontal-mix')"
                  :alt="'顶部菜单混合模式'"
                />
              </template>
              <span>{{ '顶部菜单混合模式' }}</span>
            </n-tooltip>
            <n-badge dot color="#19be6b" v-show="settingStore.navMode === 'horizontal-mix'" />
          </div>
        </div>

        <n-divider title-placement="center">{{ '导航栏风格' }}</n-divider>

        <div class="drawer-setting-item align-items-top">
          <div class="drawer-setting-item-style align-items-top">
            <n-tooltip placement="top">
              <template #trigger>
                <img
                  src="~@/assets/images/nav-theme-dark.svg"
                  :alt="'暗色侧边栏'"
                  @click="togNavTheme('dark')"
                />
              </template>
              <span>{{ '暗色侧边栏' }}</span>
            </n-tooltip>
            <n-badge dot color="#19be6b" v-if="settingStore.navTheme === 'dark'" />
          </div>

          <div class="drawer-setting-item-style">
            <n-tooltip placement="top">
              <template #trigger>
                <img
                  src="~@/assets/images/nav-theme-light.svg"
                  :alt="'白色侧边栏'"
                  @click="togNavTheme('light')"
                />
              </template>
              <span>{{ '白色侧边栏' }}</span>
            </n-tooltip>
            <n-badge dot color="#19be6b" v-if="settingStore.navTheme === 'light'" />
          </div>

          <div class="drawer-setting-item-style">
            <n-tooltip placement="top">
              <template #trigger>
                <img
                  src="~@/assets/images/header-theme-dark.svg"
                  @click="togNavTheme('header-dark')"
                  :alt="'暗色顶栏'"
                />
              </template>
              <span>{{ '暗色顶栏' }}</span>
            </n-tooltip>
            <n-badge dot color="#19be6b" v-if="settingStore.navTheme === 'header-dark'" />
          </div>
        </div>

        <n-divider title-placement="center">{{ '登录页风格' }}</n-divider>

        <div class="drawer-setting-item align-items-top">
          <div class="drawer-setting-item-style align-items-top">
            <n-tooltip placement="top">
              <template #trigger>
                <img
                  src="~@/assets/images/login-style-mgmt.svg"
                  :alt="'管理端风格'"
                  @click="togLoginStyle('mgmt')"
                />
              </template>
              <span>{{ '管理端风格' }}</span>
            </n-tooltip>
            <n-badge dot color="#19be6b" v-if="settingStore.loginStyle === 'mgmt'" />
          </div>

          <div class="drawer-setting-item-style align-items-top">
            <n-tooltip placement="top">
              <template #trigger>
                <img
                  src="~@/assets/images/login-style-vue.svg"
                  :alt="'旧版前台风格'"
                  @click="togLoginStyle('vue')"
                />
              </template>
              <span>{{ '旧版前台风格' }}</span>
            </n-tooltip>
            <n-badge dot color="#19be6b" v-if="settingStore.loginStyle === 'vue'" />
          </div>
        </div>
        <div v-if="!canLoginStyle" class="drawer-setting-tip">需要「登录页风格」权限才能切换</div>

        <n-divider title-placement="center">{{ '界面功能' }}</n-divider>

        <div class="drawer-setting-item">
          <div class="drawer-setting-item-title"> {{ '分割菜单' }} </div>
          <div class="drawer-setting-item-action">
            <n-switch
              :disabled="readonly || settingStore.navMode !== 'horizontal-mix'"
              v-model:value="settingStore.menuSetting.mixMenu"
            />
          </div>
        </div>

        <div class="drawer-setting-item">
          <div class="drawer-setting-item-title"> {{ '固定顶栏' }} </div>
          <div class="drawer-setting-item-action">
            <n-switch :disabled="readonly" v-model:value="settingStore.headerSetting.fixed" />
          </div>
        </div>

        <!--        <div class="drawer-setting-item">-->
        <!--          <div class="drawer-setting-item-title">-->
        <!--            固定侧边栏-->
        <!--          </div>-->
        <!--          <div class="drawer-setting-item-action">-->
        <!--            <n-switch :disabled="readonly" v-model:value="settingStore.menuSetting.fixed" />-->
        <!--          </div>-->
        <!--        </div>-->

        <div class="drawer-setting-item">
          <div class="drawer-setting-item-title"> {{ '固定多页签' }} </div>
          <div class="drawer-setting-item-action">
            <n-switch :disabled="readonly" v-model:value="settingStore.multiTabsSetting.fixed" />
          </div>
        </div>

        <n-divider title-placement="center">{{ '界面显示' }}</n-divider>

        <div class="drawer-setting-item">
          <div class="drawer-setting-item-title"> {{ '显示重载页面按钮' }} </div>
          <div class="drawer-setting-item-action">
            <n-switch :disabled="readonly" v-model:value="settingStore.headerSetting.isReload" />
          </div>
        </div>

        <div class="drawer-setting-item">
          <div class="drawer-setting-item-title"> {{ '显示面包屑导航' }} </div>
          <div class="drawer-setting-item-action">
            <n-switch :disabled="readonly" v-model:value="settingStore.crumbsSetting.show" />
          </div>
        </div>

        <div class="drawer-setting-item">
          <div class="drawer-setting-item-title"> {{ '显示面包屑图标' }} </div>
          <div class="drawer-setting-item-action">
            <n-switch :disabled="readonly" v-model:value="settingStore.crumbsSetting.showIcon" />
          </div>
        </div>

        <div class="drawer-setting-item">
          <div class="drawer-setting-item-title"> {{ '显示多页签' }} </div>
          <div class="drawer-setting-item-action">
            <n-switch :disabled="readonly" v-model:value="settingStore.multiTabsSetting.show" />
          </div>
        </div>
        <!--1.15废弃，没啥用，占用操作空间-->
        <!--        <div class="drawer-setting-item">-->
        <!--          <div class="drawer-setting-item-title"> 显示页脚 </div>-->
        <!--          <div class="drawer-setting-item-action">-->
        <!--            <n-switch :disabled="readonly" v-model:value="settingStore.showFooter" />-->
        <!--          </div>-->
        <!--        </div>-->

        <n-divider title-placement="center">{{ '动画' }}</n-divider>

        <div class="drawer-setting-item">
          <div class="drawer-setting-item-title"> {{ '启用动画' }} </div>
          <div class="drawer-setting-item-action">
            <n-switch :disabled="readonly" v-model:value="settingStore.isPageAnimate" />
          </div>
        </div>

        <div class="drawer-setting-item">
          <div class="drawer-setting-item-title"> {{ '动画类型' }} </div>
          <div class="drawer-setting-item-select">
            <n-select :disabled="readonly" v-model:value="settingStore.pageAnimateType" :options="animateOptions" />
          </div>
        </div>

        <div class="drawer-setting-item">
          <n-alert :type="readonly ? 'warning' : 'info'" :showIcon="false">
            <p>{{ readonly ? '无修改权限，仅可预览' : '修改后自动保存' }}</p>
          </n-alert>
        </div>
      </div>
    </n-drawer-content>
  </n-drawer>
</template>

<script lang="ts">
  import { defineComponent, reactive, toRefs, unref, watch, computed } from 'vue';
  import { useProjectSettingStore } from '@/store/projectSetting';
  import { useDesignSettingStore } from '@/store/designSetting';
  import { usePermission } from '@/hooks/usePermission';
  import { CheckOutlined } from '@vicons/antd';
  import { darkTheme } from 'naive-ui';
  import { useMessage } from 'naive-ui';
  import { animates } from '@/settings/animateSetting';

/** 动画类型中文名 */
/** 登录页风格的权限标识，与后端 model.PermLoginStyle 保持一致 */
const PERM_LOGIN_STYLE = 'login:style';

const ANIMATE_LABEL: Recordable<string> = {
  'zoom-fade': '渐变',
  'zoom-out': '闪现',
  'fade-slide': '滑动',
  fade: '消退',
  'fade-bottom': '底部消退',
  'fade-scale': '缩放消退',
};

  export default defineComponent({
    name: 'ProjectSetting',
    components: { CheckOutlined },
    props: {
      title: {
        type: String,
        default: '',
      },
      width: {
        type: Number,
        default: 280,
      },
    },
    setup(props) {
      const settingStore = useProjectSettingStore();
      const designStore = useDesignSettingStore();
      const { hasPermission } = usePermission();
      const message = useMessage();

      // 布局配置有独立权限点（config:layout）：无权限时只能预览
      const readonly = computed(() => !hasPermission(['config:layout']));
      // 登录页风格另有 login:style：两个权限都可能单独授予
      const canLoginStyle = computed(() => !readonly.value && hasPermission([PERM_LOGIN_STYLE]));

      /** 上次已保存的快照，用于跳过「加载完成后又原样写回」的无意义请求 */
      let lastSaved = '';
      let saveTimer: ReturnType<typeof setTimeout> | null = null;

      const persist = () => {
        if (readonly.value) return;
        const snapshot = settingStore.snapshot();
        if (snapshot === lastSaved) return;
        lastSaved = snapshot;
        settingStore.saveLayoutSetting().catch((err) => {
          message.error(err?.msg || '保存失败');
          // 保存失败时清空快照，下一次改动会重试
          lastSaved = '';
        });
      };

      /** 只监听会入库的字段，改动后 400ms 落库，避免连续点击刷接口 */
      const persisted = computed(() => [
        settingStore.navMode,
        settingStore.navTheme,
        settingStore.headerSetting.fixed,
        settingStore.headerSetting.isReload,
        settingStore.multiTabsSetting.show,
        settingStore.multiTabsSetting.fixed,
        settingStore.menuSetting.fixed,
        settingStore.menuSetting.mixMenu,
        settingStore.menuSetting.collapsed,
        settingStore.menuSetting.menuWidth,
        settingStore.menuSetting.minMenuWidth,
        settingStore.crumbsSetting.show,
        settingStore.crumbsSetting.showIcon,
        settingStore.showFooter,
        settingStore.isPageAnimate,
        settingStore.pageAnimateType,
      ]);

      watch(persisted, () => {
        if (saveTimer) clearTimeout(saveTimer);
        saveTimer = setTimeout(persist, 400);
      });

      const animateOptions = computed(() =>
        animates.map((item) => ({
          value: item.value,
          label: ANIMATE_LABEL[item.value] ?? item.value,
        }))
      );

      const state = reactive({
        width: props.width,
        isDrawer: false,
        placement: 'right' as const,
        appThemeList: designStore.appThemeList,
      });

      const drawerTitle = computed(() => props.title || '布局配置');

      const directionsOptions = computed(() => {
        return animateOptions.value.find(
          (item) => item.value == unref(settingStore.pageAnimateType)
        );
      });

      function openDrawer() {
        state.isDrawer = true;
        // 每次打开都拉一次最新值，多人协作时不至于看到别人的旧设置
        settingStore.loadLayoutSetting(true).then(() => {
          lastSaved = settingStore.snapshot();
        });
      }

      function closeDrawer() {
        state.isDrawer = false;
        // 关闭时立即落库，不必等防抖计时到
        if (saveTimer) clearTimeout(saveTimer);
        persist();
      }

      function togNavTheme(theme) {
        if (readonly.value) return;
        settingStore.navTheme = theme;
        if (settingStore.navMode === 'horizontal' && ['light'].includes(theme)) {
          settingStore.navTheme = 'dark';
        }
      }

      function togTheme(color) {
        if (readonly.value) return;
        designStore.appTheme = color;
      }

      function togNavMode(mode) {
        if (readonly.value) return;
        settingStore.navMode = mode;
        settingStore.menuSetting.mixMenu = false;
      }

      /**
       * 切换登录页风格
       * 风格存于 layout 配置，但改动需 login:style 权限，所以走专用接口，
       * 不随整份布局一起提交
       */
      function togLoginStyle(style: 'mgmt' | 'vue') {
        if (!canLoginStyle.value) {
          message.warning('没有「登录页风格」权限');
          return;
        }
        if (settingStore.loginStyle === style) return;
        settingStore.saveLoginStyle(style).then(
          () => message.success('登录页风格已切换'),
          (err) => message.error(err?.msg || '切换失败')
        );
      }

      return {
        ...toRefs(state),
        drawerTitle,
        settingStore,
        designStore,
        readonly,
        canLoginStyle,
        togNavTheme,
        togNavMode,
        togLoginStyle,
        togTheme,
        darkTheme,
        openDrawer,
        closeDrawer,
        animateOptions,
        directionsOptions,
      };
    },
  });
</script>

<style lang="less" scoped>
  .drawer {
    .n-divider:not(.n-divider--vertical) {
      margin: 10px 0;
    }

    &-setting-item {
      display: flex;
      align-items: center;
      padding: 12px 0;
      flex-wrap: wrap;

      &-style {
        display: inline-block;
        position: relative;
        margin-right: 16px;
        cursor: pointer;
        text-align: center;
      }

      &-title {
        flex: 1 1;
        font-size: 14px;
      }

      &-action {
        flex: 0 0 auto;
      }

      &-select {
        flex: 1;
      }

      .theme-item {
        width: 20px;
        min-width: 20px;
        height: 20px;
        cursor: pointer;
        border: 1px solid #eee;
        border-radius: 2px;
        margin: 0 5px 5px 0;
        text-align: center;
        line-height: 14px;

        .n-icon {
          color: #fff;
        }
      }
    }

    .align-items-top {
      align-items: flex-start;
      padding: 2px 0;
    }

    // 权限不足时的说明文字，跟在选项下方
    .drawer-setting-tip {
      padding: 0 0 8px;
      font-size: 12px;
      line-height: 1.5;
      color: #f0a020;
    }

    .justify-center {
      justify-content: center;
    }

  }
</style>
