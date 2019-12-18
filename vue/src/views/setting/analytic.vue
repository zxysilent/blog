
<template>
	<Card dis-hover>
		<p slot="title">
			<Icon type="ios-pulse" /> 统计设置
		</p>
		<div>
			<Form label-position="top">
				<FormItem label="网站统计代码">
					<Input v-model="model.value" style="width:600px" type="textarea" :autosize="{minRows: 15,maxRows: 20}" placeholder="Enter code "></Input>
				</FormItem>
				<div>
					<Button type="warning" :loading="saveLoading" @click="cmtSave">保&nbsp;&nbsp;&nbsp;&nbsp;存</Button>
				</div>
			</Form>
		</div>
	</Card>
</template>
<script>
import { apiOptsGet, admOptsEdit } from "@/api/opts";
export default {
	data() {
		return {
			model: { key: "analytic", value: "" },
			saveLoading: false
		};
	},
	methods: {
		cmtSave() {
			this.saveLoading = true;
			admOptsEdit(this.model).then(resp => {
				this.saveLoading = false;
				if (resp.code == 200) {
					this.$Message.success({ content: "网站统计代码,更新成功" });
				} else {
					this.$Message.error({
						content: `网站统计代码,更新失败,请重试`,
						duration: 3,
						onClose() {
							this.init();
						}
					});
				}
			});
		},
		init() {
			apiOptsGet(this.model.key).then(resp => {
				if (resp.code == 200) {
					this.model.value = resp.data;
				} else {
					this.model.value = "";
				}
			});
		}
	},
	mounted() {
		this.init();
	}
};
</script>
