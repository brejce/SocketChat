<template>
	<view style="background-color:  #ecf5ff;">
		<u-navbar :is-back="false" :title="user.name">
			<view class="slot-wrap" slot="right">
				<navigator url="../setting/setting">设置</navigator>
			</view>
		</u-navbar>
		<view class="content">
			<view>
				<u-card v-for="(msg,index) in list" :title="msg.name">
					<view slot="body">{{msg.msge}}</view>
				</u-card>
			</view>
			<view class="inputBox">
				<u-row gutter="24" justify="start" style="width: 100%;margin: 20upx;">
					<u-col span="8">
						<u-input v-model="msg" :border="true" ></u-input>
					</u-col>
					<u-col span="1">
						<u-button @click="sendGo()" size="medium">发送</u-button>
					</u-col>
				</u-row>
			</view>
		</view>			
	</view>
</template>

<script>
	export default {
		data() {
			return {
				scroll:0,
				user:null,
				list:[],
				Task:null,
				msg:'',
			}
		},
		methods: {
			sendGo(){
				if(''!=this.msg){
					let d = new Date().getTime().toString()
					var msgData ={
						name:this.user.name,
						idtime:d,
						msge:this.msg
					}
					this.Task.send({
						data:JSON.stringify(msgData),
						success:(res)=>{
							this.msg = ''
						},
						fail:(res)=>{
							getApp().globalData.SetLog('message send msg fail',res)
						}
					});
				}else{
					uni.showToast({
						title:'消息不能为空',
						duration:2000
					})
				}	
			},
			GetLocalUser(){
					try {
						const value = uni.getStorageSync('user')
						if (value) {						
							this.user = value
							// uni.setNavigationBarTitle({
							// 	title:value.name
							// });
						}
					} catch (e) {
						// error
						console.log(e);
						getApp().globalData.SetLog('message get user fail',)
					}
			},
			GetLocalList(){
				this.list = getApp().globalData.msgMap
			},
			initTask(){
				getApp().globalData.initSocketTask()
				this.Task = getApp().globalData.SocketTask
				this.Task.onError((res)=>{
					getApp().globalData.SetLog('message init task onErro',res)
				});
				this.Task.onClose((res)=>{
					getApp().globalData.SetLog('message init task onClose',res)
				});
				this.Task.onOpen((res)=>{
					getApp().globalData.SetLog('message init task onOpen',res)
				});
				this.Task.onMessage((res)=>{
					getApp().globalData.SetLog('message init task onMessage','成功')
					if(null == this.list){
						var s = [JSON.parse(res.data)]
						this.list = s
					}else{
						this.list.push(JSON.parse(res.data))
					}
				});
			}
		},
		onShow() {
			if(this.user == null){
				this.GetLocalUser()
			}
			if(this.list == null){
				this.GetLocalList()
			}
		},
		beforeDestroy() {
			this.Task.close();
		},
		onLoad() {
			this.GetLocalUser()
			this.GetLocalList()
			this.initTask()
		}
	}
</script>

<style scoped lang="scss">
.inputBox{
	position: fixed;
	bottom: 0upx;
	left: 0upx;
	right: 0upx;
	width: 100%;
	height: 200upx;
	z-index: 99;
	background: #FFFFFF;
}
.slot-wrap {
		display: flex;
		align-items: center;
		/* 如果您想让slot内容占满整个导航栏的宽度 */
		/* flex: 1; */
		/* 如果您想让slot内容与导航栏左右有空隙 */
		padding: 0 30rpx; 
	}
.content{
	margin-top: 0upx;
	// margin-left: 50upx;
	// margin-right: 50upx;
	left: 20upx;
	right: 20upx;
	margin-bottom: 200upx;
	// background: #ecf5ff;
}
</style>
