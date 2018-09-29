
<template>
    <div>
        <Card dis-hover>
            <p slot="title">
                <Icon type="ios-settings-outline" /> 基本设置
            </p>
            <div>
                <Form ref="base" :model="base" label-position="top" :rules="base">
                    <FormItem label="站点名称" prop="title">
                        <div style="display:inline-block;width:600px;">
                            <Input v-model="base.title" search enter-button="确    定" @on-search="cmtSave('title')"></Input>
                        </div>
                    </FormItem>
                    <FormItem label="Logo" prop="logo_url">
                        <div style="display:inline-block;width:600px;">
                            <Input v-model="base.logo_url" search enter-button="确    定" @on-search="cmtSave('logo_url')"></Input>
                        </div>
                    </FormItem>
                    <FormItem label="站点描述" prop="description">
                        <div style="display:inline-block;width:600px;">
                            <Input v-model="base.description" search enter-button="确    定" @on-search="cmtSave('description')" ></Input>
                        </div>
                    </FormItem>
                    <FormItem label="站点地址" prop="site_url">
                        <div style="display:inline-block;width:600px;">
                            <Input v-model="base.site_url" search enter-button="确    定" @on-search="cmtSave('site_url')"></Input>
                        </div>
                    </FormItem>
                    <FormItem label="Favicon" prop="favicon_url">
                        <div style="display:inline-block;width:600px;">
                            <Input v-model="base.favicon_url" search enter-button="确    定" @on-search="cmtSave('favicon_url')"></Input>
                        </div>
                    </FormItem>
                    <FormItem label="关键词" prop="keywords">
                        <div style="display:inline-block;width:600px;">
                            <Input v-model="base.keywords" search enter-button="确    定" @on-search="cmtSave('keywords')"></Input>
                        </div>
                    </FormItem>
                    <FormItem label="每页文章数目" prop="page_size">
                        <div style="display:inline-block;width:600px;">
                            <br>
                            <Slider v-model="page_size" :min="4" :max="12" show-stops show-tip="always" @on-change="onChange"></Slider>
                            此数目用于指定文章每页显示的文章数目.
                        </div>
                    </FormItem>
                    <FormItem label="Github" prop="github_url">
                        <div style="display:inline-block;width:600px;">
                            <Input v-model="base.github_url" search enter-button="确    定" @on-search="cmtSave('github_url')"></Input>
                        </div>
                    </FormItem>
                    <FormItem label="Weibo" prop="weibo_url">
                        <div style="display:inline-block;width:600px;">
                            <Input v-model="base.weibo_url" search enter-button="确    定" @on-search="cmtSave('weibo_url')"></Input>
                        </div>
                    </FormItem>
                    <FormItem label="备案号" prop="miitbeian">
                        <div style="display:inline-block;width:600px;">
                            <Input v-model="base.miitbeian" search enter-button="确    定" @on-search="cmtSave('miitbeian')"></Input>
                        </div>
                    </FormItem>
                </Form>
            </div>
        </Card>
    </div>
</template>
<script>
import { optsBase, optsEdit } from "@/api/opts";
export default {
  data() {
    return {
      page_size: 2,
      base: {},
      saveLoading: false
    };
  },
  methods: {
    cmtSave(key) {
      this.save_loading = true;
      optsEdit({
        key: key,
        value: this.base[key]
      }).then(resp => {
        this.saveLoading = false;
        if (resp.code == 200) {
          this.$Message.success({
            content: "更新成功"
          });
        } else {
          this.$Message.error({
            content: `更新失败`,
            duration: 3,
            onClose() {
              this.init();
            }
          });
        }
      });
    },
    onChange(val) {
      optsEdit({
        key: "page_size",
        value: val + ""
      }).then(resp => {
        this.saveLoading = false;
        if (resp.code != 200) {
          this.$Message.error({
            content: `更新失败,请重试`,
            duration: 3,
            onClose() {
              this.init();
            }
          });
        }
      });
    },
    init() {
      optsBase().then(resp => {
        if (resp.code == 200) {
          this.page_size = Number(resp.data["page_size"], 10);
          console.log(resp.data);
          this.base = resp.data;
        } else {
          this.base = {};
        }
      });
    }
  },
  mounted() {
    this.init();
  }
};
</script>
