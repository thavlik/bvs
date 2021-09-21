package node

func (s *Server) startDBSyncStack(dbPath string, postgresPort int) error {
	// We need to start postgres and cardano-db-sync
	//postgresDone := make(chan error, 1)
	//go func() {
	//	postgresDone <- s.startPostgres(dbPath, postgresPort)
	//	close(postgresDone)
	//}()

	//dbSyncDone := make(chan error, 1)
	//go func() {
	//	dbSyncDone <- s.startDBSync()
	//	close(dbSyncDone)
	//}()

	select {
	//case err := <-dbSyncDone:
	//	return fmt.Errorf("db sync: %v", err)
	//case err := <-postgresDone:
	//	return fmt.Errorf("postgres: %v", err)
	}
}
