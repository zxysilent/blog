<style lang="less">
@import "../../../styles/common.less";
</style>


<template>
    <div>
        <Row class="margin-top-20">
            <Card>
                <div id="editor">
                    <mavon-editor style="height: 100%"></mavon-editor>
                </div>
            </Card>
        </Row>
    </div>
</template>

<script>
import { mavonEditor } from "mavon-editor";
// import "mavon-editor/dist/css/index.css";
export default {
  name: "markdown-editor",
  components: {
    mavonEditor
    // or 'mavon-editor': mavonEditor
  },
  mounted() {}
};
</script>
