package link

import (
	"gorm.io/gorm/clause"
	"links-shortener/pkg/db"
)

type LinkRepository struct {
	database *db.Db
}

func NewLinkRepository(database *db.Db) *LinkRepository {
	return &LinkRepository{database: database}
}

func (repo *LinkRepository) Create(link *Link) (*Link, error) {
	result := repo.database.DB.Create(link)
	if result.Error != nil {
		return nil, result.Error
	}
	return link, nil
}

func (repo *LinkRepository) GetByHash(hash string) (*Link, error) {
	var link Link
	result := repo.database.DB.First(&link, "hash = ?", hash)
	if result.Error != nil {
		return nil, result.Error
	}
	return &link, nil
}

func (repo *LinkRepository) Update(link *Link) (*Link, error) {
	result := repo.database.DB.Clauses(clause.Returning{}).Updates(link)
	if result.Error != nil {
		return nil, result.Error
	}
	return link, nil
}

func (repo *LinkRepository) Delete(id uint) error {
	result := repo.database.DB.Delete(&Link{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *LinkRepository) GetById(id uint) (*Link, error) {
	var link Link
	result := repo.database.DB.First(&link, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &link, nil
}

func (repo *LinkRepository) GetCount() int64 {
	var count int64
	repo.database.
		Table("links").
		Where("deleted_at IS NULL").
		Count(&count)
	return count
}

func (repo *LinkRepository) GetAll(limit, offset int) []Link {
	var links []Link
	repo.database.
		Table("links").
		Where("deleted_at IS NULL").
		Order("id asc").
		Limit(limit).
		Offset(offset).
		Scan(&links)
	return links
}
