package shopifys

// ProductCreateInput 创建产品的输入结构
type ProductCreateInput struct {
	Title           string                `json:"title"`
	DescriptionHtml string                `json:"descriptionHtml,omitempty"`
	ProductType     string                `json:"productType,omitempty"`
	Vendor          string                `json:"vendor,omitempty"`
	Tags            []string              `json:"tags,omitempty"`
	Status          string                `json:"status,omitempty"`
	Images          []ProductImageInput   `json:"images,omitempty"`
	Variants        []ProductVariantInput `json:"variants,omitempty"`
	ProductOptions  []OptionCreateInput   `json:"options,omitempty"`
}
type ProductUpdateInput struct {
	Id              string   `json:"id"`
	Status          string   `json:"status,omitempty"`
	Title           string   `json:"title"`
	DescriptionHtml string   `json:"descriptionHtml,omitempty"`
	ProductType     string   `json:"productType,omitempty"`
	Tags            []string `json:"tags,omitempty"`
	Vendor          string   `json:"vendor,omitempty"`
}
type ProductImageInput struct {
	Src string `json:"src"`
	Alt string `json:"altText,omitempty"`
}

type ProductVariantInput struct {
	Price             string  `json:"price"`
	CompareAtPrice    string  `json:"compareAtPrice,omitempty"`
	Sku               string  `json:"sku,omitempty"`
	InventoryQuantity int     `json:"inventoryQuantity,omitempty"`
	RequiresShipping  bool    `json:"requiresShipping,omitempty"`
	Taxable           bool    `json:"taxable,omitempty"`
	Weight            float64 `json:"weight,omitempty"`
	WeightUnit        string  `json:"weightUnit,omitempty"`
}

// CreateMediaInput 用于添加图片等媒体资源
type CreateMediaInput struct {
	Alt              string `json:"alt"`
	MediaContentType string `json:"mediaContentType"`
	OriginalSource   string `json:"originalSource"`
}
type OptionCreateInput struct {
	Name     string                   `json:"name"`
	Position int                      `json:"position"`
	Values   []OptionValueCreateInput `json:"values,omitempty"`
}
type OptionValueCreateInput struct {
	LinkedMetafieldValue string `json:"linkedMetafieldValue"`
	Name                 string `json:"name"`
}

// ProductCreateResponse 响应结构
type ProductCreateResponse struct {
	ProductCreate struct {
		Product struct {
			ID        string `json:"id"`
			Title     string `json:"title"`
			Handle    string `json:"handle"`
			Status    string `json:"status"`
			CreatedAt string `json:"createdAt"`
			Images    struct {
				Edges []struct {
					Node struct {
						ID      string `json:"id"`
						URL     string `json:"url"`
						AltText string `json:"altText"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"images,omitempty"`
			Variants struct {
				Edges []struct {
					Node struct {
						ID    string `json:"id"`
						SKU   string `json:"sku"`
						Price string `json:"price"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"variants,omitempty"`
		} `json:"product"`
		UserErrors []struct {
			Field   []string `json:"field"`
			Message string   `json:"message"`
		} `json:"userErrors"`
	} `json:"productCreate"`
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
		Edges []struct {
			Node struct {
				ID      string `json:"id"`
				URL     string `json:"url"`
				AltText string `json:"altText"`
			} `json:"node"`
		} `json:"edges"`
	} `json:"images"`
	Variants struct {
		Edges []struct {
			Node struct {
				ID                string `json:"id"`
				Title             string `json:"title"`
				SKU               string `json:"sku"`
				Price             string `json:"price"`
				CompareAtPrice    string `json:"compareAtPrice"`
				InventoryQuantity int    `json:"inventoryQuantity"`
			} `json:"node"`
		} `json:"edges"`
	} `json:"variants"`
}

type ProductResponse struct {
	Product Product `json:"product"`
}

type VariantCreateInput struct {
	Price        string                `json:"price,omitempty"`
	OptionValues []VariantOptionValues `json:"optionValues,omitempty"`
	//Barcode       string                `json:"barcode,omitempty"`
	InventoryItem InventoryItemInput `json:"inventoryItem,omitempty"`
}

type VariantUpdateInput struct {
	Id           string                `json:"id"`
	Price        string                `json:"price,omitempty"`
	OptionValues []VariantOptionValues `json:"optionValues,omitempty"`
	//Barcode       string                `json:"barcode,omitempty"`
	InventoryItem InventoryItemInput `json:"inventoryItem,omitempty"`
}

type VariantOptionValues struct {
	Name       string `json:"name"`
	OptionName string `json:"optionName,omitempty"`
}

type InventoryItemInput struct {
	Sku     string `json:"sku,omitempty"`
	Tracked bool   `json:"tracked"`
}
