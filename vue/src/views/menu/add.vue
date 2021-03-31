
<template>
	<Card :bordered="false" dis-hover>
		<p slot="title">
			<Icon type="ios-add-circle-outline" /> 添加分类
		</p>
		<div style="max-width:520px">
			<Form ref="dataForm" :model="dataForm" :label-width="100" label-position="right" :rules="dataRules">
				<Alert type="warning" closable>保存之后,无法修改</Alert>
				<FormItem label="所属分类：" prop="pid">
					<Select v-model="dataForm.pid">
						<template v-for="item in menuAll">
							<Option v-if="item.use" style="white-space: pre;" :value="item.id" :key="item.id">{{ item.name }}</Option>
							<Option v-else style="white-space: pre;" :value="item.id" :key="item.id" disabled>{{ item.name }}</Option>
						</template>
					</Select>
				</FormItem>
				<FormItem label="分类名称：" prop="name">
					<Input v-model="dataForm.name"></Input>
				</FormItem>
				<FormItem label="跳转链接：" prop="url">
					<Input v-model="dataForm.url">
					<Tooltip slot="append" placement="top">
						<Button>说明</Button>
						<div slot="content">
							<p>外部链接请以 http(s) 开头</p>
							<p>内部链接请以 / 开头</p>
							<p>不填写则不跳转</p>
						</div>
					</Tooltip>
					</Input>
				</FormItem>
				<Row>
					<Col span="8">
					<FormItem :label-width="100" label="是否使用：" prop="use" title="是否允许添加新闻">
						<i-Switch size="large" v-model="dataForm.use"><span slot="open">使用</span>
							<span slot="close">禁用</span>
						</i-Switch>
					</FormItem>
					</Col>
					<Col span="8" style="text-align: center">
					<FormItem :label-width="90" label="是否显示：" prop="show" title="导航是否显示">
						<i-Switch size="large" v-model="dataForm.show"><span slot="open">显示</span>
							<span slot="close">隐藏</span>
						</i-Switch>
					</FormItem>
					</Col>
					<Col span="8">
					<FormItem :label-width="90" label="显示顺序：" prop="sort">
						<InputNumber :max="10000" :min="1" v-model="dataForm.sort"></InputNumber>
					</FormItem>
					</Col>
				</Row>
				<FormItem label="首页模板：" prop="itpl">
					<Input v-model="dataForm.itpl"></Input>
				</FormItem>
				<FormItem label="列表模板：" prop="ltpl">
					<Input v-model="dataForm.ltpl"></Input>
				</FormItem>
				<FormItem label="详细模板：" prop="dtpl">
					<Input v-model="dataForm.dtpl"></Input>
				</FormItem>
				<Alert type="warning" closable>图片操作完成请点击 确认上传</Alert>
				<FormItem label="封面图片：">
				</FormItem>
				<FormItem>
					<Button type="warning" :loading="loading" @click="submitAdd">提交保存</Button>
					<Button type="success" @click="resetForm()" style="margin-left: 8px">重置填写</Button>
				</FormItem>
			</Form>
		</div>
	</Card>
</template>
<script>
import { apiMenuTree, admMenuAdd } from "@/api/menu";
import { apiUploadImage } from "@/api/fetch";
export default {
	data() {
		return {
			menuAll: [],
			dataForm: {
				level: 1,
				pid: 0,
				name: "",
				url: "",
				use: true,
				show: true,
				sort: 1000,
				cover: "",
				itpl: "subidx.html",
				ltpl: "list.sml.html",
				dtpl: "detail.html"
			},
			dataRules: {
				pid: [{ type: "number", required: true, message: "请选择分类", trigger: "change", min: 0 }],
				name: [{ required: true, message: "请填写分类名称", trigger: "blur", max: 64 }]
			},
			loading: false
		};
	},
	methods: {
		init() {
			apiMenuTree({ root: true }).then((resp) => {
				if (resp.code == 200) {
					this.menuAll = resp.data;
				}
			});
		},
		// onChange(_, item) {
		// 	console.log(item);
		// 	this.dataForm.pid = item[item.length - 1].id;
		// 	this.dataForm.level = item[item.length - 1].level + 1;
		// },
		resetForm() {
			this.$refs.dataForm.resetFields();
			this.$refs.upimg.resetCrop();
			this.defValue = ["所有类别"];
		},
		submitAdd() {
			this.$refs["dataForm"].validate((valid) => {
				if (valid) {
					this.loading = true;
					admMenuAdd(this.dataForm).then((resp) => {
						this.loading = false;
						if (resp.code == 200) {
							this.$Message.success({
								content: "添加成功",
								onClose: () => {
									// this.$router.push({
									// 	name: "menu-list"
									// });
									this.resetForm();
								}
							});
						} else {
							this.$Message.error({
								content: resp.msg,
								duration: 3
							});
						}
					});
				}
			});
		}
	},
	created() {
		this.init();
	}
};
</script>
