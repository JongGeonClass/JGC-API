# API서버를 실행시킬 이미지입니다.

# 빌드 이미지
# ----------------------------------------
FROM golang:1.17 AS builder

# 프로젝트 디렉토리를 복사합니다.
ADD . /src

# 테스트를 진행한 뒤 설치합니다.
RUN cd /src && \
    go get && \
    go test && \
    GOOS=linux CGO_ENABLED=0 go build -a -installsuffix cgo

# 프로덕션 이미지
# ----------------------------------------
FROM alpine

# qemu-x86_64: Could not open '/lib64/ld-linux-x86-64.so.2': No such file or directory
# 위 오류를 해결하기 위한 모듈입니다.
RUN apk add gcompat

# working directory 생성
RUN mkdir -p /jgc

# 서버 실행파일 복사
COPY --from=builder /src/JGC-API /jgc/JGC-API

ENTRYPOINT ["jgc/JGC-API"]
