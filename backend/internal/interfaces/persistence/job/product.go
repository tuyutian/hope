package job

import (
	"context"
	"time"

	"backend/internal/domain/entity/jobs"
	jobRepo "backend/internal/domain/repo/jobs"

	"xorm.io/xorm"
)

type ProductRepoImpl struct {
	db *xorm.Engine
}

func NewProductRepository(db *xorm.Engine) jobRepo.ProductRepository {
	return &ProductRepoImpl{db: db}
}

func (j *ProductRepoImpl) First(ctx context.Context, id int) (*jobs.JobProduct, error) {
	var jobProduct jobs.JobProduct
	has, err := j.db.Context(ctx).Where("id = ?", id).Get(&jobProduct)

	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}

	return &jobProduct, nil
}

func (j *ProductRepoImpl) Create(ctx context.Context, jobProduct *jobs.JobProduct) (int, error) {
	_, err := j.db.Context(ctx).Insert(jobProduct)
	if err != nil {
		return 0, err
	}

	return jobProduct.Id, nil
}

func (j *ProductRepoImpl) UpdateJobTime(ctx context.Context, id int) error {
	_, err := j.db.Context(ctx).
		Where("id = ?", id).
		Update(&jobs.JobProduct{JobTime: time.Now().Unix()})
	if err != nil {
		return err
	}
	return nil
}

func (j *ProductRepoImpl) UpdateStatus(ctx context.Context, id int, status int) error {
	_, err := j.db.Context(ctx).
		Where("id = ?", id).
		Update(&jobs.JobProduct{IsSuccess: status})
	if err != nil {
		return err
	}
	return nil
}

func (j *ProductRepoImpl) Clear(ctx context.Context) error {
	_, _ = j.db.Context(ctx).Where("is_success > ?", 0).Delete(&jobs.JobProduct{})
	return nil
}
