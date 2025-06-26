package products

import (
	"time"
)

// UserProduct 保险用户产品表
type UserProduct struct {
	Id          int    `xorm:"pk autoincr 'id' int(11) comment('ID')" json:"id"`
	Uid         int    `xorm:"'uid' int(11) notnull comment('用户id')" json:"uid"`
	ProductId   string `xorm:"'product_id' varchar(100) default '0' comment('shopify上传成功的产品ID')" json:"product_id"`
	Title       string `xorm:"'title' varchar(255) notnull comment('标题')" json:"title"`
	ProductType string `xorm:"'product_type' varchar(255) notnull comment('产品类型')" json:"product_type"`
	Vendor      string `xorm:"'vendor' varchar(255) notnull comment('vendor')" json:"vendor"`
	Collection  string `xorm:"'collection' varchar(255) notnull comment('集合')" json:"collect"`
	Tags        string `xorm:"'tags' varchar(255) notnull comment('产品标签')" json:"tags"`
	Description string `xorm:"'description' text comment('描述')" json:"description"`
	Option1     string `xorm:"'option_1' varchar(255) notnull comment('产品属性1')" json:"option_1"`
	Option2     string `xorm:"'option_2' varchar(255) comment('产品属性2')" json:"option_2"`
	Option3     string `xorm:"'option_3' varchar(255) comment('产品属性3')" json:"option_3"`
	ImageUrl    string `xorm:"'image_url' varchar(300) notnull comment('封面图')" json:"image_url"`
	IsPublish   int    `xorm:"'is_publish' tinyint(1) default 0 notnull comment('发布Shopify：0:未发布 1:已发布 2:正在发布中 3:shopify平台已删除')" json:"is_publish"`
	PublishTime int64  `xorm:"'publish_time' int(11) default 0 notnull comment('发布时间')" json:"publish_time"`
	IsDel       int    `xorm:"'is_del' tinyint(1) default 0 notnull comment('删除状态 0 正常 1 已删除')" json:"is_del"`
	CreateTime  int64  `xorm:"'create_time' int(11) notnull comment('创建时间')" json:"create_time"`
	UpdateTime  int64  `xorm:"'update_time' int(11) notnull comment('修改时间')" json:"update_time"`
}

//func (m *UserProduct) TableName() string {
//	return "in_user_product"
//}

func (u *UserProduct) BeforeInsert() {
	now := time.Now().Unix()
	// 自动填充 创建时间、 更新时间
	u.CreateTime = now
	u.UpdateTime = now
}

func (u *UserProduct) BeforeUpdate() {
	now := time.Now().Unix()
	u.UpdateTime = now
}

// UserVariant  保险用户变体表
type UserVariant struct {
	Id            int     `xorm:"pk autoincr 'id' int(11) comment('ID')" json:"id"`
	Uid           int     `xorm:"'uid' int(11) notnull comment('用户id')" json:"uid"`
	UserProductId int     `xorm:"'user_product_id' int(11) notnull comment('保险用户产品表ID')" json:"user_product_id"`
	ProductId     string  `xorm:"'product_id' varchar(100) comment('Shopify产品ID')" json:"product_id"`
	VariantId     string  `xorm:"'variant_id' varchar(100) comment('Shopify变体ID')" json:"variant_id"`
	InventoryId   string  `xorm:"'inventory_id' varchar(100) comment('Shopify仓库ID')" json:"inventory_id"`
	SkuName       string  `xorm:"'sku_name' varchar(150) comment('SKU')" json:"sku_name"`
	ImageUrl      string  `xorm:"'image_url' varchar(300) comment('变体封面图')" json:"image_url"`
	Sku1          string  `xorm:"'sku_1' varchar(150) comment('变体属性1')" json:"sku_1"`
	Sku2          string  `xorm:"'sku_2' varchar(150) comment('变体属性2')" json:"sku_2"`
	Sku3          string  `xorm:"'sku_3' varchar(150) comment('变体属性3')" json:"sku_3"`
	Price         float64 `xorm:"'price' decimal(10,2) comment('价格设定')" json:"price"`
	CreateTime    int64   `xorm:"'create_time' int(11) notnull comment('创建时间')" json:"create_time"`
	UpdateTime    int64   `xorm:"'update_time' int(11) notnull comment('修改时间')" json:"update_time"`
}

//func (m *UserVariants) TableName() string {
//	return "in_user_variants"
//}

func (v *UserVariant) BeforeInsert() {
	now := time.Now().Unix()
	// 自动填充 创建时间、 更新时间
	v.CreateTime = now
	v.UpdateTime = now
}

func (v *UserVariant) BeforeUpdate() {
	now := time.Now().Unix()
	v.UpdateTime = now
}
