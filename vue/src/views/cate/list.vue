<template>
	<div>
		<Card dis-hover>
			<div slot="title">
				<Icon type="ios-school-outline" /> 分类管理
			</div>
			<Form inline>
				<!-- <FormItem>
				<Select v-model="cid" placeholder="请选择文章类别" style="width:200px">
					<Option v-for="item in columnAll" :value="item.id" :key="item.id">{{ item.name }}</Option>
				</Select>
			</FormItem> -->
				<FormItem>
					<Button :to="{name:'cate-add'}" style="margin-right: 8px">添加分类</Button>
					<Button type="info" @click="init" icon="md-refresh" title="刷新数据">刷&nbsp;&nbsp;新</Button>
				</FormItem>
			</Form>
			<Table stripe size="small" :columns="tabCol" :data="tabData"></Table>
			<Page :total="tabCount" :current.sync="tabPage.pi" :page-size="tabPage.ps" :page-size-opts="[8,10,12,15,20,30]" @on-change="onPiChange" @on-page-size-change="onPsChange" show-sizer show-elevator show-total></Page>
		</Card>
	</div>
</template>
<script>
import { apiCatePage, admCateDrop } from "@/api/cate";
export default {
	data() {
		return {
			tabCount: 0,
			tabPage: { pi: 1, ps: 12 },
			tabCol: [
				{ type: "index", minWidth: 60, maxWidth: 100, align: "center" },
				{ title: "分类名", minWidth: 100, maxWidth: 300, key: "name" },
				{ title: "分类介绍", minWidth: 100, maxWidth: 300, key: "intro" },
				{
					title: "Action",
					minWidth: 100,
					align: "left",
					render: (h, data) => {
						return h("a", [
							h("Icon", {
								props: { type: "md-create", size: "20", color: "#FFB800" },
								attrs: { title: "修改" },
								style: { marginRight: "15px" },
								on: {
									click: () => {
										this.$router.push({
											name: "cate-edit",
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
			apiCatePage(this.tabPage).then((resp) => {
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
			admCateDrop({ id: data.row.id }).then((resp) => {
				if (resp.code == 200) {
					this.$Message.success({
						content: "删除成功",
						onClose: () => {
							this.tabData.splice(data.index, 1);
						}
					});
				} else {
					this.$Message.error("删除失败,请重试！");
				}
			});
		}
	},
	mounted() {
		this.init();
	}
};
</script>
