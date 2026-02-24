package container

import (
	"context"
	"os"

	"github.com/OddKuru/core-accounts/internal/app/usecase/account"
	"github.com/OddKuru/core-accounts/internal/infra/config"
	"github.com/OddKuru/core-accounts/internal/infra/service/password"
	"github.com/OddKuru/core-accounts/internal/infra/storage/psql"
	ayaka "github.com/OddKuru/core-accounts/pkg/core"
	"github.com/OddKuru/core-accounts/pkg/ecosystem"
	"github.com/OddKuru/core-accounts/pkg/logger"
	"github.com/OddKuru/core-accounts/pkg/logger/log"
	"github.com/OddKuru/core-accounts/pkg/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

const (
	ConfigPathEnvName = "CONFIG_PATH"
)

type Dependency struct {
	accountsUseCase account.Account
	appLogger       ayaka.Logger
	logger          logger.Logger
	postgresPool    *pgxpool.Pool
	config          *config.Config
}

func AppContainer() (*Dependency, error) {
	deps := &Dependency{}
	err := deps.combine(
		deps.initConfig,
		deps.initLogger,
		deps.connectDatabase,
		deps.initUseCases,
	)
	if err != nil {
		return nil, err
	}

	return deps, nil
}

func (d *Dependency) AccountsUseCase() account.Account {
	return d.accountsUseCase
}

func (d *Dependency) AppLogger() ayaka.Logger {
	return d.appLogger
}

func (d *Dependency) Config() *config.Config {
	return d.config
}

func (d *Dependency) Logger() logger.Logger {
	return d.logger
}

func (d *Dependency) combine(fns ...func() error) error {
	for _, fn := range fns {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

func (d *Dependency) initConfig() error {
	value := os.Getenv(ConfigPathEnvName)
	if value == "" {
		return errors.New("config env variable not found")
	}

	cfg, err := config.NewConfig(value)
	if err != nil {
		return errors.Wrap(err, "[Dependency] config.New")
	}

	d.config = cfg
	return nil
}

func (d *Dependency) initLogger() error {
	cfg := d.config.Logger

	var cfgLvl logger.Level
	switch cfg.Level {
	case "debug":
		cfgLvl = logger.DebugLvl
	case "info":
		cfgLvl = logger.InfoLvl
	case "warn":
		cfgLvl = logger.WarnLvl
	case "error":
		cfgLvl = logger.ErrorLvl
	}

	out := log.DefaultOutput(cfg.Pretty)

	l := log.NewLogger(&log.Options{
		Out:    out,
		Pretty: cfg.Pretty,
		Level:  cfgLvl,
	})

	d.logger = l
	d.appLogger = ecosystem.NewAppLoggerWithZerolog(l.OriginalLogger())

	return nil
}

func (d *Dependency) connectDatabase() error {
	conn, err := psql.Connect(context.Background(), d.config.Psql)
	if err != nil {
		return errors.Wrap(err, "[Dependency] psql.Connect")
	}
	d.postgresPool = conn
	return nil
}

func (d *Dependency) initUseCases() error {
	timeNow := utils.TimeNow{}
	IDGen := UUIDV4{}
	passwordManager, err := password.New()
	if err != nil {
		return errors.Wrap(err, "[Dependency] password.New")
	}

	accQuery, err := psql.NewAccountQuery(d.postgresPool)
	if err != nil {
		return errors.Wrap(err, "[Dependency] psql.NewAccountQuery")
	}

	uow, err := psql.NewUnitOfWork(d.postgresPool)
	if err != nil {
		return errors.Wrap(err, "[Dependency] psql.NewUnitOfWork")
	}

	accUseCase, err := account.NewUseCase(
		d.logger,
		accQuery,
		uow,
		passwordManager,
		IDGen,
		timeNow,
	)
	if err != nil {
		return errors.Wrap(err, "[Dependency] account.NewUseCase")
	}
	d.accountsUseCase = accUseCase
	return nil
}
