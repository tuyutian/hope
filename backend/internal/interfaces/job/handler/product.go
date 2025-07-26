package handler

import (
	"context"

	"github.com/hibiken/asynq"

	"backend/internal/application/jobs"
)

type ProductHandler struct {
	productService *jobs.ProductService
}

func (p *ProductHandler) HandleProduct(ctx context.Context, task *asynq.Task) error {
	err := p.productService.UploadProduct(ctx, task)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductHandler) HandleDelProduct(ctx context.Context, task *asynq.Task) error {
	err := p.productService.DelProduct(ctx, task)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductHandler) HandleShopifyProduct(ctx context.Context, task *asynq.Task) error {
	return p.productService.HandleShopifyProduct(ctx, task)
}
