<template>
	<u-form style="width: 80%;align-items: center;margin: auto;justify-content: center;margin-top: 200upx;border: 5upx solid #ecf5ff;padding: 20upx;border-radius: 20upx;background-color: #f4f4f5;" >
			<u-form-item label="用户名">
				<u-input v-model="user.name" :border="true" />
			</u-form-item>
		<u-form-item label="密码">
			<u-input v-model="user.passwd" type="password" :passwordIcon="true" :border="true"/>
		</u-form-item>
		<u-form-item style="justify-content: center;">
			<u-button @click="Login()">
			提交
			</u-button>	
			<u-button style="font-size: 30upx;" @click="Relaun()">
			没有账号?
			</u-button>
		</u-form-item>
	</u-form>	
</template>

<script>
	export default {
		data() {
			return {
				user:{name:'',passwd:'',status:''},
				haveUser:false
			}
		},
		onLoad() {
			if(!this.haveUser){
				try {
					const value = uni.getStorageSync('user');
					if (value) {
						this.haveUser = true
						value.status = 'login'
						this.user = value
					}
				} catch (e) {
					getApp().globalData.SetLog('init Login get user fail',e);
				}
			}
		},
		methods: {
			Changeassword(){
				uni.navigateTo({
					url:'../Changeasswordp/Changeasswordp'
				})
			},
			Relaun(){
				uni.navigateTo({
					url:'../register/register'
				});
			},
			Login(){
				if('' == this.user.name){
					uni.showToast({
						title:'用户名不能为空',
						duration:2000
					});
				}else{
					if(''==this.user.passwd){
						uni.showToast({
							title:'密码不能为空',
							duration:2000
						});
					}else{
						uni.request({
							url:getApp().globalData.serverAddr+'login',
							method:'POST',
							data:this.user,
							success:(res) =>{
								switch(res.data){
									case 'yes':
									
									try{
										uni.setStorageSync('user',this.user)
									}catch(e){
										getApp().globalData.SetLog('login get user fail',e)
									}
									uni.reLaunch({
										url:'../message/message'
									});
									uni.showToast({
										title: '登录成功!',
										duration: 2000
									});
									getApp().globalData.SetLog('login success','yes')
									break;
									case 'no':
									uni.showToast({
										title: '信息有误!',
										duration: 2000
									});
									getApp().globalData.SetLog('login fail','Wrong user name or password ')
									break;
									default:
									getApp().globalData.SetLog('login fail',res)
									uni.showToast({
										title: '网络开小差了!',
										duration: 2000
									});
								}
							},
							fail(res) {
								getApp().globalData.SetLog('login request fail',res)
							}
						});
					}
				}
				
			},
		}
	}
</script>

<style scoped lang="scss">
.wrap {
		padding: 24rpx;
	}

	.u-row {
		margin: 40rpx 0;
	}

	.demo-layout {
		height: 80rpx;
		border-radius: 8rpx;
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
</style>
