package confighandler_test

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"github.com/av-belyakov/enricher_geoip/internal/confighandler"
)

const Root_Dir = "enricher_geoip"

var (
	conf *confighandler.ConfigApp

	err error
)

func TestMain(m *testing.M) {
	os.Unsetenv("GO_ENRICHERGEOIP_MAIN")

	//Подключение к NATS
	os.Unsetenv("GO_ENRICHERGEOIP_NHOST")
	os.Unsetenv("GO_ENRICHERGEOIP_NPORT")
	os.Unsetenv("GO_ENRICHERGEOIP_NSUBSCR")
	os.Unsetenv("GO_ENRICHERGEOIP_NCACHETTL")

	//Подключение к GeoIP БД
	os.Unsetenv("GO_ENRICHERGEOIP_GIPHOST")
	os.Unsetenv("GO_ENRICHERGEOIP_GIPPOST")
	os.Unsetenv("GO_ENRICHERGEOIP_GIPPATH")

	//Настройки доступа к БД в которую будут записыватся логи
	os.Unsetenv("GO_ENRICHERGEOIP_DBWLOGHOST")
	os.Unsetenv("GO_ENRICHERGEOIP_DBWLOGPORT")
	os.Unsetenv("GO_ENRICHERGEOIP_DBWLOGNAME")
	os.Unsetenv("GO_ENRICHERGEOIP_DBWLOGUSER")
	os.Unsetenv("GO_ENRICHERGEOIP_DBWLOGPASSWD")
	os.Unsetenv("GO_ENRICHERGEOIP_DBWLOGSTORAGENAME")

	//загружаем ключи и пароли
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalln(err)
	}

	os.Setenv("GO_ENRICHERGEOIP_MAIN", "development")

	conf, err = confighandler.New(Root_Dir)
	if err != nil {
		log.Fatalln(err)
	}

	os.Exit(m.Run())
}

func TestConfigHandler(t *testing.T) {
	t.Run("Тест чтения конфигурационного файла", func(t *testing.T) {
		t.Run("Тест 1. Проверка настройки NATS из файла config_dev.yml", func(t *testing.T) {
			assert.Equal(t, conf.GetNATS().Host, "192.168.9.208")
			assert.Equal(t, conf.GetNATS().Port, 4222)
			assert.Equal(t, conf.GetNATS().CacheTTL, 3600)
			assert.Equal(t, conf.GetNATS().Subscription, "object.geoip-request.test")
		})

		t.Run("Тест 2. Проверка настройки GeoIPDataBase из файла config_dev.yml", func(t *testing.T) {
			assert.Equal(t, conf.GetGeoIPDB().Host, "pg2.cloud.gcm")
			assert.Equal(t, conf.GetGeoIPDB().Port, 88)
			assert.Equal(t, conf.GetGeoIPDB().Path, "ip")
			assert.Equal(t, conf.GetGeoIPDB().RequestTimeout, 5)
		})

		t.Run("Тест 3. Проверка настройки WriteLogDataBase из файла config_dev.yml", func(t *testing.T) {
			assert.Equal(t, conf.GetLogDB().Host, "datahook.cloud.gcm")
			assert.Equal(t, conf.GetLogDB().Port, 9200)
			assert.Equal(t, conf.GetLogDB().User, "log_writer")
			assert.Equal(t, conf.GetLogDB().Passwd, os.Getenv("GO_ENRICHERGEOIP_DBWLOGPASSWD"))
			assert.Equal(t, conf.GetLogDB().NameDB, "")
			assert.Equal(t, conf.GetLogDB().StorageNameDB, "enricher_geoip")
		})
	})

	t.Run("Тест чтения переменных окружения", func(t *testing.T) {
		t.Run("Тест 1. Проверка настройки NATS", func(t *testing.T) {
			os.Setenv("GO_ENRICHERGEOIP_NHOST", "127.0.0.1")
			os.Setenv("GO_ENRICHERGEOIP_NPORT", "4242")
			os.Setenv("GO_ENRICHERGEOIP_NCACHETTL", "650")
			os.Setenv("GO_ENRICHERGEOIP_NSUBSCR", "obj.subscript.test")

			conf, err := confighandler.New(Root_Dir)
			assert.NoError(t, err)

			assert.Equal(t, conf.GetNATS().Host, "127.0.0.1")
			assert.Equal(t, conf.GetNATS().Port, 4242)
			assert.Equal(t, conf.GetNATS().CacheTTL, 650)
			assert.Equal(t, conf.GetNATS().Subscription, "obj.subscript.test")
		})

		t.Run("Тест 2. Проверка настройки GeoIPDataBase", func(t *testing.T) {
			os.Setenv("GO_ENRICHERGEOIP_GIPHOST", "examle.database.cm")
			os.Setenv("GO_ENRICHERGEOIP_GIPPOST", "9559")
			os.Setenv("GO_ENRICHERGEOIP_GIPPATH", "any_path")

			conf, err := confighandler.New(Root_Dir)
			assert.NoError(t, err)

			assert.Equal(t, conf.GetGeoIPDB().Host, "examle.database.cm")
			assert.Equal(t, conf.GetGeoIPDB().Port, 9559)
			assert.Equal(t, conf.GetGeoIPDB().Path, "any_path")
		})

		t.Run("Тест 3. Проверка настройки WriteLogDataBase", func(t *testing.T) {
			os.Setenv("GO_ENRICHERGEOIP_DBWLOGHOST", "domaniname.database.cm")
			os.Setenv("GO_ENRICHERGEOIP_DBWLOGPORT", "8989")
			os.Setenv("GO_ENRICHERGEOIP_DBWLOGUSER", "somebody_user")
			os.Setenv("GO_ENRICHERGEOIP_DBWLOGNAME", "any_name_db")
			os.Setenv("GO_ENRICHERGEOIP_DBWLOGPASSWD", "your_passwd")
			os.Setenv("GO_ENRICHERGEOIP_DBWLOGSTORAGENAME", "log_storage")

			conf, err := confighandler.New(Root_Dir)
			assert.NoError(t, err)

			assert.Equal(t, conf.GetLogDB().Host, "domaniname.database.cm")
			assert.Equal(t, conf.GetLogDB().Port, 8989)
			assert.Equal(t, conf.GetLogDB().User, "somebody_user")
			assert.Equal(t, conf.GetLogDB().Passwd, os.Getenv("GO_ENRICHERGEOIP_DBWLOGPASSWD"))
			assert.Equal(t, conf.GetLogDB().NameDB, "any_name_db")
			assert.Equal(t, conf.GetLogDB().StorageNameDB, "log_storage")
		})
	})
}
