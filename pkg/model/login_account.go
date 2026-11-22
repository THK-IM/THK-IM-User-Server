package model

import (
	"crypto/subtle"
	"github.com/thk-im/thk-im-user-server/pkg/errorx"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LoginAccount struct {
	Account      string
	UserId       int64
	PasswordHash *string
	Version      int64
}

func (d defaultUserModel) TokenVersion(uid int64) (int64, error) {
	var version int64
	err := d.db.Table("account_login_0").Select("version").Where("user_id=?", uid).Scan(&version).Error
	return version, err
}
func (d defaultUserModel) verifyPassword(tx *gorm.DB, a *LoginAccount, password string) error {
	if password == "" {
		return errorx.InvalidCredentials
	}
	if a.PasswordHash != nil {
		if bcrypt.CompareHashAndPassword([]byte(*a.PasswordHash), []byte(password)) != nil {
			return errorx.InvalidCredentials
		}
		return nil
	}
	var legacy Account
	q := tx.Table(d.genAccountTableName(a.UserId)).Where("user_id=? AND channel=? AND LOWER(TRIM(account))=?", a.UserId, ChannelDefault, a.Account)
	if err := q.Take(&legacy).Error; err != nil {
		return err
	}
	if legacy.Password == nil || subtle.ConstantTimeCompare([]byte(*legacy.Password), []byte(password)) != 1 {
		return errorx.InvalidCredentials
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if err = tx.Table("account_login_0").Where("user_id=?", a.UserId).Update("password_hash", string(hash)).Error; err != nil {
		return err
	}
	return tx.Table(d.genAccountTableName(a.UserId)).Where("id=?", legacy.Id).Update("password", nil).Error
}
func (d defaultUserModel) Authenticate(account, password string) (*LoginAccount, error) {
	var a LoginAccount
	err := d.db.Transaction(func(tx *gorm.DB) error {
		q := tx.Table("account_login_0").Clauses(clause.Locking{Strength: "UPDATE"}).Where("account=?", account).Find(&a)
		if q.Error != nil {
			return q.Error
		}
		if q.RowsAffected != 1 {
			return errorx.InvalidCredentials
		}
		return d.verifyPassword(tx, &a, password)
	})
	return &a, err
}
func (d defaultUserModel) ChangePassword(uid int64, oldPassword, newPassword string, expectedVersion int64) (int64, error) {
	var a LoginAccount
	err := d.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("account_login_0").Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id=?", uid).Take(&a).Error; err != nil {
			return err
		}
		if a.Version != expectedVersion {
			return errorx.UserTokenError
		}
		if err := d.verifyPassword(tx, &a, oldPassword); err != nil {
			return err
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		a.Version++
		return tx.Table("account_login_0").Where("user_id=?", uid).Updates(map[string]interface{}{"version": a.Version, "password_hash": string(hash)}).Error
	})
	return a.Version, err
}
