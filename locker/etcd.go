// Package locker 分布式锁实现
package locker

//type EtcdLocker struct {
//	cli *etcd.Client
//}

//
//func newEtcd(opt *Option) *EtcdLocker {
//	cli, err := etcd.New(etcd.Config{
//		Endpoints: opt.Url,
//		TLS:       opt.TLS,
//		Username:  opt.Username,
//		Password:  opt.Password,
//	})
//	if err != nil {
//		log.Panicln("Init locker with etcd failed")
//	}
//	return &EtcdLocker{cli}
//}
//
//// Lock 加锁
//// 此处的k 请使用目录形式，eg：/dir_test
//func (r *EtcdLocker) Lock(k string, ex ...time.Duration) bool {
//	ex = append(ex, time.Minute*10)
//	c := etcd.NewKV(r.cli)
//	lease := etcd.NewLease(r.cli)
//	leaseGrantResponse, err := lease.Grant(context.TODO(), int64(ex[0].Seconds()))
//	txn := c.Txn(context.Background())
//	resp, err := txn.If(etcd.Compare(etcd.CreateRevision(k), "=", 0)).
//		Then(etcd.OpPut(k, time.Now().Format(time.RFC3339), etcd.WithLease(leaseGrantResponse.ID))).
//		Else(etcd.OpGet("k")).Commit()
//	if err != nil {
//		return false
//	}
//	return resp.Succeeded
//}
//
//// Unlock 解锁
//// 此处的k 请使用目录形式，eg：/dir_test
//func (r *EtcdLocker) Unlock(k string) {
//	session, err := concurrency.NewSession(r.cli)
//	if err != nil {
//		log.Panicln("locker failed(use etcd)")
//	}
//	defer session.Close()
//	l := concurrency.NewMutex(session, k)
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
//	defer cancel()
//	_ = l.Unlock(ctx)
//}
