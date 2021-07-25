<template>
	<Card dis-hover>
		<Form inline>
			<FormItem>
				<Button :to="{name:'post-add'}" style="margin-right: 8px">添加文章</Button>
				<Select v-model="cate_id" placeholder="请选择文章类别" style="width:200px">
					<Option v-for="item in cateAll" :value="item.id" :key="item.id">{{ item.name }} [{{ item.intro }}]</Option>
				</Select>
			</FormItem>
			<!-- <FormItem>
                        <Input type="text" v-model="mult" placeholder="标题关键词"></Input>
                    </FormItem> -->
			<FormItem>
				<Button type="primary" icon="md-search" @click="search">查&nbsp;&nbsp;询</Button>
			</FormItem>
		</Form>
		<Table size="small" stripe :columns="tabCol" :data="tabData"></Table>
		<br>
		<Page :total="tabCount" :current.sync="tabPage.pi" :page-size="tabPage.ps" :page-size-opts="[8,10,12,15,20,30]" @on-change="onPiChange" @on-page-size-change="onPsChange" show-sizer show-elevator show-total></Page>
	</Card>
</template>
<script>
import { apiCateAll } from "@/api/cate";
import { apiPostPage, admPostDrop } from "@/api/post";
export default {
	data() {
		return {
			cateAll: [],
			cate_id: -1, //all
			tabCount: 0,
			tabPage: { pi: 1, ps: 12, mult: "" },
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
										window.open(process.env.VUE_APP_SRV + "/post/" + data.row.path + ".html");
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
											name: "post-edit",
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
		preinit() {
			apiCateAll().then((resp) => {
				if (resp.code == 200) {
					resp.data.unshift({ id: -1, name: " all ", intro: "全部分类" });
					this.cateAll = resp.data;
				} else {
					this.cateAll = [];
					this.$Message.error({ content: resp.msg, duration: 3 });
				}
			});
		},
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
			apiPostPage({ ...this.tabPage, cate_id: this.cate_id }).then((resp) => {
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
		search() {
			this.tabPage.pi = 1;
			this.init();
		},
		piChange(pi) {
			this.page.pi = pi;
			this.init();
		},
		psChange(ps) {
			this.page.ps = ps;
			this.init();
		},
		//删除
		emitDrop(data) {
			admPostDrop({ id: data.row.id }).then((resp) => {
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
		this.preinit();
		this.init();
	}
};
</script>
