<style lang="less">
</style>
<template>
	<Card dis-hover>
		<!-- <div slot="title">
			<Icon type="md-person" /> 用户管理
		</div> -->
		<Form inline>
			<!-- <FormItem>
				<Select v-model="cid" placeholder="请选择文章类别" style="width:200px">
					<Option v-for="item in columnAll" :value="item.id" :key="item.id">{{ item.name }}</Option>
				</Select>
			</FormItem> -->
			<FormItem>
				<Button :to="{name:'user-add'}" style="margin-right: 8px">添加用户</Button>
				<Button type="info" @click="init" icon="md-refresh" title="刷新数据">刷&nbsp;&nbsp;新</Button>
			</FormItem>
			<!-- <FormItem>
				<Button type="warning" @click="reCache" icon="ios-alert-outline" :title="'重新加载数据库缓存'+'\n'+'适用于直接修改数据库'">重载缓存</Button>
			</FormItem> -->
		</Form>
		<Table size="small" border :columns="tabCol" :data="tabData"></Table>
		<br>
		<Page :total="total" :current.sync="page.pi" :page-size="page.ps" :page-size-opts="[10,12,15]" @on-change="piChange" @on-page-size-change="psChange" show-sizer show-elevator show-total></Page>
	</Card>
</template>
<script>
import { admUserPage, admUserDrop, admUserEditReset, admUserEditLock } from "@/api/user";
export default {
	data() {
		return {
			total: 0,
			page: { pi: 1, ps: 12 },
			tabCol: [
				{ type: "index", width: 60, align: "center" },
				{ title: "登录账号", key: "num", width: 200 },
				{ title: "用户名", width: 280, key: "name" },
				{
					title: "用户角色",
					minWidth: 150,
					render: (h, data) => {
						return h("div", data.row.role.name);
					}
				},
				{ title: "用户电话", width: 150, key: "phone" },
				{ title: "用户邮箱", width: 150, key: "email" },
				{
					title: "状态",
					width: 100,
					key: "lock",
					render: (h, data) => {
						const that = this;
						return h(
							"i-switch",
							{
								props: { value: !data.row.lock, size: "large" },
								on: {
									"on-change": function (val) {
										admUserEditLock({ id: data.row.id, lock: val }).then((resp) => {
											if (resp.code != 200) {
												data.row.lock = !val;
												that.$Message.error({
													content: resp.msg,
													duration: 3,
													onClose: () => {
														that.init();
													}
												});
											}
										});
									}
								}
							},
							[
								h("span", { slot: "open", domProps: { innerHTML: "启用" } }),
								h("span", { slot: "close", domProps: { innerHTML: "禁用" } })
							]
						);
					}
				},
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
											admUserEditReset({ id: data.row.id }).then((resp) => {
												if (resp.code != 200) {
													this.$Message.error({ content: resp.msg, duration: 3 });
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
			admUserPage(this.page).then((resp) => {
				if (resp.code == 200) {
					this.tabData = resp.data.items;
					this.total = resp.data.count;
				} else {
					this.tabData = [];
					this.total = 0;
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
			admUserDrop({ id: data.row.id }).then((resp) => {
				if (resp.code == 200) {
					this.$Message.success({
						content: "删除成功",
						onClose: () => {
							this.tabData.splice(data.index, 1);
							this.total--;
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