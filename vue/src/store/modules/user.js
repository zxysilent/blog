
const user = {
  state: {},
  mutations: {
    logout(state, vm) {
        localStorage.removeItem('bearer');
    }
  }
};

export default user;
