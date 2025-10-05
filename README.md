# JSON CRUD Service

클린 아키텍처를 기반으로 구현된 JSON 문서 CRUD 서비스입니다.

## 아키텍처

이 프로젝트는 클린 아키텍처 원칙을 따라 다음과 같이 구성되어 있습니다:

```
├── cmd/server/           # 애플리케이션 진입점
├── internal/
│   ├── domain/          # 도메인 레이어
│   │   ├── entity/      # 도메인 엔티티
│   │   └── repository/  # 저장소 인터페이스
│   ├── usecase/         # 유즈케이스 레이어 (비즈니스 로직)
│   ├── infrastructure/  # 인프라스트럭처 레이어
│   │   └── repository/  # 저장소 구현체
│   └── presentation/    # 프레젠테이션 레이어
│       ├── handler/     # HTTP 핸들러
│       └── middleware/  # 미들웨어
├── pkg/                 # 공통 패키지
│   └── response/        # HTTP 응답 유틸리티
└── config/              # 설정
```

## 기능

- JSON 문서 생성 (Create)
- JSON 문서 조회 (Read)
- JSON 문서 수정 (Update)
- JSON 문서 삭제 (Delete)
- JSON 문서 목록 조회

## API 엔드포인트

### 헬스 체크
- `GET /health` - 서비스 상태 확인

### JSON 문서 관리
- `POST /documents` - 새 문서 생성
- `GET /documents` - 모든 문서 목록 조회
- `GET /documents/{id}` - 특정 문서 조회
- `PUT /documents/{id}` - 문서 수정
- `DELETE /documents/{id}` - 문서 삭제

## 사용 방법

### 서버 실행
```bash
# 빌드
go build -o bin/server cmd/server/main.go

# 실행
./bin/server

# 또는 포트 지정
PORT=3000 ./bin/server
```

### API 사용 예시

#### 1. 새 문서 생성
```bash
curl -X POST http://localhost:8080/documents \
  -H "Content-Type: application/json" \
  -d '{
    "id": "doc1",
    "data": {
      "name": "John Doe",
      "age": 30,
      "city": "Seoul"
    }
  }'
```

#### 2. 문서 조회
```bash
curl http://localhost:8080/documents/doc1
```

#### 3. 문서 수정
```bash
curl -X PUT http://localhost:8080/documents/doc1 \
  -H "Content-Type: application/json" \
  -d '{
    "data": {
      "name": "John Doe",
      "age": 31,
      "city": "Busan"
    }
  }'
```

#### 4. 모든 문서 조회
```bash
curl http://localhost:8080/documents
```

#### 5. 문서 삭제
```bash
curl -X DELETE http://localhost:8080/documents/doc1
```

## 응답 형식

### 성공 응답
```json
{
  "success": true,
  "message": "Document created successfully",
  "data": {
    "id": "doc1",
    "data": {
      "name": "John Doe",
      "age": 30,
      "city": "Seoul"
    },
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### 에러 응답
```json
{
  "success": false,
  "error": "document not found"
}
```

## 기술 스택

- **언어**: Go 1.21+
- **아키텍처**: Clean Architecture
- **저장소**: In-Memory (메모리)
- **HTTP 서버**: net/http (표준 라이브러리)

## 개발 환경 설정

1. Go 1.21 이상 설치
2. 저장소 클론
3. 의존성 설치: `go mod tidy`
4. 빌드: `go build -o bin/server cmd/server/main.go`
5. 실행: `./bin/server`

## 확장 가능성

현재는 메모리 저장소를 사용하고 있지만, 인터페이스 기반 설계로 인해 다음과 같은 확장이 가능합니다:

- 데이터베이스 저장소 (PostgreSQL, MongoDB 등)
- 파일 시스템 저장소
- 클라우드 저장소 (AWS S3, Google Cloud Storage 등)
- 캐싱 레이어 추가
- 인증/권한 시스템 추가
- 로깅 시스템 추가
- 메트릭 수집 추가