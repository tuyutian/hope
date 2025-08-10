package shopifys

// ProductCreateInput 结构保持不变（不包含 media 字段）
type ProductCreateInput struct {
	Title           string               `json:"title,omitempty"`
	BodyHTML        string               `json:"bodyHtml,omitempty"`
	Category        string               `json:"category,omitempty"`
	DescriptionHtml string               `json:"descriptionHtml,omitempty"`
	Vendor          string               `json:"vendor,omitempty"`
	ProductType     string               `json:"productType,omitempty"`
	Handle          string               `json:"handle,omitempty"`
	Status          string               `json:"status,omitempty"`
	Tags            []string             `json:"tags,omitempty"`
	ProductOptions  []ProductOptionInput `json:"productOptions,omitempty"`
	Variants        []VariantCreateInput `json:"variants,omitempty"`
	Metafields      []MetafieldInput     `json:"metafields,omitempty"`
	SEO             *SEOInput            `json:"seo,omitempty"`
}
type ProductOptionInput struct {
	Name     string                   `json:"name"`
	Position int                      `json:"position"`
	Values   []OptionValueCreateInput `json:"values,omitempty"`
}

// SEOInput SEO输入结构
type SEOInput struct {
	Title       string `json:"title,omitempty"`       // SEO标题
	Description string `json:"description,omitempty"` // SEO描述
}

// PrivateMetafieldInput 私有元字段输入结构
type PrivateMetafieldInput struct {
	Namespace  string                     `json:"namespace"`  // 命名空间
	Key        string                     `json:"key"`        // 键
	ValueInput PrivateMetafieldValueInput `json:"valueInput"` // 值
}

// PrivateMetafieldValueInput 私有元字段值输入结构
type PrivateMetafieldValueInput struct {
	Value     string `json:"value"`     // 值
	ValueType string `json:"valueType"` // 值类型
}

// InventoryItemInput 完善的库存项目输入结构
type InventoryItemInput struct {
	SKU                          string                             `json:"sku,omitempty"`                          // SKU
	Tracked                      bool                               `json:"tracked"`                                // 是否跟踪库存
	RequiresShipping             bool                               `json:"requiresShipping,omitempty"`             // 是否需要配送
	Cost                         string                             `json:"cost,omitempty"`                         // 成本
	CountryCodeOfOrigin          string                             `json:"countryCodeOfOrigin,omitempty"`          // 原产国代码
	ProvinceCodeOfOrigin         string                             `json:"provinceCodeOfOrigin,omitempty"`         // 原产省份代码
	HarmonizedSystemCode         string                             `json:"harmonizedSystemCode,omitempty"`         // 协调系统代码
	CountryHarmonizedSystemCodes []CountryHarmonizedSystemCodeInput `json:"countryHarmonizedSystemCodes,omitempty"` // 国家协调系统代码
}

// CountryHarmonizedSystemCodeInput 国家协调系统代码输入结构
type CountryHarmonizedSystemCodeInput struct {
	CountryCode          string `json:"countryCode"`          // 国家代码
	HarmonizedSystemCode string `json:"harmonizedSystemCode"` // 协调系统代码
}

// CreateMediaInput 最新版本的媒体输入结构
type CreateMediaInput struct {
	Alt              string                     `json:"alt,omitempty"`    // 替代文本
	MediaContentType FileCreateInputContentType `json:"mediaContentType"` // 媒体类型：IMAGE, VIDEO, EXTERNAL_VIDEO, MODEL_3D
	OriginalSource   string                     `json:"originalSource"`   // 图片URL或上传后的资源URL
}

// OptionCreateInput 选项创建输入结构（完善版）
type OptionCreateInput struct {
	Name     string                   `json:"name"`               // 选项名称
	Position int                      `json:"position,omitempty"` // 位置
	Values   []OptionValueCreateInput `json:"values"`             // 选项值
}

// OptionValueCreateInput 选项值创建输入结构（完善版）
type OptionValueCreateInput struct {
	Name                 string `json:"name,omitempty"`                 // 值名称
	LinkedMetafieldValue string `json:"linkedMetafieldValue,omitempty"` // 关联的元字段值
}

// MetafieldInput 元字段输入结构（完善版）
type MetafieldInput struct {
	Namespace   string `json:"namespace"`             // 命名空间
	Key         string `json:"key"`                   // 键
	Value       string `json:"value"`                 // 值
	Type        string `json:"type,omitempty"`        // 类型
	Description string `json:"description,omitempty"` // 描述
}

type ProductUpdateInput struct {
	Id              string   `json:"id"`
	Status          string   `json:"status,omitempty"`
	Category        string   `json:"category,omitempty"`
	Title           string   `json:"title,omitempty"`
	DescriptionHtml string   `json:"descriptionHtml,omitempty"`
	ProductType     string   `json:"productType,omitempty"`
	Tags            []string `json:"tags,omitempty"`
	Vendor          string   `json:"vendor,omitempty"`
}

type MutationProduct struct {
	Product struct {
		ID        string `json:"id"`
		Title     string `json:"title"`
		Handle    string `json:"handle"`
		Status    string `json:"status"`
		CreatedAt string `json:"createdAt"`

		// 新版本使用 media 字段
		Media struct {
			Nodes []MediaNode `json:"nodes"`
		} `json:"media,omitempty"`

		Variants struct {
			Nodes []struct {
				ID    string `json:"id"`
				SKU   string `json:"sku"`
				Price string `json:"price"`
			} `json:"nodes"`
		} `json:"variants,omitempty"`
	} `json:"product"`
}

// ProductCreateResponse 响应结构也需要更新
type ProductCreateResponse struct {
	ProductCreate struct {
		MutationProduct
		UserErrors []struct {
			Field   []string `json:"field"`
			Message string   `json:"message"`
		} `json:"userErrors"`
	} `json:"productCreate"`
}

type ProductUpdateResponse struct {
	ProductUpdate struct {
		MutationProduct
		UserErrors []struct {
			Field   []string `json:"field"`
			Message string   `json:"message"`
		}
	} `json:"productUpdate"`
}

type Product struct {
	ID              string   `json:"id"`
	Title           string   `json:"title"`
	Handle          string   `json:"handle"`
	Status          string   `json:"status"`
	DescriptionHtml string   `json:"descriptionHtml"`
	ProductType     string   `json:"productType"`
	Vendor          string   `json:"vendor"`
	Tags            []string `json:"tags"`
	CreatedAt       string   `json:"createdAt"`
	UpdatedAt       string   `json:"updatedAt"`
	Images          struct {
		Nodes []struct {
			ID      string `json:"id"`
			URL     string `json:"url"`
			AltText string `json:"altText"`
		} `json:"nodes"`
	} `json:"images"`
	Variants struct {
		Nodes []struct {
			ID                string `json:"id"`
			Title             string `json:"title"`
			SKU               string `json:"sku"`
			Price             string `json:"price"`
			CompareAtPrice    string `json:"compareAtPrice"`
			InventoryQuantity int    `json:"inventoryQuantity"`
		} `json:"nodes"`
	} `json:"variants"`
}

type ProductResponse struct {
	Product *Product `json:"product"`
}

// MediaNode 结构
type MediaNode struct {
	ID               string `json:"id"`
	Alt              string `json:"alt"`
	MediaContentType string `json:"mediaContentType"`
	Preview          struct {
		Status string `json:"status"`
		Url    string `json:"url"`
	} `json:"preview"`
}

// TaxonomyResponse Taxonomy 查询响应
type TaxonomyResponse struct {
	Taxonomy *Taxonomy `json:"taxonomy"`
}

// Taxonomy 分类体系
type Taxonomy struct {
	ID         string                     `json:"id"`
	Name       string                     `json:"name"`
	Categories TaxonomyCategoryConnection `json:"categories"`
}

// TaxonomyCategoryConnection 分类连接
type TaxonomyCategoryConnection struct {
	Edges []TaxonomyCategoryEdge `json:"edges"`
}

type TaxonomyCategoryEdge struct {
	Node TaxonomyCategory `json:"node"`
}

// TaxonomyCategory 分类节点
type TaxonomyCategory struct {
	ID         string                      `json:"id"`
	Name       string                      `json:"name"`
	FullName   string                      `json:"fullName"`
	IsLeaf     bool                        `json:"isLeaf"`
	IsRoot     bool                        `json:"isRoot"`
	Level      int                         `json:"level"`
	ParentID   *string                     `json:"parentId"`
	Attributes TaxonomyAttributeConnection `json:"attributes"`
}

// TaxonomyAttributeConnection 属性连接
type TaxonomyAttributeConnection struct {
	Edges []TaxonomyAttributeEdge `json:"edges"`
}

type TaxonomyAttributeEdge struct {
	Node TaxonomyAttribute `json:"node"`
}

// TaxonomyAttribute 分类属性
type TaxonomyAttribute struct {
	ID     string                  `json:"id"`
	Name   string                  `json:"name"`
	Values TaxonomyValueConnection `json:"values"`
}

// TaxonomyValueConnection 属性值连接
type TaxonomyValueConnection struct {
	Edges []TaxonomyValueEdge `json:"edges"`
}

type TaxonomyValueEdge struct {
	Node TaxonomyValue `json:"node"`
}

// TaxonomyValue 属性值
type TaxonomyValue struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// SimplifiedCategory 简化的分类结构（用于返回给前端）
type SimplifiedCategory struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"fullName"`
	Level    int    `json:"level"`
	ParentID string `json:"parentId,omitempty"`
	IsLeaf   bool   `json:"isLeaf"`
}
