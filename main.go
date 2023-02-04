package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/JongGeonClass/JGC-API/config"
	"github.com/JongGeonClass/JGC-API/migrate"
	"github.com/JongGeonClass/JGC-API/router"
	"github.com/JongGeonClass/JGC-API/util"
	"github.com/thak1411/gorn"
	"github.com/thak1411/rnlog"
)

func main() {
	// 환경 분석을 위해 플래그를 파싱합니다.
	// 기본 플래그는 로컬입니다.
	envp := flag.String("env", "native", "Environment\n- native\n- test\n- product\n")
	isMigrate := flag.Bool("migrate", false, "Migrate database")
	flag.Parse()

	// config file을 초기화 합니다. 이때 rn logger를 사용하는데 초기화 하지 않았으므로,
	// 에러 로그가 파일로 저장되지 않습니다.
	config.Init(*envp)

	conf := config.Get()

	// rn logger를 초기화 합니다.
	err := rnlog.Init(conf.LogFilePath, conf.LogFile)
	if err != nil {
		fmt.Println("Failed to initialize logger")
		os.Exit(1)
	}
	defer rnlog.Close()
	rnlog.Log(util.BarLine(120))

	// 디비는 서버가 종료될 때 닫아주어야 하기 때문에, 메인에서 생성합니다.
	db := gorn.NewDB("mysql")
	if err := db.Open(&gorn.DBConfig{
		User:      conf.DB.User,
		Password:  conf.DB.Password,
		Host:      conf.DB.Host,
		Port:      conf.DB.Port,
		Schema:    conf.DB.JGCSchema,
		PoolSize:  conf.DB.PoolSize,
		MaxConn:   conf.DB.MaxConn,
		Lifecycle: conf.DB.Lifecycle,
	}); err != nil {
		rnlog.Fatal("DB open Error: %+v\n", err)
		return
	} else {
		rnlog.Info("Successfully connected to DB(%s)", conf.DB.JGCSchema)
	}
	defer db.Close()
	// 만약 마이그레이션 로직을 실행해야 한다면
	// 마이그레이션 로직을 실행하고 종료합니다.
	if isMigrate != nil && *isMigrate {
		migrate.Migrate(db)
		return
	}

	router := router.New()

	rnlog.Info("JGC API server is running...")
	rnlog.Info("Server port: %d", conf.Port)

	if err = router.Run(conf.Port); err != nil {
		rnlog.Error("Server error: " + err.Error())
	} else {
		rnlog.Info("Server is shutting down...")
	}
}
