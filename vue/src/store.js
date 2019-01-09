const store = {
  state: {
    cachePage: [],
  },
  mutations: {
    logout(state) {
      localStorage.removeItem("bearer");
      localStorage.clear();
    }
  } //,
  // actions: {},
  // modules: {}
};

export default store;
