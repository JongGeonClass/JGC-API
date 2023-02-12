package database

import (
	"context"
	"time"

	"github.com/JongGeonClass/JGC-API/dbmodel"
	"github.com/thak1411/gorn"
	"github.com/thak1411/rnlog"
)

// 상품 디비의 인터페이스 입니다.
type ProductDatabase interface {
	ExecTx(ctx context.Context, fn func(txdb ProductDatabase) error) error
	AddProduct(ctx context.Context, product *dbmodel.Product) (int64, error)
	DeleteAllProducts(ctx context.Context) error
	GetProducts(ctx context.Context, page, pagesize int) ([]*dbmodel.PublicProduct, error)
	GetProductsCount(ctx context.Context) (int, error)
	AddBrand(ctx context.Context, brand *dbmodel.Brand) (int64, error)
	DeleteAllBrands(ctx context.Context) error
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

// 새로운 상품을 추가합니다.
// 이후 추가된 상품 아이디를 반환합니다.
func (h *ProductDB) AddProduct(ctx context.Context, product *dbmodel.Product) (int64, error) {
	ntime := time.Now()
	product.CreatedTime = ntime
	product.UpdatedTime = ntime
	return h.InsertWithLastId(ctx, "PRODUCT", product)
}

// 모든 상품을 삭제합니다.
func (h *ProductDB) DeleteAllProducts(ctx context.Context) error {
	sql := gorn.NewSql().
		DeleteFrom("PRODUCT").
		Where("id > ?", -1)
	res, err := h.Exec(ctx, sql)
	if err != nil {
		return err
	}
	if _, err := res.RowsAffected(); err != nil {
		return err
	}
	return nil
}

// 상품 목록을 가져옵니다.
//TODO: pagination 제대로 가져오기
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

// 새로운 브랜드를 추가합니다.
// 이후 추가된 브랜드 아이디를 반환합니다.
func (h *ProductDB) AddBrand(ctx context.Context, brand *dbmodel.Brand) (int64, error) {
	ntime := time.Now()
	brand.CreatedTime = ntime
	brand.UpdatedTime = ntime
	return h.InsertWithLastId(ctx, "BRAND", brand)
}

// 모든 브랜드를 삭제합니다.
func (h *ProductDB) DeleteAllBrands(ctx context.Context) error {
	sql := gorn.NewSql().
		DeleteFrom("BRAND").
		Where("id > ?", -1)
	res, err := h.Exec(ctx, sql)
	if err != nil {
		return err
	}
	if _, err := res.RowsAffected(); err != nil {
		return err
	}
	return nil
}

// 새로운 디비 객체를 연결합니다.
func NewProduct(db *gorn.DB) ProductDatabase {
	return &ProductDB{
		DB: db,
	}
}
