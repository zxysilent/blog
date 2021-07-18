// 权限按钮展示指令
export const auth = Vue => {
	Vue.directive("auth", {
		// 当被绑定的元素插入到 DOM 中时……
		update: function(el, binding, vNode) {
			const { $root: vm } = vNode.context;
			// 获取当前用户拥有的权限列表
			const access = vm.$store.getters.AuthGrant;
			let { value, modifiers } = binding;
			let type = "";
			switch (Object.prototype.toString.call(value)) {
				case "[object Array]":
					type = "Array";
					break;
				case "[object String]":
					type = "String";
					break;
				case "[object Number]":
					type = "Number";
					break;
				default:
					type = "";
					break;
			}
			if (type === "") {
				/* eslint-disable */
				console.error("v-auth必须是Array,Number,String属性,暂不支持其他属性");
				/* eslint-enable */
				return;
			}
			// 需要的权限
			const need = binding.value.toString().split(",");
			// console.log(need, access);
			if (!need.some(v => access.includes(v))) {
				try {
					el.parentNode.removeChild(el);
				} catch (error) {
					el.style.display = "none";
				}
			}
			console.log("v-auth");
		}
	});
};
