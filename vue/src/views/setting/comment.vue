
<template>
    <Card :bordered="false" dis-hover>
        <p slot="title">
            <Icon type="ios-text-outline" /> 评论设置
        </p>
        <div>
            <Form label-position="top">
                <FormItem label="评论代码">
                    <Input v-model="comment" style="width:600px" type="textarea" :autosize="{minRows: 15,maxRows: 20}" placeholder="Enter code "></Input>
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
			comment: "",
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
