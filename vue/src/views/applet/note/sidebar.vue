<template>
    <div class="note-layout-sider note-layout-sider--show-content layout-sider" style="width: 64px">
        <div role="menu" class="note-menu">
            <n-tooltip trigger="hover" placement="right">
                <template #trigger>
                    <div role="menuitem" class="note-menu-item">
                        <div class="note-menu-item-content note-menu-item-content--selected" style="padding-left: 21px">
                            <div class="note-menu-item-content__icon" style="width: 22px; height: 22px; font-size: 22px; margin-right: 19px">
                                <i role="img" class="note-icon" style="color: #fff;">
                                    <svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 12 12">
                                        <g fill="none"><path d="M6.5 1.75a.75.75 0 0 0-1.5 0V5H1.75a.75.75 0 0 0 0 1.5H5v3.25a.75.75 0 0 0 1.5 0V6.5h3.25a.75.75 0 0 0 0-1.5H6.5V1.75z" fill="currentColor"></path></g>
                                    </svg>
                                </i>
                            </div>
                        </div>
                    </div>
                </template>
                新建
            </n-tooltip>
            <n-tooltip trigger="hover" placement="right">
                <template #trigger>
                    <div role="menuitem" class="note-menu-item">
                        <div class="note-menu-item-content" style="padding-left: 21px">
                            <div class="note-menu-item-content__icon" style="width: 22px; height: 22px; font-size: 22px; margin-right: 19px">
                                <i role="img" class="note-icon">
                                    <svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 512 512">
                                        <rect x="48" y="48" width="176" height="176" rx="20" ry="20" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="32"></rect>
                                        <rect x="288" y="48" width="176" height="176" rx="20" ry="20" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="32"></rect>
                                        <rect x="48" y="288" width="176" height="176" rx="20" ry="20" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="32"></rect>
                                        <rect x="288" y="288" width="176" height="176" rx="20" ry="20" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="32"></rect>
                                    </svg>
                                </i>
                            </div>
                        </div>
                    </div>
                </template>
                最新
            </n-tooltip>
            <n-tooltip trigger="hover" placement="right">
                <template #trigger>
                    <div role="menuitem" class="note-menu-item">
                        <div class="note-menu-item-content" style="padding-left: 21px">
                            <div class="note-menu-item-content__icon" style="width: 22px; height: 22px; font-size: 22px; margin-right: 19px">
                                <i role="img" class="note-icon">
                                    <svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 512 512">
                                        <path
                                            d="M440 432H72a40 40 0 0 1-40-40V120a40 40 0 0 1 40-40h75.89a40 40 0 0 1 22.19 6.72l27.84 18.56a40 40 0 0 0 22.19 6.72H440a40 40 0 0 1 40 40v240a40 40 0 0 1-40 40z"
                                            fill="none"
                                            stroke="currentColor"
                                            stroke-linecap="round"
                                            stroke-linejoin="round"
                                            stroke-width="32"></path>
                                        <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="32" d="M32 192h448"></path>
                                    </svg>
                                </i>
                            </div>
                        </div>
                    </div>
                </template>
                目录
            </n-tooltip>
            <n-tooltip trigger="hover" placement="right">
                <template #trigger>
                    <div role="menuitem" class="note-menu-item">
                        <div class="note-menu-item-content" style="padding-left: 21px">
                            <div class="note-menu-item-content__icon" style="width: 22px; height: 22px; font-size: 22px; margin-right: 19px">
                                <i role="img" class="note-icon">
                                    <svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 512 512">
                                        <path
                                            d="M403.29 32H280.36a14.46 14.46 0 0 0-10.2 4.2L24.4 281.9a28.85 28.85 0 0 0 0 40.7l117 117a28.86 28.86 0 0 0 40.71 0L427.8 194a14.46 14.46 0 0 0 4.2-10.2v-123A28.66 28.66 0 0 0 403.29 32z"
                                            fill="none"
                                            stroke="currentColor"
                                            stroke-linecap="round"
                                            stroke-linejoin="round"
                                            stroke-width="32"></path>
                                        <path d="M352 144a32 32 0 1 1 32-32a32 32 0 0 1-32 32z" fill="currentColor"></path>
                                        <path d="M230 480l262-262a13.81 13.81 0 0 0 4-10V80" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="32"></path>
                                    </svg>
                                </i>
                            </div>
                        </div>
                    </div>
                </template>
                标签
            </n-tooltip>
            <n-tooltip trigger="hover" placement="right">
                <template #trigger>
                    <div role="menuitem" class="note-menu-item">
                        <a href="#/post/index" class="">
                            <div class="note-menu-item-content" style="padding-left: 21px">
                                <div class="note-menu-item-content__icon" style="width: 22px; height: 22px; font-size: 22px; margin-right: 19px">
                                    <i role="img" class="note-icon">
                                        <svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 512 512">
                                            <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="32" d="M249.38 336L170 256l79.38-80"></path>
                                            <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="32" d="M181.03 256H342"></path>
                                            <path d="M448 256c0-106-86-192-192-192S64 150 64 256s86 192 192 192s192-86 192-192z" fill="none" stroke="currentColor" stroke-miterlimit="10" stroke-width="32"></path>
                                        </svg>
                                    </i>
                                </div></div
                        ></a>
                    </div>
                </template>
                返回
            </n-tooltip>
        </div>
    </div>
</template>
<style lang="less">
.layout-sider {
    min-height: 100vh;
    -webkit-box-shadow: 2px 0 8px 0 rgb(29 35 41 / 5%);
    box-shadow: 2px 0 8px 0 rgb(29 35 41 / 5%);
    position: relative;
    z-index: 13;
    -webkit-transition: all 0.2s ease-in-out;
    transition: all 0.2s ease-in-out;
}

.note-layout-sider {
    flex-shrink: 0;
    box-sizing: border-box;
    position: relative;
    z-index: 1;
    color: #fff;
    transition: color 0.3s cubic-bezier(0.4, 0, 0.2, 1), border-color 0.3s cubic-bezier(0.4, 0, 0.2, 1), min-width 0.3s cubic-bezier(0.4, 0, 0.2, 1), max-width 0.3s cubic-bezier(0.4, 0, 0.2, 1),
        transform 0.3s cubic-bezier(0.4, 0, 0.2, 1), background-color 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    background-color: rgb(0, 20, 40);
    display: flex;
    justify-content: flex-end;
}
.note-menu {
    background-color: rgb(0, 20, 40);
    color: #bbb;
    overflow: hidden;
    transition: background-color 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    box-sizing: border-box;
    font-size: 14px;
    padding-bottom: 6px;
}
.note-menu .note-menu-item {
    height: 42px;
    margin-top: 6px;
    position: relative;
}
.note-menu .note-menu-item-content {
    box-sizing: border-box;
    line-height: 1.75;
    height: 100%;
    display: grid;
    grid-template-areas: "icon content arrow";
    grid-template-columns: auto 1fr auto;
    align-items: center;
    cursor: pointer;
    position: relative;
    padding-right: 18px;
    transition: background-color 0.3s cubic-bezier(0.4, 0, 0.2, 1), padding-left 0.3s cubic-bezier(0.4, 0, 0.2, 1), border-color 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}
.note-menu.note-menu--collapsed .note-menu-item-content .note-menu-item-content__icon {
    color: #bbb;
}
.note-menu .note-menu-item-content .note-menu-item-content__icon {
    grid-area: icon;
    color: #bbb;
    transition: color 0.3s cubic-bezier(0.4, 0, 0.2, 1), font-size 0.3s cubic-bezier(0.4, 0, 0.2, 1), margin-right 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    box-sizing: content-box;
    display: inline-flex;
    align-items: center;
    justify-content: center;
}
.note-menu .note-menu-item-content > * {
    z-index: 1;
}
body .note-icon {
    vertical-align: middle;
}
.note-icon {
    height: 1em;
    width: 1em;
    line-height: 1em;
    text-align: center;
    display: inline-block;
    position: relative;
    fill: currentColor;
    transform: translateZ(0);
}
.note-icon:hover {
    color: #fff;
}
.note-menu .note-menu-item-content.note-menu-item-content--selected::before {
    background-color: #2d8cf0ff;
}
.note-menu .note-menu-item-content::before {
    z-index: auto;
    content: "";
    background-color: #0000;
    position: absolute;
    left: 8px;
    right: 8px;
    top: 0;
    bottom: 0;
    pointer-events: none;
    border-radius: 3px;
    transition: background-color 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}
.note-menu .note-menu-item-content > * {
    z-index: 1;
}
</style>
