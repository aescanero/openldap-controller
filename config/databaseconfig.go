package config

type DatabaseConfig struct {
	Base       string
	Replicatls []ReplicaTls
}

func (dbIn *DatabaseConfig) ImportNotNull(db *DatabaseConfig) {
	if db.Base != "" {
		dbIn.Base = db.Base
	}
	ix := 0
	for _, rep := range db.Replicatls {
		if len(dbIn.Replicatls) <= ix {
			dbIn.Replicatls = append(dbIn.Replicatls, rep)
		} else {
			dbIn.Replicatls[ix].ImportNotNull(&rep)
		}

		ix = ix + 1
	}
}
