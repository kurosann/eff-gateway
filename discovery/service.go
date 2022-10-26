package discovery

func InitService() {
	InitEtcd()
}

func KeepAlive(key string, f EtcdEvFunc) {
	EC.KeepWatch(key, f)
}

func DelSv(key string) error {
	return EC.Del(key)
}
func PutSv(key string, value string) error {
	return EC.Put(key, value)
}
