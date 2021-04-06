<style lang="less">
</style>
<template>
	<Card dis-hover>
		<p slot="title">
			<Icon type="md-person" /> 用户列表
		</p>
		<Table size="small" border :columns="tabCol" :data="tabData"></Table>
		<br>
		<Page :total="page.total" :current.sync="page.pi" :page-size="page.ps" :page-size-opts="[10,12,15]" @on-change="piChange" @on-page-size-change="psChange" show-sizer show-elevator show-total></Page>
	</Card>
</template>
<script>
import { admUserPage, admUserDrop, admUserReset, admUserChgAtv } from "@/api/user";
export default {
	data() {
		return {
			page: { pi: 1, ps: 12, total: 0 },
			tabCol: [
				{ type: "index", width: 60, align: "center" },
				{ title: "登录账号", key: "num", width: 200 },
				{ title: "用户名", width: 280, key: "name" },
				{
					title: "用户角色",
					minWidth: 150,
					render: (h, data) => {
						if (data.row.role == 4096) {
							return h("div", "系统管理员");
						} else if (data.row.role == 2048) {
							return h("div", "内容管理员");
						}
					}
				},
				{ title: "用户电话", width: 150, key: "phone" },
				{ title: "用户邮箱", width: 150, key: "email" },
				// {
				// 	title: "状态",
				// 	width: 100,
				// 	key: "istop",
				// 	render: (h, data) => {
				// 		let taht = this;
				// 		return h(
				// 			"i-switch",
				// 			{
				// 				props: { value: this.getRole(40), size: "large" },
				// 				on: {
				// 					"on-change": function(val) {
				// 						admUserChgAtv({ id: data.row.id, flag: val }).then(resp => {
				// 							if (resp.code != 200) {
				// 								taht.$Message.error({
				// 									content: "状态修失败",
				// 									onClose: () => {
				// 										taht.init();
				// 									}
				// 								});
				// 							}
				// 						});
				// 					}
				// 				}
				// 			},
				// 			[
				// 				h("span", {
				// 					slot: "open",
				// 					domProps: {
				// 						innerHTML: "启用"
				// 					}
				// 				}),
				// 				h("span", {
				// 					slot: "close",
				// 					domProps: {
				// 						innerHTML: "禁用"
				// 					}
				// 				})
				// 			]
				// 		);
				// 	}
				// },
				{
					title: "创建日期",
					width: 150,
					render: (h, data) => {
						return h("div", data.row.ctime.replace(/T|\+08:00/g, " "));
					}
				},
				{
					title: "操作",
					key: "action",
					width: 120,
					align: "center",
					render: (h, data) => {
						return h("a", [
							h(
								"Poptip",
								{
									props: { confirm: true, title: "确定要重置密码吗？" },
									on: {
										"on-ok": () => {
											admUserReset(data.row.id).then(resp => {
												if (resp.code != 200) {
													this.$Message.error(resp.msg);
												}
											});
										}
									}
								},
								[
									h("Icon", {
										props: { type: "md-refresh", size: "20", color: "#2F4056" },
										attrs: { title: "重置密码" },
										style: { marginRight: "10px" }
									})
								]
							),
							h("Icon", {
								props: { type: "md-create", size: "20", color: "#FFB800" },
								attrs: { title: "修改" },
								style: { marginRight: "10px" },
								on: {
									click: () => {
										this.$router.push({
											name: "user-edit",
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
		init() {
			admUserPage(this.page).then(resp => {
				if (resp.code == 200) {
					this.tabData = resp.data.items;
					this.page.total = resp.data.count;
				} else {
					this.tabData = [];
					this.page.total = 0;
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
		emitDrop(data) {
			admUserDrop(data.row.id).then(resp => {
				if (resp.code == 200) {
					this.$Message.success({
						content: "删除成功",
						onClose: () => {
							this.tabData.splice(data.index, 1);
							this.page.total--;
						}
					});
				} else {
					this.$Message.error(resp.msg);
				}
			});
		}
	},
	mounted() {
		this.init();
	}
};
</script>