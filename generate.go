package goidempo

//go:generate go-enum -f cache.go
//
//go:generate mockgen -source=cache.go -destination=mocks/cache.go -package=mocks
