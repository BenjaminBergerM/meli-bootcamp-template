package section

import (
	"context"
	"database/sql"
	"errors"

	"github.com/BenjaminBergerM/bootcamp-go/internal/domain"
)

// Repository encapsulates the storage of a section.
type Repository interface {
	GetAll(ctx context.Context) ([]domain.Section, error)
	Get(ctx context.Context, id int) (domain.Section, error)
	Exists(ctx context.Context, cid int) bool
	Save(ctx context.Context, s domain.Section) (int, error)
	Update(ctx context.Context, s domain.Section) error
	Delete(ctx context.Context, id int) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetAll(ctx context.Context) ([]domain.Section, error) {
	rows, err := r.db.Query(`SELECT * FROM "main"."sections"`)
	if err != nil {
		return nil, err
	}

	var sections []domain.Section

	for rows.Next() {
		s := domain.Section{}
		_ = rows.Scan(&s.ID, &s.SectionNumber, &s.CurrentTemperature, &s.MinimumTemperature, &s.CurrentCapacity, &s.MinimumCapacity, &s.MaximumCapacity, &s.WarehouseID, &s.ProductTypeID)
		sections = append(sections, s)
	}

	return sections, nil
}

func (r *repository) Get(ctx context.Context, id int) (domain.Section, error) {

	sqlStatement := `SELECT * FROM "main"."sections" WHERE id=$1;`
	row := r.db.QueryRow(sqlStatement, id)
	s := domain.Section{}
	err := row.Scan(&s.ID, &s.SectionNumber, &s.CurrentTemperature, &s.MinimumTemperature, &s.CurrentCapacity, &s.MinimumCapacity, &s.MaximumCapacity, &s.WarehouseID, &s.ProductTypeID)
	if err != nil {
		return domain.Section{}, err
	}

	return s, nil
}

func (r *repository) Exists(ctx context.Context, sectionNumber int) bool {
	sqlStatement := `SELECT section_number FROM "main"."sections" WHERE section_number=$1;`
	row := r.db.QueryRow(sqlStatement, sectionNumber)
	err := row.Scan(&sectionNumber)
	if err != nil {
		return false
	}
	return true
}

func (r *repository) Save(ctx context.Context, s domain.Section) (int, error) {

	stmt, err := r.db.Prepare(`INSERT INTO "main"."sections"("section_number","current_temperature","minimum_temperature","current_capacity","minimum_capacity", "maximum_capacity", "warehouse_id", "product_type_id") VALUES (?,?,?,?,?,?,?,?)`)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(&s.SectionNumber, &s.CurrentTemperature, &s.MinimumTemperature, &s.CurrentCapacity, &s.MinimumCapacity, &s.MaximumCapacity, &s.WarehouseID, &s.ProductTypeID)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *repository) Update(ctx context.Context, s domain.Section) error {
	stmt, err := r.db.Prepare(`UPDATE "main"."sections" SET "section_number"=?, "current_temperature"=?, "minimum_temperature"=?, "current_capacity"=?, "minimum_capacity"=?, "maximum_capacity"=?, "warehouse_id"=?, "product_type_id"=?  WHERE "id"=?`)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(&s.SectionNumber, &s.CurrentTemperature, &s.MinimumTemperature, &s.CurrentCapacity, &s.MinimumCapacity, &s.MaximumCapacity, &s.WarehouseID, &s.ProductTypeID, &s.ID)
	if err != nil {
		return err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affect < 1 {
		return errors.New("section not found")
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, id int) error {
	stmt, err := r.db.Prepare(`DELETE FROM "main"."sections" WHERE id=?`)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affect < 1 {
		return errors.New("section not found")
	}

	return nil
}
