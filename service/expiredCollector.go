package service

type ExpiredCollector interface {
	Run()
	Reload()
}
