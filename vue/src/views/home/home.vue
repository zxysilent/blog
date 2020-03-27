<style lang="less" scoped>
@import "./home.less";
</style>
<template>
	<div class="home">
		<Row>
			<Col :md='{span:10}'>
			<Row>
				<Col span="12" class="collect">
				<Card dis-hover>
					<Row>
						<Col span="10">
						<Icon type="ios-megaphone"></Icon>
						</Col>
						<Col span="14">
						<p class="title">文章</p>
						<h3>{{collect.post}}</h3>
						</Col>
					</Row>
				</Card>
				</Col>
				<Col span="12" class="collect">
				<Card dis-hover>
					<Row>
						<Col span="10">
						<Icon type="ios-map"></Icon>
						</Col>
						<Col span="14">
						<p class="title">页面</p>
						<h3>{{collect.page}}</h3>
						</Col>
					</Row>
				</Card>
				</Col>
				<Col span="12" class="collect">
				<Card dis-hover>
					<Row>
						<Col span="10">
						<Icon type="ios-school"></Icon>
						</Col>
						<Col span="14">
						<p class="title">分类</p>
						<h3>{{collect.cate}}</h3>
						</Col>
					</Row>
				</Card>
				</Col>
				<Col span="12" class="collect">
				<Card dis-hover>
					<Row>
						<Col span="10">
						<Icon type="ios-pricetags"></Icon>
						</Col>
						<Col span="14">
						<p class="title">标签</p>
						<h3>{{collect.tag}}</h3>
						</Col>
					</Row>
				</Card>
				</Col>
			</Row>
			</Col>
			<Col :md='{span:14}'>
			<Card dis-hover style="height: 158px">
				<p slot="title">
					系统信息
				</p>
				<Steps status="finish" :current="4">
					<Step icon="ios-flash" title="系统架构" :content="sys.arch"></Step>
					<Step icon="logo-snapchat" title="操作系统" :content="sys.os"></Step>
					<Step icon="ios-thunderstorm" title="运行环境" :content="sys.version"></Step>
					<Step icon="ios-flower" title="逻辑处理器" :content="sys.num_cpu+''"></Step>
				</Steps>
			</Card>
			</Col>
		</Row>
		<Row>
			<Col :md='{span:8}'>
			<Card dis-hover style="height: 320px">
				<p slot="title">
					暂时未统计，正在进行
				</p>
				<a href="#" slot="extra" @click.prevent="refresh">
					<Icon type="ios-loop-strong"></Icon>
				</a>
				<Steps :current="2" direction="vertical" size="small">
					<Step title="已完成" content="这里是该步骤的描述信息"></Step>
					<Step title="已完成" content="这里是该步骤的描述信息"></Step>
					<Step title="进行中" content="这里是该步骤的描述信息"></Step>
					<Step title="待进行" content="这里是该步骤的描述信息"></Step>
				</Steps>
			</Card>
			</Col>
		</Row>
	</div>
</template>
<script>
import { admSys, admCollect } from "@/api/auth";
export default {
	name: "index",
	data() {
		return {
			sys: {
				arch: "--",
				num_cpu: 0,
				os: "--",
				version: "--"
			},
			collect: {
				post: 0,
				page: 0,
				cate: 0,
				tag: 0
			}
		};
	},
	methods: {
		init() {
			admSys().then(resp => {
				if (resp.code == 200) {
					this.sys = resp.data;
				} else {
					this.$Message.warning("未查询到系统信息,请重试！");
				}
			});
			admCollect().then(resp => {
				if (resp.code == 200) {
					this.collect = resp.data;
				} else {
					this.$Message.warning("未查询到统计信息,请重试！");
				}
			});
		}
	},
	mounted() {
		this.init();
	}
};
</script>