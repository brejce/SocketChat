<template>
	<u-form style="width: 80%;align-items: center;margin: auto;justify-content: center;margin-top: 200upx;border: 5upx solid #ecf5ff;padding: 20upx;border-radius: 20upx;background-color: #f4f4f5;" >
			<u-form-item label="用户名">
				<u-input v-model="user.name" :border="true" />
			</u-form-item>
		<u-form-item label="密码">
			<u-input v-model="user.passwd" type="password" :passwordIcon="true" :border="true"/>
		</u-form-item>
		<u-form-item >
			<u-row >
				<u-col>
				<u-button @click="register()">
					提交
				</u-button>	
				</u-col>
			</u-row>
		</u-form-item>
	</u-form>
</template>

<script>
	export default {
		data() {
			return {
				user:{name:'',passwd:'',status:'register'},
			}
		},
		methods: {
			register(){
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
							url:getApp().globalData.serverAddr+'register',
							method:'POST',
							data:this.user,
							success(res) {
								switch(res.data){
									case 'yes':
										try{
											uni.setStorageSync('user',this.user)
										}catch(e){
											getApp().globalData.SetLog('register set user',e)
										}
										uni.reLaunch({
											url:'../login/login'
										})
										getApp().globalData.SetLog('register success','登录成功')
									break;
									case 'no':
										uni.showToast({
											title:'用户已存在',
											duration:2000
										})
										getApp().globalData.SetLog('dis login get user fail','用户已存在')
									break;
									default:
										getApp().globalData.SetLog('dis login get user fail','网络开小差')
										uni.showToast({
											title:'网络开小差',
											duration:3000,
										})
									;
									
								}
							},
							fail(res) {
								getApp().globalData.SetLog('register fail',res)
							}
						})
					}
				}
			}
		}
	}
</script>

<style>
.content {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
	}

	.logo {
		height: 200rpx;
		width: 200rpx;
		margin-top: 200rpx;
		margin-left: auto;
		margin-right: auto;
		margin-bottom: 50rpx;
	}

	.text-area {
		display: flex;
		justify-content: center;
	}

	.title {
		font-size: 36rpx;
		color: #8f8f94;
	}
</style>
