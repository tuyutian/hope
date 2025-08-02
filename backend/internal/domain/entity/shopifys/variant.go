package shopifys

type ProductVariantInput struct {
	Price             string  `json:"price,omitempty"`
	CompareAtPrice    string  `json:"compareAtPrice,omitempty"`
	Sku               string  `json:"sku,omitempty"`
	InventoryQuantity int     `json:"inventoryQuantity,omitempty"`
	RequiresShipping  bool    `json:"requiresShipping,omitempty"`
	Taxable           bool    `json:"taxable,omitempty"`
	Weight            float64 `json:"weight,omitempty"`
	WeightUnit        string  `json:"weightUnit,omitempty"`
}

type VariantUpdateInput struct {
	Id           string                `json:"id"`
	Price        string                `json:"price,omitempty"`
	OptionValues []VariantOptionValues `json:"optionValues,omitempty"`
	//Barcode       string                `json:"barcode,omitempty"`
	InventoryItem InventoryItemInput `json:"inventoryItem,omitempty"`
}

// VariantOptionValues 变体选项值结构（完善版）
type VariantOptionValues struct {
	Name       string `json:"name"`                 // 值名称
	OptionName string `json:"optionName,omitempty"` // 选项名称
}

// VariantCreateInput 完善的变体创建输入结构
type VariantCreateInput struct {
	// 基本信息
	Price          string `json:"price,omitempty"`          // 价格
	CompareAtPrice string `json:"compareAtPrice,omitempty"` // 原价
	CostPerItem    string `json:"costPerItem,omitempty"`    // 成本价格

	// SKU和库存
	SKU               string             `json:"sku,omitempty"`               // SKU
	InventoryPolicy   string             `json:"inventoryPolicy,omitempty"`   // 库存策略 DENY, CONTINUE
	InventoryQuantity int                `json:"inventoryQuantity,omitempty"` // 库存数量
	InventoryItem     InventoryItemInput `json:"inventoryItem,omitempty"`     // 库存项目

	// 选项值
	OptionValues []VariantOptionValues `json:"optionValues,omitempty"` // 选项值

	// 物理属性
	Weight     float64 `json:"weight,omitempty"`     // 重量
	WeightUnit string  `json:"weightUnit,omitempty"` // 重量单位 GRAMS, KILOGRAMS, OUNCES, POUNDS

	// 配送和税务
	RequiresShipping bool   `json:"requiresShipping,omitempty"` // 是否需要配送
	Taxable          bool   `json:"taxable,omitempty"`          // 是否需要缴税
	TaxCode          string `json:"taxCode,omitempty"`          // 税务代码

	// 条码
	Barcode string `json:"barcode,omitempty"` // 条码

	// 媒体
	MediaSrc []string `json:"mediaSrc,omitempty"` // 媒体资源

	// 元字段
	Metafields        []MetafieldInput        `json:"metafields,omitempty"`        // 元字段
	PrivateMetafields []PrivateMetafieldInput `json:"privateMetafields,omitempty"` // 私有元字段
}
