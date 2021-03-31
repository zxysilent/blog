
<template>
	<Card :bordered="false" dis-hover>
		<p slot="title">
			<Icon type="ios-create-outline" /> 修改分类
		</p>
		<div style="max-width:520px">
			<Form ref="dataForm" :model="dataForm" :label-width="100" label-position="right" :rules="dataRules">
				<Alert type="warning">部分信息无法修改</Alert>
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
					<FormItem :label-width="100" label="是否使用：" prop="phone">
						<i-Switch size="large" v-model="dataForm.use"><span slot="open">使用</span>
							<span slot="close">禁用</span>
						</i-Switch>
					</FormItem>
					</Col>
					<Col span="8" style="text-align: center">
					<FormItem :label-width="90" label="是否显示：" prop="keywords">
						<i-Switch size="large" v-model="dataForm.show"><span slot="open">显示</span>
							<span slot="close">隐藏</span>
						</i-Switch>
					</FormItem>
					</Col>
					<Col span="8">
					<FormItem :label-width="90" label="显示顺序：" prop="keywords">
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
				<FormItem>
					<Button type="warning" :loading="loading" @click="submitEdit">提交保存</Button>
					<Button type="success" @click="resetForm()" style="margin-left: 8px">重置填写</Button>
					<Button :to="{name:'menu-list'}" style="margin-left: 8px">返回列表</Button>
				</FormItem>
			</Form>
		</div>
	</Card>
</template>
<script>
import { apiMenuGet, admMenuEdit } from "@/api/menu";
import { apiUploadImage } from "@/api/fetch";
export default {
	data() {
		return {
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
				name: [{ required: true, message: "请填写分类名称", trigger: "blur", max: 64 }]
			},
			loading: false
		};
	},
	methods: {
		onCrop(blob) {
			apiUploadImage(blob, {}).then(resp => {
				if (resp.code == 200) {
					this.dataForm.cover = resp.data;
				}
			});
		},
		onReset() {
			this.dataForm.cover = "";
		},
		init() {
			apiMenuGet(this.dataForm.id).then(resp => {
				if (resp.code == 200) {
					this.dataForm = resp.data;
				}
			});
		},
		resetForm() {
			this.init();
		},
		submitEdit() {
			this.$refs["dataForm"].validate(valid => {
				if (valid) {
					this.loading = true;
					admMenuEdit(this.dataForm).then(resp => {
						this.loading = false;
						if (resp.code == 200) {
							this.$Message.success({
								content: "修改成功",
								onClose: () => {
									this.init();
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
		this.dataForm.id = parseInt(this.$route.params.id, 10);
		this.init();
	}
};
</script>
