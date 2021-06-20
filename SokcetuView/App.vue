<script>
export default {
  globalData: {
	SocketTask:null,
    log: [],
    serverAddr: "http://192.168.1.6:8989/",
    msgMap: [],
    DisLogin() {
      var user = {};
      try {
        const value = uni.getStorageSync("user");
        if (value) {
          user = value;
          uni.request({
            url: getApp().globalData.serverAddr + "dislogin",
            method: "POST",
            data: JSON.stringify(user),
            success: (res) => {
              getApp().globalData.SetLog("dis login success", res);
            },
            fail: (res) => {
              getApp().globalData.SetLog("dis login dail", res);
            },
          });
        }
      } catch (e) {
        getApp().globalData.SetLog("dis login get user fail", e);
      }
    },
    SetLog(id, res) {
      var l = {
        ID: id,
        RES: res,
      };
      getApp().globalData.log.push(l);
    },
    GetAllMsgInit() {
      uni.request({
        url: getApp().globalData.serverAddr + "getallmessage",
        method: "POST",
        success: (res) => {
          getApp().globalData.SetLog("get message init :", "success");
          getApp().globalData.msgMap = res.data;
        },
        fail: (res) => {
          getApp().globalData.SetLog("get message fail", res);
        },
      });
    },
    initSocketTask() {
      getApp().globalData.SocketTask = uni.connectSocket({
        url: getApp().globalData.serverAddr + "message",
        success: (res) => {
          getApp().globalData.SetLog("message init task success", res);
        },
      });
    },
  },
  onLaunch: function () {
    console.log("App Launch");
    getApp().globalData.GetAllMsgInit();
  },
  onShow: function () {
    console.log("App Show");
  },
  onHide: function () {
    console.log("App Hide");
  },
};
</script>

<style lang="scss">
@import "uview-ui/index.scss";
/*每个页面公共css */
</style>
