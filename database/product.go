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
	CheckProductExists(ctx context.Context, productId int64) (bool, error)
	AddBrand(ctx context.Context, brand *dbmodel.Brand) (int64, error)
	DeleteAllBrands(ctx context.Context) error
	AddCategory(ctx context.Context, category *dbmodel.Category) (int64, error)
	DeleteAllCategories(ctx context.Context) error
	AddProductCategory(ctx context.Context, productCategoryMap *dbmodel.ProductCategoryMap) error
	DeleteAllProductCategoryMap(ctx context.Context) error
	AddCart(ctx context.Context, cart *dbmodel.Cart) error
	CheckCartHasProduct(ctx context.Context, userId, productId int64) (bool, error)
	GetCartProduct(ctx context.Context, userId, productId int64) (*dbmodel.Cart, error)
	UpdateCart(ctx context.Context, cart *dbmodel.Cart) error
	DeleteCartProduct(ctx context.Context, userId, productId int64) error
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
func (h *ProductDB) GetProducts(ctx context.Context, page, pagesize int) ([]*dbmodel.PublicProduct, error) {
	result := []*dbmodel.PublicProduct{}
	sql := gorn.NewSql().
		Select(&dbmodel.PublicProduct{}).
		From("PRODUCT").
		InnerJoin("BRAND").
		On("PRODUCT.brand_id = BRAND.id").
		InnerJoin("PRODUCT_CATEGORY_MAP").
		On("PRODUCT_CATEGORY_MAP.product_id = PRODUCT.id").
		InnerJoin("CATEGORY").
		On("PRODUCT_CATEGORY_MAP.category_id = CATEGORY.id").
		AddPlainQuery("GROUP BY PRODUCT.id").
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

// 존재하는 상품인지 확인합니다.
func (h *ProductDB) CheckProductExists(ctx context.Context, productId int64) (bool, error) {
	type ProductCount struct {
		Count int `rnsql:"COUNT(*)"`
	}
	result := &ProductCount{}
	sql := gorn.NewSql().
		Select(result).
		From("PRODUCT").
		Where("id = ?", productId)
	row := h.QueryRow(ctx, sql)
	if err := h.ScanRow(row, result); err != nil {
		return false, err
	}
	return result.Count > 0, nil
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

// 새로운 카테고리를 추가합니다.
// 이후 추가된 카테고리 아이디를 반환합니다.
func (h *ProductDB) AddCategory(ctx context.Context, category *dbmodel.Category) (int64, error) {
	return h.InsertWithLastId(ctx, "CATEGORY", category)
}

// 모든 카테고리를 삭제합니다.
func (h *ProductDB) DeleteAllCategories(ctx context.Context) error {
	sql := gorn.NewSql().
		DeleteFrom("CATEGORY").
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

// 프로덕트에 카테고리를 연결합니다.
func (h *ProductDB) AddProductCategory(ctx context.Context, productCategoryMap *dbmodel.ProductCategoryMap) error {
	return h.Insert(ctx, "PRODUCT_CATEGORY_MAP", productCategoryMap)
}

// 모든 프로덕트에 연결된 카테고리를 삭제합니다.
func (h *ProductDB) DeleteAllProductCategoryMap(ctx context.Context) error {
	sql := gorn.NewSql().
		DeleteFrom("PRODUCT_CATEGORY_MAP").
		Where("product_id > ?", -1)
	res, err := h.Exec(ctx, sql)
	if err != nil {
		return err
	}
	if _, err := res.RowsAffected(); err != nil {
		return err
	}
	return nil
}

// 장바구니에 상품을 추가합니다.
func (h *ProductDB) AddCart(ctx context.Context, cart *dbmodel.Cart) error {
	ntime := time.Now()
	cart.CreatedTime = ntime
	cart.UpdatedTime = ntime
	return h.Insert(ctx, "CART", cart)
}

// 장바구니에 상품이 담겨있는지 확인합니다.
func (h *ProductDB) CheckCartHasProduct(ctx context.Context, userId, productId int64) (bool, error) {
	type CartCount struct {
		Count int `rnsql:"COUNT(*)"`
	}
	result := &CartCount{}
	sql := gorn.NewSql().
		Select(result).
		From("CART").
		Where("user_id = ?", userId).
		Where("product_id = ?", productId)
	row := h.QueryRow(ctx, sql)
	if err := h.ScanRow(row, result); err != nil {
		return false, err
	}
	return result.Count > 0, nil
}

// 장바구니에 담긴 단일 상품 정보를 가져옵니다.
func (h *ProductDB) GetCartProduct(ctx context.Context, userId, productId int64) (*dbmodel.Cart, error) {
	result := &dbmodel.Cart{}
	sql := gorn.NewSql().
		Select(result).
		From("CART").
		Where("user_id = ?", userId).
		Where("product_id = ?", productId)
	row := h.QueryRow(ctx, sql)
	if err := h.ScanRow(row, result); err != nil {
		return nil, err
	}
	return result, nil
}

// 장바구니 정보를 업데이트합니다.
func (h *ProductDB) UpdateCart(ctx context.Context, cart *dbmodel.Cart) error {
	ntime := time.Now()
	cart.UpdatedTime = ntime
	sql := gorn.NewSql().
		Update("CART", cart).
		Where("user_id = ?", cart.UserId).
		Where("product_id = ?", cart.ProductId)
	res, err := h.Exec(ctx, sql)
	if err != nil {
		return err
	}
	if _, err := res.RowsAffected(); err != nil {
		return err
	}
	return nil
}

// 장바구니에서 상품을 삭제합니다.
func (h *ProductDB) DeleteCartProduct(ctx context.Context, userId, productId int64) error {
	sql := gorn.NewSql().
		DeleteFrom("CART").
		Where("user_id = ?", userId).
		Where("product_id = ?", productId)
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
