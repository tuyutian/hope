-- 优化后的数据表结构 - 基于性能和最佳实践
-- 生成时间: 2025-06-27

-- 用户表
CREATE TABLE `user` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `name` varchar(255) NOT NULL COMMENT 'shopify-name',
    `domain` varchar(255) NOT NULL DEFAULT '' COMMENT 'my shopify Domain网站域名',
    `plan_display_name` varchar(40) NOT NULL DEFAULT '' COMMENT 'shopify套餐版本',
    `access_token` varchar(512) NOT NULL DEFAULT '' COMMENT 'shopify-token',
    `user_token` varchar(255) NOT NULL DEFAULT '' COMMENT '用户token',
    `pwd` varchar(255) NOT NULL DEFAULT '' COMMENT '密码',
    `level` tinyint(2) NOT NULL DEFAULT 0 COMMENT 'app内部套餐等级',
    `email` varchar(100) NOT NULL DEFAULT '' COMMENT '邮箱',
    `country_name` varchar(25) NOT NULL DEFAULT '' COMMENT '国家名称',
    `country_code` varchar(25) NOT NULL DEFAULT '' COMMENT '国家简码',
    `city` varchar(25) NOT NULL DEFAULT '' COMMENT '城市',
    `trial_time` bigint(20) NOT NULL DEFAULT 0 COMMENT '试用时间',
    `currency_code` varchar(10) NOT NULL DEFAULT '' COMMENT '货币简码',
    `money_format` varchar(20) NOT NULL DEFAULT '' COMMENT '货币单位符号',
    `last_login` bigint(20) NOT NULL DEFAULT 0,
    `is_del` tinyint(1) NOT NULL DEFAULT 1 COMMENT '删除状态 1 正常 2 卸载 3关店',
    `publish_id` varchar(100) NOT NULL DEFAULT '' COMMENT '店铺publish_id',
    `steps` varchar(500) NOT NULL DEFAULT '' COMMENT '新手引导',
    `collection` text COMMENT '用户集合列表',
    `uninstall_time` bigint(20) NOT NULL DEFAULT 0 COMMENT '卸载时间',
    `create_time` bigint(20) NOT NULL COMMENT '创建时间',
    `update_time` bigint(20) NOT NULL COMMENT '最近修改时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_domain` (`domain`),
    KEY `idx_name` (`name`),
    KEY `idx_email` (`email`),
    KEY `idx_is_del` (`is_del`),
    KEY `idx_level` (`level`),
    KEY `idx_last_login` (`last_login`),
    KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 保险用户产品表
CREATE TABLE `user_product` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `user_id` bigint(20) NOT NULL COMMENT '用户id',
    `product_id` varchar(100) NOT NULL DEFAULT '' COMMENT 'shopify上传成功的产品ID',
    `title` varchar(255) NOT NULL COMMENT '标题',
    `product_type` varchar(255) NOT NULL COMMENT '产品类型',
    `vendor` varchar(255) NOT NULL COMMENT 'vendor',
    `collection` varchar(255) NOT NULL COMMENT '集合',
    `tags` varchar(500) NOT NULL DEFAULT '' COMMENT '产品标签',
    `description` text COMMENT '描述',
    `option_1` varchar(255) NOT NULL COMMENT '产品属性1',
    `option_2` varchar(255) NOT NULL DEFAULT '' COMMENT '产品属性2',
    `option_3` varchar(255) NOT NULL DEFAULT '' COMMENT '产品属性3',
    `image_url` varchar(500) NOT NULL COMMENT '封面图',
    `is_publish` tinyint(1) NOT NULL DEFAULT 0 COMMENT '发布Shopify：0:未发布 1:已发布 2:正在发布中 3:shopify平台已删除',
    `publish_time` bigint(20) NOT NULL DEFAULT 0 COMMENT '发布时间',
    `is_del` tinyint(1) NOT NULL DEFAULT 0 COMMENT '删除状态 0 正常 1 已删除',
    `create_time` bigint(20) NOT NULL COMMENT '创建时间',
    `update_time` bigint(20) NOT NULL COMMENT '修改时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id_del` (`user_id`, `is_del`),
    KEY `idx_product_id` (`product_id`),
    KEY `idx_product_type` (`product_type`),
    KEY `idx_is_publish` (`is_publish`),
    KEY `idx_publish_time` (`publish_time`),
    KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='保险用户产品表';

-- 保险用户变体表
CREATE TABLE `user_variant` (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `user_id` int(11) NOT NULL COMMENT '用户id',
    `user_product_id` int(11) NOT NULL COMMENT '保险用户产品表ID',
    `product_id` varchar(100) NOT NULL DEFAULT '' COMMENT 'Shopify产品ID',
    `variant_id` varchar(100) NOT NULL DEFAULT '' COMMENT 'Shopify变体ID',
    `inventory_id` varchar(100) NOT NULL DEFAULT '' COMMENT 'Shopify仓库ID',
    `sku_name` varchar(150) NOT NULL DEFAULT '' COMMENT 'SKU',
    `image_url` varchar(500) NOT NULL DEFAULT '' COMMENT '变体封面图',
    `sku_1` varchar(150) NOT NULL DEFAULT '' COMMENT '变体属性1',
    `sku_2` varchar(150) NOT NULL DEFAULT '' COMMENT '变体属性2',
    `sku_3` varchar(150) NOT NULL DEFAULT '' COMMENT '变体属性3',
    `price` decimal(12,2) NOT NULL DEFAULT 0.00 COMMENT '价格设定',
    `create_time` bigint(20) NOT NULL COMMENT '创建时间',
    `update_time` bigint(20) NOT NULL COMMENT '修改时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_user_product_id` (`user_product_id`),
    KEY `idx_product_id` (`product_id`),
    KEY `idx_variant_id` (`variant_id`),
    KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='保险用户变体表';

-- 保险用户基础配置表
CREATE TABLE `user_cart_setting` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `user_id` bigint(20) NOT NULL COMMENT '用户id',
    `plan_title` varchar(100) NOT NULL DEFAULT '' COMMENT '保险标题(内部)',
    `addon_title` varchar(100) NOT NULL DEFAULT '' COMMENT '保险标题',
    `enabled_desc` varchar(200) NOT NULL DEFAULT '' COMMENT '按钮打开文案',
    `disabled_desc` varchar(200) NOT NULL DEFAULT '' COMMENT '按钮关闭文案',
    `foot_text` varchar(100) NOT NULL DEFAULT '' COMMENT '保险底部',
    `foot_url` varchar(255) NOT NULL DEFAULT '' COMMENT '保险跳转',
    `in_color` varchar(50) NOT NULL DEFAULT '' COMMENT '打开颜色',
    `out_color` varchar(50) NOT NULL DEFAULT '' COMMENT '关闭颜色',
    `other_money` decimal(12,2) NOT NULL DEFAULT 0.00 COMMENT '其他金额',
    `show_cart` tinyint(1) NOT NULL DEFAULT 0 COMMENT '购物车状态 0 关闭 1 打开',
    `show_cart_icon` tinyint(1) NOT NULL DEFAULT 0 COMMENT '购物车图标 0 关闭 1 打开',
    `icon_url` text COMMENT '选中url(json)',
    `select_button` tinyint(1) NOT NULL DEFAULT 0 COMMENT '购物车图标 0 滑动 1 勾选',
    `product_type` varchar(100) NOT NULL DEFAULT '' COMMENT '产品type',
    `product_collection` varchar(100) NOT NULL DEFAULT '' COMMENT '产品选中集合',
    `pricing_type` tinyint(1) NOT NULL DEFAULT 0 COMMENT '购物车图标 0 金额 1百分比',
    `pricing_select` text COMMENT '金额计算范围',
    `tiers_select` text COMMENT '百分比计算范围',
    `create_time` bigint(20) NOT NULL COMMENT '创建时间',
    `update_time` bigint(20) NOT NULL COMMENT '修改时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='保险用户基础配置表';

-- 用户订单主表
CREATE TABLE `user_order` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `user_id` bigint(20) NOT NULL COMMENT '用户id',
    `order_id` varchar(50) NOT NULL DEFAULT '' COMMENT 'Shopify订单ID',
    `order_name` varchar(50) NOT NULL DEFAULT '' COMMENT '订单编号（#xxx）',
    `order_created_at` bigint(20) NOT NULL DEFAULT 0 COMMENT '订单创建时间',
    `order_completion_at` bigint(20) NOT NULL DEFAULT 0 COMMENT '订单完成时间',
    `financial_status` varchar(50) NOT NULL DEFAULT '' COMMENT '支付状态',
    `total_price_amount` decimal(12,2) NOT NULL DEFAULT 0.00 COMMENT '订单总金额',
    `refund_price_amount` decimal(12,2) NOT NULL DEFAULT 0.00 COMMENT '退款总金额',
    `insurance_amount` decimal(12,2) NOT NULL DEFAULT 0.00 COMMENT '保险金额',
    `currency` varchar(10) NOT NULL DEFAULT '' COMMENT '货币类型',
    `sku_num` int(11) NOT NULL DEFAULT 0 COMMENT 'sku购买数量',
    `is_del` tinyint(1) NOT NULL DEFAULT 0 COMMENT '删除状态 0 正常 1 已删除',
    `create_time` bigint(20) NOT NULL COMMENT '创建时间',
    `update_time` bigint(20) NOT NULL COMMENT '修改时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_order_id` (`order_id`),
    KEY `idx_user_id_del` (`user_id`, `is_del`),
    KEY `idx_order_created_at` (`order_created_at`),
    KEY `idx_financial_status` (`financial_status`),
    KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户订单主表';

-- 订单详情表
CREATE TABLE `user_order_info` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `user_id` bigint(20) NOT NULL COMMENT '用户id',
    `user_order_id` bigint(20) NOT NULL COMMENT '主表ID',
    `sku` varchar(100) NOT NULL DEFAULT '' COMMENT 'SKU',
    `variant_id` varchar(100) NOT NULL DEFAULT '' COMMENT '变体ID',
    `variant_title` varchar(255) NOT NULL DEFAULT '' COMMENT '变体标题',
    `quantity` int(11) NOT NULL DEFAULT 0 COMMENT '购买数量',
    `unit_price_amount` decimal(12,2) NOT NULL DEFAULT 0.00 COMMENT '单价金额',
    `currency` varchar(10) NOT NULL DEFAULT '' COMMENT '货币类型',
    `refund_num` int(11) NOT NULL DEFAULT 0 COMMENT '退款数量',
    `is_insurance` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否是保险产品',
    `create_time` bigint(20) NOT NULL COMMENT '创建时间',
    `update_time` bigint(20) NOT NULL COMMENT '修改时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_user_order_id` (`user_order_id`),
    KEY `idx_sku` (`sku`),
    KEY `idx_variant_id` (`variant_id`),
    KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单详情表';

-- 用户订单记录统计表
CREATE TABLE `order_summary` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `user_id` bigint(20) NOT NULL COMMENT '用户ID',
    `today` int(11) NOT NULL DEFAULT 0 COMMENT '当天0点时间戳',
    `orders` int(11) NOT NULL DEFAULT 0 COMMENT '订单数',
    `sales` decimal(12,2) NOT NULL DEFAULT 0.00 COMMENT '销售金额',
    `refund` decimal(12,2) NOT NULL DEFAULT 0.00 COMMENT '退款金额',
    `create_time` bigint(20) NOT NULL COMMENT '创建时间',
    `update_time` bigint(20) NOT NULL COMMENT '修改时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_id_today` (`user_id`, `today`),
    KEY `idx_today` (`today`),
    KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户订单记录统计表';

-- 用户订单同步记录表
CREATE TABLE `job_order` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `order_id` varchar(100) NOT NULL DEFAULT '' COMMENT 'shopify 订单id',
    `name` varchar(100) NOT NULL DEFAULT '' COMMENT '店铺name',
    `job_time` bigint(20) NOT NULL COMMENT '队列时间(毫秒时间戳)',
    `is_success` tinyint(1) NOT NULL DEFAULT 0 COMMENT '处理状态 0 未处理完成 1 处理成功',
    `create_time` bigint(20) NOT NULL COMMENT '创建时间',
    `update_time` bigint(20) NOT NULL COMMENT '修改时间',
    PRIMARY KEY (`id`),
    KEY `idx_order_id` (`order_id`),
    KEY `idx_job_time_status` (`job_time`, `is_success`),
    KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户订单同步记录表';

-- 用户上传记录表
CREATE TABLE `job_product` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `user_id` bigint(20) NOT NULL COMMENT '用户id',
    `user_product_id` bigint(20) NOT NULL COMMENT '用户产品ID',
    `job_time` bigint(20) NOT NULL COMMENT '队列时间(毫秒时间戳)',
    `is_success` tinyint(1) NOT NULL DEFAULT 0 COMMENT '处理状态 0 未处理完成 1 处理成功',
    `create_time` bigint(20) NOT NULL COMMENT '创建时间',
    `update_time` bigint(20) NOT NULL COMMENT '修改时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_user_product_id` (`user_product_id`),
    KEY `idx_job_time_status` (`job_time`, `is_success`),
    KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户上传记录表';