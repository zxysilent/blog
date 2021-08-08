<style lang="less" scoped>
.ivu-form-inline .ivu-form-item {
	margin-bottom: 10px;
}
</style>
<template>
	<Card dis-hover>
		<Form inline>
			<!-- <FormItem>
				<Select v-model="cid" placeholder="请选择类别" style="width:200px">
					<Option v-for="item in columnAll" :value="item.id" :key="item.id">{{ item.name }}</Option>
				</Select>
			</FormItem> -->
			<FormItem>
				<Button :to="{name:'dict-add'}" style="margin-right: 8px">添加字典</Button>
				<Button type="info" @click="init" icon="md-refresh" title="刷新数据">刷&nbsp;&nbsp;新</Button>
			</FormItem>
		</Form>
		<Table size="small" :columns="tabCol" :data="tabData"></Table>
		<Page :total="tabCount" :current.sync="tabPage.pi" :page-size="tabPage.ps" :page-size-opts="[8,10,12,15,20,30]" @on-change="onPiChange" @on-page-size-change="onPsChange" show-sizer show-elevator show-total></Page>
	</Card>
</template>
<script>
import { apiDictPage, admDictDrop } from "@/api/dict";
export default {
	data() {
		return {
			tabPage: { pi: 1, ps: 12 },
			tabCount: 0,
			tabCol: [
				{ type: "index", width: 80, align: "center" },
				{ title: "字典名", minWidth: 100, maxWidth: 200, key: "key" },
				{
					title: "类型",
					minWidth: 80,
					maxWidth: 120,
					key: "inner",
					render: (h, params) => {
						const row = params.row;
						const color = row.inner ? "success" : "primary";
						const text = row.inner ? "内置" : "普通";
						return h(
							"Tag",
							{
								props: {  color: color, }
							},
							text
						);
					}
				},
				{ title: "字典值", minWidth: 200, key: "value", tooltip: true },
				{ title: "字典说明", minWidth: 150, key: "intro", tooltip: true },
				{
					title: "操作",
					width: 200,
					align: "center",
					render: (h, data) => {
						return h("a", [
							h("Icon", {
								props: { type: "md-create", size: "20", color: "#FFB800" },
								attrs: { title: "修改" },
								style: { marginRight: "15px" },
								on: {
									click: () => {
										this.$router.push({
											name: "dict-edit",
											params: { key: data.row.key }
										});
									}
								}
							}),
							h(
								"Poptip",
								{
									props: { confirm: true, title: "确定要删除吗？" },
									on: {
										"on-ok": () => {
											this.emitDrop(data);
										}
									}
								},
								[
									h("Icon", {
										props: { type: "ios-trash", size: "20", color: "#FF5722" },
										attrs: { title: "删除" }
									})
								]
							)
						]);
					}
				}
			],
			tabData: []
		};
	},
	methods: {
		onPiChange(pi) {
			this.tabPage.pi = pi;
			this.init();
		},
		onPsChange(ps) {
			this.tabPage.ps = ps;
			this.init();
		},
		init() {
			apiDictPage(this.tabPage).then((resp) => {
				if (resp.code == 200) {
					this.tabData = resp.data.items;
					this.tabCount = resp.data.count;
				} else {
					this.tabData = [];
					this.tabCount = 0;
					this.$Message.error({ content: resp.msg, duration: 3 });
				}
			});
		},
		emitDrop(data) {
			admDictDrop({ id: data.row.id }).then((resp) => {
				if (resp.code == 200) {
					this.$Message.success({
						content: "删除成功",
						onClose: () => {
							this.tabData.splice(data.index, 1);
						}
					});
				} else {
					this.$Message.error({ content: resp.msg, duration: 3 });
				}
			});
		}
	},
	mounted() {
		this.init();
	}
};
</script>
