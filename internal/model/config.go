package model

import (
	"blog/pkg/rndm"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// 配置值类型
const (
	CfgTypeString = "string" //字符串
	CfgTypeBool   = "bool"   //布尔
	CfgTypeInt    = "int"    //整数
	CfgTypeFloat  = "float"  //浮点数
	CfgTypeJson   = "json"   //JSON 文本
)

// 内置配置主键（主键即业务键）
const (
	ConfigIdAuth   = "auth"   //认证配置
	ConfigIdBasic  = "basic"  //站点基础配置
	ConfigIdBlog   = "blog"   //博客运行配置
	ConfigIdLogin  = "login"  //登录页配置
	ConfigIdLayout = "layout" //布局配置（管理端导航、主题、多页签等外观项）
	ConfigIdSql    = "sql"    //SQL 日志配置
	ConfigIdUpload = "upload" //上传配置
)

// ConfigGroup 配置分组
// 分组名直接存中文，与权限点分组保持一致的展示习惯
// 后台按分组折叠展示，列表也支持按分组过滤
const (
	ConfigGroupAuth   = "认证" //认证、登录
	ConfigGroupSite   = "站点" //站名、备案、页脚等
	ConfigGroupSystem = "系统" //日志、上传等运行参数
	ConfigGroupUI     = "界面" //管理端布局、主题
	ConfigGroupOther  = "其它" //未指定分组时的兜底值
)

// BuiltinConfigIds 内置配置主键集合
// 每一项都有独立的 API 路径与独立权限点，便于按配置分权管理
var BuiltinConfigIds = []string{
	ConfigIdBasic,
	ConfigIdLogin,
	ConfigIdBlog,
	ConfigIdAuth,
	ConfigIdUpload,
	ConfigIdLayout,
	ConfigIdSql,
}

// ConfigPerm 内置配置对应的权限标识
// 例：basic -> config:basic，可单独授予「只能改站点基础信息」的角色
func ConfigPerm(id string) string {
	return "config:" + id
}

// 登录页风格
// mgmt：当前管理端的左右分栏 + 品牌区布局（默认）
// vue：旧版前台的居中卡片 + 浅灰底布局
// 风格是纯展示层选项，前端按 style 渲染不同登录组件，不影响登录流程本身
const (
	LoginStyleMgmt    = "mgmt" //管理端风格-默认
	LoginStyleVue     = "vue"  //旧版前台风格
	LoginStyleDefault = LoginStyleMgmt
)

// PermLoginStyle 登录页风格切换权限
// 风格决定对外登录入口的观感，属于站点门面，单独授权而不共用 config:update
const PermLoginStyle = "login:style"

// 认证配置取值范围与默认值
// 过期时间限制在 1 小时~30 天：过短影响使用，过长等于永久有效
// 私钥限制在 16~64 字符：过短易被猜解，过长无收益且占容量
const (
	AuthTokenExpMin     = 1                                 //过期时间最小-小时
	AuthTokenExpMax     = 720                               //过期时间最大-小时(30 天)
	AuthTokenExpDefault = 24                                //过期时间默认值-小时
	AuthSecretMin       = 16                                //私钥最小长度
	AuthSecretMax       = 64                                //私钥最大长度
	AuthSecretDefault   = "blog-token-secret-please-change" //私钥默认值-上线前必须修改
)

// Config 运行时配置
// 主键 Id 同时作为唯一业务键：创建时允许自定义（如 basic/blog/login），
// 未指定时由 BeforeInsert 生成 ULID
type Config struct {
	Id       string `xorm:"VARCHAR(64) PK" json:"id"`       //主键，兼作配置键
	Group    string `xorm:"VARCHAR(64) INDEX" json:"group"` //分组-便于检索与折叠展示
	Name     string `xorm:"VARCHAR(255)" json:"name"`       //显示名
	Value    string `xorm:"TEXT" json:"value"`              //值
	Min      int    `xorm:"INT(11)" json:"min"`             //数字最小值/字符串最小长度
	Max      int    `xorm:"INT(11)" json:"max"`             //数字最大值/字符串最大长度
	Type     string `xorm:"VARCHAR(64)" json:"type"`        //值类型
	Describe string `xorm:"VARCHAR(255)" json:"describe"`   //描述
	Editable bool   `xorm:"BOOL DEFAULT 0" json:"editable"` //是否允许修改
	Sort     int    `xorm:"INT(11)" json:"sort"`            //排序
	Created  int64  `xorm:"BIGINT" json:"created"`          //创建时间
	Updated  int64  `xorm:"BIGINT" json:"updated"`          //修改时间
}

func (Config) TableName() string {
	return "config"
}

// BeforeInsert 插入前生成主键
func (s *Config) BeforeInsert() {
	if s.Id == "" {
		s.Id = rndm.ULID()
	}
}

// IsPkZero 主键是否为空
func (s Config) IsPkZero() bool {
	return s.Id == ""
}

// GetId 返回主键
func (s Config) GetId() string {
	return s.Id
}

// Json 反序列化配置值为目标结构
func (c *Config) Json(v any) error {
	return json.Unmarshal([]byte(c.Value), v)
}

// Int 配置值转 int
func (c *Config) Int() int {
	v, _ := strconv.Atoi(c.Value)
	return v
}

// Bool 配置值转 bool
func (c *Config) Bool() bool {
	v, _ := strconv.ParseBool(c.Value)
	return v
}

// Float 配置值转 float64
func (c *Config) Float() float64 {
	v, _ := strconv.ParseFloat(c.Value, 64)
	return v
}

// Valid 按类型校验配置值
func (c *Config) Valid() error {
	c.Id = strings.TrimSpace(c.Id)
	c.Group = strings.TrimSpace(c.Group)
	c.Name = strings.TrimSpace(c.Name)
	c.Describe = strings.TrimSpace(c.Describe)
	c.Type = strings.ToLower(strings.TrimSpace(c.Type))
	if err := ValidateBizId(c.Id, "配置主键"); err != nil {
		return err
	}
	// 分组留空时归入「其它」，保证列表里不会出现空分组
	if c.Group == "" {
		c.Group = ConfigGroupOther
	}
	if len(c.Group) > 64 {
		return errors.New("配置分组长度不能超过 64")
	}
	if len(c.Name) > 255 || len(c.Describe) > 255 {
		return errors.New("名称或描述长度不能超过 255")
	}
	if c.Min > c.Max && c.Max != 0 {
		return errors.New("最小值不能大于最大值")
	}
	switch c.Type {
	case CfgTypeBool:
		if _, err := strconv.ParseBool(c.Value); err != nil {
			return errors.New("布尔值非法")
		}
	case CfgTypeInt:
		n, err := strconv.ParseInt(c.Value, 10, 64)
		if err != nil {
			return errors.New("整数非法")
		}
		if c.Max != 0 && (n < int64(c.Min) || n > int64(c.Max)) {
			return fmt.Errorf("整数超出范围 %d~%d", c.Min, c.Max)
		}
	case CfgTypeFloat:
		f, err := strconv.ParseFloat(c.Value, 64)
		if err != nil {
			return errors.New("浮点数非法")
		}
		if c.Max != 0 && (f < float64(c.Min) || f > float64(c.Max)) {
			return fmt.Errorf("浮点数超出范围 %d~%d", c.Min, c.Max)
		}
	case CfgTypeJson:
		if !json.Valid([]byte(c.Value)) {
			return errors.New("JSON 格式非法")
		}
	case CfgTypeString:
		if c.Max != 0 && (len(c.Value) < c.Min || len(c.Value) > c.Max) {
			return fmt.Errorf("字符串长度超出范围 %d~%d", c.Min, c.Max)
		}
	default:
		return fmt.Errorf("不支持的配置类型 %s", c.Type)
	}
	// 内置认证配置额外校验取值区间，避免把非法过期时间或短私钥写进库
	if c.Id == ConfigIdAuth {
		mod := &Auth{}
		if err := c.Json(mod); err != nil {
			return errors.New("认证配置格式非法")
		}
		mod.Normalize()
		if err := mod.Valid(); err != nil {
			return err
		}
	}
	// 上传配置同理：体积与尺寸边界必须校验，否则一个请求就能写满磁盘
	if c.Id == ConfigIdUpload {
		mod := &Upload{}
		if err := c.Json(mod); err != nil {
			return errors.New("上传配置格式非法")
		}
		if err := mod.Valid(); err != nil {
			return err
		}
	}
	// SQL 日志配置：阈值过小会让正常查询全被标成慢 SQL
	if c.Id == ConfigIdSql {
		mod := &Sql{}
		if err := c.Json(mod); err != nil {
			return errors.New("SQL 日志配置格式非法")
		}
		if err := mod.Valid(); err != nil {
			return err
		}
	}
	// 布局配置：布局枚举与宽度区间必须校验，否则前端会渲染出异常布局
	// 登录页风格也在这里（layout_login_style），一并校验
	if c.Id == ConfigIdLayout {
		mod := &Layout{}
		if err := c.Json(mod); err != nil {
			return errors.New("布局配置格式非法")
		}
		if err := mod.Valid(); err != nil {
			return err
		}
	}
	return nil
}

// NormalizeValue 规范化配置值
func (c *Config) NormalizeValue() error {
	switch c.Type {
	case CfgTypeBool:
		b, err := strconv.ParseBool(strings.TrimSpace(c.Value))
		if err != nil {
			return errors.New("布尔值非法")
		}
		c.Value = strconv.FormatBool(b)
	case CfgTypeInt:
		n, err := strconv.ParseInt(strings.TrimSpace(c.Value), 10, 64)
		if err != nil {
			return errors.New("整数非法")
		}
		c.Value = strconv.FormatInt(n, 10)
	case CfgTypeFloat:
		f, err := strconv.ParseFloat(strings.TrimSpace(c.Value), 64)
		if err != nil {
			return errors.New("浮点数非法")
		}
		c.Value = strconv.FormatFloat(f, 'f', -1, 64)
	case CfgTypeJson:
		var buf bytes.Buffer
		if err := json.Compact(&buf, []byte(c.Value)); err != nil {
			return errors.New("JSON 格式非法")
		}
		c.Value = buf.String()
	case CfgTypeString:
	default:
		return fmt.Errorf("不支持的配置类型 %s", c.Type)
	}
	return nil
}

// ToBasic 反序列化为站点基础配置
func (c *Config) ToBasic() *Basic {
	mod := &Basic{}
	_ = c.Json(mod)
	return mod
}

// ToBlog 反序列化为博客运行配置
func (c *Config) ToBlog() *Blog {
	mod := &Blog{}
	_ = c.Json(mod)
	return mod
}

// ToLogin 反序列化为登录页配置
func (c *Config) ToLogin() *Login {
	mod := LoginDefault()
	_ = c.Json(mod)
	return mod
}

// ToLayout 反序列化为布局配置
// 解析失败或字段缺失时回落默认值，保证管理端总能拿到一份可用布局
func (c *Config) ToLayout() *Layout {
	mod := LayoutDefault()
	_ = c.Json(mod)
	mod.Normalize()
	return mod
}

// ToAuth 反序列化为认证配置
// 解析失败或字段缺失时回落默认值，保证签名永远有可用的私钥
func (c *Config) ToAuth() *Auth {
	mod := AuthDefault()
	_ = c.Json(mod)
	mod.Normalize()
	if mod.TokenExp < AuthTokenExpMin || mod.TokenExp > AuthTokenExpMax {
		mod.TokenExp = AuthTokenExpDefault
	}
	if len(mod.TokenSecret) < AuthSecretMin || len(mod.TokenSecret) > AuthSecretMax {
		mod.TokenSecret = AuthSecretDefault
	}
	return mod
}

// ToUpload 反序列化为上传配置
// 逐项回落：某一项非法时只重置该项，不因为单个脏值作废整份配置
// 其中 ImageMaxSize 不得超过 FileMaxSize，否则以文件上限为准
func (c *Config) ToUpload() *Upload {
	mod := UploadDefault()
	_ = c.Json(mod)
	if mod.FileMaxSize < UploadFileMaxSizeMin || mod.FileMaxSize > UploadFileMaxSizeMax {
		mod.FileMaxSize = UploadFileMaxSizeDefault
	}
	if mod.ImageMaxSize < UploadImageMaxSizeMin || mod.ImageMaxSize > UploadImageMaxSizeMax {
		mod.ImageMaxSize = UploadImageMaxSizeDefault
	}
	if mod.ImageMaxSize > mod.FileMaxSize {
		mod.ImageMaxSize = mod.FileMaxSize
	}
	if mod.ImageMaxWidth < UploadImagePxMin || mod.ImageMaxWidth > UploadImagePxMax {
		mod.ImageMaxWidth = UploadImageWidthDefault
	}
	if mod.ImageMaxHeight < UploadImagePxMin || mod.ImageMaxHeight > UploadImagePxMax {
		mod.ImageMaxHeight = UploadImageHeightDefault
	}
	return mod
}

// ToSql 反序列化为 SQL 日志配置
// 阈值非法时回落默认，保证不会因为一个脏值把慢日志关掉或刷屏
func (c *Config) ToSql() *Sql {
	mod := SqlDefault()
	_ = c.Json(mod)
	if mod.SlowMs > 0 && (mod.SlowMs < SqlSlowMsMin || mod.SlowMs > SqlSlowMsMax) {
		mod.SlowMs = SqlSlowMsDefault
	}
	return mod
}

// ------------------------------------------------------ 内置配置结构 ------------------------------------------------------

// Basic 站点基础配置
type Basic struct {
	Title       string `json:"title"`       //标题
	SiteUrl     string `json:"site_url"`    //网站地址
	LogoUrl     string `json:"logo_url"`    //Logo地址
	Author      string `json:"author"`      //网站作者
	Copyright   string `json:"copyright"`   //版权
	Keywords    string `json:"keywords"`    //网站关键词
	Description string `json:"description"` //描述
	BeianMiit   string `json:"beian_miit"`  //ICP备案
	BeianNism   string `json:"beian_nism"`  //公安备案
}

// Blog 博客运行配置
type Blog struct {
	PageSize  int    `json:"page_size"`  //分页大小
	Analytic  string `json:"analytic"`   //统计脚本
	Comment   string `json:"comment"`    //评论脚本
	GithubUrl string `json:"github_url"` //github地址
	WeiboUrl  string `json:"weibo_url"`  //微博地址
}

// Login 登录页配置
// 只存标题与描述；风格（mgmt / vue）属于外观开关，与导航模式、导航风格同属
// 布局配置（Layout.LoginStyle），在「布局配置」抽屉里一起切换。
type Login struct {
	Title string `json:"title"` //标题
	Desc  string `json:"desc"`  //描述
}

// LoginDefault 登录页配置默认值
func LoginDefault() *Login {
	return &Login{
		Title: "欢迎登录",
		Desc:  "便捷、高效、可靠",
	}
}

// LoginSite 登录页对外数据
// 标题与描述来自 login 配置，风格来自 layout 配置；前端登录页只读这一个结构，
// 不需要自己拼两个配置。
type LoginSite struct {
	Title string `json:"title"` //标题
	Desc  string `json:"desc"`  //描述
	Style string `json:"style"` //风格 mgmt/vue
}

// 管理端导航模式
const (
	NavModeVertical      = "vertical"       //左侧菜单
	NavModeHorizontal    = "horizontal"     //顶部菜单
	NavModeHorizontalMix = "horizontal-mix" //顶部菜单混合
)

// 管理端导航风格
const (
	NavThemeDark       = "dark"        //暗色侧边栏
	NavThemeLight      = "light"       //白色侧边栏
	NavThemeHeaderDark = "header-dark" //暗色顶栏
)

// 路由动画类型
const (
	PageAnimateZoomFade  = "zoom-fade"
	PageAnimateZoomOut   = "zoom-out"
	PageAnimateFadeSlide = "fade-slide"
	PageAnimateFade      = "fade"
	PageAnimateFadeBot   = "fade-bottom"
	PageAnimateFadeScale = "fade-scale"
)

// 菜单宽度取值范围
// 过窄放不下图标与文字，过宽会挤占内容区，这里给出可用区间
const (
	MenuWidthMin     = 48
	MenuWidthMax     = 400
	MenuWidthDefault = 200
	MenuMinWidthMin  = 48
	MenuMinWidthMax  = 200
)

// Layout 布局配置
// 管理端导航、主题、多页签等外观项的持久化载体：以前这些值只活在浏览器内存里，
// 刷新即回默认；现在存进 config 表（主键 layout），由后端统一下发，
// 换设备、换浏览器表现一致，也便于按权限控制谁能改。
//
// 字段名统一 snake_case，与前端 projectSetting 的驼峰字段由前端映射；
// 移动端标记 isMobile 与权限模式 permissionMode 属于运行时/编译期概念，不入库。
type Layout struct {
	NavMode         string `json:"nav_mode"`          //导航模式 vertical/horizontal/horizontal-mix
	NavTheme        string `json:"nav_theme"`         //导航风格 dark/light/header-dark
	LoginStyle      string `json:"login_style"`       //登录页风格 mgmt/vue
	HeaderFixed     bool   `json:"header_fixed"`      //固定顶栏
	HeaderReload    bool   `json:"header_reload"`     //显示重载页面按钮
	MultiTabsShow   bool   `json:"multi_tabs_show"`   //显示多页签
	MultiTabsFixed  bool   `json:"multi_tabs_fixed"`  //固定多页签
	MenuFixed       bool   `json:"menu_fixed"`        //固定侧边栏
	MenuMix         bool   `json:"menu_mix"`          //分割菜单
	MenuCollapsed   bool   `json:"menu_collapsed"`    //折叠菜单
	MenuWidth       int    `json:"menu_width"`        //菜单宽度
	MenuMinWidth    int    `json:"menu_min_width"`    //折叠后菜单宽度
	CrumbShow       bool   `json:"crumb_show"`        //显示面包屑
	CrumbShowIcon   bool   `json:"crumb_show_icon"`   //显示面包屑图标
	ShowFooter      bool   `json:"show_footer"`       //显示页脚
	PageAnimate     bool   `json:"page_animate"`      //启用路由动画
	PageAnimateType string `json:"page_animate_type"` //路由动画类型
}

// LayoutDefault 布局配置默认值
// 与前端 settings/projectSetting.ts 的初始值保持一致，
// 这样「数据库没配过」与「前端写死默认」表现完全相同
func LayoutDefault() *Layout {
	return &Layout{
		NavMode:         NavModeVertical,
		NavTheme:        NavThemeDark,
		LoginStyle:      LoginStyleDefault,
		HeaderFixed:     true,
		HeaderReload:    true,
		MultiTabsShow:   true,
		MultiTabsFixed:  true,
		MenuFixed:       true,
		MenuMix:         false,
		MenuCollapsed:   false,
		MenuWidth:       MenuWidthDefault,
		MenuMinWidth:    64,
		CrumbShow:       true,
		CrumbShowIcon:   false,
		ShowFooter:      true,
		PageAnimate:     true,
		PageAnimateType: PageAnimateZoomFade,
	}
}

// projectAnimateTypes 合法动画类型集合
var projectAnimateTypes = map[string]struct{}{
	PageAnimateZoomFade:  {},
	PageAnimateZoomOut:   {},
	PageAnimateFadeSlide: {},
	PageAnimateFade:      {},
	PageAnimateFadeBot:   {},
	PageAnimateFadeScale: {},
}

// Normalize 把越界/未知取值收敛到合法范围
// 只做「纠正」不做「拒绝」：布局参数错了顶多不好看，不该让整份配置不可用
func (s *Layout) Normalize() {
	s.NavMode = strings.TrimSpace(strings.ToLower(s.NavMode))
	switch s.NavMode {
	case NavModeVertical, NavModeHorizontal, NavModeHorizontalMix:
	default:
		s.NavMode = NavModeVertical
	}
	s.NavTheme = strings.TrimSpace(strings.ToLower(s.NavTheme))
	switch s.NavTheme {
	case NavThemeDark, NavThemeLight, NavThemeHeaderDark:
	default:
		s.NavTheme = NavThemeDark
	}
	s.LoginStyle = strings.TrimSpace(strings.ToLower(s.LoginStyle))
	if s.LoginStyle != LoginStyleMgmt && s.LoginStyle != LoginStyleVue {
		s.LoginStyle = LoginStyleDefault
	}
	if _, ok := projectAnimateTypes[s.PageAnimateType]; !ok {
		s.PageAnimateType = PageAnimateZoomFade
	}
	if s.MenuWidth < MenuWidthMin || s.MenuWidth > MenuWidthMax {
		s.MenuWidth = MenuWidthDefault
	}
	if s.MenuMinWidth < MenuMinWidthMin || s.MenuMinWidth > MenuMinWidthMax {
		s.MenuMinWidth = 64
	}
	// 只有顶部菜单混合模式才有「分割菜单」概念，其余模式强制关闭
	if s.NavMode != NavModeHorizontalMix {
		s.MenuMix = false
	}
}

// Valid 校验布局配置
// 写入时对枚举做严格校验：静默纠正会让使用者以为设置生效了，实际没有
func (s *Layout) Valid() error {
	mode := strings.TrimSpace(strings.ToLower(s.NavMode))
	switch mode {
	case NavModeVertical, NavModeHorizontal, NavModeHorizontalMix:
	default:
		return fmt.Errorf("导航模式仅支持 %s / %s / %s", NavModeVertical, NavModeHorizontal, NavModeHorizontalMix)
	}
	theme := strings.TrimSpace(strings.ToLower(s.NavTheme))
	switch theme {
	case NavThemeDark, NavThemeLight, NavThemeHeaderDark:
	default:
		return fmt.Errorf("导航风格仅支持 %s / %s / %s", NavThemeDark, NavThemeLight, NavThemeHeaderDark)
	}
	style := strings.TrimSpace(strings.ToLower(s.LoginStyle))
	switch style {
	case LoginStyleMgmt, LoginStyleVue:
	default:
		return fmt.Errorf("登录页风格仅支持 %s 或 %s", LoginStyleMgmt, LoginStyleVue)
	}
	if _, ok := projectAnimateTypes[s.PageAnimateType]; !ok {
		return errors.New("路由动画类型非法")
	}
	if s.MenuWidth < MenuWidthMin || s.MenuWidth > MenuWidthMax {
		return fmt.Errorf("菜单宽度超出范围 %d~%d", MenuWidthMin, MenuWidthMax)
	}
	if s.MenuMinWidth < MenuMinWidthMin || s.MenuMinWidth > MenuMinWidthMax {
		return fmt.Errorf("折叠后菜单宽度超出范围 %d~%d", MenuMinWidthMin, MenuMinWidthMax)
	}
	return nil
}

// Auth 认证配置
// 存于 config 表（主键 auth），可在线修改，改完立即生效，共三项：
//   - TokenExp    令牌过期时间（小时），只影响新签发的令牌，已签发的不变
//   - TokenSecret 令牌与图形验证码的加密私钥，修改后所有旧令牌立即失效
//   - Audit       是否记录操作审计日志，关闭后写操作与登录均不再留痕
//   - Docs        是否开放 /swagger 接口文档，关闭后该路径直接返回 404
//
// TokenExp/TokenSecret 历史上写在 conf.toml，现统一收敛到数据库配置，理由：
// 容器化部署时改 toml 需要重启，而改成配置可以在线调整。
type Auth struct {
	TokenExp    int    `json:"token_exp"`    //令牌过期时间-小时
	TokenSecret string `json:"token_secret"` //加密私钥
	Audit       bool   `json:"audit"`        //是否记录操作日志
	Docs        bool   `json:"docs"`         //是否开放接口文档
}

// AuthDefault 认证配置默认值
// 审计默认开启：敏感操作留痕是安全基线，需要时可在后台关闭
// 文档默认开启：本地开发即开箱查看，正式上线可在后台关掉
func AuthDefault() *Auth {
	return &Auth{
		TokenExp:    AuthTokenExpDefault,
		TokenSecret: AuthSecretDefault,
		Audit:       true,
		Docs:        true,
	}
}

// Normalize 规范化：私钥去首尾空白
func (a *Auth) Normalize() {
	a.TokenSecret = strings.TrimSpace(a.TokenSecret)
}

// Valid 校验取值范围
// 配置入库前与启动同步时都会调用，避免脏数据导致签名行为不可预期
func (a *Auth) Valid() error {
	if a.TokenExp < AuthTokenExpMin || a.TokenExp > AuthTokenExpMax {
		return fmt.Errorf("令牌过期时间必须在 %d~%d 小时之间", AuthTokenExpMin, AuthTokenExpMax)
	}
	n := len(a.TokenSecret)
	if n < AuthSecretMin || n > AuthSecretMax {
		return fmt.Errorf("加密私钥长度必须在 %d~%d 字符之间", AuthSecretMin, AuthSecretMax)
	}
	if strings.ContainsAny(a.TokenSecret, " \t\r\n") {
		return errors.New("加密私钥不能包含空白字符")
	}
	return nil
}

// 上传配置取值范围与默认值
// 体积单位统一为 MB，尺寸单位统一为像素；上限不开放到无穷，避免一个请求把磁盘写满
const (
	UploadFileMaxSizeMin      = 1     //通用文件上限最小-MB
	UploadFileMaxSizeMax      = 512   //通用文件上限最大-MB
	UploadFileMaxSizeDefault  = 32    //通用文件上限默认值-MB
	UploadImageMaxSizeMin     = 1     //图片上限最小-MB
	UploadImageMaxSizeMax     = 100   //图片上限最大-MB
	UploadImageMaxSizeDefault = 10    //图片上限默认值-MB
	UploadImagePxMin          = 100   //图片边长最小-px
	UploadImagePxMax          = 10000 //图片边长最大-px
	UploadImageWidthDefault   = 1920  //图片宽度上限默认值-px
	UploadImageHeightDefault  = 1920  //图片高度上限默认值-px
)

// SQL 日志配置取值范围与默认值
// 阈值下限 10ms：再低就会把正常查询全部标成慢 SQL，日志失去意义
const (
	SqlSlowMsMin     = 10    //慢 SQL 阈值最小-毫秒
	SqlSlowMsMax     = 60000 //慢 SQL 阈值最大-毫秒(60s)
	SqlSlowMsDefault = 300   //慢 SQL 阈值默认值-毫秒
)

// Sql SQL 日志配置
// 存于 config 表（主键 sql），可在线修改，下一次查询立即生效
//   - Show   是否把每条 SQL 都按 debug 级别打印
//   - SlowMs 慢 SQL 阈值（毫秒），超过则按 warn 级别打印；<=0 表示不记录
//
// 历史上这两项写在 conf.toml（show_sql 与硬编码的 300ms），
// 现改为数据库配置，排查线上问题时可以只开慢日志而不重启。
type Sql struct {
	Show   bool  `json:"show"`    //是否打印全部 SQL
	SlowMs int64 `json:"slow_ms"` //慢 SQL 阈值-毫秒，<=0 关闭慢日志
}

// SqlDefault SQL 日志配置默认值
// 默认不打全部 SQL、慢日志 300ms：生产环境按默认值即可
func SqlDefault() *Sql {
	return &Sql{
		Show:   false,
		SlowMs: SqlSlowMsDefault,
	}
}

// Valid 校验取值范围
// SlowMs <= 0 是有意义的值（关闭慢日志），所以只在大于 0 时校验上下限
func (s *Sql) Valid() error {
	if s.SlowMs > 0 && (s.SlowMs < SqlSlowMsMin || s.SlowMs > SqlSlowMsMax) {
		return fmt.Errorf("慢 SQL 阈值必须在 %d~%d 毫秒之间，或设为 0 关闭", SqlSlowMsMin, SqlSlowMsMax)
	}
	return nil
}

// SlowDuration 慢 SQL 阈值时长，<=0 表示不记录
func (s *Sql) SlowDuration() time.Duration {
	if s.SlowMs <= 0 {
		return 0
	}
	return time.Duration(s.SlowMs) * time.Millisecond
}

// Upload 上传配置
// 存于 config 表（主键 upload），可在线修改，下一次上传立即生效
//   - FileMaxSize   通用文件上传上限（MB）
//   - ImageMaxSize  图片上传上限（MB），不得超过 FileMaxSize
//   - ImageMaxWidth/ImageMaxHeight 图片边长上限（px）
//   - Crop          超出边长上限时是否自动等比缩小；为 false 则直接拒绝
type Upload struct {
	FileMaxSize    int  `json:"file_max_size"`    //通用文件上限-MB
	ImageMaxSize   int  `json:"image_max_size"`   //图片上限-MB
	ImageMaxWidth  int  `json:"image_max_width"`  //图片宽度上限-px
	ImageMaxHeight int  `json:"image_max_height"` //图片高度上限-px
	Crop           bool `json:"crop"`             //超限时自动等比缩小
}

// UploadDefault 上传配置默认值
// 默认开启自动缩小：手机拍的照动辄 4000px，拒绝会让用户无从下手
func UploadDefault() *Upload {
	return &Upload{
		FileMaxSize:    UploadFileMaxSizeDefault,
		ImageMaxSize:   UploadImageMaxSizeDefault,
		ImageMaxWidth:  UploadImageWidthDefault,
		ImageMaxHeight: UploadImageHeightDefault,
		Crop:           true,
	}
}

// FileBytes 通用文件上限（字节）
func (u *Upload) FileBytes() int64 {
	return int64(u.FileMaxSize) << 20
}

// ImageBytes 图片上限（字节）
func (u *Upload) ImageBytes() int64 {
	return int64(u.ImageMaxSize) << 20
}

// Fit 计算等比缩小后的尺寸
// 返回 false 表示无需缩放（已在限制范围内）；只缩小不放大，避免小图被拉伸糊掉
func (u *Upload) Fit(w, h int) (int, int, bool) {
	maxW, maxH := u.ImageMaxWidth, u.ImageMaxHeight
	if w < 1 || h < 1 || maxW < 1 || maxH < 1 {
		return w, h, false
	}
	if w <= maxW && h <= maxH {
		return w, h, false
	}
	// 取两轴中受限更严重的一边作为缩放比，保证两边都不越界
	ratio := math.Min(float64(maxW)/float64(w), float64(maxH)/float64(h))
	nw, nh := int(math.Round(float64(w)*ratio)), int(math.Round(float64(h)*ratio))
	if nw < 1 {
		nw = 1
	}
	if nh < 1 {
		nh = 1
	}
	return nw, nh, true
}

// Valid 校验取值范围
func (u *Upload) Valid() error {
	if u.FileMaxSize < UploadFileMaxSizeMin || u.FileMaxSize > UploadFileMaxSizeMax {
		return fmt.Errorf("文件大小上限必须在 %d~%d MB 之间", UploadFileMaxSizeMin, UploadFileMaxSizeMax)
	}
	if u.ImageMaxSize < UploadImageMaxSizeMin || u.ImageMaxSize > UploadImageMaxSizeMax {
		return fmt.Errorf("图片大小上限必须在 %d~%d MB 之间", UploadImageMaxSizeMin, UploadImageMaxSizeMax)
	}
	if u.ImageMaxSize > u.FileMaxSize {
		return errors.New("图片大小上限不能超过文件大小上限")
	}
	if u.ImageMaxWidth < UploadImagePxMin || u.ImageMaxWidth > UploadImagePxMax {
		return fmt.Errorf("图片宽度上限必须在 %d~%d px 之间", UploadImagePxMin, UploadImagePxMax)
	}
	if u.ImageMaxHeight < UploadImagePxMin || u.ImageMaxHeight > UploadImagePxMax {
		return fmt.Errorf("图片高度上限必须在 %d~%d px 之间", UploadImagePxMin, UploadImagePxMax)
	}
	return nil
}

// ConfigListArgs 配置列表参数
type ConfigListArgs struct {
	List
	Group string `query:"group" form:"group" json:"group"` //分组过滤
	Name  string `query:"name" form:"name" json:"name"`    //名称模糊匹配
	Type  string `query:"type" form:"type" json:"type"`    //类型过滤
}

// ConfigAddArgs 配置新增参数
type ConfigAddArgs struct {
	Id       string `json:"id" form:"id"`             //主键，可自定义，留空自动生成
	Group    string `json:"group" form:"group"`       //分组，留空归入「其它」
	Name     string `json:"name" form:"name"`         //显示名
	Value    string `json:"value" form:"value"`       //值
	Min      int    `json:"min" form:"min"`           //最小值
	Max      int    `json:"max" form:"max"`           //最大值
	Type     string `json:"type" form:"type"`         //值类型
	Describe string `json:"describe" form:"describe"` //描述
	Sort     int    `json:"sort" form:"sort"`         //排序
}

// ConfigGroup 配置分组及数量
type ConfigGroup struct {
	Group string `json:"group"` //分组名
	Count int    `json:"count"` //该分组下的配置数量
}

// ConfigEditArgs 配置修改参数
// 指针字段：nil 表示不修改该项，区分「没传」与「传了空值」。
// 名称/最小值/最大值/排序也在此列：配置管理页的编辑抽屉一次性提交全部元信息，
// 不再区分「改值」与「改元信息」两个入口。
type ConfigEditArgs struct {
	Value    *string `json:"value" form:"value"`       //值
	Describe *string `json:"describe" form:"describe"` //描述
	Group    *string `json:"group" form:"group"`       //分组
	Name     *string `json:"name" form:"name"`         //显示名
	Min      *int    `json:"min" form:"min"`           //数字最小值/字符串最小长度
	Max      *int    `json:"max" form:"max"`           //数字最大值/字符串最大长度
	Sort     *int    `json:"sort" form:"sort"`         //排序
}

// ConfigLoginStyleArgs 登录页风格切换参数
// 风格单独成接口，便于用 login:style 权限把「站点门面」的改动与普通配置修改分开
type ConfigLoginStyleArgs struct {
	Style string `json:"style" form:"style"` //风格 mgmt/vue
}

// LayoutEditArgs 布局配置修改参数
// 整体替换而非按字段增量：前端在抽屉里改任一开关时本地已是最新全量，
// 直接回写整份对象，避免多字段合并时产生半新半旧的状态
type LayoutEditArgs = Layout
