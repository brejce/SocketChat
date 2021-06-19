<template>
<!-- 	<view class="content">
		<text>Change Password</text>
		<text> 请记住您的用户名与密码</text>
		<text>如忘记密码请联系管理员</text>
		<text>用户名忘记将无法找回</text>
		
		<input placeholder="请输入用户名" v-model="name" />
		<input placeholder="请输入旧密码" v-model="passwd" />
		<input placeholder="请输入新密码" v-model="repasswd" />
		<button @click="ChangePasswd()">确认</button>
	</view> -->
	<u-form style="width: 80%;align-items: center;margin: auto;justify-content: center;margin-top: 200upx;border: 5upx solid #ecf5ff;padding: 20upx;border-radius: 20upx;background-color: #f4f4f5;" >
		<u-form-item label="用户名">
			<u-input v-model="name" :border="true" />
		</u-form-item>
		<u-form-item label="旧密码">
			<u-input v-model="passwd" type="password" :passwordIcon="true" :border="true"/>
		</u-form-item>
		<u-form-item label="新密码">
			<u-input v-model="repasswd" type="password" :passwordIcon="true" :border="true"/>
		</u-form-item>
		<u-form-item style="justify-content: center;">
		
				<u-button @click="ChangePasswd()">
					确认
				</u-button>	
			
		</u-form-item>
	</u-form>
</template>

<script>
	export default {
		data() {
			return {
				name:'',
				passwd:'',
				repasswd:''
			}
		},
		methods: {
			ChangePasswd(){
				var user = {
					name:this.name,
					passwd:this.passwd,
					status:this.repasswd
				}
				if('' == user.name){
					uni.showToast({
						title:'用户名不能为空',
						duration:2000
					});
				}else{
					if(''==user.passwd){
						uni.showToast({
							title:'旧密码不能为空',
							duration:2000
						});
					}else{
						if('' == user.status){
							uni.showToast({
								title:'新密码不能为空',
								duration:2000
							})
						}else{
							uni.request({
									url:getApp().globalData.serverAddr+'chanagepasswd',
									method:'POST',
									data:user,
									success:(res) =>{
										switch(res.data){
											case 'yes':
											var user = {
												name:this.name,
												passwd:this.repasswd,
												status:'chanage'
											}
											try{
												uni.setStorageSync('user',user)
											}catch(e){
												getApp().globalData.SetLog('login get user fail',e)
											}
											uni.reLaunch({
												url:'../login/login'
											});
											uni.showToast({
												title: '修改成功!',
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
											getApp().globalData.SetLog('login fail','default')
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
				}
			}
		}
	}
</script>

<style>

</style>
