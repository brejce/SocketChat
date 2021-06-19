<template>
	<view>
		<view class="container u-skeleton">
			<view class="userinfo">
				<block>
					<!--u-skeleton-circle 绘制圆形-->
					<image class="userinfo-avatar u-skeleton-circle">
						<u-icon name="eye" color="#ecf5ff" size="200"></u-icon>
					</image>
				
					<text class="u-skeleton-fillet">{{user.name}}</text>
					<text class="u-skeleton-fillet">
						{{motto.creator}}:{{motto.hitokoto}}
					</text>
				</block>
			</view>
		</view>
		<u-skeleton :loading="loading" :animation="true" bgColor="#FFF"></u-skeleton>
		<u-row>
			<u-col class="butt">
				<u-button style="margin:20upx 120upx;" type="primary" @click="dislogin()">退出登录</u-button>
				<u-button style="margin:20upx 120upx;" type="warning" @click="toChange()">修改密码</u-button>
				<u-button style="margin:20upx 120upx;" type="error" @click="clear()">重置APP</u-button>
				<u-button style="margin:20upx 120upx;" type="info" @click="showlog()">显示日志</u-button>
			</u-col>
		</u-row>
	</view>
</template>
<script>
	export default {
		data() {
			return {
				user:[],
				loading:true,
				motto:{}
			}
		},
		methods: {
			showlog(){
				uni.navigateTo({
					url:'../test/test'
				})
			},
			toChange(){
				uni.navigateTo({
					url:'../Changeasswordp/Changeasswordp'
				})
			},
			clear(){
				try {
				    uni.clearStorageSync()
					uni.reLaunch({
						url:'../login/login'
					})	
				} catch (e) {
					getApp().globalData.SetLog('重置出错',e)
				}
			},
			dislogin(){
				if( null !=getApp().globalData.SocketTask ){
					getApp().globalData.SocketTask.close();
				}
				getApp().globalData.DisLogin()
				uni.reLaunch({
					url:'../login/login'
				})
			},
			getmotto(){
				uni.request({
					url:'https://sslapi.hitokoto.cn/?encode=json',
					success: (res) => {
						// console.log(res.data)
						this.motto = res.data
					},fail: (res) => {
						console.log('getfail')
						getApp().globalData.SetLog('get motto:',res)
					}
				})
			}
		},
		onLoad() {
			try {
				const value = uni.getStorageSync('user')
				if (value) {						
					this.user = value
				}
			} catch (e) {
				// error
				console.log(e);
				getApp().globalData.SetLog('setting get user fail',)
			}
			this.getmotto()
			setTimeout(() => {
				this.loading = false;
			}, 1000)
		}
	}
</script>

<style lang="scss" scoped>
	.container {
		padding: 20rpx 60rpx;
		// align-items: center;
		align-content: center;
	}

	.userinfo {
		display: flex;
		flex-direction: column;
		align-items: center;
	}

	.userinfo-avatar {
		width: 128rpx;
		height: 128rpx;
		margin: 20rpx;
		border-radius: 50%;
	}

	.lists {
		margin: 10px 0;
	}
	.butt{
		margin: 30upx;
	}
</style>