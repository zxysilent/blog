const store = {
  state: {
    cachePage: [],
    messageCount: 0
  },
  mutations: {
    logout(state) {
      localStorage.removeItem("bearer");
      localStorage.clear();
    },
    setMessageCount(state, count) {
      state.messageCount = count;
    }
  } //,
  // actions: {},
  // modules: {}
};

export default store;
