package repo

import "com.aharakitchen/app/domain"

type AuthRepo interface {
	Login(username string, password string, ip string, ips []string) (*domain.Admin, string, error)
}

