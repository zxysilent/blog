<template>
	<div>
		<Card dis-hover>
			<div slot="title">
				<Icon type="ios-pricetags-outline" /> 标签管理
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
				<!-- <FormItem>
				<Button type="warning" @click="reCache" icon="ios-alert-outline" :title="'重新加载数据库缓存'+'\n'+'适用于直接修改数据库'">重载缓存</Button>
			</FormItem> -->
			</Form>
			<Table stripe size="small" :columns="colTag" :data="dataTag"></Table>
		</Card>
		<Modal v-model="showEdit" title="修改标签信息">
			<Form ref="editForm" :model="editForm" label-position="top" :rules="editRules">
				<FormItem label="标签名称" prop="name">
					<Input v-model="editForm.name" placeholder="请填写标签名"></Input>
				</FormItem>
				<FormItem label="标签介绍" prop="intro">
					<Input v-model="editForm.intro" placeholder="请填写标签介绍"></Input>
				</FormItem>
			</Form>
			<div slot="footer">
				<ButtonGroup>
					<Button type="warning" :loading="editLoading" @click="emitEdit">提交保存</Button>
					<Button type="info" style="margin-left: 8px" @click="showEdit=false">取消关闭</Button>
				</ButtonGroup>
			</div>
		</Modal>
	</div>
</template>
<script>
import { apiTagAll, admTagEdit, admTagDrop } from "@/api/tag";
export default {
	data() {
		return {
			showEdit: false,
			editLoading: false,
			editForm: { name: "", intro: "" },
			editRules: {
				name: [{ required: true, message: "请填写标签名", trigger: "blur", max: 64 }],
				intro: [{ required: true, message: "请填写标签介绍", trigger: "blur", max: 64 }]
			},
			colTag: [
				{ type: "index", minWidth: 60, maxWidth: 100, align: "center" },
				{ title: "标签名", minWidth: 100, maxWidth: 300, key: "name" },
				{ title: "标签介绍", minWidth: 100, maxWidth: 300, key: "intro" },
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
										this.showEdit = true;
										this.editForm = data.row;
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
			dataTag: []
		};
	},
	methods: {
		init() {
			apiTagAll().then((resp) => {
				if (resp.code == 200) {
					this.dataTag = resp.data;
				} else {
					this.dataTag = [];
					this.$Message.warning("未查询到标签信息,请重试！");
				}
			});
		},
		emitEdit() {
			this.$refs["editForm"].validate((valid) => {
				if (valid) {
					this.editLoading = true;
					admTagEdit(this.editForm).then((resp) => {
						this.editLoading = false;
						if (resp.code == 200) {
							this.$Message.success({
								content: "标签信息修改成功",
								onClose: () => {
									this.showEdit = false;
								}
							});
						} else {
							this.$Message.error({ content: `标签信息修改失败,请重试`, duration: 3 });
						}
					});
				}
			});
		},
		delete(data) {
			admTagDrop(data.row.id).then((resp) => {
				if (resp.code == 200) {
					this.$Message.success({
						content: "删除成功",
						onClose: () => {
							this.dataTag.splice(data.index, 1);
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
