
<template>
    <Card :bordered="false" dis-hover>
        <p slot="title">
            <Icon type="ios-pulse" /> 统计设置
        </p>
        <div>
            <Form label-position="top">
                <FormItem label="网站统计代码">
                    <Input v-model="analytic" style="width:600px" type="textarea" :autosize="{minRows: 15,maxRows: 20}" placeholder="Enter code "></Input>
                </FormItem>
                <div>
                    <Button type="warning" :loading="saveLoading" @click="cmtSave">保&nbsp;&nbsp;&nbsp;&nbsp;存</Button>
                </div>
            </Form>
        </div>
    </Card>
</template>
<script>
export default {
	data() {
		return {
			analytic: "",
			saveLoading: false
		};
	},
	methods: {
		cmtSave() {
			this.saveLoading = true;
			setTimeout(() => {
				this.saveLoading = false;
			}, 1500);
		},
		init() {}
	},
	mounted() {
		this.init();
	}
};
</script>
