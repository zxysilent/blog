<template>
	<Card dis-hover>
		<Table size="small" stripe :columns="colPage" :data="dataPage"></Table>
	</Card>
</template>
<script>
import { apiPageAll, admPostDrop } from "@/api/post";
export default {
	data() {
		return {
			colPage: [
				{ type: "index", width: 60, align: "center" },
				{ title: "标题", ellipsis: true, tooltip: true, key: "title" },
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
					title: "创建日期",
					width: 150,
					render: (h, data) => {
						return h("div", data.row.create_time.replace(/T|\+08:00/g, " "));
					}
				},
				{
					title: "修改日期",
					width: 150,
					render: (h, data) => {
						return h("div", data.row.update_time.replace(/T|\+08:00/g, " "));
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
										window.open(process.env.VUE_APP_SRV+"/page/"+data.row.path+".html")
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
			dataPage: []
		};
	},
	methods: {
		init() {
			apiPageAll().then(resp => {
				if (resp.code == 200) {
					this.dataPage = resp.data;
				} else {
					this.dataPage = [];
					this.$Message.warning("未查询到信息，请重试！");
				}
			});
		},
		//删除
		delete(data) {
			admPostDrop(data.row.id).then(resp => {
				if (resp.code == 200) {
					this.$Message.success({
						content: "删除成功",
						onClose: () => {
							this.dataPage.splice(data.index, 1);
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
