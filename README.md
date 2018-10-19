### 一个 go语言实现的博客 
#### 技术栈
    - golang
    - vue
    - echo 框架
    - 前台界面目前修改于 https://firekylin.org/
### 开发原因 
    - 现在自己使用博客系统不方便管理、部署和二次开发
    - 渲染速度到不到想要的效果
### 计划功能
   - [x] 文章管理
   - [x] 页面管理
   - [x] 分类管理
   - [x] 标签管理
   - [x] 系统设置
   - [ ] 评论
   - [ ] 待计划
### 使用方式
    下载代码
####  运行服务
    - 配置```conf.json```数据库信息
    - 还原 ```data```目录下 ```z_blog.sql``` 数据库
    - 安装依赖
    - ```go run```  
####  运行后台
    - 安装依赖 ``` npm install ```
    - 开发运行 ``` npm run dev ```
    - 浏览器打开 [http://127.0.0.1:8080/](http://127.0.0.1:8080/)
    - 发布 ```npm run build ``` 会自动发布到 ```res/dist```目录下面
### 效果展示
### 前台
![image](./data/imgs/front-0.png)
![image](./data/imgs/front-1.png)
![image](./data/imgs/front-2.png)
### 后台
![image](./data/imgs/backend-0.png)
![image](./data/imgs/backend-1.png)
![image](./data/imgs/backend-2.png)
![image](./data/imgs/backend-3.png)
![image](./data/imgs/backend-4.png)
![image](./data/imgs/backend-5.png)