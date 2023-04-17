<template>
    <div class="page-footer">
        <div class="page-footer-link">
            <a href="https://www.cwnu.edu.cn/" target="_blank"> 学校首页 </a>
            <a href="https://nic.cwnu.edu.cn/" target="_blank"> 网络中心 </a>
            <a href="https://github.com/zxysilent" target="_blank"> 技术支持 </a>
        </div>
        <div class="copyright">Copyright &copy;网络与信息化管理中心. 技术支持<a href="https://github.com/zxysilent" target="_blank"> zxysilent </a></div>
    </div>
</template>
<style lang="less" scoped>
.page-footer {
    //margin: 28px 0 24px 0;
    padding: 0 16px;
    text-align: center;
    a {
        font-size: 14px;
        color: #808695;
        -webkit-transition: all 0.2s ease-in-out;
        transition: all 0.2s ease-in-out;
        &:hover {
            color: #515a6e;
        }
    }
    &-link {
        display: flex;
        justify-content: center;
        margin-bottom: 8px;
        a:not(:last-child) {
            margin-right: 40px;
        }
    }
    color: #808695;
    font-size: 14px;
}
</style>
