package user_test

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var mysqlTestDB *gorm.DB

func TestMain(m *testing.M) {
	pool, resource, db := InitialDockerMySQL()
	mysqlTestDB = db

	code := m.Run()

	// Clean up resource
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func InitialDockerMySQL() (
	pool *dockertest.Pool,
	resource *dockertest.Resource,
	db *gorm.DB,
) {
	var err error
	pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	options := &dockertest.RunOptions{
		Name:       "mysql_user_test",
		Repository: "mysql",
		Tag:        "8.0",
		Env: []string{
			"MYSQL_ROOT_PASSWORD=root_password",
			"MYSQL_DATABASE=banking",
		},
		ExposedPorts: []string{"3306/tcp"},
	}

	resource, err = pool.RunWithOptions(options, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// Exponential backoff-retry for the container to be ready
	if err = pool.Retry(func() error {
		dsn := fmt.Sprintf(
			"root:root_password@tcp(%s)/banking?charset=utf8mb4&parseTime=True&loc=Local",
			resource.GetHostPort("3306/tcp"),
		)

		location, errL := time.LoadLocation("UTC")
		if errL != nil {
			return errL
		}

		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
				TablePrefix:   "banking_",
			},
			Logger: logger.Default.LogMode(logger.Info),
			NowFunc: func() time.Time {
				return time.Now().In(location)
			},
		})
		if err != nil {
			return err
		}

		sqlDB, errDB := db.DB()
		if errDB != nil {
			return errDB
		}

		return sqlDB.Ping()
	}); err != nil {
		// Clean up resource if there is an error
		if purgeErr := pool.Purge(resource); purgeErr != nil {
			log.Fatalf("Could not purge resource: %s", purgeErr)
		}
		log.Fatalf("Could not connect to docker: %s", err)
	}

	return pool, resource, db
}

func getHostPort(resource *dockertest.Resource, id string) string {
	dockerURL := os.Getenv("DOCKER_HOST")
	if dockerURL == "" {
		return resource.GetHostPort(id)
	}
	u, err := url.Parse(dockerURL)
	if err != nil {
		panic(err)
	}
	return u.Hostname() + ":" + resource.GetPort(id)
}
