<template>
  <error code="404" desc="Oh~~您的页面好像飞走了~" :src="src"/>
</template>
<script>
import error404 from '@/assets/404.svg'
import error from './error.vue'
export default {
  name: 'err404',
  components: {
    error
  },
  data () {
    return {
      src: error404
    }
  }
}
</script>
