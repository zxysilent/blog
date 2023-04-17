import { ObjectDirective } from 'vue';
import { useAdminStore } from '@/store/admin';
export const auth: ObjectDirective = {
    mounted(el: HTMLElement, binding) {
        if (binding.value == undefined) return;
        const adminStore = useAdminStore();
        if (!adminStore.Authed(binding.value)) {
            el.remove();
        }
    },
};
