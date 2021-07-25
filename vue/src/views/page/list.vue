<template>
	<Card dis-hover>
		<Form inline>
			<!-- <FormItem>
				<Select v-model="cid" placeholder="请选择文章类别" style="width:200px">
					<Option v-for="item in columnAll" :value="item.id" :key="item.id">{{ item.name }}</Option>
				</Select>
			</FormItem> -->
			<FormItem>
				<Button :to="{name:'page-add'}" style="margin-right: 8px">添加页面</Button>
				<Button type="info" @click="init" icon="md-refresh" title="刷新数据">刷&nbsp;&nbsp;新</Button>
			</FormItem>
		</Form>
		<Table size="small" stripe :columns="tabCol" :data="tabData"></Table>
		<Page :total="tabCount" :current.sync="tabPage.pi" :page-size="tabPage.ps" :page-size-opts="[8,10,12,15,20,30]" @on-change="onPiChange" @on-page-size-change="onPsChange" show-sizer show-elevator show-total></Page>
	</Card>
</template>
<script>
import { apiPagePage, admPageDrop } from "@/api/page";
export default {
	data() {
		return {
			tabCount: 0,
			tabPage: { pi: 1, ps: 12 },
			tabCol: [
				{ type: "index", width: 60, align: "center" },
				{ title: "标题", minWidth: 200, ellipsis: true, tooltip: true, key: "title" },
				{
					title: "状态",
					width: 100,
					render: (h, data) => {
						if (data.row.status == 2) {
							return h("div", "已发布");
						} else {
							return h("div", "草稿");
						}
					}
				},
				// {
				// 	title: "权限",
				// 	width: 80,
				// 	render: (h, data) => {
				// 		if (data.row.is_public) {
				// 			return h("div", "公开");
				// 		} else {
				// 			return h("div", "私有");
				// 		}
				// 	}
				// },
				// { title: "点击量", width: 80, key: "hits" },
				{
					title: "创建日期",
					width: 200,
					render: (h, data) => {
						return h("div", data.row.created.replace(/T|\+08:00/g, " "));
					}
				},
				{
					title: "修改日期",
					width: 200,
					render: (h, data) => {
						return h("div", data.row.updated.replace(/T|\+08:00/g, " "));
					}
				},
				{
					title: "操作",
					key: "action",
					width: 200,
					align: "center",
					render: (h, data) => {
						return h("a", [
							h("Icon", {
								props: { type: "md-eye", size: "20", color: "#5FB878" },
								attrs: { title: "预览" },
								style: { marginRight: "15px" },
								on: {
									click: () => {
										window.open(process.env.VUE_APP_SRV + "/page/" + data.row.path + ".html");
									}
								}
							}),
							h("Icon", {
								props: { type: "md-create", size: "20", color: "#FFB800" },
								attrs: { title: "修改" },
								style: { marginRight: "15px" },
								on: {
									click: () => {
										this.$router.push({
											name: "page-edit",
											params: { id: data.row.id }
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
			this.tabPage.pi = 1;
			this.tabPage.ps = ps;
			this.init();
		},
		init() {
			apiPagePage(this.tabPage).then((resp) => {
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
		//删除
		emitDrop(data) {
			admPageDrop({ id: data.row.id }).then((resp) => {
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
