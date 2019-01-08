<template>
    <div @click="showMessage" class="message-con">
        <Tooltip :content="value > 0 ? '有' + value + '条未读消息' : '无未读消息'" placement="bottom">
            <Badge :count="value" dot>
                <Icon type="ios-text-outline" :size="22"></Icon>
            </Badge>
        </Tooltip>
    </div>
</template>

<script>
import util from "@/util/util.js";
export default {
  name: "messageTip",
  props: {
    value: {
      type: Number,
      default: 0
    }
  },
  methods: {
    showMessage() {
      this.$router.push({
        name: "message_index"
      });
    }
  }
};
</script>
