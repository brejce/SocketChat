<template>
	<view>
		{{log}}
	</view>
</template>
<script>
	export default {
		data() {
			return {
				log:[]
			}
		},
		methods: {
	
		},
		onLoad() {
			this.log = JSON.stringify(getApp().globalData.log)
		}
	}
</script>

<style scoped lang="scss">
	.demo-layout {
		height: 70rpx;
		border-radius: 8rpx;
		margin: 20rpx 0;
	}

	.bg-purple {
		background: #d3dce6;
	}

	.bg-purple-light {
		background: #e5e9f2;
	}

	.bg-purple-dark {
		background: #99a9bf;
	}

	// H5中，电脑端文档演示时，可能会导致演示块挤出边界，特别处理。
	// 真实使用环境不会产生此问题
	/* #ifdef H5 */
	.u-demo-area /deep/ .u-row {
		display: flex;
		flex-wrap: wrap;
	}
	/* #endif */
</style>