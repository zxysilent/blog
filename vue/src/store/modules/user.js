import Cookies from "js-cookie";

const user = {
  state: {},
  mutations: {
    logout(state, vm) {
      Cookies.remove("user");
      Cookies.remove("password");
      Cookies.remove("access");
      // 清空打开的页面等数据，但是保存主题数据
    }
  }
};

export default user;
