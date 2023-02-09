# JGC-API

JongGeon Class 의 API를 제공하는 서버입니다.

---

## Build Setup

-

## Documents

이 프로젝트는 Golang Clean Architecture를 기반으로 설계되었습니다.

레이어를 나누고 각 레이어들이 다른 레이어를 침범하지 않도록 개발해야 합니다.

현재 프로젝트는 다음과 같이 레이어가 존재합니다.

### router

에 레이어는 EndPoint를 Handler로 연결해주는 역할을 합니다. middleware를 붙이는 작업이나 EndPoint관리를 편하게 하기 위해 만들었습니다.

### middleware

이 레이어는 각 EndPoint에 대한 요청을 처리하기 이전 실행되어야 할 로직들을 정리한 레이어입니다.  
토큰 인증과 같이 핸들러에서 중복되어 나타나는 로직들을 처리하기 위해 만들어졌습니다.

### handler

이 레이어는 모든 EndPoint들을 관리해주는 레이어입니다.  
**비즈니스 로직은 포함되어 있지 않으며**, 비즈니스 로직을 호출해주는 역할을 합니다. Middleware에서 처리한 데이터를 파싱하거나, 파라미터를 비즈니스 로직으로 넘겨주는 역할입니다.

### usecase

이 레이어는 handler 레이어에서 호출할 비즈니스 로직입니다.  
**핵심 로직은 이 부분에 포함되어 있으며**, 써드파티 툴이나 재사용이 가능한 로직은 **포함되면 안됩니다.**

### model

이 레이어는 재사용 되는 모델을 정의하는 부분입니다. golang에서는 언어적으로 순환참조가 안되고 순환참조가 발생하면 유지보수 자체가 복잡해지기 때문에 재사용이 가능한, 혹은 여러 군데에서 참조해야할 모델을 정의하는 부분입니다.

### dbmodel

이 레이어는 데이터베이스에 저장되는 모델을 정의하는 부분입니다. 여기서 정의한 대로 디비 테이블을 생성합니다.

### database

이 레이어는 데이터베이스를 다루는 레이어입니다. 데이터베이스에 관련된 모든 코드는 이 곳에 포함되어야 합니다. **외부에서 쿼리를 직접 날리는 일은 발생하면 안됩니다.**

### config

이 레이어는 하이퍼 파라미터를 설정하는 레이어입니다. Init함수로 넘겨주는 environment에 따라 다른 데이터를 담은 객체를 반환합니다. 하이퍼 파라미터로 하드코딩 되면 안되는 데이터나 환경마다 설정이 달라져야 하는 데이터는 모두 이 곳에 포함되어야 합니다.

### util

이 레이어는 유틸성이 강한 함수들을 모아놓는 파일입니다.

---

## Architecture

이 프로젝트는 Golang Clean Architecture를 기반으로 설계되었습니다.

`main` -> `router` -> `handler` -> `usecase` -> `database`

형식으로 이루어져 있습니다.

## Run script

### Building Api Server Docker Image

~~~shell
 $ make server
~~~

### Run Api Server

~~~shell
 $ make run
~~~

도커를 사용하지 않을 경우 아래와 같이 실행해 주세요

~~~shell
 $ make serve
~~~

### Stop Api Server

~~~shell
 $ make stop
~~~

도커를 사용하지 않을 경우 아래와 같이 실행해 주세요

~~~shell
 $ ^C
~~~
