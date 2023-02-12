package database

import (
	"context"

	"github.com/JongGeonClass/JGC-API/dbmodel"
	"github.com/thak1411/gorn"
	"github.com/thak1411/rnlog"
)

// 상품 디비의 인터페이스 입니다.
type ProductDatabase interface {
	ExecTx(ctx context.Context, fn func(txdb ProductDatabase) error) error
	GetProducts(ctx context.Context, page, pagesize int) ([]*dbmodel.PublicProduct, error)
	GetProductsCount(ctx context.Context) (int, error)
}

// 상품 디비의 구현체입니다.
type ProductDB struct {
	*gorn.DB
}

// 넘겨받은 함수로 트랜잭션을 실행합니다.
func (h *ProductDB) ExecTx(ctx context.Context, fn func(txdb ProductDatabase) error) error {
	txdb, err := h.DB.BeginTx(ctx)
	if err != nil {
		return err
	}
	newHandler := &ProductDB{
		DB: txdb,
	}
	err = fn(newHandler)
	if err != nil {
		if rbErr := txdb.RollbackTx(); rbErr != nil {
			rnlog.Error("Rollback error: %v", rbErr)
			return rbErr
		}
		return err
	}
	return txdb.CommitTx()
}

// 상품 목록을 가져옵니다.
func (h *ProductDB) GetProducts(ctx context.Context, page, pagesize int) ([]*dbmodel.PublicProduct, error) {
	result := []*dbmodel.PublicProduct{}
	sql := gorn.NewSql().
		Select(&dbmodel.PublicProduct{}).
		From("PRODUCT").
		InnerJoin("BRAND").
		On("PRODUCT.brand_id = BRAND.id").
		OrderBy("PRODUCT.id").DESC().
		LimitPage(int64(page), int64(pagesize))

	rows, err := h.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if err := h.ScanRows(rows, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// 상품 개수를 가져옵니다.
func (h *ProductDB) GetProductsCount(ctx context.Context) (int, error) {
	type ProductsCount struct {
		Count int `rnsql:"COUNT(*)"`
	}
	result := &ProductsCount{}
	sql := gorn.NewSql().
		Select(result).
		From("PRODUCT")
	row := h.QueryRow(ctx, sql)
	if err := h.ScanRow(row, result); err != nil {
		return 0, err
	}
	return result.Count, nil
}

// 새로운 디비 객체를 연결합니다.
func NewProduct(db *gorn.DB) ProductDatabase {
	return &ProductDB{
		DB: db,
	}
}
