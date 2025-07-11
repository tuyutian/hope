-- Drop existing tables first
DROP TABLE IF EXISTS `redact_history`;
DROP TABLE IF EXISTS `user`;
DROP TABLE IF EXISTS `user_product`;
DROP TABLE IF EXISTS `user_variant`;
DROP TABLE IF EXISTS `user_cart_setting`;
DROP TABLE IF EXISTS `user_order`;
DROP TABLE IF EXISTS `user_order_info`;
DROP TABLE IF EXISTS `order_summary`;
DROP TABLE IF EXISTS `job_order`;
DROP TABLE IF EXISTS `job_product`;
DROP TABLE IF EXISTS `user_setting`;
DROP TABLE IF EXISTS `app_definition`;
DROP TABLE IF EXISTS `app_config`;
DROP TABLE IF EXISTS `user_app_auth`;
DROP TABLE IF EXISTS `commission_bill`;
DROP TABLE IF EXISTS `user_subscription`;

-- 用户订阅信息表
CREATE TABLE `user_subscription`
(
    `id`                        bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `user_id`                   bigint unsigned NOT NULL COMMENT '用户ID',
    `shop_domain`               varchar(100)    NOT NULL DEFAULT '' COMMENT '店铺域名',
    `subscription_id`           varchar(100)    NOT NULL DEFAULT '' COMMENT 'Shopify订阅ID',
    `subscription_name`         varchar(100)    NOT NULL DEFAULT '' COMMENT '订阅名称',
    `subscription_status`       varchar(20)     NOT NULL DEFAULT '' COMMENT '订阅状态：ACTIVE, CANCELLED, DECLINED, EXPIRED, FROZEN, PENDING',
    `subscription_line_item_id` varchar(100)    NOT NULL DEFAULT '' COMMENT '订阅项目ID（用于创建用量扣费）',
    `pricing_type`              varchar(20)     NOT NULL DEFAULT '' COMMENT '定价类型：ANNUAL, RECURRING, ONE_TIME',
    `capped_amount`             decimal(12, 2)  NOT NULL DEFAULT 0.00 COMMENT '封顶金额',
    `currency`                  varchar(10)     NOT NULL DEFAULT '' COMMENT '货币类型',
    `balance_used`              decimal(12, 2)  NOT NULL DEFAULT 0.00 COMMENT '已使用额度',
    `terms`                     text            COMMENT '计费条款',
    `current_period_start`      bigint unsigned NOT NULL DEFAULT 0 COMMENT '当前计费周期开始时间',
    `current_period_end`        bigint unsigned NOT NULL DEFAULT 0 COMMENT '当前计费周期结束时间',
    `trial_days`                int             NOT NULL DEFAULT 0 COMMENT '试用天数',
    `test_subscription`         tinyint         NOT NULL DEFAULT 0 COMMENT '是否为测试订阅：0-否, 1-是',
    `app_installation_id`       varchar(100)    NOT NULL DEFAULT '' COMMENT 'App安装ID',
    `last_sync_time`            bigint unsigned NOT NULL DEFAULT 0 COMMENT '最后同步时间',
    `create_time`               bigint unsigned NOT NULL COMMENT '创建时间',
    `update_time`               bigint unsigned NOT NULL COMMENT '修改时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_subscription` (`user_id`, `subscription_id`),
    KEY `idx_user_id_status` (`user_id`, `subscription_status`),
    KEY `idx_subscription_id` (`subscription_id`),
    KEY `idx_subscription_line_item_id` (`subscription_line_item_id`),
    KEY `idx_shop_domain` (`shop_domain`),
    KEY `idx_subscription_status` (`subscription_status`),
    KEY `idx_current_period_end` (`current_period_end`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='用户订阅信息表';

-- 抽成收费记录表
CREATE TABLE `commission_bill`
(
    `id`                     bigint unsigned         NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `charge_id`        bigint unsigned  NOT NULL COMMENT '账单编号',
    `user_id`                bigint unsigned         NOT NULL COMMENT '用户ID',
    `user_order_id`          bigint unsigned         NOT NULL COMMENT '关联的订单ID',
    `order_name`            varchar(50)    NOT NULL DEFAULT '' COMMENT 'Shopify订单编号',
    `bill_cycle`            varchar(20)    NOT NULL COMMENT '账单周期（YYYY-MM）',
    `commission_amount`     decimal(12, 2) NOT NULL DEFAULT 0.00 COMMENT '抽成金额',
    `currency`              varchar(10)    NOT NULL DEFAULT '' COMMENT '货币类型',
    `shopify_usage_record_id` varchar(100) NOT NULL DEFAULT '' COMMENT 'Shopify用量记录ID',
    `charge_status`         tinyint        NOT NULL DEFAULT 0 COMMENT '扣费状态：0-待提交, 1-已提交, 2-提交失败',
    `error_message`         text           COMMENT '错误信息',
    `charged_at`             bigint unsigned         NOT NULL DEFAULT 0 COMMENT '扣费时间',
    `create_time`            bigint unsigned         NOT NULL COMMENT '创建时间',
    `update_time`            bigint unsigned         NOT NULL COMMENT '修改时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_charge_id` (`charge_id`),
    UNIQUE KEY `uk_user_order` (`user_id`, `user_order_id`),
    KEY `idx_user_id_status` (`user_id`, `charge_status`),
    KEY `idx_bill_cycle` (`bill_cycle`),
    KEY `idx_charge_status` (`charge_status`),
    KEY `idx_shopify_usage_record_id` (`shopify_usage_record_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='抽成收费记录表';

-- Redact 历史记录表 (仅记录最小信息，用于防重复)
CREATE TABLE `redact_history`
(
    `id`          bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `app_id`      varchar(50)     NOT NULL COMMENT 'App标识',
    `shop`        varchar(100)    NOT NULL COMMENT 'Shop域名',
    `redact_time` bigint unsigned NOT NULL COMMENT 'Redact处理时间',
    `create_time` bigint unsigned NOT NULL COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_app_shop` (`app_id`, `shop`),
    KEY `idx_redact_time` (`redact_time`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='Redact历史记录表(最小化记录)';

-- 用户表 (移除 redact 字段，因为 redact 的用户会被物理删除)
CREATE TABLE `user`
(
    `id`                bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `app_id`            varchar(50)     NOT NULL COMMENT 'App标识',
    `name`              varchar(255)    NOT NULL COMMENT 'shopify-name',
    `shop`              varchar(100)    NOT NULL COMMENT 'shopify域名',
    `real_domain`       varchar(100)    NOT NULL DEFAULT '' COMMENT '网站真实域名',
    `plan_display_name` varchar(40)     NOT NULL DEFAULT '' COMMENT 'shopify套餐版本',
    `access_token`      varchar(512)    NOT NULL DEFAULT '' COMMENT 'shopify-token',
    `password`          varchar(255)    NOT NULL DEFAULT '' COMMENT '密码',
    `plans`             int             NOT NULL DEFAULT 0 COMMENT 'app套餐id',
    `email`             varchar(100)    NOT NULL DEFAULT '' COMMENT '邮箱',
    `phone`             varchar(20)     NOT NULL DEFAULT '' COMMENT '电话号码',
    `country_name`      varchar(50)     NOT NULL DEFAULT '' COMMENT '国家名称',
    `country_code`      varchar(5)      NOT NULL DEFAULT '' COMMENT '国家简码',
    `city`              varchar(50)     NOT NULL DEFAULT '' COMMENT '城市',
    `free_trial_days`   tinyint         NOT NULL DEFAULT 0 COMMENT '试用天数',
    `trial_time`        bigint unsigned NOT NULL DEFAULT 0 COMMENT '试用时间',
    `currency_code`     varchar(10)     NOT NULL DEFAULT '' COMMENT '货币简码',
    `timezone`          int             NOT NULL DEFAULT 0 COMMENT '时区偏转分钟',
    `money_format`      varchar(20)     NOT NULL DEFAULT '' COMMENT '货币单位符号',
    `last_login`        bigint unsigned NOT NULL DEFAULT 0 COMMENT '最后登录时间',
    `is_del`            tinyint         NOT NULL DEFAULT 0 COMMENT '删除状态 0正常 1已删除',
    `publish_id`        bigint unsigned NOT NULL DEFAULT 0 COMMENT '店铺publish_id',
    `install_time`      bigint unsigned NOT NULL DEFAULT 0 COMMENT '安装时间',
    `uninstall_time`    bigint unsigned NOT NULL DEFAULT 0 COMMENT '卸载时间',
    `create_time`       bigint unsigned NOT NULL COMMENT '创建时间',
    `update_time`       bigint unsigned NOT NULL COMMENT '最近修改时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_app_shop` (`app_id`, `shop`),
    KEY `idx_app_id` (`app_id`),
    KEY `idx_shop` (`shop`),
    KEY `idx_name` (`name`),
    KEY `idx_email` (`email`),
    KEY `idx_is_del` (`is_del`),
    KEY `idx_plans` (`plans`),
    KEY `idx_last_login` (`last_login`),
    KEY `idx_create_time` (`create_time`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='用户表';

-- 保险用户产品表
CREATE TABLE `user_product`
(
    `id`           bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `app_id`       varchar(50)     NOT NULL COMMENT 'App标识',
    `user_id`      bigint unsigned NOT NULL COMMENT '用户id',
    `product_id`   bigint unsigned NOT NULL DEFAULT 0 COMMENT 'shopify上传成功的产品ID',
    `handle`       varchar(255)    NOT NULL DEFAULT '' COMMENT '产品handle',
    `title`        varchar(255)    NOT NULL COMMENT '标题',
    `product_type` varchar(255)    NOT NULL COMMENT '产品类型',
    `vendor`       varchar(255)    NOT NULL COMMENT 'vendor',
    `collection`   varchar(100)    NOT NULL COMMENT '集合',
    `tags`         varchar(500)    NOT NULL DEFAULT '' COMMENT '产品标签',
    `description`  text COMMENT '描述',
    `option_1`     varchar(255)    NOT NULL COMMENT '产品属性1',
    `option_2`     varchar(255)    NOT NULL DEFAULT '' COMMENT '产品属性2',
    `option_3`     varchar(255)    NOT NULL DEFAULT '' COMMENT '产品属性3',
    `image_url`    varchar(500)    NOT NULL COMMENT '封面图',
    `status`       tinyint         NOT NULL DEFAULT 0 COMMENT '发布Shopify：0:未发布 1:已发布 2:正在发布中 3:shopify平台已删除',
    `publish_time` bigint unsigned NOT NULL DEFAULT 0 COMMENT '发布时间',
    `is_del`       tinyint         NOT NULL DEFAULT 0 COMMENT '删除状态 0 正常 1 已删除',
    `create_time`  bigint unsigned NOT NULL COMMENT '创建时间',
    `update_time`  bigint unsigned NOT NULL COMMENT '修改时间',
    PRIMARY KEY (`id`),
    KEY `idx_app_id` (`app_id`),
    KEY `idx_user_id_del` (`user_id`, `is_del`),
    KEY `idx_app_user_status` (`app_id`, `user_id`, `status`, `is_del`),
    KEY `idx_product_id` (`product_id`),
    KEY `idx_handle` (`handle`),
    KEY `idx_product_type` (`product_type`),
    KEY `idx_status` (`status`),
    KEY `idx_publish_time` (`publish_time`),
    KEY `idx_create_time` (`create_time`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='用户产品表';

-- 保险用户变体表
CREATE TABLE `user_variant`
(
    `id`              bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `user_id`         bigint unsigned NOT NULL COMMENT '用户id',
    `user_product_id` bigint unsigned NOT NULL COMMENT '保险用户产品表ID',
    `product_id`      bigint unsigned NOT NULL DEFAULT 0 COMMENT 'Shopify产品ID',
    `variant_id`      bigint unsigned NOT NULL DEFAULT 0 COMMENT 'Shopify变体ID',
    `inventory_id`    bigint unsigned NOT NULL DEFAULT 0 COMMENT 'Shopify仓库ID',
    `sku_name`        varchar(150)    NOT NULL DEFAULT '' COMMENT 'SKU',
    `image_url`       varchar(500)    NOT NULL DEFAULT '' COMMENT '变体封面图',
    `sku_1`           varchar(150)    NOT NULL DEFAULT '' COMMENT '变体属性1',
    `sku_2`           varchar(150)    NOT NULL DEFAULT '' COMMENT '变体属性2',
    `sku_3`           varchar(150)    NOT NULL DEFAULT '' COMMENT '变体属性3',
    `price`           decimal(12, 2)  NOT NULL DEFAULT 0.00 COMMENT '价格设定',
    `create_time`     bigint unsigned NOT NULL COMMENT '创建时间',
    `update_time`     bigint unsigned NOT NULL COMMENT '修改时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_user_product_id` (`user_product_id`),
    KEY `idx_product_id` (`product_id`),
    KEY `idx_variant_id` (`variant_id`),
    KEY `idx_user_product_variant` (`user_product_id`, `variant_id`),
    KEY `idx_create_time` (`create_time`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='保险用户变体表';

-- 保险用户基础配置表
CREATE TABLE `user_cart_setting`
(
    `id`                 bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `user_id`            bigint unsigned NOT NULL COMMENT '用户id',
    `plan_title`         varchar(100)    NOT NULL DEFAULT '' COMMENT '保险标题(内部)',
    `addon_title`        varchar(100)    NOT NULL DEFAULT '' COMMENT '保险标题',
    `enabled_desc`       varchar(200)    NOT NULL DEFAULT '' COMMENT '按钮打开文案',
    `disabled_desc`      varchar(200)    NOT NULL DEFAULT '' COMMENT '按钮关闭文案',
    `foot_text`          varchar(100)    NOT NULL DEFAULT '' COMMENT '保险底部',
    `foot_url`           varchar(255)    NOT NULL DEFAULT '' COMMENT '保险跳转',
    `in_color`           varchar(50)     NOT NULL DEFAULT '' COMMENT '打开颜色',
    `out_color`          varchar(50)     NOT NULL DEFAULT '' COMMENT '关闭颜色',
    `other_money`        decimal(12, 2)  NOT NULL DEFAULT 0.00 COMMENT '其他金额',
    `show_cart`          tinyint         NOT NULL DEFAULT 0 COMMENT '购物车状态 0 关闭 1 打开',
    `show_cart_icon`     tinyint         NOT NULL DEFAULT 0 COMMENT '购物车图标 0 关闭 1 打开',
    `icon_url`           text COMMENT '选中url(json)',
    `select_button`      tinyint         NOT NULL DEFAULT 0 COMMENT '购物车图标 0 滑动 1 勾选',
    `product_type`       varchar(100)    NOT NULL DEFAULT '' COMMENT '产品type',
    `product_collection` varchar(100)    NOT NULL DEFAULT '' COMMENT '产品选中集合',
    `pricing_type`       tinyint         NOT NULL DEFAULT 0 COMMENT '购物车图标 0 金额 1百分比',
    `pricing_select`     text COMMENT '金额计算范围',
    `tiers_select`       text COMMENT '百分比计算范围',
    `create_time`        bigint unsigned NOT NULL COMMENT '创建时间',
    `update_time`        bigint unsigned NOT NULL COMMENT '修改时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_id` (`user_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='保险用户基础配置表';

-- 用户订单主表
CREATE TABLE `user_order`
(
    `id`                  bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `user_id`             bigint unsigned NOT NULL COMMENT '用户id',
    `order_id`            bigint unsigned NOT NULL DEFAULT 0 COMMENT 'Shopify订单ID',
    `order_name`          varchar(50)     NOT NULL DEFAULT '' COMMENT '订单编号（#xxx）',
    `order_created_at`    bigint unsigned NOT NULL DEFAULT 0 COMMENT '订单创建时间',
    `order_completion_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT '订单完成时间',
    `financial_status`    varchar(50)     NOT NULL DEFAULT '' COMMENT '支付状态',
    `total_price_amount`  decimal(12, 2)  NOT NULL DEFAULT 0.00 COMMENT '订单总金额',
    `refund_price_amount` decimal(12, 2)  NOT NULL DEFAULT 0.00 COMMENT '退款总金额',
    `insurance_amount`    decimal(12, 2)  NOT NULL DEFAULT 0.00 COMMENT '保险金额',
    `currency`            varchar(10)     NOT NULL DEFAULT '' COMMENT '货币类型',
    `sku_num`             int             NOT NULL DEFAULT 0 COMMENT 'sku购买数量',
    `is_del`              tinyint         NOT NULL DEFAULT 0 COMMENT '删除状态 0 正常 1 已删除',
    `create_time`         bigint unsigned NOT NULL COMMENT '创建时间',
    `update_time`         bigint unsigned NOT NULL COMMENT '修改时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_order_id` (`order_id`),
    KEY `idx_user_id_del` (`user_id`, `is_del`),
    KEY `idx_user_financial_del` (`user_id`, `financial_status`, `is_del`),
    KEY `idx_order_created_at` (`order_created_at`),
    KEY `idx_financial_status` (`financial_status`),
    KEY `idx_create_time` (`create_time`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='用户订单主表';

-- 订单详情表
CREATE TABLE `user_order_info`
(
    `id`                bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `user_id`           bigint unsigned NOT NULL COMMENT '用户id',
    `user_order_id`     bigint unsigned NOT NULL COMMENT '主表ID',
    `sku`               varchar(100)    NOT NULL DEFAULT '' COMMENT 'SKU',
    `variant_id`        bigint unsigned NOT NULL DEFAULT 0 COMMENT '变体ID',
    `variant_title`     varchar(255)    NOT NULL DEFAULT '' COMMENT '变体标题',
    `quantity`          int             NOT NULL DEFAULT 0 COMMENT '购买数量',
    `unit_price_amount` decimal(12, 2)  NOT NULL DEFAULT 0.00 COMMENT '单价金额',
    `currency`          varchar(10)     NOT NULL DEFAULT '' COMMENT '货币类型',
    `refund_num`        int             NOT NULL DEFAULT 0 COMMENT '退款数量',
    `is_insurance`      tinyint         NOT NULL DEFAULT 0 COMMENT '是否是保险产品',
    `create_time`       bigint unsigned NOT NULL COMMENT '创建时间',
    `update_time`       bigint unsigned NOT NULL COMMENT '修改时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_user_order_id` (`user_order_id`),
    KEY `idx_order_insurance` (`user_order_id`, `is_insurance`),
    KEY `idx_sku` (`sku`),
    KEY `idx_variant_id` (`variant_id`),
    KEY `idx_create_time` (`create_time`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='订单详情表';

-- 用户订单记录统计表
CREATE TABLE `order_summary`
(
    `id`          bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `user_id`     bigint unsigned NOT NULL COMMENT '用户ID',
    `today`       bigint unsigned NOT NULL DEFAULT 0 COMMENT '当天0点时间戳',
    `orders`      int             NOT NULL DEFAULT 0 COMMENT '订单数',
    `sales`       decimal(12, 2)  NOT NULL DEFAULT 0.00 COMMENT '销售金额',
    `refund`      decimal(12, 2)  NOT NULL DEFAULT 0.00 COMMENT '退款金额',
    `create_time` bigint unsigned NOT NULL COMMENT '创建时间',
    `update_time` bigint unsigned NOT NULL COMMENT '修改时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_id_today` (`user_id`, `today`),
    KEY `idx_today` (`today`),
    KEY `idx_create_time` (`create_time`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='用户订单记录统计表';

-- 用户订单同步记录表
CREATE TABLE `job_order`
(
    `id`          bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `order_id`    bigint unsigned NOT NULL DEFAULT 0 COMMENT 'shopify 订单id',
    `name`        varchar(100)    NOT NULL DEFAULT '' COMMENT '店铺name',
    `job_time`    bigint unsigned NOT NULL COMMENT '队列时间(毫秒时间戳)',
    `is_success`  tinyint         NOT NULL DEFAULT 0 COMMENT '处理状态 0 未处理完成 1 处理成功',
    `create_time` bigint unsigned NOT NULL COMMENT '创建时间',
    `update_time` bigint unsigned NOT NULL COMMENT '修改时间',
    PRIMARY KEY (`id`),
    KEY `idx_order_id` (`order_id`),
    KEY `idx_job_time_status` (`job_time`, `is_success`),
    KEY `idx_name` (`name`),
    KEY `idx_create_time` (`create_time`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='用户订单同步记录表';

-- 用户上传记录表
CREATE TABLE `job_product`
(
    `id`              bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `user_id`         bigint unsigned NOT NULL COMMENT '用户id',
    `user_product_id` bigint unsigned NOT NULL COMMENT '用户产品ID',
    `job_time`        bigint unsigned NOT NULL COMMENT '队列时间(毫秒时间戳)',
    `is_success`      tinyint         NOT NULL DEFAULT 0 COMMENT '处理状态 0 未处理完成 1 处理成功',
    `create_time`     bigint unsigned NOT NULL COMMENT '创建时间',
    `update_time`     bigint unsigned NOT NULL COMMENT '修改时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_user_product_id` (`user_product_id`),
    KEY `idx_job_time_status` (`job_time`, `is_success`),
    KEY `idx_create_time` (`create_time`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='用户上传记录表';

-- 用户自定义设置表
CREATE TABLE `user_setting`
(
    `id`          bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `user_id`     bigint unsigned NOT NULL COMMENT '用户id',
    `label`       varchar(255)    NOT NULL COMMENT '自定义设置名',
    `name`        varchar(255)    NOT NULL COMMENT '自定义设置键',
    `value`       text            NOT NULL COMMENT '配置值(JSON格式)',
    `create_time` bigint unsigned NOT NULL COMMENT '创建时间',
    `update_time` bigint unsigned NOT NULL COMMENT '修改时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_name` (`user_id`, `name`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_create_time` (`create_time`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='用户上传记录表';


-- 用户对应用授权表
CREATE TABLE `user_app_auth`
(
    `id`               bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `user_id`          bigint unsigned NOT NULL COMMENT '用户ID',
    `app_id`           varchar(50)     NOT NULL COMMENT 'App标识',
    `shop`             varchar(100)    NOT NULL DEFAULT '' COMMENT 'my shopify Domain网站域名',
    `auth_token`       varchar(255)    NOT NULL DEFAULT '' COMMENT '授权token',
    `refresh_token`    varchar(255)    NOT NULL DEFAULT '' COMMENT '刷新token',
    `token_expires_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT 'token过期时间',
    `scopes`           text            NOT NULL COMMENT '授权域',
    `status`           tinyint         NOT NULL DEFAULT 1 COMMENT '状态 1:有效 0:已撤销',
    `create_time`      bigint unsigned NOT NULL COMMENT '创建时间',
    `update_time`      bigint unsigned NOT NULL COMMENT '修改时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_app` (`user_id`, `app_id`),
    KEY `idx_app_id` (`app_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_status` (`status`),
    KEY `idx_create_time` (`create_time`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='用户应用授权表';

-- App 配置表
CREATE TABLE `app_config`
(
    `id`           bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `app_id`       varchar(50)     NOT NULL COMMENT 'App标识',
    `config_key`   varchar(100)    NOT NULL COMMENT '配置键',
    `config_value` text            NOT NULL COMMENT '配置值',
    `description`  varchar(255)    NOT NULL DEFAULT '' COMMENT '描述',
    `create_time`  bigint unsigned NOT NULL COMMENT '创建时间',
    `update_time`  bigint unsigned NOT NULL COMMENT '修改时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_app_key` (`app_id`, `config_key`),
    KEY `idx_app_id` (`app_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='App配置表';


-- App 定义表
CREATE TABLE `app_definition`
(
    `id`           bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `app_id`       varchar(50)     NOT NULL COMMENT 'App唯一标识',
    `name`         varchar(100)    NOT NULL COMMENT 'App名称',
    `description`  text COMMENT 'App描述',
    `icon_url`     varchar(255)    NOT NULL DEFAULT '' COMMENT 'App图标',
    `callback_url` varchar(255)    NOT NULL DEFAULT '' COMMENT '回调URL',
    `api_key`      varchar(100)    NOT NULL COMMENT 'API Key',
    `api_secret`   varchar(100)    NOT NULL COMMENT 'API Secret',
    `scopes`       text            NOT NULL COMMENT '授权域',
    `status`       tinyint         NOT NULL DEFAULT 1 COMMENT '状态 1:启用 0:禁用',
    `create_time`  bigint unsigned NOT NULL COMMENT '创建时间',
    `update_time`  bigint unsigned NOT NULL COMMENT '修改时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_app_id` (`app_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='App定义表';
