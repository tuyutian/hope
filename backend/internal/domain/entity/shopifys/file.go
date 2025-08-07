package shopifys

type FileCreateInputContentType string

var (
	ImageType         FileCreateInputContentType = "IMAGE"
	ExternalVideoType FileCreateInputContentType = "EXTERNAL_VIDEO"
	VideoType         FileCreateInputContentType = "VIDEO"
	FileType          FileCreateInputContentType = "FILE"
	Model3DType       FileCreateInputContentType = "MODEL_3D"
)

type FileCreateInputDuplicateMode string

var (
	AppendUUidMode FileCreateInputDuplicateMode = "APPEND_UUID"
	RiseErrorMode  FileCreateInputDuplicateMode = "RAISE_ERROR"
	ReplaceMode    FileCreateInputDuplicateMode = "REPLACE"
)

type FileCreateInput struct {
	// 用于描述文件的替代文本
	Alt string `json:"alt,omitempty"`

	// 文件的 MIME 类型。例如："image/png"、"video/mp4" 等
	ContentType FileCreateInputContentType `json:"contentType"`

	// How to handle if filename is already in use.
	DuplicateResolutionMode *FileCreateInputDuplicateMode `json:"duplicateResolutionMode"`

	// 文件源，可以是 base64 数据、数据URL 或远程 URL
	OriginalSource string `json:"originalSource"`

	// 文件名，包括扩展名。可选
	Filename *string `json:"filename,omitempty"`
}
type FileUpdateInput struct {
	ID                 string   `json:"id"`
	Alt                string   `json:"alt,omitempty"`
	FileName           string   `json:"filename,omitempty"`
	OriginalSource     string   `json:"originalSource,omitempty"`
	PreviewImageSource string   `json:"previewImageSource,omitempty"`
	ReferencesToAdd    []string `json:"referencesToAdd,omitempty"`    //要添加到文件的参考文献 ID。目前仅接受产品 ID。
	ReferencesToRemove []string `json:"referencesToRemove,omitempty"` //要从文件中删除的参考 ID。目前仅接受产品 ID。
}

// FileCreateResponse 表示文件创建响应
type FileCreateResponse struct {
	FileCreate struct {
		Files      []FileCreated `json:"files"`
		UserErrors []UserError   `json:"userErrors"`
	} `json:"fileCreate"`
}

// FileUpdateResponse 表示文件创建响应
type FileUpdateResponse struct {
	FileUpdate struct {
		Files      []FileUpdated `json:"files"`
		UserErrors []UserError   `json:"userErrors"`
	} `json:"fileUpdate"`
}

type FileCreated struct {
	ID         string      `json:"id"`
	Alt        *string     `json:"alt,omitempty"`
	FileStatus string      `json:"fileStatus"`
	FileErrors []FileError `json:"fileErrors"`
	UpdatedAt  string      `json:"updatedAt,omitempty"`
	CreatedAt  string      `json:"createdAt,omitempty"`
}
type FileUpdated struct {
	ID         string      `json:"id"`
	Alt        *string     `json:"alt,omitempty"`
	FileStatus string      `json:"fileStatus"`
	Preview    FilePreview `json:"preview,omitempty"`
	FileErrors []FileError `json:"fileErrors"`
	UpdatedAt  string      `json:"updatedAt,omitempty"`
	CreatedAt  string      `json:"createdAt,omitempty"`
}

type ImageMedia struct {
	ID         string      `json:"id"`
	Alt        *string     `json:"alt,omitempty"`
	FileStatus string      `json:"fileStatus"`
	FileErrors []FileError `json:"fileErrors"`
	Image      *Image      `json:"image,omitempty"`
	UpdatedAt  string      `json:"updatedAt,omitempty"`
	CreatedAt  string      `json:"createdAt,omitempty"`
}
type FileError struct {
	Details string `json:"details"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// FilePreview  FAILED PROCESSING READY UPLOADED
type FilePreview struct {
	Image  *Image `json:"image,omitempty"`
	Status string `json:"status"`
}

// Image 预览图片（部分 File 类型会有该字段，如图片/视频封面等）
type Image struct {
	AltText *string `json:"altText,omitempty"`
	URL     string  `json:"url"`
	Width   int     `json:"width"`
	Height  int     `json:"height"`
	ID      string  `json:"id"`
}
type StagedUploadInput struct {
	Filename   string `json:"filename"`
	MimeType   string `json:"mimeType"`
	Resource   string `json:"resource"`   // 例如 "PRODUCT_IMAGE"
	HttpMethod string `json:"httpMethod"` // "POST"
}
type StagedUploadParameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type StagedTarget struct {
	URL         string                  `json:"url"`
	ResourceURL string                  `json:"resourceUrl"`
	Parameters  []StagedUploadParameter `json:"parameters"`
}

type StagedUploadsCreateResponse struct {
	StagedUploadsCreate struct {
		StagedTargets []StagedTarget `json:"stagedTargets"`
		UserErrors    []UserError    `json:"userErrors"`
	} `json:"stagedUploadsCreate"`
}
