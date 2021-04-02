
<template>
	<Card :bordered="false" dis-hover>
		<p slot="title">
			<Icon type="ios-create-outline" /> 修改菜单
		</p>
		<div style="max-width:520px">
			<Form ref="dataForm" :model="dataForm" :label-width="100" label-position="right" :rules="dataRules">
				<!-- <Alert type="warning" closable>保存之后,无法修改</Alert> -->
				<FormItem label="所属菜单：" prop="pid">
					<Select v-model="dataForm.pid" placeholder="请选择所属菜单">
						<template v-for="item in menuAll">
							<Option v-if="item.pid==0" style="white-space: pre;" :value="item.id" :key="item.id">{{ "「"+item.name+"」"+item.title }}</Option>
						</template>
					</Select>
				</FormItem>
				<FormItem label="菜单标题：" prop="title">
					<Input v-model="dataForm.title" placeholder="请填写菜单标题"></Input>
				</FormItem>
				<FormItem label="菜单名称：" prop="name">
					<Input v-model="dataForm.name" placeholder="请填写菜单名称"></Input>
				</FormItem>
				<FormItem label="菜单图标：" prop="icon">
					<Select v-model="dataForm.icon" placeholder="请选择菜单图标" filterable>
						<Icon :type="dataForm.icon" slot="prefix" size="22" />
						<Option v-for="item in icons" :value="item" :label="item" :key="item.id">
							<span>
								<Icon :type="item" />
							</span>
							<span style="margin-left:10px;">{{item}}</span>
						</Option>
					</Select>
				</FormItem>
				<FormItem label="菜单路径：" prop="path">
					<Input v-model="dataForm.path" placeholder="请填写菜单路径"></Input>
				</FormItem>
				<!-- <FormItem label="菜单路由：" prop="ptah">
					<Input v-model="dataForm.ptah"></Input>
				</FormItem> -->

				<!-- <FormItem label="跳转链接：" prop="url">
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
				</FormItem> -->
				<FormItem label="菜单组件：" prop="comp">
					<Input v-model="dataForm.comp" placeholder="请填写菜单组件"></Input>
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

				<FormItem>
					<Button type="warning" :loading="loading" @click="emitEdit">提交保存</Button>
					<Button type="success" @click="emitReset()" style="margin-left: 8px">重置填写</Button>
				</FormItem>
			</Form>
		</div>
	</Card>
</template>
<script>
import { admMenuAll, admMenuGet, admMenuEdit } from "@/api/menu";
import { icons } from "@/utils/icons";
export default {
	data() {
		return {
			icons: icons,
			menuAll: [],
			dataForm: {
				id: 0,
				pid: 0,
				title: "",
				name: "",
				use: true,
				icon: "",
				show: true,
				comp: "",
				sort: 1000,
				inner: false
			},
			dataRules: {
				pid: [{ type: "number", required: true, message: "请选择菜单", trigger: "change", min: 0 }],
				title: [{ required: true, message: "请填写菜单标题", trigger: "blur", max: 128 }],
				name: [{ required: true, message: "请填写菜单名称", trigger: "blur", max: 128 }],
				path: [{ required: true, message: "请填写菜单路径", trigger: "blur", max: 128 }],
				icon: [{ required: true, message: "请选择菜单图标", trigger: "change", max: 128 }],
				comp: [{ required: true, message: "请填写菜单组件", trigger: "blur", max: 128 }]
			},
			loading: false
		};
	},
	methods: {
		preinit() {
			admMenuAll({ root: true, slt: true }).then((resp) => {
				if (resp.code == 200) {
					this.menuAll = resp.data;
				} else {
					this.$Message.error({
						content: resp.msg,
						duration: 3
					});
				}
			});
		},
		init() {
			admMenuGet({ id: this.dataForm.id }).then((resp) => {
				if (resp.code == 200) {
					this.dataForm = resp.data;
				} else {
					this.$Message.error({
						content: resp.msg,
						duration: 3
					});
				}
			});
		},
		emitReset() {
			this.init();
		},
		emitEdit() {
			this.$refs["dataForm"].validate((valid) => {
				if (valid) {
					this.loading = true;
					admMenuEdit(this.dataForm).then((resp) => {
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
		console.log(this.$route.params.id);
		this.dataForm.id = parseInt(this.$route.params.id, 10);
		this.preinit();
		console.log(this.$route);
		this.init();
	}
};
</script>
