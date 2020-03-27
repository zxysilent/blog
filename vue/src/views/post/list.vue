<template>
	<Card dis-hover>
		<Form inline>
			<FormItem>
				<Select v-model="cate_id" placeholder="请选择文章类别" style="width:200px">
					<Option v-for="item in cateAll" :value="item.id" :key="item.id">{{ item.name }} [{{ item.intro }}]</Option>
				</Select>
			</FormItem>
			<!-- <FormItem>
                        <Input type="text" v-model="mult" placeholder="标题关键词"></Input>
                    </FormItem> -->
			<FormItem>
				<Button type="primary" @click="search()">查&nbsp;询</Button>
			</FormItem>
		</Form>
		<Table size="small" stripe :columns="colPost" :data="dataPost"></Table>
		<br>
		<Page :total="total" :current.sync="page.pi" :page-size="page.ps" :page-size-opts="[10,12,15,20,30]" @on-change="piChange" @on-page-size-change="psChange" show-sizer show-elevator show-total></Page>
	</Card>
</template>
<script>
import { apiCateAll } from "@/api/cate";
import { apiCatePost, apiPostGet, admPostDrop } from "@/api/post";
export default {
	data() {
		return {
			mult: "",
			cateAll: [],
			cate_id: 0, //all
			post: {},
			page: { pi: 1, ps: 15 },
			total: 0,
			colPost: [
				{ type: "index", width: 60, align: "center" },
				{
					title: "标题",
					ellipsis: true,
					tooltip: true,
					render: (h, data) => {
						return h("div", data.row.title);
					}
				},
				{
					title: "状态",
					width: 80,
					render: (h, data) => {
						if (data.row.status == 3) {
							return h("div", "已发布");
						} else {
							return h("div", "草稿");
						}
					}
				},
				{
					title: "权限",
					width: 80,
					render: (h, data) => {
						if (data.row.is_public) {
							return h("div", "公开");
						} else {
							return h("div", "私有");
						}
					}
				},
				{ title: "点击量", width: 80, key: "hits" },
				{
					title: "日期",
					width: 150,
					render: (h, data) => {
						return h("div", data.row.create_time.replace(/T|\+08:00/g, " "));
					}
				},
				{
					title: "操作",
					key: "action",
					width: 150,
					align: "center",
					render: (h, data) => {
						return h("a", [
                             h("Icon", {
								props: { type: "md-eye", size: "20", color: "#5FB878" },
								attrs: { title: "预览" },
								style: { marginRight: "15px" },
								on: {
									click: () => {
										window.open(process.env.VUE_APP_SRV+"/post/"+data.row.path+".html")
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
											this.delete(data);
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
			dataPost: []
		};
	},
	methods: {
		init_() {
			apiCateAll().then(resp => {
				if (resp.code == 200) {
					resp.data.unshift({ id: 0, name: " all ", intro: "全部分类" });
					this.cateAll = resp.data;
				} else {
					this.cateAll = [];
					this.$Message.warning("未查询到分类信息,请重试！");
				}
			});
		},
		init() {
			apiCatePost(this.cate_id, this.page).then(resp => {
				if (resp.code == 200) {
					this.dataPost = resp.data.items;
					this.total = resp.data.count;
				} else {
					this.dataPost = [];
					this.total = 0;
					this.$Message.warning("未查询到数据，请重试！");
				}
			});
		},
		search() {
			this.page.pi = 1;
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
		delete(data) {
			admPostDrop(data.row.id).then(resp => {
				if (resp.code == 200) {
					this.$Message.success({
						content: "删除成功",
						onClose: () => {
							this.dataPost.splice(data.index, 1);
						}
					});
				} else {
					this.$Message.error("删除失败,请重试！");
				}
			});
		}
	},
	mounted() {
		this.init_();
		this.init();
	}
};
</script>
