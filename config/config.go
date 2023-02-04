package config

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/thak1411/rnlog"
)

// config 정보를 담을 객체입니다.
// 보통 hyper parameter를 여기서 관리하며
// 추가해야하는 정보가 생긴다면 struct에 추가한 뒤 default_config에 무조건 추가한 뒤,
// 다른 환경에 데이터에 값을 추가하거나, 추가하지 않으셔도 됩니다.
type Config struct {
	// 로그 파일의 경로입니다.
	LogFilePath string

	// 로그 파일의 이름입니다.
	LogFile string

	// 돌아가고 있는 서버의 포트와 도메인입니다.
	Port   int
	Domain string

	// 브라우저 쿠키에 관련된 데이터입니다.
	Cookies struct {
		// 유저 정보를 파싱하게 도와줄 쿠키입니다.
		PublicSessionName string

		// 유저의 인증 토큰을 담은 쿠키입니다.
		SessionName string

		// 두 쿠키의 유지 시간입니다.
		SessionTimeout time.Duration
	}

	// jwt 관련 데이터입니다.
	Jwt struct {
		// jwt 데이터의 인증 키입니다.
		SecretKey string
	}

	// 데이터 베이스 관련 데이터입니다.
	DB struct {
		// 접속할 데이터베이스의 스키마입니다.
		JGCSchema string

		// 커넥션 풀의 대기 커넥션 수입니다.
		PoolSize int

		// 커넥션 풀의 최대 연결 수입니다.
		MaxConn int

		// 디비의 유저 이름 입니다.
		User string

		// 디비의 비밀번호 입니다.
		Password string

		// 디비의 도메인 혹은 ip 입니다.
		Host string

		// 디비의 포트 입니다.
		Port int

		// 디비의 라이프 사이클 입니다. (한 커넥션이 유지되는 시간)
		// Mysql의 최대 커넥션 주기보다 짧아야 합니다.
		// 기본 값은 8시간이므로 이보다 짧은 7시간으로 설정합니다.)
		Lifecycle time.Duration
	}
}

// 객체를 복사합니다.
func (c *Config) Copy(dest *Config) {
	if dest == nil {
		rnlog.Fatal("dest is null")
		panic("dest is null")
	}
	var b bytes.Buffer

	encoder := gob.NewEncoder(&b)
	decoder := gob.NewDecoder(&b)

	encoder.Encode(c)
	decoder.Decode(dest)
}

var env = "default"

// 환경을 초기화 합니다.
// 메인에서 최초 한 번만 호출되어야 합니다.
func Init(_env string) {
	env = _env
}

// Init함수로 초기화해준 Config 객체를 반환합니다.
// 초기화 하지 않으면 default_config가 반환됩니다.
// 객체를 불러오는 우선순위는 다음과 같습니다.
// 1. default_config 데이터를 불러옵니다.
// 2. 그 위에 각 환경마다 저장되어있는 데이터를 덮어씌웁니다.
// 따라서 default_config가 설정돼 있어도, 각 환경에서 다른 데이터를 넣어준다면 default_config 데이터는 무시됩니다.
func Get() *Config {
	switch env {
	case "default":
		return GetDefault()
	case "native":
		return GetNative()
	case "test":
		return GetTest()
	case "product":
		return GetProduct()
	default:
		rnlog.Fatal("Unknown config env: %s", env)
		panic(fmt.Errorf("unknown config env: %s", env))
	}
}
