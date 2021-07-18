<style lang="less" scoped>
.ivu-form-inline .ivu-form-item {
	margin-bottom: 10px;
}
</style>
<template>
	<Card dis-hover>
		<Form inline>
			<FormItem>
				<Button v-auth="'record_show'" type="info" @click="init" icon="md-refresh" title="刷新数据">刷&nbsp;&nbsp;新</Button>
			</FormItem>
		</Form>
		<Table size="small" :columns="tabCol" :data="tabData"></Table>
		<Page :total="tabCount" :current.sync="tabPage.pi" :page-size="tabPage.ps" :page-size-opts="[8,10,12,15,20,30]" @on-change="onPiChange" @on-page-size-change="onPsChange" show-sizer show-elevator show-total></Page>
	</Card>
</template>
<script>
import { admRecordPage,admRecordDrop } from "@/api/home";
export default {
	data() {
		return {
			tabPage: { pi: 1, ps: 12 },
			tabCount: 0,
			tabCol: [
				{ type: "index", width: 80, align: "center" },
				// { title: "用户id", key: "user_id" },
				{
					title: "用户",
					width: 150,
					render: (h, data) => {
						return h("div", data.row.user.name);
					}
				},
				{ title: "模块", key: "module", width: 120 },
				{ title: "做了什么", key: "action", width: 120 },
				{
					title: "时间",
					width: 150,
					render: (h, data) => {
						return h("div", data.row.ctime.replace(/T|\+08:00/g, " "));
					}
				},
				{ title: "数据", key: "data", minWidth: 150, tooltip: true },
				{
					title: "操作",
					width: 150,
					align: "center",
					render: (h, data) => {
						let actios = [];
						if (this.$auth("record_drop")) {
							actios.push(
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
							);
						}
						return h("a", actios);
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
			admRecordPage(this.tabPage).then((resp) => {
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
			admRecordDrop({ id: data.row.id }).then((resp) => {
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
