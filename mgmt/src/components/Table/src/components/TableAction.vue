<template>
  <div class="tableAction">
    <div class="flex items-center justify-center">
      <!-- 主操作：图标模式只显示图标 + tooltip，文字模式保持原样 -->
      <template v-for="(action, index) in getActions" :key="`${index}-${action.label}`">
        <n-tooltip v-if="useIcon(action)" trigger="hover" :delay="300">
          <template #trigger>
            <n-button class="mx-1" circle quaternary size="small" @click="action.onClick">
              <template #icon>
                <n-icon :component="action.icon" />
              </template>
            </n-button>
          </template>
          {{ action.label }}
        </n-tooltip>
        <n-button v-else v-bind="omitLabel(action) as any" class="mx-1">
          {{ action.label }}
          <template #icon v-if="action.hasOwnProperty('icon')">
            <n-icon :component="action.icon" />
          </template>
        </n-button>
      </template>

      <!-- 更多：其余操作统一收进下拉，避免操作列过宽 -->
      <n-dropdown
        v-if="dropDownActions && getDropdownList.length"
        trigger="hover"
        :options="getDropdownList as any"
        @select="select as any"
      >
        <slot name="more"></slot>
        <template v-if="!$slots.more">
          <n-tooltip v-if="iconOnly" trigger="hover" :delay="300">
            <template #trigger>
              <n-button class="mx-1" circle quaternary size="small">
                <template #icon>
                  <n-icon :component="EllipsisOutlined" />
                </template>
              </n-button>
            </template>
            更多
          </n-tooltip>
          <n-button v-else v-bind="getMoreProps as any" class="mx-1" icon-placement="right">
            <div class="flex items-center">
              <span>更多</span>
              <n-icon size="14" class="ml-1">
                <DownOutlined />
              </n-icon>
            </div>
          </n-button>
        </template>
      </n-dropdown>
    </div>
  </div>
</template>

<script lang="ts">
  import { defineComponent, PropType, computed, toRaw } from 'vue';
  import { ActionItem } from '@/components/Table';
  import { usePermission } from '@/hooks/usePermission';
  import { isBoolean, isFunction } from '@/utils/is';
  import { DownOutlined, EllipsisOutlined } from '@vicons/antd';

  export default defineComponent({
    name: 'TableAction',
    components: { DownOutlined, EllipsisOutlined },
    props: {
      actions: {
        type: Array as PropType<ActionItem[]>,
        default: null,
        required: true,
      },
      dropDownActions: {
        type: Array as PropType<ActionItem[]>,
        default: null,
      },
      style: {
        type: String as PropType<string>,
        default: 'button',
      },
      select: {
        type: Function as PropType<(...args: any[]) => any>,
        default: () => {},
      },
      /**
       * 图标模式：主操作只渲染图标 + tooltip，节省横向空间
       * 默认开启；置为 false 可退回「文字按钮 + 更多」的老样式
       */
      iconOnly: {
        type: Boolean,
        default: true,
      },
    },
    setup(props) {
      const { hasPermission } = usePermission();

      const actionType =
        props.style === 'button' ? 'default' : props.style === 'text' ? 'primary' : 'default';
      const actionText =
        props.style === 'button' ? undefined : props.style === 'text' ? true : undefined;

      const getMoreProps = computed(() => {
        return {
          text: actionText,
          type: actionType,
          size: 'small',
        };
      });

      /** label 不是 n-button 的合法 prop，渲染时摘掉避免落到 DOM 属性上 */
      function omitLabel(action: Recordable): Recordable {
        const { label, ...rest } = action as Recordable;
        return rest;
      }

      /** 图标模式判断：Actions 可能带渲染时附加的内部字段，这里只看 icon 是否存在 */
      function useIcon(action: ActionItem | Recordable): boolean {
        return props.iconOnly && !!(action as ActionItem).icon;
      }

      const getDropdownList = computed(() => {
        return (toRaw(props.dropDownActions) || [])
          .filter((action) => {
            return hasPermission(action.auth as string[]) && isIfShow(action);
          })
          .map((action) => {
            const { popConfirm } = action;
            return {
              size: 'small',
              text: actionText,
              type: actionType,
              ...action,
              ...popConfirm,
              onConfirm: popConfirm?.confirm,
              onCancel: popConfirm?.cancel,
            };
          });
      });

      function isIfShow(action: ActionItem): boolean {
        const ifShow = action.ifShow;

        let isIfShow = true;

        if (isBoolean(ifShow)) {
          isIfShow = ifShow;
        }
        if (isFunction(ifShow)) {
          isIfShow = ifShow(action);
        }
        return isIfShow;
      }

      const getActions = computed(() => {
        return (toRaw(props.actions) || [])
          .filter((action) => {
            return hasPermission(action.auth as string[]) && isIfShow(action);
          })
          .map((action) => {
            const { popConfirm } = action;
            //需要展示什么风格，自己修改一下参数
            return {
              size: 'small',
              text: actionText,
              type: actionType,
              ...action,
              ...(popConfirm || {}),
              onConfirm: popConfirm?.confirm,
              onCancel: popConfirm?.cancel,
              enable: !!popConfirm,
            };
          });
      });

      return {
        getActions,
        getDropdownList,
        getMoreProps,
        useIcon,
        omitLabel,
        EllipsisOutlined,
      };
    },
  });
</script>
