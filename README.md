# Hope - Shopify 产品保险服务应用

## 🚀 快速开始

### 环境要求
- Docker & Docker Compose
- Git

### 1. 克隆项目
```bash
git clone <repository-url>
cd hope
```

### 2. 配置环境变量
```bash
# 复制环境变量模板
cp env.example .env

# 编辑 .env 文件，设置您的配置
# 重要：请修改默认密码！
```

### 3. 启动服务
```bash
# 开发环境
docker-compose up -d

# 生产环境（推荐）
docker-compose -f docker-compose.prod.yml up -d
```

### 4. 访问应用
- 前端应用：http://localhost
- API文档：http://localhost:8080
- 监控面板：http://localhost:8090

## 📋 系统概述
这是一个为 Shopify 商家提供**产品保险服务**的应用，允许商家在其店铺中销售保险产品，并跟踪相关的订单和收益数据。

## 🔐 环境变量配置

本项目使用环境变量管理敏感配置信息。详细配置说明请查看 [ENV-README.md](./ENV-README.md)。

### 必需配置
- `DB_PASSWORD`: 数据库密码
- `REDIS_PASSWORD`: Redis密码
- `MYSQL_ROOT_PASSWORD`: MySQL root密码

### 安全建议
1. 使用强密码
2. 生产环境使用 `docker-compose.prod.yml`
3. 定期更换密码
4. 不要将 `.env` 文件提交到版本控制

## 🗂️ 表结构分析和数据存储时机

### 1. **用户管理模块**

#### `user` - 用户表
**存储内容**：Shopify 商家的基本信息
**数据来源时机**：
- 商家安装应用时创建
- Shopify OAuth 认证后获取店铺信息
- 商家首次登录时同步店铺数据

```sql
-- 商家安装应用时
INSERT INTO user (name, domain, access_token, email, country_name...)
```


### 2. **产品管理模块**

#### `user_product` - 保险产品表
**存储内容**：商家创建的保险产品
**数据来源时机**：
- 商家在应用中创建保险产品时
- 产品发布到 Shopify 时更新 `product_id` 和 `is_publish` 状态

#### `user_variant` - 产品变体表
**存储内容**：保险产品的不同规格/价格变体
**数据来源时机**：
- 创建产品时同时创建变体
- 发布到 Shopify 后更新 `variant_id`

#### `user_cart_setting` - 购物车配置表
**存储内容**：保险产品在购物车中的展示配置
**数据来源时机**：
- 商家配置保险产品展示方式时
- 设置价格计算规则、颜色、文案等

### 3. **订单处理模块**

#### `user_order` - 订单主表
**存储内容**：来自 Shopify 的订单信息
**数据来源时机**：
- **Webhook 触发**：Shopify 订单创建/更新时
- **定时同步**：定期拉取 Shopify 订单数据
- 包含保险金额计算结果

#### `user_order_info` - 订单详情表
**存储内容**：订单中每个商品的详细信息
**数据来源时机**：
- 处理订单时解析 line_items
- 标识哪些是保险产品 (`is_insurance`)

#### `order_summary` - 订单统计表
**存储内容**：按天聚合的订单统计数据
**数据来源时机**：
- **每日凌晨**：统计前一天的订单数据
- 计算销售额、保险收益、退款等

### 4. **任务队列模块**

#### `job_order` - 订单同步任务表
**存储内容**：订单同步任务的执行记录
**数据来源时机**：
- Shopify Webhook 接收到订单事件时创建任务
- 定时任务创建批量同步任务
- 记录处理状态和结果

#### `job_product` - 产品上传任务表
**存储内容**：产品发布到 Shopify 的任务记录
**数据来源时机**：
- 商家点击"发布到 Shopify"时创建任务
- 批量发布产品时创建多个任务

## 🔄 数据流向时序

### 商家使用流程：
1. **安装应用** → 创建 `user` 记录
2. **配置保险设置** → 更新 `user_cart_setting`
3. **创建保险产品** → 插入 `user_product` + `user_variant`
4. **发布产品** → 创建 `job_product` 任务 → 更新产品状态
5. **客户下单** → Webhook 触发 → 创建 `job_order` 任务
6. **处理订单** → 插入 `user_order` + `user_order_info`
7. **生成统计** → 每日更新 `order_summary`

### 关键业务逻辑：
- **保险金额计算**：根据 `user_cart_setting` 中的定价规则计算保险费用
- **收益统计**：从订单中提取保险部分的收益进行统计
- **异步处理**：使用任务队列处理 Shopify API 调用，避免阻塞
- **数据同步**：定期与 Shopify 同步订单状态和产品信息

这是一个典型的 **SaaS 电商插件系统**，专注于为 Shopify 商家提供产品保险销售能力，通过佣金模式盈利。

开启插件后,自动同步 7 天的单号,同步其中有购买保险的单号
1. orders 时间筛选抽离
    1. orders 顶部统计卡片
2. protection page 设置
    1. checkbox 样式
    2. demo 的佣金计算逻辑
    3. demo 的跳转
    4. product 删除后的图片异常问题
    5. product 禁用发货设置
    6. css样式设置和应用
   
    7. 插件安装检测banner
    8. fulfillment 规则应用
    9. 佣金设置的默认值
    10. 佣金设置 other value的针对不同 type 的数据类型设置
    11. cart setting 插件设置应用
    12. Footer link URL 在 demo 里的应用
3. ~~全局loading样式~~
4. dashboard 
   1. 图表样式
   2. 新人引导默认选中规则
   3. help card
   4. 默认图表数据展示
   5. 统计环比